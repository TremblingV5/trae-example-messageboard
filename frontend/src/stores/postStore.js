import { create } from 'zustand';
import * as postService from '@/services/postService';

const usePostStore = create((set, get) => ({
  posts: [],
  currentPost: null,
  pagination: {
    page: 1,
    pageSize: 10,
    total: 0,
  totalPages: 0,
  },
  loading: false,
  error: null,

  // 获取帖子列表
  fetchPosts: async (params = {}) => {
    set({ loading: true, error: null });
    try {
      const response = await postService.getPosts({
        page: params.page || get().pagination.page,
        pageSize: params.pageSize || get().pagination.pageSize,
        keyword: params.keyword,
      });
      
      const data = response.data;
      
      set({
        posts: data.posts || [],
        pagination: {
          page: data.page || 1,
          pageSize: data.page_size || 10,
          total: data.total || 0,
          totalPages: Math.ceil((data.total || 0) / (data.page_size || 10)),
        },
        loading: false,
      });
      
      return response;
    } catch (error) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },

  // 获取帖子详情
  fetchPostById: async (id) => {
    set({ loading: true, error: null });
    try {
      const response = await postService.getPostById(id);
      set({ currentPost: response.data, loading: false });
      return response.data;
    } catch (error) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },

  // 创建帖子
  createPost: async (data) => {
    set({ loading: true, error: null });
    try {
      const response = await postService.createPost(data);
      const newPost = response.data;
      set((state) => ({
        posts: [newPost, ...state.posts],
        loading: false,
      }));
      return newPost;
    } catch (error) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },

  // 搜索帖子
  searchPosts: async (keyword, params = {}) => {
    set({ loading: true, error: null });
    try {
      const response = await postService.searchPosts({
        keyword,
        page: params.page || 1,
        pageSize: params.pageSize || 10,
      });
      
      const data = response.data;
      
      set({
        posts: data.posts || [],
        pagination: {
          page: data.page || 1,
          pageSize: data.page_size || 10,
          total: data.total || 0,
          totalPages: Math.ceil((data.total || 0) / (data.page_size || 10)),
        },
        loading: false,
      });
      
      return response;
    } catch (error) {
      set({ error: error.message, loading: false });
      throw error;
    }
  },

  // 清除当前帖子
  clearCurrentPost: () => {
    set({ currentPost: null });
  },

  // 重置状态
  reset: () => {
    set({
      posts: [],
      currentPost: null,
      pagination: {
        page: 1,
        pageSize: 10,
        total: 0,
        totalPages: 0,
      },
      loading: false,
      error: null,
    });
  },
}));

export default usePostStore;