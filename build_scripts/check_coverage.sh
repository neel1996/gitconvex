#!/bin/bash

if [[ $COVERAGE_PERCENT < $COVERAGE_THRESHOLD ]];
then
  echo "Coverage $COVERAGE_PERCENT is less that the threshold $COVERAGE_THRESHOLD" && exit 1;
else
  echo " ✔️ Good work! Coverage $COVERAGE_PERCENT meets the expectation"
fi;