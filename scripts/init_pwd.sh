echo "your user is:$USER"
echo "your password is:$PASSWORD"
echo "your minio endPoint is:$MINIO_ENDPOINT"
echo "your data dir is $DATA_DIR"

sed -i "/^\([[:space:]]*dbMysqlUserName: *\).*/s//\1$USER/;0,/\([[:space:]]*dbUserName: *\).*/s//\1 $USER/;/\([[:space:]]*accessKeyID: *\).*/s//\1 $USER/;/\([[:space:]]*endpoint: *\).*/s//\1\"abc\"/;" ../config/config.yaml
sed -i "/^\([[:space:]]*dbMysqlPassword: *\).*/s//\1$PASSWORD/;/\([[:space:]]*dbPassword: *\).*/s//\1$PASSWORD/;/\([[:space:]]*secret: *\).*/s//\1$PASSWORD/;/\([[:space:]]*secretAccessKey: *\).*/s//\1$PASSWORD/;" ../config/config.yaml

sed -i "/\([[:space:]]*endpoint: *\).*/s##\1$MINIO_ENDPOINT#;" ../config/config.yaml
sed -i "/\([[:space:]]*dbPassWord: *\).*/s//\1$PASSWORD/;" ../config/config.yaml
sed -i "/\([[:space:]]*secret: *\).*/s//\1$PASSWORD/;" ../.docker-compose_cfg/config.yaml
