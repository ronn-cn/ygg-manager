<template>
  <div class="app-container">
    <div class="filter-container">
      <!-- 新建应用按钮 -->
      <el-button class="filter-item" type="primary" @click="handleApplicationClick('create', null)"> 新建应用 </el-button> 

      <el-input v-model="listQuery.name" placeholder="应用名称" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-input v-model="listQuery.appid" placeholder="APPID" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-select v-model="listQuery.type" placeholder="应用类型" clearable style="width: 120px" class="filter-item">
        <el-option v-for="item in appTypeOptions" :key="item.text" :label="item.text" :value="item.value" />
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
      <el-table-column label="应用名称" min-width="150">
        <template slot-scope="{row}">
          <span>{{ row.name }}</span>
          <el-tag style="margin-left:5px;" v-if="row.status == 0" size="mini" type="info" effect="dark">开发</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="APPID" min-width="150" align="center" show-overflow-tooltip>
        <template slot-scope="{row}">
          <span>{{ row.appid }}</span>
        </template>
      </el-table-column>
      <el-table-column label="应用类型" min-width="110" align="center">
        <template slot-scope="{row}">
          <span>{{ getAppTypeText(row.type) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="最新版本" min-width="110" align="center">
        <template slot-scope="{row}">
          <span style="color:red;">{{ row.latest==''?'无':row.latest }}</span>
        </template>
      </el-table-column>

      <el-table-column label="创建时间" align="center" sortable prop="created_at" :formatter="formatTime"></el-table-column>

      <el-table-column label="操作" align="center" fixed="right" width="150" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-link class="el-dropdown-link" type="primary" @click="handleApplicationClick('look',row)">查看</el-link>&nbsp;
          <el-link class="el-dropdown-link" type="primary" @click="handleApplicationClick('update',row)">编辑</el-link>&nbsp;
          <el-dropdown>
            <span class="el-dropdown-link">更多<i class="el-icon-arrow-down el-icon--right"></i></span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>上传版本</el-dropdown-item>
              <el-dropdown-item v-if="row.status == 0">发布应用</el-dropdown-item>
              <el-dropdown-item @click.native="deleteAppData(row.appid)">删除应用</el-dropdown-item>
              <el-dropdown-item></el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <!-- 添加编辑应用弹窗 -->
    <el-dialog :title="dialogApplicationFormTitle" :visible.sync="dialogApplicationFormVisible" :close-on-press-escape="false" :close-on-click-modal="false">
      <el-form ref="appData" :model="appData" label-position="center" label-width="100px">
        <el-form-item label="Appid" prop="appid">
          <el-input v-model="appData.appid" placeholder="不填写，默认自动生成Appid"/>
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="appData.name" placeholder="请填写应用名称"/>
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="appData.type" class="filter-item" placeholder="Please select">
            <el-option v-for="item in appTypeOptions" :key="item.text" :label="item.text" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="依赖" prop="depend">
          <el-cascader
            placeholder="试试搜索：**应用"
            v-model="appData.dependarr"
            :options="optionsAppData(appData.appid)"
            :props="{ multiple: true }"
            separator=","
            filterable
            clearable
            >
          </el-cascader>
        </el-form-item>
        <el-form-item v-if="info.type_id == 3">
          <el-checkbox v-model="appData.status">发布为正式应用</el-checkbox>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="appData.remark" :autosize="{ minRows: 2, maxRows: 4}" type="textarea" placeholder="请填写备注信息" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogApplicationFormVisible = false"> 取消 </el-button>
        <el-button type="primary" @click="dialogApplicationFormStatus==='create'?createAppData():updateAppData()"> 提交 </el-button>
      </div>
    </el-dialog>

    <!-- 应用查看弹窗 TODO：暂时先弹窗查看，后面改为单页面查看，显示应用版本列表-->
    <el-dialog :title="dialogApplicationFormTitle" :visible.sync="dialogApplicationLookVisible">
      <el-descriptions class="margin-top" :column="3" border>
        <el-descriptions-item>
          <template slot="label">应用名称</template>
          {{ appData.name }}
        </el-descriptions-item>
        <el-descriptions-item span="2">
          <template slot="label">APPID</template>
          {{ appData.appid }}
        </el-descriptions-item>
        
        <el-descriptions-item>
          <template slot="label">类型</template>
          <span v-if="appData.type == 0">程序</span>
          <span v-else-if="appData.type == 1">视图</span>
          <span v-else-if="appData.type == 2">资源</span>
          <span v-else-if="appData.type == 3">服务</span>
        </el-descriptions-item>
        <el-descriptions-item span="2">
          <template slot="label">最新版本</template>
          {{ appData.latest }}
        </el-descriptions-item>

        <el-descriptions-item>
          <template slot="label">状态</template>
          {{ appData.status?'正式':'开发' }}
        </el-descriptions-item>
        <el-descriptions-item>
          <template slot="label">创建时间</template>
          {{ parseTime(appData.created_at) }}
        </el-descriptions-item>
        
        <el-descriptions-item>
          <template slot="label">更新时间</template>
          {{ parseTime(appData.updated_at) }}
        </el-descriptions-item>
        <el-descriptions-item span="3" label="备注">
          {{ appData.remark }}
        </el-descriptions-item>

        <el-descriptions-item span="3" label="依赖列表">
          <el-table :data="appData.dependApplist" style="width: 100%">
            <el-table-column prop="name" label="应用名称"> </el-table-column>
            <el-table-column prop="appid" label="应用ID" show-overflow-tooltip> </el-table-column>
            <el-table-column prop="latest" label="最新版本"> </el-table-column>
          </el-table>
        </el-descriptions-item>
      </el-descriptions>

      <h3>版本列表</h3>
      <el-table :data="appData.versionList" style="width: 100%">
        <el-table-column prop="version" label="版本号"> </el-table-column>
        <el-table-column prop="method" label="方法"> </el-table-column>
        <el-table-column prop="description" label="描述"> </el-table-column>
        <!-- <el-table-column prop="time" label="时间"> </el-table-column> -->
      </el-table>
    </el-dialog>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
import { getApplicationList, createApplication, updateApplication, deleteApplication } from '@/api/application'
import { getVersionList } from '@/api/version'
import waves from '@/directive/waves' // waves directive
import { parseTime } from '@/utils'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination
import { get } from "js-cookie";

export default {
  name: 'ComplexTable',
  components: { Pagination },
  directives: { waves },
  filters: { },
  computed: {
    ...mapGetters(["info", "roles"]),
  },
  data() {
    return {
      list: [],
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        name: undefined,
        appid: undefined,
        type: undefined,
      },
      // 应用类型选项
      appTypeOptions: [
        { value: 0, text: '程序' },
        { value: 1, text: '视图' },
        { value: 2, text: '资源' },
        { value: 3, text: '服务' },
      ],
      appData: {
        appid: '',
        name: '',
        type: 0,
        latest: '',
        status: 1,
        depend:'',
        dependarr:[],
      },
      dialogApplicationFormStatus: '',      // 应用表单弹窗状态
      dialogApplicationFormTitle: '',       // 应用表单弹窗标题
      dialogApplicationFormVisible: false,  // 应用添加/编辑弹窗标识
      dialogApplicationLookVisible: false,  // 应用查看弹窗标识  
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true
      getApplicationList(this.listQuery).then(response => {
        this.list = response.data.items
        this.total = response.data.total
        this.listLoading = false
      })
    },
    async getAppVersionList(appid){
      return await getVersionList({"appid": appid}).then(response => {
        var redata = response.data.items.reverse()
        if (response.data.total > 10){
          return redata.slice(0, 9)
        }
        return redata
      })
    },
    tableRowClassName({row,rowIndex}){
      row.index = rowIndex;
    },
    parseTime(time){
      return parseTime(time)
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
    getAppTypeText(typ){
      for (let i = 0; i < this.appTypeOptions.length; i++){
        if (this.appTypeOptions[i].value == typ){
          return this.appTypeOptions[i].text;
        }
      }
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    resetAppData() {
      this.appData = {
        appid: '',
        name: '',
        type: 0,
        latest: '',
        status: 1,
      }
    },
    optionsAppData(appid){
      let appOptions = []
      for (let i = 0; i < this.list.length; i++){
        if (this.list[i].appid != appid){
          let item = {
            value: this.list[i].appid,
            label: this.list[i].name+"["+this.list[i].appid+"]",
          }
          appOptions.push(item)
        }
      }
      return appOptions
    },
    // 处理应用管理的点击
    async handleApplicationClick(typ,data){
      switch(typ){
        case 'create':
          this.resetAppData()
          this.$nextTick(() => {
            this.$refs['appData'].clearValidate()
          })
          this.dialogApplicationFormStatus = 'create'
          this.dialogApplicationFormTitle = "新建应用"
          this.dialogApplicationFormVisible = true
          break;
        case 'update':
          this.dialogApplicationFormStatus = 'update'
          this.dialogApplicationFormTitle = "编辑应用"
          this.appData = data
          this.appData.dependarr = this.appData.depend.split(',')
          this.dialogApplicationFormVisible = true
          break;
        case 'look':
          this.dialogApplicationFormStatus = 'look'
          this.dialogApplicationFormTitle = "查看应用"
          this.appData = data
          this.appData.dependApplist = [];
          this.appData.versionList = await this.getAppVersionList(this.appData.appid);
          console.log("ccc: ",this.appData.versionList);
          for (var i = 0; i < this.list.length; i++) {
            if (this.appData.depend.includes(this.list[i].appid)){
              this.appData.dependApplist.push(this.list[i]);
            }
          }
          this.dialogApplicationLookVisible = true
          break;
        default:
          break;
      }
    },
    createAppData() {
      this.$refs['appData'].validate((valid) => {
        if (valid) {
          console.log("this.appData:",this.appData)
          if (this.appData.dependarr){
            this.appData.depend = this.appData.dependarr.toString()
          }
          createApplication(this.appData).then(() => {
            this.getList()
            this.dialogApplicationFormVisible = false
            this.$message({
              message: '新建应用成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    updateAppData() {
      this.$refs['appData'].validate((valid) => {
        if (valid) {
          console.log("this.appData:",this.appData)
          if (this.appData.dependarr){
            this.appData.depend = this.appData.dependarr.toString()
          }
          const tempData = Object.assign({}, this.appData)
          updateApplication(tempData).then(() => {
            this.getList()
            this.dialogApplicationFormVisible = false
            this.$message({
              message: '更新应用成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    deleteAppData(appid){
      console.log("要删除的appid:",appid)
      this.$confirm('此操作将永久删除该应用, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteApplication(appid).then(() => {
            this.getList()
            this.$message({
              type: 'success',
              message: '删除应用成功!'
            });
          });
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除应用'
        });          
      });
    },
    handleFetchPv(pv) {
      fetchPv(pv).then(response => {
        this.pvData = response.data.pvData
        this.dialogPvVisible = true
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
