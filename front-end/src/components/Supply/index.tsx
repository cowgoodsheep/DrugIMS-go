import React, { useEffect, useState } from 'react'
import { Space, Tag } from 'antd';
import MyTable from '../MyTable';
import { getSupplyList } from '../../api/Api';
import { addOneDay } from '../../utils';

const columns = [
    {
        title: '进货单ID',
        dataIndex: 'supply_id',
        key: 'supply_id',
        fixed: 'left',
        width: 100
    },
    {
        title: '供应商ID',
        dataIndex: 'user_id',
        key: 'user_id',
        fixed: 'left',
        width: 100
    },
    {
        title: '供应商名称',
        dataIndex: 'user_name',
        key: 'user_name',
        fixed: 'left',
        width: 100
    },
    {
        title: '进货药品名称',
        dataIndex: 'drug_name',
        key: 'drug_name',
        fixed: 'left',
        width: 100
    },
    {
        title: '进货日期',
        dataIndex: 'supply_date',
        key: 'supply_date',
        width: 100
    },
    {
        title: '批号',
        dataIndex: 'batch_number',
        key: 'batch_number',
        fixed: 'left',
        width: 100
    },
    {
        title: '进货数量',
        dataIndex: 'supply_quantity',
        key: 'supply_quantity',
        width: 100
    },
    {
        title: '进货单价',
        dataIndex: 'supply_price',
        key: 'supply_price',
        width: 100
    },
    {
        title: '进货总金额',
        dataIndex: 'supply_total_amount',
        key: 'supply_total_amount',
        width: 100
    },
    {
        title: '备注',
        dataIndex: 'note',
        key: 'note',
        width: 100
    },
];
export default function PublicDb({ searchValue }: { searchValue: string }) {
    const [data, setData] = useState([])
    const [loading, setLoading] = useState(false)
    useEffect(() => {
        if (searchValue === ' ') {
            searchValue = ''
        }
        getData(searchValue)
    }, [searchValue])
    const getData = async (searchValue = '') => {
        setLoading(false)
        const data = await getSupplyList(searchValue)
        data.map((v, i) => {
            data[i].supply_total_amount = data[i].supply_quantity * (data[i].supply_price || 0);
            data[i].supply_total_amount=data[i].supply_total_amount.toFixed(2)
            data[i].supply_date = addOneDay(data[i].create_time.split('T')[0])
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
