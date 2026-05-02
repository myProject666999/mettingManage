package database

import (
	"fmt"
	"log"

	"meetingmanage/config"
	"meetingmanage/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	fmt.Println("数据库连接成功")
	return nil
}

func AutoMigrate() error {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.MeetingRoom{},
		&models.Meeting{},
		&models.Participant{},
		&models.MeetingDocument{},
		&models.Announcement{},
		&models.Carousel{},
		&models.MeetingReminder{},
	)
	if err != nil {
		return err
	}

	fmt.Println("数据库迁移成功")
	return nil
}

func SeedData() error {
	// 创建默认超级管理员
	var adminCount int64
	DB.Model(&models.Admin{}).Count(&adminCount)
	if adminCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		defaultAdmin := models.Admin{
			Username: "admin",
			Password: string(hashedPassword),
			Email:    "admin@example.com",
			Role:     "super_admin",
		}

		if err := DB.Create(&defaultAdmin).Error; err != nil {
			return err
		}

		log.Println("默认超级管理员已创建: 用户名=admin, 密码=admin123")
	}

	// 创建测试用户
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		testUser := models.User{
			Username: "testuser",
			Password: string(hashedPassword),
			Email:    "testuser@example.com",
			Phone:    "13800138001",
			FullName: "测试用户",
		}

		if err := DB.Create(&testUser).Error; err != nil {
			return err
		}

		log.Println("测试用户已创建: 用户名=testuser, 密码=123456")
	}

	// 创建测试会议室
	var roomCount int64
	DB.Model(&models.MeetingRoom{}).Count(&roomCount)
	if roomCount == 0 {
		rooms := []models.MeetingRoom{
			{
				Name:        "一号会议室",
				Location:    "3楼东侧",
				Capacity:    20,
				Description: "大型会议室，配备投影仪、白板",
				Equipment:   "投影仪,白板,音响系统,视频会议设备",
				Status:      "available",
			},
			{
				Name:        "二号会议室",
				Location:    "3楼西侧",
				Capacity:    10,
				Description: "中型会议室，配备电视显示屏",
				Equipment:   "电视显示屏,白板",
				Status:      "available",
			},
			{
				Name:        "三号会议室",
				Location:    "2楼东侧",
				Capacity:    6,
				Description: "小型洽谈室",
				Equipment:   "白板",
				Status:      "available",
			},
			{
				Name:        "多功能厅",
				Location:    "5楼",
				Capacity:    50,
				Description: "大型多功能厅，可用于培训、大型会议",
				Equipment:   "专业音响系统,高清投影仪,电子白板,灯光系统",
				Status:      "available",
			},
		}

		for _, room := range rooms {
			if err := DB.Create(&room).Error; err != nil {
				return err
			}
		}

		log.Println("测试会议室已创建")
	}

	// 创建测试公告
	var announcementCount int64
	DB.Model(&models.Announcement{}).Count(&announcementCount)
	if announcementCount == 0 {
		// 获取第一个管理员作为作者
		var admin models.Admin
		DB.First(&admin)

		announcements := []models.Announcement{
			{
				Title:   "系统上线公告",
				Content: "欢迎使用会议室管理系统！本系统提供会议室预订、会议管理、参会管理等功能。如有问题请联系管理员。",
				AuthorID: admin.ID,
				IsTop:    true,
				Status:   "published",
			},
			{
				Title:   "会议室使用须知",
				Content: "1. 请提前预订会议室，避免冲突。\n2. 使用完毕请保持会议室整洁。\n3. 如遇设备问题请及时联系管理员。\n4. 取消会议请及时释放会议室。",
				AuthorID: admin.ID,
				IsTop:    false,
				Status:   "published",
			},
		}

		for _, announcement := range announcements {
			if err := DB.Create(&announcement).Error; err != nil {
				return err
			}
		}

		log.Println("测试公告已创建")
	}

	fmt.Println("数据初始化完成")
	return nil
}
