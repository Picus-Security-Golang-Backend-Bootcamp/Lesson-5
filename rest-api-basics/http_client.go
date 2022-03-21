// package main

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net"
// 	"net/http"
// 	"net/http/httptrace"
// 	"net/url"
// 	"os"
// 	"time"
// )

// func main() {
// 	//client := httpClient()
// 	//getFromApi(client)
// 	clientWithTrace()

// }

// // func httpClient() *http.Client {
// // 	client := &http.Client{
// // 		Timeout: time.Second * 10,
// // 		Transport: &http.Transport{
// // 			MaxIdleConns:       10,
// // 			IdleConnTimeout:    30 * time.Second,
// // 			DisableCompression: true,
// // 			Dial: (&net.Dialer{
// // 				Timeout: 5 * time.Second,
// // 			}).Dial,
// // 			MaxIdleConnsPerHost: 10, // istek atılan adres başına boşta bekleyebilecek connection sayısı
// // 		},
// // 	}

// // 	return client
// // }

// func getFromApi() {
// 	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer resp.Body.Close()
// 	body, _ := io.ReadAll(resp.Body)
// 	fmt.Println(string(body))
// }

// func getFromApiWithClient(client *http.Client) {

// 	req, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/posts", nil)
// 	req.Header.Add("Accept", `application/json`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)
// 	fmt.Println(string(body))
// }

// func getWithQueryString() {

// 	id := "1"
// 	params := "postId=" + url.QueryEscape(id)
// 	path := fmt.Sprintf("https://jsonplaceholder.typicode.com/comments?%s", params)

// 	client := &http.Client{
// 		Timeout: time.Second * 10,
// 		Transport: &http.Transport{
// 			MaxIdleConns:       10,
// 			IdleConnTimeout:    30 * time.Second,
// 			DisableCompression: true,
// 			Dial: (&net.Dialer{
// 				Timeout: 5 * time.Second,
// 			}).Dial,
// 			MaxIdleConnsPerHost: 10, // istek atılan adres başına boşta bekleyebilecek connection sayısı
// 		},
// 	}

// 	req, _ := http.NewRequest("GET", path, nil)

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()
// 	body, err := io.ReadAll(resp.Body)
// 	fmt.Println(string(body))
// }

// func getGithubUser() {
// 	// current user
// 	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)

// 	req.Header.Add("Accept", `application/json`)
// 	// add header for authentication
// 	req.Header.Add("Authorization", fmt.Sprintf("%s", os.Getenv("TOKEN")))

// 	// çalıştıracağınız zaman ||| GITHUB_TOKEN=aabbcc go run main.go ||| GITHUB_TOKEN değerini kendinize göre değiştirin
// }

// func clientWithPost(client *http.Client) {
// 	body := map[string]interface{}{
// 		"sku":  "PRD-0012",
// 		"name": "product name",
// 	}

// 	jsonBody, _ := json.Marshal(body)
// 	req, _ := http.NewRequest(http.MethodPost, "http://product-service.internal.myproject.com", bytes.NewBuffer(jsonBody))

// 	resp, err := client.Do(req)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer resp.Body.Close()
// }

// //https://husni.dev/how-to-reuse-http-connection-in-go/
// func clientWithTrace() {
// 	clientTrace := &httptrace.ClientTrace{
// 		GotConn: func(gci httptrace.GotConnInfo) {
// 			log.Println(gci.Reused)
// 		},
// 	}

// 	ctx := httptrace.WithClientTrace(context.Background(), clientTrace)

// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	b, _ := io.ReadAll(resp.Body)

// 	fmt.Println(len(b))
// 	resp.Body.Close()
// 	req, err = http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	resp, err = http.DefaultClient.Do(req)
// 	b, _ = io.ReadAll(resp.Body)

// 	fmt.Println(len(b))
// 	resp.Body.Close()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
