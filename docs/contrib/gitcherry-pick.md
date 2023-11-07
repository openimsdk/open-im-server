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

## Prerequisites

- [Contributor License Agreement](https://github.com/openim-sigs/cla) is considered implicit for all code within cherry pick pull requests, **unless there is a large conflict**.
- A pull request merged against the `main` branch.
- The release branch exists (example: [`release-1.18`](https://github.com/openimsdk/open-im-server/tree/release-v3.1))
- The normal git and GitHub configured shell environment for pushing to your openim-server `origin` fork on GitHub and making a pull request against a configured remote `upstream` that tracks `https://github.com/openimsdk/open-im-server.git`, including `GITHUB_USER`.
- Have GitHub CLI (`gh`) installed following [installation instructions](https://github.com/cli/cli#installation).
- A github personal access token which has permissions "repo" and "read:org". Permissions are required for [gh auth login](https://cli.github.com/manual/gh_auth_login) and not used for anything unrelated to cherry-pick creation process (creating a branch and initiating PR).

## What Kind of PRs are Good for Cherry Picks

Compared to the normal main branch's merge volume across time, the release branches see one or two orders of magnitude less PRs. This is because there is an order or two of magnitude higher scrutiny. Again, the emphasis is on critical bug fixes, e.g.,

- Loss of data
- Memory corruption
- Panic, crash, hang
- Security

A bugfix for a functional issue (not a data loss or security issue) that only affects an alpha feature does not qualify as a critical bug fix.

If you are proposing a cherry pick and it is not a clear and obvious critical bug fix, please reconsider. If upon reflection you wish to continue, bolster your case by supplementing your PR with e.g.,

- A GitHub issue detailing the problem
- Scope of the change
- Risks of adding a change
- Risks of associated regression
- Testing performed, test cases added
- Key stakeholder SIG reviewers/approvers attesting to their confidence in the change being a required backport

If the change is in cloud provider-specific platform code (which is in the process of being moved out of core openim-server), describe the customer impact, how the issue escaped initial testing, remediation taken to prevent similar future escapes, and why the change cannot be carried in your downstream fork of the openim-server project branches.

It is critical that our full community is actively engaged on enhancements in the project. If a released feature was not enabled on a particular provider's platform, this is a community miss that needs to be resolved in the `main` branch for subsequent releases. Such enabling will not be backported to the patch release branches.

## Initiate a Cherry Pick

### Before you begin

- Plan to initiate a cherry-pick against *every* supported release branch. If you decide to skip some release branch, explain your decision in a comment to the PR being cherry-picked.
- Initiate cherry-picks in order, from newest to oldest supported release branches. For example, if 3.1 is the newest supported release branch, then, before cherry-picking to 2.25, make sure the cherry-pick PR already exists for in 2.26 and 3.1. This helps to prevent regressions as a result of an upgrade to the next release.

### Steps

- Run the [cherry pick script](https://github.com/openimsdk/open-im-server/tree/main/scripts/cherry-pick.sh)

  This example applies a main branch PR #98765 to the remote branch `upstream/release-v3.1`:

  ```
  scripts/cherry-pick.sh upstream/release-v3.1 98765
  ```

  - Be aware the cherry pick script assumes you have a git remote called `upstream` that points at the openim-server github org.

    Please see our [recommended Git workflow](https://github.com/openimsdk/open-im-server/blob/main/docs/contributors/github-workflow.md#workflow).

  - You will need to run the cherry pick script separately for each patch release you want to cherry pick to. Cherry picks should be applied to all [active](https://github.com/openimsdk/open-im-server/releases) release branches where the fix is applicable.

  - If `GITHUB_TOKEN` is not set you will be asked for your github password: provide the github [personal access token](https://github.com/settings/tokens) rather than your actual github password. If you can securely set the environment variable `GITHUB_TOKEN` to your personal access token then you can avoid an interactive prompt. Refer [mislav/hub#2655 (comment)](https://github.com/mislav/hub/issues/2655#issuecomment-735836048)

- Your cherry pick PR will immediately get the `do-not-merge/cherry-pick-not-approved` label.


## Cherry Pick Review

As with any other PR, code OWNERS review (`/lgtm`) and approve (`/approve`) on cherry pick PRs as they deem appropriate.

The same release note requirements apply as normal pull requests, except the release note stanza will auto-populate from the main branch pull request from which the cherry pick originated.


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