package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	ginapp "pokemon/pkg/ginapp"
)

func ExtOnline(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	session := gapp.Session()
	session.Set("id", session.MustGet("id"))
	_ = session.Save()
	gapp.OptSrv.UserSrv.UpdateIdToken(gapp.Id())
	gapp.String("<!--consumption2exp-->")
}

func McGate(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	opType := c.DefaultQuery("op", "none")
	if opType == "none" {
		return
	}
	petId := com.StrTo(c.DefaultQuery("id", "-1")).MustInt()
	if petId <= 0 {
		gapp.String("操作失败！参数错误！")
		return
	}
	pet := gapp.OptSrv.PetSrv.GetPet(gapp.Id(), petId)
	if pet == nil {
		gapp.String("操作失败！参数错误！")
		return
	}

	id := c.MustGet("id").(int)
	user := gapp.OptSrv.UserSrv.GetUserById(id)
	if opType == "z" || opType == "change" {
		// z：设置主宠，change：设置主宠，特殊，暂不管
		if petId == user.Mbid {
			gapp.String("已经是主战！")
			return
		}
		if pet.Muchang != 0 {
			gapp.String("在牧场的宝宝不能设为主战哦！")
			return
		}

		if gapp.OptSrv.UserSrv.SetMBid(id, petId) {
			gapp.String("更改主战宝宝成功!")
		} else {
			gapp.String("操作失败！参数错误！")
		}

	} else if opType == "g" {
		// 携带宠物
		if pet.ChChengSx != "" {
			gapp.String("传承中不能取出")
		} else if pet.Muchang == 0 {
			gapp.String("此宠物已经携带！")
		} else if gapp.OptSrv.PetSrv.GetCarryPetCnt(id) >= 3 {
			gapp.String("您最多同时只能携带3个宝宝！")
			//} else if gapp.OptSrv.PetSrv.PutOut(petId) {
			//	gapp.String("操作成功!")
		} else {
			gapp.String("操作失败！参数错误！")
		}
	} else if opType == "d" {
		// 丢弃宠物
		if user.Money < 10000 {

		}
	} else if opType == "s" {
		// 寄养宠物
		if pet.Muchang == 1 {
			gapp.String("已经在牧场！")
		} else if petId == user.Mbid {
			gapp.String("该宠物为主战宠物，无法寄养！")
		} else if gapp.OptSrv.PetSrv.GetMcPetCnt(id) >= user.McPlace {
			gapp.String("您的宠物寄养空间已满，不能寄养更多宝宝！")
			//} else if gapp.OptSrv.PetSrv.PutIn(petId) {
			//	gapp.String("操作成功!")
		} else {
			gapp.String("操作失败！参数错误！")
		}
	} else {
		gapp.String("操作失败！参数错误！")
	}

}

func GetBag(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id := c.MustGet("id").(int)
	user := gapp.OptSrv.UserSrv.GetUserById(id)
	clean := c.DefaultQuery("clean", "0") == "1"
	uprops := gapp.OptSrv.PropSrv.GetCarryProps(id, clean)

	//gapp.String("nothing")
	gapp.HTML("compotent/get_bag.jet.html", gin.H{"uprops": uprops, "user": user})
}

func ShowPet(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	petId := c.DefaultQuery("id", "-1")
	pet := gapp.OptSrv.PetSrv.GetPetById(com.StrTo(petId).MustInt())
	if pet == nil {
		gapp.String("宠物id不存在！")
	}
	pet.GetM()
	gapp.HTML("compotent/bbshow.jet.html", gin.H{
		"pet": pet,
	})
}

func UsePropString(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id := c.MustGet("id").(int)
	upid := c.DefaultQuery("id", "-1")
	mfs := c.DefaultQuery("js", "")
	if mfs != "" {
		// 魔法石
		pid := c.DefaultQuery("pid", "")
		if pid != "" {
			gapp.String(gapp.OptSrv.PropSrv.UsePropByPid(id, com.StrTo(pid).MustInt()))
		}
		gapp.String("参数错误！")
	} else {
		// 普通道具
		gapp.String(gapp.OptSrv.PropSrv.UsePropById(id, com.StrTo(upid).MustInt()))
	}

}

func GetPropInfo(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	id := c.MustGet("id").(int)
	pid := com.StrTo(c.DefaultQuery("id", "-1")).MustInt()
	if pid == 0 {
		return
	}
	showType := com.StrTo(c.DefaultQuery("type", "1")).MustInt()
	var result string
	if showType == 1 {
		// 展示道具种类信息
		result = gapp.OptSrv.PropSrv.GeneratePropInfo(id, 0, pid)
	} else {
		// 展示用户道具信息
		result = gapp.OptSrv.PropSrv.GeneratePropInfo(id, pid, 0)
	}
	gapp.String(result)
}

func OffZb(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	id := c.MustGet("id").(int)
	petid := com.StrTo(c.DefaultQuery("bid", "0")).MustInt()
	if petid <= 0 {
		gapp.String("3")
		return
	}
	propid := com.StrTo(c.DefaultQuery("id", "0")).MustInt()
	if propid <= 0 {
		gapp.String("1")
		return
	}
	bagCnt := gapp.OptSrv.PropSrv.GetCarryPropsCnt(id)
	user := gapp.OptSrv.UserSrv.GetUserById(id)
	if bagCnt >= user.BagPlace {
		gapp.String("5")
		return
	}
	if gapp.OptSrv.PropSrv.OffZbBypid(id, petid, propid) {
		gapp.String("2")
	} else {
		gapp.String("4")
	}
}

func GetGameLogListLen(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.String("gamelog 队列长度为：%d", gapp.OptSrv.SysSrv.GetGameLogListLen())
}

func GetLoginLogListLen(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)
	gapp.String("loginlog 队列长度为：%d", gapp.OptSrv.SysSrv.GetLoginLogListLen())
}
