package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"project2/models"
)

func main() {
	//https://blog.csdn.net/m0_73337964/article/details/138518931
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1)/ginChat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	//迁移 scheme
	//db.AutoMigrate(models.UserBasic{})
	//db.AutoMigrate(&models.Message{})
	db.AutoMigrate(&models.Contact{})
	db.AutoMigrate(&models.GroupBasic{})
	////Create
	//user := &models.UserBasic{}
	////user.Name = "张三"
	////db.Create(user)
	//
	////Read
	//var u []models.UserBasic
	//fmt.Println("-----", db.Select("name").Where("name=?", "张三").First(&u))
	////Update
	//db.Model(user).Update("Password", "12345688")
	//db.Model(user).Updates(models.UserBasic{Phone: "12345678910", Identity: "学生"})
	//db.Model(user).Updates(map[string]interface{}{"Password": "111111", "Phone": "1111122"})

	//Delete删除
	//db.Delete(user)
}
