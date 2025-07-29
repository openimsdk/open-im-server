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

## Ⓜ️ O OpenIM

OpenIM je platforma služeb speciálně navržená pro integraci chatu, audio-video hovorů, upozornění a chatbotů AI do aplikací. Poskytuje řadu výkonných rozhraní API a webhooků, které vývojářům umožňují snadno začlenit tyto interaktivní funkce do svých aplikací. OpenIM není samostatná chatovací aplikace, ale spíše slouží jako platforma pro podporu jiných aplikací při dosahování bohatých komunikačních funkcí. Následující diagram ilustruje interakci mezi AppServer, AppClient, OpenIMServer a OpenIMSDK pro podrobné vysvětlení.

![App-OpenIM Relationship](../images/oepnim-design.png)

## 🚀 O OpenIMSDK

**OpenIMSDK** je IM SDK navržený pro**OpenIMServer**, vytvořený speciálně pro vkládání do klientských aplikací. Jeho hlavní vlastnosti a moduly jsou následující:

- 🌟 Hlavní vlastnosti:

  - 📦 Místní úložiště
  - 🔔 Zpětná volání posluchačů
  - 🛡️ API obalování
  - 🌐 Správa připojení

- 📚 hlavní moduly:

  1. 🚀 Inicializace a přihlášení
  2. 👤 Správa uživatelů
  3. 👫 Správa přátel
  4. 🤖 Skupinové funkce
  5. 💬 Zpracování konverzace

Je postaven pomocí Golang a podporuje nasazení napříč platformami, což zajišťuje konzistentní přístup na všech platformách.

👉 **[Prozkoumat GO SDK](https://github.com/openimsdk/openim-sdk-core)**

## 🌐 O OpenIMServeru

- **OpenIMServer** má následující vlastnosti:
  - 🌐 Architektura mikroslužeb: Podporuje režim clusteru, včetně brány a více služeb RPC.
  - 🚀 Různé metody nasazení: Podporuje nasazení prostřednictvím zdrojového kódu, Kubernetes nebo Docker.
  - Podpora masivní uživatelské základny: Super velké skupiny se stovkami tisíc uživatelů, desítkami milionů uživatelů a miliardami zpráv.

### Vylepšené obchodní funkce:

- **REST API**: OpenIMServer nabízí REST API pro podnikové systémy, jejichž cílem je poskytnout podnikům více funkcí, jako je vytváření skupin a odesílání push zpráv přes backendová rozhraní.
- **Webhooks**: OpenIMServer poskytuje možnosti zpětného volání pro rozšíření více obchodních formulářů. Zpětné volání znamená, že OpenIMServer odešle požadavek na obchodní server před nebo po určité události, jako jsou zpětná volání před nebo po odeslání zprávy.

👉 **[Další informace](https://docs.openim.io/guides/introduction/product)**

## :building_construction: Celková architektura

Ponořte se do srdce funkčnosti Open-IM-Server s naším diagramem architektury.

![Overall Architecture](../images/architecture-layers.png)

## :rocket: Rychlý start

Podporujeme mnoho platforem. Zde jsou adresy pro rychlou práci na webové stránce:

👉 **[Online webová ukázka OpenIM](https://web-enterprise.rentsoft.cn/)**

🤲 Pro usnadnění uživatelské zkušenosti nabízíme různá řešení nasazení. Způsob nasazení si můžete vybrat ze seznamu níže:

- **[Průvodce nasazením zdrojového kódu](https://docs.openim.io/guides/gettingStarted/imSourceCodeDeployment)**
- **[Docker Deployment Guide](https://docs.openim.io/guides/gettingStarted/dockerCompose)**
- **[Průvodce nasazením Kubernetes](https://docs.openim.io/guides/gettingStarted/k8s-deployment)**
- **[Průvodce nasazením pro vývojáře Mac](https://docs.openim.io/guides/gettingstarted/mac-deployment-guide)**

## :hammer_and_wrench: Chcete-li začít vyvíjet OpenIM

[![Open in Dev Container](https://img.shields.io/static/v1?label=Dev%20Container&message=Open&color=blue&logo=visualstudiocode)](https://vscode.dev/github/openimsdk/open-im-server)

OpenIM Naším cílem je vybudovat špičkovou open source komunitu. Máme soubor standardů v [komunitním repozitáři](https://github.com/OpenIMSDK/community).

Pokud byste chtěli přispět do tohoto úložiště Open-IM-Server, přečtěte si naši [dokumentaci pro přispěvatele](https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md).

Než začnete, ujistěte se, že jsou vaše změny vyžadovány. Nejlepší pro to je vytvořit [nová diskuze](https://github.com/openimsdk/open-im-server/discussions/new/choose) NEBO [Slack Communication](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A), nebo pokud narazíte na problém, [nahlásit jej](https://github.com/openimsdk/open-im-server/issues/new/choose) jako první.

- [OpenIM API Reference](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/api.md)
- [Protokolování OpenIM Bash](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/bash-log.md)
- [Akce OpenIM CI/CD](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/cicd-actions.md)
- [Konvence kódu OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/code-conventions.md)
- [Pokyny k zavázání OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/commit.md)
- [Průvodce vývojem OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/development.md)
- [Struktura adresáře OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/directory.md)
- [Nastavení prostředí OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/environment.md)
- [Referenční kód chybového kódu OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/error-code.md)
- [Pracovní postup OpenIM Git](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/git-workflow.md)
- [OpenIM Git Cherry Pick Guide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/gitcherry-pick.md)
- [Pracovní postup OpenIM GitHub](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/github-workflow.md)
- [standardy kódu OpenIM Go](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/go-code.md)
- [Pokyny pro obrázky OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/images.md)
- [Počáteční konfigurace OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/init-config.md)
- [Průvodce instalací OpenIM Docker](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/install-docker.md)
- [nstalace systému OpenIM OpenIM Linux](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/install-openim-linux-system.md)
- [OpenIM Linux Development Guide](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/linux-development.md)
- [Průvodce místními akcemi OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/local-actions.md)
- [Konvence protokolování OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/logging.md)
- [Offline nasazení OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/offline-deployment.md)
- [Nástroje protokolu OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/protoc-tools.md)
- [Příručka testování OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/test.md)
- [OpenIM Utility Go](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-go.md)
- [OpenIM Makefile Utilities](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-makefile.md)
- [OpenIM Script Utilities](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-scripts.md)
- [OpenIM Versioning](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/version.md)
- [Spravovat backend a monitorovat nasazení](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/prometheus-grafana.md)
- [Průvodce nasazením pro vývojáře Mac pro OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/mac-developer-deployment-guide.md)

## :busts_in_silhouette: Společenství

- 📚 [Komunita OpenIM](https://github.com/OpenIMSDK/community)
- 💕 [Zájmová skupina OpenIM](https://github.com/Openim-sigs)
- 🚀 [Připojte se k naší komunitě Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A)
- :eyes: [Připojte se k našemu wechatu](https://openim-1253691595.cos.ap-nanjing.myqcloud.com/WechatIMG20.jpeg)

## :calendar: Komunitní setkání

Chceme, aby se do naší komunity a přispívání kódu zapojil kdokoli, nabízíme dárky a odměny a vítáme vás, abyste se k nám připojili každý čtvrtek večer.

Naše konference je v [OpenIM Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A) 🎯, pak můžete vyhledat kanál Open-IM-Server a připojit se

Zaznamenáváme si každou [dvoutýdenní schůzku](https://github.com/orgs/OpenIMSDK/discussions/categories/meeting)do [diskuzí na GitHubu](https://github.com/openimsdk/open-im-server/discussions/categories/meeting), naše historické poznámky ze schůzek a také záznamy schůzek jsou k dispozici na [Dokumenty Google :bookmark_tabs:](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing).

## :eyes: Kdo používá OpenIM

Podívejte se na naši stránku [případové studie uživatelů](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md), kde najdete seznam uživatelů projektu. Neváhejte zanechat[📝komentář](https://github.com/openimsdk/open-im-server/issues/379) a podělte se o svůj případ použití.

## :page_facing_up: License

OpenIM je licencován pod licencí Apache 2.0. Úplný text licence naleznete v [LICENCE](https://github.com/openimsdk/open-im-server/tree/main/LICENSE).

Logo OpenIM, včetně jeho variací a animovaných verzí, zobrazené v tomto úložišti [OpenIM](https://github.com/openimsdk/open-im-server)v adresářích [assets/logo](./assets/logo) a [assets/logo-gif](assets/logo-gif) je chráněno autorským právem.

## 🔮 Děkujeme našim přispěvatelům!

<a href="https://github.com/openimsdk/open-im-server/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=openimsdk/open-im-server" />
</a>
