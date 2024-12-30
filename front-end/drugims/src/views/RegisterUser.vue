<!-- <template>
    <el-form :model="form" class="login-container">
        <h3>用户注册</h3>
        <el-form-item label="姓名:">
            <el-input v-model="form.user_name" type="input" placeholder="请输入用户名">
            </el-input>
        </el-form-item>
        <el-form-item label="密码:">
            <el-input v-model="form.password" type="input" placeholder="请输入密码">
            </el-input>
        </el-form-item>
        <el-form-item label="手机号码:">
            <el-input v-model="form.telephone" type="input" placeholder="请输入手机号码">
            </el-input>
        </el-form-item>
        <el-form-item label="角色:">
            <el-radio-group v-model="form.role">
                <el-radio :label="1">管理员</el-radio>
                <el-radio :label="2">客户</el-radio>
                <el-radio :label="3">供应商</el-radio>
            </el-radio-group>
        </el-form-item>
        <el-form-item>
            <el-button type="primary" @click="sendData">提交</el-button>
            <el-button @click="goToLogin">返回</el-button>
        </el-form-item>
    </el-form>
</template> -->
<template>
    <div class="login-container" style="background-color: white; height: 100vh;">
        <el-dialog title="新用户注册" v-model="dialogVisible" :close-on-click-modal="false">
            <el-form :model="registerForm" :rules="registerRules" ref="registerFormRef" label-width="100px"
                style="max-width: 600px;">
                <el-form-item label="用户名称" prop="username">
                    <el-input v-model="registerForm.user_name"></el-input>
                </el-form-item>
                <el-form-item label="密码" prop="password">
                    <el-input type="password" v-model="registerForm.password"></el-input>
                </el-form-item>
                <el-form-item label="角色" prop="role">
                    <el-select v-model="registerForm.role" placeholder="请选择">
                        <el-option label="管理员" value=1></el-option>
                        <el-option label="客户" value=2></el-option>
                        <el-option label="供应商" value=3></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="sendData">立即注册</el-button>
                    <el-button @click="goToLogin">返回</el-button>
                </el-form-item>
            </el-form>
        </el-dialog>
        <login-form-page :background-image-url="'https://www.freeimg.cn/i/2023/12/23/65867de19f3bf.jpg'"
            :title="'药品存销信息管理系统'" :sub-title="'3121004661 刘浩洋'" @finish="false">
            <el-form :model="loginForm" :rules="loginRules" ref="loginFormRef">
                <el-form-item label="用户名称" prop="username">
                    <el-input v-model="loginForm.user_name" prefix-icon="el-icon-user"></el-input>
                </el-form-item>
                <el-form-item label="密码" prop="password">
                    <el-input type="password" v-model="loginForm.password" prefix-icon="el-icon-lock"></el-input>
                </el-form-item>
            </el-form>
            <div style="margin-top: 24px;">
                <el-link type="primary" @click="goToLogin">注册用户</el-link>
            </div>
        </login-form-page>
    </div>
</template>


<script lang="ts" setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router';
import { post } from '@/api';
import { AxiosError } from "axios"
import { ElMessage } from "element-plus"

const $router = useRouter();
interface Form {
    user_name: string
    password: string
    telephone: string
    role: number
}
const dialogVisible = false
const loginForm = reactive<Form>({} as Form)
const registerForm = reactive<Form>({} as Form)

const loginRules = {
    username: [{ required: true, message: '请输入用户名!', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码！', trigger: 'blur' }]
};

const registerRules = {
    username: [{ required: true, message: '请输入用户名称', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
    role: [{ required: true, message: '请选择角色', trigger: 'change' }]
};

async function sendData() {
    try {
        const rsp = await post('/user/register', registerForm);
        console.log(rsp);
        ElMessage.success(rsp.data.msg)
    } catch (error) {
        console.log(registerForm);
        // 使用类型断言将 error 断言为 AxiosError 类型
        const axiosError = error as AxiosError<{ data: { msg: string } }>;
        console.log(axiosError);

        // 如果是 AxiosError，并且包含响应信息
        if (axiosError.isAxiosError && axiosError.response) {
            // 显示具体的后端错误信息
            const errorMsg = axiosError.response.data.data.msg;
            ElMessage.error(`注册失败：${errorMsg}`);
        } else {
            // 显示通用的错误信息
            ElMessage.error('注册失败，请检查你的信息');
        }
    }
}

const goToLogin = () => {
    $router.push('/user/login');
}

</script>
<style scoped>
.login-container {
    width: 350px;
    background-color: #fff;
    border: 1px solid#eaeaea;
    border-radius: 15px;
    padding: 35px 35px 15px 35px;
    box-shadow: 0 0 25px #cacaca;
    margin: 180px auto;

    h3 {
        text-align: center;
        margin-bottom: 20px;
        color: #505450;
    }

    /deep/.el-form-item__content {
        justify-content: center;
    }

}
</style>