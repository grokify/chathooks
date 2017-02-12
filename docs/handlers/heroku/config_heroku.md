Adding Heroku Notifications to Glip
===================================

1. Use the Heroku CLI app webhooks add on with the following command:

```bash
$ heroku addons:create deployhooks:http 
    --url=http://example.org
Adding deployhooks:http to myapp...Done.
```

See the [Heroku webhook docs](https://devcenter.heroku.com/articles/deploy-hooks#http-post-hook) for more info.
