<template>
<!-- 聚合查询错误信息 -->
  <div>
    <el-form ref="form" >
      <el-form-item label="索引" label-width="200px">
        <el-select v-model="value" placeholder="请选择" 
          style="width:400px"
          size="medium" @change="indexChange">
          <el-option
            v-for="indexName in indexNames" 
            :key="indexName"
            :label="indexName"
            :value="indexName"
            >
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>

    <!-- 分割线 -->
    <el-divider></el-divider>
    <!-- 统计图展示区 -->
    <div :class="className" :style="{height:height,width:width}" />

    
  </div>
</template>

<script>
import {
  getExceptionView
} from '@/api/exceptionView'
import echarts from 'echarts'
require('echarts/theme/macarons') // echarts theme

export default {
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
          default: '300px'
      }
  },
  data() {
    return {
      listApi: getExceptionView,
      value: '',
      indexNames: [],
      aggIndexs: []
    }
   
  },
  methods: {
    indexChange(e, value){
      debugger
      console.log(e)
      console.log(value)
    },
    // 图表初始化
    initChart(){
      this.chart = echarts.init(this.$el, 'light')
      this.chart.setOption({
        tooltip: {
          trigger: 'axis',
          axisPointer: { // 坐标轴指示器，坐标轴触发有效
              type: 'shadow' // 默认为直线，可选为：'line' | 'shadow'
          }
        },
        legend: {
          data: this.indexNames
        },
      })  
    }
  },
  async created () {
    // debugger
    // 获取异步执行后的结果
    getExceptionView().then((data) => {
        // 请求成功
        if(data.code == 0) {
          
          // responseData = data.data
          this.indexNames = data.data.indexNames
          this.aggIndexs = data.data.aggIndexs
          if (this.indexNames) {
            this.value = this.indexNames[0]
          }

        }                                     
8    })  
  }
}
</script>