# ChatHub - Product Requirements Document

## Overview

**ChatHub** is an MCP (Model Context Protocol) server that enables sharing AI assistant conversations across platforms. Users can save conversations from any MCP-compatible AI assistant and read them from any other.

## Problem Statement

AI conversations are siloed within their respective platforms:

- ChatGPT conversations stay in ChatGPT
- Claude conversations stay in Claude
- Gemini, Perplexity, Codex conversations are similarly isolated

Users cannot:

- Continue a conversation started in one AI with another AI
- Search across all their AI conversations
- Build a personal knowledge base from AI interactions
- Share conversation context between desktop and CLI tools

## Solution

ChatHub provides a unified storage layer accessible via MCP:

1. **Save** conversations from any MCP-compatible client
2. **Read** conversations from any MCP-compatible client
3. **Search** across all saved conversations
4. **Organize** by source, tags, and categories

## Target Users

- Power users working with multiple AI assistants
- Developers using both ChatGPT desktop app and Claude Code
- Researchers building knowledge bases from AI interactions
- Teams sharing AI-generated content

## Supported Platforms

### Write (Save Conversations)

| Platform | Interface | Status |
|----------|-----------|--------|
| ChatGPT Desktop | MCP Server | Primary |
| Claude.ai | MCP Server | Primary |
| Claude Code | MCP Server / CLI | Primary |
| Gemini (web) | MCP Server | Planned |
| Gemini CLI | MCP Server | Planned |
| Perplexity | MCP Server | Planned |
| Codex CLI | MCP Server | Planned |

### Read (Retrieve Conversations)

Any MCP-compatible client can read stored conversations.

## Storage Backends

| Backend | Status | Use Case |
|---------|--------|----------|
| GitHub | Primary | Version-controlled, public/private repos |
| S3-compatible | Supported | AWS S3, Cloudflare R2, MinIO, Wasabi |
| Dropbox | Supported | Personal cloud storage |
| Local filesystem | Supported | Development, offline use |
| Google Drive | Future | Personal cloud storage |
| OneDrive | Future | Enterprise integration |

## Core Features

### F1: Save Conversation

Save an AI conversation to storage with metadata.

**Inputs:**

- `title` (required): Conversation title
- `content` (required): Full conversation in Markdown
- `source` (required): Source platform identifier
- `tags` (optional): Array of tags for categorization
- `categories` (optional): Array of categories

**Behavior:**

- Generates Hugo-compatible YAML frontmatter
- Creates file at `{root}/{source}/{date}_{slug}.md`
- Returns saved file path

### F2: Read Conversation

Read a conversation from storage.

**Inputs:**

- `path` (required): Full path to conversation file

**Outputs:**

- Full Markdown content including frontmatter

### F3: List Conversations

List available conversations with filtering.

**Inputs:**

- `source` (optional): Filter by source platform
- `tags` (optional): Filter by tags
- `limit` (optional): Max results (default: 50)
- `offset` (optional): Pagination offset

**Outputs:**

- Array of conversation metadata (path, title, date, source, tags)

### F4: Search Conversations

Full-text search across conversations.

**Inputs:**

- `query` (required): Search query
- `source` (optional): Filter by source

**Outputs:**

- Matching conversations with relevant snippets

### F5: Delete Conversation

Remove a conversation from storage.

**Inputs:**

- `path` (required): Full path to conversation file

## Hugo Compatibility Extension

ChatHub conversations use Hugo-compatible YAML frontmatter, enabling direct publishing as a static site.

### Standard Hugo Fields

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Conversation title |
| `date` | datetime | Creation timestamp (RFC 3339) |
| `lastmod` | datetime | Last modification timestamp |
| `draft` | boolean | Draft status (default: false) |
| `tags` | []string | Tags for categorization |
| `categories` | []string | Categories for organization |
| `author` | string | Source AI platform |
| `description` | string | Brief summary/first message |
| `slug` | string | URL-friendly identifier |
| `weight` | int | Sort order (optional) |
| `aliases` | []string | Redirect URLs (optional) |

### ChatHub Extension Fields

| Field | Type | Description |
|-------|------|-------------|
| `source` | string | AI platform: `chatgpt`, `claude`, `gemini`, `perplexity`, `codex`, `claude-code` |
| `conversation_id` | string | Unique identifier (e.g., `conv_abc123`) |
| `participants` | []string | Conversation participants |
| `message_count` | int | Number of messages |
| `model` | string | AI model used (e.g., `gpt-4`, `claude-3-opus`) |
| `tokens` | int | Estimated token count (optional) |

### Example Frontmatter

```yaml
---
title: "Building an MCP Server in Go"
date: 2026-01-10T14:30:00Z
lastmod: 2026-01-10T15:45:00Z
draft: false
tags: ["mcp", "golang", "tutorial"]
categories: ["development"]
author: "chatgpt"
description: "Discussion about building MCP servers using the Go SDK"
slug: "mcp-server-golang"
# ChatHub extensions
source: "chatgpt"
conversation_id: "conv_abc123"
participants: ["user", "gpt-4"]
message_count: 24
model: "gpt-4"
---

# Building an MCP Server in Go

**User:** How do I build an MCP server in Go?

**ChatGPT:** To build an MCP server in Go, you can use the official MCP Go SDK...
```

### Hugo Site Integration

1. Configure ChatHub to save to Hugo `content/conversations/` directory
2. Conversations appear as posts with full metadata
3. Build searchable knowledge base with `hugo serve`

```
my-hugo-site/
├── content/
│   └── conversations/
│       ├── chatgpt/
│       │   └── 2026-01-10_mcp-server-golang.md
│       └── claude/
│           └── 2026-01-10_debugging-tips.md
├── layouts/
│   └── conversations/
│       └── single.html
└── config.toml
```

## Non-Functional Requirements

### Performance

- List operations: < 500ms for 1000 conversations
- Read operations: < 200ms for typical conversation
- Save operations: < 2s (including GitHub commit)

### Security

- GitHub token stored securely (environment variable)
- No conversation content logged
- Support for private repositories

### Reliability

- Graceful handling of network failures
- Retry logic for transient errors
- Clear error messages for user feedback

## Success Metrics

- Conversations saved per week
- Cross-platform reads (save in ChatGPT, read in Claude Code)
- User retention after 30 days

## Out of Scope (v1)

- Real-time sync between backends
- Conversation merging/threading
- Export to PDF/JSON formats
- Web UI for browsing
- Multi-user collaboration

## References

- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- [Hugo Front Matter](https://gohugo.io/content-management/front-matter/)
- [mcpkit](https://github.com/agentplexus/mcpkit) - Library-first MCP runtime
- [omnistorage](https://github.com/grokify/omnistorage) - Multi-backend storage abstraction
