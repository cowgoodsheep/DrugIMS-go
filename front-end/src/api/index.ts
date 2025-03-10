
import axios from "axios";
import { message } from 'antd'
// 创建 axios 请求实例
const serviceAxios = axios.create({
  baseURL: "http://localhost:8080", // 基础请求地址
  // baseURL: "http://1pg04db718281.vicp.fun:58602/",
  timeout: 10000, // 请求超时设置
  withCredentials: false, // 跨域请求是否需要携带 cookie
});

// 创建请求拦截
serviceAxios.interceptors.request.use(
  (config) => {
    // 如果开启 token 认证
    // if (false) {
    config.headers["Token"] = localStorage.getItem("token"); // 请求头携带 Token
    // }
    // 设置请求头
    if (!config.headers["content-type"]) { // 如果没有设置请求头

      config.headers["content-type"] = "application/json"; // 默认类型
    }
    return config;
  },
  (error) => {
    Promise.reject(error);
  }
);

// 创建响应拦截
serviceAxios.interceptors.response.use(
  (res) => {
    const { data, headers } = res;  // 获取响应数据和响应头
    const token = headers['token'] || headers['x-auth-token'];  // 根据服务器返回的 header 名称获取 token
    if (token) localStorage.setItem('token', token);
    // 处理自己的业务逻辑，比如判断 token 是否过期等等
    // 代码块
    return data;
  },
  (error) => {
    if (error.response.data.data.msg.includes('token')) window.location.href = 'http://localhost:5173/'
    message.warning(error.response.data.data.msg);
    return Promise.reject(error.response.data.data.msg);
  }
);

export default serviceAxios;