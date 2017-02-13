Adding Heroku Notifications to Glip
===================================

1. Use the Heroku command line app to add a deployment webhook to your app with the following command, replacing the example `url` with your webhook URL.

```bash
$ heroku addons:create deployhooks:http \
    --url=http://example.org
Adding deployhooks:http to myapp...Done.
```

See the [Heroku webhook docs](https://devcenter.heroku.com/articles/deploy-hooks#http-post-hook) for more info.
