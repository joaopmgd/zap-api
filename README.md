# GolangRestfulApi ZAP

This API will recover the properties for ZAP ou VivaReal

### Prerequisites

This project was built in Go version go1.11.1, so it must be a version similar with that.

NOTE: dependency management was done using Dep (dependency management tool for Go)

### Running

To run this project just enter the command bellow:

Running in the terminal:
```
make run
```

Building and Creating a Container:
```
make docker-build
```

Running the Container locally:
```
make docker-run
```



GET
- To recover and see the information about the properties, the GET request will only need a HEADER parameter called "source" which should be either "zap" or "vivareal"

Header:
```
source: zap
```
ou
```
source: vivareal
```

Pagination is implemented, the default Response has 10 Properties, if there is a need to change that just run the request with the below for example:
```
localhost:8080/properties?offset=25&limit=25
```

## Running the tests

To run the tests just execute:

```
go test
```


## Running on Docker

The endpoint and the Port are external varibles setup in the Dockerfile, if change is needed it is ok, but they are required.

Thanks ;D