# CommonChat

[![Build Status][build-status-svg]][build-status-link]
[![Go Report Card][goreport-svg]][goreport-link]
[![Docs][docs-godoc-svg]][docs-godoc-link]
[![License][license-svg]][license-link]

CommonChat is an abstraction library for chat / team messaging services like Glip and Slack. It currently includes two parts:

* Common message format - After converting a message to the `commonchat.Message` format, the libraries can be used to convert to individula chat services.
* Webhook clients - Given a webhook URL, the clients use the `commonchat.Adapter` interface to enable webhook API calls using the `commonchat.Message` format.

It is currently used with the Chathooks webhook formatting service:

[https://github.com/grokify/chathooks](https://github.com/grokify/chathooks)

 [build-status-svg]: https://api.travis-ci.org/grokify/commonchat.svg?branch=master
 [build-status-link]: https://travis-ci.org/grokify/commonchat
 [coverage-status-svg]: https://coveralls.io/repos/grokify/commonchat/badge.svg?branch=master
 [coverage-status-link]: https://coveralls.io/r/grokify/commonchat?branch=master
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/commonchat
 [goreport-link]: https://goreportcard.com/report/github.com/grokify/commonchat
 [codeclimate-status-svg]: https://codeclimate.com/github/grokify/commonchat/badges/gpa.svg
 [codeclimate-status-link]: https://codeclimate.com/github/grokify/commonchat
 [docs-godoc-svg]: https://img.shields.io/badge/docs-godoc-blue.svg
 [docs-godoc-link]: https://godoc.org/github.com/grokify/commonchat
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-link]: https://github.com/grokify/commonchat/blob/master/LICENSE.md
