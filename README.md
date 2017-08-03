# micro-github-auth-proxy
Github Auth proxy designed for micro-apps with static client and an API server


## What is this project?

## Limitations

 * Currently, for `V1`, this only supports HTTP backend. So for example if you
    have a frontend and a backend, you can host your frontend on some service that
    support HTTP endpoint. eg: S3 with static hosting. Some `heroku` or anything
    that will give you a valid web address.

## Design Decisions and architecture

### Github Auth Context Token Memoization Process

Once you have a Github Token in the http cookie, the proxy will validate that
token against Github every 2 minutes (counted when the process is running).

For example:

You went to github, authenticated and you got the token `ABC` back. The proxy
verify that this token belongs to a user that the proxy allows through.

Once that token is verified, it is stored in memory for 2 minutes to avoid
re-validating against the Github API.
