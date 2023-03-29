# WebSocket Server

## Overview
This is a WebSocket server implemented in Go. The server listens for incoming WebSocket connections from clients, handles upgrade requests, and sends HTTP requests to a lambda server. It receives responses from the lambda server and sends them back to the client.

## Features
- Handles WebSocket upgrade requests from clients.
- Listens for incoming WebSocket connections on `localhost:4999`.
- Sends HTTP requests to the lambda server located at `http://localhost:5000/run/echo`.
- Forwards responses from the lambda server back to the client over the WebSocket connection.

## Module and Dependencies
run `go get` to install the following dependencies:
```shell
go get github.com/gobwas
go get github.com/mailru/easygo
```

## Usage
simply run wscontroller.go to start the server:
```shell
go run wscontroller.go
```

## TODO
- [ ] Add support for multiple lambda servers.
- [ ] Deal with websocket status codes.
- [ ] Solve err handling.