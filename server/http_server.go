package server

import (
	"fmt"
	"gogocache/cache"
	"gogocache/client"
	"gogocache/hash"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	defaultPath     = "/gogocache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	self        string
	basePath    string
	mu          sync.Mutex
	peers       *hash.Map
	httpGetters map[string]*client.HttpGetter
}

func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers = hash.NewHash(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*client.HttpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &client.HttpGetter{BaseURL: peer + p.basePath}
	}
}

func (p *HTTPPool) PickPeer(key string) (client.PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{self: self, basePath: defaultPath}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	h.Log("%s %s", r.Method, r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(h.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]
	group := cache.GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
