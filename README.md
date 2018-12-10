hyper
=====

[![Build Status](https://travis-ci.org/vaniila/hyper.svg?branch=master)](https://travis-ci.org/vaniila/hyper)
[![GoDoc](https://godoc.org/github.com/vaniila/hyper?status.svg)](https://godoc.org/github.com/vaniila/hyper)

Package hyper implements an ease-of-use HTTP web framework for the Go
programming language. Featuring GraphQL (w/ Subscription), DataLoader
and Swagger integration.

Project Website: https://github.com/vaniila/hyper<br>
API documentation: https://godoc.org/github.com/vaniila/hyper<br>
API examples: https://github.com/vaniila/hyper/tree/master/examples

Installation
------------

    go get github.com/vaniila/hyper

Features
--------

* Built-in support for HTTP/2 protocol
* Built-in support for Opentracing
* Build GraphQL APIs
  * Support for GraphQL subscription via websocket
  * Implement custom object, argument, enum, scalar and union types
  * API to access the Dataloader interface
* Support for websocket
  * Authorization via middleware functions
  * Scales horizontally with PubSub channels
* Build RESTful with swagger integration
* Built-in request validation (Query, Body and Header)
* Support for namespaces
* Customizable middleware and HTTP error handling
* Define and use custom logger
* Identity and access management
* Configurable CORS
* Automatic crash prevention

Contributing
------------

Everyone is encouraged to help improve this project. Here are a few ways you can help:

- [Report bugs](https://github.com/vaniila/hyper/issues)
- Fix bugs and [submit pull requests](https://github.com/vaniila/hyper/pulls)
- Write, clarify, or fix documentation
- Suggest or add new features

License
-------

> Copyright (c) Vaniila, Inc. and its affiliates.
> Use of this source code is governed by a MIT license
> that can be found in the [LICENSE](./LICENSE) file.
