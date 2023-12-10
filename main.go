package main

import (
	"flag"
	"fmt"
	"gogocache/cache"
	"gogocache/server"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func startCacheServer(addr string, addrs []string, cache *cache.Group) {
	peers := server.NewHTTPPool(addr)
	peers.Set(addrs...)
	cache.RegisterPeers(peers)
	log.Println("GoGoCache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *cache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

// func main() {

// 	r := gin.Default()
// 	r.GET("/ping", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"message": "pong",
// 		})
// 	})

// 	r.POST("/startCache", createCache)

// 	r.Run() // listen and serve on 0.0.0.0:8080
// }

var addrMap = map[int]string{}

func createCache(c *gin.Context) {
	var request struct {
		Port   []int `json:"port"`
		Server int   `json:"server"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// for _, v := range request.Port {
	// 	addrMap[v] = "http://localhost:" + strconv.Itoa(v)
	// }
	// var addrs []string
	// for _, v := range addrMap {
	// 	addrs = append(addrs, v)
	// }
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	apiAddr := "http://localhost:9999"

	for _, p := range request.Port {
		cache := createGroup()
		if p == request.Server {
			go startAPIServer(apiAddr, cache)
		}
		go startCacheServer(addrMap[p], []string(addrs), cache)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cache creation successful", "ports": addrs})
}

func createGroup() *cache.Group {
	return cache.NewGroup("scores", 2<<10, cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func main() {

	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Go server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	apiAddr := "http://localhost:9999"

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], []string(addrs), gee)
}
