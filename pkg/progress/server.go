package progress

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

// HTTPHandler processes HTTP requests received on the SSH reverse tunnel
// recording the users progress
type HTTPHandler struct {
	StateProvider StateProvider
	Logger        *logrus.Logger
}

// NewHTTPHandler constructs a new HTTPHandler instance
func NewHTTPHandler(sp StateProvider, logger *logrus.Logger) HTTPHandler {
	return HTTPHandler{
		StateProvider: sp,
		Logger:        logger,
	}
}

func (hh HTTPHandler) writeOkResponse(rw http.ResponseWriter, sp *ScenarioProgress) {
	hh.Logger.Info("Got HTTP request")
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(sp)
	if err != nil {
		http.Error(rw, "Error marshaling progress to json",
			http.StatusInternalServerError)
		return
	}

	if _, err := io.WriteString(rw, string(bytes)); err != nil {
		http.Error(rw, "Error writing body", http.StatusInternalServerError)
		hh.Logger.Println("Error writing body for GET request")
		return
	}
}

// ServeHTTP handles HTTP requests to record the progress a user has made on
// a scenario
func (hh HTTPHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		scenario := req.URL.Query().Get("scenario")
		if scenario == "" {
			hh.Logger.Warn("Missing scenario name")
			http.Error(rw, "Missing scenario name", http.StatusBadRequest)
			return
		}

		progress, err := hh.StateProvider.GetProgress(scenario)
		if err != nil {
			hh.Logger.WithFields(logrus.Fields{
				"Error":    err,
				"Scenario": scenario,
			}).Error("Error getting scenario progress")
			http.Error(rw, "Error getting progress", http.StatusInternalServerError)
			return
		}

		if progress == nil {
			hh.Logger.WithFields(logrus.Fields{
				"Scenario": scenario,
			}).Info("No scenario progress found")
			http.NotFound(rw, req)
			return
		}

		hh.writeOkResponse(rw, progress)
		return
	}

	if req.Method == "POST" {
		hh.Logger.Println("Got HTTP POST Request")
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			hh.Logger.WithFields(logrus.Fields{
				"Error": err,
			}).Error("Error reading POST body")
			http.Error(rw, "Error recording progress", http.StatusInternalServerError)
			return
		}

		hh.Logger.Println("Got POST of scenario progress")
		hh.Logger.Println(string(body))

		var progress ScenarioProgress
		if err := json.Unmarshal(body, &progress); err != nil {
			hh.Logger.WithFields(logrus.Fields{
				"Error": err,
			}).Error("Error unmarshaling POST body")
			http.Error(rw, "Malformed POST body", http.StatusBadRequest)
			return
		}

		hh.writeOkResponse(rw, &progress)
		return
	}

	hh.Logger.WithFields(logrus.Fields{
		"Method": req.Method,
	}).Warn("Disallowed HTTP method")
	http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
}
