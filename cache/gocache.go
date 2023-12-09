package cache

import (
	"fmt"
	pb "gogocache/cachepb"
	"gogocache/client"
	"gogocache/singleflight"
	"log"
	"sync"
)

type Group struct {
	name   string
	getter Getter
	cache  Cache
	peers  client.PeerPicker

	loader *singleflight.Group
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, capacity int64, getter Getter) *Group {
	if getter == nil {
		panic("Getter Required")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		cache:  *newCache(capacity),
		getter: getter,
		loader: &singleflight.Group{},
	}
	groups[name] = g
	return g
}

func NewGroupRetGroups(name string, capacity int64, getter Getter) *map[string]*Group {
	if getter == nil {
		panic("Getter Required")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		cache:  *newCache(capacity),
		getter: getter,
	}
	groups[name] = g
	return &groups
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (group *Group) Get(key string) (ByteView, error) {

	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if val, flag := group.cache.Get(key); flag {
		return val, nil
	}

	return group.Load(key)
}

func (g *Group) RegisterPeers(peers client.PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (group *Group) Load(key string) (ByteView, error) {
	view, err := group.loader.Do(key, func() (interface{}, error) {
		if group.peers != nil {
			if peer, ok := group.peers.PickPeer(key); ok {
				var err error
				if value, err := group.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		return group.getFromLocal(key)
	})
	if err == nil {
		return view.(ByteView), nil
	}
	return ByteView{}, err
}

// func (g *Group) getFromPeer(peer client.PeerGetter, key string) (ByteView, error) {
// 	bytes, err := peer.Get(g.name, key)
// 	if err != nil {
// 		return ByteView{}, err
// 	}
// 	return ByteView{data: bytes}, nil
// }

func (g *Group) getFromPeer(peer client.PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{data: res.Value}, nil
}

func (group *Group) getFromLocal(key string) (ByteView, error) {
	value, err := group.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	data := ByteView{data: cloneBytes(value)}
	group.populateCache(key, data)
	return data, nil
}

func (group *Group) populateCache(key string, value ByteView) {
	group.cache.Add(key, value)
}

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
