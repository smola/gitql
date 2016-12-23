package expression

import (
	"github.com/gitql/gitql/sql"
)

type Count struct {
	UnaryExpression
}

func NewCount(e sql.Expression) *Count {
	return &Count{UnaryExpression{
		SimpleName: "count",
		Child: e,
		Copy: func(e sql.Expression) sql.Expression {
			return NewCount(e)
		},
	}}
}

func (c *Count) NewBuffer() sql.Row {
	return sql.NewMemoryRow(int32(0))
}

func (c *Count) Type() sql.Type {
	return sql.Integer
}

func (c *Count) Resolved() bool {
	if _, ok := c.Child.(*Star); ok {
		return true
	}

	return c.UnaryExpression.Resolved()
}

func (c *Count) Update(buffer, row sql.Row) {
	mr := buffer.(sql.MemoryRow)
	var inc bool
	if _, ok := c.Child.(*Star); ok {
		inc = true
	} else {
		v := c.Child.Eval(row)
		if v != nil {
			inc = true
		}
	}

	if inc {
		mr[0] = getInt32At(buffer, 0) + int32(1)
	}
}

func (c *Count) Merge(buffer, partial sql.Row) {
	mb := buffer.(sql.MemoryRow)
	mb[0] = getInt32At(buffer, 0) + getInt32At(partial, 0)
}

func (c *Count) Eval(buffer sql.Row) interface{} {
	return getInt32At(buffer, 0)
}

type First struct {
	UnaryExpression
}

func NewFirst(e sql.Expression) *First {
	return &First{UnaryExpression{
		Child: e,
		Copy: func(e sql.Expression) sql.Expression { return NewFirst(e) },
	}}
}

func (e *First) NewBuffer() sql.Row {
	return sql.NewMemoryRow(nil)
}

func (e *First) Type() sql.Type {
	return e.Child.Type()
}

func (e *First) Update(buffer, row sql.Row) {
	mr := buffer.(sql.MemoryRow)
	if mr[0] == nil {
		mr[0] = e.Child.Eval(row)
	}
}

func (e *First) Merge(buffer, partial sql.Row) {
	mb := buffer.(sql.MemoryRow)
	if mb[0] == nil {
		mp := partial.(sql.MemoryRow)
		mb[0] = mp[0]
	}
}

func (e *First) Eval(buffer sql.Row) interface{} {
	return buffer.Fields()[0]
}

func getInt32At(row sql.Row, i int) int32 {
	f := row.Fields()
	return f[i].(int32)
}
