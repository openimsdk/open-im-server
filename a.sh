#!/bin/bash
git checkout --orphan temp $1
git commit -m "v3 - main to cut out"
git rebase --onto temp $1 main-test2
git branch -D temp

