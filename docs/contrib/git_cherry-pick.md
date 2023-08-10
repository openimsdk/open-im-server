# Git Cherry-Pick Guide

- Git Cherry-Pick Guide
  - [Introduction](#introduction)
  - [What is git cherry-pick?](#what-is-git-cherry-pick)
  - [Using git cherry-pick](#using-git-cherry-pick)
  - [Applying Multiple Commits](#applying-multiple-commits)
  - [Configurations](#configurations)
  - [Handling Conflicts](#handling-conflicts)
  - [Applying Commits from Another Repository](#applying-commits-from-another-repository)

## Introduction

Author: @cubxxw

As OpenIM has progressively embarked on a standardized path, I've had the honor of initiating a significant project, `git cherry-pick`. While some may see it as merely a naming convention in the Go language, it represents more. It's a thoughtful design within the OpenIM project, my very first conscious design, and a first in laying out an extensive collaboration process and copyright management with goals of establishing a top-tier community standard.

## What is git cherry-pick?

In multi-branch repositories, transferring commits from one branch to another is common. You can either merge all changes from one branch (using `git merge`) or selectively apply certain commits. This selective application of commits is where `git cherry-pick` comes into play.

Our collaboration strategy with GitHub necessitates maintenance of multiple `release-v*` branches alongside the `main` branch. To manage this, we mainly develop on the `main` branch and selectively merge into `release-v*` branches. This ensures the `main` branch stays current while the `release-v*` branches remain stable.

Ensuring this strategy's success extends beyond just documentation; it hinges on well-engineered solutions and automation tools, like Makefile, powerful CI/CD processes, and even Prow.

## Using git cherry-pick

`git cherry-pick` applies specified commits from one branch to another.

```bash
$ git cherry-pick <commitHash>
```

As an example, consider a repository with `main` and `release-v3.1` branches. To apply commit `f` from the `release-v3.1` branch to the `main` branch:

```
# Switch to main branch
$ git checkout main

# Perform cherry-pick
$ git cherry-pick f
```

You can also use a branch name instead of a commit hash to cherry-pick the latest commit from that branch.

```bash
$ git cherry-pick release-v3.1
```

## Applying Multiple Commits

To apply multiple commits simultaneously:

```bash
$ git cherry-pick <HashA> <HashB>
```

To apply a range of consecutive commits:

```bash
$ git cherry-pick <HashA>..<HashB>
```

## Configurations

Here are some commonly used configurations for `git cherry-pick`:

- **`-e`, `--edit`**: Open an external editor to modify the commit message.
- **`-n`, `--no-commit`**: Update the working directory and staging area without creating a new commit.
- **`-x`**: Append a reference in the commit message for tracking the source of the cherry-picked commit.
- **`-s`, `--signoff`**: Add a sign-off message at the end of the commit indicating who performed the cherry-pick.
- **`-m parent-number`, `--mainline parent-number`**: When the original commit is a merge of two branches, specify which parent branch's changes should be used.

## Handling Conflicts

If conflicts arise during the cherry-pick:

- **`--continue`**: After resolving conflicts, stage the changes with `git add .` and then continue the cherry-pick process.
- **`--abort`**: Abandon the cherry-pick and revert to the previous state.
- **`--quit`**: Exit the cherry-pick without reverting to the previous state.

## Applying Commits from Another Repository

You can also cherry-pick commits from another repository:

1. Add the external repository as a remote:

   ```
   $ git remote add target git://gitUrl
   ```

2. Fetch the commits from the remote:

   ```
   $ git fetch target
   ```

3. Identify the commit hash you wish to cherry-pick:

   ```
   $ git log target/main
   ```

4. Perform the cherry-pick:

   ```
   $ git cherry-pick <commitHash>
   ```