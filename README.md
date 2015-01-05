**ResourceD** collects resources data inside host machine and perform the following:

* Serves the data as HTTP+JSON.

* Sends the data using custom programs.

ResourceD is currently alpha software. Use it at your own risk.


## Installation

Precompiled binary for darwin and linux will be provided in the future.


## Collecting and Running ResourceD

**1. Configuring readers and writers**

The quickest way to configure a reader is to use dynamic language.

1. Write your script. There is only 1 requirement to your script: **You must output the data in JSON format through STDOUT**.

2. Write ResourceD config file. See examples [here](https://github.com/resourced/resourced/tree/master/tests/data/config-reader).


**2. Running ResourceD**

This is an example on how to run ResourceD as foreground process.

```bash
RESOURCED_CONFIG_READER_DIR=$GOPATH/src/github.com/resourced/resourced/tests/data/config-reader \
go run $GOPATH/src/github.com/resourced/resourced/resourced.go
```

ResourceD accepts a few environment variables as configuration:

* `RESOURCED_ADDR`

* `RESOURCED_CERT_FILE`

* `RESOURCED_KEY_FILE`

* `RESOURCED_CONFIG_READER_DIR`

* `RESOURCED_CONFIG_WRITER_DIR`


## Reading Data from HTTP Interface

You can GET your resource data by sending HTTP request to the `Path = /your-resource` defined in your config-reader TOML.

For example, if you run ResourceD server as described above, you should be able to GET load average data via cURL:

```bash
curl -X GET -H "Content-type: application/json" http://localhost:55555/load-avg
```


## Third Party Data Source

The following is list of 3rd party data source that ResourceD readers use.

Big thanks to these authors, without whom ResourceD would not be possible.

* https://github.com/cloudfoundry/gosigar

* https://github.com/shirou/gopsutil

* https://github.com/guillermo/go.procmeminfo