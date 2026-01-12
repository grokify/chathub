package frontmatter

import (
	"testing"
	"time"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Building an MCP Server", "building-an-mcp-server"},
		{"Test!@#$%Special", "testspecial"},
		{"Multiple   Spaces", "multiple-spaces"},
		{"", ""},
		{"A Very Long Title That Should Be Truncated To Fifty Characters Maximum", "a-very-long-title-that-should-be-truncated-to-fift"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := GenerateSlug(tt.input)
			if got != tt.expected {
				t.Errorf("GenerateSlug(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGeneratePath(t *testing.T) {
	date := time.Date(2026, 1, 10, 14, 30, 0, 0, time.UTC)
	path := GeneratePath("conversations", "chatgpt", "Hello World", date)
	expected := "conversations/chatgpt/2026-01-10_hello-world.md"
	if path != expected {
		t.Errorf("GeneratePath() = %q, want %q", path, expected)
	}
}

func TestParse(t *testing.T) {
	content := []byte(`---
title: "Test Conversation"
date: 2026-01-10T14:30:00Z
source: chatgpt
tags: [mcp, golang]
---

# Test Content

Hello world!
`)

	fm, body, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if fm == nil {
		t.Fatal("Parse() returned nil frontmatter")
	}

	if fm.Title != "Test Conversation" {
		t.Errorf("Title = %q, want %q", fm.Title, "Test Conversation")
	}

	if fm.Source != "chatgpt" {
		t.Errorf("Source = %q, want %q", fm.Source, "chatgpt")
	}

	if len(fm.Tags) != 2 {
		t.Errorf("len(Tags) = %d, want 2", len(fm.Tags))
	}

	if string(body) == "" {
		t.Error("body is empty")
	}
}

func TestParseNoFrontmatter(t *testing.T) {
	content := []byte("# Just Markdown\n\nNo frontmatter here.")

	fm, body, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if fm != nil {
		t.Error("expected nil frontmatter")
	}

	if len(body) == 0 {
		t.Error("expected non-empty body")
	}
}

func TestRender(t *testing.T) {
	fm := &Frontmatter{
		Title:  "Test",
		Date:   time.Date(2026, 1, 10, 14, 30, 0, 0, time.UTC),
		Source: "chatgpt",
		Tags:   []string{"mcp", "golang"},
	}

	output, err := fm.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	// Parse it back
	parsed, _, err := Parse(output)
	if err != nil {
		t.Fatalf("Parse(Render()) error = %v", err)
	}

	if parsed.Title != fm.Title {
		t.Errorf("Title = %q, want %q", parsed.Title, fm.Title)
	}

	if parsed.Source != fm.Source {
		t.Errorf("Source = %q, want %q", parsed.Source, fm.Source)
	}
}

func TestRenderWithContent(t *testing.T) {
	fm := New("Test Conversation", "chatgpt")
	content := []byte("# Hello\n\nThis is the conversation content.")

	output, err := fm.RenderWithContent(content)
	if err != nil {
		t.Fatalf("RenderWithContent() error = %v", err)
	}

	// Parse it back
	parsed, body, err := Parse(output)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if parsed.Title != fm.Title {
		t.Errorf("Title = %q, want %q", parsed.Title, fm.Title)
	}

	if string(body) != string(content) {
		t.Errorf("body = %q, want %q", string(body), string(content))
	}
}

func TestValidSource(t *testing.T) {
	validSources := []string{"chatgpt", "claude", "claude-code", "gemini", "perplexity", "codex"}
	for _, s := range validSources {
		if !ValidSource(s) {
			t.Errorf("ValidSource(%q) = false, want true", s)
		}
	}

	invalidSources := []string{"", "unknown", "gpt4", "ChatGPT"}
	for _, s := range invalidSources {
		if ValidSource(s) {
			t.Errorf("ValidSource(%q) = true, want false", s)
		}
	}
}

func TestExtractDescription(t *testing.T) {
	content := []byte(`# Title

**User:** How do I build an MCP server?

**Assistant:** To build an MCP server...`)

	desc := ExtractDescription(content, 50)
	if desc == "" {
		t.Error("ExtractDescription() returned empty string")
	}

	// Should skip the title line
	if desc == "# Title" {
		t.Error("ExtractDescription() should skip markdown headers")
	}
}
