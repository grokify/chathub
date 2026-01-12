// Package frontmatter provides Hugo-compatible YAML frontmatter parsing and generation.
package frontmatter

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Source platform identifiers
const (
	SourceChatGPT    = "chatgpt"
	SourceClaude     = "claude"
	SourceClaudeCode = "claude-code"
	SourceGemini     = "gemini"
	SourcePerplexity = "perplexity"
	SourceCodex      = "codex"
)

// Frontmatter represents Hugo-compatible YAML frontmatter with ChatHub extensions.
type Frontmatter struct {
	// Hugo standard fields
	Title       string    `yaml:"title"`
	Date        time.Time `yaml:"date"`
	LastMod     time.Time `yaml:"lastmod,omitempty"`
	Draft       bool      `yaml:"draft,omitempty"`
	Tags        []string  `yaml:"tags,omitempty"`
	Categories  []string  `yaml:"categories,omitempty"`
	Author      string    `yaml:"author,omitempty"`
	Description string    `yaml:"description,omitempty"`
	Slug        string    `yaml:"slug,omitempty"`
	Weight      int       `yaml:"weight,omitempty"`
	Aliases     []string  `yaml:"aliases,omitempty"`

	// ChatHub extension fields
	Source         string   `yaml:"source"`
	ConversationID string   `yaml:"conversation_id,omitempty"`
	Participants   []string `yaml:"participants,omitempty"`
	MessageCount   int      `yaml:"message_count,omitempty"`
	Model          string   `yaml:"model,omitempty"`
	Tokens         int      `yaml:"tokens,omitempty"`
}

var (
	// ErrInvalidFrontmatter indicates malformed frontmatter
	ErrInvalidFrontmatter = errors.New("invalid frontmatter format")

	// frontmatterDelimiter is the YAML frontmatter delimiter
	frontmatterDelimiter = []byte("---")

	// slugRegex matches non-alphanumeric characters for slug generation
	slugRegex = regexp.MustCompile(`[^a-z0-9\s-]`)
	// whitespaceRegex matches whitespace for slug generation
	whitespaceRegex = regexp.MustCompile(`[\s_]+`)
)

// Parse extracts frontmatter from Markdown content.
// Returns the parsed frontmatter, the remaining content, and any error.
func Parse(content []byte) (*Frontmatter, []byte, error) {
	content = bytes.TrimSpace(content)

	// Check for frontmatter delimiter
	if !bytes.HasPrefix(content, frontmatterDelimiter) {
		return nil, content, nil // No frontmatter
	}

	// Find end delimiter
	rest := content[len(frontmatterDelimiter):]
	endIdx := bytes.Index(rest, frontmatterDelimiter)
	if endIdx == -1 {
		return nil, nil, ErrInvalidFrontmatter
	}

	// Extract YAML
	yamlContent := bytes.TrimSpace(rest[:endIdx])
	remaining := bytes.TrimSpace(rest[endIdx+len(frontmatterDelimiter):])

	// Parse YAML
	var fm Frontmatter
	if err := yaml.Unmarshal(yamlContent, &fm); err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrInvalidFrontmatter, err)
	}

	return &fm, remaining, nil
}

// Render generates YAML frontmatter bytes.
func (f *Frontmatter) Render() ([]byte, error) {
	yamlBytes, err := yaml.Marshal(f)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.Write(yamlBytes)
	buf.WriteString("---\n")

	return buf.Bytes(), nil
}

// RenderWithContent generates complete Markdown with frontmatter and content.
func (f *Frontmatter) RenderWithContent(content []byte) ([]byte, error) {
	fmBytes, err := f.Render()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write(fmBytes)
	buf.WriteByte('\n')
	buf.Write(content)

	return buf.Bytes(), nil
}

// GenerateSlug creates a URL-friendly slug from a title.
func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = slugRegex.ReplaceAllString(slug, "")
	slug = whitespaceRegex.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	// Limit length
	if len(slug) > 50 {
		slug = slug[:50]
		// Don't end on a hyphen
		slug = strings.TrimRight(slug, "-")
	}

	return slug
}

// GeneratePath creates a file path for a conversation.
// Format: {folder}/{source}/{date}_{slug}.md
func GeneratePath(folder, source, title string, date time.Time) string {
	slug := GenerateSlug(title)
	dateStr := date.Format("2006-01-02")
	return fmt.Sprintf("%s/%s/%s_%s.md", folder, source, dateStr, slug)
}

// GenerateConversationID creates a unique conversation ID.
func GenerateConversationID() string {
	return fmt.Sprintf("conv_%d", time.Now().UnixNano())
}

// New creates a new Frontmatter with default values.
func New(title, source string) *Frontmatter {
	now := time.Now().UTC()
	return &Frontmatter{
		Title:          title,
		Date:           now,
		LastMod:        now,
		Draft:          false,
		Author:         source,
		Slug:           GenerateSlug(title),
		Source:         source,
		ConversationID: GenerateConversationID(),
	}
}

// ValidSource checks if a source is valid.
func ValidSource(source string) bool {
	switch source {
	case SourceChatGPT, SourceClaude, SourceClaudeCode, SourceGemini, SourcePerplexity, SourceCodex:
		return true
	default:
		return false
	}
}

// ExtractDescription extracts a description from content (first non-empty line or first N chars).
func ExtractDescription(content []byte, maxLen int) string {
	lines := bytes.Split(content, []byte("\n"))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		// Skip markdown headers
		if bytes.HasPrefix(line, []byte("#")) {
			continue
		}
		// Skip bold markers at start (like "**User:**")
		text := string(line)
		if len(text) > maxLen {
			text = text[:maxLen] + "..."
		}
		return text
	}
	return ""
}
