<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { getAccounts, addAccount, updateAccount, deleteAccount } from '../../service/api'

const router = useRouter()

const accounts = ref<any[]>([])
const showModal = ref(false)
const isEdit = ref(false)

const currentAccount = ref({
  id: 0,
  name: '',
  type: 'bank',
  account_type: 'asset',
  balance: 0,
  exclude_amount: 0,
  account_number: '',
  remark: ''
})

const fetchAccounts = async () => {
  try {
    accounts.value = await getAccounts()
  } catch (error) {
    console.error('Failed to fetch accounts', error)
  }
}

onMounted(() => {
  fetchAccounts()
})

const openAdd = () => {
  isEdit.value = false
  currentAccount.value = {
    id: 0,
    name: '',
    type: 'bank',
    account_type: 'asset',
    balance: 0,
    exclude_amount: 0,
    account_number: '',
    remark: ''
  }
  showModal.value = true
}

const openEdit = (acc: any) => {
  isEdit.value = true
  currentAccount.value = { ...acc }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const saveAccount = async () => {
  if (!currentAccount.value.name.trim()) {
    alert('请输入账户名称')
    return
  }
  
  // Check for duplicate account name
  const isDuplicate = accounts.value.some(a => a.name === currentAccount.value.name && a.id !== currentAccount.value.id)
  if (isDuplicate) {
    alert('已存在相同名称的账户，请换一个名称！')
    return
  }

  // Data coercion
  currentAccount.value.balance = Number(currentAccount.value.balance)
  currentAccount.value.exclude_amount = Number(currentAccount.value.exclude_amount)

  try {
    if (isEdit.value) {
      await updateAccount(currentAccount.value.id, currentAccount.value)
    } else {
      await addAccount(currentAccount.value)
    }
    showModal.value = false
    fetchAccounts()
  } catch (error) {
    console.error('Failed to save account', error)
  }
}

const handleDelete = async () => {
  if (!confirm('确定要删除此账户吗？将会影响该账户的记账记录！')) return
  try {
    await deleteAccount(currentAccount.value.id)
    showModal.value = false
    fetchAccounts()
  } catch (error) {
    console.error('Failed to delete account', error)
  }
}

const netAsset = computed(() => {
  const b = Number(currentAccount.value.balance) || 0
  const e = Number(currentAccount.value.exclude_amount) || 0
  return (b - e).toFixed(2)
})
</script>

<template>
  <div class="subpage">
    <div class="header">
      <div class="header-left">
        <button @click="router.back()" class="back-btn">‹</button>
        <h2>账户管理</h2>
      </div>
      <button class="add-btn" @click="openAdd">+ 添加</button>
    </div>

    <!-- Account List -->
    <div class="account-list">
      <div class="account-card" v-for="acc in accounts" :key="acc.id" @click="openEdit(acc)">
        <div class="acc-icon">💳</div>
        <div class="acc-info">
          <div class="acc-name">
            {{ acc.name }}
            <span class="acc-tag" :class="acc.account_type === 'asset' ? 'tag-asset' : 'tag-liability'">
              {{ acc.account_type === 'asset' ? '资产' : '负债' }}
            </span>
          </div>
          <!-- <div class="acc-desc">{{ acc.account_number || acc.remark }}</div> -->
        </div>
        <div class="acc-balance">¥{{ acc.balance.toFixed(2) }}</div>
      </div>
    </div>

    <!-- Drawer/Modal -->
    <transition name="slide-up">
      <div class="drawer" v-if="showModal">
        <div class="drawer-header">
          <h3>{{ isEdit ? '编辑账户' : '添加账户' }}</h3>
          <button class="close-btn" @click="closeModal">✕</button>
        </div>

        <div class="drawer-content">
          <div class="form-group">
            <label>账户名称 <span class="required">*</span></label>
            <input type="text" v-model="currentAccount.name" placeholder="请输入账户名称" />
          </div>

          <div class="form-group">
            <label>账户类型 <span class="required">*</span> <span class="help-icon">?</span></label>
            <div class="type-cards">
              <div class="type-card" :class="{ active: currentAccount.account_type === 'asset' }" @click="currentAccount.account_type = 'asset'">
                <span class="radio-icon">{{ currentAccount.account_type === 'asset' ? '◉' : '◯' }}</span> 资产账户
              </div>
              <div class="type-card" :class="{ active: currentAccount.account_type === 'liability' }" @click="currentAccount.account_type = 'liability'">
                <span class="radio-icon">{{ currentAccount.account_type === 'liability' ? '◉' : '◯' }}</span> 负债账户
              </div>
            </div>
            <p class="help-text" v-if="currentAccount.account_type === 'asset'">资产账户余额通常为正数，表示实际持有金额</p>
            <p class="help-text" v-if="currentAccount.account_type === 'liability'">负债账户余额通常为正数，表示欠款金额（如信用卡）</p>
          </div>

          <div class="form-group">
            <label>账户余额 <span class="required">*</span></label>
            <div class="input-with-prefix">
              <span class="prefix">¥</span>
              <input type="number" v-model="currentAccount.balance" placeholder="0.00" />
            </div>
          </div>

          <div class="form-group">
            <label>不计入金额 <span class="help-icon">?</span></label>
            <div class="input-with-prefix">
              <span class="prefix">¥</span>
              <input type="number" v-model="currentAccount.exclude_amount" placeholder="0.00" />
            </div>
            <p class="help-text">此金额不计入净资产统计</p>
            
            <div class="formula-box">
              <div class="f-title">账户净资产 = 账户余额 - 不计入金额</div>
              <div class="f-calc">
                <span class="f-result">{{ netAsset }}</span> = ({{ (Number(currentAccount.balance) || 0).toFixed(2) }}) - ({{ (Number(currentAccount.exclude_amount) || 0).toFixed(2) }})
              </div>
            </div>
          </div>

          <div class="form-group">
            <label>卡号/账号</label>
            <input type="text" v-model="currentAccount.account_number" placeholder="请输入卡号（可选）" />
          </div>

          <div class="form-group">
            <label>备注</label>
            <input type="text" v-model="currentAccount.remark" placeholder="请输入备注" />
          </div>
        </div>

        <div class="drawer-footer">
          <button v-if="isEdit" class="btn btn-delete" @click="handleDelete">删除</button>
          <button class="btn btn-cancel" @click="closeModal">取消</button>
          <button class="btn btn-save" @click="saveAccount">{{ isEdit ? '保存' : '添加' }}</button>
        </div>
      </div>
    </transition>
    <transition name="fade">
      <div class="overlay" v-if="showModal" @click="closeModal"></div>
    </transition>
  </div>
</template>

<style scoped>
.subpage {
  min-height: 100vh;
  background: var(--color-bg);
  color: var(--color-text-primary);
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #1e1e1e;
  border-bottom: 1px solid #333;
}
.header-left {
  display: flex;
  align-items: center;
}
.back-btn {
  background: none;
  border: none;
  color: white;
  font-size: 24px;
  margin-right: 12px;
  cursor: pointer;
}
.header h2 {
  font-size: 18px;
  margin: 0;
}
.add-btn {
  background: #3b82f6;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 14px;
}

.account-list {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.account-card {
  background: #242424;
  border-radius: var(--radius-md);
  padding: 16px;
  display: flex;
  align-items: center;
  cursor: pointer;
  border: 1px solid #333;
  transition: border-color 0.3s;
}
.account-card:hover {
  border-color: #555;
}
.acc-icon {
  width: 40px;
  height: 40px;
  background: rgba(255,255,255,0.05);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  margin-right: 16px;
}
.acc-info {
  flex: 1;
}
.acc-name {
  font-size: 16px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
}
.acc-tag {
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 4px;
}
.tag-asset {
  background: rgba(74, 222, 128, 0.2);
  color: #4ade80;
}
.tag-liability {
  background: rgba(251, 146, 60, 0.2);
  color: #fb923c;
}
.acc-balance {
  font-size: 16px;
  font-weight: bold;
}

/* Drawer UI */
.overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 2000;
}
.drawer {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  width: 100%;
  max-width: 500px;
  background: #1c1c1e;
  z-index: 2001;
  display: flex;
  flex-direction: column;
  box-shadow: -5px 0 20px rgba(0,0,0,0.5);
}

.drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #333;
}
.drawer-header h3 {
  margin: 0;
  font-size: 18px;
}
.close-btn {
  background: none;
  border: none;
  color: #888;
  font-size: 20px;
  cursor: pointer;
}

.drawer-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}
.form-group {
  margin-bottom: 24px;
}
.form-group label {
  display: block;
  font-size: 14px;
  margin-bottom: 8px;
  color: #ccc;
}
.required {
  color: #ef4444;
}
.help-icon {
  display: inline-block;
  width: 14px;
  height: 14px;
  background: #444;
  color: #aaa;
  border-radius: 50%;
  text-align: center;
  line-height: 14px;
  font-size: 10px;
  cursor: help;
  margin-left: 4px;
}
.form-group input[type="text"],
.form-group input[type="number"] {
  width: 100%;
  background: #111;
  border: 1px solid #333;
  color: white;
  padding: 12px;
  border-radius: var(--radius-sm);
  font-size: 15px;
  box-sizing: border-box;
}
.form-group input:focus {
  outline: none;
  border-color: #3b82f6;
}

.type-cards {
  display: flex;
  gap: 12px;
  margin-bottom: 8px;
}
.type-card {
  flex: 1;
  background: #2a2a2a;
  border: 1px solid #333;
  border-radius: var(--radius-sm);
  padding: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 15px;
  transition: all 0.2s;
}
.type-card .radio-icon {
  margin-right: 8px;
  font-size: 18px;
  color: #666;
}
.type-card.active {
  background: rgba(59, 130, 246, 0.1);
  border-color: #3b82f6;
  color: #3b82f6;
}
.type-card.active .radio-icon {
  color: #3b82f6;
}
.help-text {
  font-size: 12px;
  color: #888;
  margin: 4px 0 0 0;
}

.input-with-prefix {
  position: relative;
  display: flex;
  align-items: center;
}
.input-with-prefix .prefix {
  position: absolute;
  left: 12px;
  color: #888;
}
.input-with-prefix input {
  padding-left: 28px !important;
}

.formula-box {
  background: #111;
  border-radius: var(--radius-sm);
  padding: 12px;
  margin-top: 12px;
}
.f-title {
  font-size: 12px;
  color: #888;
  margin-bottom: 6px;
}
.f-calc {
  font-size: 14px;
  color: #ccc;
  font-family: monospace;
}
.f-result {
  color: #3b82f6;
  font-weight: bold;
}

.drawer-footer {
  padding: 20px;
  border-top: 1px solid #333;
  display: flex;
  gap: 12px;
}
.btn {
  padding: 12px;
  border-radius: var(--radius-sm);
  font-size: 16px;
  cursor: pointer;
  flex: 1;
  text-align: center;
}
.btn-delete {
  flex: 0 0 auto;
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
  padding: 12px 20px;
}
.btn-cancel {
  background: transparent;
  color: white;
  border: 1px solid #444;
}
.btn-save {
  background: #3b82f6;
  color: white;
  border: none;
}

/* Transitions */
.slide-up-enter-active, .slide-up-leave-active {
  transition: transform 0.3s ease;
}
.slide-up-enter-from, .slide-up-leave-to {
  transform: translateX(100%);
}
@media (max-width: 768px) {
  .slide-up-enter-from, .slide-up-leave-to {
    transform: translateY(100%);
  }
}
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
