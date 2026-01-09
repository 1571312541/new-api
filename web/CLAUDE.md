[根目录](../../CLAUDE.md) > [web](../) > **src**

# Web 前端模块

> React 前端应用，提供用户界面和管理后台。

## 模块职责

Web 模块实现了 new-api 的用户界面：

1. **用户界面**：登录、注册、个人中心
2. **管理后台**：渠道、令牌、用户管理
3. **数据看板**：统计图表、使用分析
4. **Playground**：API 测试工具
5. **设置页面**：系统配置

## 技术栈

| 技术 | 说明 |
|------|------|
| React 18 | UI 框架 |
| Vite | 构建工具 |
| Semi UI | 组件库 (字节跳动) |
| React Router | 路由管理 |
| Axios | HTTP 请求 |
| i18next | 国际化 |
| VChart | 图表库 |
| Tailwind CSS | 样式工具 |

## 入口与启动

- **入口文件**：`src/main.jsx`
- **应用组件**：`src/App.jsx`

### 开发启动

```bash
cd web
bun install    # 安装依赖
bun run dev    # 开发服务器 (http://localhost:5173)
bun run build  # 生产构建
```

## 目录结构

```
web/src/
├── components/           # 组件目录
│   ├── auth/            # 认证组件
│   ├── common/          # 通用组件
│   ├── dashboard/       # 数据看板
│   ├── layout/          # 布局组件
│   ├── playground/      # Playground
│   ├── settings/        # 设置页面
│   └── setup/           # 初始化向导
├── pages/               # 页面组件
├── helpers/             # 工具函数
├── hooks/               # 自定义 Hooks
├── contexts/            # Context 上下文
├── locales/             # 国际化文件
└── assets/              # 静态资源
```

## 关键组件

### 认证组件 (`components/auth/`)

| 组件 | 说明 |
|------|------|
| `LoginForm.jsx` | 登录表单 |
| `RegisterForm.jsx` | 注册表单 |
| `TwoFAVerification.jsx` | 2FA 验证 |
| `PasswordResetForm.jsx` | 密码重置 |

### 布局组件 (`components/layout/`)

| 组件 | 说明 |
|------|------|
| `PageLayout.jsx` | 页面布局 |
| `headerbar/` | 顶部导航栏 |
| `Footer.jsx` | 页脚 |

### Dashboard (`components/dashboard/`)

| 组件 | 说明 |
|------|------|
| `index.jsx` | 主面板 |
| `StatsCards.jsx` | 统计卡片 |
| `ChartsPanel.jsx` | 图表面板 |
| `ApiInfoPanel.jsx` | API 信息 |

### Playground (`components/playground/`)

| 组件 | 说明 |
|------|------|
| `ChatArea.jsx` | 聊天区域 |
| `SettingsPanel.jsx` | 参数设置 |
| `CodeViewer.jsx` | 代码预览 |
| `DebugPanel.jsx` | 调试面板 |

## 国际化

支持的语言：
- 中文 (zh)
- 英文 (en)
- 法语 (fr)
- 日语 (ja)

语言文件位于 `src/locales/` 目录。

## 构建配置

### Vite 配置

```javascript
// vite.config.js
export default {
  plugins: [react(), semiPlugin()],
  build: {
    outDir: 'dist'
  }
}
```

### 环境变量

| 变量 | 说明 |
|------|------|
| `VITE_REACT_APP_VERSION` | 版本号 |

## 开发规范

1. **组件命名**：PascalCase，如 `UserSettings.jsx`
2. **文件命名**：与组件名一致
3. **样式**：使用 Tailwind CSS 或 Semi UI 组件样式
4. **状态管理**：使用 React Hooks
5. **国际化**：所有文本使用 `t()` 函数

## 常见问题 (FAQ)

**Q: 如何添加新页面？**
A:
1. 在 `components/` 或 `pages/` 创建组件
2. 在路由配置中添加路由
3. 添加国际化文本

**Q: 如何修改主题？**
A: 使用 Semi UI 的主题定制功能

**Q: 如何添加新语言？**
A:
1. 在 `src/locales/` 创建语言文件
2. 在 i18next 配置中注册

## 相关文件清单

| 文件/目录 | 说明 |
|-----------|------|
| `package.json` | 依赖配置 |
| `vite.config.js` | Vite 配置 |
| `tailwind.config.js` | Tailwind 配置 |
| `src/main.jsx` | 入口文件 |
| `src/App.jsx` | 应用组件 |
| `index.html` | HTML 模板 |

## 变更记录 (Changelog)

| 时间 | 操作 | 说明 |
|------|------|------|
| 2025-12-28T20:58:40+0800 | 初始化 | 首次生成模块文档 |
