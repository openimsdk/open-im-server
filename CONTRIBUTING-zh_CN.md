# 如何给OpenIM贡献代码（提交pull request）

<p align="center">
  <a href="./CONTRIBUTING.md">Englist</a> · 
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


本指南将以 [openimsdk/open-im-server](https://github.com/openimsdk/open-im-server)为例详细说明如何为  OpenIM 项目贡献代码。我们采用“一问题一分支”的策略，确保每个 Issue 都对应一个专门的分支，以便有效管理代码变更。

## 1. Fork 仓库
前往 [openimsdk/open-im-server](https://github.com/openimsdk/open-im-server) GitHub 页面，点击右上角的 "Fork" 按钮，将仓库 Fork 到你的 GitHub 账户下。

## 2. 克隆仓库
将你 Fork 的仓库克隆到本地：
```bash
git clone https://github.com/your-username/open-im-server.git
```

## 3. 设置远程上游
添加原始仓库为远程上游以便跟踪其更新：
```bash
git remote add upstream https://github.com/openimsdk/open-im-server.git
```

## 4. 创建 Issue
在原始仓库中创建一个新的 Issue，详细描述你遇到的问题或希望添加的新功能。

## 5. 创建新分支
基于主分支创建一个新分支，并使用描述性的名称与 Issue ID，例如：
```bash
git checkout -b fix-bug-123
```

## 6. 提交更改
在你的本地分支上进行更改后，提交这些更改：
```bash
git add .
git commit -m "Describe your changes in detail"
```

## 7. 推送分支
将你的分支推送回你的 GitHub Fork：
```bash
git push origin fix-bug-123
```

## 8. 创建 Pull Request
在 GitHub 上转到你的 Fork 仓库，点击 "Pull Request" 按钮。确保 PR 描述清楚，并链接到相关的 Issue。

## 9. 签署 CLA
如果这是你第一次提交 PR，你需要在 PR 的评论中回复：
```
I have read the CLA Document and I hereby sign the CLA
```
## 其他说明

如果需要将同一修改提交到两个不同的分支（例如 `main` 和 `release-v3.7`），应从对应的远程分支分别创建两个新分支。首先在一个分支上完成修改，然后使用 `cherry-pick` 命令将这些更改应用到另一个分支。之后，为每个分支独立提交 Pull Request。

