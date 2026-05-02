<template>
  <el-container class="user-layout">
    <el-header class="header">
      <div class="logo">
        <el-icon><OfficeBuilding /></el-icon>
        <span>会议室管理系统</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="nav-menu"
        mode="horizontal"
        router
      >
        <el-menu-item index="/user/home">
          <el-icon><House /></el-icon>
          <span>首页</span>
        </el-menu-item>
        <el-menu-item index="/user/meeting-rooms">
          <el-icon><OfficeBuilding /></el-icon>
          <span>会议室</span>
        </el-menu-item>
        <el-menu-item index="/user/meetings">
          <el-icon><Calendar /></el-icon>
          <span>我的会议</span>
        </el-menu-item>
        <el-menu-item index="/user/participations">
          <el-icon><User /></el-icon>
          <span>参会信息</span>
        </el-menu-item>
        <el-menu-item index="/user/announcements">
          <el-icon><Bell /></el-icon>
          <span>公告</span>
        </el-menu-item>
        <el-menu-item index="/user/reminders">
          <el-icon><AlarmClock /></el-icon>
          <span>会议提醒</span>
        </el-menu-item>
      </el-menu>
      <div class="user-info">
        <el-dropdown @command="handleCommand">
          <span class="user-dropdown">
            <el-avatar :size="32" icon="User" />
            <span class="username">{{ userStore.userInfo?.username || '用户' }}</span>
            <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">
                <el-icon><User /></el-icon>
                <span>个人信息</span>
              </el-dropdown-item>
              <el-dropdown-item command="change-password">
                <el-icon><Lock /></el-icon>
                <span>修改密码</span>
              </el-dropdown-item>
              <el-dropdown-item divided command="logout">
                <el-icon><SwitchButton /></el-icon>
                <span>退出登录</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>
    <el-main class="main">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const handleCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/user/profile')
      break
    case 'change-password':
      router.push('/user/change-password')
      break
    case 'logout':
      ElMessageBox.confirm('确定要退出登录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        userStore.logout()
        router.push('/login')
        ElMessage.success('已退出登录')
      }).catch(() => {})
      break
  }
}
</script>

<style scoped>
.user-layout {
  height: 100vh;
}

.header {
  display: flex;
  align-items: center;
  background-color: #409eff;
  padding: 0 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.logo {
  display: flex;
  align-items: center;
  color: #fff;
  font-size: 20px;
  font-weight: bold;
  margin-right: 40px;
}

.logo .el-icon {
  margin-right: 10px;
  font-size: 24px;
}

.nav-menu {
  flex: 1;
  border-bottom: none;
  background-color: transparent;
}

.nav-menu .el-menu-item {
  color: #fff;
  border-bottom: none;
}

.nav-menu .el-menu-item:hover,
.nav-menu .el-menu-item.is-active {
  background-color: rgba(255, 255, 255, 0.2);
  border-bottom: none;
}

.user-info {
  margin-left: auto;
}

.user-dropdown {
  display: flex;
  align-items: center;
  color: #fff;
  cursor: pointer;
}

.username {
  margin: 0 8px;
  font-size: 14px;
}

.main {
  background-color: #f5f7fa;
  padding: 20px;
  overflow-y: auto;
}
</style>
