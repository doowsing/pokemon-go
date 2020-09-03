package npc

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"pokemon/common/persistence"
	"pokemon/game/ginapp"
	"pokemon/game/models"
	"strings"
)

// 占卜屋道具列表
func ZhanBuWuList(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	ps := []*models.MProp{}
	persistence.GetOrm().Where("varyname=22").Order("stime").Find(&ps)
	datas := []gin.H{}
	for _, p := range ps {
		data := gin.H{"id": p.ID, "name": strings.ReplaceAll(p.Name, "石", ""), "img": fmt.Sprintf("%d.gif", p.ID)}
		datas = append(datas, data)
	}
	gapp.JSONDATAOK("", gin.H{"cards": datas})
}

// 占卜屋道具使用
func ZhanBuWuUse(c *gin.Context) {
	gapp := ginapp.NewGapp(c)
	defer ginapp.DropGapp(gapp)

	id := com.StrTo(c.Query("id")).MustInt()
	if id == 0 {
		gapp.JSONDATAOK("参数出错！", gin.H{"result": false})
		return
	}
	msg := gapp.OptSvc.PropSrv.UsePropByPid(gapp.Id(), id)
	gapp.JSONDATAOK(msg, gin.H{"result": true})

}
