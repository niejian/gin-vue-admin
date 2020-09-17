<template>
<!-- 近30天的错误信息 -->
  <div>
    <div>
    <el-form ref="form" >
      <el-form-item label="索引" label-width="200px">
        <el-select v-model="selectedIndexName" filterable placeholder="请选择" 
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
    </div>

    <!-- 分割线 -->
    <el-divider></el-divider>
    <!-- 统计图展示区 -->
    <div :class="className" :id="className" :style="{height:height,width:width}"></div>

    <!-- 异常显示(折叠面板) -->
    <!-- 新增角色弹窗 -->
    <el-dialog :title="dialogTitle" :visible.sync="show" width="80%" @close="closeDetail">
      <div style="width:100%;height:100%">
        
        <el-collapse v-for="item in items" :key="item.id">
          <el-collapse-item >
            <template slot="title">
              <span style="color:gray;font-size:20px">{{item.createDate}}</span>-
              <span style="color:green;font-size:20px">{{item.ip}}</span>-
              <span style="color:red;font-size:20px">{{item.exceptionTag}}</span>
                &nbsp;&nbsp;
                <i class="header-icon el-icon-question"></i>
            </template>
            <div v-html="item.msg"></div>
          </el-collapse-item>
        </el-collapse>
      </div>

    </el-dialog>

    
  </div>
 
</template>

<script>
import {
  getExceptionOverview,
  getExceptionView,
  exceptionDetails,
  indexException
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
      default: '400px'
    },
    option: {
      type: Object,
      default() {
        return {}
      }
    }
  },
  data() {
    return {
      queryFlag: false,
      dialogTitle: '异常详情',
      createDate: '',
      myChart: null,
      items: [],
      show: false,
      chart: "",
      listApi: getExceptionView,
      indexNames: [],
      selectedIndexName: '',
      aggIndexs: [],
      exceptionTags: [],
      // 柱状图点击时获取到的值
      clickVal: '',
      // echart模板数据
      options: {
        title: {
          text: ''
        },
        tooltip: {
          trigger: 'item'
        },
        legend: {
          data: [this.selectedIndexName]
        },
        xAxis: {
          name: '',
          nameLocation: 'start',
          nameGap: '50',
          boundaryGap: true,
          data: []
        },
        yAxis: { //纵轴标尺固定
          type: 'value',
          scale: true,
          name: '数量',
          min: 0,
          boundaryGap: [0.2, 0.2]
        },
        series: [{
          name: '数量',
          type: 'bar',
          data: []
        }]
      },
      requestData: {
        indexName: '',
        exceptionTag: '',
        createDate: ''
      }
    } 
  },
  methods: {
    indexChange(e){
      // e --> indexName
      this.selectedIndexName = e
      this.options.title.text = e
      this.requestData.indexName = e
      this.requestData.createDate = this.createDate
      // 初始化图表信息
      // 通过索引名称获取错误信息
      this.showExceptionOverview(this.selectedIndexName) 
    },
   
    // 获取横坐标信息
    getXData(aggIndexs) {
      let exs = []
      aggIndexs.forEach(aggIndex => exs.push(aggIndex.key))
      return exs
    },
    // 获取series.data信息
    getSeriesData(aggIndexs){
      let nums = []
      aggIndexs.forEach(aggIndex => nums.push(aggIndex.docCount))
      return nums
    },
    // 显示某天的异常图表信息
    initChart() {
      // 初始化图表

      this.myChart = echarts.init(document.getElementById('chart'), 'macarons')
      this.myChart.off('click')
      //this.myChart.hideLoading();
      // 图表设置数据
      this.myChart.setOption(this.options)
      
      // params.name ==> 横坐标
      this.myChart.on('click', (params) => {
        
        this.clickVal = params.name
        if (isNaN(this.clickVal) && !this.queryFlag) {
          this.showDetail()
        }else {
          //console.log("显示错误详情：exception", this.clickVal)
        }      
       
      })
    },
    // 显示选择的异常详细
    showDailyException(){
      this.show = false
      let data = {
        createDate: this.createDate,
        indexName: this.selectedIndexName
      }
      // 通过点击获取到创建时间和此时选中的索引名称
      indexException(data).then((resp) => { 
        if(resp.code == 0) {
          this.aggIndexs = resp.data.aggIndexs
          // 初始化横坐标-异常
          this.options.xAxis.data = this.getXData(this.aggIndexs)
          // 初始化纵坐标
          this.options.series[0].data = this.getSeriesData(this.aggIndexs)
          // 改变表头
          this.options.title.text = this.selectedIndexName + ":" + this.createDate
          this.initChart()   
        }
      });
    },
     // 显示异常详情 
    showDetail(){      
      this.requestData.indexName = this.selectedIndexName
      this.requestData.exceptionTag = this.clickVal
      this.requestData.createDate = this.createDate
      this.items = []
      
      exceptionDetails(this.requestData).then((resp) => {
        this.show = true
        if (resp.code == 0) {
          this.items = resp.data
          
          this.dialogTitle = resp.data[0].appName +' ：异常详情'
        }
      })
    },
    // 预览
    showExceptionOverview(indexName){
      getExceptionOverview({indexName: indexName}).then((data) => {
        //debugger
        // 请求成功
        if(data.code == 0) {
          // responseData = data.data
          // 重新获取到所有的索引名称
          this.indexNames = data.data.indexNames
          this.aggIndexs = data.data.aggIndexs
          if (this.indexNames) {
            this.myChart = echarts.init(document.getElementById('chart'),'macarons');
            this.myChart.off('click')
            this.selectedIndexName = indexName
            // 初始化横坐标-异常
            this.options.xAxis.data = this.getXData(this.aggIndexs)
            // 初始化纵坐标
            this.options.series[0].data = this.getSeriesData(this.aggIndexs)
            this.options.title.text = this.selectedIndexName
            
            // 初始化图表
            // this.myChart = echarts.init(document.getElementById('chart'), 'macarons')
            // 图表设置数据
            this.myChart.setOption(this.options)
            // 点击事件
            this.myChart.on('click', (params) => {              
              if (isNaN(params.name) == false) {
                  this.createDate = params.name
                  this.showDailyException()
              } 
              
            })

          }
        }                                     
      })
    },
    closeDetail() {
      this.show = false
      this.queryFlag = false
      // 清空相关数据
      this.items = []
      this.options.title.text = ''
      this.showDailyException()
    }
  },
  async created () {
    // debugger
    // 获取异步执行后的结果    
    getExceptionOverview().then((data) => {
        // 请求成功
        if(data.code == 0) {
          // responseData = data.data
          this.indexNames = data.data.indexNames
          this.aggIndexs = data.data.aggIndexs
          if (this.indexNames) {
            // this.chart = echarts.init(this.$el, 'light')
            // let myChart = echarts.init(document.getElementById('chart'),'macarons');
            this.selectedIndexName = this.indexNames[0]
            // 初始化横坐标-异常
            this.options.xAxis.data = this.getXData(this.aggIndexs)
            // 初始化纵坐标
            this.options.series[0].data = this.getSeriesData(this.aggIndexs)
            this.options.title.text = this.selectedIndexName
            
            // 初始化图表
            this.myChart = echarts.init(document.getElementById('chart'), 'macarons')
            this.myChart.off('click')
            // 图表设置数据
            this.myChart.setOption(this.options)
            // 点击事件
            this.myChart.on('click', (params) => {              
              if (isNaN(params.name) == false) {
                  this.createDate = params.name
              }  
              this.showDailyException()
            })

          }
        }                                     
    })  
  }
}
</script>