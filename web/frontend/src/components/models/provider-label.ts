import type { ModelProviderOption } from "@/api/models"

const PROVIDER_LABELS: Record<string, string> = {
  openai: "OpenAI",
  bedrock: "AWS Bedrock",
  anthropic: "Anthropic",
  "anthropic-messages": "Anthropic Messages",
  azure: "Azure OpenAI",
  gemini: "Google Gemini",
  deepseek: "DeepSeek",
  "coding-plan": "Alibaba Coding Plan",
  "coding-plan-anthropic": "Alibaba Coding Plan (Anthropic)",
  "qwen-portal": "Qwen (阿里云)",
  "qwen-intl": "Qwen International",
  "qwen-us": "Qwen US",
  moonshot: "Moonshot (月之暗面)",
  groq: "Groq",
  openrouter: "OpenRouter",
  nvidia: "NVIDIA",
  cerebras: "Cerebras",
  volcengine: "Volcengine (火山引擎)",
  shengsuanyun: "ShengsuanYun (神算云)",
  antigravity: "Google Code Assist",
  "github-copilot": "GitHub Copilot",
  "claude-cli": "Claude CLI (local)",
  "codex-cli": "Codex CLI (local)",
  ollama: "Ollama (local)",
  lmstudio: "LM Studio (local)",
  litellm: "LiteLLM",
  mistral: "Mistral AI",
  avian: "Avian",
  vllm: "VLLM (local)",
  zhipu: "Zhipu AI (智谱)",
  zai: "Z.ai",
  mimo: "Xiaomi MiMo",
  venice: "Venice AI",
  vivgrid: "Vivgrid",
  minimax: "MiniMax",
  longcat: "LongCat",
  modelscope: "ModelScope (魔搭社区)",
  novita: "Novita AI",
}

const PROVIDER_ALIASES: Record<string, string> = {
  qwen: "qwen-portal",
  "qwen-international": "qwen-intl",
  "dashscope-intl": "qwen-intl",
  "z.ai": "zai",
  "z-ai": "zai",
  google: "gemini",
  "google-antigravity": "antigravity",
}

export const PROVIDER_PRIORITY: Record<string, number> = {
  volcengine: 0,
  openai: 1,
  gemini: 2,
  anthropic: 3,
  bedrock: 4,
  "anthropic-messages": 5,
  zhipu: 6,
  deepseek: 7,
  openrouter: 8,
  "qwen-portal": 9,
  "qwen-intl": 10,
  "qwen-us": 11,
  moonshot: 12,
  groq: 13,
  "coding-plan": 14,
  "coding-plan-anthropic": 15,
  "github-copilot": 16,
  antigravity: 17,
  nvidia: 18,
  cerebras: 19,
  shengsuanyun: 20,
  venice: 21,
  vivgrid: 22,
  minimax: 23,
  longcat: 24,
  modelscope: 25,
  mistral: 26,
  avian: 27,
  novita: 28,
  azure: 29,
  litellm: 30,
  ollama: 31,
  vllm: 32,
  lmstudio: 33,
  "claude-cli": 34,
  "codex-cli": 35,
  zai: 36,
  mimo: 37,
}

export function getProviderKey(provider?: string): string {
  const normalized = provider?.trim().toLowerCase()
  if (!normalized) return "openai"
  return PROVIDER_ALIASES[normalized] ?? normalized
}

export function getProviderLabel(provider?: string): string {
  const prefix = getProviderKey(provider)
  return PROVIDER_LABELS[prefix] ?? prefix
}

export function findProviderOption(
  provider: string | undefined,
  options: ModelProviderOption[],
): ModelProviderOption | undefined {
  const providerKey = getProviderKey(provider)
  return options.find((option) => option.id === providerKey)
}

export function getProviderDefaultAPIBase(
  provider: string | undefined,
  options: ModelProviderOption[],
): string {
  return findProviderOption(provider, options)?.default_api_base ?? ""
}

export function getSortedProviderOptions(
  options: ModelProviderOption[],
): ModelProviderOption[] {
  return [...options].sort((a, b) => {
    const aPriority = PROVIDER_PRIORITY[a.id] ?? Number.MAX_SAFE_INTEGER
    const bPriority = PROVIDER_PRIORITY[b.id] ?? Number.MAX_SAFE_INTEGER
    if (aPriority !== bPriority) {
      return aPriority - bPriority
    }
    return getProviderLabel(a.id).localeCompare(getProviderLabel(b.id))
  })
}

export function getProviderDefaultAuthMethod(
  provider: string | undefined,
  options: ModelProviderOption[],
): string {
  return findProviderOption(provider, options)?.default_auth_method ?? ""
}

export function isProviderAuthMethodLocked(
  provider: string | undefined,
  options: ModelProviderOption[],
): boolean {
  return findProviderOption(provider, options)?.auth_method_locked === true
}
