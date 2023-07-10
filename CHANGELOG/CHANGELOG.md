# Changelog

All notable changes to this project will be documented in this file.


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
