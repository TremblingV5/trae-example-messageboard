import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import * as authService from '@/services/authService';
import * as userService from '@/services/userService';

const useAuthStore = create(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      expireAt: null,
      isAuthenticated: false,
      loading: false,
      error: null,

      // 检查 token 是否过期
      isTokenExpired: () => {
        const { expireAt } = get();
        if (!expireAt) return true;
        return Date.now() > expireAt;
      },

      // 登录
      login: async (credentials) => {
        set({ loading: true, error: null });
        try {
          const response = await authService.login(credentials);
          const { token, expire_at, user } = response.data;

          localStorage.setItem('token', token);

          // 如果登录接口返回了用户信息，直接使用
          let userData = user;
          if (!userData) {
            // 否则通过 getCurrentUser 获取
            const meResp = await authService.getCurrentUser();
            userData = meResp.data;
          }

          set({
            user: userData,
            token,
            expireAt: expire_at ? new Date(expire_at).getTime() : null,
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
          const { token, expire_at, user } = response.data;

          // 注册成功后设置 token
          if (token) {
            localStorage.setItem('token', token);
          }

          set({
            user,
            token,
            expireAt: expire_at ? new Date(expire_at).getTime() : null,
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
          expireAt: null,
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

        // 检查 token 是否过期
        if (get().isTokenExpired()) {
          get().logout();
          return;
        }

        set({ loading: true });
        try {
          const response = await authService.getCurrentUser();
          const user = response.data;
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
            expireAt: null,
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
        expireAt: state.expireAt,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);

export default useAuthStore;
