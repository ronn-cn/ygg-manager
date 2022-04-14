<template>
  <div class="app-container">
    <div class="filter-container">
      <!-- 新建系统按钮 -->
      <el-button class="filter-item"  type="primary" @click="handleSystemClick('create', null)"> 新建系统 </el-button> 
      <el-input v-model="listQuery.title" placeholder="系统名称" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
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
      <el-table-column label="系统名称" min-width="150">
        <template slot-scope="{row}">
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="OUID" align="center" min-width="150" show-overflow-tooltip>
        <template slot-scope="{row}">
          <span>{{ row.ouid }}</span>
        </template>
      </el-table-column>
      <el-table-column label="应用列表" min-width="110">
        <template slot-scope="{row}">
          <span>{{ row.list }}</span>
        </template>
      </el-table-column>
      <el-table-column label="主应用" min-width="110">
        <template slot-scope="{row}">
          <span>{{ row.main }}</span>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" align="center" sortable prop="created_at" :formatter="formatTime"></el-table-column>
      <el-table-column label="更新时间" align="center" sortable prop="updated_at" :formatter="formatTime"></el-table-column>
      <el-table-column label="操作" align="center" fixed="right" width="150" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-link class="el-dropdown-link" type="primary" @click="handleSystemClick('look',row)">查看</el-link>&nbsp;
          <el-link class="el-dropdown-link" type="primary" @click="handleSystemClick('update',row)">编辑</el-link>&nbsp;
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

    <!-- 系统新建与编辑弹窗 -->
    <el-dialog :title="dialogSystemFormTitle" :visible.sync="dialogSystemFormVisible" :close-on-click-modal="false">
      <el-form ref="systemData" :model="systemData" label-position="left" label-width="70px">
        <el-form-item label="OUID" prop="ouid">
          <el-input v-model="systemData.ouid" placeholder="不填写，默认自动生成OUID" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="systemData.name" placeholder="请输入系统名称"/>
        </el-form-item>
        <el-form-item label="应用集合" prop="list">
          <el-cascader
            placeholder="试试搜索：**应用"
            v-model="systemData.listarr"
            :options="optionsAppData()"
            :props="{ multiple: true }"
            filterable
            clearable
            @change="systemDataListChange"
            >
          </el-cascader>
        </el-form-item>
        <el-form-item label="主应用" prop="main">
          <el-select v-model="systemData.main" placeholder="请选择">
            <el-option
              v-for="item in systemData.mainOptionArr"
              :key="item.label"
              :label="item.label"
              :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="systemData.remark" :autosize="{ minRows: 2, maxRows: 4}" type="textarea" placeholder="Please input" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogSystemFormVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="dialogSystemFormStatus==='create'?createSystemData():updateSystemData()">
          提交
        </el-button>
      </div>
    </el-dialog>

    <!-- 系统查看弹窗 -->
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import { getSystemList,createSystem } from '@/api/system'
import { getApplicationList } from '@/api/application'
import waves from '@/directive/waves' // waves directive
import { parseTime } from '@/utils'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination

export default {
  name: 'ComplexTable',
  components: { Pagination },
  directives:
   { waves },
  filters: {},
  computed: {
    ...mapGetters(["info", "roles"]),
  },
  data() {
    return {
      list: [],
      appList: [],
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        name: undefined,
      },
      systemData: {
        ouid: '',
        name: '',
        list: '',
        listarr: [],
        mainOptionArr:[],
        main: '',
        remark: '',
        status: 1,
      },
      dialogSystemLookVisible: false,
      dialogSystemFormVisible: false,
      dialogSystemFormStatus: '',
      dialogSystemFormTitle: '',
      downloadLoading: false
    }
  },
  created() {
    this.getList()
    this.getAppList()
  },
  methods: {
    getList() {
      this.listLoading = true
      getSystemList(this.listQuery).then(response => {
        this.list = response.data.items
        this.total = response.data.total

        this.listLoading = false
      })
    },
    getAppList(){
      getApplicationList().then(response => {
        this.appList = response.data.items
      })
    },
    optionsAppData(){
      let appOptions = []
      console.log(this.appList)
      for (let i = 0; i < this.appList.length; i++){
          let item = {
            value: this.appList[i].appid,
            label: this.appList[i].name+"["+this.appList[i].appid+"]",
          }
          appOptions.push(item)
      }
      return appOptions
    },
    // 系统的应用列表数据改变时操作的函数
    systemDataListChange(){
      this.systemData.mainOptionArr = []
      this.systemData.list = this.systemData.listarr.toString()
      for (var i = 0; i < this.appList.length; i++){
        if(this.systemData.list.includes(this.appList[i].appid)){
          let opt = {
            label: this.appList[i].name,
            value: this.appList[i].appid,
          }
          this.systemData.mainOptionArr.push(opt)
        }
      }
      if(!this.systemData.list.includes(this.systemData.main)){
        this.systemData.main = "";
      }
    },
    tableRowClassName({row,rowIndex}){
      row.index = rowIndex;
    },
    // 格式化时间
    formatTime(row, column) {
      if(row[column.property] == 0){
        return "/"
      } else {
        return parseTime(row[column.property])
      }
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    resetSystemData() {
      this.systemData = {
        ouid: '',
        name: '',
        list: '',
        main: '',
        remark: '',
        status: 1,
      }
    },
    handleSystemClick(typ, data) {
      switch(typ){
        case 'create':
          this.resetSystemData()
          this.$nextTick(() => {
            this.$refs['systemData'].clearValidate()
          })
          this.dialogSystemFormStatus = 'create'
          this.dialogSystemFormTitle = "新建系统"
          this.dialogSystemFormVisible = true
          break;
        case 'update':
          this.dialogSystemFormStatus = 'update'
          this.dialogSystemFormTitle = "编辑系统"
          this.systemData = data
          this.systemData.listarr = data.list.split(",")
          this.dialogSystemFormVisible = true
          break;
        case 'look':
          this.dialogSystemFormStatus = 'look'
          this.dialogSystemFormTitle = "查看系统"
          this.systemData = data
          this.dialogSystemLookVisible = true
          break;
        default:
          break;
      }
    },
    createSystemData() {
      this.$refs['systemData'].validate((valid) => {
        if (valid) {
          createSystem(this.systemData).then(() => {
            this.getList()
            this.dialogSystemFormVisible = false
            this.$message({
              message: '创建系统成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleSystemUpdate(row) {
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
    handleFetchPv(pv) {
      fetchPv(pv).then(response => {
        this.pvData = response.data.pvData
        this.dialogPvVisible = true
      })
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
    },
    getSortClass: function(key) {
      const sort = this.listQuery.sort
      return sort === `+${key}` ? 'ascending' : 'descending'
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
