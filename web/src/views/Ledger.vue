<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, onActivated } from 'vue'
import { getTransactions, addTransaction, deleteTransaction, getTransactionTypes, type Transaction, type TransactionType } from '../service/api'
import AddTransaction from '../components/AddTransaction.vue'

const now = new Date()
const selectedYear = ref(now.getFullYear())
const selectedMonth = ref(now.getMonth() + 1)
const monthIncome = ref(0)
const monthExpense = ref(0)
const monthBalance = ref(0)

const currentMonthStr = computed(() => `${selectedYear.value}年${selectedMonth.value}月`)

const transactions = ref<Transaction[]>([])
const tTypes = ref<TransactionType[]>([])
const isLoading = ref(true)

// 月份导航
const prevMonth = () => {
  if (selectedMonth.value === 1) {
    selectedMonth.value = 12
    selectedYear.value--
  } else {
    selectedMonth.value--
  }
  loadData()
}
const nextMonth = () => {
  const now = new Date()
  if (selectedYear.value === now.getFullYear() && selectedMonth.value === now.getMonth() + 1) return
  if (selectedMonth.value === 12) {
    selectedMonth.value = 1
    selectedYear.value++
  } else {
    selectedMonth.value++
  }
  loadData()
}
const isCurrentMonth = computed(() => {
  const n = new Date()
  return selectedYear.value === n.getFullYear() && selectedMonth.value === n.getMonth() + 1
})

// 日期选择器
const showPicker = ref(false)
const pickerYear = ref(selectedYear.value)
const monthNames = ['1月','2月','3月','4月','5月','6月','7月','8月','9月','10月','11月','12月']

const openPicker = () => {
  pickerYear.value = selectedYear.value
  showPicker.value = true
}
const pickerYearNext = () => {
  if (pickerYear.value < new Date().getFullYear()) pickerYear.value++
}
const selectMonthYear = (m: number) => {
  selectedYear.value = pickerYear.value
  selectedMonth.value = m
  showPicker.value = false
  loadData()
}
const isFutureMonth = (m: number) => {
  const n = new Date()
  return pickerYear.value > n.getFullYear() || (pickerYear.value === n.getFullYear() && m > n.getMonth() + 1)
}

const closePicker = (e: MouseEvent) => {
  const picker = document.querySelector('.month-picker')
  const trigger = document.querySelector('.current-month')
  if (picker && !picker.contains(e.target as Node) && !trigger?.contains(e.target as Node)) {
    showPicker.value = false
  }
}
onMounted(() => {
  loadData()
  document.addEventListener('click', closePicker)
})

let isFirstMount = true
onActivated(() => {
  if (!isFirstMount) {
    loadData()
  }
  isFirstMount = false
})
onUnmounted(() => {
  document.removeEventListener('click', closePicker)
})

const loadData = async () => {
  isLoading.value = true
  try {
    const [data, typesData] = await Promise.all([
      getTransactions(selectedYear.value, selectedMonth.value),
      getTransactionTypes()
    ])
    
    tTypes.value = typesData as TransactionType[]
    
    // Reset stats
    monthIncome.value = 0
    monthExpense.value = 0
    
    // Group transactions by date
    const grouped = (data as Transaction[]).reduce((acc: any, tx: Transaction) => {
      const recordDate = tx.record_date || ''
      const dateStr = recordDate.length >= 10 ? recordDate.substring(5, 10) : '未知'
      const date = dateStr !== '未知' ? dateStr.replace('-', '月') + '日' : '未知日期'
      
      if (!acc[date]) {
        acc[date] = []
      }
      
      const shortDate = dateStr !== '未知' ? dateStr.replace('-', '/') : '-'
      
      let dayOfWeek = ''
      if (recordDate.length >= 10) {
        const d = new Date(recordDate.substring(0, 10))
        const days = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
        dayOfWeek = days[d.getDay()]
      }
      
      // Determine flow
      let flow = ''
      if (tx.type === 'income') flow = 'inflow'
      else if (tx.type === 'expense') flow = 'outflow'
      else {
        const tType = tTypes.value.find(t => t.name === tx.type)
        if (tType) flow = tType.flow
      }
      
      const signedAmount = flow === 'outflow' ? -tx.amount : tx.amount
      
      acc[date].push({
        id: tx.id,
        shortDate: shortDate,
        dayOfWeek: dayOfWeek,
        title: tx.category?.name || '-',
        account: tx.account?.name || '-',
        remark: tx.remark || '-',
        amount: signedAmount,
        type: tx.type,
      })
      
      if (flow === 'inflow') monthIncome.value += tx.amount
      else if (flow === 'outflow') monthExpense.value += tx.amount
      
      return acc
    }, {})
    
    monthBalance.value = monthIncome.value - monthExpense.value
    
    // We store the grouped array in transactions instead of raw Transactions (quirk of this component)
    transactions.value = Object.keys(grouped).map(date => ({
      date,
      items: grouped[date]
    })) as any[]
  } catch (e) {
    // Errors handled by axios interceptor globally
  } finally {
    isLoading.value = false
  }
}

const showAddModal = ref(false)

const handleAddTransaction = async (data: any) => {
  try {
    const payload = {
      ...data,
      category_id: parseInt(data.categoryId) || 0,
      account_id: parseInt(data.accountId),
      amount: parseFloat(data.amount),
      record_date: new Date(data.recordDate).toISOString()
    }
    await addTransaction(payload)
    await loadData()
  } catch (e) {
    console.error('Failed to add transaction', e)
  }
}

const showDeleteConfirm = ref(false)
const deleteTargetId = ref<number | null>(null)

const handleDelete = (id: number) => {
  deleteTargetId.value = id
  showDeleteConfirm.value = true
}

const executeDelete = async () => {
  if (deleteTargetId.value === null) return
  try {
    await deleteTransaction(deleteTargetId.value)
    await loadData()
    showDeleteConfirm.value = false
    deleteTargetId.value = null
  } catch (e) {
    console.error('Failed to delete transaction', e)
  }
}

const cancelDelete = () => {
  showDeleteConfirm.value = false
  deleteTargetId.value = null
}
</script>

<template>
  <div class="ledger-page">
    <div class="header">
      <h2>明细</h2>
      <div class="header-actions">
        <button class="icon-btn">▽</button>
        <button class="icon-btn">📄</button>
        <button class="icon-btn">🔍</button>
      </div>
    </div>

    <!-- 月份与统计卡片 -->
    <div class="card month-card">
      <div class="month-selector">
        <button class="arrow-btn" @click="prevMonth">‹</button>
        <div class="current-month" @click="openPicker" style="cursor:pointer">
          {{ currentMonthStr }} <span :style="showPicker ? 'display:inline-block;transform:rotate(180deg);transition:.2s' : 'display:inline-block;transition:.2s'">⌄</span>
        </div>
        <button class="arrow-btn" @click="nextMonth">›</button>
      </div>

      <!-- 月份展开选择器 -->
      <transition name="slide-down">
        <div v-if="showPicker" class="month-picker">
          <div class="picker-year-nav">
            <button @click="pickerYear--">‹</button>
            <span>{{ pickerYear }}年</span>
            <button @click="pickerYear++">›</button>
          </div>
          <div class="picker-months">
            <div
              v-for="(m, idx) in monthNames" :key="idx"
              class="picker-month"
              :class="{ active: pickerYear === selectedYear && (idx+1) === selectedMonth, disabled: isFutureMonth(idx+1) }"
              @click="!isFutureMonth(idx+1) && selectMonthYear(idx+1)"
            >{{ m }}</div>
          </div>
        </div>
      </transition>
      
      <div class="stat-body">
        <div class="stat-item">
          <span class="label">收入</span>
          <span class="value income">+{{ monthIncome.toFixed(2) }}</span>
        </div>
        <div class="stat-item">
          <span class="label">支出</span>
          <span class="value expense">{{ monthExpense.toFixed(2) }}</span>
        </div>
        <div class="stat-item">
          <span class="label">结余</span>
          <span class="value balance">{{ monthBalance.toFixed(2) }}</span>
        </div>
      </div>
    </div>

    <!-- 流水列表 (Table view) -->
    <div v-if="transactions && transactions.length > 0" class="transaction-table">
      <!-- 表头 -->
      <div class="t-header">
        <div class="th-date">日期</div>
        <div class="th-cat">分类</div>
        <div class="th-acc">账户</div>
        <div class="th-remark">备注</div>
        <div class="th-amount">金额</div>
        <div class="th-actions">操作</div>
      </div>

      <!-- 数据组 -->
      <div v-for="group in transactions" :key="group.date" class="date-group">
        <div class="t-row" v-for="item in group.items" :key="item.id">
          <div class="td-date">
            <div>{{ item.shortDate }}</div>
            <div style="font-size: 11px; margin-top: 2px; opacity: 0.8;">{{ item.dayOfWeek }}</div>
          </div>
          <div class="td-cat">{{ item.title }}</div>
          <div class="td-acc">{{ item.account }}</div>
          <div class="td-remark">{{ item.remark }}</div>
          <div class="td-amount" :class="{ 'income': item.amount > 0, 'expense': item.amount < 0 }">
            {{ item.amount > 0 ? '+' : '' }}{{ item.amount.toFixed(2) }}
          </div>
          <div class="td-actions">
            <button class="action-btn btn-del" @click="handleDelete(item.id)">🗑</button>
          </div>
        </div>
      </div>
    </div>

    <el-empty v-else-if="!isLoading" description="暂无数据" />

    <!-- 添加按钮 -->
    <button class="fab-add" @click="showAddModal = true">+</button>

    <AddTransaction 
      :show="showAddModal" 
      @close="showAddModal = false" 
      @submit="handleAddTransaction" 
    />

    <!-- 删除确认弹窗 -->
    <transition name="fade">
      <div class="modal-overlay" v-if="showDeleteConfirm" @click.self="cancelDelete">
        <div class="confirm-modal">
          <h3 class="confirm-title">删除确认</h3>
          <p class="confirm-desc">确定要删除这条记账记录吗？</p>
          <div class="confirm-actions">
            <button class="btn-cancel" @click="cancelDelete">取消</button>
            <button class="btn-confirm" @click="executeDelete">确认</button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<style scoped>
.ledger-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding-bottom: 80px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 10px;
}
.header h2 {
  font-size: 24px;
  font-weight: bold;
}
.header-actions {
  display: flex;
  gap: 12px;
}
.icon-btn {
  background: var(--color-card-bg);
  border: none;
  width: 36px;
  height: 36px;
  border-radius: var(--radius-sm);
  box-shadow: var(--shadow-sm);
  cursor: pointer;
}

.card {
  background: var(--color-card-bg);
  border-radius: var(--radius-lg);
  padding: 20px;
  box-shadow: var(--shadow-sm);
}

.month-selector {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}
.arrow-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 1px solid #e8e8e8;
  background: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: var(--color-text-secondary);
}
.current-month {
  font-size: 18px;
  font-weight: bold;
}

.month-picker {
  background: var(--color-card-bg);
  border: 1px solid #333;
  border-radius: 12px;
  margin: 0 0 16px;
  padding: 12px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.4);
}
.picker-year-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  font-size: 15px;
  font-weight: 600;
}
.picker-year-nav button {
  background: none;
  border: none;
  font-size: 22px;
  color: var(--color-text-secondary);
  cursor: pointer;
  padding: 0 8px;
  line-height: 1;
}
.picker-year-nav button:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}
.picker-months {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 6px;
}
.picker-month {
  text-align: center;
  padding: 8px 4px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  color: var(--color-text-secondary);
  transition: background 0.15s, color 0.15s;
}
.picker-month:hover:not(.disabled) {
  background: var(--color-card-inner);
  color: var(--color-text-primary);
}
.picker-month.active {
  background: var(--color-primary);
  color: #111;
  font-weight: 700;
}
.picker-month.disabled {
  opacity: 0.3;
  cursor: not-allowed;
}
.slide-down-enter-active, .slide-down-leave-active {
  transition: all 0.2s ease;
}
.slide-down-enter-from, .slide-down-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

.stat-body {
  display: flex;
  justify-content: space-between;
  text-align: center;
}
.stat-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}
.stat-item .label {
  color: var(--color-text-secondary);
  font-size: 13px;
}
.stat-item .value {
  font-size: 16px;
  font-weight: 600;
}
.value.income { color: var(--color-income); }
.value.expense { color: var(--color-expense); }
.value.balance { color: var(--color-text-primary); }

/* Table View */
.transaction-table {
  background: white;
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  font-size: 14px;
}

.t-header {
  display: flex;
  padding: 12px 16px;
  border-bottom: 1px solid #e8e8e8;
  color: #888;
  font-weight: bold;
  background: white;
}
.th-date { width: 60px; }
.th-cat { flex: 1; }
.th-acc { flex: 1; }
.th-remark { flex: 1.5; }
.th-amount { width: 100px; text-align: right; }
.th-actions { width: 80px; text-align: right; }

.date-group {
  display: flex;
  flex-direction: column;
}
.group-header {
  background: #f5f5f5;
  padding: 8px 16px;
  font-weight: bold;
  color: #666;
  font-size: 13px;
}
.t-row {
  display: flex;
  padding: 16px;
  align-items: center;
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.2s ease, transform 0.2s ease;
}
.t-row:hover {
  background-color: rgba(128, 128, 128, 0.15);
}
.date-group:last-child .t-row:last-child {
  border-bottom: none;
}
.td-date { width: 60px; color: #666; }
.td-cat { flex: 1; color: #333; }
.td-acc { flex: 1; color: #888; }
.td-remark { flex: 1.5; color: #888; }
.td-amount { width: 100px; text-align: right; font-weight: bold; font-size: 15px; }
.td-amount.income { color: #4ade80; }
.td-amount.expense { color: #f87171; }

.td-actions {
  width: 80px;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
.action-btn {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 1px solid #ddd;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 14px;
  color: #888;
}
.btn-del {
  border-color: #fca5a5;
  color: #ef4444;
  background: #fef2f2;
}
.action-btn:hover {
  filter: brightness(0.95);
}

.fab-add {
  position: fixed;
  bottom: 90px;
  right: 20px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: var(--color-primary);
  color: white;
  border: none;
  font-weight: bold;
  font-size: 30px;
  box-shadow: var(--shadow-primary);
  cursor: pointer;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s;
}
.fab-add:hover {
  transform: scale(1.05);
}

/* Add dark mode support for table */
@media (prefers-color-scheme: dark) {
  .transaction-table {
    background: #1c1c1e;
  }
  .t-header {
    background: #1c1c1e;
    border-bottom: 1px solid #333;
  }
  .group-header {
    background: #242426;
    color: #aaa;
  }
  .t-row {
    border-bottom: 1px solid #333;
  }
  .t-row:hover {
    background-color: rgba(255, 255, 255, 0.05);
  }
  .td-cat { color: #eee; }
  .td-amount.expense { color: #ef4444; }
  .action-btn {
    background: #2a2a2c;
    border-color: #444;
  }
  .btn-del {
    background: rgba(239, 68, 68, 0.1);
    border-color: rgba(239, 68, 68, 0.3);
  }
}

/* Mobile Layout */
@media (max-width: 768px) {
  .t-header { 
    display: none; 
  }
  
  .t-row {
    display: grid;
    grid-template-columns: 46px 1fr auto;
    grid-template-areas:
      "date cat amount"
      "date acc amount"
      "date remark actions";
    gap: 2px 8px;
    padding: 12px 10px;
    align-items: center;
  }

  .td-date { 
    grid-area: date; 
    width: auto; 
    font-size: 13px; 
    color: #888;
    text-align: center;
  }
  
  .td-cat { 
    grid-area: cat; 
    width: auto; 
    font-size: 15px; 
    font-weight: 500; 
  }
  
  .td-acc { 
    grid-area: acc; 
    width: auto; 
    font-size: 12px; 
    color: #888;
  }
  
  .td-remark { 
    grid-area: remark; 
    width: auto; 
    font-size: 12px; 
    color: #888;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  
  .td-amount { 
    grid-area: amount; 
    width: auto; 
    text-align: right; 
    font-size: 16px;
    align-self: flex-start;
    padding-top: 2px;
  }
  
  .td-actions { 
    grid-area: actions; 
    width: auto; 
    justify-content: flex-end;
  }
}

/* 确认弹窗样式 */
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
