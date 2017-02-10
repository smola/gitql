package analyzer

import (
	"fmt"
	"reflect"

	"github.com/gitql/gitql/sql"
)

const maxAnalysisIterations = 1000

type Analyzer struct {
	Rules           []Rule
	Catalog         *sql.Catalog
	CurrentDatabase string
}

type Rule struct {
	Name  string
	Apply func(*Analyzer, sql.Node) sql.Node
}

func New(catalog *sql.Catalog) *Analyzer {
	return &Analyzer{
		Rules:   DefaultRules,
		Catalog: catalog,
	}
}

func (a *Analyzer) Analyze(n sql.Node) (sql.Node, error) {
	prev := n
	cur := a.analyzeOnce(n)
	i := 0
	for !reflect.DeepEqual(prev, cur) {
		prev = cur
		cur = a.analyzeOnce(n)
		i++
		if i >= maxAnalysisIterations {
			return cur, fmt.Errorf("exceeded max analysis iterations (%d)", maxAnalysisIterations)
		}
	}

	return cur, a.validate(cur)
}

func (a *Analyzer) analyzeOnce(n sql.Node) sql.Node {
	result := n
	for _, rule := range a.Rules {
		result = rule.Apply(a, result)
	}
	return result
}

func (a *Analyzer) validate(n sql.Node) error {
	if !n.Resolved() {
		inner := mostInnerUnresolved(n)
		return fmt.Errorf("plan is not resolved: %q", inner)
	}

	return nil
}

func mostInnerUnresolved(n sql.Node) sql.Node {
	if n.Resolved() {
		return nil
	}

	for _, c := range n.Children() {
		if !c.Resolved() {
			res := mostInnerUnresolved(c)
			if res != nil {
				return res
			}
		}
	}

	return n
}
