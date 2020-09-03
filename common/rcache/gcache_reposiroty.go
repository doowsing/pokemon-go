package rcache

import (
	"fmt"
	"github.com/dgraph-io/ristretto"
	"strconv"
	"time"
)

var gCache *GCache

type GCache struct {
	c *ristretto.Cache
}

func (g *GCache) init() {
	var err error
	g.c, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
		Cost: func(value interface{}) int64 {
			cost := 1
			switch value.(type) {
			case int, int32, int64, float64, bool:
				cost = 1
				break
			case string:
				cost = len(value.(string))
			case []byte:
				cost = len(value.([]byte))
			}
			return int64(cost)
		},
	})
	if err != nil {
		panic(err)
	}
}

func (g *GCache) getCache() *ristretto.Cache {
	return g.c
}

func (g *GCache) Get(key interface{}) (interface{}, bool) {
	return g.c.Get(key)
}

func (g *GCache) Set(key interface{}, value interface{}) bool {
	return g.c.Set(key, value, 0)
}

func (g *GCache) SetWithTTL(key interface{}, value interface{}, expire time.Duration) bool {
	return g.c.SetWithTTL(key, value, 0, expire)
}

func (g *GCache) Del(key interface{}) {
	g.c.Del(key)
}

func (g *GCache) Close() {
	g.c.Close()
}

func (g *GCache) Clear() {
	g.c.Clear()
}

func init() {
	gCache = &GCache{}
	gCache.init()
}

func GetGCache() *GCache {
	return gCache
}

type TimeManager struct {
	key string
	t   int // second
}

func NewTimeManager(key string, t int) *TimeManager {
	return &TimeManager{key: key, t: t}
}

func (m *TimeManager) Get(userId int) int {
	t, ok := gCache.Get(m.key + strconv.Itoa(userId))
	if ok {
		return t.(int)
	}
	return 0
}

func (m *TimeManager) Set(userId, unixTime int) {
	gCache.Set(m.key+strconv.Itoa(userId), unixTime)
}

func (m *TimeManager) Del(userId int) {
	gCache.Del(m.key + strconv.Itoa(userId))
}

// 是否还在冷却时间
func (m *TimeManager) InCoolTime(userId, unixTime int) bool {
	t := m.Get(userId)
	if t+m.t > unixTime {
		return true
	}
	return false
}

func (m *TimeManager) NowUnix() int {
	return int(time.Now().Unix())
}

func main() {

	// set a value with a cost of 1
	gCache.Set("key", "value")

	// wait for value to pass through buffers
	time.Sleep(10 * time.Millisecond)

	value, found := gCache.Get("key")
	if !found {
		panic("missing value")
	}
	fmt.Println(value)
	gCache.Del("key")
}
