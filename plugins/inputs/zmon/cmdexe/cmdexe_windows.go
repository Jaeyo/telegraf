// +build windows

package cmdexe

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/influxdata/telegraf/internal"
	"github.com/pkg/errors"
)

func onRestart() error {
	tmpScript, err := ioutil.TempFile("", "*.cmd")
	if err != nil {
		return errors.Wrap(err, "failed to create temp script")
	}

	scriptContent := `
		net stop telegraf
		net start telegraf
	`

	if err = ioutil.WriteFile(tmpScript.Name(), []byte(scriptContent), 0644); err != nil {
		return errors.Wrap(err, "failed to write into temp script")
	}

	tmpScriptName := tmpScript.Name()
	tmpScript.Close()

	cmd := exec.Command(tmpScriptName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err = internal.RunTimeout(cmd, 10*time.Second); err != nil {
		s := stderr.String()
		return errors.Wrap(err, fmt.Sprintf("failed to execute temp script, stderr: %s", s))
	}

	return nil
}
