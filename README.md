# goxyq
Goxyq is a proxy server written in Golang with built in queues that
may be used to sequentialize an asynchronous system or API. Goxyq uses
REDIS as its storage system for queuing information.

*Disclaimer: Goxyq's current design is quite simple and pretty ad-hoc
but with enough generality to be used in a great gamut of use cases. The
system will be maintained if people show interest for it. If not...
it already served me well in a production system*


### Usage
* To install the tool just run (be sure that your $GOPATH is set):
```
go get github.com/adlrocha/goxyq
```
or just clone the repo.

* Before running the proxy be sure that you have a REDIS instance running.
You can achieve this by running the following in the repo's path:
```
./sccripts/redis_start.sh
```
* Now you are ready to run the goxyq:
```
go run server.go
```
* You can also build the tool and run it:
```
go build
./server
```

### Configuration
To modify Goxyq's configuration go to `./config/config.go`.
* `Port`: Goxyq's listening port.
* `DestinationHost`: Destination host.
* `ProxyPathPrefix`: Prefix path for which the proxy will inspect traffic.
* `QueueAtrribute`: Body attributes that will trigger the creation of a new
queue or the stack of the request.
* `HeaderBypass`: List of headers that will be bypassed by the proxy.
### Potential short-term enhancements
* Allow the use of multiple QueueAttributes.
* Use of etcd instead of REDIS as queue storage (it better fits decentralized
environments).
