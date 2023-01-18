package sockrus

import (
	fqdn "github.com/Showmax/go-fqdn"
	"github.com/sirupsen/logrus"
)

// Config serves as means to configure logger and hook.
type Config struct {
	Hostname       string       // Hostname of the machine we are logging from.
	LogLevel       logrus.Level // Log level of messages we want to send to socket.
	Service        string       // Service that is creating the logs.
	SocketAddr     string       // Address to the socket.
	SocketProtocol string       // Protocol of the socket.
	AddNewLine     bool         // Toggle to send newline after every message to socket.
}

// NewSockrus is a wrapper for initialization of logrus with sockrus hook. It
// returns new instance of logrus.Logger and logrus.Entry. All errors are fatal.
func NewSockrus(config Config) (*logrus.Logger, *logrus.Entry) {
	if config.Service == "" {
		config.Service = "unknown"
	}

	logInstance := logrus.New()
	logInstance.Level = config.LogLevel

	// Get hostname.
	if config.Hostname == "" {
		config.Hostname, _ = fqdn.FqdnHostname()
		if config.Hostname == "unknown" {
			logInstance.WithFields(logrus.Fields{
				"hostname": config.Hostname,
				"msg_type": "log",
				"service":  config.Service,
			}).Fatal("Failed to get FQDN of machine I'm running at.")
		}
	}

	hook, err := NewHook(config.SocketProtocol, config.SocketAddr, config.AddNewLine)
	if err != nil {
		logInstance.WithFields(logrus.Fields{
			"hostname": config.Hostname,
			"msg_type": "log",
			"service":  config.Service,
			"error":    err.Error(),
		}).Fatal("Failed to add Unix Socket Hook.")
	}
	logInstance.Hooks.Add(hook)

	log := logInstance.WithFields(logrus.Fields{
		"hostname": config.Hostname,
		"msg_type": "log",
		"service":  config.Service,
	})
	return logInstance, log
}
