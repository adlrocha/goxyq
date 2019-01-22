package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/adlrocha/goxyq/config"
)

// ProxyRequest main proxy handler. All requests are handled with specific prefix
// handled by this function
func ProxyRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		bypass(w, r, config.GetConfig().DestinationHost+r.URL.Path)
	} else if r.Method == "POST" {
		fmt.Fprintf(w, "Work in progress!!") // send data to client side
	} else {
		respondError(w, http.StatusMethodNotAllowed, "Method not supported by goxyq")
	}
	// r.ParseForm()       // parse arguments, you have to call this by yourself
	// fmt.Println(r.Form) // print form information in server side
	// fmt.Println("path", r.URL.Path)
	// fmt.Println("scheme", r.URL.Scheme)
	// fmt.Println(r.Form["url_long"])
	// for k, v := range r.Form {
	// 	fmt.Println("key:", k)
	// 	fmt.Println("val:", strings.Join(v, ""))
	// }

}

// AliveFunction Dummy function to check if service alive.AliveFunction
// We are building a proxy so it makes sense to check this.
func AliveFunction(w http.ResponseWriter, r *http.Request) {
	q := make(map[string]string)
	q["alive"] = "ok"
	respondJSON(w, http.StatusOK, q)
}

// Bypass - the proxy just bypasses the request.
func bypass(w http.ResponseWriter, r *http.Request, url string) {
	// path := r.URL.Path
	body := makeGetRequest(url)
	// Use interface to dynamically get different response JSON structures.
	q := make(map[string]interface{})
	json.Unmarshal(body, &q)
	respondJSON(w, http.StatusOK, q)
}

// respondJSON makes the response with payload as json format
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError makes the error response with payload as json format
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

func makePostRequest() {
	url := "https://localhost:9000/api/v1/token/create"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	q := make(map[string]string)
	err = json.Unmarshal(body, &q)
	// respondJSON(w, http.StatusOK, q)
}

func makeGetRequest(url string) (body []byte) {
	// url := "https://localhost:9000/api/v1/token/"
	req, err := http.NewRequest("GET", url, nil)
	// TODO: This should be generalized.
	// TODO: Get API KEY in config file.
	apiKey := os.Getenv("APIKEY")
	req.Header.Set("X-API-KEY", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	// q := make(map[string]string)
	// err = json.Unmarshal(body, &q)
	return body
}
