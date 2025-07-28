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

## Ⓜ️ Σχετικά με το OpenIM

Το OpenIM είναι μια πλατφόρμα υπηρεσιών σχεδιασμένη ειδικά για την ενσωμάτωση συνομιλίας, κλήσεων ήχου-βίντεο, ειδοποιήσεων και chatbots AI σε εφαρμογές. Παρέχει μια σειρά από ισχυρά API και Webhooks, επιτρέποντας στους προγραμματιστές να ενσωματώσουν εύκολα αυτές τις αλληλεπιδραστικές λειτουργίες στις εφαρμογές τους. Το OpenIM δεν είναι μια αυτόνομη εφαρμογή συνομιλίας, αλλά λειτουργεί ως πλατφόρμα υποστήριξης άλλων εφαρμογών για την επίτευξη πλούσιων λειτουργιών επικοινωνίας. Το παρακάτω διάγραμμα απεικονίζει την αλληλεπίδραση μεταξύ AppServer, AppClient, OpenIMServer και OpenIMSDK για να εξηγήσει αναλυτικά.

![App-OpenIM Relationship](../../docs/images/oepnim-design.png)

## 🚀 Σχετικά με το OpenIMSDK

Το **OpenIMSDK** είναι ένα SDK για αμεση ανταλλαγή μηνυμάτων σχεδιασμένο για το **OpenIMServer**, δημιουργήθηκε ειδικά για ενσωμάτωση σε εφαρμογές πελατών. Οι κύριες δυνατότητες και μονάδες του είναι οι εξής:

- 🌟 Κύριες Δυνατότητες:

  - 📦 Τοπική αποθήκευση
  - 🔔 Callbacks ακροατών
  - 🛡️ Περιτύλιγμα API
  - 🌐 Διαχείριση σύνδεσης

- 📚 Κύριες Μονάδες:

  1. 🚀 Αρχικοποίηση και Σύνδεση
  2. 👤 Διαχείριση Χρηστών
  3. 👫 Διαχείριση Φίλων
  4. 🤖 Λειτουργίες Ομάδας
  5. 💬 Διαχείριση Συνομιλιών

Είναι κατασκευασμένο χρησιμοποιώντας Golang και υποστηρίζει διασταυρούμενη πλατφόρμα ανάπτυξης, διασφαλίζοντας μια συνεπή εμπειρία πρόσβασης σε όλες τις πλατφόρμες.

👉 **[Εξερευνήστε το GO SDK](https://github.com/openimsdk/openim-sdk-core)**

## 🌐 Σχετικά με το OpenIMServer

- Το **OpenIMServer** έχει τις ακόλουθες χαρακτηριστικές:
  - 🌐 Αρχιτεκτονική μικροϋπηρεσιών: Υποστηρίζει λειτουργία σε σύμπλεγμα, περιλαμβάνοντας έναν πύλη και πολλαπλές υπηρεσίες rpc.
  - 🚀 Διάφοροι τρόποι ανάπτυξης: Υποστηρίζει ανάπτυξη μέσω πηγαίου κώδικα, Kubernetes, ή Docker.
  - Υποστήριξη για τεράστια βάση χρηστών: Πολύ μεγάλες ομάδες με εκατοντάδες χιλιάδες χρήστες, δεκάδες εκατομμύρια χρήστες και δισεκατομμύρια μηνύματα.

### Ενισχυμένη Επιχειρηματική Λειτουργικότητα:

- **REST API**: Το OpenIMServer προσφέρει REST APIs για επιχειρηματικά συστήματα, με στόχο την ενδυνάμωση των επιχειρήσεων με περισσότερες λειτουργικότητες, όπως η δημιουργία ομάδων και η αποστολή μηνυμάτων push μέσω backend διεπαφών.
- **Webhooks**: Το OpenIMServer παρέχει δυνατότητες επανάκλησης για την επέκταση περισσότερων επιχειρηματικών μορφών. Μια επανάκληση σημαίνει ότι το OpenIMServer στέλνει ένα αίτημα στον επιχειρηματικό διακομιστή πριν ή μετά από ένα συγκεκριμένο γεγονός, όπως επανακλήσεις πριν ή μετά την αποστολή ενός μηνύματος.

👉 **[Μάθετε περισσότερα](https://docs.openim.io/guides/introduction/product)**

## :building_construction: Συνολική Αρχιτεκτονική

Εξερευνήστε σε βάθος τη λειτουργικότητα του Open-IM-Server με το διάγραμμα αρχιτεκτονικής μας.

![Overall Architecture](../../docs/images/architecture-layers.png)

## :rocket: Γρήγορη Εκκίνηση

Υποστηρίζουμε πολλές πλατφόρμες. Εδώ είναι οι διευθύνσεις για γρήγορη εμπειρία στην πλευρά του διαδικτύου:

👉 **[Διαδικτυακή επίδειξη του OpenIM](https://web-enterprise.rentsoft.cn/)**

🤲 Για να διευκολύνουμε την εμπειρία του χρήστη, προσφέρουμε διάφορες λύσεις ανάπτυξης. Μπορείτε να επιλέξετε τη μέθοδο ανάπτυξης σας από την παρακάτω λίστα:

- **[Οδηγός Ανάπτυξης Κώδικα Πηγής](https://docs.openim.io/guides/gettingStarted/imSourceCodeDeployment)**
- **[δηγός Ανάπτυξης μέσω Docker](https://docs.openim.io/guides/gettingStarted/dockerCompose)**
- **[Οδηγός Ανάπτυξης Kubernetes](https://docs.openim.io/guides/gettingStarted/k8s-deployment)**
- **[Οδηγός Ανάπτυξης για Αναπτυξιακούς στο Mac](https://docs.openim.io/guides/gettingstarted/mac-deployment-guide)**

## :hammer_and_wrench: Για να Αρχίσετε την Ανάπτυξη του OpenIM

[![Άνοιγμα σε Dev Container](https://img.shields.io/static/v1?label=Dev%20Container&message=Open&color=blue&logo=visualstudiocode)](https://vscode.dev/github/openimsdk/open-im-server)

OpenIM Στόχος μας είναι να δημιουργήσουμε μια κορυφαίου επιπέδου ανοιχτή πηγή κοινότητας. Διαθέτουμε ένα σύνολο προτύπων, στο [Αποθετήριο Κοινότητας](https://github.com/OpenIMSDK/community).

Εάν θέλετε να συνεισφέρετε σε αυτό το αποθετήριο Open-IM-Server, παρακαλούμε διαβάστε την [τεκμηρίωση συνεισφέροντος](https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md).

Πριν ξεκινήσετε, παρακαλούμε βεβαιωθείτε ότι οι αλλαγές σας είναι ζητούμενες. Το καλύτερο για αυτό είναι να δημιουργήσετε ένα [νέα συζήτηση](https://github.com/openimsdk/open-im-server/discussions/new/choose) ή [Επικοινωνία Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A), ή αν βρείτε ένα ζήτημα, [αναφέρετέ το](https://github.com/openimsdk/open-im-server/issues/new/choose) πρώτα.

- [Αναφορά API του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/api.md)
- [Καταγραφή Bash του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/bash-log.md)
- [Ενέργειες CI/CD του OpenIMs](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/cicd-actions.md)
- [Συμβάσεις Κώδικα του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/code-conventions.md)
- [Οδηγίες Commit του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/commit.md)
- [Οδηγός Ανάπτυξης του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/development.md)
- [Δομή Καταλόγου του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/directory.md)
- [Ρύθμιση Περιβάλλοντος του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/environment.md)
- [Αναφορά Κωδικών Σφάλματος του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/error-code.md)
- [Ροή Εργασίας Git του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/git-workflow.md)
- [Οδηγός Cherry Pick του Git του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/gitcherry-pick.md)
- [Ροή Εργασίας GitHub του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/github-workflow.md)
- [Πρότυπα Κώδικα Go του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/go-code.md)
- [Οδηγίες Εικόνας του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/images.md)
- [Αρχική Διαμόρφωση του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/init-config.md)
- [Οδηγός Εγκατάστασης Docker του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/install-docker.md)
- [Οδηγός Εγκατάστασης Συστήματος Linux του Open](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/install-openim-linux-system.md)
- [Οδηγός Ανάπτυξης Linux του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/linux-development.md)
- [Οδηγός Τοπικών Δράσεων του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/local-actions.md)
- [Συμβάσεις Καταγραφής του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/logging.md)
- [Αποστολή Εκτός Σύνδεσης του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/offline-deployment.md)
- [Εργαλεία Protoc του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/protoc-tools.md)
- [Οδηγός Δοκιμών του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/test.md)
- [Χρησιμότητα Go του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-go.md)
- [Χρησιμότητες Makefile του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-makefile.md)
- [Χρησιμότητες Σεναρίου του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/util-scripts.md)
- [Έκδοση του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/version.md)
- [Διαχείριση backend και παρακολούθηση ανάπτυξης](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/prometheus-grafana.md)
- [Οδηγός Ανάπτυξης για Προγραμματιστές Mac του OpenIM](https://github.com/openimsdk/open-im-server/tree/main/docs/contrib/mac-developer-deployment-guide.md)

## :busts_in_silhouette: Κοινότητα

- 📚 [Κοινότητα OpenIM](https://github.com/OpenIMSDK/community)
- 💕 [Ομάδα Ενδιαφέροντος OpenIM](https://github.com/Openim-sigs)
- 🚀 [Εγγραφείτε στην κοινότητα Slack μας](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A)
- :eyes: [γγραφείτε στην ομάδα μας wechat (微信群)](https://openim-1253691595.cos.ap-nanjing.myqcloud.com/WechatIMG20.jpeg)

## :calendar: Συναντήσεις της κοινότητας

Θέλουμε οποιονδήποτε να εμπλακεί στην κοινότητά μας και να συνεισφέρει κώδικα. Προσφέρουμε δώρα και ανταμοιβές και σας καλωσορίζουμε να μας ενταχθείτε κάθε Πέμπτη βράδυ.

Η διάσκεψή μας είναι στο [OpenIM Slack](https://join.slack.com/t/openimsdk/shared_invite/zt-2ijy1ys1f-O0aEDCr7ExRZ7mwsHAVg9A) 🎯, στη συνέχεια μπορείτε να αναζητήσετε τη διαδικασία Open-IM-Server για να συμμετάσχετε

Κάνουμε σημειώσεις για κάθε μια [Σημειώνουμε κάθε διμηνιαία συνάντηση](https://github.com/orgs/OpenIMSDK/discussions/categories/meeting) στις [συζητήσεις του GitHub](https://github.com/openimsdk/open-im-server/discussions/categories/meeting), Οι ιστορικές μας σημειώσεις συναντήσεων, καθώς και οι επαναλήψεις των συναντήσεων είναι διαθέσιμες στο[Έγγραφα της Google :bookmark_tabs:](https://docs.google.com/document/d/1nx8MDpuG74NASx081JcCpxPgDITNTpIIos0DS6Vr9GU/edit?usp=sharing).

## :eyes: Ποιοί Χρησιμοποιούν το OpenIM

Ελέγξτε τη σελίδα με τις [μελέτες περίπτωσης χρήσης ](https://github.com/OpenIMSDK/community/blob/main/ADOPTERS.md) μας για μια λίστα των χρηστών του έργου. Μην διστάσετε να αφήσετε ένα[📝σχόλιο](https://github.com/openimsdk/open-im-server/issues/379) και να μοιραστείτε την περίπτωση χρήσης σας.

## :page_facing_up: Άδεια Χρήσης

Το OpenIM διατίθεται υπό την άδεια Apache 2.0. Δείτε τη [ΑΔΕΙΑ ΧΡΗΣΗΣ](https://github.com/openimsdk/open-im-server/tree/main/LICENSE) για το πλήρες κείμενο της άδειας.

Το λογότυπο του OpenIM, συμπεριλαμβανομένων των παραλλαγών και των κινούμενων εικόνων, που εμφανίζονται σε αυτό το αποθετήριο[OpenIM](https://github.com/openimsdk/open-im-server) υπό τις διευθύνσεις [assets/logo](../../assets/logo) και [assets/logo-gif](../../assets/logo-gif) προστατεύονται από τους νόμους περί πνευματικής ιδιοκτησίας.

## 🔮 Ευχαριστούμε τους συνεισφέροντες μας!

<a href="https://github.com/openimsdk/open-im-server/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=openimsdk/open-im-server" />
</a>
