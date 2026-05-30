import api from './api';

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