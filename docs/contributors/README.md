This section is dedicated for current or future ResourceD contributors.

## FAQ

**Q: What are the prerequisites to build ResourceD?**

1. Go programming language with version > 1.2.

2. [Vagrant](https://www.vagrantup.com/), to build Linux binary.


**Q: How to run the daemon?**

Below is an example on how to run ResourceD as foreground process.
```bash
RESOURCED_CONFIG_DIR=$GOPATH/src/github.com/resourced/resourced/tests/resourced-configs \
go run $GOPATH/src/github.com/resourced/resourced/resourced.go
```


**Q: How to run tests?**

There are a few ways to run tests:

1. It is best to run tests inside Vagrant VM, because the VM installs all dependencies:
    ```bash
    cd $GOPATH/src/github.com/resourced/resourced
    vagrant up ubuntu       # or vagrant up centos
    vagrant ssh ubuntu      # or vagrant ssh centos

    sudo su -

    # resourced code is located here
    cd /vagrant

    # test agent
    go test ./agent

    # test readers
    go test ./readers

    # test writers
    go test ./writers

    # You may wonder, why can't I just do: go test ./... on root project?
    It's because of /vendor folder debacle: https://github.com/golang/go/issues/11659
    ```

2. Some tests are runnable without dependencies, thus you can run them on laptop:
    ```bash
    cd $GOPATH/src/github.com/resourced/resourced/agent; go test

    # any tests inside lib packages are without dependencies
    cd $GOPATH/src/github.com/resourced/resourced/libstring; go test
    ```

3. To test and see if resourced is buildable inside docker:
    ```bash
    sudo docker build -t resourced . && sudo docker run -t resourced
    ```


**Q: What is the coding style?**

Please use `go fmt` everywhere. If you use SublimeText, feel free to install `GoSublime` and `GoOracle`.


**Q: What is the general architecture?**

ResourceD has 4 components: Reader, Writer, Executor, and Logger. They all run in their own goroutines.

* Reader scrapes information on your server and returns JSON data.

* Writer reads the JSON data, process it further, and sends it to other places.

* Executor executes logic based on expression performed on readers data.

* Logger tails a log file and forwards the log lines to master.


ResourceD also runs multiple network listeners:

* Metrics receiver (TCP or UDP): They are for receiving live Graphite or StatsD metrics.

* Log receiver (TCP or UDP): They are for receiving log lines, think of them as logs forwarder to master.

