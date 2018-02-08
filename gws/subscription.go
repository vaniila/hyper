package gws

import "github.com/graphql-go/graphql/language/ast"

type subscription struct {
	id, query, opname string
	fields            []string
	variables, args   map[string]interface{}
	doc               *ast.Document
	ctx               Context
}

func (v *subscription) ID() string {
	return v.id
}

func (v *subscription) Query() string {
	return v.query
}

func (v *subscription) Variables() map[string]interface{} {
	return v.variables
}

func (v *subscription) Arguments() map[string]interface{} {
	return v.args
}

func (v *subscription) OperationName() string {
	return v.opname
}

func (v *subscription) Document() *ast.Document {
	return v.doc
}

func (v *subscription) Fields() []string {
	return v.fields
}

func (v *subscription) Connection() Context {
	return v.ctx
}
