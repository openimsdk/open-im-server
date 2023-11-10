# Code conventions

- [Code conventions](#code-conventions)
  - [POSIX shell](#posix-shell)
  - [Go](#go)
  - [Directory and file conventions](#directory-and-file-conventions)
  - [Testing conventions](#testing-conventions)

## POSIX shell

- [Style guide](https://google.github.io/styleguide/shell.xml)

## Go

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go.html)
- Know and avoid [Go landmines](https://gist.github.com/lavalamp/4bd23295a9f32706a48f)
- Comment your code.
  - [Go's commenting conventions](http://blog.golang.org/godoc-documenting-go-code)
  - If reviewers ask questions about why the code is the way it is, that's a sign that comments might be helpful.
- Command-line flags should use dashes, not underscores
- Naming
  - Please consider package name when selecting an interface name, and avoid redundancy. For example, `storage.Interface` is better than `storage.StorageInterface`.
  - Do not use uppercase characters, underscores, or dashes in package names.
  - Please consider parent directory name when choosing a package name. For example, `pkg/controllers/autoscaler/foo.go` should say `package autoscaler` not `package autoscalercontroller`.
    - Unless there's a good reason, the `package foo` line should match the name of the directory in which the `.go` file exists.
    - Importers can use a different name if they need to disambiguate.Ⓜ️

## OpenIM Naming Conventions Guide

Welcome to the OpenIM Naming Conventions Guide. This document outlines the best practices and standardized naming conventions that our project follows to maintain clarity, consistency, and alignment with industry standards, specifically taking cues from the Google Naming Conventions.

### 1. General File Naming

Files within the OpenIM project should adhere to the following rules:

+ Both hyphens (`-`) and underscores (`_`) are acceptable in file names.
+ Underscores (`_`) are preferred for general files to enhance readability and compatibility.
+ For example: `data_processor.py`, `user_profile_generator.go`

### 2. Special File Types

#### a. Script and Markdown Files

+ Bash scripts and Markdown files should use hyphens (`-`) to facilitate better searchability and compatibility in web browsers.
+ For example: `deploy-script.sh`, `project-overview.md`

#### b. Uppercase Markdown Documentation

+ Markdown files with uppercase names, such as `README`, may include underscores (`_`) to separate words if necessary.
+ For example: `README_SETUP.md`, `CONTRIBUTING_GUIDELINES.md`

### 3. Directory Naming

+ Directories must use hyphens (`-`) exclusively to maintain a clean and organized file structure.
+ For example: `image-assets`, `user-data`

### 4. Configuration Files

+ Configuration files, including but not limited to `.yaml` files, should use hyphens (`-`).
+ For example: `app-config.yaml`, `logging-config.yaml`

### Best Practices

+ Keep names concise but descriptive enough to convey the file's purpose or contents at a glance.
+ Avoid using spaces in names; use hyphens or underscores instead to improve compatibility across different operating systems and environments.
+ Stick to lowercase naming where possible for consistency and to prevent issues with case-sensitive systems.
+ Include version numbers or dates in file names if the file is subject to updates, following the format: `project-plan-v1.2.md` or `backup-2023-03-15.sql`.

## Directory and file conventions

- Avoid general utility packages. Packages called "util" are suspect. Instead, derive a name that describes your desired function. For example, the utility functions dealing with waiting for operations are in the `wait` package and include functionality like `Poll`. The full name is `wait.Poll`.
- All filenames should be lowercase.
- All source files and directories should use underscores, not dashes.
  - Package directories should generally avoid using separators as much as possible. When package names are multiple words, they usually should be in nested subdirectories.

## Testing conventions

Please refer to [TESTING.md](https://github.com/openimsdk/open-im-server/tree/main/test/readme) document.
