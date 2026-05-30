import React from 'react';
import { Pagination, Empty } from 'antd';
import PostItem from './PostItem';

const PostList = ({ posts, pagination, onPageChange }) => {
  if (!posts || posts.length === 0) {
    return <Empty description="暂无帖子" />;
  }

  return (
    <div>
      {posts.map((post) => (
        <PostItem key={post.id} post={post} />
      ))}
      {pagination && pagination.totalPages > 1 && (
        <div style={{ display: 'flex', justifyContent: 'center', marginTop: 24 }}>
          <Pagination
            current={pagination.page}
            pageSize={pagination.pageSize}
            total={pagination.total}
            onChange={onPageChange}
            showTotal={(total) => `共 ${total} 条`}
          />
        </div>
      )}
    </div>
  );
};

export default PostList;
