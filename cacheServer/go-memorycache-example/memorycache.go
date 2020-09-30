package memorycache

import (
	"io/ioutil"
	"log"
	"strings"
)

type Cache struct {}

func New() *Cache {
	return &Cache{}
}

func (c *Cache) Set(key string, value []byte) {
	path := strings.Replace(key, "/", "_", -1)
	if err := ioutil.WriteFile("cache/" + path, value, 0644); err != nil {
		log.Println(err)
		return
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	path := strings.Replace(key, "/", "_", -1)
	value, err :=ioutil.ReadFile("cache/" + path)
	if err != nil {
		return nil, false
	}

	return value, true
}