package gws

import (
	"encoding/json"

	"github.com/graphql-go/graphql/gqlerrors"
)

// Constants for packet types
const (
	gqlConnectionInit      = "connection_init"
	gqlConnectionAck       = "connection_ack"
	gqlConnectionKeepAlive = "keepalive"
	gqlConnectionError     = "connection_error"
	gqlConnectionTerminate = "connection_terminate"
	gqlSubscriptionSuccess = "subscription_success"
	gqlStart               = "start"
	gqlData                = "data"
	gqlError               = "error"
	gqlComplete            = "complete"
	gqlStop                = "stop"
	gqlUnknown             = ""
	gqlInvalid             = "<invalid>"
)

// InitMessagePayload defines the parameters of a connection
// init message.
type InitMessagePayload struct {
	AuthToken string `json:"authToken"`
}

// StartMessagePayload defines the parameters of an operation that
// a client requests to be started.
type StartMessagePayload struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName"`
}

// DataMessagePayload defines the result data of an operation.
type DataMessagePayload struct {
	Data   interface{} `json:"data"`
	Errors []error     `json:"errors"`
}

// OperationMessage represents a GraphQL WebSocket message.
type OperationMessage struct {
	ID      string      `json:"id,omitempty"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (msg *OperationMessage) String() string {
	s, _ := json.Marshal(msg)
	if s != nil {
		return string(s)
	}
	return "<invalid>"
}

// Marshal helper func
func (msg *OperationMessage) Marshal() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		return []byte(gqlInvalid)
	}
	return b
}

// Unmarshal helper func
func (msg *OperationMessage) Unmarshal(b []byte) error {
	return json.Unmarshal(b, &msg)
}

// ErrorsFromGraphQLErrors convert from GraphQL errors to regular errors.
func ErrorsFromGraphQLErrors(errors []gqlerrors.FormattedError) []error {
	if len(errors) == 0 {
		return nil
	}
	out := make([]error, len(errors))
	for i := range errors {
		out[i] = errors[i]
	}
	return out
}
