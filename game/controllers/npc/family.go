package npc

import (
	"github.com/gin-gonic/gin"
	"pokemon/game/ginapp"
	"strconv"
	"unicode/utf8"
)

func FamilyPageInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	allFamily := gapp.OptSvc.NpcSrv.GetAllFamilyData()
	data := gin.H{
		"all_family":         allFamily,
		"my_family":          nil,
		"family_store_props": []gin.H{},
		"honor":              0,
		"contribution":       0,
	}
	member := gapp.OptSvc.NpcSrv.GetFamilyMember(gapp.Id())
	if member != nil && member.Authority > 0 {
		myFamily := gapp.OptSvc.NpcSrv.GetFamilyInfo(member.FamilyId, gapp.Id())
		family := gapp.OptSvc.NpcSrv.GetFamily(member.FamilyId)
		store := gapp.OptSvc.NpcSrv.GetFamilyStoreData(gapp.Id(), family.ShopLevel)
		data["my_family"] = myFamily
		data["family_store_props"] = store
		data["honor"] = member.Honor
		data["contribution"] = member.Contribution
	}
	gapp.JSONDATAOK("", data)
}

func FamilyShowInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	gapp.JSONDATAOK("", gapp.OptSvc.NpcSrv.GetFamilyInfo(id, gapp.Id()))
}

func CreateFamily(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	name := c.PostForm("title")
	content := c.PostForm("content")
	if l := utf8.RuneCountInString(name); l < 2 || l > 8 {
		gapp.JSONDATAOK("名称长度不符合要求！", gin.H{"result": false})
		return
	}
	if l := utf8.RuneCountInString(content); l < 1 || l > 200 {
		gapp.JSONDATAOK("介绍长度不符合要求！", gin.H{"result": false})
		return
	}
	ok, msg := gapp.OptSvc.NpcSrv.CreateFamily(gapp.Id(), name, content)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func ApplyFamily(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	ok, msg := gapp.OptSvc.NpcSrv.ApplyFamily(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func FireFamilyMember(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	ok, msg := gapp.OptSvc.NpcSrv.FireFamilyMember(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func ReplyApplyFamily(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	pass, err := strconv.Atoi(c.Query("pass"))
	if err != nil && pass < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	ok, msg := gapp.OptSvc.NpcSrv.ReplyApplyFamily(gapp.Id(), id, pass == 1)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func ManageFamilyAuthority(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	new_authority, err := strconv.Atoi(c.Query("new_authority"))
	if err != nil && new_authority < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	ok, msg := gapp.OptSvc.NpcSrv.ManageAuthority(gapp.Id(), id, new_authority)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func DonateFamily(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	num, err := strconv.Atoi(c.Query("num"))
	if err != nil && num < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	ok, msg := gapp.OptSvc.NpcSrv.FamilyDonate(gapp.Id(), id, num)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func PurchaseFamilyStore(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil && id < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}
	num, err := strconv.Atoi(c.Query("num"))
	if err != nil && num < 1 {
		gapp.JSONDATAOK("请检查参数", nil)
		return
	}

	ok, msg := gapp.OptSvc.NpcSrv.PurchaseFamilyStoreProp(gapp.Id(), id, num)
	honor := 0
	contribution := 0
	member := gapp.OptSvc.NpcSrv.GetFamilyMember(gapp.Id())
	if member != nil && member.Authority != 0 {
		honor = member.Honor
		contribution = member.Contribution
	}
	gapp.JSONDATAOK(msg, gin.H{"result": ok, "honor": honor, "contribution": contribution})
}

func DisbandFamily(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	ok, msg := gapp.OptSvc.NpcSrv.DisbandFamily(gapp.Id())
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func UpgradeFamily(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	ok, msg := gapp.OptSvc.NpcSrv.UpgradeFamily(gapp.Id())
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func FamilyUpgradeInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	info, msg := gapp.OptSvc.NpcSrv.GetFamilyUpgradeInfo(gapp.Id())
	gapp.JSONDATAOK(msg, info)
}

func ReceiveFamilyWelFare(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	ok, msg := gapp.OptSvc.NpcSrv.GetFamilyWelfare(gapp.Id())
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func ExitFamily(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	ok, msg := gapp.OptSvc.NpcSrv.ExitFamily(gapp.Id())
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

func UpgradeFamilyStore(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	ok, msg := gapp.OptSvc.NpcSrv.UpgradeFamilyStore(gapp.Id())
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}
