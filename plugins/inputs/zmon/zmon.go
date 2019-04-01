package zmon

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/inputs/zmon/api"
	"github.com/influxdata/telegraf/plugins/inputs/zmon/cmdexe"
	"github.com/influxdata/telegraf/plugins/inputs/zmon/common"
	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"
	"github.com/pkg/errors"
)

var _zmon *Zmon
var once sync.Once

type Zmon struct {
	ZmonAccountID  string `toml:"zmon_account_id"`
	Token          string
	ServerEndpoint string `toml:"server_endpoint"`

	agentAPI          api.AgentAPI
	commandAPI        api.CommandAPI
	telegrafConfigAPI api.TelegrafConfigAPI

	once sync.Once
}

func (zmon *Zmon) Init() {
	interval, err := parseInterval()
	if err != nil {
		log.Printf("E, ZMON, failed to parse interval")
		interval = time.Duration(5 * time.Second)
	}

	zmon.agentAPI = api.NewAgentAPI(zmon.ZmonAccountID, zmon.Token, zmon.ServerEndpoint)
	zmon.commandAPI = api.NewCommandAPI(zmon.ZmonAccountID, zmon.Token, zmon.ServerEndpoint)
	zmon.telegrafConfigAPI = api.NewTelegrafConfigAPI(zmon.ZmonAccountID, zmon.Token, zmon.ServerEndpoint)

	if err := zmon.agentAPI.NotifyAgentUp(interval); err != nil {
		log.Printf("E, ZMON, failed to notify agent up: %s", err.Error())
	}

	go func() {
		signals := make(chan os.Signal)
		signal.Notify(signals, syscall.SIGHUP, syscall.SIGTERM)

		sig := <-signals
		log.Printf("I, ZMON, get signal: %v\n", sig)
		if err := zmon.agentAPI.NotifyAgentDown(); err != nil {
			log.Printf("E, ZMON, failed to notify agent down: %s", err.Error())
		}
	}()
}

func (_ *Zmon) Description() string {
	return ""
}

const sampleConfig = `
	zmon_account_id = ""
	token = ""
	server_endpoint = "localhost:8099"
`

func (_ *Zmon) SampleConfig() string {
	return sampleConfig
}

func (zmon *Zmon) Gather(_ telegraf.Accumulator) error {
	once.Do(func() { zmon.Init() })

	command, err := zmon.commandAPI.CheckCommand()
	if err != nil {
		log.Printf("E, ZMON, failed to check command: %s\n", err.Error())
	}

	if command == nil {
		return nil
	} else if command.ConfigUpdateCommand != nil {
		return cmdexe.OnConfigUpdate(zmon.telegrafConfigAPI, command.ConfigUpdateCommand)
	} else if command.RestartCommand != nil {
		return cmdexe.OnRestart(command.RestartCommand)
	} else if command.LogUploadCommand != nil {
		return cmdexe.OnLogUpload(command.LogUploadCommand)
	} else if command.TelegrafUpgradeCommand != nil {
		return cmdexe.OnTelegrafUpgrade(command.TelegrafUpgradeCommand)
	}

	return nil
}

func parseInterval() (time.Duration, error) {
	bytes, err := ioutil.ReadFile(common.GetDefaultConfigPath())
	if err != nil {
		return 0, errors.Wrap(err, "failed to read default config file")
	}

	table, err := toml.Parse(bytes)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse config file")
	}

	agentTable, ok := table.Fields["agent"].(*ast.Table)
	if !ok {
		return 0, errors.Wrap(err, "failed to parse agent config")
	}

	intervalKV, ok := agentTable.Fields["Interval"].(*ast.KeyValue)
	if !ok {
		return 0, errors.Wrap(err, "failed to parse interval key value")
	}

	intervalV, ok := intervalKV.Value.(*ast.String)
	if !ok {
		return 0, errors.Wrap(err, "failed to parse interval value")
	}

	dur, err := time.ParseDuration(intervalV.Value)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse interval duration")
	}

	return dur, nil
}

func init() {
	_zmon = &Zmon{}
	inputs.Add("zmon", func() telegraf.Input { return _zmon })
}
