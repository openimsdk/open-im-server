#!/usr/bin/env bash
# Copyright © 2023 OpenIM. All rights reserved.
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

# input: [10023, 2323, 3434]
# output: 10023 2323 3434

# 函数功能：将列表转换为字符串，去除空格和括号
list_to_string() {
    ports_list=$*  # 获取传入的参数列表
    sub_s1=$(echo $ports_list | sed 's/ //g')  # 去除空格
    sub_s2=${sub_s1//,/ }  # 将逗号替换为空格
    sub_s3=${sub_s2#*[}  # 去除左括号及其之前的内容
    sub_s4=${sub_s3%]*}  # 去除右括号及其之后的内容
    ports_array=$sub_s4  # 将处理后的字符串赋值给变量 ports_array
}

# 函数功能：去除字符串中的空格
remove_space() {
    value=$*  # 获取传入的参数
    result=$(echo $value | sed 's/ //g')  # 去除空格
}
