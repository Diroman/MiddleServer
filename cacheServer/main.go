package main

import (
	memorycache2 "cacheServer/go-memorycache-example"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var cache = memorycache2.New()

func cached(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	path := r.URL.Path
	isCSS := strings.Index(r.URL.Path, ".css") != -1
	isMainPage := r.URL.Path == "/"

	// TODO: return save fragment info file
	isVideoInfoFile := (strings.Index(r.URL.Path, ".ts") == -1) &&
		(strings.Index(r.URL.Path, "watch") != -1)

	if isMainPage || isVideoInfoFile {
		respBody, err := getRequest(path)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.Write(respBody)
	}

	value, ok := cache.Get(path)
	if !ok {
		respBody, err := getRequest(path)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		value = respBody
		var storageTime time.Duration
		if isCSS {
			//storageTime = 60 * time.Minute
			storageTime = 5 * time.Second
		} else {
			//storageTime = 10 * time.Minute
			storageTime = 5 * time.Second
		}

		cache.Set(path, value)
	}

	if isCSS {
		w.Header().Set("Content-Type", "text/css")
	}

	w.Write(value)
}

func getRequest(path string) ([]byte, error){
	log.Println("Send request to server: " + path)
	resp, err := http.Get("http://sports-stream.pw" + path)

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return b, nil
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(cached)

	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:8080",
	}

	log.Fatal(srv.ListenAndServe())
}
