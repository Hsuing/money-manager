<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { getTransactions, getTransactionTypes, type Transaction, type TransactionType } from '../service/api'
import * as echarts from 'echarts'

const router = useRouter()
const activeTab = ref('expense')
const activeFilter = ref('currentMonth')

const filters = [
  { id: 'currentMonth', label: '\u5f53\u6708' },
  { id: 'lastMonth',    label: '\u4e0a\u6708' },
  { id: 'last3Months',  label: '\u8fd13\u6708' },
  { id: 'last6Months',  label: '\u8fd16\u6708' },
  { id: 'lastYear',     label: '\u8fd11\u5e74' },
  { id: 'thisYear',     label: '\u5f53\u5e74' }
]

const transactions = ref<Transaction[]>([])
const tTypes = ref<TransactionType[]>([])
const isLoading = ref(true)

const customDateRange = ref<[Date, Date] | null>(null)

watch(activeFilter, (newVal) => {
  if (newVal !== 'custom') {
    customDateRange.value = null
  }
})

watch(customDateRange, (newVal) => {
  if (newVal && newVal.length === 2) {
    activeFilter.value = 'custom'
  }
})

// 根据 filter ID 计算开始、结束日期
const getDateRange = (filterId: string): { start: Date; end: Date } => {
  const now = new Date()
  const y = now.getFullYear()
  const m = now.getMonth()
  switch (filterId) {
    case 'custom':
      if (customDateRange.value && customDateRange.value.length === 2) {
        const [start, end] = customDateRange.value
        return { start, end: new Date(end.getFullYear(), end.getMonth(), end.getDate() + 1) }
      }
      return { start: new Date(y, m, 1), end: new Date(y, m + 1, 1) }
    case 'currentMonth':
      return { start: new Date(y, m, 1), end: new Date(y, m + 1, 1) }
    case 'lastMonth':
      return { start: new Date(y, m - 1, 1), end: new Date(y, m, 1) }
    case 'last3Months':
      return { start: new Date(y, m - 2, 1), end: new Date(y, m + 1, 1) }
    case 'last6Months':
      return { start: new Date(y, m - 5, 1), end: new Date(y, m + 1, 1) }
    case 'lastYear':
      return { start: new Date(y - 1, m + 1, 1), end: new Date(y, m + 1, 1) }
    case 'thisYear':
      return { start: new Date(y, 0, 1), end: new Date(y + 1, 0, 1) }
    default:
      return { start: new Date(y, m, 1), end: new Date(y, m + 1, 1) }
  }
}

// \u5f53\u524d\u65f6\u6bb5\u5185\u7684\u6d41\u6c34
const filteredTransactions = computed(() => {
  const { start, end } = getDateRange(activeFilter.value)
  return transactions.value.filter(tx => {
    const d = new Date(tx.record_date)
    return d >= start && d < end
  })
})

// \u6839\u636e\u4ea4\u6613\u7c7b\u578b\u5224\u65adflow\u65b9\u5411
const getFlow = (tx: Transaction): string => {
  if (tx.type === 'income') return 'inflow'
  if (tx.type === 'expense') return 'outflow'
  const t = tTypes.value.find(t => t.name === tx.type)
  return t?.flow || ''
}

const incomeTotal = computed(() =>
  filteredTransactions.value
    .filter(tx => getFlow(tx) === 'inflow')
    .reduce((sum, tx) => sum + tx.amount, 0)
)

const expenseTotal = computed(() =>
  filteredTransactions.value
    .filter(tx => getFlow(tx) === 'outflow')
    .reduce((sum, tx) => sum + tx.amount, 0)
)

// \u5206\u7c7b\u5360\u6bd4\uff08\u6309\u5f53\u524dtab + \u65f6\u6bb5\u8fc7\u6ee4\uff09
const analysisData = computed(() => {
  const isIncome = activeTab.value === 'income'
  const filtered = filteredTransactions.value.filter(tx =>
    isIncome ? getFlow(tx) === 'inflow' : getFlow(tx) === 'outflow'
  )
  const total = filtered.reduce((sum, tx) => sum + tx.amount, 0)
  if (total === 0) return []

  const grouped = filtered.reduce((acc: Record<string, number>, tx: Transaction) => {
    const cat = tx.category?.name || '\u5176\u4ed6'
    if (!acc[cat]) acc[cat] = 0
    acc[cat] += tx.amount
    return acc
  }, {})

  return Object.entries(grouped)
    .map(([title, amount]) => ({
      title,
      amount: amount as number,
      percentage: Math.round(((amount as number) / total) * 100)
    }))
    .sort((a, b) => b.amount - a.amount)
})

const viewMode = ref<'list' | 'chart'>('chart')
const chartRef = ref<HTMLElement | null>(null)
let pieChart: echarts.ECharts | null = null

const renderChart = () => {
  if (viewMode.value !== 'chart' || !chartRef.value || analysisData.value.length === 0) {
    if (pieChart) {
      pieChart.dispose()
      pieChart = null
    }
    return
  }
  if (!pieChart) {
    pieChart = echarts.init(chartRef.value)
  }

  const chartData = analysisData.value.map(item => ({
    name: item.title,
    value: item.amount
  }))

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: ¥{c} ({d}%)',
      backgroundColor: '#fff',
      textStyle: { color: '#333' },
      confine: true
    },
    legend: {
      bottom: '0%',
      left: 'center',
      textStyle: { color: '#888' },
      type: 'scroll'
    },
    series: [
      {
        name: activeTab.value === 'income' ? '收入分类' : '支出分类',
        type: 'pie',
        radius: ['35%', '55%'],
        avoidLabelOverlap: true,
        itemStyle: {
          borderRadius: 8,
          borderColor: '#111',
          borderWidth: 2
        },
        label: {
          show: true,
          formatter: '{b}\n{d}%',
          color: '#ccc'
        },
        labelLine: {
          show: true,
          length: 8,
          length2: 8,
          smooth: true,
          lineStyle: {
            color: '#666'
          }
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        },
        data: chartData
      }
    ]
  }
  
  pieChart.setOption(option)
}

const handleResize = () => pieChart?.resize()

watch([analysisData, viewMode], () => {
  nextTick(() => {
    renderChart()
  })
}, { deep: true })

onMounted(async () => {
  window.addEventListener('resize', handleResize)
  isLoading.value = true
  try {
    const [data, types] = await Promise.all([getTransactions(), getTransactionTypes()])
    transactions.value = data as Transaction[]
    tTypes.value = types as TransactionType[]
  } catch (e) {
    // Handled by global interceptor
  } finally {
    isLoading.value = false
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (pieChart) {
    pieChart.dispose()
    pieChart = null
  }
})
</script>

<template>
  <div class="analysis-page">
    <div class="header">
      <h2>统计</h2>
      <div class="header-actions">
        <el-date-picker
          v-model="customDateRange"
          type="daterange"
          unlink-panels
          range-separator="至"
          start-placeholder="开始"
          end-placeholder="结束"
          size="default"
          style="width: 220px; height: 36px;"
          class="custom-date-picker"
        />
        <button class="icon-btn" @click="router.push('/settings/category')" title="分类管理">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="3"></circle>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path>
          </svg>
        </button>
        <button class="icon-btn" :class="{ active: viewMode === 'chart' }" @click="viewMode = viewMode === 'list' ? 'chart' : 'list'">
          <svg v-if="viewMode === 'list'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21.21 15.89A10 10 0 1 1 8 2.83"></path><path d="M22 12A10 10 0 0 0 12 2v10z"></path></svg>
          <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="8" y1="6" x2="21" y2="6"></line><line x1="8" y1="12" x2="21" y2="12"></line><line x1="8" y1="18" x2="21" y2="18"></line><line x1="3" y1="6" x2="3.01" y2="6"></line><line x1="3" y1="12" x2="3.01" y2="12"></line><line x1="3" y1="18" x2="3.01" y2="18"></line></svg>
        </button>
      </div>
    </div>

    <!-- 顶部收支切换框(大圆角边框) -->
    <div class="tab-container">
      <div 
        class="tab-item" 
        :class="{ active: activeTab === 'income', 'income-tab': activeTab === 'income' }"
        @click="activeTab = 'income'"
      >
        <span class="label">收入</span>
        <span class="value" :class="{ 'text-income': activeTab === 'income' }">¥{{ incomeTotal.toFixed(2) }}</span>
      </div>
      <div 
        class="tab-item" 
        :class="{ active: activeTab === 'expense', 'expense-tab': activeTab === 'expense' }"
        @click="activeTab = 'expense'"
      >
        <span class="label">支出</span>
        <span class="value" :class="{ 'text-expense': activeTab === 'expense' }">¥{{ expenseTotal.toFixed(2) }}</span>
      </div>
    </div>

    <!-- 时间筛选栏 (横向滚动) -->
    <div class="filter-scroll">
      <div 
        v-for="filter in filters" 
        :key="filter.id"
        class="filter-chip"
        :class="{ active: activeFilter === filter.id }"
        @click="activeFilter = filter.id"
      >
        {{ filter.label }}
      </div>
    </div>

    <!-- 分类占比列表 (网格布局或列表布局) -->
    <div v-if="viewMode === 'list'" class="analysis-list-wrapper">
      <div v-if="analysisData.length > 0" class="analysis-grid">
        <div v-for="item in analysisData" :key="item.title" class="card category-card">
          <div class="cat-title">{{ item.title }}</div>
          <div class="cat-stats">
            <span class="cat-amount">¥{{ item.amount.toFixed(2) }}</span>
            <span class="cat-percent" :class="{ 'text-income': activeTab === 'income', 'text-expense': activeTab === 'expense' }">
              {{ item.percentage }}%
            </span>
          </div>
        </div>
      </div>
      <el-empty v-else-if="!isLoading" description="暂无数据" />
    </div>

    <!-- 饼图视图 -->
    <div v-else class="card chart-container" v-loading="isLoading" element-loading-background="transparent">
      <div v-if="analysisData.length > 0" ref="chartRef" class="pie-chart"></div>
      <el-empty v-else-if="!isLoading" description="暂无数据" />
    </div>
  </div>
</template>

<style scoped>
.analysis-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
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
  display: flex;
  align-items: center;
  justify-content: center;
}
.icon-btn.active {
  background: var(--color-primary-bg);
  color: var(--color-primary);
  border: 1px solid var(--color-primary);
}

:deep(.custom-date-picker) {
  background-color: var(--color-card-bg);
  border: none;
  box-shadow: var(--shadow-sm);
  border-radius: var(--radius-sm);
}
:deep(.custom-date-picker .el-range-input) {
  color: var(--color-text-primary);
}
:deep(.custom-date-picker .el-range-separator) {
  color: var(--color-text-secondary);
}

.tab-container {
  display: flex;
  gap: 12px;
  margin-top: 8px;
}
.tab-item {
  flex: 1;
  background: var(--color-card-bg);
  border-radius: var(--radius-lg);
  padding: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  box-shadow: var(--shadow-sm);
  border: 2px solid transparent;
  transition: all 0.3s;
}
.tab-item.active.income-tab {
  border-color: var(--color-primary);
}
.tab-item.active.expense-tab {
  border-color: var(--color-expense);
}
.tab-item .label {
  color: var(--color-text-secondary);
  font-size: 14px;
}
.tab-item .value {
  font-size: 20px;
  font-weight: bold;
}
.text-income { color: var(--color-income); }
.text-expense { color: var(--color-expense); }

/* Chart Container */
.chart-container {
  padding: 16px;
  margin-top: 16px;
}
.pie-chart {
  width: 100%;
  height: 300px;
}

.filter-scroll {
  display: flex;
  overflow-x: auto;
  gap: 12px;
  padding: 4px 0;
  -ms-overflow-style: none;  /* IE and Edge */
  scrollbar-width: none;  /* Firefox */
}
.filter-scroll::-webkit-scrollbar {
  display: none;
}
.filter-chip {
  white-space: nowrap;
  background: white;
  padding: 8px 16px;
  border-radius: var(--radius-pill);
  font-size: 14px;
  color: var(--color-text-secondary);
  cursor: pointer;
  box-shadow: var(--shadow-sm);
  transition: all 0.2s;
}
.filter-chip.active {
  background: var(--color-primary);
  color: white;
}

.analysis-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 12px;
  margin-top: 8px;
}
.category-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  background: var(--color-card-bg);
  border-radius: var(--radius-md);
  padding: 16px;
  box-shadow: var(--shadow-sm);
}
.cat-title {
  font-size: 14px;
  font-weight: 500;
}
.cat-stats {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.cat-amount {
  font-size: 14px;
  color: var(--color-text-secondary);
}
.cat-percent {
  font-size: 18px;
  font-weight: bold;
}
</style>
