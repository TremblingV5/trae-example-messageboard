import React from 'react';
import { createBrowserRouter, Navigate, Outlet, useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/stores';
import { Layout } from 'antd';
import { Header } from '@/components';
import {
  Home,
  Login,
  Register,
  PostDetail,
  CreatePost,
  UserProfile,
  Settings,
} from '@/pages';

const { Content } = Layout;

// 处理 401 未授权事件，使用 React Router 导航
function AuthGuard() {
  const navigate = useNavigate();
  const logout = useAuthStore((state) => state.logout);

  React.useEffect(() => {
    const handleUnauthorized = () => {
      logout();
      navigate('/login', { replace: true });
    };

    window.addEventListener('auth:unauthorized', handleUnauthorized);
    return () => {
      window.removeEventListener('auth:unauthorized', handleUnauthorized);
    };
  }, [navigate, logout]);

  return <Outlet />;
}

// 路由守卫组件 - 需要登录才能访问
const ProtectedRoute = ({ children }) => {
  const { isAuthenticated } = useAuthStore();
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }
  
  return children;
};

// 主布局组件
const MainLayout = ({ children }) => {
  return (
    <Layout className="main-layout">
      <Header />
      <Content className="main-content">
        {children}
      </Content>
    </Layout>
  );
};

// 认证布局组件
const AuthLayout = ({ children }) => {
  return (
    <Layout className="auth-layout">
      <Content className="auth-content">
        {children}
      </Content>
    </Layout>
  );
};

// 路由配置
const router = createBrowserRouter([
  {
    element: <AuthGuard />,
    children: [
      {
        path: '/',
        element: <MainLayout><Home /></MainLayout>,
      },
      {
        path: '/login',
        element: <AuthLayout><Login /></AuthLayout>,
      },
      {
        path: '/register',
        element: <AuthLayout><Register /></AuthLayout>,
      },
      {
        path: '/posts/:id',
        element: <MainLayout><PostDetail /></MainLayout>,
      },
      {
        path: '/create',
        element: <MainLayout><ProtectedRoute><CreatePost /></ProtectedRoute></MainLayout>,
      },
      {
        path: '/users/:id',
        element: <MainLayout><UserProfile /></MainLayout>,
      },
      {
        path: '/settings',
        element: <MainLayout><ProtectedRoute><Settings /></ProtectedRoute></MainLayout>,
      },
      {
        path: '*',
        element: <Navigate to="/" replace />,
      },
    ],
  },
]);

export default router;
