**ResourceD** collects resources data inside host machine and perform the following:

* Serves the data as HTTP+JSON.

* Sends the data using custom programs.

ResourceD is currently alpha software. Use it at your own risk.


## 1. Installation

#### Running as foreground process
```bash
RESOURCED_CONFIG_READER_DIR=$GOPATH/src/github.com/resourced/resourced/tests/data/config-reader \
go run $GOPATH/src/github.com/resourced/resourced/resourced.go
```


## 2. Usage

#### GET resource data through HTTP+JSON

You can GET your resource data by sending HTTP request to the `Path = /your-resource` defined in your config-reader TOML.

For example, if you run ResourceD server as described above, you should be able to GET load average data via cURL:

```bash
curl -X GET -H "Content-type: application/json" http://localhost:55555/load-avg
```


## 3rd party data source

The following is list of 3rd party data source that ResourceD readers use.

Big thanks to these authors, without whom ResourceD would not be possible.

* https://github.com/cloudfoundry/gosigar

* https://github.com/shirou/gopsutil

* https://github.com/guillermo/go.procmeminfo