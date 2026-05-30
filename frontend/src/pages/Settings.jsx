import React, { useState, useRef } from 'react';
import { Card, Form, Input, Button, Avatar, Upload, Typography, message, Divider } from 'antd';
import { UserOutlined, UploadOutlined } from '@ant-design/icons';
import { useAuthStore } from '@/stores';
import * as userService from '@/services/userService';

const { Title, Text } = Typography;

const Settings = () => {
  const { user, updateUser } = useAuthStore();
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [avatarUploading, setAvatarUploading] = useState(false);
  const fileInputRef = useRef(null);

  React.useEffect(() => {
    if (user) {
      form.setFieldsValue({
        nickname: user.nickname || '',
        bio: user.bio || '',
      });
    }
  }, [user, form]);

  const handleUpdateProfile = async (values) => {
    if (!user?.id) return;
    setLoading(true);
    try {
      const response = await userService.updateUser(user.id, values);
      const updatedUser = response.data;
      updateUser(updatedUser);
      message.success('个人信息更新成功');
    } catch (error) {
      message.error('更新失败');
    } finally {
      setLoading(false);
    }
  };

  const handleAvatarUpload = async (file) => {
    if (!user?.id) return;

    const isImage = file.type.startsWith('image/');
    if (!isImage) {
      message.error('只能上传图片文件');
      return false;
    }

    const isLt2M = file.size / 1024 / 1024 < 2;
    if (!isLt2M) {
      message.error('图片大小不能超过 2MB');
      return false;
    }

    setAvatarUploading(true);
    try {
      const formData = new FormData();
      formData.append('avatar', file);
      const response = await userService.uploadAvatar(user.id, formData);
      const updatedUser = response.data;
      updateUser(updatedUser);
      message.success('头像更新成功');
    } catch (error) {
      message.error('头像上传失败');
    } finally {
      setAvatarUploading(false);
    }

    return false; // 阻止默认上传行为
  };

  return (
    <div style={{ maxWidth: 600, margin: '0 auto' }}>
      <Title level={3}>个人设置</Title>

      <Card style={{ marginBottom: 24 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 16, marginBottom: 24 }}>
          <Avatar size={80} src={user?.avatar} icon={!user?.avatar && <UserOutlined />} />
          <div>
            <Upload
              beforeUpload={handleAvatarUpload}
              showUploadList={false}
              accept="image/*"
            >
              <Button icon={<UploadOutlined />} loading={avatarUploading}>
                上传头像
              </Button>
            </Upload>
            <br />
            <Text type="secondary" style={{ fontSize: 12, marginTop: 4 }}>
              支持 JPG、PNG 格式，不超过 2MB
            </Text>
          </div>
        </div>

        <Divider />

        <Form
          form={form}
          onFinish={handleUpdateProfile}
          layout="vertical"
          autoComplete="off"
        >
          <Form.Item
            name="nickname"
            label="昵称"
            rules={[
              { max: 30, message: '昵称最多30个字符' },
            ]}
          >
            <Input placeholder="设置你的昵称" />
          </Form.Item>

          <Form.Item
            name="bio"
            label="个人简介"
            rules={[
              { max: 200, message: '个人简介最多200个字符' },
            ]}
          >
            <Input.TextArea rows={3} placeholder="介绍一下自己吧" showCount maxLength={200} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              保存修改
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default Settings;
