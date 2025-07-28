# Stress Test

## Usage

You need set `TestTargetUserList` and `DefaultGroupID` variables.

### Build

```bash

go build -o test/stress-test/stress-test test/stress-test/main.go
```

### Excute

```bash

tools/stress-test/stress-test -c config/
```
