package utils

import (
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
}

// goroutine 通信数据机构
type ChannelData struct {
	Kw  string
	Fid string
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
				orm.NewOrm().QueryTable(models.TableName("forums")).
					Filter("user_id", user.Id).
					Filter("kw", kw).
					Update(orm.Params{"reply_json": reply, "last_sign": time.Now().Day()})
			}
		}(user)
	}

	return nil
})

var syncScoreFromZcmu = toolbox.NewTask("sync_zcmu_grades", "0 */10 * * * *", func() error {
	// todo chunk result
	o := orm.NewOrm()
	var stus []*models.Stu
	num, err := o.QueryTable("stu").All(&stus)
	if err != nil {
		beego.Debug(err)
		return err
	}
	if num == 0 {
		beego.Debug("没有学生数据")
	}
	beego.Debug(fmt.Sprintf("开始同步学生的成绩了， 一共有%d位同学需要同步", num))
	size := 20 // 并发数
	done := make(chan string, size)
	for indx, stu := range stus {
		go func(stu *models.Stu, indx int) {
			//beego.Debug("开始登陆, 序号: ", indx, stu.Num, stu.Pwd)
			var scores []models.Score
			// 登陆尝试
			retry := 3
			for try := 0; try < retry; try++ {
				scores, err = Login(stu.Num, stu.Pwd)
				if err != nil {
					beego.Debug("第", try, "登陆", stu.Num, "登陆发生错误", err)
				} else {
					break
				}
			}
			//beego.Debug(stu.Num, "成绩的个数", len(scores))
			done <- fmt.Sprintf("[%s %s] 更新的成绩: %d", stu.Num, stu.Pwd, len(scores))
			if len(scores) > 1 {
				//beego.Debug("开始更新 ", stu.Num, "的成绩，共计 ", len(scores))
				for _, score := range scores {
					score.StuId = stu.Id
					score.CreatedAt = time.Now()
					models.InsertOrUpdateScore(&score)
				}
			}
		}(stu, indx)
	}
	for i := 0; i < size; i++ {
		ret := <-done
		beego.Debug(ret)
	}
	return nil
})
