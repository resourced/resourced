package agent

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func (a *Agent) RootGetHandler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		jsonObjects := make([]string, 0)

		for _, config := range a.ConfigStorage.Readers {
			jsonData, err := a.GetRunByPath(config.Path)
			if err == nil && jsonData != nil {
				jsonObjects = append(jsonObjects, string(jsonData))
			}
		}

		if len(jsonObjects) > 0 {
			w.WriteHeader(200)
			arrayOfJsonObjects := "[" + strings.Join(jsonObjects, ",") + "]"
			w.Write([]byte(arrayOfJsonObjects))
		} else {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist."}`)))
		}
	}
}

func (a *Agent) PathsGetHandler() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		jsonObjects := make([]string, 0)

		for _, config := range a.ConfigStorage.Readers {
			jsonObjects = append(jsonObjects, config.Path)
		}

		arrayOfJsonObjects, err := json.Marshal(jsonObjects)

		if len(jsonObjects) > 0 && err == nil {
			w.WriteHeader(200)
			w.Write([]byte(arrayOfJsonObjects))
		} else if err != nil {
			w.WriteHeader(503)
			w.Write([]byte(fmt.Sprintf(`{"Error": "%v"}`, err)))
		} else {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist."}`)))
		}
	}
}

func (a *Agent) ReadersGetHandler() map[string]func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	handlersMap := make(map[string]func(w http.ResponseWriter, r *http.Request, ps httprouter.Params))

	for _, config := range a.ConfigStorage.Readers {
		path := config.Path

		if path != "" {
			handlersMap[path] = func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	}
	return handlersMap
}

func (a *Agent) HttpRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", a.RootGetHandler())

	router.GET("/paths", a.PathsGetHandler())

	for path, handler := range a.ReadersGetHandler() {
		router.GET(path, handler)
	}

	return router
}

func (a *Agent) ListenAndServe(addr string) error {
	if addr == "" {
		addr = ":55555"
	}

	router := a.HttpRouter()
	return http.ListenAndServe(addr, router)
}

func (a *Agent) ListenAndServeTLS(addr string, certFile string, keyFile string) error {
	if addr == "" {
		addr = ":55555"
	}

	router := a.HttpRouter()
	return http.ListenAndServeTLS(addr, certFile, keyFile, router)
}
