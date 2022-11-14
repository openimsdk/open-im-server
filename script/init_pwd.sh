#!/usr/bin/env bash

source ../.env

nameList=(dbMysqlUserName, dbUserName, dbUserName, accessKeyID)
pwdList=(dbMysqlPassword, dbPassword, dbPassWord, secretAccessKey)

for i in ${nameList[*]}; do
 sed -i 's/{$i}: [a-z]/{$i}: {$USER}/g' ../config/usualConfig.yaml
done

for i in ${pwdList[*]}; do
 sed -i 's/{$i}: [a-z]/{$i}: {$PASSWORD}/g' ../config/usualConfig.yaml
done