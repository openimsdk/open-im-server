# Notes about go workspace

As openim is using go1.18's [workspace feature](https://go.dev/doc/tutorial/workspaces), once you add a new module, you need to run `go work use -r .` at root directory to update the workspace synced.

### Create a new extensions

1. Create your tools_name directory in pkg `/tools` first and cd into it.
2. Init the project.
3. Then `go work use -r .` at current directory to update the workspace.
4. Create your tools

You can execute the following commands to do things above:

```bash
# edit the CRD_NAME and CRD_GROUP to your own
export OPENIM_TOOLS_NAME=<Changeme>

# copy and paste to create a new CRD and Controller
mkdir tools/${OPENIM_TOOLS_NAME}
cd tools/${OPENIM_TOOLS_NAME}
go mod init github.com/openimsdk/open-im-server/tools/${OPENIM_TOOLS_NAME}
go mod tidy
go work use -r .
cd ../..
```