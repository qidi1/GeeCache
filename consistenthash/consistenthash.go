package consistenthash

import (
	"hash/fnv"
)

type HashFunc func(data []byte) uint64
type Map struct {
	hash      HashFunc
	nodeNames []string
}

func New(hash HashFunc) *Map {
	if hash == nil {
		fnvHash := fnv.New64()
		hash = func(data []byte) uint64 {
			fnvHash.Write(data)
			return fnvHash.Sum64()
		}
	}
	return &Map{
		hash:      hash,
		nodeNames: make([]string, 0),
	}
}

//	https://arxiv.org/abs/1406.2294
// Hash consistently chooses a hash bucket number in the range [0, numBuckets) for the given key. numBuckets must be >= 1.
func Hash(key uint64, numBuckets int) int32 {

	var b int64 = -1
	var j int64

	for j < int64(numBuckets) {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(int64(1)<<31) / float64((key>>33)+1)))
	}

	return int32(b)
}

func (m *Map) Add(names ...string) {
	for _, name := range names {
		m.nodeNames = append(m.nodeNames, name)
	}
}
func (m *Map) Get(key string) string {
	if len(m.nodeNames) == 0 {
		return ""
	}
	b := []byte(key)
	node := Hash(m.hash(b), len(m.nodeNames))
	return m.nodeNames[node]
}
