create DATABASE if not exists openIM_v3;

CREATE TABLE if not EXISTS `blacks` (
    `owner_user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `block_user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `create_time` datetime(3) DEFAULT NULL,
    `add_source` int(11) DEFAULT NULL,
    `operator_user_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `ex` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    PRIMARY KEY (`owner_user_id`,`block_user_id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE if not EXISTS `chat_logs` (
    `server_msg_id` char(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `client_msg_id` char(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `send_id` char(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `recv_id` char(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `sender_platform_id` int(11) DEFAULT NULL,
    `sender_nick_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `sender_face_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `session_type` int(11) DEFAULT NULL,
    `msg_from` int(11) DEFAULT NULL,
    `content_type` int(11) DEFAULT NULL,
    `content` varchar(3000) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `status` int(11) DEFAULT NULL,
    `send_time` datetime(3) DEFAULT NULL,
    `create_time` datetime(3) DEFAULT NULL,
    `ex` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    KEY `sendTime` (`send_time`),
    KEY `send_id` (`send_time`,`send_id`),
    KEY `recv_id` (`send_time`,`recv_id`),
    KEY `session_type` (`send_time`,`session_type`),
    KEY `session_type_alone` (`session_type`),
    KEY `content_type` (`send_time`,`content_type`),
    KEY `content_type_alone` (`content_type`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE  if not EXISTS  `conversations` (
    `owner_user_id` char(128) COLLATE utf8mb4_unicode_ci NOT NULL,
    `conversation_id` char(128) COLLATE utf8mb4_unicode_ci NOT NULL,
    `conversation_type` int(11) DEFAULT NULL,
    `user_id` char(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `group_id` char(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `recv_msg_opt` int(11) DEFAULT NULL,
    `is_pinned` tinyint(1) DEFAULT NULL,
    `is_private_chat` tinyint(1) DEFAULT NULL,
    `burn_duration` int(11) DEFAULT '30',
    `group_at_type` int(11) DEFAULT NULL,
    `attached_info` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `ex` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `max_seq` bigint(20) DEFAULT NULL,
    `min_seq` bigint(20) DEFAULT NULL,
    `create_time` datetime(3) DEFAULT NULL,
    `is_msg_destruct` tinyint(1) DEFAULT '0',
    `msg_destruct_time` bigint(20) DEFAULT '604800',
    `latest_msg_destruct_time` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`owner_user_id`,`conversation_id`),
    KEY `create_time` (`create_time`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE  if not EXISTS  `friend_requests` (
    `from_user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `to_user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `handle_result` int(11) DEFAULT NULL,
    `req_msg` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `create_time` datetime(3) DEFAULT NULL,
    `handler_user_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `handle_msg` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `handle_time` datetime(3) DEFAULT NULL,
    `ex` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    PRIMARY KEY (`from_user_id`,`to_user_id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE  if not EXISTS  `friends` (
    `owner_user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `friend_user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `remark` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `create_time` datetime(3) DEFAULT NULL,
    `add_source` int(11) DEFAULT NULL,
    `operator_user_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `ex` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    PRIMARY KEY (`owner_user_id`,`friend_user_id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE  if not EXISTS  `group_members` (
    `group_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `nickname` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `user_group_face_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `role_level` int(11) DEFAULT NULL,
    `join_time` datetime(3) DEFAULT NULL,
    `join_source` int(11) DEFAULT NULL,
    `inviter_user_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `operator_user_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `mute_end_time` datetime(3) DEFAULT NULL,
    `ex` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    PRIMARY KEY (`group_id`,`user_id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE  if not EXISTS  `group_requests` (
    `user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `group_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `handle_result` int(11) DEFAULT NULL,
    `req_msg` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `handle_msg` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `req_time` datetime(3) DEFAULT NULL,
    `handle_user_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `handle_time` datetime(3) DEFAULT NULL,
    `join_source` int(11) DEFAULT NULL,
    `inviter_user_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `ex` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    PRIMARY KEY (`user_id`,`group_id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE  if not EXISTS  `groups` (
    `group_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `notification` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `introduction` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `face_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `create_time` datetime(3) DEFAULT NULL,
    `ex` longtext COLLATE utf8mb4_unicode_ci,
    `status` int(11) DEFAULT NULL,
    `creator_user_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `group_type` int(11) DEFAULT NULL,
    `need_verification` int(11) DEFAULT NULL,
    `look_member_info` int(11) DEFAULT NULL,
    `apply_member_friend` int(11) DEFAULT NULL,
    `notification_update_time` datetime(3) DEFAULT NULL,
    `notification_user_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    PRIMARY KEY (`group_id`),
    KEY `create_time` (`create_time`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE  if not EXISTS  `object_hash` (
    `hash` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
    `engine` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
    `size` bigint(20) DEFAULT NULL,
    `bucket` longtext COLLATE utf8mb4_unicode_ci,
    `name` longtext COLLATE utf8mb4_unicode_ci,
    `create_time` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`hash`,`engine`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE  if not EXISTS  `object_info` (
    `name` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `hash` longtext COLLATE utf8mb4_unicode_ci,
    `content_type` longtext COLLATE utf8mb4_unicode_ci,
    `valid_time` datetime(3) DEFAULT NULL,
    `create_time` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`name`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE  if not EXISTS  `object_put` (
    `put_id` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL,
    `hash` longtext COLLATE utf8mb4_unicode_ci,
    `path` longtext COLLATE utf8mb4_unicode_ci,
    `name` longtext COLLATE utf8mb4_unicode_ci,
    `content_type` longtext COLLATE utf8mb4_unicode_ci,
    `object_size` bigint(20) DEFAULT NULL,
    `fragment_size` bigint(20) DEFAULT NULL,
    `put_urls_hash` longtext COLLATE utf8mb4_unicode_ci,
    `valid_time` datetime(3) DEFAULT NULL,
    `effective_time` datetime(3) DEFAULT NULL,
    `create_time` datetime(3) DEFAULT NULL,
    PRIMARY KEY (`put_id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE  if not EXISTS  `users` (
    `user_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
    `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `face_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `ex` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `create_time` datetime(3) DEFAULT NULL,
    `app_manger_level` int(11) DEFAULT '18',
    `global_recv_msg_opt` int(11) DEFAULT NULL,
    PRIMARY KEY (`user_id`),
    KEY `create_time` (`create_time`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;