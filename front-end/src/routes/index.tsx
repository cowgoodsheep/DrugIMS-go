import Login from '../pages/Login/index.tsx'
import Home from '../pages/Home/index.tsx'
import Pay from '../pages/Pay/index.tsx'
const routes = [
    // 登录注册
    {
        path: '/login',
        element: <Login />,
    },
    // home
    {
        path: '/home',
        element: <Home />,
    },
    // 支付页
    {
        path: '/pay',
        element: <Pay />,
    },
    // default route
    {
        path: '',
        element: <Login />
    }
]
export default routes 