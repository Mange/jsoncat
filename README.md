# jsoncat

> `cat`, but for JSON

## Usage

    $ jsoncat foo.json bar.json > baz.json

### Merging

By default `jsoncat` will wrap the input files in an array:

    $ cat paul.json
    {"name": "paul"}

    $ cat kat.json
    {"name": "katherine"}

    $ jsoncat paul.json kat.json
    [{"name":"paul"},{"name":"katherine"}]

In case you have root objects that you want to merge you can pass the `--merge` option.

    $ cat movies.json
    ["Shawshank Redemption", "Primer"]

    $ cat music.json
    ["The Wall", "A Violent Emotion"]

    $ jsoncat movies.json music.json
    [["Shawshank Redemption","Primer"],["The Wall","A Violent Emotion"]]

    $ jsoncat --merge movies.json music.json
    ["Shawshank Redemption","Primer","The Wall","A Violent Emotion"]

The command will fail on incompatible types:

    $ jsoncat --merge object.json array.json
		Cannot merge incompatible types: Object, Array
    $ echo $?
    1

## Install

*Option 1:* Set up [a proper `$GOPATH`][gopath] and then run `go install github.com/Mange/jsoncat`.

*Option 2:* Download the source and run `go build && cp jsoncat ~/bin`.

## Missing stuff

* `STDIN` is ignored
* Needs to support the "-" argument (for `STDIN`)

## License

Public domain.

[gopath]: http://golang.org/doc/code.html#Organization
