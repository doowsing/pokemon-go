package rcache

import (
	"pokemon/game/models"
)

const (
	sl_prize_info       = "sl_prize_info"
	today_sl_user       = "today_sl_user"
	today_is_use_ticket = "today_is_use_ticket"
	sl_die_option       = "sl_die_option"
)

// 获取扫雷奖励列表，从1关到第9关，每个人每天都不一样
func GetSaoleiAward(userId int) (map[int]*models.SaoLeiAwardInfo, error) {
	awards := make(map[int]*models.SaoLeiAwardInfo)
	err := RdbOperator.Hget(sl_prize_info, userId).Struct(&awards)
	return awards, err
}

// 设置扫雷奖励列表缓存
func SetSaoleiAward(userId int, awards map[int]*models.SaoLeiAwardInfo) {
	RdbOperator.Hset(sl_prize_info, userId, awards)
}

// 删除扫雷奖励列表缓存
func DelSaoleiAward(userId int) {
	RdbOperator.Hdel(sl_prize_info, userId)
}

// 设置玩家今日已扫雷
func SetSaoleiTodayUser(userId, value int) {
	RdbOperator.Hset(today_sl_user, userId, value)
}

// 获取玩家今日扫雷记录
func GetSaoleiTodayUser(userId int) (int, error) {
	return RdbOperator.Hget(today_sl_user, userId).Int()

}

// 设置玩家已用扫雷复活卡
func SetSaoleiTicketUser(userId, value int) {
	RdbOperator.Hset(today_is_use_ticket, userId, value)

}

// 获取玩家已用扫雷闯关卡
func GetSaoleiTicketUser(userId int) (int, error) {
	return RdbOperator.Hget(today_is_use_ticket, userId).Int()

}

// 设置玩家上一次死亡时的关卡，用于复活恢复关卡
func SetSaoleiDieUserLevel(userId, value int) {
	RdbOperator.Hset(sl_die_option, userId, value)
}

// 获得玩家上一次死亡时的关卡，用于复活恢复关卡
func GetSaoleiDieUserUserLevel(userId int) (int, error) {
	level, err := RdbOperator.Hget(sl_die_option, userId).Int()
	return level, err
}

// 删除玩家上一次死亡时的关卡，用于选择不复活的场景
func DelSaoleiDieUserUserLevel(userId int) {
	RdbOperator.Hdel(sl_die_option, userId)
}

// 清除扫雷记录
func ClearSaoleiTodayUser() {
	RdbOperator.Delete(today_sl_user)
}

// 清除扫雷复活记录
func ClearSaoleiTicketUser() {
	RdbOperator.Delete(today_is_use_ticket)
}

// 清除死亡关卡记录
func ClearSaoleiDieUserLevel() {
	RdbOperator.Delete(sl_die_option)
}

// 清除扫雷奖励记录
func ClearSaoleiAward() {
	RdbOperator.Delete(today_sl_user)
}
