<template>
  <div>
    <el-button @click="copyLink" data-clipboard-target="#code" data-clipboard-action="copy" class="tag">复制</el-button>

    <pre v-highlightA class="codecopy">
      <code v-model="code" id="code">
        {{code}}
      </code>
    </pre>
  </div>
</template>

<script>
import Clipboard from 'clipboard'

export default {
  data() {
    return {
      code: `           
# 监控的日志文件路径, 只需要关注那些实时在变化的日志文件，归档日志不需关注
logPaths:
  - path: /data/server/website/node_8081/logs/catalina.out
    appName: buying-center-web
    toUserIds: 80468295|124
  - path: /data/server/website/node_8082/logs/catalina.out
    appName: buying-center-control
    toUserIds: 80468295|124
# 监控的日志文件是否带有日期字样
#enableLogPattern: false
#logDatePattern: YYYY-MM-dd # 日志文件日期格式，默认(支持两种格式：YYYY-MM-dd， YYYYMMdd)
emails:
  - niejian@bluemoon.com.cn
  - 393357255@qq.com
userIds:
  - 80468295
ignores:
  # - Exception
  - OverlappingFileLockException
errs:
 #- Exception
  - GlobalException
  - OverlappingFileLockException
  - ArithmeticException
  - IndexOutOfBoundsException
  - UnknownHostException
  - ConnectException
  - AxisFault
  - IndexOutOfBoundsException
  - NullPointerException
  - AxisFault
enable: true
# 是否开启端口检查（无端口检查则发送告警）
healthy:
  check: false
# 要检查的应用端口信息
apps:
  - name: scm-provider
    port: 8081
    userIds: 80468295|1234
  - name: scm-provider02
    port: 8082
    userIds: 80468295|1234
es:
  enable: true
  urls:
    - http://192.168.243.16:9200
    - http://192.168.243.17:9200
    - http://192.168.243.8:9200
  username: elastic
  password: GreatAge  
      `
    }
  },
  methods:{
     copyLink() {
       
       let clipboard = new Clipboard('.tag', {
        text: function () {
          return this.code
        }
      })

      clipboard.on('success', e => {
        this.$message({message: '复制成功', showClose: true, type: 'success'})
        // 释放内存
        clipboard.destroy()
      })
      clipboard.on('error', e => {
        this.$message({message: '复制失败,请使用chrome浏览器', showClose: true, type: 'error'})
        clipboard.destroy()
      })
     
     }
  }
}
</script>
