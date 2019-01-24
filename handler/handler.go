package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/adlrocha/goxyq/config"
	"github.com/adlrocha/goxyq/log"
)

// ProxyRequest main proxy handler. All requests are handled with specific prefix
// handled by this function
func ProxyRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		bypass(w, r, config.GetConfig().DestinationHost+r.URL.Path)
	} else if r.Method == "POST" {
		processPost(w, r, config.GetConfig().DestinationHost+r.URL.Path)
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

// Auxiliary function to check equality between byte[]
func bytesEqual(a []byte, b []byte) (res bool) {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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
	resPayload := makeGetRequest(url, r.Header)
	// Use interface to dynamically get different response JSON structures.
	q := make(map[string]interface{})
	json.Unmarshal(resPayload, &q)
	if len(q) == 0 {
		l := make([]interface{}, 0)
		json.Unmarshal(resPayload, &l)
		respondJSON(w, http.StatusOK, l)
	} else {
		respondJSON(w, http.StatusOK, q)
	}
}

func processPost(w http.ResponseWriter, r *http.Request, url string) {
	recBody := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&recBody); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	resPayload := makePostRequest(url, recBody, r.Header)
	// Use interface to dynamically get different response JSON structures.
	q := make(map[string]interface{})
	json.Unmarshal(resPayload, &q)
	if len(q) == 0 {
		l := make([]interface{}, 0)
		json.Unmarshal(resPayload, &l)
		respondJSON(w, http.StatusOK, l)
	} else {
		respondJSON(w, http.StatusOK, q)
	}
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

// Bypass selected headers.
func addBypassHeader(r *http.Request, header http.Header) {
	for _, headerVar := range config.GetConfig().HeaderBypass {
		r.Header.Set(headerVar, header.Get(headerVar))
	}
}

// Make post request to the proxy destination.
func makePostRequest(url string, body map[string]interface{}, header http.Header) (recBody []byte) {
	bodyBytes, err := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	// Bypass the header of the received packet
	addBypassHeader(req, header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Debugf("response Headers: %v", resp.Header)
	recBody, _ = ioutil.ReadAll(resp.Body)
	log.Debugf("response Body: %v", string(recBody))
	return recBody
}

// Make get request to the proxy destination.
func makeGetRequest(url string, header http.Header) (body []byte) {
	req, err := http.NewRequest("GET", url, nil)
	// Bypass the header of the received packet
	addBypassHeader(req, header)
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
