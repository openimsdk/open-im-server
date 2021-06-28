SET
FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `account`
-- ----------------------------
DROP TABLE IF EXISTS `account`;
CREATE TABLE `account`
(
    `id`       int(11) NOT NULL AUTO_INCREMENT,
    `uid`      varchar(32)  NOT NULL,
    `account`  varchar(255) NOT NULL,
    `password` varchar(255) NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_account` (`account`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of account
-- ----------------------------

-- ----------------------------
-- Table structure for `black_list`
-- ----------------------------
DROP TABLE IF EXISTS `black_list`;
CREATE TABLE `black_list`
(
    `uid`                varchar(32) NOT NULL COMMENT 'uid',
    `begin_disable_time` datetime      DEFAULT NULL,
    `end_disable_time`   datetime      DEFAULT NULL,
    `ex`                 varchar(1024) DEFAULT NULL,
    PRIMARY KEY (`uid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of black_list
-- ----------------------------

-- ----------------------------
-- Table structure for `chat_log`
-- ----------------------------
DROP TABLE IF EXISTS `chat_log`;
CREATE TABLE `chat_log`
(
    `msg_id`             varchar(128)  NOT NULL,
    `send_id`            varchar(255)  NOT NULL,
    `session_type`       int(11) NOT NULL,
    `recv_id`            varchar(255)  NOT NULL,
    `content_type`       int(11) NOT NULL,
    `msg_from`           int(11) NOT NULL,
    `content`            varchar(1000) NOT NULL,
    `remark`             varchar(100) DEFAULT NULL,
    `sender_platform_id` int(11) NOT NULL,
    `send_time`          datetime      NOT NULL,
    PRIMARY KEY (`msg_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of chat_log
-- ----------------------------

-- ----------------------------
-- Table structure for `friend`
-- ----------------------------
DROP TABLE IF EXISTS `friend`;
CREATE TABLE `friend`
(
    `owner_id`    varchar(255) NOT NULL,
    `friend_id`   varchar(255) NOT NULL,
    `comment`     varchar(255) DEFAULT NULL,
    `friend_flag` int(11) NOT NULL,
    `create_time` datetime     NOT NULL,
    PRIMARY KEY (`owner_id`, `friend_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of friend
-- ----------------------------

-- ----------------------------
-- Table structure for `friend_request`
-- ----------------------------
DROP TABLE IF EXISTS `friend_request`;
CREATE TABLE `friend_request`
(
    `req_id`      varchar(255) NOT NULL,
    `user_id`     varchar(255) NOT NULL,
    `flag`        int(11) NOT NULL DEFAULT '0',
    `req_message` varchar(255) DEFAULT NULL,
    `create_time` datetime     NOT NULL,
    PRIMARY KEY (`user_id`, `req_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of friend_request
-- ----------------------------

-- ----------------------------
-- Table structure for `group`
-- ----------------------------
DROP TABLE IF EXISTS `group`;
CREATE TABLE `group`
(
    `group_id`     varchar(255) NOT NULL,
    `name`         varchar(255) DEFAULT NULL,
    `introduction` varchar(255) DEFAULT NULL,
    `notification` varchar(255) DEFAULT NULL,
    `face_url`     varchar(255) DEFAULT NULL,
    `create_time`  datetime     DEFAULT NULL,
    `ex`           varchar(255) DEFAULT NULL,
    PRIMARY KEY (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of group
-- ----------------------------

-- ----------------------------
-- Table structure for `group_member`
-- ----------------------------
DROP TABLE IF EXISTS `group_member`;
CREATE TABLE `group_member`
(
    `group_id`            varchar(255) NOT NULL,
    `uid`                 varchar(255) NOT NULL,
    `nickname`            varchar(255) DEFAULT NULL,
    `user_group_face_url` varchar(255) DEFAULT NULL,
    `administrator_level` int(11) NOT NULL,
    `join_time`           datetime     NOT NULL,
    PRIMARY KEY (`group_id`, `uid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of group_member
-- ----------------------------

-- ----------------------------
-- Table structure for `group_request`
-- ----------------------------
DROP TABLE IF EXISTS `group_request`;
CREATE TABLE `group_request`
(
    `id`                 int(11) NOT NULL AUTO_INCREMENT,
    `group_id`           varchar(255) NOT NULL,
    `from_user_id`       varchar(255) NOT NULL,
    `to_user_id`         varchar(255) NOT NULL,
    `flag`               int(10) NOT NULL DEFAULT '0',
    `req_msg`            varchar(255) DEFAULT '',
    `handled_msg`        varchar(255) DEFAULT '',
    `create_time`        datetime     NOT NULL,
    `from_user_nickname` varchar(255) DEFAULT '',
    `to_user_nickname`   varchar(255) DEFAULT NULL,
    `from_user_face_url` varchar(255) DEFAULT '',
    `to_user_face_url`   varchar(255) DEFAULT '',
    `handled_user`       varchar(255) DEFAULT '',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of group_request
-- ----------------------------

-- ----------------------------
-- Table structure for `receive`
-- ----------------------------
DROP TABLE IF EXISTS `receive`;
CREATE TABLE `receive`
(
    `user_id`     varchar(255) NOT NULL,
    `seq`         int(11) NOT NULL,
    `msg_id`      varchar(128) NOT NULL,
    `create_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_id`, `seq`) USING BTREE,
    KEY           `fk_msgid` (`msg_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of receive
-- ----------------------------

-- ----------------------------
-- Table structure for `user`
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `uid`         varchar(64) NOT NULL,
    `name`        varchar(64)   DEFAULT NULL,
    `icon`        varchar(1024) DEFAULT NULL,
    `gender`      int(11) unsigned zerofill DEFAULT NULL,
    `mobile`      varchar(32)   DEFAULT NULL,
    `birth`       varchar(16)   DEFAULT NULL,
    `email`       varchar(64)   DEFAULT NULL,
    `ex`          varchar(1024) DEFAULT NULL,
    `create_time` datetime      DEFAULT NULL,
    PRIMARY KEY (`uid`),
    UNIQUE KEY `uk_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of user
-- ----------------------------

-- ----------------------------
-- Table structure for `user_black_list`
-- ----------------------------
DROP TABLE IF EXISTS `user_black_list`;
CREATE TABLE `user_black_list`
(
    `owner_id`    varchar(255) NOT NULL,
    `block_id`    varchar(255) NOT NULL,
    `create_time` datetime     NOT NULL,
    PRIMARY KEY (`owner_id`, `block_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of user_black_list
-- ----------------------------
