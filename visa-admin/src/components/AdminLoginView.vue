<template>
  <div class="admin-login">
    <div class="admin-login__card">
      <div class="admin-login__brand">
        <div class="admin-login__mark">V</div>
        <div>
          <div class="admin-login__title">管理后台</div>
        </div>
      </div>

      <el-form @submit.prevent>
        <el-form-item label="账号">
          <el-input v-model.trim="form.account" placeholder="请输入账号" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" show-password placeholder="请输入密码" @keyup.enter="submit" />
        </el-form-item>
        <el-button class="admin-login__btn" type="primary" :loading="loading" @click="submit">登录后台</el-button>
      </el-form>

    </div>
  </div>
</template>

<script setup>
import { reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '../api/client'
import { saveAdminSession } from '../auth'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['success', 'loading'])

const form = reactive({
  account: '',
  password: '',
})

async function submit() {
  if (!form.account || !form.password) {
    ElMessage.warning('请输入管理员账号和密码')
    return
  }
  emit('loading', true)
  try {
    const data = await api.loginAdmin({
      account: form.account,
      phone: form.account,
      password: form.password,
    })
    if (!data || !data.token) {
      throw new Error('登录响应异常')
    }
    saveAdminSession(data.token, data.user || null)
    ElMessage.success('管理员登录成功')
    emit('success', data.user || null)
  } catch (error) {
    ElMessage.error(error.message || '登录失败')
  } finally {
    emit('loading', false)
  }
}
</script>

<style scoped>
.admin-login {
  min-height: 100vh;
  background: linear-gradient(180deg, #eaf1fb 0%, #f6f8fc 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  box-sizing: border-box;
}

.admin-login__card {
  width: 100%;
  max-width: 440px;
  padding: 28px;
  border-radius: 22px;
  background: #fff;
  border: 1px solid #e5e7eb;
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.08);
}

.admin-login__brand {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 24px;
}

.admin-login__mark {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  background: #1d4ed8;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  font-weight: 900;
}

.admin-login__title {
  font-size: 22px;
  font-weight: 900;
  color: #111827;
}

.admin-login__sub,
.admin-login__hint {
  color: #6b7280;
  font-size: 13px;
  line-height: 1.6;
}

.admin-login__sub {
  margin-top: 4px;
}

.admin-login__btn {
  width: 100%;
  margin-top: 4px;
}

.admin-login__hint {
  margin-top: 16px;
}
</style>
