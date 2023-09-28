# OpenIM Server Chat

## 目录结构

```bash
openim-server/
  Chart.yaml          # 包含了chart信息的YAML文件
  LICENSE             # 包含OpenIM Chart许可证的纯文本文件
  README.md           # OpenIM 可读的README文件
  values.yaml         # chart 默认的配置值
  charts/             # 包含chart依赖的其他chart
  crds/               # 自定义资源的定义
  templates/          # 模板目录， 当和values 结合时，可生成有效的Kubernetes manifest文件
  templates/NOTES.txt # 包含简要使用说明的纯文本文件
```