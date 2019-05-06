package posts

import (
	"github.com/gin-gonic/gin"
	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/jinzhu/gorm"
	"loundry_rest/database/models"
	"loundry_rest/lib/common"
	"strconv"
)

// Post type alias
type Post = models.Post

// User type alias
type User = models.User

// JSON type alias
type JSON = common.JSON

func create(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	Text := c.PostForm("Text")

	user := c.MustGet("user").(User)
	post := Post{Text: Text, User: user}
	db.NewRecord(post)
	db.Create(&post)
	// c.JSON(200, post.Serialize())
	c.JSON(200, gin.H{
		"status"	:  true,
		"message"	: "Sukses Input Data",
	})
}

func list(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	//limits := "2"

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	// limit, _ := strconv.Atoi(c.DefaultQuery("limit", limits))
	limit := 2

	var posts []Post

	query := db.Preload("User").Limit(limit).Offset(page).Order("id desc").Find(&posts)

    paginator := pagination.Paging(&pagination.Param{
        DB:      db,
        Page:    page,
        Limit:   limit,
        OrderBy: []string{"id desc"},
        ShowSQL: true,
	}, query)
	
	c.JSON(200, gin.H{
		"status"	: true,
		"data"		: paginator,
	})
}

func readById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var post Post

	// auto preloads the related model
	// http://gorm.io/docs/preload.html#Auto-Preloading
	if err := db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatusJSON(404, gin.H{
			"status": false,
			"message": "Data Tidak Ada!!",
		})
		return
	}

	c.JSON(200, post.Serialize())
}

func readByParams(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Query("id")
	var post Post

	// auto preloads the related model
	// http://gorm.io/docs/preload.html#Auto-Preloading
	if err := db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatusJSON(404, gin.H{
			"status": false,
			"message": "Data Tidak Ada!!",
		})
		return
	}

	c.JSON(200, post.Serialize())
}

func remove(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	user := c.MustGet("user").(User)

	var post Post
	if err := db.Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatusJSON(404, gin.H{
			"status": false,
			"message": "Data Tidak Ada!!",
		})
		return
	}

	if post.UserID != user.ID {
		c.AbortWithStatusJSON(403, gin.H{
			"status": false,
			"message": "Anda Tidak Berhak Menghapus Data Ini...!!",
		})
		return
	}

	db.Delete(&post)
	c.Status(204)
}

func update(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	//id := c.Param("id")

	user := c.MustGet("user").(User)

	Text := c.PostForm("Text")
	id := c.PostForm("id")

	var post Post
	if err := db.Preload("User").Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatusJSON(404, gin.H{
			"status": false,
			"message": "Data Tidak Ada!!",
		})
		return
	}

	//kalo user id yang nginput dan yg edit beda maka :
	if post.UserID != user.ID {
		c.AbortWithStatusJSON(403, gin.H{
			"status": false,
			"message": "Anda Tidak Berhak Mengedit Data Ini...!!",
		})
		return
	}

	post.Text = Text
	db.Save(&post)
	c.JSON(200, post.Serialize())
}
