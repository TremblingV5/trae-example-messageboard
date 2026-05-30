import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Input, Button, Spin, Typography } from 'antd';
import { SearchOutlined, PlusOutlined } from '@ant-design/icons';
import { usePostStore, useAuthStore } from '@/stores';
import PostList from '@/components/PostList';

const { Title } = Typography;

const Home = () => {
  const navigate = useNavigate();
  const { posts, pagination, loading, fetchPosts, searchPosts } = usePostStore();
  const { isAuthenticated } = useAuthStore();
  const [keyword, setKeyword] = useState('');

  useEffect(() => {
    fetchPosts({ page: 1 });
  }, []);

  const handleSearch = () => {
    if (keyword.trim()) {
      searchPosts(keyword, { page: 1 });
    } else {
      fetchPosts({ page: 1 });
    }
  };

  const handlePageChange = (page) => {
    if (keyword.trim()) {
      searchPosts(keyword, { page });
    } else {
      fetchPosts({ page });
    }
  };

  return (
    <div style={{ maxWidth: 800, margin: '0 auto', padding: '24px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={3} style={{ margin: 0 }}>留言板</Title>
        {isAuthenticated && (
          <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/create')}>
            发帖
          </Button>
        )}
      </div>
      <div style={{ display: 'flex', gap: 8, marginBottom: 24 }}>
        <Input
          placeholder="搜索帖子..."
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onPressEnter={handleSearch}
          allowClear
        />
        <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
          搜索
        </Button>
      </div>
      <Spin spinning={loading}>
        <PostList
          posts={posts}
          pagination={pagination}
          onPageChange={handlePageChange}
        />
      </Spin>
    </div>
  );
};

export default Home;
