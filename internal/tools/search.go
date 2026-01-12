package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/grokify/chathub/internal/frontmatter"
	"github.com/grokify/chathub/internal/storage"
)

const (
	defaultSearchLimit  = 20
	snippetContextChars = 100
)

// SearchConversations searches conversations by content or metadata.
func SearchConversations(ctx context.Context, store *storage.Storage, input SearchConversationsInput) (SearchConversationsOutput, error) {
	// Set default limit
	limit := input.Limit
	if limit <= 0 {
		limit = defaultSearchLimit
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
		return SearchConversationsOutput{}, fmt.Errorf("failed to list conversations: %w", err)
	}

	// Search through files
	query := strings.ToLower(input.Query)
	var results []SearchResult

	for _, filePath := range files {
		if !strings.HasSuffix(filePath, ".md") {
			continue
		}

		result, found := searchFile(ctx, store, filePath, query)
		if found {
			results = append(results, result)
			if len(results) >= limit {
				break
			}
		}
	}

	return SearchConversationsOutput{
		Results: results,
		Total:   len(results),
	}, nil
}

func searchFile(ctx context.Context, store *storage.Storage, filePath, query string) (SearchResult, bool) {
	content, err := store.Read(ctx, filePath)
	if err != nil {
		return SearchResult{}, false
	}

	contentLower := strings.ToLower(string(content))
	idx := strings.Index(contentLower, query)
	if idx == -1 {
		return SearchResult{}, false
	}

	// Parse frontmatter for title
	fm, body, _ := frontmatter.Parse(content)

	title := ""
	if fm != nil {
		title = fm.Title
	}

	// Extract snippet around match
	snippet := extractSnippet(string(body), idx, query, snippetContextChars)

	// Simple scoring based on position and frequency
	score := calculateScore(contentLower, query, idx)

	return SearchResult{
		Path:    filePath,
		Title:   title,
		Snippet: snippet,
		Score:   score,
	}, true
}

func extractSnippet(content string, idx int, query string, contextChars int) string {
	// Adjust index for frontmatter offset
	start := idx - contextChars
	if start < 0 {
		start = 0
	}

	end := idx + len(query) + contextChars
	if end > len(content) {
		end = len(content)
	}

	snippet := content[start:end]

	// Clean up snippet
	snippet = strings.ReplaceAll(snippet, "\n", " ")
	snippet = strings.Join(strings.Fields(snippet), " ")

	// Add ellipsis
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(content) {
		snippet = snippet + "..."
	}

	return snippet
}

func calculateScore(content, query string, firstIdx int) float64 {
	// Count occurrences
	count := strings.Count(content, query)

	// Score based on:
	// - Earlier position is better (normalized to 0-1)
	// - More occurrences is better
	positionScore := 1.0 - (float64(firstIdx) / float64(len(content)))
	frequencyScore := float64(count) / 10.0 // cap at 10 occurrences
	if frequencyScore > 1.0 {
		frequencyScore = 1.0
	}

	return (positionScore + frequencyScore) / 2.0
}
