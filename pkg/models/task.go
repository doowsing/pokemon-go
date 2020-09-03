package models

var getTask = DefaultGetTask

type Task struct {
	Id      int    `gorm:"column:id;primary_key"`
	Title   string `gorm:"column:title"`
	FromNpc string `gorm:"column:fromnpc"`
	FromMsg string `gorm:"column:frommsg"`
	OkMsg   string `gorm:"column:okmsg"`
	OkNpc   int    `gorm:"column:oknpc"`
	OkNeed  string `gorm:"column:okneed"`
	Result  string `gorm:"column:result"`
	Cid     string `gorm:"column:cid"`
	LimitLv string `gorm:"column:limitlv"`
	Hide    int    `gorm:"column:hide;default:1"`
	Xulie   int    `gorm:"column:xulie;default:0"`
	Flags   int    `gorm:"column:flags;default:0"`
	Color   int    `gorm:"column:color;default:1"`
}

func (u *Task) TableName() string {
	return "task"
}

type UTask struct {
	Id      int    `gorm:"column:id;primary_key"`
	Uid     int    `gorm:"column:uid"`
	TaskId  int    `gorm:"column:taskid"`
	State   string `gorm:"column:state;default:'0'"`
	ComSelf string `gorm:"column:comself;default:'0'"`
	Time    int    `gorm:"column:time;default:0"`
}

func (u *UTask) TableName() string {
	return "task_accept"
}

type TaskLog struct {
	Id      int `gorm:"column:Id;primary_key"`
	Uid     int `gorm:"column:uid;default:0"`
	TaskId  int `gorm:"column:taskid;default:0"`
	Xulie   int `gorm:"column:xulie;default:0"`
	Time    int `gorm:"column:time;default:0"`
	FromNpc int `gorm:"column:fromnpc;default:0"`

	MModel *Task `gorm:"-"`
}

func (u *TaskLog) TableName() string {
	return "tasklog"
}

func (u *TaskLog) GetM() *Task {
	if u.MModel == nil {
		u.MModel = getTask(u.TaskId)
	}
	return u.MModel
}

func DefaultGetTask(id int) *Task {
	return nil
}

func SetTaskFunc(f func(id int) *Task) {
	getTask = f
}
