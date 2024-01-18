# OpenIM Release Automation Design Document

This document outlines the automation process for releasing OpenIM. You can use the `make release` command for automated publishing. We will discuss how to use the `make release` command and Github Actions CICD separately, while also providing insight into the design principles involved.

## Github Actions Automation

In our CICD pipeline, we have implemented logic for automating the release process using the goreleaser tool. To achieve this, follow these steps on your local machine or server:

```bash
git clone https://github.com/openimsdk/open-im-server
cd open-im-server
git tag -a v3.6.0 -s -m "release: xxx"
# For pre-release versions: git tag -a v3.6.0-rc.0 -s -m "pre-release: xxx"
git push origin v3.6.0
```

The remaining tasks are handled by automated processes:

+ Automatically complete the release publication on Github
+ Automatically build the `v3.6.0` version image and push it to aliyun, dockerhub, and github

Through these automated steps, we achieve rapid and efficient OpenIM version releases, simplifying the release process and enhancing productivity.


Certainly, here is the continuation of the document in English:

## Local Make Release Design

There are two primary scenarios for local usage:

+ Advanced compilation and release, manually executed locally
+ Quick compilation verification and version release, manually executed locally

**These two scenarios can also be combined, for example, by tagging locally and then releasing:**

```bash
git add .
git commit -a -s -m "release(v3.6.0): ......"
git tag v3.6.0
git release
git push origin main
```

In a local environment, you can use the `make release` command to complete the release process. The main implementation logic can be found in the `/data/workspaces/open-im-server/scripts/lib/release.sh` file. First, let's explore its usage through the help information.

### Help Information

To view the help information, execute the following command:

```bash
$ ./scripts/release.sh --help
Usage: release.sh [options]
Options:
  -h, --help                Display this help message
  -se, --setup-env          Execute environment setup
  -vp, --verify-prereqs     Execute prerequisite verification
  -bc, --build-command      Execute build command
  -bi, --build-image        Execute build image (default is not executed)
  -pt, --package-tarballs   Execute tarball packaging
  -ut, --upload-tarballs    Execute tarball upload
  -gr, --github-release     Execute GitHub release
  -gc, --generate-changelog Execute changelog generation
```

### Default Behavior

If no options are provided, all operations are executed by default:

```bash
# If no options are provided, enable all operations by default
if [ "$#" -eq 0 ]; then
    perform_setup_env=true
    perform_verify_prereqs=true
    perform_build_command=true
    perform_package_tarballs=true
    perform_upload_tarballs=true
    perform_github_release=true
    perform_generate_changelog=true
    # TODO: Defaultly not enable build_image
    # perform_build_image=true
fi
```

### Environment Variable Setup

Before starting, you need to set environment variables:

```bash
export TENCENT_SECRET_KEY=OZZ****************************
export TENCENT_SECRET_ID=AKI****************************
```

### Modifying COS Account and Password

If you need to change the COS account, password, and bucket information, please modify the following section in the `/data/workspaces/open-im-server/scripts/lib/release.sh` file:

```bash
readonly BUCKET="openim-1306374445"
readonly REGION="ap-guangzhou"
readonly COS_RELEASE_DIR="openim-release"
```

### GitHub Release Configuration

If you intend to use the GitHub Release feature, you also need to set the environment variable:

```bash
export GITHUB_TOKEN="your_github_token"
```

### Modifying GitHub Release Basic Information

If you need to modify the basic information of GitHub Release, please edit the following section in the `/data/workspaces/open-im-server/scripts/lib/release.sh` file:

```bash
# OpenIM GitHub account information
readonly OPENIM_GITHUB_ORG=openimsdk
readonly OPENIM_GITHUB_REPO=open-im-server
```

This setup allows you to configure and execute the local release process according to your specific needs.


### GitHub Release Versioning Rules

Firstly, it's important to note that GitHub Releases should primarily be for pre-release versions. However, goreleaser might provide a `prerelease: auto` option, which automatically marks versions with pre-release indicators like `-rc1`, `-beta`, etc., as pre-releases.

So, if your most recent tag does not have pre-release indicators such as `-rc1` or `-beta`, even if you use `make release` for pre-release versions, goreleaser might still consider them as formal releases.

To avoid this issue, I have added the `--draft` flag to github-release. This way, all releases are created as drafts.

## CICD Release Documentation Design

The release records still require manual composition for GitHub Release. This is different from github-release.

This approach ensures that all releases are initially created as drafts, allowing you to manually review and edit the release documentation on GitHub. This manual step provides more control and allows you to curate release notes and other information before making them public.


## Makefile Section

This document aims to elaborate and explain key sections of the OpenIM Release automation design, including the Makefile section and functions within the code. Below, we will provide a detailed explanation of the logic and functions of each section.

In the project's root directory, the Makefile imports a subdirectory:

```makefile
include scripts/make-rules/release.mk
```

And defines the `release` target as follows:

```makefile
## release: release the project ✨
.PHONY: release release: release.verify release.ensure-tag
    @scripts/release.sh
```

### Importing Subdirectory

At the beginning of the Makefile, the `include scripts/make-rules/release.mk` statement imports the `release.mk` file from the subdirectory. This file contains rules and configurations related to releases to be used in subsequent operations.

### The `release` Target

The Makefile defines a target named `release`, which is used to execute the project's release operation. This target is marked as a phony target (`.PHONY`), meaning it doesn't represent an actual file or directory but serves as an identifier for executing a series of actions.

In the `release` target, two dependency targets are executed first: `release.verify` and `release.ensure-tag`. Afterward, the `scripts/release.sh` script is called to perform the actual release operation.

## Logic of `release.verify` and `release.ensure-tag`

```makefile
## release.verify: Check if a tool is installed and install it
.PHONY: release.verify
release.verify: tools.verify.git-chglog tools.verify.github-release tools.verify.coscmd tools.verify.coscli

## release.ensure-tag: ensure tag
.PHONY: release.ensure-tag
release.ensure-tag: tools.verify.gsemver
    @scripts/ensure-tag.sh
```

### `release.verify` Target

The `release.verify` target is used to check and install tools. It depends on four sub-targets: `tools.verify.git-chglog`, `tools.verify.github-release`, `tools.verify.coscmd`, and `tools.verify.coscli`. These sub-targets aim to check if specific tools are installed and attempt to install them if they are not.

The purpose of this target is to ensure that the necessary tools required for the release process are available so that subsequent operations can be executed smoothly.

### `release.ensure-tag` Target

The `release.ensure-tag` target is used to ensure that the project has a version tag. It depends on the sub-target `tools.verify.gsemver`, indicating that it should check if the `gsemver` tool is installed before executing.

When the `release.ensure-tag` target is executed, it calls the `scripts/ensure-tag.sh` script to ensure that the project has a version tag. Version tags are typically used to identify specific versions of the project for management and release in version control systems.

## Logic of `release.sh` Script

```bash
openim::golang::setup_env
openim::build::verify_prereqs
openim::release::verify_prereqs
#openim::build::build_image
openim::build::build_command
openim::release::package_tarballs
openim::release::upload_tarballs
git push origin ${VERSION}
#openim::release::github_release
#openim::release::generate_changelog
```

The `release.sh` script is responsible for executing the actual release operations. Below is the logic of this script:

1. `openim::golang::setup_env`: This function sets up some configurations for the Golang development environment.

2. `openim::build::verify_prereqs`: This function is used to verify whether the prerequisites for building are met. This includes checking dependencies, environment variables, and more.

3. `openim::release::verify_prereqs`: Similar to the previous function, this one is used to verify whether the prerequisites for the release are met. It focuses on conditions relevant to the release.

4. `openim::build::build_command`: This function is responsible for building the project's command, which typically involves compiling the project or performing other build operations.

5. `openim::release::package_tarballs`: This function is used to package tarball files required for the release. These tarballs are usually used for distribution packages during the release.

6. `openim::release::upload_tarballs`: This function is used to upload the packaged tarball files, typically to a distribution platform or repository.

7. `git push origin ${VERSION}`: This line of command pushes the version tag to the remote Git repository's `origin` branch, marking this release in the version control system.

In the comments, you can see that there are some operations that are commented out, such as `openim::build::build_image`, `openim::release::github_release`, and `openim::release::generate_changelog`. These operations are related to building images, releasing to GitHub, and generating changelogs, and they can be enabled in the release process as needed.

Let's take a closer look at the function responsible for packaging the tarball files:

```bash
function openim::release::package_tarballs() {
  # Clean out any old releases
  rm -rf "${RELEASE_STAGE}" "${RELEASE_TARS}" "${RELEASE_IMAGES}"
  mkdir -p "${RELEASE_TARS}"
  openim::release::package_src_tarball &
  openim::release::package_client_tarballs &
  openim::release::package_openim_manifests_tarball &
  openim::release::package_server_tarballs &
  openim::util::wait-for-jobs || { openim::log::error "previous tarball phase failed"; return 1; }

  openim::release::package_final_tarball & # _final depends on some of the previous phases
  openim::util::wait-for-jobs || { openim::log::error "previous tarball phase failed"; return 1; }
}
```

The `openim::release::package_tarballs()` function is responsible for packaging the tarball files required for the release. Here is the specific logic of this function:

1. `rm -rf "${RELEASE_STAGE}" "${RELEASE_TARS}" "${RELEASE_IMAGES}"`: First, the function removes any old release directories and files to ensure a clean starting state.

2. `mkdir -p "${RELEASE_TARS}"`: Next, it creates a directory `${RELEASE_TARS}` to store the packaged tarball files. If the directory doesn't exist, it will be created.

3. `openim::release::package_final_tarball &`: This is an asynchronous operation that depends on some of the previous phases. It is likely used to package the final tarball file, which includes the contents of all previous asynchronous operations.

4. `openim::util::wait-for-jobs`: It waits for all asynchronous operations to complete. If any of the previous asynchronous operations fail, an error will be returned.
