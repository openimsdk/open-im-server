# ncpu

**ncpu** is a simple utility to fetch the number of CPU cores across different operating systems.

## Introduction

In various scenarios, especially while compiling code, it's beneficial to know the number of available CPU cores to optimize the build process. However, the command to fetch the CPU core count differs between operating systems. For example, on Linux, we use `nproc`, while on macOS, it's `sysctl -n hw.ncpu`. The `ncpu` utility provides a unified way to obtain this number, regardless of the platform.

## Usage

To retrieve the number of CPU cores, simply use the `ncpu` command:

```bash
$ ncpu
```

This will return an integer representing the number of available CPU cores.

### Example:

Let's say you're compiling a project using `make`. To utilize all the CPU cores for the compilation process, you can use:

```bash
$ make -j $(ncpu) build # or any other build command
```

The above command will ensure the build process takes advantage of all the available CPU cores, thereby potentially speeding up the compilation.

## Why use `ncpu`?

- **Cross-platform compatibility**: No need to remember or detect which OS-specific command to use. Just use `ncpu`!
  
- **Ease of use**: A simple and intuitive command that's easy to incorporate into scripts or command-line operations.

- **Consistency**: Ensures consistent behavior and output across different systems and environments.

## Installation

(Include installation steps here, e.g., how to clone the repo, build the tool, or install via package manager.)
