package expression

import "github.com/mvader/gitql/sql"

type Literal struct {
	value     interface{}
	fieldType sql.Type
	name      string
}

func NewLiteral(value interface{}, fieldType sql.Type) *Literal {
	return &Literal{
		value:     value,
		fieldType: fieldType,
		name:      "literal_" + fieldType.Name(),
	}
}

func (p Literal) Type() sql.Type {
	return p.fieldType
}

func (p Literal) Eval(row sql.Row) interface{} {
	return p.value
}

func (p Literal) Name() string {
	return p.name
}
