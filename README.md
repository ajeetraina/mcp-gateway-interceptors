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


## Understanding the logs

```
[+] Running 6/6
 ✔ mcp-gateway-interceptors-interceptor              Built                                                  0.0s 
 ✔ mcp-gateway-interceptors-gateway                  Built                                                  0.0s 
 ✔ mcp-gateway-interceptors-client                   Built                                                  0.0s 
 ✔ Container mcp-gateway-interceptors-interceptor-1  Recreated                                              0.1s 
 ✔ Container mcp-gateway-interceptors-gateway-1      Recr...                                                0.1s 
 ✔ Container mcp-gateway-interceptors-client-1       Recre...                                               0.1s 
Attaching to client-1, gateway-1, interceptor-1
gateway-1  | - Reading configuration...
gateway-1  |   - Reading catalog from [https://desktop.docker.com/mcp/catalog/v2/catalog.yaml]
gateway-1  | - Reading configuration...
gateway-1  |   - Reading catalog from [https://desktop.docker.com/mcp/catalog/v2/catalog.yaml]
gateway-1  | - Configuration read in 87.988042ms
gateway-1  | - Interceptors enabled: before:exec:echo ============BEFORE==============: Query=$(jq -r ".params.arguments.query") >&2, before:http:http://interceptor:8080/before, after:http:http://interceptor:8080/after
gateway-1  | - Those servers are enabled: duckduckgo
gateway-1  | - Listing MCP tools...
gateway-1  |   - Running mcp/duckduckgo with [run --rm -i --init --security-opt no-new-privileges --cpus 1 --memory 2Gb --pull never -l docker-mcp=true -l docker-mcp-tool-type=mcp -l docker-mcp-name=duckduckgo -l docker-mcp-transport=stdio]
gateway-1  | - Configuration read in 87.988042ms
gateway-1  | - Interceptors enabled: before:exec:echo ============BEFORE==============: Query=$(jq -r ".params.arguments.query") >&2, before:http:http://interceptor:8080/before, after:http:http://interceptor:8080/after
gateway-1  | - Those servers are enabled: duckduckgo
gateway-1  | - Listing MCP tools...
gateway-1  |   - Running mcp/duckduckgo with [run --rm -i --init --security-opt no-new-privileges --cpus 1 --memory 2Gb --pull never -l docker-mcp=true -l docker-mcp-tool-type=mcp -l docker-mcp-name=duckduckgo -l docker-mcp-transport=stdio]
gateway-1  |   > duckduckgo: (2 tools)
gateway-1  |   > duckduckgo: (2 tools)
gateway-1  | > 2 tools listed in 1.026639542s
gateway-1  | - Using images: View Config   w Enable Watch
gateway-1  |   - mcp/duckduckgo@sha256:68eb20db6109f5c312a695fc5ec3386ad15d93ffb765a0b4eb1baf4328dec14f
gateway-1  | > 2 tools listed in 1.026639542s
gateway-1  | - Using images:
gateway-1  |   - mcp/duckduckgo@sha256:68eb20db6109f5c312a695fc5ec3386ad15d93ffb765a0b4eb1baf4328dec14f
gateway-1  | > Images pulled in 21.786542ms
gateway-1  | - Verifying images [mcp/duckduckgo]
gateway-1  | > Images pulled in 21.786542ms
gateway-1  | - Verifying images [mcp/duckduckgo]
gateway-1  | > Images verified in 3.680085043s
gateway-1  | > Images verified in 3.680085043s
gateway-1  | > Initialized in 4.823202669s
gateway-1  | > Start streaming server on port 9011
gateway-1  | > Initialized in 4.823202669s
gateway-1  | > Start streaming server on port 9011
gateway-1  | - Client initialized Config   w Enable Watch
gateway-1  | - Client initialized
gateway-1  |   - ============BEFORE==============: Query=null
gateway-1  |   - ============BEFORE==============: Query=null
gateway-1  |   - Calling tool search with arguments: {"query":"Docker"}
gateway-1  |   - Scanning tool call arguments for secrets...
gateway-1  |   > No secret found in arguments.
gateway-1  |   - Running mcp/duckduckgo with [run --rm -i --init --security-opt no-new-privileges --cpus 1 --memory 2Gb --pull never -l docker-mcp=true -l docker-mcp-tool-type=mcp -l docker-mcp-name=duckduckgo -l docker-mcp-transport=stdio --network mcp-gateway-interceptors_default]
gateway-1  |   - Calling tool search with arguments: {"query":"Docker"}
gateway-1  |   - Scanning tool call arguments for secrets...
gateway-1  |   > No secret found in arguments.nable Watch
gateway-1  |   - Running mcp/duckduckgo with [run --rm -i --init --security-opt no-new-privileges --cpus 1 --memory 2Gb --pull never -l docker-mcp=true -l docker-mcp-tool-type=mcp -l docker-mcp-name=duckduckgo -l docker-mcp-transport=stdio --network mcp-gateway-interceptors_default]
gateway-1  |   - Scanning tool call response for secrets...
gateway-1  |   - Scanning tool call response for secrets...
gateway-1  |   > No secret found in response.
gateway-1  |   > Calling tool search took: 1.552597793sch
gateway-1  |   > No secret found in response.
gateway-1  |   > Calling tool search took: 1.552597793s
client-1   | Found 10 search results:
client-1   | 
client-1   | 1. Docker: Accelerated Container Application Development
client-1   |    URL: https://www.docker.com/
client-1   |    Summary: Dockeris a tool that helps developers build, share, run, and verify applications using containers. Learn how to useDockerDesktop,DockerHub,DockerScout, and other products and services to accelerate your development and secure your workflows.
client-1   | 
client-1   | 2. Docker Desktop: The #1 Containerization Tool for Developers | Docker
client-1   |    URL: https://www.docker.com/products/docker-desktop/
client-1   |    Summary: DockerDesktop is a powerful platform for building, running, and managing containers on your local machine. It integrates with your development tools, supports Kubernetes, offers extensions, and enhances security and performance.
client-1   | 
client-1   | 3. Docker (software) - Wikipedia
client-1   |    URL: https://en.wikipedia.org/wiki/Docker_(software)
client-1   |    Summary: Former logoDockeris a set of platform as a service (PaaS) products that use OS-level virtualization to deliver software in packages called containers. [5] The service has both free and premium tiers. The software that hosts the containers is calledDockerEngine. [6] It was first released in 2013 and is developed byDocker, Inc. [7]Dockeris a tool that is used to automate the deployment of ...
client-1   | 
client-1   | 4. Get Started | Docker
client-1   |    URL: https://www.docker.com/get-started/
client-1   |    Summary: DockerDesktop lets you installDockerand customize your development environment with tools that enhance your tech stack and optimize your process. Learn how to useDockerCLI, IDE integrations, AI/ML, Trusted Open Source Content, andDockerHub.
client-1   | 
client-1   | 5. Install | Docker Docs
client-1   |    URL: https://docs.docker.com/engine/install/
client-1   |    Summary: Learn how to choose the best method for you to installDockerEngine. This client-server application is available on Linux, Mac, Windows, and as a static binary.
client-1   | 
client-1   | 6. Docker Personal - Sign Up for Free | Docker
client-1   |    URL: https://www.docker.com/products/personal/
client-1   |    Summary: DockerPersonal offers free access to an intuitive platform allowing developers to build, share, and run cloud-native applications.
client-1   | 
client-1   | 7. Docker Engine | Docker Docs
client-1   |    URL: https://docs.docker.com/engine/
client-1   |    Summary: Find a comprehensive overview ofDockerEngine, including how to install, storage details, networking, and more
client-1   | 
client-1   | 8. Docker Hub
client-1   |    URL: https://hub.docker.com/welcome
client-1   |    Summary: DockerHub is a central repository for finding, sharing, and managing container images and applications with ease.
client-1   | 
client-1   | 9. Get Docker | Docker Docs
client-1   |    URL: https://docs.docker.com/get-started/get-docker/
client-1   |    Summary: Download and installDockeron the platform of your choice, including Mac, Linux, or Windows.
client-1   | 
client-1   | 10. Products | Docker
client-1   |    URL: https://www.docker.com/products/
client-1   |    Summary: Dockeroffers a suite of integrated tools for building, securing, and deploying containerized applications. Learn aboutDockerDesktop,DockerHub,DockerScout, andDockerBuild Cloud, and how they can improve your software development workflow.
client-1   | 
client-1   | Found 10 search results:
client-1   | 
client-1   | 1. Docker: Accelerated Container Application Development
client-1   |    URL: https://www.docker.com/
client-1   |    Summary: Dockeris a tool that helps developers build, share, run, and verify applications using containers. Learn how to useDockerDesktop,DockerHub,DockerScout, and other products and services to accelerate your development and secure your workflows.
client-1   | 
client-1   | 2. Docker Desktop: The #1 Containerization Tool for Developers | Docker
client-1   |    URL: https://www.docker.com/products/docker-desktop/
client-1   |    Summary: DockerDesktop is a powerful platform for building, running, and managing containers on your local machine. It integrates with your development tools, supports Kubernetes, offers extensions, and enhances security and performance.
client-1   | 
client-1   | 3. Docker (software) - Wikipedia
client-1   |    URL: https://en.wikipedia.org/wiki/Docker_(software)
client-1   |    Summary: Former logoDockeris a set of platform as a service (PaaS) products that use OS-level virtualization to deliver software in packages called containers. [5] The service has both free and premium tiers. The software that hosts the containers is calledDockerEngine. [6] It was first released in 2013 and is developed byDocker, Inc. [7]Dockeris a tool that is used to automate the deployment of ...
client-1   | 
client-1   | 4. Get Started | Docker
client-1   |    URL: https://www.docker.com/get-started/
client-1   |    Summary: DockerDesktop lets you installDockerand customize your development environment with tools that enhance your tech stack and optimize your process. Learn how to useDockerCLI, IDE integrations, AI/ML, Trusted Open Source Content, andDockerHub.
client-1   | 
client-1   | 5. Install | Docker Docs
client-1   |    URL: https://docs.docker.com/engine/install/
client-1   |    Summary: Learn how to choose the best method for you to installDockerEngine. This client-server application is available on Linux, Mac, Windows, and as a static binary.
client-1   | 
client-1   | 6. Docker Personal - Sign Up for Free | Docker
client-1   |    URL: https://www.docker.com/products/personal/
client-1   |    Summary: DockerPersonal offers free access to an intuitive platform allowing developers to build, share, and run cloud-native applications.
client-1   | 
client-1   | 7. Docker Engine | Docker Docs
client-1   |    URL: https://docs.docker.com/engine/
client-1   |    Summary: Find a comprehensive overview ofDockerEngine, including how to install, storage details, networking, and more
client-1   | 
client-1   | 8. Docker Hub
client-1   |    URL: https://hub.docker.com/welcome
client-1   |    Summary: DockerHub is a central repository for finding, sharing, and managing container images and applications with ease.
client-1   | 
client-1   | 9. Get Docker | Docker Docs
client-1   |    URL: https://docs.docker.com/get-started/get-docker/
client-1   |    Summary: Download and installDockeron the platform of your choice, including Mac, Linux, or Windows.
client-1   | 
client-1   | 10. Products | Docker
client-1   |    URL: https://www.docker.com/products/
client-1   |    Summary: Dockeroffers a suite of integrated tools for building, securing, and deploying containerized applications. Learn aboutDockerDesktop,DockerHub,DockerScout, andDockerBuild Cloud, and how they can improve your software development workflow.
client-1   | 
client-1 exited with code 0
```

