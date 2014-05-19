# go-meeko-webhook-receiver #

[![Build
Status](https://drone.io/github.com/meeko-contrib/go-meeko-webhook-receiver/status.png)](https://drone.io/github.com/meeko-contrib/go-meeko-webhook-receiver/latest) [![Coverage Status](https://coveralls.io/repos/meeko-contrib/go-meeko-webhook-receiver/badge.png)](https://coveralls.io/r/meeko-contrib/go-meeko-webhook-receiver)

go-meeko helper for implementing webhook collectors.

This package takes care of some boilerplate code while implementing webhook
collectors for Meeko. All that is necessary is to pass a `http.Handler` into
`ListenAndServe` and the rest will be taken care of.

## Usage ##

```go
import "github.com/meeko-contrib/go-meeko-webhook-receiver/receiver"
```

## Agent Configuration ##

There are two variables that must be defined for `ListenAndServe` to work:

* `LISTEN_ADDRESS` - the network address to listen on, format `HOST:[PORT]`
* `ACCESS_TOKEN` - the access token that must be present in all the POST request
                   as a parameter, otherwise the request is rejected

These should be included in `.meeko/agent.json` as other variables.

## License ##

MIT, see the `LICENSE` file.
