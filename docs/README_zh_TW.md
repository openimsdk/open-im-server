# OpenIM 伺服器文檔

歡迎來到 OpenIM 文件中心！ 該中心提供全面的指南和手冊，旨在幫助您充分利用 OpenIM 體驗。

## 目錄

1. [Contrib](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib) - 開發人員貢獻和配置指南
2. [Conversions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib) - 編碼約定、日誌記錄策略和其他轉換工具

------

## 貢獻

本節為開發人員提供了有關如何貢獻程式碼、設定環境以及遵循相關流程的詳細指南。

- [Code Conventions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/code-conventions.md) - 在 OpenIM 中編寫程式碼的規則和約定。
- [Development Guide](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/development.md) - 有關如何在 OpenIM 中進行開發的指南。
- [Git Cherry Pick](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/gitcherry-pick.md) - 精挑細選操作指南。
- [Git Workflow](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/git-workflow.md) - OpenIM 中的 git 工作流程。
- [Initialization Configurations](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/init-config.md) - 設定和初始化 OpenIM 的指南。
- [Docker Installation](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/install-docker.md) - 如何在您的電腦上安裝 Docker。
- [Linux Development Environment](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/linux-development.md) - Linux 上的開發環境設定指南。
- [Local Actions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/local-actions.md) - 關於如何在當地進行某些共同行動的指南。
- [Offline Deployment](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/offline-deployment.md) - 離線部署OpenIM的方法。
- [Protoc Tools](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/protoc-tools.md) - 協議工具使用指南。
- [Go Tools](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/util-go.md) - OpenIM 在 Go 中的工具和函式庫。
- [Makefile Tools](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/util-makefile.md) - Makefile 的最佳實務和工具。
- [Script Tools](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/util-scripts.md) - 腳本的最佳實踐和工具。

## 轉換

本節介紹 OpenIM 中的各種約定和策略，包括程式碼、日誌、版本等。

- [API Conversions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/api.md) - API 轉換的指南和方法。
- [Logging Policy](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/bash-log.md) - OpenIM 中的日誌記錄策略和約定。
- [CI/CD Actions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/cicd-actions.md) - CI/CD 的程序和約定。
- [Commit Conventions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/commit.md) - OpenIM 中程式碼提交的約定。
- [Directory Conventions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/directory.md) - OpenIM 中的目錄結構和約定。
- [Error Codes](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/error-code.md) - 錯誤代碼的清單和描述。
- [Go Code Conversions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/go-code.md) - Go 程式碼的約定和轉換。
- [Docker Image Strategy](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/images.md) - OpenIM Docker 映像的管理策略，跨越多種架構和映像儲存庫。
- [Logging Conventions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/logging.md) - 有關日誌記錄的更詳細約定。
- [Version Conventions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/version.md) - OpenIM 版本的命名與管理策略。


## 對於開發者、貢獻者和社區維護者

### 開發者和貢獻者

如果您是開發人員或熱衷於做出貢獻的人：

- 熟悉我們的 [Code Conventions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/code-conventions.md) 和 [Git Workflow](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/git-workflow.md) 以確保順利貢獻。
- 深入閱讀 [Development Guide](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/development.md) ，掌握 OpenIM 的開發實務。

### 社區維護者

作為社區維護者：

- 確保貢獻符合我們文件中概述的標準。
- 定期查看 [Logging Policy](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/bash-log.md) 和 [Error Codes](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/error-code.md) 以保持更新。

## 對於用戶

使用者應特別注意：

- [Docker Installation](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/install-docker.md) - 如果您打算使用 OpenIM 的 Docker 映像，則這是必要的。
- [Docker Image Strategy](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/images.md) - 了解可用的不同影像以及如何為您的架構選擇正確的影像。