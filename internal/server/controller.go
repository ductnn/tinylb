package server

import (
	"container/list"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Controller struct {
	servers map[int]*server
	upIDs   *list.List
	downIDs *list.List
	mux     sync.Mutex
}

func NewController() *Controller {
	return &Controller{
		upIDs:   list.New(),
		downIDs: list.New(),
	}
}

func (c *Controller) SetupServers(urls ...*url.URL) {
	c.servers = make(map[int]*server, len(urls))

	for i, u := range urls {
		id := i + 1
		c.servers[id] = newServer(u, c.serverHTTPHandler(id, u))
		c.upIDs.PushBack(id)
	}
}

func (c *Controller) getNext() *server {
	id := c.getNextID()
	if id == 0 {
		return nil
	}
	return c.servers[id]
}

func (c *Controller) getNextID() int {
	defer c.mux.Unlock()
	c.mux.Lock()
	if e := c.upIDs.Front(); e != nil {
		c.upIDs.MoveToBack(e)
		return e.Value.(int)
	}
	return 0
}

func (c *Controller) HTTPHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if server := c.getNext(); server != nil {
			server.ServeHTTP(rw, req)
			return
		}
		http.Error(rw, "Service unavailable", http.StatusServiceUnavailable)
	})
}

func (c *Controller) serverHTTPHandler(id int, u *url.URL) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		c.down(id)
		c.HTTPHandler().ServeHTTP(rw, req)
	}
	return proxy
}

func (c *Controller) down(id int) {
	defer c.mux.Unlock()
	c.mux.Lock()

	c.downIDs.PushBack(id)
	for e := c.upIDs.Front(); e != nil; e = e.Next() {
		if upID := e.Value.(int); upID == id {
			c.upIDs.Remove(e)
			break
		}
	}
	log.Warnf("[%s] down", c.servers[id].url)
}

func (c *Controller) HealthCheck() []string {
	defer c.mux.Unlock()
	c.mux.Lock()

	var (
		msgs []string
		next *list.Element
	)
	for e := c.upIDs.Front(); e != nil; e = next {
		next = e.Next()
		id := e.Value.(int)
		if server := c.servers[id]; !server.IsAlive() {
			c.downIDs.PushBack(id)
			c.upIDs.Remove(e)
			log.Warnf("[%s] down", server.url)
		}
	}
	for e := c.downIDs.Front(); e != nil; e = next {
		next = e.Next()
		id := e.Value.(int)
		if server := c.servers[id]; server.IsAlive() {
			c.upIDs.PushBack(id)
			c.downIDs.Remove(e)
			log.Warnf("[%s] up", server.url)
		}
	}
	return msgs
}
