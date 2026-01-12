package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/grokify/chathub/internal/frontmatter"
	"github.com/grokify/chathub/internal/storage"
)

// SaveConversation saves an AI conversation to storage.
func SaveConversation(ctx context.Context, store *storage.Storage, input SaveConversationInput) (SaveConversationOutput, error) {
	// Validate source
	if !frontmatter.ValidSource(input.Source) {
		return SaveConversationOutput{}, fmt.Errorf("invalid source: %s", input.Source)
	}

	// Create frontmatter
	fm := frontmatter.New(input.Title, input.Source)
	fm.Tags = input.Tags
	fm.Categories = input.Categories

	// Set description (use provided or extract from content)
	if input.Description != "" {
		fm.Description = input.Description
	} else {
		fm.Description = frontmatter.ExtractDescription([]byte(input.Content), 150)
	}

	// Generate file path
	filePath := frontmatter.GeneratePath(store.Folder(), input.Source, input.Title, time.Now().UTC())

	// Render complete document
	content, err := fm.RenderWithContent([]byte(input.Content))
	if err != nil {
		return SaveConversationOutput{}, fmt.Errorf("failed to render frontmatter: %w", err)
	}

	// Save to storage
	if err := store.Save(ctx, filePath, content); err != nil {
		return SaveConversationOutput{}, fmt.Errorf("failed to save conversation: %w", err)
	}

	return SaveConversationOutput{
		Path:           filePath,
		ConversationID: fm.ConversationID,
	}, nil
}
