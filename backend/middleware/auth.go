package middleware

import (
	"context"
	"net/http"
	"strings"

	"meetingmanage/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(w, "未提供身份验证令牌")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.UnauthorizedResponse(w, "身份验证令牌格式无效")
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			utils.UnauthorizedResponse(w, "身份验证令牌无效: "+err.Error())
			return
		}

		// 将用户信息添加到上下文中
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserContextKey).(*utils.Claims)
		if !ok {
			utils.UnauthorizedResponse(w, "无法获取用户信息")
			return
		}

		// 检查用户角色
		if claims.Role != "user" {
			utils.UnauthorizedResponse(w, "需要用户权限")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserContextKey).(*utils.Claims)
		if !ok {
			utils.UnauthorizedResponse(w, "无法获取用户信息")
			return
		}

		// 检查管理员角色
		if claims.Role != "admin" && claims.Role != "super_admin" {
			utils.UnauthorizedResponse(w, "需要管理员权限")
			return
		}

		next.ServeHTTP(w, r)
	})
}
