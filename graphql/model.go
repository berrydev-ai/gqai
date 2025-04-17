package graphql

// graphqlRequest represents a GraphQL graphqlRequest.
type graphqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// Response represents a GraphQL response.
type response struct {
	Data   map[string]interface{} `json:"data,omitempty"`
	Errors []graphqlError         `json:"errors,omitempty"`
}

// Error represents a GraphQL error.
type graphqlError struct {
	Message string `json:"message"`
}
