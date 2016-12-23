package expression

import "github.com/gitql/gitql/sql"

type BooleanExpression struct{}

func (BooleanExpression) Type() sql.Type {
	return sql.Boolean
}

type Not struct {
	BooleanExpression
	UnaryExpression
}

func NewNot(child sql.Expression) *Not {
	return &Not{
		UnaryExpression: UnaryExpression{
			SimpleName: "not",
			Child:      child,
			Copy: func(e sql.Expression) sql.Expression {
				return NewNot(e)
			},
		},
	}
}

func (e Not) Eval(row sql.Row) interface{} {
	return !e.Child.Eval(row).(bool)
}
