// TODO:
// - Convert this into module/package called "collect" because use for inputs and parsers.
// - Find any way to simplify this to make more dinamyc.
// - If not set any option, ignore and no enter in infinite loop.

package inputs

import (
	"log"
	"sync"
	"time"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/audit"
	"github.com/swapbyt3s/zenit/plugins/inputs/mysql/slow"
	"github.com/swapbyt3s/zenit/plugins/inputs/os"
	"github.com/swapbyt3s/zenit/plugins/inputs/process"
	"github.com/swapbyt3s/zenit/plugins/inputs/proxysql"
	"github.com/swapbyt3s/zenit/plugins/outputs/clickhouse"
	"github.com/swapbyt3s/zenit/plugins/outputs/prometheus"
)

func Gather() {
	var wg sync.WaitGroup

	log.Printf("I! Starting Zenit %s\n", config.Version)

	wg.Add(2)

	go doCollectPlugins(&wg)
	go doCollectParsers(&wg)

	wg.Wait()
}

func doCollectPlugins(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		if config.File.OS.CPU {
			os.CPU()
		}
		if config.File.OS.Disk {
			os.Disk()
		}
		if config.File.OS.Mem {
			os.Mem()
		}
		if config.File.OS.Net {
			os.Net()
		}
		if config.File.OS.Limits {
			os.SysLimits()
		}
		if config.File.MySQL.Indexes && mysql.Check() {
			mysql.Indexes()
		}
		if config.File.MySQL.Overflow && mysql.Check() {
			mysql.Overflow()
		}
		if config.File.MySQL.Slave && mysql.Check() {
			mysql.Slave()
		}
		if config.File.MySQL.Status && mysql.Check() {
			mysql.Status()
		}
		if config.File.MySQL.Tables && mysql.Check() {
			mysql.Tables()
		}
		if config.File.MySQL.Variables && mysql.Check() {
			mysql.Variables()
		}
		if config.File.ProxySQL.Enable && proxysql.Check() {
			proxysql.ConnectionPool()
			proxysql.QueryDigest()
		}
		if config.File.Process.PerconaToolKitKill {
			process.PerconaToolKitKill()
		}
		if config.File.Process.PerconaToolKitDeadlockLogger {
			process.PerconaToolKitDeadlockLogger()
		}
		if config.File.Process.PerconaToolKitSlaveDelay {
			process.PerconaToolKitSlaveDelay()
		}
		if config.File.Prometheus.Enable {
			prometheus.Run()
		}
		accumulator.Load().Reset()
		time.Sleep(config.File.General.Interval * time.Second)
	}
}

func doCollectParsers(wg *sync.WaitGroup) {
	defer wg.Done()

	if config.File.MySQL.AuditLog.Enable {
		if config.File.General.Debug {
			log.Println("D! - Load MySQL AuditLog")
			log.Printf("D! - Read MySQL AuditLog: %s\n", config.File.MySQL.AuditLog.LogPath)
		}

		if !clickhouse.Check() {
			log.Println("E! - AuditLog require active connection to ClickHouse.")
		}

		if config.File.MySQL.AuditLog.Format == "xml-old" {
			channel_tail := make(chan string)
			channel_parser := make(chan map[string]string)
			channel_data := make(chan map[string]string)

			event := &clickhouse.Event{
				Type:    "AuditLog",
				Schema:  "zenit",
				Table:   "mysql_audit_log",
				Size:    config.File.MySQL.AuditLog.BufferSize,
				Timeout: config.File.MySQL.AuditLog.BufferTimeOut,
				Wildcard: map[string]string{
					"_time":          "'%s'",
					"command_class":  "'%s'",
					"connection_id":  "%s",
					"db":             "'%s'",
					"host":           "'%s'",
					"host_ip":        "IPv4StringToNum('%s')",
					"host_name":      "'%s'",
					"ip":             "'%s'",
					"name":           "'%s'",
					"os_login":       "'%s'",
					"os_user":        "'%s'",
					"priv_user":      "'%s'",
					"proxy_user":     "'%s'",
					"record":         "'%s'",
					"sqltext":        "'%s'",
					"sqltext_digest": "'%s'",
					"status":         "%s",
					"user":           "'%s'",
				},
				Values: []map[string]string{},
			}

			go common.Tail(config.File.MySQL.AuditLog.LogPath, channel_tail)
			go audit.Parser(config.File.MySQL.AuditLog.LogPath, channel_tail, channel_parser)
			go clickhouse.Send(event, channel_data)

			go func() {
				for channel_event := range channel_parser {
					channel_data <- channel_event
				}
			}()
		}
	}

	if config.File.MySQL.SlowLog.Enable {
		if config.File.General.Debug {
			log.Println("D! - Load MySQL SlowLog")
			log.Printf("D! - Read MySQL SlowLog: %s\n", config.File.MySQL.SlowLog.LogPath)
		}

		if !clickhouse.Check() {
			log.Println("E! - SlowLog require active connection to ClickHouse.")
		}

		channel_tail := make(chan string)
		channel_parser := make(chan map[string]string)
		channel_data := make(chan map[string]string)

		event := &clickhouse.Event{
			Type:    "SlowLog",
			Schema:  "zenit",
			Table:   "mysql_slow_log",
			Size:    config.File.MySQL.SlowLog.BufferSize,
			Timeout: config.File.MySQL.SlowLog.BufferTimeOut,
			Wildcard: map[string]string{
				"_time":         "'%s'",
				"bytes_sent":    "%s",
				"host_ip":       "IPv4StringToNum('%s')",
				"host_name":     "'%s'",
				"killed":        "%s",
				"last_errno":    "%s",
				"lock_time":     "%s",
				"query":         "'%s'",
				"query_digest":  "'%s'",
				"query_time":    "%s",
				"rows_affected": "%s",
				"rows_examined": "%s",
				"rows_read":     "%s",
				"rows_sent":     "%s",
				"schema":        "'%s'",
				"thread_id":     "%s",
				"user_host":     "'%s'",
			},
			Values: []map[string]string{},
		}

		go common.Tail(config.File.MySQL.SlowLog.LogPath, channel_tail)
		go slow.Parser(config.File.MySQL.SlowLog.LogPath, channel_tail, channel_parser)
		go clickhouse.Send(event, channel_data)

		go func() {
			for channel_event := range channel_parser {
				channel_data <- channel_event
			}
		}()
	}
}
