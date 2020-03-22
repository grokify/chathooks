# Chathooks Developer Guide

Chathooks is a webhook to chat proxy framework that converts outbound cloud app webhooks to inbound team chat webhooks.

## Creating a New Handler

1. Select one of the handlers in the `src/handlers` folder and create a duplicate, for example in th dirctory `src/handlers/myapp`.
1. Add examples `docs/handlers/myapp`.
1. Add examples filenames to `src/util/example_events.go`.
1. Add icon to `docs/icons` folder. Add icon reference to `src/config/icons/go`.
1. Add to `examples/local_send/local_send.go`.
1. Run `examples/local_send/local_send.go` to test handler.