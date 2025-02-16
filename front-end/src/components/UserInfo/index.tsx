import React, { useEffect, useState } from "react";
import { Card, Button, Row, Col } from 'antd';
import { useModel } from '../../utils';
import { getUser } from "../../api/Api";

const gridStyle: React.CSSProperties = {
  width: '25%',
  textAlign: 'center',
};

const UserInfo: React.FC = () => {
  const { setType } = useModel();
  const [data, setData] = useState(null); // 用户信息
  const [loading, setLoading] = useState(true); // 加载状态

  useEffect(() => {
    getData();
  }, []);

  const getData = async () => {
    setLoading(true);
    try {
      const userInfo = JSON.parse(localStorage.getItem('userinfo'));
      const data = await getUser(userInfo.user_id);
      setData(data);
    } catch (error) {
      console.error("获取用户信息失败：", error);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdate = () => {
    setType(0);
  };

  const handleRecharge = () => {
    // 携带用户信息跳转到充值页面
    history.push('/recharge', { userInfo: data });
  };

  const handleWithdraw = () => {
    // 携带用户信息跳转到提现页面
    history.push('/withdraw', { userInfo: data });
  };

  if (loading) {
    return <div>加载中...</div>;
  }

  if (!data) {
    return <div>用户信息加载失败</div>;
  }

  return (
    <>
      <Card title="个人基本信息">
        <Card.Grid style={gridStyle} className='B'>用户名：</Card.Grid>
        <Card.Grid hoverable={false} style={gridStyle}>{data.user_name}</Card.Grid>
        <Card.Grid className='B' style={gridStyle}>角色：</Card.Grid>
        <Card.Grid hoverable={false} style={gridStyle}>{data.role}</Card.Grid>
        <Card.Grid className='B' style={gridStyle}>电话号：</Card.Grid>
        <Card.Grid hoverable={false} style={gridStyle}>{data.telephone}</Card.Grid>
        <Card.Grid className='B' style={gridStyle}>地址：</Card.Grid>
        <Card.Grid hoverable={false} style={gridStyle}>
          {data.address ? data.address : <span style={{ color: 'gray' }}>请补充地址信息,否则无法发货</span>}
        </Card.Grid>
      </Card>
      <Button type='primary' onClick={handleUpdate} style={{ margin: '10px' }}>修改信息</Button>

      <Card title="账户信息">
        <Card.Grid style={gridStyle} className='B'>账户余额：</Card.Grid>
        <Card.Grid hoverable={false} style={gridStyle}>{data.banlance.toFixed(2)}</Card.Grid>
        <Card.Grid style={gridStyle} className='B'>冻结余额：</Card.Grid>
        <Card.Grid hoverable={false} style={gridStyle}>{data.block_banlance.toFixed(2)}</Card.Grid>
      </Card>
      <Button type='primary' onClick={handleRecharge} style={{ margin: '10px' }}>充值</Button>
      <Button type='primary' onClick={handleWithdraw} style={{ margin: '10px' }}>提现</Button>
    </>
  );
};

export default UserInfo;