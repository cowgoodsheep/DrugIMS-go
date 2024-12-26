<template>
  <div>
    <h1>接口响应</h1>
    <p v-if="message">{{ message }}</p>
    <p v-else>加载中...</p>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';

export default {
  setup() {
    const message = ref('');

    onMounted(async () => {
      try {
        const response = await fetch('http://127.0.0.1:8080/home');
        if (response.ok) {
          const data = await response.json();
          message.value = data.msg;
        } else {
          message.value = '请求失败';
        }
      } catch (error) {
        message.value = '发生错误：' + error.message;
      }
    });

    return {
      message
    };
  }
};
</script>