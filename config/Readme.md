配置文件说明

每个组件单独有一个配置文件，主要是地址及账号密码信息。





callback





rpc共同配置项说明

```
rpc:
  #api或其他rpc可通过这个ip访问到此rpc，如果为空则获取内网ip，默认为空即可	
  registerIP: ''
  #监听ip，如果为0.0.0.0，则内外网ip都监听，为空则自动获取内网ip监听
  listenIP: 0.0.0.0
  #监听端口，如果配置多个则会启动多个实例，但需和prometheus.ports保持一致
  ports: [ 10120 ]

prometheus:
  enable: true
  #prometheus
  ports: [ 20104 ]
```



api配置项说明

```
api:
  #0.0.0.0表示内外网ip都监听，不应该修改
  listenIP: 0.0.0.0
  #监听端口，如果配置多个则会启动多个实例,和prometheus.ports保持一致
  ports: [ 10002 ]

prometheus:
  enable: true
  ports: [ 20113 ]
  #怎么描述 grafanaURL地址， 外网地址，通过浏览器能访问到
  grafanaURL: http://127.0.0.1:13000/
```



log配置项说明

```
#log存放路径，如果要修改则改为全路径
storageLocation: ../../../../logs/
rotationTime: 24
remainRotationCount: 2
#3: 生产环境; 6：日志较多，调试环境
remainLogLevel: 6
isStdout: false
isJson: false
withStack: false
```



