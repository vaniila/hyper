package gws

import (
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

func operationDefinitionsWithOperation(doc *ast.Document, op string) []*ast.OperationDefinition {
	defs := []*ast.OperationDefinition{}
	for _, node := range doc.Definitions {
		if node.GetKind() == kinds.OperationDefinition {
			if def, ok := node.(*ast.OperationDefinition); ok {
				if def.Operation == op {
					defs = append(defs, def)
				}
			}
		}
	}
	return defs
}

func selectionSetsForOperationDefinitions(defs []*ast.OperationDefinition) []*ast.SelectionSet {
	sets := []*ast.SelectionSet{}
	for _, def := range defs {
		if set := def.GetSelectionSet(); set != nil {
			sets = append(sets, set)
		}
	}
	return sets
}

func getSubscriptionInfo(doc *ast.Document, vars map[string]interface{}) ([]string, map[string]interface{}) {

	var names = make([]string, 0)
	var args = make(map[string]interface{})

	defs := operationDefinitionsWithOperation(doc, "subscription")
	sets := selectionSetsForOperationDefinitions(defs)

	for _, set := range sets {
		if len(set.Selections) >= 1 {
			if field, ok := set.Selections[0].(*ast.Field); ok {
				names = append(names, field.Name.Value)
				for _, arg := range field.Arguments {
					key := arg.Name.Value
					switch o := arg.Value.GetValue().(type) {
					case *ast.Name:
						if d, ok := vars[o.Value]; ok {
							args[key] = d
						} else {
							args[key] = nil
						}
					default:
						args[key] = o
					}
				}
			}
		}
	}

	return names, args
}
