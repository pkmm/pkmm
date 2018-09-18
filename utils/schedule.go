package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/toolbox"
	"github.com/pkmm/gb/baidu"
	"pkmm/models"
	"strconv"
	"time"
)

// 初始化函数
func init() {
	toolbox.AddTask("syncUsersForumsFromOfficial", syncUsersForumsFromOfficial)
	toolbox.AddTask("sign", signForums)
	toolbox.AddTask("sync_score_from_zcmu", syncScoreFromZcmu)
	toolbox.AddTask("re_sign", reSignForums)
}

// goroutine 通信数据机构
type ChannelData struct {
	Kw  string
	Fid string
}

type ReplyJson struct {
	Ctime      int64       `json:"ctime"`
	ErrorCode  string      `json:"error_code"`
	ErrorMsg   string      `json:"error_msg"`
	Info       []string    `json:"info"`
	Logid      int64       `json:"logid"`
	ServerTime string      `json:"server_time"`
	Time       int64       `json:"time"`
	UserInfo   interface{} `json:"user_info"`
}

var syncUsersForumsFromOfficial = toolbox.NewTask("syncUsersForumsFromOfficial", "0 0 23 * * *", func() error {
	// 每天11:00 PM 更新贴吧
	//fmt.Println("begin get Userlist")
	users, total, err := models.GetAllUsers()
	if err != nil {
		return err
	}
	if total == 0 {
		return nil
	}
	//fmt.Printf("一共 %d 位用户需要更新\n", total)
	for _, user := range users {
		go func(user *models.User) {
			orm.NewOrm().Raw("update t_forums set is_deleted = 1 where user_id = ?", user.Id).Exec()
			w := baidu.NewForumWorker(user.Bduss)
			forums := w.RetrieveForums()
			size := len(forums)
			ch := make(chan ChannelData, size)
			for _, forum := range forums {
				go func(kw string, ch chan ChannelData) {
					mp := ChannelData{Kw: kw, Fid: w.GetFid(kw)}
					ch <- mp
				}(forum, ch)
			}
			for size > 0 {
				mp := <-ch
				fid, _ := strconv.Atoi(mp.Fid)
				has, _ := orm.NewOrm().QueryTable(models.TableName("forums")).Filter("user_id", user.Id).
					Filter("kw", mp.Kw).Filter("fid", fid).Count()
				if has == 0 {
					orm.NewOrm().Raw("insert into t_forums(user_id, fid, kw, last_sign, created_at, is_deleted) values(?,?,?,?,?,?)",
						user.Id, fid, mp.Kw, -1, time.Now().Unix(), 0).Exec()
				} else {
					orm.NewOrm().Raw("update t_forums set is_deleted = 0 where user_id = ? and fid = ? and kw = ?", user.Id, fid, mp.Kw).Exec()
				}
				size--
			}
		}(user)
	}
	//fmt.Println("end task")
	return nil
})

var reSignForums = toolbox.NewTask("re_sign", "0 30 */1 * * *", func() error {
	// TODO: chunk data.
	users, totalTaskCount, err := models.GetAllUsers()
	if totalTaskCount == 0 {
		return nil
	}
	if err != nil {
		return err
	}

	goroutine := 2
	taskOut := make(chan string, goroutine)
	taskIn := make(chan models.User, totalTaskCount)

	for i := 0; i < goroutine; i++ {
		go func() {
			for {
				user := <-taskIn
				forums, _ := models.GetSignFailureForums(user.Id)
				w := baidu.NewForumWorker(user.Bduss)
				forumList := baidu.ForumList{}
				for _, forum := range forums {
					forumList = append(forumList, baidu.Forum{Kw: forum.Kw, Fid: strconv.Itoa(forum.Fid)})
				}
				ret := w.SignAll(&forumList)
				for kw, reply := range *ret {
					replyJson := ReplyJson{}
					json.Unmarshal([]byte(reply), &replyJson)
					hasError := 1
					if replyJson.ErrorCode == "0" || replyJson.ErrorCode == "160002" {
						hasError = 0
						orm.NewOrm().QueryTable(models.TableName("forums")).
							Filter("user_id", user.Id).
							Filter("kw", kw).
							Update(orm.Params{"reply_json": reply, "last_sign": time.Now().Day(), "sign_status": hasError})
					}

				}
				taskOut <- "okay. end this task."
			}
		}()
	}

	go func() {
		for _, user := range users {
			taskIn <- *user
		}
	}()

	for i := int64(0); i < totalTaskCount; i++ {
		<-taskOut
	}

	return nil
})

var signForums = toolbox.NewTask("sign", "0 0 0 * * *", func() error {
	// 每天 0:00 签到贴吧
	users, total, err := models.GetAllUsers()
	if total == 0 {
		return nil
	}
	if err != nil {
		return err
	}
	for _, user := range users {
		go func(user *models.User) {
			forums, _ := models.NeedSignForumsByUserId(user.Id)
			w := baidu.NewForumWorker(user.Bduss)
			forumList := baidu.ForumList{}
			for _, forum := range forums {
				forumList = append(forumList, baidu.Forum{Kw: forum.Kw, Fid: strconv.Itoa(forum.Fid)})
			}
			ret := w.SignAll(&forumList)
			for kw, reply := range *ret {
				replyJson := ReplyJson{}
				json.Unmarshal([]byte(reply), &replyJson)
				hasError := 1
				if replyJson.ErrorCode == "0" || replyJson.ErrorCode == "160002" {
					hasError = 0
				}
				orm.NewOrm().QueryTable(models.TableName("forums")).
					Filter("user_id", user.Id).
					Filter("kw", kw).
					Update(orm.Params{"reply_json": reply, "last_sign": time.Now().Day(), "sign_status": hasError})
			}
		}(user)
	}

	return nil
})

var syncScoreFromZcmu = toolbox.NewTask("sync_zcmu_grades", "0 */30 * * * *", func() error {
	// todo chunk result
	o := orm.NewOrm()
	var stus []*models.Stu
	num, err := o.QueryTable("stu").All(&stus)
	if err != nil {
		beego.Debug(err)
		return err
	}
	if num == 0 {
		return nil
	}

	totalCount := len(stus) // 总共的任务数量, 否则会直接把500M的内存直接跑满。现在基本上24%的内存
	goroutine := 10         // 并发的数量
	chResStu := make(chan string, goroutine)
	chReqStu := make(chan models.Stu, totalCount)

	// worker 10个worker 等待工作的到来
	for i := 0; i < goroutine; i++ {
		go func() {
			for {
				// 获取任务
				stu := <-chReqStu
				b := time.Now()

				// 以下是任务的核心处理
				var scores []models.Score
				// 登陆尝试
				retry := 3
				crawl := NewCrawl(stu.Num, stu.Pwd)
				for try := 0; try < retry; try++ {
					if scores, err = crawl.Login(); err == nil {
						break
					} else {
						// todo handle error
						sendMail, _ := beego.AppConfig.Bool("mail.send_failure_sync_score")
						if sendMail {
							content := "同步成绩，出现错误 " + err.Error() + ", " + stu.Num
							SendMail("zccxxx79@gmail.com", "PKMM", "Sync Job failed.", content, []string{})
						}
					}
				}
				if len(scores) > 1 {
					for _, score := range scores {
						score.StuId = stu.Id
						models.InsertOrUpdateScore(&score)
					}
				}
				e := time.Since(b)
				// I写输出
				chResStu <- fmt.Sprintf("[%s] has lessons: %02d, Cost time: %s.", stu.Num, len(scores), e.String())
			}
		}()
	}

	// producer 生产者 就一直生产了

	go func() {
		for _, stu := range stus {
			chReqStu <- *stu
		}
	}()

	var br bytes.Buffer
	var s string
	for i := 0; i < totalCount; i++ {
		s = <-chResStu // ignore value.
		br.WriteString(s)
	}

	sendMail, _ := beego.AppConfig.Bool("mail.send_failure_sync_score")

	if sendMail {
		SendMail("zccxxx79@gmail.com", "PKMM", "Sync Job Detail.", br.String(), []string{})
	}
	
	return nil
})
