package main

import (
	"fmt"
	"net/http"
	"log"
	"html/template"
	"github.com/gorilla/websocket"
)

var tpl *template.Template
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}


func init() {
	tpl = template.Must(template.ParseGlob("public/*.html"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/ws", wsEndPoint)
	//http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	http.Handle("/public/css/", http.StripPrefix("/public/css", http.FileServer(http.Dir("public/css"))))

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Println("home")
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
		if err := conn.WriteMessage(messageType, p); err != nil {
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