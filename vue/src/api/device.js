import request from '@/utils/request'

// 查询设备列表
export function getDeviceList(query) {
  return request({
    url: '/device/get-device-list',
    method: 'get',
    params: query
  })
}

// 查询设备列表
export function queryDeviceStatus(query) {
  return request({
    url: '/device/query-device-status',
    method: 'get',
    params: query
  })
}

// 创建设备
export function createDevice(data) {
  return request({
    url: '/device/create-device',
    method: 'post',
    data
  })
}


// 更新设备
export function updateDevice(data) {
  return request({
    url: '/device/update-device',
    method: 'put',
    data
  })
}

// 更新设备PIN
export function updateDevicePin(data) {
  return request({
    url: '/device/update-device-pin',
    method: 'put',
    data
  })
}


// 删除设备
export function deleteDevice(id) {
  return request({
    url: '/device/delete-device',
    method: 'delete',
    params: { ouid:id }
  })
}


// 查询设备
export function queryDevice(id) {
  return request({
    url: '/device/query-device',
    method: 'get',
    params: { id }
  })
}

// 查询设备
export function pushUpdate(id) {
  return request({
    url: '/device/push-update?ouid='+id,
    method: 'post'
  })
}