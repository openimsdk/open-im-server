#### openIM k8s部署文档
### 1. 修改配置文件
在Open-IM-SERVER根目录下修改config/config.yaml配置文件, 请确保以下修改的所有地址必须保证k8s pod能够访问
1. 修改ETCD配置为自己的ETCD ip地址, 最好和k8s本身使用的ETCD分开
2. 修改MySQL配置
3. 修改Mongo配置
4. 修改Redis配置
5. 修改Kafka配置
6. 将rpcRegisterIP修改为空, 此地址为每个rpc注册到ETCD的地址, 置空每个rpc将会将pod地址注册到ETCD, 才能正确rpc请求(重要)
7. 如果使用minio作为对象存储, 还需要修改minio的地址
8. 其他如果使用离线推送,需要修改push离线推送配置
9. 修改demo中的imAPIURL字段为openIM api的ingress或者service地址, 需要让demo的pod能正确请求到(重要)
10. 其他非必须配置修改, 如短信,推送等

### 2. 项目根目录创建im configMap到k8s openim namespace
1. 为open-IM项目创建单独命名空间
 ```  
    kubectl create namespace openim
 ```  
2. 在项目根目录通过config/config.yaml  
 ```  
    kubectl -n openim create configmap config --from-file=config/config.yaml
 ```
    查看configmap
 ```
   kubectl -n openim get configmap
 ``` 

### 3(可选). 修改每个deployment.yml
  每个rpc的deployment在Open-IM-SERVER根目录deploy_k8s下
  给需要调度的node打上标签
 ```
    kubectl get nodes
        kubectl label node k8s-node1 role=openIMworker
 ```
  在deployment的spec.template.spec加上
 ```
    nodeSelector:
     role: openIMworker
 ``` 
   创建资源清单时添加上nodeSelector属性对应即可,
   修改每种服务数量，建议至少每种2个rpc。
   如果修改了config/config.yaml某些配置比如端口，同时需要修改对应deployment端口和ingress端口


### 4. 修改ingress.yaml配置文件
1. 需要安装ingress controller, 这里使用的是ingress-nginx, 使用其他类型的ingress controller需要更改ingress.class, 将host改为自己部署服务的host

### 5. 执行./kubectl_start.sh脚本
1. 脚本给予可执行权限
 ```
    chmod +x ./kubectl_start.sh ./kubectl_stop.sh
 ```
2. 启动k8s service和deployment 
 ```
    ./kubectl_start.sh
 ``` 
3. 启动k8s ingress
 ```
    kubectl -n openim apply -f ingress.yaml
 ```
kubectl 启动所有deployment, services, ingress

### 6. 查看k8s deployment service ingress状态

 ```
    kubectl -n openim get services
    kubectl -n openim get deployment
    kubectl -n openim get ingress
    kubectl -n openim get pods
 ```
 检测服务可达
 ```
    telnet msg-gateway.openim.xxx.com {{your_ingress_port}}
    telnet sdk-server.openim.xxx.com {{your_ingress_port}}
    telnet api.openim.xxx.com {{your_ingress_port}}
    telnet cms-api.openim.xxx.com {{your_ingress_port}}
    telnet demo.openim.xxx.com {{your_ingress_port}}
 ```

#### openIM k8s更新
1. 暂存配置文件，拉取代码
 ```
    git stash push config/config.yaml
    git pull
 ```
2. 合并配置文件, 解决冲突
 ```
    git stash pop
 ```
3. 重新生成configmap
 ```
    kubectl -n openim create configmap config --from-file=config/config.yaml
 ```
4.修改所有deployment文件的spec.template.spec.image 改为新版本后在/deploy_k8s下重新执行
 ```
    ./kubectl_stop.sh
    ./kubectl_start.sh
 ```