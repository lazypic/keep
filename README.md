# Keep

![travisCI](https://secure.travis-ci.org/lazypic/keep.svg)

Keep organize your code!

Keep helps collaboration on Git hosts by setting upstream automatically.

Keep downloads any source code into your `$KEEPPATH`.

It's inspired by $GOPATH. And you can also set `$KEEPPATH` to `$GOPATH/src`.

## Setup Environment

Keep needs an environment that is independant of host.

```
KEEPPATH # where Keep will use as root for it's cloned repositories.
```

If `$KEEPPATH` is not set, `$HOME/src` will be used by default.

Keep also needs to be set host specific Environments.

## Usage

If you download a github repo, do following.
(Replace org and repo)

```bash
$ keep github.com/:org/:repo
```

If you want to fork and download the forked repo instead, do following.

```bash
$ keep -fork github.com/:org/:repo
```

Note that the only supported host is Github by now.

## Github

### Environment

You need to set these environment variables to setup 'origin-private' remote.
'origin-private' is a special remote that is same as 'origin' but with credential
to do some administrative actions.

```
GITHUB_USER # your github id
GITHUB_AUTH # your github authentication token
```

### Create an Authentication token

Login into your account.

Goto `Setting` > `Developer Settings` > `Personal access token`

Select `Generate new token`. Turn on the `repo` checkbox.
When you generate the token, you should see your token value.

Save it in your `GITHUB_AUTH` environment variable.

