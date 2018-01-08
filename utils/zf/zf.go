package zf

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
	"os/exec"
	"encoding/base64"
	"crypto/md5"
	"encoding/hex"
	"crypto/rand"
	"pkmm/utils"
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
	fmt.Println("初始化成功")
}

func getViewState(html []byte) string {
	pattern, _ := regexp.Compile(`<input type="hidden" name="__VIEWSTATE" value="(.*?)" />`)
	viewstate := pattern.FindSubmatch(html)
	if len(viewstate) > 0 {
		return string(viewstate[1])
	}
	return ""
}

func downloadImage() string {
	rep, _ := client.Get(baseUrl + codeUrl)
	picName := UniqueId() + ".png"
	out, _ := os.Create("/root/gopath/src/pkmm/utils/zf/verifyCode/" + picName)
	io.Copy(out, rep.Body)
	defer out.Close()
	//fmt.Printf("验证码 已经保存: %s\n", picName)
	return picName
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
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

func Login(num, pwd string) [][]string {
	var err error
	rep, _ := client.Get(baseUrl)
	html, _ := ioutil.ReadAll(rep.Body)
	viewstate := getViewState(html)
	//fmt.Println(viewstate)
	//fmt.Println(rep.Cookies())
	picName := downloadImage()
	//var code string
	//fmt.Println("请输入验证码")
	//fmt.Scanln(&code)
	//fmt.Println("输入的验证码是：" + code)
	code := imgToString(picName)
	//fmt.Println("Code is => ", code, len(code))
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
	//fmt.Println(formData.Encode())
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
		return [][]string{}
	}
	html, _ = ioutil.ReadAll(rep.Body)
	//tt,_ := GbkToUtf8(html)
	//fmt.Println(string(tt))

	// 获取viewstate, 用于打开成绩页面
	newViewState := getViewState(html)
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

	return retrieveScores(false, utf8Html)

}

// 通过图片的路径去取图片然后识别验证码（python识别代码实现）
func imgToString(imageFilePath string) string {
	ans, err := exec.Command("/usr/bin/python", "/root/gopath/src/pkmm/utils/zf/verifyCode/test.py", imageFilePath).Output()
	//fmt.Println("decode verify code:", err)
	//fmt.Println(string(ans))
	if err != nil {
		//fmt.Println(err)
		panic("exit")
	}
	rs := string(ans)
	rs = strings.TrimRight(rs, "\n")
	length := len(rs)
	return rs[length-5:length-1]
}
func main() {

	//var zf_url = "http://zfxk.zjtcm.net/xscj_gc.aspx?xh=201312203501029&xm=%D5%C5%B4%AB%B3%C9&gnmkdm=N121605"
	//var client = &http.Client{}
	//var cookie = http.Cookie{Name: "ASP.NET_SessionId", Value: "2j3zjo45vpijnr55pamzjzqn"}
	//var data = url.Values{
	//	"__VIEWSTATE": {"dDwxODI2NTc3MzMwO3Q8cDxsPHhoOz47bDwyMDEzMTIyMDM1MDEwMjk7Pj47bDxpPDE+Oz47bDx0PDtsPGk8MT47aTwzPjtpPDU+O2k8Nz47aTw5PjtpPDExPjtpPDEzPjtpPDE2PjtpPDI2PjtpPDI3PjtpPDI4PjtpPDM1PjtpPDM3PjtpPDM5PjtpPDQxPjtpPDQ1Pjs+O2w8dDxwPHA8bDxUZXh0Oz47bDzlrablj7fvvJoyMDEzMTIyMDM1MDEwMjk7Pj47Pjs7Pjt0PHA8cDxsPFRleHQ7PjtsPOWnk+WQje+8muW8oOS8oOaIkDs+Pjs+Ozs+O3Q8cDxwPGw8VGV4dDs+O2w85a2m6Zmi77ya5Yy75a2m5oqA5pyv5a2m6ZmiOz4+Oz47Oz47dDxwPHA8bDxUZXh0Oz47bDzkuJPkuJrvvJo7Pj47Pjs7Pjt0PHA8cDxsPFRleHQ7PjtsPOWMu+WtpuS/oeaBr+W3peeoizs+Pjs+Ozs+O3Q8cDxwPGw8VGV4dDs+O2w86KGM5pS/54+t77ya5Yy75a2m5L+h5oGv5bel56iLMjAxM+e6pzHnj607Pj47Pjs7Pjt0PHA8cDxsPFRleHQ7PjtsPDIwMTMwOTE2Oz4+Oz47Oz47dDx0PHA8cDxsPERhdGFUZXh0RmllbGQ7RGF0YVZhbHVlRmllbGQ7PjtsPFhOO1hOOz4+Oz47dDxpPDU+O0A8XGU7MjAxNi0yMDE3OzIwMTUtMjAxNjsyMDE0LTIwMTU7MjAxMy0yMDE0Oz47QDxcZTsyMDE2LTIwMTc7MjAxNS0yMDE2OzIwMTQtMjAxNTsyMDEzLTIwMTQ7Pj47Pjs7Pjt0PHA8O3A8bDxvbmNsaWNrOz47bDx3aW5kb3cucHJpbnQoKVw7Oz4+Pjs7Pjt0PHA8O3A8bDxvbmNsaWNrOz47bDx3aW5kb3cuY2xvc2UoKVw7Oz4+Pjs7Pjt0PHA8cDxsPFZpc2libGU7PjtsPG88dD47Pj47Pjs7Pjt0PEAwPDs7Ozs7Ozs7Ozs+Ozs+O3Q8QDA8Ozs7Ozs7Ozs7Oz47Oz47dDxAMDw7Ozs7Ozs7Ozs7Pjs7Pjt0PDtsPGk8MD47aTwxPjtpPDI+O2k8ND47PjtsPHQ8O2w8aTwwPjtpPDE+Oz47bDx0PDtsPGk8MD47aTwxPjs+O2w8dDxAMDw7Ozs7Ozs7Ozs7Pjs7Pjt0PEAwPDs7Ozs7Ozs7Ozs+Ozs+Oz4+O3Q8O2w8aTwwPjtpPDE+Oz47bDx0PEAwPDs7Ozs7Ozs7Ozs+Ozs+O3Q8QDA8Ozs7Ozs7Ozs7Oz47Oz47Pj47Pj47dDw7bDxpPDA+Oz47bDx0PDtsPGk8MD47PjtsPHQ8QDA8Ozs7Ozs7Ozs7Oz47Oz47Pj47Pj47dDw7bDxpPDA+O2k8MT47PjtsPHQ8O2w8aTwwPjs+O2w8dDxAMDxwPHA8bDxWaXNpYmxlOz47bDxvPGY+Oz4+Oz47Ozs7Ozs7Ozs7Pjs7Pjs+Pjt0PDtsPGk8MD47PjtsPHQ8QDA8cDxwPGw8VmlzaWJsZTs+O2w8bzxmPjs+Pjs+Ozs7Ozs7Ozs7Oz47Oz47Pj47Pj47dDw7bDxpPDA+Oz47bDx0PDtsPGk8MD47PjtsPHQ8cDxwPGw8VGV4dDs+O2w8SkxVOz4+Oz47Oz47Pj47Pj47Pj47dDxAMDw7Ozs7Ozs7Ozs7Pjs7Pjs+Pjs+Pjs+WUUxE3p4x9UKYY9kmhjPe0w8ZNY="},
	//	"ddlXN":       {""},
	//	"ddlXQ":       {""},
	//	"Button1":     {"%B0%B4%D1%A7%C6%DA%B2%E9%D1%AF"},
	//}
	//r, _ := http.NewRequest("POST", zf_url, strings.NewReader(data.Encode()))
	//r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//r.Header.Set("Referer", zf_url)
	//r.AddCookie(&cookie)
	//t, err := client.Do(r)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//body, _ := ioutil.ReadAll(t.Body)
	////fmt.Println(string(body))
	//re, _ := regexp.Compile(`<td>(.*?)</td>`)
	//rets := re.FindAllSubmatch(body, -1)
	//
	//for k, v := range rets {
	//	fmt.Println(k, string(v[1]))
	//
	//}

	Login("201312203501029", "520asd")

}
