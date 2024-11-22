## [v3.8.2](https://github.com/openimsdk/open-im-server/releases/tag/v3.8.2) 	(2024-11-22)

### New Features
* feat: improve publish docker image workflows [#2697](https://github.com/openimsdk/open-im-server/pull/2697)
* feat: Msg filter [#2703](https://github.com/openimsdk/open-im-server/pull/2703)
* feat: provide the interface required [#2712](https://github.com/openimsdk/open-im-server/pull/2712)
* feat: add webhooks of online status and remove zookeeper configuration. [#2716](https://github.com/openimsdk/open-im-server/pull/2716)
* feat: Add More Multi Login Policy [#2770](https://github.com/openimsdk/open-im-server/pull/2770)
* feat: Push configuration can ignore case sensitivity [#2775](https://github.com/openimsdk/open-im-server/pull/2775)
* feat: support app update service [#2794](https://github.com/openimsdk/open-im-server/pull/2794)
* feat: implement merge milestone PR to target-branch. [#2796](https://github.com/openimsdk/open-im-server/pull/2796)
* feat: support app update service [#2811](https://github.com/openimsdk/open-im-server/pull/2811)
* feat: ApplicationVersion move chat [#2813](https://github.com/openimsdk/open-im-server/pull/2813)
* feat: Update login policy [#2822](https://github.com/openimsdk/open-im-server/pull/2822)
* feat: merge js sdk [#2856](https://github.com/openimsdk/open-im-server/pull/2856)
* feat: Print Panic Log [#2850](https://github.com/openimsdk/open-im-server/pull/2850)

### Bug Fixes
* fix: update load file logic. [#2700](https://github.com/openimsdk/open-im-server/pull/2700)
* fix: the message I sent is not set to read seq in mongodb [#2718](https://github.com/openimsdk/open-im-server/pull/2718)
* fix: cannot modify group member avatars [#2719](https://github.com/openimsdk/open-im-server/pull/2719)
* fix: auth package import twice [#2724](https://github.com/openimsdk/open-im-server/pull/2724)
* fix: join the group chat directly, notification type error [#2772](https://github.com/openimsdk/open-im-server/pull/2772)
* fix: change update group member level logic [#2730](https://github.com/openimsdk/open-im-server/pull/2730)
* fix: joinSource check args error. [#2773](https://github.com/openimsdk/open-im-server/pull/2773)
* fix: Change group member roleLevel can`t send notification [#2777](https://github.com/openimsdk/open-im-server/pull/2777)
* fix: client sends message status error to server [#2779](https://github.com/openimsdk/open-im-server/pull/2779)
* fix: del UserB's conversation version cache when userA set conversati… [#2785](https://github.com/openimsdk/open-im-server/pull/2785)
* fix: improve setConversationAtInfo logic. [#2782](https://github.com/openimsdk/open-im-server/pull/2782)
* fix: improve transfer Owner logic when newOwner is mute. [#2790](https://github.com/openimsdk/open-im-server/pull/2790)
* fix: improve getUserInfo logic. [#2792](https://github.com/openimsdk/open-im-server/pull/2792)
* fix: improve time condition check mehtod. [#2804](https://github.com/openimsdk/open-im-server/pull/2804)
* fix: webhook before online push [#2805](https://github.com/openimsdk/open-im-server/pull/2805)
* fix: set own read seq in MongoDB when sender send a message. [#2808](https://github.com/openimsdk/open-im-server/pull/2808)
* fix: solve err Notification when setGroupInfo. [#2806](https://github.com/openimsdk/open-im-server/pull/2806)
* fix: improve condition check. [#2815](https://github.com/openimsdk/open-im-server/pull/2815)
* fix: Write back message to Redis [#2836](https://github.com/openimsdk/open-im-server/pull/2836)
* fix: get group return repeated result [#2842](https://github.com/openimsdk/open-im-server/pull/2842)
* fix: SetConversations can update new conversation [#2838](https://github.com/openimsdk/open-im-server/pull/2838)
* fix(push): push content with jpush [#2844](https://github.com/openimsdk/open-im-server/pull/2844)
* fix #2860 migrate jpns to jpush [#2861](https://github.com/openimsdk/open-im-server/pull/2861)
* fix: concurrent write to websocket connection [#2866](https://github.com/openimsdk/open-im-server/pull/2866)
* fix: Remove admin token in redis [#2871](https://github.com/openimsdk/open-im-server/pull/2871)
* Fix Push2User webhookBeforeOfflinePush [#2862](https://github.com/openimsdk/open-im-server/pull/2862)
* fix: move workflow to correct path [#2837](https://github.com/openimsdk/open-im-server/pull/2837)
* fix: del login Policy [#2825](https://github.com/openimsdk/open-im-server/pull/2825)
* fix: Wrong Redis Error Check [#2876](https://github.com/openimsdk/open-im-server/pull/2876)

### Chores
* chore: remove unused content [#2786](https://github.com/openimsdk/open-im-server/pull/2786)

### Builds
* build: improve workflows logic. [#2801](https://github.com/openimsdk/open-im-server/pull/2801)
* build: implement version file update when release. [#2826](https://github.com/openimsdk/open-im-server/pull/2826)
* build: update mongo and kafka start logic. [#2858](https://github.com/openimsdk/open-im-server/pull/2858)
* build: create changelog tool and workflows. [#2869](https://github.com/openimsdk/open-im-server/pull/2869)
* build(deps): bump github.com/golang-jwt/jwt/v4 from 4.5.0 to 4.5.1 [#2851](https://github.com/openimsdk/open-im-server/pull/2851)

### Others
* Revert: Change group member roleLevel can`t send notification [#2789](https://github.com/openimsdk/open-im-server/pull/2789)
* Introducing OpenIM Guru on Gurubase.io [#2788](https://github.com/openimsdk/open-im-server/pull/2788)
* @lkzz made their first contribution in https://github.com/openimsdk/open-im-server/pull/2724 [#2724](https://github.com/openimsdk/open-im-server/pull/2724)
* @alilestera made their first contribution in https://github.com/openimsdk/open-im-server/pull/2773 [#2773](https://github.com/openimsdk/open-im-server/pull/2773)
* @kursataktas made their first contribution in https://github.com/openimsdk/open-im-server/pull/2788 [#2788](https://github.com/openimsdk/open-im-server/pull/2788)
* @yoyo930021 made their first contribution in https://github.com/openimsdk/open-im-server/pull/2844 [#2844](https://github.com/openimsdk/open-im-server/pull/2844)
* @wikylyu made their first contribution in https://github.com/openimsdk/open-im-server/pull/2861 [#2861](https://github.com/openimsdk/open-im-server/pull/2861)
* @storyn26383 made their first contribution in https://github.com/openimsdk/open-im-server/pull/2862 [#2862](https://github.com/openimsdk/open-im-server/pull/2862)

**Full Changelog**: [v3.8.1...v3.8.2](https://github.com/openimsdk/open-im-server/compare/v3.8.1...v3.8.2)

