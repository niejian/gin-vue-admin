<template>
  <div class="console"
       id="terminal"></div>
</template>
<script>
import Terminal from '@/view/systemTools/console/Xterm'
export default {
  name: 'console',
  props: {
    terminal: {
      type: Object,
      // eslint-disable-next-line vue/require-valid-default-prop
      default: {}
    },
    url: ''
  },
  data () {
    return {
      term: null,
      terminalSocket: null
    }
  },
  methods: {
    termOnData(data){
      debugger
      // console.log("socket send:", msg.data)
      this.terminalSocket.send(data)
    },
    onmessage(msg) {
      // console.log("on message:", msg.data)
      // this.term.write(msg.data)
    },
    runRealTerminal (data) {
      debugger
      this.terminalSocket.send(data)
      console.log('webSocket is finished')
    },
    errorRealTerminal (e) {
      console.log('error', e)
    },
    closeRealTerminal () {
      console.log('close closeRealTerminal')
      this.terminalSocket.close()
      this.term.destroy()
    },
   
    initTerm(){
      let terminalContainer = document.getElementById('terminal')
      this.term = new Terminal({
        name: 'terminal',
        rendererType: "canvas", //渲染类型
        cols: 120,
        rows: 30,
        cursorStyle: 'block', //光标样式  null | 'block' | 'underline' | 'bar'
        cursorBlink: true, // 光标闪烁
        convertEol: true, //启用时，光标将设置为下一行的开头
        disableStdin: false, //是否应禁用输入。
        scrollback: 800, //回滚
        tabStopWidth: 8, //制表宽度
        screenKeys: true,//
        theme: {
          foreground: 'yellow', //字体
          background: '#060101', //背景色
          // background: 'white', //背景色
          cursor: 'help',//设置光标
        }
      })
      this.term.open(terminalContainer)
      this.term.clear()
      this.term.writeln('Welcome to xterm.js');
      this.term.writeln('This is a local terminal emulation, without a real terminal in the back-end.');
      this.term.writeln('Type some keys and commands to play around.');
      this.term.writeln('');
      this.term.clear();
      // open websocket
      console.log(this.url)
      this.terminalSocket = new WebSocket(this.url)
      this.term.ondata = this.termOnData
      this.terminalSocket.onopen = this.runRealTerminal
      this.terminalSocket.onclose = this.closeRealTerminal
      this.terminalSocket.onerror = this.errorRealTerminal
      this.terminalSocket.onmessage = this.onmessage
      this.term.attach(this.terminalSocket)
      this.term._initialized = true
      
      console.log('mounted is going on')
    }
  },
  mounted () {
    console.log('pid : ' + this.terminal.pid + ' is on ready')
    // this.initTerm()
  },
  beforeDestroy () {
    console.log("console destory")
    this.terminalSocket.close()
    this.term.destroy()
  }
}
</script>