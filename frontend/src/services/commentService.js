import api from './api';

// 获取帖子评论列表
export const getComments = (postId) => {
  return api.get(`/posts/${postId}/comments`);
};

// 创建评论
export const createComment = (postId, data) => {
  return api.post(`/posts/${postId}/comments`, data);
};

// 点赞评论
export const voteComment = (commentId, action) => {
  return api.post(`/comments/${commentId}/vote`, { action });
};

// 获取评论点赞数
export const getCommentVotes = (commentId) => {
  return api.get(`/comments/${commentId}/votes`);
};