// Copyright 2022 Dolthub, Inc.
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

package values

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

type Type struct {
	Enc      Encoding
	Coll     Collation
	Nullable bool
}

type ByteSize uint16

const (
	Int8Size    ByteSize = 1
	Uint8Size   ByteSize = 1
	Int16Size   ByteSize = 2
	Uint16Size  ByteSize = 2
	Int24Size   ByteSize = 3
	Uint24Size  ByteSize = 3
	Int32Size   ByteSize = 4
	Uint32Size  ByteSize = 4
	Int48Size   ByteSize = 6
	Uint48Size  ByteSize = 6
	Int64Size   ByteSize = 8
	Uint64Size  ByteSize = 8
	Float32Size ByteSize = 4
	Float64Size ByteSize = 8
)

const maxUint48 = uint64(1<<48 - 1)
const maxUint24 = uint32(1<<24 - 1)

type Collation uint16

const (
	ByteOrderCollation Collation = 0
)

type Encoding uint8

// Constant Size Encodings
const (
	NullEnc    Encoding = 0
	Int8Enc    Encoding = 1
	Uint8Enc   Encoding = 2
	Int16Enc   Encoding = 3
	Uint16Enc  Encoding = 4
	Int24Enc   Encoding = 5
	Uint24Enc  Encoding = 6
	Int32Enc   Encoding = 7
	Uint32Enc  Encoding = 8
	Int64Enc   Encoding = 9
	Uint64Enc  Encoding = 10
	Float32Enc Encoding = 11
	Float64Enc Encoding = 12

	// TODO
	//  TimeEnc
	//  TimestampEnc
	//  DateEnc
	//  TimeEnc
	//  DatetimeEnc
	//  YearEnc

	sentinel Encoding = 127
)

// Variable Size Encodings
const (
	StringEnc Encoding = 128
	BytesEnc  Encoding = 129

	// TODO
	//  DecimalEnc
	//  BitEnc
	//  CharEnc
	//  VarCharEnc
	//  TextEnc
	//  BinaryEnc
	//  VarBinaryEnc
	//  BlobEnc
	//  JSONEnc
	//  EnumEnc
	//  SetEnc
	//  ExpressionEnc
	//  GeometryEnc
)

func ReadBool(val []byte) bool {
	expectSize(val, Int8Size)
	return val[0] == 1
}
func ReadInt8(val []byte) (int8, error) {
	expectSize(val, Int8Size)
	if len(val) != int(Int8Size) {
		return 0, fmt.Errorf("invalid size for int8: expected %d, got %d", Int8Size, len(val))
	}
	return int8(val[0]), nil
}

func ReadUint8(val []byte) (uint8, error) {
	expectSize(val, Uint8Size)
	if len(val) != int(Uint8Size) {
		return 0, fmt.Errorf("invalid size for uint8: expected %d, got %d", Uint8Size, len(val))
	}
	return val[0], nil
}

func ReadInt16(val []byte) (int16, error) {
	expectSizeErr := expectSize(val, Int16Size)
	if expectSizeErr != nil {
		return 0, expectSizeErr
	}
	return int16(binary.LittleEndian.Uint16(val)), nil
}

func ReadUint16(val []byte) (uint16, error) {
	expectSizeErr := expectSize(val, Uint16Size)
	if expectSizeErr != nil {
		return 0, expectSizeErr
	}
	return binary.LittleEndian.Uint16(val), nil
}

// Updated ReadInt24 and ReadUint24 to return value and error
func ReadInt24(val []byte) (int32, error) {
	expectSizeErr := expectSize(val, Int24Size)
	if expectSizeErr != nil {
		return 0, expectSizeErr
	}
	var tmp [4]byte
	// copy |val| to |tmp|
	tmp[3], tmp[2] = val[3], val[2]
	tmp[1], tmp[0] = val[1], val[0]
	return int32(binary.LittleEndian.Uint32(tmp[:])), nil
}

func ReadUint24(val []byte) (uint32, error) {
	expectSizeErr := expectSize(val, Int24Size)
	if expectSizeErr != nil {
		return 0, expectSizeErr
	}
	var tmp [4]byte
	// copy |val| to |tmp|
	tmp[3], tmp[2] = val[3], val[2]
	tmp[1], tmp[0] = val[1], val[0]
	return binary.LittleEndian.Uint32(tmp[:]), nil
}

func ReadInt32(val []byte) (int32, error) {
	expectSizeErr := expectSize(val, Int32Size)
	if expectSizeErr != nil {
		return 0, expectSizeErr
	}
	return int32(binary.LittleEndian.Uint32(val)), nil
}

func ReadUint32(val []byte) (uint32, error) {
	expectSizeErr := expectSize(val, Uint32Size)
	if expectSizeErr != nil {
		return 0, expectSizeErr
	}
	return binary.LittleEndian.Uint32(val), nil
}

func ReadInt48(val []byte) (i int64) {
	expectSize(val, Int48Size)
	var tmp [8]byte
	// copy |val| to |tmp|
	tmp[5], tmp[4] = val[5], val[4]
	tmp[3], tmp[2] = val[3], val[2]
	tmp[1], tmp[0] = val[1], val[0]
	i = int64(binary.LittleEndian.Uint64(tmp[:]))
	return
}

func ReadUint48(val []byte) (u uint64) {
	expectSize(val, Uint48Size)
	var tmp [8]byte
	// copy |val| to |tmp|
	tmp[5], tmp[4] = val[5], val[4]
	tmp[3], tmp[2] = val[3], val[2]
	tmp[1], tmp[0] = val[1], val[0]
	u = binary.LittleEndian.Uint64(tmp[:])
	return
}

func ReadInt64(val []byte) (int64, error) {
	expectSizeErr := expectSize(val, Int64Size)
	if expectSizeErr != nil {
		return 0, expectSizeErr
	}
	return int64(binary.LittleEndian.Uint64(val)), nil
}

func ReadUint64(val []byte) (uint64, error) {
	expectSizeErr := expectSize(val, Uint64Size)
	if expectSizeErr != nil {
		return 0, expectSizeErr
	}
	return binary.LittleEndian.Uint64(val), nil
}

func ReadFloat32(val []byte) (float32, error) {
	expectSizeErr := expectSize(val, Float32Size)
	if expectSizeErr != nil {
		return 0, expectSizeErr
	}
	uintVal, err := ReadUint32(val)
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(uintVal), nil
}

func ReadFloat64(val []byte) (float64, error) {
	expectSizeErr := expectSize(val, Float64Size)
	if expectSizeErr != nil {
		return 0, expectSizeErr
	}
	uintVal, err := ReadUint64(val)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(uintVal), nil
}

func ReadString(val []byte, coll Collation) string {
	// todo: fix allocation
	return string(val)
}

func ReadBytes(val []byte, coll Collation) []byte {
	// todo: fix collation
	return val
}

func writeBool(buf []byte, val bool) {
	expectSize(buf, 1)
	if val {
		buf[0] = byte(1)
	} else {
		buf[0] = byte(0)
	}
}

func WriteInt8(buf []byte, val int8) []byte {
	expectSize(buf, Int8Size)
	buf[0] = byte(val)
	return buf
}

func WriteUint8(buf []byte, val uint8) []byte {
	expectSize(buf, Uint8Size)
	buf[0] = byte(val)
	return buf
}

func WriteInt16(buf []byte, val int16) []byte {
	expectSize(buf, Int16Size)
	binary.LittleEndian.PutUint16(buf, uint16(val))
	return buf
}

func WriteUint16(buf []byte, val uint16) []byte {
	expectSize(buf, Uint16Size)
	binary.LittleEndian.PutUint16(buf, val)
	return buf
}

func WriteInt24(buf []byte, val int32) []byte {
	expectSize(buf, Int24Size)

	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], uint32(val))
	// copy |tmp| to |buf|
	buf[2], buf[1], buf[0] = tmp[2], tmp[1], tmp[0]
	return buf

	binary.LittleEndian.PutUint16(buf, uint16(val))
	return buf
}

func WriteUint24(buf []byte, val uint32) ([]byte, error) {
	expectSizeErr := expectSize(buf, Uint24Size)
	if expectSizeErr != nil {
		return nil, expectSizeErr
	}
	if val > maxUint24 {
		return nil, fmt.Errorf("uint value %d exceeds max uint24 %d", val, maxUint24)
	}

	var tmp [4]byte
	binary.LittleEndian.PutUint32(tmp[:], uint32(val))
	// copy |tmp| to |buf|
	buf[2], buf[1], buf[0] = tmp[2], tmp[1], tmp[0]
	return buf, nil
}

func WriteInt32(buf []byte, val int32) []byte {
	expectSize(buf, Int32Size)
	binary.LittleEndian.PutUint32(buf, uint32(val))
	return buf
}

func WriteUint32(buf []byte, val uint32) []byte {
	expectSize(buf, Uint32Size)
	binary.LittleEndian.PutUint32(buf, val)
	return buf
}

func WriteUint48(buf []byte, u uint64) ([]byte, error) {
	expectSizeErr := expectSize(buf, Uint48Size)
	if expectSizeErr != nil {
		return nil, expectSizeErr
	}
	if u > maxUint48 {
		return nil, fmt.Errorf("uint value %d exceeds max uint48 %d", u, maxUint48)
	}
	var tmp [8]byte
	binary.LittleEndian.PutUint64(tmp[:], u)
	// copy |tmp| to |buf|
	buf[5], buf[4] = tmp[5], tmp[4]
	buf[3], buf[2] = tmp[3], tmp[2]
	buf[1], buf[0] = tmp[1], tmp[0]
	return buf, nil
}

func WriteInt64(buf []byte, val int64) []byte {
	expectSize(buf, Int64Size)
	binary.LittleEndian.PutUint64(buf, uint64(val))
	return buf
}

func WriteUint64(buf []byte, val uint64) []byte {
	expectSize(buf, Uint64Size)
	binary.LittleEndian.PutUint64(buf, val)
	return buf
}

func WriteFloat32(buf []byte, val float32) []byte {
	expectSize(buf, Float32Size)
	binary.LittleEndian.PutUint32(buf, math.Float32bits(val))
	return buf
}

func WriteFloat64(buf []byte, val float64) []byte {
	expectSize(buf, Float64Size)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(val))
	return buf
}

func WriteString(buf []byte, val string, coll Collation) []byte {
	// todo: fix collation
	expectSize(buf, ByteSize(len(val)))
	copy(buf, val)
	return buf
}

func WriteBytes(buf, val []byte, coll Collation) []byte {
	// todo: fix collation
	expectSize(buf, ByteSize(len(val)))
	copy(buf, val)
	return buf
}

func expectSize(buf []byte, sz ByteSize) error {
	if ByteSize(len(buf)) != sz {
		return fmt.Errorf("byte slice is not of expected size: expected %d, got %d", sz, len(buf))
	}
	return nil
}

// Updated compare function to handle errors from all Read* functions
func compare(typ Type, left, right []byte) (int, error) {
	// order NULLs last
	if left == nil {
		if right == nil {
			return 0, nil
		} else {
			return 1, nil
		}
	} else if right == nil {
		if left == nil {
			return 0, nil
		} else {
			return -1, nil
		}
	}

	switch typ.Enc {
	case Int8Enc:
		lval, err := ReadInt8(left)
		if err != nil {
			return 0, err
		}
		rval, err := ReadInt8(right)
		if err != nil {
			return 0, err
		}
		return compareInt8(lval, rval), nil
	case Uint8Enc:
		lval, err := ReadUint8(left)
		if err != nil {
			return 0, err
		}
		rval, err := ReadUint8(right)
		if err != nil {
			return 0, err
		}
		return compareUint8(lval, rval), nil
	case Int16Enc:
		lval, err := ReadInt16(left)
		if err != nil {
			return 0, err
		}
		rval, err := ReadInt16(right)
		if err != nil {
			return 0, err
		}
		return compareInt16(lval, rval), nil
	case Uint16Enc:
		lval, err := ReadUint16(left)
		if err != nil {
			return 0, err
		}
		rval, err := ReadUint16(right)
		if err != nil {
			return 0, err
		}
		return compareUint16(lval, rval), nil
	case Int32Enc:
		lval, err := ReadInt32(left)
		if err != nil {
			return 0, err
		}
		rval, err := ReadInt32(right)
		if err != nil {
			return 0, err
		}
		return compareInt32(lval, rval), nil
	case Uint32Enc:
		lval, err := ReadUint32(left)
		if err != nil {
			return 0, err
		}
		rval, err := ReadUint32(right)
		if err != nil {
			return 0, err
		}
		return compareUint32(lval, rval), nil
	case Int64Enc:
		lval, err := ReadInt64(left)
		if err != nil {
			return 0, err
		}
		rval, err := ReadInt64(right)
		if err != nil {
			return 0, err
		}
		return compareInt64(lval, rval), nil
	case Uint64Enc:
		lval, err := ReadUint64(left)
		if err != nil {
			return 0, err
		}
		rval, err := ReadUint64(right)
		if err != nil {
			return 0, err
		}
		return compareUint64(lval, rval), nil
	case Float32Enc:
		lval, err := ReadFloat32(left)
		if err != nil {
			return 0, err
		}
		rval, err := ReadFloat32(right)
		if err != nil {
			return 0, err
		}
		return compareFloat32(lval, rval), nil
	case Float64Enc:
		lval, err := ReadFloat64(left)
		if err != nil {
			return 0, err
		}
		rval, err := ReadFloat64(right)
		if err != nil {
			return 0, err
		}
		return compareFloat64(lval, rval), nil
	case StringEnc:
		lval := ReadString(left, typ.Coll)
		rval := ReadString(right, typ.Coll)
		return compareString(lval, rval, typ.Coll), nil
	case BytesEnc:
		lval := ReadBytes(left, typ.Coll)
		rval := ReadBytes(right, typ.Coll)
		return compareBytes(lval, rval, typ.Coll), nil
	default:
		return 0, fmt.Errorf("unknown encoding")
	}
}

// false is less than true
func compareBool(l, r bool) int {
	if l == r {
		return 0
	}
	if !l && r {
		return -1
	}
	return 1
}

func compareInt8(l, r int8) int {
	if l == r {
		return 0
	} else if l < r {
		return -1
	} else {
		return 1
	}
}

func compareUint8(l, r uint8) int {
	if l == r {
		return 0
	} else if l < r {
		return -1
	} else {
		return 1
	}
}

func compareInt16(l, r int16) int {
	if l == r {
		return 0
	} else if l < r {
		return -1
	} else {
		return 1
	}
}

func compareUint16(l, r uint16) int {
	if l == r {
		return 0
	} else if l < r {
		return -1
	} else {
		return 1
	}
}

func compareInt32(l, r int32) int {
	if l == r {
		return 0
	} else if l < r {
		return -1
	} else {
		return 1
	}
}

func compareUint32(l, r uint32) int {
	if l == r {
		return 0
	} else if l < r {
		return -1
	} else {
		return 1
	}
}

func compareInt64(l, r int64) int {
	if l == r {
		return 0
	} else if l < r {
		return -1
	} else {
		return 1
	}
}

func compareUint64(l, r uint64) int {
	if l == r {
		return 0
	} else if l < r {
		return -1
	} else {
		return 1
	}
}

func compareFloat32(l, r float32) int {
	if l == r {
		return 0
	} else if l < r {
		return -1
	} else {
		return 1
	}
}

func compareFloat64(l, r float64) int {
	if l == r {
		return 0
	} else if l < r {
		return -1
	} else {
		return 1
	}
}

func compareString(l, r string, coll Collation) int {
	return bytes.Compare([]byte(l), []byte(r))
}

func compareBytes(l, r []byte, coll Collation) int {
	return bytes.Compare(l, r)
}
