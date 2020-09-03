package models

type Welcome struct {
	Id      int    `gorm:"column:id;primary_key"`
	Code    string `gorm:"column:code"`
	Text    string `gorm:"column:value2"`
	Content string `gorm:"column:contents"`
}

func (W *Welcome) TableName() string {
	return "welcome"
}

func (w *Welcome) IsValid() bool {
	return w.Code != ""
}

func (w *Welcome) AfterFind() (err error) {
	//w.Text = utils.ToUtf8(w.Text)
	//w.Content = utils.ToUtf8(w.Content)

	return
}

type TimeConfig struct {
	Id        int    `gorm:"column:Id;primary_key"`
	Title     string `gorm:"column:titles"`
	Day       string `gorm:"column:days"`
	StartTime string `gorm:"column:starttime"`
	EndTime   string `gorm:"column:endtime"`
}

func (this *TimeConfig) TableName() string {
	return "timeconfig"
}

type SaoLeiAwardInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Img  string `json:"img"`
}
