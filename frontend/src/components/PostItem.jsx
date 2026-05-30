import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, Avatar, Typography, Tag } from 'antd';
import { UserOutlined } from '@ant-design/icons';

const { Text, Paragraph } = Typography;

const PostItem = ({ post }) => {
  const navigate = useNavigate();

  return (
    <Card hoverable onClick={() => navigate(`/posts/${post.id}`)} style={{ marginBottom: 16 }}>
      <Typography.Title level={4} style={{ marginBottom: 8 }}>{post.title}</Typography.Title>
      <Paragraph ellipsis={{ rows: 2 }} style={{ color: '#666', marginBottom: 12 }}>
        {post.content}
      </Paragraph>
      <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
        <Avatar size="small" src={post.author?.avatar || undefined} icon={<UserOutlined />} />
        <Text type="secondary">{post.author?.nickname || post.author?.username}</Text>
        <Tag color="blue">{new Date(post.created_at).toLocaleDateString()}</Tag>
      </div>
    </Card>
  );
};

export default PostItem;
