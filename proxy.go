package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

var (
	counter int = 0
	ports   [2]string
)

func main() {
	ports[0] = "8000"
	ports[1] = "8001"
	http.HandleFunc("/", handleProxy)
	log.Fatalln(http.ListenAndServe("localhost:8080", nil))
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	r.RequestURI = ""
	var currentUrl string
	if counter == 0 {
		currentUrl = "http://localhost:" + ports[counter]
		counter++
	} else {
		currentUrl = "http://localhost:" + ports[counter]
		counter = 0
	}

	u, err := url.Parse(currentUrl + r.URL.Path)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP:", err)
	}
	r.URL = u
	resp, err := client.Do(r)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP:", err)
		return
	}
	defer resp.Body.Close()
	w.WriteHeader(resp.StatusCode)
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(content)

}
