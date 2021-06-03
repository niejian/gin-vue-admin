<template>
  <!--已配置信息展示-->
  <div id="app">
    <el-form 
    :model="conditions"
    ref="conditionsForm"
    label-width="120px"
    label-position="left"
    :inline=true
    >
      <el-form-item label="命名空间" >
        <el-input v-model="conditions.namespace" style="width:150px"></el-input>    
      </el-form-item>
      <el-form-item label="应用名称">
        <el-input v-model="conditions.appName" style="width:150px"></el-input>
      </el-form-item>
      <el-form-item>
          <el-button @click="query" type="primary" @keyup.enter="enterSearch">查询</el-button>
        </el-form-item>
        <el-form-item>
          
        </el-form-item>
    </el-form>

    <el-table :data="tableData" :cell-style="cellStyle" border stripe>
      <el-table-column type="index" ></el-table-column>
      <el-table-column label="id" v-if=false min-width="60" prop="ID" style="display:none"></el-table-column>
      <el-table-column label="命名空间" prop="namespace" ></el-table-column>
      <el-table-column label="应用名称" prop="appName" ></el-table-column>
      <el-table-column label="告警异常" prop="errs" >
        <template slot-scope="scope">
          <div v-for="err in scope.row.errs.split('|')" :key="err">

            <el-tag type="info" v-if="err">{{err}}</el-tag><br/>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="忽略异常" prop="ignores" >
        <template slot-scope="scope">
          <div v-for="ignore in scope.row.ignores.split('|')" :key="ignore">

            <el-tag type="warning" v-if="ignore">{{ignore}}</el-tag><br/>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="通知人" prop="toUserIds">
        <template slot-scope="scope">
          <div v-for="userId in scope.row.toUserIds.split('|')" :key="userId">
            <el-tag v-if="userId">{{userId}}</el-tag><br/>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="是否开启告警" >
        <template slot-scope="scope">
          <el-switch 
            v-model="scope.row.isEnable"
            :active-value="1"
            :inactive-value="0"
            :disabled="true"
            ></el-switch>
        </template>
      </el-table-column>
      <el-table-column label="是否持久化错误" >
        <template slot-scope="scope">
          <el-switch 
            v-model="scope.row.enableStore"
            :active-value="1"
            :inactive-value="0"
            :disabled="true"
            ></el-switch>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      :current-page="page"
      :page-size="pageSize"
      :page-sizes="[10, 30, 50, 100]"
      :style="{float:'right',padding:'20px'}"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
      layout="total, sizes, prev, pager, next, jumper"
    ></el-pagination>

  </div>
</template>

<script>
import {
    getConfList
    
  } from '@/api/k8sWatchdog'
export default {
  data() {
    return {
      conditions: {
        appName: '',
        namespace: ''
      },
      searchInfo: {},
      tableData: [],
      page: 1,
      total: 10,
      pageSize: 10,
    }
  },
  methods:{
    handleSizeChange(val) {
      this.pageSize = val
      this.getTableData()
    },
    handleCurrentChange(val) {
      this.page = val
      this.getTableData()
    },
    async getTableData(page = this.page, pageSize = this.pageSize) {
      const table = await getConfList({ page, pageSize, ...this.searchInfo })
      this.tableData = table.data.list
      this.total = table.data.total
      this.page = table.data.page
      this.pageSize = table.data.pageSize
    },
    cellStyle(column){
      // debugger
      if (column.columnIndex == 5) {
        // return "display:none"
      }
    },
    query(){
      this.searchInfo.appName = this.conditions.appName
      this.searchInfo.namespace = this.conditions.namespace
      this.getTableData()
    },
    enterSearch(){
      document.onkeydown = e =>{
        debugger
        //13表示回车键，baseURI是当前页面的地址，为了更严谨，也可以加别的，可以打印e看一下
        if (e.keyCode === 13 && e.target.baseURI.match(/k8s-watchdog-config/)) {
          //回车后执行搜索方法
          this.query()
        }
      }
    }
  },
  created() {
    this.getTableData()
    this.enterSearch()
  }
}
</script>

<style>

</style>