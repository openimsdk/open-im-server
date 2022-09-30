
#### 镜像
```
windows下 程序编译  win_cross_build_all_service.cmd
windows下 镜像编译  win_build_push_all_images.cmd
```

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
