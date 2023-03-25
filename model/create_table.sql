DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
                        `id` bigint(20) NOT NULL AUTO_INCREMENT,
                        `user_id` bigint(20) NOT NULL,
                        `username` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
                        `password` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
                        `email` varchar(64) COLLATE utf8mb4_general_ci NOT NULL ,
                        `location` varchar(64) COLLATE utf8mb4_general_ci DEFAULT '',
                        `position` varchar(64) COLLATE utf8mb4_general_ci DEFAULT '',
                        `company` varchar(64) COLLATE utf8mb4_general_ci DEFAULT '',
                        `homepage` varchar(64) COLLATE utf8mb4_general_ci DEFAULT '',
                        `github` varchar(64) COLLATE utf8mb4_general_ci DEFAULT '',
                        `avatar` varchar(128) COLLATE utf8mb4_general_ci DEFAULT '',
                        `introduction` varchar(256) COLLATE utf8mb4_general_ci DEFAULT '',
                        `follower` int NOT NULL DEFAULT '0',
                        `following` int NOT NULL DEFAULT '0',
                        `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                        `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `idx_email` (`email`) USING BTREE,
                        UNIQUE KEY `idx_user_id` (`user_id`) USING BTREE
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

alter table user add column `follower` int NOT NULL DEFAULT '0' after introduction;
alter table user add column `following` int NOT NULL DEFAULT '0' after follower;
alter table user add column `github` varchar(64) COLLATE utf8mb4_general_ci DEFAULT '' after homepage;

DROP TABLE IF EXISTS `admin`;
CREATE TABLE `admin` (
                        `id` bigint(20) NOT NULL AUTO_INCREMENT,
                        `admin_id` bigint(20) NOT NULL,
                        `admin_name` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
                        `password` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
                        `email` varchar(64) COLLATE utf8mb4_general_ci NOT NULL ,
                        `avatar` varchar(128) COLLATE utf8mb4_general_ci DEFAULT '',
                        `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                        `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `idx_email` (`email`) USING BTREE,
                        UNIQUE KEY `idx_user_id` (`admin_id`) USING BTREE
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

# DROP TABLE IF EXISTS `tag`;
CREATE TABLE `tag`(
                            `id` int(11) NOT NULL AUTO_INCREMENT,
                            `tag_id` int(10) unsigned NOT NULL,
                            `article_number` bigint(20) NOT NULL DEFAULT 0,
                            `follower_number` bigint(20) NOT NULL DEFAULT 0,
                            `tag_name` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,
                            `image` varchar(128) COLLATE utf8mb4_general_ci DEFAULT '',
                            `introduction` varchar(256) COLLATE utf8mb4_general_ci DEFAULT '',
                            `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                            `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                            PRIMARY KEY (`id`),
                            UNIQUE KEY `idx_tag_id` (`tag_id`),
                            UNIQUE KEY `idx_tag_name` (`tag_name`)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `tag` VALUES ('1','1','0','0','JavaScript','https://s1.ax1x.com/2022/10/15/x0NQHO.png',
                          'JavaScript是一种跨平台、面向对象的脚本语言。',
                          '2022-10-07 16:41:00','2022-10-07 16:41:00');
INSERT INTO `tag` VALUES ('2','2','0','0','Web开发','https://s1.ax1x.com/2022/10/15/x0NDUg.png',
                          'Web 开发是指构建、创建和维护网站。它包括网页设计、网页发布、网页编程和数据库管理等方面。',
                          '2022-10-07 16:46:00','2022-10-07 16:46:00');
INSERT INTO `tag` VALUES ('3','3','0','0','Go语言','https://s1.ax1x.com/2022/10/15/x0ttmT.png',
                          'Go是一种静态强类型、编译型、并发型，并具有垃圾回收功能的编程语言，可以轻松构建简单、可靠和高效的软件。',
                          '2022-10-07 16:51:00','2022-10-07 16:51:00');
INSERT INTO `tag` VALUES ('4','4','0','0','Vue.js','https://s1.ax1x.com/2022/10/15/x0tg0O.png',
                          'Vue 是一个用于构建用户界面的渐进式框架。',
                          '2022-10-07 16:52:00','2022-10-07 16:52:00');
INSERT INTO `tag` VALUES ('5','5','0','0','Python','https://s1.ax1x.com/2022/10/15/x0tI1I.png',
                          'Python 是一种具有动态语义的解释型、面向对象的高级编程语言。',
                          '2022-10-07 17:50:00','2022-10-07 17:50:00');
INSERT INTO `tag` VALUES ('6','6','0','0','Java','https://s1.ax1x.com/2022/10/15/x0tdk4.png',
                          'Java是一种广泛使用的计算机编程语言，拥有跨平台、面向对象、泛型编程的特性，广泛应用于企业级Web应用开发和移动应用开发。',
                          '2022-10-07 17:53:00','2022-10-07 17:53:00');
INSERT INTO `tag` VALUES ('7','7','0','0','TypeScript','https://s1.ax1x.com/2022/10/15/x0NMDK.png',
                          'TypeScript 是一种开源语言，它通过添加静态类型定义构建在 JavaScript之上。',
                          '2022-10-07 17:55:00','2022-10-07 17:55:00');
INSERT INTO `tag` VALUES ('8','8','0','0','Tailwind CSS','https://s1.ax1x.com/2022/10/15/x0t27D.png',
                          '一个实用程序优先的 CSS 框架，可以直接在您的标记中构建任何设计。',
                          '2022-10-07 17:57:00','2022-10-07 17:57:00');
INSERT INTO `tag` VALUES ('9','9','0','0','C++','https://s1.ax1x.com/2022/10/15/x0tlfs.png',
                          'C++ 是一种通用编程语言，它具有命令式、面向对象和通用编程特性，支持低级内存操作。',
                          '2022-10-12 21:50:00','2022-10-12 21:50:00');
INSERT INTO `tag` VALUES ('10','10','0','0','MySQL','https://s1.ax1x.com/2022/10/15/x0tXNQ.png',
                          'MySQL 是一个开源关系 SQL 数据库管理系统。',
                          '2022-10-12 22:29:00','2022-10-12 22:29:00');
INSERT INTO `tag` VALUES ('11','11','0','0','Redis','https://s1.ax1x.com/2022/10/15/x0tftH.png',
                          'Redis 是一种内存数据存储，通常用作数据库、缓存层和消息队列。它是一个键值存储，支持字符串、哈希、列表、集合、排序集合等数据结构。',
                          '2022-10-12 22:30:00','2022-10-12 22:30:00');
INSERT INTO `tag` VALUES ('12','12','0','0','CSS','https://s1.ax1x.com/2022/10/15/x0tG60.png',
                          '层叠样式表(CSS)是一种样式表语言，用于描述以HTML编写的文档的呈现方式。',
                          '2022-10-12 22:31:00','2022-10-12 22:31:00');
INSERT INTO `tag` VALUES ('13','13','0','0','HTML','https://s1.ax1x.com/2022/10/15/x0tN0U.png',
                          '超文本标记语言(HTML)是用于创建网页和 Web 应用程序的标准标记语言。',
                          '2022-10-12 22:31:00','2022-10-12 22:31:00');
INSERT INTO `tag` VALUES ('14','14','0','0','React.js','https://s1.ax1x.com/2022/10/15/x0thhd.png',
                          'React是一个自由及开放源代码的前端JavaScript工具库，基于UI组件构建用户界面。',
                          '2022-10-12 22:32:00','2022-10-12 22:32:00');
INSERT INTO `tag` VALUES ('15','15','0','0','Linux','https://s1.ax1x.com/2022/10/15/x08TdH.png',
                          'Linux是一种自由和开放源码的类UNIX操作系统。',
                          '2022-10-12 22:33:00','2022-10-12 22:33:00');
INSERT INTO `tag` VALUES ('16','16','0','0','面试','https://s1.ax1x.com/2022/10/15/x08oee.png',
                          '',
                          '2022-10-12 22:36:00','2022-10-12 22:36:00');
INSERT INTO `tag` VALUES ('17','17','0','0','算法','https://s1.ax1x.com/2022/10/30/xIAGnI.png',
                          '',
                          '2022-10-12 22:37:00','2022-10-12 22:37:00');
INSERT INTO `tag` VALUES ('18','18','0','0','Docker','https://s1.ax1x.com/2022/10/30/xIApkT.png',
                          'Docker是一个开源的引擎，可以轻松的为任何应用创建一个轻量级的、可移植的、自给自足的容器。',
                          '2022-10-30 01:35:00','2022-10-30 01:35:00');
INSERT INTO `tag` VALUES ('19','19','0','0','Nginx','https://s1.ax1x.com/2022/10/30/xIAF1J.png',
                          'Nginx是一个高性能的HTTP和反向代理web服务器。',
                          '2022-10-30 01:54:00','2022-10-30 01:54:00');
INSERT INTO `tag` VALUES ('20','20','0','0','MongoDB','https://s1.ax1x.com/2022/10/30/xIAip4.png',
                          'MongoDB 是一个基于分布式文件存储的NoSQL数据库。',
                          '2022-10-30 01:54:00','2022-10-30 01:54:00');
INSERT INTO `tag` VALUES ('21','21','0','0','Git','https://s1.ax1x.com/2022/10/30/xIAChF.png',
                          'Git 是一个开源的分布式版本控制系统，用于敏捷高效地处理任何或小或大的项目。',
                          '2022-10-30 01:54:00','2022-10-30 01:54:00');
INSERT INTO `tag` VALUES ('22','22','0','0','PostgreSQL','https://s1.ax1x.com/2022/10/30/xIAkc9.png',
                          'PostgreSQL是一个功能非常强大的、开源的客户/服务器关系型数据库管理系统。',
                          '2022-10-30 01:54:00','2022-10-30 01:54:00');
INSERT INTO `tag` VALUES ('23','23','0','0','Node.js','https://s1.ax1x.com/2022/10/30/xIAl1H.png',
                          'Node.js 是一个基于 Chrome V8 引擎的 Javascript 运行环境。',
                          '2022-10-30 01:54:00','2022-10-30 01:54:00');
INSERT INTO `tag` VALUES ('24','24','0','0','Django','https://s1.ax1x.com/2022/10/30/xIA1cd.png',
                          'Django 是一个由 Python 编写的一个开放源代码的 Web 应用框架。',
                          '2022-10-30 01:54:00','2022-10-30 01:54:00');

-- DROP TABLE IF EXISTS `article`;
CREATE TABLE `article`(
                       `id` bigint(20) NOT NULL AUTO_INCREMENT,
                       `article_id` bigint(20) NOT NULL COMMENT '文章ID',
                       `title` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '标题',
                       `content` longtext COLLATE utf8mb4_general_ci NOT NULL COMMENT '原始内容',
                       `html` longtext COLLATE utf8mb4_general_ci NOT NULL COMMENT 'html',
                       `markdown` longtext COLLATE utf8mb4_general_ci NOT NULL COMMENT 'markdown',
                       `image` varchar(128) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '头图',
                       `source` varchar(128) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '资源',
                       `subtitle` varchar(128) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '副标题',
                       `author_id` bigint(20) NOT NULL COMMENT '作者的用户ID',
                       `view_count` int DEFAULT '0' COMMENT '浏览量',
                       `likes` int DEFAULT '0' COMMENT '喜欢',
                       `comments` int DEFAULT '0' COMMENT '评论',
                       `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                       `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                       PRIMARY KEY (`id`),
                       UNIQUE KEY `idx_article_id` (`article_id`) USING BTREE,
                       KEY `idx_author_id` (`author_id`) USING BTREE
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

alter table article add column `subtitle` varchar(128) collate utf8mb4_general_ci default '' after title;

-- DROP TABLE IF EXISTS `article_tag`;
CREATE TABLE `article_tag`(
                          `id` bigint(20) NOT NULL AUTO_INCREMENT,
                          `article_id` bigint(20) NOT NULL COMMENT '文章ID',
                          `tag_id` int(10) unsigned NOT NULL COMMENT '标签ID',
                          `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                          PRIMARY KEY (`id`),
                          KEY `idx_article_id` (`article_id`) USING BTREE
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- DROP TABLE IF EXISTS `follow_tag`;
CREATE TABLE `follow_tag`(
                              `id` bigint(20) NOT NULL AUTO_INCREMENT,
                              `user_id` bigint(20) NOT NULL COMMENT '用户ID',
                              `follow_tag_id` int(10) unsigned NOT NULL COMMENT '标签ID',
                              `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                              PRIMARY KEY (`id`),
                              KEY `idx_article_id` (`user_id`) USING BTREE
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- DROP TABLE IF EXISTS `follow_user`;
CREATE TABLE `follow_user`(
                             `id` bigint(20) NOT NULL AUTO_INCREMENT,
                             `user_id` bigint(20) NOT NULL COMMENT '用户ID',
                             `follow_user_id` bigint(20) NOT NULL COMMENT '被关注用户ID',
                             `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                             `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                             PRIMARY KEY (`id`),
                             KEY `idx_article_id` (`user_id`) USING BTREE
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments`(
                          `comment_id` bigint(20) NOT NULL COMMENT '评论ID',
                          `user_id` bigint(20) NOT NULL COMMENT '用户ID',
                          `item_id` bigint(20) NOT NULL COMMENT '评论目标ID',
                          `item_type` int NOT NULL COMMENT '评论目标类型',
                          `status` int NOT NULL COMMENT '状态',
                          `likes` int DEFAULT '0' COMMENT '喜欢',
                          `comment_picture` varchar(128) DEFAULT '' COMMENT '评论图片',
                          `comment_content` text NOT NULL COMMENT '评论内容',
                          `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          PRIMARY KEY (`comment_id`)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- DROP TABLE IF EXISTS `tree_path`;
CREATE TABLE `tree_path`(
                          `ancestor` bigint(20) NOT NULL COMMENT '祖先ID',
                          `descendant` bigint(20) NOT NULL COMMENT '后代ID',
                          `distance` int NOT NULL COMMENT '深度',
                          PRIMARY KEY (`ancestor`,`descendant`,`distance`),
                          KEY `idx_descendant` (`descendant`) USING BTREE,
                          KEY `idx_distance` (`distance`) USING BTREE
#                           FOREIGN KEY (`ancestor`) REFERENCES comments(`comment_id`),
#                           FOREIGN KEY (`descendant`) REFERENCES comments(`comment_id`)
)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

select title from article where title like '%GO%'
