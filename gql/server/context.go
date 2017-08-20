package server

import (
	"context"
	"net/http"
	"time"

	"github.com/vaniila/hyper/router"
)

// ERRORS GraphQL error reference
const ERRORS = "graphql-errors"

// graphql context
type ctx struct {
	private router.Context
}

func (v *ctx) Deadline() (time.Time, bool)          { return v.private.Deadline() }
func (v *ctx) Done() <-chan struct{}                { return v.private.Done() }
func (v *ctx) Err() error                           { return v.private.Err() }
func (v *ctx) Value(key interface{}) interface{}    { return v.private.Value(key) }
func (v *ctx) Identity() router.Identity            { return v.private.Identity() }
func (v *ctx) MachineID() string                    { return v.private.MachineID() }
func (v *ctx) ProcessID() string                    { return v.private.ProcessID() }
func (v *ctx) Context() context.Context             { return v.private.Context() }
func (v *ctx) Req() *http.Request                   { return v.private.Req() }
func (v *ctx) Res() http.ResponseWriter             { return v.private.Res() }
func (v *ctx) Client() router.Client                { return v.private.Client() }
func (v *ctx) Cache() router.CacheAdaptor           { return v.private.Cache() }
func (v *ctx) Message() router.MessageAdaptor       { return v.private.Message() }
func (v *ctx) KV() router.KV                        { return v.private.KV() }
func (v *ctx) Cookie() router.Cookie                { return v.private.Cookie() }
func (v *ctx) Header() router.Header                { return v.private.Header() }
func (v *ctx) MustParam(s string) router.Value      { return v.private.MustParam(s) }
func (v *ctx) MustQuery(s string) router.Value      { return v.private.MustQuery(s) }
func (v *ctx) MustBody(s string) router.Value       { return v.private.MustBody(s) }
func (v *ctx) Param(s string) (router.Value, error) { return v.private.Param(s) }
func (v *ctx) Query(s string) (router.Value, error) { return v.private.Query(s) }
func (v *ctx) Body(s string) (router.Value, error)  { return v.private.Body(s) }
func (v *ctx) File(s string) []byte                 { return v.private.File(s) }
func (v *ctx) Abort()                               { v.private.Abort() }
func (v *ctx) IsAborted() bool                      { return v.private.IsAborted() }
func (v *ctx) Native() router.Context               { return v.private }

func (v *ctx) HasErrors() bool {
	o := v.private.KV().Get(ERRORS)
	return o != nil && len(o) > 0
}

func (v *ctx) Errors() []error {
	return nil
}

func (v *ctx) GraphQLError(e error) {
}
