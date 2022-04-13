import request from '@/utils/request'

// 查询账号列表
export function getAccountList(query) {
  return request({
    url: '/account/get-account-list',
    method: 'get',
    params: query
  })
}

// 创建账号
export function createAccount(data) {
  return request({
    url: '/account/create-account',
    method: 'post',
    data
  })
}

// 更新账号
export function updateAccount(data) {
  return request({
    url: '/account/update-account',
    method: 'put',
    data
  })
}

// 删除账号
export function deleteAccount(ouid) {
  console.log("ouid:",ouid)
  return request({
    url: '/account/delete-account',
    method: 'delete',
    params: {"ouid":ouid}
  })
}

// 登录账号
export function login(data) {
  return request({
    url: '/account/login',
    method: 'post',
    data
  })
}

export function getAccountInfo() {
  return request({
    url: '/account/get-account-info',
    method: 'get'
  })
}

export function getAccountPermission() {
  return request({
    url: '/account/get-account-permission',
    method: 'get'
  })
}

export function logout() {
  return request({
    url: '/account/logout',
    method: 'post'
  })
}

export function fetchArticle(id) {
  return request({
    url: '/vue-element-admin/article/detail',
    method: 'get',
    params: { id }
  })
}

export function fetchPv(pv) {
  return request({
    url: '/vue-element-admin/article/pv',
    method: 'get',
    params: { pv }
  })
}

export function updateArticle(data) {
  return request({
    url: '/vue-element-admin/article/update',
    method: 'post',
    data
  })
}
