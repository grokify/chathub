# ChatHub: ChatGPT Desktop ↔ Claude Code Integration

This guide explains how to share AI conversations between ChatGPT Desktop and Claude Code using ChatHub with GitHub as the storage backend.

## Overview

```
┌───────────────────────┐    ┌──────────────────────┐
│    ChatGPT Desktop    │    │      Claude Code     │
│                       │    │                      │
│  save_conversation    |    │  read_conversation   |
│  append_conversation  |    │  list_conversations  |
└────────┬──────────────┘    └────────┬─────────────┘
         │                            │
         │      MCP Protocol          │
         ▼                            ▼
    ┌─────────────────────────────────────┐
    │            ChatHub MCP              │
    │         (same binary/config)        │
    └─────────────────┬───────────────────┘
                      │
                      ▼
              ┌───────────────┐
              │    GitHub     │
              │  Repository   │
              │               │
              │ conversations/│
              │   chatgpt/    │
              │   claude/     │
              └───────────────┘
```

## Prerequisites

- Go 1.23+ (to build ChatHub)
- GitHub account
- ChatGPT Desktop app with MCP support
- Claude Code CLI

## Step 1: Create GitHub Repository

Create a private repository to store your conversations:

```bash
# Using GitHub CLI
gh repo create ai-conversations --private --description "AI conversation archive"

# Or create manually at https://github.com/new
```

## Step 2: Create GitHub Personal Access Token

1. Go to https://github.com/settings/tokens?type=beta
2. Click **Generate new token**
3. Configure:
   - **Token name**: `chathub-access`
   - **Expiration**: 90 days (or custom)
   - **Repository access**: Select `ai-conversations`
   - **Permissions**:
     - **Contents**: Read and write
     - **Metadata**: Read

4. Click **Generate token**
5. Copy the token (starts with `github_pat_`)

> **Security tip**: Use a fine-grained PAT with minimal permissions rather than a classic PAT with broad `repo` scope.

## Step 3: Build ChatHub

```bash
# Clone and build
git clone https://github.com/grokify/chathub.git
cd chathub
go build -o chathub ./cmd/chathub

# Install to a stable location
sudo cp chathub /usr/local/bin/chathub

# Verify installation
/usr/local/bin/chathub --help
```

## Step 4: Configure ChatGPT Desktop

Add ChatHub to your ChatGPT Desktop MCP configuration.

**Location**: ChatGPT Desktop Settings → MCP Servers (or config file)

```json
{
  "mcpServers": {
    "chathub": {
      "command": "/usr/local/bin/chathub",
      "env": {
        "GITHUB_TOKEN": "github_pat_YOUR_TOKEN_HERE",
        "GITHUB_OWNER": "your-github-username",
        "GITHUB_REPO": "ai-conversations",
        "CHATHUB_BACKEND": "github",
        "CHATHUB_FOLDER": "conversations"
      }
    }
  }
}
```

**Restart ChatGPT Desktop** after adding the configuration.

## Step 5: Configure Claude Code

Add ChatHub to Claude Code's MCP configuration.

**Location**: `~/.claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "chathub": {
      "command": "/usr/local/bin/chathub",
      "env": {
        "GITHUB_TOKEN": "github_pat_YOUR_TOKEN_HERE",
        "GITHUB_OWNER": "your-github-username",
        "GITHUB_REPO": "ai-conversations",
        "CHATHUB_BACKEND": "github",
        "CHATHUB_FOLDER": "conversations"
      }
    }
  }
}
```

> **Important**: Use the **same token, owner, repo, and folder** in both configurations. This is what enables cross-platform sharing.

## Step 6: Verify Setup

### Test in ChatGPT Desktop

```
You: Save this conversation as "Test conversation"

ChatGPT: I'll save this conversation for you.
[Calls save_conversation tool]
Saved to: conversations/chatgpt/2026-01-10_test-conversation.md
```

### Test in Claude Code

```
You: List my saved conversations

Claude: Let me check your conversation archive.
[Calls list_conversations tool]
Found 1 conversation:
- conversations/chatgpt/2026-01-10_test-conversation.md
  Title: "Test conversation"
  Source: chatgpt
  Date: 2026-01-10
```

```
You: Read the test conversation

Claude: [Calls read_conversation tool]
Here's the conversation from ChatGPT...
```

## Usage Examples

### Save from ChatGPT

```
Save this conversation about React hooks
```

```
Archive our discussion with tags "react" and "frontend"
```

### Continue in Claude Code

```
What ChatGPT conversations do I have about React?
```

```
Read my React hooks conversation and summarize the key points
```

```
Based on my ChatGPT discussion, help me implement the useCallback pattern
```

### Append from Either Platform

```
Add this follow-up to my React hooks conversation
```

## File Structure in GitHub

After using ChatHub, your repository will look like:

```
ai-conversations/
└── conversations/
    ├── chatgpt/
    │   ├── 2026-01-10_react-hooks.md
    │   ├── 2026-01-11_debugging-tips.md
    │   └── 2026-01-12_api-design.md
    ├── claude/
    │   └── 2026-01-10_code-review.md
    └── claude-code/
        └── 2026-01-11_refactoring.md
```

Each file contains Hugo-compatible YAML frontmatter:

```yaml
---
title: "React Hooks Discussion"
date: 2026-01-10T14:30:00Z
source: chatgpt
tags: ["react", "frontend", "hooks"]
conversation_id: "conv_1234567890"
---

**User:** How do I use useCallback effectively?

**ChatGPT:** useCallback is a React hook that...
```

## Troubleshooting

### "GITHUB_TOKEN is required" error

Ensure the token is set in the MCP server config's `env` block, not as a system environment variable.

### "Repository not found" error

- Verify `GITHUB_OWNER` matches your GitHub username exactly
- Verify `GITHUB_REPO` matches your repository name exactly
- Check that your PAT has access to the repository

### "Permission denied" error

Your PAT needs **Contents: Read and write** permission. Regenerate with correct scopes.

### ChatHub not appearing in ChatGPT/Claude

- Verify the binary path exists: `ls -la /usr/local/bin/chathub`
- Check the binary is executable: `chmod +x /usr/local/bin/chathub`
- Restart the application after config changes

### Conversations not syncing

Both clients must use the **exact same** configuration:
- Same `GITHUB_TOKEN`
- Same `GITHUB_OWNER`
- Same `GITHUB_REPO`
- Same `CHATHUB_FOLDER`

## Security Considerations

1. **Use a private repository** for personal conversations
2. **Use fine-grained PATs** with minimal permissions
3. **Don't commit tokens** to version control
4. **Rotate tokens periodically** (set expiration dates)
5. **Review repository access** in GitHub settings

## Available Tools

| Tool | Description | Use Case |
|------|-------------|----------|
| `save_conversation` | Save new conversation | End of a chat session |
| `append_conversation` | Add to existing | Continue a topic |
| `read_conversation` | Read by path | Review past discussions |
| `list_conversations` | List with filters | Find conversations |
| `search_conversations` | Full-text search | Search by content |
| `delete_conversation` | Remove conversation | Clean up drafts |

## Next Steps

- **Hugo integration**: Use the same repo as a Hugo site to publish conversations as a searchable knowledge base
- **Additional backends**: Configure S3 or Dropbox for alternative storage
- **Team sharing**: Use a shared repository for team knowledge capture
