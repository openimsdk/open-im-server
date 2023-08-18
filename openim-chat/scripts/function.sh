#!/usr/bin/env bash


# Copyright © 2023 OpenIM open source community. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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