package WEST

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"WEST"},
}

var methods = []string{
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"DELETE",
	"PATCH",
}

func Listen(host string, port int, handler http.Handler) {
	loggedHandler := LoggerHandler(handler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		defer c.Close()

		for {
			mt, message, err := c.ReadMessage()

			if err != nil {
				log.Println("read:", err)
				break
			}

			west, err := processRequest(message)

			if err != nil {
				if west != nil {
					err = c.WriteMessage(mt, []byte(west.Id+" 400 "+err.Error()))

					if err != nil {
						log.Println("write:", err)
						break
					}
				}

				continue
			}

			writer := MakeWestWriter()

			loggedHandler.ServeHTTP(writer, west.toHTTPRequest(r))

			err = c.WriteMessage(mt, []byte(west.Id+" "+strconv.Itoa(writer.status)+" "+writer.buffer.String()))

			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	})

	addr := host + ":" + strconv.Itoa(port)

	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
