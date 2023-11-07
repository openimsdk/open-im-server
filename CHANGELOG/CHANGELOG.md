# Changelog

- [Changelog](#changelog)
  - [OpenIM versioning policy](#openim-versioning-policy)
  - [command](#command)
    - [install](#install)
  - [User](#user)
  - [create next tag](#create-next-tag)
  - [Release version logs](#release-version-logs)
  - [Introduction](#introduction)
  - [Naming Format](#naming-format)
  - [Examples](#examples)
  - [Version Modifiers](#version-modifiers)
  - [Versioning Strategy](#versioning-strategy)


All notable changes to this project will be documented in this file.

+ [https://github.com/openimsdk/open-im-server/releases](https://github.com/openimsdk/open-im-server/releases)

## OpenIM versioning policy

+ [OpenIM Version](../docs/contrib/version.md)

## command

To use git-chglog you need to configure:

1. CHANGELOG templates
2. git-chglog configuration

### install 

```bash
$ go get github.com/git-chglog/git-chglog/cmd/git-chglog
```


## User

```bash
$ git-chglog --init
```

**Options**

- What is the URL of your repository?: https://github.com/openimsdk/open-im-server
- What is your favorite style?: github
- Choose the format of your favorite commit message: <type>(<scope>): <subject> -- feat(core): Add new feature
- What is your favorite template style?: standard
- Do you include Merge Commit in CHANGELOG?: n
- Do you include Revert Commit in CHANGELOG?: y
- In which directory do you output configuration files and templates?: .chglog

```bash
git-chglog --tag-filter-pattern 'v2.0.*'  -o CHANGELOG-2.0.md
```

**Other uses:**

```bash
$ git-chglog

  If <tag query> is not specified, it corresponds to all tags.
  This is the simplest example.

$ git-chglog 1.0.0..2.0.0

  The above is a command to generate CHANGELOG including commit of 1.0.0 to 2.0.0.

$ git-chglog 1.0.0

  The above is a command to generate CHANGELOG including commit of only 1.0.0.

$ git-chglog $(git describe --tags $(git rev-list --tags --max-count=1))

  The above is a command to generate CHANGELOG with the commit included in the latest tag.

$ git-chglog --output CHANGELOG.md

  The above is a command to output to CHANGELOG.md instead of standard output.

$ git-chglog --config custom/dir/config.yml

  The above is a command that uses a configuration file placed other than ".chglog/config.yml".
```


## create next tag

```bash
$ git-chglog --next-tag 2.0.0 -o CHANGELOG.md
$ git commit -am "release 2.0.0"
$ git tag 2.0.0
```

| Query          | Description                                    | Example                     |
| -------------- | ---------------------------------------------- | --------------------------- |
| `<old>..<new>` | Commit contained in `<new>` tags from `<old>`. | `$ git-chglog 1.0.0..2.0.0` |
| `<name>..`     | Commit from the `<name>` to the latest tag.    | `$ git-chglog 1.0.0..`      |
| `..<name>`     | Commit from the oldest tag to `<name>`.        | `$ git-chglog ..2.0.0`      |
| `<name>`       | Commit contained in `<name>`.                  | `$ git-chglog 1.0.0`        |


## Release version logs

+ [OpenIM CHANGELOG-V1.0](CHANGELOG-1.0.md)
+ [OpenIM CHANGELOG-V2.0](CHANGELOG-2.0.md)
+ [OpenIM CHANGELOG-V2.1](CHANGELOG-2.1.md)
+ [OpenIM CHANGELOG-V2.2](CHANGELOG-2.2.md)
+ [OpenIM CHANGELOG-V2.3](CHANGELOG-2.3.md)
+ [OpenIM CHANGELOG-V2.9](CHANGELOG-2.9.md)
+ [OpenIM CHANGELOG-V3.0](CHANGELOG-3.0.md)
+ [OpenIM CHANGELOG-V3.1](CHANGELOG-3.1.md)
+ [OpenIM CHANGELOG-V3.2](CHANGELOG-3.2.md)
+ [OpenIM CHANGELOG-V3.3](CHANGELOG-3.3.md)


## Introduction

In both the open-source and closed-source software development communities, it is important to follow a consistent and understandable versioning scheme for software projects. This ensures clear communication of changes, compatibility, and stability across different releases. One widely adopted naming convention is the Semantic Versioning 2.0.0.

## Naming Format

The most common format for version numbers is as follows:

```bash
major.minor[.patch[.build]]
```

Let's take a closer look at each component:

1. **Major Version**: This is the first number in the versioning scheme and indicates significant changes that may not be backward compatible (specific to each project).
2. **Minor Version**: The second number signifies the addition of new features while maintaining backward compatibility.
3. **Patch Version**: The third number represents bug fixes or code optimizations without introducing new features. It is generally backward compatible.
4. **Build Version**: Typically an automatically generated number that increments with each code commit.

## Examples

Here are a few examples to illustrate the versioning scheme:

1. `1.0`
2. `2.14.0.1478`
3. `3.2.1 build-354`

## Version Modifiers

Apart from the version numbers, there are also version modifiers used to indicate specific stages or statuses of a release. Some commonly used version modifiers include:

- **alpha**: An internal testing version with numerous known bugs. It is primarily used for communication among developers.
- **beta**: A testing version released to enthusiastic users for feedback and bug detection.
- **rc (release candidate)**: The final testing version before the official release.
- **ga (general availability)**: The initial stable release for public distribution.
- **r/release** (or no modifier at all): The final released version intended for general users.
- **lts (long-term support)**: Designates a version that will receive extended maintenance and bug fixes for a specified number of years.

## Versioning Strategy

To effectively manage version numbers, the following strategies are commonly employed:

- The initial version of a project can be either `0.1` or `1.0`.
- When fixing bugs, the patch version is incremented by 1.
- When adding new features, the minor version is incremented by 1, and the patch version is reset to 0.
- In the case of significant modifications, the major version is incremented by 1.
- The build version is usually automatically generated by the compilation process and follows a defined format. It does not require manual control.

By adhering to these strategies and guidelines, developers can maintain consistency and clarity in versioning their software projects. This enables users and collaborators to understand the nature of changes between different releases and ensure compatibility with their systems.

(Note: Markdown formatting has been used to structure this article. Markdown is a lightweight markup language used to format text on platforms like GitHub.)

------

**Note**: The above article is based on the given content and aims to provide a Markdown-formatted English article explaining the naming conventions for software project versions, specifically focusing on the Semantic Versioning 2.0.0.