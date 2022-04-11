import request from '@/utils/request'

// 查询日志列表
export function getRecordList(query) {
  return request({
    url: '/record/get-record-list',
    method: 'get',
    params: query
  })
}


// 删除日志
export function deleteRecord(id) {
  return request({
    url: '/record/delete-record',
    method: 'get',
    params: { id }
  })
}


// 查询日志
export function queryRecord(id) {
  return request({
    url: '/record/query-record',
    method: 'get',
    params: { id }
  })
}