package routes

import (
	"meetingmanage/controllers"
	"meetingmanage/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
	// 用户控制器
	userController := &controllers.UserController{}
	
	// 管理员控制器
	adminController := &controllers.AdminController{}
	
	// 会议室控制器
	meetingRoomController := &controllers.MeetingRoomController{}
	
	// 会议控制器
	meetingController := &controllers.MeetingController{}
	
	// 参会者控制器
	participantController := &controllers.ParticipantController{}
	
	// 会议资料控制器
	meetingDocumentController := &controllers.MeetingDocumentController{}
	
	// 公告控制器
	announcementController := &controllers.AnnouncementController{}
	
	// 轮播图控制器
	carouselController := &controllers.CarouselController{}
	
	// 会议提醒控制器
	meetingReminderController := &controllers.MeetingReminderController{}

	// 公共路由 - 不需要认证
	api := router.PathPrefix("/api").Subrouter()
	
	// 用户认证路由
	api.HandleFunc("/register", userController.Register).Methods("POST")
	api.HandleFunc("/login", userController.Login).Methods("POST")
	api.HandleFunc("/admin/login", adminController.Login).Methods("POST")
	
	// 公共信息路由
	api.HandleFunc("/announcements", announcementController.GetAll).Methods("GET")
	api.HandleFunc("/announcements/{id}", announcementController.GetByID).Methods("GET")
	api.HandleFunc("/carousels", carouselController.GetActive).Methods("GET")
	api.HandleFunc("/meeting-rooms", meetingRoomController.GetAvailable).Methods("GET")
	api.HandleFunc("/meeting-rooms/{id}", meetingRoomController.GetByID).Methods("GET")

	// 用户路由 - 需要用户认证
	userRoutes := api.PathPrefix("/user").Subrouter()
	userRoutes.Use(middleware.AuthMiddleware)
	userRoutes.Use(middleware.UserAuthMiddleware)
	
	// 用户个人信息
	userRoutes.HandleFunc("/profile", userController.GetProfile).Methods("GET")
	userRoutes.HandleFunc("/profile", userController.UpdateProfile).Methods("PUT")
	userRoutes.HandleFunc("/change-password", userController.ChangePassword).Methods("POST")
	
	// 会议室相关
	userRoutes.HandleFunc("/meeting-rooms", meetingRoomController.GetAll).Methods("GET")
	userRoutes.HandleFunc("/meeting-rooms/{id}", meetingRoomController.GetByID).Methods("GET")
	
	// 会议相关
	userRoutes.HandleFunc("/meetings", meetingController.GetUserMeetings).Methods("GET")
	userRoutes.HandleFunc("/meetings/{id}", meetingController.GetByID).Methods("GET")
	userRoutes.HandleFunc("/meetings", meetingController.Create).Methods("POST")
	userRoutes.HandleFunc("/meetings/{id}", meetingController.Update).Methods("PUT")
	userRoutes.HandleFunc("/meetings/{id}", meetingController.Cancel).Methods("DELETE")
	
	// 参会相关
	userRoutes.HandleFunc("/participants", participantController.GetUserParticipations).Methods("GET")
	userRoutes.HandleFunc("/participants/join", participantController.JoinMeeting).Methods("POST")
	userRoutes.HandleFunc("/participants/{id}", participantController.UpdateStatus).Methods("PUT")
	userRoutes.HandleFunc("/participants/{id}", participantController.LeaveMeeting).Methods("DELETE")
	
	// 会议资料相关
	userRoutes.HandleFunc("/meetings/{meeting_id}/documents", meetingDocumentController.GetByMeeting).Methods("GET")
	userRoutes.HandleFunc("/documents/{id}/download", meetingDocumentController.Download).Methods("GET")
	
	// 公告相关
	userRoutes.HandleFunc("/announcements", announcementController.GetAll).Methods("GET")
	userRoutes.HandleFunc("/announcements/{id}", announcementController.GetByID).Methods("GET")
	
	// 会议提醒相关
	userRoutes.HandleFunc("/reminders", meetingReminderController.GetUserReminders).Methods("GET")
	userRoutes.HandleFunc("/reminders/{id}", meetingReminderController.GetByID).Methods("GET")

	// 管理员路由 - 需要管理员认证
	adminRoutes := api.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middleware.AuthMiddleware)
	adminRoutes.Use(middleware.AdminAuthMiddleware)
	
	// 管理员个人信息
	adminRoutes.HandleFunc("/profile", adminController.GetProfile).Methods("GET")
	adminRoutes.HandleFunc("/profile", adminController.UpdateProfile).Methods("PUT")
	adminRoutes.HandleFunc("/change-password", adminController.ChangePassword).Methods("POST")
	
	// 管理员管理（仅超级管理员）
	adminRoutes.HandleFunc("/admins", adminController.GetAll).Methods("GET")
	adminRoutes.HandleFunc("/admins", adminController.Create).Methods("POST")
	adminRoutes.HandleFunc("/admins/{id}", adminController.GetByID).Methods("GET")
	adminRoutes.HandleFunc("/admins/{id}", adminController.Update).Methods("PUT")
	adminRoutes.HandleFunc("/admins/{id}", adminController.Delete).Methods("DELETE")
	
	// 用户管理
	adminRoutes.HandleFunc("/users", adminController.GetAllUsers).Methods("GET")
	adminRoutes.HandleFunc("/users/{id}", adminController.GetUserByID).Methods("GET")
	adminRoutes.HandleFunc("/users", adminController.CreateUser).Methods("POST")
	adminRoutes.HandleFunc("/users/{id}", adminController.UpdateUser).Methods("PUT")
	adminRoutes.HandleFunc("/users/{id}", adminController.DeleteUser).Methods("DELETE")
	
	// 会议室管理
	adminRoutes.HandleFunc("/meeting-rooms", meetingRoomController.GetAll).Methods("GET")
	adminRoutes.HandleFunc("/meeting-rooms/{id}", meetingRoomController.GetByID).Methods("GET")
	adminRoutes.HandleFunc("/meeting-rooms", meetingRoomController.Create).Methods("POST")
	adminRoutes.HandleFunc("/meeting-rooms/{id}", meetingRoomController.Update).Methods("PUT")
	adminRoutes.HandleFunc("/meeting-rooms/{id}", meetingRoomController.Delete).Methods("DELETE")
	
	// 会议管理
	adminRoutes.HandleFunc("/meetings", meetingController.GetAll).Methods("GET")
	adminRoutes.HandleFunc("/meetings/{id}", meetingController.GetByID).Methods("GET")
	adminRoutes.HandleFunc("/meetings", meetingController.Create).Methods("POST")
	adminRoutes.HandleFunc("/meetings/{id}", meetingController.Update).Methods("PUT")
	adminRoutes.HandleFunc("/meetings/{id}", meetingController.Cancel).Methods("DELETE")
	
	// 参会信息管理
	adminRoutes.HandleFunc("/participants", participantController.GetAll).Methods("GET")
	adminRoutes.HandleFunc("/meetings/{meeting_id}/participants", participantController.GetByMeeting).Methods("GET")
	adminRoutes.HandleFunc("/participants", participantController.Create).Methods("POST")
	adminRoutes.HandleFunc("/participants/{id}", participantController.Update).Methods("PUT")
	adminRoutes.HandleFunc("/participants/{id}", participantController.Delete).Methods("DELETE")
	
	// 会议资料管理
	adminRoutes.HandleFunc("/documents", meetingDocumentController.GetAll).Methods("GET")
	adminRoutes.HandleFunc("/meetings/{meeting_id}/documents", meetingDocumentController.GetByMeeting).Methods("GET")
	adminRoutes.HandleFunc("/documents", meetingDocumentController.Create).Methods("POST")
	adminRoutes.HandleFunc("/documents/{id}", meetingDocumentController.Update).Methods("PUT")
	adminRoutes.HandleFunc("/documents/{id}", meetingDocumentController.Delete).Methods("DELETE")
	
	// 公告管理
	adminRoutes.HandleFunc("/announcements", announcementController.GetAll).Methods("GET")
	adminRoutes.HandleFunc("/announcements/{id}", announcementController.GetByID).Methods("GET")
	adminRoutes.HandleFunc("/announcements", announcementController.Create).Methods("POST")
	adminRoutes.HandleFunc("/announcements/{id}", announcementController.Update).Methods("PUT")
	adminRoutes.HandleFunc("/announcements/{id}", announcementController.Delete).Methods("DELETE")
	
	// 轮播图管理
	adminRoutes.HandleFunc("/carousels", carouselController.GetAll).Methods("GET")
	adminRoutes.HandleFunc("/carousels/{id}", carouselController.GetByID).Methods("GET")
	adminRoutes.HandleFunc("/carousels", carouselController.Create).Methods("POST")
	adminRoutes.HandleFunc("/carousels/{id}", carouselController.Update).Methods("PUT")
	adminRoutes.HandleFunc("/carousels/{id}", carouselController.Delete).Methods("DELETE")
	
	// 会议提醒管理
	adminRoutes.HandleFunc("/reminders", meetingReminderController.GetAll).Methods("GET")
	adminRoutes.HandleFunc("/meetings/{meeting_id}/reminders", meetingReminderController.GetByMeeting).Methods("GET")
	adminRoutes.HandleFunc("/reminders", meetingReminderController.Create).Methods("POST")
	adminRoutes.HandleFunc("/reminders/{id}", meetingReminderController.Update).Methods("PUT")
	adminRoutes.HandleFunc("/reminders/{id}", meetingReminderController.Delete).Methods("DELETE")
}
