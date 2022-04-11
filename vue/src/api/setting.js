import request from '@/utils/request'

// 查询设置
export function getSetting() {
  return request({
    url: '/setting/get-setting',
    method: 'get'
  })
}