# Gong

<img src="http://assets.avi.io/logo.svg" width="300" />

### Build Status 

* Develop: ![Build Status](https://travis-ci.org/KensoDev/gong.svg?branch=develop)
* Master : ![Build Status](https://travis-ci.org/KensoDev/gong.svg?branch=master)

## Summary

Gong is a CLI to make working with an issue tracker (look at the supported clients) and still keeping your flow going in the terminal.

You can easily start branches off of issues, comment and also link commits to the issue URL.

## Usage

### Installation 

Head over to the [Github releases](https://github.com/KensoDev/gong/releases).
The latest releases all have executables for OSX and linux.

I **did not** test gong on windows so if you want to build for windows and
test, please let me know.

Once you download the latest release, put it in your `PATH` and you can now use
`gong`


### Currently supported clients

* Jira

If you would like to contribute a different client, please feel free to submit a PR

### Login

In order to use `gong` you first you need to login.

`gong login {client-name}`

Each of the supported clients will prompt the required fields in order to login to the system. Jira will need username, pass and a couple more while others might only need an API token.

Once you input all of the details the client will attempt to login. If succeeded it will let you know.

[![asciicast](https://asciinema.org/a/dcko3kv5xwobpf4rgj0e4ulyo.png)](https://asciinema.org/a/dcko3kv5xwobpf4rgj0e4ulyo)

### Configuration

It is possible to add the following options to your `~/.gong.json`:
```json
"branch_replacement_character": "_"
```
Replace any illegal characters in a branch name with the provided character.  
By default `-` is used
```json
"branch_pattern": "{{.IssueID}}/{{.IssueTitle}}"
```
Use this pattern to create the branch name from the issue, accept `{{.IssueID}}`, `{{.IssueTitle}}` and `{{.IssueType}}`. `{{.IssueID}}` and `{{.IssueTitle}}` are mandatory.  
By default `{{.IssueType}}/{{.IssueID}}-{{.IssueTitle}}` is used

### Start working on an issue

`gong start {issue-id} --type feature`

If you want to start working on an issue, you can type in `gong start` with the
issue id and what type of work is this (defaults to feature).

This will do a couple of things

1. Create a branch name, by default: `{type}/{issue-id}-{issue-title-sluggified}`
2. Transition the issue to a started state

[![asciicast](https://asciinema.org/a/c5libsysjmb5f8f8gizkbldzv.png)](https://asciinema.org/a/c5libsysjmb5f8f8gizkbldzv)

### `gong browse`

While working on a branch that matches the gong regular expression (look
above), you can type `gong browse` and it will open up a browser opened on the
issue automatically.

### `gong comment`

While working on a branch that matches the gong regular expression, you can
type `echo "comment" | gong comment` and it will send a comment on the ticket.

### Why a pipe?

The reason for choosing a pipe and not just have the comment as an argument is to have the ability to send **any** output to the comment.

What I find most useful is to send diffs, files, buffers from vim and more.

With this approach, I find I write much better comments to tickets. You will do the same :)

[![asciicast](https://asciinema.org/a/d0rcjavbv55lbq1xpsrqiyyu6.png)](https://asciinema.org/a/d0rcjavbv55lbq1xpsrqiyyu6)

### `gong prepare-commit-message`

This is **not** meant to be used directly, instead it is meant to be wrapped with simple wrapper git hooks.

Sample hooks can be found in `git-hooks` directory.

All you need to do is to copy them into your `.git/hooks` directory.

This will add a link to the issue to every commit. Whether you do `git commit "commit message" or edit the commit message using the editor with `git commit`

### Install commit hooks on your repository 

```
curl https://raw.githubusercontent.com/KensoDev/gong/develop/git-hooks/prepare-commit-msg > .git/hooks/prepare-commit-msg
chmod +x .git/hooks/prepare-commit-msg

curl https://raw.githubusercontent.com/KensoDev/gong/develop/git-hooks/commit-msg > .git/hooks/commit-msg
chmod +x .git/hooks/commit-msg
```

### `gong create`

Gong create will open the browser on the issue tracker create ticket flow. You
can then copy over the issue-id and run `gong start` which will create the
branch and you cn start working on your ticket.

## Issues/Feedback

If you have any issues, please open one here on Github or hit me up on twitter [@KensoDev](https://twitter.com/KensoDev)

## CHANGELOG

## 1.6.0

* Added transitions to the config and outputting the transitions to STDout to
  verify the config.

## 1.4.0

* Added the pivotal tracker client. Thanks to [@stephensxu](http://github.com/stephensxu).
  In order to create the client and connect to pivotal tracker, you run `gong login pivotal`

### 1.3.4

* Added the `create` command. Opens up the browser on the create ticket URL for
  the specific issue tracker


## Upcoming features

### `gong slack`

Send a message to a slack channel, tagging the issue you are working on

### `gong next/pick`

Show you the next items on your backlog, be able to start one without opening the browser

