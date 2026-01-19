# ChatHub - Technical Requirements Document

## Architecture Overview

```
┌──────────────────────────────────────────────────────────────────┐
│                         MCP Clients                              │
│  ChatGPT Desktop │ Claude Code │ Claude.ai │ Gemini │ Perplexity │
└──────────────────────────────────────────────────────────────────┘
                              │
                              │ MCP Protocol (JSON-RPC)
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│                      ChatHub MCP Server                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐            │
│  │    Tools     │  │  Frontmatter │  │    Config    │            │
│  │  save/read   │  │  Hugo-compat │  │  env vars    │            │
│  │  list/search │  │  YAML parse  │  │  validation  │            │
│  └──────────────┘  └──────────────┘  └──────────────┘            │
│                              │                                   │
│                     ┌────────┴────────┐                          │
│                     │   mcpruntime    │                          │
│                     │  (MCP SDK wrap) │                          │
│                     └────────┬────────┘                          │
└──────────────────────────────┼───────────────────────────────────┘
                               │
                    ┌──────────┴──────────┐
                    │    omnistorage      │
                    │  Backend Interface  │
                    └──────────┬──────────┘
                               │
        ┌──────────┬───────────┼───────────┬──────────┐
        ▼          ▼           ▼           ▼          ▼
    ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
    │ GitHub │ │   S3   │ │Dropbox │ │  File  │ │ Memory │
    └────────┘ └────────┘ └────────┘ └────────┘ └────────┘
```

## Technology Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Language | Go 1.23+ | Performance, single binary, MCP SDK support |
| MCP Runtime | mcpruntime | Library-first design, dual mode (library/server) |
| Storage | omnistorage | Backend abstraction, multiple storage providers |
| YAML | gopkg.in/yaml.v3 | Frontmatter parsing/generation |
| CLI | Built-in flags | Simple configuration, no external deps |

## Project Structure

```
chathub/
├── cmd/
│   └── chathub/
│       └── main.go          # CLI entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration loading
│   ├── frontmatter/
│   │   ├── frontmatter.go   # Parse/generate YAML frontmatter
│   │   └── frontmatter_test.go
│   ├── storage/
│   │   ├── storage.go       # OmniStorage backend setup
│   │   └── storage_test.go
│   └── tools/
│       ├── save.go          # save_conversation tool
│       ├── read.go          # read_conversation tool
│       ├── list.go          # list_conversations tool
│       ├── search.go        # search_conversations tool
│       ├── delete.go        # delete_conversation tool
│       └── tools_test.go
├── PRD.md
├── TRD.md
├── ROADMAP.md
├── README.md
├── CHANGELOG.json
├── go.mod
└── go.sum
```

## Dependencies

### Direct Dependencies

```go
require (
    github.com/agentplexus/mcpkit v0.3.1
    github.com/grokify/omnistorage v0.1.0
    github.com/grokify/omnistorage-github v0.1.0
    github.com/modelcontextprotocol/go-sdk v1.2.0
    gopkg.in/yaml.v3 v3.0.1
)
```

### Development Dependencies

```go
require (
    github.com/stretchr/testify v1.9.0 // testing assertions
)
```

## Configuration

### Environment Variables

```bash
# Core settings
CHATHUB_BACKEND=github          # github, s3, dropbox, file, memory
CHATHUB_FOLDER=conversations    # Root folder for conversations

# GitHub backend
GITHUB_TOKEN=ghp_xxxx           # Personal access token
GITHUB_OWNER=grokify            # Repository owner
GITHUB_REPO=ai-conversations    # Repository name
GITHUB_BRANCH=main              # Branch (default: main)

# S3 backend
AWS_ACCESS_KEY_ID=xxxx
AWS_SECRET_ACCESS_KEY=xxxx
S3_BUCKET=my-bucket
S3_REGION=us-west-2
S3_ENDPOINT=                    # For R2/MinIO (optional)

# Dropbox backend
DROPBOX_TOKEN=xxxx

# File backend
FILE_ROOT=/path/to/storage

# Server settings
CHATHUB_PORT=8080               # HTTP port (for HTTP transport)
CHATHUB_TRANSPORT=stdio         # stdio, http
```

### Config Struct

```go
type Config struct {
    Backend       string            // Backend type
    BackendConfig map[string]string // Backend-specific config
    Folder        string            // Root folder
    Transport     string            // stdio or http
    Port          int               // HTTP port
}

func LoadConfig() (*Config, error) {
    cfg := &Config{
        Backend:   getEnv("CHATHUB_BACKEND", "github"),
        Folder:    getEnv("CHATHUB_FOLDER", "conversations"),
        Transport: getEnv("CHATHUB_TRANSPORT", "stdio"),
        Port:      getEnvInt("CHATHUB_PORT", 8080),
    }

    cfg.BackendConfig = loadBackendConfig(cfg.Backend)
    return cfg, cfg.Validate()
}
```

## MCP Tool Definitions

### save_conversation

```go
type SaveConversationInput struct {
    Title       string   `json:"title" jsonschema:"Conversation title"`
    Content     string   `json:"content" jsonschema:"Full conversation in Markdown"`
    Source      string   `json:"source" jsonschema:"Source platform (chatgpt/claude/gemini/perplexity/codex/claude-code)"`
    Tags        []string `json:"tags,omitempty" jsonschema:"Tags for categorization"`
    Categories  []string `json:"categories,omitempty" jsonschema:"Categories for organization"`
    Description string   `json:"description,omitempty" jsonschema:"Brief summary"`
}

type SaveConversationOutput struct {
    Path string `json:"path" jsonschema:"Saved file path"`
    URL  string `json:"url,omitempty" jsonschema:"Web URL if available"`
}
```

> **Note:** Fields without `omitempty` in the `json` tag are automatically marked as required by the JSON schema generator. The `jsonschema` tag value is used as the field description.

### read_conversation

```go
type ReadConversationInput struct {
    Path string `json:"path" jsonschema:"Full path to conversation"`
}

type ReadConversationOutput struct {
    Content     string            `json:"content" jsonschema:"Full Markdown content"`
    Frontmatter map[string]any    `json:"frontmatter" jsonschema:"Parsed frontmatter"`
}
```

### list_conversations

```go
type ListConversationsInput struct {
    Source string   `json:"source,omitempty" jsonschema:"Filter by source platform"`
    Tags   []string `json:"tags,omitempty" jsonschema:"Filter by tags"`
    Limit  int      `json:"limit,omitempty" jsonschema:"Max results (default 50)"`
    Offset int      `json:"offset,omitempty" jsonschema:"Pagination offset"`
}

type ConversationSummary struct {
    Path        string   `json:"path"`
    Title       string   `json:"title"`
    Date        string   `json:"date"`
    Source      string   `json:"source"`
    Tags        []string `json:"tags,omitempty"`
    Description string   `json:"description,omitempty"`
}

type ListConversationsOutput struct {
    Conversations []ConversationSummary `json:"conversations"`
    Total         int                   `json:"total"`
    HasMore       bool                  `json:"has_more"`
}
```

### search_conversations

```go
type SearchConversationsInput struct {
    Query  string `json:"query" jsonschema:"Search query"`
    Source string `json:"source,omitempty" jsonschema:"Filter by source"`
    Limit  int    `json:"limit,omitempty" jsonschema:"Max results (default 20)"`
}

type SearchResult struct {
    Path    string  `json:"path"`
    Title   string  `json:"title"`
    Snippet string  `json:"snippet" jsonschema:"Matching text excerpt"`
    Score   float64 `json:"score" jsonschema:"Relevance score"`
}

type SearchConversationsOutput struct {
    Results []SearchResult `json:"results"`
    Total   int            `json:"total"`
}
```

### delete_conversation

```go
type DeleteConversationInput struct {
    Path string `json:"path" jsonschema:"Full path to conversation"`
}

type DeleteConversationOutput struct {
    Deleted bool   `json:"deleted"`
    Message string `json:"message,omitempty"`
}
```

## Frontmatter Implementation

### Parse Frontmatter

```go
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

    // ChatHub extensions
    Source         string   `yaml:"source"`
    ConversationID string   `yaml:"conversation_id"`
    Participants   []string `yaml:"participants,omitempty"`
    MessageCount   int      `yaml:"message_count,omitempty"`
    Model          string   `yaml:"model,omitempty"`
    Tokens         int      `yaml:"tokens,omitempty"`
}

func ParseFrontmatter(content []byte) (*Frontmatter, []byte, error)
func (f *Frontmatter) Render() ([]byte, error)
```

### Generate Slug

```go
func GenerateSlug(title string) string {
    // Lowercase, replace spaces with hyphens, remove special chars
    slug := strings.ToLower(title)
    slug = regexp.MustCompile(`[^a-z0-9\s-]`).ReplaceAllString(slug, "")
    slug = regexp.MustCompile(`[\s_]+`).ReplaceAllString(slug, "-")
    slug = strings.Trim(slug, "-")
    if len(slug) > 50 {
        slug = slug[:50]
    }
    return slug
}
```

### Generate Path

```go
func GeneratePath(folder, source, title string, date time.Time) string {
    slug := GenerateSlug(title)
    dateStr := date.Format("2006-01-02")
    return fmt.Sprintf("%s/%s/%s_%s.md", folder, source, dateStr, slug)
}
```

## Storage Integration

### Backend Initialization

```go
import (
    "github.com/grokify/omnistorage"
    _ "github.com/grokify/omnistorage-github/backend/github"
    _ "github.com/grokify/omnistorage/backend/s3"
    _ "github.com/grokify/omnistorage/backend/dropbox"
    _ "github.com/grokify/omnistorage/backend/file"
    _ "github.com/grokify/omnistorage/backend/memory"
)

func NewBackend(cfg *Config) (omnistorage.Backend, error) {
    return omnistorage.Open(cfg.Backend, cfg.BackendConfig)
}
```

### Storage Operations

```go
type Storage struct {
    backend omnistorage.Backend
    folder  string
}

func (s *Storage) Save(ctx context.Context, path string, content []byte) error {
    w, err := s.backend.NewWriter(ctx, path)
    if err != nil {
        return err
    }
    defer w.Close()
    _, err = w.Write(content)
    return err
}

func (s *Storage) Read(ctx context.Context, path string) ([]byte, error) {
    r, err := s.backend.NewReader(ctx, path)
    if err != nil {
        return nil, err
    }
    defer r.Close()
    return io.ReadAll(r)
}

func (s *Storage) List(ctx context.Context, prefix string) ([]string, error) {
    return s.backend.List(ctx, prefix)
}

func (s *Storage) Delete(ctx context.Context, path string) error {
    return s.backend.Delete(ctx, path)
}
```

## MCP Server Setup

### Using mcpruntime

```go
func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal(err)
    }

    backend, err := storage.NewBackend(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer backend.Close()

    store := storage.New(backend, cfg.Folder)

    // Create MCP runtime
    rt := mcpruntime.New(&mcp.Implementation{
        Name:    "chathub",
        Version: "0.1.0",
    }, nil)

    // Register tools
    tools.RegisterAll(rt, store)

    // Serve
    ctx := context.Background()
    switch cfg.Transport {
    case "stdio":
        if err := rt.ServeStdio(ctx); err != nil {
            log.Fatal(err)
        }
    case "http":
        http.Handle("/mcp", rt.StreamableHTTPHandler(nil))
        log.Printf("Listening on :%d", cfg.Port)
        log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
    }
}
```

### Tool Registration

```go
func RegisterAll(rt *mcpruntime.Runtime, store *storage.Storage) {
    mcpruntime.AddTool(rt, &mcp.Tool{
        Name:        "save_conversation",
        Description: "Save an AI conversation to storage",
    }, func(ctx context.Context, input SaveConversationInput) (SaveConversationOutput, error) {
        return SaveConversation(ctx, store, input)
    })

    mcpruntime.AddTool(rt, &mcp.Tool{
        Name:        "read_conversation",
        Description: "Read a conversation from storage",
    }, func(ctx context.Context, input ReadConversationInput) (ReadConversationOutput, error) {
        return ReadConversation(ctx, store, input)
    })

    // ... register other tools
}
```

## Testing Strategy

### Unit Tests

- Frontmatter parsing/generation
- Slug generation
- Path generation
- Tool handlers with mock storage

### Integration Tests

- Memory backend: Full tool flow
- File backend: Filesystem operations
- GitHub backend: API integration (with test repo)

### Test with mcpruntime Library Mode

```go
func TestSaveAndReadConversation(t *testing.T) {
    // Use memory backend for testing
    backend := memory.New()
    store := storage.New(backend, "conversations")

    rt := mcpruntime.New(&mcp.Implementation{
        Name: "chathub-test",
    }, nil)

    tools.RegisterAll(rt, store)

    // Test save
    saveResult, err := rt.CallTool(ctx, "save_conversation", map[string]any{
        "title":   "Test Conversation",
        "content": "# Test\n\nHello world",
        "source":  "chatgpt",
    })
    require.NoError(t, err)

    // Test read
    readResult, err := rt.CallTool(ctx, "read_conversation", map[string]any{
        "path": saveResult["path"],
    })
    require.NoError(t, err)
    assert.Contains(t, readResult["content"], "Hello world")
}
```

## Error Handling

### Error Types

```go
var (
    ErrConversationNotFound = errors.New("conversation not found")
    ErrInvalidSource        = errors.New("invalid source platform")
    ErrInvalidFrontmatter   = errors.New("invalid frontmatter format")
    ErrStorageError         = errors.New("storage operation failed")
)
```

### Error Responses

MCP tools return errors in the standard format:

```json
{
  "isError": true,
  "content": [
    {
      "type": "text",
      "text": "conversation not found: conversations/chatgpt/2026-01-10_test.md"
    }
  ]
}
```

## Security Considerations

1. **Token Storage**: All tokens via environment variables, never in code
2. **Input Validation**: Sanitize paths to prevent directory traversal
3. **Content Filtering**: No logging of conversation content
4. **Rate Limiting**: Respect GitHub API rate limits (5000 req/hour)

## Performance Targets

| Operation | Target | Notes |
|-----------|--------|-------|
| save_conversation | < 2s | Includes GitHub commit |
| read_conversation | < 200ms | Single file read |
| list_conversations | < 500ms | Up to 1000 files |
| search_conversations | < 1s | Full-text search |
| delete_conversation | < 1s | Includes GitHub commit |

## Deployment

### Binary Distribution

```bash
# Build
go build -o chathub ./cmd/chathub

# Run with stdio transport (for MCP clients)
./chathub

# Run with HTTP transport
CHATHUB_TRANSPORT=http CHATHUB_PORT=8080 ./chathub
```

### MCP Client Configuration

**Claude Code (~/.config/claude-code/config.json):**

```json
{
  "mcpServers": {
    "chathub": {
      "command": "/path/to/chathub",
      "env": {
        "GITHUB_TOKEN": "ghp_xxxx",
        "GITHUB_OWNER": "username",
        "GITHUB_REPO": "ai-conversations"
      }
    }
  }
}
```

**ChatGPT Desktop:**

Configure via MCP server settings with same environment variables.
