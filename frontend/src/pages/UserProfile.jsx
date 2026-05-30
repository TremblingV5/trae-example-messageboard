import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card, Avatar, Typography, Spin, List, Tag } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import * as userService from '@/services/userService';
import * as postService from '@/services/postService';

const { Title, Text, Paragraph } = Typography;

const UserProfile = () => {
  const { id } = useParams();
  const [user, setUser] = useState(null);
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const userData = await userService.getUserById(id);
        setUser(userData);
        const postsData = await postService.getPosts({ page: 1, pageSize: 20 });
        setPosts(postsData.posts || []);
      } catch (err) {
        console.error('Failed to fetch user data:', err);
      } finally {
        setLoading(false);
      }
    };
    if (id) fetchData();
  }, [id]);

  if (loading) return <Spin style={{ display: 'block', margin: '100px auto' }} />;
  if (!user) return <div style={{ textAlign: 'center', padding: 100 }}>用户不存在</div>;

  return (
    <div style={{ maxWidth: 800, margin: '0 auto', padding: '24px' }}>
      <Card>
        <div style={{ display: 'flex', alignItems: 'center', gap: 16, marginBottom: 24 }}>
          <Avatar size={80} src={user.avatar || undefined} icon={<UserOutlined />} />
          <div>
            <Title level={3} style={{ margin: 0 }}>{user.nickname || user.username}</Title>
            <Text type="secondary">@{user.username} | 注册于 {new Date(user.created_at).toLocaleDateString()}</Text>
          </div>
        </div>
      </Card>
      <Title level={4} style={{ marginTop: 24 }}>发布的帖子</Title>
      <List
        dataSource={posts.filter(p => p.author && String(p.author.id) === String(id))}
        locale={{ emptyText: '暂无帖子' }}
        renderItem={(post) => (
          <List.Item>
            <List.Item.Meta
              title={<a href={`/posts/${post.id}`}>{post.title}</a>}
              description={<Paragraph ellipsis={{ rows: 2 }}>{post.content}</Paragraph>}
            />
            <Tag>{new Date(post.created_at).toLocaleDateString()}</Tag>
          </List.Item>
        )}
      />
    </div>
  );
};

export default UserProfile;
