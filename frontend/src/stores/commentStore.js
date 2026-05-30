import { create } from 'zustand';
import * as commentService from '@/services/commentService';

const useCommentStore = create((set, get) => ({
  comments: [],
  loading: false,
  error: null,
  replyingTo: null,

  // 获取评论列表
  fetchComments: async (postId) => {
    set({ loading: true, error: null });
    try {
      const response = await commentService.getComments(postId);
      const comments = response.data || [];      
      set({
        comments,
        loading: false,
      });
      
      return comments;
    } catch (error) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },

  // 创建评论
  createComment: async (postId, data) => {
    set({ loading: true, error: null });
    try {
      const response = await commentService.createComment(postId, data);
      const newComment = response.data;
      
      // 刷新评论列表
      await get().fetchComments(postId);
      
      set({ loading: false, replyingTo: null });
      return newComment;
    } catch (error) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },

  // 点赞评论
  voteComment: async (commentId, action) => {
    try {
      const response = await commentService.voteComment(commentId, action);
      return response.data;
    } catch (error) {
      throw error;
    }
  },

  // 设置回复目标
  setReplyingTo: (commentId) => {
    set({ replyingTo: commentId });
  },

  // 取消回复
  cancelReply: () => {
    set({ replyingTo: null });
  },

  // 清除评论
  clearComments: () => {
    set({ comments: [], replyingTo: null });
  },
}));

export default useCommentStore;