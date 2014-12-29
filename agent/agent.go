package agent

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/julienschmidt/httprouter"
	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/libstring"
	"net/http"
	"os"
	"time"
)

func NewAgent() (*Agent, error) {
	agent := &Agent{}

	err := agent.setConfigStorage()
	if err != nil {
		return nil, err
	}

	err = agent.setDb()
	if err != nil {
		return nil, err
	}

	return agent, err
}

type Agent struct {
	ConfigStorage *resourced_config.ConfigStorage
	DbPath        string
	Db            *bolt.DB
}

func (a *Agent) setConfigStorage() error {
	readerDir := os.Getenv("RESOURCED_CONFIG_READER_DIR")
	writerDir := os.Getenv("RESOURCED_CONFIG_WRITER_DIR")

	configStorage, err := resourced_config.NewConfigStorage(readerDir, writerDir)
	if err == nil {
		a.ConfigStorage = configStorage
	}
	return err
}

func (a *Agent) setDb() error {
	var err error

	dbPath := os.Getenv("RESOURCED_DB")
	if dbPath == "" {
		dbPath = "~/resourced/db"

		err = os.MkdirAll(libstring.ExpandTilde("~/resourced"), 0755)
		if err != nil {
			return err
		}
	}

	a.DbPath = libstring.ExpandTilde(dbPath)

	a.Db, err = bolt.Open(a.DbPath, 0644, nil)
	if err != nil {
		return err
	}

	// Create "resources" bucket
	a.Db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucket([]byte("resources"))
		return nil
	})

	return err
}

func (a *Agent) DbBucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte("resources"))
}

func (a *Agent) HttpRouter() *httprouter.Router {
	router := httprouter.New()

	for _, config := range a.ConfigStorage.Readers {
		path := config.Path
		router.GET(path, func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			jsonData, err := a.GetRunByPath(config.Path)
			w.Header().Set("Content-Type", "application/json")

			if err == nil && jsonData != nil {
				w.WriteHeader(200)
				w.Write(jsonData)
			} else {
				w.WriteHeader(404)
				w.Write([]byte(fmt.Sprintf(`{"Error": "Run data does not exist.", "Path": "%v"}`, config.Path)))
			}
		})
	}

	return router
}

func (a *Agent) Run(config resourced_config.Config) ([]byte, error) {
	output, err := config.Run()
	if err != nil {
		return nil, err
	}

	err = a.SaveRun(config, output)

	return output, err
}

func (a *Agent) SaveRun(config resourced_config.Config, output []byte) error {
	resourcedData := make(map[string]interface{})
	resourcedData["UnixNano"] = time.Now().UnixNano()
	resourcedData["Command"] = config.Command
	resourcedData["Path"] = config.Path
	resourcedData["Interval"] = config.Interval

	runData := make(map[string]interface{})
	err := json.Unmarshal(output, &runData)
	if err != nil {
		return err
	}

	record := make(map[string]interface{})
	record["ResourceD"] = resourcedData
	record["Data"] = runData

	recordInJson, err := json.Marshal(record)
	if err != nil {
		return err
	}

	err = a.Db.Update(func(tx *bolt.Tx) error {
		return a.DbBucket(tx).Put([]byte(config.Path), recordInJson)
	})

	return err
}

func (a *Agent) GetRun(config resourced_config.Config) ([]byte, error) {
	return a.GetRunByPath(config.Path)
}

func (a *Agent) GetRunByPath(path string) ([]byte, error) {
	var data []byte

	a.Db.View(func(tx *bolt.Tx) error {
		data = a.DbBucket(tx).Get([]byte(path))
		return nil
	})

	return data, nil
}
