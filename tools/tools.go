package tools

import (
	"strconv"
	"time"
	"crypto/md5"
	"encoding/hex"
)

/*
Create a unique name
params: text - is a name
params: userID - is a user _id generate request in MongoDB
 */
func SetName(text, userID string) string {
	now := strconv.FormatInt(time.Now().Unix(), 10)
	hashText := md5.Sum([]byte(text + now + userID))
	str := hex.EncodeToString(hashText[:])
	return str
}
