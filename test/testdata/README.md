
# Test Data for OpenIM Server

This directory (`testdata`) contains various JSON formatted data files that are used for testing the OpenIM Server.

## Structure

```bash
testdata/
│
├── README.md         # 描述该目录下各子目录和文件的作用
│
├── storage/              # 存储模拟的数据库数据
│   ├── users.json   # 用户的模拟数据
│   └── messages.json # 消息的模拟数据
│
├── requests/        # 存储模拟的请求数据
│   ├── login.json   # 模拟登陆请求
│   ├── register.json # 模拟注册请求
│   └── sendMessage.json # 模拟发送消息请求
│
└── responses/       # 存储模拟的响应数据
    ├── login.json   # 模拟登陆响应
    ├── register.json # 模拟注册响应
    └── sendMessage.json # 模拟发送消息响应
```

Here is an overview of what each subdirectory or file represents:

- `db/` - This directory contains mock data mimicking the actual database contents.
  - `users.json` - Represents a list of users in the system. Each entry contains user-specific information such as user ID, username, password hash, etc.
  - `messages.json` - Contains a list of messages exchanged between users. Each message entry includes the sender's and receiver's user IDs, message content, timestamp, etc.
- `requests/` - This directory contains mock requests that a client might send to the server.
  - `login.json` - Represents a user login request. It includes fields such as username and password.
  - `register.json` - Mimics a user registration request. Contains details such as username, password, email, etc.
  - `sendMessage.json` - Simulates a message sending request from a user to another user.
- `responses/` - This directory holds the expected server responses for the respective requests.
  - `login.json` - Represents a successful login response from the server. It typically includes a session token and user-specific information.
  - `register.json` - Simulates a successful registration response from the server, usually containing the new user's ID, username, etc.
  - `sendMessage.json` - Depicts a successful message sending response from the server, confirming the delivery of the message.

## JSON Format

All the data files in this directory are in JSON format. JSON (JavaScript Object Notation) is a lightweight data-interchange format that is easy for humans to read and write and easy for machines to parse and generate.

Here is a simple example of what a JSON file might look like:

```bash
  "users": [
    {
      "id": 1,
      "username": "user1",
      "password": "password1"
    },
    {
      "id": 2,
      "username": "user2",
      "password": "password2"
    }
  ]

```

In this example, "users" is an array of user objects. Each user object has an "id", "username", and "password".
