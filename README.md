# Gong

<img src="http://assets.avi.io/logo.svg" width="300" />

### Readme driven development

This project is using readme driven development, making sure I don't miss
documentation or features. However, some things that are documented here still
don't exist.

For current usage, check the documentation.

## Summary

Working with Jira can be a delight but it can also be a huge pain if you are
like me, working in the terminal most of the day. Opening the browser and
entering a bunch of fields can be a real pain.

I created (am creating) jglow to solve a bunch of my pains in managing my flow

## Examples

1. Creating a new issue from a template, prepopulating some of the fields.
   Getting an issue-id back and opening a branch with that id.

When I want to start working on something, say someone slacked me a bug or a
problem I need to look at. It is really a pain to open the browser and start
the process of opening an issue.

This is solved by `gong issue create -t Devops`. Devops being a pre-populated
template that is selecting fields.

Typing this will open up vim for you (or emacs if you choose), you can type in
the title and the description into Vim. (First line is title, one blank line
for description)

```
This is the title


This is the description, I can type multiple
lines here
```

Saving this file and exiting vim will create the issue for you, outputting the
issue-id. Which you can then create a branch with.

2. Checking open issues that I have and picking one up.

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

