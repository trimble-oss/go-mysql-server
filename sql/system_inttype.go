// Copyright 2021 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sql

import (
	"math"
	"reflect"
	"strconv"

	"github.com/shopspring/decimal"

	"github.com/dolthub/vitess/go/sqltypes"
	"github.com/dolthub/vitess/go/vt/proto/query"
)

var systemIntValueType = reflect.TypeOf(int64(0))

// systemIntType is an internal integer type ONLY for system variables.
type systemIntType struct {
	varName     string
	lowerbound  int64
	upperbound  int64
	negativeOne bool
}

var _ SystemVariableType = systemIntType{}

// NewSystemIntType returns a new systemIntType.
func NewSystemIntType(varName string, lowerbound, upperbound int64, negativeOne bool) SystemVariableType {
	return systemIntType{varName, lowerbound, upperbound, negativeOne}
}

// Compare implements Type interface.
func (t systemIntType) Compare(a interface{}, b interface{}) (int, error) {
	as, err := t.Convert(a)
	if err != nil {
		return 0, err
	}
	bs, err := t.Convert(b)
	if err != nil {
		return 0, err
	}
	ai := as.(int64)
	bi := bs.(int64)

	if ai == bi {
		return 0, nil
	}
	if ai < bi {
		return -1, nil
	}
	return 1, nil
}

// Convert implements Type interface.
func (t systemIntType) Convert(v interface{}) (interface{}, error) {
	// String nor nil values are accepted
	switch value := v.(type) {
	case int:
		return t.Convert(int64(value))
	case uint:
		return t.Convert(int64(value))
	case int8:
		return t.Convert(int64(value))
	case uint8:
		return t.Convert(int64(value))
	case int16:
		return t.Convert(int64(value))
	case uint16:
		return t.Convert(int64(value))
	case int32:
		return t.Convert(int64(value))
	case uint32:
		return t.Convert(int64(value))
	case int64:
		if value >= t.lowerbound && value <= t.upperbound {
			return value, nil
		}
		if t.negativeOne && value == -1 {
			return value, nil
		}
	case uint64:
		if value <= math.MaxInt64 {
			return t.Convert(int64(value))
		}
		return nil, ErrInvalidSystemVariableValue.New(t.varName, v)
	case float32:
		return t.Convert(float64(value))
	case float64:
		// Float values aren't truly accepted, but the engine will give them when it should give ints.
		// Therefore, if the float doesn't have a fractional portion, we treat it as an int.
		if value == float64(int64(value)) {
			if value >= float64(t.lowerbound) && value <= float64(t.upperbound) {
				if value >= float64(math.MinInt64) && value <= float64(math.MaxInt64) {
					intVal := int64(value)
					return t.Convert(intVal)
				}
			}
			return nil, ErrInvalidSystemVariableValue.New(t.varName, v)
		}
	case decimal.Decimal:
		f, _ := value.Float64()
		return t.Convert(f)
	case decimal.NullDecimal:
		if value.Valid {
			f, _ := value.Decimal.Float64()
			return t.Convert(f)
		}
	}

	return nil, ErrInvalidSystemVariableValue.New(t.varName, v)
}

// MustConvert implements the Type interface.
func (t systemIntType) MustConvert(v interface{}) interface{} {
	value, err := t.Convert(v)
	if err != nil {
		panic(err)
	}
	return value
}

// Equals implements the Type interface.
func (t systemIntType) Equals(otherType Type) bool {
	if ot, ok := otherType.(systemIntType); ok {
		return t.varName == ot.varName && t.lowerbound == ot.lowerbound && t.upperbound == ot.upperbound && t.negativeOne == ot.negativeOne
	}
	return false
}

// MaxTextResponseByteLength implements the Type interface
func (t systemIntType) MaxTextResponseByteLength() uint32 {
	// system types are not sent directly across the wire
	return 0
}

// Promote implements the Type interface.
func (t systemIntType) Promote() Type {
	return t
}

// SQL implements Type interface.
func (t systemIntType) SQL(ctx *Context, dest []byte, v interface{}) (sqltypes.Value, error) {
	if v == nil {
		return sqltypes.NULL, nil
	}

	v, err := t.Convert(v)
	if err != nil {
		return sqltypes.Value{}, err
	}

	stop := len(dest)
	dest = strconv.AppendInt(dest, v.(int64), 10)
	val := dest[stop:]

	return sqltypes.MakeTrusted(t.Type(), val), nil
}

// String implements Type interface.
func (t systemIntType) String() string {
	return "system_int"
}

// Type implements Type interface.
func (t systemIntType) Type() query.Type {
	return sqltypes.Int64
}

// ValueType implements Type interface.
func (t systemIntType) ValueType() reflect.Type {
	return systemIntValueType
}

// Zero implements Type interface.
func (t systemIntType) Zero() interface{} {
	return int64(0)
}

// EncodeValue implements SystemVariableType interface.
func (t systemIntType) EncodeValue(val interface{}) (string, error) {
	expectedVal, ok := val.(int64)
	if !ok {
		return "", ErrSystemVariableCodeFail.New(val, t.String())
	}
	return strconv.FormatInt(expectedVal, 10), nil
}

// DecodeValue implements SystemVariableType interface.
func (t systemIntType) DecodeValue(val string) (interface{}, error) {
	parsedVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return nil, err
	}
	if parsedVal >= t.lowerbound && parsedVal <= t.upperbound {
		return parsedVal, nil
	}
	return nil, ErrSystemVariableCodeFail.New(val, t.String())
}
