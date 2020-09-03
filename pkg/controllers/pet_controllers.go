package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"pokemon/pkg/common"
	"pokemon/pkg/ginapp"
	"pokemon/pkg/models"
	"pokemon/pkg/services"
	"strconv"
)

type PetController struct {
	service *services.PetService
}

func NewPetController() *PetController {
	return &PetController{service: services.NewPetService(nil)}
}

func (ps *PetController) InitTable(c *gin.Context) {
	ps.service.InitTable()
	c.JSON(200, gin.H{"code": 200, "msg": "", "data": "成功创建表！"})
}

func (ps *PetController) ShowALLMPet(c *gin.Context) {
	mPets := ps.service.GetAllMPet()

	c.JSON(200, gin.H{"code": 200, "msg": "获取所有宠物原型", "data": mPets})
}

func (ps *PetController) ShowMSkill(c *gin.Context) {
	skillId := c.Query("id")
	_skillId, err := strconv.Atoi(skillId)
	if err != nil {
		c.JSON(200, gin.H{"code": 200, "msg": "获取技能原型出错", "data": err})
	} else {
		mskill := ps.service.GetMskill(_skillId)

		c.JSON(200, gin.H{"code": 200, "msg": "获取技能原型成功", "data": mskill})
	}

}

func (ps *PetController) ShowMPet(c *gin.Context) {
	petId := c.Query("id")
	_petId, err := strconv.Atoi(petId)
	if err != nil {
		c.JSON(200, gin.H{"code": 200, "msg": "参数出错", "data": err})
	} else {
		mpet := ps.service.GetMpet(_petId)

		c.JSON(200, gin.H{"code": 200, "msg": "获取宠物原型成功", "data": mpet})
	}

}

//
//func (ps *PetController) CreatePet(c *gin.Context) {
//	if uid, exists := c.Get("id"); !exists {
//		log.Print("玩家ID为", c.Keys)
//		c.JSON(http.StatusForbidden, models.Response{Err: common.NewErr(403, "请登录后再操作！"), Data: nil})
//	} else {
//		pidStr := c.Query("id")
//		pid, _ := strconv.Atoi(pidStr)
//		if ps.service.CreatPetById(uid.(int), pid) {
//			c.JSON(http.StatusOK, models.Response{Err: common.NewErr(200, "新增宠物成功！"), Data: nil})
//		} else {
//			c.JSON(http.StatusOK, models.Response{Err: common.NewErr(200, "新增宠物失败！"), Data: nil})
//		}
//	}
//}

func PetInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	pet := gapp.OptSrv.PetSrv.GetPetById(id)
	if pet == nil {
		gapp.JSONDATAOK("找不到该宠物", nil)
		return
	}
	pet.GetM()
	petMap := gin.H{
		"id":    pet.ID,
		"name":  pet.MModel.Name,
		"level": pet.Level,
		"wx":    pet.WxName(),
		"hp":    pet.Hp,
		"mp":    pet.Mp,
		"ac":    pet.Ac,
		"mc":    pet.Mc,
		"hits":  pet.Hits,
		"miss":  pet.Miss,
		"speed": pet.Speed,
		"czl":   pet.Czl,
		"img":   pet.MModel.ImgEffect,
	}
	gapp.JSONDATAOK("", petMap)
}

func CarryPetList(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	uid := gapp.Id()
	pets := gapp.OptSrv.PetSrv.GetCarryPets(uid)
	petsMap := []gin.H{}
	for _, pet := range pets {
		pet.GetM()
		pinfo := gin.H{
			"id":    pet.ID,
			"name":  pet.MModel.Name,
			"img":   pet.MModel.ImgCard,
			"level": pet.Level,
			"wx":    pet.WxName(),
		}
		petsMap = append(petsMap, pinfo)
	}
	user := gapp.OptSrv.UserSrv.GetUserById(uid)
	gapp.JSONDATAOK("", gin.H{"pets": petsMap, "main_id": user.Mbid})
}

func AllPetList(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	uid := gapp.Id()
	pets := gapp.OptSrv.PetSrv.GetAllPets(uid)
	petsMap := []gin.H{}
	for _, pet := range pets {
		pinfo := gin.H{
			"id":       pet.ID,
			"name":     pet.MModel.Name,
			"img":      pet.MModel.ImgCard,
			"level":    pet.Level,
			"wx":       pet.WxName(),
			"position": pet.Muchang,
		}
		petsMap = append(petsMap, pinfo)
	}
	user := gapp.OptSrv.UserSrv.GetUserById(uid)
	gapp.JSONDATAOK("", gin.H{"pets": petsMap, "main_id": user.Mbid})
}

// 放入牧场
func PutInPet(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	result, msg := gapp.OptSrv.PetSrv.PutIn(gapp.Id(), id)
	gapp.JSONDATAOK(msg, result)
}

// 携带宠物
func PutOutPet(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	result, msg := gapp.OptSrv.PetSrv.PutOut(gapp.Id(), id, c.Query("passwd"))
	gapp.JSONDATAOK(msg, result)
}

// 放生宠物
func ThrowPet(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	result, msg := gapp.OptSrv.PetSrv.Throw(gapp.Id(), id, c.Query("passwd"))
	gapp.JSONDATAOK(msg, result)
}

// 设置主宠
func SetMbid(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	pet := gapp.OptSrv.PetSrv.GetPet(gapp.Id(), id)
	if pet == nil {
		gapp.JSONDATAOK("宠物不存在！", gin.H{"result": false})
		return
	}
	if pet.Muchang != 0 {
		gapp.JSONDATAOK("宠物必须携带后才可设置为主宠！", gin.H{"result": false})
		return
	}
	if gapp.OptSrv.UserSrv.SetMBid(gapp.Id(), id) {
		gapp.JSONDATAOK("设置成功！！", gin.H{"result": true})
		return
	} else {
		gapp.JSONDATAOK("设置失败！！", gin.H{"result": false})
		return
	}
}

// 学习技能
func LearnSkill(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	ok, msg := gapp.OptSrv.PetSrv.LearnSkill(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 升级技能
func UpdateSkill(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	ok, msg := gapp.OptSrv.PetSrv.UpdateSkill(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 公屏展示宠物
func ShowPet2All(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	user := gapp.OptSrv.UserSrv.GetUserById(gapp.Id())
	userInfo := gapp.OptSrv.UserSrv.GetUserInfoById(gapp.Id())
	if userInfo.Bbshow > 0 {

		gapp.OptSrv.SysSrv.AnnoucePet2All(user.Nickname, "showpet||"+idStr)
		gapp.OptSrv.UserSrv.GetDb().Model(userInfo).Update(gin.H{"bbshow": gorm.Expr("bbshow-1")})
		gapp.JSONDATAOK("", gin.H{"result": true})
	} else {
		gapp.JSONDATAOK("展示次数不足！", gin.H{"result": false})
	}
}

func (ps *PetController) ShowPets(c *gin.Context) {

	uid, _ := c.Get("id")
	pets := ps.service.GetCarryPets(uid.(int))
	c.JSON(http.StatusOK, models.Response{Err: common.NewErr(200, ""), Data: pets})
}

func (ps *PetController) ShowExpByLv(c *gin.Context) {
	levelStr := c.Query("level")
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{Err: common.NewErr(200, "参数错误"), Data: levelStr})
	} else {
		exp, ok := ps.service.GetExpByLv(level)
		if ok != true {
			c.JSON(http.StatusOK, models.Response{Err: common.NewErr(200, "提取经验出错"), Data: err})
		} else {
			c.JSON(http.StatusOK, models.Response{Err: common.NewErr(200, "经验获取成功"), Data: exp})
		}
	}
}

func (ps *PetController) IncreatExp(c *gin.Context) {
	exp := c.Query("exp")
	userid, _ := c.Get("id")
	expNum, err := strconv.Atoi(exp)
	if err != nil {
		c.JSON(200, models.Response{common.NewErr(200, "参数错误！"), nil})
	} else {
		c.JSON(200, ps.service.IncreaceExp2MainPet(userid.(int), expNum))
	}
}
