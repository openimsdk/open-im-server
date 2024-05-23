
db.friend_version.aggregate([
    {
        "$match": {
            "d_id": "100",
        }
    },
    {
        "$project": {
            "_id": 0,
            "d_id": 0,
        }
    },
    {
        "$addFields": {
            "logs": {
                $cond: {
                    if: {
                        $or: [
                            {$lt: ["$version", 3]},
                            {$gte: ["$deleted", 3]},
                        ],
                    },
                    then: [],
                    else: "$logs",
                }
            }
        },
    },
    {
        "$addFields": {
            "logs": {
                "$filter": {
                    input: "$logs",
                    as: "l",
                    cond: { $gt: ["$$l.version", 3] }
                }
            }
        }
    },
    {
        "$addFields": {
            "log_len": {
                $size: "$logs"
            }
        }
    },
    {
        "$addFields": {
            "logs": {
                $cond: {
                    if: {$gt: ["$log_len", 1]},
                    then: [],
                    else: "$logs",
                }
            }
        }
    }
])
