## ResourceD

It collects resources data inside host machine and perform the following:

* Serves the data as HTTP+JSON.

* Sends the data using custom programs.

ResourceD is currently alpha software. Use it at your own risk.


## Running as foreground process

ResourceD is currently lacking proper init script. To run as foreground process:

```bash
RESOURCED_CONFIG_READER_DIR=$GOPATH/src/github.com/resourced/resourced/tests/data/config-reader \
go run $GOPATH/src/github.com/resourced/resourced/resourced.go
```


## Fetching resource data through HTTP+JSON

You can GET your resource data by sending HTTP request to the `Path = /your-resource` defined in your config-reader TOML.

Example, by running your server as described above, you should be able to GET load average data via cURL:

```bash
curl -X GET -H "Content-type: application/json" -H "Accept: application/json" \
http://localhost:55555/load-avg
```
