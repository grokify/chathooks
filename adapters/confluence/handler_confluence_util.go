package confluence

import (
	"github.com/grokify/glip-go-webhook"
)

func ExamplePageCreatedMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExamplePageCreatedMessageSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return Normalize(msg), nil
}

func ExamplePageCreatedMessageSource() (ConfluenceOutMessage, error) {
	return ConfluenceOutMessageFromBytes(ExamplePageCreatedMessageBytes())
}

func ExamplePageCreatedMessageBytes() []byte {
	return []byte(`{
   "page": {
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
   "user": "admin",
   "userKey": "ff80808154510724015451074c160001",
   "timestamp": 1471926079645,
   "username": "admin"
}`)
}

func ExampleCommentCreatedMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExampleCommentCreatedMessageSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return Normalize(msg), nil
}

func ExampleCommentCreatedMessageSource() (ConfluenceOutMessage, error) {
	return ConfluenceOutMessageFromBytes(ExampleCommentCreatedMessageBytes())
}

func ExampleCommentCreatedMessageBytes() []byte {
	return []byte(`{
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
}

/*

Confluence page_created

Activity: msg.Page.CreatorName created page in space [msg.pPge.SpaceKey]()
Body []()

{
   "page": {
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
   "user": "admin",
   "userKey": "ff80808154510724015451074c160001",
   "timestamp": 1471926079645,
   "username": "admin"
 }

Confluence comment_created

{
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
 }

*/
