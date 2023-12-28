package main

import (
	"Lottery/common"
	"Lottery/middleware"
	"Lottery/modules/model"
	"Lottery/modules/transport"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Bets struct {
	common.SQLModel
	PlayerId     int          `json:"player_id" gorm:"column:player_id;"`
	BetNumber    int          `json:"bet_number" gorm:"column:bet_number;"`
	BetTime      string       `json:"bet_time" gorm:"column:bet_time;"`
	ResultNumber int          `json:"result_number" gorm:"column:result_number;"`
	Player       model.Player `gorm:"foreignKey:PlayerId"`
}

func (Bets) TableName() string { return "bets" }

var betsWithUsers []struct {
	Bets
	Player model.Player
}

type BetsCreation struct {
	common.SQLModel
	PlayerId     int       `json:"player_id" gorm:"column:player_id;"`
	BetNumber    int       `json:"bet_number" gorm:"column:bet_number;"`
	BetTime      time.Time `json:"bet_time" gorm:"column:bet_time;"`
	ResultNumber int       `json:"result_number" gorm:"column:result_number;"`
}

func (BetsCreation) TableName() string { return Bets{}.TableName() }

var userBetTimes = make(map[int]time.Time)
var mutex sync.Mutex

func main() {
	dsn := "root:lotto@tcp(db:3306)/vietlott?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("Cannot connect to MySQL:", err)
	}

	log.Println("Connected:", db)
	router := gin.Default()
	router.Use(middleware.Recovery())
	done := make(chan struct{})
	go spinAndSave(db, done)
	v1 := router.Group("/v1")
	{
		v1.POST("/register", transport.CreateUser(db)) // register user
		v1.POST("/login", transport.Login(db))         // login user
		v1.GET("/player/:id", transport.GetItem(db))   //	get user by id
		v1.POST("/bets", createBets(db))
		v1.GET("/history", getListOfBets(db))
		v1.GET("/history/:id", getListOfBetsById(db))
	}

	router.Run(":3030")
}
func validateBetNumber(betnumber int) error {
	if betnumber < 0 || betnumber > 9 {
		return errors.New("Chỉ được nhập 1 số từ 0 đến 9")
	}
	return nil
}
func spinRandomInt(max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}
func spinAndSave(db *gorm.DB, finish chan struct{}) {
	var dbLock sync.Mutex
	for {
		select {
		case <-finish:
			return
		default:
			result, err := spinRandomInt(10)
			if err != nil {
				fmt.Println("Lỗi quay số:", err)
				continue
			}
			var data BetsCreation
			currentTime := time.Now()
			truncatedTime := roundUpToNextHour(currentTime)
			data.BetTime = truncatedTime
			data.ResultNumber = result
			dbLock.Lock()
			db.Model(&BetsCreation{}).Where("result_number = 0").Update("result_number", result)
			if db.RowsAffected == 0 {
				var data BetsCreation
				currentTime := time.Now()
				truncatedTime := roundUpToNextHour(currentTime)
				data.BetTime = truncatedTime
				data.ResultNumber = result

				db.Create(&data)
			}
			dbLock.Unlock()
			time.Sleep(60 * time.Second)
		}
	}

}
func roundUpToNextHour(t time.Time) time.Time {
	return t.Add(time.Hour - time.Duration(t.Minute())*time.Minute - time.Duration(t.Second())*time.Second -
		time.Duration(t.Nanosecond())*time.Nanosecond)
}
func createBets(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataItem BetsCreation

		if err := c.ShouldBind(&dataItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if dataItem.BetNumber < 1 || dataItem.BetNumber > 9 {
			errors.New("Chỉ được nhập 1 số từ 1 đến 9")
			return
		}
		// Lock the mutex to ensure thread safety when accessing the global state
		mutex.Lock()
		defer mutex.Unlock()
		var count int64
		db.Model(&BetsCreation{}).
			Where("user_id = ? AND HOUR(bet_time) = HOUR(?)", dataItem.PlayerId, time.Now()).
			Count(&count)

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bạn chỉ được phép đặt cược trong cùng một khung giờ"})
			return
		}
		// Check if the user has placed a bet before
		lastBetTime, exists := userBetTimes[dataItem.PlayerId]

		// If the user has placed a bet before, check if it's in a different hour
		if exists && lastBetTime.Hour() != time.Now().Hour() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bạn chỉ được phép đặt cược trong cùng một khung giờ"})
			return
		}

		// Update the last bet time for the user
		userBetTimes[dataItem.PlayerId] = roundUpToNextHour(time.Now())

		// Your other bet creation logic here...

		// For demonstration purposes, let's print the user's last bet time
		fmt.Println("Last bet time for user", dataItem.PlayerId, "is", lastBetTime)
		currentTime := time.Now()
		truncatedTime := roundUpToNextHour(currentTime)
		dataItem.BetTime = truncatedTime
		/*var betTime = currentTime*/
		if err := db.Create(&dataItem).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"bạn đã đặt số thành công": dataItem.BetNumber})
	}
}
func getListOfBets(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type DataPaging struct {
			Page  int   `json:"page" form:"page"`
			Limit int   `json:"limit" form:"limit"`
			Total int64 `json:"total" form:"-"`
		}

		var paging DataPaging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if paging.Page <= 0 {
			paging.Page = 1
		}

		if paging.Limit <= 0 {
			paging.Limit = 10
		}

		offset := (paging.Page - 1) * paging.Limit

		var result []Bets

		if err := db.Table(Bets{}.TableName()).Model(&Bets{}).Preload("Player").
			Count(&paging.Total).
			Offset(offset).
			Order("id desc").
			Find(&result).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		/*		if err := db.Table("bets").
				Model(&Bets{}).
				Preload("Player").
				Where("bets.id = ?", id).
				Count(&paging.Total).
				Offset(offset).
				Order("id desc").
				Find(&data).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return*/

		c.JSON(http.StatusOK, gin.H{"data": result})
	}
}
func getListOfBetsById(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type DataPaging struct {
			Page  int   `json:"page" form:"page"`
			Limit int   `json:"limit" form:"limit"`
			Total int64 `json:"total" form:"-"`
		}

		var paging DataPaging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if paging.Page <= 0 {
			paging.Page = 1
		}

		if paging.Limit <= 0 {
			paging.Limit = 10
		}

		offset := (paging.Page - 1) * paging.Limit

		/*var result []Bets*/
		var data []Bets
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		/*		if err := db.Table("bets").Where("id = ?", id).
				Count(&paging.Total).
				Offset(offset).
				Order("id desc").
				First(&data).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}*/
		/*		if err := db.Table("bets").Joins("JOIN players ON players.id = bets.player_id").
				Where("bets.id = ?", id).
				Count(&paging.Total).
				Offset(offset).
				Order("id desc").
				Find(&data).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}*/
		if err := db.Table("bets").
			Model(&Bets{}).
			Preload("Player").
			Where("bets.id = ?", id).
			Count(&paging.Total).
			Offset(offset).
			Order("id desc").
			Find(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": data})
	}
}
