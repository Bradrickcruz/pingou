// Alvo local para o teste de memória: HTTP em :9000 (/ok, /slow, /fail).
// O mesmo listener serve de alvo TCP.
package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ok/", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })
	mux.HandleFunc("/slow/", func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(2 * time.Second)
		w.Write([]byte("slow"))
	})
	mux.HandleFunc("/fail/", func(w http.ResponseWriter, _ *http.Request) { http.Error(w, "boom", 500) })
	log.Fatal(http.ListenAndServe(":9000", mux))
}
