package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Rate limiter with Token bucket algorithm

var rateLimit = 5
var bucket []int64
var mu sync.Mutex    // Mutex to avoid race conditions on the bucket

func refillBucket() {
	mu.Lock() // lock the bucket to avoid concurrent modification
	defer mu.Unlock()

	token := time.Now().Unix()
	if len(bucket) < rateLimit {
		bucket = append(bucket, token)
	}

	fmt.Println(len(bucket))
}

func applyRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// refill the bucket before limiting the reqs
		// refillBucket()

		// lock the bucket before accessing it
		mu.Lock()
		defer mu.Unlock()

		// every req consumes a single token from the bucket
		if len(bucket) > 0 {
			bucket = bucket[1:]
			bucket_length := strconv.Itoa(len(bucket))
			w.Header().Set("X-RateLimit-Remained", bucket_length)
			next.ServeHTTP(w, r)
		}

		// no more token in the bucket
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Too many requests"))
		return 

	})
}
