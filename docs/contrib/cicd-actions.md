# Continuous Integration and Automation

Every change on the OpenIM repository, either made through a pull request or direct push, triggers the continuous integration pipelines defined within the same repository. Needless to say, all the OpenIM contributions can be merged until all the checks pass (AKA having green builds).

- [Continuous Integration and Automation](#continuous-integration-and-automation)
  - [CI Platforms](#ci-platforms)
    - [GitHub Actions](#github-actions)
  - [Running locally](#running-locally)

## CI Platforms

Currently, there are two different platforms involved in running the CI processes:

- GitHub actions
- Drone pipelines on CNCF infrastructure

### GitHub Actions

All the existing GitHub Actions are defined as YAML files under the `.github/workflows` directory. These can be grouped into:

- **PR Checks**. These actions run all the required validations upon PR creation and update. Covering the DCO compliance check, `x86_64` test batteries (unit, integration, smoke), and code coverage.
- **Repository automation**. Currently, it only covers issues and epic grooming.

Everything runs on GitHub's provided runners; thus, the tests are limited to run in `x86_64` architectures.


## Running locally

A contributor should verify their changes locally to speed up the pull request process. Fortunately, all the CI steps can be on local environments, except for the publishing ones, through either of the following methods:

**User Makefile:**
```bash
root@PS2023EVRHNCXG:~/workspaces/openim/Open-IM-Server# make help 😊

Usage: make <TARGETS> <OPTIONS> ...

Targets:

all                          Run tidy, gen, add-copyright, format, lint, cover, build 🚀
build                        Build binaries by default 🛠️
multiarch                    Build binaries for multiple platforms. See option PLATFORMS. 🌍
tidy                         tidy go.mod ✨
vendor                       vendor go.mod 📦
style                        code style -> fmt,vet,lint 💅
fmt                          Run go fmt against code. ✨
vet                          Run go vet against code. ✅
lint                         Check syntax and styling of go sources. ✔️
format                       Gofmt (reformat) package sources (exclude vendor dir if existed). 🔄
test                         Run unit test. 🧪
cover                        Run unit test and get test coverage. 📊
updates                      Check for updates to go.mod dependencies 🆕
imports                      task to automatically handle import packages in Go files using goimports tool 📥
clean                        Remove all files that are created by building. 🗑️
image                        Build docker images for host arch. 🐳
image.multiarch              Build docker images for multiple platforms. See option PLATFORMS. 🌍🐳
push                         Build docker images for host arch and push images to registry. 📤🐳
push.multiarch               Build docker images for multiple platforms and push images to registry. 🌍📤🐳
tools                        Install dependent tools. 🧰
gen                          Generate all necessary files. 🧩
swagger                      Generate swagger document. 📖
serve-swagger                Serve swagger spec and docs. 🚀📚
verify-copyright             Verify the license headers for all files. ✅
add-copyright                Add copyright ensure source code files have license headers. 📄
release                      release the project 🎉
help                         Show this help info. ℹ️
help-all                     Show all help details info. ℹ️📚

Options:

DEBUG            Whether or not to generate debug symbols. Default is 0. ❓

BINS             Binaries to build. Default is all binaries under cmd. 🛠️
This option is available when using: make {build}(.multiarch) 🧰
Example: make build BINS="openim-api openim_cms_api".

PLATFORMS        Platform to build for. Default is linux_arm64 and linux_amd64. 🌍
This option is available when using: make {build}.multiarch 🌍
Example: make multiarch PLATFORMS="linux_s390x linux_mips64
linux_mips64le darwin_amd64 windows_amd64 linux_amd64 linux_arm64".

V                Set to 1 enable verbose build. Default is 0. 📝
```


How to Use Makefile to Help Contributors Build Projects Quickly 😊

The `make help` command is a handy tool that provides useful information on how to utilize the Makefile effectively. By running this command, contributors will gain insights into various targets and options available for building projects swiftly.

Here's a breakdown of the targets and options provided by the Makefile:

**Targets 😃**

1. `all`: This target runs multiple tasks like `tidy`, `gen`, `add-copyright`, `format`, `lint`, `cover`, and `build`. It ensures comprehensive project building.
2. `build`: The primary target that compiles binaries by default. It is particularly useful for creating the necessary executable files.
3. `multiarch`: A target that builds binaries for multiple platforms. Contributors can specify the desired platforms using the `PLATFORMS` option.
4. `tidy`: This target cleans up the `go.mod` file, ensuring its consistency.
5. `vendor`: A target that updates the project dependencies based on the `go.mod` file.
6. `style`: Checks the code style using tools like `fmt`, `vet`, and `lint`. It ensures a consistent coding style throughout the project.
7. `fmt`: Formats the code using the `go fmt` command, ensuring proper indentation and formatting.
8. `vet`: Runs the `go vet` command to identify common errors in the code.
9. `lint`: Validates the syntax and styling of Go source files using a linter.
10. `format`: Reformats the package sources using `gofmt`. It excludes the vendor directory if it exists.
11. `test`: Executes unit tests to ensure the functionality and stability of the code.
12. `cover`: Performs unit tests and calculates the test coverage of the code.
13. `updates`: Checks for updates to the project's dependencies specified in the `go.mod` file.
14. `imports`: Automatically handles import packages in Go files using the `goimports` tool.
15. `clean`: Removes all files generated during the build process, effectively cleaning up the project directory.
16. `image`: Builds Docker images for the host architecture.
17. `image.multiarch`: Similar to the `image` target, but it builds Docker images for multiple platforms. Contributors can specify the desired platforms using the `PLATFORMS` option.
18. `push`: Builds Docker images for the host architecture and pushes them to a registry.
19. `push.multiarch`: Builds Docker images for multiple platforms and pushes them to a registry. Contributors can specify the desired platforms using the `PLATFORMS` option.
20. `tools`: Installs the necessary tools or dependencies required by the project.
21. `gen`: Generates all the required files automatically.
22. `swagger`: Generates the swagger document for the project.
23. `serve-swagger`: Serves the swagger specification and documentation.
24. `verify-copyright`: Verifies the license headers for all project files.
25. `add-copyright`: Adds copyright headers to the source code files.
26. `release`: Releases the project, presumably for distribution.
27. `help`: Displays information about available targets and options.
28. `help-all`: Shows detailed information about all available targets and options.

**Options 😄**

1. `DEBUG`: A boolean option that determines whether or not to generate debug symbols. The default value is 0 (false).
2. `BINS`: Specifies the binaries to build. By default, it builds all binaries under the `cmd` directory. Contributors can provide a list of specific binaries using this option.
3. `PLATFORMS`: Specifies the platforms to build for. The default platforms are `linux_arm64` and `linux_amd64`. Contributors can specify multiple platforms by providing a space-separated list of platform names.
4. `V`: A boolean option that enables verbose build output when set to 1 (true). The default value is 0 (false).

With these targets and options in place, contributors can efficiently build projects using the Makefile. Happy coding! 🚀😊
