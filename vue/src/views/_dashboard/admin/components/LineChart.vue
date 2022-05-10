<template>
  <div :class="className" :style="{height:height,width:width}" />
</template>

<script>
import echarts from 'echarts'
require('echarts/theme/macarons') // echarts theme
import resize from './mixins/resize'

function getday2() {
  let days = [];
  var date = new Date();
  for(let i=0; i<=24*6;i+=24){		//今天加上前6天
    let dateItem=new Date(date.getTime() - i * 60 * 60 * 1000);
    let d= dateItem.getDate();	//获取日期
    d = addDate0(d);	//给为单数的日期补零
    let valueItem = d;	//组合
    days.push(valueItem);	//添加至数组
  }
  return days;		
}

//给日期加0
function addDate0(time) {
  if (time.toString().length == 1) {
    time = '0' + time.toString();
  }
  return time;
}

export default {
  mixins: [resize],
  props: {
    className: {
      type: String,
      default: 'chart'
    },
    width: {
      type: String,
      default: '100%'
    },
    height: {
      type: String,
      default: '350px'
    },
    autoResize: {
      type: Boolean,
      default: true
    },
    chartData: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      chart: null
    }
  },
  watch: {
    chartData: {
      deep: true,
      handler(val) {
        this.setOptions(val)
      }
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.initChart()
    })
  },
  beforeDestroy() {
    if (!this.chart) {
      return
    }
    this.chart.dispose()
    this.chart = null
  },
  methods: {
    initChart() {
      this.chart = echarts.init(this.$el, 'macarons')
      this.setOptions(this.chartData)
    },
    setOptions({ deviceTrend } = {}) {
      this.chart.setOption({
        xAxis: {
          data: getday2().reverse(),
          boundaryGap: false,
          axisTick: {
            show: false
          }
        },
        grid: {
          left: 10,
          right: 10,
          bottom: 20,
          top: 30,
          containLabel: true
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross'
          },
          padding: [5, 10]
        },
        yAxis: {
          axisTick: {
            show: false
          }
        },
        legend: {
          data: ['设备']
        },
        series: [{
          name: '设备', itemStyle: {
            normal: {
              color: '#5a5f6b',
              lineStyle: {
                color: '#5a5f6b',
                width: 2
              }
            }
          },
          smooth: true,
          type: 'line',
          data: deviceTrend,
          animationDuration: 2800,
          animationEasing: 'cubicInOut'
        }]
      })
    }
  }
}
</script>
