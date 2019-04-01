// +build !windows

package common

func GetDefaultConfigPath() string {
	return "/etc/telegraf/telegraf.conf"
}
