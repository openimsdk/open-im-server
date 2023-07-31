#  OpenIM Branch Management and Versioning

Our project, OpenIM, follows the [Semantic Versioning 2.0.0](https://semver.org/lang/zh-CN/) standards.

OpenIM, the open source project, employs a comprehensive version management system to ensure the reliability and traceability of our software. Our version management consists of three main components: the `main` branch, the `release` branch, and `tag` management.

## Main Branch

The `main` branch is where all the latest code resides. It's the hub of activity, embodying all the cutting-edge features that are currently being developed or updated. However, since it's subject to frequent changes and updates, it may not always represent the most stable version of the software. Access the `main` branch [here](https://github.com/openimsdk/openim-server/tree/main).

## Release Branch

On the other hand, we have the `release` branch. For instance, in the context of version 3.1, we maintain a `release-v3.1` branch. Unlike the `main` branch, the release branch is designed to be a continuously stable and updated version of the software. This provides a reliable option for users who prefer stability over the latest, but potentially unstable, features. Access the `release-v3.1` branch [here](https://github.com/openimsdk/openim-server/tree/release-v3.1).

## Tag Management

Finally, there's `tag` management. Despite having both `main` and `release` branches, `tag` serves a crucial role. Tags are immutable, i.e., they remain unchanged once created. Therefore, if you need a specific version of the software, you can use the corresponding tag. Check out the available tags [here](https://github.com/openimsdk/openim-server/tags).

Moreover, our Docker image versions are closely tied with these three components. For instance, a tag might correspond to the Docker image `ghcr.io/openimsdk/openim-server:v3.1.0`, a release would be `ghcr.io/openimsdk/openim-server:release-v3.0`, and the main branch might be represented as `ghcr.io/openimsdk/openim-server:main` or `ghcr.io/openimsdk/openim-server:latest`.

To find out more, or to contribute to our project, please visit our GitHub repository at [OpenIM Server](https://github.com/openimsdk/openim-server).

We believe that this approach offers a balanced blend of innovation and stability, enabling us to provide the best possible software to our users.

## OpenIM version

OpenIM manages two primary branches: `main` and `release`. The project uses Semantic Versioning 2.0.0 to tag different versions of the software, each indicating a significant milestone in the software's development.

In the OpenIM repository, the versioning adheres to the `MAJOR.MINOR.PATCH` format, where:

- `MAJOR` version changes when there are incompatible changes to the API,
- `MINOR` version changes when features are added in a backward-compatible manner, and
- `PATCH` version changes when backward-compatible bugs are fixed.

## Milestones and Branching

+ [OpenIM Milestones](https://github.com/OpenIMSDK/Open-IM-Server/milestones)
+ [OpenIM Tags](https://github.com/OpenIMSDK/Open-IM-Server/tags)
+ [OpenIM Branches](https://github.com/OpenIMSDK/Open-IM-Server/branches)

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


## Docker images version management

+ [OpenIM Docker Images Administration](https://github.com/OpenIMSDK/Open-IM-Server/blob/main/docs/conversions/images.md)
