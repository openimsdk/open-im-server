## OpenIM Logging System: Design and Usage

**PATH:** `scripts/lib/logging.sh`

### Introduction

OpenIM, an intricate project, requires a robust logging mechanism to diagnose issues, maintain system health, and provide insights. A custom-built logging system embedded within OpenIM ensures consistent and structured logs. Let's delve into the design of this logging system and understand its various functions and their usage scenarios.

### Design Overview

1. **Initialization**: The system begins by determining the verbosity level through the `OPENIM_VERBOSE` variable. If it's not set, a default value of 5 is assigned. This verbosity level dictates the depth of the log details.
2. **Log File Setup**: Logs are stored in the directory specified by `OPENIM_OUTPUT`. If this variable isn't explicitly set, it defaults to the `_output` directory relative to the script location. Each log file is named based on the date to facilitate easy identification.
3. **Logging Function**: The `echo_log()` function plays a pivotal role by writing messages to both the console (stdout) and the log file.
4. **Logging to a file**: The `echo_log()` function writes to the log file by appending the message to the file. It also adds a timestamp to the message. path: `_output/logs/*`, Enable logging by default. Set to false to disable. If you wish to turn off output to log files set `ENABLE_LOGGING=flase`. 

### Key Functions & Their Usages

1. **Error Handling**:
   - `openim::log::errexit()`: Activated when a command exits with an error. It prints a call tree showing the sequence of functions leading to the error and then calls `openim::log::error_exit()` with relevant information.
   - `openim::log::install_errexit()`: Sets up the trap for catching errors and ensures that the error handler (`errexit`) gets propagated to various script constructs like functions, expansions, and subshells.
2. **Logging Levels**:
   - `openim::log::error()`: Logs error messages with a timestamp. The log message starts with '!!!' to indicate its severity.
   - `openim::log::info()`: Provides informational messages. The display of these messages is governed by the verbosity level (`OPENIM_VERBOSE`).
   - `openim::log::progress()`: Designed for logging progress messages or creating progress bars.
   - `openim::log::status()`: Logs status messages with a timestamp, prefixing each entry with '+++' for easy identification.
   - `openim::log::success()`: Highlights successful operations with a bright green prefix. It's ideal for visually signifying operations that completed successfully.
3. **Exit and Stack Trace**:
   - `openim::log::error_exit()`: Logs an error message, dumps the call stack, and exits the script with a specified exit code.
   - `openim::log::stack()`: Prints out a stack trace, showing the call hierarchy leading to the point where this function was invoked.
4. **Usage Information**:
   - `openim::log::usage() & openim::log::usage_from_stdin()`: Both functions provide a mechanism to display usage instructions. The former accepts arguments directly, while the latter reads them from stdin.
5. **Test Function**:
   - `openim::log::test_log()`: This function is a test suite to verify that all logging functions are operating as expected.

### Usage Scenario

Imagine a situation where an OpenIM operation fails, and you need to ascertain the cause. With the logging system in place, you can:

- Check the log file for the specific day to find error messages with the '!!!' prefix.
- View the call tree and stack trace to trace back the sequence of operations leading to the failure.
- Use the verbosity level to filter out unnecessary details and focus on the crux of the issue.

This systematic and structured approach greatly simplifies the debugging process, making system maintenance more efficient.

### Conclusion

OpenIM's logging system is a testament to the importance of structured and detailed logging in complex projects. By using this logging mechanism, developers and system administrators can streamline troubleshooting and ensure the seamless operation of the OpenIM project.