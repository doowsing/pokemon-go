package services

import (
	"encoding/json"
	"fmt"
	"pokemon/pkg/models"
	"pokemon/pkg/rcache"
	"pokemon/pkg/repositories"
	"strconv"
	"strings"
)

type FightService struct {
	baseService
	repo   *repositories.FightRepositories
	rdrepo *rcache.FightRedisRepository
}

func NewFightService(osrc *OptService) *FightService {
	us := &FightService{repo: repositories.NewFightRepositories(), rdrepo: rcache.NewFightRedisRepository()}
	us.SetOptSrc(osrc)
	return us
}

func (fs *FightService) getZbAttrCache(petId int) map[string]float64 {
	speAttrStr, err := rcache.Get(PETATTR + strconv.Itoa(petId))
	if err == nil {
		var speAttr map[string]float64
		err = json.Unmarshal(speAttrStr, &speAttr)
		if err == nil {
			fmt.Printf("zbs attrs:%s\n", speAttr)
			return speAttr
		}
	}
	fmt.Printf("getzbattr error:%s\n", err)
	return nil
}

func (fs *FightService) setZbAttrCache(petId int, speAttrStr map[string]float64) bool {
	rcache.Set(PETATTR+strconv.Itoa(petId), speAttrStr)
	return true
}

func (fs *FightService) DelZbAttr(upid int) {
	rcache.Delete(PETATTR + strconv.Itoa(upid))
}

func (fs *FightService) GetZbAttr(upet *models.UPet, zbs []models.UProp) {
	if speAttr := fs.getZbAttrCache(upet.ID); speAttr != nil {
		upet.SpeAttr = speAttr
		upet.Ac = int(speAttr["ac"])
		upet.Mc = int(speAttr["mc"])
		upet.Hits = int(speAttr["hits"])
		upet.Miss = int(speAttr["miss"])
		upet.Speed = int(speAttr["speed"])
		upet.Hp = int(speAttr["hp"])
		upet.Mp = int(speAttr["mp"])
		fmt.Printf("zbs attrs:%s\n", upet.SpeAttr)
		return
	}
	if speAttrStr, err := rcache.Get(PETATTR + strconv.Itoa(upet.ID)); err == nil {
		var speAttr map[string]float64
		if err := json.Unmarshal(speAttrStr, speAttr); err == nil {
			upet.SpeAttr = speAttr
			upet.Ac = int(speAttr["ac"])
			upet.Mc = int(speAttr["mc"])
			upet.Hits = int(speAttr["hits"])
			upet.Miss = int(speAttr["miss"])
			upet.Speed = int(speAttr["speed"])
			upet.Hp = int(speAttr["hp"])
			upet.Mp = int(speAttr["mp"])
			fmt.Printf("zbs attrs:%s\n", upet.SpeAttr)
			return
		}
	}
	if zbs == nil {
		zbs = fs.OptSrc.PropSrv.GetPZbs(upet.ID)
	}
	zbAttr := fs.OptSrc.PropSrv.CountZbAttr(zbs)
	for t, v := range zbAttr {
		switch t {
		case "ac":
			upet.Ac = int(float64(upet.Ac) + v)
			break
		case "mc":
			upet.Mc = int(float64(upet.Mc) + v)
			break
		case "hp":
			upet.Hp = int(float64(upet.Hp) + v)
			break
		case "mp":
			upet.Mp = int(float64(upet.Mp) + v)
			break
		case "speed":
			upet.Speed = int(float64(upet.Speed) + v)
			break
		case "hits":
			upet.Hits = int(float64(upet.Hits) + v)
			break
		case "miss":
			upet.Miss = int(float64(upet.Miss) + v)
			break
		case "time", "crit", "dxsh", "sdmp", "hitshp", "hitsmp", "shjs", "szmp":
			upet.SpeAttr[t] += v
			break
		}
	}
	for t, v := range zbAttr {

		if !strings.Contains(t, "rate") {
			continue
		}
		switch strings.ReplaceAll(t, "rate", "") {
		case "ac":
			upet.Ac = int(float64(upet.Ac) * (1 + v))
			break
		case "mc":
			upet.Mc = int(float64(upet.Mc) * (1 + v))
			break
		case "hp":
			upet.Hp = int(float64(upet.Hp) * (1 + v))
			break
		case "mp":
			upet.Mp = int(float64(upet.Mp) * (1 + v))
			break
		case "speed":
			upet.Speed = int(float64(upet.Speed) * (1 + v))
			break
		case "hits":
			upet.Hits = int(float64(upet.Hits) * (1 + v))
			break
		case "miss":
			upet.Miss = int(float64(upet.Miss) * (1 + v))
			break
		}
	}

	// 称号
	userInfo := fs.OptSrc.UserSrv.GetUserInfoById(upet.Uid)
	if userInfo.NowAchievementTitle != "" {
		cardTitle := fs.OptSrc.UserSrv.GetCardTile(userInfo.NowAchievementTitle)
		if cardTitle != nil {
			upet.Ac += cardTitle.Ac
			upet.Mc += cardTitle.Mc
			upet.Hits += cardTitle.Hits
			upet.Miss += cardTitle.Miss
			upet.Speed += cardTitle.Speed
			upet.Hp += cardTitle.Hp
			upet.Mp += cardTitle.Mp
			upet.Ac = int(float64(upet.Ac) * (1 + float64(cardTitle.AcRate)*0.01))
			upet.Mc = int(float64(upet.Mc) * (1 + float64(cardTitle.McRate)*0.01))
			upet.Hits = int(float64(upet.Hits) * (1 + float64(cardTitle.HitsRate)*0.01))
			upet.Miss = int(float64(upet.Miss) * (1 + float64(cardTitle.MissRate)*0.01))
			upet.Speed = int(float64(upet.Speed) * (1 + float64(cardTitle.SpeedRate)*0.01))
			upet.Hp = int(float64(upet.Hp) * (1 + float64(cardTitle.HpRate)*0.01))
			upet.Mp = int(float64(upet.Mp) * (1 + float64(cardTitle.MpRate)*0.01))
			upet.SpeAttr["time"] += float64(cardTitle.Time)
			upet.SpeAttr["money"] += float64(cardTitle.AddMoney)
			upet.SpeAttr["dxsh"] += float64(cardTitle.Dxsh)
			upet.SpeAttr["sdmp"] += float64(cardTitle.Sdmp)
			upet.SpeAttr["hitshp"] += float64(cardTitle.HitsHp)
			upet.SpeAttr["hitsmp"] += float64(cardTitle.HitsMp)
			upet.SpeAttr["shjs"] += float64(cardTitle.Shjs)
			upet.SpeAttr["szmp"] += float64(cardTitle.Szmp)
		}
	}
	upet.SpeAttr["ac"] = float64(upet.Ac)
	upet.SpeAttr["mc"] = float64(upet.Mc)
	upet.SpeAttr["hits"] = float64(upet.Hits)
	upet.SpeAttr["miss"] = float64(upet.Miss)
	upet.SpeAttr["speed"] = float64(upet.Speed)
	upet.SpeAttr["hp"] = float64(upet.Hp)
	upet.SpeAttr["mp"] = float64(upet.Mp)
	fmt.Printf("zbs attrs:%s\n", upet.SpeAttr)
	fs.setZbAttrCache(upet.ID, upet.SpeAttr)
}

func (fs *FightService) OpenMap(userId, mapId int) bool {
	props := fs.OptSrc.PropSrv.GetCarryProps(userId, false)
	for _, p := range props {
		p.GetM()
		if p.MModel.VaryName == 13 && p.MModel.Effect == "openmap:"+strconv.Itoa(mapId) {
			fs.OptSrc.PropSrv.DecrPropById(p.ID, 1)
			return true
		}
	}
	return false
}

func (fs *FightService) GetFbRecord(userId, mapId int) *models.RecordFb {
	record := &models.RecordFb{}
	if fs.GetDb().Where("uid=? and inmap=?", userId, mapId).Find(record).RowsAffected > 0 {
		return record
	}
	return nil
}
