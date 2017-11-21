CREATE TABLE `t_user` (
  `id`         INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_name`  VARCHAR(40)      NOT NULL DEFAULT ''
  COMMENT '用户名',
  `password`   VARCHAR(60)      NOT NULL DEFAULT ''
  COMMENT '密码',
  `email`      VARCHAR(30)      NOT NULL DEFAULT '',
  `salt`       CHAR(10)         NOT NULL DEFAULT ''
  COMMENT '密码盐',
  `last_login` INT(11)          NOT NULL DEFAULT '0'
  COMMENT '最后登录的时间',
  `status`     TINYINT(4)       NOT NULL DEFAULT '0'
  COMMENT '状态, 0 正常 -1 禁用',
  `created_at` INT(11) UNSIGNED NOT NULL DEFAULT '0'
  COMMENT '创建时间',
  `bduss`      VARCHAR(300)              DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_name` (`user_name`)
)
  ENGINE = INNODB
  DEFAULT CHARSET = utf8;

CREATE TABLE `t_forums` (
  `id`          INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`     INT(11) UNSIGNED NOT NULL DEFAULT '0',
  `kw`          VARCHAR(100)     NOT NULL DEFAULT ''
  COMMENT '贴吧的kw值',
  `fid`         INT(10)          NOT NULL DEFAULT '-1'
  COMMENT '贴吧的fid',
  `last_sign`   TINYINT(2)       NOT NULL DEFAULT '-1'
  COMMENT '上一次签到的日期',
  `sign_status` TINYINT(2)       NOT NULL DEFAULT '-1'
  COMMENT '签到的状态，0 成功，-1为失败',
  `created_at`  INT(11) UNSIGNED NOT NULL DEFAULT '0'
  COMMENT '创建时间',
  `reply_json`  VARCHAR(500)              DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid_kw`(`user_id`, `kw`)
)
  ENGINE = INNODB
  DEFAULT CHARSET = utf8;

