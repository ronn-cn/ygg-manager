<template>
  <div class="app-container">
    <!-- 新建账号按钮 -->
    <el-button type="primary" @click="handleAccountClick('create', null)">新建账号</el-button>

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
        <template slot-scope="{ row }">
          <span>******</span>
          <el-tooltip class="item" effect="light" content="修改密码" placement="bottom-start">
            <em class="el-icon-edit" style="cursor: pointer" @click="handleAccountPasswdChangeClick(row)"></em>
          </el-tooltip>
        </template>
      </el-table-column>

      <el-table-column class-name="status-col" label="类型">
        <template slot-scope="{ row }">
          <el-tag :class="'tag tag-' + row.type.id">
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
            @change="handleAccountStatusChange(row)"
          />
        </template>
      </el-table-column>
      <el-table-column
        align="center"
        label="创建时间"
        prop="created_at"
        :formatter="formatTime"
      ></el-table-column>
      <el-table-column
        align="center"
        label="操作"
        fixed="right"
        width="150"
        class-name="small-padding fixed-width"
      >
        <template slot-scope="{ row }">
          <el-link class="el-dropdown-link" @click="handleAccountClick('look',row)"> 查看 </el-link>&nbsp;
          <el-link class="el-dropdown-link" @click="handleAccountClick('update', row)">编辑</el-link>&nbsp;
          <el-link class="el-dropdown-link" v-if="info.account != row.account" @click="deleteAccountData(row.ouid)">删除</el-link>
        </template>
      </el-table-column>
    </el-table>

    <!-- 这一段是修改密码的弹窗 -->
    <el-dialog
      title="修改密码"
      :visible.sync="dialogChangePasswdVisible"
      width="30%"
      :before-close="closeChangePasswdForm"
    >
      <el-form>
        <el-input
          v-model="passwd1"
          class="passwd-input"
          placeholder="请输入修改的密码"
          show-password
          autocomplete="off"
        />
        <el-input
          v-model="passwd2"
          class="passwd-input"
          placeholder="请输入刚才的密码"
          show-password
          autocomplete="off"
        />
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeChangePasswdForm">取 消</el-button>
        <el-button type="primary" @click="submiteChangePasswdForm(row)">确 定</el-button>
      </span>
    </el-dialog>

    <!-- 账号的编辑弹窗 -->
    <el-dialog
      :title="dialogAccountTitle"
      :visible.sync="dialogAccountFormVisible"
      width="60%"
      :close-on-press-escape="false" 
      :close-on-click-modal="false"
    >
      <el-form
        :model="accountData"
        ref="accountData"
        label-position="left"
        label-width="70px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="accountData.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item v-if="user.info.type_id != accountData.type_id" label="账号" prop="account">
          <el-input v-model="accountData.account" placeholder="请输入账号" autocomplete="off" />
        </el-form-item>
        <el-form-item label="密码" prop="passwd" v-if="dialogAccountStatus=='create'">
          <el-input v-model="accountData.passwd" placeholder="请输入密码" show-password autocomplete="off" />
        </el-form-item>
        <el-form-item  v-if="user.info.type_id != accountData.type_id" label="类型" prop="type_id">
          <el-select v-model="accountData.type_id" class="filter-item" placeholder="请选择类型">
            <el-option
              v-for="item in accountTypeOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="所属公司" prop="company">
          <el-input v-model="accountData.company" placeholder="请输入所属公司名称" />
        </el-form-item>
        <el-form-item label="联系方式" prop="contact">
          <el-input v-model="accountData.contact" placeholder="请输入联系方式" />
        </el-form-item>
        <el-form-item label="简介" prop="detail">
          <el-input
            v-model="accountData.detail"
            :autosize="{ minRows: 2, maxRows: 4 }"
            type="textarea"
            placeholder="请输入简介"
          />
        </el-form-item>
        <el-form-item>
          <el-button @click="dialogAccountFormVisible = false">取 消</el-button>
          <el-button
            type="primary"
            @click="dialogAccountStatus === 'create' ? createAccountData() : updateAccountData()">提 交</el-button
          >
        </el-form-item>
      </el-form>
    </el-dialog>

    <!-- 查看账号信息弹窗 -->

    <!-- 设备查看弹窗 -->
    <el-dialog :visible.sync="dialogAccountLookVisible" :title="dialogAccountTitle">
      <el-form label-position="left" :model="accountLookData">
        <el-form-item label="OUID">
          <span>{{ accountLookData.ouid }}</span>
        </el-form-item>
        <el-form-item label="名称">
          <span>{{ accountLookData.name }}</span>
        </el-form-item>
        <el-form-item label="用户名">
          <span>{{ accountLookData.account }}</span>
        </el-form-item>
        <el-form-item label="类型">
          <span>{{ accountLookData.type != undefined?accountLookData.type.name:"未知" }}</span>
        </el-form-item>
        <el-form-item label="联系方式">
          <span>{{ accountLookData.contact }}</span>
        </el-form-item>
        <el-form-item label="详情">
          <span>{{ accountLookData.detail }}</span>
        </el-form-item>
        <el-form-item label="所属公司">
          <span>{{ accountLookData.company }}</span>
        </el-form-item>
      </el-form>
    </el-dialog>

  </div>
</template>

<script>
import { mapGetters } from "vuex";
import { getAccountList, createAccount,updateAccount,deleteAccount } from "@/api/account";
import { sha256 } from "js-sha256";
export default {
  name: "AccountTable",
  filters: {},
  data() {
    return {
      passwd1: "",
      passwd2: "",
      list: null,
      listLoading: true,

      dialogAccountTitle: "",
      dialogAccountStatus: "",
      dialogAccountFormVisible: false,
      dialogAccountLookVisible: false,
      dialogChangePasswdVisible: false,
      accountData: {
        name: "",
        account: "",
        passwd: "",
        type_id: 4,
        contact: "",
        detail: "",
        status: 1,
        company: ""
      },
      accountLookData: {
        ouid: "",
        name: "",
        account: "",
        passwd: "",
        type: undefined,
        status: 1,
        contact: "",
        detail: "",
        company: "",
        created_at: 0,
        updated_at: 0,
      },
      accountTypeOptions: [
        {
          value: 2,
          label: "管理",
        },
        {
          value: 3,
          label: "开发",
        },
        {
          value: 4,
          label: "运维",
        },
        {
          value: 5,
          label: "生产",
        },
        {
          value: 6,
          label: "销售",
        },
        {
          value: 7,
          label: "客户",
        },
      ],
      accountPasswdData: {
        ouid: "",
        passwd: "",
      },
      companyOptions: []
    };
  },
  computed: {
    ...mapGetters(["info", "roles"]),
  },
  created() {
    this.getList();
    this.getUserInfo();
  },
  methods: {
    async getList() {
      this.listLoading = true;
      const { data } = await getAccountList();
      const items = data.items;
      this.list = items.map((v) => {
        if (v.status) {
          v.status = true;
        }
        return v;
      });
      console.log(JSON.stringify(items));

      // 判断不是超级管理员，菜单中取消创建管理员
      if (this.list.length == 0) {
        this.accountTypeOptions.splice(0, 1);
      } else {
        if (this.list[0].type_id != 1) {
          this.accountTypeOptions.splice(0, 1);
        }
      }
      this.listLoading = false;
    },
    getUserInfo() {
      console.log(this.info);
      this.user = {
        info: this.info,
        role: this.roles[0].name,
      };
    },
    tableRowClassName({ row, rowIndex }) {
      row.index = rowIndex;
    },
    // 格式化时间函数
    formatTime(row, column) {
      if (row[column.property] == 0) {
        return "/";
      } else {
        const date = new Date(row[column.property] * 1000);
        let y = date.getFullYear();
        let mo = date.getMonth() + 1;
        if (mo < 10) {
          mo = "0" + mo;
        }
        let d = date.getDate();
        if (d < 10) {
          d = "0" + d;
        }
        let h = date.getHours();
        if (h < 10) {
          h = "0" + h;
        }
        let mi = date.getMinutes();
        if (mi < 10) {
          mi = "0" + mi;
        }
        let s = date.getSeconds();
        if (s < 10) {
          s = "0" + s;
        }
        return y + "-" + mo + "-" + d + " " + h + ":" + mi + ":" + s;
      }
    },
    // 重置设备数据
    resetAccountData() {
      this.accountData = {
        name: "",
        account: "",
        passwd: "",
        type_id: 4,
        contact: "",
        detail: "",
        company_id: null,
      };
    },
    handleAccountClick(typ, data) {
      switch (typ) {
        case "create":
          console.log(this.$refs);
          this.$nextTick(() => {
            if (this.$refs["accountData"] !== undefined) {
              this.$refs["accountData"].resetFields();
            }
          });
          this.dialogAccountStatus = "create";
          this.dialogAccountTitle = "新建账号";
          this.resetAccountData();
          this.dialogAccountFormVisible = true;
          break;
        case "update":
          this.dialogAccountStatus = "update";
          this.dialogAccountTitle = "编辑账号";
          this.accountData = data;
          this.dialogAccountFormVisible = true;
          break;
        case "look":
          this.dialogAccountStatus = "look";
          this.dialogAccountTitle = "查看账号";
          this.accountLookData = data;
          this.dialogAccountLookVisible = true;
          break;
        default:
          break;
      }
    },
    createAccountData() {
      this.$refs["accountData"].validate((valid) => {
        if (valid) {
          this.accountData.passwd = sha256(this.accountData.passwd);
          createAccount(JSON.stringify(this.accountData))
            .then(() => {
              this.list.unshift(this.accountData);
              this.dialogAccountFormVisible = false;
              this.getList()
              this.$message({
                message: "创建账号成功",
                type: "success",
                duration: 2000,
              });
            })
            .catch((error) => {
              console.log(error);
            });
        }
      });
    },
    updateAccountData() {
      this.$refs["accountData"].validate((valid) => {
        if (valid) {
          let tempData = Object.assign({}, this.accountData);
          tempData.status = tempData.status ? 1:0;
          console.log("tempData:",tempData)
          updateAccount(tempData).then(() => {
            const index = this.list.findIndex(
              (v) => v.ouid === this.accountData.ouid
            );
            this.list.splice(index, 1, this.accountData);
            this.dialogAccountFormVisible = false;
            this.$message({
              message: "更新账号成功",
              type: "success",
              duration: 2000,
            });
          });
        }
      });
    },
    deleteAccountData(ouid){
      console.log("ouid1:",ouid)
      this.$confirm('此操作将永久删除该账号, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteAccount(ouid).then(() => {
            this.getList()
            this.$message({
              type: 'success',
              message: '删除账号成功!'
            });
          });
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除账号'
        });          
      });
    },
    handleAccountStatusChange(data){
      if(this.user.info.account == data.account){
        this.getList()
        this.$message({
          message: "不能禁用自己的账号",
          type: "warning",
          duration: 2000,
        });
        return
      }
      this.accountData = data;
      let tempData = Object.assign({}, this.accountData);
      tempData.status = tempData.status ? 1:0;
      console.log("tempData:",tempData)
      updateAccount(tempData).then(() => {
        const index = this.list.findIndex(
          (v) => v.ouid === this.accountData.ouid
        );
        this.list.splice(index, 1, this.accountData);
        this.$message({
          message: "更新账号状态成功",
          type: "success",
          duration: 2000,
        });
      });
    },
    handleAccountPasswdChangeClick(data){
      this.dialogChangePasswdVisible = true
      this.accountPasswdData.ouid = data.ouid;
    },
    submiteChangePasswdForm() {
      if (this.passwd1 == ''){
        this.$message({
          message: "不能输入空密码",
          type: "error",
          duration: 2000,
        });
        return
      }
      if (this.passwd1 == this.passwd2){
          this.accountPasswdData.passwd = sha256(this.passwd1)
          let tempData = Object.assign({}, this.accountPasswdData);
          updateAccount(tempData).then(() => {
            this.$message({
              message: "更新账号密码成功",
              type: "success",
              duration: 2000,
            });
          });
          this.closeChangePasswdForm();
      } else {
        this.$message({
          message: "两次输入的密码不一致",
          type: "error",
          duration: 2000,
        });
      }
    },
    closeChangePasswdForm() {
      this.passwd1 = "";
      this.passwd2 = "";
      this.dialogChangePasswdVisible = false;
    },
  },
};
</script>

<style scoped>
.el-dropdown-link {
  cursor: pointer;
  color: #409eff;
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
.tag {
  border: 0;
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
