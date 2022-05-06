import request from '@/utils/request'

// 查询注册授权码列表
export function getLicenseList(query) {
  return request({
    url: '/license/get-license-list',
    method: 'get',
    params: query
  })
}

// 创建注册授权码
export function createLicense(data) {
  return request({
    url: '/license/create-license',
    method: 'post',
    data
  })
}


// 更新注册授权码
export function updateLicense(data) {
  return request({
    url: '/license/update-license',
    method: 'post',
    data
  })
}


// 删除注册授权码
export function deleteLicense(id) {
  return request({
    url: '/license/delete-license',
    method: 'delete',
    params: {"code":id}
  })
}


// 查询注册授权码
export function queryLicense(id) {
  return request({
    url: '/license/query-license',
    method: 'get',
    params: { id }
  })
}