#!/bin/sh

set -e

go build

assertEqual() {
  if [ "x$1" != "x$2" ]; then
    echo "FAIL" > /dev/stderr
    echo "Expected:\n\"$1\"\nto equal:\n\"$2\"" > /dev/stderr
    exit 1
  fi
}

# Test simple pass-through
assertEqual "$(./jsoncat tests/bare/true.json)" "true"
assertEqual "$(./jsoncat tests/bare/17.json)" "17"
assertEqual "$(./jsoncat tests/arrays/123.json)" "[1,2,3]"

# Test appending
assertEqual "$(./jsoncat tests/bare/17.json tests/bare/true.json tests/bare/17.json)" "[17,true,17]"
assertEqual "$(./jsoncat tests/arrays/123.json tests/arrays/456.json)" "[[1,2,3],[4,5,6]]"
assertEqual "$(./jsoncat tests/objects/a.json tests/objects/b.json)" '[{"a":1},{"b":2}]'

# Test merging
assertEqual "$(./jsoncat --merge tests/arrays/123.json tests/arrays/456.json)" "[1,2,3,4,5,6]"

# Test merging objects
assertEqual "$(./jsoncat --merge tests/objects/a.json tests/objects/b.json)" '{"a":1,"b":2}'

# Test merging incompatible types
# TODO

# Test reading invalid files
# TODO

echo "SUCCESS" > /dev/stderr
