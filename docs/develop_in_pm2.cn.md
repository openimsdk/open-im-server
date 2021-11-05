## 使用 PM2 管理开发多个服务

### 安装 Node.js 和 PM2

```
curl -sL https://deb.nodesource.com/setup_14.x | sudo -E bash -

npm install pm2 -g
```


### PM2 开发

```

pm2 start pm2.yaml
```

### PM2 检查进程

```
pm2 ls

┌─────┬─────────────────────────┬─────────────┬─────────┬─────────┬──────────┬────────┬──────┬───────────┬──────────┬──────────┬──────────┐
│ id  │ name                    │ namespace   │ version │ mode    │ pid      │ uptime │ ↺    │ status    │ cpu      │ mem      │ watching │
├─────┼─────────────────────────┼─────────────┼─────────┼─────────┼──────────┼────────┼──────┼───────────┼──────────┼──────────┼──────────┤
│ 0   │ open_im_api             │ default     │ N/A     │ fork    │ 38641    │ 74s    │ 0    │ online    │ 0%       │ 32.5mb   │ disabled │
│ 1   │ open_im_auth            │ default     │ N/A     │ fork    │ 38642    │ 74s    │ 0    │ online    │ 0%       │ 30.4mb   │ disabled │
│ 3   │ open_im_friend          │ default     │ N/A     │ fork    │ 38644    │ 74s    │ 0    │ online    │ 0%       │ 37.9mb   │ disabled │
│ 4   │ open_im_group           │ default     │ N/A     │ fork    │ 40594    │ 0s     │ 46   │ online    │ 0%       │ 17.1mb   │ disabled │
│ 2   │ open_im_msg             │ default     │ N/A     │ fork    │ 38643    │ 74s    │ 0    │ online    │ 0%       │ 35.8mb   │ disabled │
│ 9   │ open_im_msg_gateway     │ default     │ N/A     │ fork    │ 38666    │ 74s    │ 0    │ online    │ 0%       │ 33.4mb   │ disabled │
│ 8   │ open_im_msg_transfer    │ default     │ N/A     │ fork    │ 38660    │ 74s    │ 0    │ online    │ 0%       │ 34.7mb   │ disabled │
│ 6   │ open_im_push            │ default     │ N/A     │ fork    │ 38647    │ 74s    │ 0    │ online    │ 0%       │ 33.7mb   │ disabled │
│ 7   │ open_im_timed_task      │ default     │ N/A     │ fork    │ 38657    │ 74s    │ 0    │ online    │ 0%       │ 27.3mb   │ disabled │
│ 5   │ open_im_user            │ default     │ N/A     │ fork    │ 38646    │ 74s    │ 0    │ online    │ 0%       │ 34.2mb   │ disabled │
└─────┴─────────────────────────┴─────────────┴─────────┴─────────┴──────────┴────────┴──────┴───────────┴──────────┴──────────┴──────────┘

```

### PM2 检查日志

```
pm2 logs
```

### PM2 删除进程

```
pm2 delete

pm2 flush
```