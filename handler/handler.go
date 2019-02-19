package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/adlrocha/goxyq/config"
	"github.com/adlrocha/goxyq/log"
	"github.com/adlrocha/goxyq/queue"
	"github.com/gorilla/mux"
)

// ProxyRequest main proxy handler. All requests are handled with specific prefix
// handled by this function
func ProxyRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Infof("Received GET packet")
		bypass(w, r, config.GetConfig().DestinationHost+r.URL.Path, "GET", nil)
	} else if r.Method == "POST" {
		log.Infof("Received POST packet")
		processPost(w, r, config.GetConfig().DestinationHost+r.URL.Path)
	} else {
		respondError(w, http.StatusMethodNotAllowed, "Method not supported by goxyq yet")
	}
}

// AliveFunction Dummy function to check if service alive.AliveFunction
// We are building a proxy so it makes sense to check this.
func AliveFunction(w http.ResponseWriter, r *http.Request) {
	log.Infof("Checking if alive...")
	q := make(map[string]string)
	q["alive"] = "ok"
	respondJSON(w, http.StatusOK, q)
}

// GetQueue Gets status of a queue
func GetQueue(w http.ResponseWriter, r *http.Request) {
	// Get queueId from request
	vars := mux.Vars(r)
	queueID := vars["queueID"]

	// Get queue
	var pool = queue.NewPool()
	storedQ, err := queue.GetQueue(pool, queueID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Queue does not exist")
		// respondError(w, http.StatusBadRequest, err.Error())
	} else {
		respondJSON(w, http.StatusOK, storedQ)
	}
}

// EmptyQueue Gets status of a queue
func EmptyQueue(w http.ResponseWriter, r *http.Request) {
	// Get queueId from request
	vars := mux.Vars(r)
	queueID := vars["queueID"]
	log.Infof("Emptying queue %v\n", queueID)

	// Get queue
	var pool = queue.NewPool()
	res, err := queue.EmptyQueue(pool, queueID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Queue could not be emptied")
		// respondError(w, http.StatusBadRequest, err.Error())
	} else {
		q := make(map[string]string)
		if res {
			q["result"] = "true"
		} else {
			q["result"] = "false"
		}
		respondJSON(w, http.StatusOK, q)
	}
}

// Bypass - the proxy just bypasses the request.
func bypass(w http.ResponseWriter, r *http.Request, url string, method string, body map[string]interface{}) {
	resPayload := make([]byte, 0)

	log.Infof("Bypassing request")
	if method == "GET" {
		resPayload = makeGetRequest(url, r.Header)
	} else if method == "POST" {
		resPayload = makePostRequest(url, body, r.Header)
	} else {
		log.Errorf("Method not supported for bypass")
	}
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
	if jobBody[queueAttribute] == nil {
		// If the queueAttribute is not found in the packet bypass it.
		bypassCode = 2
		return bypassCode
	}

	// Get queue name
	qName := jobBody[queueAttribute].(string)
	// Check if the queue already exists.
	storedQ, _ := queue.GetQueue(pool, qName)
	if storedQ == nil {
		// The queue does not exist and has to be created.

		queue.NewQueue(pool, qName)
	}
	// Add job to the queue
	bodyBytes, err := json.Marshal(jobBody)
	if err != nil {
		log.Errorf("[HANDLER] Error while converting body to bytes to store in REDIS")
		return 0
	}

	res, err := queue.CreateJob(pool, qName, bodyBytes)
	if err != nil {
		return 0
	}
	if res {
		bypassCode = 1
	} else {
		bypassCode = 0
	}
	return bypassCode
}

// Wait, run job and update Queue
func waitForJobTurn(jobBody map[string]interface{}, qName string) (res bool) {
	// Create REDIS pool.
	var pool = queue.NewPool()
	bodyBytes, _ := json.Marshal(jobBody)
	res, err := queue.WaitAndRunJob(pool, qName, bodyBytes)
	if err != nil {
		log.Errorf("[HANDLER] Error while waiting for job to run...")
		return false
	}
	return res
}

func processPost(w http.ResponseWriter, r *http.Request, url string) {
	log.Infof("Processing POST request")
	// Decode received POST
	recBody := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&recBody); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	bypassCode := manageNewJob(recBody, config.GetConfig().QueueAtrribute)
	// Once job managed, get the queue name for the Body for further processing
	if bypassCode == 1 {
		qName := recBody[config.GetConfig().QueueAtrribute].(string)
		log.Infof("[HANDLER] Job handled successfully and assigned to a queue")
		success := waitForJobTurn(recBody, qName)
		if success {
			log.Infof("[HANDLER] Waited and ready to send the request. Is the jobs turn...")
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
		bypass(w, r, url, "POST", recBody)

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
		log.Errorf("Seems like proxy destination is dead!")
		panic(err)
	}
	defer resp.Body.Close()

	recBody, _ = ioutil.ReadAll(resp.Body)
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
		log.Errorf("Seems like proxy destination is dead!")
		panic(err)
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)
	return body
}
