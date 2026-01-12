package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/grokify/chathub/internal/frontmatter"
	"github.com/grokify/chathub/internal/storage"
)

// AppendConversationInput is the input for the append_conversation tool.
type AppendConversationInput struct {
	Path    string `json:"path" jsonschema:"Path to the existing conversation"`
	Content string `json:"content" jsonschema:"Content to append (Markdown)"`
}

// AppendConversationOutput is the output for the append_conversation tool.
type AppendConversationOutput struct {
	Path         string `json:"path" jsonschema:"Updated file path"`
	MessageCount int    `json:"message_count,omitempty" jsonschema:"Updated message count"`
}

// AppendConversation appends content to an existing conversation.
func AppendConversation(ctx context.Context, store *storage.Storage, input AppendConversationInput) (AppendConversationOutput, error) {
	// Read existing content
	existing, err := store.Read(ctx, input.Path)
	if err != nil {
		return AppendConversationOutput{}, fmt.Errorf("failed to read conversation: %w", err)
	}

	// Parse frontmatter
	fm, body, err := frontmatter.Parse(existing)
	if err != nil {
		return AppendConversationOutput{}, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	if fm == nil {
		// No frontmatter - just append
		newContent := append(existing, []byte("\n\n"+input.Content)...)
		if err := store.Save(ctx, input.Path, newContent); err != nil {
			return AppendConversationOutput{}, fmt.Errorf("failed to save: %w", err)
		}
		return AppendConversationOutput{Path: input.Path}, nil
	}

	// Update frontmatter
	fm.LastMod = time.Now().UTC()
	fm.MessageCount++ // Increment message count (approximate)

	// Append new content to body
	newBody := append(body, []byte("\n\n"+input.Content)...)

	// Render updated document
	updated, err := fm.RenderWithContent(newBody)
	if err != nil {
		return AppendConversationOutput{}, fmt.Errorf("failed to render: %w", err)
	}

	// Save
	if err := store.Save(ctx, input.Path, updated); err != nil {
		return AppendConversationOutput{}, fmt.Errorf("failed to save: %w", err)
	}

	return AppendConversationOutput{
		Path:         input.Path,
		MessageCount: fm.MessageCount,
	}, nil
}
