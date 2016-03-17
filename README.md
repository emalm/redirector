redirector
==========

[![GoDoc](https://godoc.org/github.com/ematpl/redirector?status.svg)](https://godoc.org/github.com/ematpl/redirector)

Push to Cloud Foundry:

```
make
cf push redirector -p bin/linux -b binary_buildpack -m 64M -k 64M --no-route -c './redirector -listenAddr=0.0.0.0:$PORT'
cf map-route redirector cfapps.io -n em-go

cf push redirector -p bin/linux
```
