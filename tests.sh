#!/bin/bash

set -e

go build && go test

assertEqual() {
  if [[ "$1" != "$2" ]]; then
    echo "FAIL" > /dev/stderr
    echo -e "Expected:\n\"$1\"\nto equal:\n\"$2\"" > /dev/stderr
    exit 1
  fi
}

assertFailure() {
  local command="$1"
  local message="$2"

  # Verify error message on STDERR
  local actual_message=$($command 2>&1 >/dev/null)
  if ! echo $actual_message | grep -q "$message"; then
    echo "FAIL" > /dev/stderr
    echo -e "Expected $command to fail with error message containing:\n\t$message\nWas:\n\t$actual_message" > /dev/stderr
    exit 1
  fi

  # Verify exit status
  ($command >/dev/null 2>&1)
  local status_code=$?
  if [ $status_code -ne 1 ]; then
    echo "FAIL" > /dev/stderr
    echo "Expected $command to fail, but had status code $status_code" > /dev/stderr
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
assertFailure "./jsoncat --merge tests/objects/a.json tests/arrays/123.json" "incompatible type"

# Test reading invalid files
# TODO

echo "SUCCESS" > /dev/stderr
