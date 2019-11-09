# keep

keep organize your code!

keep downloads any source code into your $KEEPPATH.
It's inspired by $GOPATH. And you can also set $KEEPPATH to $GOPATH/src.

## Setup Environment

keep needs an environment that is independant of host.

```
KEEPPATH # where keep will use as root for it's cloned repositories.
```

If $KEEPPATH is not set, $HOME/src will be used by default.

keep also needs to be set host specific Environments.

## Usage

If you download a github repo, do following.
(Replace org and repo)

```
keep github.com/:org/:repo
```

If you want to fork and download the forked repo instead, do following.

```
keep -fork github.com/:org/:repo
```

Note that the only supported host is Github by now.

## Github

### Environment

You need to set these environment variables to create fork or access private repo.

```
KEEP_GITHUB_USER # your github id
KEEP_GITHUB_AUTH # your github authentication token
```

### Create an Authentication token

Login into your account.

Goto 'Setting' > 'Developer Settings' > 'Personal access token'

Select 'Generate new token'. Turn on the 'repo' checkbox.
When you generate the token, you should see your token value.

Save it in your KEEP_GITHUB_AUTH environment variable.


