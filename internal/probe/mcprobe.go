package probe

import (
	"net"
	"time"
)

func IsOnline(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
