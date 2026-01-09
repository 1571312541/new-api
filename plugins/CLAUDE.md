[根目录](../CLAUDE.md) > **plugins**

# Plugins 模块 - VSCode 扩展插件

> 本模块包含基于 Augment 的 VSCode 扩展插件，实现 AI API 请求拦截与自定义代理功能。

## 模块概述

plugins 目录包含一个修改版的 Augment VSCode 扩展插件，主要功能：

1. **请求拦截**：拦截 Augment 原生 API 请求
2. **自定义代理**：使用自定义 API Key、Provider 和模型进行请求
3. **配置管理**：通过 WebView 面板管理 API 配置
4. **套餐查询**：查询用户余额和套餐信息

## 目录结构

```
plugins/
└── extension/               # Augment 扩展插件
    ├── out/                 # 编译输出目录
    │   ├── extension.js     # 主扩展逻辑（13MB+，已打包）
    │   └── custom-panel.html # 配置面板 WebView
    ├── common-webviews/     # 公共 WebView 资源
    │   └── assets/          # CSS、JS、字体资源
    └── node_modules/        # 依赖模块
```

## 子模块

| 子模块 | 说明 |
|--------|------|
| [extension/](extension/CLAUDE.md) | Augment 扩展核心 - 请求拦截与配置管理 |

## 变更记录 (Changelog)

| 时间 | 操作 | 说明 |
|------|------|------|
| 2025-12-28 | 创建 | 首次生成模块文档 |
