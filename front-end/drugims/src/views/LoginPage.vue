<template>
    <div class="auth-container">
        <!-- 登录 -->
        <el-form v-if="isLogin" :model="loginForm" class="form-container">
            <h3>用户登录</h3>
            <el-form-item label="手机号码:">
                <el-input v-model="loginForm.telephone" placeholder="请输入手机号码"></el-input>
            </el-form-item>
            <el-form-item label="密码:">
                <el-input v-model="loginForm.password" type="password" placeholder="请输入密码"></el-input>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="sendLoginData">登录</el-button>
                <el-button @click="goToRegister">注册</el-button>
            </el-form-item>
        </el-form>

        <!-- 注册 -->
        <el-form v-else :model="registerForm" class="form-container">
            <h3>用户注册</h3>
            <el-form-item label="姓名:">
                <el-input v-model="registerForm.user_name" placeholder="请输入用户名"></el-input>
            </el-form-item>
            <el-form-item label="密码:">
                <el-input v-model="registerForm.password" type="password" placeholder="请输入密码"></el-input>
            </el-form-item>
            <el-form-item label="手机号码:">
                <el-input v-model="registerForm.telephone" placeholder="请输入手机号码"></el-input>
            </el-form-item>
            <el-form-item label="角色:">
                <el-radio-group v-model="registerForm.role">
                    <el-radio :label="2">客户</el-radio>
                    <el-radio :label="3">供应商</el-radio>
                </el-radio-group>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="sendRegisterData">点击注册</el-button>
                <el-button @click="goToLogin">返回</el-button>
            </el-form-item>
        </el-form>
    </div>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { post } from '@/api'
import { AxiosError } from "axios"
import { ElMessage } from "element-plus"

const $router = useRouter()

// 定义表单数据结构
interface Form {
    user_name: string
    password: string
    telephone?: string
    role?: number
}

// 登录表单数据
const loginForm = reactive<Form>({
    user_name: '',
    password: ''
})

// 注册表单数据
const registerForm = reactive<Form>({
    user_name: '',
    password: '',
    telephone: '',
    role: 2 // 默认角色为客户
})

// 控制显示登录还是注册表单
const isLogin = ref(true)

// 切换到注册表单
const goToRegister = () => {
    isLogin.value = false
}

// 切换到登录表单
const goToLogin = () => {
    isLogin.value = true
}

// 登录请求
async function sendLoginData() {
    try {
        const rsp = await post('/user/login', loginForm)
        console.log(rsp)
        ElMessage.success('登录成功')
        // 根据需求进行路由跳转或其他操作
        $router.push('/home')
    } catch (error) {
        const axiosError = error as AxiosError<{ data: { msg: string } }>
        if (axiosError.isAxiosError && axiosError.response) {
            ElMessage.error(`登录失败：${axiosError.response.data.data.msg}`)
        } else {
            ElMessage.error('登录失败，请检查你的信息')
        }
    }
}

// 注册请求
async function sendRegisterData() {
    try {
        const rsp = await post('/user/register', registerForm)
        console.log(rsp)
        ElMessage.success(rsp.data.msg)
        // 注册成功后切换到登录表单
        goToLogin()
    } catch (error) {
        const axiosError = error as AxiosError<{ data: { msg: string } }>
        if (axiosError.isAxiosError && axiosError.response) {
            ElMessage.error(`注册失败：${axiosError.response.data.data.msg}`)
        } else {
            ElMessage.error('注册失败，请检查你的信息')
        }
    }
}
</script>

<style scoped>
.auth-container {
    width: 350px;
    margin: 180px auto;
}

.form-container {
    background-color: #fff;
    border: 1px solid #eaeaea;
    border-radius: 15px;
    padding: 35px 35px 15px 35px;
    box-shadow: 0 0 25px #cacaca;
}

.form-container h3 {
    text-align: center;
    margin-bottom: 20px;
    color: #505450;
}

.form-container .el-form-item__content {
    display: flex;
    justify-content: center;
}
</style>