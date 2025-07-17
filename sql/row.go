// Copyright 2020-2021 Dolthub, Inc.
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
	"fmt"
	"io"
	"strings"

	"github.com/dolthub/vitess/go/vt/proto/query"

	"github.com/dolthub/go-mysql-server/sql/values"
)

// Row is a tuple of values.
type Row []interface{}

// NewRow creates a row from the given values.
func NewRow(values ...interface{}) Row {
	row := make([]interface{}, len(values))
	copy(row, values)
	return row
}

// Copy creates a new row with the same values as the current one.
func (r Row) Copy() Row {
	return NewRow(r...)
}

// Append appends all the values in r2 to this row and returns the result
func (r Row) Append(r2 Row) Row {
	row := make(Row, len(r)+len(r2))
	copy(row, r)
	copy(row[len(r):], r2)
	return row
}

// Equals checks whether two rows are equal given a schema.
func (r Row) Equals(row Row, schema Schema) (bool, error) {
	if len(row) != len(r) || len(row) != len(schema) {
		return false, nil
	}

	for i, colLeft := range r {
		colRight := row[i]
		cmp, err := schema[i].Type.Compare(colLeft, colRight)
		if err != nil {
			return false, err
		}
		if cmp != 0 {
			return false, nil
		}
	}

	return true, nil
}

// FormatRow returns a formatted string representing this row's values
func FormatRow(row Row) string {
	var sb strings.Builder
	sb.WriteRune('[')
	for i, v := range row {
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(fmt.Sprintf("%v", v))
	}
	sb.WriteRune(']')
	return sb.String()
}

// RowIter is an iterator that produces rows.
// TODO: most row iters need to be Disposable for CachedResult safety
type RowIter interface {
	// Next retrieves the next row. It will return io.EOF if it's the last row.
	// After retrieving the last row, Close will be automatically closed.
	Next(ctx *Context) (Row, error)
	Closer
}

// RowIter2 is an iterator that fills a row frame buffer with rows from its source
type RowIter2 interface {
	RowIter

	// Next2 produces the next row, and stores it in the RowFrame provided.
	// It will return io.EOF if it's the last row. After retrieving the
	// last row, Close will be automatically called.
	Next2(ctx *Context, frame *RowFrame) error
}

// RowIterToRows converts a row iterator to a slice of rows.
func RowIterToRows(ctx *Context, sch Schema, i RowIter) ([]Row, error) {
	// Check for nil parameters
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	if i == nil {
		return nil, fmt.Errorf("iterator cannot be nil")
	}

	// Handle RowIter2 type if available
	if ri2, ok := i.(RowIterTypeSelector); ok && ri2.IsNode2() && sch != nil {
		return RowIter2ToRows(ctx, sch, ri2.(RowIter2))
	}

	var rows []Row
	for {
		row, err := i.Next(ctx)
		if err == io.EOF {
			break
		}

		if err != nil {
			closeErr := i.Close(ctx)
			if closeErr != nil {
				// Log the close error but return the original error
				ctx.Warn(0, "error closing iterator: %s", closeErr.Error())
			}
			return nil, err
		}

		rows = append(rows, row)
	}

	// Handle any errors from closing the iterator
	if err := i.Close(ctx); err != nil {
		return nil, err
	}

	return rows, nil
}

func RowIter2ToRows(ctx *Context, sch Schema, i RowIter2) ([]Row, error) {
	// Validate parameters to prevent potential panics
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	if sch == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}
	if i == nil {
		return nil, fmt.Errorf("iterator cannot be nil")
	}

	var rows []Row

	for {
		f := NewRowFrame()
		if f == nil {
			return nil, fmt.Errorf("failed to create row frame")
		}

		err := i.Next2(ctx, f)

		if err == io.EOF {
			break
		}

		if err != nil {
			// Make sure we close the iterator on error
			closeErr := i.Close(ctx)
			if closeErr != nil {
				// Log the close error but return the original error
				ctx.Warn(0, "error closing iterator: %s", closeErr.Error())
			}
			return nil, err
		}

		// Check if Row2() returns a valid Row2 object
		r2 := f.Row2()
		if r2 == nil {
			return nil, fmt.Errorf("invalid row frame: Row2() returned nil")
		}

		rows = append(rows, rowFromRow2(sch, r2))
	}

	// Handle any errors from closing the iterator
	if err := i.Close(ctx); err != nil {
		return nil, err
	}

	return rows, nil
}

func rowFromRow2(sch Schema, r Row2) Row {
	row := make(Row, len(sch))
	for i, col := range sch {
		// Handle the case where the field might be nil or invalid
		field := r.GetField(i)
		if field.IsNull() || field.Val == nil {
			row[i] = nil
			continue
		}

		switch col.Type.Type() {
		case query.Type_INT8:
			val, err := values.ReadInt8(field.Val)
			if err != nil {
				// In case of error, use zero value
				row[i] = int8(0)
			} else {
				row[i] = val
			}
		case query.Type_UINT8:
			val, err := values.ReadUint8(field.Val)
			if err != nil {
				row[i] = uint8(0)
			} else {
				row[i] = val
			}
		case query.Type_INT16:
			val, err := values.ReadInt16(field.Val)
			if err != nil {
				row[i] = int16(0)
			} else {
				row[i] = val
			}
		case query.Type_UINT16:
			val, err := values.ReadUint16(field.Val)
			if err != nil {
				row[i] = uint16(0)
			} else {
				row[i] = val
			}
		case query.Type_INT32:
			val, err := values.ReadInt32(field.Val)
			if err != nil {
				row[i] = int32(0)
			} else {
				row[i] = val
			}
		case query.Type_UINT32:
			val, err := values.ReadUint32(field.Val)
			if err != nil {
				row[i] = uint32(0)
			} else {
				row[i] = val
			}
		case query.Type_INT64:
			val, err := values.ReadInt64(field.Val)
			if err != nil {
				row[i] = int64(0)
			} else {
				row[i] = val
			}
		case query.Type_UINT64:
			val, err := values.ReadUint64(field.Val)
			if err != nil {
				row[i] = uint64(0)
			} else {
				row[i] = val
			}
		case query.Type_FLOAT32:
			val, err := values.ReadFloat32(field.Val)
			if err != nil {
				row[i] = float32(0)
			} else {
				row[i] = val
			}
		case query.Type_FLOAT64:
			val, err := values.ReadFloat64(field.Val)
			if err != nil {
				row[i] = float64(0)
			} else {
				row[i] = val
			}
		case query.Type_TEXT, query.Type_VARCHAR, query.Type_CHAR:
			row[i] = values.ReadString(field.Val, values.ByteOrderCollation)
		case query.Type_BLOB, query.Type_VARBINARY, query.Type_BINARY:
			row[i] = values.ReadBytes(field.Val, values.ByteOrderCollation)
		case query.Type_INT24:
			val, err := values.ReadInt24(field.Val)
			if err != nil {
				row[i] = int32(0)
			} else {
				row[i] = val
			}
		case query.Type_UINT24:
			val, err := values.ReadUint24(field.Val)
			if err != nil {
				row[i] = uint32(0)
			} else {
				row[i] = val
			}
		case query.Type_BIT:
			fallthrough
		case query.Type_ENUM:
			fallthrough
		case query.Type_SET:
			fallthrough
		case query.Type_TUPLE:
			fallthrough
		case query.Type_GEOMETRY:
			fallthrough
		case query.Type_JSON:
			fallthrough
		case query.Type_EXPRESSION:
			fallthrough
		case query.Type_TIMESTAMP:
			fallthrough
		case query.Type_DATE:
			fallthrough
		case query.Type_TIME:
			fallthrough
		case query.Type_DATETIME:
			fallthrough
		case query.Type_YEAR:
			fallthrough
		case query.Type_DECIMAL:
			// Instead of panicking, set to nil for unimplemented types
			row[i] = nil
		default:
			// Also set to nil for unknown types
			row[i] = nil
		}
	}
	return row
}

// NodeToRows converts a node to a slice of rows.
func NodeToRows(ctx *Context, n Node) ([]Row, error) {
	// Check for nil parameters
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	if n == nil {
		return nil, fmt.Errorf("node cannot be nil")
	}

	i, err := n.RowIter(ctx, nil)
	if err != nil {
		return nil, err
	}

	if i == nil {
		return nil, fmt.Errorf("node returned nil iterator")
	}

	return RowIterToRows(ctx, nil, i)
}

// RowsToRowIter creates a RowIter that iterates over the given rows.
func RowsToRowIter(rows ...Row) RowIter {
	return &sliceRowIter{rows: rows}
}

type sliceRowIter struct {
	rows []Row
	idx  int
}

func (i *sliceRowIter) Next(*Context) (Row, error) {
	if i.idx >= len(i.rows) {
		return nil, io.EOF
	}

	r := i.rows[i.idx]
	i.idx++
	return r.Copy(), nil
}

func (i *sliceRowIter) Close(*Context) error {
	i.rows = nil
	return nil
}
