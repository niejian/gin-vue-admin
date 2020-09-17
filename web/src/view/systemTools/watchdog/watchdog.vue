<template>
  <el-tabs v-model="activeName" type="card" @tab-click="handleClick">
    <el-tab-pane label="环境初始化" name="init"></el-tab-pane>
      <el-form  v-show="initShow" :model="initRuleForm" :rules="initRules" demo-ruleForm label-width="350px" size="medium" ref="initRuleForm">
        <el-form-item label="IP" prop="ip">
          <el-input v-model="initRuleForm.ip" value="192.168.240.64" placeholder="192.168.240.64" style="width:300px"></el-input>
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="initRuleForm.username" style="width:300px"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="initRuleForm.password" type="password" style="width:300px"></el-input>
        </el-form-item>
        <el-form-item label="端口" prop="port">
          <el-input v-model.number="numberValidateForm.port" style="width:300px"></el-input>
        </el-form-item>
        <el-form-item label="远程路径" prop="remoteFilePath">
          <el-input v-model="initRuleForm.remoteFilePath" disabled style="width:300px" ></el-input>
        </el-form-item>
        <el-form-item>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
          <el-button type="primary" round @click="submitForm('initRuleForm')" icon="el-icon-edit">初始化</el-button>
          <el-button type="warning" round @click="resetForm('ruleForm')" icon="el-icon-arrow-right">执行命令行</el-button>
        </el-form-item>
      </el-form>
    <el-tab-pane label="配置传输" name="scpConfig"></el-tab-pane>
    <div id="config-div" v-show="configShow">config show</div>
  </el-tabs>  
</template>

<script>
  import {
    initEnv
  } from '@/api/watchdog'
  export default {
    data() {
      return {
        initRuleForm: {
          ip: '192.168.240.64',
          username: 'appadm',
          password: '',
          port: 22,
          remoteFilePath: '/data/watchDog'
        },
        initRules: {
          ip: [
            { required: true, message: '请输入登录IP', trigger: 'blur' }
          ],
          password: [
            { required: true, message: '请输入登录密码', trigger: 'blur' }
          ],
          username: [
            { required: true, message: '请输入登录账户', trigger: 'blur' }
          ],
          port: [
            { required: true, message: '请输入登录端口', trigger: 'blur' }
          ],
        },
        numberValidateForm: {
          port: 22
        },

        initShow: true,
        configShow: false,
        activeName: 'init'
      };
    },
    methods: {
      // 表达提交
      submitForm(formName) {

        this.$refs[formName].validate((valid) => {
          
          if (valid) {
            debugger
            this.initRuleForm.port = this.numberValidateForm.port
            this.doInit()
            
          } else {
            console.log('error submit!!');
            return false;
          }
        });
      },
      handleClick(tab, event) {
        let tableName = tab.name
        if (tableName == 'init') {
          this.initShow = true
          this.configShow = false
        }else if (tableName == 'scpConfig') {
          this.initShow = false
          this.configShow = true
        }
      },
      doInit() {
        this.initRuleForm.port = this.numberValidateForm.port
        initEnv(this.initRuleForm).then((resp) => {
          
        })
      }
    }
  };
</script>