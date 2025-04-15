package tool

import (
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

func ExtractInputSchema(rawQuery string) (map[string]any, error) {
	doc, err := parser.ParseQuery(&ast.Source{Input: rawQuery})
	if err != nil {
		return nil, err
	}

	if len(doc.Operations) == 0 {
		return nil, nil
	}

	op := doc.Operations[0]

	props := map[string]any{}
	required := []string{}

	for _, v := range op.VariableDefinitions {
		schema := map[string]any{
			"type": graphqlTypeToJSONSchemaType(v.Type),
		}
		props[v.Variable] = schema
		if v.Type.NonNull {
			required = append(required, v.Variable)
		}
	}

	result := map[string]any{
		"type":       "object",
		"properties": props,
	}
	if len(required) > 0 {
		result["required"] = required
	}
	return result, nil
}

func graphqlTypeToJSONSchemaType(t *ast.Type) string {
	if t.Elem != nil {
		return "array"
	}
	switch t.NamedType {
	case "String", "ID":
		return "string"
	case "Int":
		return "integer"
	case "Float":
		return "number"
	case "Boolean":
		return "boolean"
	default:
		return "string" // fallback
	}
}
