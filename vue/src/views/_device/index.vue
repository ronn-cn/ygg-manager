<template>
  <div class="app-container">
    <div class="filter-container">
      <el-button class="filter-item"  type="primary" @click="dialogNewAccountVisible = true"> 新建设备 </el-button> 
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
          <el-link class="el-dropdown-link" type="primary" @click="handleLookClick(row)">查看</el-link>&nbsp;
          <el-link class="el-dropdown-link" type="primary">编辑</el-link>&nbsp;
          <el-dropdown>
            <span class="el-dropdown-link">更多<i class="el-icon-arrow-down el-icon--right"></i></span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>修改密码</el-dropdown-item>
              <el-dropdown-item>禁用/启用</el-dropdown-item>
              <el-dropdown-item>删除设备</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="70px" style="width: 400px; margin-left:50px;">
        <el-form-item label="Type" prop="type">
          <el-select v-model="temp.type" class="filter-item" placeholder="Please select">
            <!-- <el-option v-for="item in calendarTypeOptions" :key="item.key" :label="item.display_name" :value="item.key" /> -->
          </el-select>
        </el-form-item>
        <el-form-item label="Date" prop="timestamp">
          <el-date-picker v-model="temp.timestamp" type="datetime" placeholder="Please pick a date" />
        </el-form-item>
        <el-form-item label="Title" prop="title">
          <el-input v-model="temp.title" />
        </el-form-item>
        <el-form-item label="Status">
          <el-select v-model="temp.status" class="filter-item" placeholder="Please select">
            <el-option v-for="item in statusOptions" :key="item" :label="item" :value="item" />
          </el-select>
        </el-form-item>
        <el-form-item label="Imp">
          <el-rate v-model="temp.importance" :colors="['#99A9BF', '#F7BA2A', '#FF9900']" :max="3" style="margin-top:8px;" />
        </el-form-item>
        <el-form-item label="Remark">
          <el-input v-model="temp.remark" :autosize="{ minRows: 2, maxRows: 4}" type="textarea" placeholder="Please input" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          Cancel
        </el-button>
        <el-button type="primary" @click="dialogStatus==='create'?createData():updateData()">
          Confirm
        </el-button>
      </div>
    </el-dialog>

    <el-dialog :visible.sync="dialogLookVisible" title="Reading statistics">
      <template slot-scope="dialogLookData">
        <el-form label-position="left" inline class="demo-table-expand">
          <el-form-item label="设备名称">
            <span>{{ dialogLookData.name }}</span>
          </el-form-item>
          <el-form-item label="所属店铺">
            <span>{{ dialogLookData.shop }}</span>
          </el-form-item>
          <el-form-item label="商品 ID">
            <span>{{ dialogLookData.id }}</span>
          </el-form-item>
          <el-form-item label="店铺 ID">
            <span>{{ dialogLookData.shopId }}</span>
          </el-form-item>
          <el-form-item label="商品分类">
            <span>{{ dialogLookData.category }}</span>
          </el-form-item>
          <el-form-item label="店铺地址">
            <span>{{ dialogLookData.address }}</span>
          </el-form-item>
          <el-form-item label="商品描述">
            <span>{{ dialogLookData.desc }}</span>
          </el-form-item>
        </el-form>
      </template>
    </el-dialog>
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
  filters: {
    statusFilter(status) {
      const statusMap = {
        published: 'success',
        draft: 'info',
        deleted: 'danger'
      }
      return statusMap[status]
    }
  },
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
      temp: {
        id: undefined,
        importance: 1,
        remark: '',
        timestamp: new Date(),
        title: '',
        type: '',
        status: 'published'
      },
      dialogLookVisible: false,
      dialogLookData: null,
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: 'Edit',
        create: 'Create'
      },
      rules: {
        type: [{ required: true, message: 'type is required', trigger: 'change' }],
        timestamp: [{ type: 'date', required: true, message: 'timestamp is required', trigger: 'change' }],
        title: [{ required: true, message: 'title is required', trigger: 'blur' }]
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
    handleViewClick(row){
      row.showPin = true;
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    handleLookClick(row){
      this.dialogLookData = row;
      console.log("row", row);
      this.dialogLookVisible = true;
    },
    handleModifyStatus(row, status) {
      this.$message({
        message: '操作Success',
        type: 'success'
      })
      row.status = status
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        importance: 1,
        remark: '',
        timestamp: new Date(),
        title: '',
        status: 'published',
        type: ''
      }
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          this.temp.id = parseInt(Math.random() * 100) + 1024 // mock a id
          this.temp.author = 'vue-element-admin'
          createArticle(this.temp).then(() => {
            this.list.unshift(this.temp)
            this.dialogFormVisible = false
            this.$notify({
              title: 'Success',
              message: 'Created Successfully',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row) // copy obj
      this.temp.timestamp = new Date(this.temp.timestamp)
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.temp)
          tempData.timestamp = +new Date(tempData.timestamp) // change Thu Nov 30 2017 16:41:05 GMT+0800 (CST) to 1512031311464
          updateArticle(tempData).then(() => {
            const index = this.list.findIndex(v => v.id === this.temp.id)
            this.list.splice(index, 1, this.temp)
            this.dialogFormVisible = false
            this.$notify({
              title: 'Success',
              message: 'Update Successfully',
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
