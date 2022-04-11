<template>
  <div class="app-container">
    <!-- 新建账号按钮 -->
    <el-button type="primary" @click="dialogNewAccountVisible = true">新建账号</el-button>

    <!-- 账号列表 -->
    <el-table
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%; margin-top: 30px"
      :row-class-name="tableRowClassName"
    >
      <el-table-column align="center" label="ID" width="80">
        <template slot-scope="{ row }">
          <span>{{ row.index }}</span>
        </template>
      </el-table-column>

      <el-table-column label="名称">
        <template slot-scope="{ row }">
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>

      <el-table-column align="center" label="账号">
        <template slot-scope="{ row }">
          <span>{{ row.account }}</span>
        </template>
      </el-table-column>
      <el-table-column align="center" label="密码">
        <el-button
          type="warning"
          size="mini"
          icon="el-icon-key"
          @click="dialogChangePasswdVisible = true"
        >修改</el-button>
      </el-table-column>

      <el-table-column class-name="status-col" label="类型">
        <template slot-scope="{ row }">
          <el-tag :class="'tag-'+row.type_id">
            {{ row.type.name }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column align="center" label="状态">
        <template slot-scope="{ row }">
          <el-switch
            v-model="row.status"
            active-color="#13ce66"
            inactive-color="#ff4949"
          />
        </template>
      </el-table-column>
      <el-table-column align="center" label="创建时间" prop="created_at" :formatter="formatTime"></el-table-column>
      <el-table-column align="center" label="操作">
        <template slot-scope="{ row }">
          <el-link class="el-dropdown-link" type="primary">查看</el-link>&nbsp;
          <el-link class="el-dropdown-link" type="primary">编辑</el-link>&nbsp;
          <el-dropdown>
            <span class="el-dropdown-link">
              更多<i class="el-icon-arrow-down el-icon--right"></i>
            </span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item>修改密码</el-dropdown-item>
              <el-dropdown-item>禁用/启用</el-dropdown-item>
              <el-dropdown-item v-if="info.type_id == 1&&row.type_id != 1">删除账号</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <!-- 这一段是修改密码的弹窗 -->
    <el-dialog
      title="修改密码"
      :visible.sync="dialogChangePasswdVisible  "
      width="30%"
      :before-close="closeChangePasswdForm"
    >
      <el-input
        v-model="passwd1"
        class="passwd-input"
        placeholder="请输入密码"
        show-password="true"
      />
      <el-input
        v-model="passwd2"
        class="passwd-input"
        placeholder="请确认密码"
        show-password="true"
      />
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeChangePasswdForm">取 消</el-button>
        <el-button type="primary" @click="submiteChangePasswdForm">确 定</el-button>
      </span>
    </el-dialog>

    <!-- 新建账号的弹窗 -->
    <el-dialog
      title="新建账号"
      :visible.sync="dialogNewAccountVisible"
      width="60%"
      :before-close="closePostAccountForm"
    >
      <el-form ref="accountForm" :model="accountForm" label-position="left" label-width="70px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="accountForm.name" placeholder="请输入名称"/>
        </el-form-item>
        <el-form-item label="账号" prop="account">
          <el-input v-model="accountForm.account" placeholder="请输入账号"/>
        </el-form-item>
        <el-form-item label="密码" prop="passwd">
          <el-input v-model="accountForm.passwd" placeholder="请输入密码" show-password/>
        </el-form-item>
        <el-form-item label="类型" prop="type_id">
          <el-select v-model="accountForm.type_id" class="filter-item" placeholder="请选择类型">
            <el-option v-for="item in accountTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="联系方式" prop="contact">
          <el-input v-model="accountForm.contact" placeholder="请输入联系方式"/>
        </el-form-item>
        <el-form-item label="简介" prop="detail">
          <el-input v-model="accountForm.detail" :autosize="{ minRows: 2, maxRows: 4}" type="textarea" placeholder="请输入简介" />
        </el-form-item>
        <el-form-item>
            <el-button @click="closePostAccountForm()">取 消</el-button>
            <el-button type="primary" @click="submitPostAccountForm">提 交</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>

    <!-- 查看账号信息弹窗 -->

    <!-- 编辑账号弹窗 -->

  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { getAccountList,postAccount } from '@/api/account'
export default {
  name: 'InlineEditTable',
  filters: {
    statusFilter(status) {
      const statusMap = {
        published: "success",
        draft: "info",
        deleted: "danger",
      };
      return statusMap[status];
    },
  },
  data() {
    return {
      passwd1: "",
      passwd2: "",
      accountForm:{
        name:'',
        account:'',
        passwd:'',
        type_id:4,
        contact:'',
        detail:'',
      },
      dialogChangePasswdVisible: false,
      dialogNewAccountVisible: false,
      list: null,
      listLoading: true,
      accountTypeOptions: [{
        value: 2,
        label: '管理'
      }, {
        value: 3,
        label: '开发'
      }, {
        value: 4,
        label: '运维'
      }, {
        value: 5,
        label: '生产'
      }, {
        value: 6,
        label: '销售'
      }, {
        value: 7,
        label: '客户'
      }],
    };
  },
  computed: {
    ...mapGetters([
      'info',
      'roles'
    ])
  },
  created() {
    this.getList(); 
    this.getUserInfo()
  },
  methods: {
    async getList() {
      this.listLoading = true;
      const { data } = await getAccountList();
      const items = data.items;
      this.list = items.map((v) => {
        if (v.status){
          v.status = true;
        }
        return v;
      });

      // 判断不是超级管理员，菜单中取消创建管理员
      if (this.list.length == 0) {
        this.accountTypeOptions.splice(0,1) 
      } else {
        if (this.list[0].type_id != 1){
          this.accountTypeOptions.splice(0,1) 
        }
      }
      this.listLoading = false;
    },
    getUserInfo() {
      console.log(this.info)
      this.user = {
        info: this.info,
        role: this.roles[0].name,
      }
    },
    tableRowClassName({row,rowIndex}){
      row.index = rowIndex;
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
    submitPostAccountForm() { //创建账号表单提交
      console.log(JSON.stringify(this.accountForm))
      postAccount(JSON.stringify(this.accountForm)).then((result)=>{
        console.log(result)
      }).catch((error)=>{
        console.log(error)
      })
      closePostAccountForm()
    },
    closePostAccountForm(){
      this.$refs['accountForm'].resetFields();
      this.dialogNewAccountVisible = false
    },
    submiteChangePasswdForm() {
      
      closeChangePasswdForm()
    },
    closeChangePasswdForm(){
      this.passwd1 = ''
      this.passwd2 = ''
      this.dialogChangePasswdVisible = false
    }
  },
};
</script>

<style scoped>

.el-dropdown-link {
  cursor: pointer;
  color: #409EFF;
}
.edit-input {
  padding-right: 100px;
}
.cancel-btn {
  position: absolute;
  right: 15px;
  top: 10px;
}
.passwd-input {
  margin-top: 10px;
}
/* 超级管理员 */
.tag-1 {
  background-color: #ff2d51;
  color: #fff;
}
/* 管理员 */
.tag-2 {
  background-color: #fff143;
  color: #000;
}
/* 开发 */
.tag-3 {
  background-color: #e9e7ef;
  color: #000;
}
/* 运维 */
.tag-4 {
  background-color: #fa8c35;
  color: #000;
}
/* 生产 */
.tag-5 {
  background-color: #21a675;
  color: #fff;
}
/* 销售 */
.tag-6 {
  background-color: #065279;
  color: #fff;
}
/* 客户 */
.tag-7 {
  background-color: #44cef6;
  color: #000;
}
</style>
