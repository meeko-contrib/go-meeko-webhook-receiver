# cider-webhook-receiver #

[![Build
Status](https://drone.io/github.com/salsita-cider/cider-webhook-receiver/status.png)](https://drone.io/github.com/salsita-cider/cider-webhook-receiver/latest) [![Coverage Status](https://coveralls.io/repos/salsita-cider/cider-webhook-receiver/badge.png)](https://coveralls.io/r/salsita-cider/cider-webhook-receiver)

Cider abstract application for implementing webhook collectors.

This package takes care of some boilerplate code while implementing webhook
collectors for Cider. All that is necessary is to pass a `http.Handler` into
`ListenAndServe` and the rest will be taken care of.

## Environment ##

There are two variables that must be defined for `ListenAndServe` to work:

* `LISTEN_ADDRESS` - the network address to listen on, format `HOST:[PORT]`
* `ACCESS_TOKEN` - the access token that must be present in all the POST request
                   as a parameter, otherwise the request is rejected

These should be included in `.cider/app.json` as other variables.

## License ##

GNU GPLv3, see the `LICENSE` file.
