# OpenIM Web E2E

Minimal local web page for exercising OpenIM HTTP message APIs from a browser.

Usage:

```bash
cd /Users/ren_yu/open-source/open-im-server/test/e2e/web
python3 -m http.server 18080
```

Open [http://127.0.0.1:18080](http://127.0.0.1:18080) after the OpenIM server is running.

Pages:

- `http://127.0.0.1:18080/` - token and message API tester
- `http://127.0.0.1:18080/chat.html` - simplest dev chat client
- `http://127.0.0.1:18080/ws-demo.html` - simplest native browser OpenIM WS client

Default local endpoints:

- API: `http://127.0.0.1:10002`
- WebSocket gateway: `ws://127.0.0.1:10001`

Defaults from this repo:

- admin user ID: `imAdmin`
- secret: `openIM123`

Notes:

- `/auth/get_admin_token` requires the admin user to already exist.
- `/msg/send_msg` requires an admin token in the `token` header.
- The included page currently sends plain text messages through `/msg/send_msg`.
- `chat.html` is a development-only client: it uses admin-backed APIs to send and search messages between two users.
- `ws-demo.html` performs a real browser WebSocket handshake and sends OpenIM WS envelopes over binary frames.
- The WS gateway in this repo currently expects `isBackground=false` in the query string; omitting it causes the handshake to fail before upgrade.

Debugging:

```bash
cd /Users/ren_yu/open-source/open-im-server
./scripts/debug_msggateway_ws.sh before-rpc
```

- `gateway` stops at the WS gateway after the envelope is decoded.
- `before-rpc` stops right before `g.msgClient.MsgClient.SendMsg(...)`.
- `after-rpc` stops right after `g.msgClient.MsgClient.SendMsg(...)` returns.
- These breakpoints only apply to `ws-demo.html`, not `chat.html` or `/msg/send_msg`.
