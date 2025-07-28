## [v3.8.3-patch.6](https://github.com/openimsdk/open-im-server/releases/tag/v3.8.3-patch.6) 	(2025-07-23)

### Bug Fixes
* fix: Add friend DB in notification sender [#3438](https://github.com/openimsdk/open-im-server/pull/3438)
* fix: remove update version file workflows have new line in 3.8.3-patch branch. [#3452](https://github.com/openimsdk/open-im-server/pull/3452)
* fix: s3 aws init [#3454](https://github.com/openimsdk/open-im-server/pull/3454)
* fix: use safe submodule init in workflows in v3.8.3-patch. [#3469](https://github.com/openimsdk/open-im-server/pull/3469)

**Full Changelog**: [v3.8.3-patch.5...v3.8.3-patch.6](https://github.com/openimsdk/open-im-server/compare/v3.8.3-patch.5...v3.8.3-patch.6)

## [v3.8.3-patch.5](https://github.com/openimsdk/open-im-server/releases/tag/v3.8.3-patch.5) 	(2025-06-10)

### New Features
* feat: optimize friend and group applications [#3396](https://github.com/openimsdk/open-im-server/pull/3396)

### Bug Fixes
* fix: solve unocrrect invite notification [Created [#3219](https://github.com/openimsdk/open-im-server/pull/3219)

### Builds
* build: update gomake version in dockerfile.[Patch branch] [#3416](https://github.com/openimsdk/open-im-server/pull/3416)

**Full Changelog**: [v3.8.3...v3.8.3-patch.5](https://github.com/openimsdk/open-im-server/compare/v3.8.3...v3.8.3-patch.5)

## [v3.8.3-patch.4](https://github.com/openimsdk/open-im-server/releases/tag/v3.8.3-patch.4) 	(2025-03-13)

### Bug Fixes
* fix: solve unocrrect invite notificationfrom #3213

**Full Changelog**: [v3.8.3-patch.3...v3.8.3-patch.4](https://github.com/openimsdk/open-im-server/compare/v3.8.3-patch.3...v3.8.3-patch.4)

## [v3.8.3-patch.3](https://github.com/openimsdk/open-im-server/releases/tag/v3.8.3-patch.3) 	(2025-03-07)

### New Features
* feat: optimizing BatchGetIncrementalGroupMember #3180

### Bug Fixes
* fix: solve uncorrect notification when set group info #3172
* fix: the sorting is wrong after canceling the administrator in group settings #3185
* fix: solve uncorrect GroupMember enter group notification type. #3188

### Refactors
* refactor: change sendNotification to sendMessage to avoid ambiguity regarding message sending behavior. #3173

**Full Changelog**: [v3.8.3-patch.2...v3.8.3-patch.3](https://github.com/openimsdk/open-im-server/compare/v3.8.3-patch.2...v3.8.3-patch.3)

## [v3.8.3-patch.2](https://github.com/openimsdk/open-im-server/releases/tag/v3.8.3-patch.2) 	(2025-02-28)

### Bug Fixes
* fix: Offline push does not have a badge && Android offline push (#3146) [#3174](https://github.com/openimsdk/open-im-server/pull/3174)

**Full Changelog**: [v3.8.3-patch.1...v3.8.3-patch.2](https://github.com/openimsdk/open-im-server/compare/v3.8.3-patch.1...v3.8.3-patch.2)

## [v3.8.3-patch.1](https://github.com/openimsdk/open-im-server/releases/tag/v3.8.3-patch.1) 	(2025-02-25)

### New Features
* feat: add backup volume && optimize log print [Created [#3121](https://github.com/openimsdk/open-im-server/pull/3121)

### Bug Fixes
* fix: seq conversion failed without exiting [Created [#3120](https://github.com/openimsdk/open-im-server/pull/3120)
* fix: check error in BatchSetTokenMapByUidPid [Created [#3123](https://github.com/openimsdk/open-im-server/pull/3123)
* fix: DeleteDoc crash [Created [#3124](https://github.com/openimsdk/open-im-server/pull/3124)
* fix: the abnormal message has no sending time, causing the SDK to be abnormal [Created [#3126](https://github.com/openimsdk/open-im-server/pull/3126)
* fix: crash caused [#3127](https://github.com/openimsdk/open-im-server/pull/3127)
* fix: the user sets the conversation timer cleanup timestamp unit incorrectly [Created [#3128](https://github.com/openimsdk/open-im-server/pull/3128)
* fix: seq conversion not reading env in docker environment [Created [#3131](https://github.com/openimsdk/open-im-server/pull/3131)

### Builds
* build: improve workflows contents. [Created [#3125](https://github.com/openimsdk/open-im-server/pull/3125)

**Full Changelog**: [v3.8.3-e-v1.1.5...v3.8.3-patch.1-e-v1.1.5](https://github.com/openimsdk/open-im-server-enterprise/compare/v3.8.3-e-v1.1.5...v3.8.3-patch.1-e-v1.1.5)