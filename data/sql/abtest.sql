CREATE TABLE `policy` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '策略名',
  `update_time` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态：0.通过，1.审核中，2.拒绝',
  `is_delete` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '状态：0.有效，1.删除',
  `ext1` varchar(255) NOT NULL DEFAULT '""',
  `ext2` int(10) unsigned NOT NULL DEFAULT '0',
  `div_model` varchar(255) NOT NULL DEFAULT '' COMMENT '线上策略分流模块名',
  `div_data` text NOT NULL COMMENT 'diversion json string data',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='自定义策略，策略数据以json的形式存入div_data中，div_model为线上已有的分流模块名';

CREATE TABLE `policy_group` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '组名',
  `policy_ids` varchar(255) NOT NULL DEFAULT '' COMMENT '策略id，最多10个，用,分割',
  `update_time` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态:0.通过，1.审核中，2.拒绝',
  `is_delete` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '状态:0.有效，1.删除',
  `ext1` varchar(255) NOT NULL DEFAULT '""',
  `ext2` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='策略组最多10个策略(纵横)';

CREATE TABLE `runtime` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `server_name` varchar(255) NOT NULL DEFAULT '' COMMENT '运行时策略对应的服务名',
  `policy_id` bigint(20) NOT NULL DEFAULT '0',
  `group_id` bigint(20) NOT NULL DEFAULT '0',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态：0.通过，1.审核中，2.拒绝，3.下线，4.上线',
  `is_delete` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '状态：0.可用，1.删除',
  `update_time` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `ext1` varchar(255) NOT NULL DEFAULT '""',
  `ext2` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `server_name` (`server_name`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8 COMMENT='运行时策略，用于线上所选策略';
