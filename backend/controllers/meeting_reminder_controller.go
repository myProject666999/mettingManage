package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"meetingmanage/database"
	"meetingmanage/middleware"
	"meetingmanage/models"
	"meetingmanage/utils"

	"gorm.io/gorm"
)

type MeetingReminderController struct{}

func (mrc *MeetingReminderController) GetAll(w http.ResponseWriter, r *http.Request) {
	var reminders []models.MeetingReminder
	if err := database.DB.Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		Order("reminder_time DESC").
		Find(&reminders).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取会议提醒列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", reminders)
}

func (mrc *MeetingReminderController) GetUserReminders(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok {
		utils.UnauthorizedResponse(w, "无法获取用户信息")
		return
	}

	// 获取用户参加的会议
	var participants []models.Participant
	if err := database.DB.Where("user_id = ?", claims.UserID).
		Preload("Meeting").Preload("Meeting.Reminders").
		Find(&participants).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取参会记录失败")
		return
	}

	// 收集所有提醒
	var reminders []models.MeetingReminder
	for _, p := range participants {
		for _, reminder := range p.Meeting.Reminders {
			reminder.Meeting = p.Meeting
			reminders = append(reminders, reminder)
		}
	}

	// 按提醒时间排序
	for i := 0; i < len(reminders); i++ {
		for j := i + 1; j < len(reminders); j++ {
			if reminders[i].ReminderTime.After(reminders[j].ReminderTime) {
				reminders[i], reminders[j] = reminders[j], reminders[i]
			}
		}
	}

	utils.SuccessResponse(w, "获取成功", reminders)
}

func (mrc *MeetingReminderController) GetByMeeting(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议ID
	meetingIDStr := r.URL.Query().Get("meeting_id")
	if meetingIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议ID参数")
		return
	}

	meetingID, err := strconv.ParseUint(meetingIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议ID格式无效")
		return
	}

	var reminders []models.MeetingReminder
	if err := database.DB.Where("meeting_id = ?", uint(meetingID)).
		Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		Order("reminder_time ASC").
		Find(&reminders).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取会议提醒列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", reminders)
}

func (mrc *MeetingReminderController) GetByID(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议提醒ID
	reminderIDStr := r.URL.Query().Get("id")
	if reminderIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议提醒ID参数")
		return
	}

	reminderID, err := strconv.ParseUint(reminderIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议提醒ID格式无效")
		return
	}

	var reminder models.MeetingReminder
	if err := database.DB.Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		First(&reminder, uint(reminderID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议提醒不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议提醒失败")
		}
		return
	}

	utils.SuccessResponse(w, "获取成功", reminder)
}

func (mrc *MeetingReminderController) Create(w http.ResponseWriter, r *http.Request) {
	var reminder models.MeetingReminder
	if err := json.NewDecoder(r.Body).Decode(&reminder); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	// 检查会议是否存在
	var meeting models.Meeting
	if err := database.DB.First(&meeting, reminder.MeetingID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议信息失败")
		}
		return
	}

	// 检查提醒时间是否有效
	if reminder.ReminderTime.Before(time.Now()) {
		utils.BadRequestResponse(w, "提醒时间不能早于当前时间")
		return
	}

	// 检查提醒时间是否晚于会议开始时间
	if reminder.ReminderTime.After(meeting.StartTime) {
		utils.BadRequestResponse(w, "提醒时间不能晚于会议开始时间")
		return
	}

	// 创建会议提醒
	if err := database.DB.Create(&reminder).Error; err != nil {
		utils.InternalServerErrorResponse(w, "创建会议提醒失败")
		return
	}

	// 加载关联信息
	database.DB.Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		First(&reminder, reminder.ID)

	utils.SuccessResponse(w, "创建成功", reminder)
}

func (mrc *MeetingReminderController) Update(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议提醒ID
	reminderIDStr := r.URL.Query().Get("id")
	if reminderIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议提醒ID参数")
		return
	}

	reminderID, err := strconv.ParseUint(reminderIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议提醒ID格式无效")
		return
	}

	var updateData struct {
		MeetingID    uint      `json:"meeting_id"`
		ReminderType string    `json:"reminder_type"`
		ReminderTime time.Time `json:"reminder_time"`
		Message      string    `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	var reminder models.MeetingReminder
	if err := database.DB.First(&reminder, uint(reminderID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议提醒不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议提醒失败")
		}
		return
	}

	if updateData.MeetingID > 0 {
		// 检查会议是否存在
		var meeting models.Meeting
		if err := database.DB.First(&meeting, updateData.MeetingID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.NotFoundResponse(w, "会议不存在")
			} else {
				utils.InternalServerErrorResponse(w, "获取会议信息失败")
			}
			return
		}
		reminder.MeetingID = updateData.MeetingID
	}

	if updateData.ReminderType != "" {
		reminder.ReminderType = updateData.ReminderType
	}

	if !updateData.ReminderTime.IsZero() {
		// 检查提醒时间是否有效
		if updateData.ReminderTime.Before(time.Now()) {
			utils.BadRequestResponse(w, "提醒时间不能早于当前时间")
			return
		}

		// 检查提醒时间是否晚于会议开始时间
		var meeting models.Meeting
		if err := database.DB.First(&meeting, reminder.MeetingID).Error; err == nil {
			if updateData.ReminderTime.After(meeting.StartTime) {
				utils.BadRequestResponse(w, "提醒时间不能晚于会议开始时间")
				return
			}
		}

		reminder.ReminderTime = updateData.ReminderTime
	}

	if updateData.Message != "" {
		reminder.Message = updateData.Message
	}

	if err := database.DB.Save(&reminder).Error; err != nil {
		utils.InternalServerErrorResponse(w, "更新会议提醒失败")
		return
	}

	// 加载关联信息
	database.DB.Preload("Meeting").Preload("Meeting.Room").Preload("Meeting.Organizer").
		First(&reminder, reminder.ID)

	utils.SuccessResponse(w, "更新成功", reminder)
}

func (mrc *MeetingReminderController) Delete(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取会议提醒ID
	reminderIDStr := r.URL.Query().Get("id")
	if reminderIDStr == "" {
		utils.BadRequestResponse(w, "缺少会议提醒ID参数")
		return
	}

	reminderID, err := strconv.ParseUint(reminderIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "会议提醒ID格式无效")
		return
	}

	var reminder models.MeetingReminder
	if err := database.DB.First(&reminder, uint(reminderID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "会议提醒不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取会议提醒失败")
		}
		return
	}

	if err := database.DB.Delete(&reminder).Error; err != nil {
		utils.InternalServerErrorResponse(w, "删除会议提醒失败")
		return
	}

	utils.SuccessResponse(w, "删除成功", nil)
}
