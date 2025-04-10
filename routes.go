package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	example := http.HandlerFunc(exampleAPI)
	mux.Handle("GET /api", applyRateLimit(example))

	// Refill the bucket every 2 seconds
	go func() {
		for {
			refillBucket()
			time.Sleep(2 * time.Second)
		}
	}()

	log.Println(http.ListenAndServe(":8090", mux))

}

func exampleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Req-IP", "10.10.10.1")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is responding ..."))
}
