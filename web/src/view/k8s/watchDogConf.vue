<template>
  <div>
    <div class="system">
      <el-form label-width="100px" :model="form" :rules="rules" ref="form" :inline="false">
        <h2>项目查询</h2><br/>
        
        <el-form-item label="命名空间" prop="ns">
          <el-select v-model="form.ns" filterable placeholder="请选择" @click="updateNsList" @change="nsChange" size="medium" style="width:500px">
            <el-option
              v-for="item in nsList"
              :key="item"
              :label="item"
              :value="item">
            </el-option>
          </el-select>
        </el-form-item>
      
        <el-form-item label="项目信息" prop="appName">
          <!-- 获取改namespace下的deploy信息，一般deploy.name 就是项目名 -->
          <el-select v-model="form.appName" filterable placeholder="请选择" 
            @change="deployChange"
            style="width:500px">
              <el-option
                v-for="item in deployList"
                :key="item"
                :label="item"
                :value="item">
              </el-option>
            </el-select>
        </el-form-item>
       
        <el-form-item label="异常信息" prop="errs">
          <common-tag :dynamicTags="form.errs"  :moduleName="errModuleName" ref="errTag"></common-tag>
        </el-form-item>

        <el-form-item label="忽略异常" prop="ignores">
          <common-tag :dynamicTags="form.ignores" :moduleName="ignoresModuleName" ref="ignoreTag"></common-tag>

        </el-form-item>

        <el-form-item label="通知人编号" prop="toUserIds">
          <common-tag :dynamicTags="form.toUserIds" :moduleName="toUserIdsModuleName" ref="toUserIdTag"></common-tag>
        </el-form-item>
        <el-form-item label="是否开启告警" icon="el-icon-question">
          <el-tooltip class="item" effect="dark" content="是否开启异常告警" placement="top-start">
            <i class="el-icon-question"></i>
          </el-tooltip>

          <!-- <i class="el-icon-question"></i> -->
          &nbsp;&nbsp;&nbsp;&nbsp;<el-switch v-model="form.isEnable"></el-switch>
        </el-form-item>
        <el-form-item label="是否存储" icon="el-icon-question">
          <el-tooltip class="item" effect="dark" content="是否开启错误信息持久化" placement="top-start">
            <i class="el-icon-question"></i>
          </el-tooltip>

          <!-- <i class="el-icon-question"></i> -->
          &nbsp;&nbsp;&nbsp;&nbsp;
          <el-switch v-model="form.enableStore"></el-switch>
        </el-form-item>


        <el-form-item>
          <el-button type="primary" @click="onSubmit('form')">立即更新</el-button>
          <!-- <el-button>取消</el-button> -->
        </el-form-item>

      </el-form>
      
    </div>
  </div>
</template>

<style scoped>
  /* .el-tag + .el-tag {
    margin-left: 10px;
  }
  .button-new-tag {
    margin-left: 10px;
    height: 32px;
    line-height: 30px;
    padding-top: 0;
    padding-bottom: 0;
  }
  .input-new-tag {
    width: 90px;
    margin-left: 10px;
    vertical-align: bottom;
  } */
</style>
<script>
import {
    listNs,
    deploys,
    addOrUpdate,
    getConfByNsAndAppName
  } from '@/api/k8sWatchdog'

import CommonTag from '@/components/common/CommonTag'

export default {
  components: { CommonTag },
  data () {
    return {
      // errs: [],
      // ignores: [],
      // toUserIds: [],
      errModuleName: "异常信息",
      ignoresModuleName: "忽略异常",
      toUserIdsModuleName: "通知人编号",
      nsList: [],
      deployList: [],
      configId: '',
      form: {
        id: -1,
        ns: '',
        appName: '',
        errs: [],
        ignores: [],
        toUserIds: [],
        enableStore: false,
        isEnable: true,
      },
      rules: {
        ns: [
          { required: true, message: '请选择命名空间', trigger: 'blur' },
        ],
        appName: [
          { required: true, message: '请选择项目', trigger: 'blur' },
        ],
        errs: [
          { required: true, message: '请填写告警异常', trigger: 'blur' },
        ],
        toUserIds: [
          { required: true, message: '请填写告警通知人编号', trigger: 'blur' },
        ]
      }
    }
  },
  methods: {
    updateNsList() {
      debugger
      // 命名空间为空，重新请求
      if (this.nsList.length == 0){
        listNs().then((resp) => {
          if (resp.code == 0) {
              // 初始化成功
              this.nsList = resp.data
            }
        })
      }
    },
    // 刷新deploy信息
    nsChange(item) {
      if (this.nsList.length == 0){
        listNs().then((resp) => {
          if (resp.code == 0) {
              // 初始化成功
              this.nsList = resp.data
            }
        })
      }
      if (item) {
        this.deployList = []
        this.deploy = ''
        this.form.appName = ''
        this.form.errs = []
        this.form.ignores = []
        this.form.toUserIds = []
        this.form.enableStore = 0
        this.form.isEnable = 1

        deploys(item).then((resp) => {
          if (resp.code == 0 && resp.data.items) {
            let datas = resp.data.items
            if (datas.length > 0) {
              var deployList = []

              datas.forEach(function(item) {
                let deployName = item.metadata.name
                deployList.push(deployName)
              })
              this.deployList = deployList
            }
          }
        })
      }
    },
    // 获取数据库中配置的信息, 有责更新，无责创建
    deployChange(item) {
      this.form.errs = []
      this.form.ignores = []
      this.form.toUserIds = []
      this.configId = ""
      this.form.enableStore=false
      this.form.isEnable=true
      this.form.appName = item

      let ns = this.form.ns
      if (ns && item) {
        getConfByNsAndAppName(ns, item).then(resp => {
          if (resp.code == 0) {
            let data = resp.data[0]
            if (data) {
              this.form.id = data.ID
              // 赋值操作
              let errs =  data.errs
              if (errs && "" != errs) {
                this.form.errs = errs.split("|")
              }
              let ignores = data.ignores
              if (ignores && "" != ignores) {
                this.form.ignores = ignores.split("|")
              }
              let toUserIds = data.toUserIds
              if (toUserIds && "" != toUserIds) {
                this.form.toUserIds = toUserIds.split("|")
              }
              this.form.enableStore = data.enableStore==1 ? true : false
              this.form.isEnable = data.isEnable == 1 ? true : false
            }
          }
        })
      }
      
    },
    onSubmit(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          let data = {}
          
          if (this.form.errs.length > 0) {
            data.errs = this.form.errs.join("|")
          }

          if (this.form.ignores.length > 0) {
            data.ignores = this.form.ignores.join("|")
          }

          if (this.form.toUserIds.length > 0) {
            data.toUserIds = this.form.toUserIds.join("|")
          }
          data.ns = this.form.ns
          data.appName = this.form.appName
          let enableStore = this.form.enableStore
          let isEnable = this.form.isEnable

          data.enableStore = enableStore ? 1 : 0
          data.isEnable = isEnable ? 1 : 0
          
          if (this.form.id && this.form.id > 0) {
            data.id = this.form.id
          }
          // 发送请求
          addOrUpdate(data).then(resp => {
            if (resp.code == 0) {
              // 请求成功
              this.$notify({
                title: '成功',
                message: '更新成功',
                type: 'success'
              });
              this.$refs['form'].resetFields();
              this.form = {
                id: -1,
                ns: '',
                appName: '',
                errs: [],
                ignores: [],
                toUserIds: [],
                enableStore: false,
                isEnable: false
              }
            }
          })
          
        } else {
          return false;
        }
      });
    },
    getErrs() {
      // 调用子组件中的方法
      // this.errs = this.$refs.errTag.getData()
      // console.log(this.errs)
    }
  },

  async created () {
    await listNs().then((resp) => {
      if (resp.code == 0) {
          // 初始化成功
          this.nsList = resp.data
        }
    })
  }
}

</script>
