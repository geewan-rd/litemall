package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var prices []float64 = []float64{27, 75, 150, 282, 70, 25, 260, 135, 18, 147}

type goods struct {
	ID            int    `gorm:"column:id;type:int;primaryKey"`
	SN            int    `gorm:"column:goods_sn;type:varchar(63)"`
	Name          string `gorm:"column:name;type:varchar(127)"`
	Price         float64
	CurrPrice     float64
	IsNew         bool
	IsHot         bool
	OnSale        bool
	MainPicURL    string
	GalleryPicURL string
	DescShort     string
	DescHTML      string
	PicURL        string
}

func main() {
	db, err := gorm.Open("mysql", "litemall:789632145@(rm-j6c8o2wmj3i1x28i38o.mysql.rds.aliyuncs.com)/litemall?charset=utf8")
	if err != nil {
		log.Print(err)
		return
	}
	f, err := excelize.OpenFile("goods_msg_gogo123.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("SheetJS")
	for _, row := range rows[1:] {
		var err error
		g := goods{}
		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}
		g.ID = id + 1200000
		g.Name = row[2]
		g.SN, err = strconv.Atoi(row[1])
		if err != nil {
			continue
		}
		g.SN = 1200000 + g.SN
		// g.Price, err = strconv.ParseFloat(row[3], 10)
		// if err != nil {
		// 	continue
		// }
		// g.CurrPrice, err = strconv.ParseFloat(row[4], 10)
		// if err != nil {
		// 	continue
		// }
		g.Price = prices[rand.Intn(4)]
		g.CurrPrice = g.Price
		g.IsNew, err = strconv.ParseBool(row[5])
		g.IsHot, err = strconv.ParseBool(row[6])
		g.MainPicURL = row[8]
		g.GalleryPicURL = strings.Replace(row[9], `'`, `"`, -1)
		g.DescShort = row[10]
		g.DescHTML = row[11]
		g.PicURL = row[1]

		db.Exec("DELETE FROM litemall_goods WHERE id=?", g.ID)
		db.Exec("DELETE FROM litemall_goods_product WHERE goods_id=?", g.ID)
		db.Exec("DELETE FROM litemall_goods_specification WHERE goods_id=?", g.ID)

		if g.Price <= 0 || g.CurrPrice <= 0 {
			continue
		}

		err = db.Exec("INSERT INTO litemall_goods (id,goods_sn,name,category_id,gallery,brief,is_on_sale,pic_url,is_new,is_hot,unit,counter_price,retail_price,detail) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?);",
			g.ID, g.SN, g.Name, 1025000, g.GalleryPicURL, g.DescShort, 1, g.MainPicURL, g.IsNew, g.IsHot, "件", g.Price, g.CurrPrice, g.DescHTML).Error
		if err != nil {
			log.Print(err)

		}
		err = db.Exec("INSERT INTO litemall_goods_product (goods_id,specifications,price,number,url) VALUES (?,?,?,?,?)",
			g.ID, `["标准"]`, g.Price, 100, g.MainPicURL).Error
		if err != nil {
			log.Print(err)

		}
		err = db.Exec("INSERT INTO litemall_goods_specification (goods_id,specification,value) VALUES (?,?,?)",
			g.ID, `规格`, "标准").Error
		if err != nil {
			log.Print(err)

		}
	}
}
