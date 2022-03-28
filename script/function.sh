#!/usr/bin/env bash
#input:[10023,2323,3434]
#output:10023 2323 3434
list_to_string(){
ports_list=$*
sub_s1=`echo $ports_list | sed 's/ //g'`
sub_s2=${sub_s1//,/ }
sub_s3=${sub_s2#*[}
sub_s4=${sub_s3%]*}
ports_array=$sub_s4
}
remove_space(){
  value=$*
  result=`echo $value | sed 's/ //g'`
}