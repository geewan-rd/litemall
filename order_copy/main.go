package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type tranOrder struct {
	ID        int    `gorm:"column:id;type:int;primarykey"`
	No        string `gorm:"column:order_no;type:varchar(100)"`
	Amount    int    `gorm:"column:amount;type:int"`
	State     int    `gorm:"column:state"`
	CreatedAt string `gorm:"column:created_at;type:timestamp"`
}

func (tranOrder) TableName() string {
	return "orders"
}

type tranCharge struct {
	ID      int    `gorm:"column:id;type:int;primarykey"`
	OrderID int    `gorm:"colume:order_id"`
	OrderNo string `gorm:"column:order_no"`
	Channel string `gorm:"colume:channel"`
}

func (tranCharge) TableName() string {
	return "charges"
}

type mallOrder struct {
	ID          int     `gorm:"column:id;type:int;primarykey"`
	SN          string  `gorm:"column:order_sn;type:varchar(63)"`
	UserID      int     `gorm:"column:user_id"`
	Status      int     `gorm:"column:order_status"`
	GoodsPrice  float64 `gorm:"column:goods_price"`
	CouponPrice float64 `gorm:"column:coupon_price"`
	OrderPrice  float64 `gorm:"column:order_price"`
	ActualPrice float64 `gorm:"column:actual_price"`
	PayID       string  `gorm:"column:pay_id"`
	PayTime     string  `gorm:"column:pay_time"`
}

func (o *mallOrder) TableName() string {
	return "litemall_order"
}

type mallOrderGoods struct {
	ID             int     `gorm:"column:id;type:int;primarykey"`
	OrderID        int     `gorm:"column:order_id"`
	GoodsID        int     `gorm:"column:goods_id"`
	GoodsName      string  `gorm:"column:goods_name"`
	GoodsSN        string  `gorm:"column:goods_sn"`
	ProductID      int     `gorm:"column:product_id"`
	Number         int     `gorm:"column:number"`
	Price          float64 `gorm:"column:price"`
	Specifications string  `gorm:"column:specifications"`
	PicURL         string  `gorm:"column:pic_url"`
}

func (o *mallOrderGoods) TableName() string {
	return "litemall_order_goods"
}

type goods struct {
	ID    int     `gorm:"column:id;type:int;primaryKey"`
	SN    string  `gorm:"column:goods_sn"`
	Name  string  `gorm:"column:name"`
	Price float64 `gorm:"column:counter_price"`
}

func (g *goods) TableName() string {
	return "litemall_goods"
}

type product struct {
	ID             int    `gorm:"column:id;type:int;primaryKey"`
	GoodsID        int    `gorm:"column:goods_id"`
	Specifications string `gorm:"column:specifications"`
	URL            string `gorm:"column:url"`
}

func (p *product) TableName() string {
	return "litemall_goods_product"
}

func main() {
	mallDB, err := gorm.Open("mysql", "litemall:789632145@(rm-j6c8o2wmj3i1x28i38o.mysql.rds.aliyuncs.com)/litemall?charset=utf8")
	if err != nil {
		log.Print(err)
		return
	}
	tranDB, err := gorm.Open("mysql", "read_only:transocks789632145@(rm-j6cben96r6794a3kkno.mysql.rds.aliyuncs.com)/fobwifi?charset=utf8")
	if err != nil {
		log.Print(err)
		return
	}
	t := time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05")
	tOrders := make([]tranOrder, 0)
	err = tranDB.Where(`created_at > ?`, t).Find(&tOrders).Error
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("order count: %d", len(tOrders))
	for _, o := range tOrders {
		if o.Amount <= 0 {
			continue
		}
		if o.State < 4 {
			continue
		}

		charge := &tranCharge{}
		err = tranDB.Where(&tranCharge{OrderID: o.ID}).First(charge).Error
		if err != nil {
			log.Printf("find charge failed:%s", err)
			continue
		}
		if charge.Channel != "expresspay" {
			continue
		}
		gs := make([]goods, 0)
		err = mallDB.Where(&goods{Price: float64(o.Amount) / 100}).Find(&gs).Error
		if err != nil || len(gs) == 0 {
			log.Printf("Find goods failed:%s", err)
			continue
		}
		g := gs[rand.Intn(len(gs))]
		p := &product{}
		err = mallDB.Where(&product{GoodsID: g.ID}).First(p).Error
		if err != nil {
			log.Printf("Find product failed:%s", err)
			continue
		}
		temp := &mallOrder{}
		err = mallDB.Where(&mallOrder{SN: o.No}).First(temp).Error
		if err == nil || temp.ID > 0 {
			continue
		}
		mOrder := &mallOrder{
			SN:          o.No,
			UserID:      1,
			Status:      301,
			GoodsPrice:  g.Price,
			OrderPrice:  g.Price,
			ActualPrice: g.Price,
			PayID:       strconv.Itoa(charge.ID),
			PayTime:     o.CreatedAt,
		}
		err = mallDB.Create(mOrder).Error
		if err != nil {
			log.Printf("Create order failed:%s", err)
			continue
		}
		log.Printf("%#v", mOrder)
		mOrderGoods := &mallOrderGoods{
			OrderID:        mOrder.ID,
			GoodsID:        g.ID,
			GoodsName:      g.Name,
			GoodsSN:        g.SN,
			ProductID:      p.ID,
			Number:         1,
			Price:          g.Price,
			Specifications: p.Specifications,
			PicURL:         p.URL,
		}
		err = mallDB.Create(mOrderGoods).Error
		if err != nil {
			log.Printf("Create order goods failed:%s", err)
			continue
		}
		log.Printf("%#v", mOrderGoods)
	}
}
