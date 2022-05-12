<template>
  <div class="dashboard-editor-container">
    <!-- 数字面板 -->
    <h3>数据统计</h3>
    <panel-group @handleSetLineChartData="handleSetLineChartData" :dashboardData="dashboardData" />
    <!-- 消息图表 -->
    <h3>新增趋势</h3>
    <el-row style="background: #fff; padding: 16px 16px 0; margin-bottom: 32px">
      <line-chart :chart-data="lineChartData" />
    </el-row>

    <!-- 在线列表and主机信息 -->
    <!-- <el-row :gutter="32">
      <el-col
        :xs="{ span: 24 }"
        :sm="{ span: 24 }"
        :md="{ span: 24 }"
        :lg="{ span: 12 }"
        :xl="{ span: 12 }"
        style="padding-right: 8px; margin-bottom: 30px"
      >
        <transaction-table />
      </el-col>
      <el-col
        :xs="{ span: 24 }"
        :sm="{ span: 24 }"
        :md="{ span: 24 }"
        :lg="{ span: 12 }"
        :xl="{ span: 12 }"
        style="padding-right: 8px; margin-bottom: 30px"
      >
        <transaction-table />
      </el-col>
    </el-row> -->
  </div>
</template>

<script>
// import GithubCorner from '@/components/GithubCorner'
import PanelGroup from "./components/PanelGroup";
import LineChart from "./components/LineChart";
import RaddarChart from "./components/RaddarChart";
import PieChart from "./components/PieChart";
import BarChart from "./components/BarChart";
import TransactionTable from "./components/TransactionTable";
import TodoList from "./components/TodoList";
import BoxCard from "./components/BoxCard";
import request from '@/utils/request'


export default {
  name: "DashboardAdmin",
  components: {
    // GithubCorner,
    PanelGroup,
    LineChart,
    RaddarChart,
    PieChart,
    BarChart,
    TransactionTable,
    TodoList,
    BoxCard,
  },
  data() {
    return {
      lineChartData: {},
      dashboardData: {},
    };
  },
  mounted(){
    this.getDashboardData();
  },
  methods: {
    handleSetLineChartData(type) {
      this.lineChartData = lineChartData[type];
    },
    getDashboardData() {
      request({
        url: "/dashboard/get-data",
        method: "get",
      }).then((response) => {
          console.log(response);
          this.dashboardData = response.data
          this.lineChartData = {
            deviceTrend: response.data["device_trend"].reverse(),
            systemTrend: response.data["system_trend"].reverse(),
            appTrend: response.data["app_trend"].reverse(),
          }
        })
        .catch((error) => {
          reject(error);
        });
    },
  },
};
</script>

<style lang="scss" scoped>
.dashboard-editor-container {
  padding: 32px;
  background-color: rgb(240, 242, 245);
  position: relative;

  .github-corner {
    position: absolute;
    top: 0px;
    border: 0;
    right: 0;
  }

  .chart-wrapper {
    background: #fff;
    padding: 16px 16px 0;
    margin-bottom: 32px;
  }
}

@media (max-width: 1024px) {
  .chart-wrapper {
    padding: 8px;
  }
}
</style>
