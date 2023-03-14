package cache

import (
	"github.com/MegaMindInKZ/task-techno.git/db"
	"math"
	"strings"
	"time"
)

const TTLMinute = 5

var LocalCache Cache

type Cache interface {
	Add(key, value string)
	Get(key string) (value string, ok bool)
	Len() int
}

func SetUp() {
	rows, err := db.DB.Query(
		"SELECT ID, ACTIVE_LINK, HISTORY_LINK FROM LINKS",
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var link db.Link
		err = rows.Scan(&link.ID, &link.ActiveLink, &link.HistoryLink)
		AddHotLinkToCache(link)
	}
}

func AddHotLinkToCache(link db.Link) {
	if isHot, key, value := isHostLink(link); isHot {
		LocalCache.Add(key, value)
	}
}

func isHostLink(link db.Link) (bool, string, string) {
	if strings.Contains(link.HistoryLink, "smartfony") {
		return true, link.HistoryLink, link.ActiveLink
	} else if strings.Contains(link.ActiveLink, "smartfony") {
		return true, link.ActiveLink, link.ActiveLink
	}
	return false, "", ""
}

type CacheEntry struct {
	value      string
	expiration int64
}

type InMemoryCache struct {
	cache   map[string]CacheEntry
	maxSize int
}

func NewCache(maxSize int) Cache {
	return &InMemoryCache{
		cache:   make(map[string]CacheEntry),
		maxSize: maxSize,
	}
}

func (c *InMemoryCache) Add(key, value string) {
	expiration := time.Now().Add(time.Minute * TTLMinute).Unix()
	c.cache[key] = CacheEntry{value, expiration}
	if len(c.cache) > c.maxSize {
		var oldestKey string
		var oldestExpiration int64 = math.MaxInt64
		for key, entry := range c.cache {
			if entry.expiration < oldestExpiration {
				oldestKey = key
				oldestExpiration = entry.expiration
			}
		}
		delete(c.cache, oldestKey)
	}
}

func (c *InMemoryCache) Get(key string) (string, bool) {
	if entry, ok := c.cache[key]; ok {
		if time.Now().Unix() < entry.expiration {
			return entry.value, true
		} else {
			delete(c.cache, key)
		}
	}
	return "", false
}

func (c *InMemoryCache) Len() int {
	return len(c.cache)
}
