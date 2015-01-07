[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/resourced/resourced/master/LICENSE.md)

**ResourceD** collects resources data inside host machine and perform the following:

* Serves the data as HTTP+JSON.

* Sends the data using custom programs.

ResourceD is currently alpha software. Use it at your own risk.


## Installation

Precompiled binary for darwin and linux will be provided in the future.


## Collecting Data and Running ResourceD

**1. Configuring readers and writers**

ResourceD data collector is called a reader. The quickest way to configure a reader is to use dynamic language.

1. Write your script. There is only one requirement to your script: **You must output the data(in JSON) through STDOUT**.

2. Write ResourceD config file. See examples [here](https://github.com/resourced/resourced/tree/master/tests/data/config-reader).


**2. Running ResourceD**

Below is an example on how to run ResourceD as foreground process.

```bash
RESOURCED_CONFIG_READER_DIR=$GOPATH/src/github.com/resourced/resourced/tests/data/config-reader \
RESOURCED_CONFIG_WRITER_DIR=$GOPATH/src/github.com/resourced/resourced/tests/data/config-writer \
go run $GOPATH/src/github.com/resourced/resourced/resourced.go
```

ResourceD accepts a few environment variables as configuration:

* `RESOURCED_ADDR`

* `RESOURCED_CERT_FILE`

* `RESOURCED_KEY_FILE`

* `RESOURCED_CONFIG_READER_DIR`

* `RESOURCED_CONFIG_WRITER_DIR`

* `RESOURCED_TAGS`


## Reading Data from HTTP Interface

You can read resource data by sending HTTP request to the `Path = /your-resource` defined in your config-reader TOML.

For example, if you run ResourceD server as described above, you should be able to GET load average data via cURL:

```bash
curl -X GET -H "Content-type: application/json" http://localhost:55555/load-avg
```


## Third Party Data Source

The following is list of 3rd party data source that ResourceD readers use.

Big thanks to these authors, without whom ResourceD would not be possible.

* https://github.com/cloudfoundry/gosigar

* https://github.com/shirou/gopsutil

* https://github.com/c9s/goprocinfo

* https://github.com/guillermo/go.procmeminfo


## Contributors

Are you a contributor, or looking to be one? [Go here!](https://github.com/resourced/resourced/tree/master/docs/contributors/README.md)