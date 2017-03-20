package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"github.com/olorin/nagiosplugin"
	"log"
	"net/url"
)

var scheme = flag.String("scheme", "wss", "Websocket scheme (default to secure)")
var addr = flag.String("addr", "echo.websocket.org", "http service address")
var path = flag.String("path", "/", "URI path")

func main() {
	flag.Parse()

	u := url.URL{Scheme: *scheme, Host: *addr, Path: *path}

	// Initialize the check - this will return an UNKNOWN result
	// until more results are added.
	check := nagiosplugin.NewCheck()
	// If we exit early or panic() we'll still output a result.
	defer check.Finish()

	// Connect to the server
	log.Printf("Connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		check.Criticalf("Cannot connect to websocket")
		log.Fatal("dial:", err)
	} else {
		log.Println("connected")
	}
	defer c.Close()

	// close the connection
	log.Println("Requesting connection closure")
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		check.Criticalf("Error while closing websocket")
		return
	} else {
		check.AddResult(nagiosplugin.OK, "Websocket server is working correctly")
	}

}
