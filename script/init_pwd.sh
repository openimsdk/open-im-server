echo "your user is:$user"
echo "your password is:$password"
echo "your minio endPoint is:$minio_endpoint"

sed -i "/^\([[:space:]]*dbMysqlUserName: *\).*/s//\1$user/;0,/\([[:space:]]*dbUserName: *\).*/s//\1 $user/;/\([[:space:]]*accessKeyID: *\).*/s//\1$user/;/\([[:space:]]*endpoint: *\).*/s//\1\"abc\"/;" ../config/config.yaml
sed -i "/^\([[:space:]]*dbMysqlPassword: *\).*/s//\1$password/;/\([[:space:]]*dbPassword: *\).*/s//\1$password/;/\([[:space:]]*secret: *\).*/s//\1$password/;/\([[:space:]]*secretAccessKey: *\).*/s//\1$PASSWORD/;" ../config/config.yaml

sed -i "/\([[:space:]]*endpoint: *\).*/s##\1$minio_endpoint#;" ../config/config.yaml
sed -i "/\([[:space:]]*dbPassWord: *\).*/s//\1$password/;" ../config/config.yaml
sed -i "/\([[:space:]]*secret: *\).*/s//\1$password/;" ../docker-compose_cfg/config.yaml
