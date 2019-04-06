package config

import (
	"time"
)

// All is a struct to contain all configuration imported or loaded from config file.
type All struct {
	General struct {
		Hostname string        `yaml:"hostname"`
		Interval time.Duration `yaml:"interval"`
		Debug    bool          `yaml:"debug"`
	}
	Parser struct {
		MySQL struct {
			AuditLog struct {
				Enable        bool   `yaml:"enable"`
				Format        string `yaml:"format"`
				LogPath       string `yaml:"log_path"`
				BufferSize    int    `yaml:"buffer_size"`
				BufferTimeOut int    `yaml:"buffer_timeout"`
			}
			SlowLog struct {
				Enable        bool   `yaml:"enable"`
				LogPath       string `yaml:"log_path"`
				BufferSize    int    `yaml:"buffer_size"`
				BufferTimeOut int    `yaml:"buffer_timeout"`
			}
		}
	}
	MySQL struct {
		DSN string `yaml:"dsn"`
		Inputs struct {
			Overflow  bool `yaml:"overflow"`
			Slave     bool `yaml:"slave"`
			Status    bool `yaml:"status"`
			Tables    bool `yaml:"tables"`
			Variables bool `yaml:"variables"`
		}
	}
	ProxySQL []struct {
		Hostname string `yaml:"hostname"`
		DSN      string `yaml:"dsn"`
		Inputs struct {
			Commands bool `yaml:"commands"`
			Pool     bool `yaml:"pool"`
			Queries  bool `yaml:"queries"`
		}
	}
	ClickHouse struct {
		DSN string `yaml:"dsn"`
	}
	Prometheus struct {
		Enable   bool   `yaml:"enable"`
		TextFile string `yaml:"textfile"`
	}
	OS struct {
		Inputs struct {
			CPU    bool `yaml:"cpu"`
			Disk   bool `yaml:"disk"`
			Limits bool `yaml:"limits"`
			Mem    bool `yaml:"mem"`
			Net    bool `yaml:"net"`
		}
	}
	Process struct {
		Inputs struct {
			PerconaToolKitDeadlockLogger     bool `yaml:"pt_deadlock_logger"`
			PerconaToolKitKill               bool `yaml:"pt_kill"`
			PerconaToolKitOnlineSchemaChange bool `yaml:"pt_online_schema_change"`
			PerconaToolKitSlaveDelay         bool `yaml:"pt_slave_delay"`
			PerconaXtraBackup                bool `yaml:"xtrabackup"`
		}
	}
}
