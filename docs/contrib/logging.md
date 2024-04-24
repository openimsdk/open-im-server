# OpenIM Logging and Error Handling Documentation

## Script Logging Documentation Link

If you wish to view the script's logging documentation, you can click on this link: [Logging Documentation](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/bash-log.md).

Below is the documentation for logging and error handling in the OpenIM Go project.

To create a standard set of documentation that is quick to read and easy to understand, we will highlight key information about the `Logger` interface and the `CodeError` interface. This includes the purpose of each interface, key methods, and their use cases. This will help developers quickly grasp how to effectively use logging and error handling within the project.

## Logging (`Logger` Interface)

### Purpose
The `Logger` interface aims to provide the OpenIM project with a unified and flexible logging mechanism, supporting structured logging formats for efficient log management and analysis.

### Key Methods

- **Debug, Info, Warn, Error**  
  Log messages of different levels to suit various logging needs and scenarios. These methods accept a context (`context.Context`), a message (`string`), and key-value pairs (`...interface{}`), allowing the log to carry rich context information.

- **WithValues**  
  Append key-value pair information to log messages, returning a new `Logger` instance. This helps in adding consistent context information.

- **WithName**  
  Set the name of the logger, which helps in identifying the source of the logs.

- **WithCallDepth**  
  Adjust the call stack depth to accurately identify the source of the log message.

### Use Cases

- Developers should choose the appropriate logging level (such as `Debug`, `Info`, `Warn`, `Error`) based on the importance of the information when logging.
- Use `WithValues` and `WithName` to add richer context information to logs, facilitating subsequent tracking and analysis.

## Error Handling (`CodeError` Interface)

### Purpose
The `CodeError` interface is designed to provide a unified mechanism for error handling and wrapping, making error information more detailed and manageable.

### Key Methods

- **Code**  
  Return the error code to distinguish between different types of errors.

- **Msg**  
  Return the error message description to display to the user.

- **Detail**  
  Return detailed information about the error for further debugging by developers.

- **WithDetail**  
  Add detailed information to the error, returning a new `CodeError` instance.

- **Is**  
  Determine whether the current error matches a specified error, supporting a flexible error comparison mechanism.

- **Wrap**  
  Wrap another error with additional message description, facilitating the tracing of the error's cause.

### Use Cases

- When defining errors with specific codes and messages, use error types that implement the `CodeError` interface.
- Use `WithDetail` to add additional context information to errors for more accurate problem localization.
- Use the `Is` method to judge the type of error for conditional branching.
- Use the `Wrap` method to wrap underlying errors while adding more contextual descriptions.

## Logging Standards and Code Examples

In the OpenIM project, we use the unified logging package `github.com/openimsdk/tools/log` for logging to achieve efficient log management and analysis. This logging package supports structured logging formats, making it easier for developers to handle log information.

### Logger Interface and Implementation

The logger interface is defined as follows:

```go
type Logger interface {
    Debug(ctx context.Context, msg string, keysAndValues ...interface{})
    Info(ctx context.Context, msg string, keysAndValues ...interface{})
    Warn(ctx context.Context, msg string, err error, keysAndValues ...interface{})
    Error(ctx context.Context, msg string, err error, keysAndValues ...interface{})
    WithValues(keysAndValues ...interface{}) Logger
    WithName(name string) Logger
    WithCallDepth(depth int) Logger
}
```

Example code: Using the `Logger` interface to log at the info level.

```go
func main() {
	logger := log.NewLogger().WithName("MyService")
	ctx := context.Background()
	logger.Info(ctx, "Service started", "port", "8080")
}
```

## Error Handling and Code Examples

We use the `github.com/openimsdk/tools/errs` package for unified error handling and wrapping.

### CodeError Interface and Implementation

The error interface is defined as follows:

```go
type CodeError interface {
    Code() int
    Msg() string
    Detail() string
    WithDetail(detail string) CodeError
    Is(err error, loose ...bool) bool
    Wrap(msg ...string) error
    error
}
```

Example code: Creating and using the `CodeError` interface to handle errors.

```go
package main

import (
	"fmt"
	"github.com/openimsdk/tools/errs"
)

func main() {
	err := errs.New(404, "Resource not found")
	err = err.WithDetail("

More details")
	if e, ok := err.(errs.CodeError); ok {
		fmt.Println(e.Code(), e.Msg(), e.Detail())
	}
}
```

### Detailed Logging Standards and Code Examples

1. **Print key information at startup**  
   It is crucial to print entry parameters and key process information at program startup. This helps understand the startup state and configuration of the program.

   **Code Example**:
   ```go
   package main

   import (
       "fmt"
       "os"
   )

   func main() {
       fmt.Println("Program startup, version: 1.0.0")
       fmt.Printf("Connecting to database: %s\n", os.Getenv("DATABASE_URL"))
   }
   ```

2. **Use `tools/log` and `fmt` for logging**  
   Logging should be done using a specialized logging library for unified management and formatted log output.

   **Code Example**: Logging an info level message with `tools/log`.
   ```go
   package main

   import (
       "context"
       "github.com/openimsdk/tools/log"
   )

   func main() {
       ctx := context.Background()
       log.Info(ctx, "Application started successfully")
   }
   ```

3. **Use standard error output for startup failures or critical information**  
   Critical error messages or program startup failures should be indicated to the user through standard error output.

   **Code Example**:
   ```go
   package main

   import (
       "fmt"
       "os"
   )

   func checkEnvironment() bool {
       return os.Getenv("REQUIRED_ENV") != ""
   }

   func main() {
       if !checkEnvironment() {
           fmt.Fprintln(os.Stderr, "Missing required environment variable")
           os.Exit(1)
       }
   }
   ```
   
   We encapsulate it into separate tools, which can output error information through the `tools/log` package.

   ```go
    package main

    import (
        util "github.com/openimsdk/open-im-server/v3/pkg/util/genutil"
    )

    func main() {
        if err := apiCmd.Execute(); err != nil {
            util.ExitWithError(err)
        }
    }
    ```

4. **Use `tools/log` package for runtime logging**  
   This ensures consistency and control over logging.

   **Code Example**: Same as the above example using `tools/log`. When `tools/log` is not initialized, consider using `fmt` for standard output.

5. **Error logs should be printed by the top-level caller**  
   This is to avoid duplicate logging of errors, typically errors are caught and logged at the application's outermost level.

   **Code Example**:
   ```go
   package main

   import (
       "github.com/openimsdk/tools/log"
       "context"
   )

   func doSomething() error {
       // An error occurs here
       return errs.Wrap(errors.New("An error occurred"))
   }

   func controller() error {
       err := doSomething()
       if err != nil {
           return err
       }
       return nil
   }

   func main() {
       err := controller()
       if err != nil {
           log.Error(context.Background(), "Operation failed", err)
       }
   }
   ```

6. **Handling logs for API RPC calls and non-RPC applications**

   For API RPC calls using gRPC, logs at the information level are printed by middleware on the gRPC server side, reducing the need to manually log in each RPC method. For non-RPC applications, it's recommended to manually log key execution paths to track the application's execution flow.

    **gRPC Server-Side Logging Middleware:**

    In gRPC, `UnaryInterceptor` and `StreamInterceptor` can intercept Unary and Stream type RPC calls, respectively. Here's an example of how to implement a simple Unary RPC logging middleware:

    ```go
    package main

    import (
        "context"
        "google.golang.org/grpc"
        "google.golang.org/grpc/codes"
        "google.golang.org/grpc/status"
        "log"
        "time"
    )

    // unaryServerInterceptor returns a new unary server interceptor that logs each request.
    func unaryServerInterceptor() grpc.UnaryServerInterceptor {
        return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
            // Record the start time of the request
            start := time.Now()
            // Call the actual RPC method
            resp, err = handler(ctx, req)
            // After the request ends, log the duration and other information
            log.Printf("Request method: %s, duration: %s, error status: %v", info.FullMethod, time.Since(start), status.Code(err))
            return resp, err
        }
    }

    func main() {
        // Create a gRPC server and add the middleware
        s := grpc.NewServer

(grpc.UnaryInterceptor(unaryServerInterceptor()))
        // Register your service

        // Start the gRPC server
        log.Println("Starting gRPC server...")
        // ...
    }
    ```

    **Logging for Non-RPC Applications:**

    For non-RPC applications, the key is to log at appropriate places in the code to maintain an execution trace. Here's a simple example showing how to log when handling a task:

    ```go
    package main

    import (
        "log"
    )

    func processTask(taskID string) {
        // Log the start of task processing
        log.Printf("Starting task processing: %s", taskID)
        // Suppose this is where the task is processed

        // Log after the task is completed
        log.Printf("Task processing completed: %s", taskID)
    }

    func main() {
        // Example task ID
        taskID := "task123"
        processTask(taskID)
    }
    ```

    In both scenarios, appropriate logging can help developers and operators monitor the health of the system, trace the source of issues, and quickly locate and resolve problems. For gRPC logging, using middleware can effectively centralize log management and control. For non-RPC applications, ensuring logs are placed at critical execution points can help understand the program's operational flow and state changes.

### When to Wrap Errors?

1. **Wrap errors generated within the function**  
   When an error occurs within a function, use `errs.Wrap` to add context information to the original error.

   **Code Example**:
   ```go
   func doSomething() error {
       // Suppose an error occurs here
       err, _ := someFunc()
       if err != nil {
         return errs.WrapMsg(err, "doSomething failed")
       }
   }
   ```

   It just works if the package is wrong:

   ```go
      func doSomething() error {
       // Suppose an error occurs here
       err, _ := someFunc()
       if err != nil {
         return errs.Wrap(err)
       }
   }
   ```

2. **Wrap errors from system calls or other packages**  
   When calling external libraries or system functions that return errors, also add context information to wrap the error.

   **Code Example**:
   ```go
   func readConfig(file string) error {
       _, err := os.ReadFile(file)
       if err != nil {
           return errs.Wrap(err, "Failed to read config file")
       }
       return nil
   }
   ```

3. **No need to re-wrap errors for internal module calls**  

   If an error has been appropriately wrapped with sufficient context information in an internal module call, there's no need to wrap it again.

    **Code Example**:
    ```go
    func doSomething() error {
        err := doAnotherThing()
        if err != nil {
            return err
        }
        return nil
    }
    ```

4. **Ensure comprehensive wrapping of errors with detailed messages**
   When wrapping errors, ensure to provide ample context information to make the error more understandable and easier to debug.

   **Code Example**:
   ```go
   func connectDatabase() error {
       err := db.Connect(config.DatabaseURL)
       if err != nil {
           return errs.Wrap(err, fmt.Sprintf("Failed to connect to database, URL: %s", config.DatabaseURL))
       }
       return nil
   }
   ```


### About WrapMsg Use

```go
// 	"github.com/openimsdk/tools/errs"
func WrapMsg(err error, msg string, kv ...any) error {
	if len(kv) == 0 {
		if len(msg) == 0 {
			return errors.WithStack(err)
		} else {
			return errors.WithMessage(err, msg)
		}
	}
	var buf bytes.Buffer
	if len(msg) > 0 {
		buf.WriteString(msg)
		buf.WriteString(" ")
	}
	for i := 0; i < len(kv); i += 2 {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(toString(kv[i]))
		buf.WriteString("=")
		buf.WriteString(toString(kv[i+1]))
	}
	return errors.WithMessage(err, buf.String())
}
```

1. **Function Signature**:
   - `err error`: The original error object.
   - `msg string`: The message text to append to the error.
   - `kv ...any`: A variable number of parameters used to pass key-value pairs. `any` was introduced in Go 1.18 and is equivalent to `interface{}`, meaning any type.

2. **Logic**:
   - If there are no key-value pairs (`kv` is empty):
     - If `msg` is also empty, use `errors.WithStack(err)` to return the original error with the call stack appended.
     - If `msg` is not empty, use `errors.WithMessage(err, msg)` to append the message text to the original error.
   - If there are key-value pairs, the function constructs a string containing the message text and all key-value pairs. The key-value pairs are added in the format `"key=value"`, separated by commas. If a message text is provided, it is added first, followed by a space.

3. **Key-Value Pair Formatting**:
   - A loop iterates over all the key-value pairs, processing one pair at a time.
   - The `toString` function (although not provided in the code, we can assume it converts any type to a string) is used to convert both keys and values to strings, and they are added to a `bytes.Buffer` in the format `"key=value"`.

4. **Result**:
   - Use `errors.WithMessage(err, buf.String())` to append the constructed message text to the original error, and return this new error object.

Next, let's demonstrate several ways to use the `WrapMsg` function:

**Example 1: No Additional Information**

```go
// "github.com/openimsdk/tools/errs"
err := errors.New("original error")
wrappedErr := errs.WrapMsg(err, "")
// wrappedErr will contain the original error and its call stack
```

**Example 2: Message Text Only**

```go
// "github.com/openimsdk/tools/errs"
err := errors.New("original error")
wrappedErr := errs.WrapMsg(err, "additional error information")
// wrappedErr will contain the original error, call stack, and "additional error information"
```

**Example 3: Message Text and Key-Value Pairs**

```go
// "github.com/openimsdk/tools/errs"
err := errors.New("original error")
wrappedErr := errs.WrapMsg(err, "problem occurred", "code", 404, "url", "webhook://example.com")
// wrappedErr will contain the original error, call stack, and "problem occurred code=404, url=http://example.com"
```

**Example 4: Key-Value Pairs Only**

```go
// "github.com/openimsdk/tools/errs"
err := errors.New("original error")
wrappedErr := errs.WrapMsg(err, "", "user", "john_doe", "action", "login")
// wrappedErr will contain the original error, call stack, and "user=john_doe, action=login"
```

> [!TIP] WThese examples demonstrate how the `errs.WrapMsg` function can flexibly handle error messages and context data, helping developers to more effectively track and debug their programs.


### Example 5: Dynamic Key-Value Pairs from Context
Suppose we have some runtime context variables, such as a user ID and the type of operation being performed, and we want to include these variables in the error message. This can help with later debugging and identifying the specific environment of the issue.

```go
// Define some context variables
userID := "user123"
operation := "update profile"
errorCode := 500
requestURL := "webhook://example.com/updateProfile"

// Create a new error
err := errors.New("original error")

// Wrap the error, including dynamic key-value pairs from the context
wrappedErr := errs.WrapMsg(err, "operation failed", "user", userID, "action", operation, "code", errorCode, "url", requestURL)
// wrappedErr will contain the original error, call stack, and "operation failed user=user123, action=update profile, code=500, url=http://example.com/updateProfile"
```

> [!TIP]In this example, the `WrapMsg` function accepts not just a static error message and additional information, but also dynamic key-value pairs generated from the code's execution context, such as the user ID, operation type, error code, and the URL of the request. Including this contextual information in the error message makes it easier for developers to understand and resolve the issue.