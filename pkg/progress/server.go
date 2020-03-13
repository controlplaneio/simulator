package progress

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// HTTPHandler processes HTTP requests received on the SSH reverse tunnel
// recording the users progress
type HTTPHandler struct{}

func writeOkResponse(rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	if _, err := io.WriteString(rw, "{}"); err != nil {
		log.Println("Error writing body for GET request")
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

		// TODO: Validate the scenario name
		// TODO: Actually retrieve the scenario progress and send it back
		writeOkResponse(rw)
		return
	}

	if req.Method == "POST" {
		log.Println("Got HTTP POST Request")
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("Error reading POST body: %-v\n", err)
			http.Error(rw, "Error recording progress", http.StatusInternalServerError)
		}

		log.Println(string(body))
		// TODO: validate the POST body
		// TODO: store/update the progress
		writeOkResponse(rw)
		return
	}

	http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
}
