package scheduled

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

type Task struct {
	ID          cron.EntryID `json:"-"`
	Title       string
	Description string
	Scheduled   string
	JobName     string
	Job         func() `json:"-"`
}

type TaskList struct {
	list []*Task
	cron *cron.Cron
}

func (ts TaskList) GetDatas() {

}

func (ts TaskList) ForceRun(title string) {
	for _, t := range ts.list {
		if t.Title == title {
			t.Job()
			entry := ts.cron.Entry(t.ID)
			entry.Schedule.Next(time.Now())
		}
	}
}

func (ts TaskList) Set(cron *cron.Cron) {
	for _, t := range ts.list {
		t.Job = name2func[t.JobName]
		id, err := cron.AddFunc(t.Scheduled, t.Job)
		if err != nil {
			log.Printf("init tasks err:%s\n", err)
			continue
		}
		t.ID = id
		log.Printf("init tasks:%s success!\n", t.Title)
	}
}
