import api from './api';

// 用户登录
export const login = (data) => {
  return api.post('/auth/login', data);
};

// 用户注册
export const register = (data) => {
  return api.post('/auth/register', data);
};

// 获取当前用户信息
export const getCurrentUser = () => {
  return api.get('/auth/me');
};

// 获取用户信息
export const getUserById = (id) => {
  return api.get(`/users/${id}`);
};

// 更新用户信息
export const updateUser = (id, data) => {
  return api.put(`/users/${id}`, data);
};

// 上传头像
export const uploadAvatar = (id, formData) => {
  return api.post(`/users/${id}/avatar`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
};