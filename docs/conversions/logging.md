##  Log Standards

### Log Standards

- The unified log package `github.com/openimsdk/open-im-server/internal/pkg/log` should be used for all logging;
- Use structured logging formats: `log.Infow`, `log.Warnw`, `log.Errorw`, etc. For example: `log.Infow("Update post function called")`;
- All logs should start with an uppercase letter and should not end with a `.`. For example: `log.Infow("Update post function called")`;
- Use past tense. For example, use `Could not delete B` instead of `Cannot delete B`;
- Adhere to log level standards:
  - Debug level logs use `log.Debugw`;
  - Info level logs use `log.Infow`;
  - Warning level logs use `log.Warnw`;
  - Error level logs use `log.Errorw`;
  - Panic level logs use `log.Panicw`;
  - Fatal level logs use `log.Fatalw`.
- Log settings:
  - Development and test environments: The log level is set to `debug`, the log format can be set to `console` / `json` as needed, and caller is enabled;
  - Production environment: The log level is set to `info`, the log format is set to `json`, and caller is enabled. (**Note**: In the early stages of going online, to facilitate troubleshooting, the log level can be set to `debug`)
- When logging, avoid outputting sensitive information, such as passwords, keys, etc.
- If you are calling a logging function in a function/method with a `context.Context` parameter, it is recommended to use `log.L(ctx).Infow()` for logging.