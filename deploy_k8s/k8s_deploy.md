#### im k8s部署文档
### 1.修改配置文件
在Open-IM-SERVER目录下修改config/config.yaml配置文件

### 2. 项目根目录创建im configMap到k8s
kubectl -n {namespace} create configmap config --from-file=config/config.yaml
namespace 为im项目的namespace

### 3. 

