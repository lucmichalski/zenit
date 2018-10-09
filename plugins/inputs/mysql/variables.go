package mysql

import (
	"database/sql"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/accumulator"
)

const QuerySQLVariables = "SHOW GLOBAL VARIABLES"

func Variables() {
	conn, err := mysql.Connect(config.File.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Error("MySQL:Variables - Impossible to connect: " + err.Error())
	}

	rows, err := conn.Query(QuerySQLVariables)
	defer rows.Close()
	if err != nil {
		log.Error("MySQL:Variables - Impossible to execute query: " + err.Error())
	}

	var a = accumulator.Load()
	var k string
	var v sql.RawBytes

	for rows.Next() {
		rows.Scan(&k, &v)
		if value, ok := mysql.ParseValue(v); ok {
			if config.File.General.Debug {
				// log.Printf("D! - MySQL:Variables - %s=%d\n", k, value)
			}

			a.Add(accumulator.Metric{
				Key:    "mysql_variables",
				Tags:   []accumulator.Tag{{"name", k}},
				Values: value,
			})
		}
	}
}
