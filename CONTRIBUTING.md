# How to Contribute to OpenIM (Submitting Pull Requests)

<p align="center">
  <a href="./CONTRIBUTING.md">English</a> · 
  <a href="./CONTRIBUTING-zh_CN.md">中文</a> · 
  <a href="docs/contributing/CONTRIBUTING-UA.md">Українська</a> · 
  <a href="docs/contributing/CONTRIBUTING-CS.md">Česky</a> · 
  <a href="docs/contributing/CONTRIBUTING-HU.md">Magyar</a> · 
  <a href="docs/contributing/CONTRIBUTING-ES.md">Español</a> · 
  <a href="docs/contributing/CONTRIBUTING-FA.md">فارسی</a> · 
  <a href="docs/contributing/CONTRIBUTING-FR.md">Français</a> · 
  <a href="docs/contributing/CONTRIBUTING-DE.md">Deutsch</a> · 
  <a href="docs/contributing/CONTRIBUTING-PL.md">Polski</a> · 
  <a href="docs/contributing/CONTRIBUTING-ID.md">Indonesian</a> · 
  <a href="docs/contributing/CONTRIBUTING-FI.md">Suomi</a> · 
  <a href="docs/contributing/CONTRIBUTING-ML.md">മലയാളം</a> · 
  <a href="docs/contributing/CONTRIBUTING-JP.md">日本語</a> · 
  <a href="docs/contributing/CONTRIBUTING-NL.md">Nederlands</a> · 
  <a href="docs/contributing/CONTRIBUTING-IT.md">Italiano</a> · 
  <a href="docs/contributing/CONTRIBUTING-RU.md">Русский</a> · 
  <a href="docs/contributing/CONTRIBUTING-PTBR.md">Português (Brasil)</a> · 
  <a href="docs/contributing/CONTRIBUTING-EO.md">Esperanto</a> · 
  <a href="docs/contributing/CONTRIBUTING-KR.md">한국어</a> · 
  <a href="docs/contributing/CONTRIBUTING-AR.md">العربي</a> · 
  <a href="docs/contributing/CONTRIBUTING-VN.md">Tiếng Việt</a> · 
  <a href="docs/contributing/CONTRIBUTING-DA.md">Dansk</a> · 
  <a href="docs/contributing/CONTRIBUTING-GR.md">Ελληνικά</a> · 
  <a href="docs/contributing/CONTRIBUTING-TR.md">Türkçe</a>
</p>

This guide will use [openimsdk/open-im-server](https://github.com/openimsdk/open-im-server) as an example to explain in detail how to contribute code to the OpenIM project. We adopt a "one issue, one branch" strategy to ensure each issue corresponds to a dedicated branch for effective code change management.

### 1. Fork the Repository
Go to the [openimsdk/open-im-server](https://github.com/openimsdk/open-im-server) GitHub page, click the "Fork" button in the upper right corner to fork the repository to your GitHub account.

### 2. Clone the Repository
Clone the repository you forked to your local machine:
```bash
git clone https://github.com/your-username/open-im-server.git
```

### 3. Set Upstream Remote
Add the original repository as a remote upstream to track updates:
```bash
git remote add upstream https://github.com/openimsdk/open-im-server.git
```

### 4. Create an Issue
Create a new issue in the original repository detailing the problem you encountered or the new feature you wish to add.

### 5. Create a New Branch
Create a new branch off the main branch with a descriptive name and Issue ID, for example:
```bash
git checkout -b fix-bug-123
```

### 6. Commit Changes
After making changes on your local branch, commit these changes:
```bash
git add .
git commit -m "Describe your changes

 in detail"
```

### 7. Push the Branch
Push your branch back to your GitHub fork:
```bash
git push origin fix-bug-123
```

### 8. Create a Pull Request
Go to your fork on GitHub and click the "Pull Request" button. Ensure the PR description is clear and links to the related issue.

### 9. Sign the CLA
If this is your first time submitting a PR, you will need to reply in the comments of the PR:
```
I have read the CLA Document and I hereby sign the CLA
```

### Programming Standards
Please refer to the following documents for detailed information on Go language programming standards:
- [Go Coding Standards](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/go-code.md)
- [Code Conventions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/code-conventions.md)

### Logging Standards
- **Do not use the standard `log` package**.
- Use the `"github.com/openimsdk/tools/log"` package for logging, which supports multiple log levels: `debug`, `info`, `warn`, `error`.
- **Error logs should only be printed in the function where they are first actively called** to prevent log duplication and ensure clear error context.

### Exception and Error Handling
- **Prohibit the use of `panic`**: The code should not use `panic` to avoid abrupt termination when encountering unrecoverable errors.
- **Error Wrapping**: Use `"github.com/openimsdk/tools/errs"` to wrap errors, maintaining the integrity of error information and facilitating debugging.
- **Error Propagation**: If a function cannot handle an error itself, it should return the error to the caller, rather than hiding or ignoring it.
