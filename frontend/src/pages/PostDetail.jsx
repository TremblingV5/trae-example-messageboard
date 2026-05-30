import React, { useEffect, useState } from 'react';
import { Card, Typography, Button, Input, Avatar, Space, Spin, Empty, message, List, Tag } from 'antd';
import {
  LikeOutlined,
  DislikeOutlined,
  CommentOutlined,
  ArrowLeftOutlined,
  UserOutlined,
  SendOutlined,
  CloseOutlined,
} from '@ant-design/icons';
import { useParams, useNavigate } from 'react-router-dom';
import { usePostStore, useCommentStore, useAuthStore } from '@/stores';

const { Title, Text, Paragraph } = Typography;
const { TextArea } = Input;

// 递归渲染评论树
const CommentTree = ({ comments, parentId = null, onReply, onVote, replyingTo, onCancelReply, onSubmitComment, submitting }) => {
  const childComments = comments.filter((c) => c.parent_id === parentId);

  if (childComments.length === 0) return null;

  return (
    <div>
      {childComments.map((comment) => (
        <div key={comment.id} style={{ marginLeft: parentId ? 32 : 0, marginTop: 12 }}>
          <Card size="small" style={{ background: '#fafafa' }}>
            <div style={{ display: 'flex', alignItems: 'flex-start', gap: 12 }}>
              <Avatar size="small" src={comment.author?.avatar} icon={!comment.author?.avatar && <UserOutlined />} />
              <div style={{ flex: 1 }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 4 }}>
                  <Text strong>{comment.author?.nickname || comment.author?.username || '匿名'}</Text>
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    {comment.created_at && new Date(comment.created_at).toLocaleString('zh-CN')}
                  </Text>
                </div>
                <Paragraph style={{ marginBottom: 8 }}>{comment.content}</Paragraph>
                <Space size="small">
                  <Button
                    type="text"
                    size="small"
                    icon={<LikeOutlined />}
                    onClick={() => onVote(comment.id, 'up')}
                  >
                    {comment.upvotes || 0}
                  </Button>
                  <Button
                    type="text"
                    size="small"
                    icon={<DislikeOutlined />}
                    onClick={() => onVote(comment.id, 'down')}
                  >
                    {comment.downvotes || 0}
                  </Button>
                  <Button
                    type="text"
                    size="small"
                    icon={<CommentOutlined />}
                    onClick={() => onReply(comment.id)}
                  >
                    回复
                  </Button>
                </Space>

                {replyingTo === comment.id && (
                  <div style={{ marginTop: 8 }}>
                    <TextArea
                      rows={2}
                      placeholder={`回复 ${comment.author?.nickname || comment.author?.username || '匿名'}...`}
                      id={`reply-input-${comment.id}`}
                    />
                    <Space style={{ marginTop: 8 }}>
                      <Button
                        type="primary"
                        size="small"
                        icon={<SendOutlined />}
                        loading={submitting}
                        onClick={() => {
                          const input = document.getElementById(`reply-input-${comment.id}`);
                          if (input && input.value.trim()) {
                            onSubmitComment(input.value.trim(), comment.id);
                          }
                        }}
                      >
                        回复
                      </Button>
                      <Button size="small" icon={<CloseOutlined />} onClick={onCancelReply}>
                        取消
                      </Button>
                    </Space>
                  </div>
                )}

                <CommentTree
                  comments={comments}
                  parentId={comment.id}
                  onReply={onReply}
                  onVote={onVote}
                  replyingTo={replyingTo}
                  onCancelReply={onCancelReply}
                  onSubmitComment={onSubmitComment}
                  submitting={submitting}
                />
              </div>
            </div>
          </Card>
        </div>
      ))}
    </div>
  );
};

const PostDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { isAuthenticated, user } = useAuthStore();
  const { currentPost, loading: postLoading, fetchPostById, clearCurrentPost } = usePostStore();
  const { comments, loading: commentLoading, fetchComments, createComment, voteComment, replyingTo, setReplyingTo, cancelReply } = useCommentStore();
  const [commentContent, setCommentContent] = useState('');
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    fetchPostById(id);
    fetchComments(id);

    return () => {
      clearCurrentPost();
    };
  }, [id, fetchPostById, fetchComments, clearCurrentPost]);

  const handleSubmitComment = async (content, parentId = null) => {
    if (!isAuthenticated) {
      message.warning('请先登录');
      return;
    }

    setSubmitting(true);
    try {
      await createComment(id, { content, parent_id: parentId });
      message.success(parentId ? '回复成功' : '评论成功');
      setCommentContent('');
      cancelReply();
    } catch (error) {
      message.error('评论失败');
    } finally {
      setSubmitting(false);
    }
  };

  const handleVote = async (commentId, action) => {
    if (!isAuthenticated) {
      message.warning('请先登录');
      return;
    }
    try {
      await voteComment(commentId, action);
      await fetchComments(id);
    } catch (error) {
      message.error('操作失败');
    }
  };

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

  if (postLoading) {
    return (
      <div style={{ textAlign: 'center', padding: 48 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!currentPost) {
    return <Empty description="帖子不存在" />;
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
        <Title level={2}>{currentPost.title}</Title>

        <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 16 }}>
          <Avatar src={currentPost.author?.avatar} icon={!currentPost.author?.avatar && <UserOutlined />} />
          <div>
            <Text strong>{currentPost.author?.nickname || currentPost.author?.username || '匿名'}</Text>
            <br />
            <Text type="secondary" style={{ fontSize: 12 }}>{formatTime(currentPost.created_at)}</Text>
          </div>
          <div style={{ marginLeft: 'auto', display: 'flex', gap: 8 }}>
            <Tag icon={<CommentOutlined />}>{currentPost.comment_count || 0} 评论</Tag>
            <Tag icon={<LikeOutlined />}>{currentPost.view_count || 0} 浏览</Tag>
          </div>
        </div>

        <Paragraph style={{ fontSize: 16, lineHeight: 1.8, whiteSpace: 'pre-wrap' }}>
          {currentPost.content}
        </Paragraph>
      </Card>

      <Card title={`评论 (${comments.length})`} style={{ marginBottom: 24 }}>
        {isAuthenticated ? (
          <div style={{ marginBottom: 16 }}>
            <TextArea
              rows={3}
              placeholder="发表评论..."
              value={commentContent}
              onChange={(e) => setCommentContent(e.target.value)}
            />
            <Button
              type="primary"
              style={{ marginTop: 8 }}
              loading={submitting}
              onClick={() => handleSubmitComment(commentContent)}
              disabled={!commentContent.trim()}
            >
              发表评论
            </Button>
          </div>
        ) : (
          <div style={{ textAlign: 'center', padding: 16, background: '#f5f5f5', borderRadius: 8, marginBottom: 16 }}>
            <Text type="secondary">请先登录后发表评论</Text>
          </div>
        )}

        <Spin spinning={commentLoading}>
          {comments.length === 0 ? (
            <Empty description="暂无评论，快来抢沙发吧" />
          ) : (
            <CommentTree
              comments={comments}
              parentId={null}
              onReply={setReplyingTo}
              onVote={handleVote}
              replyingTo={replyingTo}
              onCancelReply={cancelReply}
              onSubmitComment={handleSubmitComment}
              submitting={submitting}
            />
          )}
        </Spin>
      </Card>
    </div>
  );
};

export default PostDetail;
