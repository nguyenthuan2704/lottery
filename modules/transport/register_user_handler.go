package transport

import (
	"Lottery/modules/biz"
	"Lottery/modules/model"
	"Lottery/modules/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func CreateUser(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data model.PlayerCreation

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}
		store := storage.NewSQLStore(db)
		business := biz.NewCreateUserBusiness(store)
		if err := business.CreateNewUser(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Success"})
	}
}
