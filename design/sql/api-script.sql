CREATE DATABASE IF NOT EXISTS go_project DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE go_project;

CREATE TABLE `banner` (
    id bigint AUTO_INCREMENT PRIMARY KEY,
    city int NOT NULL DEFAULT 0 COMMENT '城市编码',
    title varchar(20) NOT NULL DEFAULT '',
    img varchar(150) NOT NULL DEFAULT '' COMMENT '图片链接',
    type tinyint NOT NULL DEFAULT 0 COMMENT '0不跳转，1小程序内部路径，2外部H5链接',
    link varchar(200) NOT NULL DEFAULT '' COMMENT '跳转链接路径',
    sort tinyint NOT NULL DEFAULT 0 COMMENT '排序(0~99),从小到大',
    begin_time bigint NOT NULL DEFAULT 0 COMMENT '开始时间',
    end_time bigint NOT NULL DEFAULT 0 COMMENT '结束时间',
    status tinyint NOT NULL DEFAULT 1 COMMENT 'off(-1),on(1)',
    create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='轮播图';

CREATE TABLE `user` (
    id bigint AUTO_INCREMENT PRIMARY KEY,
    openid varchar(50) NOT NULL UNIQUE,
    unionid varchar(50) NOT NULL DEFAULT '',
    phone_number varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
    nickname varchar(10) NOT NULL DEFAULT '' COMMENT '昵称',
    avatar_url varchar(150) NOT NULL DEFAULT '' COMMENT '头像链接',
    create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY (phone_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户信息';

CREATE TABLE `wechat_analysis` (
    ref_date varchar(10) PRIMARY KEY,
    session_cnt int NOT NULL DEFAULT 0 COMMENT '打开次数',
    visit_pv int NOT NULL DEFAULT 0 COMMENT '访问次数',
    visit_uv int NOT NULL DEFAULT 0 COMMENT '访问人数',
    visit_uv_new int NOT NULL DEFAULT 0 COMMENT '新用户数',
    share_pv int NOT NULL DEFAULT 0 COMMENT '转发次数',
    share_uv int NOT NULL DEFAULT 0 COMMENT '转发人数',
    create_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='微信小程序访问趋势';
