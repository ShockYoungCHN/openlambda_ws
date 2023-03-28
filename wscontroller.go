package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/mailru/easygo/netpoll"
	"io"
	"log"
	"net"
	"net/http"
)

type Client struct {
	conn net.Conn
	send chan []byte
}

var poller netpoll.Poller

// serveWs handles websocket requests from the clients
func WsHandler(conn net.Conn) {
	_, err := ws.Upgrade(conn)
	if err != nil {
		log.Printf("upgrade error: %s", err)
		return
	}

	client := &Client{conn, make(chan []byte, 256)}
	/*	payload, err := wsutil.ReadClientText(client.conn)
		fmt.Println(payload)*/

	// Get netpoll descriptor with EventRead|EventEdgeTriggered.
	desc := netpoll.Must(netpoll.HandleRead(conn))

	err = poller.Start(desc, func(ev netpoll.Event) {
		go send(client)
	})
	if err != nil {
		return
	}
}

// send sends the http request to the lambda server and sends the response to the client
func send(client *Client) {
	buf := bufio.NewReadWriter(bufio.NewReader(client.conn), bufio.NewWriter(client.conn))
	body, _, err := sendRequest(buf)
	if err != nil {
		log.Println(err)
		body = []byte(err.Error())
	}

	err = wsutil.WriteServerText(client.conn, body)
	if err != nil {
		log.Printf("failed to write WebSocket message: %s", err)
	}
	client.conn.Close()
}

// sendRequest sends an HTTP request to the lambda server and returns the response body
func sendRequest(buf io.ReadWriter) ([]byte, int, error) {
	funcName := "echo"
	url := "http://localhost:5000/run/" + funcName

	//todo: maybe this can be optimized?
	payload, err := wsutil.ReadClientText(buf)
	if err != nil {
		return nil, -1, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, -1, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	var respBuf bytes.Buffer
	_, err = io.Copy(&respBuf, resp.Body)
	if err != nil {
		return nil, -1, err
	}

	return respBuf.Bytes(), resp.StatusCode, nil
}

func startPoller() {
	newPoller, err := netpoll.New(nil)
	poller = newPoller
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:4999")
	// todo: get the port dynamically
	if err != nil {
		log.Fatal(err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Println(err)
		}
	}(listener)

	startPoller()
	for {
		fmt.Print("running\n")
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		WsHandler(conn)
	}
}
