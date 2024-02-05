<p align="center">
    <a href="https://openim.io">
        <img src="./assets/logo-gif/openim-logo.gif" width="60%" height="30%"/>
    </a>
</p>

<div align="center">

[![Stars](https://img.shields.io/github/stars/openimsdk/open-im-server?style=for-the-badge&logo=github&colorB=ff69b4)](https://github.com/openimsdk/open-im-server/stargazers)
[![Forks](https://img.shields.io/github/forks/openimsdk/open-im-server?style=for-the-badge&logo=github&colorB=blue)](https://github.com/openimsdk/open-im-server/network/members)
[![Codecov](https://img.shields.io/codecov/c/github/openimsdk/open-im-server?style=for-the-badge&logo=codecov&colorB=orange)](https://app.codecov.io/gh/openimsdk/open-im-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/openimsdk/open-im-server?style=for-the-badge)](https://goreportcard.com/report/github.com/openimsdk/open-im-server)
[![Go Reference](https://img.shields.io/badge/Go%20Reference-blue.svg?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/openimsdk/open-im-server/v3)
[![License](https://img.shields.io/badge/license-Apache--2.0-green?style=for-the-badge)](https://github.com/openimsdk/open-im-server/blob/main/LICENSE)
[![Slack](https://img.shields.io/badge/Slack-500%2B-blueviolet?style=for-the-badge&logo=slack&logoColor=white)](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q)
[![Best Practices](https://img.shields.io/badge/Best%20Practices-purple?style=for-the-badge)](https://www.bestpractices.dev/projects/8045)
[![Good First Issues](https://img.shields.io/github/issues/openimsdk/open-im-server/good%20first%20issue?style=for-the-badge&logo=github)](https://github.com/openimsdk/open-im-server/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)


<p align="center">
  <a href="./README.md">Englist</a> · 
  <a href="./README-zh_CN.md">中文</a> · 
  <a href="docs/readme/README-UA.md">Українська</a> · 
  <a href="docs/readme/README-CS.md">Česky</a> · 
  <a href="docs/readme/README-HU.md">Magyar</a> · 
  <a href="docs/readme/README-ES.md">Español</a> · 
  <a href="docs/readme/README-FA.md">فارسی</a> · 
  <a href="docs/readme/README-FR.md">Français</a> · 
  <a href="docs/readme/README-DE.md">Deutsch</a> · 
  <a href="docs/readme/README-PL.md">Polski</a> · 
  <a href="docs/readme/README-ID.md">Indonesian</a> · 
  <a href="docs/readme/README-FI.md">Suomi</a> · 
  <a href="docs/readme/README-ML.md">മലയാളം</a> · 
  <a href="docs/readme/README-JP.md">日本語</a> · 
  <a href="docs/readme/README-NL.md">Nederlands</a> · 
  <a href="docs/readme/README-IT.md">Italiano</a> · 
  <a href="docs/readme/README-RU.md">Русский</a> · 
  <a href="docs/readme/README-PTBR.md">Português (Brasil)</a> · 
  <a href="docs/readme/README-EO.md">Esperanto</a> · 
  <a href="docs/readme/README-KR.md">한국어</a> · 
  <a href="docs/readme/README-AR.md">العربي</a> · 
  <a href="docs/readme/README-VN.md">Tiếng Việt</a> · 
  <a href="docs/readme/README-DA.md">Dansk</a> · 
  <a href="docs/readme/README-GR.md">Ελληνικά</a> · 
  <a href="docs/readme/README-TR.md">Türkçe</a>
</p>


</div>

</p>

## 🟢 WeChatをスキャンしてグループ交流に参加
<img src="./docs/images/Wechat.jpg" width="300">


## Ⓜ️ OpenIMについて

OpenIMは、アプリケーション内でチャット、音声通話、通知、AIチャットボットなどの通信機能を統合するために特別に設計されたサービスプラットフォームです。一連の強力なAPIとWebhooksを提供することで、開発者はアプリケーションに簡単にこれらの通信機能を統合できます。OpenIM自体は独立したチャットアプリではなく、アプリケーションにサポートを提供し、豊富な通信機能を実現するプラットフォームです。以下の図は、AppServer、AppClient、OpenIMServer、OpenIMSDK間の相互作用を示しています。



![App-OpenIMの関係](./docs/images/oepnim-design.png)

## 🚀 OpenIMSDKについて

 **OpenIMSDK**は、**OpenIMServer**用に設計されたIM SDKで、クライアントアプリケーションに組み込むためのものです。主な機能とモジュールは以下の通りです：

+ 🌟 主な機能：

  - 📦  ローカルストレージ
  - 🔔 リスナーコールバック
  - 🛡️ APIのラッピング
  - 🌐 接続管理

  ## 📚 主なモジュール：

  1. 🚀 初初期化とログイン
  2. 👤 ユーザー管理
  3. 👫 友達管理
  4. 🤖 グループ機能
  5. 💬 会話処理

Golangを使用して構築され、クロスプラットフォームの導入をサポートし、すべてのプラットフォームで一貫したアクセス体験を提供します。

👉 **[GO SDKを探索する](https://github.com/openimsdk/openim-sdk-core)**

## 🌐 OpenIMServerについて

+ **OpenIMServer** には以下の特徴があります：
  - 🌐 マイクロサービスアーキテクチャ：クラスターモードをサポートし、ゲートウェイ（gateway）と複数のrpcサービスを含みます。
  - 🚀 多様なデプロイメント方法：ソースコード、kubernetes、またはdockerでのデプロイメントをサポートします。
  - 海量ユーザーサポート：十万人規模の超大型グループ、千万人のユーザー、および百億のメッセージ

### 強化されたビジネス機能：

+ **REST API**：OpenIMServerは、ビジネスシステム用のREST APIを提供しており、ビジネスにさらに多くの機能を提供することを目指しています。たとえば、バックエンドインターフェースを通じてグループを作成したり、プッシュメッセージを送信したりするなどです。
+ **Webhooks**：OpenIMServerは、より多くのビジネス形態を拡張するためのコールバック機能を提供しています。コールバックとは、特定のイベントが発生する前後に、OpenIMServerがビジネスサーバーにリクエストを送信することを意味します。例えば、メッセージ送信の前後のコールバックなどです。

👉 **[もっと詳しく知る](https://docs.openim.io/guides/introduction/product)**

## :rocket: クイックスタート

iOS/Android/H5/PC/Webでのオンライン体験：

👉 **[OpenIM online demo](https://www.openim.io/zh/commercial)**

🤲 ユーザー体験を容易にするために、私たちは様々なデプロイメントソリューションを提供しています。以下のリストから、ご自身のデプロイメント方法を選択できます：

+ **[ソースコードデプロイメントガイド](https://docs.openim.io/guides/gettingStarted/imSourceCodeDeployment)**
+ **[Docker デプロイメントガイド](https://docs.openim.io/guides/gettingStarted/dockerCompose)**
+ **[Kubernetes デプロイメントガイド](https://docs.openim.io/guides/gettingStarted/k8s-deployment)**

## :hammer_and_wrench: OpenIMの開発を始める

OpenIMの目標は、トップレベルのオープンソースコミュニティを構築することです。私たちは一連の標準を持っており、それらは[コミュニティのリポジトリ](https://github.com/OpenIMSDK/community)にあります。

このOpen-IM-Serverリポジトリに貢献したい場合は、私たちの[貢献者ドキュメント]を読んでください。(https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md)。

始める前に、変更に必要があることを確認してください。最良の方法は、[新しいディスカッション](https://github.com/openimsdk/open-im-server/discussions/new/choose)や[Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q)での通信を作成すること、または問題を発見した場合は、まずそれを[報告](https://github.com/openimsdk/open-im-server/issues/new/choose)することです。

+ [コード標準](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/go-code.md)

+ [Docker イメージ標準](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/images.md)

+ [ディレクトリ標準](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/directory.md)

+ [提出標準](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/commit.md)

+ [バージョン管理標準](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/version.md)

+ [インターフェース標準](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/interface.md)

+ [OpenIMの設定と環境変数の設定](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/environment.md)

> **Note**
>中国のユーザー向けに、私たちの[Dockerイメージ標準](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/images.md)を読んで、国内のaliyunのイメージアドレスを使用してください。OpenIMは、中国向けのgitee同期リポジトリも持っています。[gitee.com](https://gitee.com/openimsdk)でそれを見つけることができます。 

## :link: リンク

  + **[完全なドキュメント](https://doc.rentsoft.cn/)**
  + **[更新ログ](https://github.com/openimsdk/open-im-server/blob/main/CHANGELOG.md)**
  + **[FAQ](https://github.com/openimsdk/open-im-server/blob/main/FAQ.md)**
  + **[コード例](https://github.com/openimsdk/open-im-server/blob/main/examples)**

## :handshake: 社コミュニティ

  + **[GitHub Discussions](https://github.com/openimsdk/open-im-server/discussions)**
  + **[Slack 通信](https://join.slack.com/t/openimsdk/shared_invite/zt-22720d66b-o_FvKxMTGXtcnnnHiMqe9Q)**
  + **[GitHub Issues](https://github.com/openimsdk/open-im-server/issues)**

  これらのプラットフォームに参加して、問題を議論したり、提案をしたり、成功のストーリーを共有したりできます！

## :writing_hand: 貢献

  あらゆる形の貢献を歓迎します！Pull Requestを提出する前に、私たちの[貢献者ドキュメント](https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md)を必ず読んでください。
  

  + **[バグを報告する](https://github.com/openimsdk/open-im-server/issues/new?assignees=&labels=bug&template=bug_report.md&title=)**
  + **[新機能を提案する](https://github.com/openimsdk/open-im-server/issues/new?assignees=&labels=enhancement&template=feature_request.md&title=)**
  + **[Pull Requestを提出する](https://github.com/openimsdk/open-im-server/pulls)**

  あなたの貢献に感謝します。一緒に強力なインスタントメッセージングソリューションを構築しましょう！

## :closed_book: ライセンス

  OpenIMSDKは、Apache License 2.0の下で利用可能です。[LICENSEファイル](https://github.com/openimsdk/open-im-server/blob/main/LICENSE)を見て、詳細を確認してください。

## 🔮 Thanks to our contributors!

<a href="https://github.com/openimsdk/open-im-server/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=openimsdk/open-im-server" />
</a>
