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
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"sse-test/model"
	"time"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) InitializeRoutes() {
	var err error
	server.Router = mux.NewRouter()
	server.Router.HandleFunc("/event", sseHandler)
	server.Router.HandleFunc("/send", SetMiddlewareJSON(server.sendMessage)).Methods("POST")
	server.Router.HandleFunc("/getMessages", SetMiddlewareJSON(server.GetAllMessages)).Methods("GET")

	dsn := "host=localhost user=admin password=admin dbname=sse-test port=5433 sslmode=disable TimeZone=Asia/Shanghai"
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("Cannot connect to %s database", "23")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database \n", "5")
	}
	server.DB.Debug().AutoMigrate(model.Message{}) //migrations
	log.Fatal(http.ListenAndServe(":3500", server.Router))
}

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func main() {
	router := http.NewServeMux()
	var server Server
	router.HandleFunc("/event", sseHandler)
	//router.HandleFunc("/time", getTime)
	router.HandleFunc("/send", SetMiddlewareJSON(server.sendMessage))
	server.InitializeRoutes()
	//log.Fatal(http.ListenAndServe(":3500", router))

}

var msgChan chan model.Message
var testChan chan string

func getTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if testChan != nil {
		msg := time.Now().Format("15:04:05")
		testChan <- msg
	}
}

func getTime2() {
	msg := time.Now().Format("15:04:05")
	for {
		time.Sleep(time.Second)
		msg = time.Now().Format("15:04:05")
		testChan <- msg
	}
}

func (server *Server) sendMessage(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, ": error occur\n\n")
	}
	var message model.Message
	err = json.Unmarshal(body, &message)
	if err != nil {
		fmt.Fprintf(w, ": error occur in parsing json\n\n")
	}
	messageCreated, err := message.SendMessage(server.DB)
	if err != nil {
		fmt.Fprintf(w, "error eith database")
	}
	//msg := message.Text
	msgChan <- message
	JSON(w, http.StatusOK, messageCreated)
}

func sseHandler(w http.ResponseWriter, r *http.Request) {

	//go getTime2()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	msgChan = make(chan model.Message)

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
			fmt.Fprintf(w, "data:%s\n\n", message.Text)
			flusher.Flush()
		case <-r.Context().Done():
			fmt.Println("Client closed connection")
			return
		case <-timeout:
			fmt.Fprintf(w, ": nothing to sent\n\n")
		}
	}
}

func (server *Server) GetAllMessages(w http.ResponseWriter, r *http.Request) {
	var message model.Message
	messageCreated, err := message.GetAllMessages(server.DB)
	if err != nil {
		fmt.Fprintf(w, "error eith database while l9oading all messsages")
	}
	JSON(w, http.StatusOK, messageCreated)
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}
