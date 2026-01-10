// Copyright (C) moneta. 2025-present.
//
// Created at 2025-08-28, by liasica

package godoc

import (
	"fmt"
	"maps"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Schema struct {
	root map[string]any

	parent any
	key    any
}

func ConvertEnum2OneOf(input string) (err error) {
	var data []byte
	data, err = os.ReadFile(input)
	if err != nil {
		return
	}

	var root map[string]any
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		return
	}

	schema := Schema{
		root:   make(map[string]any),
		parent: nil,
		key:    nil,
	}

	maps.Copy(schema.root, root)
	schema.walk(root)

	var out []byte
	out, err = yaml.Marshal(root)
	if err != nil {
		return
	}

	return os.WriteFile(input, out, 0644)
}

func convertEnum(node map[string]any) (ok bool) {
	enumVals, ok1 := node["enum"].([]any)
	varNames, ok2 := node["x-enum-varnames"].([]any)
	enumComments, ok3 := node["x-enum-comments"].(map[string]any)

	ok = ok1 && ok2 && ok3 && len(enumVals) == len(varNames)
	if ok {
		var builder strings.Builder

		builder.WriteString(node["description"].(string))
		builder.WriteString(", 可选值:\n")
		var vals []any

		for i, val := range enumVals {
			name := varNames[i].(string)
			if val != "-" {
				vals = append(vals, val)
				builder.WriteString("\n")
				builder.WriteString(fmt.Sprintf("- `%s`: %v", val, enumComments[name]))
			}
		}
		node["description"] = strings.TrimSuffix(builder.String(), "\n")
		node["enum"] = vals
	}

	return
}

func (s *Schema) walk(node any) {
	switch n := node.(type) {
	case map[string]any:
		ok := convertEnum(n)
		if ok {
			// fmt.Println(s.key, s.parent, node)
			if s.parent != nil {
				switch m := s.parent.(type) {
				case map[string]any:
					if _, exists := m["description"]; exists {
						delete(m, "description")
					}
				}
			}
		}
		for k, v := range n {
			s.key = k
			s.parent = n
			s.walk(v)
		}
	case []any:
		for k, v := range n {
			s.key = k
			s.parent = n
			s.walk(v)
		}
	}
}
