DROP TABLE IF EXISTS `courses`;

CREATE TABLE IF NOT EXISTS `courses` (
    `id` CHAR(36) NOT NULL,
    `user_id` CHAR(36) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `content` TEXT,
    `created_at` DATETIME NOT NULL,
    `created_by` CHAR(36) NOT NULL,
    `updated_at` DATETIME,
    `updated_by` CHAR(36),
    `deleted_at` DATETIME,
    `deleted_by` CHAR(36),
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_course_user_id` FOREIGN KEY (`user_id`)
        REFERENCES `users` (`id`)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;