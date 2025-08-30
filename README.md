# Demonstrating Docker MCP Gateway Interceptors

- The Docker MCP (Model Context Protocol) Gateway is a middleware layer that sits between AI assistants and the tools they interact with. 
- It acts as a proxy that manages and controls how AI agents execute tool calls, providing a standardized interface for tool execution while adding powerful interception capabilities.

## Purpose and Architecture of Interceptors

Interceptors are middleware components that allow you to:

- Monitor, modify, or bypass tool calls in real-time
- Add logging, validation, or security checks
- Transform requests and responses
- Implement custom business logic without modifying the underlying tools


## The Right Architecture

```
ALL requests → MCP Gateway (with interceptors) → Backend/Tools
```

## Prerequisite

- Docker Desktop 4.43.X

## Clone the repo

```console
git clone https://github.com/ajeetraina/mcp-gateway-interceptors
cd  mcp-gateway-interceptors
```

## How to run

```console
docker compose up --build
```

# Types of Interceptors

There are three types of interceptors:
- `exec`, 
- `docker` and 
- `http`.

Interceptors can run `before` a tool call or `after` a tool call`.

Those which run run `before` have access to the full tool call request and
can either let the call go through or bypass the call and return a custom response.

Those which run run `after` have access to the tool call response.

## `exec`

Usage: `--interceptor=before:exec:script` or `--interceptor=after:exec:script`.

The `script` is a shell script that will run with `/bin/sh -c`. e.g:

```
--interceptor=before:exec:echo Query=$(jq -r ".params.arguments.query") >&2
```

The tool call request (`before`) or tool call response (`after`) are passed as json objects into stdin.
To return a custom response, the interceptor needs to write it to `stdout` as a json object.
Every output sent to `stderr` will be shown in the gateway's logs.

## `docker`

Usage: `--interceptor=before:docker:image arg1 arg2` or `--interceptor=after:docker:image arg1 arg2`.

e.g:

```
--interceptor=before:docker:alpine sh -c 'echo BEFORE >&2'
```

The tool call request (`before`) or tool call response (`after`) are passed as json objects into stdin.
To return a custom response, the interceptor needs to write it to `stdout` as a json object.
Every output sent to `stderr` will be shown in the gateway's logs.

## `http`

Usage: `--interceptor=before:http:http://host:port/path` or `--interceptor=after:http:http://host:port/path`.

e.g:

```
--interceptor=before:http:http://interceptor:8080/before
--interceptor=after:http:http://interceptor:8080/after
```

The tool call request (`before`) or tool call response (`after`) are passed as json objects into a `POST` request.
To return a custom response, the interceptor needs to write a non empty json object.

# Examples

Log the tool request's arguments:

```yaml
- --interceptor
- before:exec:echo Arguments=$(jq -r ".params.arguments") >&2
```

Log the tool call's response:

```yaml
- --interceptor
- after:exec:echo Response=$(jq -r ".") >&2
```

Trim down the tool's response text:

```yaml
- --interceptor
- after:exec:jq -c '.content[].text |= (.[:100])'





