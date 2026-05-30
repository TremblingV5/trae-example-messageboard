import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Form, Input, Button, Card, Typography, Upload, message } from 'antd';
import { InboxOutlined } from '@ant-design/icons';
import { usePostStore } from '@/stores';
import * as postService from '@/services/postService';

const { Title } = Typography;
const { TextArea } = Input;
const { Dragger } = Upload;

const CreatePost = () => {
  const navigate = useNavigate();
  const { createPost } = usePostStore();
  const [loading, setLoading] = useState(false);
  const [imageUrl, setImageUrl] = useState('');

  const handleSubmit = async (values) => {
    setLoading(true);
    try {
      await createPost({ ...values, image_url: imageUrl });
      message.success('发帖成功');
      navigate('/');
    } catch (err) {
      message.error('发帖失败');
    } finally {
      setLoading(false);
    }
  };

  const handleUpload = async (file) => {
    const formData = new FormData();
    formData.append('image', file);
    try {
      const res = await postService.uploadPostImage(formData);
      setImageUrl(res.image_url);
      message.success('图片上传成功');
    } catch (err) {
      message.error('图片上传失败');
    }
    return false;
  };

  return (
    <div style={{ maxWidth: 800, margin: '0 auto', padding: '24px' }}>
      <Card>
        <Title level={3}>创建帖子</Title>
        <Form onFinish={handleSubmit} layout="vertical" size="large">
          <Form.Item name="title" label="标题" rules={[{ required: true, message: '请输入标题' }]}>
            <Input placeholder="请输入帖子标题" />
          </Form.Item>
          <Form.Item name="content" label="内容" rules={[{ required: true, message: '请输入内容' }]}>
            <TextArea rows={8} placeholder="请输入帖子内容" />
          </Form.Item>
          <Form.Item label="图片（可选）">
            <Dragger accept="image/*" maxCount={1} beforeUpload={handleUpload} onRemove={() => setImageUrl('')}>
              <p className="ant-upload-drag-icon"><InboxOutlined /></p>
              <p className="ant-upload-text">点击或拖拽上传图片</p>
            </Dragger>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>发布帖子</Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default CreatePost;
