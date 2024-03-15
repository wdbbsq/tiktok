CREATE TABLE `t_user`
(
    `id`             bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `name`           varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '姓名',
    `username`       varchar(32) COLLATE utf8mb4_unicode_ci  DEFAULT NULL COMMENT '用户名',
    `password`       varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '密码',
    `follower_count` bigint(20)                              DEFAULT NULL COMMENT '粉丝数',
    `follow_count`   bigint(255)                             DEFAULT NULL COMMENT '关注数',
    `created_at`     datetime                                DEFAULT NULL COMMENT '创建时间',
    `updated_at`     datetime                                DEFAULT NULL COMMENT '修改时间',
    `deleted_at`     datetime                                DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='用户表';

CREATE TABLE `t_video`
(
    `id`             bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `title`          varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '标题',
    `play_url`       varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '播放地址',
    `cover_url`      varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '封面地址',
    `favorite_count` bigint(20)                              DEFAULT 0 COMMENT '点赞数',
    `comment_count`  bigint(20)                              DEFAULT 0 COMMENT '评论数',
    `is_favorite`    tinyint(1)                              DEFAULT NULL COMMENT '用户是否点赞',
    `user_id`        bigint(20) unsigned NOT NULL COMMENT '用户ID',
    `created_at`     datetime                                DEFAULT NULL COMMENT '创建时间',
    `updated_at`     datetime                                DEFAULT NULL COMMENT '修改时间',
    `deleted_at`     datetime                                DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='视频表';

CREATE TABLE `t_favorite`
(
    `id`         bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id`    bigint(20) unsigned NOT NULL COMMENT '用户ID',
    `video_id`   bigint(20) unsigned NOT NULL COMMENT '视频ID',
    `created_at` datetime DEFAULT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT = '点赞表';
