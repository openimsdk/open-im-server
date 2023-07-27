#  OpenIM Branch Management and Versioning

Our project, OpenIM, follows the [Semantic Versioning 2.0.0](https://semver.org/lang/zh-CN/) standards.

## OpenIM version

OpenIM manages two primary branches: `main` and `release`. The project uses Semantic Versioning 2.0.0 to tag different versions of the software, each indicating a significant milestone in the software's development.

In the OpenIM repository, the versioning adheres to the `MAJOR.MINOR.PATCH` format, where:

- `MAJOR` version changes when there are incompatible changes to the API,
- `MINOR` version changes when features are added in a backward-compatible manner, and
- `PATCH` version changes when backward-compatible bugs are fixed.

## Milestones and Branching

When a significant milestone like v3.1.0 is achieved, a new branch `release-v3.1` is created. This branch contains all the code pertaining to this stable release. All bug fixes and features intended for the next version, v3.2.0, are merged into this branch.

The release of `PATCH` versions (Z in `X.Y.Z`) are driven by bug fixes, and these can be rolled out depending on the bug's priority or over a scheduled time. On the other hand, `MINOR` versions (Y in `X.Y.Z`) are released based on the project's roadmap, milestone completion, or on a scheduled timeline. Importantly, the API of minor versions is always backward-compatible.

## Dealing with Major Bugs

In the event of a major bug discovery, the fix would selectively be merged into the previous version (e.g., v3.1 or the `release-v3.1` branch), as well as into the `main` branch. This is to ensure that users relying on the older version can still receive important bug fixes, while also keeping the main branch updated.

It's worth noting that a robust testing regime should be in place to ensure the integrity of all branches at any given time. Automated tests and code review sessions are crucial components of maintaining a healthy codebase.

To summarize, OpenIM's approach to branch management and versioning ensures a balance between introducing new features, fixing bugs, and maintaining backward compatibility. This strategy is vital for managing user expectations, supporting older versions, and paving the way for the project's continuous growth.

## Git Workflow Example

To put the above principles into practice, here's a Git workflow example that you might follow when working on a bug fix:

```
bashCopy code# Checkout the branch for the version that needs the bug fix
git checkout release-v3.1

# Create a new branch for the bug fix
git checkout -b bug/bug-name

# ... Make changes, commit your work ...

# Push the branch to your remote repository
git push origin bug/bug-name

# After the pull request is merged into the release-v3.1 branch, 
# checkout and update your main branch
git checkout main
git pull origin main

# Merge or rebase the changes from release-v3.1 into main
git merge release-v3.1

# Push the updates to the main branch
git push origin main
```

Remember, communication with your team is key throughout this process, keeping everyone up-to-date with the changes being made.