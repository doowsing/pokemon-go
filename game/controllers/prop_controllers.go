package controllers

import (
	"github.com/gin-gonic/gin"
	"pokemon/common/rcache"
	"pokemon/game/ginapp"
	"pokemon/game/services"
	"pokemon/game/utils"
	"strconv"
)

type PropController struct {
	service *services.PropService
}

func NewPropController() *PropController {
	return &PropController{service: services.NewPropService(nil)}
}

func ShowProp2All(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数出错", gin.H{"result": false})
		return
	}
	prop := gapp.OptSvc.PropSrv.GetProp(gapp.Id(), id, false)
	if prop == nil {
		gapp.JSONDATAOK("道具不存在", gin.H{"result": false})
		return
	}
	now := utils.NowUnix()
	if rcache.PropShowTimer.InCoolTime(gapp.Id(), now) {
		gapp.JSONDATAOK("操作过快！", gin.H{"result": false})
		return
	}
	rcache.PropShowTimer.Set(gapp.Id(), now)
	services.AnnounceProp2All(gapp.Id(), prop)
	gapp.JSONDATAOK("", gin.H{"result": true})
}

func PropInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数出错", nil)
		return
	}
	propMap := gapp.OptSvc.PropSrv.GetPropInfoJson(gapp.Id(), id, 0)
	if propMap == nil {
		gapp.JSONDATAOK("没有这个道具！", nil)
		return
	}
	gapp.JSONDATAOK("", propMap)
}

func MPropInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil || pid < 1 {
		gapp.JSONDATAOK("参数出错", nil)
		return
	}
	propMap := gapp.OptSvc.PropSrv.GetPropInfoJson(0, 0, pid)
	if propMap == nil {
		gapp.JSONDATAOK("没有这个道具！", nil)
		return
	}
	gapp.JSONDATAOK("", propMap)
}

func CarryPropList(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	cleanStr := c.PostForm("clean")
	props := gapp.OptSvc.PropSrv.GetCarryProps(gapp.Id(), cleanStr == "true")
	propList := []gin.H{}
	for _, p := range props {
		propList = append(propList, gin.H{
			"id":        p.ID,
			"name":      p.MModel.Name,
			"vary_id":   p.MModel.VaryName,
			"vary_name": p.MModel.VaryNameStr,
			"price":     p.MModel.SellJb,
			"sum":       p.Sums,
		})
	}
	user := gapp.OptSvc.UserSrv.GetUserById(gapp.Id())
	gapp.JSONDATAOK("", gin.H{"props": propList, "bb_max": user.BagPlace})
}

func UseProp(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		gapp.JSONDATAOK("参数出错", nil)
		return
	}
	result := gapp.OptSvc.PropSrv.UsePropById(gapp.Id(), id)
	prop := gapp.OptSvc.PropSrv.GetProp(gapp.Id(), id, true)
	gapp.JSONDATAOK(result, gin.H{"result": true, "sum": prop.Sums, "id": prop.ID})

}

// 放入仓库
func PutInProp(c *gin.Context) {

	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	numStr := c.Query("num")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	var num int
	if numStr == "" {
		num = -1
	} else {
		num, err = strconv.Atoi(numStr)
		if err != nil || num < 0 {
			gapp.JSONDATAOK("请输入正确的数字！", gin.H{"result": false})
			return
		}
	}

	result, msg := gapp.OptSvc.PropSrv.PutIn(gapp.Id(), id, num)
	gapp.JSONDATAOK(msg, result)
}

// 取出到背包
func PutOutProp(c *gin.Context) {

	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	numStr := c.Query("num")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	var num int
	if numStr == "" {
		num = -1
	} else {
		num, err = strconv.Atoi(numStr)
		if err != nil || num < 0 {
			gapp.JSONDATAOK("请输入正确的数字！", gin.H{"result": false})
			return
		}
	}
	result, msg := gapp.OptSvc.PropSrv.PutOut(gapp.Id(), id, num, c.Query("passwd"))
	gapp.JSONDATAOK(msg, result)
}

// 丢弃道具
func ThrowProp(c *gin.Context) {

	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	if gapp.OptSvc.PropSrv.Throw(gapp.Id(), id) {
		gapp.JSONDATAOK("丢弃成功！", gin.H{"result": true, "id": id})
		return
	} else {
		gapp.JSONDATAOK("丢弃失败，背包没有该道具！", gin.H{"result": false})
		return
	}
}

// 拍卖道具
func Auction(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	idStr := c.PostForm("id")
	numStr := c.PostForm("num")
	priceStr := c.PostForm("price")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	num, err := strconv.Atoi(numStr)
	if num < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	price, err := strconv.Atoi(priceStr)
	if price < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	currencyStr := c.PostForm("currency")
	nickname := c.PostForm("nickname")
	currency := "jb"
	switch currencyStr {
	case "jinbi":
		currency = "jb"
		break
	case "shuijing":
		currency = "sj"
		break
	case "yuanbao":
		currency = "yb"
		break
	default:
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	ok, msg := gapp.OptSvc.PropSrv.Auction(gapp.Id(), id, num, price, currency, nickname)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 续拍道具
func ReAuction(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	ok, msg := gapp.OptSvc.PropSrv.ReAuction(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 取回道具
func RollAuction(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	ok, msg := gapp.OptSvc.PropSrv.RollAuction(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 交易所购买
func Purchase(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	numStr := c.Query("num")
	num, err := strconv.Atoi(numStr)
	if num < 1 || err != nil {
		gapp.JSONDATAOK("请检查输入数量！", nil)
		return
	}
	ok, msg := gapp.OptSvc.PropSrv.Purchase(gapp.Id(), id, num)
	gapp.JSONDATAOK(msg, gin.H{"result": ok})
}

// 商店出售，只出售为金币
func ShopSell(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if id < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", nil)
		return
	}
	numStr := c.Query("num")
	num, err := strconv.Atoi(numStr)
	if num < 1 || err != nil {
		gapp.JSONDATAOK("请检查输入数量！", nil)
		return
	}
	ok, msg := gapp.OptSvc.PropSrv.ShopSell(gapp.Id(), id, num)
	user := gapp.OptSvc.UserSrv.GetUserById(gapp.Id())
	gapp.JSONDATAOK(msg, gin.H{"result": ok, "jb": user.Money})
}

// 商店购买
func ShopPurchase(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	pidStr := c.Query("pid")
	pid, err := strconv.Atoi(pidStr)
	if pid < 1 || err != nil {
		gapp.JSONDATAOK("参数出错！", gin.H{"result": false})
		return
	}
	numStr := c.Query("num")
	num, err := strconv.Atoi(numStr)
	if num < 1 || err != nil {
		gapp.JSONDATAOK("请检查输入数量！", gin.H{"result": false})
		return
	}
	currency := c.Query("currency")
	ok, msg := gapp.OptSvc.PropSrv.ShopPurchase(gapp.Id(), pid, num, currency)
	if ok {
		user := gapp.OptSvc.UserSrv.GetUserById(gapp.Id())
		userInfo := gapp.OptSvc.UserSrv.GetUserInfoById(gapp.Id())
		gapp.JSONDATAOK(msg, gin.H{
			"result": ok,
			"jb":     user.Money,
			"sj":     userInfo.Sj,
			"yb":     user.Yb,
			"ww":     user.Prestige,
		})
	} else {
		gapp.JSONDATAOK(msg, gin.H{"result": ok})
	}

}
