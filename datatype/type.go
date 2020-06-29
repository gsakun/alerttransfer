package datatype

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// AlarmData definition alarm data struct
type AlarmData struct {
	PodName           string
	PodCluster        string
	PodNamespace      string
	PodAlarmRule      string
	PodAlarmCondition string
	PodAlarmContent   string
	PodAlarmTime      string
	PodAlarmDesc      string
	PodAlarmPrority   string
}

// Handler use for sync alarmdata to mysql
func (data *AlarmData) Handler(db *sql.DB) error {
	var id int = 0
	querysql := "select id from alarm where pod_name=" + "\"" + data.PodName + "\"" + " and pod_cluster=" + "\"" + data.PodCluster + "\"" + " and pod_namespace=" + "\"" + data.PodNamespace + "\"" + " and pod_alarm_rule=" + "\"" + data.PodAlarmRule + "\""
	err := db.QueryRow(querysql).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Errorf("Query DB failed %s-%s-%s-%s %v", data.PodCluster, data.PodName, data.PodNamespace, data.PodAlarmRule, err)
			return fmt.Errorf("Query DB failed %s-%s-%s-%s %v", data.PodCluster, data.PodName, data.PodNamespace, data.PodAlarmRule, err)
		}
	}
	if id == 0 {
		stmt, err := db.Prepare(`INSERT alarm(pod_name,pod_cluster,pod_namespace,pod_alarm_rule,pod_alarm_content,pod_alarm_codition,pod_alarm_time,pod_alarm_desc,pod_alarm_prority,pod_alarm_count) values(?,?,?,?,?,?,?,?,?,?)`)
		if err != nil {
			log.Errorf("Insert prepare err %s-%s-%s-%s %v", data.PodCluster, data.PodName, data.PodNamespace, data.PodAlarmRule, err)
			return err
		}
		if err != nil {
			_, err = stmt.Exec(data.PodName, data.PodCluster, data.PodNamespace, data.PodAlarmRule, data.PodAlarmContent, data.PodAlarmCondition, data.PodAlarmTime, data.PodAlarmDesc, data.PodAlarmPrority, 1)
			log.Errorf("Insert exec err %s-%s-%s-%s %v", data.PodCluster, data.PodName, data.PodNamespace, data.PodAlarmRule, err)
			return err
		}
	} else {
		stmt, err := db.Prepare(`UPDATE alarm set pod_alarm_content=?,pod_alarm_codition=?,pod_alarm_time=?,pod_alarm_desc,pod_alarm_prority,pod_alarm_count=pod_alarm_count+1 where id=?`)
		if err != nil {
			log.Errorf("Update prepare err %s-%s-%s-%s %v", data.PodCluster, data.PodName, data.PodNamespace, data.PodAlarmRule, err)
			return err
		}
		_, err = stmt.Exec(data.PodAlarmContent, data.PodAlarmCondition, data.PodAlarmTime, data.PodAlarmDesc, data.PodAlarmPrority, id)
		if err != nil {
			log.Errorf("Update exec err %s-%s-%s-%s %v", data.PodCluster, data.PodName, data.PodNamespace, data.PodAlarmRule, err)
			return err
		}
	}
	return nil
}
