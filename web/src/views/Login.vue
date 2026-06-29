<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { login, register, getCaptcha } from '../service/api'
import { ElMessage } from 'element-plus'

const router = useRouter()

const isRegisterMode = ref(false)

const form = ref({
  username: '',
  password: '',
  confirmPassword: ''
})

const errorMsg = ref('')
const isLoading = ref(false)

const captchaId = ref('')
const captchaImg = ref('')
const captchaVal = ref('')

const fetchCaptcha = async () => {
  try {
    const res: any = await getCaptcha()
    if (res && res.captchaId) {
      captchaId.value = res.captchaId
      captchaImg.value = res.captchaImg
    }
  } catch (err) {
    console.error('Failed to fetch captcha', err)
  }
}

const toggleMode = () => {
  isRegisterMode.value = !isRegisterMode.value
  errorMsg.value = ''
  form.value.password = ''
  form.value.confirmPassword = ''
  captchaVal.value = ''
  if (isRegisterMode.value) {
    fetchCaptcha()
  }
}

const handleSubmit = async () => {
  errorMsg.value = ''
  
  if (isRegisterMode.value) {
    if (form.value.password !== form.value.confirmPassword) {
      errorMsg.value = '两次输入的密码不一致'
      return
    }
    if (form.value.password.length < 6) {
      errorMsg.value = '密码至少需要6位'
      return
    }
  }

  isLoading.value = true
  
  try {
    if (isRegisterMode.value) {
      await register({ 
        username: form.value.username, 
        password: form.value.password,
        captchaId: captchaId.value,
        captchaVal: captchaVal.value
      })
      ElMessage.success('注册成功，请登录！')
      toggleMode()
    } else {
      const res: any = await login({ username: form.value.username, password: form.value.password })
      if (res && res.token) {
        localStorage.setItem('token', res.token)
        router.push('/dashboard')
      }
    }
  } catch (err: any) {
    errorMsg.value = err.response?.data?.msg || err.response?.data?.error || err.message || (isRegisterMode.value ? '注册失败' : '登录失败，请检查用户名和密码')
    if (isRegisterMode.value) {
      fetchCaptcha()
      captchaVal.value = ''
    }
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="login-container">
    <!-- Background subtle shapes -->
    <div class="login-bg">
      <div class="bg-shape shape-red"></div>
      <div class="bg-shape shape-green"></div>
      <div class="bg-shape shape-pill p1"></div>
      <div class="bg-shape shape-pill p2"></div>
      <div class="bg-shape shape-pill p3"></div>
    </div>

    <!-- Left Side: Branding -->
    <div class="brand-section">
      <div class="brand-content">
        <div class="logo-wrapper">
          <img src="/logo.png" alt="Logo" class="logo-img" />
          <h1 class="logo-text">EasyAccounts</h1>
        </div>
        <div class="slogan-wrapper">
          <p class="slogan-en">Where flow goes, life shows</p>
          <p class="slogan-cn">流水知去处，岁月自成诗</p>
        </div>
      </div>
    </div>

    <!-- Right Side: Login Form -->
    <div class="form-section">
      <div class="login-card">
        <a href="#" class="register-link" @click.prevent="toggleMode">{{ isRegisterMode ? '返回登录' : '注册账户' }}</a>
        <h2>{{ isRegisterMode ? '创建账户' : '欢迎回来' }}</h2>
        <p class="subtitle">{{ isRegisterMode ? '加入我们，开始记账之旅' : '请登录您的账户' }}</p>

        <form @submit.prevent="handleSubmit" class="login-form">
          <div class="input-group">
            <span class="icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
            </span>
            <input 
              type="text" 
              v-model="form.username" 
              @input="form.username = form.username.replace(/[^a-zA-Z0-9]/g, '')"
              placeholder="请输入用户名" 
              required
            />
          </div>

          <div class="input-group">
            <span class="icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
              </svg>
            </span>
            <input 
              type="password" 
              v-model="form.password" 
              placeholder="请输入密码" 
              required
            />
          </div>

          <transition name="slide-up">
            <div v-if="isRegisterMode" style="display: flex; flex-direction: column; gap: 20px; width: 100%; margin-bottom: 24px;">
              <div class="input-group" style="margin-bottom: 0;">
                <span class="icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
                  </svg>
                </span>
                <input 
                  type="password" 
                  v-model="form.confirmPassword" 
                  placeholder="请再次输入密码" 
                  required
                />
              </div>

              <div class="captcha-group">
                <div class="input-group" style="flex: 1; margin-bottom: 0;">
                  <span class="icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                      <line x1="9" y1="3" x2="9" y2="21"></line>
                    </svg>
                  </span>
                  <input 
                    type="text" 
                    v-model="captchaVal" 
                    placeholder="图形验证码" 
                    required
                  />
                </div>
                <img v-if="captchaImg" :src="captchaImg" @click="fetchCaptcha" class="captcha-img" alt="captcha" title="点击刷新" />
              </div>
            </div>
          </transition>

          <div class="actions" v-if="!isRegisterMode">
            <a href="#" class="forgot-pwd">忘记密码?</a>
          </div>

          <div v-if="errorMsg" class="error-msg">
            {{ errorMsg }}
          </div>

          <button type="submit" class="login-btn" :disabled="isLoading">
            {{ isLoading ? '请稍候...' : (isRegisterMode ? '注 册' : '登 录') }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  min-height: 100vh;
  width: 100vw;
  background-color: #080808;
  color: #ffffff;
  position: relative;
  z-index: 1;
}

.login-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
  z-index: -1;
}

/* --- Left Brand Section --- */
.brand-section {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  background: transparent;
  overflow: hidden;
}

.brand-content {
  position: relative;
  z-index: 10;
  padding: 0 40px;
}

.logo-wrapper {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.logo-img {
  width: 48px;
  height: 48px;
  object-fit: contain;
}

.logo-text {
  font-size: 32px;
  font-weight: 700;
  letter-spacing: 0.5px;
  color: #ffffff;
  margin: 0;
}

.slogan-wrapper {
  padding-left: 64px; /* Align with text, offset by logo width + gap */
}

.slogan-en {
  font-size: 18px;
  font-style: italic;
  color: #8c8f96;
  margin: 0 0 12px 0;
  font-family: 'Inter', sans-serif;
}

.slogan-cn {
  font-size: 14px;
  color: #6c6f76;
  letter-spacing: 1px;
  margin: 0;
}

/* Exact background elements based on new screenshot */
.bg-shape {
  position: absolute;
}

.shape-red {
  width: 120vw;
  height: 100vw;
  background: #180d0e;
  top: -60vw;
  right: -50vw;
  transform: rotate(20deg);
}

.shape-green {
  width: 120vw;
  height: 100vw;
  background: #0d1510;
  bottom: -60vw;
  left: -50vw;
  transform: rotate(20deg);
}

.shape-pill {
  background-color: #10161d;
  border-radius: 100px;
  height: 60px;
  transform: rotate(-20deg);
}

.p1 {
  width: 250px;
  top: 30%;
  left: -40px;
}

.p2 {
  width: 400px;
  top: 45%;
  left: -60px;
}

.p3 {
  width: 320px;
  top: 60%;
  left: -50px;
}

/* --- Right Form Section --- */
.form-section {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: transparent;
  position: relative;
}

.form-section::before {
  content: '';
  position: absolute;
  top: 0; right: 0; bottom: 0; left: 0;
  background: radial-gradient(circle at top right, rgba(40, 20, 20, 0.05), transparent 40%);
  pointer-events: none;
}

.login-card {
  position: relative;
  width: 100%;
  max-width: 440px;
  background-color: #1e1e1e;
  border-radius: 20px;
  padding: 48px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
  z-index: 10;
}

.register-link {
  position: absolute;
  top: 48px;
  right: 48px;
  color: #8c8f96;
  text-decoration: none;
  font-size: 14px;
  font-weight: 500;
  transition: color 0.2s;
}

.register-link:hover {
  color: #ffffff;
}

.login-card h2 {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 8px 0;
  color: #ffffff;
}

.subtitle {
  font-size: 14px;
  color: #8c8f96;
  margin: 0 0 32px 0;
}

.input-group {
  position: relative;
  margin-bottom: 24px;
}

.input-group .icon {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: #6b7280;
  width: 18px;
  height: 18px;
  display: flex;
}

.input-group input {
  width: 100%;
  background-color: #141414;
  border: 1px solid #2a2a2a;
  border-radius: 8px;
  padding: 14px 16px 14px 44px;
  color: #ffffff;
  font-size: 15px;
  outline: none;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.input-group input::placeholder {
  color: #4b515d;
}

.input-group input:focus {
  background-color: rgba(255, 255, 255, 0.05);
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.actions {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 24px;
}

.forgot-pwd {
  color: #3b82f6;
  font-size: 13px;
  text-decoration: none;
  transition: color 0.2s;
}

.forgot-pwd:hover {
  color: #60a5fa;
  text-decoration: underline;
}

.login-btn {
  width: 100%;
  background-color: #3b82f6;
  color: #ffffff;
  border: none;
  border-radius: 12px;
  padding: 14px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s ease, transform 0.1s ease;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.login-btn:hover {
  background-color: #4b8df8;
}

.login-btn:active:not(:disabled) {
  transform: translateY(1px);
}

.login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.error-msg {
  color: #ef4444;
  font-size: 13px;
  margin-bottom: 16px;
  text-align: center;
}

.captcha-group {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.captcha-img {
  height: 50px;
  border-radius: 12px;
  cursor: pointer;
  border: 1px solid rgba(255, 255, 255, 0.1);
  transition: opacity 0.2s;
}

.captcha-img:hover {
  opacity: 0.8;
}

/* Responsive adjustments */
@media (max-width: 900px) {
  .login-container {
    flex-direction: column;
  }
  
  .brand-section {
    flex: none;
    padding: 60px 20px;
  }
  
  .form-section {
    padding: 40px 20px;
    align-items: flex-start;
  }
  
  .login-card {
    padding: 32px 24px;
  }
  
  .register-link {
    top: 32px;
    right: 24px;
  }
}
</style>
