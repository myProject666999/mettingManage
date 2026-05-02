<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <el-icon size="40" color="#409eff"><OfficeBuilding /></el-icon>
        <h1>会议室管理系统</h1>
        <p>Meeting Management System</p>
      </div>
      
      <el-tabs v-model="activeTab" class="login-tabs">
        <el-tab-pane label="用户登录" name="user">
          <el-form
            ref="userFormRef"
            :model="userForm"
            :rules="userRules"
            class="login-form"
          >
            <el-form-item prop="username">
              <el-input
                v-model="userForm.username"
                placeholder="请输入用户名"
                prefix-icon="User"
                size="large"
              />
            </el-form-item>
            <el-form-item prop="password">
              <el-input
                v-model="userForm.password"
                type="password"
                placeholder="请输入密码"
                prefix-icon="Lock"
                size="large"
                show-password
                @keyup.enter="handleUserLogin"
              />
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                size="large"
                :loading="userLoading"
                class="login-btn"
                @click="handleUserLogin"
              >
                登 录
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <el-tab-pane label="管理员登录" name="admin">
          <el-form
            ref="adminFormRef"
            :model="adminForm"
            :rules="adminRules"
            class="login-form"
          >
            <el-form-item prop="username">
              <el-input
                v-model="adminForm.username"
                placeholder="请输入管理员用户名"
                prefix-icon="UserFilled"
                size="large"
              />
            </el-form-item>
            <el-form-item prop="password">
              <el-input
                v-model="adminForm.password"
                type="password"
                placeholder="请输入密码"
                prefix-icon="Lock"
                size="large"
                show-password
                @keyup.enter="handleAdminLogin"
              />
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                size="large"
                :loading="adminLoading"
                class="login-btn"
                @click="handleAdminLogin"
              >
                登 录
              </el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
      
      <div class="register-link" v-if="activeTab === 'user'">
        还没有账号？
        <router-link to="/register">立即注册</router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import userApi from '@/api/user'
import adminApi from '@/api/admin'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const activeTab = ref('user')

// 用户登录表单
const userFormRef = ref<FormInstance>()
const userForm = reactive({
  username: '',
  password: ''
})
const userRules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}
const userLoading = ref(false)

// 管理员登录表单
const adminFormRef = ref<FormInstance>()
const adminForm = reactive({
  username: '',
  password: ''
})
const adminRules: FormRules = {
  username: [{ required: true, message: '请输入管理员用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}
const adminLoading = ref(false)

// 处理用户登录
const handleUserLogin = async () => {
  if (!userFormRef.value) return
  
  await userFormRef.value.validate(async (valid) => {
    if (valid) {
      userLoading.value = true
      try {
        const res = await userApi.login(userForm.username, userForm.password)
        userStore.login(res.token, res.user, 'user')
        ElMessage.success('登录成功')
        
        const redirect = route.query.redirect as string || '/user/home'
        router.push(redirect)
      } catch (error) {
        console.error('登录失败:', error)
      } finally {
        userLoading.value = false
      }
    }
  })
}

// 处理管理员登录
const handleAdminLogin = async () => {
  if (!adminFormRef.value) return
  
  await adminFormRef.value.validate(async (valid) => {
    if (valid) {
      adminLoading.value = true
      try {
        const res = await adminApi.login(adminForm.username, adminForm.password)
        userStore.login(res.token, res.admin, res.admin.role)
        ElMessage.success('登录成功')
        
        const redirect = route.query.redirect as string || '/admin/dashboard'
        router.push(redirect)
      } catch (error) {
        console.error('登录失败:', error)
      } finally {
        adminLoading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-box {
  width: 400px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 15px 35px rgba(50, 50, 93, 0.1), 0 5px 15px rgba(0, 0, 0, 0.07);
  padding: 40px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h1 {
  margin: 15px 0 10px 0;
  font-size: 24px;
  color: #303133;
}

.login-header p {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.login-tabs {
  margin-bottom: 20px;
}

.login-form {
  margin-top: 20px;
}

.login-btn {
  width: 100%;
}

.register-link {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: #909399;
}

.register-link a {
  color: #409eff;
  text-decoration: none;
}

.register-link a:hover {
  text-decoration: underline;
}
</style>
