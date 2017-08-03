# micro-github-auth-proxy
Github Auth proxy designed for micro-apps with static client and an API server



## What is this project?

## Configuration


```
  {
    "users": [
      "xyz",
      "xyz"
    ],
    "upstreams": [
      {
        "prefix": "/",
        "type": "static",
        "location": "static"
      },
      {
        "type": "server",
        "prefix": "/api",
        "location": "http://localhost:4040"
      }
    ]
  }
```

* `users` - Users that are allowed to access your app.
* `upstreams` - HTTP supported upstreams that the proxy will call

You also need a couple of ENV variables

```
  CLIENT_ID
  CLIENT_SECRET
```

These are your Github Client ID and Client secret from the Oauth page.


## Limitations

 * Currently, for `V1`, this only supports HTTP backend. So for example if you
   have a frontend and a backend, you can host your frontend on some service that
   support HTTP endpoint. eg: S3 with static hosting. Some `heroku` or anything
   that will give you a valid web address.
 * If you have a client side app calling to a server side through the proxy
   (which is what this is designed for), you will need to pass in the cookie
   for all server requests.

   If you are using `fetch`, you will need to use `{ credentials: "include" }`.
   This will make sure that the github token will get passed through to the
   proxy and it will direct the request to the server and not redirect to
   Github all over again.

## Design Decisions and architecture

### Github Auth Context Token Memoization Process

Once you have a Github Token in the http cookie, the proxy will validate that
token against Github every 2 minutes (counted when the process is running).

For example:

You went to github, authenticated and you got the token `ABC` back. The proxy
verify that this token belongs to a user that the proxy allows through.

Once that token is verified, it is stored in memory for 2 minutes to avoid
re-validating against the Github API.
