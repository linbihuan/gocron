package calendar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/linbihuan/gocron/libs/dbdriver"
)

type CalendarItem struct {
	Id                  int
	Importance          int
	Influence           string
	Level               int
	Mark                string
	Previous            string
	Actual              string
	Forecast            string
	Revised             string
	Push_status         int
	Related_assets      string
	Remark              string
	Stars               int
	Timestamp           int
	Title               string
	Accurate_flag       int
	Calendar_type       string
	Category_id         int
	Country             string
	Currency            string
	Description         string
	Event_row_id        string
	FlagURL             string
	Ticker              string
	Subscribe_status    int
	Is_has_history_data bool
	Uri                 string
	Calendar_key        string
}

type CalendarData struct {
	Items       []CalendarItem
	Count       int
	Next_cursor string
}

type Calendar struct {
	Code    int
	Message string
	Data    CalendarData
}

func Crawl() {
	fmt.Println("start")
	start := time.Now()

	urlFormat := "https://api-prod.wallstreetcn.com/apiv1/finfo/calendars?start=%d&end=%d"

	tm, _ := time.Parse("2006-01-02 -0700", time.Now().Format("2006-01-02 -0700"))
	ts := tm.Unix()

	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	var i int64
	for i = 0; i <= 30; i++ {
		start := ts + i*86400
		end := ts + (i+1)*86400 - 1
		url := fmt.Sprintf(urlFormat, start, end)
		fmt.Println(url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err)
			continue
		}

		req.Header.Set("Accept", "*/*")
		req.Header.Set("Referer", "https://wallstreetcn.com/calendar")
		req.Header.Set("Origin", "https://wallstreetcn.com")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")

		flag := false
		var resp *http.Response
		for j := 0; j < 20; j++ {
			resp, err = client.Do(req)
			if err != nil {
				fmt.Println("请求失败")
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode == 200 {
				flag = true
				break
			}

			fmt.Printf("第%d重试: %s\n", j, resp.Status)
		}

		if flag == false {
			fmt.Println("请求失败")
			continue
		}

		// 读取出来是byte
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if len(body) == 0 {
			fmt.Println("body 为空")
			continue
		}

		cal := &Calendar{}
		err = json.Unmarshal(body, &cal)
		if err != nil {
			fmt.Println(err)
			continue
		}

		save(cal.Data.Items, start)

		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
	}

	end := time.Now()
	fmt.Println("cost: ", end.Sub(start).Seconds())
	fmt.Println("end")
}

func save(items []CalendarItem, cal_date int64) {
	db, err := dbdriver.MysqlConn()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("%p, %T\n", db, db)
		return
	}

	sql := `insert into stock_calendar_v2
		(country,
		currency,
		flagURL,
		title,
		timestamp,
		cal_date,
		create_time,
		update_time)
		values (?, ?, ?, ?, ?, ?, ?, ?)
		on duplicate key update
		title=values(title),
		country=values(country),
		currency=values(currency),
		flagURL=values(flagURL),
		timestamp=values(timestamp),
		cal_date=values(cal_date),
		update_time=values(update_time)`

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("%p, %T\n", tx, tx)
		return
	}

	for _, item := range items {
		fmt.Println(item)
		tx.Exec(sql, item.Country, item.Currency, item.FlagURL, item.Title, item.Timestamp, cal_date, time.Now().UnixNano(), time.Now().UnixNano())
	}
	tx.Commit()
	db.Close()
}
