package confluence

import (
	"testing"

	"github.com/grokify/glip-webhook-proxy-go/config"
)

var ConfigurationTests = []struct {
	v    int
	want string
}{
	{8080, ":8080"}}

func TestConfluence(t *testing.T) {
	for _, tt := range ConfigurationTests {
		cfg := config.Configuration{
			Port: tt.v}

		addr := cfg.Address()
		if tt.want != addr {
			t.Errorf("Configuration.Address(%v): want %v, got %v", tt.v, tt.want, addr)
		}
	}
}

func CommentCreated() []byte {
	bytes := []byte(`{
   "comment": {
     "spaceKey": "~admin",
     "parent": {
       "spaceKey": "~admin",
       "modificationDate": 1471926079631,
       "creatorKey": "ff80808154510724015451074c160001",
       "creatorName": "admin",
       "lastModifierKey": "ff80808154510724015451074c160001",
       "self": "https://cloud-development-environment.atlassian.net/wiki/display/~admin/Some+random+test+page",
       "lastModifierName": "admin",
       "id": 16777227,
       "title": "Some random test page",
       "creationDate": 1471926079631,
       "version": 1
     },
     "modificationDate": 1471926091465,
     "creatorKey": "ff80808154510724015451074c160001",
     "creatorName": "admin",
     "lastModifierKey": "ff80808154510724015451074c160001",
     "self": "https://cloud-development-environment/wiki/display/~admin/Some+random+test+page?focusedCommentId=16777228#comment-16777228",
     "lastModifierName": "admin",
     "id": 16777228,
     "creationDate": 1471926091465,
     "version": 1
   },
   "user": "admin",
   "userKey": "ff80808154510724015451074c160001",
   "timestamp": 1471926091468,
   "username": "admin"
 }`)
	return bytes
}
