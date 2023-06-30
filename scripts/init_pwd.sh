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

echo "your user is:$USER"
echo "your password is:$PASSWORD"
echo "your minio endPoint is:$MINIO_ENDPOINT"
echo "your data dir is $DATA_DIR"

sed -i "/^\([[:space:]]*dbMysqlUserName: *\).*/s//\1$USER/;0,/\([[:space:]]*dbUserName: *\).*/s//\1 $USER/;/\([[:space:]]*accessKeyID: *\).*/s//\1 $USER/;/\([[:space:]]*endpoint: *\).*/s//\1\"abc\"/;" ../config/config.yaml
sed -i "/^\([[:space:]]*dbMysqlPassword: *\).*/s//\1$PASSWORD/;/\([[:space:]]*dbPassword: *\).*/s//\1$PASSWORD/;/\([[:space:]]*secret: *\).*/s//\1$PASSWORD/;/\([[:space:]]*secretAccessKey: *\).*/s//\1$PASSWORD/;" ../config/config.yaml

sed -i "/\([[:space:]]*endpoint: *\).*/s##\1$MINIO_ENDPOINT#;" ../config/config.yaml
sed -i "/\([[:space:]]*dbPassWord: *\).*/s//\1$PASSWORD/;" ../config/config.yaml
sed -i "/\([[:space:]]*secret: *\).*/s//\1$PASSWORD/;" ../.docker-compose_cfg/config.yaml
