package tool

type MCPTool struct {
	Name        string                                  `json:"name"`
	Description string                                  `json:"description"`
	InputSchema map[string]interface{}                  `json:"inputSchema"`
	Execute     func(input map[string]any) (any, error) `json:"-"`
	Annotations struct {
		Title           string `json:"title,omitempty"`
		ReadOnlyHint    bool   `json:"readOnlyHint"`
		DestructiveHint bool   `json:"destructiveHint"`
		IdempotentHint  bool   `json:"idempotentHint"`
		OpenWorldHint   bool   `json:"openWorldHint"`
	} `json:"annotations,omitempty"`
}
