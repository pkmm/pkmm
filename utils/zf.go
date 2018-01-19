package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
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

var client *http.Client
var cookieJar *cookiejar.Jar

func init() {
	cookieJar, _ = cookiejar.New(nil)
	client = &http.Client{
		Jar: cookieJar,
	}
}

func getViewState(html []byte) (string, error) {
	pattern, _ := regexp.Compile(`<input type="hidden" name="__VIEWSTATE" value="(.*?)" />`)
	viewstate := pattern.FindSubmatch(html)
	if len(viewstate) > 0 {
		return string(viewstate[1]), nil
	}
	return "", errors.New("解析 viewstate 失败")
}

func downloadImage() (string, error) {
	var err error
	rep, err := client.Get(baseUrl + codeUrl)
	if err != nil {
		return "", err
	}
	picName := UniqueId() + ".png"
	out, err := os.Create("/root/gopath/src/pkmm/utils/zf/verifyCode/" + picName)
	if err != nil {
		return "", err
	}
	io.Copy(out, rep.Body)
	defer out.Close()
	beego.Debug("验证码 已经保存 ", picName)
	return picName, nil
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

type Score struct {
	KCMC string
	XQ   string
	XN   string
	CJ   string
}

func retrieveScores(htmlInlocal bool, fileContent []byte) [][]string {
	if htmlInlocal {
		fileContent, _ = ioutil.ReadFile("html.txt")
	}
	// 小记： 使用(?s)标记表示.可以匹配换行符
	pattern := regexp.MustCompile(`(?s)<table .+?id="Datagrid1"[\s\S]*?>(.*?)</table>`)
	ret := pattern.FindSubmatch(fileContent)
	if len(ret) == 0 {
		return [][]string{}
	}
	table := ret[0]
	// <td>学年</td><td>学期</td><td>课程代码</td><td>课程名称</td><td>课程性质</td><td>课程归属</td><td>学分</td><td>绩点</td><td>成绩</td><td>辅修标记</td><td>补考成绩</td><td>重修成绩</td><td>学院名称</td><td>备注</td><td>重修标记</td><td>课程英文名称</td>
	pattern = regexp.MustCompile(`(?s)<td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>(.*?)</td><td>.*?</td><td>.*?</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>(.*?)</td><td>(.*?)</td><td>.*?</td><td>.*?</td><td>.*?</td><td>.*?</td>`)
	tds := pattern.FindAllSubmatch(table, -1)

	ans := make([][]string, 0, len(tds))
	for _, td := range tds {
		row := make([]string, 0, len(td))
		var i int
		for ii, x := range td {
			i = ii
			if ii == 0 {
				continue
			}
			row = append(row, string(x))
		}
		if i == 0 {
			continue
		}
		ans = append(ans, row)
	}
	return ans
}

func Login(num, pwd string) ([][]string, error) {
	var err error
	rep, _ := client.Get(baseUrl)
	html, _ := ioutil.ReadAll(rep.Body)
	viewstate, err := getViewState(html)
	if err != nil {
		return [][]string{}, err
	}
	//picName, err := downloadImage()
	//if err != nil {
	//	return [][]string{}, err
	//}
	//code, err := imgToString(picName)
	//if err != nil {
	//	return [][]string{}, err
	//}

	// 加载验证码
	rep, err = client.Get(baseUrl + codeUrl)
	defer rep.Body.Close()
	if err != err {
		return [][]string{}, errors.New("加载验证码失败")
	}
	code := Predict(rep.Body, false)

	beego.Debug("num", num, "Code is => ", code, len(code))
	formData := url.Values{
		VIEWSTATE:          {viewstate},
		"txtUserName":      {num},
		"Textbox1":         {""},
		"TextBox2":         {pwd},
		"txtSecretCode":    {code},
		"RadioButtonList1": {"%D1%A7%C9%FA"},
		"Button1":          {""},
		"lbLanguage":       {""},
		"hidPdrs":          {""},
		"hidsc":            {""},
	}
	fmt.Println(formData.Encode())
	rep, _ = client.PostForm(baseUrl+loginUrl, formData)
	html, _ = ioutil.ReadAll(rep.Body)
	defer rep.Body.Close()
	//tt,_ := GbkToUtf8(html)
	//fmt.Println(string(tt))

	r, _ := http.NewRequest(GET, "http://zfxk.zjtcm.net/xscj_gc.aspx?xh="+num+"&xm=%D5%C5%B4%AB%B3%C9&gnmkdm=N121605", nil)
	r.Header.Set("Referer", "http://zfxk.zjtcm.net/xs_main.aspx?xh="+num)
	rep, err = client.Do(r)
	if err != nil {
		//fmt.Println(err)
		return [][]string{}, err
	}
	html, _ = ioutil.ReadAll(rep.Body)
	//tt,_ := GbkToUtf8(html)
	//fmt.Println(string(tt))

	// 获取viewstate, 用于打开成绩页面
	newViewState, err := getViewState(html)
	if err != nil {
		return [][]string{}, err
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

	r, _ = http.NewRequest(POST,
		"http://zfxk.zjtcm.net/xscj_gc.aspx?xh="+num+"&xm=%D5%C5%B4%AB%B3%C9&gnmkdm=N121605",
		strings.NewReader(formData.Encode()))
	r.Header.Set("Referer", "http://zfxk.zjtcm.net/xs_main.aspx?xh="+num)
	r.Header.Set("Host", "zfxk.zjtcm.net")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded") // 很重要
	rep, err = client.Do(r)
	if err != nil {
		fmt.Println(err)
	}
	html, _ = ioutil.ReadAll(rep.Body)
	defer rep.Body.Close()
	utf8Html, _ := GbkToUtf8(html)

	//fmt.Print(string(utf8Html))
	//out, _ := os.Create("html.txt")
	//io.Copy(out, bytes.NewReader(utf8Html))
	//defer out.Close()

	return retrieveScores(false, utf8Html), nil

}

// 通过图片的路径去取图片然后识别验证码（python识别代码实现）
func imgToString(imageFilePath string) (string, error) {
	ans, err := exec.Command("/usr/bin/python", "/root/gopath/src/pkmm/utils/zf/verifyCode/test.py", imageFilePath).Output()
	//fmt.Println("decode verify code:", err)
	beego.Debug(string(ans))
	if err != nil {
		return "", errors.New("识别验证码失败")
	}
	rs := string(ans)
	return rs[:4], nil
}
func main() {
	Login("201312203501029", "520asd")
}