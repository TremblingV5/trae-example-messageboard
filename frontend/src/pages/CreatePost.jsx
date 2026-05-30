import React from 'react';
import { Form, Input, Button, Card, Typography, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { usePostStore } from '@/stores';

const { Title } = Typography;
const { TextArea } = Input;

const CreatePost = () => {
  const navigate = useNavigate();
  const { createPost, loading } = usePostStore();
  const [form] = Form.useForm();

  const handleSubmit = async (values) => {
    try {
      const newPost = await createPost(values);
      message.success('帖子创建成功');
      navigate(`/posts/${newPost.id}`, { replace: true });
    } catch (error) {
      message.error('创建帖子失败');
    }
  };

  return (
    <div style={{ maxWidth: 720, margin: '0 auto' }}>
      <Title level={3}>创建帖子</Title>

      <Card>
        <Form
          form={form}
          onFinish={handleSubmit}
          layout="vertical"
          autoComplete="off"
        >
          <Form.Item
            name="title"
            label="标题"
            rules={[
              { required: true, message: '请输入帖子标题' },
              { max: 100, message: '标题最多100个字符' },
            ]}
          >
            <Input placeholder="请输入帖子标题" size="large" />
          </Form.Item>

          <Form.Item
            name="content"
            label="内容"
            rules={[
              { required: true, message: '请输入帖子内容' },
              { min: 1, message: '内容不能为空' },
            ]}
          >
            <TextArea
              rows={12}
              placeholder="请输入帖子内容..."
              showCount
              maxLength={5000}
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              发布帖子
            </Button>
            <Button style={{ marginLeft: 12 }} onClick={() => navigate(-1)}>
              取消
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default CreatePost;
