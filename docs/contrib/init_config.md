# Init OpenIM Config

- [Init OpenIM Config](#init-openim-config)
  - [Start](#start)
  - [Define Automated Configuration](#define-automated-configuration)
  - [Define Configuration Variables](#define-configuration-variables)
    - [Bash Parsing Features](#bash-parsing-features)
  - [Reasons and Advantages of the Design](#reasons-and-advantages-of-the-design)


##  Start

With the increasing complexity of software engineering, effective configuration management has become more and more important. Yaml and other configuration files provide the necessary parameters and guidance for systems, but they also impose additional management overhead for developers. This article explores how to automate and optimize configuration management, thereby improving efficiency and reducing the chances of errors.

First, obtain the OpenIM code through the contributor documentation and initialize it following the steps below.

## Define Automated Configuration

We no longer strongly recommend modifying the same configuration file. If you have a new configuration file related to your business, we suggest generating and managing it through automation.

In the `scripts/init_config.sh` file, we defined some template files. These templates will be automatically generated to the corresponding directories when executing `make init`.

```
# Defines an associative array where the keys are the template files and the values are the corresponding output files.
declare -A TEMPLATES=(
  [""${OPENIM_ROOT}"/scripts/template/config-tmpl/env.template"]="${OPENIM_OUTPUT_SUBPATH}/bin/.env"
  [""${OPENIM_ROOT}"/scripts/template/config-tmpl/config.yaml"]="${OPENIM_OUTPUT_SUBPATH}/bin/config.yaml"
)
```

If you have your new mapping files, you can implement them by appending them to the array.

Lastly, run:

```
./scripts/init_config.sh
```

## Define Configuration Variables

In the `scripts/install/environment.sh` file, we defined many reusable variables for automation convenience.

In the provided example, the def function is a core element. This function not only provides a concise way to define variables but also offers default value options for each variable. This way, even if a specific variable is not explicitly set in an environment or scenario, it can still have an expected value.

```
function def() {
    local var_name="$1"
    local default_value="$2"
    eval "readonly $var_name=\${$var_name:-$default_value}"
}
```

### Bash Parsing Features

Since bash is a parsing script language, it executes commands in the order they appear in the script. This characteristic means we can define commonly used or core variables at the beginning of the script and then reuse or modify them later on.

For instance, we can initially set a universal password and reuse this password in subsequent variable definitions.

```
# Set a consistent password for easy memory
def "PASSWORD" "openIM123"

# Linux system user for openim
def "LINUX_USERNAME" "openim"
def "LINUX_PASSWORD" "${PASSWORD}"
```

## Reasons and Advantages of the Design

1. **Simplify Configuration Management**: Through automation scripts, we can avoid manual operations and configurations, thus reducing tedious repetitive tasks.
2. **Reduce Errors**: Manually editing yaml or other configuration files can lead to formatting mistakes or typographical errors. Automating with scripts can lower the risk of such errors.
3. **Enhanced Readability**: Using the `def` function and other bash scripting techniques, we can establish a clear, easy-to-read, and maintainable configuration system.
4. **Improved Reusability**: As demonstrated above, we can reuse variables and functions in different parts of the script, reducing redundant code and increasing overall consistency.
5. **Flexible Default Value Mechanism**: By providing default values for variables, we can ensure configurations are complete and consistent across various scenarios, while also offering customization options for advanced users.