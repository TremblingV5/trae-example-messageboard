import React from 'react';
import { Layout, Menu, Button, Avatar, Dropdown, Space } from 'antd';
import {
  HomeOutlined,
  EditOutlined,
  UserOutlined,
  SettingOutlined,
  LoginOutlined,
  UserAddOutlined,
  LogoutOutlined,
} from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuthStore } from '@/stores';

const { Header: AntHeader } = Layout;

const HeaderComponent = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, isAuthenticated, logout } = useAuthStore();

  const handleLogout = () => {
    logout();
    navigate('/login', { replace: true });
  };

  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人主页',
      onClick: () => navigate(`/users/${user?.id}`),
    },
    {
      key: 'settings',
      icon: <SettingOutlined />,
      label: '设置',
      onClick: () => navigate('/settings'),
    },
    {
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      danger: true,
      onClick: handleLogout,
    },
  ];

  const menuItems = [
    {
      key: '/',
      icon: <HomeOutlined />,
      label: '首页',
      onClick: () => navigate('/'),
    },
  ];

  return (
    <AntHeader style={{ display: 'flex', alignItems: 'center', background: '#fff', borderBottom: '1px solid #f0f0f0', padding: '0 24px' }}>
      <div style={{ fontSize: 18, fontWeight: 'bold', marginRight: 40, cursor: 'pointer', color: '#1890ff' }} onClick={() => navigate('/')}>
        留言板社区
      </div>
      <Menu
        mode="horizontal"
        selectedKeys={[location.pathname]}
        items={menuItems}
        style={{ flex: 1, border: 'none' }}
      />
      <Space size="middle">
        {isAuthenticated ? (
          <>
            <Button type="primary" icon={<EditOutlined />} onClick={() => navigate('/create')}>
              发帖
            </Button>
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <Space style={{ cursor: 'pointer' }}>
                <Avatar src={user?.avatar} icon={!user?.avatar && <UserOutlined />} />
                <span>{user?.nickname || user?.username}</span>
              </Space>
            </Dropdown>
          </>
        ) : (
          <>
            <Button icon={<LoginOutlined />} onClick={() => navigate('/login')}>
              登录
            </Button>
            <Button type="primary" icon={<UserAddOutlined />} onClick={() => navigate('/register')}>
              注册
            </Button>
          </>
        )}
      </Space>
    </AntHeader>
  );
};

export default HeaderComponent;
