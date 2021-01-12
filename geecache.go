package GeeCache

import (
	"errors"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	mainCache cache
	get       Getter
}

var (
	groups = make(map[string]*Group)
	mu     sync.RWMutex
)

func NewGroup(name string, maxBufferSize int64, get Getter) *Group {
	if name == "" {
		panic("Group name cant be empty")
	}
	if get == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	group := Group{
		name: name,
		mainCache: cache{
			cacheBytes: maxBufferSize,
		},
		get: get,
	}
	groups[name] = &group
	return &group
}
func GetGroup(name string) (group *Group, ok bool) {
	if name == "" {
		return nil, false
	}
	mu.RLock()
	defer mu.RUnlock()
	group, ok = groups[name]
	return
}
func (g *Group) Get(key string) (value ByteView, err error) {
	if value, ok := g.mainCache.get(key); ok {
		return value, nil
	}
	return g.load(key)
}
func (g *Group) load(key string) (value ByteView, err error) {
	if value, err = g.remoteLoad(key); err == nil {
		return value, nil
	}
	return g.localLoad(key)
}
func (g *Group) remoteLoad(key string) (value ByteView, err error) {
	//TODO need to be done
	return value, errors.New("")
}
func (g *Group) localLoad(key string) (value ByteView, err error) {
	bytes, err := g.get.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value = ByteView{
		bytes: copyBytes(bytes),
	}
	g.populateCache(key, value)
	return
}
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
