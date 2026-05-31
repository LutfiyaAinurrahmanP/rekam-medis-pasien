#!/usr/bin/env bash

# Simple test runner that mimics the PowerShell script's summary output
PATH_ARG="./..."
FILTER_ARG=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    -Path)
      PATH_ARG="$2"; shift 2;;
    -Filter)
      FILTER_ARG="$2"; shift 2;;
    --path)
      PATH_ARG="$2"; shift 2;;
    --filter)
      FILTER_ARG="$2"; shift 2;;
    *) shift;;
  esac
done

echo
echo "PASS $PATH_ARG"
echo

packages=$(go list "$PATH_ARG" 2>/dev/null)
package_count=0
if [[ -n "$packages" ]]; then
  package_count=$(printf '%s\n' "$packages" | awk 'END { print NR }')
  echo "Test Packages:"
  echo "$packages"
  echo
fi

if [[ -z "$FILTER_ARG" ]]; then
  testOutput=$(go test -v -count=1 "$PATH_ARG" 2>&1)
else
  testOutput=$(go test -v -count=1 "$PATH_ARG" -run "$FILTER_ARG" 2>&1)
fi

# Print the test output
echo "$testOutput"

total=$(echo "$testOutput" | grep -c "=== RUN")
passed=$(echo "$testOutput" | grep -c "\-\-\- PASS:")
failed=$(echo "$testOutput" | grep -c "\-\-\- FAIL:")

if [[ $package_count -eq 0 ]]; then
  package_count=1
fi

echo
if [[ $failed -gt 0 ]]; then
  echo -e "\e[31mTest Suites: 1 failed, $package_count total\e[0m"
else
  echo -e "\e[32mTest Suites: $package_count passed, $package_count total\e[0m"
fi

if [[ $failed -gt 0 ]]; then
  echo -e "\e[31mTests:       $failed failed, $passed passed, $total total\e[0m"
else
  echo -e "\e[32mTests:       $passed passed, $total total\e[0m"
fi

echo "Snapshots:   0 total"
echo "Ran all test suites."
echo

if [[ $failed -gt 0 ]]; then
  exit 1
fi
