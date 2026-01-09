# 套餐信息功能修改文档

> 生成时间: 2025-12-28
> 功能版本: v1.0

## 功能概述

本次修改为 Augment 插件添加了**套餐信息显示功能**，用户可以在面板和状态栏中查看当前 API 套餐的余额、总额度和过期时间。

### 主要功能点

1. **套餐信息卡片** - 在面板中新增套餐信息卡片，显示余额、总额度、过期时间
2. **状态栏集成** - 状态栏显示套餐余额，tooltip 显示完整套餐信息
3. **自动刷新** - 定时刷新和手动刷新套餐信息

---

## 涉及文件

| 文件路径 | 修改类型 | 说明 |
|---------|---------|------|
| `plugins/extension/out/extension.js` | 修改 | 添加套餐信息 API 调用和状态栏更新 |
| `plugins/extension/out/custom-panel.html` | 修改 | 添加套餐信息卡片和 JavaScript 函数 |

---

## 详细修改内容

---

### 1. extension.js - 添加命令处理

**文件**: `plugins/extension/out/extension.js`
**位置**: 约第 39322 行

**修改前**:
```javascript
case "getModelConfig":
  handleGetModelConfig(e);
  break;
case "saveModelConfig":
  await handleSaveModelConfig(e, t.configIndex, t.config);
}
```

**修改后**:
```javascript
case "getModelConfig":
  handleGetModelConfig(e);
  break;
case "saveModelConfig":
  await handleSaveModelConfig(e, t.configIndex, t.config);
  break;
// {{ AURA: Add - 套餐信息命令处理 }}
case "refreshQuotaInfo":
  await handleRefreshQuotaInfo(e);
}
```

---

### 2. extension.js - 添加套餐信息 API 函数

**文件**: `plugins/extension/out/extension.js`
**位置**: `handleRefreshBalance` 函数后

**新增代码**:
```javascript
// {{ AURA: Add - 套餐信息 API 调用函数 }}
let cachedQuotaInfo = null;

async function getApiKeyFromConfig() {
  try {
    const fs = require("fs"),
      path = require("path"),
      os = require("os"),
      configPath = path.join(os.homedir(), ".augment", "model-config.json");
    if (fs.existsSync(configPath)) {
      const config = JSON.parse(fs.readFileSync(configPath, "utf8"));
      return config.api_key || null;
    }
  } catch (e) {
    console.error("[QUOTA] 读取 API Key 失败:", e);
  }
  return null;
}

async function fetchQuotaInfo(apiKey) {
  const https = require("https");
  return new Promise((resolve, reject) => {
    const options = {
      hostname: "newapi.stonefancyx.com",
      path: "/api/usage/token/",
      method: "GET",
      headers: {
        Authorization: `Bearer ${apiKey}`,
      },
    };
    const req = https.request(options, (res) => {
      let data = "";
      res.on("data", (chunk) => (data += chunk));
      res.on("end", () => {
        try {
          const json = JSON.parse(data);
          if (json.code === true && json.data) {
            resolve(json.data);
          } else {
            reject(new Error(json.message || "获取套餐信息失败"));
          }
        } catch (e) {
          reject(new Error("解析响应失败"));
        }
      });
    });
    req.on("error", (e) => reject(e));
    req.setTimeout(10000, () => {
      req.destroy();
      reject(new Error("请求超时"));
    });
    req.end();
  });
}

async function refreshQuotaInfoInternal() {
  try {
    const apiKey = await getApiKeyFromConfig();
    if (!apiKey) {
      console.log("[QUOTA] 未配置 API Key");
      return null;
    }
    const quotaData = await fetchQuotaInfo(apiKey);
    cachedQuotaInfo = quotaData;
    // 更新状态栏
    if (globalStatusBarManager && quotaData) {
      globalStatusBarManager.updateQuotaInfo(quotaData);
    }
    return quotaData;
  } catch (e) {
    console.error("[QUOTA] 刷新套餐信息失败:", e);
    return null;
  }
}

async function handleRefreshQuotaInfo(e) {
  try {
    const apiKey = await getApiKeyFromConfig();
    if (!apiKey) {
      e.webview.postMessage({
        command: "quotaInfoLoaded",
        success: false,
        error: "⚠️ 请先配置 API Key",
      });
      return;
    }
    const quotaData = await fetchQuotaInfo(apiKey);
    cachedQuotaInfo = quotaData;
    // 更新状态栏
    if (globalStatusBarManager && quotaData) {
      globalStatusBarManager.updateQuotaInfo(quotaData);
    }
    e.webview.postMessage({
      command: "quotaInfoLoaded",
      success: true,
      data: quotaData,
    });
  } catch (t) {
    e.webview.postMessage({
      command: "quotaInfoLoaded",
      success: false,
      error: t.message,
    });
  }
}
```

---

### 5. extension.js - 状态栏类修改

**文件**: `plugins/extension/out/extension.js`
**位置**: 状态栏管理类

#### 5.1 构造函数添加 quotaInfo 属性

**修改前**:
```javascript
constructor(vscode) {
  ((this.vscode = vscode),
    (this.statusBarItem = null),
    (this.currentState = "notConfigured"),
    (this.updateTimer = null),
    (this.isUpdating = !1),
    this.init());
}
```

**修改后**:
```javascript
constructor(vscode) {
  ((this.vscode = vscode),
    (this.statusBarItem = null),
    (this.currentState = "notConfigured"),
    (this.updateTimer = null),
    (this.isUpdating = !1),
    // {{ AURA: Add - 套餐信息缓存 }}
    (this.quotaInfo = null),
    this.init());
}
```

#### 5.2 新增 updateQuotaInfo 方法

**新增代码**:
```javascript
// {{ AURA: Add - 更新套餐信息 }}
updateQuotaInfo(quotaData) {
  this.quotaInfo = quotaData;
  // 如果当前是正常状态，更新显示（包括文本和tooltip）
  if (this.currentState === "normal" && this.statusBarItem) {
    const available = Math.floor(quotaData.total_available / 100);
    this.statusBarItem.text = `🔋 ${available}`;
    this.statusBarItem.color = this.getQuotaColor(available);
    // {{ AURA: Add - 同时更新 tooltip }}
    this.statusBarItem.tooltip = this.generateTooltip(this.lastAccountInfo || {}, "点击打开积分面板");
  }
}
```

#### 5.3 新增辅助方法

**新增代码**:
```javascript
getQuotaColor(available) {
  if (available <= 0) return "#ff4444";
  if (available < 1000) return "#ffaa00";
  return void 0;
}

formatQuotaExpires(expiresAt) {
  if (!expiresAt) return "永久有效";
  const date = new Date(expiresAt * 1000);
  return date.toLocaleString();
}
```

#### 5.4 修改 generateTooltip 方法

**修改前**:
```javascript
generateTooltip(e, t) {
  return `Augment 账号信息\n        邮箱账号：${e?.email || "待获取"}\n        套餐名称：${e?.plan_name || "待获取"}\n        到期时间：${null === e?.end_date ? "无期限" : e?.end_date || "待获取"}\n        剩余积分：${e?.balance ? this.formatBalance(e.balance) : "待获取"}\n        ${t}`;
}
```

**修改后**:
```javascript
// {{ AURA: Modify - 修改 tooltip 显示，增加套餐信息 }}
generateTooltip(e, t) {
  let quotaSection = "";
  if (this.quotaInfo) {
    const available = Math.floor(this.quotaInfo.total_available / 100);
    const granted = Math.floor(this.quotaInfo.total_granted / 100);
    const expires = this.formatQuotaExpires(this.quotaInfo.expires_at);
    quotaSection = `\n        ────────────────\n        💳 套餐信息\n        当前余额：${available} / ${granted}\n        过期时间：${expires}`;
  }
  return `Augment 账号信息\n        邮箱账号：${e?.email || "待获取"}\n        套餐名称：${e?.plan_name || "待获取"}\n        到期时间：${null === e?.end_date ? "无期限" : e?.end_date || "待获取"}\n        剩余积分：${e?.balance ? this.formatBalance(e.balance) : "待获取"}${quotaSection}\n        ${t}`;
}
```

#### 5.5 修改 setNormal 方法

**修改前**:
```javascript
setNormal(e) {
  this.currentState = "normal";
  const t = this.formatBalance(e.balance);
  ((this.statusBarItem.text = `🔋 ${t}`),
    (this.statusBarItem.tooltip = this.generateTooltip(e, "点击打开积分面板")),
    (this.statusBarItem.backgroundColor = void 0),
    (this.statusBarItem.color = this.getBalanceColor(e.balance)));
}
```

**修改后**:
```javascript
// {{ AURA: Modify - 修改 setNormal，优先显示套餐余额，并保存账号信息 }}
setNormal(e) {
  this.currentState = "normal";
  // {{ AURA: Add - 保存账号信息用于后续 tooltip 更新 }}
  this.lastAccountInfo = e;
  let displayBalance;
  let balanceColor;
  if (this.quotaInfo) {
    const available = Math.floor(this.quotaInfo.total_available / 100);
    displayBalance = available;
    balanceColor = this.getQuotaColor(available);
  } else {
    displayBalance = this.formatBalance(e.balance);
    balanceColor = this.getBalanceColor(e.balance);
  }
  ((this.statusBarItem.text = `🔋 ${displayBalance}`),
    (this.statusBarItem.tooltip = this.generateTooltip(e, "点击打开积分面板")),
    (this.statusBarItem.backgroundColor = void 0),
    (this.statusBarItem.color = balanceColor));
}
```

#### 5.6 修改 refreshBalance 方法

**修改前**:
```javascript
async refreshBalance(e = !1) {
  if (this.isInitialized() && (!this.isUpdating || e)) {
    this.isUpdating = !0;
    try {
      if (!SystemApiService.isTokenConfigured(vscodeModule))
        return void this.setNotConfigured();
      this.setLoading();
      const e = await SystemApiService.getFullAccountInfo(vscodeModule);
      e.isConfigured
        ? this.setNormal({...})
        : this.setNotConfigured();
    } catch (e) {...}
  }
}
```

**修改后**:
```javascript
// {{ AURA: Modify - 在刷新余额时同时刷新套餐信息 }}
async refreshBalance(e = !1) {
  if (this.isInitialized() && (!this.isUpdating || e)) {
    this.isUpdating = !0;
    try {
      if (!SystemApiService.isTokenConfigured(vscodeModule))
        return void this.setNotConfigured();
      this.setLoading();
      // 同时刷新套餐信息
      refreshQuotaInfoInternal().catch(err => console.error("[QUOTA] 刷新失败:", err));
      const e = await SystemApiService.getFullAccountInfo(vscodeModule);
      e.isConfigured
        ? this.setNormal({...})
        : this.setNotConfigured();
    } catch (e) {...}
  }
}
```

---

### 6. custom-panel.html - 套餐信息卡片

**文件**: `plugins/extension/out/custom-panel.html`
**位置**: balance-card 上方

**新增 HTML**:
```html
<!-- {{ AURA: Add - 套餐信息卡片，显示当前余额和过期时间 }} -->
<div class="card" id="quota-info-card">
  <div class="card-header">
    <h2>💳 套餐信息</h2>
    <button class="button small" onclick="refreshQuotaInfo()">🔄 刷新</button>
  </div>
  <div class="config-section">
    <div class="info-row">
      <span class="label">当前余额:</span>
      <span class="value" id="quota-available"><span class="loading">⏳ 加载中...</span></span>
    </div>
    <div class="info-row">
      <span class="label">总额度:</span>
      <span class="value" id="quota-granted">-</span>
    </div>
    <div class="info-row">
      <span class="label">过期时间:</span>
      <span class="value" id="quota-expires">-</span>
    </div>
    <div class="config-status" id="quota-status"></div>
  </div>
</div>
```

---

### 7. custom-panel.html - JavaScript 函数

**文件**: `plugins/extension/out/custom-panel.html`

#### 7.1 新增 refreshQuotaInfo 函数

**新增代码**:
```javascript
// {{ AURA: Add - 套餐信息相关函数 }}
// 刷新套餐信息
function refreshQuotaInfo() {
  console.log('[PANEL] 刷新套餐信息...');
  document.getElementById('quota-available').innerHTML = '<span class="loading">⏳ 加载中...</span>';
  document.getElementById('quota-status').textContent = '';

  vscode.postMessage({
    command: 'refreshQuotaInfo'
  });
}
```

#### 7.2 新增 updateQuotaDisplay 函数

**新增代码**:
```javascript
// 更新套餐信息显示
function updateQuotaDisplay(quotaInfo) {
  console.log('[PANEL] 更新套餐信息显示:', quotaInfo);
  const availableEl = document.getElementById('quota-available');
  const grantedEl = document.getElementById('quota-granted');
  const expiresEl = document.getElementById('quota-expires');
  const statusEl = document.getElementById('quota-status');

  if (quotaInfo && quotaInfo.success) {
    const data = quotaInfo.data;
    // {{ AURA: Modify - 格式化余额（除以100，显示整数） }}
    const available = Math.floor(data.total_available / 100);
    const granted = Math.floor(data.total_granted / 100);

    availableEl.textContent = `${available}`;
    availableEl.style.color = parseFloat(available) > 1000 ? 'var(--vscode-testing-iconPassed)' :
                               parseFloat(available) > 100 ? 'var(--vscode-editorWarning-foreground)' :
                               'var(--vscode-errorForeground)';

    grantedEl.textContent = `${granted}`;

    // 格式化过期时间
    if (data.expires_at) {
      const expiresDate = new Date(data.expires_at * 1000);
      expiresEl.textContent = expiresDate.toLocaleString();
    } else {
      expiresEl.textContent = '永久有效';
    }

    statusEl.className = 'config-status success';
    statusEl.textContent = '✅ 套餐信息已更新';
    setTimeout(() => {
      statusEl.textContent = '';
      statusEl.className = 'config-status';
    }, 3000);
  } else {
    availableEl.innerHTML = '<span class="error">获取失败</span>';
    grantedEl.textContent = '-';
    expiresEl.textContent = '-';
    statusEl.className = 'config-status warning';
    statusEl.textContent = quotaInfo?.error || '⚠️ 无法获取套餐信息';
  }
}
```

#### 7.3 新增消息处理

**修改前**:
```javascript
case 'modelConfigSaved':
  // ...
  break;
}
```

**修改后**:
```javascript
case 'modelConfigSaved':
  // ...
  break;
// {{ AURA: Add - 套餐信息消息处理 }}
case 'quotaInfoLoaded':
  console.log('[PANEL] 收到套餐信息:', message);
  updateQuotaDisplay(message);
  break;
}
```

#### 7.4 页面加载时自动刷新

**修改前**:
```javascript
// 加载所有功能
loadCurrentToken();
refreshBalance();

// {{ AURA: Modify - 只显示 balance-card，隐藏其他面板 }}
```

**修改后**:
```javascript
// 加载所有功能
loadCurrentToken();
refreshBalance();
// {{ AURA: Add - 加载套餐信息 }}
refreshQuotaInfo();

// {{ AURA: Modify - 只显示 balance-card，隐藏其他面板 }}
```

---

## API 接口说明

### 套餐信息 API

**请求**:
```
GET https://newapi.stonefancyx.com/api/usage/token/
Authorization: Bearer {api_key}
```

**响应**:
```json
{
  "code": true,
  "data": {
    "total_available": 500000,  // 可用额度（需除以250）
    "total_granted": 1000000,   // 总额度（需除以250）
    "expires_at": 1735689600    // Unix 时间戳（秒）
  }
}
```

**数据处理**:
- `total_available / 250` = 实际可用余额
- `total_granted / 250` = 实际总额度
- `expires_at` = Unix 时间戳，转换为本地时间显示

---

## 测试验证

### 功能测试清单

- [ ] 套餐信息卡片正常显示
- [ ] 点击刷新按钮能获取最新套餐信息
- [ ] 状态栏显示套餐余额（整数格式）
- [ ] 状态栏 tooltip 显示完整套餐信息
- [ ] 定时刷新时同步更新套餐信息
- [ ] 余额颜色根据数值变化（绿色 > 1000，黄色 > 100，红色 <= 100）

---

## 版本历史

| 版本 | 日期 | 修改内容 |
|-----|------|---------|
| v1.0 | 2025-12-28 | 初始版本，添加套餐信息功能 |

---

## 代码混淆加密与打包

### 概述

为保护插件源代码，提供了混淆加密和压缩工具，可对 `extension.js` 和 `custom-panel.html` 进行处理。

### 依赖安装

```powershell
cd plugins\extension
npm install javascript-obfuscator terser html-minifier-terser cheerio --save-dev --legacy-peer-deps
```

### 构建脚本

| 脚本文件 | 用途 |
|---------|------|
| `build-obfuscate.js` | 混淆 HTML 文件中的内联 JS + 尝试混淆 extension.js |
| `build-extension-minify.js` | 仅压缩 extension.js（更稳定，推荐） |

### 使用方法

#### 1. 混淆 HTML 文件

```powershell
cd plugins\extension
node build-obfuscate.js
```

**输出示例**:
```
🔐 开始混淆加密和压缩...

📄 处理 extension.js...
   💾 已备份到: backup\extension.js.xxx.bak
   🗜️ 压缩中...
   🔒 混淆中...
📄 处理 out/custom-panel.html...
   💾 已备份到: backup\out-custom-panel.html.xxx.bak
   📝 发现 2 个内联脚本
   ✅ 完成! 50.47 KB → 83.08 KB
📄 处理 common-webviews/custom-panel.html...
   ✅ 完成! 46.76 KB → 83.90 KB

✅ 混淆加密和压缩完成！
📁 原始文件已备份到: backup
```

#### 2. 压缩 extension.js

```powershell
cd plugins\extension
node build-extension-minify.js
```

**输出示例**:
```
🗜️ 开始压缩 extension.js...

💾 已备份到: backup\extension.js.xxx.bak

✅ 压缩完成!
   原始大小: 12.50 MB
   压缩后:   8.20 MB
   减少:     34.4%
```

### 打包 VSIX

完成混淆和压缩后，使用 `vsce` 打包：

```powershell
cd plugins\extension
npx vsce package --no-dependencies --allow-star-activation --skip-license
```

**输出示例**:
```
DONE  Packaged: vscode-augment-0.696.2.vsix (1337 files, 18.26MB)
```

### 安装插件

```powershell
# 方式1: 命令行安装
code --install-extension vscode-augment-0.696.2.vsix

# 方式2: VSCode 命令面板
# Ctrl+Shift+P → Extensions: Install from VSIX... → 选择 .vsix 文件
```

### 备份与恢复

所有原始文件在处理前会自动备份到 `plugins/extension/backup/` 目录。

**恢复方法**:
```powershell
# 恢复 extension.js
copy backup\extension.js.xxx.bak out\extension.js

# 恢复 custom-panel.html
copy backup\out-custom-panel.html.xxx.bak out\custom-panel.html
```

### 混淆配置说明

`build-obfuscate.js` 中的混淆配置：

| 配置项 | 值 | 说明 |
|-------|-----|------|
| `controlFlowFlattening` | true | 控制流扁平化 |
| `deadCodeInjection` | true | 死代码注入 |
| `stringArray` | true | 字符串数组化 |
| `stringArrayEncoding` | base64 | 字符串编码方式 |
| `identifierNamesGenerator` | hexadecimal | 变量名生成方式 |
| `renameGlobals` | false | 不重命名全局变量（保护 VSCode API） |

### 注意事项

1. **extension.js 混淆可能失败** - 由于文件过大且包含特殊字符，建议使用 `build-extension-minify.js` 仅压缩
2. **混淆后文件变大是正常的** - 混淆器添加了保护代码
3. **保留 VSCode API** - 配置中已排除 `vscode`、`acquireVsCodeApi` 等关键变量
4. **测试验证** - 打包后务必测试插件功能是否正常

---

# 调试日志功能

> 添加时间: 2025-12-29
> 版本: v1.1

## 功能概述

为 Augment 插件添加调试日志功能，在 VSCode 开发者工具控制台中打印请求/响应数据，便于排查问题。

### 日志标识符

所有调试日志使用 `[AUGMENT-DEBUG]` 前缀，便于过滤和查找。

---

## 涉及文件

| 文件路径 | 修改类型 | 说明 |
|---------|---------|------|
| `plugins/extension/out/extension.js` | 修改 | 在 4 个位置添加 console.log 调试日志 |

---

## 详细修改内容

### 1. YD 函数入口 - 响应解析日志

**文件**: `plugins/extension/out/extension.js`
**位置**: 约第 6687 行
**作用**: 记录 Augment 响应解析器的输入数据

**修改前**:
```javascript
function YD(e) {
  let t = {
    text: kl("BackChatResult", "text", e.text),
```

**修改后**:
```javascript
function YD(e) {
  // {{ AURA: Add - 调试日志：记录 BackChatResult 解析输入 }}
  console.log("[AUGMENT-DEBUG] YD Input:", JSON.stringify(e, null, 2));
  let t = {
    text: kl("BackChatResult", "text", e.text),
```

---

### 2. HTTP 请求发送 - 请求日志

**文件**: `plugins/extension/out/extension.js`
**位置**: 约第 176287 行
**作用**: 记录发送到 API 的请求 URL、Headers 和 Body

**修改前**:
```javascript
(p && (e.Authorization = `Bearer ${p}`),
  await this.signRequest(n, e),
```

**修改后**:
```javascript
(p && (e.Authorization = `Bearer ${p}`),
  // {{ AURA: Add - 调试日志：记录请求详情 }}
  console.log("[AUGMENT-DEBUG] Request URL:", d.toString()),
  console.log("[AUGMENT-DEBUG] Request Headers:", JSON.stringify(e, null, 2)),
  console.log("[AUGMENT-DEBUG] Request Body:", f),
  await this.signRequest(n, e),
```

---

### 3. 响应行解析 - 原始响应日志

**文件**: `plugins/extension/out/extension.js`
**位置**: 约第 176332 行
**作用**: 记录从服务器接收的每一行原始响应数据

**修改前**:
```javascript
let e = r.indexOf("\n"),
  t = r.substring(0, e);
r = r.substring(e + 1);
try {
  let e = JSON.parse(t);
  yield a(e);
```

**修改后**:
```javascript
let e = r.indexOf("\n"),
  t = r.substring(0, e);
// {{ AURA: Add - 调试日志：记录原始响应行 }}
console.log("[AUGMENT-DEBUG] Raw Response Line:", t);
r = r.substring(e + 1);
try {
  let e = JSON.parse(t);
  // {{ AURA: Add - 调试日志：记录解析后的 JSON }}
  console.log("[AUGMENT-DEBUG] Parsed JSON:", JSON.stringify(e, null, 2));
  yield a(e);
```

---

## 使用方法

### 查看调试日志

1. 在 VSCode 中按 `Ctrl+Shift+I`（或 `Cmd+Shift+I`）打开开发者工具
2. 切换到 **Console** 标签页
3. 在过滤器中输入 `AUGMENT-DEBUG` 过滤日志
4. 使用 Augment 插件发送请求，观察日志输出

### 日志输出示例

```
[AUGMENT-DEBUG] Request URL: https://api.example.com/v1/chat-stream
[AUGMENT-DEBUG] Request Headers: {
  "Authorization": "Bearer sk-xxx",
  "Content-Type": "application/json"
}
[AUGMENT-DEBUG] Request Body: {"message":"Hello","model":"gpt-4o"}
[AUGMENT-DEBUG] Raw Response Line: {"text":"Hello"}
[AUGMENT-DEBUG] Parsed JSON: {
  "text": "Hello"
}
[AUGMENT-DEBUG] YD Input: {
  "text": "Hello"
}
```

---

## 注意事项

1. **生产环境移除** - 调试日志会影响性能，正式发布前应移除或禁用
2. **敏感信息** - 日志中可能包含 API Key，注意不要泄露
3. **日志量大** - 流式响应会产生大量日志，建议仅在调试时使用

---

## 版本历史

| 版本 | 日期 | 修改内容 |
|-----|------|---------|
| v1.0 | 2025-12-28 | 初始版本，添加套餐信息功能 |
| v1.1 | 2025-12-29 | 添加调试日志功能 |

