# How do I contribute code to OpenIM

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
This guide will explain in detail how to contribute code to the OpenIM project, using `openimsdk/open-im-server` as an example. We adopt a "one issue, one branch" strategy to ensure each issue corresponds to a dedicated branch, allowing for effective management of code changes.

### 1. Fork the Repository
Go to the `openimsdk/open-im-server` GitHub page, click the "Fork" button in the upper right corner to fork the repository to your GitHub account.

### 2. Clone the Repository
Clone the forked repository to your local machine:
```bash
git clone https://github.com/your-username/open-im-server.git
```

### 3. Set Upstream Remote
Add the original repository as a remote upstream to track updates:
```bash
git remote add upstream https://github.com/openimsdk/open-im-server.git
```

### 4. Create an Issue
Create a new issue in the original repository describing the problem you are facing or the new feature you want to add. For significant feature adjustments, propose an RFC issue to facilitate broad discussion and participation from community members.

### 5. Create a New Branch
Create a new branch based on the main branch and name it descriptively, including the Issue ID, for example:
```bash
git checkout -b fix-bug-123
```

### 6. Commit Changes
After making changes on your local branch, commit them:
```bash
git add .
git commit -m "Describe your changes in detail"
```

### 7. Push the Branch
Push your branch back to your GitHub fork:
```bash
git push origin fix-bug-123
```

### 8. Create a Pull Request
Go to your fork on GitHub, click the "Pull Request" button. Make sure the PR description is clear and links to the related Issue.
#### 🅰 Fixed issue #issueID

### 9. Sign the CLA
If this is your first time submitting a PR, you need to reply in the PR comments:
```
I have read the CLA Document and I hereby sign the CLA
```

### Additional Notes
If the same modification needs to be submitted to two different branches (e.g., main and release-v3.7), create two new branches from the corresponding remote branches. First complete the modification in one branch, then use the `cherry-pick` command to apply these changes to the other branch. After that, submit a separate Pull Request for each branch.

