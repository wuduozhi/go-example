package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:XGDB)$!%6558ca@tcp(qnxg.net:3306)/weihuda?charset=utf8")
	checkErr(err)

	sql_query := "SELECT DISTINCT(tp_bind.stuID),tp_bind.stuPASS,tp_bind.hdjwPASS,grade_notify.email FROM tp_bind JOIN grade_notify ON tp_bind.stuId=grade_notify.stuID AND grade_notify.xn=2018 and grade_notify.xq=0"

	// query
	rows, err := db.Query(sql_query)
	checkErr(err)

	for rows.Next() {
		var xh string
		var ptPass sql.NullString
		var hdjwPass sql.NullString
		var email sql.NullString
		err = rows.Scan(&xh, &ptPass, &hdjwPass, &email)
		checkErr(err)
		// fmt.Println(xh)
		// fmt.Println(email)
		pt := ptPass.String
		hdjw := hdjwPass.String
		getSpiderData(xh, pt, hdjw)
	}

}

func getSpiderData(xh, ptPass, hdjwPass string) {
	url := "http://spider.qnxg.net/bks/grade?stuid=" + xh + "&password=" + hdjwPass + "&xn=2018&xq=1"
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
