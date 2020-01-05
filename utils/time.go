package utils

import (
	"time"
)

// obtain the TimeStamp after lag second from now
func TimeStamp(lag int64) int64 {
	return time.Now().Unix() + lag
}
