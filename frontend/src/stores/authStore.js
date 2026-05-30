import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import * as authService from '@/services/authService';

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
          // api.js 拦截器已解包，response 就是 { token, expire_at }
          const response = await authService.login(credentials);
          const { token, expire_at } = response;

          localStorage.setItem('token', token);

          // 通过 getCurrentUser 获取用户信息
          const userData = await authService.getCurrentUser();

          set({
            user: userData,
            token,
            expireAt: expire_at ? new Date(expire_at).getTime() : null,
            isAuthenticated: true,
            loading: false,
          });

          return { success: true };
        } catch (error) {
          const errorMessage = error.response?.data?.message || error.message || '登录失败';
          set({ error: errorMessage, loading: false });
          return { success: false, error: errorMessage };
        }
      },

      // 注册
      register: async (userData) => {
        set({ loading: true, error: null });
        try {
          // api.js 拦截器已解包，response 就是用户信息对象
          const response = await authService.register(userData);

          // 注册成功后跳转到登录页（后端注册接口不返回 token）
          set({
            user: response,
            isAuthenticated: false,
            loading: false,
          });

          return { success: true };
        } catch (error) {
          const errorMessage = error.response?.data?.message || error.message || '注册失败';
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
          // api.js 拦截器已解包，response 就是用户信息对象
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