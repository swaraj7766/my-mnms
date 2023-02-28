package atopmib

type atopvalue struct {
	backupServerIP                string
	backupAgentBoardFwFileName    string
	backupStatus                  int
	restoreServerIP               string
	restoreAgentBoardFwFileName   string
	restoreStatus                 int
	syslogStatus                  int
	eventServerPort               int
	eventServerLevel              int
	eventLogToFlash               int
	eventServerIP                 string
	sntpClientStatus              int
	sntpUTCTimezone               int
	sntpServer1                   string
	sntpServer2                   string
	sntpServerQueryPeriod         int
	agingTimeSetting              int
	ptpState                      int
	ptpVersion                    int
	ptpSyncInterval               int
	ptpClockStratum               int
	ptpPriority1                  int
	ptpPriority2                  int
	rstpStatus                    int
	qosCOSPriorityQueue           int
	qosTOSPriorityQueue           int
	eventPortEventEmail           int
	eventPortEventRelay           int
	eventPowerEventSMTP1          int
	eventPowerEventSMTP2          int
	syslogEventsSMTP              int
	eventEmailAlertAddr           string
	eventEmailAlertAuthentication int
	eventEmailAlertAccount        string
	lldpStatus                    int
	trapServerStatus              int
	trapServerIP                  string
	trapServerPort                int
	trapServerTrapComm            string
}