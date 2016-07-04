package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/resourced/resourced/libhttp"
)

func (a *Agent) DefaultLogrusFieldsForHTTP() logrus.Fields {
	return logrus.Fields{
		"Addr":                a.GeneralConfig.Addr,
		"ResourcedMaster.URL": a.GeneralConfig.ResourcedMaster.URL,
	}
}

// AuthorizeMiddleware wraps all other handlers; returns 403 for clients that aren't authorized to connect.
func (a *Agent) AuthorizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Immediately forward request if there's no AccessTokens.
		if a.AccessTokens == nil || len(a.AccessTokens) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")

		if auth == "" {
			libhttp.BasicAuthUnauthorized(w, nil)
			return
		}

		accessTokenString, _, ok := libhttp.ParseBasicAuth(auth)
		if !ok {
			libhttp.BasicAuthUnauthorized(w, nil)
			return
		}

		if !a.IsAllowed(accessTokenString) {
			err := errors.New("You are not authorized to connect")
			libhttp.BasicAuthUnauthorized(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// HeadRootHandler returns empty body. This is useful for curl -I.
func (a *Agent) HeadRootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(200)
		w.Write([]byte(""))
	})
}

// GetPathsHandler returns paths to all data.
func (a *Agent) GetPathsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		payload := make(map[string][]string)
		payload["Readers"] = a.allReaderPaths()
		payload["Writers"] = a.allWriterPaths()
		payload["Executors"] = a.allExecutorPaths()
		payload["Loggers"] = a.allLogPaths()

		payloadBytes, err := json.Marshal(payload)

		if (len(payload["Readers"]) > 0 || len(payload["Writers"]) > 0 || len(payload["Executors"]) > 0 || len(payload["Loggers"]) > 0) && err == nil {
			w.WriteHeader(200)
			w.Write(payloadBytes)
		} else if err != nil {
			w.WriteHeader(503)
			w.Write([]byte(fmt.Sprintf(`{"Error": "%v"}`, err)))
		} else {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "There are no readers and writers data at all."}`)))
		}
	})
}

func (a *Agent) allReadersJsonString() string {
	readerJsonBytes := make([]string, 0)
	arrayOfReaderJsonString := "[]"

	for _, config := range a.Configs.Readers {
		jsonData, err := a.GetRunByPath(config.PathWithPrefix())
		if err == nil && jsonData != nil {
			readerJsonBytes = append(readerJsonBytes, string(jsonData))
		}
	}
	if len(readerJsonBytes) > 0 {
		arrayOfReaderJsonString = "[" + strings.Join(readerJsonBytes, ",") + "]"
	}

	return arrayOfReaderJsonString
}

// GetReadersHandler returns all readers data stored in memory.
func (a *Agent) GetReadersHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		arrayOfReaderJsonString := a.allReadersJsonString()

		if arrayOfReaderJsonString == "[]" {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist."}`)))

		} else {
			w.WriteHeader(200)
			w.Write([]byte(arrayOfReaderJsonString))
		}
	})
}

func (a *Agent) allWritersJsonString() string {
	writerJsonBytes := make([]string, 0)
	arrayOfWriterJsonString := "[]"

	for _, config := range a.Configs.Writers {
		jsonData, err := a.GetRunByPath(config.PathWithPrefix())
		if err == nil && jsonData != nil {
			writerJsonBytes = append(writerJsonBytes, string(jsonData))
		}
	}
	if len(writerJsonBytes) > 0 {
		arrayOfWriterJsonString = "[" + strings.Join(writerJsonBytes, ",") + "]"
	}

	return arrayOfWriterJsonString
}

// GetWritersHandler returns all writers data stored in memory.
func (a *Agent) GetWritersHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		arrayOfWriterJsonString := a.allWritersJsonString()

		if arrayOfWriterJsonString == "[]" {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist."}`)))

		} else {
			w.WriteHeader(200)
			w.Write([]byte(arrayOfWriterJsonString))
		}
	})
}

func (a *Agent) allExecutorsJsonString() string {
	executorJsonBytes := make([]string, 0)
	arrayOfExecutorJsonString := "[]"

	for _, config := range a.Configs.Executors {
		jsonData, err := a.GetRunByPath(config.PathWithPrefix())
		if err == nil && jsonData != nil {
			executorJsonBytes = append(executorJsonBytes, string(jsonData))
		}
	}
	if len(executorJsonBytes) > 0 {
		arrayOfExecutorJsonString = "[" + strings.Join(executorJsonBytes, ",") + "]"
	}

	return arrayOfExecutorJsonString
}

// GetExecutorsHandler returns all executors data stored in memory.
func (a *Agent) GetExecutorsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		arrayOfExecutorJsonString := a.allExecutorsJsonString()

		if arrayOfExecutorJsonString == "[]" {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist."}`)))

		} else {
			w.WriteHeader(200)
			w.Write([]byte(arrayOfExecutorJsonString))
		}
	})
}

func (a *Agent) allLogsJsonString() string {
	logJsonBytes := make([]string, 0)
	arrayOfLogJsonString := "[]"

	for _, config := range a.Configs.Loggers {
		jsonData, err := a.GetRunByPath(config.PathWithPrefix())
		if err == nil && jsonData != nil {
			logJsonBytes = append(logJsonBytes, string(jsonData))
		}
	}
	if len(logJsonBytes) > 0 {
		arrayOfLogJsonString = "[" + strings.Join(logJsonBytes, ",") + "]"
	}

	return arrayOfLogJsonString
}

// GetLogsHandler returns all logs data stored in memory.
func (a *Agent) GetLogsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		arrayOfExecutorJsonString := a.allLogsJsonString()

		if arrayOfExecutorJsonString == "[]" {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist."}`)))

		} else {
			w.WriteHeader(200)
			w.Write([]byte(arrayOfExecutorJsonString))
		}
	})
}

// GetDataTypePathsHandler returns path to all data by type.
func (a *Agent) GetDataTypePathsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		payload := make([]string, 0)

		for path, _ := range a.ResultDB.Items() {
			if strings.HasPrefix(path, "/"+vars["dataType"]+"/") {
				payload = append(payload, path)
			}
		}

		payloadBytes, err := json.Marshal(payload)

		if len(payload) > 0 && err == nil {
			w.WriteHeader(200)
			w.Write(payloadBytes)
		} else if err != nil {
			w.WriteHeader(503)
			w.Write([]byte(fmt.Sprintf(`{"Error": "%v"}`, err)))
		} else {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "There are no readers data at all."}`)))
		}
	})
}

// all /r/readers
func (a *Agent) allReaderPaths() []string {
	payload := make([]string, 0)

	for path, _ := range a.ResultDB.Items() {
		if strings.HasPrefix(path, "/r/") {
			payload = append(payload, path)
		}
	}

	return payload
}

func (a *Agent) allWriterPaths() []string {
	payload := make([]string, 0)

	for path, _ := range a.ResultDB.Items() {
		if strings.HasPrefix(path, "/w/") {
			payload = append(payload, path)
		}
	}

	return payload
}

func (a *Agent) allExecutorPaths() []string {
	payload := make([]string, 0)

	for path, _ := range a.ResultDB.Items() {
		if strings.HasPrefix(path, "/x/") {
			payload = append(payload, path)
		}
	}

	return payload
}

func (a *Agent) allLogPaths() []string {
	payload := make([]string, len(a.Configs.Loggers)+1)

	for i, config := range a.Configs.Loggers {
		if config.Path != "" {
			payload[i] = config.PathWithPrefix()
		}
	}

	payload[len(a.Configs.Loggers)] = "/logs/tcp"

	return payload
}

// GetLogPathsHandler returns all log paths.
func (a *Agent) GetLogPathsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		payload := a.allLogPaths()

		payloadBytes, err := json.Marshal(payload)

		if len(payload) > 0 && err == nil {
			w.WriteHeader(200)
			w.Write(payloadBytes)
		} else if err != nil {
			w.WriteHeader(503)
			w.Write([]byte(fmt.Sprintf(`{"Error": "%v"}`, err)))
		} else {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "There are no writers data at all."}`)))
		}
	})
}

// GetByPathHandler returns data by a specific path.
func (a *Agent) GetByPathHandler(dataType string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		path := "/" + dataType + "/" + vars["path"]

		jsonData, err := a.GetRunByPath(path)

		if err == nil && jsonData != nil {
			w.WriteHeader(200)
			w.Write(jsonData)
		} else if err != nil {
			w.WriteHeader(503)
			w.Write([]byte(fmt.Sprintf(`{"Error": "%v"}`, err)))
		} else {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist.", "Path": "%v"}`, path)))
		}
	})
}

// GetLogsTCPHandler renders logs data in JSON.
func (a *Agent) GetLogsTCPHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data := a.LogPayload(a.TCPLogDB, "")
		dataInBytes, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`{"Error": "%v"}`, err.Error())))
		} else {
			w.WriteHeader(200)
			w.Write(dataInBytes)
		}
	})
}

// HttpRouter returns HTTP router.
func (a *Agent) HttpRouter() *mux.Router {
	router := mux.NewRouter()

	router.Handle("/", a.AuthorizeMiddleware(a.HeadRootHandler())).Methods("HEAD")
	router.Handle("/", a.AuthorizeMiddleware(a.GetPathsHandler())).Methods("GET")

	router.Handle("/paths", a.AuthorizeMiddleware(a.GetPathsHandler())).Methods("GET")

	router.Handle("/r", a.AuthorizeMiddleware(a.GetReadersHandler())).Methods("GET")

	router.Handle("/w", a.AuthorizeMiddleware(a.GetWritersHandler())).Methods("GET")

	router.Handle("/x", a.AuthorizeMiddleware(a.GetExecutorsHandler())).Methods("GET")

	router.Handle("/logs", a.AuthorizeMiddleware(a.GetLogsHandler())).Methods("GET")
	router.Handle("/logs/paths", a.AuthorizeMiddleware(a.GetLogPathsHandler())).Methods("GET")
	router.Handle("/logs/tcp", a.AuthorizeMiddleware(a.GetLogsTCPHandler())).Methods("GET")

	router.Handle("/{dataType}/paths", a.AuthorizeMiddleware(a.GetDataTypePathsHandler())).Methods("GET")

	router.Handle("/r/{path}", a.AuthorizeMiddleware(a.GetByPathHandler("r"))).Methods("GET")
	router.Handle("/w/{path}", a.AuthorizeMiddleware(a.GetByPathHandler("w"))).Methods("GET")
	router.Handle("/x/{path}", a.AuthorizeMiddleware(a.GetByPathHandler("x"))).Methods("GET")
	router.Handle("/logs/{path}", a.AuthorizeMiddleware(a.GetByPathHandler("logs"))).Methods("GET")

	return router
}
