source ../.env


# Check if PASSWORD only contains letters and numbers
if [[ "$PASSWORD" =~ ^[a-zA-Z0-9]+$ ]]
then
    echo "PASSWORD is valid."
else
    echo "ERR: PASSWORD should only contain letters and numbers. " $PASSWORD
    exit
fi


echo "your user is:$USER"
echo "your password is:$PASSWORD"
echo "your minio endPoint is:$MINIO_ENDPOINT"
echo "your data dir is $DATA_DIR"


#!/bin/bash

# Specify the config file
config_file='../config/config.yaml'

# Load variables from .env file
source ../.env

# Replace the password and username field for mysql
sed -i "/mysql:/,/database:/ s/password:.*/password: $PASSWORD/" $config_file
sed -i "/mysql:/,/database:/ s/username:.*/username: $USER/" $config_file

# Replace the password and username field for mongo
sed -i "/mongo:/,/maxPoolSize:/ s/password:.*/password: $PASSWORD/" $config_file
sed -i "/mongo:/,/maxPoolSize:/ s/username:.*/username: $USER/" $config_file

# Replace the password field for redis
sed -i '/redis:/,/password:/s/password: .*/password: '${PASSWORD}'/' $config_file

# Replace accessKeyID and secretAccessKey for minio
sed -i "/minio:/,/isDistributedMod:/ s/accessKeyID:.*/accessKeyID: $USER/" $config_file
sed -i "/minio:/,/isDistributedMod:/ s/secretAccessKey:.*/secretAccessKey: $PASSWORD/" $config_file
sed -i '/minio:/,/endpoint:/s|endpoint: .*|endpoint: '${MINIO_ENDPOINT}'|' $config_file
