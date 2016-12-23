package expression

import (
	"github.com/gitql/gitql/sql"
)

type Alias struct {
	UnaryExpression
}

func NewAlias(child sql.Expression, name string) *Alias {
	return &Alias{UnaryExpression{
		SimpleName: name,
		Child: child,
		Copy: func(e sql.Expression) sql.Expression {
			return NewAlias(e, name)
		},
	}}
}

func (e *Alias) Type() sql.Type {
	return e.Child.Type()
}

func (e *Alias) Eval(row sql.Row) interface{} {
	return e.Child.Eval(row)
}

func (e *Alias) Name() string {
	return e.SimpleName
}
