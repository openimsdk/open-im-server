# OpenIM Protoc Tool

## Introduction

OpenIM is passionate about ensuring that its suite of tools is custom-tailored to cater to the unique needs of its users. That commitment led us to develop and release our custom Protoc tool, version v1.0.0.

### Why a Custom Version?

There are several reasons to choose our custom Protoc tool over generic open-source versions:

- **Specialized Features**: OpenIM's Protoc tool has been enriched with features and plugins that are optimized for the OpenIM ecosystem. This makes it more aligned with the needs of OpenIM users.
- **Optimized Performance**: Built from the ground up with OpenIM's infrastructure in mind, our tool guarantees faster and more efficient operations.
- **Enhanced Compatibility**: Our Protoc tool ensures full compatibility with OpenIM's offerings, minimizing potential conflicts and integration challenges.
- **Rich Output Support**: Unlike generic tools, our custom tool provides a wide array of output options including C++, C#, Java, Kotlin, Objective-C, PHP, Python, Ruby, and more. This allows developers to generate code for their preferred platform with ease.

## Download

+ https://github.com/OpenIMSDK/Open-IM-Protoc

Access the official release of the Protoc tool on the OpenIM repository here: [OpenIM Protoc Tool v1.0.0 Release](https://github.com/OpenIMSDK/Open-IM-Protoc/releases/tag/v1.0.0)

### Direct Download Links:

- **Windows**: [Download for Windows](https://github.com/OpenIMSDK/Open-IM-Protoc/releases/download/v1.0.0/windows.zip)
- **Linux**: [Download for Linux](https://github.com/OpenIMSDK/Open-IM-Protoc/releases/download/v1.0.0/linux.zip)

## Installation

For Windows:

1. Navigate to the Windows download link provided above and download the version suitable for your system.
2. Extract the contents of the zip file.
3. Add the path of the extracted tool to your `PATH` environment variable to run the Protoc tool directly from the command line.

For Linux:

1. Navigate to the Linux download link provided above and download the version suitable for your system.
2. Extract the contents of the zip file.
3. Use `chmod +x ./*` to make the extracted files executable.
4. Add the path of the extracted tool to your `PATH` environment variable to run the Protoc tool directly from the command line.

## Usage

The OpenIM Protoc tool provides a multitude of options for parsing `.proto` files and generating output:

```

./protoc [OPTION] PROTO_FILES
```

Some of the key options include:

- `--proto_path=PATH`: Specify the directory to search for imports.
- `--version`: Show version info.
- `--encode=MESSAGE_TYPE`: Convert a text-format message of a given type from standard input to binary on standard output.
- `--decode=MESSAGE_TYPE`: Convert a binary message of a given type from standard input to text format on standard output.
- `--cpp_out=OUT_DIR`: Generate C++ header and source.
- `--java_out=OUT_DIR`: Generate Java source file.

... and many more. For a full list of options, run `./protoc --help` or refer to the official documentation.