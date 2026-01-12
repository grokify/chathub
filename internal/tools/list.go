package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/grokify/chathub/internal/frontmatter"
	"github.com/grokify/chathub/internal/storage"
)

const defaultListLimit = 50

// ListConversations lists conversations with optional filtering.
func ListConversations(ctx context.Context, store *storage.Storage, input ListConversationsInput) (ListConversationsOutput, error) {
	// Set default limit
	limit := input.Limit
	if limit <= 0 {
		limit = defaultListLimit
	}

	// Get file list
	var files []string
	var err error

	if input.Source != "" {
		files, err = store.ListBySource(ctx, input.Source)
	} else {
		files, err = store.ListConversations(ctx)
	}
	if err != nil {
		return ListConversationsOutput{}, fmt.Errorf("failed to list conversations: %w", err)
	}

	// Filter to only .md files
	var mdFiles []string
	for _, f := range files {
		if strings.HasSuffix(f, ".md") {
			mdFiles = append(mdFiles, f)
		}
	}

	total := len(mdFiles)

	// Apply pagination
	start := input.Offset
	if start > len(mdFiles) {
		start = len(mdFiles)
	}
	end := start + limit
	if end > len(mdFiles) {
		end = len(mdFiles)
	}

	pagedFiles := mdFiles[start:end]
	hasMore := end < len(mdFiles)

	// Build summaries
	conversations := make([]ConversationSummary, 0, len(pagedFiles))
	for _, filePath := range pagedFiles {
		summary, err := getConversationSummary(ctx, store, filePath)
		if err != nil {
			// Log error but continue with other files
			continue
		}
		conversations = append(conversations, summary)
	}

	return ListConversationsOutput{
		Conversations: conversations,
		Total:         total,
		HasMore:       hasMore,
	}, nil
}

func getConversationSummary(ctx context.Context, store *storage.Storage, filePath string) (ConversationSummary, error) {
	content, err := store.Read(ctx, filePath)
	if err != nil {
		return ConversationSummary{}, err
	}

	fm, _, err := frontmatter.Parse(content)
	if err != nil {
		// Return basic info from path if frontmatter fails
		return ConversationSummary{
			Path: filePath,
		}, nil
	}

	if fm == nil {
		return ConversationSummary{
			Path: filePath,
		}, nil
	}

	return ConversationSummary{
		Path:        filePath,
		Title:       fm.Title,
		Date:        fm.Date.Format("2006-01-02"),
		Source:      fm.Source,
		Tags:        fm.Tags,
		Description: fm.Description,
	}, nil
}
