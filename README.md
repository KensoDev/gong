# Gong

<img src="http://assets.avi.io/logo.svg" width="300" />

## Summary

Working with Jira can be a delight but it can also be a huge pain if you are
like me, working in the terminal most of the day. Opening the browser and
entering a bunch of fields can be a real pain.

I created (am creating) jglow to solve a bunch of my pains in managing my flow

## Usage

### Login

In order to use `gong` you first you need to login.

`gong login` will prompt you 3 things.

```
Type your username please: avi.zurel@DOMAIN.com
Please type your password. Don't worry, nothing will show on the screen: ************************
What is the instance URL for Jira (no need for https://): ORG.atlassian.net
Successfully logged in to Jira, congrats!
```

Login will check your credentials against Jira, if details are correct, it will
save the login details to disk. By default this goes to `$HOME/.gong.ini`

### Start working on an issue

`gong start {issue-id} --type feature`

If you want to start working on an issue, you can type in `gong start` with the
issue id and what type of work is this (defaults to feature).

This will do a couple of things

1. Create a branch name `{type}/{issue-id}-{issue-title-sluggified}`
2. Transition the issue to a started state

### `gong browse`

While working on a branch that matches the gong regular expression (look
above), you can type `gong browse` and it will open up a browser opened on the
issue automatically.

### `gong comment`

While working on a branch that matches the gong regular expression, you can
type `gong comment "some comment"` and it will send a comment on the ticket. 

This is a perfect way to streamline communication

## Work in progress

This is very much a work in progress and I am adding more features.

## Upcoming features

### gong slack

Send a message to a slack channel, tagging the issue you are working on

### `gong create`

Create a ticket, automatically giving you an id and starting to work on the
issue.

### `gong next/pick`

Show you the next items on your backlog, be able to start one without opening the browser
