#!/usr/bin/env bash

source ./style_info.cfg
source ./path_info.cfg
source ./function.sh

#define database attributes
address=$(cat $config_path | grep -w dbMysqlAddress)
list_to_string ${address}
hostAndPort=($ports_array)
DATABASE_HOST=$(echo $hostAndPort | awk -F '[:]' '{print $1}')
DATABASE_PORT=$(echo $hostAndPort | awk -F '[:]' '{print $NF}')
DATABASE_USERNAME=$(cat $config_path | grep -w dbMysqlUserName | awk -F '[:]' '{print $NF}')
DATABASE_PWD=`eval echo $(cat $config_path | grep -w dbMysqlPassword | awk -F '[:]' '{print $NF}')`
DATABASE_NAME=$(cat $config_path | grep -w dbMysqlDatabaseName | awk -F '[:]' '{print $NF}')
SQL_FILE="../config/mysql_sql_file/openIM.sql"


create_data_sql="create database IF NOT EXISTS $DATABASE_NAME"
set_character_code_sql="alter database $DATABASE_NAME character set utf8mb4 collate utf8mb4_general_ci"

echo -e "${SKY_BLUE_PREFIX}start to create database.....$COLOR_SUFFIX"
mysql -h $DATABASE_HOST -P $DATABASE_PORT -u $DATABASE_USERNAME -p$DATABASE_PWD -e "$create_data_sql"

if [ $? -eq 0 ]; then
  echo -e "${SKY_BLUE_PREFIX}create database ${DATABASE_NAME} successfully$COLOR_SUFFIX"
  mysql -h $DATABASE_HOST -P $DATABASE_PORT -u $DATABASE_USERNAME -p$DATABASE_PWD -e "$set_character_code_sql"
else
  echo -e "${RED_PREFIX}create database failed or exists the database$COLOR_SUFFIX\n"
fi

echo -e "${SKY_BLUE_PREFIX}start to source openIM.sql .....$COLOR_SUFFIX"
mysql -h $DATABASE_HOST -P $DATABASE_PORT -u $DATABASE_USERNAME -p$DATABASE_PWD -D $DATABASE_NAME <$SQL_FILE
if [ $? -eq 0 ]; then
  echo -e "${SKY_BLUE_PREFIX}source openIM.sql successfully$COLOR_SUFFIX"
else
  echo -e "${RED_PREFIX}source openIM.sql failed$COLOR_SUFFIX\n"
fi
