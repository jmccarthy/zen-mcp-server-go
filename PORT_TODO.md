# Porting Plan: Python MCP Server to Go

This document outlines the major tasks required to port the existing Python-based Model Context Protocol (MCP) server to Go. It summarizes the current architecture and provides a high-level plan for recreating each component in Go while preserving functionality.

## 1. Assess the Existing Server
- Review `server.py` for the main entry point, tool registry, provider configuration, and logging setup. The file defines the `TOOLS` dictionary and dispatches requests to specific tool implementations.【F:server.py†L133-L147】【F:server.py†L320-L377】
- Study `utils/conversation_memory.py` to understand how conversation threads are stored in Redis for cross‑tool continuation. This module provides thread creation, persistence, and reconstruction.【F:utils/conversation_memory.py†L1-L44】【F:utils/conversation_memory.py†L160-L207】
- Inspect providers under `providers/` for integrations with Gemini, OpenAI, OpenRouter, and custom endpoints. Each provider exposes methods to generate content and count tokens.
- Review utility modules such as `utils/file_utils.py` for file access safeguards and token management.

## 2. Outline the Go Architecture
1. **Project Layout**
   - Create a Go module (`go.mod`) and organize packages similar to the Python modules:
     - `cmd/zen-mcp` – main package running the server.
     - `internal/server` – JSON‑RPC handling and tool dispatching.
     - `internal/tools` – individual tool implementations (chat, analyze, etc.).
     - `internal/providers` – model provider interfaces and clients.
     - `internal/conversation` – conversation memory backed by Redis.
     - `internal/util` – helper functions (file handling, token counting).

2. **JSON‑RPC and Stdio**
   - Implement a JSON‑RPC 2.0 server over stdin/stdout similar to Python’s `mcp.server.stdio`. Packages like `github.com/sourcegraph/jsonrpc2` can simplify request handling.
   - Mirror the initialization handshake and tool discovery messages defined in `server.py`.

3. **Configuration & Logging**
   - Port constants from `config.py` and support environment variable overrides.
   - Use a structured logger (e.g., `log/slog` or `zap`) and `lumberjack` for rotating log files to replicate the logging setup (daily rotation and size limits).【F:server.py†L100-L127】

4. **Tool Interface**
   - Define a `Tool` interface with methods such as `Name() string`, `Description() string`, `InputSchema() map[string]any`, and `Execute(ctx context.Context, args map[string]any) ([]TextContent, error)`.
   - Recreate each tool from `tools/` in Go. Use `struct` types with JSON tags for request validation.
   - Include helper functions for assembling prompts from templates located in `prompts/tool_prompts.py`.

5. **Providers**
   - Design a provider interface mirroring `providers.base.ModelProvider` with methods for capability lookup, token counting, and content generation.
   - Implement clients for Gemini and OpenAI using Go’s `net/http`. Respect temperature constraints and model‑specific options as seen in `providers/base.py`.
   - Add an OpenRouter client and a generic “custom” provider for user‑defined endpoints.

6. **Conversation Memory**
   - Recreate Redis‑backed conversation threads. Each thread stores turns, files, model metadata, and is referenced via a UUID. Support thread creation, retrieval, turn addition, and expiration logic as described in `utils/conversation_memory.py`.

7. **File Utilities & Security**
   - Port `utils/file_utils.py` to securely read files from the workspace while preventing directory traversal. Implement directory exclusion and token budget checks.
   - Provide token estimation functions similar to `utils/token_utils.py` (may leverage a tokenizer library or approximate using word counts).

8. **Prompt Handling**
   - Translate system prompts and templates from the `prompts` package to Go constants or template files.
   - Ensure each tool combines system prompts with user content and optional context files, just like in `ChatTool`’s `prepare_prompt` method.【F:tools/chat.py†L44-L71】【F:tools/chat.py†L116-L147】

9. **Testing Strategy**
   - Replicate Python unit tests in Go using the `testing` package. Focus on:
     - Tool dispatch and unknown tool handling.
     - Provider registration and validation.
     - Conversation continuation across tools.
     - File security and token limit enforcement.
   - Add integration tests that simulate end‑to‑end JSON‑RPC exchanges.

10. **Documentation & Examples**
    - Update the README with instructions for building and running the Go server.
    - Provide example configuration files analogous to `claude_config_example.json`.

## 3. Migration Steps
1. Start by scaffolding the Go module and basic JSON‑RPC server.
2. Incrementally port each tool, beginning with simpler ones like `get_version` and `chat`.
3. Implement provider clients and plug them into the tool execution flow.
4. Add Redis conversation support and file utilities.
5. Recreate logging and configuration logic.
6. Gradually replace Python tests with Go tests to confirm parity.
7. Once feature complete, deprecate the Python entry point and document how to run the Go version.

## 4. Open Questions / Future Work
- Determine how to handle token counting accurately in Go—may require a third‑party tokenizer or a small Python helper service.
- Evaluate concurrency needs: Go’s goroutines can handle parallel tool execution if needed.
- Consider providing a compatibility layer so existing MCP clients can switch seamlessly to the Go server.

---
This plan should serve as a starting point for a structured port of the Python MCP server to Go, retaining all existing functionality while leveraging Go’s performance and concurrency advantages.
