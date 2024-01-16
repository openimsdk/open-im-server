# Development of a Go-Based Conformity Checker for Project File and Directory Naming Standards

### 1. Project Overview

#### Project Name

- `GoConformityChecker`

#### Functionality Description

- Checks if the file and subdirectory names in a specified directory adhere to specific naming conventions.
- Supports specific file types (e.g., `.go`, `.yml`, `.yaml`, `.md`, `.sh`, etc.).
- Allows users to specify directories to be checked and directories to be ignored.
- More read https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/code-conventions.md

#### Naming Conventions

- Go files: Only underscores are allowed.
- YAML, YML, and Markdown files: Only hyphens are allowed.
- Directories: Only underscores are allowed.

### 2. File Structure

- `main.go`: Entry point of the program, handles command-line arguments.
- `checker/checker.go`: Contains the core logic.
- `config/config.go`: Parses and stores configuration information.

### 3. Core Code Design

#### main.go

- Parses command-line arguments, including the directory to be checked and directories to be ignored.
- Calls the `checker` module for checking.

#### config.go

- Defines a configuration structure, such as directories to check and ignore.

#### checker.go

- Iterates through the specified directory.
- Applies different naming rules based on file types and directory names.
- Records files or directories that do not conform to the standards.

### 4. Pseudocode Example

#### main.go

```go
package main

import (
    "flag"
    "fmt"
    "GoConformityChecker/checker"
)

func main() {
    // Parse command-line arguments
    var targetDir string
    var ignoreDirs string
    flag.StringVar(&targetDir, "target", ".", "Directory to check")
    flag.StringVar(&ignoreDirs, "ignore", "", "Directories to ignore")
    flag.Parse()

    // Call the checker
    err := checker.CheckDirectory(targetDir, ignoreDirs)
    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

#### checker.go

```go
package checker

import (
    // Import necessary packages
)

func CheckDirectory(targetDir, ignoreDirs string) error {
    // Iterate through the directory, applying rules to check file and directory names
    // Return any found errors or non-conformities
    return nil
}
```

### 5. Implementation Details

- **File and Directory Traversal**: Use Go's `path/filepath` package to traverse directories and subdirectories.
- **Naming Rules Checking**: Apply different regex expressions for naming checks based on file extensions.
- **Error Handling and Reporting**: Record files or directories that do not conform and report to the user.

### 6. Future Development and Extensions

- Support more file types and naming rules.
- Provide more detailed error reports, such as showing line numbers and specific naming mistakes.
- Add a graphical or web interface for non-command-line users.

The above is an overview of the entire project's design. Following this design, specific coding implementation can begin. Note that the actual implementation may need adjustments based on real-world conditions.