# After s3 switches the storage engine, convert the data

- build
```shell
go build -o s3convert main.go
```

- start
```shell
./s3convert -config <config dir path> -name <old s3 name>
# ./s3convert -config ./../../config -name minio
```
