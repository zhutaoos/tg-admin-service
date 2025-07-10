package key_utils

import (
	"strconv"
)

func GetUserKey(uid uint) string {
	return ":user:" + strconv.Itoa(int(uid))
}

// GetTokenKey 存放 userJwt struct
func GetTokenKey(token string) string {
	return ":token:" + token
}

func GetUidToken(uid int) string {
	return ":uid-token:" + strconv.Itoa(uid)
}
