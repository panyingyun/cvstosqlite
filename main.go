package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/panyingyun/cvstosqlite/database"
	"github.com/panyingyun/cvstosqlite/model"

	"github.com/axgle/mahonia"

	log "github.com/Sirupsen/logrus"

	"github.com/codegangsta/cli"
)

//前100个给国民使用，从101~110给姜华的温湿度节点使用
func run(c *cli.Context) error {
	//Get Params
	csvfile := c.String("csv")
	sqlitedb := c.String("db")
	fmt.Println("csvfile = ", csvfile)
	fmt.Println("sqlitedb = ", sqlitedb)

	//Read CVS
	nodes, _ := readCVS(csvfile)
	db, err := database.NewDBEngine(sqlitedb)
	fmt.Println("err = ", err)
	db.InsertAllNodeData(nodes)
	//quit when receive end signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.Infof("signal received signal %v", <-sigChan)
	log.Warn("shutting down server")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "convert to sqlite"
	app.Usage = "auto convert csv file to sqlite"
	app.Copyright = "panyy@lowan-cn.com"
	app.Version = "0.1.0"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "csv",
			Usage:  "CSV",
			Value:  "common.csv",
			EnvVar: "CSV",
		},
		cli.StringFlag{
			Name:   "db",
			Usage:  "DB",
			Value:  "sqlite.db",
			EnvVar: "DB",
		},
	}
	app.Run(os.Args)
}

//读取CSV文件，预处理数据，返回到一个slice,用于统计
func readCVS(file string) ([]model.Node, error) {
	locs := make([]model.Node, 0)
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("(ErrorCode:0001) error is " + err.Error())
		return locs, errors.New("Read File error")
	}
	decode := mahonia.NewDecoder("gbk")
	if decode == nil {
		fmt.Println("(ErrorCode:0002) error is " + err.Error())
		return locs, errors.New("tmahonia.NewDecoder")
	}
	rcsv := csv.NewReader(strings.NewReader(decode.ConvertString(string(buf))))
	//允许双引号字符，否则报错 (ErrorCode:0003) error is line 29752, column 15: extraneous " in field
	rcsv.LazyQuotes = true

	ret, err := rcsv.ReadAll()
	if err != nil {
		fmt.Println("(ErrorCode:0003) error is " + err.Error())
		return locs, errors.New("ReadAll error")
	}
	fmt.Println("Total Count = ", len(ret))
	for _, v := range ret {
		var loc model.Node
		loc.RFU1 = v[0]
		loc.RFU2 = v[1]
		locs = append(locs, loc)
		fmt.Println(loc)
	}
	//fmt.Println(locs)
	return locs, nil
}
