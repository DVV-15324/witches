package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"time"
)

func main() {
	// 1. GET đơn giản
	getSimple()
	// 2. GET + query params
	getWithQuery()
	// 3. Path param
	getWithPathParam()
	// 4. POST JSON
	postJSON()
	// 5. GET + header
	getWithHeader()
	// 6. Full production template
	fullTemplate()
}

func getSimple() {
	resp, err := http.Get("https://api.example.com/users")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("GET simple:", string(body))
}

func getWithQuery() {

	baseURL := "https://api.example.com/users"

	req, _ := http.NewRequest("GET", baseURL, nil)

	q := req.URL.Query()
	q.Add("name", "vudinh")
	q.Add("age", "20")

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("GET query:", string(body))
}

func getWithPathParam() {

	id := "123"
	url := fmt.Sprintf("https://api.example.com/users/%s", id)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Path param:", string(body))
}

func postJSON() {

	data := map[string]any{
		"name": "vudinh",
		"age":  20,
	}

	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest(
		"POST",
		"https://api.example.com/users",
		bytes.NewBuffer(jsonData),
	)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("POST JSON:", string(body))
}

func getWithHeader() {

	req, _ := http.NewRequest("GET", "https://api.example.com/users", nil)

	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("X-API-Key", "abc123")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Header:", string(body))
}

func fullTemplate() {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, _ := http.NewRequest("GET", "https://api.example.com/users", nil)

	q := req.URL.Query()
	q.Add("name", "vudinh")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Full template:", string(body))
}
