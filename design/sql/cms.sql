CREATE TABLE `admin_role` (
    id int AUTO_INCREMENT PRIMARY KEY,
    name varchar(32) NOT NULL DEFAULT '',
    authority json,
    create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员角色';

CREATE TABLE `admin_user` (
    id int AUTO_INCREMENT PRIMARY KEY,
    username varchar(32) NOT NULL UNIQUE,
    password varchar(64) NOT NULL DEFAULT '',
    role_id int NOT NULL DEFAULT 0 COMMENT '0-super',
    status tinyint NOT NULL DEFAULT 1 COMMENT 'off(-1),on(1)',
    create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY(role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci AUTO_INCREMENT=1000 COMMENT='管理员账号';

INSERT INTO `admin_user` (id,username,password) VALUES
(1,'admin','jZae727K08KaOmKSgOaGzww_XVqGr_PKEgIMkjrc');
-- 受保护的超管账号admin，初始密码: 123456
