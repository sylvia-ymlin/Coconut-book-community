/**
 * BookCommunity API Client
 * 原生 JavaScript 实现
 */

const API_BASE_URL = 'http://localhost:8080/douyin';

class BookCommunityAPI {
  constructor() {
    this.baseURL = API_BASE_URL;
    this.token = localStorage.getItem('token');
  }

  // 设置token
  setToken(token) {
    this.token = token;
    localStorage.setItem('token', token);
  }

  // 清除token
  clearToken() {
    this.token = null;
    localStorage.removeItem('token');
  }

  // 通用请求方法
  async request(url, options = {}) {
    const headers = {
      'Content-Type': 'application/json',
      ...options.headers
    };

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    try {
      const response = await fetch(`${this.baseURL}${url}`, {
        ...options,
        headers
      });

      const data = await response.json();

      // 处理401错误
      if (response.status === 401) {
        this.clearToken();
        throw new Error('未授权，请重新登录');
      }

      return data;
    } catch (error) {
      console.error('API请求失败:', error);
      throw error;
    }
  }

  // 用户注册
  async register(username, password) {
    const data = await this.request(
      `/user/register/?username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`,
      { method: 'POST' }
    );

    if (data.status_code === 0) {
      this.setToken(data.token);
    }
    return data;
  }

  // 用户登录
  async login(username, password) {
    const data = await this.request(
      `/user/login/?username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`,
      { method: 'POST' }
    );

    if (data.status_code === 0) {
      this.setToken(data.token);
    }
    return data;
  }

  // 获取用户信息
  async getUserInfo(userId) {
    return this.request(`/user/?user_id=${userId}`);
  }

  // 获取推荐
  async getRecommendations(topK = 10) {
    return this.request(`/recommend?top_k=${topK}`);
  }

  // 搜索图书
  async searchBooks(query, topK = 10) {
    return this.request(`/search?q=${encodeURIComponent(query)}&top_k=${topK}`);
  }

  // 获取图书详情
  async getBookDetail(isbn) {
    return this.request(`/book/${isbn}`);
  }
}

// 导出单例
const api = new BookCommunityAPI();

// 如果在浏览器环境中，挂载到window对象
if (typeof window !== 'undefined') {
  window.BookCommunityAPI = api;
}

// Node.js环境导出
if (typeof module !== 'undefined' && module.exports) {
  module.exports = BookCommunityAPI;
}
