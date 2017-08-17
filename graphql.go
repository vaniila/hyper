package hyper

import (
	"encoding/json"
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
}

// Payload struct
type Payload struct {
	RawQuery        string                 `json:"query"`
	RawVariables    string                 `json:"variables"`
	ParsedQuery     string                 `json:"-"`
	ParsedVariables map[string]interface{} `json:"-"`
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
			default:
				b, _ := ioutil.ReadAll(c.Req().Body)
				json.Unmarshal(b, &payload)
			}
		default:
			switch {
			case c.MustQuery("query").Has() || c.MustQuery("variables").Has():
				payload.RawQuery = c.MustQuery("query").String()
				payload.RawVariables = c.MustQuery("variables").String()
			default:
				r := c.Req().URL.RawQuery
				b := []byte(r)
				if err := json.Unmarshal(b, &payload); err != nil {
					payload.RawQuery = r
				}
			}
		}
		payload.Parse()
		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  payload.ParsedQuery,
			VariableValues: payload.ParsedVariables,
			Context:        c.Context(),
		})
		if result.HasErrors() {
			c.Status(http.StatusForbidden)
		}
		c.Json(result)
	}
}
