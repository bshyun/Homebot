package mydb

import (
	"../myerror"
	"../svcutil"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net"
	"time"
)

type MyDB struct {
	sql.DB
	userid string
}

func NewMyDB(msgr_id string, cfg *svcutil.DSN_config) (*MyDB, error) {
	if cfg == nil {
		return nil, fmt.Errorf(myerror.ERROR_ARG_NIL)
	}

	v, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	var m MyDB
	if len(msgr_id) == 0 {
		m = MyDB{*v, ""}
	} else {
		query := fmt.Sprintf("select user_id from user_info where user_info.`msgr_id` like '%s'", msgr_id)
		result, err1 := v.Query(query)
		if err != nil {
			fmt.Println(err1)
			v.Close()
			return nil, err1
		}
		defer result.Close()

		var user_id string
		for result.Next() {
			if err1 = result.Scan(&user_id); err1 != nil {
				fmt.Println(err1)
				v.Close()
				return nil, err1
			}
		}

		m = MyDB{*v, user_id}
	}

	return &m, nil

}

func (db *MyDB) SetUserID(ID int) error {
	if db == nil {
		return fmt.Errorf(myerror.ERROR_ARG_NIL)
	}

	msgr_id := fmt.Sprintf("%s_%d", "telegram", ID)

	query := fmt.Sprintf("select user_info.user_id from user_info where user_info.msgr_id='%s'", msgr_id)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer result.Close()

	var user_id string
	for result.Next() {
		if err = result.Scan(&user_id); err != nil {
			return err
		}

		db.userid = user_id
		break
	}
	if len(user_id) == 0 {
		return fmt.Errorf("Not found user id")
	}

	return nil
}

/*
func (db *MyDB) Close() {
	db.Close()
}
*/

func (db *MyDB) GetMyIP(userid string, mod_name string) (string, error) {
	if db == nil {
		return "", fmt.Errorf(myerror.ERROR_ARG_NIL)
	}

	//fmt.Println("mod name:", mod_name)

	var query string
	if len(db.userid) == 0 {
		query = fmt.Sprintf("select mod_ip from mod_netinfo where mod_sn = (select mod_sn from mod_info where (mod_info.mod_alias = '%s' or mod_info.mod_name = '%s') and mod_info.user_id = (select user_info.user_id from user_info where user_info.msgr_id='%s'))", mod_name, mod_name, userid)
	} else {
		query = fmt.Sprintf("select mod_ip from mod_netinfo where mod_sn = (select mod_sn from mod_info where (mod_info.mod_alias = '%s' or mod_info.mod_name = '%s') and mod_info.user_id = %s)", mod_name, mod_name, db.userid)
	}
	//query := fmt.Sprintf("select * from user_info where msgr_id='%s'", userid)
	//result, err := db.Query("select * from user_info where msgr_id='", userid, "'")
	result, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer result.Close()

	var mod_ip string
	for result.Next() {
		/*
			var (
				user_id   int
				msgr_id   string
				msgr_id2  string
				last      string
				first     string
				msgr_name string
				register  string
				phone     string
				addr      string
				dep       int
			)
		*/
		if err = result.Scan(&mod_ip); err != nil {
			//if err = result.Scan(&user_id); err != nil {
			//if err = result.Scan(&user_id, &msgr_id, &msgr_id2, &last, &first, &msgr_name, &register, &phone, &addr, &dep); err != nil {
			log.Fatalln(err)
			return "", nil
		}
		//fmt.Println("Result: ", " ", user_id, msgr_id, msgr_id2, last, first, msgr_name, register, phone, addr, dep, "|||")
	}

	return mod_ip, nil
}

func (db *MyDB) FindMe(mod_sn string) bool {
	if db == nil {
		return false
	}

	query := fmt.Sprintf("select mod_ip from mod_netinfo where mod_sn like '%s'", mod_sn)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer result.Close()

	var rowcount int = 0
	var mod_ip string
	for result.Next() {
		if err = result.Scan(&mod_ip); err != nil {
			log.Fatalln(err)
			return false
		}
		rowcount++
	}

	if rowcount < 1 {
		return false
	}

	return true
}

func (db *MyDB) UpdateModIP(mod_addr net.Addr, mod_sn string) error {
	if db == nil {
		return fmt.Errorf(myerror.ERROR_ARG_NIL)
	}

	//fmt.Printf("UpdateModIP: [addr:%s, sn:%s]\n", mod_addr, mod_sn)

	rip, _, erra := net.SplitHostPort(mod_addr.String())
	if erra != nil {
		return erra
	}

	stmt, errst := db.Prepare("update mod_netinfo set mod_ip = ?, last_update_date = ? where mod_sn = ?")
	if errst != nil {
		fmt.Println("Update db error:", errst)
		return nil
	}
	defer stmt.Close()
	now := time.Now()
	y, mo, d := now.Date()
	h, mi, s := now.Clock()
	date := fmt.Sprintf("%d-%d-%d %d:%d:%d", y, mo, d, h, mi, s)
	_, err := stmt.Exec(rip, date, mod_sn)
	if err != nil {
		fmt.Println("Update db error:", err)
		return err
	}
	/*
		iid, _ := result.LastInsertId()
		rid, _ := result.RowsAffected()
		fmt.Println(iid, rid)
	*/

	return nil
}

func (db *MyDB) Close() error {
	//fmt.Println("before db close")
	err := db.DB.Close()
	//fmt.Println("after db close")

	if err != nil {
		return nil
	}
	return nil
}

func (db *MyDB) GetMyCommand(ID int) (string, error) {
	if db == nil {
		return "", fmt.Errorf(myerror.ERROR_ARG_NIL)
	}

	userid := fmt.Sprint("telegram_", ID)
	var query string
	if len(db.userid) == 0 {
		query = fmt.Sprintf("select mod_sn from mod_info where mod_info.user_id = (select user_id from user_info where user_info.`msgr_id` like '%s')", userid)
	} else {
		query = fmt.Sprintf("select mod_sn from mod_info where mod_info.user_id = %s", db.userid)
	}
	result, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer result.Close()

	rowcnt := 0
	var mod_sn string
	var command string
	//command := make([]string, 1, 100)
	for result.Next() {
		if err = result.Scan(&mod_sn); err != nil {
			log.Fatalln(err)
			return "", err
		}

		query1 := fmt.Sprintf("select command from mod_command where mod_command.mod_sn = '%s'", mod_sn)
		result1, err1 := db.Query(query1)
		if err1 != nil {
			fmt.Println(err1)
			return "", err1
		}
		defer result1.Close()
		for result1.Next() {
			var res string
			if err1 = result1.Scan(&res); err1 != nil {
				log.Fatalln(err1)
				return "", err1
			}

			command += res + ","

			//fmt.Println(command[rowcnt])
			rowcnt++
		}
	}

	/*
		if rowcnt > 0 {
			command = command[:len(command)-1]
		}
	*/

	//fmt.Println(command)

	return command, nil
}

func (db *MyDB) GetMyDev(ID int) (string, error) {
	if db == nil {
		return "", fmt.Errorf(myerror.ERROR_ARG_NIL)
	}

	//fmt.Println("db.userid: ", db.userid)
	userid := fmt.Sprint("telegram_", ID)
	var query string
	if len(db.userid) == 0 {
		query = fmt.Sprintf("select mod_sn from mod_info where mod_info.user_id = (select user_id from user_info where user_info.`msgr_id` like '%s')", userid)
	} else {
		query = fmt.Sprintf("select mod_sn from mod_info where mod_info.user_id = %s", db.userid)
	}
	result, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer result.Close()

	rowcnt := 0
	var mod_sn string
	var command string
	for result.Next() {
		if err = result.Scan(&mod_sn); err != nil {
			log.Fatalln(err)
			return "", err
		}

		query1 := fmt.Sprintf("select mod_sn, mod_name, mod_alias from mod_info where mod_info.mod_sn = '%s'", mod_sn)
		result1, err1 := db.Query(query1)
		if err1 != nil {
			fmt.Println(err1)
			return "", err1
		}
		defer result1.Close()
		for result1.Next() {
			var res1, res2, res3 string
			if err1 = result1.Scan(&res1, &res2, &res3); err1 != nil {
				log.Fatalln(err1)
				return "", err1
			}

			command += res2 + ":" + res3 + ":" + res1 + " "

			//fmt.Println(command[rowcnt])
			rowcnt++
		}
	}

	/*
		if rowcnt > 0 {
			command = command[:len(command)-1]
		}
	*/

	//fmt.Println(command)

	return command, nil
}
