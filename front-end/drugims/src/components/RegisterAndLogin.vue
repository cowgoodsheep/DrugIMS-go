<template>
  <div class="register">
    <h2>Register</h2>
    <form @submit.prevent="registerUser">
      <div>
        <label for="username">Username:</label>
        <input type="text" v-model="username" id="username" required />
      </div>
      <div>
        <label for="password">Password:</label>
        <input type="password" v-model="password" id="password" required />
      </div>
      <button type="submit">Register</button>
    </form>
    <div v-if="errorMessage" class="error">{{ errorMessage }}</div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'RegisterAndLogin',
  data() {
    return {
      username: '',
      password: '',
      errorMessage: ''
    };
  },
  methods: {
    async registerUser() {
      try {
        const response = await axios.post('localhost:8080/user/register', {
          username: this.username,
          password: this.password
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