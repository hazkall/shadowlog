# ShadowLog Application

## Overview
ShadowLog is a Go-based application designed to generate fake HTTP logs. It includes features for OpenTelemetry tracing and metrics collection, providing insights into the application's behavior over time. The application operates by running a server and periodically making HTTP GET requests to a specified endpoint, logging the results.

## Features
OpenTelemetry Integration: Includes tracing and metrics to monitor application performance and behavior.
Customizable Interval: The frequency of HTTP requests can be customized.
Configurable Server Port: The server can run on a port specified by the user.
Automated HTTP Requests: Periodically sends HTTP GET requests to a specified endpoint and logs the responses.

## Installation

```bash
git clone https://github.com/hazkall/shadowlog.git
cd shadowlog

go build -o shadowlog .

./shadowlog run --interval 2 --port 3000
```

## Docker Compose

This Docker Compose configuration sets up a full observability stack for the ShadowLog application, including OpenTelemetry, Prometheus, and Jaeger. The setup is designed to facilitate monitoring, tracing, and logging of the ShadowLog service.

```bash

docker compose up

```

ShadowLog Application: http://localhost:3000
Prometheus: http://localhost:9090
Jaeger UI: http://localhost:16686

## Docker Build

```bash

docker build --target server -t shadowlog .

```

## Usage
The application can be executed using the shadowlog command with the run subcommand. You can specify the interval for making requests and the port on which the server should run.

Command-Line Flags
--interval, -i: The interval (in seconds) at which HTTP GET requests are made to the server. This flag is required. (Default: 2 seconds)
--port, -p: The port number on which the server will run. This flag is required. (Default: 3000)

## License
This project is licensed under the MIT License.
