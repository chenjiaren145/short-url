# 前端升级计划 - 公开短链接服务

## 概述

将短链接服务从内部工具升级为公开服务，核心改造点：用户认证 + 链接管理。

---

## 优先级说明

- **P0**: 必须实现，核心功能
- **P1**: 重要功能，用户体验关键
- **P2**: 增强功能，锦上添花
- **P3**: 优化改进，可后续迭代

---

## 改造计划

### Phase 1: 用户认证系统 (P0)

| 序号 | 任务 | 涉及文件 | 说明 |
|------|------|----------|------|
| 1.1 | 创建类型定义 | `src/types/auth.ts` | User, LoginRequest, RegisterRequest, AuthResponse |
| 1.2 | 创建认证 API | `src/services/auth.ts` | login(), register(), logout(), getCurrentUser() |
| 1.3 | 创建 AuthContext | `src/contexts/AuthContext.tsx` | 全局用户状态管理，提供 user, login, logout |
| 1.4 | 创建登录页面 | `src/pages/Login.tsx` | 表单 + 调用 login API |
| 1.5 | 创建注册页面 | `src/pages/Register.tsx` | 表单 + 调用 register API |
| 1.6 | Token 管理 | `src/utils/token.ts` | 存储/读取/清除 localStorage token |
| 1.7 | 创建路由守卫 | `src/components/ProtectedRoute.tsx` | 未登录跳转登录页 |
| 1.8 | 更新入口文件 | `src/main.tsx` | 包裹 AuthProvider，配置路由 |
| 1.9 | 创建导航栏 | `src/components/Navbar.tsx` | 显示用户信息 + 登出按钮 |

### Phase 2: 链接管理功能 (P1)

| 序号 | 任务 | 涉及文件 | 说明 |
|------|------|----------|------|
| 2.1 | 更新类型定义 | `src/types/url.ts` | 添加 id, userId, createdAt, updatedAt |
| 2.2 | 更新 API 服务 | `src/services/api.ts` | getMyUrls(), deleteUrl(), updateUrl() |
| 2.3 | 创建仪表盘页面 | `src/pages/Dashboard.tsx` | 展示用户所有链接，替换现有主页逻辑 |
| 2.4 | 创建链接删除功能 | `src/components/URLList.tsx` | 添加删除按钮 + 确认弹窗 |
| 2.5 | 创建链接编辑功能 | `src/components/URLEditModal.tsx` | 编辑原始 URL（可选：自定义短码） |
| 2.6 | 分页/无限滚动 | `src/pages/Dashboard.tsx` | 链接数量多时优化性能 |

### Phase 3: 用户体验优化 (P2)

| 序号 | 任务 | 涉及文件 | 说明 |
|------|------|----------|------|
| 3.1 | 全局错误处理 | `src/utils/errorHandler.ts` | 统一处理 API 错误，401 自动跳转登录 |
| 3.2 | Loading 状态优化 | 各组件 | 添加骨架屏或 Spinner |
| 3.3 | Toast 通知 | `src/components/Toast.tsx` | 操作成功/失败提示 |
| 3.4 | 表单验证优化 | `src/pages/Login.tsx`, `Register.tsx` | 实时验证 + 错误提示 |
| 3.5 | 首页改造 | `src/pages/Home.tsx` | 公开首页，展示服务介绍 + 快速生成入口 |

### Phase 4: 高级功能 (P3)

| 序号 | 任务 | 涉及文件 | 说明 |
|------|------|----------|------|
| 4.1 | 配额显示 | `src/components/QuotaDisplay.tsx` | 显示已用/总额度 |
| 4.2 | 链接搜索过滤 | `src/pages/Dashboard.tsx` | 按关键词搜索链接 |
| 4.3 | 链接排序 | `src/pages/Dashboard.tsx` | 按时间/访问量排序 |
| 4.4 | 数据导出 | `src/pages/Dashboard.tsx` | 导出 CSV/JSON |
| 4.5 | 暗色模式 | `tailwind.config.js`, 全局组件 | 主题切换 |

---

## 文件结构预览

```
frontend/src/
├── components/
│   ├── Navbar.tsx           # P1.9
│   ├── ProtectedRoute.tsx   # P1.7
│   ├── Toast.tsx            # P2.3
│   ├── URLEditModal.tsx     # P2.5
│   ├── URLForm.tsx          # 已有，可能需修改
│   └── URLList.tsx          # 已有，需改造
├── contexts/
│   └── AuthContext.tsx      # P1.3
├── pages/
│   ├── Dashboard.tsx        # P2.3
│   ├── Home.tsx             # P2.5
│   ├── Login.tsx            # P1.4
│   └── Register.tsx         # P1.5
├── services/
│   ├── api.ts               # 已有，需改造
│   └── auth.ts              # P1.2
├── types/
│   ├── auth.ts              # P1.1
│   └── url.ts               # 已有，需改造
├── utils/
│   ├── errorHandler.ts      # P3.1
│   └── token.ts             # P1.6
├── App.tsx                  # 路由配置
└── main.tsx                 # 入口 + Provider
```

---

## 路由规划

| 路径 | 组件 | 权限 |
|------|------|------|
| `/` | Home | 公开 |
| `/login` | Login | 公开（已登录跳转 Dashboard） |
| `/register` | Register | 公开（已登录跳转 Dashboard） |
| `/dashboard` | Dashboard | 需登录 |

---

## 后端 API 依赖

前端改造需要后端配合实现以下接口：

| 接口 | 方法 | 说明 |
|------|------|------|
| `/auth/register` | POST | 用户注册 |
| `/auth/login` | POST | 用户登录 |
| `/auth/me` | GET | 获取当前用户信息 |
| `/urls` | GET | 获取当前用户链接列表 |
| `/urls/:id` | PUT | 更新链接 |
| `/urls/:id` | DELETE | 删除链接 |

---

## 学习要点

- **React Context**: 全局状态管理
- **React Router**: 路由配置 + 守卫
- **Token 认证**: JWT 存储与请求拦截
- **表单处理**: 受控组件 + 验证
- **状态管理**: 组件间数据流

---

## 建议实施顺序

```
Phase 1 (P0) → Phase 2 (P1) → Phase 3 (P2) → Phase 4 (P3)
```

建议先完成 P0 全部任务，验证可行后再继续 P1。
