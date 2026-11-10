# Rolling out monotonic read sequences

Read progress is shared by a user across devices. Redis updates, cache fills, and
MongoDB writes must retain the greatest confirmed `read_seq` for each
`(user_id, conversation_id)`. Read-receipt notifications continue to use the
existing asynchronous queue and MQ pipeline.

## MongoDB preflight

The server creates a unique compound index named `user_conversation_unique` on
`seq_user`. The explicit name allows it to coexist with the previous non-unique
index on MongoDB 5.0 and later (the documented deployment uses MongoDB 6.0).
The old index is not automatically dropped.

Before deploying, check for duplicate rows in the intended OpenIM database:

```javascript
db.seq_user.aggregate([
  { $group: {
      _id: { user_id: "$user_id", conversation_id: "$conversation_id" },
      count: { $sum: 1 }
  } },
  { $match: { count: { $gt: 1 } } },
  { $limit: 20 }
])
```

Resolve duplicates through a reviewed data migration before rollout, preserving
the greatest confirmed `read_seq` and reviewing the independent `min_seq` and
`max_seq` bounds. Coordinate the migration with writers. The application does not
delete or merge historical records automatically; index creation fails explicitly
if duplicates remain. Check index build time and permissions for the deployment.
If an equivalent unique index has already been installed manually, reconcile its
name with `user_conversation_unique` before rollout.

## Redis compatibility and limits

Read-sequence keys keep the existing Redis hash `value` field and 30-day write
TTL. Their dedicated monotonic fill path replaces generic RocksCache fills;
updating a value also invalidates an older RocksCache fill owner. Complete the
rollout across all writers, since old processes can still perform unconditional
writes. Other cached fields retain their existing cache behavior.

Atomic maximum updates tolerate concurrent devices, duplicate receipts, and
out-of-order persistence. They do not repair a previously lost persistence event
or guarantee delivery through an in-memory queue. Reliable retry/compensation and
reconciliation of existing Redis/MongoDB differences remain separate work. Do not
mark all messages as read to repair a discrepancy; use confirmed read progress.
