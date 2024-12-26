<template>
  <div class="register">
    <h2>用户注册</h2>
    <form @submit.prevent="registerUser">
      <div>
        <label for="telephone">手机号:</label>
        <input type="text" v-model="telephone" id="telephone" required />
      </div>
      <div>
        <label for="user_name">用户名:</label>
        <input type="text" v-model="user_name" id="user_name" required />
      </div>
      <div>
        <label for="password">密码:</label>
        <input type="password" v-model="password" id="password" required />
      </div>
      <button type="submit">点击注册</button>
    </form>
    <div v-if="errorMessage" class="error">{{ errorMessage }}</div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'RegisterUser',
  data() {
    return {
      telephont: '',
      user_name: '',
      password: '',
      errorMessage: ''
    };
  },
  methods: {
    async registerUser() {
      try {
        const formData = new FormData();
        formData.append('telephone', this.telephone);
        formData.append('user_name', this.user_name);
        formData.append('password', this.password);

        const response = await axios.post('http://127.0.0.1:8080/user/register', formData, {
          headers: {
          'Content-Type': 'multipart/form-data'
          }
        });
        console.log('User registered:', response.data);
        // 可以在这里添加更多的逻辑，比如跳转到登录页面
      } catch (error) {
        this.errorMessage = error.response ? error.response.data.message : 'Registration failed';
      }
    }
  }
};
</script>

<style scoped>
.register {
  max-width: 300px;
  margin: 0 auto;
}
.error {
  color: red;
}
</style>