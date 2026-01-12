package tools

import (
	"context"
	"fmt"

	"github.com/grokify/chathub/internal/storage"
)

// DeleteConversation deletes a conversation from storage.
func DeleteConversation(ctx context.Context, store *storage.Storage, input DeleteConversationInput) (DeleteConversationOutput, error) {
	// Check if file exists
	exists, err := store.Exists(ctx, input.Path)
	if err != nil {
		return DeleteConversationOutput{}, fmt.Errorf("failed to check existence: %w", err)
	}

	if !exists {
		return DeleteConversationOutput{
			Deleted: false,
			Message: fmt.Sprintf("conversation not found: %s", input.Path),
		}, nil
	}

	// Delete the file
	if err := store.Delete(ctx, input.Path); err != nil {
		return DeleteConversationOutput{}, fmt.Errorf("failed to delete conversation: %w", err)
	}

	return DeleteConversationOutput{
		Deleted: true,
		Message: fmt.Sprintf("deleted: %s", input.Path),
	}, nil
}
