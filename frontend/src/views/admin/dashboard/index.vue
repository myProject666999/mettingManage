<template>
  <div class="admin-dashboard">
    <!-- 统计卡片 -->
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card shadow="hover" class="stats-card">
          <div class="stats-content">
            <div class="stats-icon blue">
              <el-icon><User /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-value">{{ stats.users }}</div>
              <div class="stats-label">用户总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stats-card">
          <div class="stats-content">
            <div class="stats-icon green">
              <el-icon><OfficeBuilding /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-value">{{ stats.meetingRooms }}</div>
              <div class="stats-label">会议室数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stats-card">
          <div class="stats-content">
            <div class="stats-icon orange">
              <el-icon><Calendar /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-value">{{ stats.meetings }}</div>
              <div class="stats-label">会议总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stats-card">
          <div class="stats-content">
            <div class="stats-icon red">
              <el-icon><AlarmClock /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-value">{{ stats.reminders }}</div>
              <div class="stats-label">待提醒</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 快捷操作 -->
    <el-card class="quick-actions-card" shadow="never" style="margin-top: 20px;">
      <template #header>
        <span class="card-title">快捷操作</span>
      </template>
      <el-row :gutter="20">
        <el-col :span="4">
          <div class="action-item" @click="navigateTo('/admin/meeting-rooms')">
            <el-icon size="30" color="#409eff"><OfficeBuilding /></el-icon>
            <span>会议室管理</span>
          </div>
        </el-col>
        <el-col :span="4">
          <div class="action-item" @click="navigateTo('/admin/meetings')">
            <el-icon size="30" color="#67c23a"><Calendar /></el-icon>
            <span>会议管理</span>
          </div>
        </el-col>
        <el-col :span="4">
          <div class="action-item" @click="navigateTo('/admin/users')">
            <el-icon size="30" color="#e6a23c"><User /></el-icon>
            <span>用户管理</span>
          </div>
        </el-col>
        <el-col :span="4">
          <div class="action-item" @click="navigateTo('/admin/announcements')">
            <el-icon size="30" color="#f56c6c"><Bell /></el-icon>
            <span>公告管理</span>
          </div>
        </el-col>
        <el-col :span="4">
          <div class="action-item" @click="navigateTo('/admin/carousels')">
            <el-icon size="30" color="#909399"><Picture /></el-icon>
            <span>轮播图管理</span>
          </div>
        </el-col>
        <el-col :span="4">
          <div class="action-item" @click="navigateTo('/admin/admins')">
            <el-icon size="30" color="#409eff"><UserFilled /></el-icon>
            <span>管理员管理</span>
          </div>
        </el-col>
      </el-row>
    </el-card>
    
    <!-- 最近会议 -->
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <span class="card-title">最近会议</span>
          </template>
          <el-table :data="recentMeetings" v-if="recentMeetings.length > 0">
            <el-table-column prop="title" label="会议标题" />
            <el-table-column prop="room.name" label="会议室" />
            <el-table-column prop="startTime" label="开始时间">
              <template #default="scope">
                {{ formatDate(scope.row.startTime) }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ getStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-else description="暂无会议" />
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <span class="card-title">最新公告</span>
          </template>
          <div v-if="recentAnnouncements.length > 0">
            <div
              v-for="announcement in recentAnnouncements"
              :key="announcement.id"
              class="announcement-item"
            >
              <div class="announcement-title">
                <el-tag v-if="announcement.isTop" type="danger" size="small">置顶</el-tag>
                <span>{{ announcement.title }}</span>
              </div>
              <span class="announcement-time">{{ formatDate(announcement.createdAt) }}</span>
            </div>
          </div>
          <el-empty v-else description="暂无公告" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import adminApi from '@/api/admin'
import dayjs from 'dayjs'

const router = useRouter()

const stats = reactive({
  users: 0,
  meetingRooms: 0,
  meetings: 0,
  reminders: 0
})

const recentMeetings = ref<any[]>([])
const recentAnnouncements = ref<any[]>([])

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const getStatusType = (status: string) => {
  const typeMap: Record<string, string> = {
    scheduled: 'primary',
    ongoing: 'success',
    completed: 'info',
    cancelled: 'danger'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    scheduled: '待开始',
    ongoing: '进行中',
    completed: '已结束',
    cancelled: '已取消'
  }
  return textMap[status] || '未知'
}

const navigateTo = (path: string) => {
  router.push(path)
}

const fetchData = async () => {
  try {
    // 获取用户列表
    const usersRes = await adminApi.getUsers()
    stats.users = usersRes?.length || 0
    
    // 获取会议室列表
    const roomsRes = await adminApi.getMeetingRooms()
    stats.meetingRooms = roomsRes?.length || 0
    
    // 获取会议列表
    const meetingsRes = await adminApi.getMeetings()
    stats.meetings = meetingsRes?.length || 0
    recentMeetings.value = (meetingsRes || []).slice(0, 5)
    
    // 获取提醒列表
    const remindersRes = await adminApi.getReminders()
    stats.reminders = (remindersRes || []).filter((r: any) => !r.isSent).length
    
    // 获取公告列表
    const announcementsRes = await adminApi.getAnnouncements()
    recentAnnouncements.value = (announcementsRes || []).slice(0, 5)
  } catch (error) {
    console.error('获取数据失败:', error)
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.admin-dashboard {
  max-width: 100%;
}

.stats-card {
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.stats-card:hover {
  transform: translateY(-5px);
}

.stats-content {
  display: flex;
  align-items: center;
}

.stats-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
}

.stats-icon.blue {
  background-color: rgba(64, 158, 255, 0.1);
  color: #409eff;
}

.stats-icon.green {
  background-color: rgba(103, 194, 58, 0.1);
  color: #67c23a;
}

.stats-icon.orange {
  background-color: rgba(230, 162, 60, 0.1);
  color: #e6a23c;
}

.stats-icon.red {
  background-color: rgba(245, 108, 108, 0.1);
  color: #f56c6c;
}

.stats-icon .el-icon {
  font-size: 30px;
}

.stats-info {
  flex: 1;
}

.stats-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.stats-label {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.card-title {
  font-size: 16px;
  font-weight: bold;
}

.quick-actions-card .action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  cursor: pointer;
  border-radius: 8px;
  transition: background-color 0.2s;
}

.quick-actions-card .action-item:hover {
  background-color: #f5f7fa;
}

.quick-actions-card .action-item span {
  margin-top: 10px;
  font-size: 14px;
  color: #606266;
}

.announcement-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #ebeef5;
}

.announcement-item:last-child {
  border-bottom: none;
}

.announcement-title {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
}

.announcement-title span {
  font-size: 14px;
  color: #303133;
}

.announcement-time {
  font-size: 12px;
  color: #909399;
  flex-shrink: 0;
}
</style>
