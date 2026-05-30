import React, { useEffect, useState } from 'react';
import { Card, Avatar, Typography, List, Tag, Spin, Empty, Pagination, Button } from 'antd';
import { UserOutlined, EditOutlined, MessageOutlined, EyeOutlined, ArrowLeftOutlined } from '@ant-design/icons';
import { useParams, useNavigate } from 'react-router-dom';
import * as userService from '@/services/userService';
import { usePostStore } from '@/stores';

const { Title, Text, Paragraph } = Typography;

const UserProfile = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const { posts, pagination, fetchPosts } = usePostStore();

  useEffect(() => {
    const loadUser = async () => {
      try {
        const response = await userService.getUserById(id);
        setUser(response.data);
      } catch (error) {
        console.error('Failed to load user:', error);
      } finally {
        setLoading(false);
      }
    };

    loadUser();
    fetchPosts({ page: 1, author_id: id });
  }, [id, fetchPosts]);

  const formatTime = (timeStr) => {
    if (!timeStr) return '';
    return new Date(timeStr).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: 48 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!user) {
    return <Empty description="用户不存在" />;
  }

  return (
    <div>
      <Button
        type="link"
        icon={<ArrowLeftOutlined />}
        onClick={() => navigate('/')}
        style={{ marginBottom: 16, paddingLeft: 0 }}
      >
        返回首页
      </Button>

      <Card style={{ marginBottom: 24 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
          <Avatar size={64} src={user.avatar} icon={!user.avatar && <UserOutlined />} />
          <div>
            <Title level={3} style={{ marginBottom: 4 }}>
              {user.nickname || user.username}
            </Title>
            <Text type="secondary">
              {user.bio || '这个人很懒，什么都没留下'}
            </Text>
            <br />
            <Text type="secondary" style={{ fontSize: 12 }}>
              注册时间: {formatTime(user.created_at)}
            </Text>
          </div>
        </div>
      </Card>

      <Title level={4}>TA 的帖子</Title>

      {posts.length === 0 ? (
        <Empty description="暂无帖子" />
      ) : (
        <>
          <List
            dataSource={posts}
            renderItem={(post) => (
              <Card
                hoverable
                style={{ marginBottom: 12 }}
                onClick={() => navigate(`/posts/${post.id}`)}
              >
                <Title level={5} style={{ marginBottom: 4 }}>{post.title}</Title>
                <Paragraph ellipsis={{ rows: 2 }} style={{ color: '#666', marginBottom: 8 }}>
                  {post.content}
                </Paragraph>
                <div style={{ display: 'flex', gap: 12, color: '#999', fontSize: 12 }}>
                  <span>{formatTime(post.created_at)}</span>
                  <Tag icon={<MessageOutlined />} style={{ fontSize: 12 }}>{post.comment_count || 0}</Tag>
                  <Tag icon={<EyeOutlined />} style={{ fontSize: 12 }}>{post.view_count || 0}</Tag>
                </div>
              </Card>
            )}
          />
          {pagination.total > pagination.pageSize && (
            <div style={{ textAlign: 'center', marginTop: 16 }}>
              <Pagination
                current={pagination.page}
                pageSize={pagination.pageSize}
                total={pagination.total}
                onChange={(page) => fetchPosts({ page, author_id: id })}
                showTotal={(total) => `共 ${total} 条`}
              />
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default UserProfile;
