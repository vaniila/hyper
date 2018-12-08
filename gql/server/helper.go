package server

import (
	"context"

	"github.com/vaniila/hyper/gql/graphql"
	"github.com/vaniila/hyper/router"
)

// FromContext reads router context from context.Context
func FromContext(c context.Context) graphql.Context {
	return &ctx{c.Value(router.RequestContext).(router.Context)}
}
