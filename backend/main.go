package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"meetingmanage/config"
	"meetingmanage/database"
	"meetingmanage/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到 .env 文件，将使用系统环境变量")
	}

	// 初始化配置
	config.InitConfig()

	// 连接数据库
	err = database.ConnectDB()
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移数据库表
	err = database.AutoMigrate()
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化数据
	err = database.SeedData()
	if err != nil {
		log.Fatalf("数据初始化失败: %v", err)
	}

	// 创建路由器
	router := mux.NewRouter()

	// 设置路由
	routes.SetupRoutes(router)

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 启动服务器
	fmt.Printf("服务器正在端口 %s 上运行...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
