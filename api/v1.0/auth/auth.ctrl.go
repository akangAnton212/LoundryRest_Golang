package auth

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"loundry_rest/database/models"
	"loundry_rest/lib/common"
	"golang.org/x/crypto/bcrypt"
)

// User is alias for models.User
type User = models.User

func hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func checkHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(data common.JSON) (string, error) {

	//  token is valid for 7days
	date := time.Now().Add(time.Hour * 24 * 7)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": data,
		"exp":  date.Unix(),
	})

	// get path from root dir
	pwd, _ := os.Getwd()
	keyPath := pwd + "/jwtsecret.key"

	key, readErr := ioutil.ReadFile(keyPath)
	if readErr != nil {
		return "", readErr
	}
	tokenString, err := token.SignedString(key)
	return tokenString, err
}

func register(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	Username := c.PostForm("Username")
	DisplayName := c.PostForm("DisplayName")
	Password := c.PostForm("Password")

	var exists User
	if err := db.Where("username = ?", Username).First(&exists).Error; err == nil {
		c.AbortWithStatus(409)
		return
	}

	hash, hashErr := hash(Password)
	if hashErr != nil {
		c.AbortWithStatus(500)
		return
	}

	// create user
	user := User{
		Username:     Username,
		DisplayName:  DisplayName,
		PasswordHash: hash,
	}

	db.Create(&user)

	c.JSON(200, common.JSON{
		"status":  true,
		"message": "Sukses Input Data",
	})
}

func login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	Username := c.PostForm("Username")
	Password := c.PostForm("Password")

	// check existancy
	var user User
	if err := db.Where("username = ?", Username).First(&user).Error; err != nil {
		c.AbortWithStatus(404) // user not found
		return
	}

	if !checkHash(Password, user.PasswordHash) {
		c.AbortWithStatus(401)
		return
	}

	serialized := user.Serialize()
	token, _ := generateToken(serialized)

	//c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)

	c.JSON(200, common.JSON{
		"status": true,
		"token": token,
	})
}

// check API will renew token when token life is less than 3 days, otherwise, return null for token
func check(c *gin.Context) {
	userRaw, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(401)
		return
	}

	user := userRaw.(User)

	tokenExpire := int64(c.MustGet("token_expire").(float64))
	now := time.Now().Unix()
	diff := tokenExpire - now

	fmt.Println(diff)
	if diff < 60*60*24*3 {
		// renew token
		token, _ := generateToken(user.Serialize())
		c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)
		c.JSON(200, common.JSON{
			"token": token,
			"user":  user.Serialize(),
		})
		return
	}

	c.JSON(200, common.JSON{
		"token": nil,
		"user":  user.Serialize(),
	})
}
