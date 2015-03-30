# ongaku 音楽

[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/jmoiron/contact-form/master/LICENSE)

A stand-alone web server which runs in a directory and provides a simple
interface to play music files with html5.  This was made in order to make
it easy to cast music to a ChromeCast using the Chrome "cast this tab"
feature.

It uses `go generate` with `github.com/rakyll/statik` to embed the static
components into the executable.  It uses [jPlayer](http://jPlayer.org) to
play media.

It should be usable from any OS that runs Go, but has only been tested on
OSX and Linux.

## usage

```sh
$ ongaku -port=1339
2001/01/01 21:00:00 Listening on :1339
```

If you go to `:1339` on the machine you ran `ongaku` on, you should get
a (very) spartan interface to play music files.
