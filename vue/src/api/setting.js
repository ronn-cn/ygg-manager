import request from '@/utils/request'

// 查询设置
export function getSetting() {
  return request({
    url: '/setting/get-setting',
    method: 'get'
  })
}

export function saveSetting(data) {
  return request({
    url: '/setting/save-setting',
    method: 'post',
    data
  })
}