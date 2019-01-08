package hyper

import (
	"context"

	"github.com/vaniila/hyper/router"
)

// Context reads router context from context.Context
func Context(c context.Context) router.Context {
	return c.Value(router.RequestContext).(router.Context)
}

// Parse reads router context from context.Context
func Parse(c context.Context) (router.Context, bool) {
	if o, ok := c.Value(router.RequestContext).(router.Context); ok {
		return o, true
	}
	return nil, false
}
