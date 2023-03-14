package common

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
)

func MyMd5(str string) string  {

	str = AdminSecret+str
	has := md5.New()
	has.Write([]byte(str))
	b := has.Sum(nil)
	return hex.EncodeToString(b)
}
//权限判断
func HasPower(pids string ,pid int)  bool{
	 id := strconv.Itoa(pid)

	ids := ","+id+","

	return strings.Contains(pids,ids)
}