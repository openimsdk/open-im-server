db.friend_version.updateMany(
    {
        "d_id": "100"
    },
    [
        {
            $addFields: {
                elem_index: {
                    $indexOfArray: [
                        "$logs.e_id",
                        "1000"
                    ]
                }
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
                    $cond: {
                        if: {
                            $lt: ["$elem_index", 0]
                        },
                        then: {
                            $concatArrays: [
                                "$logs",
                                [
                                    {
                                        e_id: "1000",
                                        last_update: new Date(),
                                        version: "$version",
                                        deleted: false
                                    }
                                ]
                            ]
                        },
                        else: {
                            $map: {
                                input: {
                                    $range: [0, {
                                        $size: "$logs"
                                    }]
                                },
                                as: "i",
                                in: {
                                    $cond: {
                                        if: {
                                            $eq: ["$$i", "$elem_index"]
                                        },
                                        then: {
                                            e_id: "1000",
                                            last_update: new Date(),
                                            version: "$version",
                                            deleted: false
                                        },
                                        else: {
                                            $arrayElemAt: ["$logs", "$$i"]
                                        }
                                    },

                                },

                            },

                        },

                    },

                },

            },

        },
        {
            $unset: ["elem_index"]
        },
    ]
)
