package tools

import (
	"context"

	"github.com/agentplexus/mcpkit/runtime"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/grokify/chathub/internal/storage"
)

// RegisterAll registers all ChatHub tools with the MCP runtime.
func RegisterAll(rt *runtime.Runtime, store *storage.Storage) {
	// save_conversation
	runtime.AddTool[SaveConversationInput, SaveConversationOutput](rt, &mcp.Tool{
		Name:        "save_conversation",
		Description: "Save an AI conversation to storage with Hugo-compatible frontmatter",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input SaveConversationInput) (*mcp.CallToolResult, SaveConversationOutput, error) {
		output, err := SaveConversation(ctx, store, input)
		return nil, output, err
	})

	// read_conversation
	runtime.AddTool[ReadConversationInput, ReadConversationOutput](rt, &mcp.Tool{
		Name:        "read_conversation",
		Description: "Read a conversation from storage, returning content and metadata",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ReadConversationInput) (*mcp.CallToolResult, ReadConversationOutput, error) {
		output, err := ReadConversation(ctx, store, input)
		return nil, output, err
	})

	// list_conversations
	runtime.AddTool[ListConversationsInput, ListConversationsOutput](rt, &mcp.Tool{
		Name:        "list_conversations",
		Description: "List conversations with optional filtering by source platform",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListConversationsInput) (*mcp.CallToolResult, ListConversationsOutput, error) {
		output, err := ListConversations(ctx, store, input)
		return nil, output, err
	})

	// search_conversations
	runtime.AddTool[SearchConversationsInput, SearchConversationsOutput](rt, &mcp.Tool{
		Name:        "search_conversations",
		Description: "Search conversations by content or metadata",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input SearchConversationsInput) (*mcp.CallToolResult, SearchConversationsOutput, error) {
		output, err := SearchConversations(ctx, store, input)
		return nil, output, err
	})

	// delete_conversation
	runtime.AddTool[DeleteConversationInput, DeleteConversationOutput](rt, &mcp.Tool{
		Name:        "delete_conversation",
		Description: "Delete a conversation from storage",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input DeleteConversationInput) (*mcp.CallToolResult, DeleteConversationOutput, error) {
		output, err := DeleteConversation(ctx, store, input)
		return nil, output, err
	})

	// append_conversation
	runtime.AddTool[AppendConversationInput, AppendConversationOutput](rt, &mcp.Tool{
		Name:        "append_conversation",
		Description: "Append content to an existing conversation",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input AppendConversationInput) (*mcp.CallToolResult, AppendConversationOutput, error) {
		output, err := AppendConversation(ctx, store, input)
		return nil, output, err
	})
}
