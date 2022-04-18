<template>
  <div class="app-container">
    <div class="filter-container">
      <!-- 新建系统按钮 -->
      <el-button class="filter-item"  type="primary" @click="handleSystemClick('create', null)"> 新建系统 </el-button> 
      <el-input v-model="listQuery.name" placeholder="系统名称" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
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
      <!-- type="expand" -->
      <el-table-column label="应用列表" min-width="110">  
        <template slot-scope="{row}">
          <!-- <span>{{ row.list }}</span> -->
          <div v-for="app in row.listarr" :key="app">
            <el-tag style="margin:2px 0" size="small" v-if="app!=''">{{app}}</el-tag>
          </div>
        </template>
      </el-table-column>
      <!-- <el-table-column label="主应用" min-width="110">
        <template slot-scope="{row}">
          <span>{{ row.main }}</span>
        </template>
      </el-table-column> -->
      <!-- <el-table-column label="创建时间" align="center" sortable prop="created_at" :formatter="formatTime"></el-table-column> -->
      <el-table-column label="更新时间" align="center" sortable prop="updated_at" :formatter="formatTime"></el-table-column>
      <el-table-column label="操作" align="center" fixed="right" width="150" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-link class="el-dropdown-link" type="primary" @click="handleSystemClick('look',row)">查看</el-link>&nbsp;
          <el-link class="el-dropdown-link" type="primary" @click="handleSystemClick('update',row)">编辑</el-link>&nbsp;
          <el-dropdown>
            <span class="el-dropdown-link">更多<i class="el-icon-arrow-down el-icon--right"></i></span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>推送更新</el-dropdown-item>
              <el-dropdown-item>设为开发</el-dropdown-item>
              <el-dropdown-item>删除系统</el-dropdown-item>
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
          <el-input v-model="systemData.remark" :autosize="{ minRows: 2, maxRows: 4}" type="textarea" placeholder="请填写备注信息" />
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
    <el-dialog :title="dialogSystemFormTitle" :visible.sync="dialogSystemLookVisible">
      <el-descriptions class="margin-top" :column="3" border>
        <el-descriptions-item>
          <template slot="label">
            系统名称
          </template>
          {{ systemData.name }}
        </el-descriptions-item>
        <el-descriptions-item span="2">
          <template slot="label">
            OUID
          </template>
          {{ systemData.ouid }}
        </el-descriptions-item>
        <el-descriptions-item>
          <template slot="label">
            状态
          </template>
          {{ systemData.status?'正式':'开发' }}
        </el-descriptions-item>
        
        <el-descriptions-item>
          <template slot="label">
            创建时间
          </template>
          {{ parseTime(systemData.created_at) }}
        </el-descriptions-item>
        
        <el-descriptions-item>
          <template slot="label">
            更新时间
          </template>
          {{ parseTime(systemData.updated_at) }}
        </el-descriptions-item>
        <el-descriptions-item span="3" label="备注">
          {{ systemData.remark }}
        </el-descriptions-item>
        
        <el-descriptions-item span="3" label="应用列表">
          <el-table :data="systemData.applist" style="width: 100%">
            <el-table-column prop="name" label="应用名称"> </el-table-column>
            <el-table-column prop="appid" label="应用ID" show-overflow-tooltip> </el-table-column>
            <el-table-column prop="latest" label="最新版本"> </el-table-column>
            <el-table-column label="主应用标记" align="center">
              <template slot-scope="{row}">
                <el-button v-if="row.main" type="danger" size="mini" icon="el-icon-check" circle></el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
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
        applist: [],
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
      // this.systemDataListChange()
      return appOptions
    },
    // 系统的应用列表数据改变时操作的函数
    systemDataListChange(){
      console.log("systemDataListChange:",this.systemData.mainOptionArr)
      if(this.systemData.mainOptionArr){
        this.systemData.mainOptionArr.length = 0;
      } else {
        this.systemData.mainOptionArr = [];
      }
      this.systemData.list = this.systemData.listarr.toString()
      if(!this.systemData.list.includes(this.systemData.main)){
        this.systemData.main = "";
      }
      for (var i = 0; i < this.appList.length; i++){
        if(this.systemData.list.includes(this.appList[i].appid)){
          let opt = {
            label: this.appList[i].name,
            value: this.appList[i].appid,
          }
          this.systemData.mainOptionArr.push(opt)
        }
      }
    },
    tableRowClassName({row,rowIndex}){
      row.index = rowIndex;
      row.listarr =  row.list.split(",");
    },
    // 格式化时间
    formatTime(row, column) {
      if(row[column.property] == 0){
        return "/"
      } else {
        return parseTime(row[column.property])
      }
    },
    parseTime(time){
      return parseTime(time)
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
        listarr: [],
        applist: [],
        mainOptionArr:[],
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
          this.systemData.ouid = data.ouid
          this.systemData.name = data.name
          this.systemData.list = data.list
          this.systemData.main = data.main
          this.systemData.remark = data.remark
          this.systemData.listarr = data.list.split(",")
          this.systemDataListChange()
          this.dialogSystemFormVisible = true
          break;
        case 'look':
          this.dialogSystemFormStatus = 'look'
          this.dialogSystemFormTitle = "查看系统"
          this.systemData = data
          this.systemData.applist = [];
          // 查询应用列表 
          for (let i = 0; i < this.appList.length; i++){
            this.appList[i].main = false
            if (this.systemData.list.includes(this.appList[i].appid)){
              let item = this.appList[i]
              if (this.appList[i].appid == this.systemData.main){
                item.main = true
              }
              this.systemData.applist.push(item)
            }
          }
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
    updateSystemData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.systemData)
          updateSystem(tempData).then(() => {
            this.getList()
            this.dialogSystemFormVisible = false
            this.$message({
              message: '更新系统成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
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
