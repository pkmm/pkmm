package utils

import (
	"github.com/astaxie/beego/toolbox"
	"pkmm/models"
	"github.com/pkmm/gb/baidu"
	"strconv"
	"fmt"
	"time"
	"github.com/astaxie/beego/orm"
)

// 初始化函数
func init() {
	toolbox.AddTask("getForums", getForums)
	toolbox.AddTask("sign", signForums)
}

type ChannelData struct {
	Kw  string
	Fid string
}

var getForums = toolbox.NewTask("getForums", "0 30 23,18 * * *", func() error {
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
				forum := models.Forum{UserId: user.Id, Fid: fid, Kw: mp.Kw, LastSign: -1}
				num, _ := orm.NewOrm().QueryTable(models.TableName("forums")).Filter("user_id", user.Id).Filter("kw", mp.Kw).Count()
				if num == 0 {
					models.AddForum(&forum)
				}
				size--
			}
		}(user)
	}
	//fmt.Println("end task")
	return nil
})

var signForums = toolbox.NewTask("sign", "0 0 0,12 * * *", func() error {

	users, total, err := models.GetAllUsers()
	if total == 0 {
		return nil
	}
	if err != nil {
		return err
	}
	for _, user := range users {
		go func(user *models.User) {
			forums, _ := models.GetForumsByUserId(user.Id)
			w := baidu.NewForumWorker(user.Bduss)
			forumList := baidu.ForumList{}
			for _, forum := range forums {
				forumList = append(forumList, baidu.Forum{Kw: forum.Kw, Fid: strconv.Itoa(forum.Fid)})
			}
			ret := w.SignAll(&forumList)
			for kw, reply := range *ret {
				fmt.Println(kw)
				orm.NewOrm().QueryTable(models.TableName("forums")).Filter("user_id", user.Id).
					Filter("kw", kw).Update(orm.Params{"reply_json": reply, "last_sign": time.Now().Day()})
			}
		}(user)
	}

	return nil
})
