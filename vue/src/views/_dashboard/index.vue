<template>
  <div class="dashboard-container">
    <component :is="currentRole" />
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import adminDashboard from './admin'
import editorDashboard from './editor'

export default {
  name: 'Dashboard',
  components: { adminDashboard, editorDashboard },
  data() {
    return {
      currentRole: 'adminDashboard'
    }
  },
  computed: {
    ...mapGetters([
      'roles'
    ])
  },
  created() {
    if (this.roles.length > 0){
      switch (this.roles[0].id){
        case 1: // 超级管理员
        case 2: // 管理员
          this.currentRole = 'adminDashboard'
          break;
        case 3: // 开发
          break
        default:
          this.currentRole = 'editorDashboard'
          break
      }
    }
    console.log()
  }
}
</script>
