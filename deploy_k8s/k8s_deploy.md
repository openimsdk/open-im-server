#### im k8s部署文档
### 1.修改配置文件
在Open-IM-SERVER目录下修改config/config.yaml配置文件

### 2. 项目根目录创建im configMap
kubectl create configmap config config/config.yaml

### 3. 