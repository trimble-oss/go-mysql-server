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

package function

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/expression"
)

func TestAsGeoJSON(t *testing.T) {
	t.Run("convert point to geojson", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(expression.NewLiteral(sql.Point{X: 1, Y: 2}, sql.PointType{}))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [2]float64{1, 2}, "type": "Point"}}, v)
	})
	t.Run("convert linestring to geojson", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(expression.NewLiteral(sql.LineString{Points: []sql.Point{{X: 1, Y: 2}, {X: 3, Y: 4}}}, sql.LineStringType{}))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [][2]float64{{1, 2}, {3, 4}}, "type": "LineString"}}, v)
	})
	t.Run("convert polygon to geojson", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(expression.NewLiteral(sql.Polygon{Lines: []sql.LineString{{Points: []sql.Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 0}}}}}, sql.PolygonType{}))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [][][2]float64{{{0, 0}, {1, 0}, {1, 1}, {0, 0}}}, "type": "Polygon"}}, v)
	})
	t.Run("convert multipoint to geojson", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(expression.NewLiteral(sql.MultiPoint{Points: []sql.Point{{X: 1, Y: 2}, {X: 3, Y: 4}}}, sql.MultiPointType{}))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [][2]float64{{1, 2}, {3, 4}}, "type": "MultiPoint"}}, v)
	})
	t.Run("convert multilinestring to geojson", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(expression.NewLiteral(sql.MultiLineString{Lines: []sql.LineString{{Points: []sql.Point{{X: 1, Y: 2}, {X: 3, Y: 4}}}}}, sql.MultiLineStringType{}))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [][][2]float64{{{1, 2}, {3, 4}}}, "type": "MultiLineString"}}, v)
	})
	t.Run("convert multipolygon to geojson", func(t *testing.T) {
		require := require.New(t)
		line := sql.LineString{Points: []sql.Point{{X: 0, Y: 0}, {X: 1, Y: 2}, {X: 3, Y: 4}, {X: 0, Y: 0}}}
		poly := sql.Polygon{Lines: []sql.LineString{line}}
		f, err := NewAsGeoJSON(expression.NewLiteral(sql.MultiPolygon{Polygons: []sql.Polygon{poly}}, sql.MultiPolygonType{}))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [][][][2]float64{{{{0, 0}, {1, 2}, {3, 4}, {0, 0}}}}, "type": "MultiPolygon"}}, v)
	})
	t.Run("convert empty geometrycollection to geojson", func(t *testing.T) {
		require := require.New(t)
		g := sql.GeomColl{}
		f, err := NewAsGeoJSON(expression.NewLiteral(g, sql.GeomCollType{}))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)

		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"geometries": []interface{}{}, "type": "GeometryCollection"}}, v)
	})
	t.Run("convert geometrycollection to geojson", func(t *testing.T) {
		require := require.New(t)
		point := sql.Point{X: 1, Y: 2}
		line := sql.LineString{Points: []sql.Point{{X: 1, Y: 2}, {X: 3, Y: 4}}}
		poly := sql.Polygon{Lines: []sql.LineString{{Points: []sql.Point{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: 0}}}}}
		mpoint := sql.MultiPoint{Points: []sql.Point{point, point}}
		mline := sql.MultiLineString{Lines: []sql.LineString{line, line}}
		mpoly := sql.MultiPolygon{Polygons: []sql.Polygon{poly, poly}}
		gColl := sql.GeomColl{}
		g := sql.GeomColl{Geoms: []sql.GeometryValue{point, line, poly, mpoint, mline, mpoly, gColl}}
		f, err := NewAsGeoJSON(expression.NewLiteral(g, sql.GeomCollType{}))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)

		pointjson := map[string]interface{}{"coordinates": [2]float64{1, 2}, "type": "Point"}
		linejson := map[string]interface{}{"coordinates": [][2]float64{{1, 2}, {3, 4}}, "type": "LineString"}
		polyjson := map[string]interface{}{"coordinates": [][][2]float64{{{0, 0}, {1, 1}, {1, 0}, {0, 0}}}, "type": "Polygon"}
		mpointjson := map[string]interface{}{"coordinates": [][2]float64{{1, 2}, {1, 2}}, "type": "MultiPoint"}
		mlinejson := map[string]interface{}{"coordinates": [][][2]float64{{{1, 2}, {3, 4}}, {{1, 2}, {3, 4}}}, "type": "MultiLineString"}
		mpolyjson := map[string]interface{}{"coordinates": [][][][2]float64{{{{0, 0}, {1, 1}, {1, 0}, {0, 0}}}, {{{0, 0}, {1, 1}, {1, 0}, {0, 0}}}}, "type": "MultiPolygon"}
		mgeomjson := map[string]interface{}{"geometries": []interface{}{}, "type": "GeometryCollection"}

		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"geometries": []interface{}{pointjson, linejson, polyjson, mpointjson, mlinejson, mpolyjson, mgeomjson}, "type": "GeometryCollection"}}, v)
	})
	t.Run("convert point with floats to geojson", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(expression.NewLiteral(sql.Point{X: 123.45, Y: 5.6789}, sql.PointType{}))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [2]float64{123.45, 5.6789}, "type": "Point"}}, v)
	})
	t.Run("convert point with low precision", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(
			expression.NewLiteral(sql.Point{X: 0.123456789, Y: 0.987654321}, sql.PointType{}),
			expression.NewLiteral(3, sql.Int64),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [2]float64{0.123, 0.988}, "type": "Point"}}, v)
	})
	t.Run("convert point with high precision", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(
			expression.NewLiteral(sql.Point{X: 0.123456789, Y: 0.987654321}, sql.PointType{}),
			expression.NewLiteral(20, sql.Int64),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [2]float64{0.123456789, 0.987654321}, "type": "Point"}}, v)
	})
	t.Run("convert point with bounding box", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(
			expression.NewLiteral(sql.Point{X: 123.45678, Y: 456.789}, sql.PointType{}),
			expression.NewLiteral(2, sql.Int64),
			expression.NewLiteral(1, sql.Int64),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [2]float64{123.46, 456.79}, "type": "Point", "bbox": [4]float64{123.46, 456.79, 123.46, 456.79}}}, v)
	})
	t.Run("convert linestring with bounding box", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(
			expression.NewLiteral(sql.LineString{Points: []sql.Point{{X: 100, Y: 2}, {X: 1, Y: 200}}}, sql.LineStringType{}),
			expression.NewLiteral(2, sql.Int64),
			expression.NewLiteral(1, sql.Int64),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [][2]float64{{100, 2}, {1, 200}}, "type": "LineString", "bbox": [4]float64{1, 2, 100, 200}}}, v)
	})
	t.Run("convert polygon with bounding box", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(
			expression.NewLiteral(sql.Polygon{Lines: []sql.LineString{{Points: []sql.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 0, Y: 0}}}}}, sql.PolygonType{}),
			expression.NewLiteral(2, sql.Int64),
			expression.NewLiteral(1, sql.Int64),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [][][2]float64{{{0, 0}, {0, 1}, {1, 1}, {0, 0}}}, "type": "Polygon", "bbox": [4]float64{0, 0, 1, 1}}}, v)
	})
	t.Run("convert multipoint with bounding box", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(
			expression.NewLiteral(sql.MultiPoint{Points: []sql.Point{{X: 100, Y: 2}, {X: 1, Y: 200}}}, sql.MultiPointType{}),
			expression.NewLiteral(2, sql.Int64),
			expression.NewLiteral(1, sql.Int64),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [][2]float64{{100, 2}, {1, 200}}, "type": "MultiPoint", "bbox": [4]float64{1, 2, 100, 200}}}, v)
	})
	t.Run("convert multilinestring with bounding box", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(
			expression.NewLiteral(sql.MultiLineString{Lines: []sql.LineString{{Points: []sql.Point{{X: 1, Y: 2}, {X: 3, Y: 4}}}}}, sql.MultiLineStringType{}),
			expression.NewLiteral(2, sql.Int64),
			expression.NewLiteral(1, sql.Int64),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [][][2]float64{{{1, 2}, {3, 4}}}, "type": "MultiLineString", "bbox": [4]float64{1, 2, 3, 4}}}, v)
	})
	t.Run("convert multipolygon with bounding box", func(t *testing.T) {
		require := require.New(t)
		line := sql.LineString{Points: []sql.Point{{X: 0, Y: 0}, {X: 1, Y: 2}, {X: 3, Y: 4}, {X: 0, Y: 0}}}
		poly := sql.Polygon{Lines: []sql.LineString{line}}
		f, err := NewAsGeoJSON(
			expression.NewLiteral(sql.MultiPolygon{Polygons: []sql.Polygon{poly}}, sql.MultiPolygonType{}),
			expression.NewLiteral(2, sql.Int64),
			expression.NewLiteral(1, sql.Int64),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"coordinates": [][][][2]float64{{{{0, 0}, {1, 2}, {3, 4}, {0, 0}}}}, "type": "MultiPolygon", "bbox": [4]float64{0, 0, 3, 4}}}, v)
	})
	t.Run("convert empty geometrycollection to geojson with bounding box", func(t *testing.T) {
		require := require.New(t)
		g := sql.GeomColl{}
		f, err := NewAsGeoJSON(
			expression.NewLiteral(g, sql.GeomCollType{}),
			expression.NewLiteral(2, sql.Int64),
			expression.NewLiteral(1, sql.Int64),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)

		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"geometries": []interface{}{}, "type": "GeometryCollection"}}, v)
	})
	t.Run("convert geometrycollection to geojson with bounding box", func(t *testing.T) {
		require := require.New(t)
		point := sql.Point{X: 1, Y: 2}
		line := sql.LineString{Points: []sql.Point{{X: 1, Y: 2}, {X: 3, Y: 4}}}
		poly := sql.Polygon{Lines: []sql.LineString{{Points: []sql.Point{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: 0}}}}}
		mpoint := sql.MultiPoint{Points: []sql.Point{point, point}}
		mline := sql.MultiLineString{Lines: []sql.LineString{line, line}}
		mpoly := sql.MultiPolygon{Polygons: []sql.Polygon{poly, poly}}
		gColl := sql.GeomColl{}
		g := sql.GeomColl{Geoms: []sql.GeometryValue{point, line, poly, mpoint, mline, mpoly, gColl}}
		f, err := NewAsGeoJSON(
			expression.NewLiteral(g, sql.GeomCollType{}),
			expression.NewLiteral(2, sql.Int64),
			expression.NewLiteral(1, sql.Int64))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)

		pointjson := map[string]interface{}{"coordinates": [2]float64{1, 2}, "type": "Point"}
		linejson := map[string]interface{}{"coordinates": [][2]float64{{1, 2}, {3, 4}}, "type": "LineString"}
		polyjson := map[string]interface{}{"coordinates": [][][2]float64{{{0, 0}, {1, 1}, {1, 0}, {0, 0}}}, "type": "Polygon"}
		mpointjson := map[string]interface{}{"coordinates": [][2]float64{{1, 2}, {1, 2}}, "type": "MultiPoint"}
		mlinejson := map[string]interface{}{"coordinates": [][][2]float64{{{1, 2}, {3, 4}}, {{1, 2}, {3, 4}}}, "type": "MultiLineString"}
		mpolyjson := map[string]interface{}{"coordinates": [][][][2]float64{{{{0, 0}, {1, 1}, {1, 0}, {0, 0}}}, {{{0, 0}, {1, 1}, {1, 0}, {0, 0}}}}, "type": "MultiPolygon"}
		mgeomjson := map[string]interface{}{"geometries": []interface{}{}, "type": "GeometryCollection"}

		require.Equal(sql.JSONDocument{Val: map[string]interface{}{"bbox": [4]float64{0, 0, 3, 4}, "geometries": []interface{}{pointjson, linejson, polyjson, mpointjson, mlinejson, mpolyjson, mgeomjson}, "type": "GeometryCollection"}}, v)
	})
	t.Run("convert point with srid 0 and flag 2", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(
			expression.NewLiteral(sql.Point{X: 1, Y: 2}, sql.PointType{}),
			expression.NewLiteral(1, sql.Int64),
			expression.NewLiteral(2, sql.Int64),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		obj := map[string]interface{}{
			"coordinates": [2]float64{1, 2},
			"type":        "Point",
		}
		require.Equal(sql.JSONDocument{Val: obj}, v)
	})
	t.Run("convert point with srid 4326 and flag 2", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(
			expression.NewLiteral(sql.Point{SRID: 4326, X: 1, Y: 2}, sql.PointType{}),
			expression.NewLiteral(1, sql.Int64),
			expression.NewLiteral(2, sql.Int64),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		obj := map[string]interface{}{
			"crs": map[string]interface{}{
				"type": "name",
				"properties": map[string]interface{}{
					"name": "EPSG:4326",
				},
			},
			"coordinates": [2]float64{1, 2},
			"type":        "Point",
		}
		require.Equal(sql.JSONDocument{Val: obj}, v)
	})
	t.Run("convert point with srid 4326 and flag 4", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(
			expression.NewLiteral(sql.Point{SRID: 4326, X: 1, Y: 2}, sql.PointType{}),
			expression.NewLiteral(1, sql.Int64),
			expression.NewLiteral(4, sql.Int64),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		obj := map[string]interface{}{
			"crs": map[string]interface{}{
				"type": "name",
				"properties": map[string]interface{}{
					"name": "urn:ogc:def:crs:EPSG::4326",
				},
			},
			"coordinates": [2]float64{1, 2},
			"type":        "Point",
		}
		require.Equal(sql.JSONDocument{Val: obj}, v)
	})
	t.Run("convert null is null", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(
			expression.NewLiteral(nil, sql.Null),
			expression.NewLiteral(2, sql.Int64),
			expression.NewLiteral(1, sql.Int64),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(nil, v)
	})
	t.Run("convert null precision is null", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(
			expression.NewLiteral(sql.Point{X: 1, Y: 2}, sql.PointType{}),
			expression.NewLiteral(nil, sql.Null),
			expression.NewLiteral(1, sql.Int64),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(nil, v)
	})
	t.Run("convert null flag is null", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(
			expression.NewLiteral(sql.Point{X: 1, Y: 2}, sql.PointType{}),
			expression.NewLiteral(2, sql.Int64),
			expression.NewLiteral(nil, sql.Null),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(nil, v)
	})
}

func TestGeomFromGeoJSON(t *testing.T) {
	t.Run("convert point from geojson", func(t *testing.T) {
		require := require.New(t)
		f, err := NewGeomFromGeoJSON(expression.NewLiteral(`{"type":"Point", "coordinates":[1,2]}`, sql.Blob))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.Point{SRID: 4326, X: 1, Y: 2}, v)
	})
	t.Run("convert linestring from geojson", func(t *testing.T) {
		require := require.New(t)
		f, err := NewGeomFromGeoJSON(expression.NewLiteral(`{"type":"LineString", "coordinates":[[1,2],[3,4]]}`, sql.Blob))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.LineString{SRID: 4326, Points: []sql.Point{{4326, 1, 2}, {4326, 3, 4}}}, v)
	})
	t.Run("convert polygon from geojson", func(t *testing.T) {
		require := require.New(t)
		f, err := NewGeomFromGeoJSON(expression.NewLiteral(`{"type":"Polygon", "coordinates":[[[0,0],[1,1],[0,1],[0,0]]]}`, sql.Blob))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.Polygon{SRID: 4326, Lines: []sql.LineString{{4326, []sql.Point{{4326, 0, 0}, {4326, 1, 1}, {4326, 0, 1}, {4326, 0, 0}}}}}, v)
	})
	t.Run("convert multipoint from geojson", func(t *testing.T) {
		require := require.New(t)
		f, err := NewGeomFromGeoJSON(expression.NewLiteral(`{"type":"MultiPoint", "coordinates":[[1,2],[3,4]]}`, sql.Blob))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.MultiPoint{SRID: 4326, Points: []sql.Point{{4326, 1, 2}, {4326, 3, 4}}}, v)
	})
	t.Run("convert multilinestring from geojson", func(t *testing.T) {
		require := require.New(t)
		f, err := NewGeomFromGeoJSON(expression.NewLiteral(`{"type":"MultiLineString", "coordinates":[[[0,0],[1,1],[0,1],[0,0]]]}`, sql.Blob))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.MultiLineString{SRID: 4326, Lines: []sql.LineString{{4326, []sql.Point{{4326, 0, 0}, {4326, 1, 1}, {4326, 0, 1}, {4326, 0, 0}}}}}, v)
	})
	t.Run("convert mutlipolygon from geojson", func(t *testing.T) {
		require := require.New(t)
		f, err := NewGeomFromGeoJSON(expression.NewLiteral(`{"type":"MultiPolygon", "coordinates":[[[[0,0],[1,1],[0,1],[0,0]]]]}`, sql.Blob))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.MultiPolygon{SRID: 4326, Polygons: []sql.Polygon{{SRID: 4326, Lines: []sql.LineString{{SRID: 4326, Points: []sql.Point{{SRID: 4326, X: 0, Y: 0}, {SRID: 4326, X: 1, Y: 1}, {SRID: 4326, X: 0, Y: 1}, {SRID: 4326, X: 0, Y: 0}}}}}}}, v)
	})
	t.Run("convert empty geometrycollection from geojson", func(t *testing.T) {
		require := require.New(t)
		f, err := NewGeomFromGeoJSON(expression.NewLiteral(`{"type":"GeometryCollection", "geometries":[]}`, sql.Blob))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		require.Equal(sql.GeomColl{SRID: 4326, Geoms: []sql.GeometryValue{}}, v)
	})
	t.Run("convert geometrycollection to geojson", func(t *testing.T) {
		require := require.New(t)

		s := `{"type":"GeometryCollection", "geometries":[` +
			`{"type":"Point", "coordinates":[1,2]},` +
			`{"type":"LineString", "coordinates":[[1,2],[3,4]]},` +
			`{"type":"Polygon", "coordinates":[[[0,0],[1,1],[1,0],[0,0]]]},` +
			`{"type":"MultiPoint", "coordinates":[[1,2],[1,2]]},` +
			`{"type":"MultiLineString", "coordinates":[[[1,2],[3,4]],[[1,2],[3,4]]]},` +
			`{"type":"MultiPolygon", "coordinates":[[[[0,0],[1,1],[1,0],[0,0]]],[[[0,0],[1,1],[1,0],[0,0]]]]},` +
			`{"type":"GeometryCollection", "geometries":[]}` +
			`]}`
		f, err := NewGeomFromGeoJSON(expression.NewLiteral(s, sql.Blob))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)

		point := sql.Point{SRID: 4326, X: 1, Y: 2}
		line := sql.LineString{SRID: 4326, Points: []sql.Point{{SRID: 4326, X: 1, Y: 2}, {SRID: 4326, X: 3, Y: 4}}}
		poly := sql.Polygon{SRID: 4326, Lines: []sql.LineString{{SRID: 4326, Points: []sql.Point{{SRID: 4326, X: 0, Y: 0}, {SRID: 4326, X: 1, Y: 1}, {SRID: 4326, X: 1, Y: 0}, {SRID: 4326, X: 0, Y: 0}}}}}
		mpoint := sql.MultiPoint{SRID: 4326, Points: []sql.Point{point, point}}
		mline := sql.MultiLineString{SRID: 4326, Lines: []sql.LineString{line, line}}
		mpoly := sql.MultiPolygon{SRID: 4326, Polygons: []sql.Polygon{poly, poly}}
		gColl := sql.GeomColl{SRID: 4326, Geoms: []sql.GeometryValue{}}
		g := sql.GeomColl{SRID: 4326, Geoms: []sql.GeometryValue{point, line, poly, mpoint, mline, mpoly, gColl}}
		require.Equal(g, v)
	})
	t.Run("reject dimensions greater than 2 with flag 1", func(t *testing.T) {
		require := require.New(t)
		f, err := NewGeomFromGeoJSON(
			expression.NewLiteral(`{"type":"Polygon", "coordinates":[[[0,0],[1,1],[0,1],[0,0,0]]]}`, sql.Blob),
			expression.NewLiteral(1, sql.Int32),
		)
		require.NoError(err)

		_, err = f.Eval(sql.NewEmptyContext(), nil)
		require.Error(err)
	})
	t.Run("accept dimensions greater than 2 with flag 2", func(t *testing.T) {
		require := require.New(t)
		f, err := NewGeomFromGeoJSON(
			expression.NewLiteral(`{"type":"Polygon", "coordinates":[[[0,0],[1,1],[0,1],[0,0,0]]]}`, sql.Blob),
			expression.NewLiteral(2, sql.Int32),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.Equal(sql.Polygon{SRID: 4326, Lines: []sql.LineString{{4326, []sql.Point{{4326, 0, 0}, {4326, 1, 1}, {4326, 0, 1}, {4326, 0, 0}}}}}, v)
	})
	t.Run("srid 0 swaps x and y", func(t *testing.T) {
		require := require.New(t)
		f, err := NewGeomFromGeoJSON(
			expression.NewLiteral(`{"type":"Point", "coordinates":[1,2]}`, sql.Blob),
			expression.NewLiteral(1, sql.Int32),
			expression.NewLiteral(0, sql.Int32),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.Equal(sql.Point{0, 1, 2}, v)
	})
	t.Run("srid 0 swaps x and y", func(t *testing.T) {
		require := require.New(t)
		f, err := NewGeomFromGeoJSON(
			expression.NewLiteral(`{"type":"LineString", "coordinates":[[1,2],[3,4]]}`, sql.Blob),
			expression.NewLiteral(1, sql.Int32),
			expression.NewLiteral(0, sql.Int32),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.Equal(sql.LineString{SRID: 0, Points: []sql.Point{{0, 1, 2}, {0, 3, 4}}}, v)
	})
	t.Run("srid 0 swaps x and y", func(t *testing.T) {
		require := require.New(t)
		f, err := NewGeomFromGeoJSON(
			expression.NewLiteral(`{"type":"Polygon", "coordinates":[[[0,0],[1,1],[0,1],[0,0]]]}`, sql.Blob),
			expression.NewLiteral(1, sql.Int32),
			expression.NewLiteral(0, sql.Int32),
		)
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.Equal(sql.Polygon{SRID: 0, Lines: []sql.LineString{{0, []sql.Point{{0, 0, 0}, {0, 1, 1}, {0, 0, 1}, {0, 0, 0}}}}}, v)
	})
	t.Run("check return type", func(t *testing.T) {
		require := require.New(t)
		f, err := NewAsGeoJSON(expression.NewLiteral(sql.Point{X: 1, Y: 2}, sql.PointType{}))
		require.NoError(err)

		v, err := f.Eval(sql.NewEmptyContext(), nil)
		require.NoError(err)
		typ := f.Type()

		_, err = typ.Convert(v)
		require.NoError(err)
	})
}
