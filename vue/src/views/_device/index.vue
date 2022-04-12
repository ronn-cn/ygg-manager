<template>
  <div class="app-container">
    <div class="filter-container">
      <el-button class="filter-item"  type="primary" @click="handleDeviceClick('create',null)"> 新建设备 </el-button> 
      <el-input v-model="listQuery.name" placeholder="设备名称" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-input v-model="listQuery.ouid" placeholder="OUID" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-select v-model="listQuery.status" placeholder="状态" clearable style="width: 90px" class="filter-item">
        <el-option v-for="item in statusOptions" :key="item.value" :label="item.text" :value="item.value" />
      </el-select>
      <el-select v-model="listQuery.assembly" placeholder="系统" clearable class="filter-item" style="width: 130px">
        <!-- <el-option v-for="item in calendarTypeOptions" :key="item.key" :label="item.display_name+'('+item.key+')'" :value="item.key" /> -->
      </el-select>
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter"> 搜 索 </el-button>
    </div>
    <el-table
      v-loading="listLoading"
      :data="list"
      border
      stripe
      :row-class-name="tableRowClassName"
    >
      <el-table-column label="序号" align="center" width="80">  
        <template slot-scope="{row}">
          <span>{{ row.index }}</span>
        </template>
      </el-table-column>
      <el-table-column label="设备名称" width="150" align="center">
        <template slot-scope="{row}">
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="OUID" min-width="150" show-overflow-tooltip>
        <template slot-scope="{row}">
          <span>{{ row.ouid }}</span>
        </template>
      </el-table-column>
      <el-table-column label="PIN码" min-width="80" align="center">
        <template slot-scope="{row}">
          <span v-if="row.showPin">{{ row.pin }}</span>
          <span v-else>******</span>
          </template>
      </el-table-column>
      <el-table-column label="状态"  min-width="80" 
      :filters="statusOptions"
      :filter-method="filterStatus">
        <template slot-scope="{row}">
          <span>{{ row.status.name }}</span>
        </template>
      </el-table-column>
      <el-table-column
       label="系统" 
       min-width="120" 
      :filters="[{ text: '家', value: '家' }, { text: '公司', value: '公司' }]"
      :filter-method="filterAssembly">
        <template slot-scope="{row}">
          <span>{{ row.assembly.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" align="center" sortable prop="created_at" :formatter="formatTime">
        <!-- <template slot-scope="{row}">
          <span>{{ row.created_at }}</span>
        </template> -->
      </el-table-column>
      <el-table-column label="最后在线时间" class-name="status-col" sortable prop="last_time" :formatter="formatTime">
        <!-- <template slot-scope="{row}">
          <span>{{ row.last_time }}</span>
        </template> -->
      </el-table-column>
      <el-table-column label="操作" align="center" fixed="right" width="150" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-link class="el-dropdown-link" type="primary" @click="handleDeviceClick('look',row)">查看</el-link>&nbsp;
          <el-link class="el-dropdown-link" type="primary" @click="handleDeviceClick('update',row)">编辑</el-link>&nbsp;
          <el-dropdown>
            <span class="el-dropdown-link">更多<i class="el-icon-arrow-down el-icon--right"></i></span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>修改密码</el-dropdown-item>
              <el-dropdown-item>禁用/启用</el-dropdown-item>
              <el-dropdown-item>删除</el-dropdown-item>
              <el-dropdown-item>下载</el-dropdown-item>
              <el-dropdown-item>更新</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

      <!-- 新建设备弹窗 -->
    <el-dialog :title="dialogDeviceTitle" :visible.sync="dialogDeviceFormVisible">
      <el-form ref="deviceDataForm" :model="deviceData" label-position="left" label-width="70px" style="width: 400px; margin-left:50px;">
        <el-form-item label="设备OUID">
          <el-input v-model="deviceData.ouid" placeholder="不填写，默认自动生成OUID" />
        </el-form-item>
        <el-form-item label="设备名称">
          <el-input v-model="deviceData.name" placeholder="请输入设备名称" />
        </el-form-item>
        <el-form-item label="设备PIN码">
          <el-input v-model="deviceData.pin" placeholder="不填写，默认自动生成6位PIN码"/>
        </el-form-item>
        <el-form-item label="设备系统">
          <el-select v-model="deviceData.assembly" class="filter-item" placeholder="Please select">
            <!-- <el-option v-for="item in statusOptions" :key="item" :label="item" :value="item" /> -->
          </el-select>
        </el-form-item>
        <el-form-item label="设备状态">
          <el-select v-model="deviceData.status" class="filter-item" placeholder="Please select">
            <!-- <el-option v-for="item in statusOptions" :key="item" :label="item" :value="item" /> -->
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="deviceData.remark" :autosize="{ minRows: 2, maxRows: 4}" type="textarea" placeholder="请输入设备备注信息" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogDeviceFormVisible = false"> 取消 </el-button>
        <el-button type="primary" @click="dialogDeviceStatus==='create'?createData():updateData()"> 确认 </el-button>
      </div>
    </el-dialog>

    <!-- 设备查看弹窗 -->
    <el-dialog :visible.sync="dialogDeviceLookVisible" :title="dialogDeviceTitle">
      <el-form label-position="left" :model="deviceData">
        <el-form-item label="设备名称">
          <span>{{ deviceData.name }}</span>
        </el-form-item>
        <el-form-item label="设备OUID">
          <span>{{ deviceData.ouid }}</span>
        </el-form-item>
        <el-form-item label="设备PIN码">
          <span>{{ deviceData.pin }}</span>
        </el-form-item>
        <el-form-item label="设备系统">
          <span>{{ deviceData.assembly.name }}</span>
        </el-form-item>
        <el-form-item label="设备状态">
          <span>{{ deviceData.status.name }}</span>
        </el-form-item>
        <el-form-item label="授权注册码">
          <span v-if = "deviceData.license != undefined">{{ deviceData.license.code }}</span>
        </el-form-item>
        <el-form-item label="设备备注">
          <span>{{ deviceData.remark }}</span>
        </el-form-item>
      </el-form>
    </el-dialog>

    <!-- 设备编辑弹窗 -->
  </div>
</template>

<script>
import { getDeviceList, createDevice, updateDevice } from '@/api/device'
import waves from '@/directive/waves' // waves directive
import { parseTime } from '@/utils'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination

export default {
  name: 'ComplexTable',
  components: { Pagination },
  directives: { waves },
  filters: {},
  data() {
    return {
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        name: undefined,
        ouid: undefined,
        status: undefined,
        assembly: undefined
      },
      statusOptions: [{
        value: 1,
        text: '注册'
      },{
        value: 2,
        text: '质检中'
      },{
        value: 3,
        text: '质检合格'
      },{
        value: 4,
        text: '质检不合格'
      },{
        value: 5,
        text: '销售中'
      },{
        value: 6,
        text: '已售出'
      },{
        value: 7,
        text: '待安装'
      },{
        value: 8,
        text: '安装中'
      },{
        value: 9,
        text: '安装完成'
      },{
        value: 10,
        text: '正常'
      },{
        value: 11,
        text: '停用'
      },{
        value: 12,
        text: '维护'
      },{
        value: 13,
        text: '开发'
      }],
      assemblyOptions: [],  //动态获取
      dialogDeviceLookVisible: false,
      dialogDeviceFormVisible: false,
      dialogDeviceStatus: '',
      dialogDeviceTitle: '',
      deviceData: {
        id: 0,
        name: '',
        ouid: '',
        pin: '',
        assembly: '',
        status: 0,
        license: undefined,
        product_json: undefined,
        install_json: undefined,
        remark: '',
        created_at: 0,
        updated_at: 0,
        installed_at: 0,
        slod_time: 0,
        last_time: 0,
      },
      downloadLoading: false
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true
      getDeviceList(this.listQuery).then(response => {
        this.list = response.data.items
        this.total = response.data.total
        
        this.listLoading = false
      })
    },
    tableRowClassName({row,rowIndex}){
      row.index = rowIndex;
      row.showPin = false;
    },
    // 格式化时间
    formatTime(row, column) {
      if(row[column.property] == 0){
        return "/"
      } else {
        const date = new Date(row[column.property]*1000)
        let y = date.getFullYear()
        let mo = date.getMonth() + 1
        if (mo < 10){ mo = '0' + mo }
        let d = date.getDate()
        if (d < 10){ d = '0' + d }
        let h = date.getHours()
        if (h < 10){ h = '0' + h }
        let mi = date.getMinutes()
        if (mi < 10){ mi = '0' + mi }
        let s = date.getSeconds()
        if (s < 10){ s = '0' + s }
        return y + '-' + mo + '-' + d + ' ' + h + ':' + mi + ':' + s
      }
    },
    filterStatus(value, row) {
      return row.status_id === value;
    },
    filterAssembly(value, row) {
      return row.assembly.name === value;
    },
    // 显示密码
    handleViewClick(row){
      row.showPin = true;
    },
    // 处理筛选器
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    handleDeviceClick(typ,data){
      switch(typ){
        case 'create':
          this.resetDeviceData()
          this.$nextTick(() => {
            this.$refs['deviceDataForm'].clearValidate()
          })
          this.dialogDeviceStatus = 'create'
          this.dialogDeviceTitle = "新建设备"
          this.dialogDeviceFormVisible = true
          break;
        case 'update':
          this.dialogDeviceStatus = 'update'
          this.dialogDeviceTitle = "编辑设备"
          this.deviceData = data
          this.dialogDeviceFormVisible = true
          break;
        case 'look':
          this.dialogDeviceStatus = 'look'
          this.dialogDeviceTitle = "查看设备"
          this.deviceData = data
          this.dialogDeviceLookVisible = true
          break;
        default:
          break;
      }
    },
    // 重置设备数据
    resetDeviceData() {
      this.deviceData = {
        id: 0,
        name: '',
        ouid: '',
        pin: '',
        assembly: '',
        status: 0,
        license: undefined,
        product_json: undefined,
        install_json: undefined,
        remark: '',
        created_at: 0,
        updated_at: 0,
        installed_at: 0,
        slod_time: 0,
        last_time: 0,
      }
    },
    // 处理新建设备提交
    createData() {
      this.$refs['deviceDataForm'].validate((valid) => {
        if (valid) {
          // this.temp.id = parseInt(Math.random() * 100) + 1024 // mock a id
          // this.temp.author = 'vue-element-admin'
          createDevice(this.deviceData).then(() => {
            this.list.unshift(this.deviceData)
            this.dialogDeviceFormVisible = false
            this.$notify({
              title: '请求成功',
              message: '创建设备成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    // 处理更新设备提交
    updateData() {
      this.$refs['deviceDataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.deviceData)
          // tempData.timestamp = +new Date(tempData.timestamp) // change Thu Nov 30 2017 16:41:05 GMT+0800 (CST) to 1512031311464
          updateDevice(tempData).then(() => {
            const index = this.list.findIndex(v => v.id === this.temp.id)
            this.list.splice(index, 1, this.deviceData)
            this.dialogDeviceFormVisible = false
            this.$notify({
              title: '请求成功',
              message: '更新设备成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleDelete(row, index) {
      this.$notify({
        title: 'Success',
        message: 'Delete Successfully',
        type: 'success',
        duration: 2000
      })
      this.list.splice(index, 1)
    },
    handleDownload() {
      this.downloadLoading = true
      import('@/vendor/Export2Excel').then(excel => {
        const tHeader = ['timestamp', 'title', 'type', 'importance', 'status']
        const filterVal = ['timestamp', 'title', 'type', 'importance', 'status']
        const data = this.formatJson(filterVal)
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: 'table-list'
        })
        this.downloadLoading = false
      })
    },
    formatJson(filterVal) {
      return this.list.map(v => filterVal.map(j => {
        if (j === 'timestamp') {
          return parseTime(v[j])
        } else {
          return v[j]
        }
      }))
    }
  }
}
</script>

<style lang="scss" scoped>
.filter-container{
  .filter-item{
    margin-right: 20px;
  }
}

.el-dropdown-link {
  cursor: pointer;
  color: #409EFF;
}
</style>