import { LogoutOutlined, UserOutlined, } from '@ant-design/icons';
import type { ProSettings } from '@ant-design/pro-components';
import { ProConfigProvider, ProLayout, SettingDrawer, } from '@ant-design/pro-components';
import { Dropdown, Button, Tooltip, Input, message } from 'antd';
import React, { useState, useEffect, useRef } from 'react';
import defaultProps from './defaultProps.jsx';
import MyPageContainer from '../../components/MyPageContainer';
import { aiChat } from "../../api/Api";
import { useNavigate } from 'react-router-dom';
import marked from 'marked';
import './index.scss';

const Home = () => {
  // 菜单模式设置
  const [settings, setSetting] = useState<Partial<ProSettings> | undefined>({
    fixSiderbar: true,
    layout: 'mix',
  });

  const navigate = useNavigate();
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


  const handleChatOpen = () => { // 打开聊天框
    setIsChatOpen(true);
  };

  const handleChatClose = () => { // 关闭聊天
    setIsChatOpen(false);
  };

  const handleUserInput = (e) => { // 用户输入
    setUserInput(e.target.value);
  };

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
            render: (props, dom) => {
              return (
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