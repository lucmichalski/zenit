package commands

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const query = "SELECT * FROM stats_mysql_commands_counters;"

type InputProxySQLCommands struct{}

func (l *InputProxySQLCommands) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - InputProxySQLCommands - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	for host := range config.File.Inputs.ProxySQL {
		if !config.File.Inputs.ProxySQL[host].Commands {
			return
		}

		log.Info(fmt.Sprintf("Plugin - InputProxySQLCommands - Hostname: %s", config.File.Inputs.ProxySQL[host].Hostname))

		var a = metrics.Load()
		var p = mysql.GetInstance("proxysql")
		p.Connect(config.File.Inputs.ProxySQL[host].DSN)

		rows := p.Query(query)

		for i := range rows {
			a.Add(metrics.Metric{
				Key: "zenit_proxysql_commands",
				Tags: []metrics.Tag{
					{"hostname", config.File.Inputs.ProxySQL[host].Hostname},
					{"name", rows[i]["Command"]},
				},
				Values: []metrics.Value{
					{"total_time_us", common.StringToInt64(rows[i]["Total_Time_us"])},
					{"total_cnt", common.StringToInt64(rows[i]["Total_cnt"])},
					{"cnt_100us", common.StringToInt64(rows[i]["cnt_100us"])},
					{"cnt_500us", common.StringToInt64(rows[i]["cnt_500us"])},
					{"cnt_1ms", common.StringToInt64(rows[i]["cnt_1ms"])},
					{"cnt_5ms", common.StringToInt64(rows[i]["cnt_5ms"])},
					{"cnt_10ms", common.StringToInt64(rows[i]["cnt_10ms"])},
					{"cnt_50ms", common.StringToInt64(rows[i]["cnt_50ms"])},
					{"cnt_100ms", common.StringToInt64(rows[i]["cnt_100ms"])},
					{"cnt_500ms", common.StringToInt64(rows[i]["cnt_500ms"])},
					{"cnt_1s", common.StringToInt64(rows[i]["cnt_1s"])},
					{"cnt_5s", common.StringToInt64(rows[i]["cnt_5s"])},
					{"cnt_10s", common.StringToInt64(rows[i]["cnt_10s"])},
					{"cnt_infs", common.StringToInt64(rows[i]["cnt_infs"])},
				},
			})

			log.Debug(fmt.Sprintf("Plugin - InputProxySQLCommands - %#v", rows[i]))
		}
	}
}

func init() {
	inputs.Add("InputProxySQLCommands", func() inputs.Input { return &InputProxySQLCommands{} })
}
