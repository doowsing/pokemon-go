package scheduled

import "github.com/robfig/cron/v3"

func Run() {
	// 新建一个定时任务对象
	// 根据cron表达式进行时间调度，cron可以精确到秒，大部分表达式格式也是从秒开始。
	//crontab := cron.New()  默认从分开始进行时间调度
	crontab := cron.New(cron.WithSeconds(), cron.WithChain(cron.Recover(cron.DefaultLogger)), cron.WithLogger(cron.DefaultLogger)) //精确到秒
	//定义定时器调用的任务函数
	//定时任务
	Tasks.Set(crontab)
	// 启动定时器
	crontab.Start()
	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	defer crontab.Stop()
	// 定时任务是另起协程执行的,这里使用 select 简答阻塞.实际开发中需要
	// 根据实际情况进行控制
	select {} //阻塞主线程停止
}
