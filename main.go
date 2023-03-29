/*
package main

import (

	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"

)

	type Client struct {
		name   string
		events chan *DashBoard
	}

	type DashBoard struct {
		User uint
	}

	func main() {
		app := fiber.New()
		app.Get("/sse", adaptor.HTTPHandler(handler(dashboardHandler)))
		app.Listen(":3000")
	}

	func handler(f http.HandlerFunc) http.Handler {
		return http.HandlerFunc(f)
	}

	func dashboardHandler(w http.ResponseWriter, r *http.Request) {
		client := &Client{name: r.RemoteAddr, events: make(chan *DashBoard, 10)}
		go updateDashboard(client)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		timeout := time.After(1 * time.Second)
		select {
		case ev := <-client.events:
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			enc.Encode(ev)
			fmt.Fprintf(w, "data: %v\n\n", buf.String())
			fmt.Printf("data: %v\n", buf.String())
		case <-timeout:
			fmt.Fprintf(w, ": nothing to sent\n\n")
		}

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	func updateDashboard(client *Client) {
		for {
			db := &DashBoard{
				User: uint(rand.Uint32()),
			}
			client.events <- db
		}
	}
*/
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type RequestData struct {
	Text string `json:"text"`
}

var msgChan chan string

func getTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if msgChan != nil {
		msg := time.Now().Format("15:04:05")
		msgChan <- msg
	}
}

func getTime2() {
	msg := time.Now().Format("15:04:05")
	for {
		time.Sleep(time.Second)
		msg = time.Now().Format("15:04:05")
		msgChan <- msg
	}
}
func sseHandler(w http.ResponseWriter, r *http.Request) {

	//go getTime2()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	msgChan = make(chan string)

	defer func() {
		close(msgChan)
		msgChan = nil
		fmt.Println("Client closed connection")
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Println("Could not init http.Flusher")
	}

	timeout := time.After(1 * time.Second)

	for {
		select {
		case message := <-msgChan:
			fmt.Fprintf(w, "data:%s\n\n", message)
			flusher.Flush()
		case <-r.Context().Done():
			fmt.Println("Client closed connection")
			return
		case <-timeout:
			fmt.Fprintf(w, ": nothing to sent\n\n")
		}
	}
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, ": error occur\n\n")
	}
	var requestData RequestData
	msg := ""
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		fmt.Fprintf(w, ": error occur in parsing json\n\n")
	}
	msg = requestData.Text
	msgChan <- msg
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/event", sseHandler)
	//router.HandleFunc("/time", getTime)
	router.HandleFunc("/send", sendMessage)

	log.Fatal(http.ListenAndServe(":3500", router))
}
