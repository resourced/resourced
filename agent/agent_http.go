package agent

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/libhttp"
)

// AuthorizeMiddleware wraps all other handlers; returns 403 for clients that aren't authorized to connect.
func (a *Agent) AuthorizeMiddleware(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Immediately forward request if there's no AccessTokens.
		if a.AccessTokens == nil || len(a.AccessTokens) == 0 {
			h(w, r, ps)
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

		// Forward request to given handle
		h(w, r, ps)
		return
	}
}

// RootHeadHandler returns empty payload.
func (a *Agent) RootHeadHandler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(200)
		w.Write([]byte(""))
	}
}

// RootGetHandler returns data from all readers and writers.
func (a *Agent) RootGetHandler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

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

		if arrayOfReaderJsonString == "[]" && arrayOfWriterJsonString == "[]" && arrayOfExecutorJsonString == "[]" {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist."}`)))

		} else {
			w.WriteHeader(200)
			w.Write([]byte(fmt.Sprintf(`{"Readers": %v, "Writers": %v, "Executors": %v}`, arrayOfReaderJsonString, arrayOfWriterJsonString, arrayOfExecutorJsonString)))
		}
	}
}

// ReadersGetHandler returns function that handles all readers.
func (a *Agent) ReadersGetHandler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

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

		for _, config := range a.Configs.Writers {
			jsonData, err := a.GetRunByPath(config.PathWithPrefix())
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

// ExecutorsGetHandler returns function that handles all Executors.
func (a *Agent) ExecutorsGetHandler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		writerJsonBytes := make([]string, 0)
		arrayOfWriterJsonString := "[]"

		for _, config := range a.Configs.Executors {
			jsonData, err := a.GetRunByPath(config.PathWithPrefix())
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
		payload["Readers"] = make([]string, len(a.Configs.Readers))
		payload["Writers"] = make([]string, len(a.Configs.Writers))
		payload["Executors"] = make([]string, len(a.Configs.Executors))

		for i, config := range a.Configs.Readers {
			if config.Path != "" {
				payload["Readers"][i] = config.PathWithPrefix()
			}
		}

		for i, config := range a.Configs.Writers {
			if config.Path != "" {
				payload["Writers"][i] = config.PathWithPrefix()
			}
		}

		for i, config := range a.Configs.Executors {
			if config.Path != "" {
				payload["Executors"][i] = config.PathWithPrefix()
			}
		}

		payloadBytes, err := json.Marshal(payload)

		if (len(payload["Readers"]) > 0 || len(payload["Writers"]) > 0 || len(payload["Executors"]) > 0) && err == nil {
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

		// all /r/readers + /r/graphite
		payload := make([]string, len(a.Configs.Readers)+1)

		for i, config := range a.Configs.Readers {
			if config.Path != "" {
				payload[i] = config.PathWithPrefix()
			}
		}

		payload[len(a.Configs.Readers)] = "/r/graphite"

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

		payload := make([]string, len(a.Configs.Writers))

		for i, config := range a.Configs.Writers {
			if config.Path != "" {
				payload[i] = config.PathWithPrefix()
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

// ExecutorPathsGetHandler returns function that shows all the executors paths.
func (a *Agent) ExecutorPathsGetHandler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		payload := make([]string, len(a.Configs.Executors))

		for i, config := range a.Configs.Executors {
			if config.Path != "" {
				payload[i] = config.PathWithPrefix()
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

// handlerByPath returns a function that handle reader/writer/executor.
func (a *Agent) handlerByPath(path string, config resourced_config.Config) func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	for _, config := range a.Configs.Readers {
		if config.Path != "" {
			path := config.PathWithPrefix()
			handlersMap[path] = a.handlerByPath(path, config)
		}
	}
	return handlersMap
}

// MapWritersGetHandlers returns functions that handle writers paths.
func (a *Agent) MapWritersGetHandlers() map[string]func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	handlersMap := make(map[string]func(w http.ResponseWriter, r *http.Request, ps httprouter.Params))

	for _, config := range a.Configs.Writers {
		if config.Path != "" {
			path := config.PathWithPrefix()
			handlersMap[path] = a.handlerByPath(path, config)
		}
	}
	return handlersMap
}

// MapExecutorsGetHandlers returns functions that handle executors paths.
func (a *Agent) MapExecutorsGetHandlers() map[string]func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	handlersMap := make(map[string]func(w http.ResponseWriter, r *http.Request, ps httprouter.Params))

	for _, config := range a.Configs.Executors {
		if config.Path != "" {
			path := config.PathWithPrefix()
			handlersMap[path] = a.handlerByPath(path, config)
		}
	}
	return handlersMap
}

// ReadersGraphiteGetHandler returns renders graphite readers in JSON.
func (a *Agent) ReadersGraphiteGetHandler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		dataInBytes, err := a.GraphiteDB.ToJson()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf(`{"Error": "%v"}`, err.Error())))
		} else {
			w.WriteHeader(200)
			w.Write(dataInBytes)
		}
	}
}

// HttpRouter returns HTTP router.
func (a *Agent) HttpRouter() *httprouter.Router {
	router := httprouter.New()

	router.HEAD("/", a.AuthorizeMiddleware(a.RootHeadHandler()))
	router.GET("/", a.AuthorizeMiddleware(a.RootGetHandler()))
	router.GET("/paths", a.AuthorizeMiddleware(a.PathsGetHandler()))

	router.GET("/r", a.AuthorizeMiddleware(a.ReadersGetHandler()))
	router.GET("/r/paths", a.AuthorizeMiddleware(a.ReaderPathsGetHandler()))

	router.GET("/r/graphite", a.AuthorizeMiddleware(a.ReadersGraphiteGetHandler()))

	router.GET("/w", a.AuthorizeMiddleware(a.WritersGetHandler()))
	router.GET("/w/paths", a.AuthorizeMiddleware(a.WriterPathsGetHandler()))

	router.GET("/x", a.AuthorizeMiddleware(a.ExecutorsGetHandler()))
	router.GET("/x/paths", a.AuthorizeMiddleware(a.ExecutorPathsGetHandler()))

	for path, handler := range a.MapReadersGetHandlers() {
		router.GET(path, a.AuthorizeMiddleware(handler))
	}

	for path, handler := range a.MapWritersGetHandlers() {
		router.GET(path, a.AuthorizeMiddleware(handler))
	}

	for path, handler := range a.MapExecutorsGetHandlers() {
		router.GET(path, a.AuthorizeMiddleware(handler))
	}

	return router
}
