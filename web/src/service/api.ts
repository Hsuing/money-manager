import axios from 'axios'
import { ElMessage } from 'element-plus'

export interface Category {
  id: number
  name: string
  type: string
  parent_id: number
  icon: string
  is_exclude: boolean
}

export interface Account {
  id: number
  name: string
  type: string
  account_type: string
  balance: number
  exclude_amount: number
  account_number: string
  remark: string
}

export interface Transaction {
  id: number
  account_id: number
  category_id: number
  amount: number
  type: string
  remark: string
  record_date: string
  account?: Account
  category?: Category
}

export interface MonthlyStat {
  month: number
  income: number
  expense: number
  balance: number
}

export interface DashboardData {
  total_assets: number
  net_assets: number
  year: number
  year_income: number
  year_expense: number
  year_balance: number
  month: number
  month_income: number
  month_expense: number
  month_balance: number
  monthly_stats: MonthlyStat[]
}

export interface TransactionType {
  id: number
  name: string
  flow: string
  is_exclude: boolean
  is_disabled: boolean
}

// 创建 axios 实例
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL, // 后端地址，从环境变量读取
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器：自动带上 token
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => Promise.reject(error)
)

// \u54cd\u5e94\u62e6\u622a\u5668\uff1a\u5904\u7406 401 \u767b\u51fa & \u63d0\u53d6\u6807\u51c6 response.data
api.interceptors.response.use(
  response => {
    const res = response.data
    // \u5982\u679c\u662f\u4e0b\u8f7d\u6587\u4ef6\u6216\u8005\u975e\u6807\u51c6\u7ed3\u6784\uff0c\u76f4\u63a5\u8fd4\u56de
    if (response.config.responseType === 'blob' || typeof res !== 'object' || res === null || !('code' in res)) {
      return res
    }

    if (res.code === 0) {
      return res.data
    } else {
      ElMessage.error(res.msg || 'API Error')
      return Promise.reject(new Error(res.msg || 'API Error'))
    }
  },
  error => {
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    } else {
      ElMessage.error(error.response?.data?.msg || error.response?.data?.error || error.message || '\u7f51\u7edc\u9519\u8bef')
    }
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

export const login = (data: any): Promise<any> => api.post('/login', data)
export const register = (data: any): Promise<any> => api.post('/register', data)
export const getCaptcha = (): Promise<any> => api.get('/captcha')
export const getDashboard = (year?: number): Promise<any> => api.get('/dashboard', { params: { year, _t: Date.now() } })
export const getTransactions = (year?: number, month?: number): Promise<any> => 
  api.get('/transactions', { params: { year, month } })
export const addTransaction = (data: any): Promise<any> => api.post('/transactions', data)
export const deleteTransaction = (id: string|number): Promise<any> => api.delete(`/transactions/${id}`)
export const getCategories = (): Promise<any> => api.get('/categories')
export const addCategory = (data: any): Promise<any> => api.post('/categories', data)
export const updateCategory = (id: string|number, data: any): Promise<any> => api.put(`/categories/${id}`, data)
export const deleteCategory = (id: string|number): Promise<any> => api.delete(`/categories/${id}`)
export const getAccounts = (): Promise<any> => api.get('/accounts')
export const addAccount = (data: any): Promise<any> => api.post('/accounts', data)
export const updateAccount = (id: string|number, data: any): Promise<any> => api.put(`/accounts/${id}`, data)
export const deleteAccount = (id: string|number): Promise<any> => api.delete(`/accounts/${id}`)

export const getTransactionTypes = (): Promise<any> => api.get('/transaction_types')
export const addTransactionType = (data: any): Promise<any> => api.post('/transaction_types', data)
export const updateTransactionType = (id: string|number, data: any): Promise<any> => api.put(`/transaction_types/${id}`, data)
export const deleteTransactionType = (id: string|number): Promise<any> => api.delete(`/transaction_types/${id}`)

export default api
