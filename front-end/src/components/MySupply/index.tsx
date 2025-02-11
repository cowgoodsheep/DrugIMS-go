import React, { useEffect, useState } from 'react'
import { Space, Tag, Popconfirm, Button } from 'antd';
import MyTable from '../MyTable';
import { deleteUser, getUserSupplyList } from '../../api/Api';
import { useModel } from '../../utils';
import { formatDate } from '../../utils';

export default function PublicDb({ searchValue, change }: { searchValue: string }) {
    const [data, setData] = useState([])
    const [loading, setLoading] = useState(false)
    const { setType } = useModel()
    useEffect(() => {
        if (searchValue === ' ') {
            searchValue = ''
        }
        getData()
    }, [searchValue])
    useEffect(() => {
        getData()
    }, [change])
    const getData = async () => {
        setLoading(false)
        const data = await getUserSupplyList()
        data.map((v, i) => {
            data[i].supply_total_amount = data[i].supply_quantity * (data[i].supply_price || 0);
            data[i].supply_total_amount=data[i].supply_total_amount.toFixed(2)
            data[i].create_time = formatDate(data[i].create_time)
        })
        setLoading(false)
        setData(data)
    }

    const columns = [
        {
            title: '进货ID',
            dataIndex: 'supply_id',
            key: 'supply_id',
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
            dataIndex: 'create_time',
            key: 'create_time',
            fixed: 'left',
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
            width: 100,
        }
    ];
    return (
        <>
            <MyTable loading={loading} columns={columns} data={data} />
        </>
    )
}
