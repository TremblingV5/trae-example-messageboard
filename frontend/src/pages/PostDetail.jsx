import React, { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Spin, Button, Card, Typography, Divider } from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';
import { usePostStore, useCommentStore, useAuthStore } from '@/stores';
import CommentList from '@/components/CommentList';
import CommentForm from '@/components/CommentForm';

const { Title, Text, Paragraph } = Typography;

const PostDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { currentPost, loading, fetchPostById } = usePostStore();
  const { comments, fetchComments, createComment } = useCommentStore();
  const { isAuthenticated } = useAuthStore();

  useEffect(() => {
    if (id) {
      fetchPostById(id);
      fetchComments(id);
    }
  }, [id]);

  const handleComment = async (content, parentId = null) => {
    await createComment(id, { content, parent_id: parentId });
  };

  if (loading) return <Spin style={{ display: 'block', margin: '100px auto' }} />;
  if (!currentPost) return <div style={{ textAlign: 'center', padding: 100 }}>帖子不存在</div>;

  return (
    <div style={{ maxWidth: 800, margin: '0 auto', padding: '24px' }}>
      <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/')} style={{ marginBottom: 16 }}>
        返回首页
      </Button>
      <Card>
        <Title level={2}>{currentPost.title}</Title>
        <div style={{ marginBottom: 16, color: '#888' }}>
          {currentPost.author && (
            <Text type="secondary">
              作者：{currentPost.author.nickname || currentPost.author.username} |
              发布于：{new Date(currentPost.created_at).toLocaleString()}
            </Text>
          )}
        </div>
        {currentPost.image_url && (
          <div style={{ marginBottom: 16 }}>
            <img src={currentPost.image_url} alt="帖子图片" style={{ maxWidth: '100%', borderRadius: 8 }} />
          </div>
        )}
        <Paragraph style={{ fontSize: 16, lineHeight: 1.8 }}>{currentPost.content}</Paragraph>
      </Card>
      <Divider />
      <Title level={4}>评论 ({comments.length})</Title>
      {isAuthenticated && <CommentForm onSubmit={handleComment} />}
      <CommentList comments={comments} />
    </div>
  );
};

export default PostDetail;
