# Contributing to goScaffold

First off, thank you for considering contributing to `goScaffold`. It's people like you that make goScaffold such a great tool.

## Where do I go from here?

If you've noticed a bug or have a feature request, make sure to check our Issues to see if someone else has already created a ticket. If not, go ahead and make one!

## Fork & create a branch

If this is something you think you can fix, then fork goScaffold and create a branch with a descriptive name.

## Get the test suite running

Make sure you're using Go 1.26 or newer.
```bash
go test ./...
```
Ensure all tests pass before making your changes.

## Implement your fix or feature

At this point, you're ready to make your changes! Feel free to ask for help; everyone is a beginner at first.

## Make a Pull Request

At this point, you should switch back to your master branch and make sure it's up to date with goScaffold's master branch:

```bash
git remote add upstream https://github.com/arthurgray2k/goScaffold.git
git checkout main
git pull upstream main
```

Then update your feature branch from your local copy of master, and push it! Finally, go to GitHub and make a Pull Request.
