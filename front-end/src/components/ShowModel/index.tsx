import React, { useMemo } from 'react'
import { Input, Select, DatePicker, Button, Upload, message, Form, Modal, Spin } from 'antd';
import { useModel } from '../../utils/index';
import FormItem from '../FormItem';
import { updateUserInfo, addDrug, supplyDrug, updateDrug, adminUpdateUser } from '../../api/Api';
export default function MyModel() {
  const { open, setOpen, titleMap, type } = useModel()
  const onFinish = async (values) => {
    switch (type) {
      case 0:
        const data = await updateUserInfo(values)
        const info = JSON.parse(localStorage.getItem('userinfo'))
        Object.keys(values).forEach((key) => {
          if (values[key]) {
            info[key] = values[key]
          }
        })
        let newInfo = { ...info }
        localStorage.setItem('userinfo', JSON.stringify(newInfo))
        message.success('修改成功！')
        location.reload()
        setOpen(false)
        break;
      case 1:
        if (values.stock_lower_limit > values.stock_upper_limit) {
          message.warning('上限数量必须大于下限数量!')
          return
        }
        await addDrug(values)
        message.success('添加成功！')
        location.reload()
        setOpen(false)
        break;
      case 2:
        await supplyDrug(values)
        message.success('进货成功！')
        location.reload()
        setOpen(false)
        break;
      case 3:
        if (values.stock_lower_limit > values.stock_upper_limit) {
          message.warning('上限数量必须大于下限数量!')
          return
        }
        values.note = values.note || '无'
        await updateDrug(values)
        message.success('修改成功！')
        location.reload()
        setOpen(false)
        break;
      case 4:
        await adminUpdateUser(values)
        message.success('修改成功！')
        location.reload()
        setOpen(false)
        break;
      default:
        break;
    }
  };
  const Item = useMemo(() => {
    switch (type) {
      case 0:
        return (<>
          <FormItem label='用户名' name='user_name'></FormItem>
          <FormItem label='密码' name='password'></FormItem>
          <FormItem label='电话号' name='telephone'></FormItem>
          <FormItem label='地址' name='address'></FormItem>
        </>)
      case 1:
        return (<>
          <FormItem label='药品图片(链接)' name='img' need={false} type='string'></FormItem>
          <FormItem label='药品名称' name='drug_name' need={true} type='string'></FormItem>
          <FormItem label='药品说明' name='drug_description' need={true} type='string'></FormItem>
          <FormItem label='生产厂家' name='manufacturer' need={true} type='string'></FormItem>
          <FormItem label='单位' name='unit' need={true} type='string'></FormItem>
          <FormItem label='规格' name='specification' need={true} type='string'></FormItem>
          <FormItem label='库存下限' name='stock_lower_limit' need={true} type='number'></FormItem>
          <FormItem label='库存上限' name='stock_upper_limit' need={true} type='number'></FormItem>
          <FormItem label='售价' name='sale_price' need={true} type='number'></FormItem>
        </>)
      case 2:
        return (<>
          <FormItem label='批号' name='batch_number' need={true} type='string'></FormItem>
          <Form.Item label="生产日期" name='production_date'
            rules={[{
              required: true,
            }]}
          >
            <DatePicker />
          </Form.Item>
          <FormItem label='进货单价' name='supply_price' need={true} type='number'></FormItem>
          <FormItem label='进货数量' name='supply_quantity' need={true} type='number'></FormItem>
          <FormItem label='备注' name='note' type='string'></FormItem>
        </>)
      case 3:
        return (<>
          <FormItem label='药品图片(链接)' name='img' type='string'></FormItem>
          <FormItem label='药品名称' name='drug_name' type='string'></FormItem>
          <FormItem label='药品说明' name='drug_description' type='string'></FormItem>
          <FormItem label='生产厂家' name='manufacturer' type='string'></FormItem>
          <FormItem label='单位' name='unit' type='string'></FormItem>
          <FormItem label='规格' name='specification' type='string'></FormItem>
          <FormItem label='库存下限' name='stock_lower_limit' type='number'></FormItem>
          <FormItem label='库存上限' name='stock_upper_limit' type='number'></FormItem>
          <FormItem label='售价' name='sale_price' type='number'></FormItem>
        </>)
      case 4:
        return (<>
          <FormItem label='用户名' name='user_name' type='string'></FormItem>
          <FormItem label='密码' name='password' type='string'></FormItem>
          <Form.Item label="角色" name='role'>
            <Select>
              <Select.Option value="客户">客户</Select.Option>
              <Select.Option value="供应商">供应商</Select.Option>
              <Select.Option value="管理员">管理员</Select.Option>
            </Select>
          </Form.Item>
          <FormItem label='电话号' name='telephone' type='string'></FormItem>
          <FormItem label='地址' name='address' type='string'></FormItem>
        </>)
      default:
        break;
    }
  }, [type])
  return (
    <div>    <Modal
      title={titleMap[type]}
      open={open}
      onCancel={() => setOpen(false)}
      footer={null}
    >
      <Form
        name="basic"
        labelCol={{
          span: 5,
        }}
        wrapperCol={{
          span: 16,
        }}
        style={{
          maxWidth: 600,
        }}
        initialValues={{
          remember: true,
        }}
        onFinish={onFinish}
        autoComplete="on"
      >
        {Item}
        <Form.Item
          wrapperCol={{
            offset: 10,
            span: 2,
          }}
        >
          <Button type="primary" htmlType="submit">
            确认
          </Button>
        </Form.Item>
      </Form>
    </Modal></div>
  )
}
