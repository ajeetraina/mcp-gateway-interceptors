# Demonstrating Docker MCP Gateway Interceptors

The interceptors are configurable filters that sit between AI clients and MCP tools, allowing you to inspect what tools are being called and with what data. 
They're essentially middleware hooks that let MCP Horror Stories: The GitHub Prompt Injection Data Heist you inspect, modify, or block tool calls in real-time.

The key insight from the GitHub attack prevention example is that interceptors work on MCP tool calls, not on general HTTP requests. When an AI agent makes its first GitHub tool call through the Gateway, the script records that repository in a session file. Any subsequent attempts to access a different repository get blocked with a security alert.

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
```

