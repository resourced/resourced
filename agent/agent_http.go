package agent

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	resourced_config "github.com/resourced/resourced/config"
	"net/http"
	"strings"
)

// RootGetHandler returns function that handles all readers and writers.
func (a *Agent) RootGetHandler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		readerJsonBytes := make([]string, 0)
		arrayOfReaderJsonString := "[]"

		for _, config := range a.ConfigStorage.Readers {
			jsonData, err := a.GetRunByPath(a.pathWithPrefix(config))
			if err == nil && jsonData != nil {
				readerJsonBytes = append(readerJsonBytes, string(jsonData))
			}
		}
		if len(readerJsonBytes) > 0 {
			arrayOfReaderJsonString = "[" + strings.Join(readerJsonBytes, ",") + "]"
		}

		writerJsonBytes := make([]string, 0)
		arrayOfWriterJsonString := "[]"

		for _, config := range a.ConfigStorage.Writers {
			jsonData, err := a.GetRunByPath(a.pathWithPrefix(config))
			if err == nil && jsonData != nil {
				writerJsonBytes = append(writerJsonBytes, string(jsonData))
			}
		}
		if len(writerJsonBytes) > 0 {
			arrayOfWriterJsonString = "[" + strings.Join(writerJsonBytes, ",") + "]"
		}

		if arrayOfReaderJsonString == "[]" && arrayOfWriterJsonString == "[]" {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist."}`)))

		} else {
			w.WriteHeader(200)
			w.Write([]byte(fmt.Sprintf(`{"Readers": %v, "Writers": %v}`, arrayOfReaderJsonString, arrayOfWriterJsonString)))
		}
	}
}

// ReadersGetHandler returns function that handles all readers.
func (a *Agent) ReadersGetHandler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		readerJsonBytes := make([]string, 0)
		arrayOfReaderJsonString := "[]"

		for _, config := range a.ConfigStorage.Readers {
			jsonData, err := a.GetRunByPath(a.pathWithPrefix(config))
			if err == nil && jsonData != nil {
				readerJsonBytes = append(readerJsonBytes, string(jsonData))
			}
		}
		if len(readerJsonBytes) > 0 {
			arrayOfReaderJsonString = "[" + strings.Join(readerJsonBytes, ",") + "]"
		}

		if arrayOfReaderJsonString == "[]" {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist."}`)))

		} else {
			w.WriteHeader(200)
			w.Write([]byte(arrayOfReaderJsonString))
		}
	}
}

// WritersGetHandler returns function that handles all writers.
func (a *Agent) WritersGetHandler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		writerJsonBytes := make([]string, 0)
		arrayOfWriterJsonString := "[]"

		for _, config := range a.ConfigStorage.Writers {
			jsonData, err := a.GetRunByPath(a.pathWithPrefix(config))
			if err == nil && jsonData != nil {
				writerJsonBytes = append(writerJsonBytes, string(jsonData))
			}
		}
		if len(writerJsonBytes) > 0 {
			arrayOfWriterJsonString = "[" + strings.Join(writerJsonBytes, ",") + "]"
		}

		if arrayOfWriterJsonString == "[]" {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist."}`)))

		} else {
			w.WriteHeader(200)
			w.Write([]byte(arrayOfWriterJsonString))
		}
	}
}

// PathsGetHandler returns function that shows all the paths.
func (a *Agent) PathsGetHandler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		payload := make(map[string][]string)
		payload["Readers"] = make([]string, len(a.ConfigStorage.Readers))
		payload["Writers"] = make([]string, len(a.ConfigStorage.Writers))

		for i, config := range a.ConfigStorage.Readers {
			if config.Path != "" {
				payload["Readers"][i] = a.pathWithPrefix(config)
			}
		}

		for i, config := range a.ConfigStorage.Writers {
			if config.Path != "" {
				payload["Writers"][i] = a.pathWithPrefix(config)
			}
		}

		payloadBytes, err := json.Marshal(payload)

		if (len(payload["Readers"]) > 0 || len(payload["Writers"]) > 0) && err == nil {
			w.WriteHeader(200)
			w.Write(payloadBytes)
		} else if err != nil {
			w.WriteHeader(503)
			w.Write([]byte(fmt.Sprintf(`{"Error": "%v"}`, err)))
		} else {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "There are no readers and writers data at all."}`)))
		}
	}
}

// ReaderPathsGetHandler returns function that shows all the readers paths.
func (a *Agent) ReaderPathsGetHandler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		payload := make([]string, len(a.ConfigStorage.Readers))

		for i, config := range a.ConfigStorage.Readers {
			if config.Path != "" {
				payload[i] = a.pathWithPrefix(config)
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
	}
}

// WriterPathsGetHandler returns function that shows all the writers paths.
func (a *Agent) WriterPathsGetHandler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		payload := make([]string, len(a.ConfigStorage.Writers))

		for i, config := range a.ConfigStorage.Writers {
			if config.Path != "" {
				payload[i] = a.pathWithPrefix(config)
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
			w.Write([]byte(fmt.Sprintf(`{"Error": "There are no writers data at all."}`)))
		}
	}
}

// readerOrWriterGetHandler returns a function that handle reader/writer.
func (a *Agent) readerOrWriterGetHandler(path string, config resourced_config.Config) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if path != "" {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			w.Header().Set("Content-Type", "application/json")

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
		}
	}

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist.", "Path": "%v"}`, path)))
	}
}

// MapReadersGetHandlers returns functions that handle readers paths.
func (a *Agent) MapReadersGetHandlers() map[string]func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	handlersMap := make(map[string]func(w http.ResponseWriter, r *http.Request, ps httprouter.Params))

	for _, config := range a.ConfigStorage.Readers {
		if config.Path != "" {
			path := a.pathWithPrefix(config)
			handlersMap[path] = a.readerOrWriterGetHandler(path, config)
		}
	}
	return handlersMap
}

// MapWritersGetHandlers returns functions that handle writers paths.
func (a *Agent) MapWritersGetHandlers() map[string]func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	handlersMap := make(map[string]func(w http.ResponseWriter, r *http.Request, ps httprouter.Params))

	for _, config := range a.ConfigStorage.Writers {
		if config.Path != "" {
			path := a.pathWithPrefix(config)
			handlersMap[path] = a.readerOrWriterGetHandler(path, config)
		}
	}
	return handlersMap
}

// HttpRouter returns HTTP router.
func (a *Agent) HttpRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", a.RootGetHandler())
	router.GET("/paths", a.PathsGetHandler())

	router.GET("/r", a.ReadersGetHandler())
	router.GET("/r/paths", a.ReaderPathsGetHandler())

	router.GET("/w", a.WritersGetHandler())
	router.GET("/w/paths", a.WriterPathsGetHandler())

	for readerPath, readerHandler := range a.MapReadersGetHandlers() {
		router.GET(readerPath, readerHandler)
	}

	for writerPath, writerHandler := range a.MapWritersGetHandlers() {
		router.GET(writerPath, writerHandler)
	}

	return router
}

// ListenAndServe runs HTTP server.
func (a *Agent) ListenAndServe(addr string) error {
	if addr == "" {
		addr = ":55555"
	}

	router := a.HttpRouter()
	return http.ListenAndServe(addr, router)
}

// ListenAndServe runs HTTPS server.
func (a *Agent) ListenAndServeTLS(addr string, certFile string, keyFile string) error {
	if addr == "" {
		addr = ":55555"
	}

	router := a.HttpRouter()
	return http.ListenAndServeTLS(addr, certFile, keyFile, router)
}
