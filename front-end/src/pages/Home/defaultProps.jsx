import React from 'react';
import logo from './R-C.jpg'
let userinfo = localStorage.getItem('userinfo')
let role
if((userinfo)){
  userinfo = JSON.parse(userinfo)
  role= userinfo.role
  }
const defaultProps = {
    title: 'drugims',
    logo: logo,
    route: {
        path: '/',
        routes: [
            {
                path: '/user',
                name: '用户信息',
                access: 'canAdmin',
                hideInMenu:role!=='1'
            },
            {
                path: '/drug',
                name: '药品信息',
                access: 'canAdmin',
            },
            {
                name: '药品库存',
                path: '/inventory',
                hideInMenu:role!=='1'
            },    {
                name: '销售信息',
                path: '/sellMsg',
                hideInMenu:role!=='1'

            }, 
            {
                name: '入库信息',
                path: '/addMsg',
                hideInMenu:role!=='1'

            }, 
            {
                name: '我的购买记录',
                path: '/myDrug',
                hideInMenu:role!=='2'
            }, 
            {
                name: '我的进货记录',
                path: '/myinput',
                hideInMenu:role!=='3'
            }, 
        ],
    },
    location: {
        pathname: '/',
    },
};
export default defaultProps;