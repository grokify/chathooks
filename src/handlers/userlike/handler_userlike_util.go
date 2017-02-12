package userlike

import (
	"github.com/grokify/glip-go-webhook"
)

func ExampleOfflineMessageReceiveMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExampleOfflineMessageReceiveSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return NormalizeOfflineMessageReceive(msg), nil
}

func ExampleOfflineMessageReceiveSource() (UserlikeOfflineMessageReceiveOutMessage, error) {
	return UserlikeOfflineMessageReceiveOutMessageFromBytes(ExampleOfflineMessageReceiveBytes())
}

func ExampleOfflineMessageReceiveBytes() []byte {
	return []byte(`{
 "_event": "receive",
 "_type": "offline_message",
 "browser_name": "Chrome",
 "browser_os": "Mac OS X",
 "browser_version": "32",
 "chat_widget": {
   "id": 2,
   "name": "Website"
 },
 "client_email": "support@userlike.com",
 "client_name": "Userlike Support",
 "created_at": "2014-12-20 14:50:23",
 "custom": {},
 "data_privacy": null,
 "id": 3,
 "loc_city": "Cologne",
 "loc_country": "Germany",
 "loc_lat": 50.9333000183105,
 "loc_lon": 6.94999980926514,
 "marked_read": true,
 "message": "We are happy to welcome you as a Userlike user!",
 "page_impressions": 5,
 "screenshot_oid": null,
 "screenshot_url": null,
 "status": "new",
 "topic": "Support",
 "url": "http://www.userlike.com",
 "visits": 1
}
`)
}

func ExampleUserlikeChatMetaStartOutMessageGlip() (glipwebhook.GlipWebhookMessage, error) {
	msg, err := ExampleUserlikeChatMetaStartOutMessageSource()
	if err != nil {
		return glipwebhook.GlipWebhookMessage{}, err
	}
	return NormalizeChatMetaStart(msg), nil
}

func ExampleUserlikeChatMetaStartOutMessageSource() (UserlikeChatMetaStartOutMessage, error) {
	return UserlikeChatMetaStartOutMessageFromBytes(ExampleUserlikeChatMetaStartOutMessageBytes())
}

func ExampleUserlikeChatMetaStartOutMessageBytes() []byte {
	return []byte(`{
 "_event": "start",
 "_type": "chat_meta",
 "browser_name": "Safari",
 "browser_os": "Mac OS X",
 "browser_version": "8",
 "chat_widget": {
 "id": 9,
 "name": "Testing David"
 },
 "chat_widget_goal": {
 "id": null,
 "name": null
 },
 "client_additional01_name": null,
 "client_additional01_value": null,
 "client_additional02_name": null,
 "client_additional02_value": null,
 "client_additional03_name": null,
 "client_additional03_value": null,
 "client_email": "david@optixx.org",
 "client_name": "Jo",
 "client_uuid": "nEitxHDooFzsPhMoVn0QCw8E.L3mmogIcp+FGwXos5a4NUdZ9/uQbSCBx0wDIRVFWM+o",
 "created_at": "2014-12-29 11:26:24",
 "custom": {
    "basket": {
      "item01": {
        "desc": "33X Optical Zoom Camcorder Mini DV",
        "id": "2acefe58-91e5-11e1-beba-000c2979313a",
        "price": 139.99,
        "url": "http://application/en/electronics/34-camcorder.html"
      },
      "item02": {
        "desc": "Home Theater System",
        "id": "31aca2f2-91e5-11e1-beba-000c2979313a",
        "long": "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
        "price": 499.99,
        "url": "https://application/en/electronics/39-home-theater.html"
      }
    },
    "id": "428614f0-91e5-11e1-beba-000c2979313a",
    "ref": "3efd5462e"
  },
 "data_privacy": false,
 "duration": "00:00:13",
 "ended_at": null,
 "feedback_message": null,
 "id": 71,
 "inital_url": "https://devel.userlike.local/en/",
 "loc_city": null,
 "loc_country": null,
 "loc_lat": null,
 "loc_lon": null,
 "marked_read": false,
 "messages": [],
  "notes": [],
  "operator_created": {
    "email": "david@userlike.com",
    "first_name": "David",
    "id": 5,
    "last_name": "Voswinkel",
    "name": "David Voswinkel",
    "operator_group": {
      "id": 14,
      "name": "Testing David"
    }
  },
  "operator_created_id": 5,
  "operator_current": {
    "email": "david@userlike.com",
    "first_name": "David",
    "id": 5,
    "last_name": "Voswinkel",
    "name": "David Voswinkel",
    "operator_group": {
      "id": 14,
      "name": "Testing David"
    }
   },
  "operator_current_id": 5,
  "page_impressions": 2,
  "post_survey_option": null,
  "rate": null,
  "referrer": null,
  "status": "new",
  "topic": null,
  "url": "https://devel.userlike.local/en/debug/9",
  "visits": 10,
  "was_proactive": false
}
`)
}
