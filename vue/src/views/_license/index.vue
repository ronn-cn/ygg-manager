<template>
  <div class="app-container">
    <el-button type="primary" @click="handleLicenseClick('create',null)">新建密钥</el-button>

    <el-table v-loading="listLoading" :data="list" border fit highlight-current-row style="width: 100%; margin-top:30px;">
      <el-table-column align="center" label="ID" width="80">
        <template slot-scope="{row}">
          <span>{{ row.id }}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="授权码">
        <template slot-scope="{row}">
          <span>{{ row.code }}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="数量（使用/总计）">
        <template slot-scope="{row}">
          <span>{{ row.use_count+"/"+ row.count}}</span>
        </template>
      </el-table-column>

      <el-table-column label="允许系统">
        <template slot-scope="{row}">
        </template>
      </el-table-column>
      
      <el-table-column label="创建时间" align="center" sortable prop="created_at" :formatter="formatTime"></el-table-column>
      <el-table-column label="到期时间" align="center" sortable prop="expires_at" :formatter="formatTime"></el-table-column>
      <el-table-column align="center" label="状态">
        
      </el-table-column>
      <el-table-column align="center" label="备注">
        <template slot-scope="{row}">
          <span>{{ row.remark }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" fixed="right" width="150" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-link class="el-dropdown-link" type="primary" @click="handleSystemClick('look',row)">查看</el-link>&nbsp;
          <!-- <el-link class="el-dropdown-link" type="primary" @click="handleSystemClick('update',row)">编辑</el-link>&nbsp; -->
          <el-dropdown>
            <span class="el-dropdown-link">更多<i class="el-icon-arrow-down el-icon--right"></i></span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>使失效</el-dropdown-item>
              <el-dropdown-item>删除</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    
      <!-- 新建设备弹窗 -->
    <el-dialog :title="dialogLicenseTitle" :visible.sync="dialogLicenseFormVisible">
      <el-form ref="licenseData" :model="licenseData" label-position="left" label-width="100px">
        <el-form-item label="密钥Code">
          <el-input v-model="licenseData.ouid" placeholder="不填写，默认自动生成Code" />
        </el-form-item>
        <el-form-item label="授权注册数量">
          <el-input v-model="licenseData.count" placeholder="请输入设备名称" />
        </el-form-item>
        <el-form-item label="授权的系统">
          <el-cascader
            placeholder="试试搜索"
            v-model="licenseData.permitarr"
            :options="optionsSystemData()"
            :props="{ multiple: true }"
            filterable
            clearable
            @change="licenseDataPermitChange"
            >
          </el-cascader>
        </el-form-item>
        <el-form-item label="到期时间">
          <el-date-picker
            v-model="licenseData.expires_value"
            type="datetime"
            placeholder="选择日期时间"
            align="right"
            :picker-options="pickerOptions">
          </el-date-picker>
        </el-form-item>
        <!-- <el-form-item label="厂商信息">

        </el-form-item> -->
        <el-form-item label="备注">
          <el-input v-model="licenseData.remark" :autosize="{ minRows: 2, maxRows: 4}" type="textarea" placeholder="请输入设备备注信息" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogLicenseFormVisible = false"> 取消 </el-button>
        <el-button type="primary" @click="dialogLicenseStatus==='create'?createLicenseData():updateLicenseData()"> 确认 </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getLicenseList, createLicense } from '@/api/license'
import { getSystemList } from '@/api/system'

export default {
  name: 'InlineEditTable',
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
      list: [],
      systemList: [],
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20
      },
      
      pickerOptions: {
        shortcuts: [{
          text: '一天',
          onClick(picker) {
            const date = new Date();
            date.setTime(date.getTime() + 3600 * 1000 * 24);
            picker.$emit('pick', date);
          }
        }, {
          text: '一个周',
          onClick(picker) {
            const date = new Date();
            date.setTime(date.getTime() + 3600 * 1000 * 24 * 7);
            picker.$emit('pick', date);
          }
        }, {
          text: '一个月',
          onClick(picker) {
            const date = new Date();
            date.setTime(date.getTime() + 3600 * 1000 * 24 * 30);
            picker.$emit('pick', date);
          }
        }, {
          text: '三个月',
          onClick(picker) {
            const date = new Date();
            date.setTime(date.getTime() + 3600 * 1000 * 24 * 90);
            picker.$emit('pick', date);
          }
        }, {
          text: '六个月',
          onClick(picker) {
            const date = new Date();
            date.setTime(date.getTime() + 3600 * 1000 * 24 * 180);
            picker.$emit('pick', date);
          }
        }]
      },

      licenseData: {
        code: '',
        count: 0,
        use_count: 0,
        permit: '',
        permitarr: [],
        permitOptions: [],
        vendor_json: '',
        expires_at: 0,
        expires_value: '',
        remark: ''
      },
      dialogLicenseStatus: '',
      dialogLicenseTitle: '',
      dialogLicenseFormVisible: false,
    }
  },
  created() {
    this.getList()
    this.getSystemList()
  },
  methods: {
    getList() {
      this.listLoading = true
      getLicenseList(this.listQuery).then(response => {
        this.list = response.data.items
        console.log("list:",this.list)
        this.total = response.data.total
        
        this.listLoading = false
      })
    },
    getSystemList(){
      getSystemList().then(response => {
        this.systemList = response.data.items
      })
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

    resetLicenseData(){
      this.licenseData = {
        code: '',
        count: 0,
        use_count: 0,
        permit: '',
        permitarr: [],
        permitOptions: [],
        vendor_json: '',
        expires_at: 0,
        expires_value: '',
        remark: ''
      }
    },
    handleLicenseClick(typ,data){
      switch(typ){
        case 'create':
          this.resetLicenseData()
          
          const date = new Date();
          date.setTime(date.getTime() + 3600 * 1000 * 24);
          this.licenseData.expires_value = date;

          this.$nextTick(() => {
            this.$refs['licenseData'].clearValidate()
          })
          this.dialogLicenseStatus = 'create'
          this.dialogLicenseTitle = "新建密钥"
          this.licenseData.count = 100
          this.dialogLicenseFormVisible = true
          break;
        case 'look':
          this.dialogLicenseStatus = 'look'
          this.dialogLicenseTitle = "查看密钥"
          this.licenseData = data
          this.dialogLicenseLookVisible = true
          break;
        default:
          break;
      }
    },
    optionsSystemData(){
      let systemOptions = []
      for (let i = 0; i < this.systemList.length; i++){
          let item = {
            value: this.systemList[i].ouid,
            label: this.systemList[i].name,
          }
          systemOptions.push(item)
      }
      return systemOptions
    },
    // 系统的应用列表数据改变时操作的函数
    licenseDataPermitChange(){
      console.log("licenseDataPermitChange:",this.licenseData.permitOptions)
      if(this.licenseData.permitOptions){
        this.licenseData.permitOptions.length = 0;
      } else {
        this.licenseData.permitOptions = [];
      }
      
      this.licenseData.permit = this.licenseData.permitarr.toString()
      for (var i = 0; i < this.systemList.length; i++){
        if(this.licenseData.permit.includes(this.systemList[i].ouid)){
          let opt = {
            label: this.systemList[i].name,
            value: this.systemList[i].ouid,
          }
          this.licenseData.permitOptions.push(opt)
        }
      }
    },

    // 创建授权数据
    createLicenseData(){
      this.$refs['licenseData'].validate((valid) => {
        if (valid) {
          var date = new Date(this.licenseData.expires_value)
          this.licenseData.expires_at = int(date.getTime()/1000)
          createLicense(this.licenseData).then(() => {
            this.getList()
            this.dialogLicenseFormVisible = false
            this.$message({
              message: '创建密钥成功',
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

<style scoped>
.edit-input {
  padding-right: 100px;
}
.cancel-btn {
  position: absolute;
  right: 15px;
  top: 10px;
}
.el-dropdown-link {
  cursor: pointer;
  color: #409EFF;
}
</style>
