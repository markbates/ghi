# GHI (GitHub Issues)

## Offline GitHub Issues Client

I travel a lot, and if you're reading this, so don't you! Sometimes the interwebs are hard to come by. For these cases I needed a good way to store my GitHub issues offline so I could review them whilst traveling, allowing me to be more effective.

"Wait, Mark, I know there are tools already that do this! Why the hell are you building another one!?" - Calm down. Wow! No need to get defensive. I built it because I tried about a dozen of them, and most of them I couldn't even get to install let alone run. This one is written in `Go` which means it's stupid easy (in theory) to run on anything, without a bazillion dependencies (I'm looking at you NPM!!).

## Simple

In addition to an app I could actually install, I wanted one that was stupid easy to use. I wanted to fetch all of the issues for a repo. I wanted to list them all by state, and I wanted to see the details of a particular one (along with comments). That's it. Simple. Oh, and I wanted to do it all in my terminal window. Who has the time for fancy-shmancy GUIs? Not this fella, that's for sure!

## Install

Assuming you already have [Go](https://golang.org) installed, all you need to do is the following:

```
$ go get github.com/markbates/ghi
```

That's it! I'll probably add cross-platform pre-built binaries at some point, if the demand is great enough. :)

## Usage

```
$ ghi help
```

That'll explain most of it. I hate repeating myself, so I won't.

### First Time Use

The first time you run this, you should read the `help` on the `fetch` command first.

```
$ ghi help fetch
```

## One Last Thing!

I probably should've mentioned this earlier, but this client is purely *READONLY*. Don't expect to be able to update, create, delete, wave at, give gifts to, or in any other way interactive with issues.
