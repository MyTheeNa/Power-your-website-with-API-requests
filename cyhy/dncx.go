package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	url         = "https://cyhy.xyz/serverice/ip"
	bodyData    = "ip=122.155.17.103"
	numRequests = 1000
	concurrency = 100
)

var (
	successCount int
	failCount    int
	mu           sync.Mutex
)

func makeRequest(client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyData)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		incrementFailCount()
		return
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "th-TH,th;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Sec-CH-UA", "\"Not)A;Brand\";v=\"99\", \"Google Chrome\";v=\"127\", \"Chromium\";v=\"127\"")
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Referer", "https://cyhy.xyz/serverice/ip")
	req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		incrementFailCount()
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		incrementFailCount()
		return
	}

	if resp.StatusCode == http.StatusOK {
		incrementSuccessCount()
	} else {
		incrementFailCount()
	}
}

func incrementSuccessCount() {
	mu.Lock()
	successCount++
	mu.Unlock()
}

func incrementFailCount() {
	mu.Lock()
	failCount++
	mu.Unlock()
}

func drawUI() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(0, 0, 'R', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(1, 0, 'e', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(2, 0, 'q', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(3, 0, 'u', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(4, 0, 'e', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(5, 0, 's', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(6, 0, 't', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(7, 0, 's', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(8, 0, ':', termbox.ColorGreen, termbox.ColorDefault)

	termbox.SetCell(10, 0, 'S', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(11, 0, 'u', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(12, 0, 'c', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(13, 0, 'c', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(14, 0, 'e', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(15, 0, 's', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(16, 0, 's', termbox.ColorGreen, termbox.ColorDefault)
	termbox.SetCell(17, 0, ':', termbox.ColorGreen, termbox.ColorDefault)

	termbox.SetCell(19, 0, 'F', termbox.ColorRed, termbox.ColorDefault)
	termbox.SetCell(20, 0, 'a', termbox.ColorRed, termbox.ColorDefault)
	termbox.SetCell(21, 0, 'i', termbox.ColorRed, termbox.ColorDefault)
	termbox.SetCell(22, 0, 'l', termbox.ColorRed, termbox.ColorDefault)
	termbox.SetCell(23, 0, 's', termbox.ColorRed, termbox.ColorDefault)
	termbox.SetCell(24, 0, ':', termbox.ColorRed, termbox.ColorDefault)

	termbox.Flush()
}

func updateUI() {
	for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		drawUI()
		termbox.SetCell(18, 0, rune(successCount+'0'), termbox.ColorGreen, termbox.ColorDefault)
		termbox.SetCell(25, 0, rune(failCount+'0'), termbox.ColorRed, termbox.ColorDefault)
		termbox.Flush()
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	var wg sync.WaitGroup
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	semaphore := make(chan struct{}, concurrency)

	go updateUI()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		semaphore <- struct{}{}
		go func() {
			defer func() { <-semaphore }()
			makeRequest(client, &wg)
		}()
	}

	wg.Wait()
	close(semaphore)
}
