<!-- vscode-markdown-toc -->

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->
# Install Docker


The installation command is as follows:

```bash
$ curl -fsSL https://get.docker.com | bash -s docker --mirror aliyun
``

## 2.2 Start Docker

```bash
$ systemctl start docker
```

## 2.3 Test Docker

```bash
$ docker run hello-world
```

## 2.4 Configure Docker Acceleration

```bash
$ mkdir -p /etc/docker
$ tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": ["https://registry.docker-cn.com"]
}
EOF
$ systemctl daemon-reload
$ systemctl restart docker
```

## 2.5 Install Docker Compose

```bash
$ sudo curl -L "https://github.com/docker/compose/releases/download/latest/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
$ sudo chmod +x /usr/local/bin/docker-compose
```

## 2.6 Test Docker Compose

```bash
$ docker-compose --version
```
