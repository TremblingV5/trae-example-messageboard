import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import * as authService from '@/services/authService';

const useAuthStore = create(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      loading: false,
      error: null,

      // 登录
      login: async (credentials) => {
        set({ loading: true, error: null });
        try {
          const response = await authService.login(credentials);
          const { token, expire_at } = response.data;
          
          localStorage.setItem('token', token);
          
          // 获取用户信息
          const userResp = await authService.getUserById(credentials.username);
          const user = userResp.data;
          
          set({
            user,
            token,
            isAuthenticated: true,
            loading: false,
          });
          
          return { success: true };
        } catch (error) {
          const errorMessage = error.response?.data?.message || '登录失败';
          set({ error: errorMessage, loading: false });
          return { success: false, error: errorMessage };
        }
      },

      // 注册
      register: async (userData) => {
        set({ loading: true, error: null });
        try {
          const response = await authService.register(userData);
          const user = response.data;
          
          set({
            user,
            isAuthenticated: true,
            loading: false,
          });
          
          return { success: true };
        } catch (error) {
          const errorMessage = error.response?.data?.message || '注册失败';
          set({ error: errorMessage, loading: false });
          return { success: false, error: errorMessage };
        }
      },

      // 退出登录
      logout: () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        set({
          user: null,
          token: null,
          isAuthenticated: false,
          error: null,
        });
      },

      // 获取当前用户信息
      fetchCurrentUser: async () => {
        const token = localStorage.getItem('token');
        if (!token) {
          return;
        }

        set({ loading: true });
        try {
          const user = await authService.getCurrentUser();
          set({
            user,
            token,
            isAuthenticated: true,
            loading: false,
          });
        } catch (error) {
          set({
            user: null,
            token: null,
            isAuthenticated: false,
            loading: false,
          });
        }
      },

      // 更新用户信息
      updateUser: (userData) => {
        set((state) => ({
          user: { ...state.user, ...userData },
        }));
      },

      // 清除错误
      clearError: () => {
        set({ error: null });
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);

export default useAuthStore;