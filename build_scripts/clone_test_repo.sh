#!/bin/bash

git clone https://github.com/neel1996/gitconvex-test.git $GITCONVEX_TEST_REPO

cd $GITCONVEX_TEST_REPO && \
git config user.name "test" && \
git config user.email "test@test.com"

git init $GITCONVEX_TEST_REPO/no_head

mkdir no_head_for_commit && \
touch ./no_head_for_commit/no_head.txt  && \
git init $GITCONVEX_TEST_REPO/no_head_for_commit