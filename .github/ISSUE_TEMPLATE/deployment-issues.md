---
name: Deployment issues
about: Deployment issues
title: ''
labels: ''
assignees: ''

---

If you are deploying OpenIM for the first time



```
git clone https://github.com/OpenIMSDK/Open-IM-Server.git --recursive
```

screenshot here

```
cd Open-IM-Server/script ; chmod +x *.sh ; ./env_check.sh
```

screenshot here

```
cd .. ; docker-compose up -d
```

screenshot here

```
cd script ; ./docker_check_service.sh
```

screenshot here
