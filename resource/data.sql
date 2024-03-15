INSERT INTO `t_user` (`id`, `name`, `username`, `password`, `follower_count`, `follow_count`, `created_at`,
                      `updated_at`,
                      `deleted_at`)
VALUES (1, '张三', 'zhangsan', 'password123', 100, 50, '2024-02-01 10:00:00', '2024-02-01 10:00:00', NULL),
       (2, '李四', 'lisi', 'password456', 150, 80, '2024-02-02 10:00:00', '2024-02-02 10:00:00', NULL),
       (3, '王五', 'wangwu', 'password789', 200, 120, '2024-02-03 10:00:00', '2024-02-03 10:00:00', NULL),
       (4, '赵十', 'zhaoshi', 'password0987', 300, 180, '2024-02-10 10:00:00', '2024-02-10 10:00:00', NULL);

INSERT INTO `t_video` (`id`, `title`, `play_url`, `cover_url`, `favorite_count`, `comment_count`, `is_favorite`,
                       `user_id`,
                       `created_at`, `updated_at`, `deleted_at`)
VALUES (1, '视频1', 'http://example.com/play1', 'http://example.com/cover1', 0, 0, 0, 0, '2024-02-01 10:00:00',
        '2024-02-01 10:00:00', NULL),
       (2, '视频2', 'http://example.com/play2', 'http://example.com/cover2', 0, 0, 0, 0, '2024-02-02 10:00:00',
        '2024-02-02 10:00:00', NULL),
       (3, '视频3', 'http://example.com/play3', 'http://example.com/cover3', 0, 0, 0, 0, '2024-02-03 10:00:00',
        '2024-02-03 10:00:00', NULL),
       (4, '视频10', 'http://example.com/play10', 'http://example.com/cover10', 0, 0, 0, 0, '2024-02-10 10:00:00',
        '2024-02-10 10:00:00', NULL);

