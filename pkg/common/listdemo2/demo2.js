
db.friend_version.updateMany(
    {
        "d_id": "100"
    },
    [
        {
            $addFields: {
                update_elem_ids: ["1000", "1001","1003", "2000"]
            }
        },
        {
            $set: {
                version: {
                    $add: ["$version", 1]
                },
                last_update: new Date(),
            }
        },
        {
            $set: {
                logs: {
                    $filter: {
                        input: "$logs",
                        as: "log",
                        cond: {
                            "$not": {
                                $in: ["$$log.e_id", "$update_elem_ids"]
                            }
                        }
                    }
                },

            },

        },
        {
            $set: {
                logs: {
                    $concatArrays: [
                        "$logs",
                        [
                            {
                                e_id: "1003",
                                last_update: ISODate("2024-05-25T06:32:10.238Z"),
                                version: "$version",
                                deleted: false
                            },

                        ]
                    ]
                }
            }
        },
        {
            $unset: ["update_elem_ids"]
        },

    ]
)



