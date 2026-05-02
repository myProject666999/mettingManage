package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"meetingmanage/database"
	"meetingmanage/middleware"
	"meetingmanage/models"
	"meetingmanage/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminController struct{}

func (ac *AdminController) Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	// 查找管理员
	var admin models.Admin
	if err := database.DB.Where("username = ?", loginData.Username).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.UnauthorizedResponse(w, "用户名或密码错误")
		} else {
			utils.InternalServerErrorResponse(w, "登录失败")
		}
		return
	}

	// 检查密码
	if !admin.CheckPassword(loginData.Password) {
		utils.UnauthorizedResponse(w, "用户名或密码错误")
		return
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(admin.ID, admin.Username, admin.Role)
	if err != nil {
		utils.InternalServerErrorResponse(w, "令牌生成失败")
		return
	}

	utils.SuccessResponse(w, "登录成功", map[string]interface{}{
		"admin": admin,
		"token": token,
	})
}

func (ac *AdminController) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok {
		utils.UnauthorizedResponse(w, "无法获取管理员信息")
		return
	}

	var admin models.Admin
	if err := database.DB.First(&admin, claims.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "管理员不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取管理员信息失败")
		}
		return
	}

	utils.SuccessResponse(w, "获取成功", admin)
}

func (ac *AdminController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok {
		utils.UnauthorizedResponse(w, "无法获取管理员信息")
		return
	}

	var updateData struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	var admin models.Admin
	if err := database.DB.First(&admin, claims.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "管理员不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取管理员信息失败")
		}
		return
	}

	// 检查邮箱是否已被其他管理员使用
	if updateData.Email != "" && updateData.Email != admin.Email {
		var existingAdmin models.Admin
		if err := database.DB.Where("email = ? AND id != ?", updateData.Email, admin.ID).First(&existingAdmin).Error; err == nil {
			utils.BadRequestResponse(w, "邮箱已被使用")
			return
		}
		admin.Email = updateData.Email
	}

	if err := database.DB.Save(&admin).Error; err != nil {
		utils.InternalServerErrorResponse(w, "更新管理员信息失败")
		return
	}

	utils.SuccessResponse(w, "更新成功", admin)
}

func (ac *AdminController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.UserContextKey).(*utils.Claims)
	if !ok {
		utils.UnauthorizedResponse(w, "无法获取管理员信息")
		return
	}

	var passwordData struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&passwordData); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	if len(passwordData.NewPassword) < 6 {
		utils.BadRequestResponse(w, "新密码长度不能少于6位")
		return
	}

	var admin models.Admin
	if err := database.DB.First(&admin, claims.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "管理员不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取管理员信息失败")
		}
		return
	}

	// 检查旧密码
	if !admin.CheckPassword(passwordData.OldPassword) {
		utils.BadRequestResponse(w, "旧密码错误")
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalServerErrorResponse(w, "密码加密失败")
		return
	}

	admin.Password = string(hashedPassword)
	if err := database.DB.Save(&admin).Error; err != nil {
		utils.InternalServerErrorResponse(w, "更新密码失败")
		return
	}

	utils.SuccessResponse(w, "密码修改成功", nil)
}

func (ac *AdminController) GetAll(w http.ResponseWriter, r *http.Request) {
	var admins []models.Admin
	if err := database.DB.Find(&admins).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取管理员列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", admins)
}

func (ac *AdminController) GetByID(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取管理员ID
	adminIDStr := r.URL.Query().Get("id")
	if adminIDStr == "" {
		utils.BadRequestResponse(w, "缺少管理员ID参数")
		return
	}

	adminID, err := strconv.ParseUint(adminIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "管理员ID格式无效")
		return
	}

	var admin models.Admin
	if err := database.DB.First(&admin, uint(adminID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "管理员不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取管理员信息失败")
		}
		return
	}

	utils.SuccessResponse(w, "获取成功", admin)
}

func (ac *AdminController) Create(w http.ResponseWriter, r *http.Request) {
	var admin models.Admin
	if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	// 检查用户名是否已存在
	var existingAdmin models.Admin
	if err := database.DB.Where("username = ?", admin.Username).First(&existingAdmin).Error; err == nil {
		utils.BadRequestResponse(w, "用户名已存在")
		return
	}

	// 检查邮箱是否已存在
	if err := database.DB.Where("email = ?", admin.Email).First(&existingAdmin).Error; err == nil {
		utils.BadRequestResponse(w, "邮箱已被注册")
		return
	}

	// 创建管理员
	if err := database.DB.Create(&admin).Error; err != nil {
		utils.InternalServerErrorResponse(w, "管理员创建失败")
		return
	}

	utils.SuccessResponse(w, "创建成功", admin)
}

func (ac *AdminController) Update(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取管理员ID
	adminIDStr := r.URL.Query().Get("id")
	if adminIDStr == "" {
		utils.BadRequestResponse(w, "缺少管理员ID参数")
		return
	}

	adminID, err := strconv.ParseUint(adminIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "管理员ID格式无效")
		return
	}

	var updateData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	var admin models.Admin
	if err := database.DB.First(&admin, uint(adminID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "管理员不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取管理员信息失败")
		}
		return
	}

	// 检查用户名是否已被其他管理员使用
	if updateData.Username != "" && updateData.Username != admin.Username {
		var existingAdmin models.Admin
		if err := database.DB.Where("username = ? AND id != ?", updateData.Username, admin.ID).First(&existingAdmin).Error; err == nil {
			utils.BadRequestResponse(w, "用户名已被使用")
			return
		}
		admin.Username = updateData.Username
	}

	// 检查邮箱是否已被其他管理员使用
	if updateData.Email != "" && updateData.Email != admin.Email {
		var existingAdmin models.Admin
		if err := database.DB.Where("email = ? AND id != ?", updateData.Email, admin.ID).First(&existingAdmin).Error; err == nil {
			utils.BadRequestResponse(w, "邮箱已被使用")
			return
		}
		admin.Email = updateData.Email
	}

	if updateData.Role != "" {
		admin.Role = updateData.Role
	}

	if err := database.DB.Save(&admin).Error; err != nil {
		utils.InternalServerErrorResponse(w, "更新管理员信息失败")
		return
	}

	utils.SuccessResponse(w, "更新成功", admin)
}

func (ac *AdminController) Delete(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取管理员ID
	adminIDStr := r.URL.Query().Get("id")
	if adminIDStr == "" {
		utils.BadRequestResponse(w, "缺少管理员ID参数")
		return
	}

	adminID, err := strconv.ParseUint(adminIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "管理员ID格式无效")
		return
	}

	var admin models.Admin
	if err := database.DB.First(&admin, uint(adminID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "管理员不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取管理员信息失败")
		}
		return
	}

	if err := database.DB.Delete(&admin).Error; err != nil {
		utils.InternalServerErrorResponse(w, "删除管理员失败")
		return
	}

	utils.SuccessResponse(w, "删除成功", nil)
}

func (ac *AdminController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		utils.InternalServerErrorResponse(w, "获取用户列表失败")
		return
	}

	utils.SuccessResponse(w, "获取成功", users)
}

func (ac *AdminController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取用户ID
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		utils.BadRequestResponse(w, "缺少用户ID参数")
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "用户ID格式无效")
		return
	}

	var user models.User
	if err := database.DB.First(&user, uint(userID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "用户不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取用户信息失败")
		}
		return
	}

	utils.SuccessResponse(w, "获取成功", user)
}

func (ac *AdminController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		utils.BadRequestResponse(w, "用户名已存在")
		return
	}

	// 检查邮箱是否已存在
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		utils.BadRequestResponse(w, "邮箱已被注册")
		return
	}

	// 创建用户
	if err := database.DB.Create(&user).Error; err != nil {
		utils.InternalServerErrorResponse(w, "用户创建失败")
		return
	}

	utils.SuccessResponse(w, "创建成功", user)
}

func (ac *AdminController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取用户ID
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		utils.BadRequestResponse(w, "缺少用户ID参数")
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "用户ID格式无效")
		return
	}

	var updateData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		FullName string `json:"full_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.BadRequestResponse(w, "请求体解析失败")
		return
	}

	var user models.User
	if err := database.DB.First(&user, uint(userID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "用户不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取用户信息失败")
		}
		return
	}

	// 检查用户名是否已被其他用户使用
	if updateData.Username != "" && updateData.Username != user.Username {
		var existingUser models.User
		if err := database.DB.Where("username = ? AND id != ?", updateData.Username, user.ID).First(&existingUser).Error; err == nil {
			utils.BadRequestResponse(w, "用户名已被使用")
			return
		}
		user.Username = updateData.Username
	}

	// 检查邮箱是否已被其他用户使用
	if updateData.Email != "" && updateData.Email != user.Email {
		var existingUser models.User
		if err := database.DB.Where("email = ? AND id != ?", updateData.Email, user.ID).First(&existingUser).Error; err == nil {
			utils.BadRequestResponse(w, "邮箱已被使用")
			return
		}
		user.Email = updateData.Email
	}

	if updateData.Phone != "" {
		user.Phone = updateData.Phone
	}

	if updateData.FullName != "" {
		user.FullName = updateData.FullName
	}

	if err := database.DB.Save(&user).Error; err != nil {
		utils.InternalServerErrorResponse(w, "更新用户信息失败")
		return
	}

	utils.SuccessResponse(w, "更新成功", user)
}

func (ac *AdminController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// 从URL参数中获取用户ID
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		utils.BadRequestResponse(w, "缺少用户ID参数")
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(w, "用户ID格式无效")
		return
	}

	var user models.User
	if err := database.DB.First(&user, uint(userID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(w, "用户不存在")
		} else {
			utils.InternalServerErrorResponse(w, "获取用户信息失败")
		}
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		utils.InternalServerErrorResponse(w, "删除用户失败")
		return
	}

	utils.SuccessResponse(w, "删除成功", nil)
}
