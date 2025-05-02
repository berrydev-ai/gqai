package graphql

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

type Operation struct {
	Name          string
	Doc           *ast.QueryDocument
	Raw           string
	OperationType string
}

func LoadOperations(config *GraphQLConfig) (map[string]*Operation, error) {
	opMap := make(map[string]*Operation)

	err := filepath.WalkDir(config.SingleProject.Documents[0], func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".graphql" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		doc, parseErr := parser.ParseQuery(&ast.Source{
			Name:  filepath.Base(path),
			Input: string(data),
		})
		if parseErr != nil {
			return parseErr
		}

		for _, op := range doc.Operations {
			opMap[op.Name] = &Operation{
				Name:          op.Name,
				Doc:           doc,
				Raw:           string(data),
				OperationType: string(op.Operation),
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load operations: %v", err)
	}

	return opMap, nil
}
