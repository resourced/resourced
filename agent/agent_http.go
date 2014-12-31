package agent

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func (a *Agent) HttpRouter() *httprouter.Router {
	router := httprouter.New()

	// Root Path
	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	})

	// /paths Path
	router.GET("/paths", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		jsonObjects := make([]string, 0)

		for _, config := range a.ConfigStorage.Readers {
			jsonObjects = append(jsonObjects, config.Path)
		}

		arrayOfJsonObjects, err := json.Marshal(jsonObjects)

		if len(jsonObjects) > 0 && err == nil {
			w.WriteHeader(200)
			w.Write([]byte(arrayOfJsonObjects))
		} else {
			w.WriteHeader(404)
			w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist."}`)))
		}
	})

	// Readers' Path
	for _, config := range a.ConfigStorage.Readers {
		path := config.Path

		if path != "" {
			router.GET(path, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
				w.Header().Set("Content-Type", "application/json")

				jsonData, err := a.GetRunByPath(path)

				if err == nil && jsonData != nil {
					w.WriteHeader(200)
					w.Write(jsonData)
				} else {
					w.WriteHeader(404)
					w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist.", "Path": "%v"}`, path)))
				}
			})
		}
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
