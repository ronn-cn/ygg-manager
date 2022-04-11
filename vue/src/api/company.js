import request from '@/utils/request'

// 查询商户列表
export function getCompanyList(query) {
  return request({
    url: '/company/get-company-list',
    method: 'get',
    params: query
  })
}

// 创建商户
export function createCompany(data) {
  return request({
    url: '/company/create-company',
    method: 'post',
    data
  })
}


// 更新商户
export function updateCompany(data) {
  return request({
    url: '/company/update-company',
    method: 'post',
    data
  })
}


// 删除商户
export function deleteCompany(id) {
  return request({
    url: '/company/delete-company',
    method: 'get',
    params: { id }
  })
}


// 查询商户
export function queryCompany(id) {
  return request({
    url: '/company/query-company',
    method: 'get',
    params: { id }
  })
}