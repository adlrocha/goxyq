package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/adlrocha/goxyq/config"
	"github.com/adlrocha/goxyq/log"
	"github.com/adlrocha/goxyq/queue"
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

// Manage the arriving job
func manageNewJob(jobBody map[string]interface{}, queueAttribute string) (bypassCode uint8) {
	// Create REDIS pool.
	var pool = queue.NewPool()
	// Check if the queue already exists.
	storedQ, _ := queue.GetQueue(pool, queueAttribute)
	if storedQ == nil {
		// The queue does not exist and has to be created.
		if jobBody[queueAttribute] == nil {
			// If the queueAttribute is not found in the packet bypass it.
			bypassCode = 2
			return bypassCode
		}
		qName := jobBody[queueAttribute].(string)
		queue.NewQueue(pool, qName)
	}
	// Add job to the queue
	bodyBytes, err := json.Marshal(jobBody)
	if err != nil {
		log.Errorf("[HANDLER] Error while converting body to bytes to store in REDIS")
		return 0
	}
	// TODO: TODO: TODO: Refactor so its pretier
	res, err := queue.CreateJob(pool, queueAttribute, bodyBytes)
	if err != nil {
		return 0
	}
	if res {
		bypassCode = 1
	} else {
		bypassCode = 0
	}
	// TODO: TODO: TODO: End of refactor.
	return bypassCode
}

// Wait, run job and update Queue
func waitForJobTurn(jobBody map[string]interface{}, queueAttribute string) (res bool) {
	// Create REDIS pool.
	var pool = queue.NewPool()
	bodyBytes, _ := json.Marshal(jobBody)
	res, err := queue.WaitAndRunJob(pool, config.GetConfig().QueueAtrribute, bodyBytes)
	if err != nil {
		log.Errorf("[HANDLER] Error while waiting for job to run...")
		return false
	}
	return res
}

func processPost(w http.ResponseWriter, r *http.Request, url string) {
	// Decode received POST
	recBody := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&recBody); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	// TODO: Manage the job
	bypassCode := manageNewJob(recBody, config.GetConfig().QueueAtrribute)
	if bypassCode == 1 {
		log.Debugf("[HANDLER] Job handled successfully and assigned to a queue")
		success := waitForJobTurn(recBody, config.GetConfig().QueueAtrribute)
		if success {
			log.Debugf("[HANDLER] Waited and ready to send the request. Is the jobs turn...")
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
		} else {
			log.Errorf("[HANDLER] Error waiting for job to run")
		}
	} else if bypassCode == 2 {
		log.Debugf("[HANDLER] QueueAttribute not find while managing job. Bypass request")
		bypass(w, r, url)

	} else {
		log.Errorf("[HANDLER] Error managing new job")
	}

	// The queue is empty. Make post request

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
