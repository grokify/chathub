# ChatHub v0.1.0 Release Notes

**Release Date:** 2026-01-11

ChatHub v0.1.0 is the initial release of the MCP server for cross-platform AI conversation sharing.

## Highlights

- **Cross-Platform Sharing**: Save conversations from ChatGPT Desktop and read them in Claude Code, or vice versa
- **Hugo-Compatible**: Conversations use standard Hugo frontmatter for easy static site publishing
- **GitHub Storage**: Version-controlled conversation storage via GitHub repositories
- **Full CRUD**: Complete conversation lifecycle with save, read, list, search, delete, and append operations

## Features

### MCP Tools

| Tool | Description |
|------|-------------|
| `save_conversation` | Save a new conversation with Hugo-compatible frontmatter |
| `read_conversation` | Read a conversation by path |
| `list_conversations` | List conversations with optional source filtering |
| `search_conversations` | Search conversations by content |
| `delete_conversation` | Delete a conversation |
| `append_conversation` | Append content to an existing conversation |

### Storage

- GitHub backend (primary) with automatic commits
- Organized by source: `conversations/{source}/{date}_{slug}.md`
- Automatic slug generation from titles

### Hugo Integration

Conversations use Hugo-compatible YAML frontmatter:

```yaml
---
title: "Building an MCP Server"
date: 2026-01-10T14:30:00Z
tags: ["mcp", "golang"]
source: "chatgpt"
conversation_id: "conv_abc123"
---
```

## Installation

```bash
go install github.com/grokify/chathub/cmd/chathub@v0.1.0
```

## Configuration

Set environment variables:

```bash
export GITHUB_TOKEN="ghp_your_token"
export GITHUB_OWNER="your-username"
export GITHUB_REPO="ai-conversations"
```

Configure your MCP client (e.g., Claude Code):

```json
{
  "mcpServers": {
    "chathub": {
      "command": "/path/to/chathub",
      "env": {
        "GITHUB_TOKEN": "ghp_your_token",
        "GITHUB_OWNER": "your-username",
        "GITHUB_REPO": "ai-conversations"
      }
    }
  }
}
```

## Dependencies

- [mcpkit](https://github.com/agentplexus/mcpkit) v0.3.1 - Library-first MCP runtime
- [omnistorage](https://github.com/grokify/omnistorage) v0.1.0 - Multi-backend storage abstraction
- [omnistorage-github](https://github.com/grokify/omnistorage-github) v0.1.0 - GitHub storage backend
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) v1.2.0

## Known Limitations

- GitHub backend only (S3, Dropbox, local filesystem planned for future releases)
- Simple text-based search (full-text indexing planned)
- No tag-based filtering in list operations (planned for v0.2.0)

## What's Next

See [ROADMAP.md](ROADMAP.md) for planned features including:

- Additional storage backends (S3, Dropbox, local filesystem)
- Tag-based filtering
- Conversation threading
- Web UI for browsing

## Contributors

- John Wang (@grokify)

## License

MIT License - see [LICENSE](LICENSE) for details.
