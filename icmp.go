// Thanks go to https://github.com/Soontao/cachet-monitor/blob/master/tcp.go
package cachet

import (
	"fmt"

	"github.com/chenjiandongx/yap"
)

// Investigating template
var defaultICMPInvestigatingTpl = MessageTemplate{
	Subject: `{{ .Monitor.Name }} - {{ .SystemName }}`,
	Message: `{{ .Monitor.Name }} check **failed** (server time: {{ .now }})

{{ .FailReason }}`,
}

// Fixed template
var defaultICMPFixedTpl = MessageTemplate{
	Subject: `{{ .Monitor.Name }} - {{ .SystemName }}`,
	Message: `**Resolved** - {{ .now }}

Down seconds: {{ .downSeconds }}s`,
}

// ICMPMonitor struct
type ICMPMonitor struct {
	AbstractMonitor `mapstructure:",squash"`
}

// CheckICMPAlive func
func CheckICMPAlive(ip string, timeout int) (bool, error) {

	pg, err := yap.NewPinger();
	if err != nil {
		defer pg.Close()
	}

	response := pg.Call(yap.Request{
		Target: ip, 
		Count: 1,
		Timeout: timeout * 1000})

	if response.Error != nil {
		return false, response.Error
	}

	return true, nil

}

// test if it available
func (m *ICMPMonitor) test() bool {
	if alive, e := CheckICMPAlive(m.Target, int(m.Timeout)); alive {
		return true
	} else {
		m.lastFailReason = fmt.Sprintf("ICMP check failed: %v", e)
		return false
	}
}

// Validate configuration
func (m *ICMPMonitor) Validate() []string {

	// set incident temp
	m.Template.Investigating.SetDefault(defaultICMPInvestigatingTpl)
	m.Template.Fixed.SetDefault(defaultICMPFixedTpl)

	// super.Validate()
	errs := m.AbstractMonitor.Validate()

	if m.Target == "" {
		errs = append(errs, "Target is required")
	}

	return errs
}
