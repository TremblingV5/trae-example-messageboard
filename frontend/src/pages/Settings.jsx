import React, { useState } from 'react';
import { Card, Form, Input, Button, Avatar, Upload, Typography, message } from 'antd';
import { UserOutlined, UploadOutlined } from '@ant-design/icons';
import { useAuthStore } from '@/stores';
import * as userService from '@/services/userService';

const { Title, Text } = Typography;

const Settings = () => {
  const { user, updateUser } = useAuthStore();
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (values) => {
    setLoading(true);
    try {
      await userService.updateUser(user.id, values);
      updateUser(values);
      message.success('更新成功');
    } catch (err) {
      message.error('更新失败');
    } finally {
      setLoading(false);
    }
  };

  const handleAvatarUpload = async (file) => {
    const formData = new FormData();
    formData.append('avatar', file);
    try {
      const res = await userService.uploadAvatar(user.id, formData);
      updateUser({ avatar: res.data?.avatar_url || res.avatar_url });
      message.success('头像更新成功');
    } catch (err) {
      message.error('头像上传失败');
    }
    return false;
  };

  return (
    <div style={{ maxWidth: 600, margin: '0 auto', padding: '24px' }}>
      <Card>
        <Title level={3}>个人设置</Title>
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <Avatar size={100} src={user?.avatar || undefined} icon={<UserOutlined />} />
          <div style={{ marginTop: 12 }}>
            <Upload beforeUpload={handleAvatarUpload} showUploadList={false} accept="image/*">
              <Button icon={<UploadOutlined />}>更换头像</Button>
            </Upload>
          </div>
        </div>
        <Form onFinish={handleSubmit} layout="vertical" initialValues={{ nickname: user?.nickname || '' }}>
          <Form.Item name="nickname" label="昵称">
            <Input placeholder="请输入昵称" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>保存修改</Button>
          </Form.Item>
        </Form>
        <div style={{ marginTop: 24, padding: '12px 0', borderTop: '1px solid #f0f0f0' }}>
          <Text type="secondary">用户名：{user?.username}</Text>
        </div>
      </Card>
    </div>
  );
};

export default Settings;
