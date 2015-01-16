[![GoDoc](https://godoc.org/github.com/resourced/resourced?status.svg)](http://godoc.org/github.com/resourced/resourced) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/resourced/resourced/master/LICENSE.md)

**ResourceD** collects resources data inside host machine and perform the following:

* Serves the data as HTTP+JSON.

* Sends the data using custom programs.


## Installation

Precompiled binary for darwin and linux will be provided in the future.


## Collecting Data and Running ResourceD

**1. Configuring readers and writers**

ResourceD data collector is called a reader. The quickest way to configure a reader is to use a script.

1. Write your script. There is only one requirement to your script: **You must output the JSON data through STDOUT**

2. Write ResourceD config file. [See examples here](https://github.com/resourced/resourced/tree/master/tests/data/config-reader).


**2. Running ResourceD**

Below is an example on how to run ResourceD as foreground process.

```bash
RESOURCED_CONFIG_READER_DIR=$GOPATH/src/github.com/resourced/resourced/tests/data/config-reader \
RESOURCED_CONFIG_WRITER_DIR=$GOPATH/src/github.com/resourced/resourced/tests/data/config-writer \
go run $GOPATH/src/github.com/resourced/resourced/resourced.go
```

ResourceD accepts a few environment variables as configuration:

* **RESOURCED_ADDR:** The HTTP server host and port. Default: ":55555"

* **RESOURCED_CERT_FILE:** Path to cert file. Default: ""

* **RESOURCED_KEY_FILE:** Path to key file. Default: ""

* **RESOURCED_CONFIG_READER_DIR:** Path to readers config directory. Default: ""

* **RESOURCED_CONFIG_WRITER_DIR:** Path to writers config directory. Default: ""

* **RESOURCED_TAGS:** Comma separated tags. Default: []


## Reading Data from HTTP Interface

If you run ResourceD server as described above, you should be able to GET load average data via cURL:

```bash
curl -X GET -H "Content-type: application/json" http://localhost:55555/r/load-avg
```

### RESTful Endpoints

* **GET** `/` Displays full JSON data of all readers and writers.

* **GET** `/paths` Displays paths to all readers and writers data.

* **GET** `/r` Displays full JSON data of all readers.

* **GET** `/r/paths` Displays paths to all readers data.

* **GET** `/w` Displays full JSON data of all writers.

* **GET** `/w/paths` Displays paths to all writers data.



## Third Party Data Source

Here are list of 3rd party data source that ResourceD use. Big thanks to these authors, without whom this project would not be possible.

* https://github.com/cloudfoundry/gosigar

* https://github.com/shirou/gopsutil

* https://github.com/c9s/goprocinfo

* https://github.com/guillermo/go.procmeminfo


## Contributors

Are you a contributor, or looking to be one? [Go here!](https://github.com/resourced/resourced/tree/master/docs/contributors/README.md)
