<template>
  <div class="app-container">
      <el-form ref="setting" :model="setting" label-position="left" label-width="100px">
        <el-form-item label="OUID" prop="ouid">
          <el-input v-model="setting.ouid" placeholder="本机的OUID"/>
        </el-form-item>
        <el-form-item label="PIN" prop="pin">
          <el-input v-model="setting.pin" show-password  placeholder="本机的PIN码"></el-input>
        </el-form-item>
        <el-form-item label="私钥" prop="prikey">
          <el-input v-model="setting.prikey" show-password  placeholder="本程序的私有密钥"/>
        </el-form-item>
        <el-form-item label="公钥" prop="pubkey">
          <el-input v-model="setting.pubkey" placeholder="本程序的公共密钥"/>
        </el-form-item>
        <el-form-item label="端口号" prop="port">
          <el-input v-model="setting.port" placeholder="默认运行端口8800"/>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="resetSetting()"> 恢复默认 </el-button>
        <el-button type="primary" @click="saveData()"> 保存 </el-button>
      </div>
  </div>
</template>

<script>
import { getSetting, saveSetting } from '@/api/setting'
import waves from '@/directive/waves' // waves directive
import Pagination from '@/components/Pagination' // secondary package based on el-pagination


export default {
  name: 'ComplexTable',
  components: { Pagination },
  directives: { waves },
  filters: {
    statusFilter(status) {
      const statusMap = {
        published: 'success',
        draft: 'info',
        deleted: 'danger'
      }
      return statusMap[status]
    },
    typeFilter(type) {
      return calendarTypeKeyValue[type]
    }
  },
  data() {
    return {
      initLoading: true,
      setting: {
        ouid: '',
        pin: '',
        prikey: '',
        pubkey: '',
        port: '8800'
      },
    }
  },
  created() {
    this.init()
  },
  methods: {
    init() {
      this.initLoading = true
      getSetting().then(response => {
        this.setting = response.data
        console.log("setting:",this.setting)
        this.initLoading = false
      })
    },
    resetSetting() {
      this.setting.port = '8800'
    },
    saveData() {
      this.$refs['setting'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.setting)
          saveSetting(tempData).then(() => {
            this.$message({
              message: '保存成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
  }
}
</script>

<style lang="scss" scoped>
.filter-container{
  .filter-item{
    margin-right: 20px;
  }
}
</style>
