<template>
  <div class="app-container">
    <div class="filter-container">

      <!-- 新建设备按钮 -->
      <el-button class="filter-item" type="primary" @click="handleDeviceClick('create', null)">新建设备</el-button>

      <!-- 设备搜索 -->
      <el-input v-model="listQuery.name" placeholder="设备名称" style="width: 200px" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-input v-model="listQuery.ouid" placeholder="OUID" style="width: 200px" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-select v-model="listQuery.status_id" placeholder="状态" clearable style="width: 90px" class="filter-item" >
        <el-option v-for="item in statusOptions" :key="item.value" :label="item.text" :value="item.value" />
      </el-select>
      <el-select v-model="listQuery.system_ouid" placeholder="系统" clearable class="filter-item" style="width: 130px" >
        <el-option v-for="item in systemOptions" :key="item.text" :label="item.text" :value="item.value" />
      </el-select>
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">搜 索</el-button>
    </div>

    <!-- 设备列表 -->
    <el-table
      v-loading="listLoading"
      :data="list"
      border
      stripe
      :row-class-name="tableRowClassName"
    >
      <el-table-column label="序号" align="center" width="80">
        <template slot-scope="{ row }">
          <span>{{ row.index }}</span>
        </template>
      </el-table-column>
      <el-table-column label="设备名称" min-width="150">
        <template slot-scope="{ row }">
          <span>{{ row.name }}</span>
          <el-tag style="margin-left:5px;" v-if="row.status_id == 7" size="mini" type="info" effect="dark">开发</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="设备型号" min-width="150">
        <template slot-scope="{ row }">
          <span>{{ row.model}}</span>
        </template>
      </el-table-column>
      <el-table-column label="OUID" min-width="150" show-overflow-tooltip>
        <template slot-scope="{ row }">
          <span>{{ row.ouid }}</span>
        </template>
      </el-table-column>
      <el-table-column label="PIN码" min-width="80" align="center">
        <template slot-scope="{ row }">
          <span v-if="row.showPin">{{ row.pin }}</span>
          <span v-else>******</span>
          <!-- 修改PIN码图标按钮 -->
          <el-tooltip class="item" effect="light" content="修改PIN码" placement="bottom-start" >
            <em class="el-icon-edit"  style="cursor: pointer" @click="handleDevicePinChangeClick(row)"></em>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column label="状态" min-width="80" align="center" :filters="statusOptions" :filter-method="filterStatus">
        <template slot-scope="{ row }">
          <span>{{ row.status ? row.status.name : "" }}</span>
        </template>
      </el-table-column>
      <el-table-column label="系统" min-width="120" align="center" :filters="systemOptions" :filter-method="filterSystem" >
        <template slot-scope="{ row }">
          <span>{{ row.system ? row.system.name : "" }}</span>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" align="center" sortable prop="created_at" :formatter="formatTime" ></el-table-column>
      <el-table-column label="最后在线时间" class-name="status-col" sortable prop="last_online_time" :formatter="formatTime"></el-table-column>
      <el-table-column label="操作" align="center" fixed="right" width="150" class-name="small-padding fixed-width">
        <template slot-scope="{ row }">
          <el-link class="el-dropdown-link" type="primary" @click="handleDeviceClick('look', row)">查看</el-link>&nbsp;
          <el-link class="el-dropdown-link" type="primary" @click="handleDeviceClick('update', row)">编辑</el-link>&nbsp;
          <!-- 下拉菜单激活方式为click -->
          <el-dropdown trigger="click"> 
            <span class="el-dropdown-link">更多<i class="el-icon-arrow-down el-icon--right"></i></span>
            <el-dropdown-menu slot="dropdown"><el-dropdown-item @click.native="pushUpdateToDevice(row.ouid)">更新</el-dropdown-item>
              <el-dropdown-item v-if="row.status_id == 1" @click.native="updateDeviceStatus(row, 2)">审核</el-dropdown-item>
              <el-dropdown-item v-else-if="row.status_id != 4" @click.native="updateDeviceStatus(row, 4)">停用</el-dropdown-item>
              <el-dropdown-item @click.native="deleteDeviceData(row.ouid)">删除</el-dropdown-item>
              <el-dropdown-item>安装信息</el-dropdown-item>
              <el-dropdown-item>产品信息</el-dropdown-item>
              <!-- <el-dropdown-item>修改密码</el-dropdown-item> -->
              <!-- <el-dropdown-item>下载</el-dropdown-item> -->
            </el-dropdown-menu>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页组件 -->
    <pagination v-show="total > 0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <!-- 新建设备弹窗 -->
    <el-dialog
      :title="dialogDeviceTitle"
      :visible.sync="dialogDeviceFormVisible"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <el-form ref="deviceData" :model="deviceData" label-position="left" label-width="100px" >
        <el-form-item label="设备OUID">
          <el-input v-model="deviceData.ouid" placeholder="若不填写，默认自动生成OUID"  auto-complete="off" />
        </el-form-item>
        <el-form-item label="设备型号">
          <el-input v-model="deviceData.model" placeholder="请输入设备型号" auto-complete="off" />
        </el-form-item>
        <el-form-item label="设备名称">
          <el-input v-model="deviceData.name" placeholder="请输入设备名称" auto-complete="off" />
        </el-form-item>
        <el-form-item label="设备PIN码">
          <el-input v-model="deviceData.pin" placeholder="若不填写，默认自动生成6位PIN码" show-password  auto-complete="off" />
        </el-form-item>
        <el-form-item label="设备系统">
          <el-select v-model="deviceData.system_ouid" class="filter-item" placeholder="请选择设备系统">
            <el-option
              v-for="item in systemOptions"
              :key="item.text"
              :label="item.text"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="设备状态">
          <el-select
            v-model="deviceData.status_id"
            class="filter-item"
            placeholder="请选择设备状态"
          >
            <el-option
              v-for="item in statusOptions"
              :key="item.text"
              :label="item.text"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="deviceData.remark"
            :autosize="{ minRows: 2, maxRows: 4 }"
            type="textarea"
            placeholder="请输入设备备注信息"
          />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogDeviceFormVisible = false"> 取消 </el-button>
        <el-button type="primary" @click=" dialogDeviceStatus === 'create' ? createDeviceData() : updateDeviceData()"> 确认 </el-button>
      </div>
    </el-dialog>

    <!-- 设备查看弹窗 -->
    <el-dialog
      :visible.sync="dialogDeviceLookVisible"
      :title="dialogDeviceTitle"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <el-descriptions class="margin-top" :column="3" border>
        <el-descriptions-item label="设备名称">
          {{ deviceData.name }}
        </el-descriptions-item>
        <el-descriptions-item span="2" label="OUID">
          {{ deviceData.ouid }}
        </el-descriptions-item>
        <el-descriptions-item label="PIN码">
          <span>{{ deviceData.pin }}</span>
        </el-descriptions-item>

        <el-descriptions-item label="系统">
          <span v-if="deviceData.system != undefined">{{
            deviceData.system.name
          }}</span>
        </el-descriptions-item>

        <el-descriptions-item label="状态">
          <span v-if="deviceData.status != undefined">{{
            deviceData.status.name
          }}</span>
        </el-descriptions-item>

        <el-descriptions-item label="授权码">
          <span v-if="deviceData.license != undefined">{{
            deviceData.license.code
          }}</span>
        </el-descriptions-item>

        <el-descriptions-item label="创建时间">
          {{ parseTime(deviceData.created_at) }}
        </el-descriptions-item>

        <el-descriptions-item label="更新时间">
          {{ parseTime(deviceData.updated_at) }}
        </el-descriptions-item>

        <el-descriptions-item span="3" label="备注">
          {{ deviceData.remark }}
        </el-descriptions-item>

        <!-- <el-descriptions-item span="3" label="产品信息">
          {{ deviceData.remark }}
        </el-descriptions-item>

        
        <el-descriptions-item span="3" label="安装信息">
          {{ deviceData.remark }}
        </el-descriptions-item> -->
      </el-descriptions>
    </el-dialog>

    <!-- 这一段是修改PIN的弹窗 -->
    <el-dialog
      title="修改PIN码"
      :visible.sync="dialogChangePinVisible"
      width="30%"
      :before-close="closeChangePinForm"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <el-form>
        <el-input
          v-model="pin1"
          class="passwd-input"
          placeholder="请输入修改的PIN码"
          show-password
          maxlength="6"
          autocomplete="off"
        />
        <el-input
          v-model="pin2"
          class="passwd-input"
          placeholder="请输入刚才的PIN码"
          show-password
          maxlength="6"
          autocomplete="off"
        />
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="closeChangePinForm">取 消</el-button>
        <el-button type="primary" @click="submiteChangePinForm()">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import {
  getDeviceList,
  createDevice,
  updateDevice,
  updateDevicePin,
  deleteDevice,
  pushUpdate,
} from "@/api/device";
import { getSystemList } from "@/api/system";
import { queryDeviceStatus } from "@/api/device";
import waves from "@/directive/waves"; // waves directive
import { parseTime } from "@/utils";
import Pagination from "@/components/Pagination"; // secondary package based on el-pagination

export default {
  name: "ComplexTable",
  components: { Pagination },
  directives: { waves },
  filters: {},
  data() {
    return {
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        name: undefined,
        ouid: undefined,
        status_id: undefined,
        system_ouid: undefined,
      },
      // TODO：修改为api获取
      statusOptions: [],
      systemOptions: [], //动态获取
      dialogDeviceLookVisible: false,
      dialogDeviceFormVisible: false,
      dialogChangePinVisible: false,
      dialogDeviceStatus: "",
      dialogDeviceTitle: "",
      deviceData: {
        name: "",
        model: "",
        ouid: "",
        pin: "",
        system_ouid: undefined,
        status_id: 7,
        license: undefined,
        product_json: undefined,
        install_json: undefined,
        remark: "",
        created_at: 0,
        updated_at: 0,
        installed_at: 0,
        slod_time: 0,
        last_online_time: 0,
        last_login_time: 0,
      },
      downloadLoading: false,
      pin1: "",
      pin2: "",
      devicePinData: {
        ouid: "",
        pin: "",
      }
    };
  },
  created() {
    this.getList();
    this.getSystemList();
    this.getStatusList();
  },
  methods: {
    getList() {
      this.listLoading = true;
      getDeviceList(this.listQuery).then((response) => {
        this.list = response.data.items;
        this.total = response.data.total;

        this.listLoading = false;
      });
    },
    getSystemList() {
      getSystemList().then((response) => {
        console.log("response.data.items:", response.data.items);
        for (let i = 0; i < response.data.items.length; i++) {
          let opt = {
            value: response.data.items[i].ouid,
            text: response.data.items[i].name,
          };
          this.systemOptions.push(opt);
        }
      });
    },
    getStatusList() {
      queryDeviceStatus().then((response) => {
        console.log("response.data:", response.data);
        for (let i = 0; i < response.data.length; i++) {
          let opt = {
            value: response.data[i].id,
            text: response.data[i].name,
          };
          this.statusOptions.push(opt);
        }
      });
    },
    tableRowClassName({ row, rowIndex }) {
      row.index = rowIndex;
      row.showPin = false;
    },
    // 格式化时间
    formatTime(row, column) {
      if (row[column.property] == 0 || row[column.property] == null || row[column.property] == undefined) {
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
    parseTime(time) {
      return parseTime(time);
    },
    filterStatus(value, row) {
      return row.status_id === value;
    },
    filterSystem(value, row) {
      if (row.system) {
        return row.system.ouid === value;
      } else {
        return "";
      }
    },
    // 显示密码
    handleViewClick(row) {
      row.showPin = true;
    },
    // 处理筛选器
    handleFilter() {
      this.listQuery.page = 1;
      console.log("查询条件：", this.listQuery);
      this.getList();
    },
    handleDeviceClick(typ, data) {
      switch (typ) {
        case "create":
          this.resetDeviceData();
          this.$nextTick(() => {
            this.$refs["deviceData"].clearValidate();
          });
          this.dialogDeviceStatus = "create";
          this.dialogDeviceTitle = "新建设备";
          this.dialogDeviceFormVisible = true;
          break;
        case "update":
          this.dialogDeviceStatus = "update";
          this.dialogDeviceTitle = "编辑设备";
          this.deviceData = data;
          this.dialogDeviceFormVisible = true;
          break;
        case "look":
          this.dialogDeviceStatus = "look";
          this.dialogDeviceTitle = "查看设备";
          this.deviceData = data;
          this.dialogDeviceLookVisible = true;
          break;
        default:
          break;
      }
    },
    // 重置设备数据
    resetDeviceData() {
      this.deviceData = {
        name: "",
        ouid: "",
        pin: "",
        system_ouid: undefined,
        status_id: 7,
        license: undefined,
        product_json: undefined,
        install_json: undefined,
        remark: "",
        created_at: 0,
        updated_at: 0,
        installed_at: 0,
        slod_time: 0,
        last_time: 0,
      };
    },
    // 处理新建设备提交
    createDeviceData() {
      this.$refs["deviceData"].validate((valid) => {
        if (valid) {
          createDevice(this.deviceData).then(() => {
            this.getList();
            this.dialogDeviceFormVisible = false;
            this.$message({
              message: "创建设备成功",
              type: "success",
              duration: 2000,
            });
          });
        }
      });
    },
    // 处理更新设备提交
    updateDeviceData() {
      this.$refs["deviceData"].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.deviceData);
          updateDevice(tempData).then(() => {
            this.getList();
            this.dialogDeviceFormVisible = false;
            this.$message({
              message: "更新设备成功",
              type: "success",
              duration: 2000,
            });
          });
        }
      });
    },
    updateDeviceStatus(row, status) {
      this.deviceData = row
      this.deviceData.status_id = status; // 设置为4是停用设备
      const tempData = Object.assign({}, this.deviceData);
      updateDevice(tempData).then(() => {
        this.getList();
        this.$message({
          message: "已操作",
          type: "success",
          duration: 2000,
        });
      });
    },
    deleteDeviceData(ouid) {
      console.log("要删除的ouid:", ouid);
      this.$confirm("此操作将永久删除该设备, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => {
          deleteDevice(ouid).then(() => {
            this.getList();
            this.$message({
              type: "success",
              message: "删除设备成功!",
            });
          });
        })
        .catch(() => {
          this.$message({
            type: "info",
            message: "已取消删除设备",
          });
        });
    },
    pushUpdateToDevice(ouid) {
      pushUpdate(ouid).then(() => {
        this.$message({
          type: "success",
          message: "推送更新完成!",
        });
      });
    },
    handleDelete(row, index) {
      this.$notify({
        title: "Success",
        message: "Delete Successfully",
        type: "success",
        duration: 2000,
      });
      this.list.splice(index, 1);
    },
    handleDownload() {
      this.downloadLoading = true;
      import("@/vendor/Export2Excel").then((excel) => {
        const tHeader = ["timestamp", "title", "type", "importance", "status"];
        const filterVal = [
          "timestamp",
          "title",
          "type",
          "importance",
          "status",
        ];
        const data = this.formatJson(filterVal);
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: "table-list",
        });
        this.downloadLoading = false;
      });
    },
    formatJson(filterVal) {
      return this.list.map((v) =>
        filterVal.map((j) => {
          if (j === "timestamp") {
            return parseTime(v[j]);
          } else {
            return v[j];
          }
        })
      );
    },
    
    handleDevicePinChangeClick(data){
      this.dialogChangePinVisible = true
      this.devicePinData.ouid = data.ouid;
    },
    submiteChangePinForm() {
      if (this.pin1 == ''){
        this.$message({
          message: "不能输入空PIN码",
          type: "error",
          duration: 2000,
        });
        return
      }
      if (this.pin1 == this.pin2){
          this.devicePinData.pin = this.pin2
          let tempData = Object.assign({}, this.devicePinData);
          updateDevicePin(tempData).then(() => {
            this.getList();
            this.$message({
              message: "更新账号密码成功",
              type: "success",
              duration: 2000,
            });
          });
          this.closeChangePinForm();
      } else {
        this.$message({
          message: "两次输入的密码不一致",
          type: "error",
          duration: 2000,
        });
      }
    },
    closeChangePinForm() {
      this.pin1 = "";
      this.pin2 = "";
      this.dialogChangePinVisible = false;
    }
  },
};
</script>

<style lang="scss" scoped>
.filter-container {
  .filter-item {
    margin-right: 20px;
  }
}

.el-dropdown-link {
  cursor: pointer;
  color: #409eff;
}

.demo-table-expand {
  font-size: 0;
  .el-form-item {
    margin-right: 0;
    margin-bottom: 0;
    label {
      color: #99a9bf;
    }
  }
}

.passwd-input {
  margin-top: 10px;
}
</style>
