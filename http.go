package GeeCache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const basicpath = "_geecache"

var httpPoolMade = false

type HttpPool struct {
	self     string
	basepath string
}

func NewHttpPool(self string) *HttpPool {
	if httpPoolMade {
		log.Println("groupcache:NewHttpPool must be called only once")
	}
	return &HttpPool{
		self:     self,
		basepath: basicpath,
	}
}
func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}
func (httpPool *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, httpPool.basepath) {
		panic("HttpPool servering unexecpt path " + r.URL.Path)
	}
	httpPool.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	strs := strings.SplitN(r.URL.Path[len(httpPool.basepath):], "/", 2)
	if len(strs) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := strs[0]
	key := strs[1]
	group, ok := GetGroup(groupName)
	if !ok {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, "no such key: "+key, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

var db_test = map[string]string{
	"tom":  "630",
	"jack": "390",
	"ice":  "439",
}

func main() {
	NewGroup("test", 2<<10,
		GetterFunc(func(key string) ([]byte, error) {
			log.Printf("[SlowDB] search key %v", key)
			if v, ok := db_test[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	addr := "localhost:9999"
	peers := NewHttpPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
