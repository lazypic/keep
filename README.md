# fork

fork simplify the fork workflow.

## Note

The only supported host is Github by now.

## Github

### Environment

You need to set these environment variables.

```
FORK_GITHUB_USER # your github id
FORK_GITHUB_AUTH # your github authentication token
```

### Create an Authentication token

Login into your account.

Goto 'Setting' > 'Developer Settings' > 'Personal access token'

Select 'Generate new token'. Turn on the 'repo' checkbox.
When you generate the token, you should see your token value.

Save it in your FORK_GITHUB_AUTH environment variable.

## Usage

Replace org and repo as you want in the following command.

```
fork github.com/:org/:repo
```

