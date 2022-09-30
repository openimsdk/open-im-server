#### 配置
```
请根据需要配置 values.yaml
#创建 k8s namespace
kubectl create namespace openim-ns

#创建/覆盖 configmap配置
kubectl -n openim-ns create configmap config --from-file=config/config.yaml

#查看 
kubectl -n openim get configmap

```
#### 安装
```
cd helm
helm install my-openim -f openim/values.yaml -n openim-ns openim
```
#### 更新
```
cd helm
helm upgrade my-openim -f openim/values.yaml -n openim-ns openim
```
#### 卸载
```
cd helm
helm uninstall  my-openim -n openim
```

#### 镜像
```
镜像编译  deploy_k8s/build_push_all_images.sh
windows镜像编译  deploy_k8s/win_build_push_all_images.cmd

```