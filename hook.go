package sockrus

import (
	logrus_logstash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
	"net"
)

// Hook represents a connection to a socket
type Hook struct {
	formatter  logrus.Formatter
	protocol   string
	address    string
	addNewline bool
}

// NewHook establish a socket connection.
// Protocols allowed are: "udp", "tcp", "unix" (corresponds to SOCK_STREAM),
// "unixdomain" (corresponds to SOCK_DGRAM) or "unixpacket" (corresponds to SOCK_SEQPACKET).
//
// For TCP and UDP, address must have the form `host:port`.
//
// For Unix networks, the address must be a file system path.
func NewHook(protocol, address string, addNewline bool) (*Hook, error) {
	logstashFormatter := logrus_logstash.DefaultFormatter(logrus.Fields{})
	return &Hook{
		protocol:   protocol,
		address:    address,
		formatter:  logstashFormatter,
		addNewline: addNewline,
	}, nil
}

// Fire send log to the defined socket
func (h *Hook) Fire(entry *logrus.Entry) error {
	var err error
	dataBytes, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}

	conn, err := net.Dial(h.protocol, h.address)
	if err != nil {
		return nil
	}
	defer conn.Close()
	// Add new line to every message if desired.
	if h.addNewline {
		dataBytes = append(dataBytes, "\n"...)
	}
	_, _ = conn.Write(dataBytes) // #nosec
	return nil
}

// Levels return an array of handled logging levels
func (h *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}
