# Changelog

All notable changes to this project will be documented in this file.

+ [https://github.com/OpenIMSDK/Open-IM-Server/releases](https://github.com/OpenIMSDK/Open-IM-Server/releases)

## command

```bash
git-chglog --tag-filter-pattern 'v2.0.*'  -o CHANGELOG-2.0.md
```

## create next tag

```bash
git-chglog --next-tag 2.0.0 -o CHANGELOG.md
git commit -am "release 2.0.0"
git tag 2.0.0
```

| Query          | Description                                    | Example                     |
| -------------- | ---------------------------------------------- | --------------------------- |
| `<old>..<new>` | Commit contained in `<new>` tags from `<old>`. | `$ git-chglog 1.0.0..2.0.0` |
| `<name>..`     | Commit from the `<name>` to the latest tag.    | `$ git-chglog 1.0.0..`      |
| `..<name>`     | Commit from the oldest tag to `<name>`.        | `$ git-chglog ..2.0.0`      |
| `<name>`       | Commit contained in `<name>`.                  | `$ git-chglog 1.0.0`        |


## Release version logs

+ [OpenIM CHANGELOG-V1.0](CHANGELOG-1.0.md)
+ [OpenIM CHANGELOG-V2.0](CHANGELOG-2.0.md)
+ [OpenIM CHANGELOG-V2.1](CHANGELOG-2.1.md)
+ [OpenIM CHANGELOG-V2.2](CHANGELOG-2.2.md)
+ [OpenIM CHANGELOG-V2.3](CHANGELOG-2.3.md)
+ [OpenIM CHANGELOG-V2.9](CHANGELOG-2.9.md)
+ [OpenIM CHANGELOG-V3.0](CHANGELOG-3.0.md)