/*
 Navicat MySQL Data Transfer

 Source Server         : localhost_3306
 Source Server Type    : MySQL
 Source Server Version : 50726
 Source Host           : localhost:3306
 Source Schema         : test

 Target Server Type    : MySQL
 Target Server Version : 50726
 File Encoding         : 65001

 Date: 15/06/2020 19:26:57
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for my_class
-- ----------------------------
DROP TABLE IF EXISTS `my_class`;
CREATE TABLE `my_class`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `class` tinyint(1) NOT NULL,
  `grade` tinyint(1) NOT NULL,
  `teacher` varchar(255) CHARACTER SET utf8 COLLATE utf8_unicode_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 3 CHARACTER SET = utf8 COLLATE = utf8_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of my_class
-- ----------------------------
INSERT INTO `my_class` VALUES (1, 1, 3, 1, '女老师');
INSERT INTO `my_class` VALUES (2, 2, 3, 1, NULL);

SET FOREIGN_KEY_CHECKS = 1;
