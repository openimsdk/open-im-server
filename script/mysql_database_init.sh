#!/usr/bin/env bash

source ./style_info.cfg

#define database attributes
DATABASE_HOST="127.0.0.1" #your mysql host
DATABASE_USERNAME="root"  #your mysql username
DATABASE_PWD="123456"     #your mysql password
DATABASE_NAME="openIM"
SQL_FILE="../config/mysql_sql_file/openIM.sql"

create_data_sql="create database $DATABASE_NAME"
set_character_code_sql="alter database $DATABASE_NAME character set utf8mb4 collate utf8mb4_general_ci"

echo -e "${SKY_BLUE_PREFIX}start to create database.....$COLOR_SUFFIX"
mysql -h $DATABASE_HOST -u $DATABASE_USERNAME -p$DATABASE_PWD -e "$create_data_sql"

if [ $? -eq 0 ]; then
  echo -e "${SKY_BLUE_PREFIX}create database OpenIM successfully$COLOR_SUFFIX"
  mysql -h $DATABASE_HOST -u $DATABASE_USERNAME -p$DATABASE_PWD -e "$set_character_code_sql"
else
  echo -e "${RED_PREFIX}create database failed or exists the database$COLOR_SUFFIX"
fi

echo -e "${SKY_BLUE_PREFIX}start to source openIM.sql .....$COLOR_SUFFIX"
mysql -h $DATABASE_HOST -u $DATABASE_USERNAME -p$DATABASE_PWD -D $DATABASE_NAME <$SQL_FILE
if [ $? -eq 0 ]; then
  echo -e "${SKY_BLUE_PREFIX}source openIM.sql successfully$COLOR_SUFFIX"
else
  echo -e "${RED_PREFIX}source openIM.sql failed$COLOR_SUFFIX"
fi
