<p align="center">
    <a href="https://openim.io">
        <img src="../../assets/logo-gif/openim-logo.gif" width="60%" height="30%"/>
    </a>
</p>

<div align="center">

[![Stars](https://img.shields.io/github/stars/openimsdk/open-im-server?style=for-the-badge&logo=github&colorB=ff69b4)](https://github.com/openimsdk/open-im-server/stargazers)
[![Forks](https://img.shields.io/github/forks/openimsdk/open-im-server?style=for-the-badge&logo=github&colorB=blue)](https://github.com/openimsdk/open-im-server/network/members)
[![Codecov](https://img.shields.io/codecov/c/github/openimsdk/open-im-server?style=for-the-badge&logo=codecov&colorB=orange)](https://app.codecov.io/gh/openimsdk/open-im-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/openimsdk/open-im-server?style=for-the-badge)](https://goreportcard.com/report/github.com/openimsdk/open-im-server)
[![Go Reference](https://img.shields.io/badge/Go%20Reference-blue.svg?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/openimsdk/open-im-server/v3)
[![License](https://img.shields.io/badge/license-Apache--2.0-green?style=for-the-badge)](https://github.com/openimsdk/open-im-server/blob/main/LICENSE)
[![Slack](https://img.shields.io/badge/Slack-500%2B-blueviolet?style=for-the-badge&logo=slack&logoColor=white)](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A)
[![Best Practices](https://img.shields.io/badge/Best%20Practices-purple?style=for-the-badge)](https://www.bestpractices.dev/projects/8045)
[![Good First Issues](https://img.shields.io/github/issues/openimsdk/open-im-server/good%20first%20issue?style=for-the-badge&logo=github)](https://github.com/openimsdk/open-im-server/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22)
[![Language](https://img.shields.io/badge/Language-Go-blue.svg?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)

<p align="center">
  <a href="../../README.md">English</a> · 
  <a href="../../README_zh_CN.md">中文</a> · 
  <a href="./README_uk.md">Українська</a> · 
  <a href="./README_cs.md">Česky</a> · 
  <a href="./README_hu.md">Magyar</a> · 
  <a href="./README_es.md">Español</a> · 
  <a href="./README_fa.md">فارسی</a> · 
  <a href="./README_fr.md">Français</a> · 
  <a href="./README_de.md">Deutsch</a> · 
  <a href="./README_pl.md">Polski</a> · 
  <a href="./README_id.md">Indonesian</a> · 
  <a href="./README_fi.md">Suomi</a> · 
  <a href="./README_ml.md">മലയാളം</a> · 
  <a href="./README_ja.md">日本語</a> · 
  <a href="./README_nl.md">Nederlands</a> · 
  <a href="./README_it.md">Italiano</a> · 
  <a href="./README_ru.md">Русский</a> · 
  <a href="./README_pt_BR.md">Português (Brasil)</a> · 
  <a href="./README_eo.md">Esperanto</a> · 
  <a href="./README_ko.md">한국어</a> · 
  <a href="./README_ar.md">العربي</a> · 
  <a href="./README_vi.md">Tiếng Việt</a> · 
  <a href="./README_da.md">Dansk</a> · 
  <a href="./README_el.md">Ελληνικά</a> · 
  <a href="./README_tr.md">Türkçe</a>
</p>

</div>

</p>

## Ⓜ️ À propos de OpenIM

OpenIM est une plateforme de services conçue spécifiquement pour intégrer des fonctionnalités de communication telles que le chat, les appels audio et vidéo, les notifications, ainsi que les robots de chat IA dans les applications. Elle offre une série d'API puissantes et de Webhooks, permettant aux développeurs d'incorporer facilement ces caractéristiques interactives dans leurs applications. OpenIM n'est pas en soi une application de chat autonome, mais sert de plateforme supportant d'autres applications pour réaliser des fonctionnalités de communication enrichies. L'image ci-dessous montre les relations d'interaction entre AppServer, AppClient, OpenIMServer et OpenIMSDK pour illustrer spécifiquement.

![Relation App-OpenIM](../../images/oepnim-design.png)

## 🚀 À propos de OpenIMSDK

**OpenIMSDK** est un SDK IM conçu pour **OpenIMServer** spécialement créé pour être intégré dans les applications clientes. Ses principales fonctionnalités et modules comprennent :

- 🌟 Fonctionnalités clés :

  - 📦 Stockage local
  - 🔔 Rappels de l'écouteur
  - 🛡️ Encapsulation d'API
  - 🌐 Gestion de la connexion

  ## 📚 Modules principaux ：

  1. 🚀 Initialisation et connexion
  2. 👤 Gestion des utilisateurs
  3. 👫 Gestion des amis
  4. 🤖 Fonctionnalités de groupe
  5. 💬 Traitement des conversations

Il est construit avec Golang et supporte le déploiement multiplateforme, assurant une expérience d'accès cohérente sur toutes les plateformes。

👉 **[Explorer le SDK GO](https://github.com/openimsdk/openim-sdk-core)**

## 🌐 À propos de OpenIMServer

- **OpenIMServer** présente les caractéristiques suivantes ：
  - 🌐 Architecture microservices : prend en charge le mode cluster, incluant le gateway (passerelle) et plusieurs services rpc。
  - 🚀 Divers modes de déploiement : supporte le déploiement via le code source, Kubernetes ou Docker。
  - Support d'une masse d'utilisateurs : plus de cent mille pour les super grands groupes, des millions d'utilisateurs, et des milliards de messages。

### Fonctionnalités commerciales améliorées :

- **REST API**：OpenIMServer fournit une REST API pour les systèmes commerciaux, visant à accorder plus de fonctionnalités, telles que la création de groupes via l'interface backend, l'envoi de messages push, etc。
- **Webhooks**：OpenIMServer offre des capacités de rappel pour étendre davantage les formes d'entreprise. Un rappel signifie que OpenIMServer enverra une requête au serveur d'entreprise avant ou après qu'un événement se soit produit, comme un rappel avant ou après l'envoi d'un message。

👉 **[En savoir plus](https://docs.openim.io/guides/introduction/product)**

## :building_construction: Architecture globale

Plongez dans le cœur de la fonctionnalité d'Open-IM-Server avec notre diagramme d'architecture.

![Architecture globale](../../images/architecture-layers.png)

## :rocket: Démarrage rapide

Nous prenons en charge de nombreuses plateformes. Voici les adresses pour une expérience rapide du côté web :

👉 **[Démo web en ligne OpenIM](https://www.openim.io/zh/commercial)**

🤲 Pour faciliter l'expérience utilisateur, nous proposons plusieurs solutions de déploiement. Vous pouvez choisir votre méthode de déploiement selon la liste ci-dessous ：

- **[Guide de déploiement du code source](https://docs.openim.io/guides/gettingStarted/imSourceCodeDeployment)**
- **[Guide de déploiement Docker](https://docs.openim.io/guides/gettingStarted/dockerCompose)**
- **[Guide de déploiement Kubernetes](https://docs.openim.io/guides/gettingStarted/k8s-deployment)**
- **[Guide de déploiement pour développeur Mac](https://docs.openim.io/guides/gettingstarted/mac-deployment-guide)**

## :hammer_and_wrench: Commencer à développer avec OpenIM

Chez OpenIM, notre objectif est de construire une communauté open source de premier plan. Nous avons un ensemble de standards, disponibles dans le[ dépôt communautaire](https://github.com/OpenIMSDK/community)。
Si vous souhaitez contribuer à ce dépôt Open-IM-Server, veuillez lire notre[ document pour les contributeurs](https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md)。

Avant de commencer, assurez-vous que vos modifications sont nécessaires. La meilleure manière est de créer une[ nouvelle discussion ](https://github.com/openimsdk/open-im-server/discussions/new/choose) ou une [ communication Slack,](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A)，ou si vous identifiez un problème, de[ signaler d'abord ](https://github.com/openimsdk/open-im-server/issues/new/choose)。

- [Référence de l'API OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/api.md)
- [Journalisation Bash OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/bash-log.md)
- [Actions CI/CD OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/cicd-actions.md)
- [Conventions de code OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/code-conventions.md)
- [Directives de commit OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/commit.md)
- [Guide de développement OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/development.md)
- [Structure de répertoire OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/directory.md)
- [Configuration de l'environnement OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/environment.md)
- [Référence des codes d'erreur OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/error-code.md)
- [Workflow Git OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/git-workflow.md)
- [Guide Cherry Pick Git OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/gitcherry-pick.md)
- [Workflow GitHub OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/github-workflow.md)
- [Normes de code Go OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/go-code.md)
- [Directives d'image OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/images.md)
- [Configuration initiale OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/init-config.md)
- [Guide d'installation Docker OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/install-docker.md)
- [Installation du système Linux OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/install-openim-linux-system.md)
- [Guide de développement Linux OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/linux-development.md)
- [Guide des actions locales OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/local-actions.md)
- [Conventions de journalisation OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/logging.md)
- [Déploiement hors ligne OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/offline-deployment.md)
- [Outils Protoc OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/protoc-tools.md)
- [Guide de test OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/test.md)
- [Utilitaire Go OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-go.md)
- [Utilitaires Makefile OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-makefile.md)
- [Utilitaires de script OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-scripts.md)
- [Versionnement OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/version.md)
- [Gérer le déploiement du backend et la surveillance](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/prometheus-grafana.md)
- [Guide de déploiement pour développeur Mac pour OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/mac-developer-deployment-guide.md)

> ## :calendar: Réunions de la Communauté

Nous voulons que tout le monde s'implique dans notre communauté et contribue au code, nous offrons des cadeaux et des récompenses, et nous vous invitons à nous rejoindre chaque jeudi soir.
Notre conférence se trouve dans le [ Slack OpenIM ](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A) 🎯, ensuite vous pouvez rechercher le pipeline Open-IM-Server pour rejoindre

Nous prenons des notes de chaque [réunion bihebdomadaire ](https://github.com/orgs/OpenIMSDK/discussions/categories/meeting) dans les [discussions GitHub](https://github.com/openimsdk/open-im-server/discussions/categories/meeting), Nos notes de réunion historiques, ainsi que les rediffusions des réunions sont disponibles sur [ Google Docs :bookmark_tabs:](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing).

## :eyes: Qui Utilise OpenIM

Consultez notre page [ études de cas d'utilisateurs ](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md) pour une liste des utilisateurs du projet. N'hésitez pas à laisser un [📝commentaire](https://github.com/openimsdk/open-im-server/issues/379) et partager votre cas d'utilisation.

## :page_facing_up: License

OpenIM est sous licence Apache 2.0. Voir [LICENSE](https://github.com/openimsdk/open-im-server/tree/main/LICENSE) pour le texte complet de la licence.

Le logo OpenIM, y compris ses variations et versions animées, affiché dans ce dépôt[OpenIM](https://github.com/openimsdk/open-im-server) sous les répertoires [assets/logo](../../assets/logo) et [assets/logo-gif](assets/logo-gif) sont protégés par les lois sur le droit d'auteur.

## 🔮 Merci à nos contributeurs !

<a href="https://github.com/openimsdk/open-im-server/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=openimsdk/open-im-server" />
</a>
