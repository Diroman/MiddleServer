package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/djherbis/times.v1"
)

var storageTime = 10 * time.Minute

func getCreateTime(path string) time.Time {
	t, err := times.Stat("../cacheServer/cache/" + path)
	if err != nil {
		log.Fatal(err.Error())
	}

	return t.ModTime()
}

func main() {
	ticker := time.NewTicker(5*time.Second)

	for {
		<-ticker.C
		files, err := ioutil.ReadDir("../cacheServer/cache")
		if err != nil {
			log.Println(err)
			continue
		}

		for _, file := range files {
			if time.Now().Add(-storageTime).After(getCreateTime(file.Name())) {
				if err := os.Remove("../cacheServer/cache/" + file.Name()); err != nil {
					log.Println(err)
				}
			}
		}
	}
}
