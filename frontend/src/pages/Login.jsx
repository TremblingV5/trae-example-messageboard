import React from 'react';
import { Form, Input, Button, Card, Typography, message } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '@/stores';

const { Title, Text } = Typography;

const Login = () => {
  const navigate = useNavigate();
  const { login, loading, error, clearError } = useAuthStore();
  const [form] = Form.useForm();

  const handleSubmit = async (values) => {
    clearError();
    const result = await login(values);
    if (result.success) {
      message.success('登录成功');
      navigate('/', { replace: true });
    } else {
      message.error(result.error);
    }
  };

  return (
    <Card style={{ boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
      <div style={{ textAlign: 'center', marginBottom: 24 }}>
        <Title level={2}>登录</Title>
        <Text type="secondary">欢迎回到留言板社区</Text>
      </div>

      <Form
        form={form}
        onFinish={handleSubmit}
        size="large"
        autoComplete="off"
      >
        <Form.Item
          name="username"
          rules={[{ required: true, message: '请输入用户名' }]}
        >
          <Input prefix={<UserOutlined />} placeholder="用户名" />
        </Form.Item>

        <Form.Item
          name="password"
          rules={[{ required: true, message: '请输入密码' }]}
        >
          <Input.Password prefix={<LockOutlined />} placeholder="密码" />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={loading} block>
            登录
          </Button>
        </Form.Item>

        <div style={{ textAlign: 'center' }}>
          <Text type="secondary">
            还没有账号？ <Link to="/register">立即注册</Link>
          </Text>
        </div>
      </Form>
    </Card>
  );
};

export default Login;
