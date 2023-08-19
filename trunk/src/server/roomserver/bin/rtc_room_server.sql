DROP DATABASE IF EXISTS rtc_room_server;
CREATE DATABASE rtc_room_server;
USE rtc_room_server;
-- ----------------------------
-- Table structure for `rtc_talk_sessions`
-- ----------------------------
DROP TABLE IF EXISTS `rtc_talk_sessions`;
CREATE TABLE `rtc_talk_sessions` (
    `id` bigint(32) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `session_id` varchar(32)  DEFAULT '' COMMENT '会话ID',
    `sub_session_ids` varchar(1024)  DEFAULT '' COMMENT '子会话ID',
    `app_id` varchar(32) DEFAULT '' COMMENT '应用ID',
    `user_id` varchar(32) DEFAULT '' COMMENT '用户ID',
    `user_name` varchar(32) DEFAULT '' COMMENT '用户名',
    `talk_type` tinyint(1) DEFAULT 0 COMMENT '通话类型',
    `user_ip` varchar(20) DEFAULT '' COMMENT '用户 IP',
    `time_out` tinyint(1) DEFAULT 0 COMMENT '会话状态',
    `create_time` datetime DEFAULT NULL COMMENT '创建时间，统一为服务器的系统时间',
    `join_time` datetime DEFAULT NULL COMMENT '加入时间',
    `leave_time` datetime DEFAULT NULL COMMENT '退出时间',
    `duration` bigint DEFAULT 0 COMMENT '时长',
    `os_name` varchar(32) DEFAULT '' COMMENT '操作系统名称',
    `browser` varchar(32) DEFAULT '' COMMENT '浏览器',
    `sdk_info` varchar(32) DEFAULT '' COMMENT 'SDK版本',
    PRIMARY KEY(id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='通话会话表';

-- ----------------------------
-- Table structure for `rtc_talk_subsessions`
-- ----------------------------
DROP TABLE IF EXISTS `rtc_talk_subsessions`;
CREATE TABLE `rtc_talk_subsessions` (
    `id` bigint(32) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `sub_session_id` varchar(32)  DEFAULT '' COMMENT '子会话ID',
    `status` tinyint(1) DEFAULT -1 COMMENT '会话状态',
    `remote_user_id` varchar(32) DEFAULT '' COMMENT '对端用户ID',
    `remote_user_name` varchar(32) DEFAULT '' COMMENT '对端用户名',
    `remote_user_ip` varchar(32) DEFAULT '' COMMENT '对端用户 IP',
    `ice_ip` varchar(32) DEFAULT '' COMMENT 'ice IP',
    `connect_type` varchar(10) DEFAULT '' COMMENT 'ice IP',
    `begin_time` datetime DEFAULT NULL COMMENT '起始时间',
    `connect_time` datetime DEFAULT NULL COMMENT '连接时间',
    `finish_time` datetime DEFAULT NULL COMMENT '结束时间',
    `duration` bigint DEFAULT 0 COMMENT '通话时长',
    `cost`  int DEFAULT 0 COMMENT '费用-单位分',
     PRIMARY KEY(id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='通话子会话表';

-- ----------------------------
-- Table structure for `rtc_talk_operations`
-- ----------------------------
DROP TABLE IF EXISTS `rtc_talk_operations`;
CREATE TABLE `rtc_talk_operations` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `app_id` varchar(32) DEFAULT '' COMMENT '应用ID',
    `session_id` varchar(32) DEFAULT '' COMMENT '所属会话ID',
    `sub_session_ids` varchar(100) DEFAULT '' COMMENT '所属子会话ID',
    `user_id` varchar(32) DEFAULT '' COMMENT '用户ID',
    `to_user_id` varchar(32) DEFAULT '' COMMENT '对端用户ID',
    `talk_type` tinyint(1) DEFAULT NULL COMMENT '通话类型',
    `operate_type` tinyint(2) DEFAULT NULL COMMENT '操作类型',
    `operate_time` datetime DEFAULT NULL COMMENT '操作时间',
    PRIMARY KEY (`id`) 
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='通话操作表';

-- ----------------------------
-- Table structure for `rtc_talk_quality_stats_infos`
-- ----------------------------
DROP TABLE IF EXISTS `rtc_talk_quality_stats_infos`;
CREATE TABLE `rtc_talk_quality_stats_infos` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `remote_user_id` varchar(32) DEFAULT '' COMMENT '对端用户名',
    `sub_session_id` varchar(32) DEFAULT '' COMMENT '所属会话ID',
    `create_time` datetime DEFAULT NULL COMMENT '时间',    
    `aud_send_lost_rate` float(3,2) DEFAULT '0' COMMENT '音频发送丢包率',
    `aud_send_bitrate` int(4) DEFAULT '0' COMMENT '音频发送比特率',
    `aud_recv_lost_rate` float(3,2) DEFAULT '0' COMMENT '音频接收丢包率',
    `aud_recv_bitrate` int(4) DEFAULT '0' COMMENT '音频接收比特率',    
    `vid_send_lost_rate` float(3,2) DEFAULT '0' COMMENT '视频发送丢包率',
    `vid_send_bitrate` int(4) DEFAULT '0' COMMENT '视频发送比特率',
    `send_width` int(4) DEFAULT '0' COMMENT '视频发送宽',
    `send_height` int(4) DEFAULT '0' COMMENT '视频发送高',
    `send_framerate_sent` int(4) DEFAULT '0' COMMENT '视频实际发送帧率',
    `vid_recv_lost_rate` float(3,2) DEFAULT '0' COMMENT '视频接收丢包率',
    `vid_recv_bitrate` int(4) DEFAULT '0' COMMENT '视频接收比特率',
    `recv_framerate_recv` int(4) DEFAULT '0' COMMENT '实际视频接收帧率',    
    PRIMARY KEY (`id`) 
)ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT COMMENT='通话质量统计信息';
