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
	"strings"

	"github.com/shopspring/decimal"

	"github.com/dolthub/vitess/go/sqltypes"
	"github.com/dolthub/vitess/go/vt/proto/query"
)

var systemBoolValueType = reflect.TypeOf(int8(0))

// systemBoolType is an internal boolean type ONLY for system variables.
type systemBoolType struct {
	varName string
}

var _ SystemVariableType = systemBoolType{}

// NewSystemBoolType returns a new systemBoolType.
func NewSystemBoolType(varName string) SystemVariableType {
	return systemBoolType{varName}
}

// Compare implements Type interface.
func (t systemBoolType) Compare(a interface{}, b interface{}) (int, error) {
	as, err := t.Convert(a)
	if err != nil {
		return 0, err
	}
	bs, err := t.Convert(b)
	if err != nil {
		return 0, err
	}
	ai := as.(int8)
	bi := bs.(int8)

	if ai == bi {
		return 0, nil
	}
	if ai < bi {
		return -1, nil
	}
	return 1, nil
}

// Convert implements Type interface.
func (t systemBoolType) Convert(v interface{}) (interface{}, error) {
	// Nil values are not accepted
	switch value := v.(type) {
	case bool:
		if value {
			return int8(1), nil
		}
		return int8(0), nil
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
		if value == 0 || value == 1 {
			return int8(value), nil
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
		if value >= float64(math.MinInt64) && value <= float64(math.MaxInt64) {
			if value == float64(int64(value)) {
				intVal := int64(value)
				if intVal >= math.MinInt8 && intVal <= math.MaxInt8 {
					return int8(intVal), nil
				}
			}
		}
		return nil, ErrInvalidSystemVariableValue.New(t.varName, v)
	case decimal.Decimal:
		f, _ := value.Float64()
		return t.Convert(f)
	case decimal.NullDecimal:
		if value.Valid {
			f, _ := value.Decimal.Float64()
			return t.Convert(f)
		}
	case string:
		switch strings.ToLower(value) {
		case "on", "true":
			return int8(1), nil
		case "off", "false":
			return int8(0), nil
		}
	}

	return nil, ErrInvalidSystemVariableValue.New(t.varName, v)
}

// MustConvert implements the Type interface.
func (t systemBoolType) MustConvert(v interface{}) interface{} {
	value, err := t.Convert(v)
	if err != nil {
		panic(err)
	}
	return value
}

// Equals implements the Type interface.
func (t systemBoolType) Equals(otherType Type) bool {
	if ot, ok := otherType.(systemBoolType); ok {
		return t.varName == ot.varName
	}
	return false
}

// MaxTextResponseByteLength implements the Type interface
func (t systemBoolType) MaxTextResponseByteLength() uint32 {
	// system types are not sent directly across the wire
	return 0
}

// Promote implements the Type interface.
func (t systemBoolType) Promote() Type {
	return t
}

// SQL implements Type interface.
func (t systemBoolType) SQL(ctx *Context, dest []byte, v interface{}) (sqltypes.Value, error) {
	if v == nil {
		return sqltypes.NULL, nil
	}

	v, err := t.Convert(v)
	if err != nil {
		return sqltypes.Value{}, err
	}

	stop := len(dest)
	dest = strconv.AppendInt(dest, int64(v.(int8)), 10)
	val := dest[stop:]

	return sqltypes.MakeTrusted(t.Type(), val), nil
}

// String implements Type interface.
func (t systemBoolType) String() string {
	return "system_bool"
}

// Type implements Type interface.
func (t systemBoolType) Type() query.Type {
	return sqltypes.Int8
}

// ValueType implements Type interface.
func (t systemBoolType) ValueType() reflect.Type {
	return systemBoolValueType
}

// Zero implements Type interface.
func (t systemBoolType) Zero() interface{} {
	return int8(0)
}

// EncodeValue implements SystemVariableType interface.
func (t systemBoolType) EncodeValue(val interface{}) (string, error) {
	expectedVal, ok := val.(int8)
	if !ok {
		return "", ErrSystemVariableCodeFail.New(val, t.String())
	}
	if expectedVal == 0 {
		return "0", nil
	}
	return "1", nil
}

// DecodeValue implements SystemVariableType interface.
func (t systemBoolType) DecodeValue(val string) (interface{}, error) {
	if val == "0" {
		return int8(0), nil
	} else if val == "1" {
		return int8(1), nil
	}
	return nil, ErrSystemVariableCodeFail.New(val, t.String())
}
