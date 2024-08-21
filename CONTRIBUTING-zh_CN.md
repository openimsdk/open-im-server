

# 如何给 OpenIM 贡献代码（提交 Pull Request）

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

本指南将以 [openimsdk/open-im-server](https://github.com/openimsdk/open-im-server) 为例，详细说明如何为 OpenIM 项目贡献代码。我们采用“一问题一分支”的策略，确保每个 Issue 都对应一个专门的分支，以便有效管理代码变更。

### 1. Fork 仓库
前往 [openimsdk/open-im-server](https://github.com/openimsdk/open-im-server) GitHub 页面，点击右上角的 "Fork" 按钮，将仓库 Fork 到你的 GitHub 账户下。

### 2. 克隆仓库
将你 Fork 的仓库克隆到本地：
```bash
git clone https://github.com/your-username/open-im-server.git
```

### 3. 设置远程上游
添加原始仓库为远程上游以便跟踪其更新：
```bash
git remote add upstream https://github.com/openimsdk/open-im-server.git
```

### 4. 创建 Issue
在原始仓库中创建一个新的 Issue，详细描述你遇到的问题或希望添加

的新功能。

### 5. 创建新分支
基于主分支创建一个新分支，并使用描述性的名称与 Issue ID，例如：
```bash
git checkout -b fix-bug-123
```

### 6. 提交更改
在你的本地分支上进行更改后，提交这些更改：
```bash
git add .
git commit -m "Describe your changes in detail"
```

### 7. 推送分支
将你的分支推送回你的 GitHub Fork：
```bash
git push origin fix-bug-123
```

### 8. 创建 Pull Request
在 GitHub 上转到你的 Fork 仓库，点击 "Pull Request" 按钮。确保 PR 描述清楚，并链接到相关的 Issue。

### 9. 签署 CLA
如果这是你第一次提交 PR，你需要在 PR 的评论中回复：
```
I have read the CLA Document and I hereby sign the CLA
```

### 编程规范
请参考以下文档以了解关于 Go 语言编程规范的详细信息：
- [Go 编码规范](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/go-code.md)
- [代码约定](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/code-conventions.md)

### 日志规范
- **禁止使用标准的 `log` 包**。
- 应使用 `"github.com/openimsdk/tools/log"` 包来打印日志，该包支持多种日志级别：`debug`、`info`、`warn`、`error`。
- **错误日志应仅在首次调用的函数中打印**，以防止日志重复，并确保错误的上下文清晰。

### 异常及错误处理
- **禁止使用 `panic`**：程序中不应使用 `panic`，以避免在遇到不可恢复的错误时突然终止。
- **错误包裹**：使用 `"github.com/openimsdk/tools/errs"` 来包裹错误，保持错误信息的完整性并增加调试便利。
- **错误传递**：如果函数本身不能处理错误，应将错误返回给调用者，而不是隐藏或忽略这些错误。
