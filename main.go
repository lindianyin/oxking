package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gocolly/colly/v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

type ZhanKu struct {
	ResultCode string `json:"resultCode"`
	ErrorCode  string `json:"errorCode"`
	Msg        string `json:"msg"`
	Data       string `json:"data"`
	Datas      []struct {
		ObjectType int `json:"objectType"`
		Content    struct {
			Id         int    `json:"id"`
			IdStr      string `json:"idStr"`
			ObjectType int    `json:"objectType"`
			Title      string `json:"title"`
			Cover      string `json:"cover"`
			Cover1X    string `json:"cover1x"`
			Cover2X    string `json:"cover2x"`
			Cover3X    string `json:"cover3x"`
			TypeStr    string `json:"typeStr"`
			Cate       int    `json:"cate"`
			CateStr    string `json:"cateStr"`
			SubCate    int    `json:"subCate"`
			SubCateStr string `json:"subCateStr"`
			Creator    int    `json:"creator"`
			CreatorObj struct {
				Id                 int           `json:"id"`
				Status             int           `json:"status"`
				Avatar             string        `json:"avatar"`
				Username           string        `json:"username"`
				MemberType         int           `json:"memberType"`
				ZcoolStudioType    int           `json:"zcoolStudioType"`
				City               int           `json:"city"`
				CityName           string        `json:"cityName"`
				Profession         int           `json:"profession"`
				ProfessionName     string        `json:"professionName"`
				GuanzhuStatus      int           `json:"guanzhuStatus"`
				FocusStatus        int           `json:"focusStatus"`
				Indusry            []interface{} `json:"indusry"`
				IndusryName        string        `json:"indusryName"`
				MemberHonors       []interface{} `json:"memberHonors"`
				PageUrl            string        `json:"pageUrl"`
				ContentCount       int           `json:"contentCount"`
				ContentCountStr    string        `json:"contentCountStr"`
				PopularityCount    int           `json:"popularityCount"`
				PopularityCountStr string        `json:"popularityCountStr"`
				FansCount          int           `json:"fansCount"`
				FansCountStr       string        `json:"fansCountStr"`
				AvatarPath         string        `json:"avatarPath"`
			} `json:"creatorObj"`
			PublishTimeDiffStr string        `json:"publishTimeDiffStr"`
			TimeTitleStr       string        `json:"timeTitleStr"`
			CreateTime         int64         `json:"createTime"`
			UpdateTime         int64         `json:"updateTime"`
			PublishTime        int64         `json:"publishTime"`
			Recommend          int           `json:"recommend"`
			RecommendTime      int64         `json:"recommendTime"`
			ViewCount          int           `json:"viewCount"`
			CommentCount       int           `json:"commentCount"`
			ViewCountStr       string        `json:"viewCountStr"`
			CommentCountStr    string        `json:"commentCountStr"`
			RecommendCountStr  string        `json:"recommendCountStr"`
			RecommendCount     int           `json:"recommendCount"`
			EventId            int           `json:"eventId"`
			EventObj           interface{}   `json:"eventObj"`
			HasVideoUrl        int           `json:"hasVideoUrl"`
			Status             int           `json:"status"`
			Tags               []interface{} `json:"tags"`
			PageUrl            string        `json:"pageUrl"`
			Show               int           `json:"show"`
			Top                int           `json:"top"`
			Locked             int           `json:"locked"`
			New3FireIcon       int           `json:"new3FireIcon"`
		} `json:"content"`
		DataSource string `json:"data_source"`
	} `json:"datas"`
}

func GetList() []*ZhanKu {
	var res []*ZhanKu
	for i := 1; i <= *count; i++ {
		//urlFmt := "https://www.zcool.com.cn/p1/discover/list?cate=609&city=0&college=0&has_video=0&recommend_level=2&sort=9&sub_cate=637&p=%d&ps=20&column=5"
		urlFmt := "https://www.zcool.com.cn/p1/discover/list?cate=609&city=0&college=0&has_video=0&recommend_level=2&sort=9&p=%d&ps=16&column=4"
		url := fmt.Sprintf(urlFmt, i)
		// 创建一个 HTTP GET 请求
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}

		// 发送 HTTP 请求并获取响应
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// 读取响应内容
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		// 打印响应内容
		//fmt.Println(string(body))
		zk := &ZhanKu{}
		err = json.Unmarshal(body, zk)
		if err != nil {
			return res
		}
		res = append(res, zk)
		//for _, data := range zk.Datas {
		//	fmt.Printf("%v\n", data.Content.PageUrl)
		//}
	}
	return res
}

var htmlFile = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
<table>
	%s
</table>
</body>
</html>`

var row = `
  <tr>
      <td>%s</td>
      <td><img src="%s"></td>
  </tr>
`

type Item struct {
	Name string
	Url  string
}

var Imgs []Item

func DownLoadFile() {
	os.Mkdir("images", 0755)
	sort.Slice(Imgs, func(i, j int) bool {
		return Imgs[i].Name < Imgs[j].Name
	})
	for _, v := range Imgs {
		resp, err := http.Get(v.Url)
		if err != nil {
			return
		}
		s := ".webp"
		if strings.Contains(v.Url, "jpg?") {
			s = ".jpg"
		} else if strings.Contains(v.Url, "png?") {
			s = ".png"
		}
		f, _ := os.Create("images/" + v.Name + s)
		io.Copy(f, resp.Body)
		f.Close()
		resp.Body.Close()
	}

}

var count *int

func main() {
	count = flag.Int("c", 10, "count")
	flag.Parse()
	list := GetList()
	rows := ""
	for _, ku := range list {
		for _, data := range ku.Datas {
			c := colly.NewCollector()
			// Find and visit all links
			c.OnHTML(".photoInformationContent", func(e *colly.HTMLElement) {
				attr := e.DOM.Nodes[0].FirstChild.FirstChild.Attr[0].Val

				rows += fmt.Sprintf(row, data.Content.Title, attr)
				fn := fmt.Sprintf("%v_%v_%v", data.Content.Title, data.Content.CreatorObj.Username, e.Index)
				Imgs = append(Imgs, Item{
					Name: fn,
					Url:  attr,
				})
				fmt.Printf("%v %v\n", fn, attr)
			})

			c.OnRequest(func(r *colly.Request) {
				fmt.Println("Visiting", r.URL)
			})

			c.Visit(data.Content.PageUrl)
		}

	}
	html := fmt.Sprintf(htmlFile, rows)
	file, err := os.Create("index.html")
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(html)
	DownLoadFile()
}
