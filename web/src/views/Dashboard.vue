<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch, nextTick, onActivated, onDeactivated } from 'vue'
import { useRouter } from 'vue-router'
import { getDashboard, type DashboardData } from '../service/api'
import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'

const router = useRouter()

const selectedYear = ref(new Date().getFullYear())
const minYear = 2000
const maxYear = new Date().getFullYear()

const totalAssets = ref(0)
const netAssets = ref(0)
const yearIncome = ref(0)
const yearExpense = ref(0)
const yearBalance = ref(0)
const currentYear = ref(new Date().getFullYear())
const currentMonth = ref(new Date().getMonth() + 1)
const monthIncome = ref(0)
const monthExpense = ref(0)
const monthBalance = ref(0)
const monthlyStats = ref<any[]>([])

const chartRef = ref<HTMLElement | null>(null)
let trendChart: echarts.ECharts | null = null

const renderChart = () => {
  if (!chartRef.value) return
  if (!trendChart) {
    trendChart = echarts.init(chartRef.value)
  }

  const chartData = [...monthlyStats.value].reverse()
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'line', lineStyle: { type: 'dashed', color: '#666' } },
      backgroundColor: '#fff',
      textStyle: { color: '#333' },
      confine: true,
      triggerOn: 'mousemove|click'
    },
    legend: {
      data: ['收入', '支出', '结余'],
      bottom: 0,
      textStyle: { color: '#888' },
      icon: 'circle',
      itemWidth: 10,
      itemHeight: 10
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '12%',
      top: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: chartData.map(item => `${item.month}月`),
      axisLine: { lineStyle: { color: '#333' } },
      axisLabel: { color: '#888' },
      triggerEvent: true,
      boundaryGap: false
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: '#222', type: 'dashed' } },
      axisLabel: { color: '#888' }
    },
    series: [
      {
        name: '收入',
        type: 'line',
        data: chartData.map(item => item.income),
        itemStyle: { color: '#10b981' },
        lineStyle: { width: 3 },
        smooth: true,
        showSymbol: true,
        symbol: 'circle',
        symbolSize: 6,
        label: { show: true, position: 'top', color: '#10b981', formatter: (p: any) => p.value === 0 ? '' : p.value }
      },
      {
        name: '支出',
        type: 'line',
        data: chartData.map(item => item.expense),
        itemStyle: { color: '#ef4444' },
        lineStyle: { width: 3 },
        smooth: true,
        showSymbol: true,
        symbol: 'circle',
        symbolSize: 6,
        label: { show: true, position: 'top', color: '#ef4444', formatter: (p: any) => p.value === 0 ? '' : `-${p.value}` }
      }
    ]
  }
  trendChart.setOption(option)

  // Handle clicking on xAxis labels specifically for mobile
  trendChart.off('click')
  trendChart.on('click', (params) => {
    if (params.componentType === 'xAxis' && trendChart) {
      trendChart.dispatchAction({
        type: 'showTip',
        seriesIndex: 0,
        dataIndex: params.dataIndex || 0
      })
    }
  })
}

const handleResize = () => trendChart?.resize()

watch(monthlyStats, () => {
  nextTick(() => {
    renderChart()
  })
}, { deep: true })

const showLegend = ref(false)

const viewMode = ref<'list' | 'grid'>('list')
const toggleViewMode = () => {
  viewMode.value = viewMode.value === 'list' ? 'grid' : 'list'
}

const maxIncome = computed(() => {
  if (!monthlyStats.value.length) return -1
  const max = Math.max(...monthlyStats.value.map(s => s.income))
  return max > 0 ? max : -1
})
const maxExpense = computed(() => {
  if (!monthlyStats.value.length) return -1
  const max = Math.max(...monthlyStats.value.map(s => s.expense))
  return max > 0 ? max : -1
})
const maxBalance = computed(() => {
  if (!monthlyStats.value.length) return -1
  return Math.max(...monthlyStats.value.map(s => s.balance))
})
const minBalance = computed(() => {
  if (!monthlyStats.value.length) return 0
  const min = Math.min(...monthlyStats.value.map(s => s.balance))
  return min < 0 ? min : 0
})

const closeYearPicker = (e: MouseEvent) => {
  const picker = document.querySelector('.year-picker')
  const trigger = document.querySelector('.current-year')
  if (picker && !picker.contains(e.target as Node) && !trigger?.contains(e.target as Node)) {
    showYearPicker.value = false
  }
}

onMounted(() => {
  loadDashboard()
  document.addEventListener('click', closeYearPicker)
})

let isFirstMount = true
onActivated(() => {
  window.addEventListener('resize', handleResize)
  if (!isFirstMount) {
    loadDashboard()
  }
  isFirstMount = false
  nextTick(() => {
    trendChart?.resize()
  })
})

onDeactivated(() => {
  window.removeEventListener('resize', handleResize)
})

onUnmounted(() => {
  document.removeEventListener('click', closeYearPicker)
  if (trendChart) {
    trendChart.dispose()
    trendChart = null
  }
})

const handleAIClick = () => {
  ElMessage.warning('该功能处于开发中')
}

const prevYear = () => {
  if (selectedYear.value > minYear) {
    selectedYear.value--
    loadDashboard()
  }
}
const nextYear = () => {
  if (selectedYear.value < maxYear) {
    selectedYear.value++
    loadDashboard()
  }
}

const showYearPicker = ref(false)
const yearList = computed(() => {
  const years = []
  for (let y = maxYear; y >= minYear; y--) years.push(y)
  return years
})
const selectYear = (y: number) => {
  selectedYear.value = y
  showYearPicker.value = false
  loadDashboard()
}

async function loadDashboard() {
  try {
    const data: DashboardData = await getDashboard(selectedYear.value)
    totalAssets.value = data.total_assets
    netAssets.value = data.net_assets
    yearIncome.value = data.year_income
    yearExpense.value = data.year_expense
    yearBalance.value = data.year_balance
    currentYear.value = data.year
    currentMonth.value = data.month
    monthIncome.value = data.month_income
    monthExpense.value = data.month_expense
    monthBalance.value = data.month_balance
    // Server returns Jan -> current month, let's reverse it to show latest first
    monthlyStats.value = (data.monthly_stats && data.monthly_stats.length > 0)
      ? [...data.monthly_stats].reverse()
      : [{ month: data.month, income: data.month_income, expense: data.month_expense, balance: data.month_balance }]
  } catch (e) {
    console.error('Failed to load dashboard data', e)
  }
}
</script>

<template>
  <div class="dashboard-page">
    <div class="header">
      <div class="title">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" class="logo">
          <!-- Logo approximation -->
          <rect x="2" y="10" width="12" height="6" rx="2" fill="#2563EB" transform="rotate(-15 8 13)"/>
          <rect x="6" y="14" width="12" height="6" rx="2" fill="#10B981" transform="rotate(-15 12 17)"/>
          <rect x="10" y="6" width="12" height="6" rx="2" fill="#EF4444" transform="rotate(-15 16 9)"/>
        </svg>
        <h2>EasyAccounts</h2>
      </div>
    </div>

    <!-- 顶部资产卡片 -->
    <div class="card asset-card">
      <div class="asset-header">
        <span class="asset-label">
          总资产 
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="eye-icon"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
        </span>
        <button class="account-btn" @click="router.push('/settings/account')">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
          <span>账户</span>
        </button>
      </div>
      <div class="asset-amount">¥ {{ totalAssets.toFixed(2) }}</div>
    </div>

    <!-- 年度统计卡片 -->
    <div class="card stat-card">
      <div class="stat-header">
        <button class="icon-circle" @click="prevYear" :disabled="selectedYear <= minYear">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="5" y1="12" x2="19" y2="12"></line></svg>
        </button>
        <div class="current-year" @click="showYearPicker = !showYearPicker" style="cursor:pointer;position:relative;">
          {{ selectedYear }}年度
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" :style="showYearPicker ? 'transform:rotate(180deg);transition:.2s' : 'transition:.2s'"><polyline points="6 9 12 15 18 9"></polyline></svg>
        </div>
        <button class="icon-circle" @click="nextYear" :disabled="selectedYear >= maxYear" :style="selectedYear >= maxYear ? 'opacity:0.3;cursor:not-allowed' : ''">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
        </button>
      </div>

      <!-- 年份下拉选择器 -->
      <transition name="slide-down">
        <div v-if="showYearPicker" class="year-picker">
          <div
            v-for="y in yearList"
            :key="y"
            class="year-option"
            :class="{ active: y === selectedYear }"
            @click="selectYear(y)"
          >{{ y }}年</div>
        </div>
      </transition>
      
      <div class="stat-body">
        <div class="stat-item">
          <span class="label">收入</span>
          <span class="value income">{{ yearIncome.toFixed(2) }}</span>
        </div>
        <div class="stat-divider"></div>
        <div class="stat-item">
          <span class="label">支出</span>
          <span class="value expense">{{ yearExpense.toFixed(2) }}</span>
        </div>
        <div class="stat-divider"></div>
        <div class="stat-item">
          <span class="label">结余</span>
          <span class="value balance">{{ yearBalance.toFixed(2) }}</span>
        </div>
      </div>
    </div>

    <!-- 趋势图表卡片 -->
    <div class="card chart-card">
      <div class="chart-header">
        <h3>{{ selectedYear }}年度趋势</h3>
      </div>
      <div ref="chartRef" class="trend-chart"></div>
    </div>

    <!-- 月度概览卡片 -->
    <div class="card monthly-card">
      <div class="monthly-header">
        <h3>月度概览</h3>
        <div class="actions">
          <button class="action-btn" @click="showLegend = true">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="16" x2="12" y2="12"></line><line x1="12" y1="8" x2="12.01" y2="8"></line></svg>
          </button>
          <button class="action-btn" @click="toggleViewMode">
            <svg v-if="viewMode === 'grid'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
            <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="8" y1="6" x2="21" y2="6"></line><line x1="8" y1="12" x2="21" y2="12"></line><line x1="8" y1="18" x2="21" y2="18"></line><line x1="3" y1="6" x2="3.01" y2="6"></line><line x1="3" y1="12" x2="3.01" y2="12"></line><line x1="3" y1="18" x2="3.01" y2="18"></line></svg>
          </button>
        </div>
      </div>
      
      <div :class="['month-container', viewMode === 'grid' ? 'month-grid' : 'month-list']">
        <div :class="viewMode === 'grid' ? 'month-card-item' : 'month-row'" v-for="stat in monthlyStats" :key="stat.month">
          <div class="month-name">{{ stat.month }}月</div>
          
          <div class="month-stats" v-if="viewMode === 'list'">
            <span class="m-val" :class="[stat.income === maxIncome ? 'badge-income' : 'm-income']">+{{ stat.income.toFixed(2) }}</span>
            <span class="m-val" :class="[stat.expense === maxExpense ? 'badge-expense' : 'm-expense']">{{ stat.expense > 0 ? '-' : '' }}{{ stat.expense.toFixed(2) }}</span>
            <span class="m-val" :class="[
              stat.balance === maxBalance ? 'badge-balance' : 
              (stat.balance < 0 ? 'badge-negative' : 'm-balance')
            ]">{{ stat.balance > 0 ? '+' : '' }}{{ stat.balance.toFixed(2) }}</span>
          </div>

          <div class="month-stats-grid" v-else>
            <div class="stat-line">
              <span class="s-label">收</span>
              <span class="m-val" :class="[stat.income === maxIncome ? 'badge-income' : 'm-income']">+{{ stat.income.toFixed(2) }}</span>
            </div>
            <div class="stat-line">
              <span class="s-label">支</span>
              <span class="m-val" :class="[stat.expense === maxExpense ? 'badge-expense' : 'm-expense']">{{ stat.expense > 0 ? '-' : '' }}{{ stat.expense.toFixed(2) }}</span>
            </div>
            <div class="stat-line">
              <span class="s-label">结</span>
              <span class="m-val" :class="[
                stat.balance === maxBalance ? 'badge-balance' : 
                (stat.balance < 0 ? 'badge-negative' : 'm-balance')
              ]">{{ stat.balance > 0 ? '+' : '' }}{{ stat.balance.toFixed(2) }}</span>
            </div>
          </div>

          <div class="arrow" v-if="viewMode === 'list'">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg>
          </div>
        </div>
      </div>
    </div>

    <!-- 悬浮 AI 按钮 -->
    <button class="fab-ai" @click="handleAIClick">AI</button>

    <!-- 标记说明 Modal -->
    <transition name="fade">
      <div class="modal-overlay" v-if="showLegend" @click="showLegend = false">
        <div class="legend-modal" @click.stop>
          <div class="legend-header">
            <h3>标记说明</h3>
            <button class="close-btn" @click="showLegend = false">✕</button>
          </div>
          <div class="legend-body">
            <div class="legend-item">
              <div class="legend-badge badge-income">+{{ maxIncome > 0 ? maxIncome.toFixed(2) : '0.00' }}</div>
              <div class="legend-text">收入最高的月份</div>
            </div>
            <div class="legend-item">
              <div class="legend-badge badge-expense">{{ maxExpense > 0 ? '-' : '' }}{{ maxExpense > 0 ? maxExpense.toFixed(2) : '0.00' }}</div>
              <div class="legend-text">支出最高的月份</div>
            </div>
            <div class="legend-item">
              <div class="legend-badge badge-balance">{{ maxBalance > 0 ? '+' : '' }}{{ maxBalance !== -1 ? maxBalance.toFixed(2) : '0.00' }}</div>
              <div class="legend-text">结余最高的月份</div>
            </div>
            <div class="legend-item">
              <div class="legend-badge badge-negative">{{ minBalance < 0 ? minBalance.toFixed(2) : '-0.00' }}</div>
              <div class="legend-text">支出大于收入（负结余）</div>
            </div>
          </div>
          <div class="legend-footer">
            <button class="btn-primary" @click="showLegend = false">知道了</button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.header {
  padding-top: 10px;
  padding-bottom: 5px;
}

.title {
  display: flex;
  align-items: center;
  gap: 10px;
}
.title h2 {
  font-size: 20px;
  color: var(--color-text-primary);
  font-weight: bold;
  letter-spacing: 0.5px;
}
.logo {
  width: 26px;
  height: 26px;
}

.card {
  background: var(--color-card-bg);
  border-radius: var(--radius-lg);
  padding: 20px;
}

/* 顶部蓝色资产卡片 */
.asset-card {
  background: #60a5fa; /* Flat solid blue like the screenshot */
  color: white;
  position: relative;
  border-radius: var(--radius-lg);
  box-shadow: 0 4px 15px rgba(96, 165, 250, 0.2);
}
.asset-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}
.asset-label {
  font-size: 14px;
  opacity: 0.9;
  display: flex;
  align-items: center;
  gap: 8px;
}
.eye-icon {
  opacity: 0.8;
}
.account-btn {
  background: rgba(255, 255, 255, 0.25);
  border: none;
  border-radius: 12px;
  color: white;
  padding: 8px 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  font-size: 11px;
}
.asset-amount {
  font-size: 40px;
  font-weight: 800;
  margin-top: 16px;
  letter-spacing: -0.5px;
}

/* 年度统计卡片 */
.stat-card {
  padding: 24px 20px;
}
.stat-header {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}
.icon-circle {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 1px solid var(--color-text-secondary);
  background: transparent;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}
.current-year {
  background: var(--color-card-inner);
  padding: 6px 16px;
  border-radius: var(--radius-pill);
  font-weight: 500;
  font-size: 15px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.year-picker {
  background: #1a1a1e;
  border: 1px solid #333;
  border-radius: 12px;
  margin: 0 auto 16px;
  max-height: 200px;
  overflow-y: auto;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 4px;
  padding: 8px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.4);
}
.year-option {
  text-align: center;
  padding: 8px 4px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  color: var(--color-text-secondary);
  transition: background 0.15s, color 0.15s;
}
.year-option:hover {
  background: #2a2a2e;
  color: var(--color-text-primary);
}
.year-option.active {
  background: var(--color-primary);
  color: #111;
  font-weight: 700;
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
  align-items: center;
}
.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  flex: 1;
}
.stat-divider {
  width: 1px;
  height: 35px;
  background: var(--color-border);
}
.stat-item .label {
  color: var(--color-text-secondary);
  font-size: 13px;
  margin-bottom: 4px;
}
.stat-item .value {
  font-size: 16px;
  font-weight: 600;
}
.value.income { color: var(--color-income); }
.value.expense { color: var(--color-expense); }
.value.balance { color: var(--color-text-primary); }

/* Chart Card */
.chart-card {
  margin-bottom: 16px;
  overflow: hidden;
}
.chart-header {
  text-align: center;
  margin-bottom: 16px;
}
.chart-header h3 {
  font-size: 16px;
  color: #fff;
  font-weight: bold;
  margin: 0;
}
.trend-chart {
  width: 100%;
  height: 240px;
}

/* 月度概览 */
.monthly-card {
  padding: 20px;
}
.monthly-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.monthly-header h3 {
  font-size: 16px;
  font-weight: 600;
}
.actions {
  display: flex;
  gap: 8px;
}
.action-btn {
  background: var(--color-card-inner);
  border: none;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  cursor: pointer;
  color: var(--color-text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Monthly List */
.month-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.month-row {
  display: flex;
  align-items: center;
  background: var(--color-card-inner);
  padding: 12px 10px;
  border-radius: 12px;
  gap: 8px;
}
.month-name {
  font-weight: 700;
  font-size: 14px;
  width: 36px;
  white-space: nowrap;
}
.month-stats {
  display: flex;
  flex: 1;
  gap: 6px;
  justify-content: flex-end;
  align-items: center;
}
.m-val {
  font-size: 12px;
  padding: 4px 6px;
  border-radius: 6px;
  white-space: nowrap;
}
.m-income { color: var(--color-income); }
.m-expense { 
  background: var(--color-expense-bg); 
  color: white; 
}
.m-balance { 
  border: 1px solid var(--color-expense-border); 
  color: var(--color-primary);
}
.m-val.badge-income {
  background: #67c23a;
  color: #fff;
}
.m-val.badge-expense {
  background: #f87171;
  color: white;
}
.m-val.badge-balance {
  background: #60a5fa;
  color: white;
  border-color: #60a5fa;
}
.m-val.badge-negative {
  background: transparent;
  border: 1px solid #ef4444;
  color: #ef4444;
}
.arrow {
  color: var(--color-text-secondary);
  display: flex;
  align-items: center;
  margin-left: 8px;
}

/* Grid Layout for Months */
.month-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}
@media (max-width: 600px) {
  .month-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.month-card-item {
  background: #000;
  border-radius: 12px;
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  border: 1px solid #333;
  transition: transform 0.2s, box-shadow 0.2s;
  cursor: pointer;
}
.month-card-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.4);
}
.month-card-item .month-name {
  font-size: 18px;
  font-weight: bold;
  text-align: left;
  border-bottom: 1px solid #2a2a2c;
  padding-bottom: 8px;
}
.month-stats-grid {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.stat-line {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.s-label {
  font-size: 12px;
  color: #666;
}
.month-stats-grid .m-val {
  min-width: 60px;
  text-align: right;
  font-size: 12px;
  padding: 2px 6px;
}

/* FAB */
.fab-ai {
  position: fixed;
  bottom: 100px;
  right: 20px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: var(--color-primary);
  color: #111;
  border: none;
  font-weight: 800;
  font-size: 20px;
  box-shadow: var(--shadow-primary);
  cursor: pointer;
  z-index: 100;
  transition: transform 0.2s;
}
.fab-ai:hover {
  transform: scale(1.05);
}

/* Legend Modal */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
}
.legend-modal {
  background: #1c1c1e;
  border-radius: 12px;
  width: 90%;
  max-width: 440px;
  padding: 24px;
  border: 1px solid #333;
  box-shadow: 0 10px 40px rgba(0,0,0,0.5);
}
.legend-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.legend-header h3 {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
}
.close-btn {
  background: none;
  border: none;
  color: #888;
  font-size: 20px;
  cursor: pointer;
  padding: 0;
}
.close-btn:hover {
  color: white;
}
.legend-body {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 24px;
}
.legend-item {
  display: flex;
  align-items: center;
  gap: 20px;
  background: #2a2a2c;
  padding: 16px;
  border-radius: 8px;
}
.legend-badge {
  font-size: 13px;
  font-weight: bold;
  padding: 6px 16px;
  border-radius: 6px;
  text-align: center;
  min-width: 100px;
}
.badge-income {
  background: #67c23a; /* Green */
  color: #fff;
}
.badge-expense {
  background: #f87171; /* Red */
  color: white;
}
.badge-balance {
  background: #60a5fa; /* Blue */
  color: white;
}
.badge-negative {
  background: transparent;
  border: 1px solid #ef4444; /* Red border */
  color: #ef4444;
}
.legend-text {
  font-size: 14px;
  color: #eee;
}
.legend-footer {
  display: flex;
  justify-content: flex-end;
}
.btn-primary {
  background: #3b82f6;
  color: white;
  border: none;
  padding: 10px 24px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}
.btn-primary:hover {
  opacity: 0.9;
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
