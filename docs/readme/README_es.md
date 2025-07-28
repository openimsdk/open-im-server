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

## Ⓜ️ Acerca de OpenIM

OpenIM es una plataforma de servicio diseñada específicamente para integrar chat, llamadas de audio y video, notificaciones y chatbots de IA en aplicaciones. Proporciona una gama de potentes API y Webhooks, lo que permite a los desarrolladores incorporar fácilmente estas características interactivas en sus aplicaciones. OpenIM no es una aplicación de chat independiente, sino que sirve como una plataforma para apoyar a otras aplicaciones en lograr funcionalidades de comunicación enriquecidas. El siguiente diagrama ilustra la interacción entre AppServer, AppClient, OpenIMServer y OpenIMSDK para explicar en detalle.

![Relación App-OpenIM](../../docs/images/oepnim-design.png)

## 🚀 Acerca de OpenIMSDK

**OpenIMSDK** es un SDK de mensajería instantánea diseñado para **OpenIMServer**, creado específicamente para su incorporación en aplicaciones cliente. Sus principales características y módulos son los siguientes:

- 🌟 Características Principales:

  - 📦 Almacenamiento local
  - 🔔 Callbacks de escuchas
  - 🛡️ Envoltura de API
  - 🌐 Gestión de conexiones

- 📚 Módulos Principales:

  1. 🚀 Inicialización y acceso
  2. 👤 Gestión de usuarios
  3. 👫 Gestión de amigos
  4. 🤖 Funciones de grupo
  5. 💬 Manejo de conversaciones

Está construido con Golang y soporta despliegue multiplataforma, asegurando una experiencia de acceso consistente en todas las plataformas.

👉 **[Explora el SDK de GO](https://github.com/openimsdk/openim-sdk-core)**

## 🌐 Acerca de OpenIMServer

- **OpenIMServer** tiene las siguientes características:
  - 🌐 Arquitectura de microservicios: Soporta modo cluster, incluyendo un gateway y múltiples servicios rpc.
  - 🚀 Métodos de despliegue diversos: Soporta el despliegue a través de código fuente, Kubernetes o Docker.
  - Soporte para una base de usuarios masiva: Grupos super grandes con cientos de miles de usuarios, decenas de millones de usuarios y miles de millones de mensajes.

### Funcionalidad Empresarial Mejorada:

- **API REST**: OpenIMServer ofrece APIs REST para sistemas empresariales, destinadas a empoderar a las empresas con más funcionalidades, como la creación de grupos y el envío de mensajes push a través de interfaces de backend.
- **Webhooks**: OpenIMServer proporciona capacidades de callback para extender más formas de negocio. Un callback significa que OpenIMServer envía una solicitud al servidor empresarial antes o después de un cierto evento, como callbacks antes o después de enviar un mensaje.

👉 **[Aprende más](https://docs.openim.io/guides/introduction/product)**

## :building_construction: Arquitectura General

Adéntrate en el corazón de la funcionalidad de Open-IM-Server con nuestro diagrama de arquitectura.

![Arquitectura General](../../docs/images/architecture-layers.png)

## :rocket: Inicio Rápido

:rocket: Inicio Rápido
Apoyamos muchas plataformas. Aquí están las direcciones para una experiencia rápida en el lado web:

👉 **[ Demostración web en línea de OpenIM](https://web-enterprise.rentsoft.cn/)**

🤲 Para facilitar la experiencia del usuario, ofrecemos varias soluciones de despliegue. Puedes elegir tu método de despliegue de la lista a continuación:

- **[Guía de Despliegue de Código Fuente](https://docs.openim.io/guides/gettingStarted/imSourceCodeDeployment)**
- **[Guía de Despliegue con Docker](https://docs.openim.io/guides/gettingStarted/dockerCompose)**
- **[Guía de Despliegue con Kubernetes](https://docs.openim.io/guides/gettingStarted/k8s-deployment)**
- **[Guía de Despliegue para Desarrolladores en Mac](https://docs.openim.io/guides/gettingstarted/mac-deployment-guide)**

## :hammer_and_wrench: Para Comenzar a Desarrollar en OpenIM

[![Abrir en Contenedor de Desarrollo](https://img.shields.io/static/v1?label=Dev%20Container&message=Open&color=blue&logo=visualstudiocode)](https://vscode.dev/github/openimsdk/open-im-server)

Nuestro objetivo en OpenIM es construir una comunidad de código abierto de nivel superior. Tenemos un conjunto de estándares,
en el [repositorio de la Comunidad.](https://github.com/OpenIMSDK/community).

Si te gustaría contribuir a este repositorio de Open-IM-Server, por favor lee nuestra [documentación para colaboradores](https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md).

Antes de comenzar, asegúrate de que tus cambios sean demandados. Lo mejor para eso es crear una [nueva discusión](https://github.com/openimsdk/open-im-server/discussions/new/choose) O [Comunicación en Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A), o si encuentras un problema, [repórtalo](https://github.com/openimsdk/open-im-server/issues/new/choose) primero.

- [Referencia de API de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/api.md)
- [Registro de Bash de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/bash-log.md)
- [Acciones de CI/CD de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/cicd-actions.md)
- [Convenciones de Código de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/code-conventions.md)
- [Guías de Commit de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/commit.md)
- [Guía de Desarrollo de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/development.md)
- [Estructura de Directorios de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/directory.md)
- [Configuración de Entorno de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/environment.md)
- [Referencia de Códigos de Error de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/error-code.md)
- [Flujo de Trabajo de Git de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/git-workflow.md)
- [Guía de Cherry Pick de Git de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/gitcherry-pick.md)
- [Flujo de Trabajo de GitHub de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/github-workflow.md)
- [Estándares de Código Go de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/go-code.md)
- [Guías de Imágenes de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/images.md)
- [Configuración Inicial de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/init-config.md)
- [Guía de Instalación de Docker de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/install-docker.md)
- [Instalación del Sistema Linux de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/install-openim-linux-system.md)
- [Guía de Desarrollo Linux de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/linux-development.md)
- [Guía de Acciones Locales de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/local-actions.md)
- [Convenciones de Registro de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/logging.md)
- [Despliegue sin Conexión de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/offline-deployment.md)
- [Herramientas Protoc de OpenIMM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/protoc-tools.md)
- [Guía de Pruebas de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/test.md)
- [Utilidades Go de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-go.md)
- [Utilidades de Makefile de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-makefile.md)
- [Utilidades de Script de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-scripts.md)
- [Versionado de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/version.md)
- [Gestión de backend y despliegue de monitoreo](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/prometheus-grafana.md)
- [Guía de Despliegue para Desarrolladores Mac de OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/mac-developer-deployment-guide.md)

## :busts_in_silhouette: Comunidad

- 📚 [Comunidad de OpenIM](https://github.com/OpenIMSDK/community)
- 💕 [Grupo de Interés de OpenIM](https://github.com/Openim-sigs)
- 🚀 [Únete a nuestra comunidad de Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A)
- :eyes: [Únete a nuestro wechat (微信群)](https://openim-1253691595.cos.ap-nanjing.myqcloud.com/WechatIMG20.jpeg)

## :calendar: Reuniones de la Comunidad

Queremos que cualquiera se involucre en nuestra comunidad y contribuya con código, ofrecemos regalos y recompensas, y te damos la bienvenida para que te unas a nosotros cada jueves por la noche.

Nuestra conferencia está en [OpenIM Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A) 🎯, luego puedes buscar el pipeline de Open-IM-Server para unirte

Tomamos notas de cada [reunión quincenal](https://github.com/orgs/OpenIMSDK/discussions/categories/meeting) en [discusiones de GitHub](https://github.com/openimsdk/open-im-server/discussions/categories/meeting), Nuestras notas de reuniones históricas, así como las repeticiones de las reuniones están disponibles en [Google Docs :bookmark_tabs:](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing).

## :eyes: Quiénes Están Usando OpenIM

Consulta nuestros [estudios de caso de usuarios](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md) página para obtener una lista de los usuarios del proyecto. No dudes en dejar un [📝comentario](https://github.com/openimsdk/open-im-server/issues/379) y compartir tu caso de uso.

## :page_facing_up: Licencia

OpenIM está bajo la licencia Apache 2.0. Consulta [LICENSE](https://github.com/openimsdk/open-im-server/tree/main/LICENSE) para ver el texto completo de la licencia.

El logotipo de OpenIM, incluyendo sus variaciones y versiones animadas, que se muestran en este repositorio [OpenIM](https://github.com/openimsdk/open-im-server) en los directorios [assets/logo](../../assets/logo) y [assets/logo-gif](assets/logo-gif) están protegidos por las leyes de derechos de autor.

## 🔮 iGracias a nuestros colaboradores!

<a href="https://github.com/openimsdk/open-im-server/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=openimsdk/open-im-server" />
</a>
