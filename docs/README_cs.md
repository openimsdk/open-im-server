# Dokumenty serveru OpenIM

Vítejte v centru dokumentace OpenIM! Toto centrum poskytuje komplexní řadu průvodců a manuálů navržených tak, aby vám pomohly co nejlépe využít vaše zkušenosti s OpenIM.

## Obsah

1. [Contrib](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib) - Pokyny pro přispívání a konfigurace pro vývojáře
2. [Conversions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib) - Konvence kódování, zásady protokolování a další transformační nástroje

------

## Contrib

Tato část nabízí vývojářům podrobný návod, jak přispívat kódem, nastavovat prostředí a sledovat související procesy.

- [Code Conventions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/code-conventions.md) - Pravidla a konvence pro psaní kódu v OpenIM.
- [Development Guide](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/development.md) - Návod, jak provádět vývoj v rámci OpenIM.
- [Git Cherry Pick](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/gitcherry-pick.md) - Pokyny pro operace sběru třešní.
- [Git Workflow](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/git-workflow.md) - Pracovní postup git v OpenIM.
- [Initialization Configurations](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/init-config.md) - Pokyny k nastavení a inicializaci OpenIM.
- [Docker Installation](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/install-docker.md) - Jak nainstalovat Docker na váš počítač.
- [Linux Development Environment](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/linux-development.md) - Průvodce nastavením vývojového prostředí na Linuxu.
- [Local Actions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/local-actions.md) - Pokyny k provádění určitých společných akcí na místní úrovni.
- [Offline Deployment](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/offline-deployment.md) - Metody nasazení OpenIM offline.
- [Protoc Tools](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/protoc-tools.md) - Průvodce používáním protokolových nástrojů.
- [Go Tools](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/util-go.md) - Nástroje a knihovny v OpenIM for Go.
- [Makefile Tools](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/util-makefile.md) - Osvědčené postupy a nástroje pro Makefile.
- [Script Tools](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/util-scripts.md) - Doporučené postupy a nástroje pro skripty.

## Konverze

Tato část představuje různé konvence a zásady v rámci OpenIM, zahrnující kód, protokoly, verze a další.

- [API Conversions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/api.md) - Pokyny a metody pro konverze API.
- [Logging Policy](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/bash-log.md) - Zásady a konvence protokolování v OpenIM.
- [CI/CD Actions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/cicd-actions.md) - Postupy a konvence pro CI/CD.
- [Commit Conventions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/commit.md) - Konvence pro odevzdání kódu v OpenIM.
- [Directory Conventions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/directory.md) - Struktura adresářů a konvence v rámci OpenIM.
- [Error Codes](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/error-code.md) - Seznam a popisy chybových kódů.
- [Go Code Conversions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/go-code.md) - Konvence a převody pro kód Go.
- [Docker Image Strategy](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/images.md) - Strategie správy pro obrazy OpenIM Docker, zahrnující různé architektury a úložiště obrazů.
- [Logging Conventions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/logging.md) - Další podrobné konvence o protokolování.
- [Version Conventions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/version.md) - Strategie pojmenování a správy verzí OpenIM.


## Pro vývojáře, přispěvatele a správce komunity

### Vývojáři a přispěvatelé

Pokud jste vývojář nebo někdo, kdo má zájem přispívat:

- Seznamte se s našimi [Code Conventions](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/code-conventions.md) a [Git Workflow](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/git-workflow.md), abyste zajistili hladké příspěvky.
- Ponořte se do [Development Guide](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/development.md), abyste se seznámili s vývojovými postupy v OpenIM.

### Správci komunity

Jako správce komunity:

- Zajistěte, aby příspěvky odpovídaly standardům uvedeným v naší dokumentaci.
- Pravidelně kontrolujte [Logging Policy](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/bash-log.md) a [Error Codes](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/error-code.md) , abyste zůstali aktuální.

## Pro uživatele

Uživatelé by měli věnovat zvláštní pozornost:

- [Docker Installation](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/install-docker.md) - Nezbytné, pokud plánujete používat obrazy OpenIM Docker.
- [Docker Image Strategy](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/images.md) - Chcete-li porozumět různým dostupným obrázkům a jak vybrat ten správný pro vaši architekturu.