package ws

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

// checkOrigin will check origin and return true if its allowed
func checkOrigin(r *http.Request) bool {

	// Grab the request origin
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:9090":
		return true
	default:
		return false
	}
}

var upgrader = websocket.Upgrader{
	//CheckOrigin:     checkOrigin,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Manager struct {
	posts PostList
	sync.RWMutex
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		posts:    make(PostList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

// setupEventHandlers configures and adds all handlers
func (m *Manager) setupEventHandlers() {
	m.handlers[EventCreateComment] = func(e Event, c *Client) error {
		fmt.Println(e)
		return nil
	}
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *Manager) routeEvent(event Event, c *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new ws connection")

	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("query parameter 'post_id' cannot be empty"))
		return
	}

	// upgrade regular http connection into ws
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error occurred while upgrading connection into ws: %s", err)
		return
	}

	client := NewClient(conn, m, postID)

	m.AddClient(postID, client)

	go client.ReadMessages()
	go client.writeMessages()
}

func (m *Manager) AddClient(postId string, c *Client) {
	m.Lock()
	defer m.Unlock()

	if post, ok := m.posts[postId]; ok {
		post[c] = true
	} else {
		m.posts[postId] = make(ClientList)
		m.posts[postId][c] = true
	}
}

func (m *Manager) RemoveClient(c *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.posts[c.PostID][c]; ok {
		err := c.connection.Close()
		if err != nil {
			log.Printf("error occurred while closing ws connection: %s", err)
			return
		}

		delete(m.posts[c.PostID], c)
	}
}
