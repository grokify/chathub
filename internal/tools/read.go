package tools

import (
	"context"
	"fmt"

	"github.com/grokify/chathub/internal/frontmatter"
	"github.com/grokify/chathub/internal/storage"
)

// ReadConversation reads a conversation from storage.
func ReadConversation(ctx context.Context, store *storage.Storage, input ReadConversationInput) (ReadConversationOutput, error) {
	// Read from storage
	content, err := store.Read(ctx, input.Path)
	if err != nil {
		return ReadConversationOutput{}, fmt.Errorf("failed to read conversation: %w", err)
	}

	// Parse frontmatter
	fm, body, err := frontmatter.Parse(content)
	if err != nil {
		return ReadConversationOutput{}, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	output := ReadConversationOutput{
		Content: string(content), // Return full content including frontmatter
	}

	// Extract metadata from frontmatter
	if fm != nil {
		output.Title = fm.Title
		output.Date = fm.Date.Format("2006-01-02T15:04:05Z")
		output.Source = fm.Source
		output.Tags = fm.Tags
		output.Description = fm.Description

		output.Metadata = map[string]string{
			"conversation_id": fm.ConversationID,
			"author":          fm.Author,
			"slug":            fm.Slug,
		}
		if fm.Model != "" {
			output.Metadata["model"] = fm.Model
		}
	} else {
		// No frontmatter - just return content
		output.Content = string(body)
	}

	return output, nil
}
