## Error Code Standards

Error codes are one of the important means for users to locate and solve problems. When an application encounters an exception, users can quickly locate and resolve the problem based on the error code and the description and solution of the error code in the documentation.

### Error Code Naming Standards

- Follow CamelCase notation;
- Error codes are divided into two levels. For example, `InvalidParameter.BindError`, separated by a `.`. The first-level error code is platform-level, and the second-level error code is resource-level, which can be customized according to the scenario;
- The second-level error code can only use English letters or numbers ([a-zA-Z0-9]), and should use standard English word spelling, standard abbreviations, RFC term abbreviations, etc.;
- The error code should avoid multiple definitions of the same semantics, for example: `InvalidParameter.ErrorBind`, `InvalidParameter.BindError`.

### First-Level Common Error Codes

| Error Code       | Error Description                                            | Error Type |
| ---------------- | ------------------------------------------------------------ | ---------- |
| InternalError    | Internal error                                               | 1          |
| InvalidParameter | Parameter error (including errors in parameter type, format, value, etc.) | 0          |
| AuthFailure      | Authentication / Authorization error                         | 0          |
| ResourceNotFound | Resource does not exist                                      | 0          |
| FailedOperation  | Operation failed                                             | 2          |

> Error Type: 0 represents the client, 1 represents the server, 2 represents both the client / server.