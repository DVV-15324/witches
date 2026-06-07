package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func HttpServerInit() {

	http.HandleFunc("/init-1", HttpInit1)
	http.HandleFunc("/init-2", HttpInit2)
	http.HandleFunc("/init-3", HttpInit3)
	http.HandleFunc("/init-4", HttpInit4)
	http.HandleFunc("/init-5", HttpInit5)

	http.ListenAndServe(":8080", nil)
}

func HttpInit1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

type HttpProfile struct {
	Name string `json:"name"`
}

func HttpInit2(w http.ResponseWriter, r *http.Request) {
	p := &HttpProfile{
		Name: "Vu",
	}
	b, _ := json.Marshal(p)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
func HttpInit3(w http.ResponseWriter, r *http.Request) {
	p := &HttpProfile{}

	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, _ := json.Marshal(p)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func HttpInit4(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	p := &HttpProfile{
		Name: name,
	}
	b, _ := json.Marshal(p)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func HttpInit5(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	name := parts[len(parts)-1]
	p := &HttpProfile{
		Name: name,
	}
	json.NewDecoder(r.Body).Decode(&p)
	b, _ := json.Marshal(p)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
