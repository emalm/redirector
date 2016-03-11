redirector
==========

Push to Cloud Foundry:

```
make
cf push redirector -p bin/linux -b binary_buildpack --no-route -c './redirector -listenAddr=0.0.0.0:$PORT'
cf map-route redirector cfapps.io -n em-go
```
