package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type RoomConf struct {
	Id        int32  `json:"id"`
	Name      string `json:"name"`
	Level     int32  `json:"level"`
	Type      int32  `json:"type"`
	Needblind int64  `json:"needblind"`
	Blind     int64  `json:"blind"`
	Minmoney  int64  `json:"minmoney"`
	Maxmoney  int64  `json:"maxmoney"`
	Seat      int32  `json:"seat"`
	Audience  int32  `json:"audience"`
	Props     int64  `json:"props"`
	Fee       int64  `json:"fee"`
	Fast      bool   `json:"fast"`
	Status    int8   `json:"status"`
	Allin     bool   `json:"allin"`
	Pcheat    bool   `json:"cheat"`
	LeadChips int64  `json:"lead_chips"`
	Ext       string `json:"ext"`
	Quick     string `json:"quick"`
	ShowNo    string `json:"showno"`
	Ver       int32
	// RoomConf  *RoomConf
}

type GameServer struct {
	// sock network.INet
	// *pbinner.RegisterServer

	Type     int32
	Id       int32
	Level    int32
	Index    int32
	Ip       string
	Port     int32
	Retire   bool
	AutoExit bool
	HttpAddr string

	// 是否全量上报
	Report    bool
	UserCount int32
	// Tables    map[int64]*Table
	// Tables map[int64]map[int64]*Table // map[blind]map[tid]*Table

	maxId     int32
	currMaxId int32
	ids       map[int32]struct{}
	deleteIds map[int32]struct{}
}

func TestAtomic(t *testing.T) {

	t.Log(NewGameServer() == (*GameServer)(nil))
	t.Log(XX() == nil)
}
func XX() (xx *RoomConf) {
	return
}

func NewGameServer() (xx *GameServer) {
	return nil
}

func TestPointerCompare(t *testing.T) {
	now := time.Now()
	t.Log(time.Since(now))
	pointer := make([]*RoomConf, 10)
	for i := range pointer {
		pointer[i] = &RoomConf{
			Id: int32(i) + 1,
		}
	}
	t.Log(time.Since(now))
	left := 1
	right := 1
	sum := 0
	for i := 0; i < 1000000000; i++ {
		if pointer[left] == pointer[right] {
			sum++
		}
	}
	t.Log(time.Since(now))
	left = 1
	right = 1
	for i := 0; i < 1000000000; i++ {
		if *pointer[left] == *pointer[right] {
			sum++
		}
	}
	t.Log(time.Since(now))
}

func TestDB(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local", "domino", "yAUpZwWnjfrPBsWD", "192.168.1.129:3306", "xx_joyfun", "utf8")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Log(err)
		return
	}
	var roomsConf []*RoomConf
	db.Table("room").Where("level=?", 101).Order("blind ASC").Find(&roomsConf)
	if len(roomsConf) == 0 {
		return
	}
	for _, roomConf := range roomsConf {
		roomConf.Type = roomConf.Id
		t.Logf("%+v", roomConf)
	}
}

func TestDB2(t *testing.T) {
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local", "domino", "yAUpZwWnjfrPBsWD", "192.168.1.129:3306", "domino", "utf8")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Log(err)
		return
	}
	// limit := new(RiskLimit)
	var s *struct {
		Price decimal.Decimal
	}
	db.Table("payment").Select("SUM(price) as price").Where("mid = ? AND time >= ?", 7504663, 1641870377).Find(&s)
	f, ok := s.Price.Float64()
	fmt.Println(f, ok)
}
