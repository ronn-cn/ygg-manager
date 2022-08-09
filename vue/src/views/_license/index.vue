<template>
  <div class="app-container">
    <el-button type="primary" @click="handleLicenseClick('create',null)">新建密钥</el-button>

    <el-table v-loading="listLoading" :data="list" border fit highlight-current-row style="width: 100%; margin-top:30px;"
      :row-class-name="tableRowClassName">
      <el-table-column align="center" label="ID" width="80">
        <template slot-scope="{row}">
          <span>{{ row.index }}</span>
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
          <div v-for="sys in getSystemName(row.permit)" :key="sys">
            <el-tag style="margin:2px 0" size="small" v-if="sys!=''">{{sys}}</el-tag>
          </div>
        </template>
      </el-table-column>
      
      <el-table-column label="创建时间" align="center" sortable prop="created_at" :formatter="formatTime"></el-table-column>
      <el-table-column label="到期时间" align="center" sortable prop="expires_at" :formatter="formatTime"></el-table-column>
      <el-table-column align="center" label="状态">
        <template slot-scope="{row}">
          <span v-if="row.status == '正常'" style="color:green;">{{ row.status }}</span>
          <span v-else style="color:red;">{{ row.status }}</span>
        </template>
      </el-table-column>
      <el-table-column align="center" label="备注">
        <template slot-scope="{row}">
          <span>{{ row.remark }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" fixed="right" width="150" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-link class="el-dropdown-link" type="primary" @click="handleLicenseClick('look',row)">查看</el-link>&nbsp;
          <el-dropdown>
            <span class="el-dropdown-link">更多<i class="el-icon-arrow-down el-icon--right"></i></span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>使失效</el-dropdown-item>
              <el-dropdown-item @click.native="deleteLicenseData(row.code)">删除</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    
    <!-- 新建设备弹窗 -->
    <el-dialog :title="dialogLicenseTitle" :visible.sync="dialogLicenseFormVisible" :close-on-press-escape="false" :close-on-click-modal="false">
      <el-form ref="licenseData" :model="licenseData" label-position="left" label-width="100px">
        <el-form-item label="密钥Code">
          <el-input v-model="licenseData.code" placeholder="不填写，默认自动生成Code" />
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
        <el-button type="primary" @click="createLicenseData()"> 确认 </el-button>
      </div>
    </el-dialog>

    
    <!-- 查看密钥弹窗 -->
    <el-dialog :visible.sync="dialogLicenseLookVisible" :title="dialogLicenseTitle">
      <el-descriptions class="margin-top" :column="3" border>
        <el-descriptions-item span="2" label="Code">
          {{ licenseData.code }}
        </el-descriptions-item>
        <el-descriptions-item label="数量（使用/总计）">
          <span>{{ licenseData.use_count+"/"+ licenseData.count }}</span>
        </el-descriptions-item>
        
        <el-descriptions-item label="状态">
          <span v-if="licenseData.status == '正常'" style="color:green;">{{ licenseData.status }}</span>
          <span v-else style="color:red;">{{ licenseData.status }}</span>
        </el-descriptions-item>
        
        <el-descriptions-item label="创建时间">
          {{ parseTime(licenseData.created_at) }}
        </el-descriptions-item>

        <el-descriptions-item label="到期时间">
          {{ parseTime(licenseData.expires_at) }}
        </el-descriptions-item>
        
        <el-descriptions-item span="3" label="备注">
          <span>{{ licenseData.remark }}</span>
        </el-descriptions-item>
        
        <el-descriptions-item span="3" label="允许系统">
          <el-table :data="licenseData.permitOptions" style="width: 100%">
            <el-table-column prop="label" label="系统名称"> </el-table-column>
            <el-table-column prop="value" label="系统ID" show-overflow-tooltip> </el-table-column>
          </el-table>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script>
import { getLicenseList, createLicense,deleteLicense } from '@/api/license'
import { getSystemList } from '@/api/system'
import { parseTime } from '@/utils'

export default {
  name: 'InlineEditTable',
  filters: { },
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
      dialogLicenseLookVisible: false,
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
          
        
        this.listLoading = false
      })
    },
    tableRowClassName({row,rowIndex}){
      row.index = rowIndex;
      row.status = "正常";
      if(row.use_count >= row.count) {  
        row.status = "满额";
      }
      if (row.expires_at <= parseInt(new Date().getTime()/1000)) {
        row.status = "已过期";
      }
      row.permitNames = [];
      for (var i = 0; i < this.systemList.length; i++){
        if(row.permit.includes(this.systemList[i].ouid)){
          row.permitNames.push(this.systemList[i].name)
        }
      }
    },
    getSystemList(){
      getSystemList().then(response => {
        this.systemList = response.data.items
      })
    },
    getSystemName(permit){
      let names = [];
      for (var i = 0; i < this.systemList.length; i++){
        if(permit.includes(this.systemList[i].ouid)){
          names.push(this.systemList[i].name)
        }
      }
      return names
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
          data.permitOptions = [];
          for (var i = 0; i < this.systemList.length; i++){
            if(data.permit.includes(this.systemList[i].ouid)){
              var opt = {
                label: this.systemList[i].name,
                value: this.systemList[i].ouid,
              }
              data.permitOptions.push(opt)
            }
          }
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
      
      if (this.licenseData.permitarr){
        this.licenseData.permit = this.licenseData.permitarr.toString()
      }
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
          this.licenseData.expires_at = parseInt(date.getTime()/1000)
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
    },
    deleteLicenseData(code){
      console.log("要删除的code",code)
      this.$confirm('此操作将永久删除该密钥, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteLicense(code).then(() => {
            this.getList()
            this.$message({
              type: 'success',
              message: '删除密钥成功!'
            });
          });
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除密钥'
        });          
      });
    },
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
