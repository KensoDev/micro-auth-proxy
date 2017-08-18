# micro-auth-proxy

Transparent auth reverse proxy designed for micro-apps with static client and an API server.

This project is built to serve **internal** projects at your company. Those can
include dashboard, internal tools and others.

It is built with ACL and Docker deployment in mind and fits a very broad use
case that can really serve any company that requires authentication for
internal tools without messing with auth for every single project.

## Overall architecture

![Auth proxy architecture](http://assets.avi.io/authproxy-diagram.png)

This is built to protext "upstreams" behind an auth backend (Currently only supports github).

Say you have 2 upstreams configured:

```
// REDACTED
  "upstreams": [
    {
      "type": "http",
      "prefix": "/",
      "location": "https://bdf682c0.ngrok.io"
    },
    {
      "type": "http",
      "prefix": "/api/",
      "location": "https://6cc1b022.ngrok.io"
    }
  ]
// REDACTED
```

Any request starting with `/api/` will be routed through to your backend server. Anything else `/` will go to the client.

Because request are no longer routed directly from the client to the server, you will not have to define CORS or any sort of absolute URL on the client. You simply call `/api` if the client is behind the proxy and it will do the rest.

### Features

* Limiting users by username on Github
* Restrict to a single http method. Some users will only be able to do GET and the rest can do anything. This comes in REALLY handy when you want to limit view vs edit without worrying about it on your backend.
* Memorize the tokens, we don't DDOs gihub. Once a token has been verified it will be memorized and will not be checked.
* Save the token in the cookie for the user (See limitations regarding this feature).
* Docker friendly: Upstreams can be only visible to the main docker container and not accessible to the public. This makes the proxy the only gate to the code and you have to authenticate first.

## Configuration 

### Config file

```
{
  "authContext": "github",
  "users": [
    {
      "username": "KensoDev"
    },
    {
      "username": "KensoDev2",
      "restrict": "GET"
    }
  ],
  "upstreams": [
    {
      "type": "http",
      "prefix": "/",
      "location": "https://bdf682c0.ngrok.io"
    },
    {
      "type": "http",
      "prefix": "/api/",
      "location": "https://6cc1b022.ngrok.io"
    }
  ]
}

```

* `users` - Users that are allowed to access your app.
* `upstreams` - HTTP supported upstreams that the proxy will call
* `authContext` - Authentication context you want to use. Currently only
  supports `github`, see design docs for the interface implemented.

### Env Vars

### For Github

ENV vars you'll need to config in order for the proxy to work

```
  CLIENT_ID
  CLIENT_SECRET
```

### For Auth0

// TODO

These are your Github Client ID and Client secret from the Oauth page.

## Usage

```
usage: authproxy --listen-port=LISTEN-PORT --config-location=CONFIG-LOCATION [<flags>]

Flags:
  --help                     Show context-sensitive help (also try --help-long and --help-man).
  --listen=LISTEN-PORT  Which port should the proxy listen on
  --config=CONFIG-LOCATION
                             Proxy Config Location
```

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
* From your client side, you can either call the PROXY url as the backend or simply use the prefix `/api` you defined in your configuration.


## Design Decisions and architecture

### The authentication context

The authentication context is the class used in order to check whether a user
is valid, whether the action they try to do is allowed or not.

The interface supports a token and a username and returns simple objects such
as strings and booleans.

This was done by design in order to allow you to extend it and add many other
contexts beyond Github such as Auth0, Google and so on.


`IsAccessTokenValidAndUserAuthorized(accessToken string) bool`

This method is used to validate the token and check whether the user is
authorized. Checking whether they're included in the list of users or not


`GetUserName(accessToken sting) string`

Passing in an access token, this should return the username as a string


`GetHTTPEndpointPrefix() string`

What endpoint does the external service used to authentication needs exposing.

For example, Github wants you to expose a `/callback` which they send a `code`
to. You can read more about it [here in the docs.](https://developer.github.com/v3/guides/basics-of-authentication/)


`ServeHTTP(w http.ResponseWriter, req *http.Request)`

Your auth context needs to expose an HTTP handler in order to handle the
callback from the service.



### Github Auth Context Token Memoization Process

Once you have a Github Token in the http cookie, the proxy will validate that
token against Github. Make sure the user is authorized (from the user list) and will not validate again.

The cookie expiration is set for 24 hours, when you authenticate again, you will get revalidated.