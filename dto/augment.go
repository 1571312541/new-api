package dto

// AugmentChatStreamRequest Augment 插件 chat-stream 请求结构
// 参考：plugins/extension/out/extension.js chatStream 函数
type AugmentChatStreamRequest struct {
	Model              string                  `json:"model,omitempty"`
	Path               string                  `json:"path,omitempty"`
	Prefix             string                  `json:"prefix,omitempty"`
	SelectedCode       string                  `json:"selected_code,omitempty"`
	Suffix             string                  `json:"suffix,omitempty"`
	Message            string                  `json:"message"`
	ChatHistory        []AugmentChatHistory    `json:"chat_history,omitempty"`
	Lang               string                  `json:"lang,omitempty"`
	Blobs              []AugmentBlob           `json:"blobs,omitempty"`
	ToolDefinitions    []AugmentToolDefinition `json:"tool_definitions,omitempty"`
	Nodes              []any                   `json:"nodes,omitempty"`
	Mode               string                  `json:"mode,omitempty"`                 // "CHAT" | "AGENT"
	ThirdPartyOverride string                  `json:"third_party_override,omitempty"` // Base64 encoded JSON
	ConversationId     string                  `json:"conversation_id,omitempty"`
}

// AugmentChatHistory 聊天历史记录
type AugmentChatHistory struct {
	Role    string `json:"role"`    // "user" | "assistant"
	Content string `json:"content"`
}

// AugmentBlob 代码块/文件内容
type AugmentBlob struct {
	Path    string `json:"path,omitempty"`
	Content string `json:"content,omitempty"`
}

// AugmentToolDefinition 工具定义
type AugmentToolDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters,omitempty"`
}

// AugmentThirdPartyOverride 第三方覆盖配置（解码后的结构）
type AugmentThirdPartyOverride struct {
	ProviderModelName string `json:"providerModelName,omitempty"`
	ApiKey            string `json:"apiKey,omitempty"`
	BaseUrl           string `json:"baseUrl,omitempty"`
}
