# OpenIM Configuration Guide

<!-- vscode-markdown-toc -->
* 1. [Directory Structure and File Descriptions](#DirectoryStructureandFileDescriptions)
	* 1.1. [Directory Structure](#DirectoryStructure)
	* 1.2. [Directory Structure Explanation](#DirectoryStructureExplanation)
* 2. [File Descriptions](#FileDescriptions)
	* 2.1. [Files in the Root Directory](#FilesintheRootDirectory)
	* 2.2. [Files in the `templates/` Directory](#FilesinthetemplatesDirectory)
* 3. [Configuration File Generation](#ConfigurationFileGeneration)
	* 3.1. [How to Use `init-config.sh` Script](#HowtoUseinit-config.shScript)
	* 3.2. [Examples of Operations](#ExamplesofOperations)
	* 3.3. [Points to Note](#PointstoNote)
* 4. [Example Directory](#ExampleDirectory)
	* 4.1. [Overview](#Overview)
	* 4.2. [Structure](#Structure)
	* 4.3. [How to Use These Examples](#HowtoUseTheseExamples)
	* 4.4. [Tips for Using Example Files:](#TipsforUsingExampleFiles:)
* 5. [Configuration Item Descriptions](#ConfigurationItemDescriptions)
* 6. [Version Management and Upgrading](#VersionManagementandUpgrading)
	* 6.1. [Pulling the Latest Code](#PullingtheLatestCode)
	* 6.2. [Generating the Latest Example Configuration Files](#GeneratingtheLatestExampleConfigurationFiles)
	* 6.3. [Comparing Configuration File Differences](#ComparingConfigurationFileDifferences)
	* 6.4. [Updating Configuration Files](#UpdatingConfigurationFiles)
	* 6.5. [Updating Binary Files and Restarting Services](#UpdatingBinaryFilesandRestartingServices)
	* 6.6. [Best Practices for Version Management](#BestPracticesforVersionManagement)
* 7. [How to Contribute](#HowtoContribute)
	* 7.1. [OpenIM Configuration Item Descriptions](#OpenIMConfigurationItemDescriptions)
	* 7.2. [Modifying Template Files](#ModifyingTemplateFiles)
	* 7.3. [Updating Configuration Center Scripts](#UpdatingConfigurationCenterScripts)
	* 7.4. [Configuration File Generation Process](#ConfigurationFileGenerationProcess)
	* 7.5. [Contribution Guidelines](#ContributionGuidelines)
	* 7.6. [Submission and Review](#SubmissionandReview)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->


##  1. <a name='DirectoryStructureandFileDescriptions'></a>Directory Structure and File Descriptions

This document details the structure of the `config` directory, aiding users in understanding and managing configuration files.

###  1.1. <a name='DirectoryStructure'></a>Directory Structure

```bash
$ tree config
├── alertmanager.yml
├── config.yaml
├── email.tmpl
├── instance-down-rules.yml
├── notification.yaml
├── prometheus.yml
├── Readme.md
└── templates
    ├── alertmanager.yml.template
    ├── config.yaml.template
    ├── email.tmpl.template
    ├── env.template
    ├── instance-down-rules.yml.template
    ├── notification.yaml.template
    ├── open-im-ng-example.conf
    ├── prometheus-dashboard.yaml
    └── prometheus.yml.template
```

###  1.2. <a name='DirectoryStructureExplanation'></a>Directory Structure Explanation

- **Root Directory (`config/`)**: Contains actual configuration files and the `templates` subdirectory.
- **`templates/` Subdirectory**: Stores configuration templates for generating or updating configuration files in the root directory.

##  2. <a name='FileDescriptions'></a>File Descriptions

###  2.1. <a name='FilesintheRootDirectory'></a>Files in the Root Directory

- **`alertmanager.yml`**: Configuration file for AlertManager, managing and setting up the alert system.
- **`config.yaml`**: The main application configuration file, covering service settings.
- **`email.tmpl`**: Template file for email notifications, defining email format and content.
- **`instance-down-rules.yml`**: Instance downtime rules configuration file for the monitoring system.
- **`notification.yaml`**: Configuration file for notification settings, defining different types of notifications.
- **`prometheus.yml`**: Configuration file for the Prometheus monitoring system, setting monitoring metrics and rules.

###  2.2. <a name='FilesinthetemplatesDirectory'></a>Files in the `templates/` Directory

- **`alertmanager.yml.template`**: Template for AlertManager configuration.
- **`config.yaml.template`**: Main configuration template for the application.
- **`email.tmpl.template`**: Template for email notifications.
- **`env.template`**: Template for environmental variable configurations, setting environment-related configurations.
- **`instance-down-rules.yml.template`**: Template for instance downtime rules.
- **`notification.yaml.template`**: Template for notification settings.
- **`open-im-ng-example.conf`**: Example configuration file for the application.
- **`prometheus-dashboard.yaml`**: Prometheus dashboard configuration file, specific to the OpenIM application.
- **`prometheus.yml.template`**: Template for Prometheus configuration.

##  3. <a name='ConfigurationFileGeneration'></a>Configuration File Generation

Configuration files can be automatically generated using the `make init` command or the `./scripts/init-config.sh` script. These scripts conveniently extract templates from the `templates` directory and generate or update actual configuration files in the root directory.

###  3.1. <a name='HowtoUseinit-config.shScript'></a>How to Use `init-config.sh` Script

```bash
$ ./scripts/init-config.sh --help
Usage: init-config.sh [options]
Options:
  -h, --help             Show this help message
  --force                Overwrite existing files without prompt
  --skip                 Skip generation if file exists
  --examples             Generate example files
  --clean-config         Clean all configuration files
  --clean-examples       Clean all example files
```

###  3.2. <a name='ExamplesofOperations'></a>Examples of Operations

- Generate all template configuration files:

  ```bash
  $ ./scripts/init-config.sh --examples
  ```

- Force overwrite existing configuration files:

  ```bash
  $ ./scripts/init-config.sh --force
  ```

###  3.3. <a name='PointstoNote'></a>Points to Note

- **Template files should not be directly modified**: Files in the `template` directory are templates included in source code management. Direct modification may lead to version conflicts or management issues.
- **Operations for Windows Users**: Windows users can use the `cp` command to copy files from the `template` directory to the `config/` directory and then modify the configuration items as needed.

##  4. <a name='ExampleDirectory'></a>Example Directory

Welcome to our project's `examples` directory! This directory contains a range of example files, showcasing various configurations and settings of our software. These examples are intended to provide you with templates that can serve as a starting point for your own configurations.

###  4.1. <a name='Overview'></a>Overview

In this directory, you'll find examples suitable for a variety of use cases. Each file is a template with default values and configurations, demonstrating best practices and typical scenarios. Whether you're just getting started or looking to implement complex settings, these examples should help you get on the right track.

###  4.2. <a name='Structure'></a>Structure

Here's a quick overview of the contents in this directory:

- `env-example.yaml`: Demonstrates how to set up environmental variables.
- `openim-example.yaml`: Example configuration file for the OpenIM application.
- `prometheus-example.yml`: Example configuration for monitoring with Prometheus.
- `alertmanager-example.yml`: Template for setting up Alertmanager configuration.

###  4.3. <a name='HowtoUseTheseExamples'></a>How to Use These Examples

To use these examples, simply copy the relevant files to your working directory and rename them as needed (for example, removing the `-example` suffix). Then, modify the files according to your needs.

###  4.4. <a name='TipsforUsingExampleFiles:'></a>Tips for Using Example Files:

1. **Read Comments**: Each file contains comments explaining the various sections and settings. Make sure to read these comments for a better understanding of how to customize the file.
2. **Check Required Changes**: Some examples might require mandatory changes before they can be used effectively (such as setting specific environmental variables).
3. **Version Compatibility**: Ensure that the example files are compatible with the version of the software you are using.

##  5. <a name='ConfigurationItemDescriptions'></a>Configuration Item Descriptions

##  6. <a name='VersionManagementandUpgrading'></a>Version Management and Upgrading

When managing and upgrading the `config` directory's versions, it is crucial to ensure that the configuration files in both the local `config/` and `config/templates/` directories are kept in sync. This process can ensure that your configuration files are consistent with the latest standard templates, while also maintaining custom settings.

###  6.1. <a name='PullingtheLatestCode'></a>Pulling the Latest Code

First, ensure that your local repository is in sync with the remote repository. This can be achieved by pulling the latest code:

```bash
$ git pull
```

###  6.2. <a name='GeneratingtheLatestExampleConfigurationFiles'></a>Generating the Latest Example Configuration Files

Next, generate the latest example configuration files. This can be done by running the `init-config.sh` script, using the `--examples` option to generate example files, and the `--skip` option to avoid overwriting existing configuration files:

```bash
$ ./scripts/init-config.sh --examples --skip
```

###  6.3. <a name='ComparingConfigurationFileDifferences'></a>Comparing Configuration File Differences

Once the latest example configuration files are generated, you need to compare the configuration files in the `config/` and `config/templates/` directories to find any potential differences. This step ensures that you can identify and integrate any important updates or changes. Tools like `diff` can be helpful in completing this step:

```bash
$ diff -ur config/ config/templates/
```

###  6.4. <a name='UpdatingConfigurationFiles'></a>Updating Configuration Files

Based on the comparison results, manually update the configuration files in the `config/` directory to reflect the latest configurations in `config/templates/`. During this process, ensure to retain any custom configuration settings.

###  6.5. <a name='UpdatingBinaryFilesandRestartingServices'></a>Updating Binary Files and Restarting Services

After updating the configuration files, the next step is to update any related binary files. This typically involves downloading and installing the latest version of the application or service. Depending on the specific application or service, this might involve running specific update scripts or directly downloading the latest version from official sources.

Once the binary files are updated, the services need to be restarted to apply the new configurations. Make sure to conduct necessary checks before restarting to ensure the correctness of the configurations.

###  6.6. <a name='BestPracticesforVersionManagement'></a>Best Practices for Version Management

- **Record Changes**: When committing changes to a version control system, ensure to log detailed change logs.
- **Stay Synced**: Regularly sync with the remote repository to ensure that your local configurations are in line with the latest developments.
- **Backup**: Backup your current configurations before making any significant changes, so that you can revert to a previous state if necessary.

By following these steps and best practices, you can ensure effective management and smooth upgrading of your `config` directory.

##  7. <a name='HowtoContribute'></a>How to Contribute

If you have an understanding of the logic behind OpenIM's configuration generation, then you will clearly know where to make modifications to contribute code.

###  7.1. <a name='OpenIMConfigurationItemDescriptions'></a>OpenIM Configuration Item Descriptions

First, it is recommended to read the [OpenIM Configuration Items Document](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/environment.md). This will help you understand the roles of various configuration items and how they affect the operation of OpenIM.

###  7.2. <a name='ModifyingTemplateFiles'></a>Modifying Template Files

To contribute to OpenIM, focus on the `./deployments/templates` directory. This contains various configuration template files, which are the basis for generating the final configuration files.

When making modifications, ensure that your changes align with OpenIM's configuration requirements and logic. This may involve adding new template files or modifying existing files to reflect new configuration options or structural changes.

###  7.3. <a name='UpdatingConfigurationCenterScripts'></a>Updating Configuration Center Scripts

In addition to modifying template files, pay attention to the `./scripts/install/environment.sh` script. In this script, you may need to add or modify environment variables.

This script is responsible for defining environment variables that influence configuration generation. Therefore, any new configuration items or modifications to existing items need to be reflected here.

###  7.4. <a name='ConfigurationFileGenerationProcess'></a>Configuration File Generation Process

The essence of the `make init` command is to use the environment variables defined in `/scripts/install/environment.sh` to render the template files in the `./deployments/templates` directory, thereby generating the final configuration files.

When contributing code, ensure that your changes work smoothly in this process and do not cause errors during configuration file generation.

###  7.5. <a name='ContributionGuidelines'></a>Contribution Guidelines

- **Code Review**: Ensure your changes have passed code review. This typically means that the code should be clear, easy to understand, and adhere to the project's coding style and best practices.
- **Testing**: Before submitting changes, conduct thorough tests to ensure new or modified configurations work as expected and do not negatively impact existing functionalities.
- **Documentation**: If you have added a new configuration option or made significant changes to an existing one, update the relevant documentation to assist other users and developers in understanding and utilizing these changes.

###  7.6. <a name='SubmissionandReview'></a>Submission and Review

After completing your changes, submit your code to the OpenIM repository in the form of a Pull Request (PR). The PR will be reviewed by the project maintainers and you may be asked to make further modifications or provide additional information.