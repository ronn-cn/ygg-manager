<template>
  <div class="app-container">
    <div class="filter-container">
      <!-- 新建系统按钮 -->
      <el-button class="filter-item" type="primary" @click="handleSystemClick('create', null)"> 新建系统 </el-button>
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
      <el-table-column label="系统名称" min-width="150" show-overflow-tooltip>
        <template slot-scope="{row}">
          <span>{{ row.name }}</span>
          <el-tag v-if="row.status == 0" style="margin-left:5px;" size="mini" type="info" effect="dark">开发</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="OUID" align="center" min-width="150" show-overflow-tooltip>
        <template slot-scope="{row}">
          <span>{{ row.ouid }}</span>
        </template>
      </el-table-column>
      <el-table-column label="应用列表" min-width="110" show-overflow-tooltip>
        <template slot-scope="{row}">
          <!-- <el-tag style="margin:5px" v-for="app in row.listarr" size="small" :key="app.appid">{{"[ "+app.name+" ]"}}</el-tag> -->
          <el-tag v-for="name in returnAppNames(row.list)" :key="name" style="margin:5px" size="small">{{ "[ "+name+" ]" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="更新时间" align="center" sortable prop="updated_at" :formatter="formatTime" />
      <el-table-column label="操作" align="center" fixed="right" width="150" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-link class="el-dropdown-link" type="primary" @click="handleSystemClick('look',row)">查看</el-link>&nbsp;
          <el-link class="el-dropdown-link" type="primary" @click="handleSystemClick('update',row)">编辑</el-link>&nbsp;
          <el-dropdown>
            <span class="el-dropdown-link">更多<i class="el-icon-arrow-down el-icon--right" /></span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item @click.native="pushSystemData(row)">推送更新</el-dropdown-item>
              <el-dropdown-item v-if="row.status == 0">发布系统</el-dropdown-item>
              <el-dropdown-item @click.native="deleteSystemData(row.ouid)">删除系统</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <!-- 系统新建与编辑弹窗 -->
    <el-dialog :title="dialogSystemFormTitle" :visible.sync="dialogSystemFormVisible" :close-on-press-escape="false" :close-on-click-modal="false" width="65%; height:70%">
      <el-form ref="systemData" :model="systemData" label-position="left" label-width="70px">
        <el-form-item label="OUID" prop="ouid">
          <el-input v-model="systemData.ouid" placeholder="不填写，默认自动生成OUID" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="systemData.name" placeholder="请输入系统名称" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="systemData.remark" :autosize="{ minRows: 2, maxRows: 4}" type="textarea" placeholder="请填写备注信息" />
        </el-form-item>
        <el-form-item label="选择应用">
          <!-- TODO:等应用特别多的时候需要添加上搜索功能 2.还需要添加一栏显示视图、程序等应用类型-->
          <!-- <el-input style="width: 200px;" v-model="appSearch" size="mini" placeholder="输入关键字搜索"/> -->
          <el-table
            ref="multipleTable"
            :data="activeAppList"
            border
            height="300"
            tooltip-effect="dark"
            style="width: 100%"
            @select="selectItem"
            @select-all="selectAll"
            @selection-change="selectionApplistChanged"
          >
            <el-table-column type="selection" width="55" />
            <el-table-column prop="name" label="应用名称" show-overflow-tooltip />
            <el-table-column prop="appid" label="应用ID" show-overflow-tooltip />
            <el-table-column label="选择版本">
              <template slot-scope="{row}">
                <el-select v-model="row.version" placeholder="选择版本" :disabled="row.disabled">
                  <el-option v-for="item in row.versionOptions" :key="item" :label="item" :value="item" />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column prop="" label="开机启动" width="120">
              <template slot-scope="{row}">
                <el-switch :value="row.appid==systemData.main?true:false" :disabled="row.disabled" @change="switchChange(row)" />
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogSystemFormVisible = false">取消</el-button>
        <el-button type="primary" @click="dialogSystemFormStatus==='create'?createSystemData():updateSystemData()">提交</el-button>
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
            <el-table-column prop="name" label="应用名称" />
            <el-table-column prop="appid" label="应用ID" show-overflow-tooltip />
            <el-table-column prop="latest" label="最新版本" />
            <el-table-column label="标记" align="center">
              <template slot-scope="{row}">
                <el-tag v-if="row.main" type="danger">主应用</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { getSystemList, createSystem, updateSystem, deleteSystem, pushSystem } from '@/api/system'
import { getApplicationList, queryApplicationVersion } from '@/api/application'
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
    ...mapGetters(['info', 'roles'])
  },
  data() {
    return {
      list: [],
      appList: [], // 这是一个全局
      activeAppList: [],
      multipleAppList: [],
      appSearch: undefined,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        name: undefined
      },
      systemData: {
        ouid: '',
        name: '',
        list: '',
        listarr: [],
        main: '',
        remark: '',
        status: 1
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
    // 获取所有应用数据
    getAppList() {
      getApplicationList().then(response => {
        this.appList = response.data.items
        for (let i = 0; i < this.appList.length; i++) {
          const item = {
            appid: this.appList[i].appid,
            name: this.appList[i].name,
            version: 'latest',
            type: this.appList[i].type,
            disabled: true,
            versionOptions: ['latest']
          }
          queryApplicationVersion({ 'appid': this.appList[i].appid }).then(response => {
            var appversion = response.data.items
            const len = appversion.length > 10 ? 10 : appversion.length
            for (let j = len - 1; j >= 0; j--) {
              item.versionOptions.push(appversion[j].version)
            }
            this.activeAppList.push(item) // 生成提供给应用表格的数据源
          })
        }
      })
    },
    // 返回应用版本的选项
    async returnAppVersionOptions(appid) {
      const versionOptions = []
      await queryApplicationVersion({ 'appid': appid }).then(response => {
        var appversion = response.data.items
        for (let i = 0; i < appversion.length; i++) {
          const item = {
            value: appversion[i].version,
            label: appversion[i].version
          }
          versionOptions.push(item)
        }
      })
      return versionOptions
    },
    selectionApplistChanged(val) {
      this.multipleAppList = val
    },
    selectItem(item, row) {
      console.log(item, row)
      row.disabled = !row.disabled
    },
    selectAll(item) {
      if (item.length > 0) { // 说明是全选
        for (let i = 0; i < this.activeAppList.length; i++) {
          this.activeAppList[i].disabled = false
        }
      } else { // 否则是取消全选
        for (let i = 0; i < this.activeAppList.length; i++) {
          this.activeAppList[i].disabled = true
        }
      }
    },
    switchChange(row) {
      if (this.systemData.main == row.appid) {
        this.systemData.main = ''
      } else {
        this.systemData.main = row.appid
      }
    },
    tableRowClassName({ row, rowIndex }) {
      row.index = rowIndex
      row.listarr = JSON.parse(row.list)
    },
    returnAppNames(list) {
      const appnames = []
      for (let i = 0; i < this.activeAppList.length; i++) {
        if (list.includes(this.activeAppList[i].appid)) {
          appnames.push(this.activeAppList[i].name)
        }
      }
      return appnames
    },
    // 格式化时间
    formatTime(row, column) {
      if (row[column.property] == 0) {
        return '/'
      } else {
        return parseTime(row[column.property])
      }
    },
    parseTime(time) {
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
        main: '',
        remark: '',
        status: 1
      }
    },
    handleSystemClick(typ, data) {
      switch (typ) {
        case 'create':
          this.resetSystemData()
          this.$nextTick(() => {
            this.$refs['systemData'].clearValidate()
          })
          this.dialogSystemFormStatus = 'create'
          this.dialogSystemFormTitle = '新建系统'
          this.dialogSystemFormVisible = true
          this.$refs.multipleTable.clearSelection()
          break
        case 'update':
          this.dialogSystemFormStatus = 'update'
          this.dialogSystemFormTitle = '编辑系统'
          this.systemData.ouid = data.ouid
          this.systemData.name = data.name
          this.systemData.list = data.list
          this.systemData.main = data.main
          this.systemData.remark = data.remark
          this.systemData.listarr = JSON.parse(data.list)

          setTimeout(() => {
            this.$refs.multipleTable.clearSelection()
            if (this.systemData.listarr) {
              this.systemData.listarr.forEach(item => {
                for (let i = 0; i < this.activeAppList.length; i++) {
                  if (this.activeAppList[i].appid == item.appid) {
                    console.log(this.activeAppList[i])
                    this.activeAppList[i].disabled = false
                    this.activeAppList[i].version = item.version
                    this.$refs.multipleTable.toggleRowSelection(this.activeAppList[i], true)
                  }
                }
              })
            } else {
              this.$refs.multipleTable.clearSelection()
            }
          }, 200)
          this.dialogSystemFormVisible = true
          break
        case 'look':
          this.dialogSystemFormStatus = 'look'
          this.dialogSystemFormTitle = '查看系统'
          this.systemData = data
          this.systemData.applist = []
          // 查询应用列表
          for (let i = 0; i < this.appList.length; i++) {
            this.appList[i].main = false
            if (this.systemData.list.includes(this.appList[i].appid)) {
              const item = this.appList[i]
              if (this.appList[i].appid == this.systemData.main) {
                item.main = true
              }
              this.systemData.applist.push(item)
            }
          }
          this.dialogSystemLookVisible = true
          break
        default:
          break
      }
    },
    createSystemData() {
      this.systemData.listarr = []
      for (let i = 0; i < this.multipleAppList.length; i++) {
        const item = {
          appid: this.multipleAppList[i].appid,
          version: this.multipleAppList[i].version
        }
        this.systemData.listarr.push(item)
      }
      this.systemData.list = JSON.stringify(this.systemData.listarr)
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
      this.systemData.listarr = []
      for (let i = 0; i < this.multipleAppList.length; i++) {
        const item = {
          appid: this.multipleAppList[i].appid,
          version: this.multipleAppList[i].version
        }
        this.systemData.listarr.push(item)
      }
      this.systemData.list = JSON.stringify(this.systemData.listarr)

      this.$refs['systemData'].validate((valid) => {
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
    },
    deleteSystemData(ouid) {
      console.log('要删除的ouid:', ouid)
      this.$confirm('此操作将永久删除该系统, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteSystem(ouid).then(() => {
          this.getList()
          this.$message({
            type: 'success',
            message: '删除系统成功!'
          })
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除系统'
        })
      })
    },
    pushSystemData(row) {
      console.log('要推送的系统:', row.ouid)
      pushSystem(row).then(() => {
        this.$message({
          type: 'success',
          message: '推送系统成功!'
        })
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
