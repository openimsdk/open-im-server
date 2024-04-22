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

mongosh <<EOF
use admin
var rootUsername = '$MONGO_INITDB_ROOT_USERNAME';
var rootPassword = '$MONGO_INITDB_ROOT_PASSWORD';
var authResult = db.auth(rootUsername, rootPassword);
if (authResult) {
  print('Authentication successful for root user: ' + rootUsername);
} else {
  print('Authentication failed for root user: ' + rootUsername + ' with password: ' + rootPassword);
  quit(1);
}

var dbName = '$MONGO_INITDB_DATABASE';
db = db.getSiblingDB(dbName);
var openimUsername = '$MONGO_OPENIM_USERNAME';
var openimPassword = '$MONGO_OPENIM_PASSWORD';
var createUserResult = db.createUser({
  user: openimUsername,
  pwd: openimPassword,
  roles: [
    { role: 'readWrite', db: dbName }
  ]
});

if (createUserResult.ok == 1) {
  print('User creation successful. User: ' + openimUsername + ', Database: ' + dbName);
} else {
  print('User creation failed for user: ' + openimUsername + ' in database: ' + dbName);
  quit(1);
}
EOF


