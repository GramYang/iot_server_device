package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var c *cache.Cache

func SetUp() {
	c = cache.New(22*time.Hour, 24*time.Hour)
	c.Flush()
}

func SaveToken(username, token string) {
	c.Set(username, token, cache.DefaultExpiration)
}

func GetToken(username string) string {
	token, found := c.Get(username)
	if found {
		return token.(string)
	} else {
		return ""
	}
}

func DeleteToken(username string) {
	c.Delete(username)
}
