<template>
  <div class="app-container">
    <el-button type="primary" @click="handleCompanyClick('create',null)">新建商户</el-button>

    <el-table v-loading="listLoading" :data="list" border fit highlight-current-row style="width: 100%; margin-top:30px;"
      :row-class-name="tableRowClassName">
      <el-table-column align="center" label="序号" width="80">
        <template slot-scope="{row}">
          <span>{{ row.index }}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="名称">
        <template slot-scope="{row}">
          <span>{{row.name}}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="类型">
        <template slot-scope="{row}">
          <span>{{row.type}}</span>
        </template>
      </el-table-column>

      <el-table-column label="地址">
        <template slot-scope="{row}">
          <span>{{row.address}}</span>
        </template>
      </el-table-column>

      <el-table-column class-name="status-col" label="电话">
        <template slot-scope="{row}">
          <span>{{row.telephone}}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="邮箱">
        <template slot-scope="{row}">
          <span>{{row.email}}</span>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" align="center" sortable prop="created_at" :formatter="formatTime"></el-table-column>
      <el-table-column label="操作" align="center" fixed="right" width="150" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-link class="el-dropdown-link" type="primary" @click="handleCompanyClick('look',row)">查看</el-link>&nbsp;
          <el-link class="el-dropdown-link" type="primary" @click="handleCompanyClick('update',row)">编辑</el-link>&nbsp;
          <el-link class="el-dropdown-link" type="primary" @click="handleCompanyClick('delete',row)">删除</el-link>&nbsp;
        </template>
      </el-table-column>
    </el-table>

    <!-- 新建商户弹窗 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogFormVisible">
      <el-form ref="companyData" :model="companyData" label-position="left" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="companyData.name" placeholder="请填写公司的名称" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="companyData.type" placeholder="请选择公司类型">
            <el-option v-for="item in companyTypeOptions" :key="item" :label="item" :value="item"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="companyData.address" placeholder="请填写公司地址" />
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="companyData.telephone" placeholder="请填写公司电话" />
        </el-form-item>
        <el-form-item label="Email">
          <el-input v-model="companyData.email" placeholder="请填写公司Email" />
        </el-form-item>
        <el-form-item label="网站">
          <el-input v-model="companyData.website" placeholder="请填写公司官网" />
        </el-form-item>
        <el-form-item label="简介">
          <el-input v-model="companyData.description" :autosize="{ minRows: 2, maxRows: 4}" type="textarea" placeholder="请填写公司简介" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false"> 取消 </el-button>
        <el-button type="primary" @click="dialogStatus==='create'?createData():updateData()"> 确认 </el-button>
      </div>
    </el-dialog>
    
    <!-- 查看商户弹窗 -->
    <el-dialog :visible.sync="dialogLookVisible" :title="dialogTitle">
      <el-descriptions class="margin-top" :column="3" border>
        <el-descriptions-item span="3" label="名称">
          {{ companyData.name }}
        </el-descriptions-item>
        <el-descriptions-item span="3" label="地址">
          <span>{{ companyData.address }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="类型">
          {{ companyData.type }}
        </el-descriptions-item>
        
        <el-descriptions-item label="电话">
          <span>{{ companyData.telephone }}</span>
        </el-descriptions-item>
        
        <el-descriptions-item label="Email">
          <span>{{ companyData.email }}</span>
        </el-descriptions-item>

        <el-descriptions-item label="网站">
          <span>{{ companyData.website }}</span>
        </el-descriptions-item>
        
        <el-descriptions-item label="创建时间">
          {{ parseTime(companyData.created_at) }}
        </el-descriptions-item>

        <el-descriptions-item label="更新时间">
          {{ parseTime(companyData.updated_at) }}
        </el-descriptions-item>
        
        <el-descriptions-item span="3" label="简介">
          <span>{{ companyData.description }}</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script>
import { getCompanyList, createCompany, updateCompany, deleteCompany } from '@/api/company'
import { parseTime } from '@/utils'

export default {
  name: 'InlineEditTable',
  filters: { },
  data() {
    return {
      list: [],
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 10
      },
      companyTypeOptions: [
        "生产商",
        "制造商",
        "经销商",
        "使用客户"
      ],

      dialogTitle: '',
      dialogStatus: '',
      dialogFormVisible: false,
      dialogLookVisible: false,
      companyData:{
        name: '',
        type: '',
        address: '',
        telephone: '',
        email: '',
        website: '',
        description: '',
      }
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true
      getCompanyList(this.listQuery).then(response => {
        this.list = response.data.items
        this.total = response.data.total
        this.listLoading = false
      })
    },
    tableRowClassName({ row, rowIndex }) {
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
    parseTime(time){
      return parseTime(time)
    },
    resetCompanyData(){
      this.companyData = {
        name: '',
        type: '',
        address: '',
        telephone: '',
        email: '',
        website: '',
        description: '',
      }
    },
    handleCompanyClick(typ, data){
      switch(typ){
        case 'create':
          this.resetCompanyData()
          this.$nextTick(() => {
            this.$refs['companyData'].clearValidate()
          })
          this.dialogStatus = 'create'
          this.dialogTitle = "新建商户"
          this.dialogFormVisible = true
          break;
        case 'update':
          this.dialogStatus = 'update'
          this.dialogTitle = "更新商户信息"
          this.companyData = data
          this.dialogFormVisible = true
          break;
        case 'look':
          this.dialogStatus = 'look'
          this.dialogTitle = "查看商户信息"
          this.companyData = data
          this.dialogLookVisible = true
          break;
        case "delete":
          deleteCompany(data.id).then(() =>{
            this.getList()
            this.$message({
              message: '删除商户成功',
              type: 'success',
              duration: 2000
            })
          })
          break;
        default:
          break;
      }
    },
    createData(){
      this.$refs['companyData'].validate((valid) => {
        if (valid) {
          createCompany(this.companyData).then(() => {
            this.getList()
            this.dialogFormVisible = false
            this.$message({
              message: '创建商户成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    updateData(){
      this.$refs['companyData'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.companyData)
          updateCompany(tempData).then(() => {
            this.getList()
            this.dialogFormVisible = false
            this.$message({
              message: '更新商户信息成功',
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
  color: #409eff;
}
</style>
