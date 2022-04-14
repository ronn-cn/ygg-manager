import request from '@/utils/request'

// 查询系统列表
export function getSystemList(query) {
  return request({
    url: '/system/get-system-list',
    method: 'get',
    params: query
  })
}

// 创建系统
export function createSystem(data) {
  return request({
    url: '/system/create-system',
    method: 'post',
    data
  })
}


// 更新系统
export function updateSystem(data) {
  return request({
    url: '/system/update-system',
    method: 'post',
    data
  })
}


// 删除系统
export function deleteSystem(id) {
  return request({
    url: '/system/delete-system',
    method: 'get',
    params: { id }
  })
}


// 查询系统
export function querySystem(id) {
  return request({
    url: '/system/query-system',
    method: 'get',
    params: { id }
  })
}