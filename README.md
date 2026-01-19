# ChatHub

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

**ChatHub** is an MCP (Model Context Protocol) server for sharing AI assistant conversations across platforms. Save conversations from ChatGPT, Claude, Gemini, Perplexity, or any MCP-compatible client and read them from any other.

## Features

- **Cross-platform sharing**: Save from ChatGPT, read from Claude Code
- **Multiple storage backends**: GitHub, S3, Dropbox, local filesystem
- **Hugo-compatible**: Conversations use standard Hugo frontmatter for static site publishing
- **MCP standard**: Works with any MCP-compatible AI assistant

## Installation

```bash
go install github.com/grokify/chathub/cmd/chathub@latest
```

Or build from source:

```bash
git clone https://github.com/grokify/chathub.git
cd chathub
go build -o chathub ./cmd/chathub
```

## Quick Start

### 1. Configure environment

```bash
# Required for GitHub backend
export GITHUB_TOKEN="ghp_your_token"
export GITHUB_OWNER="your-username"
export GITHUB_REPO="ai-conversations"

# Optional
export CHATHUB_FOLDER="conversations"  # default
export CHATHUB_BACKEND="github"        # default
```

### 2. Configure your MCP client

**Claude Code** (`~/.claude/claude_desktop_config.json`):

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

**ChatGPT Desktop**: Configure via MCP server settings with the same environment variables.

### 3. Use it

In ChatGPT or Claude:

```
Save this conversation about MCP servers
```

The AI will call `save_conversation` and store it to your GitHub repo.

In Claude Code:

```
What conversations do I have about MCP?
```

Claude will call `list_conversations` and show your saved conversations.

## MCP Tools

| Tool | Description |
|------|-------------|
| `save_conversation` | Save a new conversation with Hugo-compatible frontmatter |
| `append_conversation` | Append content to an existing conversation |
| `read_conversation` | Read a conversation by path |
| `list_conversations` | List conversations with optional source filtering |
| `search_conversations` | Search conversations by content |
| `delete_conversation` | Delete a conversation |

## Example Prompts

### Save (Create)

```
Save this conversation as "MCP Server Tutorial"
```

```
Save our discussion about debugging with tags "golang" and "debugging"
```

```
Archive this chat to my conversation hub
```

### Append (Add to Existing)

```
Add this follow-up to my MCP server conversation
```

```
Append our latest discussion to the OAuth2 conversation
```

```
Continue the debugging tips conversation with today's session
```

### List

```
What conversations do I have saved?
```

```
Show me my recent ChatGPT conversations
```

```
List all conversations tagged with "tutorial"
```

### Read

```
Read my conversation about MCP servers from January 10th
```

```
Show me the debugging tips conversation
```

```
Open conversations/chatgpt/2026-01-10_mcp-server-tutorial.md
```

### Search

```
Search my conversations for "authentication"
```

```
Find discussions where I talked about error handling
```

```
What did I learn about MCP protocols?
```

### Delete

```
Delete the conversation at conversations/chatgpt/2026-01-10_test.md
```

```
Remove my draft conversation about testing
```

### Cross-Platform Workflow

**In ChatGPT:**

```
Let's discuss how to implement OAuth2 in Go. When we're done, save this conversation.
```

**Later, in Claude Code:**

```
Read my ChatGPT conversation about OAuth2 implementation
```

```
Continue from where I left off with OAuth2 - what were the key points?
```

## Storage Backends

| Backend | Config | Use Case |
|---------|--------|----------|
| `github` | `GITHUB_TOKEN`, `GITHUB_OWNER`, `GITHUB_REPO` | Version-controlled storage |
| `s3` | `S3_BUCKET`, `S3_REGION`, `S3_ENDPOINT` | AWS S3, R2, MinIO |
| `dropbox` | `DROPBOX_TOKEN` | Personal cloud storage |
| `file` | `FILE_ROOT` | Local filesystem |
| `memory` | (none) | Testing |

Set `CHATHUB_BACKEND` to select a backend (default: `github`).

## Hugo Integration

ChatHub conversations use Hugo-compatible YAML frontmatter:

```yaml
---
title: "Building an MCP Server"
date: 2026-01-10T14:30:00Z
tags: ["mcp", "golang"]
categories: ["development"]
author: "chatgpt"
source: "chatgpt"
conversation_id: "conv_1234567890"
---

# Building an MCP Server

**User:** How do I build an MCP server in Go?

**ChatGPT:** To build an MCP server...
```

To publish as a Hugo site:

1. Configure ChatHub to save to your Hugo `content/conversations/` directory
2. Add a layout for conversations
3. Run `hugo serve`

## File Organization

Conversations are organized by source:

```
conversations/
├── chatgpt/
│   ├── 2026-01-10_mcp-server-tutorial.md
│   └── 2026-01-11_debugging-tips.md
├── claude/
│   └── 2026-01-10_code-review.md
└── gemini/
    └── 2026-01-12_research-notes.md
```

## HTTP Transport

For HTTP/SSE transport instead of stdio:

```bash
export CHATHUB_TRANSPORT=http
export CHATHUB_PORT=8080
./chathub
```

Access the MCP endpoint at `http://localhost:8080/mcp`.

## Dependencies

- [mcpkit](https://github.com/agentplexus/mcpkit) - Library-first MCP runtime
- [omnistorage](https://github.com/grokify/omnistorage) - Multi-backend storage abstraction
- [omnistorage-github](https://github.com/grokify/omnistorage-github) - GitHub storage backend
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) - Official MCP SDK

## Guides

- [ChatGPT Desktop ↔ Claude Code Integration](README_CHATGPT_CLAUDECODE.md) - Step-by-step setup for cross-platform sharing

## License

MIT License - see [LICENSE](LICENSE) for details.

 [build-status-svg]: https://github.com/grokify/chathub/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/grokify/chathub/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/grokify/chathub/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/grokify/chathub/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/chathub
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/chathub
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/chathub
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/chathub
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fchathub
 [loc-svg]: https://tokei.rs/b1/github/grokify/chathub
 [repo-url]: https://github.com/grokify/chathub
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/chathub/blob/master/LICENSE