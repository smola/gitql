package expression

import (
	"fmt"

	"github.com/gitql/gitql/sql"
)

type Cast struct {
	UnaryExpression
	Target   sql.Type
	CastFunc func(interface{}) interface{}
}

func NewCast(child sql.Expression, target sql.Type) *Cast {
	return &Cast{UnaryExpression{child}, target, nil}
}

func (e *Cast) Resolved() bool {
	return e.CastFunc != nil && e.Child.Resolved()
}

func (e *Cast) Type() sql.Type {
	return e.Target
}

func (e *Cast) Eval(row sql.Row) interface{} {
	return e.CastFunc(e.Child.Eval(row))
}

func (e *Cast) Name() string {
	return fmt.Sprintf("cast(%s as %s)", e.Child.Name(), e.Target.Name())
}

func (e *Cast) TransformUp(f func(sql.Expression) sql.Expression) sql.Expression {
	c := e.Child.TransformUp(f)
	n := NewCast(c, e.Target)

	return f(n)
}
