import { LogoutOutlined, UserOutlined, } from '@ant-design/icons';
import type { ProSettings } from '@ant-design/pro-components';
import { ProConfigProvider, ProLayout, SettingDrawer, } from '@ant-design/pro-components';
import { Dropdown, Button, Tooltip, Input, message, Space, Badge, Modal } from 'antd';
import { WarningOutlined } from '@ant-design/icons';
import React, { useState, useEffect, useRef } from 'react';
import defaultProps from './defaultProps.jsx';
import MyPageContainer from '../../components/MyPageContainer';
import { aiChat, riskManage } from "../../api/Api";
import { useNavigate } from 'react-router-dom';
import './index.scss';

const Home = () => {
  // 菜单模式设置
  const [settings, setSetting] = useState<Partial<ProSettings> | undefined>({
    fixSiderbar: true,
    layout: 'mix',
  });

  const navigate = useNavigate();
  const role = JSON.parse(localStorage.getItem("userinfo")).role;
  const path = sessionStorage.getItem('pathname') || '/userMsg'
  const [pathname, setPathname] = useState(path);// 历史记录
  const user_name = JSON.parse(localStorage.getItem('userinfo')).user_name
  const [chatHistory, setChatHistory] = useState([]); // 聊天记录
  const [userInput, setUserInput] = useState(''); // 用户输入
  const [isChatOpen, setIsChatOpen] = useState(false); // 控制聊天框显示
  const [shouldSendRequest, setShouldSendRequest] = useState(false); // 是否需要发送请求
  const chatBoxRef = useRef(null); // 引用聊天框的 DOM 元素

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (chatBoxRef.current && !chatBoxRef.current.contains(event.target)) {
        setIsChatOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isChatOpen]);

  useEffect(() => {
    sessionStorage.setItem('pathname', pathname)
  }, [pathname])

  useEffect(() => {
    if (shouldSendRequest) {
      // 发送请求
      (async () => {
        try {
          const data = await aiChat(chatHistory);
          if (data.data.msg) {
            // 添加 AI 消息到聊天记录
            setChatHistory((prev) => [...prev, { role: 'assistant', content: data.data.msg }]);
          }
        } catch (error) {
          message.error('获取 AI 回复失败, Error: ' + error);
        }
        setShouldSendRequest(false); // 重置标志
      })();
    }
  }, [chatHistory, shouldSendRequest]);

  const handleSend = async () => {
    if (!userInput.trim()) return;
    // 添加用户消息到聊天记录
    setChatHistory((prev) => [...prev, { role: 'user', content: userInput }]);
    setShouldSendRequest(true); // 设置标志，表示需要发送请求
    setUserInput(''); // 清空输入框
  };

  const handleClearHistory = () => { // 清除历史记录
    setChatHistory([]);
  };

  const [riskWarnings, setRiskWarnings] = useState([]); // 风控警告信息
  const [hasUnprocessedWarning, setHasUnprocessedWarning] = useState(0); // 是否有未处理的警告
  const [isRiskModalOpen, setIsRiskModalOpen] = useState(false); // 控制风控警告弹窗显示
  const fetchRiskManage = async () => {
    const data = await riskManage();
    if (data && data.length > 0) {
      setRiskWarnings(data);
      setHasUnprocessedWarning(data.length);
    } else {
      setRiskWarnings([]);
      setHasUnprocessedWarning(false);
    }
  };
  useEffect(() => {
    if (role === "管理员") { // 管理员登录才加载
      fetchRiskManage(); // 初始加载
      const intervalId = setInterval(fetchRiskManage, 60000); // 每分钟刷新一次
      return () => {
        clearInterval(intervalId); // 清除定时器
      };
    }
  }, []);

  return (
    <div id="test-pro-layout" style={{ height: '100vh', }}>
      <ProConfigProvider hashed={false}>
        <ProLayout
          {...settings}
          prefixCls="my-prefix"
          {...defaultProps}
          location={{
            pathname,
          }}
          siderMenuType="group"
          menu={{
            collapsedShowGroupTitle: true,
          }}
          avatarProps={{
            src: 'https://gw.alipayobjects.com/zos/antfincdn/efFD%24IOql2/weixintupian_20170331104822.jpg',
            size: 'small',
            title: user_name,
            render: (_, dom) => {
              return (
                <Space>
                  {role === '管理员' && hasUnprocessedWarning > 0 && <Badge count={hasUnprocessedWarning} style={{ backgroundColor: 'red' }} />}
                  {role === '管理员' && (
                    <Button
                      onClick={() => setIsRiskModalOpen(true)}
                      type="text"
                      icon={<WarningOutlined />}
                    />
                  )}
                  <Dropdown
                    menu={{
                      items: [
                        {
                          key: 'usermsg',
                          icon: <UserOutlined />,
                          label: <div onClick={() => { setPathname('/userMsg') }}>个人信息</div>,
                        },
                        {
                          key: 'logout',
                          icon: <LogoutOutlined />,
                          label: <div onClick={() => { localStorage.clear(); navigate('/login'); }}>退出登录</div>,
                        },
                      ],
                    }}
                  >
                    {dom}
                  </Dropdown>
                </Space>
              );
            },
          }}
          menuFooterRender={(props) => {
            if (props?.collapsed) return undefined;
            return (
              <div
                style={{
                  textAlign: 'center',
                  paddingBlockStart: 12,
                }}
              >
                <div>© 2025 CowGoodSheep</div>
              </div>
            );
          }}
          onMenuHeaderClick={(e) => { }}
          menuItemRender={(item, dom) =>
          (<div
            onClick={() => {
              setPathname(item.path || '/userMsg');
            }}
          >
            {dom}
          </div>)
          }
        >
          <MyPageContainer pathname={pathname} />

          <div className="chat-button-container">
            <Tooltip title={'有问题可以咨询我哦'} placement="left">
              <Button
                onClick={() => setIsChatOpen(true)}
                className="chat-button"
              />
            </Tooltip>
          </div>

          {isChatOpen && (
            <div className="chat-box" ref={chatBoxRef}>
              <div className="chat-header">
                <Button onClick={() => setIsChatOpen(false)}>关闭</Button>
              </div>
              <div className="chat-content">
                {chatHistory.map((message, index) => (
                  <div key={index} className={`message ${message.role}`} data-user_name={user_name}>
                    {message.content}
                  </div>
                ))}
              </div>
              <div className="chat-footer">
                <Input
                  value={userInput}
                  onChange={(e) => setUserInput(e.target.value)}
                  onPressEnter={handleSend}
                  placeholder="请输入您的问题..."
                />
                <Button onClick={handleSend}>发送</Button>
                <Button onClick={handleClearHistory}>清除历史</Button>
              </div>
            </div>
          )}

          <Modal
            title="风控警告"
            visible={isRiskModalOpen}
            onCancel={() => setIsRiskModalOpen(false)}
            footer={null}
            width={600}
          >
            <div>
              {riskWarnings.map((warning, index) => (
                <div key={index} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '10px 0', borderBottom: '1px solid #eaeaea' }}>
                  <div>
                    <strong>{warning.warning_title}</strong>
                    <p>{warning.warning_description}</p>
                  </div>
                </div>
              ))}
              <div style={{ textAlign: 'right', padding: '10px 0' }}>
                <Button type="primary" onClick={fetchRiskManage}>
                  刷新警告列表
                </Button>
              </div>
            </div>
          </Modal>

          <SettingDrawer
            pathname={pathname}
            enableDarkTheme
            getContainer={() => document.getElementById('test-pro-layout')}
            settings={settings}
            onSettingChange={(changeSetting) => {
              setSetting(changeSetting);
            }}
            disableUrlParams={false}
          />
        </ProLayout>
      </ProConfigProvider>
    </div>
  );
};
export default Home