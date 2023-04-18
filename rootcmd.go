package mnms

import (
	"strings"

	"github.com/qeof/q"
)

//retrieveRootCmd retrieve cmd of root ,add root name and run command
func retrieveRootCmd(cmddata map[string]CmdInfo) {
	if !QC.IsRoot {
		return
	}

	for k, v := range cmddata {
		if strings.HasPrefix(k, "config local syslog ") {
			cmd := v
			cmd.Name = QC.Name
			if cmd.Command == "" {
				cmd.Command = k
				q.Q("set command", cmd)
			}
			cmddata[k] = cmd
			go func(k string, cmdinfo CmdInfo) {
				defer func() {
					InsertCmd(k, cmdinfo)
				}()
				configOfLocalSyslogCmd(&cmdinfo)
			}(k, cmd)

		}
	}
}

func configOfLocalSyslogCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command

	if strings.HasPrefix(cmd, "config local syslog path ") && QC.IsRoot {
		return SyslogSetPathCmd(cmdinfo)
	}
	if strings.HasPrefix(cmd, "config local syslog maxsize ") && QC.IsRoot {
		return SyslogSetMaxSizeCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "config local syslog compress ") && QC.IsRoot {
		return SyslogSetCompressCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "config local syslog read") && QC.IsRoot {
		return ReadSyslogCmd(cmdinfo)
	}
	q.Q("unrecognized", cmd, len(cmd))
	cmdinfo.Status = "error: unknown command"
	return cmdinfo
}
