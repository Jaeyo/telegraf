package cmdexe

import (
	"io/ioutil"

	"github.com/influxdata/telegraf/plugins/inputs/zmon/api"
	"github.com/influxdata/telegraf/plugins/inputs/zmon/common"
	"github.com/influxdata/telegraf/plugins/inputs/zmon/protodata"
	"github.com/pkg/errors"
)

func OnConfigUpdate(telegrafConfigAPI api.TelegrafConfigAPI, command *protodata.ConfigUpdateCommand) error {
	telegrafConfig, err := telegrafConfigAPI.GetTelegrafConfig(int(command.TelegrafConfigID))
	if err != nil {
		return errors.Wrap(err, "failed to get telegraf config")
	}

	if err := ioutil.WriteFile(common.GetDefaultConfigPath(), []byte(telegrafConfig.Config), 0644); err != nil {
		return errors.Wrap(err, "failed to overwrite config file")
	}

	if err := onRestart(); err != nil {
		return errors.Wrap(err, "failed to reload self")
	}

	return nil
}

func OnRestart(command *protodata.RestartCommand) error {
	if err := onRestart(); err != nil {
		return errors.Wrap(err, "failed to reload self")
	}

	return nil
}

func OnLogUpload(command *protodata.LogUploadCommand) error {
	// TODO
	return nil
}

func OnTelegrafUpgrade(command *protodata.TelegrafUpgradeCommand) error {
	// TODO
	return nil
}
