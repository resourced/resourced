[![GoDoc](https://godoc.org/github.com/resourced/resourced?status.svg)](http://godoc.org/github.com/resourced/resourced)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/resourced/resourced/master/LICENSE.md)

**ResourceD:** A single binary daemon that collects and report your server data with ease.


## Installation

1. Download the binary release [here](https://github.com/resourced/resourced/releases).

2. Use supervisor/upstart/systemd to daemonize. [Click here for examples](https://github.com/resourced/resourced/tree/master/tests/data/script-init).


## Running the Server

ResourceD requires only 2 environment variables to run.
```bash
RESOURCED_CONFIG_DIR=$GOPATH/src/github.com/resourced/resourced/tests/data/resourced-configs \
$GOPATH/bin/resourced
```

Once you executed the command above, open this URL: [http://localhost:55555/paths](http://localhost:55555/paths).
```bash
curl -X GET -H "Content-type: application/json" http://localhost:55555/r/load-avg
```


## Configuration

ResourceD accepts a few environment variables as configuration:

* **RESOURCED_CONFIG_DIR:** Path to parent config directory (Required), its children must be: readers, writers, and executors. Default: ""

* **RESOURCED_LOG_LEVEL:** Log level. Default: "info"

* **RESOURCED_ADDR:** The HTTP server host and port. Default: ":55555"

* **RESOURCED_CERT_FILE:** Path to cert file. Default: ""

* **RESOURCED_KEY_FILE:** Path to key file. Default: ""

* **RESOURCED_TAGS:** Comma separated tags. Default: ""

* **RESOURCED_ALLOWED_NETWORKS:** Comma separated list of CIDR's that are allowed to use ResourceD's API endpoint. Default: ""


## Data Gathering

ResourceD readers gather data on your server. The easiest way to create a reader is to use a scripting language.

1. Write the script following this one requirement: **Output the JSON data through STDOUT**

2. Write config file. [Click here for examples](https://github.com/resourced/resourced/tree/master/tests/data/config-reader).

For more info, [follow this link](https://github.com/resourced/resourced/tree/master/docs/users/DATA-GATHERER.md).


### RESTful Endpoints

* **GET** `/` Displays full JSON data of all readers and writers.

* **GET** `/paths` Displays paths to all readers and writers data.

* **GET** `/r` Displays full JSON data of all readers.

* **GET** `/r/paths` Displays paths to all readers data.

* **GET** `/w` Displays full JSON data of all writers.

* **GET** `/w/paths` Displays paths to all writers data.


## Third Party Data Source

Here are list of 3rd party data source that ResourceD use.
Big thanks to these authors, without whom this project would not be possible.

* https://github.com/cloudfoundry/gosigar

* https://github.com/shirou/gopsutil

* https://github.com/c9s/goprocinfo

* https://github.com/guillermo/go.procmeminfo


## Contributors

Are you a contributor, or looking to be one? [Go here!](https://github.com/resourced/resourced/tree/master/docs/contributors/README.md)
