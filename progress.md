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
