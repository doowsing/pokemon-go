package rcache

import "strconv"

func SetGroupID(userId int, groupId string) {
	GetGCache().Set(GroupKey+strconv.Itoa(userId), groupId)
}

func GetGroupID(userId int) string {
	if data, ok := GetGCache().Get(GroupKey + strconv.Itoa(userId)); ok {
		return data.(string)
	}
	return ""
}

func DelGroupID(userId int) {
	GetGCache().Del(GroupKey + strconv.Itoa(userId))
}
