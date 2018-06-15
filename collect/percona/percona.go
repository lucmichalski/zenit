package percona

import (
  "fmt"
  "gitlab.com/swapbyt3s/zenit/lib"
)

func GetRunningProcess(){
  fmt.Printf("os.process.mysqld %d\n", lib.PGrep("mysqld"))
  fmt.Printf("os.process.proxysql %d\n", lib.PGrep("proxysql"))
  fmt.Printf("os.process.pt_kill %d\n", lib.PGrep("pt-kill"))
  fmt.Printf("os.process.pt_deadlock_logger %d\n", lib.PGrep("pt-deadlock-logger"))
  fmt.Printf("os.process.pt_slave_delay %d\n", lib.PGrep("pt-slave-delay"))
}
