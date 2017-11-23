package flatten

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestFlatten(t *testing.T) {
	cases := []struct {
		test   string
		want   map[string]interface{}
		prefix string
		style  SeparatorStyle
	}{
		{
			`{
				"foo": {
					"jim":"bean"
				},
				"fee": "bar",
				"n1": {
					"alist": [
						"a",
						"b",
						"c",
						{
							"d": "other",
							"e": "another"
						}
					]
				},
				"number": 1.4567,
				"bool":   true
			}`,
			map[string]interface{}{
				"foo.jim":      "bean",
				"fee":          "bar",
				"n1.alist.0":   "a",
				"n1.alist.1":   "b",
				"n1.alist.2":   "c",
				"n1.alist.3.d": "other",
				"n1.alist.3.e": "another",
				"number":       1.4567,
				"bool":         true,
			},
			"",
			DotStyle,
		},
		{
			`{
				"foo": {
					"jim":"bean"
				},
				"fee": "bar",
				"n1": {
					"alist": [
					"a",
					"b",
					"c",
					{
						"d": "other",
						"e": "another"
					}
					]
				}
			}`,
			map[string]interface{}{
				"foo[jim]":        "bean",
				"fee":             "bar",
				"n1[alist][0]":    "a",
				"n1[alist][1]":    "b",
				"n1[alist][2]":    "c",
				"n1[alist][3][d]": "other",
				"n1[alist][3][e]": "another",
			},
			"",
			RailsStyle,
		},
		{
			`{ "a": { "b": "c" }, "e": "f" }`,
			map[string]interface{}{
				"p:a.b": "c",
				"p:e":   "f",
			},
			"p:",
			DotStyle,
		},
	}

	for i, test := range cases {
		var m interface{}
		err := json.Unmarshal([]byte(test.test), &m)
		if err != nil {
			t.Errorf("%d: failed to unmarshal test: %v", i+1, err)
		}
		got, err := Flatten(m.(map[string]interface{}), test.prefix, test.style)
		if err != nil {
			t.Errorf("%d: failed to flatten: %v", i+1, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}

func TestFlattenString(t *testing.T) {
	cases := []struct {
		test   string
		want   string
		prefix string
		style  SeparatorStyle
	}{
		{
			`{ "a": "b" }`,
			`{ "a": "b" }`,
			"",
			DotStyle,
		},
		{
			`{ "a": { "b" : { "c" : { "d" : "e" } } }, "number": 1.4567, "bool": true }`,
			`{ "a.b.c.d": "e", "bool": true, "number": 1.4567 }`,
			"",
			DotStyle,
		},
	}

	for i, test := range cases {
		got, err := FlattenString(test.test, test.prefix, test.style)
		if err != nil {
			t.Errorf("%d: failed to flatten: %v", i+1, err)
		}

		nixws := func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}

		if got != strings.Map(nixws, test.want) {
			t.Errorf("%d: mismatch, got: %v want: %v", i+1, got, test.want)
		}
	}
}

func TestFlattenArray(t *testing.T) {
	nested := map[string]interface{}{
		"keyOne":   "b",
		"keyTwo":   []string{"e1", "e2", "e3"},
		"keyThree": 1.4567,
	}

	flat, err := FlattenToInterface(nested, "", UppercaseStyle)
	assert.NoError(t, err)
	assert.NotNil(t, flat)
}

func TestFlattenArrayOfStruct(t *testing.T) {
	nested := map[string]interface{}{
		"keyOne": "b",
		"keyTwo": []map[string]interface{}{
			{
				"k1": "v1",
				"k2": "v2",
			},
			{
				"k3": "v3",
				"k4": "v4",
			},
		},
		"keyThree": 1.4567,
	}

	flat, err := FlattenToInterface(nested, "", UppercaseStyle)
	assert.NoError(t, err)
	assert.NotNil(t, flat)

	fmt.Println(flat)
}
