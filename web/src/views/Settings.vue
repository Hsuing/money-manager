<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const dataMenus = [
  { id: 'income-expense', title: '收支管理', desc: '管理收入和支出类型', icon: '⇆', color: '#60a5fa' },
  { id: 'account', title: '账户管理', desc: '管理银行卡、现金等账户', icon: '💳', color: '#60a5fa' },
  { id: 'category', title: '分类管理', desc: '管理收支分类', icon: '🏷️', color: '#60a5fa' },
  { id: 'template', title: '快记模板', desc: '快速记账预填模板', icon: '📄', color: '#60a5fa' },
  { id: 'schedule', title: '定时记账', desc: '周期性自动生成真实流水', icon: '⏰', color: '#60a5fa' }
]

const systemMenus = [
  { id: 'ai', title: 'AI+ 设置', desc: 'Token 统计与 MCP 状态', badge: '未配置', badgeClass: 'badge-warning', icon: '✨', color: '#4b5563' },
  { id: 'system', title: '系统设置', desc: '鉴权 / 邮件 / 提醒 / 备份 / 版本', icon: '⚙️', color: '#60a5fa' },
  { id: 'backup', title: '数据备份', desc: '备份与恢复数据库', icon: '📁', color: '#60a5fa' }
]

const handleMenuClick = (id: string) => {
  if (id === 'ai' || id === 'system') {
    router.push('/settings/system')
  } else {
    router.push(`/settings/${id}`)
  }
}

const showLogoutConfirm = ref(false)

const handleLogout = () => {
  showLogoutConfirm.value = true
}

const confirmLogout = () => {
  showLogoutConfirm.value = false
  localStorage.removeItem('token')
  router.push('/login')
}

const cancelLogout = () => {
  showLogoutConfirm.value = false
}
</script>

<template>
  <div class="settings-page">
    <div class="header">
      <h2>设置</h2>
    </div>

    <!-- 数据管理 -->
    <div class="section">
      <h3 class="section-title">数据管理</h3>
      <div class="menu-grid">
        <div 
          v-for="menu in dataMenus" 
          :key="menu.id" 
          class="menu-card"
          @click="handleMenuClick(menu.id)"
        >
          <div class="icon-box" :style="{ backgroundColor: menu.color }">
            {{ menu.icon }}
          </div>
          <div class="menu-info">
            <div class="menu-title">{{ menu.title }}</div>
            <div class="menu-desc">{{ menu.desc }}</div>
          </div>
          <div class="arrow">›</div>
        </div>
      </div>
    </div>

    <!-- 系统管理 -->
    <div class="section">
      <h3 class="section-title">系统管理</h3>
      <div class="menu-grid">
        <div 
          v-for="menu in systemMenus" 
          :key="menu.id" 
          class="menu-card"
          @click="handleMenuClick(menu.id)"
        >
          <div class="icon-box" :style="{ backgroundColor: menu.color, color: menu.color === '#4b5563' ? '#d1d5db' : '#fff' }">
            {{ menu.icon }}
          </div>
          <div class="menu-info">
            <div class="menu-title-row">
              <span class="menu-title">{{ menu.title }}</span>
              <span v-if="menu.badge" class="badge" :class="menu.badgeClass">{{ menu.badge }}</span>
            </div>
            <div class="menu-desc">{{ menu.desc }}</div>
          </div>
          <div class="arrow">›</div>
        </div>
      </div>
    </div>

    <!-- 其他 -->
    <div class="section">
      <h3 class="section-title">其他</h3>
      <div class="other-actions">
        <button class="action-btn" @click="() => {}">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><circle cx="12" cy="12" r="10"></circle><path d="M12 16v-4"></path><path d="M12 8h.01"></path></svg>
          关于
        </button>
        <button class="action-btn" @click="() => {}">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M4 22h16a2 2 0 0 0 2-2V4a2 2 0 0 0-2-2H8a2 2 0 0 0-2 2v16a2 2 0 0 1-2 2Zm0 0a2 2 0 0 1-2-2v-9c0-1.1.9-2 2-2h2"></path><path d="M18 14h-8"></path><path d="M15 18h-5"></path><path d="M10 6h8v4h-8V6Z"></path></svg>
          公告
        </button>
        <button class="action-btn logout-btn" @click="handleLogout">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="icon"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg>
          退出登录
        </button>
      </div>
    </div>
  </div>

  <!-- 退出确认弹窗 -->
  <transition name="fade">
    <div class="modal-overlay" v-if="showLogoutConfirm" @click.self="cancelLogout">
      <div class="confirm-modal">
        <h3 class="confirm-title">退出登录</h3>
        <p class="confirm-desc">确定要退出登录吗？</p>
        <div class="confirm-actions">
          <button class="btn-cancel" @click="cancelLogout">取消</button>
          <button class="btn-confirm" @click="confirmLogout">确认</button>
        </div>
      </div>
    </div>
  </transition>
</template>

<style scoped>
.settings-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.header {
  padding-top: 10px;
}
.header h2 {
  font-size: 24px;
  font-weight: bold;
}

.section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.section-title {
  font-size: 14px;
  color: var(--color-text-secondary);
  font-weight: 500;
  margin-left: 4px;
}

.menu-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 12px;
}

.menu-card {
  background: var(--color-card-bg);
  border-radius: var(--radius-md);
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 16px;
  cursor: pointer;
  border: 1px solid var(--color-border);
  transition: transform 0.2s, background-color 0.2s;
}
.menu-card:active {
  transform: scale(0.98);
  background: #222;
}

.icon-box {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: white;
  flex-shrink: 0;
}

.menu-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.menu-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
}
.menu-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--color-text-primary);
}
.menu-desc {
  font-size: 12px;
  color: var(--color-text-secondary);
}

.badge {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 500;
}
.badge-warning {
  background: rgba(217, 119, 6, 0.15);
  color: #fbbf24;
  border: 1px solid rgba(217, 119, 6, 0.3);
}

.arrow {
  color: var(--color-text-secondary);
  font-size: 20px;
}

.other-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--color-card-bg);
  border: 1px solid var(--color-border);
  color: var(--color-text-primary);
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}
.action-btn:hover {
  background: rgba(255, 255, 255, 0.05);
}
.action-btn .icon {
  width: 16px;
  height: 16px;
}

.logout-btn {
  color: #ef4444;
  border-color: rgba(239, 68, 68, 0.3);
  background: rgba(239, 68, 68, 0.05);
}
.logout-btn:hover {
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.5);
}
</style>

<style>
/* \u9000\u51fa\u786e\u8ba4\u5f39\u7a97\uff08\u975e scoped\uff0c\u907f\u514d\u88ab\u7236\u7ea7\u5bb9\u5668\u622a\u65ad\uff09 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.65);
  z-index: 3000;
  display: flex;
  align-items: center;
  justify-content: center;
}
.confirm-modal {
  background: #1c1c1e;
  border-radius: 16px;
  width: 280px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}
.confirm-title {
  font-size: 17px;
  font-weight: 600;
  text-align: center;
  color: #fff;
  padding: 20px 20px 8px;
  margin: 0;
}
.confirm-desc {
  font-size: 14px;
  color: #aaa;
  text-align: center;
  padding: 0 20px 20px;
  margin: 0;
}
.confirm-actions {
  display: flex;
  border-top: 1px solid #333;
}
.btn-cancel,
.btn-confirm {
  flex: 1;
  padding: 14px;
  border: none;
  background: transparent;
  font-size: 16px;
  cursor: pointer;
  transition: background 0.15s;
}
.btn-cancel {
  color: #aaa;
  border-right: 1px solid #333;
}
.btn-cancel:hover {
  background: rgba(255,255,255,0.05);
}
.btn-confirm {
  color: #3b82f6;
  font-weight: 500;
}
.btn-confirm:hover {
  background: rgba(59,130,246,0.1);
}
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
