# 前后端联调指南 Frontend Integration Guide

## 📋 目录

1. [快速开始](#快速开始)
2. [API 端点](#api-端点)
3. [认证流程](#认证流程)
4. [前端示例代码](#前端示例代码)
5. [错误处理](#错误处理)
6. [开发工具](#开发工具)
7. [常见问题](#常见问题)

---

## 🚀 快速开始

### 1. 启动后端服务

```bash
# 克隆项目
git clone https://github.com/sylvia-ymlin/Coconut-book-community.git
cd Coconut-book-community

# 启动依赖服务
docker-compose up -d

# 复制配置文件
cp config/conf/example.yaml config/conf/config.yaml

# 启动后端
go run main.go
```

**后端地址：** `http://localhost:8080`

### 2. 验证服务

```bash
# 健康检查
curl http://localhost:8080/health

# 预期响应
{
  "status": "healthy",
  "service": "BookCommunity API"
}
```

### 3. 访问 API 文档

**Swagger UI:** http://localhost:8080/swagger/index.html

---

## 📡 API 端点

### Base URL
```
http://localhost:8080/douyin
```

### 用户相关 API

#### 1. 用户注册
```http
POST /douyin/user/register/
```

**请求参数 (Query):**
- `username` (string, required) - 用户名
- `password` (string, required) - 密码

**请求示例:**
```bash
curl -X POST "http://localhost:8080/douyin/user/register/?username=testuser&password=password123"
```

**响应示例:**
```json
{
  "status_code": 0,
  "status_msg": "注册成功",
  "user_id": 1,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 2. 用户登录
```http
POST /douyin/user/login/
```

**请求参数 (Query):**
- `username` (string, required)
- `password` (string, required)

**响应示例:**
```json
{
  "status_code": 0,
  "status_msg": "登录成功",
  "user_id": 1,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 3. 获取用户信息
```http
GET /douyin/user/
```

**请求参数 (Query):**
- `user_id` (int, required) - 用户ID
- `token` (string, optional) - 认证token

**Headers:**
```
Authorization: Bearer <your-jwt-token>
```

**响应示例:**
```json
{
  "status_code": 0,
  "status_msg": "success",
  "user": {
    "id": 1,
    "name": "testuser",
    "follow_count": 10,
    "follower_count": 20,
    "is_follow": false
  }
}
```

### 推荐相关 API

#### 4. 获取个性化推荐
```http
GET /douyin/recommend
```

**请求参数 (Query):**
- `top_k` (int, optional) - 返回结果数量，默认10

**Headers:**
```
Authorization: Bearer <your-jwt-token>
```

**响应示例:**
```json
{
  "status_code": 0,
  "books": [
    {
      "isbn": "9787111544937",
      "title": "深入理解计算机系统",
      "author": "Randal E. Bryant",
      "cover_url": "https://example.com/cover.jpg",
      "rating": 9.7,
      "reason": "基于你的阅读历史推荐",
      "publisher": "机械工业出版社",
      "pub_date": "2016-11",
      "summary": "从程序员的视角看计算机系统..."
    }
  ],
  "total": 10,
  "message": "当前为模拟推荐数据"
}
```

#### 5. 搜索图书
```http
GET /douyin/search
```

**请求参数 (Query):**
- `q` (string, required) - 搜索关键词
- `top_k` (int, optional) - 返回结果数量，默认10

**请求示例:**
```bash
curl "http://localhost:8080/douyin/search?q=golang&top_k=5"
```

**响应示例:**
```json
{
  "status_code": 0,
  "books": [...],
  "total": 5,
  "query": "golang"
}
```

#### 6. 获取图书详情
```http
GET /douyin/book/:isbn
```

**路径参数:**
- `isbn` (string, required) - 图书ISBN号

**请求示例:**
```bash
curl "http://localhost:8080/douyin/book/9787111544937"
```

---

## 🔐 认证流程

### JWT Token 使用

1. **获取 Token**
   - 用户注册或登录后，响应中包含 `token` 字段
   - 保存 token 到浏览器 localStorage/sessionStorage

2. **使用 Token**
   - 在受保护的 API 请求中添加 Authorization Header
   - 格式: `Authorization: Bearer <token>`

3. **Token 刷新**
   - Token 过期后需要重新登录
   - 当前实现未包含 refresh token (可扩展)

### 认证示例 (JavaScript)

```javascript
// 1. 登录并保存token
async function login(username, password) {
  const response = await fetch(
    `http://localhost:8080/douyin/user/login/?username=${username}&password=${password}`,
    { method: 'POST' }
  );
  const data = await response.json();

  if (data.status_code === 0) {
    // 保存token
    localStorage.setItem('token', data.token);
    localStorage.setItem('user_id', data.user_id);
    return data;
  }
  throw new Error(data.status_msg);
}

// 2. 使用token调用受保护API
async function getRecommendations(topK = 10) {
  const token = localStorage.getItem('token');

  const response = await fetch(
    `http://localhost:8080/douyin/recommend?top_k=${topK}`,
    {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    }
  );
  return await response.json();
}

// 3. 退出登录
function logout() {
  localStorage.removeItem('token');
  localStorage.removeItem('user_id');
}
```

---

## 💻 前端示例代码

### React + Axios

#### 1. 安装依赖
```bash
npm install axios
```

#### 2. API 客户端配置

```typescript
// src/api/client.ts
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/douyin';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器 - 添加token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 响应拦截器 - 处理错误
apiClient.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      // Token过期，跳转登录
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default apiClient;
```

#### 3. API 服务封装

```typescript
// src/api/user.ts
import apiClient from './client';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  status_code: number;
  status_msg: string;
  user_id: number;
  token: string;
}

export const userAPI = {
  // 注册
  register: (data: LoginRequest) =>
    apiClient.post<LoginResponse>(
      `/user/register/?username=${data.username}&password=${data.password}`
    ),

  // 登录
  login: (data: LoginRequest) =>
    apiClient.post<LoginResponse>(
      `/user/login/?username=${data.username}&password=${data.password}`
    ),

  // 获取用户信息
  getUserInfo: (userId: number) =>
    apiClient.get(`/user/?user_id=${userId}`)
};
```

```typescript
// src/api/books.ts
import apiClient from './client';

export interface Book {
  isbn: string;
  title: string;
  author: string;
  cover_url: string;
  rating: number;
  reason: string;
  publisher?: string;
  pub_date?: string;
  summary?: string;
}

export const booksAPI = {
  // 获取推荐
  getRecommendations: (topK: number = 10) =>
    apiClient.get<{ status_code: number; books: Book[] }>(
      `/recommend?top_k=${topK}`
    ),

  // 搜索图书
  search: (query: string, topK: number = 10) =>
    apiClient.get<{ status_code: number; books: Book[] }>(
      `/search?q=${encodeURIComponent(query)}&top_k=${topK}`
    ),

  // 获取图书详情
  getBookDetail: (isbn: string) =>
    apiClient.get<{ status_code: number; book: Book }>(
      `/book/${isbn}`
    )
};
```

#### 4. React组件示例

```tsx
// src/components/BookList.tsx
import React, { useEffect, useState } from 'react';
import { booksAPI, Book } from '../api/books';

const BookList: React.FC = () => {
  const [books, setBooks] = useState<Book[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadRecommendations();
  }, []);

  const loadRecommendations = async () => {
    try {
      setLoading(true);
      const response = await booksAPI.getRecommendations(10);
      if (response.status_code === 0) {
        setBooks(response.books);
      } else {
        setError('获取推荐失败');
      }
    } catch (err) {
      setError('网络错误');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div>加载中...</div>;
  if (error) return <div>错误: {error}</div>;

  return (
    <div className="book-list">
      {books.map((book) => (
        <div key={book.isbn} className="book-card">
          <img src={book.cover_url} alt={book.title} />
          <h3>{book.title}</h3>
          <p>{book.author}</p>
          <p>评分: {book.rating}</p>
          <p className="reason">{book.reason}</p>
        </div>
      ))}
    </div>
  );
};

export default BookList;
```

```tsx
// src/components/SearchBox.tsx
import React, { useState } from 'react';
import { booksAPI, Book } from '../api/books';

const SearchBox: React.FC = () => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<Book[]>([]);
  const [searching, setSearching] = useState(false);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;

    try {
      setSearching(true);
      const response = await booksAPI.search(query, 10);
      if (response.status_code === 0) {
        setResults(response.books);
      }
    } catch (err) {
      console.error('搜索失败:', err);
    } finally {
      setSearching(false);
    }
  };

  return (
    <div className="search-box">
      <form onSubmit={handleSearch}>
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="搜索图书..."
        />
        <button type="submit" disabled={searching}>
          {searching ? '搜索中...' : '搜索'}
        </button>
      </form>

      <div className="search-results">
        {results.map((book) => (
          <div key={book.isbn} className="result-item">
            <h4>{book.title}</h4>
            <p>{book.author}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default SearchBox;
```

### Vue 3 + Composition API

```typescript
// src/composables/useBooks.ts
import { ref } from 'vue';
import { booksAPI, Book } from '../api/books';

export function useBooks() {
  const books = ref<Book[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const loadRecommendations = async (topK: number = 10) => {
    try {
      loading.value = true;
      error.value = null;
      const response = await booksAPI.getRecommendations(topK);
      if (response.status_code === 0) {
        books.value = response.books;
      } else {
        error.value = '获取推荐失败';
      }
    } catch (err) {
      error.value = '网络错误';
      console.error(err);
    } finally {
      loading.value = false;
    }
  };

  const search = async (query: string, topK: number = 10) => {
    try {
      loading.value = true;
      error.value = null;
      const response = await booksAPI.search(query, topK);
      if (response.status_code === 0) {
        books.value = response.books;
      }
    } catch (err) {
      error.value = '搜索失败';
    } finally {
      loading.value = false;
    }
  };

  return {
    books,
    loading,
    error,
    loadRecommendations,
    search
  };
}
```

```vue
<!-- src/components/BookList.vue -->
<template>
  <div class="book-list">
    <div v-if="loading">加载中...</div>
    <div v-else-if="error">{{ error }}</div>
    <div v-else class="books-grid">
      <div v-for="book in books" :key="book.isbn" class="book-card">
        <img :src="book.cover_url" :alt="book.title" />
        <h3>{{ book.title }}</h3>
        <p>{{ book.author }}</p>
        <p>评分: {{ book.rating }}</p>
        <p class="reason">{{ book.reason }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useBooks } from '../composables/useBooks';

const { books, loading, error, loadRecommendations } = useBooks();

onMounted(() => {
  loadRecommendations(10);
});
</script>
```

---

## ⚠️ 错误处理

### 状态码说明

| status_code | 说明 |
|-------------|------|
| 0 | 成功 |
| 500 | 服务器内部错误 |
| 401 | Token过期或未授权 |
| 400 | 参数错误 |

### 错误消息

```typescript
const ERROR_MESSAGES: Record<string, string> = {
  '服务器内部错误': 'Internal server error',
  '用户未登录': 'User not logged in',
  '用户不存在': 'User not found',
  '用户已存在': 'User already exists',
  '密码错误': 'Incorrect password',
  '参数错误': 'Invalid parameters'
};
```

### 错误处理示例

```typescript
try {
  const response = await userAPI.login({ username, password });
  if (response.status_code === 0) {
    // 成功
    localStorage.setItem('token', response.token);
  } else {
    // 业务错误
    console.error('登录失败:', response.status_msg);
    alert(response.status_msg);
  }
} catch (error) {
  // 网络错误或其他异常
  console.error('请求失败:', error);
  alert('网络错误，请稍后重试');
}
```

---

## 🛠️ 开发工具

### 1. Postman Collection

导入以下 JSON 到 Postman:

```json
{
  "info": {
    "name": "BookCommunity API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "User Register",
      "request": {
        "method": "POST",
        "url": "http://localhost:8080/douyin/user/register/?username=testuser&password=password123"
      }
    },
    {
      "name": "User Login",
      "request": {
        "method": "POST",
        "url": "http://localhost:8080/douyin/user/login/?username=testuser&password=password123"
      }
    },
    {
      "name": "Get Recommendations",
      "request": {
        "method": "GET",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{token}}"
          }
        ],
        "url": "http://localhost:8080/douyin/recommend?top_k=10"
      }
    },
    {
      "name": "Search Books",
      "request": {
        "method": "GET",
        "url": "http://localhost:8080/douyin/search?q=golang&top_k=5"
      }
    }
  ]
}
```

### 2. curl 测试脚本

```bash
#!/bin/bash
# test_api.sh

BASE_URL="http://localhost:8080/douyin"

echo "=== 1. 健康检查 ==="
curl -s http://localhost:8080/health | jq

echo -e "\n=== 2. 用户注册 ==="
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/user/register/?username=testuser&password=password123")
echo $REGISTER_RESPONSE | jq
TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.token')

echo -e "\n=== 3. 获取推荐 ==="
curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/recommend?top_k=5" | jq

echo -e "\n=== 4. 搜索图书 ==="
curl -s "$BASE_URL/search?q=golang&top_k=3" | jq
```

### 3. 浏览器开发者工具

在浏览器控制台中测试 API:

```javascript
// 测试注册
fetch('http://localhost:8080/douyin/user/register/?username=testuser&password=password123', {
  method: 'POST'
})
  .then(res => res.json())
  .then(data => {
    console.log('注册响应:', data);
    if (data.status_code === 0) {
      localStorage.setItem('token', data.token);
    }
  });

// 测试搜索
fetch('http://localhost:8080/douyin/search?q=计算机&top_k=5')
  .then(res => res.json())
  .then(data => console.log('搜索结果:', data));

// 测试推荐（需要token）
const token = localStorage.getItem('token');
fetch('http://localhost:8080/douyin/recommend?top_k=10', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
})
  .then(res => res.json())
  .then(data => console.log('推荐结果:', data));
```

---

## ❓ 常见问题

### Q1: CORS 错误

**问题:** 浏览器控制台显示 CORS 错误

**解决方案:**
后端已配置 CORS，允许以下来源:
- `http://localhost:3000` (React默认端口)
- `http://localhost:5173` (Vite默认端口)

如需其他端口，修改 `internal/server/server.go`:
```go
corsConfig.AllowOrigins = []string{
    "http://localhost:3000",
    "http://localhost:5173",
    "http://localhost:YOUR_PORT",  // 添加你的端口
}
```

### Q2: Token 过期

**问题:** API 返回 401 Unauthorized

**解决方案:**
```typescript
// 在 axios 拦截器中处理
apiClient.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      // 清除token并跳转登录
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

### Q3: 图片无法显示

**问题:** 图书封面图片404

**原因:** Mock数据中的cover_url可能是示例URL

**解决方案:**
使用占位图或提供默认图片:
```typescript
const coverUrl = book.cover_url || 'https://via.placeholder.com/150x200?text=No+Cover';
```

### Q4: 中文乱码

**问题:** 中文内容显示乱码

**解决方案:**
确保请求头包含正确的编码:
```typescript
headers: {
  'Content-Type': 'application/json; charset=UTF-8'
}
```

---

## 📖 更多资源

- **Swagger API 文档**: http://localhost:8080/swagger/index.html
- **健康检查**: http://localhost:8080/health
- **项目仓库**: https://github.com/sylvia-ymlin/Coconut-book-community

---

## 🎯 下一步

1. ✅ 启动后端服务
2. ✅ 使用 Postman/curl 测试 API
3. ✅ 集成到前端项目
4. ✅ 实现认证流程
5. ✅ 处理错误和边界情况

**Happy Coding! 🚀**
