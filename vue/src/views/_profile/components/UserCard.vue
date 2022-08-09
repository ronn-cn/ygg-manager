<template>
  <el-card style="margin-bottom: 20px">
    <div slot="header" class="clearfix">
      <span>账号信息</span>
    </div>

    <div class="user-profile">
      <div class="box-center">
        <!-- <div class="user-name text-center">{{ user.info.name }}</div> -->
        <!-- <div class="user-role text-center text-muted">account</div> -->
        <!-- <div class="user-role text-center text-muted">{{ user.role }}</div> -->
      </div>
    </div>

    <div class="user-bio">
      <div class="user-education user-bio-section">
        <div class="user-bio-section-header">
          <svg-icon icon-class="user" /><span>账号</span>
        </div>
        <div class="user-bio-section-body">
          <div class="text-muted">{{ user.info.account }}</div>
        </div>
      </div>
      
      <div class="user-education user-bio-section">
        <div class="user-bio-section-header">
          <svg-icon icon-class="people" /><span>昵称</span>
        </div>
        <div class="user-bio-section-body">
          <div class="text-muted">{{ user.info.name }}</div>
        </div>
      </div>
      
      <div class="user-education user-bio-section">
        <div class="user-bio-section-header">
          <svg-icon icon-class="nested" /><span>类型</span>
        </div>
        <div class="user-bio-section-body">
          <div class="text-muted">{{ user.role }}</div>
        </div>
      </div>

      <div class="user-education user-bio-section">
        <div class="user-bio-section-header">
          <svg-icon icon-class="icon" /><span>OUID</span>
        </div>
        <div class="user-bio-section-body">
          <div class="text-muted">
            {{ user.info.ouid }}
          </div>
        </div>
      </div>

      <div class="user-education user-bio-section">
        <div class="user-bio-section-header">
          <svg-icon icon-class="message" /><span>联系方式</span>
        </div>
        <div class="user-bio-section-body">
          <div class="text-muted">{{ user.info.contact }}</div>
        </div>
      </div>

      <div class="user-education user-bio-section">
        <div class="user-bio-section-header">
          <svg-icon icon-class="documentation" /><span>介绍</span>
        </div>
        <div class="user-bio-section-body">
          <div class="text-muted">
            {{ user.info.detail }}
          </div>
        </div>
      </div>

      <div class="user-education user-bio-section">
        <div class="user-bio-section-header">
          <svg-icon icon-class="lock" /><span>密码</span>
        </div>
        <div class="user-bio-section-body">
          <div class="text-muted">
            <el-button type="success" plain @click="changePasswordDialog = true">修改密码</el-button>
          </div>
        </div>
      </div>
    </div>
    <el-dialog title="修改密码" :visible.sync="changePasswordDialog">
      <el-form :model="form">
        <el-form-item label="输入旧密码" :label-width="formLabelWidth">
          <el-input v-model="password_old" placeholder="请输入旧密码" show-password></el-input>
        </el-form-item>
        <el-form-item label="输入新密码" :label-width="formLabelWidth">
          <el-input v-model="password_new" placeholder="请输入新密码，请切记" show-password></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="cancelChangePassword()">取 消</el-button>
        <el-button type="primary" @click="confirmChangePassword()">确 定</el-button>
      </div>
    </el-dialog>
  </el-card>
</template>

<script>
import PanThumb from "@/components/PanThumb";
import {updateAccount} from "@/api/account"
import { sha256 } from "js-sha256";

export default {
  components: { PanThumb },
  props: {
    user: {
      type: Object,
      default: () => {
        return {
          info: "",
          role: "管理员",
        };
      },
    },
  },
  data(){
    return {
      accountData: {
        account:"",
        ouid:"",
        passwd: "",
      },
      password_old: '',
      password_new: '',
      changePasswordDialog: false,
    }
  },
  methods: {
    // 取消改变密码
    cancelChangePassword(){
      this.password_old = '',
      this.password_new = '',
      this.changePasswordDialog = false
    },
    // 确认改变密码
    confirmChangePassword(){
      if (this.password_old != "" ){
        if (this.password_old == this.password_new){
          this.accountData.account = this.user.info.account
          this.accountData.ouid = this.user.info.ouid
          this.accountData.passwd = sha256(this.password_new)
          let tempData = Object.assign({}, this.accountData);
          updateAccount(tempData).then(() => {
            this.$message({
              message: "更新账号密码成功",
              type: "success",
              duration: 2000,
            });
            this.password_old = '',
            this.password_new = '',
            this.changePasswordDialog = false
          });
        } else {
          this.$message({
            message: "两次输入密码不一致",
            type: "error",
            duration: 2000,
          });
        }
      } else {
        this.$message({
          message: "密码不能设置为空",
          type: "error",
          duration: 2000,
        });
      }
    }
  }
};
</script>

<style lang="scss" scoped>
.box-center {
  margin: 0 auto;
  display: table;
}

.text-muted {
  color: #777;
  word-wrap: break-word;
}

.user-profile {
  .user-name {
    font-weight: bold;
  }

  .box-center {
    padding-top: 10px;
  }

  .user-role {
    padding-top: 10px;
    font-weight: 400;
    font-size: 14px;
  }

  .box-social {
    padding-top: 30px;

    .el-table {
      border-top: 1px solid #dfe6ec;
    }
  }

  .user-follow {
    padding-top: 20px;
  }
}

.user-bio {
  margin: 0 20%;
  margin-top: 20px;
  color: #606266;

  span {
    padding-left: 4px;
  }

  .user-bio-section {
    font-size: 14px;
    padding: 15px 0;

    .user-bio-section-header {
      border-bottom: 1px solid #dfe6ec;
      padding-bottom: 10px;
      margin-bottom: 10px;
      font-weight: bold;
    }
  }
}
</style>
