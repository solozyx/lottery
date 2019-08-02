/*
SQLyog Ultimate v11.24 (32 bit)
MySQL - 5.5.46-0ubuntu0.14.04.2 : Database - testbbs
*********************************************************************
*/


/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`lottery` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `lottery`;

/*Table structure for table `lt_blackip` */

DROP TABLE IF EXISTS `lt_blackip`;

CREATE TABLE `lt_blackip` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `ip` varchar(50) NOT NULL DEFAULT '' COMMENT 'IP地址',
  `blacktime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '黑名单限制到期时间',
  `sys_created` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `sys_updated` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip` (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

/*Data for the table `lt_blackip` */

insert  into `lt_blackip`(`id`,`ip`,`blacktime`,`sys_created`,`sys_updated`) values (1,'127.0.0.1',0,0,1532606350);

/*Table structure for table `lt_code` */

DROP TABLE IF EXISTS `lt_code`;

CREATE TABLE `lt_code` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `gift_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '奖品ID，关联lt_gift表',
  `code` varchar(255) NOT NULL DEFAULT '' COMMENT '虚拟券编码',
  `sys_created` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `sys_updated` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `sys_status` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '状态，0正常，1作废，2已发放',
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`),
  KEY `gift_id` (`gift_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;

/*Data for the table `lt_code` */

insert  into `lt_code`(`id`,`gift_id`,`code`,`sys_created`,`sys_updated`,`sys_status`) values (1,4,'abc\r',1532602694,0,0),(2,4,'aa\r',1532602694,0,0),(3,4,'cs',1532602694,0,0),(4,4,'332',1532602970,0,2);

/*Table structure for table `lt_gift` */

DROP TABLE IF EXISTS `lt_gift`;

CREATE TABLE `lt_gift` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '奖品名称',
  `prize_num` int(11) NOT NULL DEFAULT '-1' COMMENT '奖品数量，0 无限量，>0限量，<0无奖品',
  `left_num` int(11) NOT NULL DEFAULT '0' COMMENT '剩余数量',
  `prize_code` varchar(50) NOT NULL DEFAULT '' COMMENT '0-9999表示100%，0-0表示万分之一的中奖概率',
  `prize_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '发奖周期，D天',
  `img` varchar(255) NOT NULL DEFAULT '' COMMENT '奖品图片',
  `displayorder` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '位置序号，小的排在前面',
  `gtype` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '奖品类型，0 虚拟币，1 虚拟券，2 实物-小奖，3 实物-大奖',
  `gdata` varchar(255) NOT NULL DEFAULT '' COMMENT '扩展数据，如：虚拟币数量',
  `time_begin` int(11) NOT NULL DEFAULT '0' COMMENT '开始时间',
  `time_end` int(11) NOT NULL DEFAULT '0' COMMENT '结束时间',
  `prize_data` mediumtext COMMENT '发奖计划，[[时间1,数量1],[时间2,数量2]]',
  `prize_begin` int(11) NOT NULL DEFAULT '0' COMMENT '发奖计划周期的开始',
  `prize_end` int(11) NOT NULL DEFAULT '0' COMMENT '发奖计划周期的结束',
  `sys_status` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '状态，0 正常，1 删除',
  `sys_created` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `sys_updated` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `sys_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '操作人IP',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

/*Data for the table `lt_gift` */

insert  into `lt_gift`(`id`,`title`,`prize_num`,`prize_code`,`prize_time`,`img`,`displayorder`,`gtype`,`gdata`,`time_begin`,`time_end`,`prize_data`,`sys_status`,`sys_created`,`sys_updated`,`sys_ip`) values (1,'T恤',10,'1-100',30,'https://p0.ssl.qhmsg.com/t016c44d161c478cfe0.png',1,2,'',1532592420,1564128420,'',0,1532592429,1532593773,'::1'),(2,'360手机N7',1,'0-0',30,'https://p0.ssl.qhmsg.com/t016ff98b934914aca6.png',0,3,'',1532592420,1564128420,'',0,1532592474,0,''),(3,'手机充电器',10,'200-1000',30,'https://p0.ssl.qhmsg.com/t01ec4648d396ad46bf.png',3,2,'',1532592420,1564128420,'',0,1532592558,1532593828,'::1'),(4,'优惠券',100,'2000-5000',1,'https://p0.ssl.qhmsg.com/t01f84f00d294279957.png',4,1,'',1532592420,1564128420,'',0,1532599140,0,'::1');

/*Table structure for table `lt_result` */

DROP TABLE IF EXISTS `lt_result`;

CREATE TABLE `lt_result` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `gift_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '奖品ID，关联lt_gift表',
  `gift_name` varchar(255) NOT NULL DEFAULT '' COMMENT '奖品名称',
  `gift_type` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '奖品类型，同lt_gift. gtype',
  `uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `username` varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
  `prize_code` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '抽奖编号（4位的随机数）',
  `gift_data` varchar(255) NOT NULL DEFAULT '' COMMENT '获奖信息',
  `sys_created` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `sys_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '用户抽奖的IP',
  `sys_status` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '状态，0 正常，1删除，2作弊',
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  KEY `gift_id` (`gift_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

/*Data for the table `lt_result` */

insert  into `lt_result`(`id`,`gift_id`,`gift_name`,`gift_type`,`uid`,`username`,`prize_code`,`gift_data`,`sys_created`,`sys_ip`,`sys_status`) values (1,1,'T恤',2,1,'yifan',1,'',0,'',0);

/*Table structure for table `lt_user` */

DROP TABLE IF EXISTS `lt_user`;

CREATE TABLE `lt_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
  `blacktime` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '黑名单限制到期时间',
  `realname` varchar(50) NOT NULL DEFAULT '' COMMENT '联系人',
  `mobile` varchar(50) NOT NULL DEFAULT '' COMMENT '手机号',
  `address` varchar(255) NOT NULL DEFAULT '' COMMENT '联系地址',
  `sys_created` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `sys_updated` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `sys_ip` varchar(50) NOT NULL DEFAULT '' COMMENT 'IP地址',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

/*Data for the table `lt_user` */

insert  into `lt_user`(`id`,`username`,`blacktime`,`realname`,`mobile`,`address`,`sys_created`,`sys_updated`,`sys_ip`) values (1,'wangyi',0,'一凡Sir','11111111111','abcdefg',0,1532595094,'');

/*Table structure for table `lt_userday` */

DROP TABLE IF EXISTS `lt_userday`;

CREATE TABLE `lt_userday` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `day` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '日期，如：20180725',
  `num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '次数',
  `sys_created` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `sys_updated` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uid_day` (`uid`,`day`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Data for the table `lt_userday` */

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
