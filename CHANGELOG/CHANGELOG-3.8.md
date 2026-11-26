## [v3.8.3-patch.17](https://github.com/openimsdk/open-im-server/releases/tag/v3.8.3-patch.17) 	(2026-10-10)

### New Features
* feat: optimize friend and group applications [#3396](https://github.com/openimsdk/open-im-server/pull/3396)
* feat: add filtering for invalid messages and invalid conversations to… [#3483](https://github.com/openimsdk/open-im-server/pull/3483)
* feat: enable redis aof-use-rdb-preamble && disable auto rdb [Created [#3535](https://github.com/openimsdk/open-im-server/pull/3535)
* feat: replace LongConn with ClientConn interface and simplify message handling [#3645](https://github.com/openimsdk/open-im-server/pull/3645)
* feat: add error code for handled friend requests and improve error handling in friend operations  [#3671](https://github.com/openimsdk/open-im-server/pull/3671)
* feat: update protocol support botPlatform [#3697](https://github.com/openimsdk/open-im-server/pull/3697)
* feat: gomake upgrade [#3701](https://github.com/openimsdk/open-im-server/pull/3701)
* feat: enhance configuration files with detailed comments for clarity [#3707](https://github.com/openimsdk/open-im-server/pull/3707)
* feat: SendBusinessNotification supported configuration parameters [#3763](https://github.com/openimsdk/open-im-server/pull/3763)
* feat: afterMsgSaveDB webhook [#3788](https://github.com/openimsdk/open-im-server/pull/3788)

### Bug Fixes
* fix: solve unocrrect invite notification [Created [#3219](https://github.com/openimsdk/open-im-server/pull/3219)
* fix: Add friend DB in notification sender [#3438](https://github.com/openimsdk/open-im-server/pull/3438)
* fix: remove update version file workflows have new line in 3.8.3-patch branch. [#3452](https://github.com/openimsdk/open-im-server/pull/3452)
* fix: s3 aws init [#3454](https://github.com/openimsdk/open-im-server/pull/3454)
* fix: use safe submodule init in workflows in v3.8.3-patch. [#3469](https://github.com/openimsdk/open-im-server/pull/3469)
* fix: correctly aggregate read seqs [#3482](https://github.com/openimsdk/open-im-server/pull/3482)
* fix: import friends send notification in v3.8.3-patch [#3488](https://github.com/openimsdk/open-im-server/pull/3488)
* fix: solve redis config db field in v3.8.3-patch [#3490](https://github.com/openimsdk/open-im-server/pull/3490)
* fix: fix incorrect kicked logic and PCAndOther Login policy In v3.8.3-patch [#3511](https://github.com/openimsdk/open-im-server/pull/3511)
* fix: solve batch incorrect error in Find DocIDs in v3.8.3-patch branch. [#3515](https://github.com/openimsdk/open-im-server/pull/3515)
* fix: fill in the most recent sendTime for a gap message to prevent th… [#3523](https://github.com/openimsdk/open-im-server/pull/3523)
* fix: switch kafka & etcd image namespace to bitnamilegacy in v3.8.3-patch branch. [#3557](https://github.com/openimsdk/open-im-server/pull/3557)
* fix: GetSortedConversationList nil pointer when chatlog not found in v3.8.3-patch branch. [#3569](https://github.com/openimsdk/open-im-server/pull/3569)
* fix: prometheus 3.8.3 [#3579](https://github.com/openimsdk/open-im-server/pull/3579)
* Fix: Resolved the issue of incorrect generation of conversationID in 3.8.3-patch branch. [#3582](https://github.com/openimsdk/open-im-server/pull/3582)
* fix: solve incorrect time.Unix and logger asyncwrite in 3.8.3-patch [#3585](https://github.com/openimsdk/open-im-server/pull/3585)
* fix: full id version [#3589](https://github.com/openimsdk/open-im-server/pull/3589)
* fix: incorrect redis fields. [#3602](https://github.com/openimsdk/open-im-server/pull/3602)
* fix: group notification [#3619](https://github.com/openimsdk/open-im-server/pull/3619)
* fix: resolve deadlock in cache eviction and improve GetBatch implementation [#3638](https://github.com/openimsdk/open-im-server/pull/3638)
* Fix/group member rejoin bug 3.8.3 [#3646](https://github.com/openimsdk/open-im-server/pull/3646)
* fix(group): set max_seq to 0 when join group [#3647](https://github.com/openimsdk/open-im-server/pull/3647)
* fix(group): move member count retrieval after member deletion for accurate updates [#3650](https://github.com/openimsdk/open-im-server/pull/3650)
* fix: Mongo Malloc upsert overwrites min_seq initialization [#3659](https://github.com/openimsdk/open-im-server/pull/3659)
* fix(core): enforce admin access for queries [#3782](https://github.com/openimsdk/open-im-server/pull/3782)
* fix(core): filter invalid and outdated conversation messages [#3792](https://github.com/openimsdk/open-im-server/pull/3792)
* fix(core): trigger after-join webhooks for invited members [#3797](https://github.com/openimsdk/open-im-server/pull/3797)
* fix(core): update MinIO and etcd images and settings [#3812](https://github.com/openimsdk/open-im-server/pull/3812)

### Chores
* chore: update tools ver [#3729](https://github.com/openimsdk/open-im-server/pull/3729)

### Refactors
* refactor: simplify platformID handling and enhance UserConnContext structure [#3648](https://github.com/openimsdk/open-im-server/pull/3648)
* refactor(ci): actions and dockerfile  [#3733](https://github.com/openimsdk/open-im-server/pull/3733)
* refactor(group): optimize group notification and cache handling [#3705](https://github.com/openimsdk/open-im-server/pull/3705)
* refactor(msg): update regex pattern for conversationID to include a t… [#3736](https://github.com/openimsdk/open-im-server/pull/3736)
* refactor(core): support kubernetes discovery registry [#3749](https://github.com/openimsdk/open-im-server/pull/3749)
* refactor(changelog): modernize generator and workflows [#3752](https://github.com/openimsdk/open-im-server/pull/3752)
* refactor(core): select discovery backend from configured type [#3772](https://github.com/openimsdk/open-im-server/pull/3772)

### Builds
* build: update gomake version in dockerfile.[Patch branch] [#3416](https://github.com/openimsdk/open-im-server/pull/3416)
* build: improve publish docker image workflow in v3.8.3-patch branch. [#3551](https://github.com/openimsdk/open-im-server/pull/3551)

### Others
* bugfix(permission):fix no permission error when set group notification. [#3675](https://github.com/openimsdk/open-im-server/pull/3675)
* CopyStructFields has been replaced due to performance issue. [#3673](https://github.com/openimsdk/open-im-server/pull/3673)
* bugfix(conversation):removed unexpectedly called functions and itself… [#3667](https://github.com/openimsdk/open-im-server/pull/3667)
* ci: action matrix release [#3762](https://github.com/openimsdk/open-im-server/pull/3762)

**Full Changelog**: [v3.8.3...v3.8.3-patch.17](https://github.com/openimsdk/open-im-server/compare/v3.8.3...v3.8.3-patch.17)

