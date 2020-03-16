package progress

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// HTTPHandler processes HTTP requests received on the SSH reverse tunnel
// recording the users progress
type HTTPHandler struct {
	StateProvider StateProvider
}

// NewHTTPHandler constructs a new HTTPHandler instance
func NewHTTPHandler(sp StateProvider) HTTPHandler {
	return HTTPHandler{
		StateProvider: sp,
	}
}

func writeOkResponse(rw http.ResponseWriter, sp *ScenarioProgress) {
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(sp)
	if err != nil {
		http.Error(rw, "Error marshaling progress to json", http.StatusInternalServerError)
		return
	}

	if _, err := io.WriteString(rw, string(bytes)); err != nil {
		http.Error(rw, "Error writing body", http.StatusInternalServerError)
		log.Println("Error writing body for GET request")
		return
	}
}

// ServeHTTP handles HTTP requests to record the progress a user has made on
// a scenario
func (hh HTTPHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		scenario := req.URL.Query().Get("scenario")
		if scenario == "" {
			http.Error(rw, "Missing scenario name", http.StatusBadRequest)
			return
		}

		progress, err := hh.StateProvider.GetProgress(scenario)
		if err != nil {
			http.Error(rw, "Error getting progress", http.StatusInternalServerError)
			return
		}

		if progress == nil {
			http.NotFound(rw, req)
			return
		}

		writeOkResponse(rw, progress)
		return
	}

	if req.Method == "POST" {
		log.Println("Got HTTP POST Request")
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("Error reading POST body: %-v\n", err)
			http.Error(rw, "Error recording progress", http.StatusInternalServerError)
			return
		}

		log.Println("Got POST of scenario progress")
		log.Println(string(body))

		var progress ScenarioProgress
		if err := json.Unmarshal(body, &progress); err != nil {
			log.Printf("Error unmarshaling POST body: %-v\n", err)
			http.Error(rw, "Malformed POST body", http.StatusBadRequest)
			return
		}

		writeOkResponse(rw, &progress)
		return
	}

	http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
}
