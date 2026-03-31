CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `nickname` varchar(50) NOT NULL,
  `email` varchar(100) NOT NULL DEFAULT '',
  `phone` varchar(20) NOT NULL DEFAULT '',
  `avatar` varchar(255) NOT NULL DEFAULT '',
  `dept_id` bigint unsigned DEFAULT NULL,
  `role_id` bigint unsigned DEFAULT NULL,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `departments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` bigint unsigned DEFAULT NULL,
  `name` varchar(50) NOT NULL,
  `sort` int NOT NULL DEFAULT 0,
  `leader` varchar(50) NOT NULL DEFAULT '',
  `phone` varchar(20) NOT NULL DEFAULT '',
  `email` varchar(100) NOT NULL DEFAULT '',
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `code` varchar(50) NOT NULL,
  `description` varchar(255) NOT NULL DEFAULT '',
  `permissions` json DEFAULT NULL,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `approval_flows` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `code` varchar(50) NOT NULL,
  `description` varchar(255) NOT NULL DEFAULT '',
  `nodes` json NOT NULL,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `approvals` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `type` varchar(20) NOT NULL,
  `content` json DEFAULT NULL,
  `applicant_id` bigint unsigned NOT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'pending',
  `current_node` int NOT NULL DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_applicant` (`applicant_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `approval_nodes` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `approval_id` bigint unsigned NOT NULL,
  `name` varchar(100) NOT NULL,
  `type` varchar(20) NOT NULL,
  `approver_ids` json DEFAULT NULL,
  `status` varchar(20) NOT NULL DEFAULT 'pending',
  `comment` varchar(500) NOT NULL DEFAULT '',
  `approver_id` bigint unsigned DEFAULT NULL,
  `sort` int NOT NULL DEFAULT 0,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_approval_id` (`approval_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `notices` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `content` text NOT NULL,
  `type` tinyint NOT NULL DEFAULT 1,
  `summary` varchar(500) NOT NULL DEFAULT '',
  `cover` varchar(255) NOT NULL DEFAULT '',
  `is_top` tinyint NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `publisher_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_publisher` (`publisher_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `notice_reads` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `notice_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_notice_user` (`notice_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `schedules` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `description` text,
  `start_time` datetime(3) NOT NULL,
  `end_time` datetime(3) NOT NULL,
  `is_all_day` tinyint NOT NULL DEFAULT 0,
  `priority` tinyint NOT NULL DEFAULT 1,
  `location` varchar(255) NOT NULL DEFAULT '',
  `color` varchar(20) NOT NULL DEFAULT '#1677FF',
  `creator_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_creator` (`creator_id`),
  KEY `idx_time` (`start_time`, `end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `schedule_participants` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `schedule_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_schedule_user` (`schedule_id`, `user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
