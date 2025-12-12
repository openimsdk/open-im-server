#!/usr/bin/env bash
set -e

echo "Creating OpenIM user in database ${MONGO_INITDB_DATABASE}..."

mongosh -u "${MONGO_INITDB_ROOT_USERNAME}" -p "${MONGO_INITDB_ROOT_PASSWORD}" --authenticationDatabase admin <<EOF
use ${MONGO_INITDB_DATABASE}
if (!db.getUser("${MONGO_OPENIM_USERNAME}")) {
  db.createUser({
    user: "${MONGO_OPENIM_USERNAME}",
    pwd: "${MONGO_OPENIM_PASSWORD}",
    roles: [{role: "readWrite", db: "${MONGO_INITDB_DATABASE}"}]
  })
  print("OpenIM user created successfully")
} else {
  print("User ${MONGO_OPENIM_USERNAME} already exists")
}
EOF

echo "OpenIM user setup completed"
