<template>
  <div id="app">
    <el-form ref="form" :model="form" :inline="false" :rules="rules">
      <el-form-item label="索引" label-width="100px" prop="indexName">
        <el-select v-model="form.indexName" filterable placeholder="请选择" 
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

      <el-form-item label="群聊名称" label-width="100px" prop="groupName">
        <el-input v-model="form.groupName" style="width:400px"></el-input>
      </el-form-item>

      <el-form-item label="群发成员" label-width="100px" prop="toUserIds">
        <common-tag v-bind:dynamicTags="form.toUserIds" :moduleName="toUserIdsModuleName" ref="toUserIdTag"></common-tag>
      </el-form-item>

      <el-form-item label="cron表达式" label-width="100px" prop="sendTime">
        <el-input v-model="form.sendTime" style="width:400px"></el-input> <el-link type="primary" href="https://cron.qqe2.com" target="_blank">在线验证cron表达式</el-link>
      </el-form-item>
      <el-form-item label="是否开启" label-width="100px" icon="el-icon-question">
        <!-- <i class="el-icon-question"></i> -->
        
        <el-switch v-model="form.isEnable"></el-switch>
      </el-form-item>
      <el-form-item>
          <el-input style="display:none" v-model="form.chatId"></el-input>
          <el-button type="primary" @click="onSubmit('form')">立即更新</el-button>
          <!-- <el-button>取消</el-button> -->
        </el-form-item>
        
    </el-form>
  </div>
</template>

<script>
// import { isValidCron } from 'cron-validator'
// import { cronParse } from 'cron-parser'

import {
  getExceptionView,
  addOrUpdateErrorWarn,
  getConfInfoByIndexName,
  getUserInfo
} from '@/api/exceptionView'
import CommonTag from '@/components/common/CommonTag'

export default {
  components: { CommonTag },
  data(){
    return {
      //isDisable: false,
      indexNames: [],
      toUserIdsModuleName: "群发成员",
      aggIndexs: [],
      form: {
        indexName: '',
        toUserIds: [],
        groupName: '',
        sendTime: '',
        chatId: "-1",
        isEnable: true
      },
      rules: {
        indexName: [
          { required: true, message: '请选择项目名称', trigger: 'blur' },
        ],
        groupName: [
          { required: true, message: '请输入群聊名称', trigger: 'blur' },
        ],
        sendTime: [
          { required: true, message: '请填写cron表达式', trigger: 'blur' },
          {validator: function(rule, val, callback){
            let cronParse = require('cron-parser');
            try {
              const interval = cronParse.parseExpression(val)
              console.log('cronDate:', interval.next().toDate())
              callback()
            } catch (e) {
              callback('非Cron表达式格式，请检查！' + e.message)
            }
          
          }, trigger: 'blur'
          }
        ],
        toUserIds: [
          { required: true, message: '请填写告警通知人编号', trigger: 'blur' },
          // {
          //   validator: this.validUser, trigger: 'blur'
          // },
          {validator: function(rule, value, callback){
            if (value.length < 2) {
              callback(new Error("至少填两个员工编号"))
            }else{
              callback()
            }
          }, trigger: 'blur' },
          
        ]
      }

    }
  },
  methods: {
    
    getUser(userIds){
      getUserInfo(userIds).then((resp) => {
        if (resp.code != 0) {
          return "员工号：" + userId + " 不存在，请重新输入"
        }
        return ""

      })
      
      
    },
    onSubmit(formName){
      this.$refs[formName].validate((valid) => {
        if (valid) {
         
          // let toUserIds = []
          // let userIds = ''
          let data = this.form
          data.isEnable = this.form.isEnable?1:0
          // if (this.form.toUserIds.length > 0) {
          //   toUserIds = this.form.toUserIds
          //   userIds = toUserIds.join("|")
          //   data.toUserIds = userIds 
          // }
          // data.toUserIds = data.toUserIds.join("|")
          this.addOrUpdate(data)   
        }
      })
    },
    addOrUpdate(data){
      // 创建或更新
      addOrUpdateErrorWarn(data).then((resp) => {
        if(resp.code == 0) {
          this.$notify({
            title: '成功',
            message: '更新成功',
            type: 'success'
          }); 
          this.$refs['form'].resetFields();
          this.formReset()
        }else{
          debugger
          this.$notify({
            title: '失败',
            message: resp.msg,
            type: 'error'
          }); 
        }
        
      })

    },
    formReset(){
      this.$refs['form'].resetFields();
      //this.isDisable=false
      this.form = {
        indexName: '',
        chatGroupName: '',
        sendTime: '',
        chatId: "-1",
        isEnable: true
      } 
      this.form.toUserIds = []  
    },
    indexChange(e){
      // 通过索引名找到已经创建信息
      if(e) {
        this.formReset()
        this.form.indexName = e
        getConfInfoByIndexName(e).then((resp) => {
          if(resp.code == 0 && null != resp.data) {
            //this.isDisable = true
            let userIds = resp.data.toUserIds
            this.form = resp.data
            this.form.isEnable = resp.data.isEnable == 1 ? true : false
            this.form.toUserIds = userIds.split("|")
            
          }
        })
      }
    }
  },
  async created() {
    getExceptionView().then((data) => {
      if(data.code == 0) {
        // responseData = data.data
        this.indexNames = (data.data.indexNames).sort((a, b) => a < b?-1:1)
        
      }
    })
  }
}
</script>

<style>

</style>