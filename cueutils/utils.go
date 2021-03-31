package cueutils

import (
	"cuelang.org/go/cue/ast"
)

func GetAcceptedValues(node ast.Node) ([]string, error) {
	switch v := node.(type) {
	case *ast.Ident:
		return []string{v.Name}, nil

	case *ast.ListLit:
		return []string{"list"}, nil
	}

	return []string{"None"}, nil
}
