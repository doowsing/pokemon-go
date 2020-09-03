package services

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	phpserialize "github.com/techleeone/gophp/serialize"
	"github.com/unknwon/com"
	"math/rand"
	common2 "pokemon/common"
	"pokemon/common/persistence"
	"pokemon/common/rcache"
	"pokemon/common/rpc-client/rpc-group"
	"pokemon/game/models"
	"pokemon/game/services/common"
	"pokemon/game/utils"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SysService struct {
	BaseService
}

func NewSysService(osrc *OptService) *SysService {
	srv := &SysService{}
	srv.SetOptSrc(osrc)
	return srv
}

func (ss *SysService) InitRdModels() bool {
	ps := NewPetService(nil)
	if err := rcache.SetMPets(ps.GetAllMPet()); err != nil {
		fmt.Println("初始化宠物模型到redis失败，error:", err)
	}
	if err := rcache.SetMSkill(ps.GetAllMSkill()); err != nil {
		fmt.Println("初始化技能模型到redis失败，error:", err)
	}
	prs := NewPropService(nil)
	MProps, ok := prs.Repo.GetAllMPropFromMysql()
	if ok {
		if err := rcache.SetMProps(MProps); err != nil {
			fmt.Println("初始化道具模型到redis失败，error:", err)
		}
	}

	//fs := fight.NewFightService(nil)

	//MMaps, ok := fs.repo.GetAllMapFromMysql()
	//if ok {
	//	if _, err := rcache.SetMMaps(MMaps); err != nil {
	//		fmt.Println("初始化地图模型到redis失败，error:", err)
	//	}
	//}
	//MGpcs, ok := fs.repo.GetAllGpcFromMysql()
	//if ok {
	//	if _, err := rcache.SetMGpcs(MGpcs); err != nil {
	//		fmt.Println("初始化怪物模型到redis失败，error:", err)
	//	}
	//}
	//MGpcGroup, ok := fs.repo.GetAllGpcGroupFromMysql()
	//if ok {
	//	if _, err := rcache.SetMGpcGroup(MGpcGroup); err != nil {
	//		fmt.Println("初始化怪物队列模型到redis失败，error:", err)
	//	}
	//}
	return true
}
func (ss *SysService) InitSettings() {

}

func (ss *SysService) LoginLog(uname, ip string, t int) {
	ss.GetDb().Create(&models.LoginLog{UName: uname, IP: ip, Time: t})
}

var GameLogQueue = make(chan *models.GameLog)

var vipProps = []int{1427, 1474, 1475, 1476, 1477, 1478, 1479, 1480, 1481, 1482, 1483, 1484, 1485}

func SelfGameLog(suid int, note string, vary int) {
	GameLogQueue <- &models.GameLog{SUid: suid, BUid: suid, Note: note, Category: vary, Time: utils.NowUnix()}
}

func GameLog(suid, buid int, note string, vary int) {
	GameLogQueue <- &models.GameLog{SUid: suid, BUid: buid, Note: note, Category: vary, Time: utils.NowUnix()}
}

func AnnounceSys(note string) {
	str := fmt.Sprintf("[系统公告]%s", note)
	//AddMsgQueue(str)
	common2.AnnounceChat(str)
	fmt.Print(str + "\n")
}

func AnnounceAll(username, note string) {
	str := fmt.Sprintf("[系统公告]恭喜玩家 %s %s", username, note)
	//AddMsgQueue(str)
	common2.AnnounceChat(str)
	fmt.Print(str + "\n")
}

func AnnouncePet2All(userId int, pet *models.UPet) {
	pet.GetM()
	str := fmt.Sprintf("show pet :%s", pet.MModel.Name)
	//AddMsgQueue(str)
	common2.ShowPetChat(userId, pet.ID, pet.MModel.Name)
	fmt.Print(str + "\n")
}

func AnnounceProp2All(userId int, prop *models.UProp) {
	prop.GetM()
	str := fmt.Sprintf("show prop :%s", prop.MModel.Name)
	//AddMsgQueue(str)
	common2.ShowPropChat(userId, prop.ID, prop.MModel.Name, prop.MModel.PropsColor)
	fmt.Print(str + "\n")
}

func (ss *SysService) SendEmail(userId int, now time.Time, msg string) {
	msgInfo := &models.EmailMsg{
		Uid:     userId,
		Time:    now,
		Content: msg,
	}
	ss.GetDb().Create(msgInfo)
}

func (ss *SysService) SaveGameLog() {
	rkey := "gameLogList"
	fmt.Print("开始收集插入游戏日志！")
	for true {
		logLen, _ := rcache.LLen(rkey)
		if logLen < 1 {
			time.Sleep(time.Millisecond * 500)
			continue
		}
		if logLen > 50 {
			logLen = 50
		}
		logListStr, _ := rcache.LRange(rkey, 0, logLen)
		logList := []*models.GameLog{}
		for _, s := range logListStr {
			//fmt.Println("string:%s", string(s))
			var logMap map[string]interface{}
			if err := json.Unmarshal(s, &logMap); err == nil {
				newLog := &models.GameLog{}
				ltime := logMap["time"]
				newLog.Time = int(ltime.(float64))
				note, ok := logMap["note"].(string)
				if ok {
					newLog.Note = note
				} else {
					newLog.Note = ""
				}
				newLog.Note = strings.ReplaceAll(newLog.Note, "\\n", "\n")
				newLog.Category = int(logMap["category"].(float64))
				newLog.SUid = int(logMap["suid"].(float64))
				newLog.BUid = int(logMap["buid"].(float64))
				logList = append(logList, newLog)
			} else {
				fmt.Println("反编码失败！err:%s", err)
			}
		}
		for _, log := range logList {
			//fmt.Println(log.String())
			GameLogQueue <- log
		}
		rcache.LTrim(rkey, logLen+1, -1)
	}
}

func (ss *SysService) GetGameLogListLen() int {
	rkey := "gameLogList"
	logLen, _ := rcache.LLen(rkey)
	return logLen
}

func (ss *SysService) GetLoginLogListLen() int {
	rkey := "LoginLogList"
	logLen, _ := rcache.LLen(rkey)
	return logLen
}

func (ss *SysService) SaveLoginLog() {
	rkey := "LoginLogList"
	fmt.Print("开始收集插入登录日志！")
	for true {
		logLen, _ := rcache.LLen(rkey)
		if logLen < 1 {
			time.Sleep(time.Millisecond * 500)
			continue
		}
		if logLen > 50 {
			logLen = 50
		}
		logListStr, _ := rcache.LRange(rkey, 0, logLen)
		logList := []*models.LoginLog{}
		for _, s := range logListStr {
			//fmt.Println("string:%s", string(s))
			var logMap map[string]interface{}
			if err := json.Unmarshal(s, &logMap); err == nil {
				newLog := &models.LoginLog{}
				newLog.Time = int(logMap["time"].(float64))
				newLog.UName = logMap["username"].(string)
				newLog.IP = logMap["ip"].(string)
				logList = append(logList, newLog)
			} else {
				fmt.Println("反编码失败！string:%s", string(s))
				//panic(err)
				continue
			}
		}
		ts := persistence.GetLogOrm().Begin()

		for _, log := range logList {
			ts.Create(log)
		}

		rcache.LTrim(rkey, logLen+1, -1)
		ts.Commit()
	}
}

func (ss *SysService) PrepareTestDate() {
	rkey := "gameLogList"
	fmt.Println("开始插入模拟日志队列。。。")
	for i := 0; i < 1000; i++ {
		newLog := &models.GameLog{}
		newLog.SUid = i
		newLog.BUid = i
		newLog.Note = fmt.Sprintf("这是第%d条记录", i)
		newLog.Category = 5
		newLog.Time = utils.NowUnix() + rand.Int()
		str, _ := json.Marshal(newLog)
		cnt, err := rcache.RPush(rkey, string(str))
		if err != nil {
			fmt.Println("插入队列失败！")
		} else {
			fmt.Printf("队列长度:%d\n", cnt)
		}
	}
	fmt.Println("插入模拟日志队列完毕！")
}

func (ss *SysService) AutoInsertGameLog() {
	tmpSlice := []*models.GameLog{}
	exeute := false
	for {
		select {
		case log := <-GameLogQueue:
			tmpSlice = append(tmpSlice, log)
		case <-time.After(time.Millisecond * 500):
			if len(tmpSlice) > 0 {
				exeute = true
			}
		}
		if len(tmpSlice) > 10 || (len(tmpSlice) > 0 && exeute) {
			ts := persistence.GetLogOrm().Begin()
			for _, log := range tmpSlice {
				fmt.Println(log.String())
				ts.Create(log)
			}
			ts.Commit()
			tmpSlice = []*models.GameLog{}
			exeute = false
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func (ss *SysService) GetPublicContent() string {
	info := common.GetWelcome("public")
	if info != nil {
		return info.Content
	}
	return ""
}

func (ss *SysService) GetPublicRankLists(update bool) map[string][]gin.H {
	rankList := make(map[string][]gin.H)
	if !update {
		data, err := rcache.GetRankList()
		if err == nil && len(data) == 4 {

			for k, v := range data {
				v1 := []gin.H{}
				for _, vv := range v {
					v1 = append(v1, vv)
				}
				rankList[k] = v1
			}
			return rankList
		}
	}
	// 等级排行榜
	levelRank := []gin.H{}
	sql := `SELECT b.id as id,b.name as name,b.level as level,u.nickname as nickname
								FROM userbb as b left join player as u on u.id = b.uid
							   ORDER BY Level DESC,nowexp DESC
							   LIMIT 0,30`
	levelList := []struct {
		Id       int
		Name     string
		Level    int
		Nickname string
	}{}
	ss.GetDb().Raw(sql).Scan(&levelList)
	for i, element := range levelList {
		levelRank = append(levelRank, gin.H{
			"rank":     i + 1,
			"petid":    element.Id,
			"petname":  element.Name,
			"level":    element.Level,
			"nickname": element.Nickname,
		})
	}
	// 神圣成长排行榜
	ssczRank := []gin.H{}
	sql = `SELECT userbb.id as id,bb.name as name, player.nickname as nickname,userbb.czl as czl FROM userbb 
left join bb on userbb.bid=bb.id
left join player on userbb.uid=player.id 
WHERE bb.wx = 7 ORDER by userbb.czl+0 desc limit 15`
	ssczList := []struct {
		Id       int
		Name     string
		Czl      string
		Nickname string
	}{}
	ss.GetDb().Raw(sql).Scan(&ssczList)
	for i, element := range ssczList {
		ssczRank = append(ssczRank, gin.H{
			"rank":     i + 1,
			"petid":    element.Id,
			"petname":  element.Name,
			"czl":      element.Czl,
			"nickname": element.Nickname,
		})
	}
	// 普通成长排行榜
	czRank := []gin.H{}
	czList := []struct {
		Id       int
		Name     string
		Czl      string
		Nickname string
	}{}
	czRankLen := 50
	czRankSetting := common.GetWelcome("czpb_count")
	if czRankSetting != nil {
		czRankLen, _ = strconv.Atoi(czRankSetting.Content)
	}
	sql = `SELECT userbb.id as id,bb.Name as name, player.nickname as nickname,userbb.czl as czl FROM userbb 
left join bb on userbb.bid=bb.id 
left join player on userbb.uid=player.id 
WHERE bb.wx != 7 ORDER by userbb.czl+0 desc limit ?`
	ss.GetDb().Raw(sql, czRankLen).Scan(&czList)
	for i, element := range czList {
		czRank = append(czRank, gin.H{
			"rank":     i + 1,
			"petid":    element.Id,
			"petname":  element.Name,
			"czl":      element.Czl,
			"nickname": element.Nickname,
		})
	}
	// 成长增长排行榜
	czzcRank := []gin.H{}

	now := time.Now()
	var month string
	month = strconv.Itoa(int(now.Month()))
	if now.Month() < 10 {
		month = "0" + month
	}
	key := fmt.Sprintf("%d%smonthcczz", now.Year(), month)
	record := common.GetWelcome(key)
	if record != nil {
		recordList, err := phpserialize.UnMarshal([]byte(record.Content))
		if err == nil {
			if dataList, ok := recordList.([]interface{}); ok {
				for _, d := range dataList {
					userRecord := d.(map[string]interface{})
					uid := com.StrTo(userRecord["uid"].(string)).MustInt()
					maxCzl := com.StrTo(userRecord["max_czl"].(string)).MustFloat64()
					userCzl := struct {
						Czl      float64
						Nickname string
					}{}
					//sql = "SELECT max(bb.czl+0) as czl,u.nickname as nickname from userbb bb left join player u on bb.uid=u.id where bb.uid= ? and bb.wx != 7 limit 1"
					sql = "select nickname, mczl.czl from player, (SELECT max(bb.czl+0) as czl FROM userbb bb where bb.uid= ? and bb.wx != 7 LIMIT 1) mczl WHERE player.id=? limit 1"
					ss.GetDb().Raw(sql, uid, uid).Scan(&userCzl)
					addCzl := userCzl.Czl - maxCzl
					if addCzl < 0 {
						addCzl = 0
					}
					czzcRank = append(czzcRank, gin.H{
						"nickname": userCzl.Nickname,
						"czl":      addCzl,
					})
				}
			}
		}
	}
	sort.Slice(czzcRank, func(i, j int) bool {
		return czzcRank[i]["czl"].(float64) > czzcRank[j]["czl"].(float64)
	})
	for i := 0; i < len(czzcRank); i++ {
		czzcRank[i]["rank"] = i + 1
		czzcRank[i]["czl"] = utils.CzlStr(czzcRank[i]["czl"].(float64))
	}

	rankList["level"] = levelRank
	rankList["sscz"] = ssczRank
	rankList["cz"] = czRank
	rankList["czzc"] = czzcRank
	rcache.SetRankList(&rankList)
	return rankList

}

func (ss *SysService) GetConsumptionList(start, end time.Time) []gin.H {
	userList := []gin.H{}
	sql := fmt.Sprintf(`select sum(yblog.yb) fee,player.nickname as nickname from yblog left join player on yblog.nickname=player.name
	where yblog.buytime >= %d 
	and yblog.buytime < %d
	group by player.id order by fee desc limit %d`, start.Unix(), end.Unix(), 50)
	conList := []struct {
		Fee      int
		Nickname string
	}{}
	ss.GetDb().Raw(sql).Scan(&conList)
	for i, con := range conList {
		color := ""
		if con.Fee >= 90000 {
			color = "red"
		} else if con.Fee >= 60000 {
			color = "blue"
		} else if con.Fee >= 30000 {
			color = "green"
		}
		userList = append(userList, gin.H{
			"rank":     i + 1,
			"nickname": con.Nickname,
			"color":    color,
		})
	}
	return userList
}

func (ss *SysService) GetUserConsumption(start, end time.Time, account string) int {
	con := struct {
		fee int
	}{}
	ss.GetDb().Raw("SELECT SUM(yb) as fee FROM yblog WHERE buytime >= ? AND buytime < ? AND nickname = ?", start.Unix(), end.Unix(), account).Scan(&con)
	return con.fee
}

func (ss *SysService) GetChatUserInfo(userId int) gin.H {
	user := ss.OptSvc.UserSrv.GetUserById(userId)
	userInfo := ss.OptSvc.UserSrv.GetUserInfoById(userId)
	img := ""
	if userInfo.NowAchievementTitle != "" {
		cardTitle := ss.OptSvc.UserSrv.GetCardTile(userInfo.NowAchievementTitle)
		if cardTitle != nil {
			img = cardTitle.Img
		}
	}
	isVip := false
	vipP := &models.UProp{}
	ss.OptSvc.GetDb().Where("pid in (?) and uid=?", vipProps, userId).First(vipP)
	if vipP.ID > 0 {
		isVip = true
	}

	info := gin.H{
		"id":         user.ID,
		"nickname":   user.Nickname,
		"is_merge":   userInfo.Merge > 0,
		"is_vip":     isVip,
		"img":        img,
		"group_uuid": rpc_group.GetGroupID(userId),
	}
	return info
}

// 消费排行榜信息：
// 返回： 消费活动是否开启了， 排行榜列表， 消费活动时间， 本人消费元宝数
func (ss *SysService) GetConsumptionInfo(userId int) (bool, []gin.H, string, int) {
	userList := []gin.H{}
	timeSetting := common.GetTimeConfig("consumption")
	startTime, err := utils.StrParseTime(timeSetting.StartTime + " 00:00:00")
	if err != nil {
		return false, userList, "", 0
	}
	endTime, err := utils.StrParseTime(timeSetting.EndTime + " 00:00:00")
	if err != nil {
		return false, userList, "", 0
	}
	now := time.Now()

	if !(startTime.Sub(now).Seconds() < 0 && endTime.Sub(now).Seconds() > 0) {
		return false, userList, "", 0
	}
	userList = ss.GetConsumptionList(startTime, endTime)
	timeSet := timeSetting.StartTime + " ~ " + timeSetting.StartTime
	user := ss.OptSvc.UserSrv.GetUserById(userId)
	userCon := ss.GetUserConsumption(startTime, endTime, user.Account)
	return true, userList, timeSet, userCon
}
