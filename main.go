package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"stats.io/pkg/persistency"
	"stats.io/pkg/stats"
)

var (
	StatHandler *stats.Stats
)

func main() {
	dbReader := persistency.NewPersistent()
	data := dbReader.LoadFileToMem()
	log.Println("Loaded data length:", len(*data))
	StatHandler = stats.NewStats(data)
	go func(p *persistency.Pers, s *stats.Stats) {
		ticker := time.Tick(100 * time.Millisecond)
		for range ticker {
			p.DumpToFile(s)
		}
	}(dbReader, StatHandler)
	go StatHandler.PeriodicCleanup()
	http.HandleFunc("/", handleRequest)
	log.Println("Starting server on 8080")
	http.ListenAndServe(":8080", nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	count := StatHandler.AddEntry()
	w.WriteHeader(http.StatusOK)

	// Write the response body
	response := fmt.Sprintf("%d", count)
	w.Write([]byte(response))
}
