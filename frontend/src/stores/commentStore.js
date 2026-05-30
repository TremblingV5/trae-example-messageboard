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
      // api.js 拦截器已解包，response 就是评论数组
      const comments = await commentService.getComments(postId);
      set({
        comments: comments || [],
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
      // api.js 拦截器已解包
      await commentService.createComment(postId, data);

      // 刷新评论列表
      await get().fetchComments(postId);

      set({ loading: false, replyingTo: null });
    } catch (error) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },

  // 点赞评论
  voteComment: async (commentId, action) => {
    try {
      // api.js 拦截器已解包，response 就是 { vote_count, voted, value }
      const result = await commentService.voteComment(commentId, action);
      return result;
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