# Go Service Router
Simple service router written in Go and backed with Consul.

|    |    |
|----|----|
| Build Status | [![CircleCI](https://circleci.com/gh/nicholasjackson/router.svg?style=svg)](https://circleci.com/gh/nicholasjackson/router) |
| Stories Ready | [![Stories in Ready](https://badge.waffle.io/nicholasjackson/router.svg?label=ready&title=Ready)](http://waffle.io/nicholasjackson/router) |

Service router is an implementation of a service router built in Go, it can act as a proxy internally for your microservice environment and is designed to be fault tollerant and highly efficient.

# How to Build

```bash
$ cd _build
$ ./minke build_image
$ ./minke cucumber
```

# Running
```bash
$ cd _build
$ ./minke build_image
$ ./minke run
```

## Flow
Periodically (default 30s) service router will make a call to the backend and retrieve registered services.
Calling the service router with a url `http://myapi.servicerouter/myfunction?params=1` will look up services which have the name *myapi* and call the backend found in the registry passing header, path, and querystring.  It will return the response to the client.

## Load Balancing
The router will load balance all endpoints found for an api in a round robin fashion.

## Circuit breaking / retries
The router implements a circuit breaking pattern so that if a backend is not available the circuit will be opened and future calls not attempted until the breaker times out.  On an error the router will try the next endpoint in the loadbalanced list.

## Restricting endpoint by tags
If a call specifies the header `X-API-ENDPOINT-TAG: tag_name` then the endpoints returned by the registry will be filtered by this tag.  This is useful if multiple endpoints exist on a cluster which may be of different versions, it is possible to tag these endpoints and dynamically route requests based on this tag.  This may be used as part of a Canary deployment strategy where the old and the newly deployed endpoint exists at the same time.  
If the header is not present then the router will not filter the service list and will return all backends which match the service name.

## StatsD
The router exports metrics related to all operations to a StatsD backend using UDP.  