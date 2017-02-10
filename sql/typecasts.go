package sql

import (
	"fmt"
	"strconv"
	"time"
)

// TypeCastRegistry holds cast functions between types.
type TypeCastRegistry struct {
	s []*typeCast
}

// NewTypeCastRegistry returns a new TypeCastRegistry.
func NewTypeCastRegistry() *TypeCastRegistry {
	return &TypeCastRegistry{}
}

// RegisterTypeCast registers a new cast function between the given types.
func (r *TypeCastRegistry) RegisterTypeCast(from, to Type, f func(interface{}) interface{}) {
	r.s = append(r.s, &typeCast{from, to, f})
}

// TypeCast returns a cast function the given types.
func (r *TypeCastRegistry) TypeCast(from, to Type) (func(interface{}) interface{}, error) {
	if from == to {
		return identity, nil
	}

	for _, tc := range r.s {
		if tc.From == from && tc.To == to {
			return tc.Cast, nil
		}
	}

	return nil, fmt.Errorf("type cast does not exist from %s to %s", from.Name(), to.Name())
}

type typeCast struct {
	From, To Type
	Cast     func(interface{}) interface{}
}

func RegisterDefaultTypeCasts(r *TypeCastRegistry) {
	for _, tc := range typeCastFunctions {
		r.RegisterTypeCast(tc.From, tc.To, tc.Cast)
	}
}

var typeCastFunctions []*typeCast = []*typeCast{
	{From: Integer, To: BigInteger, Cast: func(v interface{}) interface{} { return int64(v.(int32)) }},
	{From: Integer, To: Float, Cast: func(v interface{}) interface{} { return float64(v.(int32)) }},
	{From: Integer, To: Boolean, Cast: func(v interface{}) interface{} { return v.(int32) != 0 }},
	{From: Integer, To: String, Cast: func(v interface{}) interface{} { return strconv.FormatInt(int64(v.(int32)), 10) }},
	{From: Integer, To: TimestampWithTimezone, Cast: func(v interface{}) interface{} { return time.Unix(int64(v.(int32)), 0) }},

	{From: BigInteger, To: Integer, Cast: func(v interface{}) interface{} { return int32(v.(int64)) }},
	{From: BigInteger, To: Float, Cast: func(v interface{}) interface{} { return float64(v.(int64)) }},
	{From: BigInteger, To: Boolean, Cast: func(v interface{}) interface{} { return v.(int64) != 0 }},
	{From: BigInteger, To: String, Cast: func(v interface{}) interface{} { return strconv.FormatInt(v.(int64), 0) }},
	{From: BigInteger, To: TimestampWithTimezone, Cast: func(v interface{}) interface{} { return time.Unix(v.(int64), 0) }},

	//TODO: handle overflow for explicit cast, disallow implicit
	{From: Float, To: Integer, Cast: func(v interface{}) interface{} { return int32(v.(float64)) }},
	//TODO: handle overflow for explicit cast, disallow implicit
	{From: Float, To: BigInteger, Cast: func(v interface{}) interface{} { return int64(v.(float64)) }},
	{From: Float, To: Boolean, Cast: func(v interface{}) interface{} { return v.(float64) != 0 }},
	//TODO: add float formatting functions (and chose a default?)
	{From: Float, To: String, Cast: func(v interface{}) interface{} { return strconv.FormatFloat(v.(float64), 'g', -1, 64) }},

	{From: Boolean, To: Integer, Cast: func(v interface{}) interface{} { return btoi32(v) }},
	{From: Boolean, To: BigInteger, Cast: func(v interface{}) interface{} { return int64(btoi32(v)) }},
	{From: Boolean, To: Float, Cast: func(v interface{}) interface{} { return float64(btoi32(v)) }},
	{From: Boolean, To: String, Cast: func(v interface{}) interface{} { return strconv.FormatBool(v.(bool)) }},

	{From: String, To: Integer, Cast: func(v interface{}) interface{} { return int32(atoi(32, v)) }},
	{From: String, To: BigInteger, Cast: func(v interface{}) interface{} { return atoi(64, v) }},
	{From: String, To: Float, Cast: func(v interface{}) interface{} { return atof(64, v) }},
	{From: String, To: Boolean, Cast: func(v interface{}) interface{} { return atob(v) }},
	//TODO: {From: String, To: TimestampWithTimezone},

	//TODO: {From: TimestampWithTimezone, To: Integer},
	//TODO: {From: TimestampWithTimezone, To: BigInteger},
	//TODO: {From: TimestampWithTimezone, To: Float},
	//TODO: {From: TimestampWithTimezone, To: Boolean},
	//TODO: {From: TimestampWithTimezone, To: String},
}

func identity(v interface{}) interface{} {
	return v
}

func btoi32(v interface{}) int32 {
	if v.(bool) {
		return int32(1)
	}

	return int32(0)
}

func atoi(size int, v interface{}) int64 {
	s := v.(string)
	i, err := strconv.ParseInt(s, 10, size)
	if err != nil {
		return 0
	}

	return i
}

func atof(size int, v interface{}) float64 {
	s := v.(string)
	f, err := strconv.ParseFloat(s, size)
	if err != nil {
		return float64(0)
	}

	return f
}

func atob(v interface{}) bool {
	s := v.(string)
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}

	return b
}
