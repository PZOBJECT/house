import axios from 'axios'
import { ElMessage } from 'element-plus'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

api.interceptors.response.use(
  (res) => {
    const data = res.data
    if (data.code !== 0) {
      ElMessage.error(data.msg || '请求失败')
      return Promise.reject(new Error(data.msg))
    }
    return data
  },
  (err) => {
    ElMessage.error('网络错误，请稍后重试')
    return Promise.reject(err)
  }
)

export function getRooms(rented) {
  return api.get('/rooms', { params: { rented } })
}

export function updateRoom(roomNo, data) {
  return api.put(`/rooms/${roomNo}`, data)
}

export function generateBills(year, month) {
  return api.post('/bills/generate', { year, month })
}

export function getBillList(params) {
  return api.get('/bills/list', { params })
}

export function updateBill(id, data) {
  return api.put(`/bills/${id}`, data)
}

export function payBill(id) {
  return api.put(`/bills/${id}/pay`)
}

export function unpayBill(id) {
  return api.put(`/bills/${id}/unpay`)
}

export function getReceipt(id) {
  return api.get(`/bills/${id}/receipt`)
}
