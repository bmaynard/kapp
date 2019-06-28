package clusterapply

import "crypto/md5"

var messagesNotified = make(map[[16]byte]int)

func notifyChange(msg string) bool {
	hash := md5.Sum([]byte(msg))

	_, present := messagesNotified[hash]

	if present && messagesNotified[hash]%10 != 0 {
		messagesNotified[hash]++
		return false
	}

	messagesNotified[hash] = 1
	return true
}
