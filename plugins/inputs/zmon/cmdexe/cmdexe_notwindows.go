// +build !windows

package cmdexe

import (
	"os"
	"syscall"

	"github.com/pkg/errors"
)

func onRestart() error {
	p, err := os.FindProcess(syscall.Getpid())
	if err != nil {
		return errors.Wrap(err, "failed to find process")
	}
	return p.Signal(syscall.SIGHUP)
}
