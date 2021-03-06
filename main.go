package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("public/*.html"))
}

func main() {
	// create hub
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", index)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	//http.HandleFunc("/ws", wsEndPoint)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func index(w http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(p))
		reply := string(p) + " receive message from client and send it back!"
		if err := conn.WriteMessage(messageType, []byte(reply)); err != nil {
			log.Println(err)
			return
		}
	}
}

func wsEndPoint(w http.ResponseWriter, req *http.Request) {
	fmt.Println("wsEndPoint")
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("client successfully connected...")
	reader(conn)
}
