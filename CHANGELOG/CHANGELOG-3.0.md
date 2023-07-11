# Version logging for OpenIM

**3.0 Major refactoring**

<!-- BEGIN MUNGE: GENERATED_TOC -->

- [Version logging for OpenIM](#version-logging-for-openim)
  - [\[v3.0\]](#v30)
  - [v3.0.0 - 2023-07-10](#v300---2023-07-10)
  - [v2.9.0+1.839643f - 2023-07-07](#v2901839643f---2023-07-07)
  - [v2.9.0+2.35f07fe - 2023-07-06](#v290235f07fe---2023-07-06)
  - [v2.9.0+1.b5072b1 - 2023-07-05](#v2901b5072b1---2023-07-05)
  - [v2.9.0+3.2667a3a - 2023-07-05](#v29032667a3a---2023-07-05)
  - [v2.9.0+7.04818ca - 2023-07-05](#v290704818ca---2023-07-05)
  - [v2.9.0 - 2023-07-04](#v290---2023-07-04)
  - [v0.0.0+1.3714b4f - 2023-07-04](#v00013714b4f---2023-07-04)
  - [v0.0.0+635.8b92c90 - 2023-07-04](#v0006358b92c90---2023-07-04)
  - [v0.0.0+1.78a6d03 - 2023-07-04](#v000178a6d03---2023-07-04)
  - [v0.0.0+2.e057c18 - 2023-07-04](#v0002e057c18---2023-07-04)
  - [v0.0.0+630.b55ac4a - 2023-07-04](#v000630b55ac4a---2023-07-04)
    - [Reverts](#reverts)
    - [Pull Requests](#pull-requests)
  - [v2.3.3 - 2022-09-18](#v233---2022-09-18)
  - [v2.3.2 - 2022-09-09](#v232---2022-09-09)
  - [v2.3.0-rc2 - 2022-07-29](#v230-rc2---2022-07-29)
  - [v2.3.0-rc1 - 2022-07-25](#v230-rc1---2022-07-25)
  - [v2.3.0-rc0 - 2022-07-15](#v230-rc0---2022-07-15)
  - [v2.2.0 - 2022-07-01](#v220---2022-07-01)
  - [v2.1.0 - 2022-06-17](#v210---2022-06-17)
    - [Pull Requests](#pull-requests-1)
  - [v2.0.10 - 2022-05-13](#v2010---2022-05-13)
  - [v2.0.9 - 2022-04-29](#v209---2022-04-29)
    - [Reverts](#reverts-1)
    - [Pull Requests](#pull-requests-2)
  - [v2.0.8 - 2022-04-24](#v208---2022-04-24)
    - [Pull Requests](#pull-requests-3)
  - [v2.0.7 - 2022-04-08](#v207---2022-04-08)
    - [Pull Requests](#pull-requests-4)
  - [v2.0.6 - 2022-04-01](#v206---2022-04-01)
    - [Pull Requests](#pull-requests-5)
  - [v2.0.5 - 2022-03-24](#v205---2022-03-24)
  - [v2.04 - 2022-03-18](#v204---2022-03-18)
  - [v2.0.3 - 2022-03-11](#v203---2022-03-11)
  - [v2.0.2 - 2022-03-04](#v202---2022-03-04)
    - [Pull Requests](#pull-requests-6)
  - [v2.0.1 - 2022-02-25](#v201---2022-02-25)
  - [v2.0.0 - 2022-02-23](#v200---2022-02-23)
  - [v1.0.7 - 2021-12-17](#v107---2021-12-17)
  - [v1.0.6 - 2021-12-10](#v106---2021-12-10)
  - [v1.0.5 - 2021-12-03](#v105---2021-12-03)
  - [v1.0.4 - 2021-11-25](#v104---2021-11-25)
  - [v1.0.3 - 2021-11-12](#v103---2021-11-12)
  - [v1.0.1 - 2021-11-04](#v101---2021-11-04)
  - [v1.0.0 - 2021-10-28](#v100---2021-10-28)
    - [Reverts](#reverts-2)


<!-- END MUNGE: GENERATED_TOC -->

<a name="unreleased"></a>
## [v3.0]


<a name="v3.0.0"></a>
## [v3.0.0] - 2023-07-10

<a name="v2.9.0+1.839643f"></a>
## [v2.9.0+1.839643f] - 2023-07-07

<a name="v2.9.0+2.35f07fe"></a>
## [v2.9.0+2.35f07fe] - 2023-07-06

<a name="v2.9.0+1.b5072b1"></a>
## [v2.9.0+1.b5072b1] - 2023-07-05

<a name="v2.9.0+3.2667a3a"></a>
## [v2.9.0+3.2667a3a] - 2023-07-05

<a name="v2.9.0+7.04818ca"></a>
## [v2.9.0+7.04818ca] - 2023-07-05

<a name="v2.9.0"></a>
## [v2.9.0] - 2023-07-04

<a name="v0.0.0+1.3714b4f"></a>
## [v0.0.0+1.3714b4f] - 2023-07-04

<a name="v0.0.0+635.8b92c90"></a>
## [v0.0.0+635.8b92c90] - 2023-07-04

<a name="v0.0.0+1.78a6d03"></a>
## [v0.0.0+1.78a6d03] - 2023-07-04

<a name="v0.0.0+2.e057c18"></a>
## [v0.0.0+2.e057c18] - 2023-07-04

<a name="v0.0.0+630.b55ac4a"></a>
## [v0.0.0+630.b55ac4a] - 2023-07-04
### Reverts
- update etcd to v3.5.2 ([#206](https://github.com/OpenIMSDK/Open-IM-Server/issues/206))

### Pull Requests
- Merge branch 'tuoyun'


<a name="v2.3.3"></a>
## [v2.3.3] - 2022-09-18

<a name="v2.3.2"></a>
## [v2.3.2] - 2022-09-09

<a name="v2.3.0-rc2"></a>
## [v2.3.0-rc2] - 2022-07-29

<a name="v2.3.0-rc1"></a>
## [v2.3.0-rc1] - 2022-07-25

<a name="v2.3.0-rc0"></a>
## [v2.3.0-rc0] - 2022-07-15

<a name="v2.2.0"></a>
## [v2.2.0] - 2022-07-01

<a name="v2.1.0"></a>
## [v2.1.0] - 2022-06-17
### Pull Requests
- Merge branch 'tuoyun'
- Merge branch 'tuoyun'
- Merge branch 'tuoyun'


<a name="v2.0.10"></a>
## [v2.0.10] - 2022-05-13

<a name="v2.0.9"></a>
## [v2.0.9] - 2022-04-29
### Reverts
- update etcd to v3.5.2 ([#206](https://github.com/OpenIMSDK/Open-IM-Server/issues/206))

### Pull Requests
- Merge branch 'tuoyun'
- Merge branch 'tuoyun'
- Merge branch 'tuoyun'
- Merge branch 'tuoyun'
- Merge branch 'tuoyun'
- Merge branch 'tuoyun'


<a name="v2.0.8"></a>
## [v2.0.8] - 2022-04-24
### Pull Requests
- Merge branch 'tuoyun'
- Merge branch 'tuoyun'


<a name="v2.0.7"></a>
## [v2.0.7] - 2022-04-08
### Pull Requests
- Merge branch 'tuoyun'
- Merge branch 'tuoyun'
- Merge branch 'tuoyun'


<a name="v2.0.6"></a>
## [v2.0.6] - 2022-04-01
### Pull Requests
- Merge branch 'tuoyun'


<a name="v2.0.5"></a>
## [v2.0.5] - 2022-03-24

<a name="v2.04"></a>
## [v2.04] - 2022-03-18

<a name="v2.0.3"></a>
## [v2.0.3] - 2022-03-11

<a name="v2.0.2"></a>
## [v2.0.2] - 2022-03-04
### Pull Requests
- Merge branch 'tuoyun'
- Merge branch 'tuoyun'
- Merge branch 'tuoyun'
- Merge branch 'tuoyun'


<a name="v2.0.1"></a>
## [v2.0.1] - 2022-02-25

<a name="v2.0.0"></a>
## [v2.0.0] - 2022-02-23

<a name="v1.0.7"></a>
## [v1.0.7] - 2021-12-17

<a name="v1.0.6"></a>
## [v1.0.6] - 2021-12-10

<a name="v1.0.5"></a>
## [v1.0.5] - 2021-12-03

<a name="v1.0.4"></a>
## [v1.0.4] - 2021-11-25

<a name="v1.0.3"></a>
## [v1.0.3] - 2021-11-12

<a name="v1.0.1"></a>
## [v1.0.1] - 2021-11-04

<a name="v1.0.0"></a>
## v1.0.0 - 2021-10-28
### Reverts
- friend modify
- update


[Unreleased]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v3.0.0...HEAD
[v3.0.0]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.9.0+1.839643f...v3.0.0
[v2.9.0+1.839643f]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.9.0+2.35f07fe...v2.9.0+1.839643f
[v2.9.0+2.35f07fe]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.9.0+1.b5072b1...v2.9.0+2.35f07fe
[v2.9.0+1.b5072b1]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.9.0+3.2667a3a...v2.9.0+1.b5072b1
[v2.9.0+3.2667a3a]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.9.0+7.04818ca...v2.9.0+3.2667a3a
[v2.9.0+7.04818ca]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.9.0...v2.9.0+7.04818ca
[v2.9.0]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v0.0.0+1.3714b4f...v2.9.0
[v0.0.0+1.3714b4f]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v0.0.0+635.8b92c90...v0.0.0+1.3714b4f
[v0.0.0+635.8b92c90]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v0.0.0+1.78a6d03...v0.0.0+635.8b92c90
[v0.0.0+1.78a6d03]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v0.0.0+2.e057c18...v0.0.0+1.78a6d03
[v0.0.0+2.e057c18]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v0.0.0+630.b55ac4a...v0.0.0+2.e057c18
[v0.0.0+630.b55ac4a]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.3.3...v0.0.0+630.b55ac4a
[v2.3.3]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.3.2...v2.3.3
[v2.3.2]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.3.0-rc2...v2.3.2
[v2.3.0-rc2]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.3.0-rc1...v2.3.0-rc2
[v2.3.0-rc1]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.3.0-rc0...v2.3.0-rc1
[v2.3.0-rc0]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.2.0...v2.3.0-rc0
[v2.2.0]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.1.0...v2.2.0
[v2.1.0]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.0.10...v2.1.0
[v2.0.10]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.0.9...v2.0.10
[v2.0.9]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.0.8...v2.0.9
[v2.0.8]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.0.7...v2.0.8
[v2.0.7]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.0.6...v2.0.7
[v2.0.6]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.0.5...v2.0.6
[v2.0.5]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.04...v2.0.5
[v2.04]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.0.3...v2.04
[v2.0.3]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.0.2...v2.0.3
[v2.0.2]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.0.1...v2.0.2
[v2.0.1]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v2.0.0...v2.0.1
[v2.0.0]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v1.0.7...v2.0.0
[v1.0.7]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v1.0.6...v1.0.7
[v1.0.6]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v1.0.5...v1.0.6
[v1.0.5]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v1.0.4...v1.0.5
[v1.0.4]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v1.0.3...v1.0.4
[v1.0.3]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v1.0.1...v1.0.3
[v1.0.1]: https://github.com/OpenIMSDK/Open-IM-Server/compare/v1.0.0...v1.0.1
