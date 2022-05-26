#### openIM k8s部署文档
### 1. 修改配置文件
在Open-IM-SERVER目录下修改config/config.yaml配置文件, 将MySQL, Kafka, MongoDB等配置修改。

### 2. 项目根目录创建im configMap到k8s openim namespace
kubectl create namespace openim
kubectl -n openim create configmap config --from-file=config/config.yaml
openim 为im项目的namespace, 可选
查看configmap
kubectl -n openim get configmap

### 3(可选). 修改每个deployment.yml
kubectl get nodes
kubectl label node k8s-node1 role=kube-Node
应需要调度的node打上标签
nodeSelector:
  node: kube-Node
创建资源清单时添加上nodeSelector属性对应即可
修改每种服务数量，建议至少每种2个rpc。
如果修改了config/config.yaml某些配置比如端口，同时需要修改对应deployment端口和ingress端口


### 4. 修改ingress.yaml配置文件
需要安装ingress controller 这里使用的是ingress-nginx 其他ingress需要修改配置文件
进行域名修改等操作

### 5. 执行./kubectl_start.sh脚本
chmod +x ./kubectl_start.sh ./kubectl_stop.sh
./kubectl_start.sh
kubectl -n openim apply -f ingress.yaml
kubectl 启动所有deployment，services，ingress

### 6. 查看k8s deployment service ingress状态
kubectl -n openim get services
kubectl -n openim get deployment
kubectl -n openim get ingress
kubectl -n openim get pods
