import React from 'react';
import { Avatar, Typography, Button, Space, message } from 'antd';
import { UserOutlined, LikeOutlined, MessageOutlined } from '@ant-design/icons';
import { useCommentStore, useAuthStore } from '@/stores';
import CommentList from './CommentList';

const { Text, Paragraph } = Typography;

const CommentItem = ({ comment }) => {
  const { voteComment } = useCommentStore();
  const { isAuthenticated } = useAuthStore();

  const handleVote = async (action) => {
    if (!isAuthenticated) {
      message.warning('请先登录');
      return;
    }
    try {
      await voteComment(comment.id, action);
    } catch (err) {
      message.error('操作失败');
    }
  };

  return (
    <div style={{ padding: '12px 0', borderBottom: '1px solid #f0f0f0' }}>
      <div style={{ display: 'flex', gap: 8 }}>
        <Avatar size="small" src={comment.author?.avatar || undefined} icon={<UserOutlined />} />
        <div style={{ flex: 1 }}>
          <Space size={4}>
            <Text strong>{comment.author?.nickname || comment.author?.username}</Text>
            <Text type="secondary" style={{ fontSize: 12 }}>{new Date(comment.created_at).toLocaleString()}</Text>
          </Space>
          <Paragraph style={{ margin: '4px 0 8px' }}>{comment.content}</Paragraph>
          <Space size={4}>
            <Button type="text" size="small" icon={<LikeOutlined />} onClick={() => handleVote('upvote')}>
              {comment.vote_count > 0 ? comment.vote_count : '赞'}
            </Button>
            <Button type="text" size="small" icon={<MessageOutlined />}>
              回复
            </Button>
          </Space>
          {comment.children && comment.children.length > 0 && (
            <div style={{ marginLeft: 24, marginTop: 8 }}>
              <CommentList comments={comment.children} />
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default CommentItem;
