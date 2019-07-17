package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	SubscriptionThingsAll            = "all"
	SubscriptionSubmissionsAll       = "submissions_all"
	SubscriptionSubmissionsSubreddit = "submissions_subreddit"
	SubscriptionSubmissionsUser      = "submissions_user"
	SubscriptionCommentsAll          = "comments_all"
	SubscriptionCommentsSubreddit    = "comments_subreddit"
	SubscriptionCommentsSubmission   = "comments_submission"
	SubscriptionCommentsUser         = "comments_user"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Disable CORS for testing
		},
	}
)

func serveWS(hub *Hub, subscription string) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		query := getSubscriptionQuery(c, subscription)
		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), subscriptionType: subscription, subscriptionQuery: query}
		client.hub.register <- client
		fmt.Printf("Client connected: %s", c.RealIP())

		go client.writePump()
		go client.readPump()
		return nil
	}
}

func newRouter(hub *Hub) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.File("/ui", "static/index.html")
	e.GET("/", serveWS(hub, SubscriptionThingsAll))
	e.GET("/submissions", serveWS(hub, SubscriptionSubmissionsAll))
	e.GET("/comments", serveWS(hub, SubscriptionCommentsAll))
	e.GET("/r/:subreddit", serveWS(hub, SubscriptionSubmissionsSubreddit))
	e.GET("/r/:subreddit/comments", serveWS(hub, SubscriptionCommentsSubreddit))
	e.GET("/submission/:submissionId/comments", serveWS(hub, SubscriptionCommentsSubmission))
	e.GET("/u/:username/submissions", serveWS(hub, SubscriptionSubmissionsUser))
	e.GET("/u/:username/comments", serveWS(hub, SubscriptionCommentsUser))
	return e
}

func NewWebsocketServer(host string, port int) error {
	hub := newHub()
	go hub.run()
	firehose := newFirehose()
	go broadCastMessages(firehose, hub)
	httpServer := newRouter(hub)
	return httpServer.Start(fmt.Sprintf("%s:%d", host, port))
}
