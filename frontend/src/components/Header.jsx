import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Layout, Button, Dropdown, Avatar, Space, Typography } from 'antd';
import { UserOutlined, LogoutOutlined, SettingOutlined, HomeOutlined } from '@ant-design/icons';
import { useAuthStore } from '@/stores';

const { Header: AntHeader } = Layout;
const { Text } = Typography;

const Header = () => {
  const navigate = useNavigate();
  const { user, isAuthenticated, logout } = useAuthStore();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const userMenuItems = [
    { key: 'settings', icon: <SettingOutlined />, label: '个人设置', onClick: () => navigate('/settings') },
    { key: 'profile', icon: <UserOutlined />, label: '我的主页', onClick: () => user && navigate(`/users/${user.id}`) },
    { type: 'divider' },
    { key: 'logout', icon: <LogoutOutlined />, label: '退出登录', onClick: handleLogout },
  ];

  return (
    <AntHeader style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '0 24px', background: '#fff', boxShadow: '0 2px 8px rgba(0,0,0,0.06)' }}>
      <Link to="/" style={{ fontSize: 18, fontWeight: 'bold', color: '#1890ff', textDecoration: 'none' }}>
        <HomeOutlined style={{ marginRight: 8 }} />留言板
      </Link>
      <Space>
        {isAuthenticated ? (
          <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
            <Space style={{ cursor: 'pointer' }}>
              <Avatar size="small" src={user?.avatar || undefined} icon={<UserOutlined />} />
              <Text>{user?.nickname || user?.username}</Text>
            </Space>
          </Dropdown>
        ) : (
          <>
            <Button type="link" onClick={() => navigate('/login')}>登录</Button>
            <Button type="primary" onClick={() => navigate('/register')}>注册</Button>
          </>
        )}
      </Space>
    </AntHeader>
  );
};

export default Header;
