package transport

import (
	"Lottery/modules/biz"
	"Lottery/modules/model"
	"Lottery/modules/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func Login(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var loginUserData model.PlayerLogin

		if err := c.ShouldBind(&loginUserData); err != nil {
			/*panic(common.ErrInvalidRequest(err))*/
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		store := storage.NewSQLStore(db)
		business := biz.NewLoginBusiness(store)
		/*		err := business.Login(c.Request.Context(), &loginUserData)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err.Error(),
					})
					panic(err)
				}*/
		player, err := business.Login(c.Request.Context(), &loginUserData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		/*		response := loginUserData{
				Id:          player.Id,
				PhoneNumber: player.PhoneNumber,
			}*/
		c.JSON(http.StatusOK, gin.H{"data": player.Id})
		/*c.JSON(http.StatusOK, common.SimpleSuccessResponse(loginUserData))*/
	}
}
