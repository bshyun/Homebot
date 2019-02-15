use homebot;

insert into user_info (msgr_id, msgr_id2, last_name, first_name, msgr_name, register_date, user_phone, user_addr, deprecated) values
(
	"telegram_111111", "111111", "last", "first", "winot", "2016-08-05 13:33:00", "010-1234-5678", "동해물과 백두산이", 0
);

insert into mod_info (mod_sn, user_id, mod_name, mod_alias, register_date) values
(
	"26f7262e-b4df-42dd-b95a-d21df9632e6e", 0, "temper", NULL, "2016-08-05 13:34:00"
);

insert into mod_netinfo (mod_sn, user_id, mod_ip, mod_port, mod_net, mod_mac, last_update_date) values
(
	"26f7262e-b4df-42dd-b95a-d21df9632e6e", 0, "192.168.39.1", 12120, "tcp", "00:50:56:c0:00:08", "2016-08-05 13:35:00"
);
