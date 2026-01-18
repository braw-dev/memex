# Memex AI

Memex is an L1 cache for LLM coding tools like Claude Code.

It aims to speed up workflows and reduce token usage through smart caching of repeat/similar prompts and requests.

Memex can be run locally or on an external server. When running on an external server the cache can be shared between people using the proxy.

At a high level, Memex works by intercepting requests and using multiple data points (system prompt, user prompt similarity, file tree) to see if a request is in the cache. If so, it returns it immediately effectively bypassing the slow request to the upstream LLM.

You can also enable Memex to proxy and cache MCP server calls. Memex can be configured to respect `Cache-Control` headers or to force cache all requests for a certain time.

As Memex is a proxy it is shared between tools and agents (when they are configured to use it). If Claude Code makes a request to access the MCP documentation server for a library, that request is cached enabling other agents or tools (e.g. Cursor) to immediately receive the response when asking for the same docs.

## Cache Hits

Caching is a hard problem to solve and can have issues such as returning stale results. To make it clear when a response is being served by Memex we add a footer to each message saying so. If you are debugging why your responses seem stale please check for the presence of this message and consider clearing your cache with the [Hot Command](#hot-commands).

## Hot Commands

You can use some hot commands to alter the behaviour with Memex. For example, including `!memex:skip` anywhere in your user message, will bypass the Memex cache. You can embed these hot commands in your custom rules or commands as well.

If the hot command is the only text in the message Memex will reply with an acknowledgement or status message. If there is other content in the message Memex will forward this to the upstream LLM as normal.

| Command | Purpose |
| -- | ------ |
| `!memex:skip` | Bypass the Memex cache |
| `!memex:bust` | Clear/bust the Memex cache. This will remove all cache entries, essentially starting from scratch. If there is other content in the message this will be forwarded on as normal. |

## Frequently Asked Questions (FAQ)

### How does this differ from caching built in to the LLM providers?

LLM caching typically only caches user inputs for a short period of time (e.g. 5 minutes [for Anthropic](https://platform.claude.com/docs/en/build-with-claude/prompt-caching#how-prompt-caching-works)). It is used to cache system prompts to speed up LLM usage and reduce costs (although there is a cost associated with cache reads and writes).

Memex caches full requests locally. If it finds a cache hit it does not send the request to the LLM. This is significantly faster (milliseconds instead of seconds) and does not eat into your credits/usage limits.

An example using a shared Memex as part of a team editing a codebase:

- You are working on a team project with Claude and ask how authentication works.
- The request is proxied through Memex which has not seen this before so forwards the request on to the upstream model.
- On their own machine, a colleague also asks how authentication works (likely in different words).
- Memex identifies that this request is similar to one it has seen before and returns the cached result. The request does not hit the upstream model.

An example using Memex locally:

- You
