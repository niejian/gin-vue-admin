<!-- 动态标签公共组件 -->
<template>
  <div>
    <el-tag
      :key="tag"
      v-for="tag in dynamicTags"
      closable
      :disable-transitions="false"
      @close="handleClose(tag)">
      {{tag}}
    </el-tag>
    <el-input
      class="input-new-tag"
      v-if="inputVisible"
      v-model="inputValue"
      ref="saveTagInput"
      size="small"
      @keyup.enter.native="handleInputConfirm"
      @blur="handleInputConfirm"
    >
    </el-input>
    <el-button v-else class="button-new-tag" size="small" @click="showInput">&nbsp;&nbsp;点击添加&nbsp;&nbsp;+&nbsp;&nbsp;</el-button>
  </div>
</template>

<script>
export default {
  // 数据传递
  props: [
    'dynamicTags',
    'moduleName'
    // 'inputVisible',
    // 'inputValue'
  ],
  data() {
    return {
      inputVisible: false,
      inputValue: '',
    }
  },
  methods: {
    handleClose(tag) {
      debugger
      let index = this.dynamicTags.indexOf(tag)
      this.dynamicTags.splice(index, 1);
      // 强制更新视图
      this.$forceUpdate();
      //this.$set(this.dynamicTags, index, null)
    },

    showInput() {
      debugger
      this.inputVisible = true;
      this.$nextTick(_ => {
        this.$refs.saveTagInput.$refs.input.focus();
      });
    },

    handleInputConfirm() {
      let inputValue = this.inputValue;
      // 校验数据的唯一性
      if (inputValue && this.dynamicTags.includes(inputValue)) {

        this.$notify.error({
          title: '错误',
          message: "【" + this.moduleName + '】已存在' + inputValue + ", 请重新输入！"
        });
        this.inputValue = ''
        return
      }
      if (inputValue) {
        debugger
        this.dynamicTags.push(inputValue);
      }
      this.inputVisible = false;
      this.inputValue = '';
    },

    getData() {
      return this.dynamicTags
    }
  }

}
</script>

<style scoped>
  .el-tag + .el-tag {
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
  }
</style>