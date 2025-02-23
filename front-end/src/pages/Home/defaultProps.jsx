import React from 'react';
import logo from './R-C.jpg'
let userinfo = localStorage.getItem('userinfo')
let role
if ((userinfo)) {
    userinfo = JSON.parse(userinfo)
    role = userinfo.role
}
const defaultProps = {
    title: 'DrugIMS',
    logo: logo,
    route: {
        path: '/',
        routes: [
            // 总
            {
                path: '/drug',
                name: '药品信息',
                access: 'canAdmin',
            },
            // 管理员
            {
                path: '/user',
                name: '用户信息',
                access: 'canAdmin',
                hideInMenu: role !== '管理员'
            },
            {
                name: '药品库存',
                path: '/stock',
                hideInMenu: role !== '管理员'
            },
            {
                name: '销售信息',
                path: '/saleInfo',
                hideInMenu: role !== '管理员'
            },
            {
                name: '入库信息',
                path: '/supplyInfo',
                hideInMenu: role !== '管理员'
            },
            {
                name: '审批信息',
                path: '/approval',
                hideInMenu: role !== '管理员'
            },
            {
                name: '统计信息',
                path: '/statisticsInfo',
                hideInMenu: role !== '管理员'
            },
            // 客户
            {
                name: '我的订单',
                path: '/myBuyRecord',
                hideInMenu: role !== '客户'
            },
            // 供应商
            {
                name: '我的进货',
                path: '/myinput',
                hideInMenu: role !== '供应商'
            },
        ],
    },
    location: {
        pathname: '/',
    },
};
export default defaultProps;