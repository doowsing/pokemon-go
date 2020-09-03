package datastore

import (
	"pokemon/pkg/models"
)

var Mmap = &mmapStore{newStore()}

type mmapStore struct {
	*store
}

func (s *mmapStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	var all []models.Map
	db.Find(&all)
	for _, v := range all {
		if len(v.Levels) > 1 {
			for _, gpc := range tmpGpcs {
				if gpc.Boss < 4 && (gpc.Level >= v.Levels[0] || gpc.Level <= v.Levels[1]) {
					gpc1 := gpc
					v.Gpcs = append(v.Gpcs, &gpc1)
				}
			}
		}
		v1 := v
		s.store.data.Store(v.ID, &v1)
	}
}

func (s *mmapStore) Get(id int) *models.Map {
	if data, exist := s.store.data.Load(id); exist {
		return data.(*models.Map)
	}
	return nil
}

var Mgpc = &mgpcStore{newStore()}
var tmpGpcs []models.Gpc

type mgpcStore struct {
	*store
}

func (s *mgpcStore) Update(newTime int) {
	if !s.checkTime(newTime) {
		return
	}
	defer s.upTime()
	var all []models.Gpc
	tmpGpcs = all
	db.Find(&all)
	for _, v := range all {
		v1 := v
		s.store.data.Store(v.ID, &v1)
	}
}

func (s *mgpcStore) Get(id int) *models.Gpc {
	if data, exist := s.store.data.Load(id); exist {
		return data.(*models.Gpc)
	}
	return nil
}

func InitMap(newTime int) {
	Mgpc.Update(newTime)
	Mmap.Update(newTime)
}

type Fuben struct {
	ID    int
	Name  string
	Time  int
	Level int
	GwIds []int
}

var fubenStore = map[int]*Fuben{
	11: {
		ID:    11,
		Name:  "辉煌的大道",
		Time:  86400,
		Level: 30,
		GwIds: []int{1417, 1418, 1419, 1420, 1421, 1422, 1423, 1424, 1425, 1426, 1427, 1428},
	},
	151: {
		ID:    151,
		Name:  "伊苏王",
		Time:  36000,
		Level: 30,
		GwIds: []int{1434, 1435, 1436, 1437, 1438, 1439, 1440, 1441, 1442, 1443, 1444, 1445, 1446, 1447, 1448, 1449, 1450, 1451, 1452, 1453, 1454, 1455},
	},
	12: {
		ID:    12,
		Name:  "火龙王的宫殿",
		Time:  36000,
		Level: 50,
		GwIds: []int{158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 184, 185, 186, 187},
	},
	13: {
		ID:    13,
		Name:  "史芬克斯密穴",
		Time:  43200,
		Level: 70,
		GwIds: []int{188, 189, 190, 191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209, 210, 211, 212, 213, 214},
	},
	14: {
		ID:    14,
		Name:  "玲珑城",
		Time:  54000,
		Level: 85,
		GwIds: []int{263, 264, 265, 266, 267, 268, 269, 270, 271, 272, 273, 274, 275, 276, 277, 278, 279, 280, 281, 282, 283, 284, 285, 286, 287, 289, 290, 291, 292},
	},
	50: {
		ID:    50,
		Name:  "厄非斯深渊",
		Time:  54000,
		Level: 90,
		GwIds: []int{429, 430, 431, 432, 433, 434, 435, 436, 437, 438, 439, 440, 441, 442, 443, 444, 445, 446, 447, 448, 449, 450, 451, 452, 453, 454, 455},
	},
	124: {
		ID:    124,
		Name:  "阿尔提密林",
		Time:  60000,
		Level: 1,
		GwIds: []int{505, 506, 507, 508, 509, 510, 511, 512, 513},
	},
	127: {
		ID:    127,
		Name:  "菲拉苛地域",
		Time:  72000,
		Level: 1,
		GwIds: []int{774, 775, 776, 777, 778, 779, 780, 781, 782, 783, 784, 785, 786, 787, 789, 790},
	},
	143: {
		ID:    143,
		Name:  "熔岩地宫",
		Time:  72000,
		Level: 1,
		GwIds: []int{1145, 1146, 1147, 1148, 1149, 1150, 1151, 1152, 1153, 1154, 1155, 1156, 1157, 1158, 1159},
	},
	144: {
		ID:    144,
		Name:  "幻魔之境",
		Time:  72000,
		Level: 1,
		GwIds: []int{1160, 1161, 1162, 1163, 1164, 1165, 1166, 1167, 1168, 1169, 1170, 1171, 1172, 1173, 1174, 1175, 1176, 1177, 1178, 1179, 1180, 1181, 1182, 1183, 1184},
	},
}

func GetFbSetting(mapId int) *Fuben {
	if fuben, ok := fubenStore[mapId]; ok {
		return fuben
	}
	return nil
}
