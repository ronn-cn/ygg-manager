import request from '@/utils/request'

// 查询应用列表
export function getApplicationList(query) {
  return request({
    url: '/application/get-application-list',
    method: 'get',
    params: query
  })
}

// 创建应用
export function createApplication(data) {
  return request({
    url: '/application/create-application',
    method: 'post',
    data
  })
}


// 更新应用
export function updateApplication(data) {
  return request({
    url: '/application/update-application',
    method: 'put',
    data
  })
}


// 删除应用
export function deleteApplication(id) {
  return request({
    url: '/application/delete-application',
    method: 'delete',
    params: { appid: id }
  })
}


// 查询应用
export function queryApplication(id) {
  return request({
    url: '/application/query-application',
    method: 'get',
    params: { appid: id }
  })
}