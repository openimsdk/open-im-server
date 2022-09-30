#### 配置
```
请根据需要配置 values.yaml
#创建 k8s namespace
kubectl create namespace openim-ns
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
