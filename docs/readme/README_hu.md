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

## Ⓜ️ Az OpenIM-ről

Az OpenIM egy szolgáltatási platform, amelyet kifejezetten a csevegés, az audio-video hívások, az értesítések és az AI chatbotok alkalmazásokba történő integrálására terveztek. Számos hatékony API-t és Webhookot kínál, lehetővé téve a fejlesztők számára, hogy ezeket az interaktív szolgáltatásokat könnyen beépítsék alkalmazásaikba. Az OpenIM nem egy önálló csevegőalkalmazás, hanem platformként szolgál más alkalmazások támogatására a gazdag kommunikációs funkciók elérésében. A következő diagram az AppServer, az AppClient, az OpenIMServer és az OpenIMSDK közötti interakciót szemlélteti részletesen.

![App-OpenIM Relationship](../images/oepnim-design.png)

## 🚀 Az OpenIMSDK-ról

Az **OpenIMSDK** egy **OpenIMServer** számára készült azonnali üzenetküldő SDK, amelyet kifejezetten ügyfélalkalmazásokba való beágyazáshoz hoztak létre. Fő jellemzői és moduljai a következők:

- 🌟 Főbb jellemzők:

  - 📦 Helyi raktár
  - 🔔 Hallgatói visszahívások
  - 🛡️ API-csomagolás
  - 🌐 Kapcsolatkezelés

- 📚 Fő modulok:

  1. 🚀 Inicializálás és bejelentkezés
  2. 👤 Felhasználókezelés
  3. 👫 Barátkezelés
  4. 🤖 Csoportfunkciók
  5. 💬 Beszélgetéskezelés

Golang használatával készült, és támogatja a többplatformos telepítést, biztosítva a konzisztens hozzáférési élményt minden platformon.

👉 **[Fedezze fel a GO SDK-t](https://github.com/openimsdk/openim-sdk-core)**

## 🌐 Az OpenIMServerről

- **OpenIMServer** a következő jellemzőkkel rendelkezik:
  - 🌐 Mikroszolgáltatási architektúra: Támogatja a fürt módot, beleértve az átjárót és több rpc szolgáltatást.
  - 🚀 Változatos telepítési módszerek: Támogatja a forráskódon, Kubernetesen vagy Dockeren keresztül történő telepítést.
  - Hatalmas felhasználói bázis támogatása: Szuper nagy csoportok több százezer felhasználóval, több tízmillió felhasználóval és több milliárd üzenettel.

### Továbbfejlesztett üzleti funkcionalitás:

- **REST API**: Az OpenIMServer REST API-kat kínál az üzleti rendszerek számára, amelyek célja, hogy a vállalkozásokat több funkcióval ruházza fel, mint például csoportok létrehozása és push üzenetek küldése háttérfelületeken keresztül.
- **Webhooks**: Az OpenIMServer visszahívási lehetőségeket biztosít több üzleti forma kiterjesztéséhez. A visszahívás azt jelenti, hogy az OpenIMServer kérelmet küld az üzleti szervernek egy bizonyos esemény előtt vagy után, például visszahívásokat üzenet küldése előtt vagy után.

👉 **[Tudj meg többet](https://docs.openim.io/guides/introduction/product)**

## :building_construction: Általános építészet

Merüljön el az Open-IM-Server funkcióinak szívében az architektúra diagramunk segítségével.

![Overall Architecture](../images/architecture-layers.png)

## :rocket: Gyors indítás

Számos platformot támogatunk. Íme a címek a gyors weboldali használathoz:

👉 **[OpenIM online webdemó](https://web-enterprise.rentsoft.cn/)**

🤲 A felhasználói élmény megkönnyítése érdekében különféle telepítési megoldásokat kínálunk. Az alábbi listából választhatja ki a telepítési módot:

- **[Forráskód-telepítési útmutató](https://docs.openim.io/guides/gettingStarted/imSourceCodeDeployment)**
- **[Docker telepítési útmutató](https://docs.openim.io/guides/gettingStarted/dockerCompose)**
- **[Kubernetes telepítési útmutató](https://docs.openim.io/guides/gettingStarted/k8s-deployment)**
- **[Mac fejlesztői telepítési útmutató](https://docs.openim.io/guides/gettingstarted/mac-deployment-guide)**

## :hammer_and_wrench: Az OpenIM fejlesztésének megkezdéséhez

[![Open in Dev Container](https://img.shields.io/static/v1?label=Dev%20Container&message=Open&color=blue&logo=visualstudiocode)](https://vscode.dev/github/openimsdk/open-im-server)

OpenIM Célunk egy felső szintű nyílt forráskódú közösség felépítése. Van egy szabványkészletünk a [Közösségi adattárban](https://github.com/OpenIMSDK/community).

Ha hozzá szeretne járulni ehhez az Open-IM-Server adattárhoz, kérjük, olvassa el [közreműködői dokumentációnkat](https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md).

Mielőtt elkezdené, győződjön meg arról, hogy a változtatásokra van-e igény. Erre a legjobb egy [új beszélgetés](https://github.com/openimsdk/open-im-server/discussions/new/choose) VAGY [Slack Communication](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A)létrehozása, vagy ha problémát talál, először [jelentse](https://github.com/openimsdk/open-im-server/issues/new/choose) first.

- [OpenIM API referencia](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/api.md)
- [OpenIM Bash naplózás](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/bash-log.md)
- [OpenIM CI/CD műveletek](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/cicd-actions.md)
- [OpenIM Code-egyezmények](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/code-conventions.md)
- [OpenIM Commit Guidelines](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/commit.md)
- [OpenIM fejlesztési útmutató](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/development.md)
- [OpenIM címtárszerkezet](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/directory.md)
- [OpenIM környezet beállítása](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/environment.md)
- [OpenIM hibakód hivatkozás](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/error-code.md)
- [OpenIM Git Workflow](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/git-workflow.md)
- [OpenIM Git Cherry Pick Guide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/gitcherry-pick.md)
- [OpenIM GitHub munkafolyamat](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/github-workflow.md)
- [OpenIM Go Code szabványok](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/go-code.md)
- [OpenIM képre vonatkozó irányelvek](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/images.md)
- [OpenIM kezdeti konfiguráció](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/init-config.md)
- [OpenIM Docker telepítési útmutató](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/install-docker.md)
- [OpenIM OpenIM Linux rendszertelepítés](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/install-openim-linux-system.md)
- [OpenIM Linux fejlesztési útmutató](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/linux-development.md)
- [OpenIM helyi műveletek útmutatója](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/local-actions.md)
- [OpenIM naplózási egyezmények](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/logging.md)
- [OpenIM offline telepítés](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/offline-deployment.md)
- [OpenIM Protoc Tools](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/protoc-tools.md)
- [OpenIM tesztelési útmutató](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/test.md)
- [OpenIM Utility Go](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-go.md)
- [OpenIM Makefile Utilities](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-makefile.md)
- [OpenIM Script Utilities](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-scripts.md)
- [OpenIM verzió](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/version.md)
- [A háttérrendszer kezelése és a telepítés figyelése](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/prometheus-grafana.md)
- [Mac Developer Deployment Guide for OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/mac-developer-deployment-guide.md)

## :busts_in_silhouette: Közösség

- 📚 [OpenIM közösség](https://github.com/OpenIMSDK/community)
- 💕 [OpenIM érdeklődési csoport](https://github.com/Openim-sigs)
- 🚀 [Csatlakozz a Slack közösségünkhöz](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A)
- :eyes: [Csatlakozz a wechathez](https://openim-1253691595.cos.ap-nanjing.myqcloud.com/WechatIMG20.jpeg)

## :calendar: Közösségi Találkozók

Szeretnénk, ha bárki bekapcsolódna közösségünkbe és hozzájárulna kódunkhoz, ajándékokat és jutalmakat kínálunk, és szeretettel várjuk, hogy csatlakozzon hozzánk minden csütörtök este.

Konferenciánk az [OpenIM Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A) 🎯alatt van, akkor kereshet az Open-IM-Server folyamatban a csatlakozáshoz

A [GitHub-beszélgetések](https://github.com/orgs/OpenIMSDK/discussions/categories/meeting)minden [kéthetente történő megbeszélésről](https://github.com/openimsdk/open-im-server/discussions/categories/meeting) jegyzeteket készítünk. A találkozók történeti feljegyzései, valamint az értekezletek visszajátszásai a [Google Dokumentumok :bookmark_tabs:](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing) webhelyen érhetők el.

## :eyes: Kik használják az OpenIM-et

Tekintse meg [felhasználói esettanulmányok](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md) oldalunkat a projekt felhasználóinak listájáért. Ne habozzon, hagyjon [📝megjegyzést](https://github.com/openimsdk/open-im-server/issues/379), és ossza meg használati esetét.

## :page_facing_up: Engedély

Az OpenIM licence az Apache 2.0 licence alá tartozik. A teljes licencszövegért lásd: [LICENSE](https://github.com/openimsdk/open-im-server/tree/main/LICENSE).

Az ebben az [OpenIM](https://github.com/openimsdk/open-im-server) tárolóban az [assets/logo](./assets/logo) és [assets/logo-gif](assets/logo-gif) könyvtárak alatt megjelenő OpenIM logót, beleértve annak változatait és animált változatait, szerzői jogi törvények védik.

## 🔮 Köszönjük közreműködőinknek!

<a href="https://github.com/openimsdk/open-im-server/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=openimsdk/open-im-server" />
</a>
