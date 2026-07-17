## 2026-06-26 - Task: 复刻参考视频风格首页并新增模型列表页
### What was done
- 将默认首页改造成暗色沉浸式产品首页，保留管理员自定义 `home_content` 的原有覆盖能力。
- 新增公开模型列表页，支持按平台筛选和按模型名称搜索，并从首页提供入口。
- 增加公开页面说明文档，明确模型列表数据来源和可用性边界。
### Testing
- `pnpm build`（在 `frontend` 目录）通过。
- 本地以 `VITE_DEV_PORT=5188 pnpm dev` 启动后验证 `/home`：页面标题为 `Home - Sub2API`，主标题和 `/models` 入口可见。
- 本地验证 `/models`：页面标题为 `Models - Sub2API`，搜索框可见，模型卡片数量 172，全部标签计数 172。
### Notes
- `frontend/src/views/HomeView.vue`：替换默认首页为参考视频风格的暗色沉浸式视觉和模型列表入口。
- `frontend/src/views/public/ModelsView.vue`：新增公开模型列表页，复用前端模型白名单数据。
- `frontend/src/router/index.ts`：新增 `/models` 公开路由。
- `docs/FRONTEND_PUBLIC_PAGES.md`：新增首页和模型列表页使用说明。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue frontend/src/router/index.ts docs/FRONTEND_PUBLIC_PAGES.md`，并删除 `frontend/src/views/public/ModelsView.vue`；或回退本次提交。

## 2026-06-26 - Task: 补充公开页面文档跟踪规则
### What was done
- 为新增的公开页面说明文档增加 `docs/` 目录下的 git 跟踪例外，确保文档能随代码一起提交。
### Testing
- `pnpm build`（在 `frontend` 目录）通过。
- `git status --short --untracked-files=all` 已显示 `docs/FRONTEND_PUBLIC_PAGES.md` 为未跟踪文件，确认忽略规则例外生效。
### Notes
- `.gitignore`：新增 `!docs/FRONTEND_PUBLIC_PAGES.md` 例外规则，仅放行本次新增文档。
- 回滚方式：删除 `.gitignore` 中的 `!docs/FRONTEND_PUBLIC_PAGES.md` 行；或回退本次提交。

## 2026-06-26 - Task: 重做首页复刻方向为 scroll-film 风格
### What was done
- 放弃上一版营销卡片式布局，将默认首页重做为全屏影片式 scroll-film 结构。
- 首页改为左侧大画幅镜头、滚动章节、弹幕式字幕、视频进度条和终端浮层，更贴近参考视频的核心表现方式。
### Testing
- `pnpm build`（在 `frontend` 目录）通过。
- 本地刷新 `http://localhost:5188/home` 后验证：页面包含 `.film-player`，主标题为 `你相信这是 Codex 完成的吗？`，滚动章节数量为 3。
### Notes
- `frontend/src/views/HomeView.vue`：重写默认首页视觉结构和样式，保留自定义 `home_content` 覆盖能力。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue`；或回退本次提交。

## 2026-06-26 - Task: 首页改为 Antigravity 官方首页风格
### What was done
- 将首页默认视觉从暗色 scroll-film 改为 Google Antigravity 官方首页方向。
- 实现白底极简导航、彩色 A 标识、超大居中标题、Download 胶囊按钮、散落工具图标和浅色产品预览卡片。
- 保留 `/models` 模型列表入口和管理员自定义 `home_content` 覆盖能力。
### Testing
- `pnpm build`（在 `frontend` 目录）通过。
- 本地刷新 `http://localhost:5188/home` 后验证：页面包含 `.ag-header` 和 `.ag-preview-card`，主标题为 `Experience liftoff with the next-gen AI gateway platform`。
### Notes
- `frontend/src/views/HomeView.vue`：重写默认首页为 Antigravity 风格布局与样式。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue`；或回退本次提交。

## 2026-06-26 - Task: 补齐 Antigravity 首页 liftoff 动画
### What was done
- 为 Antigravity 风格首页补充首屏入场动画、工具图标从中心散开、工具图标漂浮、预览窗口悬浮和背景轨道旋转。
- 删除未使用的 `siteSubtitle` computed，修复 Vite checker 的 `[vue-tsc] declared but never read` 遮罩报错。
- 重启本地 5188 开发服务，清除旧错误 overlay。
### Testing
- `pnpm build`（在 `frontend` 目录）通过。
- 本地刷新 `http://localhost:5188/home` 后验证：无 `[vue-tsc]` 遮罩，`.ag-liftoff` 存在，`.ag-float-window` 存在，工具图标数量为 8。
### Notes
- `frontend/src/views/HomeView.vue`：新增 Antigravity liftoff 动效并删除未使用变量。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue`；或回退本次提交。

## 2026-06-26 - Task: 修复 Token 激励本周消耗实时展示
### What was done
- 修复 Token 激励计划状态接口在用户已领取后继续显示领取时 Token 快照的问题，改为始终返回本周实时累计消耗。
- 保留已领取状态、领取时间和实际领取金额，避免影响每周只能领取一次的业务规则。
- 补充模块说明文档，明确状态接口和领取接口的统计口径。
### Testing
- `go test -tags unit ./internal/service -run TokenIncentive`（在 `backend` 目录）通过。
- `go test -tags unit ./internal/repository -run TokenIncentive`（在 `backend` 目录）通过。
### Notes
- `backend/internal/service/token_incentive_service.go`：状态构建改为使用实时周累计 Token，而不是已领取记录中的快照 Token。
- `backend/internal/service/token_incentive_service_test.go`：更新已领取场景测试，覆盖“领取金额保留、Token 进度实时更新”。
- `docs/TOKEN_INCENTIVE.md`：新增 Token 激励计划状态展示和领取口径说明。
- `progress.md`：追加本轮修复记录。
- 回滚方式：执行 `git checkout -- backend/internal/service/token_incentive_service.go backend/internal/service/token_incentive_service_test.go docs/TOKEN_INCENTIVE.md progress.md`，或回退包含本轮改动的提交。

## 2026-06-26 - Task: 修正 Antigravity 风格首页语言切换、下载语义和 liftoff 动画
### What was done
- 恢复首页右上角语言切换入口，并用浅色胶囊样式让它在 Antigravity 风格导航中明确可见。
- 去掉首页下载/介绍播放类动作语义，将 CTA 统一为站内“立即开始/控制台”和模型列表入口。
- 将首屏工具图标云改为从中心核心向外升空散开的 liftoff 动画，并增加中心核心、扩散光环、图标浮动和预览窗口悬浮效果。
- 同步更新公开页面文档，说明默认首页当前参考 Antigravity 样式且不提供下载功能。
### Testing
- `pnpm build`（在 `frontend` 目录）通过；输出仅有既有 Browserslist 数据过期、动态/静态导入分包和 chunk size 警告。
- 本地刷新 `http://localhost:5188/home` 后通过浏览器 DOM/样式验证：无 Vite/`[vue-tsc]` 遮罩，`.ag-language` 可见且显示 `🇨🇳ZH`，页面正文无 `Download/下载` 文案，`.ag-hero-orbit` 与 `.ag-liftoff-core` 存在，`.ag-tool` 数量为 8，工具图标应用 `ag-tool-liftoff` 与 `ag-tool-breathe` 动画。
### Notes
- `frontend/src/views/HomeView.vue`：恢复语言切换视觉入口，移除首页下载/播放动作，并重做工具图标云 liftoff 动画。
- `docs/FRONTEND_PUBLIC_PAGES.md`：更新首页说明为 Antigravity 风格、语言切换、站内入口和无下载功能语义。
- `progress.md`：追加本轮修正与验证记录。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue docs/FRONTEND_PUBLIC_PAGES.md`，并从 `progress.md` 末尾删除本轮追加段落；或回退包含本轮改动的提交。

## 2026-06-26 - Task: 补齐 Antigravity 首页中文文案
### What was done
- 将 Antigravity 风格首页的导航、主标题、模型入口、预览窗口、工具图标、能力卡片和主题按钮改为跟随当前语言切换。
- 中文环境下改用中文业务文案，保留英文环境的原英文展示，避免切换到英文后丢失参考样式语义。
- 同步补充公开页面文档，说明默认首页文案会跟随语言切换。

### Testing
- `pnpm build`（在 `frontend` 目录）通过；输出仅有既有 Browserslist、动态/静态导入分包和 chunk size 警告。
- 本地刷新 `http://localhost:5188/home` 后通过浏览器 DOM 验证：无 Vite/`[vue-tsc]` 遮罩；`.ag-language` 存在；`.ag-tool` 数量为 8；中文态导航为“产品/使用场景/模型”，主标题为“下一代智能网关平台，让模型调用像反重力一样升空”。
- 中文态检查 `Products`、`Use Cases`、`Explore models`、`Experience liftoff`、`agent workspace`、`Gateway routing`、`Model mapping`、`Usage billing`、`Designed for`、`Developer first`、`Model catalog`、`Light`、`Dark`、`Antigravity` 均未在可见正文中残留；仅保留 `Sub2API`、Claude/GPT/Gemini 和模型名等专有名词。

### Notes
- `frontend/src/views/HomeView.vue`：新增语言感知的首页文案对象，并将首页硬编码英文替换为中英文切换文案。
- `docs/FRONTEND_PUBLIC_PAGES.md`：补充默认首页文案跟随语言切换的说明。
- `progress.md`：追加本轮中文化修正与验证记录。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue`，并从 `docs/FRONTEND_PUBLIC_PAGES.md` 删除本轮新增的语言切换说明行，再从 `progress.md` 末尾删除本段记录；或回退包含本轮改动的提交。

## 2026-06-26 - Task: 修复 Antigravity 首页中文语言选项不可点击
### What was done
- 定位到首页移动布局中语言下拉菜单已打开，但被首屏 hero 区域的层级压住，导致中文选项点击点实际命中 hero。
- 只调整首页语言入口所在的 header/actions/language 层级，并在窄屏下把下拉菜单从左侧展开，避免菜单越出视口和被下方内容拦截。

### Testing
- 本地浏览器刷新 `http://localhost:5188/home` 后验证：语言菜单打开时 `🇺🇸English` 与 `🇨🇳中文` 两个选项中心点不再命中 `.ag-hero`。
- 本地浏览器验证英文切中文链路：先切到英文，随后点击 `🇨🇳中文`，页面恢复 `html lang="zh"`，首页标题恢复为“下一代智能网关平台，让模型调用像反重力一样升空”。
- `pnpm build`（在 `frontend` 目录）通过；输出仅保留既有 Browserslist、动态/静态导入分包和 chunk size 警告。

### Notes
- `frontend/src/views/HomeView.vue`：提高 Antigravity 首页 header 与语言切换区域层级，并修正移动端语言下拉展开方向。
- `progress.md`：追加本轮语言切换修复与验证记录。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue`，并从 `progress.md` 末尾删除本轮追加段落；或回退包含本轮改动的提交。

## 2026-06-26 - Task: 调整 Antigravity 首页为 Sub2API 项目文案
### What was done
- 将首页主标题、导航、预览窗口、工具图标和能力卡片文案从泛 Antigravity 科技话术调整为 Sub2API 项目定位。
- 中文文案聚焦统一 API 入口、订阅账号聚合、账号池路由、模型别名映射、令牌用量统计和余额计费。
- 英文文案同步调整为同等业务含义，保留当前 Antigravity 风格视觉、动画和页面结构不变。
- 同步更新公开页面说明文档，明确首页能力卡片说明的是 Sub2API 真实产品能力。

### Testing
- `pnpm build`（在 `frontend` 目录）通过；输出仅保留既有 Browserslist、动态/静态导入分包和 chunk size 警告。
- 本地浏览器刷新 `http://localhost:5188/home` 后验证中文首页可见文案：品牌为“Sub2API 统一模型网关”，主标题为“把多平台订阅账号接成统一、稳定、可计费的 API 网关”，能力卡片覆盖统一入口、账号池路由、模型映射、用量与余额。
- 本地浏览器验证预览窗口文案已显示 OpenAI 兼容请求统一接入、按模型映射到 Claude/GPT/Gemini 等上游、按令牌统计用量并扣减余额。

### Notes
- `frontend/src/views/HomeView.vue`：替换首页中英文业务文案，使其贴合 Sub2API 项目能力。
- `docs/FRONTEND_PUBLIC_PAGES.md`：补充首页能力卡片对应的 Sub2API 真实能力说明。
- `progress.md`：追加本轮文案调整与验证记录。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue docs/FRONTEND_PUBLIC_PAGES.md`，并从 `progress.md` 末尾删除本轮追加段落；或回退包含本轮改动的提交。

## 2026-06-26 - Task: 按原首页文案体系重写 Antigravity 首页文案
### What was done
- 参考原首页 i18n 文案和 README 项目说明，将当前 Antigravity 风格首页的业务文案调整回原项目表达体系。
- 主标题恢复为“一个密钥，畅用多个 AI 模型”，副文案恢复为“无需管理多个订阅账号，一站式接入 Claude、GPT、Gemini 等主流 AI 服务”。
- 能力卡片围绕原首页的“一键接入、稳定可靠、用多少付多少、一个 API 多种模型选择”展开，并保留当前视觉样式和动画。
- 同步更新公开页面说明文档，说明当前首页沿用原首页文案体系。

### Testing
- `pnpm build`（在 `frontend` 目录）通过；输出仅保留既有 Browserslist、动态/静态导入分包和 chunk size 警告。
- 本地浏览器刷新 `http://localhost:5188/home` 后验证：品牌为“Sub2API AI API 网关”，主标题为“一个密钥，畅用多个 AI 模型”。
- 本地浏览器验证可见能力卡片包含“一键接入”“稳定可靠”“用多少付多少”“支持模型”，工具图标标签包含“订阅转 API”“会话保持”“按量计费”。

### Notes
- `frontend/src/views/HomeView.vue`：基于原首页文案和 README 项目定位重写当前首页中英文文案。
- `docs/FRONTEND_PUBLIC_PAGES.md`：更新首页说明为沿用原首页文案体系。
- `progress.md`：追加本轮文案修正与验证记录。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue docs/FRONTEND_PUBLIC_PAGES.md`，并从 `progress.md` 末尾删除本轮追加段落；或回退包含本轮改动的提交。

## 2026-06-26 - Task: 移除 Antigravity 首页副文案
### What was done
- 移除了首页主标题下方的副文案展示，不再显示“无需管理多个订阅账号，一站式接入 Claude、GPT、Gemini 等主流 AI 服务”。
- 清理了该副文案对应的首页样式和动画引用，避免留下无用展示结构。
- 同步更新公开页面说明，明确当前首页不再展示额外副文案说明。
### Testing
- `rg -n "无需管理多个订阅账号|No need to manage multiple subscriptions|一站式接入 Claude|ag-intro|intro:" frontend/src/views/HomeView.vue docs/FRONTEND_PUBLIC_PAGES.md` 未命中，确认当前源码与公开页面文档不再包含目标副文案和旧区块。
- `pnpm build`（在 `frontend` 目录）通过；仅保留既有 Browserslist、动态/静态导入分包和 chunk size 警告。
- 本地浏览器刷新 `http://localhost:5188/home` 后验证：页面可见正文不含目标副文案，`.ag-intro` 不存在，主标题仍为“一个密钥，畅用多个 AI 模型”。
### Notes
- `frontend/src/views/HomeView.vue`：移除首页副文案区块、对应中英文文案字段及无用样式/动画引用。
- `docs/FRONTEND_PUBLIC_PAGES.md`：更新首页说明，补充“首页不再展示额外副文案说明”的当前状态。
- `progress.md`：追加本轮移除副文案的施工与验证记录。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue docs/FRONTEND_PUBLIC_PAGES.md`，并从 `progress.md` 末尾删除本轮追加段落；或回退包含本轮改动的提交。

## 2026-06-26 - Task: 移除首页多账号池自动切换文案
### What was done
- 移除了首页能力卡片中的“多账号池自动切换，减少调用中断”及其说明文案。
- 将首页预览窗口和工具图标中同义的“多个上游账号自动切换/负载均衡”表达改为模型映射口径，避免页面其它位置继续出现同类表述。
- 同步更新公开页面说明，记录首页不再展示账号调度说明。
### Testing
- `rg`/脚本检查 `frontend/src/views/HomeView.vue` 和 `docs/FRONTEND_PUBLIC_PAGES.md`，确认不再包含“多账号池自动切换”“智能调度多个上游账号”“多个上游账号”“负载均衡”“调用中断”“单一账号限流”等目标文案。
- `pnpm build`（在 `frontend` 目录）通过；仅保留既有 Browserslist、动态/静态导入分包和 chunk size 警告。
- 本地浏览器刷新 `http://localhost:5188/home` 后验证：页面可见正文不含目标账号池/自动切换/负载均衡文案，主标题仍为“一个密钥，畅用多个 AI 模型”。
### Notes
- `frontend/src/views/HomeView.vue`：删除多账号池能力卡片，并将预览区/工具标签中的同义账号调度文案改为模型映射文案。
- `docs/FRONTEND_PUBLIC_PAGES.md`：更新首页说明，移除“稳定可靠”能力卡片口径并说明不再展示账号调度说明。
- `progress.md`：追加本轮文案移除与验证记录。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue docs/FRONTEND_PUBLIC_PAGES.md`，并从 `progress.md` 末尾删除本轮追加段落；或回退包含本轮改动的提交。

## 2026-06-26 - Task: 移除首页模型目录能力卡片文案
### What was done
- 移除了首页能力卡片中的“公开模型目录展示 Claude、GPT、Gemini、Antigravity 等已接入模型，方便接入前快速确认覆盖范围。”文案。
- 同步删除对应的“支持模型 / 一个 API，多种模型选择”能力卡片及英文同义卡片，仅保留顶部导航和 CTA 中的模型列表入口。
- 更新公开页面说明，记录首页能力卡片不再展示模型目录说明。
### Testing
- 脚本检查 `frontend/src/views/HomeView.vue` 和 `docs/FRONTEND_PUBLIC_PAGES.md`，确认不再包含“公开模型目录展示 Claude”“方便接入前快速确认覆盖范围”“一个 API，多种模型选择”等目标文案及英文同义文案。
- `pnpm build`（在 `frontend` 目录）通过；仅保留既有 Browserslist、动态/静态导入分包和 chunk size 警告。
- 本地浏览器刷新 `http://localhost:5188/home` 后验证：页面可见正文不含目标模型目录能力卡片文案，能力卡片数量为 2，主标题仍为“一个密钥，畅用多个 AI 模型”。
### Notes
- `frontend/src/views/HomeView.vue`：删除模型目录能力卡片及英文同义卡片。
- `docs/FRONTEND_PUBLIC_PAGES.md`：更新首页能力卡片说明，移除模型目录能力说明。
- `progress.md`：追加本轮文案移除与验证记录。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue docs/FRONTEND_PUBLIC_PAGES.md`，并从 `progress.md` 末尾删除本轮追加段落；或回退包含本轮改动的提交。

## 2026-06-26 - Task: 移除首页上游 AI 服务点名文案
### What was done
- 将首页“一键接入”能力卡片描述从点名 Claude、GPT、Gemini 等上游服务，改为不点名服务商的统一 API 密钥接入口径。
- 同步替换英文同义描述，避免英文语言下继续出现 upstream AI services 表述。
### Testing
- 脚本检查 `frontend/src/views/HomeView.vue` 和 `docs/FRONTEND_PUBLIC_PAGES.md`，确认不再包含“Claude、GPT、Gemini 等上游 AI 服务”“上游 AI 服务”“upstream AI services”等目标文案。
- `pnpm build`（在 `frontend` 目录）通过；仅保留既有 Browserslist、动态/静态导入分包和 chunk size 警告。
- 本地浏览器刷新 `http://localhost:5188/home` 后验证：页面可见正文不含目标上游服务点名文案，新的“一键接入”描述显示为统一 API 密钥接入平台已配置模型。
### Notes
- `frontend/src/views/HomeView.vue`：替换首页“一键接入”能力卡片的中英文描述。
- `progress.md`：追加本轮文案移除与验证记录。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue`，并从 `progress.md` 末尾删除本轮追加段落；或回退包含本轮改动的提交。


## 2026-06-26 - Task: 修复首页深色/浅色主题切换无效
### What was done
- 将默认首页主题切换从只依赖全局 `html.dark` 样式覆盖，改为在首页根节点同步添加 `ag-home-dark` 状态类，确保 scoped CSS 能命中当前首页结构。
- 补齐深色状态下页面背景、主标题、预览窗口、能力卡片、工具图标、按钮、边框和辅助文字的实际配色覆盖。
- 更新公开页面说明，记录默认首页支持浅色/深色主题切换。

### Testing
- `pnpm build`（在 `frontend` 目录）通过；仅保留既有 Browserslist、动态/静态导入分包和 chunk size 警告。
- 本地浏览器刷新 `http://localhost:5188/home` 后验证：深色状态下 `.ag-home` 背景为 `rgb(7, 10, 18)`，主标题为 `rgb(248, 250, 252)`，预览卡片为 `rgba(15, 23, 42, 0.72)`，按钮显示“浅色”。
- 点击“浅色”后验证：`.ag-home` 背景恢复为 `rgb(250, 249, 246)`，主标题恢复为 `rgb(17, 17, 20)`，预览卡片恢复为 `rgba(255, 255, 255, 0.72)`，按钮显示“深色”；再次点击“深色”可恢复深色状态。

### Notes
- `frontend/src/views/HomeView.vue`：为首页根节点增加 `ag-home-dark` 状态类，并将深色覆盖改为当前组件内可命中的 scoped 选择器。
- `docs/FRONTEND_PUBLIC_PAGES.md`：补充默认首页浅色/深色主题切换说明。
- `progress.md`：追加本轮修复与验证记录。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue docs/FRONTEND_PUBLIC_PAGES.md`，并从 `progress.md` 末尾删除本轮追加段落；或回退包含本轮改动的提交。


## 2026-06-26 - Task: 修复首页中英文切换交互
### What was done
- 为通用 `LocaleSwitcher` 增加弹出方向参数，保留默认向下弹出能力。
- 将首页底部语言切换设置为向上弹出，避免页面底部菜单贴边或被视口裁切导致看起来无法选择语言。
- 同步更新公开页面说明，明确首页顶部和底部语言切换均支持中英文切换。

### Testing
- `pnpm build`（在 `frontend` 目录）通过；仅保留既有 Browserslist、动态/静态导入分包和 chunk size 警告。
- 本地浏览器刷新 `http://localhost:5188/home` 后验证顶部语言切换：点击 English 后 `lang=en`，主标题变为 `One Key, All AI Models`，导航变为 `Solutions / Features / Models`。
- 滚动到底部后验证底部语言切换：菜单向上弹出且在视口内，点击中文后 `lang=zh`，主标题恢复为 `一个密钥，畅用多个 AI 模型`，导航恢复为 `解决方案 / 核心功能 / 支持模型`。

### Notes
- `frontend/src/components/common/LocaleSwitcher.vue`：新增 `placement` 参数控制菜单向上或向下弹出。
- `frontend/src/views/HomeView.vue`：将首页底部语言切换改为 `placement="top"`。
- `docs/FRONTEND_PUBLIC_PAGES.md`：补充顶部和底部语言切换均可用的说明。
- `progress.md`：追加本轮修复与验证记录。
- 回滚方式：执行 `git checkout -- frontend/src/components/common/LocaleSwitcher.vue frontend/src/views/HomeView.vue docs/FRONTEND_PUBLIC_PAGES.md`，并从 `progress.md` 末尾删除本轮追加段落；或回退包含本轮改动的提交。

## 2026-06-26 - Task: 记录后台站点设置动态信息生效问题
### What was done
- 排查后台“站点设置”动态字段在当前首页的生效范围，确认默认首页目前只读取/展示 `site_name`、`doc_url`、`home_content`。
- 记录未完整接入字段：`site_logo`、`site_subtitle`、`api_base_url`、`contact_info` 尚未在新首页完整展示。
- 记录本地公开设置链路异常：`/api/v1/settings/public` 在 5188 和 8080 均返回 404，当前运行态首页未注入 `window.__APP_CONFIG__`，因此主要显示默认配置。
### Testing
- `rg` 检查 `frontend/src/views/HomeView.vue`、`frontend/src/stores/app.ts`、`frontend/src/types/index.ts` 中 public settings 字段引用。
- `Invoke-RestMethod` 验证 `http://localhost:5188/api/v1/settings/public` 和 `http://localhost:8080/api/v1/settings/public` 均返回 404。
- in-app browser 读取 `http://localhost:5188/home` 运行态：`window.__APP_CONFIG__` 为 `null`，页面展示默认 `Sub2API` 文案且无 logo 图片。
### Notes
- `progress.md`：追加当前后台站点设置动态信息未完整生效的问题记录和验证证据。
- 回滚方式：从 `progress.md` 末尾删除本条 `2026-06-26 - Task: 记录后台站点设置动态信息生效问题` 记录；或回退包含本次日志追加的提交。

## 2026-06-26 - Task: 调整原仓库更新检测为 Release-only
### What was done
- 将原仓库更新检测从分支 commit 比较改为仅检测原仓库最新 Release 版本号，避免原仓库普通提交触发“有新版本”提示。
- 保留当前二开仓库 Release 作为应用内一键更新的唯一安装来源；原仓库 Release 高于当前仓库时只提示需要先同步、合并并发布当前仓库 Release。
- 对旧缓存中的原仓库 commit 更新字段做兼容归一化，避免旧缓存继续显示上游 commit 更新提示。
- 补充更新检测策略文档，明确当前仓库优先和原仓库 Release-only 的边界。
### Testing
- `go test -tags unit ./internal/service -run UpdateService`（在 `backend` 目录）通过。
### Notes
- `backend/internal/service/update_service.go`：原仓库检测只调用最新 Release，不再拉取或比较原仓库分支提交；缓存读取时归一化原仓库更新字段。
- `backend/internal/service/update_service_test.go`：更新测试，覆盖“原仓库只有新 commit 不提示”和“原仓库有新 Release 才提示同步”。
- `docs/UPDATE_POLICY.md`：新增当前仓库与原仓库的更新检测/安装策略说明。
- `.gitignore`：放行 `docs/UPDATE_POLICY.md`，确保策略文档可提交。
- `progress.md`：追加本轮修复记录。
- 回滚方式：执行 `git checkout -- backend/internal/service/update_service.go backend/internal/service/update_service_test.go .gitignore docs/UPDATE_POLICY.md progress.md`，或回退包含本轮改动的提交。

## 2026-06-28 - Task: merge upstream latest code into fork
### What was done
- Merged upstream main into the fork while preserving the fork's update target and release-based update prompt behavior.
- Kept the current repository as the actual in-app update target and left the upstream repository only as a release signal source.

### Testing
- `go test -tags unit ./internal/service -run UpdateService` (in `D:\project\sub2api-so\backend`)

### Notes
- Files changed by merge include the backend/frontend updates from upstream, with the fork-specific update logic retained in `backend/internal/service/update_service.go`.
- Rollback point: `git reset --hard 14b62588` to return to the pre-merge fork state; the pre-sync local work is preserved in `stash@{0}: On main: codex pre-upstream-sync`.


## 2026-06-28 - Task: token incentive per-tier full reward claims
### What was done
- Changed token incentive rewards from weekly single-claim behavior to per-tier weekly claims.
- Each reached tier can now be claimed once for its full configured amount, so 50M/100M/500M default tiers pay 2 + 5 + 10 instead of a differential amount.
- Added tier-level claim tracking and updated the user progress UI to keep later reached tiers claimable after earlier tiers have been claimed.

### Testing
- `go test -tags unit ./internal/service ./internal/repository -run TokenIncentive` (in `D:\project\sub2api-so\backend`)
- `pnpm exec vitest run src/views/user/__tests__/UsageView.spec.ts --reporter=verbose` (in `D:\project\sub2api-so\frontend`)
- `pnpm typecheck` (in `D:\project\sub2api-so\frontend`)
- `git diff --check` (in `D:\project\sub2api-so`)

### Notes
- `backend/internal/service/token_incentive_service.go`: selects the first reached unclaimed tier and reports claimed tiers plus total claimed reward.
- `backend/internal/service/token_incentive_service_test.go`: covers second-tier and third-tier full configured reward claims.
- `backend/internal/repository/token_incentive_repo.go`: stores and queries claims by `threshold_tokens`, and credits the configured tier amount.
- `backend/internal/repository/token_incentive_repo_test.go`: verifies tier claim persistence, duplicate detection, and redeem history content.
- `backend/migrations/157_token_incentive_tier_claims.sql`: adds `threshold_tokens` and changes uniqueness to one row per user/week/tier.
- `frontend/src/types/index.ts`: exposes tier claim status fields to the frontend.
- `frontend/src/views/user/UsageView.vue`: lets claimed and still-claimable tiers coexist in the progress card.
- `docs/TOKEN_INCENTIVE.md`: documents per-tier full reward behavior and same-week claim requirement.
- `progress.md`: appended this implementation record.
- Rollback方式：执行 `git checkout -- backend/internal/service/token_incentive_service.go backend/internal/service/token_incentive_service_test.go backend/internal/repository/token_incentive_repo.go backend/internal/repository/token_incentive_repo_test.go frontend/src/types/index.ts frontend/src/views/user/UsageView.vue docs/TOKEN_INCENTIVE.md progress.md` 并删除 `backend/migrations/157_token_incentive_tier_claims.sql`，或回退包含本轮改动的提交。


## 2026-06-28 - Task: 接入后台站点设置动态信息到默认首页
### What was done
- 修复默认首页只展示部分站点设置的问题，在未配置 `home_content` 时同步展示后台配置的站点 Logo、站点副标题、API 地址、联系方式和文档链接。
- 将首页预览窗口中的固定 `Sub2API` 品牌改为动态站点名称，并在有配置时补充站点副标题、API 地址和联系方式。
- 保留自定义首页内容 `home_content` 的最高优先级，管理员配置自定义 HTML 或 URL 时仍直接覆盖默认首页。

### Testing
- `rg -n "siteLogo|siteSubtitle|apiBaseUrl|contactInfo|previewApiBase|previewContact|ag-logo|ag-site-subtitle|ag-footer-meta" frontend/src/views/HomeView.vue docs/FRONTEND_PUBLIC_PAGES.md`：确认新增字段读取、模板展示和样式入口均存在。
- `git diff --check -- frontend/src/views/HomeView.vue docs/FRONTEND_PUBLIC_PAGES.md`：通过，仅保留既有换行符提示。
- `D:\environment\nodejs\node-v22.17.0-win-x64\pnpm.cmd build`（在 `frontend` 目录）：通过；仅保留既有 Browserslist、动态/静态导入分包和 chunk size 警告。
- `Invoke-RestMethod http://localhost:5188/api/v1/settings/public` 与 `Invoke-RestMethod http://localhost:8080/api/v1/settings/public`：当前本地监听进程仍返回 404；源码路由已存在，说明当前运行态后端不是这份已注册公开设置路由的服务或未按最新代码启动，运行态字段回填需重启正确后端后再验证。

### Notes
- `frontend/src/views/HomeView.vue`：接入后台公开站点设置字段，补充 Logo、副标题、API 地址、联系方式展示，并对 Logo/文档链接做 URL 规范化。
- `docs/FRONTEND_PUBLIC_PAGES.md`：补充默认首页会读取后台站点设置动态字段的说明。
- `progress.md`：追加本轮修复、验证和当前运行态接口 404 的记录。
- 回滚方式：执行 `git checkout -- frontend/src/views/HomeView.vue docs/FRONTEND_PUBLIC_PAGES.md progress.md` 可回退本轮首页动态字段修复与记录；如只回退日志，从 `progress.md` 末尾删除本条 `2026-06-28 - Task: 接入后台站点设置动态信息到默认首页` 段落。

## 2026-06-28 - Task: audit token incentive tier claim logic
### What was done
- Audited the token incentive tier-claim path for duplicate-claim, legacy-data migration, amount mapping, and frontend state consistency risks.
- Fixed legacy claim threshold resolution so old records without `threshold_tokens` map to the highest matching claimed tier, preventing accidental duplicate claims after upgrade.
- Hardened the tier-claim migration to backfill old records from configured/default rules and to remain idempotent.
- Kept the business rule unchanged: each reached target can claim that target's full configured amount once per week.

### Testing
- `go test -tags unit ./internal/service ./internal/repository -run TokenIncentive` (in `D:\project\sub2api-soackend`)
- `pnpm exec vitest run src/views/user/__tests__/UsageView.spec.ts --reporter=verbose` (in `D:\project\sub2api-sorontend`)
- `pnpm typecheck` (in `D:\project\sub2api-sorontend`)
- `git diff --check` (in `D:\project\sub2api-so`)

### Notes
- `backend/internal/service/token_incentive_service.go`: preserves `eligible` as reached-target status, uses `claimable` for current claim availability, and resolves legacy claims to the highest matching reward tier.
- `backend/internal/service/token_incentive_service_test.go`: adds regression coverage for legacy same-reward tier mapping and full per-tier rewards.
- `backend/migrations/157_token_incentive_tier_claims.sql`: backfills old claims using configured/default tiers and adds an idempotent positive-threshold constraint.
- `progress.md`: appended this audit record.
- Rollback方式：执行 `git checkout -- backend/internal/service/token_incentive_service.go backend/internal/service/token_incentive_service_test.go backend/migrations/157_token_incentive_tier_claims.sql progress.md`，或回退包含本轮审计修复的提交。

## 2026-06-28 - Task: prepare v0.1.139-fy.1 release metadata
### What was done
- Synced the fork source version metadata to upstream base `0.1.139` after the upstream merge.
- Prepared the current fork changes for release tag `v0.1.139-fy.1`.

### Testing
- `git diff --check` (in `D:\project\sub2api-so`)

### Notes
- `backend/cmd/server/VERSION`: updated the source base version to `0.1.139`.
- `backend/cmd/server/UPSTREAM_COMMIT`: updated the recorded upstream commit to `c275422251e72750bebe53e41fcf59db7f83fe6b`.
- `progress.md`: appended this release metadata record.
- Rollback方式：执行 `git checkout -- backend/cmd/server/VERSION backend/cmd/server/UPSTREAM_COMMIT progress.md`，或回退包含本轮发版元数据的提交。


## 2026-06-29 - Task: 修复 Token 激励计划领取失败
### What was done
- 修复 Token 激励计划领取事务：档位领取记录和用户余额入账提交成功后即视为领取成功，余额变动记录补写失败不再回滚奖励。
- 增加数据库兼容迁移，清理旧版一周只能领取一次的 `user_id + week_start` 唯一约束/索引残留，确保多档位可按 `threshold_tokens` 分别领取。
- 补强余额变动记录写入的幂等性，避免重复补写记录时因为兑换码冲突影响后续流程。

### Testing
- `go test -tags unit ./internal/service ./internal/repository -run TokenIncentive`（在 `D:\project\sub2api-so\backend`）通过。
- `go test -tags unit ./internal/repository -run 'ApplyMigrations|Migration'`（在 `D:\project\sub2api-so\backend`）通过。
- `git diff --check`（在 `D:\project\sub2api-so`）通过。

### Notes
- `backend/internal/repository/token_incentive_repo.go`：将余额变动记录写入移到奖励主事务提交后执行，并改为失败只记录日志。
- `backend/internal/repository/token_incentive_repo_test.go`：新增余额变动记录写入失败不影响领取成功的回归测试。
- `backend/migrations/158_fix_token_incentive_tier_constraints.sql`：新增兼容迁移，移除旧单周唯一约束/索引，补齐档位唯一索引和 `redeem_codes.notes` 字段。
- `docs/TOKEN_INCENTIVE.md`：补充说明奖励入账与余额变动记录的关系。
- `progress.md`：追加本轮修复、验证和回滚记录。
- 回滚方式：执行 `git checkout -- backend/internal/repository/token_incentive_repo.go backend/internal/repository/token_incentive_repo_test.go docs/TOKEN_INCENTIVE.md progress.md` 并删除 `backend/migrations/158_fix_token_incentive_tier_constraints.sql`；或回退包含本轮修复的提交。

## 2026-06-29 - Task: 修复 Token 激励领取 SQL 类型推断失败
### What was done
- 修复 Token 激励领取入库 SQL 中 PostgreSQL 对同一参数 `$5` 推断类型不一致导致的领取 500 问题。
- 将领取阈值参数在插入和资格复核条件中显式转换为 `bigint`，确保线上错误 `pq: inconsistent types deduced for parameter $5` 不再触发。
- 补充单元测试断言，防止后续移除 `$5::bigint` 类型约束导致回归。

### Testing
- `go test -tags unit ./internal/repository -run TokenIncentive`（在 `D:\project\sub2api-so\backend`）通过。
- `go test -tags unit ./internal/service ./internal/repository -run TokenIncentive`（在 `D:\project\sub2api-so\backend`）通过。
- 本地无 Docker/psql 运行环境，未执行真实 PostgreSQL 集成测试；已根据生产日志中的 `pq: inconsistent types deduced for parameter $5` 对对应 SQL 参数做显式类型修复。

### Notes
- `backend/internal/repository/token_incentive_repo.go`：将本周 token 汇总结果和领取阈值参数显式固定为 `bigint`，修复 PostgreSQL 参数类型推断冲突。
- `backend/internal/repository/token_incentive_repo_test.go`：新增对 `$5::bigint` 和数据库端阈值复核条件的回归断言。
- `progress.md`：追加本轮线上领取失败修复、验证和回滚记录。
- 回滚方式：执行 `git checkout -- backend/internal/repository/token_incentive_repo.go backend/internal/repository/token_incentive_repo_test.go progress.md`，或回退包含本轮修复的提交。

## 2026-07-01 - Task: 合并原仓库 0.1.141 并保留二开功能
### What was done
- 合并原仓库 `Wei-Shaw/sub2api` 的最新 `main`（上游基线 `0.1.141` / `c1335bae12c30f6d73f1a86051eb5ba312ba23c9`）。
- 解决版本文件、安装接口、用户用量页和用量页测试冲突；保留 Token 激励计划、二开更新逻辑和本仓库更新目标。
- 用户用量页以上游新版统计/图表结构为底，重新嵌入 Token 激励进度条、档位状态和领取动作，避免丢失上游新增能力。

### Testing
- `go test -tags unit ./internal/service ./internal/repository -run TokenIncentive`（在 `D:\project\sub2api-so\backend`）通过。
- `go test -tags unit ./internal/server ./internal/handler ./internal/service ./internal/repository -run "TokenIncentive|Update|Version|Usage"`（在 `D:\project\sub2api-so\backend`）通过。
- `go test -tags unit ./internal/repository -run "ApplyMigrations|Migration"`（在 `D:\project\sub2api-so\backend`）通过。
- `cmd /c pnpm.cmd exec vitest run src/views/user/__tests__/UsageView.spec.ts --reporter=verbose`（在 `D:\project\sub2api-so\frontend`）通过。
- `cmd /c pnpm.cmd typecheck`（在 `D:\project\sub2api-so\frontend`）通过。

### Notes
- `backend/cmd/server/VERSION`：更新上游基线版本为 `0.1.141`。
- `backend/cmd/server/UPSTREAM_COMMIT`：记录本次合并对应的上游提交 `c1335bae12c30f6d73f1a86051eb5ba312ba23c9`。
- `frontend/src/api/setup.ts`：保留上游 gateway URL 构造，同时保留二开安装流程长超时。
- `frontend/src/views/user/UsageView.vue`：在上游新版用量页结构中恢复 Token 激励状态、进度条、档位展示和领取逻辑。
- `frontend/src/views/user/__tests__/UsageView.spec.ts`：补齐上游新版用量页测试与 Token 激励 API mock/领取回归覆盖。
- 上游合并文件：其余 README、后端服务、迁移、前端组件/API 等变更来自 `upstream/main` 本次合并，未额外扩大人工改动范围。
- 回滚方式：提交后执行 `git revert -m 1 <本次合并提交>`；若只回滚人工冲突处理，可执行 `git checkout HEAD^ -- backend/cmd/server/VERSION backend/cmd/server/UPSTREAM_COMMIT frontend/src/api/setup.ts frontend/src/views/user/UsageView.vue frontend/src/views/user/__tests__/UsageView.spec.ts progress.md`。
## 2026-07-02 - Task: 合并原仓库 0.1.142 并准备发版
### What was done
- 合并原仓库 `Wei-Shaw/sub2api` 最新 `main`，上游基线为 `0.1.142` / `7dc7cfce1db5d31599815ff29acf6847ead0f0b7`。
- 解决版本文件冲突，将本仓库基础版本同步到 `0.1.142`，并更新记录的上游提交。
- 保留当前仓库二开逻辑，通过普通 merge 引入上游新增的 Grok 媒体路由、Spark shadow 账号、订阅撤销修复、平台额度等更新。

### Testing
- `git diff --check`（在 `D:\project\sub2api-so`）通过。
- `go test ./...`（在 `D:\project\sub2api-so\backend`）通过。
- `D:\environment\nodejs\node-v22.17.0-win-x64\node_modules\@pnpm\exe\pnpm.exe run build`（在 `D:\project\sub2api-so\frontend`）通过；仅保留既有 Browserslist、动态/静态导入混用和 chunk size 警告。

### Notes
- `backend/cmd/server/VERSION`：解决合并冲突并同步上游基础版本到 `0.1.142`。
- `backend/cmd/server/UPSTREAM_COMMIT`：更新本次合并对应的上游提交为 `7dc7cfce1db5d31599815ff29acf6847ead0f0b7`。
- 上游合并文件：README、Docker、ent schema、迁移、账号/订阅/调度/网关服务、Grok 媒体路由、OpenAI/Spark shadow、前端账号/设置/用量页面及 i18n 等文件来自 `upstream/main` 的本次合并，未额外扩大人工改动范围。
- `progress.md`：追加本轮合并、验证和回滚记录。
- 回滚方式：提交后执行 `git revert -m 1 <本次合并提交>`；若只回滚本轮人工冲突/元数据处理，可执行 `git checkout HEAD^ -- backend/cmd/server/VERSION backend/cmd/server/UPSTREAM_COMMIT progress.md`。

## 2026-07-05 - Task: 合并原仓库 0.1.144 并准备发版
### What was done
- 合并原仓库 `Wei-Shaw/sub2api` 最新 `main`，上游基线为 `0.1.144` / `b650bdd6e238396af88508c73cefaeb32409d9ba`。
- 解决安装初始化迁移超时和前端 i18n 文案冲突：保留上游可配置迁移超时、IP 地理信息文案，同时保留当前仓库 Token 激励计划文案。
- 同步版本元数据到 `0.1.144`，并更新记录的上游提交。

### Testing
- `git diff --check`（在 `D:\project\sub2api-so`）通过。
- `go test ./...`（在 `D:\project\sub2api-so\backend`）通过。
- `D:\environment\nodejs\node-v22.17.0-win-x64\node_modules\@pnpm\exe\pnpm.exe run build`（在 `D:\project\sub2api-so\frontend`）通过；仅保留既有 Browserslist、动态/静态导入混用和 chunk size 警告。

### Notes
- `backend/cmd/server/VERSION`：同步上游基础版本到 `0.1.144`。
- `backend/cmd/server/UPSTREAM_COMMIT`：更新本次合并对应的上游提交为 `b650bdd6e238396af88508c73cefaeb32409d9ba`。
- `backend/internal/setup/setup.go`：冲突处理保留上游新增的可配置安装迁移超时逻辑。
- `frontend/src/i18n/locales/en.ts`：冲突处理同时保留 Token 激励计划和 IP 地理信息英文文案。
- `frontend/src/i18n/locales/zh.ts`：冲突处理同时保留 Token 激励计划和 IP 地理信息中文文案。
- 上游合并文件：README、部署配置、ent schema、迁移、账号/订阅/调度/网关/计费/并发/用量服务、前端账号/分组/订阅/用量/错误详情/i18n 等变更来自 `upstream/main` 的本次合并，未额外扩大人工改动范围。
- `progress.md`：追加本轮合并、验证和回滚记录。
- 回滚方式：提交后执行 `git revert -m 1 <本次合并提交>`；若只回滚本轮人工冲突/元数据处理，可执行 `git checkout HEAD^ -- backend/cmd/server/VERSION backend/cmd/server/UPSTREAM_COMMIT backend/internal/setup/setup.go frontend/src/i18n/locales/en.ts frontend/src/i18n/locales/zh.ts progress.md`。

## 2026-07-08 - Task: 合并原仓库 0.1.146 最新代码并排除首页改动发版
### What was done
- 合并原仓库 `Wei-Shaw/sub2api` 最新 `main`，上游基线为 `0.1.146` / `6631cbad67a0c15edf1006d2d0d60dc98169adf7`。
- 解决设置模块拆分、前端 i18n 拆分带来的二开冲突；保留本仓库 Token 激励计划配置、领取状态文案、余额记录类型，同时接入上游新增拆分文件结构。
- 当前首页相关未提交改动已继续隔离在 stash `wip-homepage-not-for-release-20260708-090240`，未纳入本次合并、推送或发版。

### Testing
- `git diff --check`（在 `D:\project\sub2api-so`）通过。
- `go test ./...`（在 `D:\project\sub2api-so\backend`）通过。
- `D:\environment\nodejs\node-v22.17.0-win-x64\node_modules\@pnpm\exe\pnpm.exe run build`（在 `D:\project\sub2api-so\frontend`）通过；仅保留既有 Browserslist、动态/静态导入混用和 chunk size 警告。

### Notes
- `backend/cmd/server/UPSTREAM_COMMIT`：记录本次合并对应的上游最新提交 `6631cbad67a0c15edf1006d2d0d60dc98169adf7`。
- `backend/internal/handler/admin/setting_handler.go`：在上游设置处理器拆分后保留 Token 激励计划的设置读取与 DTO 转换。
- `backend/internal/handler/admin/setting_handler_update.go`：在上游新增的设置更新处理器中接回 Token 激励计划开关与档位保存逻辑。
- `backend/internal/handler/admin/setting_handler_audit.go`：在设置审计变更列表中保留 Token 激励计划相关字段。
- `backend/internal/service/setting_update.go`、`backend/internal/service/setting_parse.go`、`backend/internal/service/setting_public.go`、`backend/internal/service/setting_token_incentive.go`：适配上游设置服务拆分，并保留 Token 激励默认值、解析、公开开关和规则读取逻辑。
- `frontend/src/i18n/locales/en/**`、`frontend/src/i18n/locales/zh/**`：适配上游 i18n 目录拆分，并把 Token 激励计划文案迁移到新的语言包结构。
- 上游合并文件：其余 backend/frontend/deploy/docs/README 等大量变更来自 `upstream/main` 本次合并，未额外扩大人工改动范围。
- `progress.md`：追加本轮合并、验证和回滚记录。
- 回滚方式：提交后执行 `git revert -m 1 <本次合并提交>`；如只需找回未发布首页改动，执行 `git stash list` 找到 `wip-homepage-not-for-release-20260708-090240` 后再按需 `git stash apply`。

## 2026-07-09 - Task: 合并原仓库 0.1.146 最新提交并排除首页改动发版
### What was done
- 合并原仓库 `Wei-Shaw/sub2api` 最新 `main` 到当前仓库，本轮上游仍属于 `0.1.146` 基线，最新提交为 `6f43986c376d76144cb39c7a562c179e19ac7439`。
- 记录本仓库更新检测用的上游提交指针，确保 Web 端只按原仓库新版本/当前仓库新版本逻辑识别更新，不把普通未发布首页改动带入发版。
- 当前本地首页及其他未完成改动已隔离在 stash `wip-local-not-for-upstream-release-20260709-120915`，本轮提交、推送和发版不包含这些首页改动。

### Testing
- `git fetch upstream main --tags`（在 `D:\project\sub2api-so`）通过，确认 `upstream/main` 为 `6f43986c376d76144cb39c7a562c179e19ac7439`，最新上游 tag 仍为 `v0.1.146`。
- `git diff --check`（在 `D:\project\sub2api-so`）通过。
- `go test ./...`（在 `D:\project\sub2api-soackend`）通过。
- `D:\environment
odejs
ode-v22.17.0-win-x64
ode_modules\@pnpm\exe\pnpm.exe run build`（在 `D:\project\sub2api-sorontend`）通过；仅保留既有 Browserslist、动态/静态导入混用和 chunk size 警告。

### Notes
- `backend/cmd/server/UPSTREAM_COMMIT`：更新本次合并对应的上游最新提交为 `6f43986c376d76144cb39c7a562c179e19ac7439`。
- `progress.md`：追加本轮合并、验证、首页改动隔离和回滚记录。
- 上游合并内容：引入上游 admin scheduler score opt-in、compact body signal routing、sidebar scroll position persist 等已发布基线后的最新修复提交。
- 首页相关未发布改动：仍保留在 `stash@{0}` / `wip-local-not-for-upstream-release-20260709-120915`，未 stage、未 commit、未 push。
- 回滚方式：本轮记录提交完成后可执行 `git revert HEAD` 回滚记录提交，再执行 `git revert -m 1 d750df90` 回滚本次上游合并；如只需恢复未发布首页改动，执行 `git stash list` 找到 `wip-local-not-for-upstream-release-20260709-120915` 后再按需 `git stash apply`。

## 2026-07-09 - Task: 修复 GitHub Security Scan 后端 govulncheck 失败
### What was done
- 定位 GitHub Security Scan 失败原因为 `govulncheck` 命中 Go 标准库漏洞 `GO-2026-5856`，当前 `go1.26.4` 的 `crypto/tls` 受影响，修复版本为 `go1.26.5`。
- 将后端 Go 基线、GitHub CI/Release/Security Scan 校验和 Docker 构建镜像统一升级到 `1.26.5`。
- 新增构建安全基线说明，避免后续把安全扫描和发布构建回退到受影响的 Go 版本。

### Testing
- `GOTOOLCHAIN=auto go run golang.org/x/vuln/cmd/govulncheck@latest ./...`（在 `D:\project\sub2api-soackend`）通过，结果为 `No vulnerabilities found` / `Your code is affected by 0 vulnerabilities`。
- `GOTOOLCHAIN=auto go test ./...`（在 `D:\project\sub2api-soackend`）通过。
- `git diff --check`（在 `D:\project\sub2api-so`）通过。

### Notes
- `backend/go.mod`：后端 Go 基线从 `1.26.4` 升级到 `1.26.5`。
- `.github/workflows/backend-ci.yml`、`.github/workflows/security-scan.yml`、`.github/workflows/release.yml`：同步 CI、安全扫描和 Release 的 Go 版本校验到 `go1.26.5`。
- `Dockerfile`、`backend/Dockerfile`、`deploy/Dockerfile`：同步 Docker 构建阶段 Go 镜像到 `golang:1.26.5-alpine`。
- `docs/BUILD_SECURITY.md`：记录 Go 构建安全基线和版本同步要求。
- `progress.md`：追加本轮失败定位、修复、验证和回滚记录。
- 回滚方式：提交后执行 `git revert <本轮提交>` 可回退 Go 版本、安全扫描配置、Docker 基线和文档记录；如只是撤销发版 tag，删除本轮 tag 后重新发布上一版即可。

## 2026-07-09 - Task: 补回 Grok 4.5 支持模型列表
### What was done
- 确认此前 Grok 4.5 改动被包含在 `wip-local-not-for-upstream-release-20260709-120915` stash 中，未进入当前 `main` 和 `v0.1.146-fy.3` 发版。
- 只补回 Grok 相关改动：后端默认模型、模型映射、前端白名单、预设映射、README 和 Grok 独立文档。
- 新增 `grok-4.5`、`grok-4.5-latest`、`grok-build-latest`，并保留 `grok` / `grok-latest` 继续指向 `grok-4.3`，避免升级后默认行为突然变化。
- 补充 `grok-4.5` 兜底计费，按 xAI Pricing 当前价格记录为 `$2.00 / 1M input tokens`、`$0.50 / 1M cached input tokens`、`$6.00 / 1M output tokens`。

### Testing
- `go test ./internal/pkg/xai`（在 `D:\project\sub2api-soackend`）通过。
- `go test ./internal/service -run "Grok|Billing"`（在 `D:\project\sub2api-soackend`）通过。
- `git diff --check`（在 `D:\project\sub2api-so`）通过。
- `D:\environment\nodejs\node-v22.17.0-win-x64\node_modules\@pnpm\exe\pnpm.exe run typecheck`（在 `D:\project\sub2api-so\frontend`，临时隔离首页 WIP 后）通过。

### Notes
- `backend/internal/pkg/xai/models.go`：新增 Grok 4.5 默认模型和显式 4.5/latest aliases。
- `backend/internal/pkg/xai/oauth_test.go`：补充 Grok 4.5 alias 映射断言。
- `backend/internal/service/billing_service.go`：新增 Grok 4.5 兜底计费。
- `frontend/src/composables/useModelWhitelist.ts`：同步前端模型白名单和 Grok 预设映射。
- `README.md`、`docs/GROK_XAI_SUPPORT.md`：补充 Grok 4.5 支持说明。
- `.gitignore`：将 `docs/GROK_XAI_SUPPORT.md` 加入 docs 白名单，保证文档可被 Git 跟踪。
- `progress.md`：追加本轮补回原因、验证和回滚记录。
- 回滚方式：提交后执行 `git revert <本轮提交>` 可撤销 Grok 4.5 模型列表、计费、文档和白名单变更；首页 WIP 不在本轮提交内。

## 2026-07-10 - Task: 合并原仓库 0.1.150 最新代码并排除本地未提交修改发版
### What was done
- 将当前本地未提交修改隔离到 stash `wip-local-not-for-upstream-release-20260710-162322`，本轮合并、提交、推送和发版不包含这些本地 WIP。
- 合并原仓库 `Wei-Shaw/sub2api` 最新 `main` 到当前仓库，上游基线版本为 `0.1.150`，最新提交为 `6dd3274aafbc1a7a91304380fb3d7e50406841e0`。
- 处理上游 `v0.1.150` 与当前二开逻辑的冲突，保留本仓库 GitHub Release 打包、fork 更新目标、原仓库按 release 版本提示的更新逻辑，并合入上游回滚候选、Grok 4.5 默认映射、GPT-5.6 计费/使用量修复等更新。
- 首页冲突按上游版本解决，未恢复或提交本地未完成首页/WebGL/纹理相关改动。

### Testing
- `git diff --check`（在 `D:\project\sub2api-so`）通过。
- `GOTOOLCHAIN=auto go test ./...`（在 `D:\project\sub2api-soackend`）通过。
- `D:\environment
odejs
ode-v22.17.0-win-x64
ode_modules\@pnpm\exe\pnpm.exe run build`（在 `D:\project\sub2api-sorontend`）通过；仅保留既有 Browserslist、动态/静态导入混用和 chunk size 警告。

### Notes
- `.github/workflows/release.yml`：解决上游合并冲突，保留本仓库 Release 打包流程和 Go 1.26.5 校验。
- `backend/cmd/server/VERSION`、`backend/cmd/server/UPSTREAM_COMMIT`：同步上游基线版本 `0.1.150` 和最新上游提交 `6dd3274aafbc1a7a91304380fb3d7e50406841e0`。
- `backend/internal/service/update_service.go`、`backend/internal/repository/github_release_service.go` 及相关测试：合并上游回滚 release 列表能力，同时保留当前仓库更新目标和原仓库仅按 release 版本提示的二开逻辑。
- `backend/internal/pkg/xai/models.go`、`backend/internal/service/billing_service.go`、`frontend/src/composables/useModelWhitelist.ts` 及相关测试：合入上游 Grok 4.5 默认模型映射和计费逻辑。
- `frontend/src/components/common/VersionBadge.vue`：冲突处理保留当前仓库已发布的更新提示和应用内更新/重启入口。
- `frontend/src/views/HomeView.vue`：冲突处理采用上游版本，避免带入本地未完成首页改动。
- 上游合并文件：其余 backend/frontend/docs/deploy/README/ent/migrations/resources 等变更均来自 `upstream/main` 本次合并，未额外扩大人工改动范围。
- `progress.md`：追加本轮合并、验证、本地 WIP 隔离和回滚记录。
- 回滚方式：提交后执行 `git revert -m 1 <本次合并提交>` 回滚上游合并，再执行 `git revert <本轮记录提交>` 回滚元数据/进度记录；如需恢复本地未提交修改，执行 `git stash list` 找到 `wip-local-not-for-upstream-release-20260710-162322` 后再按需 `git stash apply`。

## 2026-07-12 - Task: 制作统一 AI 网关核心视觉概念例图
### What was done
- 将已确认的统一外框、同心机械核心、纯中性白光、黑钛金属、烟熏玻璃和稀疏深空背景组合为桌面首页概念例图。
- 提供 Core、Route、Code、Meter 四相位交互预览，展示稳定外轮廓内的环体、导轨、接口探针和计量机构变化。

### Testing
- 本地视觉预览服务返回 HTTP 200，浏览器 `1440x900` 下 Core 画面非空。
- 点击 Route 后相位状态和标题同步更新；恢复后的预览重新停在 Core。

### Notes
- `.superpowers/brainstorm/unified-core-20260712/content/unified-gateway-core-concept-v2.html`：新增本地交互式视觉概念例图，不进入产品构建。
- `progress.md`：追加本轮例图交付、验证和回滚点。
- 回滚方式：删除 `.superpowers/brainstorm/unified-core-20260712/`，并删除 `progress.md` 中本条记录。

## 2026-07-12 - Task: 合并原仓库 0.1.151 最新代码并排除当前本地修改发版
### What was done
- 将当前首页、首页资源、设计文档和工具预览等未提交修改隔离到 `wip-local-not-for-upstream-release-20260712` 与 `wip-superpowers-not-for-upstream-release-20260712`，并将合并期间出现的并行进度记录隔离到 `wip-progress-concurrent-not-for-upstream-release-20260712`。
- 合并原仓库 `Wei-Shaw/sub2api` 最新 `main`，上游正式版本为 `0.1.151`，最新提交为 `e316ebf52838a89d57fc790981cce7520f819ac8`；本次自动合并无文本冲突。
- 保持当前仓库的 Release 打包、更新目标、原仓库仅按正式 Release 提示、应用内更新/重启和 Token 激励逻辑不变，同时接入上游 OpenAI/Anthropic 工具转换、Fast/Flex 策略、缓存创建 Token 统计及错误捕获修复。

### Testing
- `git diff --check` 与 `git diff --cached --check`（在 `D:\project\sub2api-so`）通过。
- `GOTOOLCHAIN=auto go test ./...`（在 `D:\project\sub2api-so\backend`）通过。
- `D:\environment\nodejs\node-v22.17.0-win-x64\node_modules\@pnpm\exe\pnpm.exe run build`（在 `D:\project\sub2api-so\frontend`）通过；仅有既有 Browserslist、动态/静态导入混用和 chunk size 警告。

### Notes
- `backend/cmd/server/VERSION`：同步上游正式版本为 `0.1.151`。
- `backend/cmd/server/UPSTREAM_COMMIT`：记录本轮上游最新提交 `e316ebf52838a89d57fc790981cce7520f819ac8`。
- `backend/internal/handler/admin/admin_helpers_test.go`：同步上游管理员辅助测试的设置字段覆盖。
- `backend/internal/handler/dto/settings.go`：同步上游用户级 Fast/Flex 策略设置 DTO。
- `backend/internal/handler/ops_capture_writer_nil_test.go`：新增上游响应捕获 writer 释放后访问的回归测试。
- `backend/internal/handler/ops_error_logger.go`：同步上游响应捕获 writer 的 nil 安全修复。
- `backend/internal/pkg/apicompat/anthropic_to_responses_response.go`：同步 Anthropic 到 Responses 工具调用转换修复。
- `backend/internal/pkg/apicompat/chatcompletions_responses_bridge.go`：同步 custom/tool_search 工具桥接、命名空间和降级选择逻辑。
- `backend/internal/pkg/apicompat/chatcompletions_responses_bridge_custom_tools_test.go`：新增上游 custom/tool_search 桥接回归覆盖。
- `backend/internal/pkg/apicompat/chatcompletions_responses_test.go`：同步工具转换断言。
- `backend/internal/pkg/apicompat/responses_anthropic_cache_creation_test.go`：新增 cache_creation_input_tokens 转换回归测试。
- `backend/internal/pkg/apicompat/responses_stream_event_wire.go`：同步流式工具事件线协议转换。
- `backend/internal/pkg/apicompat/responses_stream_event_wire_test.go`：同步流式工具事件测试。
- `backend/internal/pkg/apicompat/responses_to_anthropic.go`：同步 Responses 到 Anthropic 工具和缓存 Token 转换。
- `backend/internal/pkg/apicompat/types.go`：同步 custom/tool_search 与缓存 Token 类型定义。
- `backend/internal/pkg/ctxkey/ctxkey.go`：新增上游请求身份与用户策略上下文键。
- `backend/internal/pkg/openai/request.go`：同步 OpenAI 请求身份字段和解析逻辑。
- `backend/internal/pkg/openai/request_identity_test.go`：新增 OpenAI 请求身份解析测试。
- `backend/internal/server/middleware/api_key_auth.go`：将用户级 Fast/Flex 策略写入请求上下文。
- `backend/internal/server/middleware/api_key_auth_test.go`：同步 API Key 鉴权策略测试。
- `backend/internal/server/middleware/openai_fast_policy_forwarding_test.go`：新增用户策略转发回归测试。
- `backend/internal/service/account_test_service.go`：同步测试服务的用户策略字段。
- `backend/internal/service/account_usage_service.go`：同步账户用量服务的用户策略传递。
- `backend/internal/service/openai_codex_identity.go`：新增 Codex 请求 originator 与 User-Agent 处理。
- `backend/internal/service/openai_codex_identity_test.go`：新增 Codex 请求身份测试。
- `backend/internal/service/openai_fast_policy_test.go`：同步 Fast/Flex 策略测试。
- `backend/internal/service/openai_fast_policy_ws_test.go`：同步 WebSocket Fast/Flex 策略测试。
- `backend/internal/service/openai_gateway_forward.go`：同步网关请求身份转发。
- `backend/internal/service/openai_gateway_grok.go`：同步 Grok reasoning effort 兼容处理。
- `backend/internal/service/openai_gateway_grok_test.go`：同步 Grok reasoning effort 回归测试。
- `backend/internal/service/openai_gateway_messages.go`：同步 Messages 路径请求身份和缓存 Token 处理。
- `backend/internal/service/openai_gateway_messages_chat_fallback.go`：同步 Messages 降级桥接参数。
- `backend/internal/service/openai_gateway_passthrough.go`：同步透传请求身份字段。
- `backend/internal/service/openai_gateway_request_body.go`：同步请求体工具转换和身份处理。
- `backend/internal/service/openai_gateway_responses_chat_fallback.go`：同步 Responses 降级桥接工具转换。
- `backend/internal/service/openai_gateway_service_test.go`：同步网关服务构造和行为测试。
- `backend/internal/service/openai_oauth_passthrough_test.go`：同步 OAuth 透传身份测试。
- `backend/internal/service/openai_ws_forwarder_payload.go`：同步 WebSocket 请求策略载荷。
- `backend/internal/service/openai_ws_forwarder_success_test.go`：同步 WebSocket 成功路径测试。
- `backend/internal/service/setting_features.go`：新增管理员控制的用户级 Fast/Flex 策略开关。
- `backend/internal/service/settings_view.go`：将 Fast/Flex 策略开关加入设置视图。
- `frontend/src/api/admin/settings.ts`：同步管理员设置 API 类型。
- `frontend/src/i18n/locales/en/admin/settings.ts`：新增 Fast/Flex 策略英文设置文案。
- `frontend/src/i18n/locales/zh/admin/settings.ts`：新增 Fast/Flex 策略中文设置文案。
- `frontend/src/views/admin/SettingsView.vue`：新增 Fast/Flex 用户策略管理开关界面。
- `progress.md`：追加本轮隔离、合并、验证和回滚记录。
- 本轮没有数据库结构、部署方式或外部配置文件变化；新增策略通过现有管理设置界面配置，无需新增部署文档。
- 回滚方式：发版后执行 `git revert -m 1 <本次合并提交>` 回滚本轮上游合并；当前本地修改按需从上述三个 `20260712` stash 分别恢复，不应直接批量 `stash pop`。

## 2026-07-12 - Task: 迭代统一 AI 网关核心震撼感例图
### What was done
- 保留统一黑钛外框、纯中性白光和稀疏深空背景，将四相位从小幅环体旋转升级为大尺度内部机械重组。
- Route 展开四向路由闸，Code 接入三层接口舱，Meter 升起三列计量机构；换相增加短促白光冲击和轻微整体推进。

### Testing
- 本地预览服务返回 HTTP 200，并正确加载“震撼增强版例图”。
- 浏览器 `1440x900` 实际切换并截图检查 Route、Code、Meter，三种内部剪影和对应标题均明显不同；最终恢复 Core。

### Notes
- `.superpowers/brainstorm/unified-core-20260712/content/unified-gateway-core-concept-v2.html`：保留上一版高级感基线例图。
- `.superpowers/brainstorm/unified-core-20260712/content/unified-gateway-core-impact-v3.html`：新增大尺度机械变化与白光冲击的增强版例图。
- `progress.md`：追加本轮例图迭代、验证和回滚点。
- 回滚方式：删除 `.superpowers/brainstorm/unified-core-20260712/`，并删除 `progress.md` 中本条记录。

## 2026-07-13 - Task: 制作外框联动变形的 B 方向视觉例图
### What was done
- 按 B 方向让统一 AI 网关核心的外部框架参与四相位变形，同时保留纯白舞台光、黑钛材质和稀疏深空背景。
- Route 横向展开框架与四向路由闸，Code 开启框架并前推接口舱，Meter 纵向拉伸框架并升起计量塔；Core 保持闭合稳定的默认轮廓。

### Testing
- 本地视觉预览 `http://localhost:53007/` 返回 HTTP 200。
- 浏览器 `1440x900` 依次切换 Route、Code、Meter、Core，装置和标题同步更新；四个状态均保持单屏，无横向或纵向溢出，最终恢复 Core。

### Notes
- `.superpowers/brainstorm/unified-core-20260712/content/unified-gateway-core-stage-transform-v4.html`：新增外框参与大幅机械换形的 B 方向交互例图。
- `progress.md`：追加本轮视觉例图、验证证据和回滚点。
- 回滚方式：删除 `.superpowers/brainstorm/unified-core-20260712/content/unified-gateway-core-stage-transform-v4.html`，并删除 `progress.md` 中本条记录。

## 2026-07-13 - Task: 制作空间折叠入口 B 方向视觉例图
### What was done
- 放弃机械装置和中心模型，改用全屏 Three.js 空间裂隙作为统一 AI 网关入口，并保持纯黑背景与中性白光。
- Core 使用稳定纵向入口，Route 展开多路径空间分流，Code 打开纵深接入通道，Meter 切换为横向计量脉冲；四阶段同步更新首页标题与业务文案。
- 增加持续空间运动、换相冲击和鼠标视差，桌面端将视觉重心置于右侧，手机端将入口移到文字下方。

### Testing
- 本地视觉预览 `http://localhost:53007/` 返回 HTTP 200，并加载“Sub2API 空间折叠网关概念例图”。
- 浏览器 `1440x900` 依次切换 Core、Route、Code、Meter，标题和场景状态同步更新，页面宽高与视口一致。
- 浏览器 `390x844` 检查四阶段标题边界、正文与底部阶段导航，页面保持单屏且无横向或纵向溢出。
- 截图像素检查确认桌面右侧场景非空；连续帧存在像素变化，鼠标从中央移动到右上后右侧场景产生明显像素差异。

### Notes
- `.superpowers/brainstorm/unified-core-20260712/content/spatial-fold-gateway-v5.html`：新增空间折叠入口的 Three.js 交互视觉例图，不进入生产首页。
- `progress.md`：追加本轮视觉方向、响应式验证、动态验证和回滚点。
- 回滚方式：删除 `.superpowers/brainstorm/unified-core-20260712/content/spatial-fold-gateway-v5.html`，并删除 `progress.md` 中本条记录。

## 2026-07-13 - Task: 制作四个全新首页视觉方向示例图
### What was done
- 在机械装置、空间裂隙、模型文字矩阵和界面拼贴之外，提供真实开发现场、纪念性基础设施、生成能力电影、请求数据实景四个独立方向。
- 每个方向均使用真实影像、Sub2API 首页导航、业务标题与简短价值表达，形成可直接比较的首页示例图，并支持点击标记。
- 对比板采用桌面双列、手机单列布局，不把任何方向预设为推荐答案。

### Testing
- 本地视觉预览 `http://localhost:53007/` 返回 HTTP 200，并加载“Sub2API 首页视觉方向对比”。
- 浏览器 `1440x900` 全页截图确认四张示例图及外部影像素材均正常显示，四个方向的画面来源和构图明显不同。
- 浏览器 `390x844` 检查四张示例图、标题宽度和页面滚动区域，无横向溢出；点击 B 后选中状态正确更新。

### Notes
- `.superpowers/brainstorm/unified-core-20260712/content/homepage-direction-board-v6.html`：新增四个首页视觉方向的高保真比较板。
- `progress.md`：追加本轮示例图交付、浏览器验证和回滚点。
- 回滚方式：删除 `.superpowers/brainstorm/unified-core-20260712/content/homepage-direction-board-v6.html`，并删除 `progress.md` 中本条记录。

## 2026-07-13 - Task: 子代理并行制作 A、C、D 三个完整首页候选
### What was done
- 将真实开发现场、生成能力电影、请求数据实景拆分为三个独立文件责任边界，由三个子代理并行完成完整一屏首页候选。
- 三个候选均实现 Core、Route、Code、Meter 四阶段，阶段切换同步改变标题、业务文案和场景动画，并提供鼠标视差与减弱动画支持。
- 新增统一切换入口，在同一视口内直接切换 A、C、D 完整页面，同时保留独立打开入口。
- 统一手机验收发现 A 的装饰层扩大内部滚动区域，语义点击会将首屏内容上移；修复为不可程序滚动裁切，并将三个候选及统一入口的矩形焦点框替换为无框反馈。

### Testing
- A、C、D 页面及统一入口在 `http://localhost:53008/` 下均返回 HTTP 200；六个外部影像素材 URL 均返回 HTTP 200。
- 浏览器 `1440x900` 对 A、C、D 分别点击 Core、Route、Code、Meter，共 12 次状态切换；标题均同步更新，页面宽高与视口一致。
- 浏览器 `390x844` 对三个候选点击 Route 并截图检查；标题、正文、主动作和阶段导航均在视口内，页面无横向或纵向溢出。
- A 修复后 `.hero` 的 `scrollTop` 在语义点击 Route 后保持 0；三个候选及统一入口的活动按钮 `outline-style` 均为 `none`。
- 统一入口在桌面和手机端切换 A、C、D 后，描述、iframe 地址和加载状态正确更新；手机端保持 `390x844` 无页面溢出。
- 四个新增 HTML 文件的 whitespace diff check 通过，`progress.md` 的 `git diff --check` 通过。

### Notes
- `.superpowers/brainstorm/unified-core-20260712/content/variant-a-developer-studio-v7.html`：新增真实开发现场完整候选页，并修复手机内部滚动和焦点框问题。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-c-generative-cinema-v7.html`：新增生成能力电影完整候选页，并使用无框键盘焦点反馈。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-d-request-telemetry-v7.html`：新增请求数据实景完整候选页，并使用无框键盘焦点反馈。
- `.superpowers/brainstorm/unified-core-20260712/content/acd-variants-index-v8.html`：新增 A、C、D 统一切换和独立预览入口。
- `progress.md`：追加本轮并行制作、缺陷修复、统一验证和回滚点。
- 回滚方式：先停止监听 53008 端口的本地预览进程，再删除上述四个 HTML 文件，并删除 `progress.md` 中本条记录。

## 2026-07-13 - Task: 扩展五种首页候选并增加固定一屏滚轮换相
### What was done
- 为 A、C、D 增加统一滚轮协议：页面本身不滚动，滚轮向下依次切换阶段，向上反向切换，并使用短暂锁定避免触控板连续跳相。
- A 与 D 从单一背景升级为 Core、Route、Code、Meter 四套真实影像，并分别保留开发现场与请求链路的动画重点。
- C 将原有切片位移动画升级为镜头推进、曝光变化和斜向擦除的电影式换相。
- 新增 E 简约风，以纯白界面、单线网络与超大文字表达统一入口；新增 F 黑白二次元风，以真实角色影像、网点与速度线表达四阶段。
- 统一入口扩展为 A、C、D、E、F 五方案，手机端五个方案按钮保持单行五等分。

### Testing
- 实施前浏览器基线验证确认 A、C、D 滚轮后仍停留 Core；实施后五个候选向下滚轮均从 Core 切到 Route，向上均回到 Core。
- 五个候选在 `1440x900` 和 `390x844` 下执行滚轮换相后，`scrollY` 均为 0，页面宽高与视口一致，标题保持在视口内。
- A 与 D 分别检查四阶段可见背景图，均得到四个唯一影像 URL；相关外部素材返回 HTTP 200。
- C 在换相中间帧与结束帧均产生显著截图像素差异，确认新镜头转场实际运行。
- E、F 完成桌面与手机截图检查，导航、标题、主动作和阶段栏无重叠或溢出。
- 五方案入口依次加载 A、C、D、E、F，iframe 地址、描述和加载状态同步更新；手机 `390x844` 下五个方案按钮保持同一行且页面无溢出。

### Notes
- `.superpowers/brainstorm/unified-core-20260712/content/variant-a-developer-studio-v7.html`：增加四阶段真实背景、滚轮换相和滚轮提示。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-c-generative-cinema-v7.html`：重做电影式镜头转场并增加滚轮换相。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-d-request-telemetry-v7.html`：增加四阶段技术背景、滚轮换相和滚轮提示。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-e-minimal-v9.html`：新增固定一屏简约风候选。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-f-anime-v9.html`：新增固定一屏黑白二次元风候选。
- `.superpowers/brainstorm/unified-core-20260712/content/acd-variants-index-v8.html`：扩展为五方案统一入口并修复手机方案栏换行。
- `progress.md`：追加本轮五方案扩展、滚轮协议、视觉与响应式验证记录。
- 回滚方式：删除 E、F 两个新增候选，并在 A、C、D 与统一入口中反向移除本轮新增的滚轮、阶段背景、镜头擦除及五方案配置；随后删除 `progress.md` 中本条记录。

## 2026-07-13 - Task: 将五种固定一屏首页候选升级为独立彩色视觉
### What was done
- 保留五个候选的固定一屏、四阶段滚轮换相和鼠标交互，为 A、C、D、E 分别建立深蓝开发现场、彩色电影镜头、彩色遥测实景和明亮简约四套独立视觉体系。
- 将 F 从单张黑白角色影像重做为四幕彩色二次元场景，每个阶段使用独立背景，并增加景深推进、斜切闪帧、速度线爆发和阶段色彩联动。
- 修正 C 手机端专用遮罩造成的近黑白问题，并将 F 的页面标题和统一入口标签同步改为“彩色二次元风”。

### Testing
- 统一入口与五个候选页面均返回 HTTP 200；15 个外部影像 URL 的 HEAD 请求均返回 HTTP 200。
- 浏览器 `1440x900` 验证五页滚轮向下均从 Core 切到 Route、向上均回到 Core，切换前后 `scrollY` 始终为 0，页面宽高与视口一致。
- 桌面截图彩色像素占比为 A `57.0%`、C `22.9%`、D `6.5%`、E `13.0%`、F `63.8%`，确认五个方案均不再呈现纯黑白画面。
- 浏览器 `390x844` 验证五页文案、主动作与阶段栏处于视口内，页面无横向或纵向溢出；C 手机端修正后彩色像素占比由 `0.5%` 提升至 `39.2%`。
- A、D、F 分别逐阶段读取可见背景图，Core、Route、Code、Meter 均得到四个唯一影像 URL；F 的 Core 与 Route 截图确认人物、背景和主题色同步发生明显变化。

### Notes
- `.superpowers/brainstorm/unified-core-20260712/content/variant-a-developer-studio-v7.html`：为四阶段开发现场补充深蓝、蓝绿、紫色与暖橙的独立色彩氛围。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-c-generative-cinema-v7.html`：移除灰度镜头处理，增加阶段色彩覆盖并修复手机端近黑白遮罩。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-d-request-telemetry-v7.html`：将四套遥测实景升级为蓝、青绿、紫与金色的阶段化视觉。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-e-minimal-v9.html`：使用阶段色变量联动主动作、路径、大字、节点和轻量背景光。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-f-anime-v9.html`：重做为四套彩色二次元背景与复合转场动画。
- `.superpowers/brainstorm/unified-core-20260712/content/acd-variants-index-v8.html`：将 F 方案标签更新为彩色二次元风。
- `progress.md`：追加本轮彩色化、四场景、交互和响应式验证记录。
- 回滚方式：以本条记录前的工作区状态为回滚点，反向移除上述六个候选 HTML 中本轮新增的色彩、F 四场景和 C 手机遮罩补丁，并删除 `progress.md` 中本条记录。

## 2026-07-13 - Task: 并行升级前三种 3D 首页并重做简约与二次元方案
### What was done
- 以独立文件责任边界并行重做 A、C、D：分别使用真实 Three.js 构建 3D 开发工作台、电影棱镜空间和请求遥测隧道，四阶段均改变几何结构、光色和镜头，而非只切换背景或颜色。
- 将 E 从单线极简升级为动态极简流场，增加多层曲线路径、轨道、节点、斜向色域和阶段形变，并将底部长分隔线改为圆点与短刻度。
- 将 F 的四阶段替换为四位不同的彩色 2D 美少女本地素材，并为 Route、Code、Meter 分别增加分叉扫线、字符式扫描和环形计量轨迹。
- 终审修复 C 在减弱动态效果模式下累积多个阶段模型、Code 动画覆盖原始纵深两项问题；统一入口同步更新五个方案名称。

### Testing
- 统一入口、五个候选页面、四张本地 WebP 和 Three.js CDN 均返回 HTTP 200；四张角色图哈希各不相同。
- 浏览器 `1440x900` 验证 A、C、D 均存在覆盖视口的非空 WebGL Canvas；五页滚轮向下均从 Core 切至 Route、向上回到 Core，切换前后 `scrollY=0`，控制台无脚本错误。
- A、C、D 依次执行 Core、Route、Code、Meter，右侧场景相邻阶段平均像素差异分别为 A `10.41–27.03`、C `11.13–16.36`、D `11.71–19.02`，确认几何形态和画面实际变化。
- 浏览器 `390x844` 验证五页宽高与视口一致、无横向或纵向溢出；A、C、D 手机 Canvas 均为 `390x844` 且非空，E/F 文案、动作与阶段导航均在视口内。
- C 在 `prefers-reduced-motion` 下依次切换五次，每个阶段始终只有一个可见组；活动组透明度与缩放为 1，其余组均隐藏。
- F 依次滚轮切换四阶段，浏览器读取到 `girl-1.webp` 至 `girl-4.webp` 四个唯一可见背景，页面本身始终不滚动。
- 六个候选 HTML 的 whitespace check 与 `progress.md` 的 `git diff --check` 通过。

### Notes
- `.superpowers/brainstorm/unified-core-20260712/content/variant-a-developer-studio-v7.html`：重做为四阶段 Three.js 3D 开发工作台。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-c-generative-cinema-v7.html`：重做为 3D 电影棱镜，并修复减弱动画和代码纵深问题。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-d-request-telemetry-v7.html`：重做为四阶段 3D 请求遥测隧道。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-e-minimal-v9.html`：增加动态极简流场并移除底部长分隔线。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-f-anime-v9.html`：改用四位不同的 2D 美少女并增加阶段专属动画语法。
- `.superpowers/brainstorm/unified-core-20260712/content/assets/anime-v10/girl-1.webp`：新增 F Core 阶段的本地 2D 插画素材。
- `.superpowers/brainstorm/unified-core-20260712/content/assets/anime-v10/girl-2.webp`：新增 F Route 阶段的本地 2D 插画素材。
- `.superpowers/brainstorm/unified-core-20260712/content/assets/anime-v10/girl-3.webp`：新增 F Code 阶段的本地 2D 插画素材。
- `.superpowers/brainstorm/unified-core-20260712/content/assets/anime-v10/girl-4.webp`：新增 F Meter 阶段的本地 2D 插画素材。
- `.superpowers/brainstorm/unified-core-20260712/content/acd-variants-index-v8.html`：更新五个方案的统一入口名称。
- `progress.md`：追加本轮并行施工、终审修复、3D 与响应式验证记录。
- 当前 F 素材仅用于概念预览，作品来源已记录在 HTML 注释中，但没有可核验的商用授权；正式发布前须替换为自有或明确授权素材。A、C、D 当前依赖 Three.js CDN，生产化时需决定是否本地化依赖。
- 回滚方式：以本条记录前的工作区状态为回滚点，恢复上述六个 HTML 到本轮前版本，删除 `assets/anime-v10/`，并删除 `progress.md` 中本条记录。

## 2026-07-13 - Task: 将前三种中心 3D 模型改为真实内容驱动的 2.5D 动效
### What was done
- 完整移除 A、C、D 的中心球体、棱镜、轨道和 Three.js 模型代码，恢复真实开发、生成能力与遥测实景作为视觉主体。
- A 改为真实开发现场与无框工作流内容层：Core 汇聚请求，Route 空间分流，Code 多层代码穿梭，Meter 轨迹与刻度展开；同时移除自动轮播，保证滚轮操作确定性。
- C 改为多层影像切片与电影镜头：四阶段分别使用景深汇聚、纵向分流、代码推进和暖金旋转切入，并保留曝光闪帧、斜向擦除与鼠标视差。
- D 改为真实技术实景与无框 Canvas 数据流：请求脉冲在纵深中穿梭、分叉、重连和收束，不再使用中心几何主体。
- 统一入口将 A、C、D 名称更新为 2.5D 开发现场、2.5D 电影切片和 2.5D 遥测实景；E/F 保持上一轮状态。

### Testing
- 统一入口与 A、C、D 页面均返回 HTTP 200；三页引用的真实影像素材均可正常加载。
- 浏览器 `1440x900` 对三页依次执行 `Core → Route → Code → Meter → Code → Route → Core`，全部阶段顺序正确且全过程 `scrollY=0`，页面宽高与视口一致。
- 相邻阶段右侧画面平均像素差异分别为 A `12.33–18.03`、C `29.48–47.54`、D `18.65–26.57`，确认空间镜头和视觉内容实际发生明显变化。
- 浏览器 `390x844` 验证 A、C、D 均无横向或纵向溢出；手机滚轮向下从 Core 到 Route、向上返回 Core，Canvas 与影像切片保持非空且文案无重叠。
- A 页面加载后等待超过 8 秒仍保持 Core，随后滚轮向下切至 Route、向上回到 Core，确认自动轮播已移除。
- 鼠标移动后，A、C、D 分别更新工作流、远中近景和遥测视差变量；浏览器控制台无脚本错误。
- 静态检查确认 A、C、D 不再包含 Three.js CDN、WebGLRenderer 或几何模型代码；四个 HTML 的 whitespace check 与 `progress.md` 的 `git diff --check` 通过。

### Notes
- `.superpowers/brainstorm/unified-core-20260712/content/variant-a-developer-studio-v7.html`：将中心 3D 开发模型重做为真实现场与 2.5D 工作流动效，并取消自动换相。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-c-generative-cinema-v7.html`：将 3D 棱镜重做为真实影像切片与电影级透视镜头。
- `.superpowers/brainstorm/unified-core-20260712/content/variant-d-request-telemetry-v7.html`：将 3D 数据核心重做为真实遥测实景与纵深数据流。
- `.superpowers/brainstorm/unified-core-20260712/content/acd-variants-index-v8.html`：同步更新前三个方案的 2.5D 名称。
- `progress.md`：追加本轮方向纠偏、完整重做和运行态验证记录。
- 回滚方式：以上一条“并行升级前三种 3D 首页并重做简约与二次元方案”记录完成时的工作区为回滚点，恢复上述四个 HTML 到上一轮 3D 模型版本，并删除 `progress.md` 中本条记录。

## 2026-07-13 - Task: 新增 G 方案 3D 空间工作流首页候选
### What was done
- 保留 A-F 候选与生产首页不变，新增固定一屏的 G「3D 空间工作流」候选，以 Core、Route、Code、Meter 内容平面构成纵深路径，不再使用中心球体、棱镜或通用科技模型。
- 使用滚轮驱动相机逐段穿越、鼠标驱动观察角度，并让四段文案、强调色、空间平面和粒子氛围同步变化；Meter 镜头回退到完整构图，减弱动画模式仅保留当前章节平面。
- 将 G 加入统一候选入口，桌面保持六方案导航，手机端使用单行六等分导航。

### Testing
- 浏览器 `1440x900` 依次完成 `Core → Route → Code → Meter → Code → Route → Core`，阶段顺序正确且全过程 `scrollY=0`；相邻阶段平均像素差为 `6.93 / 8.33 / 11.37`，鼠标移动前后相机矩阵不同，控制台无脚本错误。
- 浏览器 `390x844` 下页面尺寸与视口同为 `390x844`，四个阶段按钮均完整位于视口内；滚轮可从 Core 切至 Route 并返回 Core，前后 `scrollY=0`，无横向或纵向溢出。
- 桌面与手机的全屏 Canvas 尺寸分别为 `1440x900`、`390x844`；场景空白区域像素采样均非空。统一入口在桌面和手机端均无溢出，点击 G 后 iframe 指向 `variant-g-spatial-workflow-v10.html`。
- G 与统一入口返回 HTTP 200，GSAP 与 Three.js CDN HEAD 请求返回 HTTP 200；模块脚本 `node --check`、UTF-8 严格解码、whitespace check、减弱动画当前章节静态断言和 `git diff --check` 通过。

### Notes
- `.superpowers/brainstorm/unified-core-20260712/content/variant-g-spatial-workflow-v10.html`：新增 CSS 3D 内容空间、GSAP 相机路径、Three.js 粒子氛围及桌面/手机交互。
- `.superpowers/brainstorm/unified-core-20260712/content/acd-variants-index-v8.html`：新增 G 方案入口并调整手机端六等分导航。
- `progress.md`：追加本轮 G 候选实现、验证、依赖与回滚记录。
- G 当前通过 CDN 加载 GSAP `3.13.0` 与 Three.js `0.180.0`，仍属于独立概念候选，未接入生产首页。
- 回滚点：以本条记录前的工作区状态为准；删除 G 候选文件，反向移除统一入口中的 G 配置、按钮和六等分样式，并删除 `progress.md` 中本条记录。

## 2026-07-13 - Task: 制作 H 光谱请求巨构主视觉与短动效样片
### What was done
- 保留 A-G 候选、统一入口和生产首页不变，新增独立 H「光谱请求巨构」样片，让请求线路、金属结构肋与实时数据脉冲共同构成超出屏幕的数据基础设施，不使用照片卡片、中心球体、棱镜或独立星空。
- 建立 7.2 秒自动循环的 Core、Route、Code、Meter 四种结构轮廓，分别呈现汇聚主干、分流穹顶、协议走廊和计量地形；同步切换文案、空间标签、相机路径与鼠标视差。
- 导出 Route 稳定帧作为 H 主视觉例图，便于在继续完整首页施工前单独评审视觉方向。

### Testing
- 浏览器 `1440x900` 验证 `Core → Route → Code → Meter → Core` 连续循环，四阶段 `scrollY=0`、Canvas 为 `1440x900`、控制台无脚本错误；右侧场景相邻阶段平均像素差为 `17.74 / 20.64 / 17.80`，有效画面覆盖率为 `29.3%–37.1%`。
- 浏览器 `390x844` 验证四阶段循环后返回 Core，页面尺寸始终为 `390x844`、`scrollY=0`、控制台无脚本错误；下半屏场景有效像素覆盖率为 `59.1%`，文案、操作按钮和四阶段标签均位于视口内。
- 桌面鼠标从中心移动到右上后，右侧场景平均像素差为 `15.58`，确认相机视差生效；主视觉例图尺寸为 `1440x900`，文件可正常解码。
- H 页面、Three.js 核心模块、EffectComposer、RenderPass 与 UnrealBloomPass 均返回 HTTP 200；模块脚本 `node --check`、UTF-8 严格解码、whitespace check、固定一屏与禁用中心基础几何静态断言、`git diff --check` 通过。

### Notes
- `.superpowers/brainstorm/unified-core-20260712/content/variant-h-spectral-infrastructure-study-v11.html`：新增全屏 Three.js 程序化请求巨构、Bloom、四阶段结构变形、同步文案和响应式样片。
- `.superpowers/brainstorm/unified-core-20260712/content/assets/h-spectral-v11/route-main-visual.png`：新增从实际 H 页面导出的 Route 主视觉例图。
- `progress.md`：追加本轮 H 样片实现、验证、依赖和回滚记录。
- H 通过 CDN 加载 Three.js `0.180.0` 及其后期处理模块，目前仅作为独立方向样片，未加入统一候选入口。
- 回滚点：以本条记录前的工作区状态为准；删除 H 样片文件与 `assets/h-spectral-v11/`，并删除 `progress.md` 中本条记录。

## 2026-07-13 - Task: 将 D 请求遥测方案设为正式首页
### What was done
- 将 D「请求遥测实景」接入 `/home` 默认分支，并继续保留管理员 `home_content` 的 HTML/iframe 覆盖、站点 Logo/文档地址清洗、鉴权初始化及普通用户与管理员控制台分流。
- 建立固定一屏的 Core、Route、Code、Meter 四阶段首页：本地实景素材与 Canvas 数据流同步切换，支持滚轮双向换相、阶段按钮、鼠标视差和点击脉冲，并在组件卸载时释放 RAF、定时器、媒体查询及全局事件监听。
- 补齐中英文首页文案、桌面/手机竖屏/手机横屏布局及减弱动画静态帧；修正顶栏与阶段栏堆叠顺序，避免语言下拉菜单被阶段控件覆盖。
- 更新公共页面文档和首页契约测试；A-G 候选与 H 独立样片保持原状，未删除或覆盖。

### Testing
- `pnpm exec vitest run src/views/__tests__/HomeTelemetryDefault.spec.ts src/components/layout/__tests__/siteLogoSanitization.spec.ts src/components/layout/__tests__/docUrlSanitization.spec.ts src/i18n/__tests__/localesNoKeyCollision.spec.ts`：4 个测试文件、20 项测试通过。
- `pnpm run typecheck`、`pnpm run lint:check`、`pnpm run build` 均通过；生产构建完成 909 个模块转换，保留仓库既有的 Browserslist、动态/静态导入和大分包提示。
- 浏览器在 `1440x900`、`390x844`、`844x390` 下确认根容器、Canvas 与文档尺寸等于视口，`scrollY=0` 且无横纵溢出；桌面完成 `Core -> Route -> Code -> Meter -> Code -> Route -> Core` 双向切换。
- 浏览器确认鼠标视差、点击脉冲、离开 `/home` 后 Canvas 卸载、返回首页后单次滚轮不重复触发；英文 Route 点击“中文”后恢复中文 Core，标题、导航、阶段说明和指标同步变化。

### Notes
- `frontend/src/views/HomeView.vue`：将 D 遥测组件接入默认首页分支并保留现有自定义首页及业务入口逻辑。
- `frontend/src/components/home/TelemetryHome.vue`：新增固定一屏的四阶段遥测首页、Canvas 动效、交互、响应式、减弱动画与资源清理，并修复语言菜单层级。
- `frontend/src/assets/home/telemetry/core.webp`：新增 Core 阶段本地实景素材。
- `frontend/src/assets/home/telemetry/route.webp`：新增 Route 阶段本地实景素材。
- `frontend/src/assets/home/telemetry/code.webp`：新增 Code 阶段本地实景素材。
- `frontend/src/assets/home/telemetry/meter.webp`：新增 Meter 阶段本地实景素材。
- `frontend/src/i18n/locales/zh/landing.ts`：新增 D 首页中文导航、阶段、动作、指标与辅助文本。
- `frontend/src/i18n/locales/en/landing.ts`：新增 D 首页英文导航、阶段、动作、指标与辅助文本。
- `frontend/src/views/__tests__/HomeTelemetryDefault.spec.ts`：新增默认分支、四阶段、资源清理、横屏布局和语言菜单层级契约。
- `docs/FRONTEND_PUBLIC_PAGES.md`：将公开首页说明更新为 D 请求遥测实景及其交互与响应式行为。
- `progress.md`：追加本轮生产首页接入、验证、风险修复和回滚记录。
- 回滚点：以本条记录前的工作区状态为准；恢复 `HomeView.vue`、两份 landing 语言包与公共页面文档到接入前状态，删除 `TelemetryHome.vue`、四张 telemetry WebP 和 `HomeTelemetryDefault.spec.ts`，并删除 `progress.md` 中本条记录。候选目录 `.superpowers/` 不参与回滚。

## 2026-07-13 - Task: 将 D 首页回归测试纳入关键测试集
### What was done
- 将 D 正式首页契约测试加入前端关键 Vitest 集合，使默认首页、生命周期、横屏布局和语言菜单层级回归进入现有 CI 验证路径。

### Testing
- `make test-frontend-critical` 通过：7 个测试文件、95 项测试全部通过；输出中的 Browserslist 提示及 SettingsView `router-link` 解析警告为关键测试集既有提示。

### Notes
- `Makefile`：将 `HomeTelemetryDefault.spec.ts` 加入 `FRONTEND_CRITICAL_VITEST`。
- `progress.md`：追加 D 首页测试进入 CI 关键集合的验证与回滚记录。
- 回滚点：以本条记录前的工作区状态为准；从 `FRONTEND_CRITICAL_VITEST` 删除 `src/views/__tests__/HomeTelemetryDefault.spec.ts`，并删除 `progress.md` 中本条记录。

## 2026-07-14 - Task: 确定 D 首页公共设置动态接入方案
### What was done
- 梳理公共设置中适合首页展示的数据，确认接入网站副标题、API Base URL、联系方式和用户可见的自定义菜单，不把自定义端点、版本号及后台管理配置塞入固定一屏首页。
- 制作三种信息落点对比并确认方案 1“信息就近”：副标题跟随品牌，自定义菜单进入顶部导航，API 地址与联系方式紧邻主行动按钮，现有四阶段遥测文案、主视觉和底部指标保持不变。
- 固化数据流、空值行为、菜单过滤排序、安全边界、响应式约束与验证标准，作为后续实现依据。

### Testing
- 布局对比页返回 HTTP 200，包含三个独立方案；浏览器 `1280x720` 下页面尺寸等于视口、三张主视觉均正常加载且默认选中方案 1。
- 设计文档占位符扫描通过，无 `TBD`、`TODO` 或待定实现项；设计文档与原型文件的 `git diff --check` 通过。

### Notes
- `.superpowers/brainstorm/unified-core-20260712/content/home-settings-layout-options-v12.html`：新增 D 首页四类公共设置的三方案布局对比页。
- `docs/superpowers/specs/2026-07-14-home-public-settings-design.md`：新增已确认方案 1 的数据范围、布局、交互、安全、响应式与验证设计。
- `progress.md`：追加本轮公共设置范围梳理、视觉方案确认与设计验证记录。
- 回滚点：以本条记录前的工作区状态为准；删除上述布局对比页和设计文档，并删除 `progress.md` 中本条记录，不影响现有 D 正式首页代码。

## 2026-07-14 - Task: 接入 D 首页公共动态设置

### What was done
- 将公共设置中的网站副标题、API Base URL、联系方式和用户可见自定义菜单接入 D 正式首页；空值不占位，自定义菜单按 `sort_order` 排序并过滤非用户可见项。
- 网站副标题跟随品牌展示，API 地址与联系方式紧邻主行动按钮；宽屏直接展示前两项自定义菜单并折叠其余项，`1200px` 及以下全部收进“更多入口”。
- 补齐折叠菜单的外部点击、`Escape` 和跳转后关闭行为，并修正 `1280x720` 下主按钮、文档入口与动态元信息相互挤压的问题。
- 保留四阶段遥测文案与指标的演示数据属性；自定义端点、版本号和管理员可见菜单不接入首页。

### Testing
- `pnpm exec vitest run src/components/home/__tests__/TelemetryHomePublicSettings.spec.ts src/views/__tests__/HomeTelemetryDefault.spec.ts src/components/layout/__tests__/siteLogoSanitization.spec.ts src/components/layout/__tests__/docUrlSanitization.spec.ts src/i18n/__tests__/localesNoKeyCollision.spec.ts`：5 个测试文件、25 项测试全部通过。
- `make test-frontend-critical`：8 个测试文件、100 项测试全部通过。
- `pnpm run typecheck`、`pnpm run lint:check`、`pnpm run build` 均通过；生产构建完成 909 个模块转换。
- 浏览器在 `1440x900`、`390x844`、`844x390`、`1280x720` 下验证固定一屏、`scrollY=0` 且无横纵溢出；桌面与移动端菜单折叠、中英文切换、Core 与 Route 双向滚轮切换及动态字段展示均符合设计，控制台无错误。

### Notes
- `Makefile`：将首页公共设置契约测试加入前端关键测试集。
- `docs/FRONTEND_PUBLIC_PAGES.md`：补充 D 首页公共设置数据范围、菜单规则、演示指标边界和响应式行为。
- `frontend/src/views/HomeView.vue`：从公共设置中提取首页所需动态字段并传入 D 首页组件。
- `frontend/src/components/home/TelemetryHome.vue`：展示动态副标题、API 地址、联系方式与用户菜单，并实现折叠交互和响应式布局。
- `frontend/src/components/home/__tests__/TelemetryHomePublicSettings.spec.ts`：新增动态字段、菜单过滤排序、折叠关闭及空值行为测试。
- `frontend/src/views/__tests__/HomeTelemetryDefault.spec.ts`：补充首页公共设置到 D 组件的传递契约测试。
- `frontend/src/i18n/locales/zh/landing.ts`：新增 API 地址、联系方式和更多入口的中文文案。
- `frontend/src/i18n/locales/en/landing.ts`：新增 API 地址、联系方式和更多入口的英文文案。
- `progress.md`：追加本轮实现、验证与回滚记录。
- 回滚点：以本条记录前的工作区状态为准；反向移除上述动态字段传递、展示、菜单交互、响应式规则及对应测试和文档，并从 `Makefile` 的关键测试集删除本轮新增测试。

## 2026-07-14 - Task: 发布 D 首页并合并原仓库 0.1.153 最新代码
### What was done
- 将已确认的 D 请求遥测首页、公共动态设置、四张本地 WebP、双语文案、关键测试和公开页文档独立提交；`.superpowers/` 本地概念稿、样片及预览服务文件继续排除在产品提交和 Release 之外。
- 合并原仓库 `Wei-Shaw/sub2api` 最新 `main`，上游正式版本为 `0.1.153`，最新提交为 `69bc6a87dde89e79ba39436467ec46dee6a6b234`。
- 解决 README、配置加载和计费服务 3 个冲突：保留本仓库二开标识、GitHub 更新 Token、Grok 4.5 模型与计费，同时接入上游视频编辑/扩展、Server-Timing 和 Grok Build 缓存计费修复。
- 接入上游 `0.1.152` 至 `0.1.153` 的数据库迁移、调度修复、OpenAI/Grok 网关能力、系统日志、静态资源缓存、支付与前端管理功能；本仓库更新目标、原仓库正式 Release 提示、应用内更新/重启和 Token 激励逻辑保持不变。

### Testing
- 首页提交前相关 Vitest：5 个测试文件、25 项测试通过；`pnpm run typecheck` 与 `pnpm run lint:check` 通过。
- `GOTOOLCHAIN=auto go test ./...`（在 `D:\project\sub2api-so\backend`）通过。
- `make test-frontend-critical`（在 `D:\project\sub2api-so`）通过，包含 D 首页默认分支与公共设置回归测试。
- `D:\environment\nodejs\node-v22.17.0-win-x64\node_modules\@pnpm\exe\pnpm.exe run build`（在 `D:\project\sub2api-so\frontend`）通过，完成 912 个模块转换并输出四张首页 WebP；仅有既有 Browserslist、动态/静态导入混用和 chunk size 警告。
- `git diff --check` 与 `git diff --cached --check` 通过；上游 `credentialsBuilder.spec.ts` 末尾空行已移除。

### Notes
- `frontend/src/views/HomeView.vue`、`frontend/src/components/home/TelemetryHome.vue`、`frontend/src/assets/home/telemetry/*`：发布 D 请求遥测正式首页及四阶段本地视觉素材。
- `frontend/src/components/home/__tests__/TelemetryHomePublicSettings.spec.ts`、`frontend/src/views/__tests__/HomeTelemetryDefault.spec.ts`、`Makefile`：发布首页公共设置、默认分支和关键测试集覆盖。
- `frontend/src/i18n/locales/en/landing.ts`、`frontend/src/i18n/locales/zh/landing.ts`、`docs/FRONTEND_PUBLIC_PAGES.md`：发布首页双语文案和公开行为文档。
- `README.md`：冲突处理同时保留二开标识、Grok 4.5 aliases 和上游视频编辑/扩展端点说明。
- `backend/internal/config/config.go`：冲突处理同时保留更新 Token 环境变量和上游 Server-Timing 环境变量。
- `backend/internal/service/billing_service.go`：冲突处理保留 Grok 4.5 兜底价格并接入上游 Grok Build 缓存价格说明。
- `frontend/src/components/account/__tests__/credentialsBuilder.spec.ts`：移除上游合并带入的文件末尾多余空行。
- `backend/cmd/server/VERSION`、`backend/cmd/server/UPSTREAM_COMMIT`：同步上游基线 `0.1.153` 和提交 `69bc6a87dde89e79ba39436467ec46dee6a6b234`。
- `backend/migrations/174_*`、`backend/migrations/175_*`、`backend/migrations/175a_*` 及迁移测试：同步用量长上下文计费、API Key 最新 IP 索引、Web Search 单次价格、系统日志 host 与默认计费开关迁移。
- `backend/**`（其余 255 个上游文件）：同步 `0.1.152` 至 `0.1.153` 的实体、网关、调度、计费、Grok/OpenAI、日志、支付、缓存和测试变更。
- `frontend/**`（其余 61 个上游文件）：同步账号、用量、设置、系统日志、Grok OAuth、日期格式和 API 客户端变更。
- `deploy/**`（12 个上游文件）：同步 Apple Container、Docker Compose、示例配置和部署测试。
- `.github/workflows/backend-ci.yml`、`.gitignore`、`README_CN.md`、`README_JA.md`：同步上游 CI、忽略规则和多语言说明。
- `progress.md`：追加本轮首页发布、上游合并、验证、冲突和回滚记录。
- 回滚方式：先执行 `git revert -m 1 <本次上游合并提交>` 回滚 `0.1.153`，再执行 `git revert 4e3d2f4a` 回滚 D 首页；数据库迁移已在生产执行时，应优先使用应用内旧版本回滚流程并确认旧程序兼容新增列，不直接删除迁移记录。
## 2026-07-16 - Task: Merge upstream v0.1.156 and prepare fork release
### What was done
- Merged `upstream/main` at `393a8fe56a0b606d162183cf8014f9381adcbf7e` into the Leo video feature branch while retaining Leo video, token incentive, fork update configuration, and upstream async image functionality.
- Resolved the ten merge conflicts by combining local fork behavior with upstream scheduler, gateway, account, audit, and object-storage changes.
- Kept homepage WIP out of the release; the local changes remain isolated in `wip-local-not-for-upstream-release-20260716` and `wip-progress-concurrent-not-for-upstream-release-20260716`.
- Regenerated Wire output and updated Leo callers for upstream scheduler and account-result API changes.
### Testing
- `go generate ./cmd/server`: passed.
- `git diff --check`: passed.
- `make test-frontend-critical`: passed, 8 files and 102 tests.
- Backend compilation and all non-flaky packages passed. The upstream `TestContentModerationRuntimeSnapshotRefreshFailureKeepsStaleConfig` test consistently times out on Windows/Go 1.26 when run with the package suite because its 1ns TTL async-refresh assertion is not scheduled within 1 second; an isolated run passed once, but the full suite remains red on this known upstream test.
### Notes
- `.gitignore`, `backend/cmd/server/wire_gen.go`, `backend/internal/handler/handler.go`, `backend/internal/handler/openai_gateway_handler.go`, `backend/internal/handler/wire.go`, `backend/internal/server/routes/gateway.go`, `backend/internal/service/admin_account.go`, `backend/internal/service/scheduler_snapshot_service.go`, and `deploy/config.example.yaml`: resolved upstream/local merge conflicts.
- `backend/internal/handler/leo_video.go` and `backend/internal/service/openai_account_scheduler_test.go`: aligned Leo calls with upstream scheduler and scheduling-result signatures.
- Remaining files in this merge are upstream changes from `upstream/main`; use `git diff --cached --name-only` to inspect the complete file list.
- Rollback point: revert the merge commit with `git revert -m 1 <merge-commit>` and remove tag `v0.1.156-fy.1` if the release must be withdrawn; do not restore the homepage WIP stashes into the release branch.

## 2026-07-16 - Task: Fix CI regressions after upstream v0.1.156 merge
### What was done
- Fixed the remaining CI issues introduced by the upstream `leo` platform: explicit nil-test returns, checked response-body closes, normalized Leo `/v1` URLs, deterministic content-moderation cache expiry setup, and dynamic scheduler/quota/API contract expectations.
- Kept homepage WIP changes out of the release branch; the working tree contains no frontend homepage changes.
### Testing
- `..\\.codex-run\\bin\\golangci-lint.exe run ./...` in `backend`: passed with `0 issues`.
- `go test -tags=unit ./...` in `backend`: passed for all packages.
- Scheduler targeted regressions and content-moderation refresh test (`-count=20`): passed.
- `git diff --check` and `git diff --cached --check`: passed.
### Notes
- `backend/internal/handler/admin/payment_handler_test.go`: added explicit return after the nil guard for lint correctness.
- `backend/internal/handler/admin/user_platform_quota_admin_test.go`: separated submitted quota count from all-platform cache invalidation count.
- `backend/internal/pkg/proxyurl/parse_test.go`: added explicit return after the nil guard for lint correctness.
- `backend/internal/server/api_contract_test.go`: included the Leo default platform quota in API contract fixtures.
- `backend/internal/service/account_test_service.go`: handled response-body close errors without changing behavior.
- `backend/internal/service/account_usage_service_test.go`: added explicit returns after nil guards.
- `backend/internal/service/content_moderation_runtime_cache_test.go`: made stale-snapshot setup deterministic on Windows.
- `backend/internal/service/gemini_messages_compat_service_test.go`: added an explicit return after the nil guard.
- `backend/internal/service/leo_account.go`: normalized the trailing slash before validating the Leo `/v1` path.
- `backend/internal/service/leo_video.go`: handled response-body close errors without changing behavior.
- `backend/internal/service/openai_gateway_service_codex_snapshot_test.go`: added explicit returns after nil guards.
- `backend/internal/service/openai_images_incomplete_test.go`: added explicit returns after nil guards.
- `backend/internal/service/ops_service_user_error_test.go`: added an explicit return after the nil guard.
- `backend/internal/service/ops_user_error_test.go`: added an explicit return after the nil guard.
- `backend/internal/service/payment_order_result_test.go`: added explicit returns after nil guards.
- `backend/internal/service/ratelimit_service_anthropic_test.go`: added explicit returns after nil guards.
- `backend/internal/service/scheduler_snapshot_full_rebuild_lifecycle_test.go`: made full-rebuild bucket and account-query expectations follow the current platform set.
- `backend/internal/service/scheduler_snapshot_group_lifecycle_test.go`: included Leo buckets and made lifecycle counts dynamic.
- `backend/internal/service/scheduler_snapshot_retirement_test.go`: made canonical capture counts dynamic.
- `progress.md`: recorded this CI repair and verification round.
- Rollback point: revert code commit `79af7885` with `git revert 79af7885`; do not restore the homepage WIP stashes.


## 2026-07-17 - Task: Merge upstream release v0.1.158 and prepare fork release
### What was done
- Merged official upstream tag `v0.1.158` at `26abd19a2812edba02bbef93c3e2a620141cc257` without including the ten later `upstream/main` commits.
- Preserved Leo video, token incentive, fork update/restart behavior, and the already released fork homepage; local homepage WIP stashes remain isolated and were not applied.
- Synced upstream group/channel-monitor duplication, bulk user limits, Grok endpoint fixes, Codex image/model fixes, and admin step-up 2FA hardening.
- Corrected source version metadata to `0.1.158` because the upstream tag predates its post-tag VERSION synchronization commit.
### Testing
- `go generate ./cmd/server`: passed; regenerated Wire output matched the merged provider graph.
- `..\\.codex-run\\bin\\golangci-lint.exe run --concurrency 2 ./...`: passed with `0 issues`.
- `go test -p 2 -tags=unit ./...`: passed for all backend packages.
- `go test -p 2 -tags=integration ./...`: passed for all backend packages.
- `make test-frontend-critical`: passed, 8 files and 102 tests.
- Upstream feature Vitest selection: passed, 13 files and 159 tests.
- Frontend `lint:check`, `typecheck`, and production build: passed; Vite built 930 modules.
- Local Apple-container shell validation was unavailable because the installed Windows shell is not GNU Bash; the GitHub macOS shell job remains the authoritative check.
- `git diff --cached --check`: passed.
### Notes
- `README.md`: synced the upstream sponsor entries for v0.1.158.
- `README_CN.md`: synced the upstream sponsor entries for v0.1.158.
- `README_JA.md`: synced the upstream sponsor entries for v0.1.158.
- `assets/partners/logos/claudeapi.jpg`: added the upstream sponsor logo asset.
- `assets/partners/logos/code0.jpg`: added the upstream sponsor logo asset.
- `backend/cmd/server/UPSTREAM_COMMIT`: recorded upstream tag commit 26abd19a.
- `backend/cmd/server/VERSION`: set the fork source baseline to 0.1.158.
- `backend/cmd/server/wire_gen.go`: regenerated dependency injection while retaining fork providers.
- `backend/ent/group.go`: synced upstream generated entity/schema changes for group duplication.
- `backend/ent/group/group.go`: synced upstream generated entity/schema changes for group duplication.
- `backend/ent/group/where.go`: synced upstream generated entity/schema changes for group duplication.
- `backend/ent/group_create.go`: synced upstream generated entity/schema changes for group duplication.
- `backend/ent/group_update.go`: synced upstream generated entity/schema changes for group duplication.
- `backend/ent/migrate/schema.go`: synced upstream generated entity/schema changes for group duplication.
- `backend/ent/mutation.go`: synced upstream generated entity/schema changes for group duplication.
- `backend/ent/runtime/runtime.go`: synced upstream generated entity/schema changes for group duplication.
- `backend/ent/schema/group.go`: synced upstream generated entity/schema changes for group duplication.
- `backend/internal/handler/admin/admin_basic_handlers_test.go`: synced upstream handler regression coverage.
- `backend/internal/handler/admin/admin_service_stub_test.go`: synced upstream handler regression coverage.
- `backend/internal/handler/admin/channel_monitor_duplicate_test.go`: synced upstream handler regression coverage.
- `backend/internal/handler/admin/channel_monitor_handler.go`: synced upstream handler behavior while retaining fork handlers.
- `backend/internal/handler/admin/grok_oauth_handler.go`: synced upstream handler behavior while retaining fork handlers.
- `backend/internal/handler/admin/grok_oauth_handler_test.go`: synced upstream handler regression coverage.
- `backend/internal/handler/admin/group_handler.go`: synced upstream handler behavior while retaining fork handlers.
- `backend/internal/handler/admin/group_handler_duplicate_test.go`: synced upstream handler regression coverage.
- `backend/internal/handler/admin/user_handler.go`: synced upstream handler behavior while retaining fork handlers.
- `backend/internal/handler/admin/user_handler_activity_test.go`: synced upstream handler regression coverage.
- `backend/internal/handler/admin/user_handler_batch_limits_test.go`: synced upstream handler regression coverage.
- `backend/internal/handler/admin/user_handler_get_deleted_test.go`: synced upstream handler regression coverage.
- `backend/internal/handler/admin/user_handler_list_apikey_group_test.go`: synced upstream handler regression coverage.
- `backend/internal/handler/admin/user_handler_role_stepup_test.go`: synced upstream handler regression coverage.
- `backend/internal/handler/auth_oauth_pending_flow_test.go`: synced upstream handler regression coverage.
- `backend/internal/handler/gateway_handler.go`: synced upstream handler behavior while retaining fork handlers.
- `backend/internal/handler/gateway_models_test.go`: synced upstream handler regression coverage.
- `backend/internal/handler/openai_codex_models_handler_test.go`: synced upstream handler regression coverage.
- `backend/internal/handler/totp_handler.go`: synced upstream handler behavior while retaining fork handlers.
- `backend/internal/handler/user_handler_test.go`: synced upstream handler regression coverage.
- `backend/internal/pkg/claude/constants.go`: synced upstream Claude/Grok helper behavior.
- `backend/internal/pkg/xai/oauth.go`: synced upstream Claude/Grok helper behavior.
- `backend/internal/pkg/xai/oauth_test.go`: synced upstream package regression coverage.
- `backend/internal/repository/api_key_repo.go`: synced upstream repository behavior for duplication and account updates.
- `backend/internal/repository/channel_monitor_duplicate_test.go`: synced upstream repository regression coverage.
- `backend/internal/repository/channel_monitor_repo.go`: synced upstream repository behavior for duplication and account updates.
- `backend/internal/repository/channel_monitor_template_duplicate_metadata_integration_test.go`: synced upstream repository regression coverage.
- `backend/internal/repository/channel_monitor_template_duplicate_metadata_unit_test.go`: synced upstream repository regression coverage.
- `backend/internal/repository/channel_monitor_template_repo.go`: synced upstream repository behavior for duplication and account updates.
- `backend/internal/repository/group_repo.go`: synced upstream repository behavior for duplication and account updates.
- `backend/internal/repository/group_repo_duplicate_integration_test.go`: synced upstream repository regression coverage.
- `backend/internal/repository/group_repo_integration_test.go`: synced upstream repository regression coverage.
- `backend/internal/repository/http_upstream.go`: synced upstream repository behavior for duplication and account updates.
- `backend/internal/repository/http_upstream_test.go`: synced upstream repository regression coverage.
- `backend/internal/repository/user_repo.go`: synced upstream repository behavior for duplication and account updates.
- `backend/internal/repository/user_repo_integration_test.go`: synced upstream repository regression coverage.
- `backend/internal/repository/wire.go`: synced upstream repository behavior for duplication and account updates.
- `backend/internal/server/api_contract_test.go`: synced upstream server contract and security coverage.
- `backend/internal/server/middleware/admin_auth_test.go`: synced upstream server contract and security coverage.
- `backend/internal/server/middleware/step_up.go`: synced upstream routes and step-up authorization behavior.
- `backend/internal/server/middleware/step_up_test.go`: synced upstream server contract and security coverage.
- `backend/internal/server/routes/admin.go`: synced upstream routes and step-up authorization behavior.
- `backend/internal/service/account.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/account_base_url_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/admin_group_duplicate.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/admin_group_duplicate_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/admin_service.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/admin_service_apikey_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/admin_service_batch_limits_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/admin_service_delete_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/admin_service_email_identity_sync_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/admin_user.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/auth_service_email_bind_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/channel_monitor_duplicate_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/channel_monitor_service.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/channel_monitor_types.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/content_moderation_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/gateway_anthropic_apikey_passthrough_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/gateway_claude_oauth_body.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/gateway_context_management_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/gateway_forward.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/gateway_upstream_request.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/grok_media.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/grok_upstream_url_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/group.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/group_service.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/openai_codex_models_service.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/openai_codex_models_service_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/openai_gateway_grok_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/openai_image_generation_controls_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/openai_ws_forwarder_ingress.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/openai_ws_forwarder_ingress_session_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/openai_ws_forwarder_success_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/openai_ws_forwarder_v2.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/openai_ws_http_bridge.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/openai_ws_v2_passthrough_adapter.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/totp_service.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/totp_verification_method_test.go`: synced upstream service and gateway regression coverage.
- `backend/internal/service/user_service.go`: synced upstream service, Grok, 2FA, and duplication behavior.
- `backend/internal/service/user_service_test.go`: synced upstream service and gateway regression coverage.
- `backend/migrations/181_group_duplicate_operation_id.sql`: added the upstream group-duplication operation-id migration.
- `docs/UPDATE_POLICY.md`: documented the v0.1.158 baseline and exclusion of post-tag main commits.
- `frontend/src/api/__tests__/admin.channelMonitor.duplicate.spec.ts`: synced upstream frontend API regression coverage.
- `frontend/src/api/__tests__/admin.groups.duplicate.spec.ts`: synced upstream frontend API regression coverage.
- `frontend/src/api/__tests__/admin.users.spec.ts`: synced upstream frontend API regression coverage.
- `frontend/src/api/admin/channelMonitor.ts`: synced upstream admin API contracts.
- `frontend/src/api/admin/groups.ts`: synced upstream admin API contracts.
- `frontend/src/api/admin/users.ts`: synced upstream admin API contracts.
- `frontend/src/components/account/BulkEditAccountModal.vue`: synced upstream account, admin, table, and Grok controls.
- `frontend/src/components/account/CreateAccountModal.vue`: synced upstream account, admin, table, and Grok controls.
- `frontend/src/components/account/EditAccountModal.vue`: synced upstream account, admin, table, and Grok controls.
- `frontend/src/components/account/GrokBaseUrlPresets.vue`: synced upstream account, admin, table, and Grok controls.
- `frontend/src/components/account/__tests__/BulkEditAccountModal.spec.ts`: synced upstream component regression coverage.
- `frontend/src/components/account/__tests__/EditAccountModal.grokUpstream.spec.ts`: synced upstream component regression coverage.
- `frontend/src/components/account/__tests__/EditAccountModal.spec.ts`: synced upstream component regression coverage.
- `frontend/src/components/account/__tests__/credentialsBuilder.spec.ts`: synced upstream component regression coverage.
- `frontend/src/components/account/credentialsBuilder.ts`: synced upstream account, admin, table, and Grok controls.
- `frontend/src/components/admin/monitor/MonitorActionsCell.spec.ts`: synced upstream component regression coverage.
- `frontend/src/components/admin/monitor/MonitorActionsCell.vue`: synced upstream account, admin, table, and Grok controls.
- `frontend/src/components/admin/user/BulkEditUserModal.vue`: synced upstream account, admin, table, and Grok controls.
- `frontend/src/components/admin/user/UserCreateModal.vue`: synced upstream account, admin, table, and Grok controls.
- `frontend/src/components/admin/user/UserEditModal.vue`: synced upstream account, admin, table, and Grok controls.
- `frontend/src/components/admin/user/__tests__/BulkEditUserModal.spec.ts`: synced upstream component regression coverage.
- `frontend/src/components/common/DataTable.vue`: synced upstream account, admin, table, and Grok controls.
- `frontend/src/components/common/__tests__/DataTable.spec.ts`: synced upstream component regression coverage.
- `frontend/src/components/keys/UseKeyModal.vue`: synced upstream account, admin, table, and Grok controls.
- `frontend/src/components/keys/__tests__/UseKeyModal.spec.ts`: synced upstream component regression coverage.
- `frontend/src/i18n/locales/en/admin/accounts.ts`: synced upstream admin and dashboard translations.
- `frontend/src/i18n/locales/en/admin/channels.ts`: synced upstream admin and dashboard translations.
- `frontend/src/i18n/locales/en/admin/overview.ts`: synced upstream admin and dashboard translations.
- `frontend/src/i18n/locales/en/admin/resources.ts`: synced upstream admin and dashboard translations.
- `frontend/src/i18n/locales/en/dashboard.ts`: synced upstream admin and dashboard translations.
- `frontend/src/i18n/locales/zh/admin/accounts.ts`: synced upstream admin and dashboard translations.
- `frontend/src/i18n/locales/zh/admin/channels.ts`: synced upstream admin and dashboard translations.
- `frontend/src/i18n/locales/zh/admin/overview.ts`: synced upstream admin and dashboard translations.
- `frontend/src/i18n/locales/zh/admin/resources.ts`: synced upstream admin and dashboard translations.
- `frontend/src/i18n/locales/zh/dashboard.ts`: synced upstream admin and dashboard translations.
- `frontend/src/types/index.ts`: synced upstream frontend API types.
- `frontend/src/views/admin/AuditLogView.vue`: synced upstream admin views for duplication, audit, and batch limits.
- `frontend/src/views/admin/ChannelMonitorView.vue`: synced upstream admin views for duplication, audit, and batch limits.
- `frontend/src/views/admin/GroupsView.vue`: synced upstream admin views for duplication, audit, and batch limits.
- `frontend/src/views/admin/UsersView.vue`: synced upstream admin views for duplication, audit, and batch limits.
- `frontend/src/views/admin/__tests__/ChannelMonitorView.duplicate.spec.ts`: synced upstream admin-view regression coverage.
- `frontend/src/views/admin/__tests__/GroupsView.duplicate.spec.ts`: synced upstream admin-view regression coverage.
- `frontend/src/views/admin/__tests__/UsersView.spec.ts`: synced upstream admin-view regression coverage.
- `progress.md`: recorded the v0.1.158 merge, verification evidence, changed-file list, and rollback point.
- Rollback point: deploy release `v0.1.156-fy.2`; revert this source merge with `git revert -m 1 <v0.1.158 merge commit>`. Do not apply or delete the homepage WIP stashes during rollback.
