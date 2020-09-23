<template>
  <el-tabs v-model="activeName" type="card" @tab-click="handleClick">
    <el-tab-pane label="环境初始化" name="init">
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
          <el-button type="warning" round @click="exec('initRuleForm')" ref="console" icon="el-icon-arrow-right">执行命令行</el-button>
        </el-form-item>
      </el-form>
    </el-tab-pane>
    <el-dialog :title="title" :visible.sync="showTerminal" width="80%" @close="close" @open="open">
      <my-terminal :terminal="terminal" :url="url" ref="console"></my-terminal>
    </el-dialog>

    <el-dialog title="watchDog.yaml" :visible.sync="showdemo" width="80%" >
      <config-demo></config-demo>
    </el-dialog>
      
      
    <el-tab-pane label="配置传输" name="scpConfig">
      <el-form :model="configRuleForm" :rules="configFormRule" label-width="350px" size="medium" ref="configRuleForm">
        <el-form-item label="IP" prop="cip">
          <el-input v-model="configRuleForm.cip" placeholder="192.168.240.64" style="width:300px"></el-input>
        </el-form-item>
        <el-form-item label="用户名" prop="cusername">
          <el-input v-model="configRuleForm.cusername" style="width:300px"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="cpassword">
          <el-input v-model="configRuleForm.cpassword" type="password" style="width:300px"></el-input>
        </el-form-item>
        <el-form-item label="端口" prop="cport">
          <el-input v-model.number="configRuleForm.cport" style="width:300px"></el-input>
        </el-form-item>
        <el-form-item label="远程路径" prop="cremoteFilePath">
          <el-input v-model="configRuleForm.cremoteFilePath" disabled style="width:300px" ></el-input>
        </el-form-item>
        <el-form-item label="文件路径" prop="localFile">
          <!-- <el-input type="file"  v-model="configRuleForm.localFile" style="width:300px"></el-input> -->
          <el-upload
            style="width: 350px;float:left"
            class="upload-demo"
            ref="upload"
            action=""
            :http-request="httpRequest"
            :list-type="listType"
            :before-upload="beforeUpload"
            :limit="limit"
            :auto-upload="autoUpload">
            <el-button slot="trigger" size="small" type="primary">选取文件</el-button>
            <el-button style="margin-left: 10px;" size="small" type="info" @click="configDemo">配置示例</el-button>
            <el-button style="margin-left: 10px;" size="small" type="danger" @click="download">模板下载</el-button>
            <!-- <el-link type="primary" @click="download" download="">配置模板下载</el-link> -->
            <div slot="tip" class="el-upload__tip" style="width:250px">请上传名为watchDog.yaml的配置文件</div>
          </el-upload>

        </el-form-item>
        <!-- <el-link type="primary" @click="download" download="">配置模板下载</el-link> -->
      
        <el-form-item>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
          <el-button type="success" round @click="submitUpload" icon="el-icon-upload">配置上传</el-button>
         
        </el-form-item>

      </el-form>
    </el-tab-pane>
  </el-tabs>  
</template>

<script>
  import {
    initEnv,
    try2Connect,
    downloadConfig
  } from '@/api/watchdog'
  import Console from '@/view/systemTools/console/console'
  import ConfigDemo from './configDemo.vue'
  import { store } from '@/store/index'
  import axios from 'axios'; 
  export default {
    data() {
      
      return {
        limit: 1,
        autoUpload: false,
        listType: 'text',
        title: '',
        url: "",
        showTerminal: false,
        showdemo: false,
        terminal: {
          pid: 1,
          name: 'terminal',
          // cols: 400,
          // rows: 400,
          cursorStyle: 'underline', //光标样式
          cursorBlink: true, // 光标闪烁
          convertEol: true, //启用时，光标将设置为下一行的开头
          disableStdin: false, //是否应禁用输入。
          theme: {
            foreground: 'yellow', //字体
            background: '#060101', //背景色
            // background: 'white', //背景色
            cursor: 'help',//设置光标
          }
        },
        initRuleForm: {
          ip: '192.168.240.64',
          username: 'appadm',
          password: 'bluemoon2016#',
          port: 22,
          remoteFilePath: '/data/watchDog'
        },
        configRuleForm:{
          cip: '192.168.240.64',
          cusername: 'appadm',
          cpassword: 'bluemoon2016#',
          cport: 22,
          localFile: '',
          cremoteFilePath: '/data/watchDog'
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
        configFormRule: {
          cip: [
            { required: true, message: '请输入登录IP', trigger: 'blur' }
          ],
          cpassword: [
            { required: true, message: '请输入登录密码', trigger: 'blur' }
          ],
          cusername: [
            { required: true, message: '请输入登录账户', trigger: 'blur' }
          ],
          cport: [
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
    components: {
      "my-terminal": Console,
      "config-demo": ConfigDemo
    },
    methods: {
      // 表达提交
      submitForm(formName) {

        this.$refs[formName].validate((valid) => {
          
          if (valid) {
            if (formName == 'initRuleForm') {
              this.initRuleForm.port = this.numberValidateForm.port
              this.doInitEnv()
            }else if (formName == 'configRuleForm') {
              this.configRuleForm.cport = this.numberValidateForm.port
            }     
          } else {
              console.log('error submit!!');
              return false;
            }
        });
      },
      // param是自带参数。 this.$refs.upload.submit() 会自动调用 httpRequest方法.在里面取得file
      httpRequest(param){
        // console.log(param)
        // debugger
        let fileObj = param.file // 相当于input里取得的files
        let fd = new FormData()// FormData 对象
        fd.append('file', fileObj)// 文件对象

        fd.append('username', this.configRuleForm.cusername)
        fd.append('password', this.configRuleForm.cpassword)
        fd.append('ip', this.configRuleForm.cip)
        fd.append('port', this.configRuleForm.cport)
        const token = store.getters['user/token']
        const user = store.getters['user/userInfo']
        
        // 文件上传
        let config = {
          headers: {
            'Content-Type': 'multipart/form-data',
            'x-token': token,
            'x-user-id': user.ID
          }
        }

        // 请求地址
        let url = process.env.VUE_APP_BASE_API + '/watchdog/upload'

        axios.post(url, fd, config).then(res => {
          if(res.data.code === 0){
            // 提交成功
            this.$message({
              message: '操作成功',
              type: 'success'
            });
            // 清空列表
            this.$refs.upload.clearFiles()
          }else {
            // 提交失败
            this.$message({
              message: res.data.msg,
              type: 'error'
            });
          }
        })
      },
      beforeUpload(file) {     
        console.log(file)
      },
      submitUpload() {
        // console.log(this.$refs)  
        this.$refs.upload.submit();
      },
      download() {
        debugger
        downloadConfig(this.configRuleForm).then(res => {
          // 判断是否连接成功
          if (res.data.code == 7){
            this.$message({
              message: res.data.msg,
              type: 'error'
            });
            return
          }
          // 从response中获取到文件信息并下载下来
          let blob = new Blob([res.data], { type: 'application/octet-stream' })
          let href = window.URL.createObjectURL(blob)

          const link = document.createElement('a')
          link.style.display = 'none'
          link.href = href
          link.download = 'watchDog.yaml'
          document.body.appendChild(link)
          link.click()
          document.body.removeChild(link)
          window.URL.revokeObjectURL(link)
        })
      },
      // 配置示例
      configDemo(){
        this.showdemo = true

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
      doInitEnv() {
        this.initRuleForm.port = this.numberValidateForm.port
        initEnv(this.initRuleForm).then((resp) => {
          if (resp.code == 0) {
            // 初始化成功
            this.$message({
              message: '初始化成功',
              type: 'success'
            });
          }
        })
      },

      // 执行命令行
      exec(formName) {
        this.$refs[formName].validate((valid) => {

          if (valid) {
            if ('initRuleForm' == formName) {
              this.initRuleForm.port = this.numberValidateForm.port
              this.showTerminal = true
              // 获取当前访问地址
              let addr = window.location.href
              // console.log(addr)
              if (addr) {
                let ns = addr.split(":")[1] // ==> //localhost
                this.url = "ws:"+ns+":8888/ws/1?cols=188&rows=50&username="+this.initRuleForm.username+"&host="+this.initRuleForm.ip + "&password="+escape(this.initRuleForm.password)  + "&port=22"

              }

              
            }            
          }
          
        })

      },
      open(){
        try2Connect(this.initRuleForm).then( resp => {
          //this.terminalFunc(this.initRuleForm)
          // 重置终端
          this.title = this.initRuleForm.ip
          this.$refs.console.initTerm()      
          
        })
      },
      // 终端相关
      close() {
        this.url = ''
        // this.$refs.console.closeRealTerminal()
        this.$refs.console.closeRealTerminal()
      }
    }
  };
</script>