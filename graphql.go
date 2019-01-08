package hyper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/graphql-go/graphql"

	"github.com/vaniila/hyper/router"
)

// GQLQueries parameters
var GQLQueries = []router.Param{
	Query("query").
		Format(Text).
		Require(false),
	Query("variables").
		Format(Text).
		Require(false),
}

// GQLBodies parameters
var GQLBodies = []router.Param{
	Body("query").
		Format(Text).
		Require(false),
	Body("variables").
		Format(Text).
		Require(false),
	Body("file").
		Format(File).
		Require(false),
}

// Payload struct
type Payload struct {
	RawQuery        string                 `json:"-"`
	RawVariables    string                 `json:"-"`
	ParsedQuery     string                 `json:"query"`
	ParsedVariables map[string]interface{} `json:"variables"`
}

// Parse to read raw query and variables
func (v *Payload) Parse() {
	v.ParsedQuery = v.RawQuery
	v.ParsedVariables = make(map[string]interface{})
	json.Unmarshal([]byte(v.RawVariables), &v.ParsedVariables)
}

// GraphQL handles graphql
func GraphQL(schema graphql.Schema) router.HandlerFunc {
	return func(c router.Context) {
		var payload = new(Payload)
		switch c.Req().Method {
		case "PUT", "POST", "PATCH", "CONNECT":
			switch {
			case c.MustBody("query").Has() || c.MustBody("variables").Has():
				payload.RawQuery = c.MustBody("query").String()
				payload.RawVariables = c.MustBody("variables").String()
				payload.Parse()
			default:
				b, _ := ioutil.ReadAll(c.Req().Body)
				if err := json.Unmarshal(b, &payload); err != nil {
					payload.ParsedQuery = string(b[:])
				}
			}
		default:
			switch {
			case c.MustQuery("query").Has() || c.MustQuery("variables").Has():
				payload.RawQuery = c.MustQuery("query").String()
				payload.RawVariables = c.MustQuery("variables").String()
				payload.Parse()
			default:
				r := c.Req().URL.RawQuery
				b := []byte(r)
				if err := json.Unmarshal(b, &payload); err != nil {
					payload.ParsedQuery = r
				}
			}
		}

		span := c.StartSpan("HTTP GraphQL Execution")
		span.LogKV("graphql-query", payload.ParsedQuery)
		for k, v := range payload.ParsedVariables {
			span.LogKV(fmt.Sprintf("graphql-variable:%s", k), v)
		}
		defer span.Finish()
		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  payload.ParsedQuery,
			VariableValues: payload.ParsedVariables,
			Context:        c.Context(),
		})
		if result.HasErrors() {
			c.Status(http.StatusForbidden)
		}
		span.LogKV("graphql-result", result)
		c.Json(result)
	}
}
