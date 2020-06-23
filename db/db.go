package db

import (
	"database/sql"
	"fmt"

	"github.com/gsakun/alerttransfer/datatype"
	log "github.com/sirupsen/logrus"
)

var DB *sql.DB

func Init(database string, maxidle, maxopen int) {
	var err error
	DB, err = sql.Open("mysql", database)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}
	DB.SetMaxIdleConns(maxidle)
	DB.SetMaxOpenConns(maxopen)
	err = DB.Ping()
	if err != nil {
		log.Fatalln("ping db fail:", err)
	}
}

func Dataprocess(data *datatype.NetData) {
	if data.ProtoType == "flow" {
		sqln := fmt.Sprintf("select id from server where ip=\"%s\" and port=%d", data.ClientIp, data.ClientPort)
		log.Println(sqln)
		var id int64
		err := DB.QueryRow(sqln).Scan(&id)
		log.Infoln(err)
		if id != 0 {
			return
		}
	}

	sql := fmt.Sprintf("select id,byte_in,byte_out,type from server where ip=\"%s\" and port=%d", data.ServerIp, data.ServerPort)
	log.Println(sql)
	var id int64
	var in, out int64
	var dtype string
	err := DB.QueryRow(sql).Scan(&id, &in, &out, &dtype)
	log.Infoln(err)
	if id == 0 {
		log.Println("START INSERT SERVER DATA")
		stmt, err := DB.Prepare(`INSERT server(ip,port,proc,byte_in,byte_out,last_time,type) values(?,?,?,?,?,?,?)`)
		if err != nil {
			log.Errorln(err)
			return
		}
		result, err := stmt.Exec(data.ServerIp, data.ServerPort, "", data.SourceBytes, data.DestBytes, data.Time, data.ProtoType)
		if err != nil {
			log.Errorln(err)
			return
		}

		/*		sql = fmt.Sprintf(
					"insert  server(ip, port, proc, byte_in,byte_out,last_time,type) values ('%s', %d, '%s', %d, %d, '%s','%s') ",
					data.ServerIp,
					data.ServerPort,
					"",
					data.SourceBytes,
					data.DestBytes,
					data.Time,
					data.ProtoType,
				)
				log.Println(sql)
				result, err := DB.Exec(sql)
				if err != nil {
					log.Errorln("exec", sql, "fail", err)
					return
				}*/
		num, err := result.LastInsertId()
		if err != nil {
			log.Errorln(err)
			log.Errorln("return insertserverid err")
			return
		}
		id = num
		log.Println("SUCCESS INSERT SERVER DATA")

	} else {
		if data.ProtoType == "flow" && data.ProtoType != dtype {
			return
		}
		log.Println("START UPDATE SERVER DATA")
		stmt, err := DB.Prepare(`UPDATE server set byte_in=?,byte_out=?,type=?, last_time=? where id=?`)
		if err != nil {
			log.Errorln(err)
			return
		}
		_, err = stmt.Exec(in+data.SourceBytes, out+data.DestBytes, data.ProtoType, data.Time, id)
		if err != nil {
			log.Errorln(err)
			return
		}
		/*		sql = fmt.Sprintf("update server set byte_in='%s', byte_out='%s',last_time='%s'",
					in+data.SourceBytes,
					out+data.DestBytes,
					data.Time,
				)
				log.Println(sql)
				_, err := DB.Exec(sql)
				if err != nil {
					log.Errorln("exec", sql, "fail", err)
				}*/
		log.Println("SUCCESS UPDATE SERVER DATA")
	}
	log.Println("START INSERT CLIENT DATA")
	sql2 := fmt.Sprintf("select id,byte_in,byte_out,type from client where ip=\"%s\" and port=%d and server_id=%d", data.ClientIp, data.ClientPort, id)
	log.Println(sql2)
	var clientid int64
	var clientin, clientout int64
	var dt string
	err = DB.QueryRow(sql2).Scan(&clientid, &clientin, &clientout, &dt)
	log.Infoln(err)
	if clientid == 0 {
		stmt, err := DB.Prepare(`INSERT client(ip,port,proc,byte_in,byte_out,last_time,server_id,type,count) values(?,?,?,?,?,?,?,?,?)`)
		if err != nil {
			log.Errorln(err)
			return
		}
		result, err := stmt.Exec(data.ClientIp, data.ClientPort, "", data.DestBytes, data.SourceBytes, data.Time, id, data.ProtoType, 1)
		if err != nil {
			log.Errorln(err)
			return
		}

		/*		sql2 = fmt.Sprintf(
					"insert client(ip, port, proc, byte_in,byte_out,last_time,type,count) values ('%s', %d, '%s', %d, %d, '%s','%s','%s')",
					data.ClientIp,
					data.ClientPort,
					"",
					data.DestBytes,
					data.SourceBytes,
					data.Time,
					data.ProtoType,
					0,
				)
				log.Println(sql2)
				result, err := DB.Exec(sql2)
				if err != nil {
					log.Errorln("exec", sql2, "fail", err)
					return
				}*/
		num, err := result.LastInsertId()
		if err != nil {
			log.Errorln(err)
			log.Errorln("return insertclientid err")
			return
		}
		clientid = num
		log.Println("SUCCESS INSERT ClIENT DATA")

	} else {
		if data.ProtoType == "flow" && data.ProtoType != dt {
			return
		}
		log.Println("START UPDATE CLIENT DATA")
		stmt, err := DB.Prepare(`UPDATE client set byte_in=?,byte_out=?,type=? ,last_time=?,port=?,count=count+1 where id=?`)
		if err != nil {
			log.Errorln(err)
			return
		}
		_, err = stmt.Exec(clientin+data.DestBytes, clientout+data.SourceBytes, data.ProtoType, data.Time, data.ClientPort, clientid)
		if err != nil {
			log.Errorln(err)
			return
		}

		/*		sql2 = fmt.Sprintf("update client set byte_in=%d, byte_out=%d,last_time='%s',count=count+1",
					clientin+data.DestBytes,
					clientout+data.SourceBytes,
					data.Time,
				)
				log.Println(sql2)
				_, err := DB.Exec(sql)
				if err != nil {
					log.Errorln("exec", sql, "fail", err)
				}*/
		log.Println("SUCCESS UPDATE CLIENT DATA")
	}
}
