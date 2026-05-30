import api from './api';

// 获取帖子列表（分页）
export const getPosts = (params) => {
  return api.get('/posts', { params });
};

// 获取帖子详情
export const getPostById = (id) => {
  return api.get(`/posts/${id}`);
};

// 创建帖子
export const createPost = (data) => {
  return api.post('/posts', data);
};

// 搜索帖子
export const searchPosts = (params) => {
  return api.get('/posts/search', { params });
};

// 上传帖子图片
export const uploadPostImage = (formData) => {
  return api.post('/posts/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};