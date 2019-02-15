use homebot;

CREATE TABLE `mod_command` (
	  `mod_sn` char(100) NOT NULL DEFAULT '',
	  `command` text NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

CREATE TABLE `user_info` (
	  `user_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	  `msgr_id` char(22) NOT NULL,
	  `msgr_id2` char(22) DEFAULT NULL,
	  `last_name` char(12) DEFAULT NULL,
	  `first_name` char(12) DEFAULT NULL,
	  `msgr_name` char(10) DEFAULT NULL,
	  `register_date` datetime DEFAULT NULL,
	  `user_phone` char(20) DEFAULT NULL,
	  `user_addr` char(100) DEFAULT NULL,
	  `deprecated` tinyint(3) unsigned DEFAULT NULL,
	  PRIMARY KEY (`user_id`)
) ENGINE=MyISAM AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

CREATE TABLE `mod_netinfo` (
	  `mod_sn` char(100) NOT NULL DEFAULT '',
	  `user_id` int(10) unsigned NOT NULL,
	  `mod_ip` char(30) NOT NULL,
	  `mod_port` int(10) unsigned DEFAULT NULL,
	  `mod_net` char(10) DEFAULT NULL,
	  `mod_mac` char(20) DEFAULT NULL,
	  `last_update_date` datetime DEFAULT NULL,
	  PRIMARY KEY (`mod_sn`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

CREATE TABLE `mod_info` (
	  `mod_sn` char(100) NOT NULL DEFAULT '',
	  `user_id` int(10) unsigned NOT NULL,
	  `mod_name` char(10) NOT NULL,
	  `mod_alias` char(100) DEFAULT NULL,
	  `register_date` datetime DEFAULT NULL,
	  PRIMARY KEY (`mod_sn`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;


show tables;
