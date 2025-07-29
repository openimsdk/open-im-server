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

## :busts_in_silhouette: Fællesskab

- 📚 [OpenIM-fællesskab](https://github.com/OpenIMSDK/community)
- 💕 [OpenIM-interessegruppe](https://github.com/Openim-sigs)
- 🚀 [Deltag i vores Slack-fællesskab](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A)
- :eyes: [Deltag i vores WeChat (微信群)](https://openim-1253691595.cos.ap-nanjing.myqcloud.com/WechatIMG20.jpeg)
- 👫 [Deltag i vores Reddit](https://www.reddit.com/r/OpenIMessaging)
- 💬 [Følg vores Twitter-konto](https://twitter.com/openimsdk)

## Ⓜ️ Om OpenIM

OpenIM er en serviceplatform designet specifikt til integration af chat, lyd-videoopkald, notifikationer og AI-chatbots i applikationer. Den tilbyder en række kraftfulde API'er og Webhooks, som gør det let for udviklere at integrere disse interaktive funktioner i deres applikationer. OpenIM er ikke en selvstændig chatapplikation, men fungerer snarere som en platform, der understøtter andre applikationer i at opnå omfattende kommunikationsfunktionaliteter. Følgende diagram illustrerer interaktionen mellem AppServer, AppClient, OpenIMServer og OpenIMSDK for at forklare detaljeret.

![App-OpenIM Relationship](../images/oepnim-design.png)

## 🚀 Om OpenIMSDK

**OpenIMSDK** er en IM SDK designet til **OpenIMServer**, skabt specifikt til indlejring i klientapplikationer. Dens vigtigste funktioner og moduler er som følger:

- 🌟 Hovedfunktioner:

  - 📦 Lokal lagring
  - 🔔 Lytter-callbacks
  - 🛡️ API-indkapsling
  - 🌐 Forbindelsesstyring

  ## 📚 Hovedmoduler:

  1. 🚀 Initialisering og login
  2. 👤 Brugerstyring
  3. 👫 Venstyring
  4. 🤖 Gruppefunktioner
  5. 💬 Håndtering af samtaler

Det er bygget ved hjælp af Golang og understøtter tværplatformsudrulning, hvilket sikrer en konsekvent adgangsoplevelse på tværs af alle platforme.

👉 **[Udforsk GO SDK](https://github.com/openimsdk/openim-sdk-core)**

## 🌐 Om OpenIMServer

- **OpenIMServer** har følgende karakteristika:
  - 🌐 Mikroservicarkitektur: Understøtter klyngetilstand, inklusive en gateway og flere rpc-tjenester.
  - 🚀 Forskellige udrulningsmetoder: Understøtter udrulning via kildekode, Kubernetes eller Docker.
  - Støtte til massiv brugerbase: Super store grupper med hundredtusinder af brugere, titusinder af brugere og milliarder af beskeder.

### Forbedret forretningsfunktionalitet:

- **REST API**：OpenIMServer tilbyder REST API'er til forretningssystemer, med det formål at give virksomheder flere funktioner, såsom at oprette grupper og sende push-beskeder gennem backend-grænseflader.
- **Webhooks**：OpenIMServer giver mulighed for callback-funktionalitet for at udvide flere forretningsformer. Et callback betyder, at OpenIMServer sender en anmodning til forretningsserveren før eller efter en bestemt begivenhed, som callbacks før eller efter at have sendt en besked.

👉 **[Lær mere](https://docs.openim.io/guides/introduction/product)**

## :building_construction: Samlet Arkitektur

Dyk ned i hjertet af Open-IM-Servers funktionalitet med vores arkitekturdiagram.

![Overall Architecture](../images/architecture-layers.png)

## :rocket: Hurtig start

Vi understøtter mange platforme. Her er adresserne for hurtig oplevelse på websiden:

👉 **[OpenIM online demo](https://www.openim.io/zh/commercial)**

🤲 For at lette brugeroplevelsen tilbyder vi forskellige udrulningsløsninger. Du kan vælge din udrulningsmetode fra listen nedenfor:

- **[Vejledning til udrulning af kildekode](https://docs.openim.io/guides/gettingStarted/imSourceCodeDeployment)**
- **[Vejledning til Docker-udrulning](https://docs.openim.io/guides/gettingStarted/dockerCompose)**
- **[Vejledning til Kubernetes-udrulning](https://docs.openim.io/guides/gettingStarted/k8s-deployment)**
- **[Vejledning til Mac-udviklerudrulning](https://docs.openim.io/guides/gettingstarted/mac-deployment-guide)**

## :hammer_and_wrench: For at starte udviklingen af OpenIM

[![Open in Dev Container](https://img.shields.io/static/v1?label=Dev%20Container&message=Open&color=blue&logo=visualstudiocode)](https://vscode.dev/github/openimsdk/open-im-server)

OpenIM Vores mål er at bygge et topniveau åben kildekode-fællesskab. Vi har et sæt standarder i [Community-repositoriet](https://github.com/OpenIMSDK/community).

Hvis du gerne vil bidrage til dette Open-IM-Server-repositorium, bedes du læse vores [dokumentation for bidragydere](https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md).

Før du starter, skal du sikre dig, at dine ændringer er efterspurgte. Det bedste for det er at oprette en [ny diskussion](https://github.com/openimsdk/open-im-server/discussions/new/choose) ELLER [Slack-kommunikation](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A), eller hvis du finder et problem, [rapportere det](https://github.com/openimsdk/open-im-server/issues/new/choose) først.

- [OpenIM API-referencer](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/api.md)
- [OpenIM Bash-logging](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/bash-log.md)
- [OpenIM CI/CD-handlinger](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/cicd-actions.md)
- [OpenIM kodekonventioner](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/code-conventions.md)
- [OpenIM commit-retningslinjer](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/commit.md)
- [OpenIM udviklingsguide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/development.md)
- [OpenIM mappestruktur](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/directory.md)
- [OpenIM miljøopsætning](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/environment.md)
- [OpenIM fejlkode-reference](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/error-code.md)
- [OpenIM Git-arbejdsgang](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/git-workflow.md)
- [OpenIM Git Cherry Pick-guide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/gitcherry-pick.md)
- [OpenIM GitHub-arbejdsgang](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/github-workflow.md)
- [OpenIM Go kode-standarder](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/go-code.md)
- [OpenIM billedretningslinjer](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/images.md)
- [OpenIM initialkonfiguration](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/init-config.md)
- [OpenIM Docker installationsguide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/install-docker.md)
- [OpenIM OpenIM Linux-systeminstallation](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/install-openim-linux-system.md)
- [OpenIM Linux-udviklingsguide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/linux-development.md)
- [OpenIM lokale handlingsguide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/local-actions.md)
- [OpenIM logningskonventioner](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/logging.md)
- [OpenIM offline-udrulning](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/offline-deployment.md)
- [OpenIM Protoc-værktøjer](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/protoc-tools.md)
- [OpenIM testguide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/test.md)
- [OpenIM Utility Go](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-go.md)
- [OpenIM Makefile-værktøjer](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-makefile.md)
- [OpenIM skriptværktøjer](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-scripts.md)
- [OpenIM versionsstyring](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/version.md)
- [Administrer backend og overvåg udrulning](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/prometheus-grafana.md)
- [Mac-udviklerudrulningsguide for OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/mac-developer-deployment-guide.md)

## :calendar: Fællesskabsmøder

Vi ønsker, at alle involverer sig i vores fællesskab og bidrager med kode, vi tilbyder gaver og belønninger, og vi byder dig velkommen til at deltage hver torsdag aften.

Vores konference er på [OpenIM Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A) 🎯, derefter kan du søge Open-IM-Server pipeline for at deltage.

Vi tager [notater](https://github.com/orgs/OpenIMSDK/discussions/categories/meeting) af hvert fjortendages møde i [GitHub-diskussioner](https://github.com/openimsdk/open-im-server/discussions/categories/meeting), Vores historiske mødenotater samt genudsendelser af møderne er tilgængelige på [Google Docs](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing) 📑.

## :eyes: Hvem Bruger OpenIM

Tjek vores side med [brugercasestudier](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md) for en liste over projektbrugerne. Tøv ikke med at efterlade en 📝[kommentar](https://github.com/openimsdk/open-im-server/issues/379) og dele dit brugstilfælde.

## :page_facing_up: Licens

OpenIM er licenseret under Apache 2.0-licensen. Se [LICENSE](https://github.com/openimsdk/open-im-server/tree/main/LICENSE) for den fulde licens tekst.

OpenIM-logoet, inklusive dets variationer og animerede versioner, vist i dette repositorium [OpenIM](https://github.com/openimsdk/open-im-server) under mapperne [assets/logo](../../assets/logo) og [assets/logo-gif](../../assets/logo-gif), er beskyttet af ophavsretslove.

## 🔮 Tak til vores bidragydere!

<a href="https://github.com/openimsdk/open-im-server/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=openimsdk/open-im-server" />
</a>
