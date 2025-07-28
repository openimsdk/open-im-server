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

