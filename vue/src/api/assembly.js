import request from '@/utils/request'

// 查询系统列表
export function getAssemblyList(query) {
  return request({
    url: '/assembly/get-assembly-list',
    method: 'get',
    params: query
  })
}

// 创建系统
export function createAssembly(data) {
  return request({
    url: '/assembly/create-assembly',
    method: 'post',
    data
  })
}


// 更新系统
export function updateAssembly(data) {
  return request({
    url: '/assembly/update-assembly',
    method: 'post',
    data
  })
}


// 删除系统
export function deleteAssembly(id) {
  return request({
    url: '/assembly/delete-assembly',
    method: 'get',
    params: { id }
  })
}


// 查询系统
export function queryAssembly(id) {
  return request({
    url: '/assembly/query-assembly',
    method: 'get',
    params: { id }
  })
}