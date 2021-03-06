package utils

import (
	"bytes"
	"errors"
	"github.com/astaxie/beego"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"pkmm/models"
	"regexp"
	"strings"
)

const (
	baseUrl  = "http://zfxk.zjtcm.net/"
	loginUrl = "default2.aspx"
	codeUrl  = "CheckCode.aspx"

	POST      = "POST"
	GET       = "GET"
	VIEWSTATE = "__VIEWSTATE"
)

func getViewState(html []byte) (string, error) {
	pattern, _ := regexp.Compile(`<input type="hidden" name="__VIEWSTATE" value="(.*?)" />`)
	viewstate := pattern.FindSubmatch(html)
	if len(viewstate) > 0 {
		return string(viewstate[1]), nil
	}
	return "", errors.New("解析 viewstate 失败")
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func retrieveScores(fileContent []byte) []models.Score {
	// 小记： 使用(?s)标记表示.可以匹配换行符
	pattern := regexp.MustCompile(`(?s)<table .+?id="Datagrid1"[\s\S]*?>(.*?)</table>`)
	ret := pattern.FindSubmatch(fileContent)
	if len(ret) == 0 {
		return []models.Score{}
	}
	table := ret[0]
	// <td>学年</td><td>学期</td><td>课程代码</td><td>课程名称</td><td>课程性质</td><td>课程归属</td><td>学分</td><td>绩点</td><td>成绩</td><td>辅修标记</td><td>补考成绩</td><td>重修成绩</td><td>学院名称</td><td>备注</td><td>重修标记</td><td>课程英文名称</td>
	pattern = regexp.MustCompile(`(?s)<td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>.*?</td><td>.*?</td><td>.*?</td>`)
	tds := pattern.FindAllSubmatch(table, -1)

	var scores []models.Score
	for index, row := range tds {
		if index == 0 {
			continue
		}
		score := models.Score{
			Xn:   string(row[1]),
			Xq:   string(row[2]),
			Kcmc: string(row[3]),
			Type: string(row[4]),
			Xf:   string(row[5]),
			Jd:   string(row[6]),
			Cj:   string(row[7]),
			Bkcj: string(row[8]),
			Cxcj: string(row[9]),
		}
		scores = append(scores, score)
	}
	//fmt.Println(scores)
	return scores
}

// 1. 打开登陆页
func (this *Crawl) openLoginPage() (string, error) {
	rep, err := this.Client.Get(baseUrl)
	defer rep.Body.Close()
	if err != nil {
		return "", errors.New("获取登陆页面失败")
	}
	html, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return "", errors.New("解析登陆页面失败")
	}
	viewState, err := getViewState(html)
	if err != nil {
		return "", errors.New("解析登录页的viewstate失败")
	}
	return viewState, nil
}

// 2. 获取验证码
func (this *Crawl) getCode() (string, error) {
	rep, err := this.Client.Get(baseUrl + codeUrl)
	defer rep.Body.Close()
	if err != err {
		return "", errors.New("加载验证码失败")
	}
	code, err := Predict(rep.Body, false)
	if err != nil {
		return "", err
	}
	return code, nil
}

// 3. 登陆后的主页
func (this *Crawl) GetMainPage() (string, error) {
	viewstate, err := this.openLoginPage()
	if err != nil {
		return "", err
	}
	code, err := this.getCode()
	if err != nil {
		return "", err
	}
	beego.Debug("num", this.Num, "Code is => ", code, len(code))
	formData := url.Values{
		VIEWSTATE:          {viewstate},
		"txtUserName":      {this.Num},
		"Textbox1":         {""},
		"TextBox2":         {this.Pwd},
		"txtSecretCode":    {code},
		"RadioButtonList1": {"%D1%A7%C9%FA"},
		"Button1":          {""},
		"lbLanguage":       {""},
		"hidPdrs":          {""},
		"hidsc":            {""},
	}
	beego.Debug(formData.Encode())
	rep, err := this.Client.PostForm(baseUrl+loginUrl, formData)
	defer rep.Body.Close()
	if err != nil {
		return "", err
	}
	html, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return "", err
	}
	tt, err := GbkToUtf8(html)
	if err != nil {
		return "", errors.New("转码登陆后的页面失败")
	}
	return string(tt), nil
}

func ValidAccount(num, pwd string) (bool, string) {
	crawl := NewCrawl(num, pwd)
	html, err := crawl.GetMainPage()
	if err != nil {
		return false, err.Error()
	} else {
		reg := regexp.MustCompile("验证码不正确")
		result := reg.FindString(html)
		if result != "" {
			return false, "验证码不正确"
		} else {
			reg = regexp.MustCompile("密码错误")
			result = reg.FindString(html)
			if result != "" {
				return false, "密码错误"
			} else {
				reg = regexp.MustCompile("用户名不存在或未按照要求参加教学活动")
				result = reg.FindString(html)
				if result != "" {
					return false, "用户名不存在或未按照要求参加教学活动"
				}
				return true, ""
			}
		}
	}
}

type Crawl struct {
	Client *http.Client
	Num    string
	Pwd    string
}

func NewCrawl(num, pwd string) *Crawl {
	crawl := &Crawl{}
	tmpJar, _ := cookiejar.New(nil)
	crawl.Client = &http.Client{
		Jar: tmpJar,
	}
	crawl.Num = num
	crawl.Pwd = pwd
	return crawl
}

func (this *Crawl) Login() ([]models.Score, error) {
	var err error
	var scores []models.Score
	rep, err := this.Client.Get(baseUrl)
	if err != nil {
		return scores, err
	}
	html, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return scores, err
	}
	viewstate, err := getViewState(html)
	if err != nil {
		return scores, err
	}

	// 加载验证码
	rep, err = this.Client.Get(baseUrl + codeUrl)
	defer rep.Body.Close()
	if err != err {
		return scores, errors.New("加载验证码失败")
	}
	code, _ := Predict(rep.Body, false)

	//	beego.Debug("num", num, "Code is => ", code, len(code))
	formData := url.Values{
		VIEWSTATE:          {viewstate},
		"txtUserName":      {this.Num},
		"Textbox1":         {""},
		"TextBox2":         {this.Pwd},
		"txtSecretCode":    {code},
		"RadioButtonList1": {"%D1%A7%C9%FA"},
		"Button1":          {""},
		"lbLanguage":       {""},
		"hidPdrs":          {""},
		"hidsc":            {""},
	}
	//fmt.Println(formData.Encode())
	rep, err = this.Client.PostForm(baseUrl+loginUrl, formData)
	if err != nil {
		return scores, err
	}
	html, err = ioutil.ReadAll(rep.Body)
	if err != nil {
		return scores, err
	}
	defer rep.Body.Close()
	//tt, _ := GbkToUtf8(html)
	//beego.Debug(string(tt))

	r, err := http.NewRequest(GET, "http://zfxk.zjtcm.net/xscj_gc.aspx?xh="+this.Num+"&xm=%D5%C5%B4%AB%B3%C9&gnmkdm=N121605", nil)
	if err != nil {
		return scores, err
	}
	r.Header.Set("Referer", "http://zfxk.zjtcm.net/xs_main.aspx?xh="+this.Num)
	rep, err = this.Client.Do(r)
	if err != nil {
		//fmt.Println(err)
		return scores, err
	}
	html, err = ioutil.ReadAll(rep.Body)
	if err != nil {
		return scores, err
	}
	//tt,_ := GbkToUtf8(html)
	//fmt.Println(string(tt))

	// 获取viewstate, 用于打开成绩页面
	newViewState, err := getViewState(html)
	if err != nil {
		return scores, err
	}
	//fmt.Println(newViewState)
	//return
	//return
	var ddlXN = ""
	var ddlXQ = ""
	formData = make(url.Values)
	formData.Set(VIEWSTATE, newViewState)
	formData.Set("ddlXN", ddlXN)
	formData.Set("ddlXQ", ddlXQ)
	formData.Set("Button2", "%D4%DA%D0%A3%D1%A7%CF%B0%B3%C9%BC%A8%B2%E9%D1%AF")

	r, err = http.NewRequest(POST,
		"http://zfxk.zjtcm.net/xscj_gc.aspx?xh="+this.Num+"&xm=%D5%C5%B4%AB%B3%C9&gnmkdm=N121605",
		strings.NewReader(formData.Encode()))
	if err != nil {
		return scores, err
	}
	r.Header.Set("Referer", "http://zfxk.zjtcm.net/xs_main.aspx?xh="+this.Num)
	r.Header.Set("Host", "zfxk.zjtcm.net")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded") // 很重要
	rep, err = this.Client.Do(r)
	if err != nil {
		return scores, err
	}
	html, err = ioutil.ReadAll(rep.Body)
	if err != nil {
		return scores, err
	}
	defer rep.Body.Close()
	utf8Html, err := GbkToUtf8(html)
	if err != nil {
		return scores, err
	}

	//fmt.Print(string(utf8Html))
	//out, _ := os.Create("html.txt")
	//io.Copy(out, bytes.NewReader(utf8Html))
	//defer out.Close()

	return retrieveScores(utf8Html), nil

}
