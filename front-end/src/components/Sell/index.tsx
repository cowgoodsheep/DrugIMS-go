import React, { useEffect, useState } from 'react'
import {  Space, Tag } from 'antd';
import MyTable from '../MyTable';
import { getSaleList } from '../../api/Api';
import { addOneDay } from '../../utils';

const columns = [
    {
        title: '销售ID',
        dataIndex: 'sale_id',
        key: 'sale_id',
        fixed: 'left',
        width: 100
    },
    {
        title: '药品名称',
        dataIndex: 'drug_name',
        key: 'drug_name',
        fixed: 'left',
        width: 100
    },
    {
        title: '客户名称',
        dataIndex: 'user_name',
        key: 'user_name',
        width: 100
    },
    {
        title: '销售日期',
        dataIndex: 'sale_date',
        key: 'sale_date',
        width: 100
    }, {
        title: '销售数量',
        dataIndex: 'sale_quantity',
        key: 'sale_quantity',
        width: 100
    }, {
        title: '销售单价',
        dataIndex: 'sale_unit_price',
        key: 'sale_unit_price',
        width: 100
    }, 
    {
        title: '销售金额',
        dataIndex: 'sale_amount',
        key: 'sale_amount',
        width: 100
    }, 
];
export default function PublicDb({ searchValue }:{searchValue:string}) {
    const [data, setData] = useState([])
    const [loading, setLoading] = useState(false)
    useEffect(() => {
        if (searchValue === ' ') {
            searchValue = ''
        }
        getData(searchValue)
    }, [searchValue])
    const getData =async (searchValue = '') => {
        setLoading(true)
        const data = await getSaleList(searchValue)
        data.map((v,i)=>{
            data[i].sale_date = addOneDay(data[i].sale_date.split('T')[0])
        })
        setLoading(false)
       setData(data)
    }

    














    return (
        <>
            <MyTable loading={loading} columns={columns} data={data} />
        </>
    )
}
