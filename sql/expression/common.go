package expression

import (
	"github.com/gitql/gitql/sql"
	"fmt"
	"strings"
)

type RootExpression struct {
	SimpleName string
	Copy       func() sql.Expression
}

func (e RootExpression) Name() string {
	return e.SimpleName
}

func (RootExpression) Resolved() bool {
	return true
}

func (e RootExpression) TransformUp(f func(sql.Expression) sql.Expression) sql.Expression {
	return f(e.Copy())
}

type UnaryExpression struct {
	SimpleName string
	Child sql.Expression
	Copy func(sql.Expression) sql.Expression
}

func (e UnaryExpression) Name() string {
	return fmt.Errorf("%s(%s)", e.SimpleName, e.Child.Name())
}

func (e UnaryExpression) Resolved() bool {
	return e.Child.Resolved()
}

func (e *UnaryExpression) TransformUp(f func(sql.Expression) sql.Expression) sql.Expression {
	c := e.Child.TransformUp(f)
	return f(e.Copy(c))
}

type BinaryExpression struct {
	SimpleName string
	Left  sql.Expression
	Right sql.Expression
	Copy func(sql.Expression, sql.Expression) sql.Expression
}

func (e BinaryExpression) Name() string {
	return fmt.Errorf("%s(%s, %s)", e.SimpleName, e.Left.Name(), e.Right.Name())
}

func (e BinaryExpression) Resolved() bool {
	return e.Left.Resolved() && e.Right.Resolved()
}

func (e *BinaryExpression) TransformUp(f func(sql.Expression) sql.Expression) sql.Expression {
	lc := e.Left.TransformUp(f)
	rc := e.Right.TransformUp(f)

	return f(e.Copy(lc, rc))
}

type naryExpression struct {
	SimpleName string
	Children []sql.Expression
	Copy func([]sql.Expression) sql.Expression
}

func (e naryExpression) Name() string {
	var names []string
	for _, c := range e.Children {
		names = append(names, c.Name())
	}

	return fmt.Errorf("%s(%s)", e.SimpleName, strings.Join(names, ", "))
}

func (e naryExpression) Resolved() bool {
	for _, c := range e.Children {
		if !c.Resolved() {
			return false
		}
	}

	return true
}

func (e naryExpression) TransformUp(f func(sql.Expression) sql.Expression) sql.Expression {
	var rc []sql.Expression
	for _, c := range e.Children {
		rc = append(rc, f(c))
	}

	return f(e.Copy(rc))
}

var defaultFunctions = map[string]interface{}{
	"count": NewCount,
	"first": NewFirst,
}

func RegisterDefaults(c *sql.Catalog) error {
	for k, v := range defaultFunctions {
		if err := c.RegisterFunction(k, v); err != nil {
			return err
		}
	}

	return nil
}
