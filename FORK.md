# AndChat OpenIM Fork Policy

This fork backs the AndChat in-house OpenIM deployment.

## Remotes

```text
upstream = https://github.com/openimsdk/open-im-server.git
origin   = git@github.com:yuanlu-and/open-im-server.git
```

## Branches

```text
upstream/main
- pristine OpenIM source

origin/main
- mirror of upstream/main

origin/andchat-prod
- production branch for the AndChat OpenIM deployment
```

Keep `main` close to upstream. Put AndChat deployment patches on `andchat-prod`.

## What Belongs In This Fork

Acceptable changes:

- deployment and configuration defaults
- Docker/build adjustments
- observability and logging
- small OpenIM bug fixes
- chat-engine behavior that must live inside OpenIM
- narrow extension points that can be proposed upstream

Avoid adding:

- phone OTP auth
- phone-number discovery
- app profiles
- billing
- AI memory or retrieval logic
- product-specific permissions
- reaction source-of-truth tables unless reactions become OpenIM-native

Those features belong in `andchat-backend`.

## Current Divergence

No runtime behavior changes yet.

This file is the initial fork policy document for the `andchat-prod` branch.

## Upgrade Flow

```bash
git fetch upstream --tags
git checkout main
git merge --ff-only upstream/main
git push origin main

git checkout andchat-prod
git merge origin/main
git push origin andchat-prod
```

Record every intentional production patch in this file or in a dedicated section below.
