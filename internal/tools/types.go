// Package tools provides MCP tool implementations for ChatHub.
package tools

// SaveConversationInput is the input for the save_conversation tool.
type SaveConversationInput struct {
	Title       string   `json:"title" jsonschema:"Conversation title"`
	Content     string   `json:"content" jsonschema:"Full conversation in Markdown"`
	Source      string   `json:"source" jsonschema:"Source platform (chatgpt/claude/gemini/perplexity/codex/claude-code)"`
	Tags        []string `json:"tags,omitempty" jsonschema:"Tags for categorization"`
	Categories  []string `json:"categories,omitempty" jsonschema:"Categories for organization"`
	Description string   `json:"description,omitempty" jsonschema:"Brief summary"`
}

// SaveConversationOutput is the output for the save_conversation tool.
type SaveConversationOutput struct {
	Path           string `json:"path" jsonschema:"Saved file path"`
	ConversationID string `json:"conversation_id" jsonschema:"Unique conversation ID"`
}

// ReadConversationInput is the input for the read_conversation tool.
type ReadConversationInput struct {
	Path string `json:"path" jsonschema:"Full path to conversation"`
}

// ReadConversationOutput is the output for the read_conversation tool.
type ReadConversationOutput struct {
	Content     string            `json:"content" jsonschema:"Full Markdown content"`
	Title       string            `json:"title" jsonschema:"Conversation title"`
	Date        string            `json:"date" jsonschema:"Creation date"`
	Source      string            `json:"source" jsonschema:"Source platform"`
	Tags        []string          `json:"tags,omitempty" jsonschema:"Tags"`
	Description string            `json:"description,omitempty" jsonschema:"Brief summary"`
	Metadata    map[string]string `json:"metadata,omitempty" jsonschema:"Additional metadata"`
}

// ListConversationsInput is the input for the list_conversations tool.
type ListConversationsInput struct {
	Source string `json:"source,omitempty" jsonschema:"Filter by source platform"`
	Limit  int    `json:"limit,omitempty" jsonschema:"Max results (default 50)"`
	Offset int    `json:"offset,omitempty" jsonschema:"Pagination offset"`
}

// ConversationSummary represents a conversation in list results.
type ConversationSummary struct {
	Path        string   `json:"path"`
	Title       string   `json:"title"`
	Date        string   `json:"date"`
	Source      string   `json:"source"`
	Tags        []string `json:"tags,omitempty"`
	Description string   `json:"description,omitempty"`
}

// ListConversationsOutput is the output for the list_conversations tool.
type ListConversationsOutput struct {
	Conversations []ConversationSummary `json:"conversations"`
	Total         int                   `json:"total"`
	HasMore       bool                  `json:"has_more"`
}

// SearchConversationsInput is the input for the search_conversations tool.
type SearchConversationsInput struct {
	Query  string `json:"query" jsonschema:"Search query"`
	Source string `json:"source,omitempty" jsonschema:"Filter by source"`
	Limit  int    `json:"limit,omitempty" jsonschema:"Max results (default 20)"`
}

// SearchResult represents a search result.
type SearchResult struct {
	Path    string  `json:"path"`
	Title   string  `json:"title"`
	Snippet string  `json:"snippet" jsonschema:"Matching text excerpt"`
	Score   float64 `json:"score" jsonschema:"Relevance score"`
}

// SearchConversationsOutput is the output for the search_conversations tool.
type SearchConversationsOutput struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
}

// DeleteConversationInput is the input for the delete_conversation tool.
type DeleteConversationInput struct {
	Path string `json:"path" jsonschema:"Full path to conversation"`
}

// DeleteConversationOutput is the output for the delete_conversation tool.
type DeleteConversationOutput struct {
	Deleted bool   `json:"deleted"`
	Message string `json:"message,omitempty"`
}
