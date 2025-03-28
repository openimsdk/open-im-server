# Stress Test

## Usage

You need set `TestTargetUserList` and `DefaultGroupID` variables.

### Build

```bash
go build -o _output/bin/tools/linux/amd64/stress-test tools/stress-test/main.go

# or

go build -o tools/stress-test/stress-test tools/stress-test/main.go
```

### Excute

```bash
_output/bin/tools/linux/amd64/stress-test -c config/

#or

tools/stress-test/stress-test -c config/
```
