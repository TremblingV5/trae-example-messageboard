import React, { useEffect, useState } from 'react';
import { Card, List, Pagination, Input, Button, Empty, Spin, Typography, Tag } from 'antd';
import { SearchOutlined, EditOutlined, MessageOutlined, EyeOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { usePostStore, useAuthStore } from '@/stores';

const { Title, Text, Paragraph } = Typography;
const { Search } = Input;

const Home = () => {
  const navigate = useNavigate();
  const { isAuthenticated } = useAuthStore();
  const { posts, pagination, loading, fetchPosts, searchPosts } = usePostStore();
  const [keyword, setKeyword] = useState('');

  useEffect(() => {
    fetchPosts({ page: 1 });
  }, [fetchPosts]);

  const handlePageChange = (page) => {
    if (keyword) {
      searchPosts(keyword, { page });
    } else {
      fetchPosts({ page });
    }
  };

  const handleSearch = (value) => {
    setKeyword(value);
    if (value.trim()) {
      searchPosts(value, { page: 1 });
    } else {
      fetchPosts({ page: 1 });
    }
  };

  const formatTime = (timeStr) => {
    if (!timeStr) return '';
    const date = new Date(timeStr);
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={3} style={{ margin: 0 }}>帖子列表</Title>
        <div style={{ display: 'flex', gap: 12 }}>
          <Search
            placeholder="搜索帖子..."
            allowClear
            onSearch={handleSearch}
            style={{ width: 260 }}
            prefix={<SearchOutlined />}
          />
          {isAuthenticated && (
            <Button type="primary" icon={<EditOutlined />} onClick={() => navigate('/create')}>
              创建帖子
            </Button>
          )}
        </div>
      </div>

      <Spin spinning={loading}>
        {posts.length === 0 ? (
          <Empty description="暂无帖子" />
        ) : (
          <List
            dataSource={posts}
            renderItem={(post) => (
              <Card
                hoverable
                style={{ marginBottom: 16 }}
                onClick={() => navigate(`/posts/${post.id}`)}
              >
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                  <div style={{ flex: 1 }}>
                    <Title level={4} style={{ marginBottom: 8 }}>
                      {post.title}
                    </Title>
                    <Paragraph
                      ellipsis={{ rows: 2 }}
                      style={{ color: '#666', marginBottom: 12 }}
                    >
                      {post.content}
                    </Paragraph>
                    <div style={{ display: 'flex', gap: 16, color: '#999', fontSize: 13 }}>
                      <span>作者: {post.author?.nickname || post.author?.username || '匿名'}</span>
                      <span>{formatTime(post.created_at)}</span>
                    </div>
                  </div>
                  <div style={{ display: 'flex', gap: 16, marginLeft: 16, color: '#999' }}>
                    <Tag icon={<MessageOutlined />}>{post.comment_count || 0} 评论</Tag>
                    <Tag icon={<EyeOutlined />}>{post.view_count || 0} 浏览</Tag>
                  </div>
                </div>
              </Card>
            )}
          />
        )}
      </Spin>

      {pagination.total > 0 && (
        <div style={{ textAlign: 'center', marginTop: 24 }}>
          <Pagination
            current={pagination.page}
            pageSize={pagination.pageSize}
            total={pagination.total}
            onChange={handlePageChange}
            showTotal={(total) => `共 ${total} 条`}
          />
        </div>
      )}
    </div>
  );
};

export default Home;
