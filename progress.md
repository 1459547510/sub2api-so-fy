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


## 2026-07-17 - Task: Merge upstream release v0.1.159 and prepare fork release
### What was done
- Merged official upstream tag `v0.1.159` at `2a75d7d2387587d86ca3c5e5cd8ca96cf3d104c6` after explicit approval for its authentication and client-IP security changes.
- Synced the API-key alpha search scheduling fix, trusted forwarded client-IP handling, Grok Free function-tool cache routing, account upstream-site links, and lazy Stripe loading.
- Preserved Leo video, token incentive, and fork update/restart behavior; local homepage WIP stashes remained isolated and were not applied.
- Corrected source version metadata to `0.1.159` because the upstream tag still contains the prior `0.1.158` VERSION value.
### Testing
- `go generate ./cmd/server`: passed; regenerated Wire output had no diff and retained the fork providers.
- `..\\.codex-run\\bin\\golangci-lint.exe run --concurrency 2 ./...`: passed with `0 issues`.
- `go test -p 2 -tags=unit ./...`: passed for all backend packages.
- `go test -p 2 -tags=integration ./...`: passed for all backend packages.
- v0.1.159 frontend Vitest selection: passed, 3 files and 13 tests.
- `make test-frontend-critical`: passed, 8 files and 102 tests.
- Frontend ESLint, `vue-tsc --noEmit`, build typecheck, and Vite production build: passed; Vite built 930 modules.
- `git diff --check`: passed; final cached checks run before commit.
### Notes
- `backend/cmd/server/UPSTREAM_COMMIT`: recorded upstream tag commit `2a75d7d2`.
- `backend/cmd/server/VERSION`: set the fork source baseline to `0.1.159`.
- `backend/internal/handler/admin/audit_log_handler.go`: routed audit client-IP resolution through the shared trusted-forwarding policy.
- `backend/internal/pkg/ip/ip.go`: added shared request client-IP resolution for trusted proxy deployments.
- `backend/internal/pkg/ip/ip_test.go`: added trusted and untrusted forwarded-IP regression coverage.
- `backend/internal/server/middleware/api_key_auth.go`: reused the shared client-IP resolver for API-key ACL checks.
- `backend/internal/server/middleware/api_key_auth_google.go`: reused the shared client-IP resolver for Google-compatible API-key ACL checks.
- `backend/internal/server/middleware/audit_log.go`: aligned audit-log IP capture with the shared trust setting.
- `backend/internal/server/middleware/session_binding.go`: aligned session IP binding with trusted forwarded client IPs.
- `backend/internal/server/middleware/session_binding_test.go`: added session-binding regression coverage for proxy trust modes.
- `backend/internal/server/router.go`: wired the shared trusted-client-IP policy into middleware construction.
- `backend/internal/service/account.go`: exposed API-key account upstream-site metadata used by the admin link.
- `backend/internal/service/openai_account_scheduler_test.go`: covered API-key scheduling for OpenAI alpha search.
- `backend/internal/service/openai_alpha_search.go`: restored API-key account scheduling and compatible failover for alpha search.
- `backend/internal/service/openai_alpha_search_test.go`: added pure and mixed-group alpha-search regressions.
- `backend/internal/service/openai_gateway_grok.go`: enabled Grok Free cache routing for Responses function tools.
- `backend/internal/service/openai_gateway_grok_cache.go`: handled function-tool and built-in web-search name conflicts in cache routing.
- `backend/internal/service/openai_images_test.go`: updated image gateway regression fixtures for current routing behavior.
- `backend/internal/service/openai_ws_http_bridge.go`: applied Grok Free function-tool cache routing to the WebSocket bridge.
- `docs/UPDATE_POLICY.md`: advanced the documented upstream baseline to v0.1.159.
- `frontend/src/components/payment/StripePaymentInline.vue`: switched Stripe SDK loading to the lazy path.
- `frontend/src/i18n/locales/en/admin/accounts.ts`: added English account upstream-link labels.
- `frontend/src/i18n/locales/en/admin/settings.ts`: updated English trusted-client-IP setting guidance.
- `frontend/src/i18n/locales/zh/admin/accounts.ts`: added Chinese account upstream-link labels.
- `frontend/src/i18n/locales/zh/admin/settings.ts`: updated Chinese trusted-client-IP setting guidance.
- `frontend/src/views/admin/AccountsView.vue`: linked API-key account names to their configured upstream sites.
- `frontend/src/views/admin/__tests__/AccountsView.sparkShadow.spec.ts`: added account upstream-link view coverage.
- `frontend/src/views/user/StripePaymentView.vue`: lazy-loaded the Stripe SDK on the dedicated payment page.
- `frontend/src/views/user/StripePopupView.vue`: lazy-loaded the Stripe SDK in popup payment flow.
- `frontend/src/views/user/__tests__/StripePaymentView.spec.ts`: updated Stripe payment view regression setup.
- `frontend/src/views/user/__tests__/stripeLazyLoading.spec.ts`: added Stripe lazy-loading regressions.
- `frontend/vite.config.ts`: split the on-demand Stripe dependency into its own vendor chunk.
- `progress.md`: recorded the v0.1.159 merge, verification evidence, changed-file list, and rollback point.
- Rollback point: deploy release `v0.1.158-fy.1`; revert this source merge with `git revert -m 1 <v0.1.159 merge commit>`. Do not apply or delete the homepage WIP stashes during rollback.


## 2026-07-17 - Task: Fix false fork branch update notice and missing version translations
### What was done
- Fixed the update check so a fork default branch behind the running Release is not reported as a newer branch commit.
- Preserved branch notices when the branch is actually ahead or has commits absent from the running build.
- Added the missing Chinese and English text for fork Release, upstream Release, synchronization, and branch update states.
- Documented that branch update notices require the default branch to be ahead of the running version.
### Testing
- Added a backend regression test and confirmed it failed before the fix because `HasNewCommit` was incorrectly true.
- `go test -tags=unit -run TestUpdateService ./internal/service`: passed.
- `go test -p 2 -tags=unit ./internal/service`: passed.
- `..\\.codex-run\\bin\\golangci-lint.exe run --concurrency 2 ./...`: passed with `0 issues`.
- VersionBadge locale-key presence check: passed for all 9 added keys in Chinese and English.
- `localesNoKeyCollision.spec.ts`: passed, 6 tests.
- `localesMessageCompile.spec.ts` could not collect because the existing frontend manifest does not declare `@intlify/message-compiler` as a direct dependency; no dependency changes were made for this scoped fix.
- Frontend ESLint, `vue-tsc --noEmit`, build typecheck, and Vite production build: passed; Vite built 930 modules.
- `git diff --check`: passed before the documentation and progress-log append; final checks follow.
### Notes
- `backend/internal/service/update_service.go`: compares the running commit with the fork branch before setting the branch-update flag.
- `backend/internal/service/update_service_test.go`: covers a fork branch that is behind the installed Release commit.
- `frontend/src/i18n/locales/en/misc.ts`: added English VersionBadge update-state messages.
- `frontend/src/i18n/locales/zh/misc.ts`: added Chinese VersionBadge update-state messages.
- `docs/UPDATE_POLICY.md`: documented directional branch-update detection.
- `progress.md`: recorded the false-update fix, translation completion, verification, and rollback point.
- Rollback point: revert these uncommitted files before release, or revert the eventual fix commit with `git revert <fix commit>`.


## 2026-07-17 - Task: 设计 Leo 异步视频工作台
### What was done
- 确定用户端采用双栏视频工作台，支持纯文本、图片 URL、本地图片和多任务异步队列。
- 确定 Sub2API 持久化 LeoStudio 任务映射、后台同步状态，并通过冻结、结算和释放保证多任务只计费一次。
- 确定同宿主机本地磁盘图片方案、回环临时 URL、文件清理策略、API Key 隔离和数据库边界。
### Testing
- 对照 LeoStudio `f822735629c51f15d115e3e60b161a93ec2e20ff` 的异步创建、查询、取消和图片 URL 契约完成逐项复核。
- 规格占位符、矛盾、范围和歧义检查通过；明确了创建响应不确定时不重提、标记失败并释放冻结的 at-most-once 策略。
### Notes
- `docs/superpowers/specs/2026-07-17-leo-async-video-workbench-design.md`: 新增已批准范围下的异步协议、数据模型、计费、本地图片和客户端设计。
- `progress.md`: 追加本轮设计决策、验证与回滚记录。
- 回滚点：删除上述设计文档，并删除 `progress.md` 中本条记录；尚未修改数据库、后端协议或客户端代码。


## 2026-07-17 - Task: 制定 Leo 异步视频工作台实施计划
### What was done
- 将已确认设计拆分为数据模型、异步协议、账务、本地图片、后台协调、网关和客户端十个可独立验证的 TDD 任务。
- 固定每个任务的文件边界、失败测试、最小实现、验证命令和提交点，并保留同步 Leo 回归检查。
### Testing
- 逐项对照设计规格复核任务覆盖，包含数据库、重启恢复、多任务、三类图片输入、幂等计费、权限、客户端和文档。
- 计划占位符、类型命名和路径一致性检查通过；中英文导航与工作台文案文件已精确定位。
### Notes
- `docs/superpowers/plans/2026-07-17-leo-async-video-workbench.md`: 新增完整实施计划。
- `progress.md`: 追加本轮计划拆分、验证与回滚记录。
- 回滚点：删除上述实施计划，并删除 `progress.md` 中本条记录；设计提交 `a12565b3` 不受影响。


## 2026-07-19 - Task: Merge upstream release v0.1.161 with database and authentication changes
### What was done
- Merged official upstream tag `v0.1.161` at `19149ca196eeae4a4482e5299dc6fa4ba0b06c8c` after explicit approval for its database migrations and authentication changes.
- Accepted migrations `181_prompt_audit.sql` through `184_auth_cache_invalidation_outbox.sql`; verified duplicate `181_*` files are tracked and executed by full filename.
- Preserved Token Incentive, Leo video, fork update/restart behavior, and the fork container image target.
- Resolved the semantic Leo conflict introduced when upstream replaced legacy content moderation with the unified security-audit coordinator.
- Kept the homepage WIP stashes unapplied and left `.superpowers/` untracked.
### Testing
- `go generate ./cmd/server`: passed; generated Wire output matched the merge and retained Token Incentive, Leo, prompt-audit, and auth-cache invalidation providers.
- Migration/auth audit: confirmed full-filename migration identity, v0.1.161 defaults (`prompt audit`, `step-up 2FA`, and `session IP/UA binding` all off), nullable security-setting updates, and SHA-256-only auth outbox keys.
- Core backend package tests passed, including repository, admin handler, securityaudit, middleware, routes, handler, service, and server packages.
- `go test -p 2 -tags=unit ./...`: all packages passed except one pre-existing 500 ms timing-wheel test flaked under concurrent load; the unchanged test then passed 10/10 in isolation and the complete `internal/service` package passed with `-parallel 1`.
- `go test -tags=integration -run '^(TestMigrationsRunner_|TestAuthCacheInvalidationTriggers_)' -v -timeout 2m ./internal/repository`: exited successfully but explicitly skipped because Docker is unavailable on this Windows host; CI remains the real PostgreSQL/Redis integration gate.
- `..\\.codex-run\\bin\\golangci-lint.exe run --concurrency 2 ./...`: passed with `0 issues`.
- Frontend `lint:check` and `vue-tsc --noEmit`: passed.
- Full frontend Vitest suite: passed after restoring the frozen pnpm 10 lockfile dependencies; the previously missing i18n compiler suite also passed 2/2 in isolation.
- Frontend production build: passed; Vite transformed 949 modules and generated embedded assets.
- `go build -o ..\\.codex-run\\sub2api-v0.1.161-fy.1.exe ./cmd/server` and `--version`: passed and reported `Sub2API 0.1.161`.
- `git diff --check` passed for source; the immutable upstream source-freeze `.patch` is excluded from cached whitespace checks because it intentionally preserves historical patch whitespace.
### Notes
- `Dockerfile`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `README.md`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `assets/partners/logos/code0.jpg`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/cmd/cleanup-ingress-reject-logs/README.md`: added the upstream v0.1.161 implementation or verification asset.
- `backend/cmd/cleanup-ingress-reject-logs/main.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/cmd/cleanup-ingress-reject-logs/main_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/cmd/server/UPSTREAM_COMMIT`: recorded upstream v0.1.161 commit 19149ca196eeae4a4482e5299dc6fa4ba0b06c8c.
- `backend/cmd/server/VERSION`: set the fork source baseline to 0.1.161.
- `backend/cmd/server/main.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/cmd/server/wire.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/cmd/server/wire_gen.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/cmd/server/wire_gen_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/config/config.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/config/config_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/account_handler.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/account_handler_mixed_channel_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/admin_basic_handlers_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/admin_service_stub_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/grok_import_probe.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/grok_import_probe_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/grok_oauth_handler.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/grok_oauth_handler_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/ops_auth_cache_health_handler.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/handler/admin/ops_ingress_reject_handler.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/handler/admin/ops_ingress_reject_handler_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/handler/admin/setting_handler.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/setting_handler_audit.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/setting_handler_stepup_switch_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/handler/admin/setting_handler_update.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/user_handler.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/user_handler_activity_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/user_handler_batch_limits_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/user_handler_get_deleted_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/user_handler_list_apikey_group_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/admin/user_handler_role_stepup_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/batch_image_handler.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/content_moderation_helper.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/dto/settings.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/gateway_handler.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/gateway_handler_chat_completions.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/gateway_handler_responses.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/gemini_v1beta_handler.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/grok_media.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/grok_media_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/handler.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/image_task_handler.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/no_account_error.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/no_account_error_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/openai_alpha_search.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/openai_chat_completions.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/openai_embeddings.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/openai_gateway_count_tokens.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/openai_gateway_credential_failover_loop_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/openai_gateway_credential_failover_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/openai_gateway_handler.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/openai_gateway_handler_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/openai_grok_image_intent_gate_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/openai_images.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/ops_error_logger.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/ops_error_logger_attribution_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/ops_error_logger_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/handler/ops_ingress_reject_capture_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/handler/security_audit_errors.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/handler/security_audit_errors_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/handler/security_audit_helper.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/handler/security_audit_helper_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/handler/security_audit_media_submit_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/handler/security_audit_order_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/handler/wire.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/pkg/apicompat/anthropic_to_responses_response.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/pkg/apicompat/anthropic_to_responses_stream_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/pkg/ip/ip.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/pkg/ip/ip_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/pkg/logger/logger.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/pkg/xai/billing.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/account_repo.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/account_repo_model_availability_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/repository/account_repo_sort_integration_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/account_repo_upstream_billing_probe_update_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/api_key_cache.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/api_key_cache_subscriber_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/repository/api_key_repo.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/api_key_repo_integration_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/auth_cache_invalidation_outbox_integration_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/repository/auth_cache_invalidation_outbox_repo.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/repository/auth_cache_invalidation_outbox_repo_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/repository/http_upstream.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/http_upstream_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/migrations_runner.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/migrations_runner_notx_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/migrations_schema_integration_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/ops_error_where_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/ops_ingress_reject_repo.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/repository/ops_ingress_reject_repo_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/repository/ops_repo.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/ops_repo_args_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/ops_repo_get_error_log_by_id_integration_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/ops_repo_lookup_deleted_key_audit_integration_test.go`: synced the upstream v0.1.161 removal.
- `backend/internal/repository/scheduler_cache.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/scheduler_cache_unit_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/user_repo.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/user_repo_delete_atomicity_integration_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/repository/wire.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/securityaudit/coordinator.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/coordinator_legacy.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/coordinator_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_config.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_config_integration_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_config_store.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_config_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_enqueue.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_event_repository.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_guard.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_guard_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_handler.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_handler_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_issue_summary.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_logging.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_logging_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_metrics.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_metrics_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_module.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_outbound_security.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_outbound_security_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_payload_store.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_payload_store_integration_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_qwen3guard.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_qwen3guard_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_repository.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_repository_integration_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_repository_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_scanner.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_service.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_service_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_snapshot.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_snapshot_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_types.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_worker.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/securityaudit/prompt_worker_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/server/api_contract_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/http.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/http_ingress_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/server/middleware/api_key_auth.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/api_key_auth_google.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/api_key_auth_google_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/api_key_auth_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/audit_log.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/audit_log_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/client_request_id.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/client_request_id_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/ingress_reject.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/server/middleware/ingress_reject_access_sampler.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/server/middleware/ingress_reject_access_sampler_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/server/middleware/ingress_reject_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/server/middleware/invalid_auth_abuse_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/server/middleware/logger.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/middleware.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/request_access_logger_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/request_logger.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/request_metadata.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/server/middleware/session_binding.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/session_binding_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/step_up.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/middleware/step_up_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/router.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/routes/admin.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/routes/gateway.go`: merged upstream text-body limits and Grok video content routing while retaining Leo unsupported-endpoint guards.
- `backend/internal/server/routes/gateway_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/server/routes/ops_ingress_reject_routes_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/server/routes/prompt_audit_route_coverage_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/account.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/account_grok_media_eligibility_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/account_service.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/account_service_delete_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/account_test_service.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/account_test_service_grok_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/admin_account.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/admin_account_upstream_billing_probe_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/admin_service.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/antigravity_gateway_service.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/antigravity_subscription_service.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/antigravity_subscription_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/api_key_auth_cache_impl.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/api_key_service.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/api_key_service_cache_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/auth_cache_invalidation_outbox.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/auth_cache_invalidation_outbox_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/channel_monitor_checker.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/channel_monitor_checker_body_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/domain_constants.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/error_policy_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/gateway_model_availability.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/gateway_model_availability_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/gateway_multiplatform_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/gateway_non_streaming_response_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/gateway_request.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/gateway_request_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/gemini_chat_completions_compat_service.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/gemini_error_policy_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/gemini_messages_compat_service.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/gemini_multiplatform_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/grok_media.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/grok_media_content_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/grok_quota_fetcher.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/grok_quota_fetcher_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/grok_quota_service.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/grok_quota_service_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/grok_upstream_url.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/grok_upstream_url_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/image_generation_intent.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/image_generation_intent_explicit_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/invalid_auth_abuse_limiter.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/invalid_auth_abuse_limiter_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/openai_account_runtime_block_fastpath.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_account_runtime_block_fastpath_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_account_scheduler_test.go`: retained Leo platform isolation coverage and added upstream Grok media eligibility coverage.
- `backend/internal/service/openai_alpha_search.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_embeddings.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_gateway_cc_pipeline.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_gateway_chat_completions.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_gateway_chat_completions_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_gateway_forward.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_gateway_grok.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_gateway_grok_cache.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_gateway_grok_cache_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_gateway_grok_cache_tool_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/openai_gateway_grok_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_gateway_model_availability.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_gateway_passthrough.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_gateway_scheduling.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_gateway_upstream_errors.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_images.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_images_responses.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_stream_read_error.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/openai_ws_client_read.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_ws_v2/passthrough_relay.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_ws_v2/passthrough_relay_internal_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_ws_v2/passthrough_relay_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_ws_v2_passthrough_adapter.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/openai_ws_v2_passthrough_lifecycle_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/ops_cleanup_executor.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ops_cleanup_service.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ops_ingress_reject.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/ops_ingress_reject_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/ops_log_runtime_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ops_models.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ops_port.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ops_queue_sanitize_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/ops_repo_mock_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ops_runtime_snapshot_test.go`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/service/ops_service.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ops_service_user_error_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ops_settings.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ops_settings_advanced_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ops_settings_models.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ops_system_log_sink.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ops_system_log_sink_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ops_upstream_context.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ratelimit_service.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ratelimit_service_model_not_found_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/ratelimit_session_window_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/scheduler_snapshot_batch_query_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/scheduler_snapshot_full_rebuild_lifecycle_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/setting_features.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/setting_parse.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/setting_update.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/settings_view.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/subscription_assign_idempotency_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/subscription_service.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/upstream_billing_probe.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/upstream_billing_probe_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/user_subscription_daily_quota_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/service/wire.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/web/embed_on.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/internal/web/embed_test.go`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `backend/migrations/181_prompt_audit.sql`: added the upstream v0.1.161 implementation or verification asset.
- `backend/migrations/182_prompt_audit_full_prompt.sql`: added the upstream v0.1.161 implementation or verification asset.
- `backend/migrations/183_ops_ingress_reject_aggregates.sql`: added the upstream v0.1.161 implementation or verification asset.
- `backend/migrations/184_auth_cache_invalidation_outbox.sql`: added the upstream v0.1.161 implementation or verification asset.
- `backend/scripts/finalize-ingress-reject-cleanup.sql`: added the upstream v0.1.161 implementation or verification asset.
- `deploy/Caddyfile`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `deploy/EDGE_SECURITY.md`: added the upstream v0.1.161 implementation or verification asset.
- `deploy/README.md`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `deploy/build_image.sh`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `deploy/config.example.yaml`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `deploy/docker-compose.yml`: retained the fork image target and synced the upstream Redis persistence command fix.
- `docs/UPDATE_POLICY.md`: advanced the documented upstream baseline to v0.1.161 and retained fork-first release policy.
- `frontend/package.json`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/pnpm-lock.yaml`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/App.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/api/admin/ops.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/api/admin/settings.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/components/account/AccountUsageCell.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/components/account/BulkEditAccountModal.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/components/account/CreateAccountModal.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/components/account/UpstreamBillingRateCell.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/components/account/__tests__/AccountUsageCell.spec.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/components/account/__tests__/BulkEditAccountModal.spec.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/components/account/__tests__/CreateAccountModal.spec.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/components/account/__tests__/UpstreamBillingRateCell.spec.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/components/common/DataTable.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/components/common/__tests__/DataTable.spec.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/components/layout/AppSidebar.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/features/prompt-audit/PromptAuditView.vue`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/features/prompt-audit/__tests__/PromptAuditView.spec.ts`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/features/prompt-audit/__tests__/api.spec.ts`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/features/prompt-audit/__tests__/components.spec.ts`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/features/prompt-audit/__tests__/integrationSurface.spec.ts`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/features/prompt-audit/__tests__/viewModel.spec.ts`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/features/prompt-audit/api.ts`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/features/prompt-audit/components/EndpointPool.vue`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/features/prompt-audit/components/EventDetailDialog.vue`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/features/prompt-audit/components/EventWorkspace.vue`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/features/prompt-audit/components/FilterDeleteDialog.vue`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/features/prompt-audit/components/PolicyPanel.vue`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/features/prompt-audit/components/RuntimeOverview.vue`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/features/prompt-audit/types.ts`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/features/prompt-audit/viewModel.ts`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/i18n/__tests__/wsModeLocaleDesc.spec.ts`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/i18n/locales/en/admin/accounts.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/i18n/locales/en/admin/index.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/i18n/locales/en/admin/ops.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/i18n/locales/en/admin/promptAudit.ts`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/i18n/locales/en/admin/settings.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/i18n/locales/en/common.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/i18n/locales/en/misc.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/i18n/locales/zh/admin/accounts.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/i18n/locales/zh/admin/index.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/i18n/locales/zh/admin/ops.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/i18n/locales/zh/admin/promptAudit.ts`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/i18n/locales/zh/admin/settings.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/i18n/locales/zh/common.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/i18n/locales/zh/misc.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/main.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/router/index.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/style.css`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/types/index.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/utils/__tests__/branding.spec.ts`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/utils/branding.ts`: added the upstream v0.1.161 implementation or verification asset.
- `frontend/src/views/admin/AccountsView.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/views/admin/BackupView.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/views/admin/SettingsView.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/views/admin/__tests__/AccountsView.bulkEdit.spec.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/views/admin/__tests__/AccountsView.usageWindowsHint.spec.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/views/admin/__tests__/SettingsView.spec.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/views/admin/ops/components/OpsErrorDetailModal.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/views/admin/ops/components/OpsErrorLogTable.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/views/admin/ops/components/OpsSettingsDialog.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/views/admin/orders/AdminPaymentPlansView.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/views/admin/orders/PlanEditDialog.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/src/views/public/LegalDocumentView.vue`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `frontend/vite.config.ts`: synced upstream v0.1.161 behavior while preserving fork-owned features.
- `openspec/changes/add-openai-compatible-prompt-audit/.openspec.yaml`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/changes/add-openai-compatible-prompt-audit/README.md`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/changes/add-openai-compatible-prompt-audit/design.md`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/changes/add-openai-compatible-prompt-audit/implementation-evidence.md`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/changes/add-openai-compatible-prompt-audit/implementation-guide.md`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/changes/add-openai-compatible-prompt-audit/proposal.md`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/changes/add-openai-compatible-prompt-audit/source-baseline.md`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/changes/add-openai-compatible-prompt-audit/source-feature-map.md`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/changes/add-openai-compatible-prompt-audit/source-freeze/MANIFEST.md`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/changes/add-openai-compatible-prompt-audit/source-freeze/aicodex-prompt-audit-tracked.patch`: added the immutable upstream prompt-audit source-freeze patch; its historical whitespace is intentionally preserved.
- `openspec/changes/add-openai-compatible-prompt-audit/source-freeze/aicodex-prompt-audit-untracked.tar.gz`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/changes/add-openai-compatible-prompt-audit/specs/prompt-input-audit/spec.md`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/changes/add-openai-compatible-prompt-audit/specs/prompt-input-guard/spec.md`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/changes/add-openai-compatible-prompt-audit/specs/security-audit-console/spec.md`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/changes/add-openai-compatible-prompt-audit/tasks.md`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/changes/add-openai-compatible-prompt-audit/verification.md`: added the upstream v0.1.161 implementation or verification asset.
- `openspec/config.yaml`: added the upstream v0.1.161 implementation or verification asset.
- `backend/internal/handler/leo_video.go`: moved the fork Leo video path to the v0.1.161 unified security-audit coordinator before scheduling.
- `backend/internal/handler/leo_video_test.go`: added regression coverage proving a blocked Leo prompt never reaches account scheduling.
- `progress.md`: recorded the v0.1.161 merge, verification evidence, changed-file list, exclusions, and rollback point.
- Rollback point: deploy release `v0.1.159-fy.1`; revert this source merge with `git revert -m 1 <v0.1.161 merge commit>`. Do not apply or delete the homepage WIP stashes during rollback.

## 2026-07-17 - Task: Persist Leo asynchronous video jobs
### What was done
- Added the durable `video_jobs` model for Sub2API-owned public IDs, Leo account affinity, request snapshots, lifecycle state, results, and billing state.
- Added API-key-scoped reads and lists plus atomic allowed-state transitions so concurrent workers cannot advance the same job from a stale state.
- Registered the repository for dependency injection and generated the Ent access layer.
### Testing
- `go test ./internal/repository -run TestVideoJobRepositoryCreateListAndTransition -count=1`: passed.
- `go test ./internal/repository -run 'VideoJob|Migration' -count=1`: passed.
- `go test ./ent/schema -count=1`: passed.
- `go test ./migrations -count=1`: passed.
- `go test ./internal/service -run VideoJob -count=1`: passed.
### Notes
- `backend/ent/schema/video_job.go`: defined the durable video job fields and indexes.
- `backend/ent/schema/video_job_schema_test.go`: verified required fields and unique/query indexes.
- `backend/migrations/182_video_jobs.sql`: added the additive PostgreSQL table and indexes.
- `backend/migrations/video_jobs_migration_test.go`: verified the embedded migration contract.
- `backend/internal/service/video_job.go`: added video job states, domain data, transition data, and repository contract.
- `backend/internal/repository/video_job_repo.go`: implemented create, scoped reads, ordered lists, active scans, and conditional transitions.
- `backend/internal/repository/video_job_repo_test.go`: covered persistence, API key isolation, ordering, and transition conflicts.
- `backend/internal/repository/wire.go`: registered the video job repository provider.
- `backend/ent/videojob.go`: generated the VideoJob entity model.
- `backend/ent/videojob/videojob.go`: generated VideoJob field and order helpers.
- `backend/ent/videojob/where.go`: generated VideoJob predicates.
- `backend/ent/videojob_create.go`: generated VideoJob create builders.
- `backend/ent/videojob_delete.go`: generated VideoJob delete builders.
- `backend/ent/videojob_query.go`: generated VideoJob query builders.
- `backend/ent/videojob_update.go`: generated VideoJob update builders.
- `backend/ent/client.go`: registered the generated VideoJob client.
- `backend/ent/ent.go`: registered VideoJob mutation metadata.
- `backend/ent/hook/hook.go`: registered generated VideoJob hooks.
- `backend/ent/intercept/intercept.go`: registered generated VideoJob interceptors.
- `backend/ent/migrate/schema.go`: added the generated VideoJob table specification.
- `backend/ent/mutation.go`: added the generated VideoJob mutation implementation.
- `backend/ent/predicate/predicate.go`: added the generated VideoJob predicate type.
- `backend/ent/runtime/runtime.go`: initialized generated VideoJob runtime metadata.
- `backend/ent/tx.go`: exposed the generated VideoJob transactional client.
- `backend/go.sum`: recorded the Ent generator command dependency checksum.
- `progress.md`: recorded this tested task and rollback procedure.
- Rollback point: `d4daa84ed6bfdd338eb6c6c9f311b3b0e68487ea`; after this task is committed, run `git revert $(git log --format=%H --grep='^feat: persist leo video jobs$' -n 1)` from the repository root.

## 2026-07-17 - Task: Call LeoStudio asynchronous video jobs
### What was done
- Added exact `Prefer: respond-async` detection without treating parameterized or assigned values as opt-in.
- Added authenticated LeoStudio create, status, and pending-cancel clients with model mapping and typed responses.
- Kept transport and upstream errors detached from Gin and removed configured Leo credentials and sensitive query values from surfaced messages.
### Testing
- `go test ./internal/service -run 'PrefersLeoRespondAsync|LeoAsync' -count=1`: passed.
- `go test ./internal/service -run 'PrefersLeoRespondAsync|LeoAsync|ForwardLeoVideo|LeoVideo' -count=1`: passed, including existing synchronous Leo forwarding coverage.
- `git diff --check`: passed.
### Notes
- `backend/internal/service/leo_video_async.go`: implemented the narrow LeoStudio asynchronous protocol client and typed errors.
- `backend/internal/service/leo_video_async_test.go`: covered preference parsing, model mapping, Bearer authentication, create/get/cancel decoding, and secret redaction.
- `progress.md`: recorded this tested task and rollback procedure.
- Rollback point: `7605fd45`; run `git revert $(git log --format=%H --grep='^feat: call leo async video jobs$' -n 1)` from the repository root.

## 2026-07-17 - Task: Freeze and settle asynchronous video billing
### What was done
- Added idempotent video-specific balance hold and release commands using the existing usage dedup tables and frozen-balance accounting.
- Added submission-time video price/rate snapshots and terminal settlement that records actual output cost through normal usage billing before releasing the full hold.
- Added a cost override to normal OpenAI usage recording so asynchronous settlement does not recalculate against changed group pricing.
### Testing
- `go test ./internal/service ./internal/repository -run 'VideoJobBilling|VideoBalance|CostOverride' -count=1`: passed.
- `go test ./internal/service -count=1`: passed.
- `go test ./internal/repository -count=1`: passed.
- `git diff --check`: passed.
### Notes
- `backend/internal/service/usage_billing.go`: added video hold command/result types and repository contract methods.
- `backend/internal/repository/usage_billing_repo.go`: implemented idempotent video hold/release transactions and frozen-balance updates.
- `backend/internal/repository/video_balance_hold_repo_test.go`: verified reserve and release SQL behavior.
- `backend/internal/service/video_job_billing.go`: implemented price snapshots, reserve, actual-cost settlement, release, and loader contracts.
- `backend/internal/service/video_job_billing_test.go`: covered snapshot hold, idempotent completion, and failure release behavior.
- `backend/internal/service/openai_gateway_usage.go`: added the optional cost override path while preserving normal pricing.
- `backend/internal/service/openai_gateway_record_usage_test.go`: verified snapshot cost wins over current group pricing.
- `backend/internal/service/batch_image_settlement_test.go`: extended the billing fake for the additive video hold interface.
- `progress.md`: recorded this tested task and rollback procedure.
- Rollback point: `f0d930bb`; run `git revert $(git log --format=%H --grep='^feat: settle async leo video billing$' -n 1)` from the repository root.

## 2026-07-17 - Task: Create and access durable Leo video jobs
### What was done
- Added the task orchestration service for validation, Leo account selection, model mapping, hold-before-upstream submission, durable mapping, bounded pre-202 failover, and ambiguous-failure handling.
- Added API-key-scoped list/detail access and pending-only cancellation with upstream cancellation and idempotent hold release.
- Extended the request snapshot parser with aspect ratio and audio fields and allowed account affinity to change only while a job remains unaccepted.
### Testing
- `go test ./internal/service -run TestVideoJobService -count=1`: passed.
- `go test ./internal/service -run 'VideoJobService|VideoJobBilling|LeoVideo' -count=1`: passed.
- `go test ./internal/repository -run VideoJob -count=1`: passed.
- `go test ./ent/schema -run VideoJob -count=1`: passed.
- `git diff --check`: passed.
### Notes
- `backend/internal/service/video_job_service.go`: implemented durable submission, failover, access, and cancellation orchestration plus the gateway account-selector adapter.
- `backend/internal/service/video_job_service_test.go`: covered validation, ordering, mapping, failover, ambiguity, isolation, and cancellation.
- `backend/internal/service/video_job.go`: added public cancel conflict, account transition support, and opaque video job ID generation.
- `backend/ent/schema/video_job.go`: allowed pre-acceptance account affinity updates.
- `backend/ent/videojob_create.go`: regenerated account create builder metadata after the schema change.
- `backend/ent/videojob_update.go`: regenerated account update setter used by pre-acceptance failover.
- `backend/internal/repository/video_job_repo.go`: applied conditional account updates during failover.
- `backend/internal/service/leo_video.go`: parsed aspect ratio and audio request fields.
- `progress.md`: recorded this tested task and rollback procedure.
- Rollback point: `ab263a10`; run `git revert $(git log --format=%H --grep='^feat: manage leo async video jobs$' -n 1)` from the repository root.

## 2026-07-17 - Task: Reconcile Leo video jobs in the background
### What was done
- Added a restart-safe runtime that scans up to 50 active jobs, polls the fixed Leo account, preserves transient failures for retry, and advances completed jobs through settling before completion.
- Added terminal failure/cancel release handling and loop shutdown that waits for in-flight upstream polling.
- Wired the shared billing service, video job service, runtime startup, and cleanup ordering into the server and regenerated Wire output.
### Testing
- `go test ./internal/service -run TestVideoJobRuntime -count=1`: passed.
- `go test ./internal/service -run 'TestVideoJobRuntime|VideoJobService|VideoJobBilling' -count=1`: passed.
- `go test ./cmd/server -count=1`: passed.
- `go generate ./cmd/server`: passed.
- `git diff --check`: passed.
### Notes
- `backend/internal/service/video_job_runtime.go`: implemented active-job polling, settlement transitions, retry behavior, and Start/Stop lifecycle.
- `backend/internal/service/video_job_runtime_test.go`: covered completion, settling, failure release, transient retry, restart scan, and in-flight Stop waiting.
- `backend/internal/service/video_job_service.go`: added the production constructor used by Wire.
- `backend/internal/service/video_job_service_test.go`: added the active-job fake repository method for runtime coverage.
- `backend/internal/service/wire.go`: provided shared video billing/runtime services and interface bindings.
- `backend/cmd/server/wire.go`: stopped the video runtime before infrastructure shutdown.
- `backend/cmd/server/wire_gen.go`: regenerated runtime construction and cleanup wiring.
- `backend/cmd/server/wire_gen_test.go`: updated cleanup fixture arguments.
- `progress.md`: recorded this tested task and rollback procedure.
- Rollback point: `80c07244`; run `git revert $(git log --format=%H --grep='^feat: reconcile leo video jobs$' -n 1)` from the repository root.

## 2026-07-17 - Task: Store temporary Leo video input images locally
### What was done
- Added local disk storage under `<data_dir>/video-inputs` with content-based PNG/JPEG/WebP validation, a 10 MiB limit, opaque URL-safe tokens, and atomic file permissions.
- Added terminal one-hour and orphan 24-hour cleanup plus loopback-only internal reads that reject non-loopback clients before token lookup.
- Added an API-key-protected multipart upload handler returning an upload ID and loopback URL.
### Testing
- `go test ./internal/service ./internal/handler -run 'VideoInput' -count=1`: passed.
- `git diff --check`: passed.
### Notes
- `backend/internal/service/video_input_store.go`: implemented local image save/open, token URL handling, terminal marking, and cleanup.
- `backend/internal/service/video_input_store_test.go`: covered MIME detection, size limits, opaque names, terminal cleanup, and orphan cleanup.
- `backend/internal/handler/video_input.go`: implemented authenticated upload and loopback-only streaming handlers.
- `backend/internal/handler/video_input_test.go`: covered multipart upload, loopback read, non-loopback 404, and missing API key.
- `progress.md`: recorded this tested task and rollback procedure.
- Rollback point: `d6ea34e4`; run `git revert $(git log --format=%H --grep='^feat: host temporary leo video inputs$' -n 1)` from the repository root.

## 2026-07-17 - Task: Expose Leo asynchronous video gateway routes
### What was done
- Added `Prefer: respond-async` dispatch before the existing synchronous Leo flow, preserving the synchronous response and accounting path.
- Added public generation, upload, list, detail, cancel, and loopback input routes with Leo platform gating, API-key auth, no-store responses, and public-only job DTOs.
- Wired video job and local input services into the OpenAI gateway handler without changing existing direct constructor call sites.
### Testing
- `go test ./internal/handler -run 'LeoVideoAsync|VideoInput' -count=1`: passed.
- `go test ./internal/server/routes -count=1`: passed.
- `go test ./internal/handler -run 'LeoVideo|VideoInput' -count=1`: passed.
- `go generate ./cmd/server`: passed.
- `go test ./cmd/server -count=1`: passed.
- `git diff --check`: passed.
### Notes
- `backend/internal/handler/leo_video_async.go`: implemented async generation, upload, list, detail, cancel, DTO, and internal input endpoints.
- `backend/internal/handler/leo_video_async_test.go`: covered 202 public mapping and API-key isolation.
- `backend/internal/handler/leo_video.go`: dispatched exact async preference while preserving sync behavior.
- `backend/internal/handler/openai_gateway_handler.go`: added video service fields and setter injection.
- `backend/internal/handler/wire.go`: added video input/gateway providers and service injection.
- `backend/internal/server/routes/gateway.go`: registered Leo-only `/videos/jobs`, `/videos/uploads`, cancel, and loopback routes.
- `backend/cmd/server/wire_gen.go`: regenerated handler construction for video dependencies.
- `progress.md`: recorded this tested task and rollback procedure.
- Rollback point: `6a208440`; run `git revert $(git log --format=%H --grep='^feat: expose leo asynchronous video gateway routes$' -n 1)` from the repository root.

## 2026-07-17 - Task: Add Leo video client API
### What was done
- Added typed browser calls for local image upload, asynchronous generation, job listing, job detail, and pending-job cancellation.
- Kept LeoStudio credentials out of the browser contract; every request uses the selected Sub2API API Key.
### Testing
- `node node_modules/vitest/vitest.mjs run src/api/__tests__/videoGeneration.spec.ts --reporter=verbose`: passed (3 tests).
### Notes
- `frontend/src/api/videoGeneration.ts`: added typed Leo video request/response wrappers and shared error parsing.
- `frontend/src/api/__tests__/videoGeneration.spec.ts`: covered multipart upload headers, async preference, API-key auth, list query, and DELETE cancellation.
- `progress.md`: recorded this tested task and rollback procedure.
- Rollback point: `2f586842`; run `git revert $(git log --format=%H --grep='^feat: add leo video client api$' -n 1)` from the repository root.

## 2026-07-17 - Task: Add Leo video generation workbench
### What was done
- Added a responsive two-column user workbench for text, remote-image, and local-image Leo video generation.
- Added active Leo API Key filtering, asynchronous job polling, pending cancellation, completed video preview, download/open actions, and empty states.
- Added `/video-generation`, the user sidebar entry, and English/Chinese navigation and workbench messages.
### Testing
- `node node_modules/vitest/vitest.mjs run src/views/user/__tests__/VideoGenerationView.spec.ts src/api/__tests__/videoGeneration.spec.ts --reporter=verbose`: passed (8 tests).
- `node node_modules/typescript/bin/tsc --noEmit -p tsconfig.json`: passed.
- `git diff --check`: passed.
- Targeted ESLint could not start because the existing workspace dependency `vue-eslint-parser` is missing; no dependency was added in this task.
### Notes
- `frontend/src/views/user/VideoGenerationView.vue`: implemented the responsive video workbench and async lifecycle controls.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: covered Leo Key filtering, text/URL submission, preview, cancellation, and no-Key state.
- `frontend/src/router/index.ts`: registered the authenticated `/video-generation` route.
- `frontend/src/components/layout/AppSidebar.vue`: added a video-generation navigation icon and item.
- `frontend/src/i18n/locales/en/common.ts`: added English navigation label.
- `frontend/src/i18n/locales/zh/common.ts`: added Chinese navigation label.
- `frontend/src/i18n/locales/en/dashboard.ts`: added English workbench messages.
- `frontend/src/i18n/locales/zh/dashboard.ts`: added Chinese workbench messages.
- `progress.md`: recorded this tested task and rollback procedure.
- Rollback point: `b3b72ec5`; run `git revert $(git log --format=%H --grep='^feat: add leo video workbench$' -n 1)` from the repository root.

## 2026-07-17 - Task: Wire local video input lifecycle cleanup
### What was done
- Shared one local input store across upload handling, job submission, cancellation, and the background runtime.
- Marked local inputs terminal when jobs complete, fail, or cancel, and added startup/daily orphan and terminal cleanup execution.
### Testing
- `go test ./internal/service -run 'TestVideoJob(Runtime|Service)' -count=1`: passed.
- `go generate ./cmd/server`: passed.
- `go test ./cmd/server -count=1`: passed.
- `go test ./internal/service ./internal/handler ./internal/server/routes ./cmd/server -count=1`: passed.
- `git diff --check`: passed.
### Notes
- `backend/internal/service/video_input_store.go`: added terminal marking helper.
- `backend/internal/service/video_job_runtime.go`: added shared-store cleanup scheduling and terminal marking.
- `backend/internal/service/video_job_service.go`: marked canceled local inputs terminal.
- `backend/internal/service/video_job_runtime_test.go`: verified delayed terminal-input cleanup.
- `backend/internal/service/video_job_service_test.go`: verified cancellation terminal marking.
- `backend/internal/service/wire.go`: provided one shared input store and injected it into job service/runtime.
- `backend/internal/handler/wire.go`: reused the shared store for upload handling.
- `backend/cmd/server/wire_gen.go`: regenerated shared-store wiring.
- `progress.md`: recorded this tested task and rollback procedure.
- Rollback point: `aa502fd3`; run `git revert $(git log --format=%H --grep='^fix: wire local video input lifecycle cleanup$' -n 1)` from the repository root.

## 2026-07-17 - Task: Document and verify Leo async video workbench
### What was done
- Updated the Leo channel guide with async submission, workbench usage, local-input lifecycle, same-host deployment, billing, and operational boundaries.
### Testing
- `backend: go test ./... -count=1`: passed.
- `backend: go vet ./...`: passed.
- `frontend: node node_modules/vitest/vitest.mjs run --exclude src/i18n/__tests__/localesMessageCompile.spec.ts --reporter=dot`: passed; the excluded pre-existing locale compiler test requires missing `@intlify/message-compiler`.
- `frontend: node node_modules/vue-tsc/bin/vue-tsc.js --noEmit`: passed.
- `frontend: node node_modules/vite/bin/vite.js build`: passed; only existing chunk-size and dynamic-import warnings were emitted.
- `frontend: node node_modules/eslint/bin/eslint.js ...`: unavailable because the existing workspace dependency `vue-eslint-parser` is missing.
- `git diff --check`: passed.
### Notes
- `docs/LEO_VIDEO_CHANNEL.md`: reconciled sync compatibility with the async workbench and documented setup, API, local storage, cleanup, billing, and limits.
- `progress.md`: recorded the final tested task and rollback procedure.
- Rollback point: `74850777`; run `git revert $(git log --format=%H --grep='^docs: explain leo async video workbench$' -n 1)` from the repository root.

## 2026-07-19 - Task: Add admin video generation menu switch
### What was done
- Added the `video_generation_enabled` system setting with a default-enabled, explicit-false opt-out policy across admin, public settings, and frontend injection.
- Added an administrator toggle under System Settings > Feature Flags and wired the user sidebar video generation entry to the shared feature flag registry.
- Kept the switch limited to sidebar visibility; the `/video-generation` route and Leo APIs remain available.
### Testing
- `backend: go test -tags unit ./internal/service ./internal/handler/dto ./internal/server -run 'VideoGenerationMenuSwitch|PublicSettingsInjectionPayload|TestAPIContracts' -count=1`: passed.
- `backend: go test ./... -count=1`: passed.
- `backend: go vet ./...`: passed.
- `frontend: .\\node_modules\\.bin\\vitest.cmd run --exclude src/i18n/__tests__/localesMessageCompile.spec.ts --reporter=dot`: passed; the excluded existing test requires the missing `@intlify/message-compiler` dependency.
- `frontend: .\\node_modules\\.bin\\vitest.cmd run src/components/layout/__tests__/AppSidebar.spec.ts src/views/admin/__tests__/SettingsView.spec.ts`: passed (34 tests).
- `frontend: .\\node_modules\\.bin\\vue-tsc.cmd --noEmit`: passed.
- `frontend: .\\node_modules\\.bin\\vite.cmd build`: passed with existing dynamic-import and chunk-size warnings.
- Targeted ESLint could not start because the existing workspace dependency `vue-eslint-parser` is missing; no dependency was added in this task.
- `git diff --check`: passed; Git reported LF-to-CRLF working-copy notices for `docs/LEO_VIDEO_CHANNEL.md` and `progress.md`.
### Notes
- `backend/internal/handler/admin/setting_handler.go`: returned the video menu setting from the administrator settings endpoint.
- `backend/internal/handler/admin/setting_handler_audit.go`: included video menu changes in administrator setting audit diffs.
- `backend/internal/handler/admin/setting_handler_update.go`: accepted, preserved, and returned the optional video menu setting during updates.
- `backend/internal/handler/dto/settings.go`: added the video menu field to admin and public response DTOs.
- `backend/internal/handler/setting_handler.go`: returned the setting from the public settings endpoint.
- `backend/internal/server/api_contract_test.go`: updated administrator settings response contracts with the default-enabled field.
- `backend/internal/service/domain_constants.go`: defined the `video_generation_enabled` setting key and its boundary.
- `backend/internal/service/setting_parse.go`: initialized and parsed the default-enabled setting.
- `backend/internal/service/setting_public.go`: loaded and injected the public video menu setting.
- `backend/internal/service/setting_service_public_test.go`: covered the absent-setting default and explicit disable behavior.
- `backend/internal/service/setting_service_update_test.go`: covered persistence of an explicit false value.
- `backend/internal/service/setting_update.go`: persisted the setting in system setting updates.
- `backend/internal/service/settings_view.go`: exposed the field in service-level admin and public settings views.
- `docs/LEO_VIDEO_CHANNEL.md`: documented the administrator switch, default, and non-authorization boundary.
- `frontend/src/api/admin/settings.ts`: added administrator read and update typings.
- `frontend/src/components/layout/AppSidebar.vue`: attached the video generation entry to the shared opt-out flag.
- `frontend/src/components/layout/__tests__/AppSidebar.spec.ts`: verified registry wiring plus unloaded, missing, disabled, and enabled runtime states.
- `frontend/src/i18n/locales/en/admin/settings.ts`: added English administrator labels and boundary text.
- `frontend/src/i18n/locales/zh/admin/settings.ts`: added Chinese administrator labels and boundary text.
- `frontend/src/stores/__tests__/app.spec.ts`: updated the complete public settings fixture with the default-enabled field.
- `frontend/src/stores/app.ts`: preserved the default-enabled behavior in the API fallback payload.
- `frontend/src/types/index.ts`: added the field to public settings typings.
- `frontend/src/utils/featureFlags.ts`: registered the video menu as an opt-out feature flag.
- `frontend/src/views/admin/SettingsView.vue`: added the toggle, default form value, and save payload field.
- `frontend/src/views/admin/__tests__/SettingsView.spec.ts`: verified loading and submitting the video menu toggle.
- `progress.md`: recorded this tested task and rollback procedure.
- Rollback point: `8067ef3c`; run `git revert $(git log --format=%H --grep='^feat: add video generation menu setting$' -n 1)` from the repository root after the task commit is created.

## 2026-07-19 - Task: Integrate the Leo video workbench and menu switch into v0.1.161
### What was done
- Ported the complete asynchronous Leo video job chain, user workbench, and administrator video-generation menu switch onto the `v0.1.161-fy.1` fork baseline.
- Resolved dependency-injection conflicts so v0.1.161 prompt auditing, Grok eligibility probing, authentication cache invalidation, Token Incentive, update/restart behavior, and the new video runtime remain wired together.
- Added the image-only `/videos/uploads` route to the v0.1.161 prompt-audit route manifest with an explicit no-prompt reason; video generation prompts remain audited at job submission.
- Confirmed the duplicate `182_*` migration prefix is safe because migrations are tracked by full filename, so `182_prompt_audit_full_prompt.sql` and `182_video_jobs.sql` execute independently.
### Testing
- Targeted backend video, settings, API-contract, Ent schema, migration, repository, service, handler, and route tests: passed across 7 packages.
- Targeted frontend video API, workbench, sidebar, administrator settings, and app-store tests: passed across 5 test files.
- `go test ./internal/server/routes -run 'TestEveryGatewayPOSTRouteIsClassifiedForPromptAuditCoverage|Leo|Video' -count=1`: passed after classifying the upload route.
- `go test -p 2 -tags=unit -timeout 10m ./...`: passed for all backend packages.
- `..\\.codex-run\\bin\\golangci-lint.exe run --concurrency 2 ./...`: passed with `0 issues`.
- `pnpm run test:run`: passed for the complete frontend Vitest suite.
- `pnpm run lint:check`: passed.
- `pnpm run typecheck`: passed.
- `pnpm run build`: passed.
- PostgreSQL migration integration entrypoint exited successfully but explicitly skipped because Docker is unavailable on this Windows host; GitHub CI remains the real database integration gate.
- `go build -o ..\\.codex-run\\sub2api-v0.1.161-fy.2.exe ./cmd/server` and `--version`: passed and reported `Sub2API 0.1.161`.
- `git diff --check`: passed before this final progress append; final staged checks follow.
### Notes
- The preceding 12 task entries in `progress.md` list every backend, frontend, migration, generated Ent, documentation, and test file introduced by the video implementation commits.
- `backend/internal/server/routes/prompt_audit_route_coverage_test.go`: classified the image-upload-only route while preserving prompt-audit enforcement for video generation submissions.
- `progress.md`: recorded the v0.1.161 integration, conflict resolution, verification evidence, release target, and rollback procedure.
- Rollback point: deploy `v0.1.161-fy.1`; for source rollback, run `git revert --no-commit v0.1.161-fy.1..HEAD`, review the staged reversal, and commit it without applying or deleting the homepage stash or `.superpowers/`.

## 2026-07-19 - Task: Settle Leo video billing only after local output is saved
### What was done
- Added durable local MP4 download, size and `ftyp` validation, atomic publication, and restart-safe reuse under `pricing.data_dir/video-outputs/`.
- Changed async settlement so a completed upstream job is charged only after its first video is saved locally; empty results fail immediately, save failures retry up to three times, and every terminal failure releases the frozen balance without creating usage.
- Added an API Key-scoped content endpoint with Range support and changed the user workbench to fetch the saved video with Bearer authentication and play a Blob URL instead of the upstream CDN URL.
- Rejected synchronous HTTP-success responses that contain no usable video URL, preventing empty success responses from being billed.
### Testing
- `backend: go generate ./cmd/server`: passed and regenerated Wire output.
- `backend: go test ./internal/service ./internal/handler -count=1`: passed.
- `backend: go test ./... -count=1`: passed.
- `backend: go vet ./...`: passed.
- `frontend: .\node_modules\.bin\vitest.cmd run src/views/user/__tests__/VideoGenerationView.spec.ts src/api/__tests__/videoGeneration.spec.ts`: passed (9 tests).
- `frontend: .\node_modules\.bin\vitest.cmd run --exclude src/i18n/__tests__/localesMessageCompile.spec.ts --reporter=dot`: passed; the excluded existing test requires the missing `@intlify/message-compiler` dependency.
- `frontend: .\node_modules\.bin\vue-tsc.cmd --noEmit`: passed.
- `frontend: .\node_modules\.bin\eslint.cmd src/api/videoGeneration.ts src/api/__tests__/videoGeneration.spec.ts src/views/user/VideoGenerationView.vue src/views/user/__tests__/VideoGenerationView.spec.ts`: passed.
- `frontend: .\node_modules\.bin\vite.cmd build`: passed with existing dynamic-import and chunk-size warnings.
- `git diff --check`: passed.
### Notes
- `backend/internal/service/video_output_store.go`: added local MP4 download, validation, result rewriting, idempotent reuse, and file opening.
- `backend/internal/service/video_output_store_test.go`: covered valid download, invalid MP4, missing URL, and idempotent reuse.
- `backend/internal/service/video_job_runtime.go`: required local save before settlement and added bounded save retries with failure release.
- `backend/internal/service/video_job_runtime_test.go`: verified save-before-charge, empty-result release, retry exhaustion, and persisted local output.
- `backend/internal/service/video_job_billing.go`: stopped treating empty result data as one billable video.
- `backend/internal/service/video_job_billing_test.go`: verified empty results cannot create a settlement.
- `backend/internal/service/video_job_service_test.go`: persisted transition result and actual cost in the shared repository test double.
- `backend/internal/service/leo_video.go`: rejected synchronous successful responses without a usable video URL before committing the response.
- `backend/internal/service/leo_video_test.go`: covered synchronous empty-success rejection.
- `backend/internal/service/wire.go`: provided the shared output store and injected it into the video runtime.
- `backend/internal/handler/leo_video_async.go`: added API Key ownership checks and authenticated MP4 serving.
- `backend/internal/handler/leo_video_async_test.go`: verified content ownership isolation and MP4 delivery.
- `backend/internal/handler/openai_gateway_handler.go`: attached the shared output store to the gateway handler.
- `backend/internal/handler/wire.go`: injected the shared output store into the gateway handler.
- `backend/internal/server/routes/gateway.go`: registered the authenticated video content route in both gateway route layouts.
- `backend/cmd/server/wire_gen.go`: regenerated application wiring for the shared output store.
- `frontend/src/api/videoGeneration.ts`: added Bearer-authenticated saved-video download.
- `frontend/src/api/__tests__/videoGeneration.spec.ts`: verified the content URL and Authorization header.
- `frontend/src/views/user/VideoGenerationView.vue`: switched preview and download to managed Blob URLs and revoked them on selection changes and unmount.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: verified local Blob playback and URL cleanup instead of CDN playback.
- `docs/LEO_VIDEO_CHANNEL.md`: documented local output storage, authenticated access, settlement order, retries, and deployment limits.
- `progress.md`: recorded this implementation, verification evidence, file list, and rollback point.
- Rollback point: `74c877d5`; after creating the task commit, run `git revert <task-commit-hash>` from the repository root.

## 2026-07-20 - Task: Integrate local Leo video output persistence into v0.1.161
### What was done
- Recovered the uncommitted video-output worktree changes that were not included in `v0.1.161-fy.2` and integrated them into the active v0.1.161 release branch.
- Preserved v0.1.161 prompt auditing, Grok eligibility probing, authentication cache invalidation, Token Incentive, and fork update wiring while injecting the shared video output store.
- Kept unrelated Baota daily-report commits, homepage preview worktrees, stashes, and `.superpowers/` outside the release.
- Fixed three static-analysis findings in the new runtime tests by explicitly checking interface type assertions before inspecting test doubles.
### Testing
- Source-worktree backend video tests passed across service, handler, and routes; targeted frontend video API and workbench tests passed.
- Final-branch backend video, output storage, authentication, route, and prompt-audit coverage tests passed.
- `go test -p 2 -tags=unit -timeout 10m ./...`: passed for all backend packages.
- `pnpm run test:run`: passed for the complete frontend Vitest suite.
- `pnpm run lint:check`, `pnpm run typecheck`, and `pnpm run build`: passed.
- `..\\.codex-run\\bin\\golangci-lint.exe run --concurrency 2 ./...`: passed with `0 issues` after the test-only type-assertion fix.
- `go build -o ..\\.codex-run\\sub2api-v0.1.161-fy.3.exe ./cmd/server` and `--version`: passed and reported `Sub2API 0.1.161`.
- `git diff --check`: passed before this final progress append; final staged checks follow.
### Notes
- `backend/internal/handler/openai_gateway_handler.go`: retained v0.1.161 security and Grok dependencies while attaching `VideoOutputStore`.
- `backend/internal/handler/wire.go`: merged output-store injection into the v0.1.161 gateway provider.
- `backend/cmd/server/wire_gen.go`: regenerated final dependency injection with prompt audit, Token Incentive, and video output persistence together.
- `backend/internal/service/video_job_runtime_test.go`: explicitly checked test-double type assertions to satisfy static analysis without changing behavior.
- `progress.md`: recorded recovery of the missed local changes, final-branch verification, exclusions, release target, and rollback procedure.
- Rollback point: deploy `v0.1.161-fy.2`; for source rollback after release, run `git revert --no-commit v0.1.161-fy.2..v0.1.161-fy.3`, review the staged reversal, and commit it without applying or deleting unrelated worktrees or stashes.

## 2026-07-20 - Task: Sync LeoStudio Seedance guidance and multi-image support
### What was done
- Advanced the existing local LeoStudio checkout from `083180d3` to upstream `7af385af` and reviewed its new video guidance contract.
- Added Sub2 request typings and passthrough coverage for start/end frames, up to four image references, and full Seedance `guidances` including existing video/audio asset IDs.
- Changed the user video workbench from one image to up to four local files or remote URLs and submitted them through `image_urls` while retaining legacy `image_url` API compatibility.
- Tracked every local image URL used by top-level and nested guidance fields so all associated temporary inputs enter cleanup after completion, failure, or cancellation without a database migration.
- Preserved LeoStudio `400` and `422` asynchronous guidance validation responses as Sub2 client errors instead of reporting them as upstream gateway failures.
### Testing
- `leostudio-admin: go test ./internal/leonardo ./internal/service ./internal/server -count=1`: passed.
- `backend: go test ./internal/service ./internal/handler -run 'LeoVideo|VideoInput|VideoJob' -count=1`: passed.
- `backend: go test ./... -count=1`: passed.
- `backend: go vet ./...`: passed.
- `frontend: .\node_modules\.bin\vitest.cmd run src/views/user/__tests__/VideoGenerationView.spec.ts src/api/__tests__/videoGeneration.spec.ts`: passed (10 tests).
- `frontend: .\node_modules\.bin\vitest.cmd run --exclude src/i18n/__tests__/localesMessageCompile.spec.ts --reporter=dot`: passed; the excluded existing test requires the missing `@intlify/message-compiler` dependency.
- `frontend: .\node_modules\.bin\vue-tsc.cmd --noEmit`: passed.
- `frontend: .\node_modules\.bin\eslint.cmd src/api/videoGeneration.ts src/api/__tests__/videoGeneration.spec.ts src/views/user/VideoGenerationView.vue src/views/user/__tests__/VideoGenerationView.spec.ts src/i18n/locales/en/dashboard.ts src/i18n/locales/zh/dashboard.ts`: passed.
- `frontend: .\node_modules\.bin\vite.cmd build`: passed with existing dynamic-import and chunk-size warnings.
- `git diff --check`: passed.
### Notes
- `backend/internal/handler/leo_video_async.go`: collected all local guidance tokens and preserved upstream asynchronous validation status codes.
- `backend/internal/handler/leo_video_async_test.go`: covered multi-input tracking and upstream guidance validation responses.
- `backend/internal/service/leo_video.go`: recognized top-level and nested image guidance URLs while keeping the request payload intact.
- `backend/internal/service/leo_video_test.go`: verified new guidance passthrough and URL collection order with deduplication.
- `backend/internal/service/video_input_store.go`: supported bounded multi-token lifecycle tracking in the existing job field.
- `backend/internal/service/video_input_store_test.go`: verified all local guidance images are marked terminal and cleaned.
- `frontend/src/api/videoGeneration.ts`: typed the latest start/end frame, image array, and full guidance request fields.
- `frontend/src/api/__tests__/videoGeneration.spec.ts`: verified new guidance fields are serialized unchanged.
- `frontend/src/views/user/VideoGenerationView.vue`: added up to four local or remote reference images with preview and cleanup.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: covered multi-URL and multi-file job submission.
- `frontend/src/i18n/locales/en/dashboard.ts`: updated English multi-image workbench labels and validation text.
- `frontend/src/i18n/locales/zh/dashboard.ts`: updated Chinese multi-image workbench labels and validation text.
- `docs/LEO_VIDEO_CHANNEL.md`: synchronized the LeoStudio version, Base URL, guidance fields, multi-image limits, and local lifecycle behavior.
- `progress.md`: recorded the protocol synchronization, verification evidence, file list, and rollback point.
- External reference checkout: `C:\Users\feiyu\.codex\tmp\leostudio-admin-review-083180d3` is clean at upstream commit `7af385af` and has no Sub2-specific edits.
- Rollback point: `0a4a064f`; after creating the task commit, run `git revert <task-commit-hash>` from the repository root.

## 2026-07-20 - Task: Integrate Leo multi-image guidance into the v0.1.161 release branch
### What was done
- Integrated the reviewed LeoStudio `7af385af` guidance and multi-image compatibility changes into the active fork release branch.
- Confirmed legacy `image_url`, synchronous video forwarding, asynchronous jobs, local input cleanup, local output persistence, menu gating, and v0.1.161 security wiring remain present.
- Kept Baota daily-report commits, homepage preview worktrees, stashes, and `.superpowers/` outside this release.
- Confirmed this change reuses the existing `video_jobs.local_input_name` field and requires no database migration.
### Testing
- Final-branch targeted backend Leo video, multi-input, job, route, and prompt-audit coverage tests passed.
- Final-branch targeted frontend multi-image workbench, video API, and sidebar tests passed.
- `go test -p 2 -tags=unit -timeout 10m ./...`: passed for all backend packages.
- `pnpm run test:run`: passed for the complete frontend Vitest suite.
- `pnpm run lint:check`, `pnpm run typecheck`, and `pnpm run build`: passed.
- `..\\.codex-run\\bin\\golangci-lint.exe run --concurrency 2 ./...`: passed with `0 issues`.
- `go build -o ..\\.codex-run\\sub2api-v0.1.161-fy.4.exe ./cmd/server` and `--version`: passed and reported `Sub2API 0.1.161`.
- `git diff --check`: passed before this final progress append; final staged checks follow.
### Notes
- The preceding multi-image guidance task entry lists all 13 implementation and documentation files changed by the feature commit.
- `progress.md`: recorded final-branch compatibility checks, complete verification, exclusions, release target, and rollback procedure.
- Rollback point: deploy `v0.1.161-fy.3`; for source rollback after release, revert the `v0.1.161-fy.3..v0.1.161-fy.4` commit range without applying or deleting unrelated worktrees or stashes.

## 2026-07-20 - Task: Fix Leo video reference submission and provider error exposure
### What was done
- Restored the user workbench's verified single-image `image_url` start-frame submission and changed multi-image submissions to explicit ordered `guidances.image_reference` entries.
- Unified asynchronous video errors under the public `Video provider` label across immediate failures, persisted jobs, runtime failures, and legacy failed-job reads.
- Kept the existing external `image_urls` API compatibility and all unrelated channel-pricing worktree changes untouched.
### Testing
- `go test ./internal/service ./internal/handler -run 'LeoVideo|VideoJob|SanitizeVideoProvider' -count=1`: passed.
- `go test ./internal/service ./internal/handler -count=1`: passed.
- `vitest run src/views/user/__tests__/VideoGenerationView.spec.ts src/api/__tests__/videoGeneration.spec.ts`: passed, 11 tests.
- Targeted ESLint for the changed video view and tests passed; `vue-tsc --noEmit` passed.
- Direct `vite build` passed with existing dynamic-import and chunk-size warnings.
- `git diff --check`: passed before this progress append.
- No production deployment or live credit-consuming video request was performed.
### Notes
- `frontend/src/views/user/VideoGenerationView.vue`: chooses the verified single-image contract and ordered multi-image guidance contract.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: covers single remote image and ordered remote/local multi-image payloads.
- `backend/internal/service/leo_video_async.go`: centralizes provider-name sanitization and uses provider-neutral asynchronous error messages.
- `backend/internal/service/leo_video_async_test.go`: verifies all named upstream aliases are hidden.
- `backend/internal/service/video_job_service.go`: sanitizes failed creation messages before persistence.
- `backend/internal/service/video_job_runtime.go`: sanitizes upstream terminal job errors before persistence.
- `backend/internal/handler/leo_video_async.go`: sanitizes immediate responses and legacy persisted errors at the user API boundary.
- `backend/internal/handler/leo_video_async_test.go`: verifies failed-job responses do not expose upstream names.
- `docs/LEO_VIDEO_CHANNEL.md`: documents the single/multi-image workbench contract and provider-neutral errors.
- `progress.md`: records this implementation, verification evidence, exclusions, and rollback command.
- Rollback point: `a5fb54df6f5860b91f2f5947bac504eb9d604471`; run `git restore --source=a5fb54df6f5860b91f2f5947bac504eb9d604471 -- backend/internal/handler/leo_video_async.go backend/internal/handler/leo_video_async_test.go backend/internal/service/leo_video_async.go backend/internal/service/leo_video_async_test.go backend/internal/service/video_job_runtime.go backend/internal/service/video_job_service.go docs/LEO_VIDEO_CHANNEL.md frontend/src/views/user/VideoGenerationView.vue frontend/src/views/user/__tests__/VideoGenerationView.spec.ts progress.md` from the Sub2API repository root.

## 2026-07-20 - Task: Add a dedicated video API integration documentation page
### What was done
- Added the authenticated `/video-generation/api-docs` subpage and shared workbench/API documentation tabs without adding another sidebar item.
- Documented Bearer authentication, local image upload, asynchronous video creation, single- and multi-image request bodies, job operations, lifecycle states, common errors, and provider privacy boundaries in Chinese and English.
- Added copyable API examples and responsive layouts whose code blocks and request table scroll within their own containers on narrow screens.
### Testing
- `vitest run src/components/video/__tests__/ApiCodeBlock.spec.ts src/components/video/__tests__/VideoSectionTabs.spec.ts src/views/user/__tests__/VideoApiDocsView.spec.ts src/views/user/__tests__/VideoGenerationView.spec.ts --reporter=dot`: passed, 4 files and 10 tests.
- Targeted ESLint for the new video documentation components, views, tests, route, and locale files: passed.
- `vue-tsc --noEmit`: passed.
- Direct `vite build`: passed with the existing dynamic-import and chunk-size warnings.
- In-app browser checks at `1440x1000` and `390x844`: passed with no page-level horizontal overflow, significant overlap, or malformed curl continuation markers; code blocks and the field table retained internal scrolling on mobile.
- Browser console contained no application errors; the temporary isolated preview only reported expected warnings for routes intentionally omitted from its minimal router.
- `git diff --check`: passed before this progress append; the temporary preview files were removed and the preview server was stopped.
- No production deployment or credit-consuming video request was performed.
### Notes
- `frontend/src/components/video/ApiCodeBlock.vue`: provides labeled, scrollable API examples with a clipboard action.
- `frontend/src/components/video/EndpointTitle.vue`: renders compact HTTP method and path headings.
- `frontend/src/components/video/SectionHeading.vue`: keeps documentation section headings consistent.
- `frontend/src/components/video/VideoSectionTabs.vue`: links the generation workbench and API documentation subpage.
- `frontend/src/components/video/__tests__/ApiCodeBlock.spec.ts`: verifies copied content and feedback state.
- `frontend/src/components/video/__tests__/VideoSectionTabs.spec.ts`: verifies both links and active-tab semantics.
- `frontend/src/views/user/VideoApiDocsView.vue`: implements the complete independent API integration documentation page and executable examples.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: covers the documented endpoint surface, image contracts, privacy boundary, and curl formatting.
- `frontend/src/views/user/VideoGenerationView.vue`: displays the shared video-page tabs above the existing workbench.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: stubs the shared tabs while retaining the workbench submission regression coverage.
- `frontend/src/router/index.ts`: registers the authenticated API documentation route.
- `frontend/src/i18n/locales/en/dashboard.ts`: adds English tab and API documentation copy.
- `frontend/src/i18n/locales/zh/dashboard.ts`: adds Chinese tab and API documentation copy.
- `docs/LEO_VIDEO_CHANNEL.md`: records the new subpage scope and its privacy boundary.
- `progress.md`: records the documentation implementation, verification evidence, changed files, and rollback point.
- Rollback point for the complete Leo video fix and API documentation sequence: `a5fb54df6f5860b91f2f5947bac504eb9d604471`. After committing this documentation task separately, use `git revert <api-docs-task-commit>` for an API-docs-only rollback.

## 2026-07-20 - Task: Add custom API Key mode to the video workbench
### What was done
- Added saved and custom Sub2 API Key modes to the user video workbench, including an in-memory password input with visibility control.
- Routed video submission, local reference image upload, job listing, cancellation, preview, and download through the active Key while preventing Key changes during submission.
- Cleared Key-scoped jobs, polling, and video Blob state when the Key source or custom value changes so results from different Keys cannot be mixed.
- Kept custom Keys out of browser storage and user account persistence; no backend authentication or database contract changed.
### Testing
- `frontend: .\node_modules\.bin\vitest.cmd run src/views/user/__tests__/VideoGenerationView.spec.ts src/api/__tests__/videoGeneration.spec.ts`: passed (13 tests).
- `frontend: .\node_modules\.bin\vitest.cmd run --exclude src/i18n/__tests__/localesMessageCompile.spec.ts --reporter=dot`: passed; the excluded existing test requires the missing `@intlify/message-compiler` dependency.
- `frontend: .\node_modules\.bin\vue-tsc.cmd --noEmit`: passed.
- `frontend: .\node_modules\.bin\eslint.cmd src/views/user/VideoGenerationView.vue src/views/user/__tests__/VideoGenerationView.spec.ts src/i18n/locales/en/dashboard.ts src/i18n/locales/zh/dashboard.ts`: passed.
- `frontend: .\node_modules\.bin\vite.cmd build`: passed with existing dynamic-import and chunk-size warnings.
- `git diff --check`: passed.
### Notes
- `frontend/src/views/user/VideoGenerationView.vue`: added the custom Key controls, effective-Key routing, and Key-scoped state reset behavior.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: covered custom-Key submission, upload, listing, cancellation, download, non-persistence, and mode isolation.
- `frontend/src/i18n/locales/en/dashboard.ts`: added English custom-Key controls and clarified Sub2 API Key wording.
- `frontend/src/i18n/locales/zh/dashboard.ts`: added Chinese custom-Key controls and clarified Sub2 API Key wording.
- `docs/LEO_VIDEO_CHANNEL.md`: documented custom-Key memory-only handling and the operations governed by the active Key.
- `progress.md`: recorded this implementation, verification evidence, file list, and rollback point.
- Rollback point: `efc132e5`; after creating the task commit, run `git revert <task-commit-hash>` from the repository root.

## 2026-07-20 - Task: Hide upstream provider names from video errors
### What was done
- Sanitized synchronous video validation errors, asynchronous submission errors, and failed job responses before returning them to users.
- Replaced Leonardo and LeoStudio naming variants with the neutral `video provider` label; logs retain the original internal error while user-visible persisted task errors are sanitized.
- Confirmed and documented that successful local video settlement creates a usage record and that custom-Key usage belongs to the Key owner and group.
### Testing
- `backend: go test ./internal/service ./internal/handler -run 'LeoVideo|PublicVideoError' -count=1`: passed.
- `backend: go test ./... -count=1`: passed.
- `backend: go vet ./...`: passed.
- Existing video billing and runtime tests passed within the full suite, including success usage recording and zero-usage failure release paths.
- `git diff --check`: passed.
### Notes
- `backend/internal/service/leo_video.go`: added the public video error name sanitizer and applied it to synchronous upstream validation responses.
- `backend/internal/service/leo_video_test.go`: covered provider-name variants and synchronous validation response sanitization.
- `backend/internal/handler/leo_video_async.go`: sanitized asynchronous submission and persisted failed-job messages at the user API boundary.
- `backend/internal/handler/leo_video_async_test.go`: verified submission and failed-job responses do not expose upstream provider names.
- `backend/internal/handler/openai_gateway_handler.go`: sanitized Leo client messages selected by administrator error-passthrough rules without changing other platforms.
- `docs/LEO_VIDEO_CHANNEL.md`: documented error-name sanitization and custom-Key usage ownership.
- `progress.md`: recorded this implementation, verification evidence, file list, and rollback point.
- Full uncommitted worktree rollback point: `efc132e5`; run `git restore --source=efc132e5 --worktree -- docs/LEO_VIDEO_CHANNEL.md frontend/src/i18n/locales/en/dashboard.ts frontend/src/i18n/locales/zh/dashboard.ts frontend/src/views/user/VideoGenerationView.vue frontend/src/views/user/__tests__/VideoGenerationView.spec.ts backend/internal/service/leo_video.go backend/internal/service/leo_video_test.go backend/internal/handler/leo_video_async.go backend/internal/handler/leo_video_async_test.go backend/internal/handler/openai_gateway_handler.go progress.md`.

## 2026-07-20 - Task: Add Seedance channel video pricing
### What was done
- Added channel-level per-second video pricing for `seedance-2.0` and `seedance-2.0-fast`, with independent exact-model entries and required 480p, 720p, and 1080p tiers.
- Made matching channel video pricing override the existing group video prices for both synchronous and asynchronous requests while retaining group prices as the fallback and preserving the group video multiplier.
- Froze the selected channel, billing model, pricing source, model mapping, three prices, and multiplier in asynchronous billing snapshot v2, and attributed successful usage to the frozen channel mapping.
- Added Leo channel administration controls, defaults, validation, and localized text without a database migration; explicit zero channel prices remain valid.
### Testing
- `backend: go generate ./cmd/server`: passed and regenerated Wire output.
- `backend: go test ./internal/service ./internal/handler -run 'Channel|Pricing|Video|Leo' -count=1`: passed.
- `backend: go test -tags=unit ./internal/service -count=1`: passed.
- `backend: go test ./... -count=1`: passed.
- `backend: go vet ./...`: passed.
- `frontend: .\node_modules\.bin\vitest.cmd run src/components/admin/channel/__tests__/types.spec.ts src/components/admin/channel/__tests__/PricingEntryCard.spec.ts src/views/user/__tests__/VideoGenerationView.spec.ts src/api/__tests__/videoGeneration.spec.ts`: passed (28 tests).
- `frontend: .\node_modules\.bin\vue-tsc.cmd --noEmit`: passed.
- Targeted frontend ESLint for every changed source and test file passed.
- `frontend: .\node_modules\.bin\vite.cmd build`: passed with existing dynamic-import and chunk-size warnings.
- `git diff --check`: passed.
### Notes
- `backend/cmd/server/wire_gen.go`: regenerated dependency wiring for the model pricing resolver in asynchronous video billing.
- `backend/internal/handler/leo_video.go`: attached resolved channel mapping fields to synchronous Leo video usage.
- `backend/internal/service/channel.go`: enabled video billing mode and resolution-tier interval validation.
- `backend/internal/service/channel_service.go`: required exact one-model video entries with complete non-negative resolution prices and rejected video mode in account statistics rules.
- `backend/internal/service/channel_service_test.go`: covered valid and invalid channel video pricing plus account statistics restrictions.
- `backend/internal/service/channel_test.go`: covered video billing mode and interval acceptance.
- `backend/internal/service/model_pricing_resolver.go`: resolved video tiers from channel pricing and converted complete tiers into video price configuration.
- `backend/internal/service/model_pricing_resolver_test.go`: covered channel video resolution and explicit zero prices.
- `backend/internal/service/openai_gateway_usage.go`: gave channel video prices precedence in synchronous cost calculation.
- `backend/internal/service/openai_gateway_record_usage_test.go`: verified synchronous channel price precedence and video multiplier behavior.
- `backend/internal/service/video_job_billing.go`: added channel price resolution, snapshot v2 freezing, group fallback, and frozen channel usage attribution for asynchronous jobs.
- `backend/internal/service/video_job_billing_test.go`: covered billing model sources, channel precedence, group fallback, zero pricing, and settlement attribution.
- `backend/internal/service/wire.go`: injected the shared model pricing resolver into asynchronous video billing.
- `frontend/src/constants/channel.ts`: registered video billing mode, Seedance model defaults, and supported resolutions.
- `frontend/src/components/admin/channel/PricingEntryCard.vue`: added the fixed three-tier USD-per-second editor for Leo model pricing.
- `frontend/src/components/admin/channel/types.ts`: added video pricing form defaults, normalization, model splitting, and validation.
- `frontend/src/components/admin/channel/__tests__/types.spec.ts`: covered video form serialization, defaults, syncing, and validation.
- `frontend/src/components/admin/channel/__tests__/PricingEntryCard.spec.ts`: covered Leo-only video mode visibility and price editing.
- `frontend/src/views/admin/ChannelsView.vue`: exposed Leo channel pricing and created independent Seedance video entries.
- `frontend/src/i18n/locales/en/admin/channels.ts`: added English video pricing labels and validation messages.
- `frontend/src/i18n/locales/zh/admin/channels.ts`: added Chinese video pricing labels and validation messages.
- `docs/LEO_VIDEO_CHANNEL.md`: documented configuration, exact-model rules, pricing precedence, formula, snapshot freezing, and usage attribution.
- `progress.md`: recorded implementation, verification evidence, changed files, and rollback guidance.
- Rollback point: `efc132e5`; after creating the task commit, run `git revert <task-commit-hash>` from the repository root.

## 2026-07-20 - Task: Integrate video pricing, custom Key workbench, and API documentation for v0.1.161-fy.5
### What was done
- Merged the Leo channel pricing and custom API Key workbench changes into the active fork release branch.
- Preserved the verified single-image `image_url` contract, ordered multi-image `guidances.image_reference` contract, provider-name sanitization, API documentation page, and existing v0.1.161 security and billing behavior.
- Added the API documentation page files that were present locally but not yet tracked, regenerated Wire, and validated the Windows service binary used for local verification.
- Kept `.superpowers/`, Baota daily-report worktrees, homepage preview worktrees, and unrelated stashes outside the release.
### Testing
- `backend: go generate ./cmd/server`: passed.
- `backend: go test ./internal/service ./internal/handler ./internal/server/routes -run 'Channel|Pricing|Video|Leo|SanitizeVideoProvider|PromptAuditCoverage' -count=1`: passed.
- `backend: go test -p 2 -tags=unit -timeout 10m ./...`: passed for all backend packages.
- `frontend: .\\node_modules\\.bin\\vitest.cmd run src/components/admin/channel/__tests__/types.spec.ts src/components/admin/channel/__tests__/PricingEntryCard.spec.ts src/views/user/__tests__/VideoGenerationView.spec.ts src/views/user/__tests__/VideoApiDocsView.spec.ts src/api/__tests__/videoGeneration.spec.ts --reporter=dot`: passed, 5 files and 30 tests.
- `frontend: .\\node_modules\\.bin\\vitest.cmd run`: passed for the complete frontend suite.
- `frontend: .\\node_modules\\.bin\\eslint.cmd . --ext .vue,.js,.jsx,.cjs,.mjs,.ts,.tsx,.cts,.mts`: passed.
- `frontend: .\\node_modules\\.bin\\vue-tsc.cmd --noEmit`: passed.
- `frontend: .\\node_modules\\.bin\\vite.cmd build`: passed with existing dynamic-import and chunk-size warnings.
- `backend: ..\\.codex-run\\bin\\golangci-lint.exe run --concurrency 2 ./...`: passed with `0 issues`.
- `backend: go build -o ..\\.codex-run\\sub2api-v0.1.161-fy.5.exe ./cmd/server` and `--version`: passed and reported `Sub2API 0.1.161`.
- `git diff --check`: passed before this progress append; final staged check follows.
### Notes
- `backend/cmd/server/wire_gen.go`: regenerated dependency wiring for asynchronous video billing.
- `backend/internal/handler/leo_video.go`: attaches resolved channel pricing context to synchronous video usage.
- `backend/internal/handler/leo_video_async.go`: returns sanitized asynchronous video errors and API-key-scoped jobs.
- `backend/internal/handler/leo_video_async_test.go`: covers provider-name sanitization and uses case-insensitive matching for the public label.
- `backend/internal/handler/openai_gateway_handler.go`: applies Leo provider-name sanitization to configured passthrough errors.
- `backend/internal/service/channel.go`: validates video billing mode and resolution tiers.
- `backend/internal/service/channel_service.go`: validates exact-model channel video pricing and rejects invalid statistics pricing combinations.
- `backend/internal/service/channel_service_test.go`: covers channel video pricing validation.
- `backend/internal/service/channel_test.go`: covers video billing mode and interval validation.
- `backend/internal/service/leo_video.go`: sanitizes synchronous public video errors.
- `backend/internal/service/leo_video_async.go`: handles asynchronous Leo video requests and provider-neutral errors.
- `backend/internal/service/leo_video_async_test.go`: covers asynchronous provider-name sanitization.
- `backend/internal/service/leo_video_test.go`: covers synchronous provider-name sanitization.
- `backend/internal/service/model_pricing_resolver.go`: resolves channel video tiers and group fallback prices.
- `backend/internal/service/model_pricing_resolver_test.go`: covers channel tier resolution and explicit zero prices.
- `backend/internal/service/openai_gateway_record_usage_test.go`: verifies synchronous channel-price precedence and video multipliers.
- `backend/internal/service/openai_gateway_usage.go`: applies channel video pricing to synchronous usage calculation.
- `backend/internal/service/video_job_billing.go`: freezes channel pricing and billing metadata in asynchronous snapshots.
- `backend/internal/service/video_job_billing_test.go`: covers asynchronous pricing precedence, fallback, snapshots, and settlement attribution.
- `backend/internal/service/video_job_runtime.go`: sanitizes terminal video job errors.
- `backend/internal/service/video_job_service.go`: sanitizes persisted video job creation errors.
- `backend/internal/service/wire.go`: injects the shared model pricing resolver.
- `docs/LEO_VIDEO_CHANNEL.md`: documents channel pricing, custom Keys, image contracts, API docs, and billing behavior.
- `frontend/src/components/admin/channel/PricingEntryCard.vue`: edits fixed Seedance resolution-tier prices.
- `frontend/src/components/admin/channel/__tests__/PricingEntryCard.spec.ts`: tests channel video pricing controls.
- `frontend/src/components/admin/channel/__tests__/types.spec.ts`: tests video pricing form normalization and validation.
- `frontend/src/components/admin/channel/types.ts`: defines video pricing defaults, normalization, and validation.
- `frontend/src/components/video/ApiCodeBlock.vue`: renders copyable API examples.
- `frontend/src/components/video/EndpointTitle.vue`: renders API endpoint headings.
- `frontend/src/components/video/SectionHeading.vue`: renders API documentation section headings.
- `frontend/src/components/video/VideoSectionTabs.vue`: links the generation workbench and API docs.
- `frontend/src/components/video/__tests__/ApiCodeBlock.spec.ts`: tests API example copying.
- `frontend/src/components/video/__tests__/VideoSectionTabs.spec.ts`: tests video page tab navigation.
- `frontend/src/constants/channel.ts`: registers video billing mode, Seedance models, and resolutions.
- `frontend/src/i18n/locales/en/admin/channels.ts`: adds English channel pricing labels.
- `frontend/src/i18n/locales/en/dashboard.ts`: adds English video workbench and API docs text.
- `frontend/src/i18n/locales/zh/admin/channels.ts`: adds Chinese channel pricing labels.
- `frontend/src/i18n/locales/zh/dashboard.ts`: adds Chinese video workbench and API docs text.
- `frontend/src/router/index.ts`: registers the authenticated API docs route.
- `frontend/src/views/admin/ChannelsView.vue`: exposes Leo channel pricing administration.
- `frontend/src/views/user/VideoApiDocsView.vue`: implements the API integration documentation page.
- `frontend/src/views/user/VideoGenerationView.vue`: supports saved/custom Keys and verified image payload contracts.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: tests documented endpoints and request examples.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: tests custom Keys, Key isolation, single-image, multi-image, upload, cancel, and download flows.
- `progress.md`: records this release integration and verification evidence.
- Rollback point: redeploy `v0.1.161-fy.4`; for source rollback after this release commit, run `git revert <release-commit-hash>` from the repository root.

## 2026-07-20 - Task: Fix embedded Leo local-image routing and verify production video generation
### What was done
- Fixed the embedded frontend middleware so `/internal/video-inputs/:token` reaches the loopback-only image handler instead of the SPA fallback.
- Added a regression test that requires the internal route to return the backend `image/jpeg` response unchanged.
- Deployed `0.1.161-fy.5-hotfix.2` to `api.fflink.top` with a timestamped binary backup and automatic rollback, then verified an uploaded JPEG byte-for-byte and completed one real 4-second Seedance video job.
- Confirmed the existing `/video-generation/api-docs` child page and workbench/API-docs tabs remain deployed and cover the complete public asynchronous workflow.
### Testing
- `backend: go test -p 2 -tags=unit -timeout 10m ./...`: passed.
- `backend: go test -tags=embed ./internal/web ./internal/handler ./internal/server/routes -count=1`: passed.
- `backend: go vet -tags=embed ./internal/web ./internal/handler ./internal/server/routes`: passed.
- `frontend: .\\node_modules\\.bin\\vite.cmd build`: passed with existing warnings.
- `frontend: .\\node_modules\\.bin\\vitest.cmd run src\\views\\user\\__tests__\\VideoApiDocsView.spec.ts src\\components\\video\\__tests__\\VideoSectionTabs.spec.ts --reporter=dot`: passed, 2 files and 2 tests.
- Production missing-token checks returned `404` locally and through `https://api.fflink.top`; `/health` returned `200`.
- Production upload `yMre6R-bQw2aboHEPQS6l1zfjWKHrR2Y` returned `image/jpeg`; source and downloaded sizes were `101224` bytes and both SHA256 values were `de90b579d2fd84dfad51c4dee7789f2832ce1f6eaacbf8a50167e5ddc8fefed5`.
- Production job `vidjob_iUSug795Mygsu-Kbvu_nZ6HyE-hhd9iv` transitioned `pending -> running -> completed`; authenticated content download returned `video/mp4`, `797507` bytes, SHA256 `227fd36d942aa815abc3a6c94ab5f96cad3789d5f993a678a0a8a9804000634f`, and settled at actual cost `0.3200000000`.
- Deployed binary SHA256 is `5fa793a1d01397514e03ffb6d4875e81492af67a9f76efe81427cb9ccdc814f8`; `git diff --check` passed before this log append.
### Notes
- `backend/internal/web/embed_on.go`: bypasses the embedded SPA for the internal video-input route.
- `backend/internal/web/embed_test.go`: covers backend image response preservation through the embedded middleware.
- `docs/LEO_VIDEO_CHANNEL.md`: documents the embed bypass requirement and production acceptance checks.
- `progress.md`: records implementation, local verification, production deployment, real-task evidence, and rollback instructions.
- The first `hotfix.1` candidate was automatically rolled back because a Windows CRLF migration payload did not match the existing database checksum; no migration or database checksum was changed. `hotfix.2` was built from a Git-archived LF source tree and matched the database migration checksum.
- The existing API Key initially returned `GROUP_NOT_ALLOWED`; the operator changed the Leo group from exclusive to public before the successful real test. This task did not modify authentication or group configuration.
- Production rollback: `install -m 0755 /opt/sub2api/sub2api.backup-before-video-input-hotfix.20260720T085733Z /opt/sub2api/sub2api && systemctl restart sub2api`.
- Source rollback point: `02561d62c5fd24c897685ed98ed6c94a0a3ec48c`; run `git restore --source=02561d62c5fd24c897685ed98ed6c94a0a3ec48c --worktree -- backend/internal/web/embed_on.go backend/internal/web/embed_test.go docs/LEO_VIDEO_CHANNEL.md progress.md` from the repository root.

## 2026-07-20 - Task: Prepare the embedded video-input hotfix release v0.1.161-fy.6
### What was done
- Re-audited the embedded frontend bypass and regression test, confirmed the change is limited to local video-input routing, and selected `v0.1.161-fy.6` as the next release.
- Confirmed the release contains no database, migration, authentication, homepage, Baota report, or unrelated worktree changes.
- Rebuilt the embedded application and verified the local release candidate reports `Sub2API 0.1.161-fy.6`.
### Testing
- `backend: go test -tags=embed ./internal/web ./internal/handler ./internal/server/routes -count=1`: passed.
- `backend: go vet -tags=embed ./internal/web ./internal/handler ./internal/server/routes`: passed.
- `backend: go test -p 2 -tags=unit -timeout 10m ./...`: passed for all backend packages.
- `frontend: .\\node_modules\\.bin\\vite.cmd build`: passed with existing dynamic-import and chunk-size warnings.
- `backend: go build -tags embed -trimpath -ldflags "-s -w -X main.Version=0.1.161-fy.6" -o ..\\.codex-run\\sub2api-v0.1.161-fy.6.exe ./cmd/server` and `--version`: passed.
- `backend: ..\\.codex-run\\bin\\golangci-lint.exe run --build-tags embed --concurrency 2 ./internal/web/...`: reported three pre-existing `errcheck` findings at `embed_on.go:246-248`; `git blame` attributes all three unchanged lines to `500f241692`, and this task does not modify them.
- `git diff --check`: passed before this progress append; final staged check follows.
### Notes
- `backend/internal/web/embed_on.go`: bypasses the embedded SPA for `/internal/video-inputs/`.
- `backend/internal/web/embed_test.go`: verifies backend image MIME and bytes survive embedded middleware routing.
- `docs/LEO_VIDEO_CHANNEL.md`: documents the embed routing requirement and production acceptance checks.
- `progress.md`: records the implementation, production evidence, release verification, residual lint gap, and rollback point.
- Rollback point: deploy `v0.1.161-fy.5`; for source rollback after release, run `git revert <v0.1.161-fy.6-commit>` from the repository root.

## 2026-07-20 - Task: Merge upstream v0.1.162 and release v0.1.162-fy.1
### What was done
- Merged the upstream `v0.1.162` release tag into the fork without taking post-release commits, and synchronized the embedded application base version to `0.1.162`.
- Resolved seven merge conflicts while preserving the fork's Telemetry homepage, token incentive settings, fork-aware update checks, in-app update/restart path, video generation and pricing, local video input handling, Prompt Audit behavior, and authentication cache behavior.
- Integrated upstream configurable client-IP handling, image-storage settings, Grok count-token estimation and cache behavior, updater lifecycle fixes, subscription display changes, and related frontend/deployment updates.
- Preserved both updater credentials: the fork's configured GitHub token remains available for private release downloads, while `UPDATE_GITHUB_TOKEN` remains an API-only fallback with redirect leakage protection.
- Added the missing Viper default for `update.github_token`, updated rollback API timeout expectations, regenerated Wire output, and aligned embedded static-resource tests with the upstream SVG logo migration.
### Testing
- `backend: go test -p 2 -tags=unit -timeout 10m ./...`: passed for all packages.
- `backend: go test ./internal/repository -run 'GitHubRelease' -count=1`: passed.
- `backend: go test ./internal/service -run 'TokenIncentive|SettingService|UpdateService|Video|ImageStorage|Grok.*CountTokens' -count=1`: passed.
- `backend: go test ./internal/handler ./internal/server/routes ./internal/server -run 'Leo|Video|CountTokens|Update|PromptAudit|APIContract|Forwarded|TokenIncentive' -count=1`: passed.
- `backend: go test -tags=embed ./internal/web ./internal/handler ./internal/server/routes -count=1`: passed.
- `backend: go vet -tags=embed ./internal/web ./internal/handler ./internal/server/routes`: passed.
- `backend: ..\\.codex-run\\bin\\golangci-lint.exe run --concurrency 2 ./...`: passed with `0 issues`.
- Frontend full Vitest suite, ESLint, `vue-tsc --noEmit`, and Vite production build: passed; Vite reported only the existing dynamic-import and large-chunk warnings.
- `backend: go build -tags embed -trimpath -ldflags "-s -w -X main.Version=0.1.162-fy.1" -o ..\\.codex-run\\sub2api-v0.1.162-fy.1.exe ./cmd/server` and `--version`: passed and reported `0.1.162-fy.1`.
- `git diff --check`, `git diff --cached --check`, merge-marker scan, and unmerged-path scan: passed before this log append.
### Notes
- `README.md`, `README_JA.md`: retain fork attribution and second-development notices while adopting the upstream SVG logo.
- `frontend/src/views/HomeView.vue`: retains the fork's Telemetry homepage instead of restoring the removed upstream homepage implementation.
- `backend/internal/repository/github_release_service.go`: combines fork/upstream release checks and credential handling without exposing API credentials across redirects.
- `backend/internal/server/routes/gateway.go`: combines the fork's Leo routing policy with upstream Grok local token estimation.
- `backend/internal/service/setting_service_update_test.go`: retains token-incentive coverage and adds upstream forwarded-client-IP coverage.
- `backend/internal/service/wire.go`, `backend/cmd/server/wire_gen.go`: retain video providers and add upstream image-storage settings providers.
- `backend/internal/config/config.go`: makes `UPDATE_GITHUB_TOKEN` reachable through Viper while retaining all upstream v0.1.162 configuration changes.
- `frontend/src/api/__tests__/admin.system.rollback.spec.ts`: verifies the upstream 15-minute update and rollback request timeout.
- `backend/internal/web/embed_test.go`: requests the shipped `logo.svg` asset and verifies its SVG content type.
- All remaining changed source, test, frontend, deployment, asset, and documentation paths are the staged upstream `v0.1.162` release delta; the authoritative per-file list is available with `git diff --name-status v0.1.161-fy.6..v0.1.162-fy.1` after tagging.
- `progress.md`: records the merge decisions, verification evidence, release target, and rollback point.
- No database migration file is included in this release delta.
- Rollback point: redeploy `v0.1.161-fy.6`; for source rollback after release, run `git revert -m 1 <v0.1.162-fy.1-merge-commit>` from the repository root.

## 2026-07-21 - Task: Fix video workbench preview and add explicit start/end frames
### What was done
- Stabilized completed-video previews so two-second job polling no longer revokes and downloads the same Blob URL repeatedly; switching jobs or Keys still invalidates the previous preview, and media failures now expose a working retry action.
- Changed local reference selection to append one image at a time in order, added independent local and remote start/end-frame controls, and uploaded the reference/start/end file snapshot concurrently before mapping it to `guidances.image_reference`, `start_frame_url`, and `end_frame_url`.
- Extended the existing video API documentation child page and bilingual text with the explicit start/end-frame contract and a combined request example.
- Deployed the embedded frontend as `0.1.162-fy.2-hotfix.2` to `api.fflink.top` with automatic rollback and two timestamped binary backups; no database, migration, authentication, group, or Leo configuration was changed.
### Testing
- `frontend: vitest run src/views/user/__tests__/VideoGenerationView.spec.ts src/views/user/__tests__/VideoApiDocsView.spec.ts src/api/__tests__/videoGeneration.spec.ts`: passed, 3 files and 17 tests.
- Frontend full Vitest suite passed; changed-file ESLint, `vue-tsc --noEmit`, and Vite production build passed. Vite reported only the existing dynamic-import and large-chunk warnings.
- The preview regression test proves refreshing an unchanged completed job keeps `downloadVideoOutput` at one call and does not revoke the active Blob URL; the upload regression proves all three upload promises start before task creation and verifies every output field.
- Server output `vidjob_jMrHMC4tCe6EqFqkmcVBvgb9PoVMXZP4.mp4` passed `ffprobe` and full `ffmpeg` decode: H.264 High, `yuv420p`, 720x1280, 24 fps, 15.04 seconds, 4,484,597 bytes.
- LF-archive candidate `go test -tags=embed ./internal/web -count=1` passed before both production builds. Final deployed SHA256 is `464b76016c0f956f4a5dbb78e9f04a2ba0e87802b59581eac136570ab9b8ed69`.
- Post-deploy service state is `active`; local and public `/health` returned `200`. Public `VideoGenerationView-D0z6THuP.js` returned `200` and contains both `start_frame_url` and `end_frame_url`.
- The in-app browser reached the production login redirect but had no authenticated session, so authenticated click-through was not claimed as verified. Server media, API payload tests, production assets, and health checks were verified independently.
- `git diff --check`: passed before this log append.
### Notes
- `frontend/src/views/user/VideoGenerationView.vue`: stabilizes preview ownership, accumulates local references, adds start/end-frame inputs, snapshots concurrent uploads, and maps the API payload.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: covers sequential references, concurrent frame/reference uploads, field mapping, and stable polling previews.
- `frontend/src/api/__tests__/videoGeneration.spec.ts`: verifies start/end-frame fields are serialized unchanged by the API client.
- `frontend/src/views/user/VideoApiDocsView.vue`: documents explicit frame fields and the combined frame/reference request body.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: verifies the child page includes both explicit frame fields.
- `frontend/src/i18n/locales/zh/dashboard.ts`: adds Chinese frame, retry, and API documentation text.
- `frontend/src/i18n/locales/en/dashboard.ts`: adds English frame, retry, and API documentation text.
- `docs/LEO_VIDEO_CHANNEL.md`: documents one-at-a-time reference selection and concurrent frame/reference submission.
- `progress.md`: records implementation, verification, deployment, file ownership, and rollback evidence.
- Production rollback: `install -m 0755 /opt/sub2api/sub2api.backup-before-video-ui-hotfix.20260721T015653Z /opt/sub2api/sub2api && systemctl restart sub2api`. The earlier baseline backup is `/opt/sub2api/sub2api.backup-before-video-ui-hotfix.20260721T013956Z`.
- Source rollback: from the repository root, run `git restore --worktree -- docs/LEO_VIDEO_CHANNEL.md frontend/src/api/__tests__/videoGeneration.spec.ts frontend/src/i18n/locales/en/dashboard.ts frontend/src/i18n/locales/zh/dashboard.ts frontend/src/views/user/VideoApiDocsView.vue frontend/src/views/user/VideoGenerationView.vue frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts frontend/src/views/user/__tests__/VideoGenerationView.spec.ts progress.md`.

## 2026-07-21 - Task: Release video workbench fixes as v0.1.162-fy.2
### What was done
- Reviewed the complete local video-workbench delta and confirmed it contains the stable preview fix, explicit start/end frames, ordered reference uploads, API documentation, bilingual copy, and regression coverage as one release unit.
- Selected `v0.1.162-fy.2` as the next fork release and excluded the untracked `.superpowers/` directory and local build artifacts from the commit.
- Confirmed the release contains no backend protocol, database migration, authentication, homepage, group, pricing, or deployment configuration changes.
### Testing
- `frontend: .\\node_modules\\.bin\\vitest.cmd run --reporter=dot`: passed for the full frontend suite, including 12 video-workbench tests and 4 video API-client tests.
- Changed-file ESLint and `.\\node_modules\\.bin\\vue-tsc.cmd --noEmit`: passed.
- `frontend: .\\node_modules\\.bin\\vite.cmd build`: passed with only the existing dynamic-import, stale Browserslist-data, and large-chunk warnings.
- `backend: go test -tags=embed ./internal/web ./internal/handler ./internal/server/routes -count=1`: passed.
- `backend: go vet -tags=embed ./internal/web ./internal/handler ./internal/server/routes`: passed.
- `backend: ..\\.codex-run\\bin\\golangci-lint.exe run --concurrency 2 ./...`: passed with `0 issues`.
- `backend: go build -tags embed -trimpath -ldflags "-s -w -X main.Version=0.1.162-fy.2" -o ..\\.codex-run\\sub2api-v0.1.162-fy.2.exe ./cmd/server` and `--version`: passed and reported `0.1.162-fy.2`.
- `git diff --check`: passed before this log append; final staged checks follow.
### Notes
- `frontend/src/views/user/VideoGenerationView.vue`: contains the preview lifecycle fix and explicit frame/reference controls and payload mapping.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: covers stable polling previews and concurrent ordered uploads.
- `frontend/src/api/__tests__/videoGeneration.spec.ts`: verifies frame fields remain intact in API requests.
- `frontend/src/views/user/VideoApiDocsView.vue`: documents combined frame and reference requests.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: verifies frame fields appear in API examples.
- `frontend/src/i18n/locales/zh/dashboard.ts`: provides Chinese video-workbench and API-documentation copy.
- `frontend/src/i18n/locales/en/dashboard.ts`: provides English video-workbench and API-documentation copy.
- `docs/LEO_VIDEO_CHANNEL.md`: documents ordered local references and concurrent frame submission.
- `progress.md`: records implementation evidence and formal release verification.
- Rollback point: redeploy `v0.1.162-fy.1`; for source rollback after release, run `git revert <v0.1.162-fy.2-commit>` from the repository root.

## 2026-07-21 - Task: Resolve the v0.1.162-fy.2 security gate and release v0.1.162-fy.3
### What was done
- Investigated both failed `v0.1.162-fy.2` security runs and traced them to the newly published high-severity axios advisory `GHSA-gcfj-64vw-6mp9`; backend `govulncheck` remained successful.
- Upgraded the direct axios dependency from `1.16.0` to the fixed `1.18.1` release instead of adding or extending an audit exception.
- Regenerated the dependency lock with CI-compatible pnpm 9, retained existing Rollup platform constraints, and removed package-manager metadata churn unrelated to axios.
- Kept the already-public `v0.1.162-fy.2` tag immutable and selected `v0.1.162-fy.3` for the auditable corrective release.
### Testing
- `pnpm@9.15.9 install --lockfile-only --frozen-lockfile` with an isolated modules directory: passed.
- `pnpm audit --prod --audit-level=high --json` plus `tools/check_pnpm_audit_exceptions.py`: passed; high/critical findings are limited to the two existing, unexpired `xlsx` exceptions and axios has no remaining high finding.
- `frontend: .\\node_modules\\.bin\\vue-tsc.cmd --noEmit`: passed.
- `frontend: .\\node_modules\\.bin\\vitest.cmd run --reporter=dot --maxWorkers=2 --minWorkers=1`: passed for the full suite. The first concurrent run was discarded after the local esbuild service and workers were terminated without an assertion failure.
- `frontend: .\\node_modules\\.bin\\vite.cmd build`: passed with only the existing dynamic-import, stale Browserslist-data, and large-chunk warnings.
- `backend: go test -tags=embed ./internal/web ./internal/handler ./internal/server/routes -count=1`: passed.
- `backend: go vet -tags=embed ./internal/web ./internal/handler ./internal/server/routes`: passed.
- `backend: go build -tags embed -trimpath -ldflags "-s -w -X main.Version=0.1.162-fy.3" -o ..\\.codex-run\\sub2api-v0.1.162-fy.3.exe ./cmd/server` and `--version`: passed and reported `0.1.162-fy.3`.
- `git diff --check`: passed before this log append; final staged checks follow.
### Notes
- `frontend/package.json`: raises the direct axios range to `^1.18.1`.
- `frontend/pnpm-lock.yaml`: locks axios `1.18.1` and only its required new proxy-agent dependency graph.
- `progress.md`: records the failed security signal, remediation decision, validation evidence, and corrective release target.
- Emergency rollback point: redeploy `v0.1.162-fy.2`; this restores axios `1.16.0` and reopens `GHSA-gcfj-64vw-6mp9`, so rollback is not recommended except to isolate an axios compatibility regression. Source rollback after release is `git revert <v0.1.162-fy.3-commit>`.

## 2026-07-21 - Task: Prevent mixed reference and frame inputs in the video workbench
### What was done
- Made local-file and remote-URL reference inputs mutually exclusive with start/end frames: selecting either mode disables the conflicting controls without silently clearing the user's current selection.
- Added file-handler guards, submit eligibility validation, and a final submission guard so programmatic state changes cannot upload or create a mixed request.
- Kept both valid workflows unchanged: up to four ordered reference images, or a start frame and end frame uploaded concurrently without reference guidance.
- Corrected the bilingual API child page and Leo channel documentation so examples and field descriptions no longer recommend combining frames with reference images.
### Testing
- `frontend: .\node_modules\.bin\vitest.cmd run src/views/user/__tests__/VideoGenerationView.spec.ts src/views/user/__tests__/VideoApiDocsView.spec.ts src/api/__tests__/videoGeneration.spec.ts`: passed, 3 files and 20 tests.
- The workbench tests cover four local references disabling both frame inputs, a local frame disabling and rejecting reference selection, remote reference URLs disabling frame URLs, a programmatic mixed URL state producing no upload or job creation, valid ordered references, and valid concurrent start/end-frame uploads.
- `frontend: .\node_modules\.bin\vue-tsc.cmd --noEmit`: passed.
- Changed-file ESLint passed for the workbench, API docs, tests, and bilingual locale files.
- `frontend: npm run build`: passed; Vite reported only the existing stale Browserslist data, dynamic-import, and large-chunk warnings.
- `git diff --check`: passed before this log append, with only the existing LF-to-CRLF notice for `docs/LEO_VIDEO_CHANNEL.md`.
### Notes
- `frontend/src/views/user/VideoGenerationView.vue`: disables conflicting local and URL controls and rejects mixed frame/reference state at selection and submission time.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: replaces the invalid mixed-upload expectation and adds bidirectional local/URL exclusivity regression coverage.
- `frontend/src/views/user/VideoApiDocsView.vue`: changes the frame-pair example to contain only start and end frames.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: asserts no example combines frame fields with `image_reference`.
- `frontend/src/i18n/locales/zh/dashboard.ts`: adds Chinese exclusivity guidance and corrects API field descriptions.
- `frontend/src/i18n/locales/en/dashboard.ts`: adds English exclusivity guidance and corrects API field descriptions.
- `docs/LEO_VIDEO_CHANNEL.md`: documents the separate frame and reference-image request modes.
- `progress.md`: records implementation, verification, file ownership, and rollback evidence.
- No backend protocol, database, pricing, authentication, deployment, or paid video-generation action was performed.
- Source rollback: from the repository root, run `git restore --worktree -- docs/LEO_VIDEO_CHANNEL.md frontend/src/i18n/locales/en/dashboard.ts frontend/src/i18n/locales/zh/dashboard.ts frontend/src/views/user/VideoApiDocsView.vue frontend/src/views/user/VideoGenerationView.vue frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts frontend/src/views/user/__tests__/VideoGenerationView.spec.ts progress.md`.

## 2026-07-21 - Task: Prepare v0.1.162-fy.4 release
### What was done
- Locked the corrective release version to `v0.1.162-fy.4` on `codex/leo-video-channel` for the video image-input exclusivity fix.
- Revalidated the embedded backend path and produced a local release binary carrying version `0.1.162-fy.4` before publishing the branch and tag.
### Testing
- `backend: go test -tags=embed ./internal/web ./internal/handler ./internal/server/routes -count=1`: passed.
- `backend: go vet -tags=embed ./internal/web ./internal/handler ./internal/server/routes`: passed.
- `backend: go build -tags embed -trimpath -ldflags "-s -w -X main.Version=0.1.162-fy.4" -o ..\\.codex-run\\sub2api-v0.1.162-fy.4.exe ./cmd/server`: passed.
- `backend: ..\\.codex-run\\sub2api-v0.1.162-fy.4.exe --version`: reported `Sub2API 0.1.162-fy.4`.
- `frontend: pnpm audit --audit-level=high`: passed with no high-severity advisory output.
- `git diff --check`: passed; Git only reported the existing LF-to-CRLF checkout notice for two Markdown files.
### Notes
- `progress.md`: records the release version, backend embed validation, security audit, and rollback point.
- The release contains only the eight files listed in the preceding implementation record; `.superpowers/` and local build output are excluded.
- Release rollback point: redeploy tag `v0.1.162-fy.3`; source rollback after publication is `git revert v0.1.162-fy.4`.

## 2026-07-21 - Task: Fix video preview playback and preserve downloads
### What was done
- Identified the production-only playback failure as a CSP mismatch: the workbench downloads completed videos into Blob URLs, while the page policy omitted `media-src blob:` and therefore blocked the `<video>` element.
- Added the narrow `media-src 'self' blob:` permission to the default policy and the security middleware's required directives so existing custom CSP configurations are upgraded without relaxing script, connection, frame, or authentication rules.
- Changed the video workbench so a media playback error no longer hides the already downloaded file; the download action remains available and the retry action now revokes the old Blob and fetches the completed output again.
### Testing
- `backend: go test ./internal/server/middleware ./internal/config -count=1`: passed.
- `frontend: .\node_modules\.bin\vitest.cmd run src/views/user/__tests__/VideoGenerationView.spec.ts src/api/__tests__/videoGeneration.spec.ts`: passed, 2 files and 20 tests.
- The new frontend regression proves that a playback error keeps the completed MP4 download visible and that retry performs a second content request after revoking the previous Blob URL.
- `frontend: .\node_modules\.bin\vue-tsc.cmd --noEmit` and changed-file ESLint: passed.
- `frontend: npm run build`: passed with only the existing stale Browserslist data, dynamic-import, and large-chunk warnings.
- `backend: go test -tags=embed ./internal/web ./internal/server/middleware ./internal/config -count=1`: passed.
- `backend: go build -tags embed -trimpath -o ..\.codex-run\sub2api-video-preview-fix.exe ./cmd/server`: passed.
- `git diff --check`: passed before this log append, with only the existing LF-to-CRLF notice for `docs/LEO_VIDEO_CHANNEL.md`.
### Notes
- `backend/internal/config/config.go`: adds Blob media playback to the default CSP policy.
- `backend/internal/server/middleware/security_headers.go`: automatically injects the required Blob media source into old custom CSP policies.
- `backend/internal/server/middleware/security_headers_test.go`: verifies default, missing, and already-present media directives.
- `deploy/config.example.yaml`: documents the production CSP media directive.
- `frontend/src/views/user/VideoGenerationView.vue`: preserves downloads after playback failure and makes retry refetch the output.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: covers download retention and actual retry behavior.
- `docs/LEO_VIDEO_CHANNEL.md`: documents the Blob media CSP requirement and failure behavior.
- `progress.md`: records root cause, implementation, verification, file ownership, and rollback evidence.
- No video generation, billing, database, authentication, Leo request, or production deployment action was performed.
- Source rollback: from the repository root, run `git restore --worktree -- backend/internal/config/config.go backend/internal/server/middleware/security_headers.go backend/internal/server/middleware/security_headers_test.go deploy/config.example.yaml docs/LEO_VIDEO_CHANNEL.md frontend/src/views/user/VideoGenerationView.vue frontend/src/views/user/__tests__/VideoGenerationView.spec.ts progress.md`.

## 2026-07-21 - Task: Prepare v0.1.162-fy.5 release
### What was done
- Locked the video preview and CSP corrective release to `v0.1.162-fy.5` on `codex/leo-video-channel`.
- Revalidated the complete frontend suite, the embedded backend path, and the dependency security gate before preparing the release commit and tag.
### Testing
- `frontend: .\\node_modules\\.bin\\vitest.cmd run --reporter=dot --maxWorkers=2 --minWorkers=1`: passed.
- `frontend: pnpm run build`: passed; output contained only existing non-blocking warnings.
- `frontend: pnpm audit --audit-level=high`: passed with no high-severity advisory output.
- `backend: go test -tags=embed ./internal/web ./internal/server/middleware ./internal/server/routes -count=1`: passed.
- `backend: go vet -tags=embed ./internal/web ./internal/server/middleware ./internal/server/routes`: passed.
- `backend: go build -tags embed -trimpath -ldflags "-s -w -X main.Version=0.1.162-fy.5" -o ..\\.codex-run\\sub2api-v0.1.162-fy.5.exe ./cmd/server`: passed.
- `backend: ..\\.codex-run\\sub2api-v0.1.162-fy.5.exe --version`: reported `Sub2API 0.1.162-fy.5`.
### Notes
- `progress.md`: records the release version, final validation evidence, scope, and rollback point.
- The release contains only the eight files listed in this and the preceding implementation record; `.superpowers/` and `.codex-run/` remain excluded.
- Release rollback point: redeploy tag `v0.1.162-fy.4`; source rollback after publication is `git revert v0.1.162-fy.5`.
## 2026-07-23 - Task: Fix video production API key quota settlement

### What was done
- Fixed the asynchronous Leo video settlement path so a successfully completed video forwards the configured API key quota updater into the shared billing recorder.
- Preserved the existing settlement rules: quota and rate-limit usage update only after a valid completed output is billed; failed, canceled, or incomplete jobs do not consume the key quota.
- Regenerated the Wire server injector and documented the API key quota behavior for the Leo video channel.

### Testing
- `backend: go test ./internal/service -count=1`: passed.
- `backend: go test ./internal/handler -run 'TestLeoVideo|Test.*Video' -count=1`: passed.
- `backend: go test ./internal/repository -run 'TestVideoJob|Test.*UsageBilling' -count=1`: passed.
- `backend: go test ./cmd/server -run '^$' -count=1`: passed.
- `backend: go generate ./cmd/server`: passed; Wire regenerated `backend/cmd/server/wire_gen.go`.
- `backend: go build -tags embed -trimpath -o ..\\.codex-run\\sub2api-video-quota-fix.exe ./cmd/server`: passed.
- `git diff --check`: passed; Git reported only the existing LF-to-CRLF checkout notice for `docs/LEO_VIDEO_CHANNEL.md`.

### Notes
- `backend/internal/service/video_job_billing.go`: carries `APIKeyQuotaUpdater` into completed-video usage recording.
- `backend/internal/service/wire.go`: injects the application `APIKeyService` into video billing.
- `backend/cmd/server/wire_gen.go`: regenerated provider call with the API key service dependency.
- `backend/internal/service/video_job_billing_test.go`: verifies video settlement forwards the quota updater while remaining idempotent.
- `docs/LEO_VIDEO_CHANNEL.md`: documents successful video settlement updating a custom API key's `quota_used`.
- `progress.md`: records this implementation and verification evidence.
- Source rollback: from the repository root, run `git restore --worktree -- backend/cmd/server/wire_gen.go backend/internal/service/video_job_billing.go backend/internal/service/video_job_billing_test.go backend/internal/service/wire.go docs/LEO_VIDEO_CHANNEL.md progress.md`.
- No production deployment, database migration, pricing change, authentication change, or paid video-generation action was performed.


## 2026-07-23 - Task: Merge upstream v0.1.163 and prepare v0.1.163-fy.1
### What was done
- Merged the official upstream tag `v0.1.163` at commit `d0bdd7e771636a8d315f542cafd39484f39bd60c` with a non-fast-forward merge; the three later untagged `upstream/main` commits were intentionally excluded.
- Preserved the fork's Leo video pricing validation, fork container image, and axios `1.18.1` security floor while integrating upstream OpenAI reasoning policy, Redis ACL, gateway, billing, scheduler, and mobile-layout changes.
- Resolved five conflicts in the group service, Docker Compose, frontend dependency manifest, pnpm lockfile, and group administration page without dropping either valid feature path.
### Testing
- `frontend: pnpm install --frozen-lockfile --offline`: passed.
- Conflict-focused frontend Vitest passed: 4 files and 26 tests covering Leo image/video pricing, OpenAI reasoning policy, model lists, and group duplication.
- `backend: go test ./internal/service -run 'Test.*Group|Test.*Leo.*Price|TestValidateLeoVideoPrices' -count=1`: passed.
- Go `1.26.5` matched CI; `go test -tags=unit ./...` and `go test -tags=integration ./...` both passed.
- `frontend: pnpm run lint:check`, `pnpm run typecheck`, the complete Vitest suite, and `pnpm run build`: passed.
- The repository frontend audit exception checker passed against a fresh production `pnpm audit --audit-level=high` JSON report.
- `backend: go vet ./...`: passed.
- The embedded release build reported `Sub2API 0.1.163-fy.1` with upstream commit `d0bdd7e771636a8d315f542cafd39484f39bd60c`.
- `git diff --check`: passed.
### Notes
- Upstream migration `185_group_reasoning_effort_policy.sql` adds group reasoning-policy fields; deployment must allow normal startup migrations before serving traffic.
- Rolling the binary back to `v0.1.162-fy.5` leaves the additive migration columns in place and does not require destructive schema rollback.
- `.superpowers/`, `.codex-run/`, audit output, and local binaries remain excluded.
- Changed files:
- `README.md`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/cmd/server/main.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/ent/group.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/ent/group/group.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/ent/group/where.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/ent/group_create.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/ent/group_update.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/ent/migrate/schema.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/ent/mutation.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/ent/runtime/runtime.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/ent/schema/group.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/go.mod`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/go.sum`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/config/config.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/config/config_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/domain/reasoning_effort.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/handler/admin/group_handler.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/handler/admin/group_handler_reasoning_effort_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/handler/admin/usage_handler.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/handler/admin/usage_handler_request_type_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/handler/dto/mappers.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/handler/dto/types.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/handler/gateway_handler.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/handler/openai_chat_completions.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/handler/openai_gateway_handler.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/pkg/apicompat/responses_client_tools.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/pkg/apicompat/responses_client_tools_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/repository/api_key_repo.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/repository/group_repo.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/repository/proxy_probe_service.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/repository/proxy_probe_service_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/repository/redis.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/repository/redis_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/repository/scheduler_cache.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/repository/scheduler_cache_last_used_unit_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/repository/scheduler_cache_unit_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/server/api_contract_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/account_test_service.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/account_test_service_openai_compact_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/account_test_service_openai_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/admin_group.go`: merged Leo video price validation with the upstream OpenAI reasoning-policy normalization and sanitization.
- `backend/internal/service/admin_group_duplicate.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/admin_group_duplicate_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/admin_service.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/admin_service_group_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/api_key_auth_cache.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/api_key_auth_cache_impl.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/api_key_auth_cache_version_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/api_key_service_cache_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/channel_monitor_runner.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/channel_monitor_runner_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/gateway_anthropic_passthrough.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/gateway_forward_as_chat_completions.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/gateway_forward_as_chat_completions_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/gateway_forward_as_responses.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/gateway_forward_as_responses_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/gateway_non_streaming_response_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/grok_media.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/grok_upstream_errors.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/grok_upstream_errors_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/group.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_account_runtime_block_fastpath.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_account_scheduler.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_account_scheduler_compact_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_account_scheduler_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_account_scheduler_upstream_cost_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_apikey_responses_probe.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_apikey_responses_probe_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_codex_identity.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_codex_identity_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_codex_transform.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_codex_transform_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_cc_pipeline.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_chat_completions_raw.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_forward.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_grok.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_grok_cache.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_grok_cache_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_grok_chat_bridge.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_grok_chat_bridge_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_grok_compact.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_grok_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_grok_tool_protocol.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_grok_tool_protocol_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_response_handling.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_response_handling_image_usage_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_scheduling.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_service.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_service_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_gateway_upstream_errors.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_oauth_passthrough_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_reasoning_effort_policy.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_reasoning_effort_policy_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_responses_namespace.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_responses_namespace_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_ws_forwarder.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_ws_forwarder_ingress.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_ws_forwarder_payload.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_ws_forwarder_success_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_ws_http_bridge.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/openai_ws_v2_passthrough_adapter.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/scheduler_snapshot_service.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/upstream_models.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/service/upstream_models_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/setup/handler.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/setup/setup.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/internal/setup/setup_test.go`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `backend/migrations/185_group_reasoning_effort_policy.sql`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `deploy/.env.example`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `deploy/config.example.yaml`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `deploy/docker-compose.dev.yml`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `deploy/docker-compose.local.yml`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `deploy/docker-compose.standalone.yml`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `deploy/docker-compose.yml`: kept the fork image target while accepting upstream Redis ACL username support.
- `docs/PAYMENT.md`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `docs/PAYMENT_CN.md`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `docs/screenshots/mobile-account-actions-menu.png`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/api/setup.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/admin/group/GroupRPMOverridesModal.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/admin/group/GroupRateMultipliersModal.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/admin/group/ReasoningEffortPolicyFields.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/admin/usage/UsageFilters.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/admin/usage/__tests__/UsageFilters.spec.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/charts/EndpointDistributionChart.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/charts/GroupDistributionChart.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/charts/ModelDistributionChart.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/common/Select.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/layout/AppHeader.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/payment/SubscriptionPlanCard.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/payment/__tests__/SubscriptionPlanCard.spec.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/payment/__tests__/validity.spec.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/payment/validity.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/components/user/dashboard/UserDashboardCharts.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/i18n/locales/en/admin/ops.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/i18n/locales/en/admin/overview.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/i18n/locales/en/landing.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/i18n/locales/en/misc.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/i18n/locales/zh/admin/accounts.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/i18n/locales/zh/admin/ops.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/i18n/locales/zh/admin/overview.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/i18n/locales/zh/landing.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/i18n/locales/zh/misc.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/main.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/types/index.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/utils/__tests__/device.spec.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/utils/__tests__/floatingPanel.spec.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/utils/__tests__/formatMultiplier.spec.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/utils/device.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/utils/floatingPanel.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/utils/formatters.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/admin/AccountsView.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/admin/GroupsView.vue`: runs both Leo video price completeness checks and upstream reasoning-policy form validation.
- `frontend/src/views/admin/PromoCodesView.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/admin/SettingsView.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/admin/__tests__/groupsReasoningEffort.spec.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/admin/groupsReasoningEffort.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/admin/ops/components/OpsAlertEventsCard.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/admin/ops/components/OpsAlertRulesCard.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/admin/ops/components/OpsDashboardSkeleton.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/admin/ops/components/OpsErrorDetailsModal.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/admin/ops/components/OpsOpenAITokenStatsCard.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/admin/ops/components/OpsRequestDetailsModal.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/admin/ops/components/OpsSystemLogTable.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/admin/ops/components/OpsThroughputTrendChart.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/admin/ops/components/__tests__/OpsThroughputTrendChart.spec.ts`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/setup/SetupWizardView.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/user/KeysView.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `frontend/src/views/user/PaymentView.vue`: imported from the signed upstream v0.1.163 release tag without a manual conflict.
- `progress.md`: records the upstream boundary, conflict decisions, verification evidence, changed-file inventory, migration note, and rollback point.
- Source rollback immediately after this merge commit: `git revert -m 1 HEAD`. Release rollback point: redeploy `v0.1.162-fy.5`.

## 2026-07-23 - Task: Constrain video workbench parameters by LeoStudio model capability
### What was done
- Added a model capability matrix to the video workbench so each model only exposes supported resolutions, discrete durations, and resolution-specific aspect ratios.
- Restricted `seedance-2.0` to the production-verified `480p` and `720p` modes; kept `480p`, `720p`, and `1080p` for `seedance-2.0-fast`.
- Reset model parameters on model changes, repaired unsupported aspect ratios after resolution changes, and added a pre-upload submission guard with a fixed request parameter snapshot.
- Updated the workbench API documentation text and Leo channel documentation with the exact parameter matrix and the distinction between UI constraints and direct API callers.

### Testing
- From `frontend/`, `.\node_modules\.bin\vitest.cmd run src/views/user/__tests__/VideoGenerationView.spec.ts src/api/__tests__/videoGeneration.spec.ts` passed: 2 files, 25 tests.
- From `frontend/`, `.\node_modules\.bin\vue-tsc.cmd --noEmit` passed.
- Changed-file ESLint passed for the workbench component, its test, and the Chinese and English dashboard locale files.
- From `frontend/`, `npm run build` passed; Vite reported only the repository's existing chunking and bundle-size warnings.
- `git diff --check` passed. No paid video generation or production deployment was performed.

### Notes
- `frontend/src/views/user/VideoGenerationView.vue`: added model-specific option filtering, defaults, parameter locking, and the final submission guard.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: added regression coverage for model switching, discrete durations, aspect fallback, invalid-state blocking, and exact request payloads.
- `frontend/src/i18n/locales/zh/dashboard.ts`: documented the Chinese parameter matrix and added the invalid-combination error message.
- `frontend/src/i18n/locales/en/dashboard.ts`: documented the English parameter matrix and added the invalid-combination error message.
- `docs/LEO_VIDEO_CHANNEL.md`: recorded the production workbench capability matrix and fallback behavior.
- `progress.md`: appended this implementation, verification, changed-file inventory, and rollback record.
- Rollback point: `2050e1ed9`. Before commit, run `git restore --source=2050e1ed9 -- docs/LEO_VIDEO_CHANNEL.md frontend/src/i18n/locales/en/dashboard.ts frontend/src/i18n/locales/zh/dashboard.ts frontend/src/views/user/VideoGenerationView.vue frontend/src/views/user/__tests__/VideoGenerationView.spec.ts progress.md`.

## 2026-07-23 - Task: Realign video workbench parameters with the latest LeoStudio documentation
### What was done
- Fetched LeoStudio remote references without touching its existing uncommitted work; `origin/main` remains at `7af385af0fd8996dab1853c8ec965d4c1179bb08`.
- Compared the locally updated Sub2 integration document with the current LeoStudio model registry and restored `1080p` for `seedance-2.0`, so both Seedance models now expose the registered `480p`, `720p`, and `1080p` modes.
- Kept the registered defaults and constraints: `480p`, 8 seconds, discrete durations from 4 through 15, all seven aspect ratios at 480p/1080p, and no `9:21` at 720p.
- Documented the updated guidance contract: direct API callers may omit image-reference `order`, while the workbench continues to send explicit continuous order values.

### Testing
- From `frontend/`, `.\node_modules\.bin\vitest.cmd run src/views/user/__tests__/VideoGenerationView.spec.ts src/api/__tests__/videoGeneration.spec.ts` passed: 2 files, 25 tests.
- From `frontend/`, `.\node_modules\.bin\vue-tsc.cmd --noEmit` passed.
- Changed-file ESLint passed for the workbench component, its test, and both dashboard locale files.
- From `frontend/`, `npm run build` passed; Vite reported only the repository's existing chunking and bundle-size warnings.
- `git diff --check` passed. No paid 1080p generation or production deployment was performed.

### Notes
- `frontend/src/views/user/VideoGenerationView.vue`: aligned the standard Seedance model with the current three-resolution registry and its 1080p aspect map.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: updated resolution expectations and retained an invalid-duration bypass test for the final submission guard.
- `frontend/src/i18n/locales/zh/dashboard.ts`: aligned Chinese API documentation text with the current resolution and reference-order contract.
- `frontend/src/i18n/locales/en/dashboard.ts`: aligned English API documentation text with the current resolution and reference-order contract.
- `docs/LEO_VIDEO_CHANNEL.md`: replaced the retired production exception with the current LeoStudio registry and documented omitted reference orders.
- `progress.md`: appended this follow-up's evidence, changed-file inventory, and rollback record.
- Full uncommitted feature rollback point: `2050e1ed9`. Run `git restore --source=2050e1ed9 -- docs/LEO_VIDEO_CHANNEL.md frontend/src/i18n/locales/en/dashboard.ts frontend/src/i18n/locales/zh/dashboard.ts frontend/src/views/user/VideoGenerationView.vue frontend/src/views/user/__tests__/VideoGenerationView.spec.ts progress.md`.

## 2026-07-23 - Task: Correct Fast resolution options from LeoStudio feat/web-admin
### What was done
- Rechecked all LeoStudio remote branches and found the actual latest capability update at `origin/feat-web-admin` commit `2fd5c21b01a049817962812cf4675ade7727cc12` (`feat: align video models with Leonardo specs`, 2026-07-23 11:30 +08:00).
- Corrected the workbench matrix to match that release: `seedance-2.0` supports `480p/720p/1080p`, while `seedance-2.0-fast` supports only `480p/720p`; both default to `720p`.
- Kept the existing strict duration and aspect-ratio validation, and retained the current page scope of the two already-configured Seedance models rather than exposing new Mini, Happy Horse, or Grok models without matching Sub2 pricing and request contracts.
- Updated the Chinese and English API help text and Leo channel documentation to identify the exact source commit and Fast limitation.

### Testing
- From `frontend/`, `.\node_modules\.bin\vitest.cmd run src/views/user/__tests__/VideoGenerationView.spec.ts src/api/__tests__/videoGeneration.spec.ts` passed: 2 files, 25 tests.
- From `frontend/`, `.\node_modules\.bin\vue-tsc.cmd --noEmit` passed.
- Changed-file ESLint passed for the workbench component, its test, and both dashboard locale files.
- From `frontend/`, `npm run build` passed; Vite reported only the repository's existing chunking and bundle-size warnings.
- `git diff --check` passed. No paid video generation or production deployment was performed.

### Notes
- `frontend/src/views/user/VideoGenerationView.vue`: removed Fast `1080p` and aligned both model defaults to `720p`.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: verifies Fast options are exactly `480p/720p`, ordinary Seedance retains `1080p`, and supported payloads remain valid.
- `frontend/src/i18n/locales/zh/dashboard.ts`: documented the corrected Fast resolution list.
- `frontend/src/i18n/locales/en/dashboard.ts`: documented the corrected Fast resolution list.
- `docs/LEO_VIDEO_CHANNEL.md`: recorded the `2fd5c21b` capability matrix and `720p` defaults.
- `progress.md`: appended this correction, verification evidence, changed-file inventory, and rollback record.
- Full uncommitted feature rollback point: `2050e1ed9`. Run `git restore --source=2050e1ed9 -- docs/LEO_VIDEO_CHANNEL.md frontend/src/i18n/locales/en/dashboard.ts frontend/src/i18n/locales/zh/dashboard.ts frontend/src/views/user/VideoGenerationView.vue frontend/src/views/user/__tests__/VideoGenerationView.spec.ts progress.md`.

## 2026-07-23 - Task: Release LeoStudio parameter constraints as v0.1.163-fy.2
### What was done
- Prepared the current six-file LeoStudio parameter-constraint change for the `v0.1.163-fy.2` release.
- Kept the release scoped to the current branch and excluded the unrelated untracked `.superpowers/` directory.

### Testing
- From `frontend/`, Vitest passed: 2 files, 25 tests.
- From `frontend/`, `vue-tsc --noEmit` passed.
- Changed-file ESLint passed for the workbench component, its test, and both dashboard locale files.
- From `frontend/`, `pnpm run build` passed.
- `git diff --check` passed.

### Notes
- `docs/LEO_VIDEO_CHANNEL.md`: included the final LeoStudio capability documentation in the release.
- `frontend/src/i18n/locales/en/dashboard.ts`: included the English workbench/API guidance changes.
- `frontend/src/i18n/locales/zh/dashboard.ts`: included the Chinese workbench/API guidance changes.
- `frontend/src/views/user/VideoGenerationView.vue`: included model-specific parameter constraints and submission validation.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: included regression coverage for parameter filtering and payload validation.
- `progress.md`: recorded this release preparation and verification evidence.
- Rollback point before this release commit: `2050e1ed9`; restore the six listed files from that commit, or revert the release commit after it is created.

## 2026-07-23 - Task: Add Seedance 2.0 Mini model and channel pricing
### What was done
- Added `seedance-2.0-mini` to Leo account model candidates, whitelist/preset mappings, and the video generation workbench.
- Aligned Mini workbench capabilities with the latest LeoStudio matrix: `720p` only, `16:9` only, and the existing `4`–`15` second duration choices; omitted resolution now defaults to `720p` for video requests.
- Added model-aware Leo channel pricing: regular Seedance models retain `480p/720p/1080p` tiers, while Mini uses a single `720p` USD-per-second tier. Mini requests using unsupported resolutions are rejected before scheduling or forwarding.
- Updated the Leo integration documentation and Chinese/English UI help text with Mini capabilities and pricing rules.

### Testing
- From `backend/`, targeted Leo service and handler tests passed, including Mini resolution capability, account candidates, channel pricing validation, resolver extraction, and unsupported-resolution request rejection.
- From `backend/`, `go test ./internal/service ./internal/handler -count=1` passed.
- From `frontend/`, Vitest passed for channel pricing, account mappings, and video workbench coverage: 5 files, 89 tests.
- From `frontend/`, the full `npm run test:run` suite completed successfully (Vitest exit code 0).
- From `frontend/`, `npm run typecheck` passed.
- From `frontend/`, `npm run lint:check -- --no-fix` passed.
- From `frontend/`, `npm run build` passed; Vite reported only existing chunk-size/dynamic-import warnings.
- `git diff --check` passed. No paid video generation or production deployment was performed.

### Notes
- `backend/internal/service/video_billing_resolution.go`: added model-aware resolution capability and default helpers.
- `backend/internal/service/leo_video.go`: applies the 720p default when resolution is omitted.
- `backend/internal/service/video_job_service.go`: validates Mini resolution before creating an async job.
- `backend/internal/handler/leo_video.go`: rejects unsupported Mini resolution on the synchronous endpoint.
- `backend/internal/service/leo_account.go`: adds Mini to Leo default model candidates.
- `backend/internal/service/channel_service.go`: validates one 720p Mini pricing tier while preserving three tiers for other Seedance models.
- `backend/internal/service/model_pricing_resolver.go`: extracts Mini channel pricing with unsupported tiers unset.
- `backend/internal/service/video_job_billing.go`: snapshots Mini pricing without requiring unsupported tiers.
- `backend/internal/service/*_test.go` and `backend/internal/handler/leo_video_test.go`: add Mini capability, pricing, account, and request-validation regression coverage.
- `frontend/src/constants/channel.ts`: adds Mini to Leo video pricing model order.
- `frontend/src/composables/useModelWhitelist.ts`: adds Mini whitelist and account mapping preset.
- `frontend/src/components/account/CreateAccountModal.vue`: prepopulates Mini account mapping.
- `frontend/src/views/user/VideoGenerationView.vue`: exposes Mini with its strict resolution/aspect matrix.
- `frontend/src/components/admin/channel/types.ts`: makes video pricing intervals model-aware.
- `frontend/src/components/admin/channel/PricingEntryCard.vue` and `frontend/src/views/admin/ChannelsView.vue`: display, normalize, and serialize Mini's single 720p pricing tier.
- `frontend/src/i18n/locales/en/dashboard.ts`, `frontend/src/i18n/locales/zh/dashboard.ts`, `frontend/src/i18n/locales/en/admin/channels.ts`, `frontend/src/i18n/locales/zh/admin/channels.ts`: document Mini capabilities and model-specific pricing validation.
- `frontend/src/components/admin/channel/__tests__/*` and `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: cover Mini pricing and workbench behavior.
- `docs/LEO_VIDEO_CHANNEL.md`: documents Mini model parameters, mapping, and pricing.
- `progress.md`: records this implementation and verification evidence.
- Rollback point: this task is part of the existing uncommitted worktree; before reverting, save `git diff` and selectively reverse the Mini-related hunks so earlier user changes in overlapping files remain intact. No commit or deployment was created.

## 2026-07-23 - Task: Verify and push Seedance 2.0 Mini support
### What was done
- Revalidated the existing Mini model, model-aware pricing, resolution validation, and Leo video workbench changes before pushing them to `codex/leo-video-channel`.
- Applied the required Go formatting correction to the new pricing resolver regression test.
- Kept the unrelated untracked `.superpowers/` directory outside the commit; this task creates no release tag.

### Testing
- From `backend/`, targeted Mini/Leo tests passed.
- From `backend/`, `go test ./internal/service ./internal/handler -count=1` passed.
- From `frontend/`, the full `pnpm run test:run` suite passed.
- From `frontend/`, `pnpm run typecheck` passed.
- From `frontend/`, `pnpm run lint:check -- --no-fix` passed.
- From `frontend/`, `pnpm run build` passed.
- `gofmt` and `git diff --check` passed for the current changes.

### Notes
- `backend/internal/handler/leo_video.go`: retained unsupported-resolution request rejection.
- `backend/internal/handler/leo_video_test.go`: retained handler validation coverage.
- `backend/internal/service/channel_service.go`: retained model-specific video pricing validation.
- `backend/internal/service/channel_service_test.go`: retained pricing validation coverage.
- `backend/internal/service/leo_account.go`: retained the Mini account model candidate.
- `backend/internal/service/leo_account_test.go`: retained account candidate coverage.
- `backend/internal/service/leo_video.go`: retained the Mini default resolution behavior.
- `backend/internal/service/model_pricing_resolver.go`: retained Mini pricing extraction.
- `backend/internal/service/model_pricing_resolver_test.go`: retained resolver coverage and applied Go formatting.
- `backend/internal/service/video_billing_resolution.go`: retained model-aware resolution capabilities.
- `backend/internal/service/video_billing_resolution_test.go`: retained resolution capability coverage.
- `backend/internal/service/video_job_billing.go`: retained Mini billing snapshot handling.
- `backend/internal/service/video_job_service.go`: retained async request validation.
- `docs/LEO_VIDEO_CHANNEL.md`: retained Mini model and pricing documentation.
- `frontend/src/components/account/CreateAccountModal.vue`: retained Mini account mapping defaults.
- `frontend/src/components/admin/channel/PricingEntryCard.vue`: retained Mini pricing tier display.
- `frontend/src/components/admin/channel/__tests__/PricingEntryCard.spec.ts`: retained pricing display coverage.
- `frontend/src/components/admin/channel/__tests__/types.spec.ts`: retained pricing type coverage.
- `frontend/src/components/admin/channel/types.ts`: retained model-aware pricing interval normalization.
- `frontend/src/composables/useModelWhitelist.ts`: retained Mini whitelist and mapping preset.
- `frontend/src/constants/channel.ts`: retained Mini Leo model registration.
- `frontend/src/i18n/locales/en/admin/channels.ts`: retained English pricing validation text.
- `frontend/src/i18n/locales/en/dashboard.ts`: retained English Mini workbench guidance.
- `frontend/src/i18n/locales/zh/admin/channels.ts`: retained Chinese pricing validation text.
- `frontend/src/i18n/locales/zh/dashboard.ts`: retained Chinese Mini workbench guidance.
- `frontend/src/views/admin/ChannelsView.vue`: retained Mini pricing form behavior.
- `frontend/src/views/user/VideoGenerationView.vue`: retained Mini capability filtering and defaults.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: retained Mini workbench regression coverage.
- `progress.md`: recorded this verification and push scope.
- Rollback point before the push commit: `45b0e1813`; revert the new commit after it is created, or restore the listed files from that commit.

## 2026-07-23 - Task: Align current Seedance video parameters with LeoStudio

### What was done
- Synchronized the Sub2 request constraints for `seedance-2.0`, `seedance-2.0-fast`, and `seedance-2.0-mini` with LeoStudio commit `2fd5c21b01a049817962812cf4675ade7727cc12`; Happy Horse and Grok were intentionally not added.
- Added shared synchronous and asynchronous validation for model-specific resolutions, resolution-specific aspect ratios, whole-second durations from 4 through 15, the 8-second default, the 5000-character prompt limit, and Seedance guidance counts.
- Kept Mini fixed to `720p` and `16:9`, kept Fast limited to `480p` and `720p`, and validated requests again after account model mapping so aliases cannot bypass the effective upstream model constraints.
- Updated the workbench prompt limit and the user-facing API documentation without changing the existing channel pricing structure.

### Testing
- From `backend/`, targeted Seedance validation, handler, async handler, billing-resolution, and video-job tests passed for both `internal/service` and `internal/handler`.
- From `backend/`, `go test ./internal/service ./internal/handler` passed.
- From `frontend/`, the two video view suites passed: 2 files and 23 tests.
- From `frontend/`, `vue-tsc --noEmit` and targeted ESLint checks passed.
- From `frontend/`, the production Vue/Vite build passed; only existing Browserslist, dynamic-import, and chunk-size warnings were reported.
- `gofmt` and `git diff --check` passed. No paid video generation or production deployment was performed.

### Notes
- `.gitignore`: allowed the new Leo video specification document to be tracked while retaining the existing `docs/*` policy.
- `backend/internal/handler/leo_video.go`: applies shared validation to synchronous requests and returns mapped-model validation failures as HTTP 400.
- `backend/internal/handler/leo_video_async.go`: classifies local specification failures as HTTP 400 for asynchronous requests.
- `backend/internal/handler/leo_video_test.go`: covers Fast resolution, Mini resolution/aspect, and duration rejection at the synchronous entry point.
- `backend/internal/handler/leo_video_async_test.go`: verifies the asynchronous endpoint rejects a non-16:9 Mini request.
- `backend/internal/service/leo_video.go`: validates the request against the account-mapped synchronous upstream model before forwarding.
- `backend/internal/service/leo_video_async.go`: validates the request against the account-mapped asynchronous upstream model before forwarding.
- `backend/internal/service/leo_video_model_specs.go`: defines the three Seedance capability matrices and shared request validation.
- `backend/internal/service/leo_video_model_specs_test.go`: covers defaults, valid combinations, rejected combinations, prompt/guidance limits, and mapped-model validation.
- `backend/internal/service/video_billing_resolution.go`: uses strict model capabilities instead of treating unknown resolutions as 480p during request validation.
- `backend/internal/service/video_billing_resolution_test.go`: covers Fast, Mini, and unknown-resolution behavior.
- `backend/internal/service/video_job_service.go`: validates async jobs before billing and persists normalized duration, resolution, and aspect values.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: records the supported Sub2 Seedance matrix and LeoStudio alignment point.
- `frontend/src/i18n/locales/en/dashboard.ts`: documents the English 5000-character prompt limit.
- `frontend/src/i18n/locales/zh/dashboard.ts`: documents the Chinese 5000-character prompt limit.
- `frontend/src/views/user/VideoGenerationView.vue`: applies the model prompt limit while retaining the synchronized resolution, duration, and aspect options.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: verifies the prompt limit in addition to the existing Mini `720p` and `16:9` UI constraint.
- `progress.md`: records this implementation and its verification evidence.
- Rollback point: `c0e60c7b3bd90f1289f49f4f46f87a57b0849ed8`. Before committing, restore the tracked files listed above from that revision and remove `backend/internal/service/leo_video_model_specs.go`, `backend/internal/service/leo_video_model_specs_test.go`, and `docs/LEO_VIDEO_MODEL_SPECS.md`; remove only this final progress entry.

## 2026-07-24 - Task: Limit Seedance 2.0 1080p duration to 12 seconds

### What was done
- Limited standard `seedance-2.0` requests at `1080p` to whole-second durations from 4 through 12 while keeping standard 480p/720p, Fast, and Mini at 4 through 15 seconds.
- Synchronized the video workbench so the 1080p duration selector ends at 12 seconds and an existing 13-15 second selection falls back to the valid 8-second default when switching resolutions.
- Enforced the same rule in backend request validation, including account-mapped model aliases, so direct API calls cannot bypass the workbench constraint.
- Updated the formal model specification and the Chinese/English API documentation shown on the video documentation subpage.

### Testing
- From `backend/`, `go test ./internal/service ./internal/handler -count=1` passed; coverage includes 1080p duration 12 acceptance, duration 13 rejection, mapped-model rejection, and continued 15-second acceptance for standard 720p and Fast 720p.
- From `frontend/`, the video generation and API documentation suites passed: 2 files and 23 tests. Coverage includes 1080p options ending at 12, automatic fallback after a resolution change, and prevention of programmatically injected 1080p/13-second submissions before upload or job creation.
- From `frontend/`, `vue-tsc --noEmit`, `vue-tsc -b`, and targeted ESLint checks passed.
- The production Vite build passed. Only existing Browserslist, dynamic-import, and chunk-size warnings were reported.
- `gofmt` and `git diff --check` passed. No paid video generation, commit, push, tag, or deployment was performed.

### Notes
- `backend/internal/service/leo_video_model_specs.go`: adds the standard model's 1080p-specific maximum and enforces it after model and resolution normalization.
- `backend/internal/service/leo_video_model_specs_test.go`: covers the 12/13-second boundary, unaffected 15-second modes, and mapped-model validation.
- `backend/internal/handler/leo_video_test.go`: verifies the HTTP entry point rejects a standard 1080p 13-second request with HTTP 400.
- `frontend/src/views/user/VideoGenerationView.vue`: filters duration options by resolution, validates the selected combination, and resets an invalid duration after resolution changes.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: covers duration option filtering, fallback behavior, and client-side rejection before network calls.
- `frontend/src/i18n/locales/en/dashboard.ts`: updates the English API duration description.
- `frontend/src/i18n/locales/zh/dashboard.ts`: updates the Chinese API duration description.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: records the 4-12 second range for standard 1080p.
- `progress.md`: records this implementation, verification evidence, file list, and rollback point.
- Rollback point: `52468162b71774e2066874e982b24948f12520e6`. Run `git restore --source=52468162b71774e2066874e982b24948f12520e6 -- backend/internal/handler/leo_video_test.go backend/internal/service/leo_video_model_specs.go backend/internal/service/leo_video_model_specs_test.go docs/LEO_VIDEO_MODEL_SPECS.md frontend/src/i18n/locales/en/dashboard.ts frontend/src/i18n/locales/zh/dashboard.ts frontend/src/views/user/VideoGenerationView.vue frontend/src/views/user/__tests__/VideoGenerationView.spec.ts progress.md` to revert this task; the unrelated `.superpowers/` directory remains untouched.

## 2026-07-24 - Task: Fix frontend security scan PostCSS advisory
### What was done
- Upgraded the frontend `postcss` dependency from `8.5.6` to the patched `8.5.22` release and synchronized the pnpm lockfile.
- Kept the existing `xlsx` exceptions unchanged because they remain scoped, documented, and unexpired; no new audit suppression was added.
- This is a security-fix commit only; no new Release tag was created in this task.

### Testing
- `pnpm audit --prod --audit-level=high --json` still reports the documented `xlsx` advisories, and `tools/check_pnpm_audit_exceptions.py` passed with exit code 0.
- From `frontend/`, frozen installation with pnpm 9 passed.
- From `frontend/`, the full test suite passed.
- From `frontend/`, typecheck, lint, and production build passed.
- `git diff --check` passed. No production deployment was performed.

### Notes
- `frontend/package.json`: raises the direct PostCSS devDependency to the patched range.
- `frontend/pnpm-lock.yaml`: locks PostCSS 8.5.22 and its required transitive `nanoid` update.
- `progress.md`: records the security scan diagnosis, remediation, and verification evidence.
- Rollback point before this security-fix commit: `38adf812d`; restore the three listed files from that revision or revert the new commit after it is created.

## 2026-07-24 - Task: Release Seedance 2.0 1080p duration cap as v0.1.163-fy.5
### What was done
- Revalidated and prepared the current 1080p duration-limit changes for `v0.1.163-fy.5`.
- Kept the release scoped to the current branch and excluded the unrelated untracked `.superpowers/` directory.

### Testing
- From `backend/`, `go test ./internal/service ./internal/handler -count=1` passed.
- From `backend/`, `gofmt` passed for the changed Go files.
- From `frontend/`, the video generation and API test suites passed.
- From `frontend/`, `pnpm run typecheck` passed.
- From `frontend/`, targeted ESLint passed for the changed Vue and locale files.
- From `frontend/`, `pnpm run build` passed.
- `git diff --check` passed. No paid video generation or production deployment was performed.

### Notes
- `backend/internal/handler/leo_video_test.go`: records HTTP rejection coverage for standard 1080p duration 13.
- `backend/internal/service/leo_video_model_specs.go`: enforces the standard 1080p 4-12 second limit.
- `backend/internal/service/leo_video_model_specs_test.go`: covers the 1080p boundary, unaffected modes, and mapped-model validation.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: documents the resolution-specific duration range.
- `frontend/src/i18n/locales/en/dashboard.ts`: updates English duration guidance.
- `frontend/src/i18n/locales/zh/dashboard.ts`: updates Chinese duration guidance.
- `frontend/src/views/user/VideoGenerationView.vue`: filters duration options and resets invalid values after resolution changes.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: covers the 1080p duration selector and client-side guard.
- `progress.md`: records this release preparation and verification evidence.
- Rollback point before the release commit: `52468162b71774e2066874e982b24948f12520e6`; revert the release commit after it is created, or restore the nine listed files from that revision.

## 2026-07-24 - Task: Merge upstream v0.1.164 and prepare fork release

### What was done
- Fetched `Wei-Shaw/sub2api` `upstream/main` at `cb24522dd` (`v0.1.164`) and merged it into the fork branch with a no-fast-forward merge.
- Preserved the local Leo video platform and token-incentive changes while integrating upstream Composite groups, Ollama Cloud usage, payment, routing, and migration changes.
- Resolved platform, generated Ent/Wire, gateway route, admin group, frontend type, icon, color, locale, and group view conflicts by retaining both platform implementations.
- Updated the scheduler composite-platform test expectation to include the existing local Leo scheduler bucket.

### Testing
- From `backend/`, `go test ./...` passed with exit code 0.
- From `frontend/`, `pnpm exec vitest run`, `pnpm typecheck`, and `pnpm run build` passed with exit code 0.
- `gofmt` and `git diff --cached --check` passed; no unresolved merge files or conflict markers remain.

### Notes
- `backend/cmd/server/wire_gen.go`, `backend/ent/client.go`, `backend/internal/domain/constants.go`, `backend/internal/handler/admin/group_handler.go`, `backend/internal/server/routes/gateway.go`, `backend/internal/server/routes/gateway_test.go`, `backend/internal/service/admin_group.go`, `backend/internal/service/domain_constants.go`, and `backend/internal/service/composite_platform_test.go`: resolve upstream Composite and local Leo integration points.
- `frontend/src/components/common/GroupBadge.vue`, `frontend/src/components/common/PlatformIcon.vue`, `frontend/src/i18n/locales/en/admin/overview.ts`, `frontend/src/i18n/locales/zh/admin/overview.ts`, `frontend/src/types/index.ts`, `frontend/src/utils/platformColors.ts`, and `frontend/src/views/admin/GroupsView.vue`: retain both platform options in the frontend.
- All other changed files are the upstream `main` synchronization set; the complete file list is recorded by the merge commit.
- `.superpowers/` remains untracked and is excluded from the merge.
- Rollback point before this merge: `79c37efaac634cbb73db825d3f4b72ac7e505927`; before commit use `git merge --abort`, or after the merge commit use `git revert -m 1 <merge_commit>`.

## 2026-07-24 - Task: Disable the user responsible for an OpenAI cyber-policy account revocation

### What was done
- Added a post-persistence rule for OpenAI Pro/ProLite OAuth accounts explicitly marked `Token revoked (401)` that correlates the account with upstream `cyber_policy` errors from the preceding 30 days.
- Selects the non-admin user with the most hits, breaking ties by the latest hit and then the lowest user ID, and disables that user while invalidating API-key authentication caches.
- Excludes generic 401 responses, OpenAI Plus/Team/Free accounts, non-OpenAI accounts, local `cyber_policy_session_blocked` rejections, deleted users, and administrators; repeated observations do not update an already disabled user.

### Testing
- From `backend/`, the targeted service and repository tests passed, covering revoked-token enforcement, generic 401 exclusion, administrator exclusion, idempotency, the 30-day window, deterministic ranking, and exact `cyber_policy` filtering.
- From `backend/`, `go test ./internal/service ./internal/repository -count=1` passed.
- From `backend/`, `go test ./... -count=1` passed.
- The additional full `-tags unit` run reached the 184-second command timeout without reporting a test failure; the corresponding tagged compile check, `go test -tags unit ./internal/service ./internal/repository -run '^$' -count=1`, passed.
- `gofmt` and `git diff --check` passed. No production database write, user ban, deployment, commit, or push was performed.

### Notes
- `.gitignore`: allows the focused cyber-policy revocation behavior document to be tracked without changing the default `docs/` ignore policy.
- `backend/internal/repository/ops_repo_cyber_policy.go`: queries the deterministic non-admin candidate for one account and time window.
- `backend/internal/repository/ops_repo_cyber_policy_test.go`: verifies exact upstream cyber-policy filtering and count/latest-hit ordering.
- `backend/internal/service/ops_cyber_policy_ban.go`: validates the persisted OpenAI revoked-token account state and applies the user disable action.
- `backend/internal/service/ops_cyber_policy_ban_test.go`: covers triggering, exclusions, administrator protection, and repeated-observation behavior.
- `backend/internal/service/account.go`: provides the normalized Pro/ProLite `plan_type` eligibility check.
- `backend/internal/service/ops_port.go`: exposes the narrow candidate query contract and result model.
- `backend/internal/service/ops_repo_mock_test.go`: adds the candidate query hook required by Ops service tests.
- `backend/internal/service/ops_service.go`: runs the rule only after successful single or batched Ops error persistence.
- `docs/CYBER_POLICY_REVOCATION_BAN.md`: documents the trigger, ranking, exclusions, time window, and Ops-monitoring dependency.
- `progress.md`: records this implementation, verification evidence, file list, and rollback point.
- Rollback point: `122aeb81dcd10e2411a65aa9878d052789066fbf`. Run `git restore --source=122aeb81dcd10e2411a65aa9878d052789066fbf -- .gitignore backend/internal/service/ops_port.go backend/internal/service/ops_repo_mock_test.go backend/internal/service/ops_service.go progress.md` and remove the four new Go files plus `docs/CYBER_POLICY_REVOCATION_BAN.md` to revert this task; leave the unrelated `.superpowers/` directory untouched.

## 2026-07-24 - Task: Adjust production Seedance video pricing and publish notice

### What was done
- Updated production channel `Seedance 2 视频专用渠道` (ID `5`) so the user-facing Seedance prices are `0.12/0.25/0.60` USD/s for standard 480p/720p/1080p, `0.10/0.20` USD/s for Fast 480p/720p, and `0.17` USD/s for Mini 720p.
- Preserved the existing Fast 1080p entry at `0.25` USD/s because the deployed channel validator requires three tiers for non-Mini entries; Fast 1080p remains unavailable in the model capability matrix and was not included in the user announcement.
- Published active all-user popup announcement ID `19`, explaining that the increase is caused by higher account-pool maintenance costs and that the new prices take effect immediately.

### Testing
- The production Admin API update returned channel ID `5` with all six requested prices and the preserved Fast 1080p compatibility tier.
- A separate `GET /api/v1/admin/channels/5` readback returned standard `0.12/0.25/0.60`, Fast `0.10/0.20/0.25`, and Mini `0.17` with channel status `active`.
- `GET /api/v1/admin/announcements/19` returned the expected title and price table with status `active`, notify mode `popup`, and empty targeting for all users.

### Notes
- `docs/LEO_VIDEO_CHANNEL.md`: records the production price snapshot and explains the non-user-facing Fast 1080p compatibility tier.
- `progress.md`: records the production configuration change, verification evidence, and rollback values.
- Production state changed through the Admin API only; no database schema, source behavior, deployment, commit, or push was performed. Temporary request payload files were removed after verification, and the administrator credential was not written to disk or logs.
- Rollback point: update channel ID `5` through `PUT /api/v1/admin/channels/5`, restoring standard prices to `0.10/0.20/0.55`, Fast to `0.08/0.16/0.25`, and Mini 720p to `0.14`; then archive announcement ID `19` with `PUT /api/v1/admin/announcements/19` and `{"status":"archived"}`. Keep the existing channel association, model mapping, and group ID `25` unchanged.

## 2026-07-24 - Task: Verify cyber-policy revocation ban with simulated trigger only

### What was done
- Re-ran the Pro/ProLite cyber-policy revocation scenarios using in-memory fake repositories and SQL mocks only.
- Confirmed the simulated top-hit user is selected, while Plus accounts, generic OpenAI 401 responses, administrators, and repeated observations are excluded or handled idempotently.

### Testing
- From `backend/`, the verbose targeted service and repository test run passed.
- No production database connection, real user status update, API call, deployment, commit, or push was performed.

### Notes
- `backend/internal/service/ops_cyber_policy_ban_test.go`: simulated trigger and exclusion coverage.
- `backend/internal/repository/ops_repo_cyber_policy_test.go`: simulated ranking and exact error-type filtering coverage.
- `progress.md`: records this simulation-only verification.
- Rollback point: remove this final progress entry; source rollback remains the prior cyber-policy task rollback point `122aeb81dcd10e2411a65aa9878d052789066fbf`.

## 2026-07-24 - Task: Persist cyber-policy revocation audit records and expose them in security audit UI

### What was done
- Added a synchronous append-only audit event after a Pro/ProLite cyber-policy revocation rule successfully disables the attributed user.
- Retained the latest matching request summary, request IDs, account and plan identifiers, hit count, model/path, client metadata, revocation request IDs, outcome, and the existing redacted input excerpt when available.
- Added readable rule-event details to the admin audit log page and moved its navigation entry under Security Audit.

### Testing
- From `backend/`, `go test ./... -count=1` passed.
- From `frontend/`, the full Vitest suite passed, `vue-tsc --noEmit` passed, and the production Vite build passed.
- `gofmt` and `git diff --check` passed. No production database connection, real user ban, deployment, commit, or push was performed.

### Notes
- `backend/internal/service/audit_log.go`: names the cyber-policy revocation audit action.
- `backend/internal/service/audit_log_service.go`: adds synchronous audit persistence for enforcement events.
- `backend/internal/service/ops_cyber_policy_ban.go`: writes the redacted rule audit event after a successful in-memory/production user update path.
- `backend/internal/service/ops_port.go`: carries the latest attributed request metadata.
- `backend/internal/repository/ops_repo_cyber_policy.go`: returns the latest matching Ops request and retained moderation excerpt.
- `backend/internal/service/ops_service.go`, `backend/internal/service/wire.go`, `backend/cmd/server/wire_gen.go`: inject the audit service into Ops.
- `backend/internal/repository/ops_repo_cyber_policy_test.go`, `backend/internal/service/ops_cyber_policy_ban_test.go`: cover request metadata and audit persistence with SQL mocks/fake repositories.
- `frontend/src/views/admin/AuditLogView.vue`, `frontend/src/components/layout/AppSidebar.vue`, `frontend/src/i18n/locales/en/admin/audit.ts`, `frontend/src/i18n/locales/zh/admin/audit.ts`, `frontend/src/features/prompt-audit/__tests__/integrationSurface.spec.ts`: expose and label the rule event in Security Audit.
- `docs/CYBER_POLICY_REVOCATION_BAN.md`: documents audit retention and redaction behavior.
- `progress.md`: records the implementation, verification evidence, file list, and rollback guidance.
- Rollback point: no commit was created; remove the audit-specific hunks in the files listed above while preserving the earlier cyber-policy rule changes in shared files such as `ops_service.go`, `ops_port.go`, and `progress.md`.

## 2026-07-24 - Task: Prepare cyber-policy revocation enforcement release

### What was done
- Prepared the cyber-policy revocation attribution, user disable action, synchronous audit record, and Security Audit UI changes for release as `v0.1.164-fy.2`.
- Kept the unrelated `.superpowers/` workspace artifacts outside the release commit.

### Testing
- From `backend/`, `go test ./internal/service -run '^TestOllamaCloudUsageRefreshSingleflightAndRunnerDeduplicateSharedGroup$' -count=10` passed after an earlier unrelated timing failure.
- From `backend/`, `go test ./... -count=1` passed.
- From `frontend/`, `pnpm.cmd test:run`, `pnpm.cmd typecheck`, `pnpm.cmd lint:check`, and `pnpm.cmd run build` passed.
- `git diff --check` passed.

### Notes
- The release file set is the cyber-policy enforcement, audit persistence, Security Audit UI, supporting tests and documentation already listed in the preceding task entries, plus this `progress.md` release record.
- `.superpowers/` remains untracked and is not part of the release.
- Rollback point before this release: `122aeb81dcd10e2411a65aa9878d052789066fbf`; after publication, revert the release commit or reinstall `v0.1.164-fy.1`.

## 2026-07-24 - Task: Exempt authenticated administrators from content auditing
### What was done
- Restored production administrator user `xl` (ID `1`) through the risk-control unban operation and verified the account role remains `admin` with status `active`.
- Added one trusted-role bypass at the unified gateway security-audit entry so an authenticated administrator skips local keyword checks, legacy moderation, prompt audit, violation counting, hit notifications, and moderation-triggered automatic disabling.
- Kept the bypass limited to requests where the authenticated user ID, API key owner ID, and loaded user record ID match; ordinary users continue through the existing moderation path.

### Testing
- Production request `ec622e1f-cc79-429e-aecf-d6e53c5edf60` was confirmed as a `keyword_block` for administrator user ID `1`, while its moderation record showed `auto_banned=false`; the follow-up unban and user readback both returned status `active`.
- From `backend/`, focused administrator-bypass and regular-user regression tests passed.
- From `backend/`, the complete `go test ./... -count=1` suite passed.
- `gofmt` and `git diff --check` passed for the changed Go files and repository diff.

### Notes
- `.gitignore`: allows the content-moderation behavior document to be tracked.
- `backend/internal/handler/security_audit_helper.go`: bypasses the unified audit coordinator only for a matching authenticated administrator identity.
- `backend/internal/handler/security_audit_helper_test.go`: proves both audit engines remain untouched for administrators and that regular users are still blocked.
- `docs/CONTENT_MODERATION.md`: documents the administrator bypass and its trusted identity boundary.
- `progress.md`: records the production recovery, implementation, verification evidence, file list, and rollback point.
- Rollback point before this task: `a9351087d`. Revert the task commit or restore the four source/documentation files from that revision and remove `docs/CONTENT_MODERATION.md`; leave the unrelated `.superpowers/` directory untouched. Production code rollback is `v0.1.164-fy.2`.

## 2026-07-24 - Task: Deploy the administrator content-audit bypass to production
### What was done
- Published release `v0.1.164-fy.3` from commit `9e25f5308` and installed it on `api.fflink.top` through the built-in checksum-verifying updater.
- Restarted the production service after installation and kept the existing database schema and risk-control configuration unchanged.
- Preserved the user's concurrent, unstaged gateway and moderation edits outside both the release commit and tag.

### Testing
- The exact release commit passed `go test ./... -count=1` from an isolated clean Git worktree; the main working tree's concurrent edits were not part of this verification or release.
- GitHub published both `checksums.txt` and the 35,323,253-byte Linux amd64 archive; the production updater reported a completed update, which requires its internal SHA256 verification to pass before binary replacement.
- After restart, public `GET /health` returned HTTP `200` and the admin version endpoint returned `0.1.164-fy.3`.
- Production user ID `1` read back as role `admin` and status `active`. Content moderation remained enabled in `pre_block` mode with automatic banning enabled for configured non-admin traffic.

### Notes
- `progress.md`: records release publication, deployment, production verification, and rollback instructions.
- No migration, database write, moderation configuration update, user-role change, or ordinary-user audit exemption was included in the deployment.
- Production rollback: call `POST /api/v1/admin/system/rollback` with `{"version":"0.1.164-fy.2"}`, then call `POST /api/v1/admin/system/restart`; source rollback is tag `v0.1.164-fy.2`.

## 2026-07-24 - Task: Retain cyber-policy request excerpts

### What was done
- Added current-turn input retention for new upstream `cyber_policy` hits on OpenAI Responses, Chat Completions, Anthropic-compatible Messages, and Responses WebSocket requests.
- Reused the existing moderation log `input_excerpt` field so the retained text is visible in Risk Control and available to the existing cyber-policy revocation audit record without a database migration or frontend change.
- Limited retained text to 4,000 Unicode characters, applied the existing credential redaction, and omitted image/base64 content by storing extracted text instead of the raw request body. Historical empty excerpts remain unchanged.

### Testing
- From `backend/`, `go test ./internal/service -run "TestRecordCyberPolicyEvent_WritesLogWhenEnabled" -count=1 -v` passed, covering persistence, secret redaction, and the 4,000-character limit.
- From `backend/`, `go test ./internal/handler -run "Test(CyberPolicyInputExcerptIfMarked|RecordCyberPolicyIfMarked|ClearCyberPolicyTurnState)" -count=1 -v` passed, covering mark-gated extraction, current input selection, image omission, idempotency, and WebSocket turn reset behavior.
- From `backend/`, `go test ./... -count=1` passed.
- `gofmt` was applied to the changed Go files and `git diff --check` passed before this log append.

### Notes
- `backend/internal/handler/openai_chat_completions.go`: passes the marked Chat Completions input excerpt into cyber-policy recording.
- `backend/internal/handler/openai_gateway_handler.go`: passes Responses and Messages excerpts and tracks the current WebSocket turn excerpt without retaining the raw request body.
- `backend/internal/handler/openai_gateway_cyber_test.go`: verifies extraction only after a cyber mark and excludes image payloads.
- `backend/internal/service/content_moderation.go`: persists the redacted, length-limited cyber-policy input excerpt.
- `backend/internal/service/content_moderation_cyber_test.go`: verifies excerpt persistence, credential redaction, and length limiting.
- `docs/CYBER_POLICY_REVOCATION_BAN.md`: documents new-hit retention, visibility, redaction, limits, and the historical-data boundary.
- `progress.md`: records the implementation, verification evidence, changed files, and rollback point.
- Rollback point: no commit was created. Run `git restore --source=HEAD -- backend/internal/handler/openai_chat_completions.go backend/internal/handler/openai_gateway_cyber_test.go backend/internal/handler/openai_gateway_handler.go backend/internal/service/content_moderation.go backend/internal/service/content_moderation_cyber_test.go docs/CYBER_POLICY_REVOCATION_BAN.md progress.md` to remove only this uncommitted task; leave `.superpowers/` untouched.

## 2026-07-24 - Task: Create customer-facing Seedance API guide

### What was done
- Created a standalone Chinese Seedance API guide that can be sent directly to customers.
- Documented Bearer authentication, asynchronous task creation, image upload, task polling, cancellation, MP4 download, supported model parameters, Python usage, and retry guidance.
- Removed internal provider, account, and administrator configuration details while retaining customer-relevant troubleshooting for invalid aspect ratios and interrupted chunked responses.

### Testing
- Read the completed file explicitly as UTF-8 and confirmed the Chinese content is intact.
- Verified the guide contains the production Base URL, `Prefer: respond-async`, the `Auto` restriction, and `incomplete chunked read` troubleshooting.
- Verified Markdown fenced code blocks are balanced and all fenced JSON examples parse successfully.
- `git diff --check` passed after the documentation and progress-log changes.

### Notes
- `docs/SEEDANCE_API_CLIENT_GUIDE_CN.md`: adds the standalone customer-facing Seedance API guide.
- `progress.md`: records the document scope, verification evidence, changed files, and rollback instructions.
- Rollback point: delete `docs/SEEDANCE_API_CLIENT_GUIDE_CN.md` and remove this appended task block from `progress.md`; preserve all earlier unrelated working-tree changes.

## 2026-07-24 - Task: Align customer Seedance API guide with current runtime behavior

### What was done
- Updated the customer guide with the current request-field constraints, upload lifecycle, completed-job and list response examples, and list limit behavior.
- Corrected the distinction between local HTTP 400 validation failures and upstream HTTP 422 rejections, and documented oversized upload and incomplete request-body errors.
- Clarified asynchronous status handling, authenticated output download, lack of idempotency keys, and balance reservation, settlement, and release behavior.

### Testing
- Compared the documented routes, statuses, request validation, upload behavior, list limits, output handling, cancellation rules, and billing lifecycle against the current backend implementation.
- Verified all fenced JSON examples parse successfully and all Markdown fenced code blocks are balanced.
- Verified the documented model matrix matches the current Seedance model specifications and `git diff --check` passes.

### Notes
- `docs/SEEDANCE_API_CLIENT_GUIDE_CN.md`: aligns the customer-facing contract and examples with the current REST API implementation.
- `progress.md`: records the documentation correction, verification evidence, changed files, and rollback instructions.
- Rollback point: restore the preceding version of `docs/SEEDANCE_API_CLIENT_GUIDE_CN.md` and remove this appended task block from `progress.md`; preserve all unrelated working-tree changes.

## 2026-07-24 - Task: Release cyber-policy request excerpts as v0.1.164-fy.4

### What was done
- Published commit `179b71842` as `v0.1.164-fy.4`, excluding the unrelated local Seedance guide and `.superpowers/` artifacts.
- Installed the checksum-verified Linux release on `api.fflink.top` through the built-in updater and restarted the service.

### Testing
- Before release, focused service and handler tests, the full `go test ./... -count=1` suite, targeted race detection, `gofmt -d`, and `git diff --check` passed.
- The release exposed both `checksums.txt` and the 35,323,880-byte Linux amd64 archive, and the production updater completed its checksum-verifying replacement path.
- After restart, public `GET /health` returned HTTP `200` and the admin version endpoint returned `0.1.164-fy.4`.

### Notes
- `progress.md`: records the isolated release, production deployment, verification evidence, and rollback point.
- No database migration or risk-control configuration was changed.
- Production rollback: call `POST /api/v1/admin/system/rollback` with `{"version":"0.1.164-fy.3"}`, then call `POST /api/v1/admin/system/restart`; source rollback is tag `v0.1.164-fy.3`.

## 2026-07-26 - Task: Align current Leo video integration with LeoStudio feat/web-admin

### What was done
- Compared the latest remote LeoStudio `feat/web-admin` commit `3ed1f43438325e56635f4435ff23b4c91c4b2db9` with the current Sub2 contract without modifying the local LeoStudio checkout.
- Synchronized exposed `seedance-2.0-mini` capabilities to LeoStudio's current Standard/HD matrix: `480p` and `720p`, each with `16:9`, `1:1`, and `9:16`.
- Updated frontend workbench selectors, backend validation, channel pricing forms, model pricing resolution, and billing snapshots for Mini's two supported tiers.
- Kept the user-requested standard `seedance-2.0` 1080p maximum of 12 seconds and kept Happy Horse/Grok out of Sub2 because their addition was previously deferred.
- Preserved compatibility with existing Mini 720p-only channel pricing: 720p remains billable through the old entry, while Mini 480p is rejected until its 480p price is configured; new pricing entries require both supported tiers.
- Refreshed Leo model specification, channel operations, and Chinese/English API parameter documentation to the upstream capability source.

### Testing
- From `backend/`, `go test ./internal/service ./internal/handler -count=1` ran with the service package passing; after updating the two stale Mini aspect assertions, `go test ./internal/handler -count=1` passed. The handler failure in the first combined run was the expected old `9:16` rejection assertion, not a runtime failure.
- From `frontend/`, the workbench, API docs, channel pricing types, and pricing card suites passed: 4 files and 40 tests.
- From `frontend/`, `vue-tsc --noEmit` and targeted ESLint checks passed.
- From `frontend/`, the production Vite build passed with the repository's existing Browserslist, dynamic-import, and chunk-size warnings.
- `gofmt` and `git diff --check` passed. No paid video generation, LeoStudio checkout mutation, commit, push, tag, or deployment was performed.

### Notes
- `backend/internal/service/leo_video_model_specs.go`: expands Mini validation to 480p/720p and three supported aspect ratios.
- `backend/internal/service/leo_video_model_specs_test.go`: covers Mini Standard/HD acceptance and unsupported aspect rejection.
- `backend/internal/service/video_billing_resolution.go`: exposes Mini's 480p and 720p pricing tiers.
- `backend/internal/service/video_billing_resolution_test.go`: verifies the updated Mini capability and pricing tier list.
- `backend/internal/service/model_pricing_resolver.go`: resolves model-specific tier pointers and reads legacy Mini 720p-only entries.
- `backend/internal/service/model_pricing_resolver_test.go`: covers new two-tier pricing and legacy 720p-only compatibility.
- `backend/internal/service/video_job_billing.go`: requires the price for the requested resolution, allowing legacy Mini 720p while blocking unpriced 480p.
- `backend/internal/service/channel_service_test.go`: validates the new Mini 480p+720p channel pricing shape.
- `backend/internal/handler/leo_video_test.go`: updates the stale unsupported Mini aspect regression to 21:9.
- `backend/internal/handler/leo_video_async_test.go`: updates the asynchronous stale unsupported Mini aspect regression to 21:9.
- `frontend/src/views/user/VideoGenerationView.vue`: exposes Mini 480p/720p and synchronized aspect options.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: verifies the Mini workbench matrix.
- `frontend/src/components/admin/channel/types.ts`: generates Mini 480p/720p pricing intervals.
- `frontend/src/components/admin/channel/__tests__/types.spec.ts`: covers Mini two-tier pricing validation.
- `frontend/src/components/admin/channel/__tests__/PricingEntryCard.spec.ts`: covers Mini two-price rendering.
- `frontend/src/i18n/locales/en/dashboard.ts`: updates English Mini resolution and aspect descriptions.
- `frontend/src/i18n/locales/zh/dashboard.ts`: updates Chinese Mini resolution and aspect descriptions.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: records the latest upstream source and Mini matrix.
- `docs/LEO_VIDEO_CHANNEL.md`: updates Mini pricing, model matrix, upstream commit, and legacy pricing behavior.
- `progress.md`: appends this implementation and verification record without rewriting prior entries.
- Rollback point: `HEAD` before this task. Run `git restore --source=HEAD -- backend/internal/handler/leo_video_async_test.go backend/internal/handler/leo_video_test.go backend/internal/service/channel_service_test.go backend/internal/service/leo_video_model_specs.go backend/internal/service/leo_video_model_specs_test.go backend/internal/service/model_pricing_resolver.go backend/internal/service/model_pricing_resolver_test.go backend/internal/service/video_billing_resolution.go backend/internal/service/video_billing_resolution_test.go backend/internal/service/video_job_billing.go docs/LEO_VIDEO_CHANNEL.md docs/LEO_VIDEO_MODEL_SPECS.md frontend/src/components/admin/channel/__tests__/PricingEntryCard.spec.ts frontend/src/components/admin/channel/__tests__/types.spec.ts frontend/src/components/admin/channel/types.ts frontend/src/i18n/locales/en/dashboard.ts frontend/src/i18n/locales/zh/dashboard.ts frontend/src/views/user/VideoGenerationView.vue frontend/src/views/user/__tests__/VideoGenerationView.spec.ts` to revert only the code/document changes; remove only this appended block from `progress.md`, preserving the pre-existing uncommitted progress entries and `.superpowers/`.

## 2026-07-26 - Task: Synchronize LeoStudio media and audio references

### What was done
- Rechecked LeoStudio `feat/web-admin` at remote commit `3ed1f43438325e56635f4435ff23b4c91c4b2db9` and corrected the prior scope: Mini remains `720p` and `16:9`; Mini `480p` was not retained.
- Added backend validation and passthrough coverage for `guidances.video_reference_base[].video` and `guidances.audio_reference[].audio`, including UUID/absolute HTTP(S) URL formats, `UPLOADED` URL types, video duration rejection, audio UUID duration limits, URL duration omission, and the required visual reference for audio.
- Included video/audio reference URLs in local input-token lifecycle tracking without treating them as image moderation inputs.
- Added video/audio request examples and synchronized Chinese/English API documentation with the latest media-reference contract.

### Testing
- Targeted backend media/audio, request parsing, and input-store tests passed.
- `go test ./internal/service ./internal/handler -count=1` passed.
- Frontend API-docs, video-generation, and video-generation workbench tests passed (27 tests); `npm run typecheck` and targeted ESLint passed.
- `gofmt` and `git diff --check` passed before the final build rerun.

### Notes
- `backend/internal/service/leo_video.go`: separates image URLs from video/audio reference URLs and keeps all reference tokens tracked.
- `backend/internal/service/leo_video_model_specs.go`: validates the latest LeoStudio video/audio guidance contract.
- `backend/internal/service/leo_video_model_specs_test.go`: covers accepted and rejected media/audio guidance inputs.
- `backend/internal/service/leo_video_test.go`: verifies JSON passthrough and reference URL parsing.
- `backend/internal/service/video_input_store_test.go`: verifies media/audio local input token retention.
- `frontend/src/views/user/VideoApiDocsView.vue`: adds video and audio reference request examples.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: verifies the new examples render.
- `frontend/src/i18n/locales/en/dashboard.ts`: documents media/audio guidance rules in English.
- `frontend/src/i18n/locales/zh/dashboard.ts`: documents media/audio guidance rules in Chinese.
- `docs/LEO_VIDEO_CHANNEL.md`: records the upstream contract and URL/UUID lifecycle.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: records media/audio field and duration rules while keeping Mini at 720p/16:9.
- `progress.md`: records this synchronization and its verification evidence.
- Rollback point: no commit was created. Run `git restore --source=HEAD -- backend/internal/service/leo_video.go backend/internal/service/leo_video_model_specs.go backend/internal/service/leo_video_model_specs_test.go backend/internal/service/leo_video_test.go backend/internal/service/video_input_store_test.go docs/LEO_VIDEO_CHANNEL.md docs/LEO_VIDEO_MODEL_SPECS.md frontend/src/i18n/locales/en/dashboard.ts frontend/src/i18n/locales/zh/dashboard.ts frontend/src/views/user/VideoApiDocsView.vue frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts` and remove only this appended block from `progress.md`; preserve unrelated `.superpowers/` and earlier progress entries.

## 2026-07-26 - Task: Add customer-facing video and audio reference uploads

### What was done
- Extended `POST /v1/videos/uploads` to accept legacy image uploads plus dedicated video/audio fields and `file + media_type`, returning an opaque `media_url` with compatibility URL fields.
- Added server-side MP4/MOV, MP3, and PCM 16/24-bit WAV validation. Enforced 10 MiB image, 100 MiB video, and 15 MiB audio limits; enforced 2-30 second audio duration for both WAV and MP3 frame streams.
- Updated the video workbench so customers select reference videos and audio files directly. The page previews selected media, rejects unsupported files promptly, enforces the LeoStudio reference counts, uploads media in parallel, and builds only `video_reference_base`/`audio_reference` guidance with `type: "UPLOADED"`.
- Prevented audio-only reference submissions, kept image and start/end frame compatibility, and kept provider UUIDs out of the customer workflow.
- Updated the standalone API documentation page, English/Chinese copy, and `docs/` contract notes with upload fields, limits, media guidance, and the audio visual-reference requirement.

### Testing
- From `backend/`, `go test ./internal/service ./internal/handler -count=1` passed, including new MP4/MOV/MP3/WAV store and multipart handler coverage.
- From `frontend/`, the three focused Vitest files passed: 30 tests total, including media upload guidance, unsupported-format prompts, audio visual-reference blocking, and UUID absence.
- From `frontend/`, `npm run typecheck`, targeted ESLint, and `npm run build` passed. The build retained the repository's existing dynamic-import and large-chunk warnings.
- `gofmt` and `git diff --check` passed. No production API key was used; no commit, push, or deployment was performed.

### Notes
- `backend/internal/handler/video_input.go`: parses image/video/audio multipart fields, maps validation errors, and returns media metadata.
- `backend/internal/handler/video_input_test.go`: covers video and audio multipart uploads and response fields.
- `backend/internal/service/video_input_store.go`: stores typed media and validates containers, WAV encoding, MP3 duration, limits, and restart MIME detection.
- `backend/internal/service/video_input_store_test.go`: covers valid and invalid media formats, duration bounds, size limits, and restart behavior.
- `frontend/src/api/videoGeneration.ts`: adds typed media upload responses and media-kind multipart uploads.
- `frontend/src/api/__tests__/videoGeneration.spec.ts`: verifies dedicated video/audio multipart fields.
- `frontend/src/views/user/VideoGenerationView.vue`: adds media selectors, previews, client validation, parallel uploads, and guidance assembly.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: verifies media guidance, format warnings, audio pairing, and UUID absence.
- `frontend/src/views/user/VideoApiDocsView.vue` and its spec: documents all upload fields and media guidance examples.
- `frontend/src/i18n/locales/en/dashboard.ts` and `frontend/src/i18n/locales/zh/dashboard.ts`: add media limits, format, and validation copy.
- `docs/LEO_VIDEO_CHANNEL.md` and `docs/LEO_VIDEO_MODEL_SPECS.md`: record the customer media upload contract.
- `progress.md`: records this task without rewriting prior history.
- Rollback point: no commit was created. Because this worktree already contained earlier Leo parameter changes, do not run a whole-file `git restore` or reset; review `git diff` and reverse only the new video/audio upload hunks while preserving the earlier changes and `.superpowers/`.

## 2026-07-26 - Task: Confirm independent start/end frame selection

### What was done
- Confirmed the customer-facing video page treats start-frame and end-frame images as independent inputs that can be uploaded together, while keeping reference images mutually exclusive with either frame input.
- Added regression assertions for end-frame availability after selecting a start frame and for submit availability when both frames are present.

### Testing
- From `frontend/`, `npm run test:run -- src/views/user/__tests__/VideoGenerationView.spec.ts` passed: 24 tests.

### Notes
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: locks the intended frame-selection and submit-button behavior.
- `progress.md`: records this verification.
- Rollback point: revert only the two added assertions in `VideoGenerationView.spec.ts` and remove this appended progress block; no production code or prior Leo changes are affected.

## 2026-07-26 - Task: Clarify frame-mode behavior in the web workbench and API docs

### What was done
- Updated the English and Chinese workbench guidance so customers can see that start and end frames may be submitted together, while reference images cannot be combined with either frame.
- Synchronized the same rule in the customer-facing API documentation description and the formal Leo video model specification.

### Testing
- Frontend focused Vitest suites passed: 30 tests across video API, workbench, and API-doc pages.
- Frontend `npm run typecheck`, targeted ESLint, and `npm run build` passed. Build output retained only the repository's existing dynamic-import, chunk-size, and Browserslist warnings.
- Backend `go test ./internal/service ./internal/handler -count=1` passed.
- `git diff --check` passed.

### Notes
- `frontend/src/i18n/locales/en/dashboard.ts`: clarifies English frame/reference input guidance.
- `frontend/src/i18n/locales/zh/dashboard.ts`: clarifies Chinese frame/reference input guidance.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: records independent frame submission and reference-image exclusion.
- `progress.md`: records the release preparation verification.
- Rollback point: revert the latest documentation-only hunks in the two dashboard locale files and `LEO_VIDEO_MODEL_SPECS.md`, then remove this appended block; preserve the earlier Leo implementation and `.superpowers/`.

## 2026-07-26 - Task: Publish Leo video reference workflow v0.1.164-fy.5

### What was done
- Committed the Leo video reference-media workflow, frame-mode clarification, web/API documentation, tests, and release preparation as `722c3c524`.
- Pushed `codex/leo-video-channel` to `origin` and published tag `v0.1.164-fy.5`.
- GitHub Release completed with the Linux amd64 archive and `checksums.txt`.

### Testing
- Full backend `go test ./... -count=1` passed.
- Frontend focused tests, typecheck, targeted ESLint, and production build passed before release.
- Release archive download returned HTTP 200 with the expected `application/octet-stream` content and 35,333,993-byte length.
- `https://api.fflink.top/health` returned HTTP 200 with `{"status":"ok"}`.
- Production binary replacement was not performed: this environment has no authenticated SSH or admin updater credential. No production state was changed.

### Notes
- `progress.md`: records the commit, tag, release assets, verification evidence, and deployment boundary.
- Rollback point: source rollback is `git revert 722c3c524`; release rollback target is `v0.1.164-fy.4`. The production server remains on its existing version because no online replacement was attempted.

## 2026-07-27 - Task: Merge upstream v0.1.165 and prepare fork release

### What was done
- Merged `upstream/main` at `7d3a896fc` into `codex/leo-video-channel`, advancing the upstream base to `0.1.165` while retaining the fork-owned Leo video channel, token incentive, updater, and other custom behavior.
- Integrated upstream ChatGPT Live gateway support, Claude Opus 5 compatibility, request-driven Ollama Cloud usage refresh, session ID persistence, announcement preview/styling, email-alias registration deduplication, scheduler and gateway fixes, and migrations 187 through 190.
- Resolved `backend/internal/service/admin_group.go` by retaining both the fork's Leo video price validation and upstream's `AllowLive` create/update normalization.
- Regenerated `frontend/pnpm-lock.yaml` so the fork's Axios/PostCSS requirements and upstream's PostCSS security override resolve together at PostCSS 8.5.23.
- Added the two missing Live capability mocks required by the existing GroupsView tests and aligned the existing Leo pricing validation message with its test contract.

### Testing
- `go test ./... -count=1` passed.
- The unchanged Redis subscriber timing test passed 10 consecutive isolated repetitions after one transient failure observed only while the full backend and frontend suites were competing for resources.
- Frontend Vitest passed: 197 files and 1,366 tests.
- `go test -tags unit ./internal/repository -count=1` passed.
- The first tagged service run exposed the Leo error-message mismatch; after the wording-only correction, `go test -tags unit ./internal/service -count=1` passed.
- `pnpm.cmd typecheck`, `pnpm.cmd lint:check`, and `pnpm.cmd run build` passed.
- `gofmt` and `git diff --check` passed; no unresolved conflict markers or sensitive configuration files were found.
- Control-flow parity appears preserved: fork-owned Leo validation remains active, upstream Live normalization is active, and full backend/frontend regression suites pass.

### Notes
- `README.md`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `README_CN.md`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `README_JA.md`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `assets/partners/logos/apikl.png`: 按上游 v0.1.165 清理已移除的资产。
- `assets/partners/logos/miyaip.png`: 按上游 v0.1.165 清理已移除的资产。
- `assets/partners/logos/tokeneum.png`: 按上游 v0.1.165 清理已移除的资产。
- `backend/cmd/server/VERSION`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/ent/group.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/ent/group/group.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/ent/group/where.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/ent/group_create.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/ent/group_update.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/ent/migrate/schema.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/ent/mutation.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/ent/runtime/runtime.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/ent/schema/group.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/config/config.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/domain/constants.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/admin/account_ollama_cloud_usage_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/admin/group_handler.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/auth_oauth_pending_flow_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/batch_image_handler.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/dto/mappers.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/dto/types.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/gateway_handler.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/gateway_handler_chat_completions.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/gateway_handler_responses.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/gemini_v1beta_handler.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/grok_media.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/openai_alpha_search.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/openai_chat_completions.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/openai_embeddings.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/openai_gateway_handler.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/openai_images.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/openai_images_failover_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/handler/openai_live.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/handler/openai_live_test.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/handler/user_handler_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/pkg/claude/constants.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/platform/liveattestation/attestation.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/platform/liveattestation/attestation_darwin.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/platform/liveattestation/attestation_unsupported.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/platform/liveattestation/attestation_unsupported_test.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/repository/account_repo_ollama_cloud_usage.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/account_repo_ollama_cloud_usage_integration_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/account_repo_ollama_cloud_usage_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/api_key_repo.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/batch_image_repo.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/concurrency_cache.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/concurrency_cache_integration_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/concurrency_cache_live_test.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/repository/gateway_cache.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/gateway_cache_live_test.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/repository/group_repo.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/integration_harness_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/migrations_schema_integration_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/usage_log_repo_insert.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/usage_log_repo_query.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/usage_log_repo_request_type_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/usage_log_session_id_integration_test.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/repository/usage_log_session_id_unit_test.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/repository/user_repo.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/repository/user_repo_email_alias_test.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/server/api_contract_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/server/middleware/admin_auth_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/server/routes/admin.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/server/routes/gateway.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/server/routes/prompt_audit_route_coverage_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/account.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/admin_group.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/admin_group_duplicate.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/admin_group_duplicate_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/admin_service.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/admin_service_apikey_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/admin_service_delete_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/admin_service_email_identity_sync_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/admin_service_group_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/api_key_auth_cache.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/api_key_auth_cache_impl.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/auth_oauth_email_flow.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/auth_service.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/auth_service_email_bind_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/auth_service_register_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/batch_image.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/batch_image_processor_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/batch_image_public.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/batch_image_public_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/batch_image_settlement.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/batch_image_settlement_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/bedrock_request.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/billing_service.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/claude_opus5_test.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/service/content_moderation_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/gateway_anthropic_apikey_passthrough_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/gateway_anthropic_passthrough.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/gateway_forward.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/gateway_upstream_response.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/gateway_usage_billing.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/gemini_chat_completions_compat_service.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/gemini_chat_completions_compat_service_test.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/service/gemini_messages_compat_service.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/group.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/ollama_cloud_usage.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/ollama_cloud_usage_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_account_runtime_block_fastpath.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_account_runtime_block_fastpath_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_codex_transform.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_gateway_apikey_item_id_test.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/service/openai_gateway_forward.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_gateway_grok.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_gateway_grok_cache.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_gateway_grok_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_gateway_scheduling.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_gateway_service.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_gateway_usage.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_live.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/service/openai_live_attestation.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/service/openai_live_lifecycle_test.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/service/openai_live_test.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/service/openai_live_types.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/service/openai_oauth_passthrough_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_responses_item_id.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/service/openai_responses_namespace.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_responses_namespace_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_responses_rejected_field_retry_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_upstream_transport_error.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_upstream_transport_error_handle_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/openai_ws_forwarder_success_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/pricing_service.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/pricing_service_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/registration_email_alias.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/service/registration_email_alias_test.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/service/session_id.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/service/session_id_test.go`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/service/usage_log.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/user_service.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/internal/service/user_service_test.go`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `backend/migrations/187_add_usage_log_session_id.sql`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/migrations/188_allow_live_usage_request_type.sql`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/migrations/189_add_group_allow_live.sql`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/migrations/190_add_users_email_alias_dedup_index_notx.sql`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/resources/model-pricing/model_prices_and_context_window.json`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `deploy/config.example.yaml`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/package.json`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/pnpm-lock.yaml`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/api/__tests__/admin.accounts.ollamaCloudUsage.spec.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/api/admin/groups.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/components/account/AccountStatusIndicator.vue`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/components/admin/usage/UsageFilters.vue`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/components/admin/usage/UsageTable.vue`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/components/common/AnnouncementBell.vue`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/components/common/AnnouncementPopup.vue`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/components/common/__tests__/AnnouncementPopup.spec.ts`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `frontend/src/composables/useModelWhitelist.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/i18n/locales/en/admin/overview.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/i18n/locales/en/admin/resources.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/i18n/locales/en/admin/settings.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/i18n/locales/en/dashboard.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/i18n/locales/zh/admin/overview.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/i18n/locales/zh/admin/resources.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/i18n/locales/zh/admin/settings.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/i18n/locales/zh/dashboard.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/styles/announcement-markdown.css`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `frontend/src/types/index.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/utils/errorBadges.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/utils/usageRequestType.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/views/admin/AnnouncementsView.vue`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/views/admin/GroupsView.vue`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/views/admin/SettingsView.vue`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/views/admin/UsageView.vue`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/views/admin/__tests__/SettingsView.spec.ts`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/views/user/AffiliateView.vue`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/views/user/UsageView.vue`: 同步上游 v0.1.165 的实现或测试改动，并保留当前分支持有功能。
- `frontend/src/views/user/__tests__/AffiliateView.spec.ts`: 新增上游 v0.1.165 的实现、迁移或回归测试资产。
- `backend/internal/service/channel_service.go`: 对齐 Leo 视频价格区间校验错误文案与既有测试预期，不改变计费规则。
- `frontend/src/views/admin/__tests__/GroupsView.columnSettings.spec.ts`: 补充 Live 能力查询 mock，覆盖合并后的分组页测试依赖。
- `frontend/src/views/admin/__tests__/GroupsView.duplicate.spec.ts`: 补充 Live 能力查询 mock，覆盖合并后的分组复制测试依赖。
- `progress.md`: 追加本轮合并、验证、文件清单和回滚记录。
- `.superpowers/`: remains untracked and is explicitly excluded from the merge commit and release.
- Rollback point: before commit, run `git merge --abort`; after commit, run `git revert -m 1 <merge_commit>`. The exact pre-merge source point is `a5cabf05253d06f358df019d52c0ffd3c720945a`.

## 2026-07-27 - Task: Publish upstream merge as v0.1.165-fy.1

### What was done
- Created merge commit `b4791c850a225fb9dffe10c5b6dd5cff9944461d` with parents `a5cabf05253d06f358df019d52c0ffd3c720945a` and upstream `7d3a896fcd912e6141b9bbbe90ca5bd08ff49ea3`.
- Pushed `codex/leo-video-channel`, created annotated tag `v0.1.165-fy.1`, and published the GitHub Release with the Linux amd64 archive and `checksums.txt`.

### Testing
- Both public Release asset URLs returned HTTP 200.
- Downloaded `sub2api_0.1.165-fy.1_linux_amd64.tar.gz`; its size is 35,411,988 bytes and SHA256 is `fb92909529d85fc6b087bc45c23ae1df6dbe4dae4b979960fd280cb1701711df`, exactly matching `checksums.txt`.
- Listed the archive successfully and confirmed it contains executable `sub2api` with an uncompressed size of 114,483,362 bytes.
- Confirmed the remote annotated tag exists and the branch push completed successfully.

### Notes
- `progress.md`: records the merge commit, tag, public release assets, checksum verification, and rollback point.
- `.superpowers/` remains untracked and was not included in the merge commit, tag, or release.
- Rollback point: source rollback is `git revert -m 1 b4791c850a225fb9dffe10c5b6dd5cff9944461d`; release rollback target is `v0.1.164-fy.5`. To withdraw this tag after removing the GitHub Release, run `git push origin :refs/tags/v0.1.165-fy.1`.

## 2026-07-27 - Task: Remove internal UUID wording from customer video docs

### What was done
- Removed the customer-facing English and Chinese wording that discussed provider UUIDs from the video API documentation upload description.
- Kept the customer instruction focused on uploading supported media and using the returned `media_url`.

### Testing
- `VideoApiDocsView.spec.ts`: 1 test passed.
- Targeted ESLint for both dashboard locale files passed.
- Confirmed no UUID wording remains in the customer-facing video docs or video locale strings.

### Notes
- `frontend/src/i18n/locales/en/dashboard.ts`: removes provider UUID wording from the English upload description.
- `frontend/src/i18n/locales/zh/dashboard.ts`: removes customer UUID wording from the Chinese upload description.
- `progress.md`: records this customer-facing copy correction.
- Rollback point: revert the next copy-fix commit and remove this appended block; preserve internal Leo compatibility documentation and the prior `v0.1.165-fy.1` release.

## 2026-07-27 - Task: Publish customer-doc UUID wording fix as v0.1.165-fy.2

### What was done
- Published the Web documentation copy correction as commit `967787768` and tag `v0.1.165-fy.2`.
- GitHub Release completed with the Linux amd64 archive and `checksums.txt`; the internal UUID compatibility path remains out of customer-facing copy.

### Testing
- Release workflow run 51 completed successfully.
- Both release assets are present on the public release page.
- The API docs test, locale ESLint, and frontend production build passed before release.

### Notes
- `progress.md`: records the copy-fix release and release verification.
- Rollback point: source rollback is `git revert 967787768`; release rollback target is `v0.1.165-fy.1`.

## 2026-07-27 - Task: Remove upstream implementation details from customer video surfaces

### What was done
- Removed the upstream brand, internal upload address example, and internal-field privacy notice from the video generation page and API documentation.
- Changed the customer model/specification and API guide wording to require uploaded media URLs only; customer materials no longer describe UUID compatibility fields or upstream provider terminology.
- Sanitized completed asynchronous video results before returning them to API clients, recursively removing provider metadata, UUIDs, source CDN URLs, provider task/account identifiers, and credential-shaped fields while retaining the platform job and local content URLs.

### Testing
- `pnpm test:run -- src/views/user/__tests__/VideoApiDocsView.spec.ts src/views/user/__tests__/VideoGenerationView.spec.ts src/api/__tests__/videoGeneration.spec.ts`: passed.
- `pnpm typecheck`: passed.
- `go test ./internal/handler ./internal/service -run 'TestLeoVideo|TestVideoOutputStore|TestVideoJob' -count=1`: passed.
- `git diff --check`: passed.
- Searched customer video views, API client types, locale strings, and customer docs for UUID, LeoStudio, upstream/provider wording, and internal upload paths; no matches remain.

### Notes
- `backend/internal/handler/leo_video_async.go`: recursively strips upstream metadata from public completed-job results.
- `backend/internal/handler/leo_video_async_test.go`: adds public-result leakage regression coverage and updates provider-name sanitization assertions.
- `frontend/src/views/user/VideoApiDocsView.vue`: removes the internal privacy callout and internal upload URL example.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: prevents upstream/internal terms from returning to rendered API docs.
- `frontend/src/api/videoGeneration.ts`: removes the provider field from the public result type.
- `frontend/src/i18n/locales/en/dashboard.ts`: removes upstream video wording from the English customer UI.
- `frontend/src/i18n/locales/zh/dashboard.ts`: removes upstream video wording from the Chinese customer UI.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: documents only the customer upload-URL contract and platform capability rules.
- `docs/SEEDANCE_API_CLIENT_GUIDE_CN.md`: removes internal upload URLs and upstream/provider wording from the customer integration guide.
- Rollback point: changes are uncommitted; reverse only the hunks listed above (or apply a saved patch in reverse) and preserve unrelated working-tree changes. Do not restore whole files because they contain other pending work.

## 2026-07-27 - Task: Detail video and audio reference upload conditions

### What was done
- Added the verified media constraints to the customer video workbench and API documentation: readable MP4/MOV video containers up to 100 MiB and at most three references; readable MP3 or PCM 16/24-bit WAV audio up to 15 MiB, 2–30 seconds, and at most one reference.
- Added the media-format details to the model specification: ISO Base Media `ftyp` requirement for video, MP3 frame or RIFF/WAVE PCM validation for audio, and the required audio-plus-visual reference pairing.
- Added immediate client-side rejection and a localized message for reference videos over 100 MiB; backend validation remains the final enforcement layer.

### Testing
- `pnpm test:run -- src/views/user/__tests__/VideoGenerationView.spec.ts src/views/user/__tests__/VideoApiDocsView.spec.ts src/api/__tests__/videoGeneration.spec.ts`: passed.
- `pnpm typecheck`: passed.
- `go test ./internal/service ./internal/handler -run 'TestVideoInput|TestLeoVideo' -count=1`: passed.
- `git diff --check`: passed.

### Notes
- `frontend/src/views/user/VideoGenerationView.vue`: rejects reference videos over 100 MiB before upload.
- `frontend/src/i18n/locales/zh/dashboard.ts`: details Chinese video/audio file, size, count, encoding, duration, and pairing rules.
- `frontend/src/i18n/locales/en/dashboard.ts`: details English video/audio file, size, count, encoding, duration, and pairing rules.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: covers the immediate oversized-video warning.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: records the byte/container/encoding validation contract without exposing upstream identifiers.
- Rollback point: changes are uncommitted; reverse only this task's hunks in the listed files and preserve earlier pending work in the same files. Do not restore whole files.

## 2026-07-27 - Task: Reconcile complete reference-media upload conditions

### What was done
- Rechecked the local LeoStudio media uploader, guidance normalizer, upload tests, and the current official Seedance API guides to separate documented limits from local validation and live observations.
- Expanded the model specification with direct URL requirements, queue-time URL reachability, MP4/MOV `ftyp` validation, supported audio response types, URL-mode duration backfill, and the absence of a published universal video duration/frame-rate/codec/dimension whitelist.
- Recorded the verified dimension observation that `640x360` was rejected with `INVALID_HEIGHT` while `864x496` succeeded, without turning that observation into an unsupported full allowlist.
- Updated Chinese and English customer-facing page/API copy with the complete, platform-neutral media conditions; no upstream brand, provider identifier, UUID, account, or internal URL was added.

### Testing
- `pnpm exec vitest run src/views/user/__tests__/VideoGenerationView.spec.ts src/views/user/__tests__/VideoApiDocsView.spec.ts --reporter=verbose`: 2 files and 25 tests passed.
- `pnpm exec vue-tsc --noEmit`: passed.
- `go test ./internal/service ./internal/handler -run 'TestVideoInput|TestLeoVideo' -count=1`: passed.
- `git diff --check`: passed.
- Direct inspection of LeoStudio `internal/leonardo/media.go`, `internal/service/video_guidance.go`, upload tests, and the official Seedance 2.0/ Fast/ Mini guides completed; no explicit upstream frame-rate, codec, duration, or complete reference-video dimension whitelist was found.

### Notes
- `docs/LEO_VIDEO_MODEL_SPECS.md`: records the verified local/upload and direct-URL media conditions, including the observed dimension compatibility fact.
- `frontend/src/i18n/locales/en/dashboard.ts`: expands English workbench and API documentation media requirements.
- `frontend/src/i18n/locales/zh/dashboard.ts`: expands Chinese workbench and API documentation media requirements.
- `progress.md`: records this source reconciliation and verification evidence.
- Rollback point: changes are uncommitted; reverse only this task's hunks in the three listed files (or apply the saved patch in reverse) and preserve earlier pending work. Do not restore whole files.

## 2026-07-27 - Task: Finalize LeoStudio model and reference-media alignment

### What was done
- Completed the final capability alignment for Seedance 2.0, Fast, Mini, Happy Horse, and Grok, including model-specific resolution, duration, aspect-ratio, prompt, frame, and guidance limits.
- Kept customer workbench/API documentation limited to platform upload responses and supported media conditions; removed the remaining upstream-identifier wording from the public model specification.
- Completed public asynchronous-result sanitization and deterministic billing-resolution normalization for the expanded model set.

### Testing
- `go test ./internal/service ./internal/handler -run 'TestLeoVideo|TestValidateLeoVideo|TestVideoInput|TestNormalizeVideoBilling' -count=1`: passed.
- `pnpm test:run -- src/views/user/__tests__/VideoApiDocsView.spec.ts src/views/user/__tests__/VideoGenerationView.spec.ts src/api/__tests__/videoGeneration.spec.ts`: 31 tests passed.
- `pnpm exec vue-tsc --noEmit`: passed.
- `git diff --check`: passed.
- Audited customer video views, API types, bilingual locale strings, and public docs for UUID, internal media paths, provider credentials, and upstream task identifiers; no customer-surface matches remain.

### Notes
- `backend/internal/handler/leo_video_async.go`: strips private provider-shaped fields from completed public results.
- `backend/internal/handler/leo_video_async_test.go`: covers result sanitization and public error wording.
- `backend/internal/handler/leo_video_test.go`: updates Leo video handler expectations.
- `backend/internal/service/leo_video.go`: sanitizes synchronous video-service wording.
- `backend/internal/service/leo_video_async.go`: sanitizes asynchronous video-service wording.
- `backend/internal/service/leo_video_async_test.go`: updates asynchronous sanitization assertions.
- `backend/internal/service/leo_video_model_specs.go`: defines the latest model and media guidance validation matrix.
- `backend/internal/service/leo_video_model_specs_test.go`: covers model-specific parameter and guidance boundaries.
- `backend/internal/service/leo_video_test.go`: updates video validation tests.
- `backend/internal/service/video_billing_resolution.go`: normalizes nonstandard billing resolutions to available tiers.
- `backend/internal/service/video_billing_resolution_test.go`: covers billing normalization.
- `frontend/src/api/videoGeneration.ts`: aligns public request/result types and hides provider metadata.
- `frontend/src/api/__tests__/videoGeneration.spec.ts`: verifies public API typing behavior.
- `frontend/src/views/user/VideoGenerationView.vue`: exposes model-filtered parameters and media upload validation.
- `frontend/src/views/user/VideoApiDocsView.vue`: documents supported models, media uploads, and public request examples.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: covers workbench parameter and upload constraints.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: guards against upstream details in rendered docs.
- `frontend/src/i18n/locales/en/dashboard.ts`: updates English video guidance and API documentation copy.
- `frontend/src/i18n/locales/zh/dashboard.ts`: updates Chinese video guidance and API documentation copy.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: records the public model and media contract.
- `docs/LEO_VIDEO_CHANNEL.md`: records channel configuration and operational capability alignment.
- `progress.md`: appends this final verification record.
- Rollback point: changes remain uncommitted; export the target file hunks with `git diff HEAD -- <file> > rollback.patch` and apply only the selected hunks in reverse with `git apply -R rollback.patch`; do not restore whole files because they contain earlier pending work.

## 2026-07-27 - Task: Prepare Leo model and media capability release v0.1.165-fy.3

### What was done
- Consolidated the pending Leo video model and reference-media capability changes for release: model-specific limits, Happy Horse/Grok support, Mini resolution/aspect expansion, prompt enhancement rules, frame/reference counts, media validation, and billing-tier normalization.
- Sanitized synchronous and asynchronous public video responses and error messages so provider identifiers, credentials, upstream task IDs, and private media fields do not reach customers.
- Updated the customer workbench, API types, bilingual copy, API documentation, and Leo operational specifications to match the validated public contract.

### Testing
- `go test ./internal/service ./internal/handler -run 'TestLeoVideo|TestValidateLeoVideo|TestVideoInput|TestNormalizeVideoBilling|TestVideoOutputStore|TestVideoJob' -count=1` passed.
- Frontend Vitest passed with exit code 0; the command executed the complete repository suite while including the three Leo-focused files.
- `pnpm.cmd typecheck` passed.
- `pnpm.cmd lint:check` passed.
- `pnpm.cmd run build` passed; only the repository's existing Browserslist, dynamic-import, and chunk-size warnings remained.
- `gofmt -d` for all modified Go files produced no output, and `git diff --check` passed.

### Notes
- `backend/internal/handler/leo_video_async.go`: strips private fields from public completed-job results.
- `backend/internal/handler/leo_video_async_test.go`: covers public result sanitization and error wording.
- `backend/internal/handler/leo_video_test.go`: updates synchronous video handler expectations.
- `backend/internal/service/leo_video.go`: sanitizes synchronous video-service errors and validates requests.
- `backend/internal/service/leo_video_async.go`: sanitizes asynchronous video-service errors.
- `backend/internal/service/leo_video_async_test.go`: covers asynchronous sanitization behavior.
- `backend/internal/service/leo_video_model_specs.go`: defines model-specific resolution, duration, prompt, frame, and media limits.
- `backend/internal/service/leo_video_model_specs_test.go`: covers the expanded model validation matrix.
- `backend/internal/service/leo_video_test.go`: updates video request validation coverage.
- `backend/internal/service/video_billing_resolution.go`: normalizes model resolution values to billable tiers.
- `backend/internal/service/video_billing_resolution_test.go`: covers billing resolution normalization.
- `frontend/src/api/videoGeneration.ts`: aligns public video request and result types.
- `frontend/src/api/__tests__/videoGeneration.spec.ts`: verifies public API typing behavior.
- `frontend/src/i18n/locales/en/dashboard.ts`: updates English video capability and media guidance.
- `frontend/src/i18n/locales/zh/dashboard.ts`: updates Chinese video capability and media guidance.
- `frontend/src/views/user/VideoApiDocsView.vue`: updates customer-facing API examples and constraints.
- `frontend/src/views/user/VideoGenerationView.vue`: exposes model-filtered parameters and media validation.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: prevents internal provider details from returning to API docs.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: covers workbench model and upload constraints.
- `docs/LEO_VIDEO_CHANNEL.md`: records operational capability and billing alignment.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: records the public model and media contract.
- `progress.md`: records this release preparation and verification evidence.
- `.superpowers/`: remains untracked and excluded from the commit and release.
- Rollback point: before commit, preserve `HEAD` at `0c7690880e0`; after commit, use `git revert <release_commit>` or restore the prior release tag `v0.1.165-fy.2`. Do not restore whole files without preserving unrelated pending work.

## 2026-07-27 - Task: Publish Leo model and media capability release v0.1.165-fy.3

### What was done
- Published commit `3799e5d321edea6b488869fd9c82d9ddffa4c239` and annotated tag `v0.1.165-fy.3` to `origin`.
- GitHub Release completed with the Linux amd64 archive and `checksums.txt`.

### Testing
- The public Linux asset returned HTTP 200 and downloaded successfully.
- `sub2api_0.1.165-fy.3_linux_amd64.tar.gz` is 35,416,050 bytes with SHA256 `0710a0b198b6201037f769b0607192304cfa3d21b4aa344618e4d719bd72b85b`, exactly matching `checksums.txt`.
- The archive listing contains executable `sub2api` with an uncompressed size of 114,503,842 bytes.
- Confirmed `v0.1.165-fy.3` points to commit `3799e5d321edea6b488869fd9c82d9ddffa4c239`; `.superpowers/` remains untracked and was not included.

### Notes
- `progress.md`: records the release commit, tag, package verification, and rollback point.
- Rollback point: source rollback is `git revert 3799e5d321edea6b488869fd9c82d9ddffa4c239`; release rollback target is `v0.1.165-fy.2`. After removing the GitHub Release, withdraw the tag with `git push origin :refs/tags/v0.1.165-fy.3`.

## 2026-07-27 - Task: Apply rightmost video retail prices to legacy billing tiers

### What was done
- Mapped the rightmost table values to the existing three video billing tiers: Happy Horse `0.15/0.15/0.19` and Grok Imagine 1.5 `0.10/0.17/0.17` USD/s for `480p/720p/1080p`.
- Changed billing normalization so Grok `400p` and `544p` use the low tier, `720p` uses the middle tier, and `960p` uses the high tier.
- Documented the compatibility mapping and the unsupported Happy Horse 480p placeholder tier.

### Testing
- `go test ./internal/service -run 'TestNormalizeVideoBillingResolutionLeo|TestLeoVideoPricingResolutions|TestVideoPriceConfigFromResolvedPricing|TestVideoJobBilling' -count=1`: passed.
- `git diff --check`: passed; only existing line-ending warnings were reported.

### Notes
- `backend/internal/service/video_billing_resolution.go`: aligns 544p with the low compatibility tier.
- `backend/internal/service/video_billing_resolution_test.go`: covers the 544p mapping.
- `docs/LEO_VIDEO_CHANNEL.md`: records the rightmost-price mapping and its limitations.
- `progress.md`: records this pricing change and verification status.
- Rollback point: revert the latest hunks in the three files, or use `git diff HEAD -- <file> > rollback.patch` followed by `git apply -R rollback.patch`; preserve `.superpowers/` and unrelated work.

## 2026-07-27 - Task: Configure production pricing for new video models

### What was done
- Confirmed production is running `0.1.166-fy.1`.
- Updated production channel `5` (`Seedance 2 视频专用渠道`) with rightmost-table pricing for `happy-horse-1.1` and `grok-imagine-1.5`.
- Preserved all three existing Seedance pricing entries and left group `25` fallback prices unchanged.

### Testing
- Read back channel `5` after the update: five unique video model pricing entries are present and both new entries have the expected three compatibility tiers.
- No video generation task was submitted in this pricing operation.

### Notes
- Production configuration target: channel `5`; rollback by removing only the two added model pricing entries and restoring the pre-update channel payload.
- `progress.md`: records the production pricing change and rollback target.

## 2026-07-27 - Task: Repair UTF-8 channel metadata after pricing update

### What was done
- Restored production channel `5` name to `Seedance 2 视频专用渠道`.
- Restored its description to `Seedance 2.0 / Fast Leo 视频专用渠道`.
- Resent the update with explicit UTF-8 bytes and preserved all five model pricing entries.

### Testing
- Read back the channel after repair: Chinese name/description are intact; Happy Horse and Grok prices remain `0.15/0.15/0.19` and `0.10/0.17/0.17`.

### Notes
- Root cause: PowerShell string request bodies were encoded with the console code page instead of UTF-8, converting Chinese characters to `?`.
- Rollback point: restore the prior channel metadata only; pricing entries were not changed by the repair.

## 2026-07-27 - Task: Merge upstream v0.1.166 and prepare fork release

### What was done
- Merged upstream `main` at official release `v0.1.166` (`59ce11c78000bde5bdd74930b5885753037a5841`) into the fork release branch.
- Preserved the fork-specific Leo video, Token incentive, Cyber audit, update/restart, and related database migration and documentation paths; upstream release fixes for panel rate limiting, Antigravity, WebSocket turn billing, provider compatibility, settings, payments, usage statistics, Caddy SSE handling, and security dependencies were integrated.
- Resolved the `req/v3` checksum conflict in favor of upstream `v3.59.0`, and combined OpenAI WebSocket turn model/billing mapping with the fork's Cyber input excerpt auditing.
- Synchronized source metadata to `0.1.166` and upstream commit `59ce11c78000bde5bdd74930b5885753037a5841`; the planned fork release tag is `v0.1.166-fy.1`.

### Testing
- `go test ./...` (in `backend`): passed after resolving the video billing merge state.
- `pnpm.cmd test:run` (in `frontend`): passed.
- `pnpm.cmd typecheck` (in `frontend`): passed.
- `pnpm.cmd lint:check` (in `frontend`): passed.
- `pnpm.cmd run build` (in `frontend`): passed and regenerated embedded frontend assets; only existing Browserslist, dynamic-import, and chunk-size warnings remained.
- `gofmt -w backend/internal/handler/openai_gateway_handler.go` and `git diff --check`: passed.

### Notes
- `backend/go.sum`: uses the upstream `github.com/imroc/req/v3 v3.59.0` checksums and retains adjacent dependency entries.
- `backend/internal/handler/openai_gateway_handler.go`: combines upstream WebSocket per-turn channel/billing mapping with fork Cyber input excerpt recording.
- `backend/internal/service/video_billing_resolution.go`: preserves the verified 544p-to-480p compatibility billing mapping.
- `backend/internal/service/video_billing_resolution_test.go`: covers the 544p compatibility billing mapping.
- `backend/cmd/server/VERSION`: records the upstream base version `0.1.166`.
- `backend/cmd/server/UPSTREAM_COMMIT`: records upstream release commit `59ce11c78000bde5bdd74930b5885753037a5841`.
- `docs/UPDATE_POLICY.md`: advances the documented upstream synchronization baseline to `v0.1.166`.
- `docs/LEO_VIDEO_CHANNEL.md`: documents the current compatibility pricing tiers.
- `progress.md`: records this merge and verification evidence.
- Other staged files are the upstream `v0.1.166` backend, frontend, deployment, resource, README, and test changes included by the merge; `.superpowers/` remains untracked and excluded.
- Rollback point: before commit, preserve `HEAD` at `e9f68492e052fcbc12958d58974009f10509601c`; after commit, use `git revert -m 1 <merge_commit>` for the upstream merge and `git revert <release_metadata_commit>` for release metadata/log changes. To withdraw the release, delete the GitHub Release and push `git push origin :refs/tags/v0.1.166-fy.1`.

## 2026-07-27 - Task: Publish fork release v0.1.166-fy.1

### What was done
- Pushed merge commit `41718f7c026fa396371e43b5e80a1a1af0030294` to `origin/codex/leo-video-channel` and annotated tag `v0.1.166-fy.1` to `origin`.
- GitHub Actions Release workflow completed successfully and published the Linux amd64 archive and `checksums.txt`.

### Testing
- Release workflow `30255449166`: completed with conclusion `success`.
- Public Release API returned a non-draft, non-prerelease release with both expected assets.
- Downloaded `sub2api_0.1.166-fy.1_linux_amd64.tar.gz` successfully; size `35,488,010` bytes.
- SHA256 `79fbb57f74cade487a3e00a84a62757481058e7e45a3fc5d7dda010eeaa8bbff` exactly matched `checksums.txt`.
- Archive listing contains executable `sub2api`; remote tag and branch both resolve to `41718f7c026fa396371e43b5e80a1a1af0030294`.

### Notes
- `progress.md`: records the published release and package verification.
- `.superpowers/`: remains untracked and excluded from the commit and release.
- Rollback point: source rollback is `git revert 41718f7c026fa396371e43b5e80a1a1af0030294`; release rollback is to remove the GitHub Release and push `git push origin :refs/tags/v0.1.166-fy.1`.

## 2026-07-28 - Task: Add Happy Horse and Grok model mappings to production Leo account

### What was done
- Updated production Leo account `1682` without changing its API key, base URL, pool settings, status, concurrency, or existing Seedance mappings.
- Added `happy-horse-1.1 -> happy-horse-1.1` and `grok-imagine-1.5 -> grok-imagine-1.5` to the account model mapping.
- Documented the required production mapping and the distinction between model exposure and upstream capability.

### Testing
- Admin account readback returned all five mappings; account remained `active`, `schedulable=true`, with no error message.
- Admin available-models endpoint returned all five video models.
- The supplied ordinary API Key returned HTTP 200 from `/v1/models` and exposed both newly mapped models.
- No video generation task was submitted during this change.

### Notes
- `docs/LEO_VIDEO_CHANNEL.md`: documents the production Leo account mapping and its validation boundary.
- `progress.md`: records this production configuration change and verification evidence.
- Rollback point: update account `1682` with the prior three-entry `model_mapping` (`seedance-2.0`, `seedance-2.0-fast`, `seedance-2.0-mini`) while preserving the current credentials; do not delete or regenerate the account key.

## 2026-07-28 - Task: Live-test Happy Horse and Grok Imagine production models

### What was done
- Submitted one minimal Happy Horse job and one minimal Grok Imagine job through the production video API using the supplied ordinary API Key.
- Retried Grok once with a direct, no-redirect PNG after the first public Wikimedia URL failed during guidance upload.

### Testing
- Happy Horse: `720p`, `3s`, `16:9`, `audio=false`; completed successfully. Content endpoint returned HTTP `200`, `video/mp4`, `2,628,184` bytes.
- Grok first attempt: `400p`, `3s`, `16:9`, Wikimedia start frame; failed with `Video service guidance upload failed`.
- Grok retry: `400p`, `3s`, `16:9`, direct `placehold.co` PNG start frame; completed successfully. Content endpoint returned HTTP `200`, `video/mp4`, `3,063,693` bytes.
- No additional retries were made after the successful Grok result.

### Notes
- `progress.md`: records the production generation test and the isolated public-image URL failure.
- Rollback point: no source or production configuration rollback is required for this test; the two failed/completed test jobs remain in the API Key's video job history.

## 2026-07-28 - Task: Document per-model video API requests in the Web API docs

### What was done
- Added a dedicated Web API documentation section with copy-ready `curl` requests for `seedance-2.0`, `seedance-2.0-fast`, `seedance-2.0-mini`, `happy-horse-1.1`, and `grok-imagine-1.5`.
- Added per-model parameter descriptions for resolution, duration, aspect ratio, start-frame requirements, and guidance limits.
- Kept Happy Horse reference-video guidance documented as unsupported because the production gateway currently rejects it with the model limit set to zero.
- Added regression coverage that requires all five model request examples to remain present.

### Testing
- `pnpm.cmd exec vitest run src/views/user/__tests__/VideoApiDocsView.spec.ts src/views/user/__tests__/VideoGenerationView.spec.ts`: passed, 26 tests.
- `pnpm.cmd typecheck`: passed.
- ESLint passed for the changed API docs component and test.
- `git diff --check`: passed.

### Notes
- `frontend/src/views/user/VideoApiDocsView.vue`: adds the per-model API request section and examples.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: verifies all five model examples are rendered.
- `frontend/src/i18n/locales/zh/dashboard.ts`: adds Chinese model-example descriptions.
- `frontend/src/i18n/locales/en/dashboard.ts`: adds English model-example descriptions.
- `progress.md`: records this documentation update and verification evidence.
- Rollback point: revert the latest documentation and test hunks in the four files above; no production API or model configuration was changed.

## 2026-07-28 - Task: Verify embedded frontend build after model API docs update

### What was done
- Rebuilt the frontend bundle so the per-model API documentation is compile-checked in the embedded web assets.

### Testing
- `pnpm.cmd run build`: passed; Vite transformed 975 modules and emitted the embedded frontend bundle. Existing dynamic-import and chunk-size warnings remain non-blocking.

### Notes
- `progress.md`: records the final build verification for the documentation change.
- Rollback point: use the prior frontend build output and revert the documentation-only changes recorded in the preceding task entry.

## 2026-07-28 - Task: Publish video group model announcement

### What was done
- Published production announcement ID `20` for the two newly exposed video models: `happy-horse-1.1` and `grok-imagine-1.5`.
- Activated the announcement as an all-user popup and included the supported resolutions, duration range, and required image input constraints.

### Testing
- `GET /api/v1/admin/announcements/20` returned the expected Chinese title and model content with `status=active`, `notify_mode=popup`, and empty targeting for all users.
- The temporary announcement payload was removed after the verified production write; no administrator credential was written to the repository.

### Notes
- `progress.md`: records the production announcement, readback evidence, and rollback operation.
- Production state changed through the Admin API only; no source behavior, deployment, commit, or push was performed.
- Rollback point: archive announcement ID `20` with `PUT /api/v1/admin/announcements/20` and `{"status":"archived"}`.

## 2026-07-28 - Task: Publish fork release v0.1.166-fy.2

### What was done
- Pushed commit `60fabd8b10b3912e807867a0180c059dfa9afe1f` to `origin/codex/leo-video-channel` and annotated tag `v0.1.166-fy.2` to `origin`.
- Published the updated per-model video API documentation and model request examples through the GitHub Release workflow.

### Testing
- Targeted Vitest for `VideoApiDocsView` and `VideoGenerationView`: 26 tests passed.
- `pnpm.cmd typecheck` and `pnpm.cmd lint:check` passed.
- `pnpm.cmd run build` passed; Vite transformed 975 modules and emitted the embedded frontend bundle. Existing Browserslist, dynamic-import, and chunk-size warnings remain non-blocking.
- Public Release assets returned successfully; `sub2api_0.1.166-fy.2_linux_amd64.tar.gz` is `35,488,759` bytes.
- SHA256 `f20775158233ca9ee8ca1a480f71e4907c1481b6bb7b63fdc752165d26d6c11c` exactly matched `checksums.txt`.
- Archive listing contains executable `sub2api`; remote tag and branch both resolve to `60fabd8b10b3912e807867a0180c059dfa9afe1f`.

### Notes
- `progress.md`: records the release commit, tag, package verification, and rollback point.
- `.superpowers/`: remains untracked and excluded from the commit and release.
- Rollback point: source rollback is `git revert 60fabd8b10b3912e807867a0180c059dfa9afe1f`; release rollback is to remove the GitHub Release and push `git push origin :refs/tags/v0.1.166-fy.2`.

## 2026-07-28 - Task: Clarify per-model video parameters and reference-video format

### What was done
- Added a Web API documentation matrix covering each public model's resolution, duration/default, aspect-ratio combinations, prompt limit, and supported reference inputs.
- Clarified the reference-video flow: upload with multipart `video`, read `media_url`, and place it at `guidances.video_reference_base[].video` with `type: "UPLOADED"`; documented the Seedance-only restriction and three-video limit.
- Added common request-field rules and reference-video/audio examples to the standalone model specification document.
- Corrected the Web API upload response example so an image upload is represented as `image/png` instead of `video/mp4`.

### Testing
- `frontend/node_modules/.bin/vitest.cmd run src/views/user/__tests__/VideoApiDocsView.spec.ts --reporter=verbose`: passed, 1 test.
- `frontend/node_modules/.bin/vue-tsc.cmd --noEmit`: passed.
- `git diff --check`: passed; only the existing LF/CRLF normalization warning for the Markdown file remains.

### Notes
- `frontend/src/views/user/VideoApiDocsView.vue`: adds the visible model matrix and upload-format explanation.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: verifies the five-row matrix and uploaded video object format.
- `frontend/src/i18n/locales/zh/dashboard.ts`: adds Chinese matrix and reference-video wording.
- `frontend/src/i18n/locales/en/dashboard.ts`: adds English matrix and reference-video wording.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: adds detailed standalone request and media-format documentation.
- `progress.md`: records this documentation change and verification evidence.
- Rollback point: revert the latest hunks in the five files above; no production API, account, pricing, or model configuration was changed.

## 2026-07-28 - Task: Publish fork release v0.1.166-fy.3

### What was done
- Pushed commit `1c460cfabaf688e61c51e8325d461525cdf3b35f` to `origin/codex/leo-video-channel` and annotated tag `v0.1.166-fy.3` to `origin`.
- Published the updated per-model video parameter matrix and reference-media format documentation through the GitHub Release workflow.

### Testing
- Targeted Vitest for `VideoApiDocsView` and `VideoGenerationView`: 26 tests passed.
- `pnpm.cmd typecheck`, `pnpm.cmd lint:check`, `pnpm.cmd run build`, and `git diff --check` passed before release.
- GitHub Release workflow run `30327285187` completed successfully.
- Public Release assets returned successfully; `sub2api_0.1.166-fy.3_linux_amd64.tar.gz` is `35,490,700` bytes.
- SHA256 `cf745ef3cd2a3f3f6336b233f45bfb4d71f9842096167b1f64ed0a767949d17a` exactly matched `checksums.txt`.

### Notes
- `progress.md`: records the release commit, tag, workflow result, package verification, and rollback point.
- `.superpowers/`: remains untracked and excluded from the commit and release.
- Rollback point: source rollback is `git revert 1c460cfabaf688e61c51e8325d461525cdf3b35f`; release rollback is to remove the GitHub Release and push `git push origin :refs/tags/v0.1.166-fy.3`.

## 2026-07-28 - Task: Integrate LeoStudio LTX 2.3 video models

### What was done
- Added `ltxv-2.3-pro` and `ltxv-2.3-fast` to the Leo model capability registry, account defaults, channel-model lists, user workbench, public API documentation, and bilingual labels.
- Enforced the current model contracts: 1080p/1440p/2160p, fixed 16:9, Pro at 6/8/10 seconds, Fast at even durations from 6 through 20 seconds, start/end frames, generated audio, and prompt enhancement. Unsupported image, video, and audio references, `seed`, and `mode` are rejected before dispatch.
- Added the LTX single-tier billing compatibility rule: 1080p is the only configured USD-per-second tier, and requests/results at 1440p and 2160p normalize to it for reserve and settlement. Removed the stale internal assumption that every single-tier model must provide 720p.
- Updated the customer-facing video workbench and API docs without exposing provider account IDs, provider task IDs, UUIDs, or credentials.

### Testing
- `go test ./internal/service -run 'TestVideoJobBillingPrepareLTX|TestVideoJobSettlementLTX|TestVideoPriceConfigFromResolvedPricing|TestLeoVideoPricingResolutions|TestNormalizeVideoBillingResolutionLeo' -count=1`: passed.
- `go test ./internal/service -count=1`: passed.
- `go test ./internal/handler -count=1`: passed.
- `pnpm.cmd exec vitest run src/views/user/__tests__/VideoGenerationView.spec.ts src/views/user/__tests__/VideoApiDocsView.spec.ts src/api/__tests__/videoGeneration.spec.ts src/components/account/__tests__/CreateAccountModal.spec.ts src/components/admin/channel/__tests__/types.spec.ts`: 5 files and 62 tests passed.
- `pnpm.cmd exec vue-tsc --noEmit`: passed.
- `pnpm.cmd run build`: passed; Vite transformed 975 modules and refreshed the embedded frontend assets. Existing dynamic-import and chunk-size warnings remain non-blocking.
- `git diff --check`: passed; Git reported only LF-to-CRLF line-ending warnings for Markdown documents.

### Modified files
- Backend: `backend/internal/service/billing_service_test.go`, `leo_account.go`, `leo_account_test.go`, `leo_video.go`, `leo_video_model_specs.go`, `leo_video_model_specs_test.go`, `model_pricing_resolver.go`, `model_pricing_resolver_test.go`, `video_billing_resolution.go`, `video_billing_resolution_test.go`, `video_job_billing.go`, and `video_job_billing_test.go`.
- Frontend: `frontend/src/api/videoGeneration.ts`, `components/account/CreateAccountModal.vue`, `components/account/__tests__/CreateAccountModal.spec.ts`, `components/admin/channel/types.ts`, `components/admin/channel/__tests__/types.spec.ts`, `composables/useModelWhitelist.ts`, `constants/channel.ts`, `i18n/locales/en/dashboard.ts`, `i18n/locales/zh/dashboard.ts`, `views/user/VideoGenerationView.vue`, `views/user/VideoApiDocsView.vue`, `views/user/__tests__/VideoGenerationView.spec.ts`, and `views/user/__tests__/VideoApiDocsView.spec.ts`.
- Documentation: `docs/LEO_VIDEO_CHANNEL.md` and `docs/LEO_VIDEO_MODEL_SPECS.md`.

### Notes
- `.superpowers/` remains untracked and excluded from this task.
- Rollback point: working-tree base is `be80e81db6a63682db4157d9e74fcfb5948d9575`. After committing this task, roll back with `git revert <ltx_integration_commit>`; before committing, reverse only this task's reviewed diff rather than resetting the shared worktree.

## 2026-07-29 - Task: Restore model-specific video billing duration limits

### What was done
- Restored the shared video billing limit to 15 seconds so Grok and legacy video paths cannot inherit the LTX Fast 20-second capability.
- Added model-aware Leo duration normalization for LTX Fast and applied it consistently to direct Leo forwarding, async video-job reserve/settlement, OpenAI video cost calculation, and usage-log metadata.
- Added regression coverage for the 15-second shared limit, LTX Fast 20-second billing, Grok duration handling, direct Leo output metadata, async reserve/settlement, and usage recording.

### Testing
- `go test ./internal/service -run 'Test(NormalizeVideoBillingResolutionLeo|LeoVideoPricingResolutions|Calculate.*VideoCost|VideoJobBillingPrepareLTX|VideoJobSettlementLTX|ForwardLeoVideoPreservesLTXFastTwentySecondDuration|ParseGrokMediaVideoRequestClampsDurationToFifteenSeconds|LTXFastVideoUsageKeepsTwentySecondDuration)$' -count=1`: passed.
- `go test ./internal/service -run 'TestVideoJobBillingPrepareLTXRequiresOnly1080PCompatibilityPrice|TestLTXFastVideoUsageKeepsTwentySecondDuration|TestVideoJobSettlementLTXUses1080PCompatibilityPrice' -count=1`: passed after correcting the expected 20-second reserve amount.
- `go test ./internal/service -count=1`: passed (96.276s).
- `go test ./internal/handler -count=1`: passed (28.974s).
- `git diff --check`: passed; only existing Markdown LF/CRLF warnings remain.

### Notes
- `backend/internal/service/video_billing_resolution.go`: restores the shared 15-second cap and adds Leo model-aware duration normalization.
- `backend/internal/service/video_billing_resolution_test.go`: verifies shared 15-second and LTX Fast 20-second normalization behavior.
- `backend/internal/service/billing_service.go`: lets video cost calculation use the upstream capability model for duration limits.
- `backend/internal/service/billing_service_test.go`: updates the legacy 15-second assertion and adds LTX Fast 20-second cost coverage.
- `backend/internal/service/leo_video.go`: preserves valid LTX Fast 20-second direct response metadata.
- `backend/internal/service/leo_video_test.go`: covers direct LTX Fast 20-second forwarding.
- `backend/internal/service/openai_gateway_usage.go`: uses the upstream Leo model for video cost and usage-duration normalization.
- `backend/internal/service/openai_gateway_record_usage_test.go`: covers 20-second LTX Fast usage records and charges.
- `backend/internal/service/openai_gateway_grok_test.go`: verifies Grok duration remains capped at 15 seconds.
- `backend/internal/service/video_job_billing.go`: applies the upstream model limit during async reserve and settlement.
- `backend/internal/service/video_job_billing_test.go`: covers 20-second LTX Fast reserve and settlement.
- Rollback point: these files also contain the preceding uncommitted LTX integration changes; do not run `git restore` on them. To roll back this round safely, reverse only the hunks under the task heading above (or save the current `git diff` first and apply the inverse selectively).

## 2026-07-29 - Task: Publish fork release v0.1.166-fy.4

### What was done
- Pushed commit `f1fefc776db07ee5c7f0211e52b4dc8ab596f850` to `origin/codex/leo-video-channel` and annotated tag `v0.1.166-fy.4` to `origin`.
- Published the LeoStudio LTX 2.3 model integration and model-specific video billing duration fix while retaining upstream release `v0.1.166` as the version baseline.

### Testing
- `go test ./internal/service -count=1`: passed.
- `go test ./internal/handler -count=1`: passed.
- Targeted frontend Vitest suite: 5 files and 62 tests passed.
- `pnpm.cmd typecheck`, `pnpm.cmd lint:check`, `pnpm.cmd run build`, Go formatting verification, and `git diff --check` passed before release.
- Public Release assets returned successfully; `sub2api_0.1.166-fy.4_linux_amd64.tar.gz` is `35,500,561` bytes and contains the executable `sub2api` entry.
- SHA256 `c18ed07592571a55e590955d6f5474b11f1bf8d8137586d5f20dbe84281d6a0d` exactly matched `checksums.txt`.

### Notes
- `progress.md`: records the release commit, tag, verification evidence, package checksum, and rollback point.
- `.superpowers/`: remains untracked and excluded from the commit and release.
- Rollback point: source rollback is `git revert f1fefc776db07ee5c7f0211e52b4dc8ab596f850`; release rollback is to remove the GitHub Release and push `git push origin :refs/tags/v0.1.166-fy.4`.

## 2026-07-29 - Task: Merge upstream v0.1.168 and publish fork release v0.1.168-fy.1

### What was done
- Merged upstream formal release `v0.1.168` into the fork without merging unpublished upstream `main` changes beyond that tag.
- Preserved the fork's Leo video generation, async video billing, video menu switch, Token Incentive Plan, and related migrations while adding upstream Passkey authentication, Model Plaza, Kimi K3 support, setup bypass, scoped update fixes, and security/audit fixes.
- Resolved 22 conflicts across dependency injection, settings contracts, backend/frontend feature flags, and admin settings; also migrated the cyber-policy ban test path to the upstream explicit `UserUpdateFields` repository update contract and restored Leo platform color mappings.

### Testing
- `go test ./... -count=1`: passed for all backend packages.
- `pnpm.cmd test:run`: passed for the complete frontend Vitest suite.
- `pnpm.cmd typecheck`: passed.
- `pnpm.cmd lint:check`: passed.
- `pnpm.cmd run build`: passed; Vite transformed 994 modules and refreshed the embedded frontend assets. Existing dynamic-import, chunk-size, and Browserslist warnings remain non-blocking.
- Go formatting, merge-marker checks, and `git diff --cached --check`: passed.

### Notes
- `backend/cmd/server/wire_gen.go`: retains Token Incentive and adds Model Plaza/Passkey handler wiring.
- `backend/internal/handler/`, `backend/internal/service/`, and `frontend/src/views/admin/SettingsView.vue`: combine upstream settings and Model Plaza fields with fork video and Token Incentive fields.
- `backend/internal/service/ops_cyber_policy_ban.go`, `backend/internal/service/ops_cyber_policy_ban_test.go`: migrate the user update call and test stub to `UserUpdateFields`.
- `frontend/src/i18n/locales/en/admin/settings.ts`, `frontend/src/i18n/locales/zh/admin/settings.ts`: retain both video and Model Plaza settings translations.
- `frontend/src/utils/platformColors.ts`: adds missing Leo strong-border and accent mappings required by the merged platform type.
- Other files changed by this task are the upstream `v0.1.168` release files included in the merge commit, plus the existing fork files retained by the merge.
- `.superpowers/`: remains untracked and excluded from the commit and release.
- Rollback point: revert the merge commit after it is created; release rollback is to remove `v0.1.168-fy.1` and push `git push origin :refs/tags/v0.1.168-fy.1`.

## 2026-07-29 - Task: Verify fork release v0.1.168-fy.1

### What was done
- Confirmed the public GitHub Release for `v0.1.168-fy.1` contains the Linux amd64 package and checksum manifest generated from the merge tag.

### Testing
- Release asset returned HTTP `200`.
- `sub2api_0.1.168-fy.1_linux_amd64.tar.gz` is `36,053,635` bytes and contains the executable `sub2api` entry.
- SHA256 `724a533607c63b986ad507bbaab0149f7b20892450ee9d20a78f73b22ea55d3e` exactly matched `checksums.txt`.

### Notes
- `progress.md`: records the final public asset and checksum verification.
- `.superpowers/`: remains untracked and excluded from the commit and release.
- Rollback point: source rollback is `git revert aede017fb`; release rollback is to remove the GitHub Release and push `git push origin :refs/tags/v0.1.168-fy.1`.

## 2026-07-29 - Task: Add native multi-resolution pricing for LTX 2.3 video models

### What was done
- Added independent `1080p`, `1440p`, and `2160p` USD-per-second channel pricing for `ltxv-2.3-fast` and `ltxv-2.3-pro` while preserving the legacy resolution mapping for all other video models.
- Extended synchronous usage billing and async video-job reserve/settlement snapshots to retain and charge the exact LTX resolution tier. New async snapshots use version 3; existing snapshots remain readable and no database schema changed.
- Updated the admin channel form and operating documentation for the six requested prices: Fast `0.06/0.21/0.24` and Pro `0.09/0.18/0.36` for `1080p/1440p/2160p`.

### Testing
- Targeted backend LTX resolution, pricing resolver, channel validation, forwarding, usage, reserve, and settlement tests passed.
- `go test ./internal/service -count=1` passed in 99.262 seconds; `go test ./internal/handler -count=1` passed in 35.143 seconds.
- `go test ./... -count=1` passed for all backend packages.
- Targeted channel pricing Vitest passed all 14 tests; `pnpm.cmd typecheck`, `pnpm.cmd lint:check`, and `pnpm.cmd run build` passed.
- `git diff --check` passed with only the existing Markdown LF-to-CRLF warning.

### Notes
- `backend/internal/service/billing_service.go`: adds native 1440p and 2160p video unit prices and model-aware resolution selection.
- `backend/internal/service/billing_service_test.go`: verifies LTX Fast 2160p cost at the requested rate.
- `backend/internal/service/channel_service_test.go`: verifies a valid LTX three-tier channel pricing entry.
- `backend/internal/service/leo_video.go`: records the native LTX output resolution returned by Leo.
- `backend/internal/service/leo_video_test.go`: verifies 2160p LTX forwarding metadata and the 20-second Fast duration.
- `backend/internal/service/model_pricing_resolver.go`: extracts native LTX 1440p and 2160p channel prices.
- `backend/internal/service/model_pricing_resolver_test.go`: verifies all three native LTX prices are preserved.
- `backend/internal/service/openai_gateway_record_usage_test.go`: verifies 2160p LTX usage metadata and cost.
- `backend/internal/service/openai_gateway_usage.go`: applies model-aware LTX resolution normalization during synchronous usage billing.
- `backend/internal/service/video_billing_resolution.go`: defines native LTX pricing tiers while retaining legacy resolution compatibility for other models.
- `backend/internal/service/video_billing_resolution_test.go`: verifies LTX tier enumeration and resolution normalization.
- `backend/internal/service/video_job_billing.go`: stores and settles 1440p and 2160p prices in version 3 async billing snapshots.
- `backend/internal/service/video_job_billing_test.go`: verifies native LTX reserve and settlement amounts plus snapshot compatibility.
- `frontend/src/components/admin/channel/types.ts`: exposes the three native LTX tiers in the admin pricing form.
- `frontend/src/components/admin/channel/__tests__/types.spec.ts`: verifies creation and validation of LTX three-tier pricing entries.
- `docs/LEO_VIDEO_CHANNEL.md`: documents native LTX tiers, exact production prices, and snapshot behavior.
- `progress.md`: records implementation scope, verification evidence, modified files, and rollback point.
- `.superpowers/` remains untracked and excluded from this task.
- Rollback point: working-tree base is `d622c6053`; after commit, use `git revert <ltx_pricing_commit>` to restore the previous single-tier compatibility behavior.

## 2026-07-29 - Task: Deploy and production-test LTX 2.3 video models

### What was done
- Published commit `c9ec2cb43` as release `v0.1.168-fy.2` and deployed it to production.
- Added exact Leo account `1682` mappings for `ltxv-2.3-fast` and `ltxv-2.3-pro`, then added both models to channel `5` without changing its existing five pricing entries.
- Configured Fast prices at `$0.06/$0.21/$0.24` and Pro prices at `$0.09/$0.18/$0.36` per second for `1080p/1440p/2160p`.
- Submitted live Fast and Pro generation jobs plus a Seedance Fast control job. All three were blocked by the only Leo upstream video generation service; no output video was produced and the test API key was not charged.

### Testing
- Production version readback returned `0.1.168-fy.2`; `/v1/models` exposed both `ltxv-2.3-fast` and `ltxv-2.3-pro`.
- Account `1682` remained `active` and schedulable, and both exact model mappings read back correctly.
- Channel `5` contained seven pricing entries and all six LTX prices read back exactly as configured.
- Fast job `vidjob_D-ysl6LNHogwFON5NXVvsoZ8Oh11emo-` reached `running` and then failed with `Video service authentication failed`.
- Pro job `vidjob_dK7nvJyFgVqigpEvwevHYSURG5sormj8` and Seedance Fast control job `vidjob_deOCnTwJK5TOP1RvdmmnUoJCDxwtAcB-` failed with `Video service request failed`.
- Leo account health still returned `LeoStudio health check passed`, but group `25` had no second active Leo account for failover. Test API key `215` remained at `quota_used = 36.27` after the failures.
- Failed-job account concurrency drained from `3` to `0` by `2026-07-29 17:16:10 +08:00` under the production 30-minute slot TTL.
- Release package `sub2api_0.1.168-fy.2_linux_amd64.tar.gz` was `36,053,307` bytes with SHA256 `26172a3bd1dbc51a944af14a77af13b8b6c97336e4d7d41eb94566999c2b5b17`.

### Notes
- `progress.md`: records the production release, exact configuration, live generation outcomes, no-charge evidence, concurrency recovery, and rollback point.
- `.superpowers/`: remains untracked and excluded from the commit and release.
- Remaining blocker: repair or refresh the Leo upstream video credentials, or add another active Leo account, before repeating a successful-output test.
- Rollback point: run `git revert c9ec2cb43`, publish and deploy the resulting rollback release (or redeploy `v0.1.168-fy.1`), then remove the two LTX mappings from account `1682` and the two LTX pricing entries from channel `5`.

## 2026-07-29 - Task: Rename public LTX 2.3 model IDs

### What was done
- Renamed the public LTX model IDs from `ltxv-2.3-fast` and `ltxv-2.3-pro` to `ltx-2.3-fast` and `ltx-2.3-pro` across the API model list, video workbench, account defaults, channel pricing, and operating documentation.
- Preserved LeoStudio compatibility by keeping account mapping targets as `ltxv-2.3-fast` and `ltxv-2.3-pro`; public requests use the new names while the upstream service still receives its required model IDs.
- Kept internal normalization for existing task and billing records that contain the former upstream IDs. No database schema changed.

### Testing
- Targeted backend Leo account, model capability, pricing, usage, reserve, and settlement tests passed.
- `go test ./... -count=1`: passed for all backend packages.
- Targeted frontend account, channel pricing, API documentation, and video workbench suites passed: 4 files and 57 tests.
- `pnpm.cmd test:run`, `pnpm.cmd typecheck`, and `pnpm.cmd lint:check`: passed.
- `pnpm.cmd build`: passed; Vite transformed 994 modules and emitted the embedded frontend bundle. Existing dynamic-import and chunk-size warnings remain non-blocking.
- `git diff --check`: passed with only the existing Markdown LF-to-CRLF warnings.

### Notes
- `backend/internal/service/leo_account.go`: defines the public LTX IDs, upstream aliases, and compatibility normalization.
- `backend/internal/service/leo_account_test.go`: verifies the new public defaults and old upstream alias compatibility.
- `backend/internal/service/leo_video_model_specs.go`: keys LTX capabilities by the new public IDs and normalizes aliases before lookup.
- `backend/internal/service/leo_video_model_specs_test.go`: verifies public model validation and upstream alias lookup.
- `backend/internal/service/video_billing_resolution.go`: recognizes both public and upstream LTX IDs for native resolution billing.
- `backend/internal/service/video_billing_resolution_test.go`: covers pricing tiers and normalization for both naming forms.
- `backend/internal/service/billing_service_test.go`: updates LTX billing coverage to the public model ID.
- `backend/internal/service/channel_service_test.go`: updates LTX channel pricing validation coverage to the public model ID.
- `backend/internal/service/leo_video_test.go`: updates LTX forwarding coverage to the public model ID.
- `backend/internal/service/model_pricing_resolver_test.go`: updates native-tier resolution coverage to the public model ID.
- `backend/internal/service/openai_gateway_record_usage_test.go`: updates LTX usage and charge coverage to the public model ID.
- `backend/internal/service/video_job_billing_test.go`: updates LTX reserve and settlement coverage to the public model ID.
- `frontend/src/components/account/CreateAccountModal.vue`: maps new public LTX IDs to the unchanged LeoStudio upstream IDs by default.
- `frontend/src/components/account/__tests__/CreateAccountModal.spec.ts`: verifies both new default mappings.
- `frontend/src/components/admin/channel/types.ts`: recognizes new public LTX IDs as native three-tier video models.
- `frontend/src/components/admin/channel/__tests__/types.spec.ts`: verifies new public LTX channel pricing entries.
- `frontend/src/composables/useModelWhitelist.ts`: exposes the new LTX IDs in model selection.
- `frontend/src/constants/channel.ts`: lists the new public LTX IDs in channel pricing order.
- `frontend/src/i18n/locales/en/dashboard.ts`: updates English LTX model labels and examples.
- `frontend/src/i18n/locales/zh/dashboard.ts`: updates Chinese LTX model labels and examples.
- `frontend/src/views/user/VideoApiDocsView.vue`: updates public API documentation and examples to the new IDs.
- `frontend/src/views/user/VideoGenerationView.vue`: updates workbench model values and capability checks to the new IDs.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: verifies the new IDs in the API documentation.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: verifies the new IDs in the workbench.
- `docs/LEO_VIDEO_CHANNEL.md`: documents the public-to-upstream mapping and new pricing model names.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: documents the new public request model IDs.
- `progress.md`: records the implementation, verification evidence, changed files, and rollback point.
- `.superpowers/` remains untracked and excluded from this task.
- Rollback point: working-tree base is `ab6b0e7d1`; after commit, use `git revert <ltx_public_id_rename_commit>` to restore the former public IDs. Existing production account, channel, and announcement configuration must be restored separately if already migrated.

## 2026-07-29 - Task: Release and deploy public LTX model ID rename

### What was done
- Pushed commit `fda87c078` and annotated tag `v0.1.168-fy.3`, then deployed the checksum-verified release to `api.fflink.top` through the built-in updater and restarted the service.
- Renamed the two public mappings on Leo account `1682` to `ltx-2.3-fast -> ltxv-2.3-fast` and `ltx-2.3-pro -> ltxv-2.3-pro`, preserving the other five mappings and the upstream model IDs.
- Renamed only the two LTX model identifiers in channel `5`; all seven pricing entries and the six LTX prices remained unchanged.
- Updated active popup announcement `21` with the new public model IDs using an explicit UTF-8 request body.

### Testing
- Release asset `sub2api_0.1.168-fy.3_linux_amd64.tar.gz` is `36,053,752` bytes and contains the `sub2api` executable.
- SHA256 `f94a7a9dc126bf6e1aa7e933ff1896bdcab6d452893cfa8b00ae9832b811729d` exactly matched the published `checksums.txt`.
- Production `/health` returned `status=ok`, and the admin version endpoint returned `0.1.168-fy.3` after restart.
- `/v1/models` returned `ltx-2.3-fast` and `ltx-2.3-pro` and did not return either former `ltxv-*` public ID.
- Account `1682` readback remained `active` and schedulable with both new-public-to-old-upstream mappings and no old public mapping keys.
- Channel `5` readback retained seven entries. Fast prices remained `0.06/0.21/0.24` and Pro prices remained `0.09/0.18/0.36` for `1080p/1440p/2160p`.
- Announcement `21` readback remained `active` with `notify_mode=popup`, contained both new public IDs, contained no former `ltxv-*` ID, and preserved valid Chinese text.

### Notes
- `progress.md`: records the release, production configuration changes, verification evidence, and rollback procedure.
- `.superpowers/` remains untracked and excluded from the release and follow-up commit.
- Production rollback: while `0.1.168-fy.3` is running, rename account `1682` mapping keys and channel `5` pricing model IDs back to `ltxv-2.3-fast/pro`, restore those names in announcement `21`, then call `POST /api/v1/admin/system/rollback` with `{"version":"0.1.168-fy.2"}` and restart the service. Source rollback is `git revert fda87c078` followed by a new release.

## 2026-07-30 - Task: Add scheduled email alerts for new account errors

### What was done
- Added a configurable scheduled scan that sends an Ops email only when an account newly enters the error state or its error reason changes.
- Persisted the current error baseline so unchanged errors are not repeated across scans or service restarts; the first scan establishes the baseline without sending historical errors.
- Added the admin switch and Cron schedule field, defaulting to disabled with a five-minute schedule, and documented configuration and detection behavior.

### Testing
- `go test ./internal/service -count=1` passed (`98.533s`).
- Targeted account error reminder, HTML escaping, default schedule, and config persistence tests passed.
- `frontend/node_modules/.bin/vue-tsc.CMD --noEmit` passed.
- ESLint passed for the changed Ops API, notification card, and English/Chinese locale files.
- Vitest locale compilation and key-collision suites passed: 2 files, 8 tests.
- `git diff --check` passed before the progress entry was appended; a final diff check is required after this entry.

### Notes
- `backend/internal/service/domain_constants.go`: adds the internal settings key for the persisted account error baseline.
- `backend/internal/service/ops_settings_models.go`: adds account error alert enablement and schedule fields to the Ops email report config.
- `backend/internal/service/ops_settings.go`: saves, normalizes, and defaults the new alert configuration.
- `backend/internal/service/ops_scheduled_report_service.go`: schedules scans, detects new errors, sends detail emails, and persists the baseline after delivery.
- `backend/internal/service/ops_scheduled_report_service_test.go`: verifies baseline, deduplication, changed reasons, recovery, escaping, defaults, and config persistence.
- `frontend/src/api/admin/ops.ts`: exposes the new report configuration fields to the admin client.
- `frontend/src/views/admin/ops/components/OpsEmailNotificationCard.vue`: adds validation and controls for the new alert schedule.
- `frontend/src/i18n/locales/zh/admin/ops.ts`: adds the Chinese field label.
- `frontend/src/i18n/locales/en/admin/ops.ts`: adds the English field label.
- `docs/ACCOUNT_ERROR_EMAIL_ALERTS.md`: documents setup, defaults, detection semantics, and retry behavior.
- `progress.md`: records implementation, verification evidence, changed files, and rollback instructions.
- Before commit, roll back this task with `git restore -- backend/internal/service/domain_constants.go backend/internal/service/ops_settings_models.go backend/internal/service/ops_settings.go backend/internal/service/ops_scheduled_report_service.go backend/internal/service/ops_scheduled_report_service_test.go frontend/src/api/admin/ops.ts frontend/src/views/admin/ops/components/OpsEmailNotificationCard.vue frontend/src/i18n/locales/zh/admin/ops.ts frontend/src/i18n/locales/en/admin/ops.ts progress.md` and `Remove-Item -LiteralPath docs/ACCOUNT_ERROR_EMAIL_ALERTS.md`; after commit, use `git revert <account_error_email_alert_commit>`.
- `.gitignore`: whitelists only the new account error alert document under the otherwise ignored `docs/` directory.
- Rollback addendum: include `.gitignore` in the `git restore -- ...` command above so the document whitelist is reverted with the feature.
- Final `git diff --check` completed after the progress entry and passed; only line-ending conversion warnings were reported.

## 2026-07-30 - Task: Release account error email alerts

### What was done
- Published the account error email alert feature commit `9e503395f1` on `codex/leo-video-channel` and created annotated tag `v0.1.168-fy.4` without including the unrelated local OpenAI, account usage, or `.superpowers/` changes.
- Confirmed the GitHub Release was generated as the latest release and points to commit `9e503395f1`.
- Downloaded and verified the Linux amd64 release package and confirmed the archive contains the `sub2api` executable.

### Testing
- On a clean detached worktree containing only the staged feature patch, `go test ./internal/service -count=1` passed in `104.455s`.
- On the same clean candidate, `pnpm.cmd test:run`, `pnpm.cmd typecheck`, `pnpm.cmd lint:check`, and `pnpm.cmd run build` all passed; the production build completed in `42.46s`.
- `git diff --cached --check` passed before the feature commit.
- Release package `sub2api_0.1.168-fy.4_linux_amd64.tar.gz` is `36,057,551` bytes with SHA256 `5e80991f2634b0dc6a147b38e2da49d164d5dc3bc25f5a9a59e553902bc3abf9`, exactly matching `checksums.txt`.
- `tar -tzf` listed the expected `sub2api` executable.

### Notes
- `progress.md`: records the isolated release scope, clean-candidate verification, GitHub Release result, package checksum, and rollback point.
- The published tag remains on feature commit `9e503395f1`; this release-record commit is intentionally branch-only.
- Unrelated local OpenAI, account usage, and `.superpowers/` changes remain uncommitted and were excluded from the feature commit, tag, and release package.
- Rollback point: run `git revert 9e503395f1`, publish the resulting rollback release, or redeploy `v0.1.168-fy.3` if an immediate binary rollback is required.

## 2026-07-30 - Task: Rebuild OpenAI Codex reverse-proxy identity fingerprinting

### What was done
- Added a stable OpenAI OAuth device fingerprint domain that preserves official inbound installation identities, honors an operator-managed `openai_device_id`, and otherwise derives a deterministic account device without using tokens or proxy addresses.
- Added controlled device generations through `openai_device_profile_id`; token refreshes and normal proxy rotation keep the same device, while an intentional profile-generation change rotates managed identifiers together.
- Extended session, conversation, and prompt-cache isolation from downstream API Key only to upstream account plus downstream API Key, preventing identifiers from crossing scheduled accounts.
- Applied the same device/session behavior to Responses HTTP, passthrough, Messages and Chat bridges, WebSocket v2 and ingress, alpha/search, Live, account tests, compact probes, and usage probes. OpenAI API-key upstream behavior remains unchanged.
- Kept the existing OpenAI HTTP/2 transport and did not apply the Anthropic Node.js TLS profile to Codex, because that would contradict the declared `codex_cli_rs` application identity.

### Testing
- `go test ./internal/service -run 'Test(OpenAICodex|ApplyOpenAICodex|IsolateOpenAIAccountSessionID|EnsureCodexIdentityHeaders|EnforceCodexIdentityHeaders)' -count=1`: passed.
- HTTP/WS/compatibility/compact/alpha regression group: passed (`27.460s`).
- Live/account/probe/session regression group: passed (`2.441s`).
- Session/prompt-cache/Codex identity/Chat Completions regression group: passed (`8.312s`).
- `go test ./internal/service -count=1`: passed (`98.747s`).
- `git diff --check`: passed before and after documentation/progress updates.

### Notes
- `.gitignore`: allows the Codex fingerprint document under the repository's ignored-by-default `docs/` directory.
- `backend/internal/service/openai_codex_identity.go`: defines device resolution precedence, deterministic account/window mapping, and body/header synchronization.
- `backend/internal/service/openai_codex_identity_test.go`: verifies official identity preservation, managed identity stability, token/proxy independence, account rotation, and body/header alignment.
- `backend/internal/service/openai_gateway_service.go`: adds account-aware session and prompt-cache isolation while retaining the existing hash format.
- `backend/internal/service/openai_gateway_service_session_isolation_test.go`: verifies account and downstream-tenant isolation.
- `backend/internal/service/openai_gateway_forward.go`: applies the unified fingerprint at the final Responses HTTP outbound boundary.
- `backend/internal/service/openai_gateway_passthrough.go`: applies the same body, session, installation, and window mapping to passthrough requests.
- `backend/internal/service/openai_gateway_chat_completions.go`: uses the selected upstream account when deriving bridged session IDs.
- `backend/internal/service/openai_gateway_messages.go`: uses the selected upstream account when deriving Messages bridge session IDs.
- `backend/internal/service/openai_alpha_search.go`: adds account-aware sessions and stable synthetic device/window identity to alpha/search.
- `backend/internal/service/openai_live.go`: gives Live calls stable account devices and per-call session/window lifetimes.
- `backend/internal/service/openai_ws_forwarder_payload.go`: aligns OAuth WebSocket handshake sessions and device headers with HTTP forwarding.
- `backend/internal/service/openai_ws_forwarder_v2.go`: maps WebSocket v2 prompt-cache and client metadata into the account device domain.
- `backend/internal/service/openai_ws_forwarder_ingress.go`: applies the same mapping to each WebSocket ingress turn.
- `backend/internal/service/account_test_service.go`: makes account and compact probes use the account's stable synthetic fingerprint.
- `backend/internal/service/account_usage_service.go`: makes Codex usage probes use the account's stable synthetic fingerprint.
- `backend/internal/service/openai_agent_identity_compat_test.go`: updates Agent Identity expectations for upstream-account isolation and mapped prompt-cache keys.
- `backend/internal/service/openai_compat_model_test.go`: updates OAuth Messages continuation expectations for account-aware sessions.
- `backend/internal/service/openai_ws_forwarder_ingress_session_test.go`: verifies account-aware ingress session headers.
- `backend/internal/service/openai_ws_forwarder_success_test.go`: verifies account-aware OAuth WebSocket session headers.
- `docs/OPENAI_CODEX_FINGERPRINT.md`: documents identity precedence, lifetimes, covered paths, transport boundaries, operation, and verification.
- `progress.md`: records implementation, verification evidence, changed files, and rollback instructions.
- No database schema, authentication credential format, proxy binding, or OpenAI TLS transport was changed.
- Rollback point: `9e503395f1eed0fa0fb43fce6bd42563eb652b82`. Before commit, restore the listed source/test files plus `.gitignore` and `progress.md` with `git restore -- <paths>` and remove `docs/OPENAI_CODEX_FINGERPRINT.md`; after a dedicated commit, use `git revert <codex_fingerprint_commit>`.

## 2026-07-30 - Task: Document video-group integration for other Sub2API instances

### What was done
- Added a Chinese secondary-development guide that defines the supported Sub2API-to-Sub2API video-group integration boundary, including dedicated API Key binding, configuration-only synchronous cascading, source-porting scope, and the contract changes required for full asynchronous cascading.
- Documented the current incompatibilities instead of presenting account health success as end-to-end compatibility: public string job IDs, RFC 3339 timestamps, authenticated relative content URLs, and cross-host media uploads all require explicit handling.

### Testing
- Cross-checked the guide against the current Leo route dispatch, account Base URL validation, asynchronous client DTOs, `video_jobs` schema, task runtime, protected content handler, and video output downloader.
- Verified all documented endpoint paths, public model mappings, permission behavior, task states, and migration references with repository source and existing video documents.
- `git diff --check` passed after this progress entry; only line-ending conversion warnings were reported.

### Notes
- `.gitignore`: allows the new tracked video-group secondary-development guide under the repository's ignored-by-default `docs/` directory.
- `docs/SUB2API_VIDEO_GROUP_SECONDARY_DEVELOPMENT_CN.md`: documents the integration model, synchronous setup, required source capabilities, asynchronous contract changes, security boundaries, acceptance checks, and rollback constraints.
- `progress.md`: records the documentation scope, evidence, changed files, and rollback method.
- No application code, API behavior, database schema, authentication path, pricing logic, or deployment configuration was changed.
- Before commit, roll back this task with `git restore -- .gitignore progress.md` and `Remove-Item -LiteralPath docs/SUB2API_VIDEO_GROUP_SECONDARY_DEVELOPMENT_CN.md`; after a dedicated commit, use `git revert <video_group_secondary_development_doc_commit>`.

## 2026-07-30 - Task: Redact internal provider identifiers from video integration guide

### What was done
- Reworked the public-facing video integration guide to use neutral terms such as “视频平台”“上游视频服务”和“视频上游账号”。
- Removed the internal provider name, provider-specific aliases, platform values, implementation field names, provider-named source filenames, and internal migration filenames from the guide while retaining the public API contract and integration steps.

### Testing
- Confirmed the guide contains no case-insensitive match for the removed provider identifier or its service name.
- Confirmed the guide no longer contains provider-specific model aliases or internal group/account/task field names.
- `git diff --check` passed after the redaction update; only existing line-ending conversion warnings were reported.

### Notes
- `docs/SUB2API_VIDEO_GROUP_SECONDARY_DEVELOPMENT_CN.md`: redacts internal provider and implementation identifiers from the externally shareable guide.
- `progress.md`: records the redaction scope, verification evidence, changed files, and rollback method.
- No application code, API behavior, database schema, authentication path, pricing logic, or deployment configuration was changed.
- Before commit, roll back this redaction with `git restore -- progress.md` and restore the previous version of `docs/SUB2API_VIDEO_GROUP_SECONDARY_DEVELOPMENT_CN.md`; after a dedicated commit, use `git revert <video_group_guide_redaction_commit>`.

## 2026-07-30 - Task: Remove provider-named model examples from video integration guide

### What was done
- Replaced the remaining concrete provider-named model mapping example with a neutral rule: use the upstream `/v1/models` result and perform any alias conversion only at the final hop.

### Testing
- Re-ran the redaction scan; the guide contains no target provider name, provider service name, provider alias, or internal group/account/task field name.
- `git diff --check` passed after the final guide edit; only existing line-ending conversion warnings were reported.

### Notes
- `docs/SUB2API_VIDEO_GROUP_SECONDARY_DEVELOPMENT_CN.md`: removes the final provider-specific model example from the externally shareable guide.
- `progress.md`: records the final redaction adjustment and verification evidence.
- No application code or runtime behavior was changed.
- Before commit, roll back this final adjustment with `git restore -- progress.md` and restore the previous version of `docs/SUB2API_VIDEO_GROUP_SECONDARY_DEVELOPMENT_CN.md`; after a dedicated commit, use `git revert <video_group_model_example_redaction_commit>`.

## 2026-07-30 - Task: Rename video integration guide

### What was done
- Renamed the externally shareable guide to `docs/对接二开文档.md` as requested.
- Updated the `docs/` ignore whitelist to track the new filename; document content and runtime behavior are unchanged.

### Testing
- Confirmed the new file exists with the same 14,362-byte content size as before the rename.
- `git diff --check` passed after the rename and whitelist update; only existing line-ending conversion warnings were reported.

### Notes
- `.gitignore`: replaces the old guide whitelist entry with `!docs/对接二开文档.md`.
- `docs/对接二开文档.md`: renamed from the previous English-named guide; content unchanged.
- `progress.md`: records the rename, verification evidence, and rollback method.
- Before commit, roll back this rename with `Move-Item -LiteralPath 'docs\对接二开文档.md' -Destination 'docs\SUB2API_VIDEO_GROUP_SECONDARY_DEVELOPMENT_CN.md'`, then restore the previous `.gitignore` whitelist line and `progress.md`; after a dedicated commit, use `git revert <video_integration_guide_rename_commit>`.

## 2026-07-30 - Task: Split legacy and new OpenAI Codex fingerprint modes

### What was done
- Kept pre-existing OpenAI OAuth accounts on the legacy device, window, session, conversation, and prompt-cache behavior through a runtime-only compatibility marker applied when unversioned accounts are loaded.
- Made newly created OpenAI OAuth accounts default to the `v1` fingerprint mode, while allowing an explicit per-account `legacy` or `v1` override in `accounts.extra`.
- Kept the split limited to OpenAI OAuth traffic; API-key accounts and other platforms remain unchanged.

### Testing
- `gofmt` passed for all modified service, repository, and test files.
- `git diff --check` passed; only existing line-ending conversion warnings were reported.
- Targeted Go tests could not compile because pre-existing unrelated edits in `backend/internal/service/billing_service.go` and `backend/internal/service/video_job_billing.go` contain top-level statements outside functions. Those files were not changed in this task.

### Notes
- `backend/internal/service/account.go`: adds the runtime-only legacy compatibility marker.
- `backend/internal/service/openai_codex_identity.go`: gates the new fingerprint contract by account mode and exposes repository marking for unversioned legacy accounts.
- `backend/internal/service/openai_gateway_service.go`: falls back to the previous API-key/session namespace for legacy accounts.
- `backend/internal/service/admin_account.go`: defaults new OpenAI OAuth accounts to `v1` and propagates the mode to new Spark shadows.
- `backend/internal/service/account_service.go`: defaults the alternate account creation service to `v1` for new OpenAI OAuth accounts.
- `backend/internal/repository/account_repo.go`: marks unversioned persisted OpenAI OAuth accounts as legacy in memory without writing the marker back during reads.
- `backend/internal/service/openai_codex_identity_test.go`: covers legacy no-op behavior alongside v1 identity behavior.
- `backend/internal/service/openai_codex_fingerprint_mode_test.go`: covers new-account defaults and explicit legacy overrides.
- `backend/internal/service/openai_gateway_service_session_isolation_test.go`: opts the account-isolation fixture into v1 mode.
- `docs/OPENAI_CODEX_FINGERPRINT.md`: documents the staged legacy/v1 rollout and per-account override.
- `progress.md`: records this implementation and verification status.
- Rollback point: restore the files listed above and the documentation change with `git restore -- <paths>` before commit; after a dedicated commit, use `git revert <codex_fingerprint_mode_commit>`. This rollback removes the mode split but does not alter account data.

## 2026-07-30 - Task: Verify OpenAI Codex fingerprint mode split

### What was done
- Corrected the mode coverage test to pass the normalized `Extra` map used by the administrator account creation path.
- Verified that runtime legacy marking only affects unversioned OpenAI OAuth accounts, preserves explicit `v1`, and ignores OpenAI API-key accounts.

### Testing
- `gofmt -w` passed for the fingerprint-related Go files.
- Targeted fingerprint mode and Codex identity tests passed.
- `go test ./internal/service -count=1` passed (`98.072s`).
- `git diff --check` passed; only existing line-ending conversion warnings were reported.

### Notes
- `backend/internal/service/openai_codex_fingerprint_mode_test.go`: fixes the creation-path fixture and adds legacy-marker boundary assertions.
- `progress.md`: records the final verification result and corrected test coverage.
- Correction to the previous task entry: `backend/internal/service/account.go` was not changed; the runtime compatibility marker is implemented in `backend/internal/service/openai_codex_identity.go`.
- Rollback point: restore `backend/internal/service/openai_codex_fingerprint_mode_test.go` and `progress.md` before commit; after a dedicated commit, use `git revert <codex_fingerprint_mode_verification_commit>`.

## 2026-07-30 - Task: Align video pricing with native resolutions

### What was done
- Updated Happy Horse pricing to use only `720p` and `1080p`, and Grok Imagine 1.5 pricing to use only `400p`, `544p`, `720p`, and `960p`; legacy three-tier entries remain readable during migration but are no longer accepted for new channel saves.
- Added native resolution prices to synchronous and asynchronous video billing snapshots, preserving LTX and Seedance model-specific tiers.
- Updated the admin pricing editor and the Leo video channel guide so no compatibility or unsupported resolution is presented.
- Committed as `26c82af5b` and pushed branch `codex/leo-video-channel`; created and re-created annotated tag `v0.1.168-fy.5` at the same commit after the first release trigger produced no visible Release.

### Testing
- Targeted native-resolution backend billing, channel validation, resolver, async reserve/settlement, and legacy snapshot compatibility tests passed.
- `go test -p 2 -tags=unit -timeout 10m ./... -count=1` passed in a clean detached worktree (`317.4s`).
- Frontend channel pricing tests passed (`16` tests); `pnpm.cmd typecheck`, `pnpm.cmd lint:check`, and `pnpm.cmd build` passed.
- A clean-tag Release reproduction passed frontend dependency installation/build and Linux amd64 `go build -tags embed`; the binary was produced successfully.
- `git diff --cached --check` and `git diff --check` passed; the dirty worktree service-wide test failure is isolated to the unrelated OpenAI fingerprint test `TestBuildAccountForCreatePreservesExplicitLegacyFingerprintMode`.
- Production readback before migration: `/health` returned `ok`, version `0.1.168-fy.3`, and channel `5` contained seven entries with legacy Happy Horse/Grok tiers.
- Production updater checks still return fork version `0.1.168-fy.3` and `GitHub API returned 404`; no production binary restart or pricing mutation was performed.

### Notes
- `backend/internal/service/billing_service.go`: adds native per-second video price slots.
- `backend/internal/service/billing_service_test.go`: verifies native Happy Horse and Grok direct costs.
- `backend/internal/service/channel_service_test.go`: validates exact native tier sets and rejects compatibility tiers.
- `backend/internal/service/model_pricing_resolver.go`: resolves native prices and reads legacy entries during migration.
- `backend/internal/service/model_pricing_resolver_test.go`: covers native and legacy resolver behavior.
- `backend/internal/service/video_billing_resolution.go`: defines native model resolution sets and normalization.
- `backend/internal/service/video_billing_resolution_test.go`: covers native resolution support and normalization.
- `backend/internal/service/video_job_billing.go`: snapshots native prices and settles native result tiers.
- `backend/internal/service/video_job_billing_test.go`: covers native reserve/settlement and v3 compatibility snapshots.
- `frontend/src/components/admin/channel/types.ts`: generates only model-supported pricing tiers.
- `frontend/src/components/admin/channel/__tests__/types.spec.ts`: verifies Happy Horse and Grok editor tiers.
- `docs/LEO_VIDEO_CHANNEL.md`: documents native pricing tiers and removes the stale three-tier editor instruction.
- `progress.md`: records this implementation, verification, release attempt, and production blocker.
- Rollback point: source rollback is `git revert 26c82af5b`; release rollback is to remove/recreate `v0.1.168-fy.5` only if necessary and keep production at `0.1.168-fy.3` until a verified Release asset is available. The untracked temporary reproduction directory is `C:\Users\feiyu\AppData\Local\Temp\sub2api-fy5-release-repro` and contains only generated build files.

## 2026-07-30 - Task: Isolate inbound Codex installation IDs and stop deterministic model failover

### What was done
- `v1` now maps inbound `X-Codex-Installation-ID` and body `client_metadata.x-codex-installation-id` with the selected upstream ChatGPT account identity. The mapping is stable for one account, differs across accounts, and is aligned between headers and body.
- `legacy` mode keeps the previous pass-through behavior for compatibility.
- Codex `response.failed` events containing `model is not supported when using Codex with a ChatGPT account` are now treated as deterministic account/model errors and do not fan out across the OAuth pool. Capacity and transient processing errors retain bounded failover behavior.
- Updated `docs/OPENAI_CODEX_FINGERPRINT.md` with the mapping and failover boundary.

### Testing
- `gofmt` passed for the changed Go files.
- Targeted identity, session-isolation, stream-failure, transient-failover, and Codex plan-gating tests passed.
- Full `go test ./internal/service -count=1` is the remaining verification step for this change.

### Notes
- Changed files: `backend/internal/service/openai_codex_identity.go`, `backend/internal/service/openai_codex_identity_test.go`, `backend/internal/service/openai_gateway_passthrough.go`, `backend/internal/service/openai_codex_failover_policy_test.go`, and `docs/OPENAI_CODEX_FINGERPRINT.md`.
- Rollback before commit: `git restore -- <the five paths above> progress.md`; after a dedicated commit, use `git revert <codex_identity_failover_commit>`.

## 2026-07-31 - Task: Release staged OpenAI Codex fingerprint modes

### What was done
- Published commit `49f845e6b` on `codex/leo-video-channel` and created annotated tag `v0.1.168-fy.6`.
- Included the staged OpenAI Codex fingerprint compatibility changes, account creation defaults, runtime legacy marking, related tests, and documentation; excluded `.superpowers/` temporary files.
- Confirmed the GitHub Release page and Linux amd64 asset were generated for `v0.1.168-fy.6`.

### Testing
- On a clean candidate worktree containing only the staged patch, `gofmt` passed for 23 Go files and `go test ./internal/service -count=1` passed in `103.497s`.
- On the same candidate, `pnpm.cmd test:run`, `pnpm.cmd typecheck`, `pnpm.cmd lint:check`, and `pnpm.cmd run build` all passed; the production build completed in `28.99s`.
- `git diff --cached --check` passed before commit.
- Release package `sub2api_0.1.168-fy.6_linux_amd64.tar.gz` is `36,069,030` bytes with SHA256 `71fddf825aefe80311788b09f2bcd1841446a2d84f9e995a3eeec147d1ce6119`, exactly matching `checksums.txt`.
- `tar -tzf` listed only the expected `sub2api` executable.

### Notes
- `progress.md`: records the release commit, isolated verification, package checksum, and rollback point.
- The published tag remains on feature commit `49f845e6b`; the later release-record commit is branch-only.
- `.superpowers/` remains untracked and excluded from the commit, tag, and Release package.
- Rollback point: run `git revert 49f845e6b`, publish the resulting rollback release, or redeploy `v0.1.168-fy.5` if an immediate binary rollback is required.

## 2026-07-31 - Task: Verify Codex installation isolation and failover policy

### What was done
- Completed the release gate for account-scoped inbound Codex installation IDs and deterministic plan/model failure handling.
- Verified that v1 mode isolates the same inbound installation across upstream accounts, legacy mode preserves existing identifiers, and plan-gated model failures do not fan out across the OAuth pool.

### Testing
- `gofmt -l` reported no formatting issues across the four changed Go files.
- Targeted Codex identity and failover tests passed with `go test ./internal/service -run 'Test(EnsureCodexIdentityHeaders|EnforceCodexIdentityHeaders|OpenAICodex|ApplyOpenAICodex|OpenAIStreamFailed)' -count=1` (`5.986s`).
- `go test -p 2 -tags=unit -timeout 10m ./... -count=1` passed across the complete backend unit suite; the service package completed in `172.726s`.
- `git diff --cached --check` passed before the verification candidate was created.

### Notes
- `backend/internal/service/openai_codex_identity.go`: scopes inbound installation IDs to the selected upstream OAuth account in v1 mode.
- `backend/internal/service/openai_codex_identity_test.go`: verifies account isolation, legacy preservation, and header/body alignment.
- `backend/internal/service/openai_gateway_passthrough.go`: prevents deterministic Codex plan/model failures from triggering pool-wide failover.
- `backend/internal/service/openai_codex_failover_policy_test.go`: covers deterministic and transient stream/HTTP failure classification.
- `docs/OPENAI_CODEX_FINGERPRINT.md`: documents installation mapping and the failover boundary.
- `progress.md`: records the completed verification evidence and rollback point.
- Before commit, roll back with `git restore -- backend/internal/service/openai_codex_identity.go backend/internal/service/openai_codex_identity_test.go backend/internal/service/openai_gateway_passthrough.go docs/OPENAI_CODEX_FINGERPRINT.md progress.md` and remove `backend/internal/service/openai_codex_failover_policy_test.go`; after commit, use `git revert <codex_identity_failover_commit>`.

## 2026-07-31 - Task: Release Codex installation isolation and failover policy

### What was done
- Published commit `004b9fa3f` on `codex/leo-video-channel` and created annotated tag `v0.1.168-fy.7`.
- Confirmed the GitHub Release page and Linux amd64 asset were generated for `v0.1.168-fy.7`.
- Kept `.superpowers/` temporary files outside the commit, tag, and Release package.

### Testing
- Clean-candidate targeted Codex identity and failover tests passed in `5.986s`.
- Clean-candidate `go test -p 2 -tags=unit -timeout 10m ./... -count=1` passed across the complete backend unit suite.
- Release package `sub2api_0.1.168-fy.7_linux_amd64.tar.gz` is `36,067,125` bytes with SHA256 `fc2ad9f65644b5a7261ee544816a36f40bd3245ac1696d48f6f66165a1343c07`, exactly matching `checksums.txt`.
- `tar -tzf` listed only the expected `sub2api` executable.

### Notes
- `progress.md`: records the release commit, package verification, excluded temporary files, and rollback point.
- The published tag remains on feature commit `004b9fa3f`; the later release-record commit is branch-only.
- Rollback point: run `git revert 004b9fa3f`, publish the resulting rollback release, or redeploy `v0.1.168-fy.6` if an immediate binary rollback is required.

## 2026-07-31 - Task: Immediately fail over OpenAI model-capacity errors

### What was done
- Changed the explicit `Selected model is at capacity` response from same-account retry to immediate failover to another eligible account.
- Preserved the existing bounded same-account retry policy for other transient processing errors and left account selection, authentication, and concurrency behavior unchanged.
- Added an end-to-end two-account regression that verifies the failed account is excluded before the request is replayed.

### Testing
- Targeted OpenAI service capacity and stream failover tests passed with `go test ./internal/service -run 'Test(OpenAIGatewayService_Forward_ModelCapacityErrorTriggersImmediateAccountFailover|OpenAIStreamingResponseFailedBeforeOutputCapacityErrorReturnsFailover|OpenAIStreamFailedTransientProcessingStillFailsOver)$' -count=1`.
- Targeted handler account-switch tests passed with `go test -tags=unit ./internal/handler -run 'TestOpenAIGatewayHandlerResponses_(ModelCapacityImmediatelySwitchesAccount|FailoverContinuesForConnectedClient)$' -count=1`.
- Package-level regression passed with `go test -p 2 -tags=unit -timeout 10m ./internal/service ./internal/handler -count=1`; service completed in `161.630s` and handler completed in `29.525s`.
- `git diff --check` passed; only line-ending conversion warnings for `docs/OPENAI_CODEX_FINGERPRINT.md` and `progress.md` were reported.

### Notes
- `backend/internal/service/openai_gateway_upstream_errors.go`: identifies the explicit model-capacity response and prevents same-account replay while retaining next-account failover.
- `backend/internal/service/openai_gateway_service_codex_cli_only_test.go`: verifies HTTP capacity failures request immediate account failover.
- `backend/internal/handler/openai_responses_failover_cancel_test.go`: verifies a pool-mode capacity failure calls account 1 once and then account 2.
- `docs/OPENAI_CODEX_FINGERPRINT.md`: documents immediate next-account failover for the explicit capacity response.
- `progress.md`: records the implementation, verification evidence, changed files, and rollback point.
- Before commit, roll back with `git restore -- backend/internal/service/openai_gateway_upstream_errors.go backend/internal/service/openai_gateway_service_codex_cli_only_test.go backend/internal/handler/openai_responses_failover_cancel_test.go docs/OPENAI_CODEX_FINGERPRINT.md progress.md`; after a dedicated commit, use `git revert <openai_capacity_failover_commit>`.

## 2026-07-31 - Task: Release immediate OpenAI capacity failover

### What was done
- Published commit `2c3c50c05` on `codex/leo-video-channel` and created annotated tag `v0.1.168-fy.8`.
- Confirmed the GitHub Release page and Linux amd64 asset were generated for `v0.1.168-fy.8`.
- Kept `.superpowers/` temporary files outside the commit, tag, and Release package.

### Testing
- Clean-candidate service and handler targeted tests passed in `5.994s` and `5.969s`.
- Clean-candidate package regression passed with service in `163.042s` and handler in `29.323s`.
- Release package `sub2api_0.1.168-fy.8_linux_amd64.tar.gz` is `36,066,718` bytes with SHA256 `937e92e08a2fd6d711a0435d2f5a61ca456f372889e5c018cc796d45139a9de7`, exactly matching `checksums.txt`.
- `tar -tzf` listed only the expected `sub2api` executable.

### Notes
- `progress.md`: records the release commit, package verification, excluded temporary files, and rollback point.
- The published tag remains on feature commit `2c3c50c05`; the later release-record commit is branch-only.
- Rollback point: run `git revert 2c3c50c05`, publish the resulting rollback release, or redeploy `v0.1.168-fy.7` if an immediate binary rollback is required.

## 2026-08-03 - Task: Hide Leo provider identity from public clients

### What was done
- Added a public `video` platform alias at HTTP response boundaries while retaining the internal `leo` identifier for scheduling, billing, routing, and admin operations.
- Replaced public video UI labels and empty-state copy with neutral `Video` wording, while keeping compatibility with legacy `leo` values in browser state.
- Expanded synchronous and asynchronous video error sanitization to cover bare `leo` as well as provider-branded variants.
- Documented the public/internal platform boundary and added regression coverage for public DTOs, channels, model plaza, quota output, and provider error sanitization.

### Testing
- `go test ./internal/service -run 'Test(PublicVideoErrorMessageHidesProviderNames|PublicPlatformIDHidesInternalVideoProvider|SanitizeVideoProviderMessageHidesUpstreamNames)$' -count=1` passed.
- `go test -tags=unit ./internal/handler ./internal/handler/dto -run 'Test(PublicPlatformMappers|BuildPlatformSections|ToModelPlazaGroupDTO|GetMyPlatformQuotas_HidesVideoProvider|GetMyPlatformQuotas_D14)' -count=1` passed.
- `go test -tags=unit ./internal/service ./internal/handler ./internal/handler/dto -count=1` passed: service `162.970s`, handler `29.849s`, dto `0.173s`.
- `frontend/node_modules/.bin/vue-tsc --noEmit` passed.
- `frontend/node_modules/.bin/vitest run src/views/user/__tests__/VideoGenerationView.spec.ts src/components/admin/user/__tests__/UserPlatformQuotaModal.spec.ts` passed: 2 files, 37 tests.
- `git diff --check` passed.

### Notes
- `backend/internal/service/leo_video.go`: defines the public `video` platform alias and sanitizes synchronous provider errors including bare `leo`.
- `backend/internal/service/leo_video_async.go`: sanitizes asynchronous provider errors including bare `leo`.
- `backend/internal/service/leo_video_test.go`: verifies public platform alias and synchronous error redaction.
- `backend/internal/service/leo_video_async_test.go`: verifies asynchronous provider-name redaction.
- `backend/internal/handler/dto/mappers.go`: adds public API key, group, user, and subscription platform mappers while preserving admin mappers.
- `backend/internal/handler/dto/public_platform_mapper_test.go`: verifies public DTO aliasing and admin-value preservation.
- `backend/internal/handler/api_key_handler.go`: applies the public alias to user API key and available-group responses.
- `backend/internal/handler/auth_handler.go`: applies the public alias to login user payloads.
- `backend/internal/handler/user_handler.go`: applies the public alias to user profile and quota responses.
- `backend/internal/handler/subscription_handler.go`: applies the public alias to user subscription responses.
- `backend/internal/handler/usage_handler.go`: applies the public alias to user dashboard platform statistics.
- `backend/internal/handler/payment_handler.go`: applies the public alias to public plan and checkout responses.
- `backend/internal/handler/available_channel_handler.go`: applies the public alias to available-channel sections, groups, and models.
- `backend/internal/handler/available_channel_handler_test.go`: covers public channel aliasing.
- `backend/internal/handler/model_plaza_handler.go`: applies the public alias to model-plaza groups and models.
- `backend/internal/handler/model_plaza_handler_test.go`: covers public model-plaza aliasing.
- `backend/internal/handler/user_platform_quotas_handler_test.go`: covers public quota aliasing.
- `frontend/src/i18n/locales/en/dashboard.ts`: removes the user-facing `Leo` empty-state label.
- `frontend/src/types/index.ts`: accepts the public `video` group platform value.
- `frontend/src/api/admin/users.ts`: accepts the public `video` quota value during mixed-version rollout.
- `frontend/src/utils/platformColors.ts`: adds neutral `Video` labeling and styling for public platform values.
- `frontend/src/components/common/PlatformIcon.vue`: renders the public video alias with the video icon.
- `frontend/src/components/common/GroupBadge.vue`: styles both legacy and public video platform values.
- `frontend/src/components/user/dashboard/UserDashboardStats.vue`: displays public and legacy video platform values as `Video`.
- `frontend/src/views/user/VideoGenerationView.vue`: accepts both `video` and legacy `leo` API key group values.
- `docs/LEO_VIDEO_CHANNEL.md`: documents the public/internal platform boundary.
- `progress.md`: records this implementation, verification evidence, and rollback point.
- Rollback before commit: run `git restore --` on the listed tracked files and separately remove the untracked `backend/internal/handler/dto/public_platform_mapper_test.go`; after commit, use `git revert <commit>`.

## 2026-08-03 - Task: Final public platform verification

### What was done
- Completed the final static audit and production frontend build for the public `video` platform alias.
- Confirmed the pre-existing untracked `.superpowers/` directory remains untouched and build output remains outside the tracked diff.

### Testing
- `frontend/node_modules/.bin/vite build` passed with 994 modules transformed.
- `git diff --check` passed; only existing LF/CRLF conversion warnings for `docs/LEO_VIDEO_CHANNEL.md` and `progress.md` remain.
- Public user source and generated frontend output contain no user-facing `Leo`, `LeoStudio`, `Leonardo`, or old Leo empty-state text outside tests and internal compatibility identifiers.

### Notes
- Source verification covered the public user views, public DTO response boundaries, and generated frontend bundle.
- `progress.md`: records the final build verification and unchanged temporary files.
- Rollback point: use the rollback procedure in the preceding task entry; after commit, use `git revert <commit>`.

## 2026-08-03 - Task: Merge upstream v0.1.170 and prepare fork release

### What was done
- Fetched `upstream/main` at `27e8f69a9` and merged the upstream v0.1.170 code into `codex/leo-video-channel`.
- Preserved the fork-only Leo video and token incentive modules, public `video` platform alias, custom Telemetry homepage, and current user-facing二开内容 while accepting upstream security, routing, billing, authentication, and admin improvements.
- Resolved conflicts in generated Ent hooks, gateway route guards, OpenAI slot acquisition, Codex probe identity tests, stream failover classification, public channel tests, and `HomeView.vue`.
- Adapted the Leo video slot caller to the upstream three-state slot result and retained upstream Responses subpath validation alongside the Leo feature gate.

### Testing
- `go test ./...` passed for all backend packages.
- `frontend/node_modules/.bin/eslint.cmd . --ext .vue,.js,.jsx,.cjs,.mjs,.ts,.tsx,.cts,.mts` passed.
- `frontend/node_modules/.bin/vue-tsc.cmd --noEmit` passed.
- Frontend homepage regression tests passed: 3 files, 17 tests.
- `frontend/node_modules/.bin/vite.cmd build` passed with 997 modules transformed and generated `backend/internal/web/dist` output.
- `git diff --cached --check` passed and no unresolved merge paths or conflict markers remain.

### Notes
- `backend/ent/client.go`: retained `VideoJob` hook and interceptor registrations while merging upstream generated Ent changes.
- `backend/internal/handler/openai_gateway_handler.go`: combined fork release timing compatibility with upstream profit-control slot result handling.
- `backend/internal/server/routes/gateway.go`: combined Leo platform rejection with upstream Responses path guard.
- `backend/internal/handler/leo_video.go`: adapted the Leo caller to the upstream slot acquisition result type.
- `backend/internal/service/account_usage_service.go`: retained Codex fingerprint probe headers and upstream load-shed identity normalization.
- `backend/internal/service/openai_codex_identity_test.go`: preserved fork UUID coverage and upstream config coverage.
- `backend/internal/service/openai_gateway_passthrough.go`: preserved non-retryable Codex plan gating and upstream rate-limit failover handling.
- `frontend/src/views/HomeView.vue`: retained Telemetry/custom-menu default homepage and added upstream compact homepage support with theme/i18n state.
- `frontend/src/views/__tests__/HomeView.compact.spec.ts`: aligned the default-home assertion with the retained Telemetry homepage.
- All other staged files are the upstream v0.1.170 merge result, including migrations 192/193 and release workflow updates.
- `.superpowers/` remains untracked and excluded from the merge commit and release package.
- Rollback point: after the merge commit is created, run `git revert <merge_commit>`; the pre-merge fork point is `905d63216`.

## 2026-08-06 - Task: Configurable cyber_policy exclusive-group revocation

### What was done
- Added a disabled-by-default `cyber_policy_revoke_group_id` setting with backend validation restricted to active OpenAI exclusive groups.
- Added server-side revocation after upstream `cyber_policy` logging for regular users, with administrator exemption and per-user auth-cache invalidation. The rule removes only the selected `user_allowed_groups` relation and does not disable users or alter other groups, API keys, account status, or account pools.
- Added the risk-control Web selector and bilingual labels, plus operational documentation.

### Testing
- `go test ./internal/service -run 'TestRecordCyberPolicyEvent|TestContentModerationConfig' -count=1` passed.
- `go test ./internal/service ./internal/handler/admin -count=1` passed.
- `frontend/node_modules/.bin/vue-tsc.cmd --noEmit` passed.
- `git diff --check` passed.

### Notes
- `backend/internal/service/content_moderation.go`: stores, validates, exposes, and applies the configured group revocation with administrator exemption.
- `backend/internal/service/content_moderation_test.go`: records test user group removals for assertions.
- `backend/internal/service/content_moderation_cyber_revoke_test.go`: covers regular-user revocation and administrator exemption.
- `backend/internal/handler/admin/content_moderation_handler.go`: forwards the new admin configuration field.
- `frontend/src/api/admin/riskControl.ts`: adds the configuration field to read/write types.
- `frontend/src/views/admin/RiskControlView.vue`: adds the OpenAI exclusive-group selector and payload binding.
- `frontend/src/i18n/locales/zh/admin/channels.ts`: adds Chinese selector labels and guidance.
- `frontend/src/i18n/locales/en/admin/channels.ts`: adds English selector labels and guidance.
- `docs/CYBER_POLICY_REVOCATION_BAN.md`: documents the configurable revocation rule.
- `docs/CONTENT_MODERATION.md`: documents the administrator exemption and group-only effect.
- Rollback: run `git restore -- backend/internal/service/content_moderation.go backend/internal/service/content_moderation_test.go backend/internal/handler/admin/content_moderation_handler.go frontend/src/api/admin/riskControl.ts frontend/src/views/admin/RiskControlView.vue frontend/src/i18n/locales/zh/admin/channels.ts frontend/src/i18n/locales/en/admin/channels.ts docs/CYBER_POLICY_REVOCATION_BAN.md docs/CONTENT_MODERATION.md` and remove `backend/internal/service/content_moderation_cyber_revoke_test.go`.

## 2026-08-06 - Task: cyber_policy revocation frontend verification

### What was done
- Completed the focused locale regression check for the risk-control page after adding the group selector.

### Testing
- `frontend/node_modules/.bin/vitest.cmd run src/i18n/__tests__/riskControlLocales.spec.ts` passed: 1 file, 3 tests.

### Notes
- No additional source files changed in this verification pass.
- Rollback point is unchanged from the preceding task entry.

## 2026-08-06 - Task: cyber_policy revocation config regression coverage

### What was done
- Added direct regression coverage for the default-disabled setting and save/reload behavior with a valid active OpenAI exclusive group.

### Testing
- `go test ./internal/service -run 'TestContentModerationConfig_(CyberPolicyRevokeGroupDefaultsDisabled|SavesValidCyberPolicyRevokeGroup)|TestRecordCyberPolicyEvent_(RemovesConfiguredExclusiveGroupForRegularUser|AdministratorsKeepConfiguredExclusiveGroup)' -count=1` passed.

### Notes
- `backend/internal/service/content_moderation_cyber_revoke_test.go`: adds configuration persistence coverage alongside revocation behavior coverage.
- Rollback point is unchanged from the original implementation entry.

## 2026-08-06 - Task: cyber_policy revocation final lint verification

### What was done
- Completed the focused frontend lint pass after the configuration and locale updates.

### Testing
- `frontend/node_modules/.bin/eslint.cmd src/views/admin/RiskControlView.vue src/api/admin/riskControl.ts src/i18n/locales/en/admin/channels.ts src/i18n/locales/zh/admin/channels.ts` passed.
- Final `git diff --check` passed.

### Notes
- No additional source files changed in this verification pass.
- Rollback point is unchanged from the original implementation entry.

## 2026-08-06 - Task: Merge upstream v0.1.171 and prepare fork release

### What was done
- Merged `upstream/main` at `c123caddd` (upstream release `v0.1.171`) into `codex/leo-video-channel` after committing the configurable `cyber_policy` group-revocation feature at `7716a380e`.
- Preserved the fork-only Leo video workflow, Token incentive plan, public `video` platform alias, Telemetry homepage, custom update/release behavior, and the new risk-control revocation setting.
- Accepted upstream captcha gates, Codex version synchronization, billing quantization, subscription renewal locking, OpenAI account recovery, scheduler cancellation, model-plaza pricing, authentication, payment, and sponsor updates.
- Resolved 12 content conflicts by combining fork behavior with upstream v0.1.171 behavior; no database migration conflicts or new migration files were introduced.

### Testing
- `go test ./...` passed for all backend packages.
- `frontend/node_modules/.bin/vue-tsc.cmd --noEmit` passed.
- `frontend/node_modules/.bin/eslint.cmd . --ext .vue,.js,.jsx,.cjs,.mjs,.ts,.tsx,.cts,.mts` passed.
- Focused frontend regression suite passed: 11 files, 116 tests covering the Telemetry/default homepage, compact homepage, video generation, usage/Token incentive page, admin settings, captcha gates, model plaza, and risk-control locales.
- `frontend/node_modules/.bin/vite.cmd build` passed with 1005 modules transformed.
- `git diff --cached --check` passed before the progress-log update; the final staged check is repeated before commit.

### Notes
- `backend/cmd/server/wire_gen.go`: registers both the fork video-job runtime and the upstream Codex-version synchronization service in cleanup wiring.
- `backend/internal/config/config.go`: combines video `blob:` media support with upstream Aliyun and Tencent captcha CSP domains.
- `backend/internal/handler/model_plaza_handler.go`: preserves the public `video` alias while exposing upstream independent image-rate fields.
- `backend/internal/server/middleware/security_headers_test.go`: verifies both media blob support and Tencent captcha CSP support.
- `backend/internal/service/account_test_service.go`: preserves account-scoped Codex fingerprint headers and applies upstream account-UA identity normalization.
- `backend/internal/service/account_usage_service.go`: preserves usage-probe fingerprint isolation and applies upstream account-UA identity normalization.
- `backend/internal/service/openai_alpha_search.go`: preserves fingerprint isolation and adopts upstream unified Codex identity/version handling.
- `backend/internal/service/openai_codex_identity_test.go`: preserves fork fingerprint lifecycle tests and adds upstream client-version normalization tests.
- `backend/internal/service/openai_gateway_forward.go`: combines device fingerprint mapping with upstream unified OAuth identity enforcement.
- `backend/internal/service/openai_gateway_passthrough.go`: combines passthrough fingerprint mapping with upstream unified OAuth identity enforcement.
- `backend/internal/service/openai_ws_forwarder_payload.go`: combines WebSocket fingerprint mapping with upstream unified OAuth identity enforcement.
- `deploy/config.example.yaml`: documents the combined captcha and video-media CSP policy.
- `.github/workflows/release.yml`: remains the fork Linux amd64 Release workflow and was verified to support the planned `v0.1.171-fy.1` tag and checksum asset.
- All remaining staged files are the direct upstream v0.1.171 merge result; `progress.md` records the merge decisions, verification evidence, and rollback point.
- `.superpowers/` remains untracked and excluded from the merge commit and release package.
- Rollback point: after the merge commit is created, run `git revert <merge_commit>`; the pre-merge fork point is `7716a380e`.

## 2026-08-07 - Task: Sync LeoStudio latest video models

### What was done
- Synchronized the LeoStudio model registry across server validation, account defaults, channel pricing, the video workbench, and the public API documentation for Hailuo 03, Gemini Omni Flash, Kling 2.x/3.x and O-series, and Veo 3.1 variants.
- Enforced model-specific resolution, duration, aspect-ratio, generated-audio, frame, image, video, and audio-reference rules, including `auto` pairing, Hailuo audio limits, and GENERATED-only Kling O-series video references.
- Completed multi-audio upload handling in the workbench, kept uploaded media payloads to public `media_url`/`UPLOADED` fields, and kept local video upload disabled for models that cannot accept `UPLOADED` video.
- Updated the English/Chinese API documentation page, model matrix/examples, formal Leo specifications/channel docs, native pricing tiers, and regression tests without exposing upstream credentials, UUIDs, or internal task identifiers to customers.

### Testing
- Backend: `go test ./...` passed for all packages after updating the Leo default-model list assertion.
- Frontend: full Vitest run passed; focused video/API-docs/pricing suite passed 47 tests.
- Frontend: `vue-tsc --noEmit`, targeted ESLint, and `vite build` passed; build transformed 1005 modules with only existing chunk-size/dynamic-import warnings.
- `git diff --check` passed.

### Notes
- `backend/internal/service/leo_account.go`: exposes the latest public Leo model defaults.
- `backend/internal/service/leo_account_test.go`: updates default-model regression coverage.
- `backend/internal/service/leo_video_model_specs.go`: adds latest model capabilities and media validation, including Seedance Mini audio support.
- `backend/internal/service/leo_video_model_specs_test.go`: covers new model parameters and guidance constraints.
- `backend/internal/service/video_billing_resolution.go`: maps native model resolutions to billing tiers.
- `frontend/src/views/user/VideoGenerationView.vue`: filters model parameters, supports multiple audio references, validates media constraints, and blocks unsupported uploads.
- `frontend/src/views/user/VideoApiDocsView.vue`: documents all current models and customer-visible request examples.
- `frontend/src/constants/channel.ts`: adds the latest Leo model catalog.
- `frontend/src/components/admin/channel/types.ts`: adds native pricing-tier mappings.
- `frontend/src/i18n/locales/en/dashboard.ts`: updates English model and media guidance.
- `frontend/src/i18n/locales/zh/dashboard.ts`: updates Chinese model and media guidance.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: adds workbench model/media regression tests.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: verifies all model examples and public exposure boundaries.
- `frontend/src/components/admin/channel/__tests__/types.spec.ts`: verifies native pricing tiers.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: records the complete public model and media contract.
- `docs/LEO_VIDEO_CHANNEL.md`: records channel mapping, pricing, workbench, and upload operations.
- `progress.md`: records this implementation and verification round.
- Rollback: `git restore -- backend/internal/service/leo_account.go backend/internal/service/leo_account_test.go backend/internal/service/leo_video_model_specs.go backend/internal/service/leo_video_model_specs_test.go backend/internal/service/video_billing_resolution.go frontend/src/views/user/VideoGenerationView.vue frontend/src/views/user/VideoApiDocsView.vue frontend/src/constants/channel.ts frontend/src/components/admin/channel/types.ts frontend/src/i18n/locales/en/dashboard.ts frontend/src/i18n/locales/zh/dashboard.ts frontend/src/views/user/__tests__/VideoGenerationView.spec.ts frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts frontend/src/components/admin/channel/__tests__/types.spec.ts docs/LEO_VIDEO_MODEL_SPECS.md docs/LEO_VIDEO_CHANNEL.md progress.md`

## 2026-08-09 - Task: Audit Leo video public boundary and latest model docs

### What was done
- Unified synchronous and asynchronous Leo video result redaction so customer-facing responses retain public media fields while removing provider names, UUIDs, generation identifiers, account identifiers, source URLs, and credential-like fields.
- Corrected the single-image Gemini Omni Flash workbench payload to use `guidances.image_reference`, matching the model capability that rejects start frames.
- Removed duplicate API documentation examples and aligned Grok's `auto` resolution/aspect-ratio pairing in the model matrix and localized guidance.
- Updated regression coverage for public result redaction, Gemini image guidance, unique API examples, and the affected service response assertion.

### Testing
- `go test ./internal/handler ./internal/service -count=1` passed.
- `frontend/node_modules/.bin/vitest.cmd run src/views/user/__tests__/VideoGenerationView.spec.ts src/views/user/__tests__/VideoApiDocsView.spec.ts --reporter=verbose` passed: 2 files, 31 tests.
- `frontend/node_modules/.bin/vue-tsc.cmd --noEmit` passed.
- `frontend/node_modules/.bin/eslint.cmd src/views/user/VideoGenerationView.vue src/views/user/VideoApiDocsView.vue src/views/user/__tests__/VideoGenerationView.spec.ts src/views/user/__tests__/VideoApiDocsView.spec.ts` passed.
- `frontend/node_modules/.bin/vite.cmd build` passed with 1005 modules transformed; only existing chunk-size and dynamic-import warnings remain.
- `git diff --check` passed.

### Notes
- `backend/internal/service/leo_video.go`: centralizes public result metadata filtering for synchronous and asynchronous paths.
- `backend/internal/service/leo_video_test.go`: updates the response assertion for the public redacted body.
- `backend/internal/handler/leo_video_async.go`: reuses the service-level redaction helper.
- `backend/internal/handler/leo_video_integration_test.go`: verifies synchronous responses hide provider metadata.
- `docs/LEO_VIDEO_CHANNEL.md`: documents the public synchronous redaction boundary.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: documents Grok `auto` resolution support.
- `frontend/src/views/user/VideoGenerationView.vue`: sends Gemini single-image inputs as image guidance when start frames are unsupported.
- `frontend/src/views/user/VideoApiDocsView.vue`: fixes the Gemini example and removes duplicate model examples.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: covers Gemini single-image request construction.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: covers example uniqueness and public Gemini guidance.
- `frontend/src/i18n/locales/en/dashboard.ts`: updates English Grok `auto` guidance.
- `frontend/src/i18n/locales/zh/dashboard.ts`: updates Chinese Grok `auto` guidance.
- `progress.md`: records this audit, verification evidence, and rollback point.
- `.superpowers/` remains an existing untracked directory and is excluded from the commit.
- Rollback: `git restore -- backend/internal/service/leo_video.go backend/internal/service/leo_video_test.go backend/internal/handler/leo_video_async.go backend/internal/handler/leo_video_integration_test.go docs/LEO_VIDEO_CHANNEL.md docs/LEO_VIDEO_MODEL_SPECS.md frontend/src/views/user/VideoGenerationView.vue frontend/src/views/user/VideoApiDocsView.vue frontend/src/views/user/__tests__/VideoGenerationView.spec.ts frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts frontend/src/i18n/locales/en/dashboard.ts frontend/src/i18n/locales/zh/dashboard.ts progress.md`.

## 2026-08-09 - Task: Enable latest Leo video models in channel sync

### What was done
- Fixed the admin channel "Sync Latest Models" endpoint so the `leo` video platform returns the current Leo video model registry instead of being rejected as an unsupported LiteLLM platform.
- Kept the response isolated from the global registry by returning a copy of the model slice.
- Updated the admin API comment and regression coverage for the full Leo model list and copy isolation.

### Testing
- `go test -tags unit ./internal/handler/admin -run TestSyncPricingModels -count=1 -v` passed.
- `go test ./internal/service -run 'TestDefaultModelsListCandidateIDsLeo' -count=1` passed.
- `frontend/node_modules/.bin/vitest.cmd run src/components/admin/channel/__tests__/types.spec.ts src/components/admin/channel/__tests__/PricingEntryCard.spec.ts --reporter=verbose` passed: 2 files, 21 tests.
- `frontend/node_modules/.bin/vue-tsc.cmd --noEmit` passed.
- `frontend/node_modules/.bin/eslint.cmd src/api/admin/channels.ts` passed.
- `frontend/node_modules/.bin/vite.cmd build` passed with 1005 modules transformed; only existing chunk-size and dynamic-import warnings remain.
- `git diff --check` passed.

### Notes
- `backend/internal/handler/admin/channel_handler.go`: serves the latest Leo video model registry from the sync endpoint.
- `backend/internal/handler/admin/channel_handler_test.go`: verifies Leo sync status, model list, and copy isolation.
- `frontend/src/api/admin/channels.ts`: documents that Leo sync uses the video registry.
- `progress.md`: records this fix, verification evidence, and rollback point.
- `.superpowers/` remains an existing untracked directory and is excluded from the commit.
- Rollback: `git restore -- backend/internal/handler/admin/channel_handler.go backend/internal/handler/admin/channel_handler_test.go frontend/src/api/admin/channels.ts progress.md`.

## 2026-08-09 - Task: Remove Leo provider names from customer-facing gateway errors

### What was done
- Replaced the two customer-facing gateway rejection messages that named the internal Leo platform with neutral platform/video API wording.
- Kept internal routing, admin labels, logs, and upstream error sanitization unchanged; only public gateway response text was adjusted.
- Added route regression assertions covering unsupported video upload/job paths and unsupported capabilities for video groups.

### Testing
- `go test ./internal/server/routes ./internal/handler` passed.
- `git diff --check` passed.
- Source audit confirmed the remaining `Leo` provider-name regexes are internal sanitizers and are not emitted directly by customer gateway responses.

### Notes
- `backend/internal/server/routes/gateway.go`: removes `Leo` from public unsupported-capability and video-only route messages.
- `backend/internal/server/routes/gateway_test.go`: verifies customer responses contain the neutral message and no `leo` token.
- `progress.md`: records this audit, fix, and verification evidence.
- Rollback: `git restore -- backend/internal/server/routes/gateway.go backend/internal/server/routes/gateway_test.go progress.md`.

## 2026-08-09 - Task: Add concurrent image `n` fan-out with account-pool scheduling

### What was done
- Added non-streaming `n>1` image fan-out for OpenAI-compatible OpenAI and Grok image endpoints. Each requested image is dispatched as an independent single-image child request and the successful JSON responses are combined into one standard `data[]` response.
- Kept account scheduling, account concurrency slots, failover, and pool-mode same-account retry in the child path; usage is recorded per successful child account/result so multi-account dispatch does not collapse billing attribution.
- Added JSON and multipart `n=1` rewriting for child requests, recursive usage merging, and regression coverage for request rewriting, response merging, concurrent handler dispatch, and race detection.
- Added public API documentation for `n` behavior without exposing upstream account or provider-internal identifiers.

### Testing
- `go test ./internal/service -run 'TestRewriteOpenAIImagesN_RewritesJSONAndMultipart|TestMergeOpenAIImageResponses_AppendsDataAndSumsUsage' -count=1` passed.
- `go test -tags unit ./internal/handler -run 'TestOpenAIGatewayHandlerImages_NSplitsIntoConcurrentSingleImageRequests|TestOpenAIGatewayHandlerImages_NSplitSupportsPoolModeAPIKey' -count=1 -v` passed; observed three child calls plus two child calls on one pool-mode API-key account, with combined image output.
- `go test -race -tags unit ./internal/handler -run TestOpenAIGatewayHandlerImages_NSplitsIntoConcurrentSingleImageRequests -count=1` passed.
- `go test ./internal/handler ./internal/service -run '^$' -count=1` and `go test -tags unit ./internal/handler -run '^$' -count=1` passed.
- `git diff --check` passed.

### Notes
- `backend/internal/handler/openai_images.go`: routes non-streaming multi-image requests to the fan-out path.
- `backend/internal/handler/openai_images_split.go`: schedules child OpenAI image tasks, enforces slots/failover, merges output, and records per-account usage.
- `backend/internal/handler/grok_media.go`: routes Grok image generation/edit multi-image requests to fan-out.
- `backend/internal/handler/grok_images_split.go`: schedules and merges Grok image child tasks with pool-mode handling.
- `backend/internal/handler/openai_images_failover_test.go`: verifies three concurrent child requests and combined output.
- `backend/internal/service/openai_images.go`: rewrites child `n` values and merges OpenAI-compatible image responses.
- `backend/internal/service/openai_images_n_split_test.go`: covers JSON/multipart rewriting and usage/data merge behavior.
- `docs/IMAGE_BATCH_REQUESTS.md`: documents the public `n` contract and scheduling/billing behavior.
- `progress.md`: records this implementation and verification round.
- `.superpowers/` remains an existing untracked directory and is excluded from the commit.
- Rollback: `git restore -- backend/internal/handler/openai_images.go backend/internal/handler/openai_images_split.go backend/internal/handler/grok_media.go backend/internal/handler/grok_images_split.go backend/internal/handler/openai_images_failover_test.go backend/internal/service/openai_images.go backend/internal/service/openai_images_n_split_test.go docs/IMAGE_BATCH_REQUESTS.md progress.md`.

## 2026-08-09 - Task: Harden concurrent image `n` fan-out and pool-mode billing

### What was done
- Bound every pool-mode child request to its selected account before retry so configured same-account retries cannot drift to another account.
- Derived a distinct, deterministic billing request ID for each successful child image so one client request with `n>1` records and deducts every returned image without losing request-level idempotency.
- Limited `n` to integers from 1 through 10 across JSON, multipart, OpenAI-compatible, and Grok image routes to prevent unbounded goroutine fan-out.
- Made the new image batching document trackable and aligned its public validation contract with the implemented limit.

### Testing
- `go test ./internal/service ./internal/handler -count=1` passed.
- `go test -tags unit ./internal/handler -count=1` passed.
- `go test -race -tags unit ./internal/handler -run 'TestOpenAIGatewayHandlerImages_NSplitsIntoConcurrentSingleImageRequests|TestOpenAIGatewayHandlerImages_NSplitSupportsPoolModeAPIKey' -count=1` passed.
- `go vet ./internal/service ./internal/handler` passed.
- The pool-mode regression simulated an initial HTTP 502, observed three `n=1` upstream calls on the same account, and returned two combined images.
- `git diff --check` passed.

### Notes
- `.gitignore`: tracks the formal image batching API document.
- `backend/internal/handler/openai_images_split.go`: binds pool sessions and separates child billing request identities.
- `backend/internal/handler/grok_images_split.go`: applies the same pool retry and billing identity behavior to Grok image requests.
- `backend/internal/handler/grok_media.go`: validates the public image count and accepts an explicit usage-record parent context.
- `backend/internal/handler/openai_images_failover_test.go`: covers pool retry after a transient failure, sticky binding, and distinct child billing IDs.
- `backend/internal/service/openai_images.go`: enforces integer `n` values from 1 through 10 for JSON and multipart image requests.
- `backend/internal/service/grok_media.go`: parses and validates Grok image counts without changing video request behavior.
- `backend/internal/service/openai_images_n_split_test.go`: covers safe fan-out limits in addition to child request rewriting and response merging.
- `docs/IMAGE_BATCH_REQUESTS.md`: documents `n=1-10`, concurrent dispatch, partial success, pool scheduling, and public response boundaries.
- `progress.md`: records this hardening and verification round.
- `.superpowers/` remains an existing untracked directory and was not modified.
- Rollback the complete image `n` feature with `git restore -- .gitignore backend/internal/handler/openai_images.go backend/internal/handler/grok_media.go backend/internal/handler/openai_images_failover_test.go backend/internal/service/openai_images.go backend/internal/service/grok_media.go`, then run `Remove-Item -LiteralPath backend/internal/handler/openai_images_split.go,backend/internal/handler/grok_images_split.go,backend/internal/service/openai_images_n_split_test.go,docs/IMAGE_BATCH_REQUESTS.md`; keep this append-only progress history as the audit record.

## 2026-08-09 - Task: Merge upstream v0.1.172 and publish fork release

### What was done
- Merged `upstream/main` at `fb0475656` (upstream tag `v0.1.172`) into `codex/leo-video-channel` on top of fork commit `6c3490659`.
- Resolved 21 merge conflicts while retaining the fork's Leo video routes and pricing, token-incentive settings, account scheduling thresholds, mandatory media billing, and OpenAI/Grok image `n=1-10` account-pool fan-out.
- Integrated upstream Grok Voice, Web Search, async video billing, per-model video prices, response-model auditing, CAPTCHA/security updates, and migrations `194`, `195`, and `217-220`.
- Updated the unit-tag image fan-out cache stub for the expanded upstream Grok video billing cache contract.
- Prepared release tag `v0.1.172-fy.1`; the tag-triggered workflow builds and publishes the Linux amd64 archive and SHA-256 checksums.

### Testing
- `go test ./... -count=1` passed for all backend packages.
- `go vet ./...` passed.
- `go test -race -tags unit ./internal/handler -run 'TestOpenAIGatewayHandlerImages_NSplitsIntoConcurrentSingleImageRequests|TestOpenAIGatewayHandlerImages_NSplitSupportsPoolModeAPIKey' -count=1` passed.
- Full frontend Vitest passed: 221 files and 1558 tests.
- Conflict-focused frontend Vitest passed: 4 files and 48 tests.
- `vue-tsc --noEmit` passed.
- Full frontend ESLint check passed.
- Vite production build passed; only the existing chunk-size warning remains.
- `git diff --cached --check` passed and no unmerged paths remain.

### Notes
- `.gitignore`: merged the upstream v0.1.172 change for this file within the fork release set.
- `README.md`: merged the upstream v0.1.172 change for this file within the fork release set.
- `README_CN.md`: merged the upstream v0.1.172 change for this file within the fork release set.
- `README_JA.md`: merged the upstream v0.1.172 change for this file within the fork release set.
- `assets/partners/logos/haoai.png`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/cmd/server/VERSION`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/cmd/server/wire_gen.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/group.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/group/group.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/group/where.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/group_create.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/group_update.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/migrate/schema.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/mutation.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/runtime/runtime.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/schema/group.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/schema/usage_log.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/usagelog.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/usagelog/usagelog.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/usagelog/where.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/usagelog_create.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/ent/usagelog_update.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/go.sum`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/config/config.go`: merged the upstream Tencent CAPTCHA CSP allowlist while retaining blob media support for the fork video surface.
- `backend/internal/config/config_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/domain/constants.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/domain/constants_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/account_handler.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/dashboard_handler.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/dashboard_handler_request_type_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/dashboard_query_cache.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/dashboard_snapshot_v2_handler.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/grok_import_probe.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/grok_import_probe_handler_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/grok_import_probe_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/grok_oauth_handler.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/grok_oauth_handler_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/group_handler.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/setting_handler.go`: returns upstream Grok and scheduling settings together with the fork video and token-incentive settings.
- `backend/internal/handler/admin/setting_handler_audit.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/setting_handler_partial_payload_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/setting_handler_update.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/usage_handler.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/admin/usage_query_cache.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/auth_captcha_request_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/auth_oauth_pending_flow.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/auth_oauth_pending_flow_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/dto/mappers.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/dto/mappers_usage_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/dto/settings.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/dto/types.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/gateway_handler.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/gateway_web_search.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/grok_audio.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/grok_audio_billing_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/grok_media.go`: merged async Grok video billing with the fork image fan-out usage context and Leo media behavior.
- `backend/internal/handler/grok_media_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/openai_gateway_handler.go`: makes image, video, search, web-search, and voice usage writes mandatory when the async queue overflows.
- `backend/internal/handler/openai_gateway_handler_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/openai_images.go`: retains n>1 image fan-out while adding the upstream OpenAI images endpoint context.
- `backend/internal/handler/openai_images_failover_test.go`: keeps fork image fan-out coverage compatible with the expanded upstream GatewayCache interface.
- `backend/internal/handler/setting_handler.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/setting_handler_public_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/usage_handler_request_type_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/handler/usage_record_submit_task_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/antigravity/claude_types.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/antigravity/claude_types_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/antigravity/client.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/apicompat/responses_to_anthropic_invalid_blocks_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/apicompat/responses_to_anthropic_request.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/ctxkey/ctxkey.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/ip/ip.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/ip/ip_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/openai/request.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/proxyutil/dialer.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/proxyutil/dialer_timeout_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/redissession/store.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/redissession/store_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/usagestats/usage_log_types.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/xai/billing.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/xai/billing_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/xai/cli_identity.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/xai/cli_identity_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/xai/models.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/xai/models_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/xai/oauth.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/xai/oauth_redis_fallback_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/xai/oauth_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/xai/quota.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/xai/quota_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/xai/sso_device.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/pkg/xai/sso_device_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/api_key_repo.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/api_key_repo_messages_dispatch_unit_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/gateway_cache.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/grok_oauth_client.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/grok_oauth_client_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/group_repo.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/http_upstream.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/http_upstream_dial_timeout_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/migrations_runner.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/migrations_runner_notx_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/migrations_schema_integration_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/tencent_captcha_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/tencent_captcha_service_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/usage_log_repo_insert.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/usage_log_repo_query.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/usage_log_repo_request_type_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/usage_log_repo_stats.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/usage_log_repo_stats_integration_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/usage_log_repo_trend.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/usage_log_session_id_unit_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/user_subscription_repo.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/repository/user_subscription_repo_integration_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/server/api_contract_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/server/middleware/api_key_auth_google_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/server/middleware/api_key_auth_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/server/middleware/security_headers.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/server/middleware/security_headers_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/server/routes/admin.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/server/routes/gateway.go`: keeps Leo upload/job/internal routes and adds upstream public video, Grok Voice, and Web Search routes.
- `backend/internal/server/routes/gateway_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/server/routes/prompt_audit_route_coverage_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account_credentials_redact.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account_grok_media_eligibility.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account_grok_media_eligibility_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account_scheduling_threshold_eval.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account_scheduling_threshold_eval_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account_scheduling_threshold_integration_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account_scheduling_threshold_reason.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account_scheduling_threshold_reason_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account_test_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account_test_service_grok_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account_usage_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account_usage_service_batch_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/account_wildcard_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/admin_account.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/admin_group.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/admin_group_duplicate.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/admin_group_duplicate_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/admin_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/antigravity_gateway_claude.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/antigravity_gateway_compat.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/antigravity_gateway_compat_stream.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/antigravity_gateway_gemini.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/antigravity_gateway_streaming.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/antigravity_gateway_upstream.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/api_key_auth_cache.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/api_key_auth_cache_impl.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/api_key_auth_cache_profit_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/billing_search_audio_cost_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/billing_service.go`: combines Leo resolution pricing with upstream per-model video, search, and audio billing.
- `backend/internal/service/channel_plaza.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/channel_plaza_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/credentials_sanitize.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/credentials_sanitize_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/dashboard_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/domain_constants.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/gateway_anthropic_passthrough.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/gateway_forward.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/gateway_hotpath_optimization_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/gateway_multiplatform_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/gateway_scheduling.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/gateway_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/gateway_upstream_response.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/gateway_usage_billing.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/gateway_usage_billing_request_id_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/gemini_messages_compat_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/gemini_multiplatform_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_audio.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_audio_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_base_url_mode_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_credential_failure.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_credential_failure_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_free_quota_gate.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_free_quota_gate_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_media.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_media_video_billing_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_model_quota_block.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_oauth_reconciliation_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_oauth_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_oauth_service_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_observed_models.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_p2_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_quota_fetcher.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_quota_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_quota_service_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_search_count.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_search_count_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_spending_reauth.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_stream_idle.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_stream_idle_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_team_rate_limit.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_team_rate_limit_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_token_refresher.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_token_refresher_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_upstream_errors.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_upstream_failure.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_upstream_failure_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_upstream_headers.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_upstream_headers_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_upstream_url.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/grok_upstream_url_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/group.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/media_price_config.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/model_rate_limit.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/oauth_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_account_scheduler.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_account_scheduler_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_alpha_search.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_alpha_search_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_capacity_shed_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_codex_identity.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_codex_identity_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_codex_models_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_codex_models_service_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_codex_pat_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_codex_pat_service_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_codex_version_consistency_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_codex_version_sync_service_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_compact_service_tier_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_compat_model_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_cyber_session_block_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_chat_completions.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_chat_completions_raw.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_count_tokens.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_count_tokens_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_forward.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_grok.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_grok_405_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_grok_cache.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_grok_cache_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_grok_chat_bridge.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_grok_search_billing_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_grok_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_messages.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_passthrough.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_record_usage_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_request_body.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_response_flush_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_response_handling.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_scheduling.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_search_surcharge_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_service_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_upstream_errors.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_gateway_usage.go`: combines fork-wide video billing with upstream response-model auditing and search/audio surcharges.
- `backend/internal/service/openai_images_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_messages_dispatch.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_messages_dispatch_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_oauth_passthrough_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_responses_tool_schema.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_responses_tool_schema_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_routing_hint.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_routing_hint_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_ws_forwarder_ingress.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_ws_forwarder_ingress_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_ws_forwarder_payload.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_ws_forwarder_success_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_ws_forwarder_v2.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_ws_http_bridge.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_ws_http_bridge_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_ws_pool.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_ws_pool_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_ws_state_store_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_ws_v2/passthrough_relay.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_ws_v2/passthrough_relay_internal_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/openai_ws_v2_passthrough_adapter.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/ops_system_log_sink.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/ops_system_log_sink_backoff_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/ratelimit_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/ratelimit_service_model_not_found_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/ratelimit_service_scheduling_threshold_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/setting_features.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/setting_gateway_runtime.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/setting_parse.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/setting_public.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/setting_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/setting_service_platform_threshold_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/setting_update.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/settings_view.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/subscription_assign_idempotency_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/subscription_daily_midnight_reset_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/subscription_expiry_service_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/subscription_monthly_window_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/subscription_reset_quota_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/subscription_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/temp_unsched.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/tencent_captcha_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/tencent_captcha_service_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/tencent_captcha_settings_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/token_refresh_service.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/upstream_models.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/upstream_models_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/upstream_response_model.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/upstream_response_model_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/usage_log.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/user_subscription.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/user_subscription_daily_quota_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/user_subscription_port.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/video_billing.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/video_billing_resolution.go`: separates generic three-tier resolution fallback from model-aware Leo high-resolution billing.
- `backend/internal/service/video_billing_test.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/websearch_config.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/service/wire.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/internal/testutil/stubs.go`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/migrations/194_add_usage_log_upstream_response_model.sql`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/migrations/195_add_usage_log_upstream_model_mismatch_index_notx.sql`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/migrations/217_group_video_model_prices.sql`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/migrations/218_group_audio_voice_pricing.sql`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/migrations/219_group_search_price_per_1k.sql`: merged the upstream v0.1.172 change for this file within the fork release set.
- `backend/migrations/220_clear_non_grok_video_generation_config.sql`: merged the upstream v0.1.172 change for this file within the fork release set.
- `deploy/config.example.yaml`: documents the merged CAPTCHA CSP domains while retaining blob media support.
- `frontend/pnpm-lock.yaml`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/api/__tests__/admin.grok.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/api/admin/accounts.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/api/admin/dashboard.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/api/admin/grok.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/api/admin/settings.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/api/admin/usage.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/CaptchaChallenge.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/TencentCaptchaGate.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/__tests__/TencentCaptchaGate.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/account/AccountStatusIndicator.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/account/AccountUsageCell.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/account/CreateAccountModal.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/account/EditAccountModal.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/account/GrokQuotaProbeCell.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/account/OAuthAuthorizationFlow.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/account/TempUnschedStatusModal.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/account/__tests__/AccountUsageCell.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/account/__tests__/CreateAccountModal.grok.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/admin/account/AccountTestModal.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/admin/account/ReAuthAccountModal.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/admin/account/__tests__/AccountTestModal.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/admin/account/__tests__/ReAuthAccountModal.grok.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/admin/proxy/ImportDataModal.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/admin/usage/UsageFilters.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/admin/usage/UsageTable.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/admin/usage/__tests__/UsageFilters.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/admin/usage/__tests__/UsageTable.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/auth/PendingOAuthCreateAccountForm.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/common/BaseDialog.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/common/ImageUpload.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/common/PlatformTypeBadge.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/common/__tests__/BaseDialog.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/common/__tests__/PlatformTypeBadge.grok.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/keys/UseKeyModal.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/keys/__tests__/UseKeyModal.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/modelPlaza/PlazaModelPricingTable.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/components/modelPlaza/__tests__/PlazaModelPricingTable.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/composables/__tests__/useGrokOAuth.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/composables/__tests__/useModelWhitelist.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/composables/useAccountOAuth.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/composables/useGrokOAuth.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/composables/useModelWhitelist.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/i18n/locales/en/admin/accounts.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/i18n/locales/en/admin/overview.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/i18n/locales/en/admin/resources.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/i18n/locales/en/admin/settings.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/i18n/locales/en/common.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/i18n/locales/en/dashboard.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/i18n/locales/zh/admin/accounts.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/i18n/locales/zh/admin/overview.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/i18n/locales/zh/admin/resources.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/i18n/locales/zh/admin/settings.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/i18n/locales/zh/common.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/i18n/locales/zh/dashboard.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/types/index.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/utils/__tests__/billingMode.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/utils/billingMode.ts`: keeps explicit video/token modes authoritative, recognizes image usage, and retains the fork token fallback.
- `frontend/src/utils/tencentCaptcha.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/admin/AccountsView.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/admin/GroupsView.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/admin/SettingsView.vue`: loads and saves both token-incentive rules and account scheduling thresholds.
- `frontend/src/views/admin/UsageView.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/admin/__tests__/GroupsView.columnSettings.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/admin/__tests__/GroupsView.duplicate.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/admin/__tests__/SettingsView.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/admin/__tests__/UsageView.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/admin/__tests__/groupsVideoModelPricing.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/admin/groupsVideoModelPricing.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/admin/ops/OpsDashboard.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/admin/ops/components/OpsErrorDetailsModal.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/admin/ops/utils/__tests__/opsErrorParams.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/admin/ops/utils/opsErrorParams.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/auth/EmailVerifyView.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/auth/ForgotPasswordView.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/auth/LoginView.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/auth/RegisterView.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/auth/__tests__/TencentCaptchaForgotPassword.spec.ts`: merged the upstream v0.1.172 change for this file within the fork release set.
- `frontend/src/views/user/UsageView.vue`: merged the upstream v0.1.172 change for this file within the fork release set.
- `progress.md`: records this merge, verification evidence, complete changed-file inventory, and rollback point.
- `.superpowers/`: remains an existing untracked directory and is excluded from the merge and release.
- Rollback after publication with `git revert -m 1 v0.1.172-fy.1`, then push the revert commit and publish a follow-up fork tag.

## 2026-08-09 - Task: Merge upstream v0.1.173 and publish fork release

### What was done
- Merged `upstream/main` at `48eb3766d` after upstream released `v0.1.173`.
- Resolved five content conflicts by retaining both upstream behavior and the fork's token incentive, video generation, homepage telemetry, and concurrent image fan-out behavior.
- Regenerated Wire dependency injection so channel-monitor v2 and the fork runtimes are initialized and stopped together.

### Testing
- `go test ./... -count=1` passed.
- `go vet ./...` passed.
- `node node_modules/vitest/vitest.mjs run --reporter=default --maxWorkers=2 --minWorkers=1` passed; the initial unconstrained Windows worker run ended with `EPIPE`, so the stable bounded-worker run is the release evidence.
- `node node_modules/vue-tsc/bin/vue-tsc.js --noEmit` passed.
- `pnpm.exe lint:check` passed after restoring dependencies with `pnpm.exe install --frozen-lockfile`.
- `pnpm.exe build` passed; Vite transformed 1031 modules and generated the production frontend bundle.

### Notes
- `Makefile`: keeps the fork homepage critical tests and adds the upstream channel-monitor v2 critical tests.
- `backend/cmd/server/VERSION`: merged the upstream v0.1.173 change into the fork release set.
- `backend/cmd/server/wire.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/cmd/server/wire_gen.go`: regenerates dependency injection with both the fork token/video services and upstream channel-monitor v2 services.
- `backend/cmd/server/wire_gen_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/handler/admin/setting_handler.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/handler/admin/setting_handler_audit.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/handler/admin/setting_handler_update.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/handler/auth_oauth_pending_flow_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/handler/channel_monitor_user_handler.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/handler/channel_monitor_v2_handler.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/handler/channel_monitor_v2_handler_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/handler/dto/settings.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/handler/handler.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/handler/setting_handler.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/handler/wire.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/repository/channel_monitor_v2_aggregation.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/repository/channel_monitor_v2_repo.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/repository/channel_monitor_v2_repo_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/repository/migrations_runner.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/repository/user_repo.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/repository/user_repo_email_alias_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/repository/user_repo_integration_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/repository/wire.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/server/api_contract_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/server/routes/admin.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/server/routes/channel_monitor_feature_gate_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/server/routes/user.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/admin_service_delete_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/antigravity_gateway_gemini.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/antigravity_gateway_service_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/antigravity_gateway_streaming.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/auth_oauth_email_flow.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/auth_oauth_email_flow_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/auth_service.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/auth_service_register_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/channel_monitor_const.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/channel_monitor_probe_retirement_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/channel_monitor_runner.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/channel_monitor_service.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/channel_monitor_v2.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/channel_monitor_v2_aggregator.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/channel_monitor_v2_aggregator_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/channel_monitor_v2_error_taxonomy.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/channel_monitor_v2_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/domain_constants.go`: keeps token-incentive setting keys while adding the upstream email-domain quota and channel-monitor v2 keys.
- `backend/internal/service/gemini_chat_completions_compat_service.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/gemini_error_policy_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/gemini_image_output_accounting.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/gemini_image_output_accounting_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/gemini_messages_compat_service.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/openai_gateway_service_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/openai_images.go`: keeps image `n` rewriting and adopts upstream client-disconnect detachment for billable image requests.
- `backend/internal/service/openai_images_upstream_context_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/registration_email_policy.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/registration_email_policy_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/setting_features.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/setting_parse.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/setting_public.go`: exposes both fork video/token flags and upstream channel-monitor mode/privacy settings.
- `backend/internal/service/setting_service.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/setting_service_public_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/setting_update.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/settings_view.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/upstream_response_model.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/upstream_response_model_bench_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/upstream_response_model_test.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/user_service.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/internal/service/wire.go`: merged the upstream v0.1.173 change into the fork release set.
- `backend/migrations/194_channel_monitor_v2.sql`: merged the upstream v0.1.173 change into the fork release set.
- `backend/migrations/195_channel_monitor_mode.sql`: merged the upstream v0.1.173 change into the fork release set.
- `backend/migrations/196_channel_monitor_v2_ignored_error_categories.sql`: merged the upstream v0.1.173 change into the fork release set.
- `backend/migrations/197_channel_monitor_v2_seed_popular_models.sql`: merged the upstream v0.1.173 change into the fork release set.
- `backend/migrations/198_channel_monitor_v2_health_thresholds.sql`: merged the upstream v0.1.173 change into the fork release set.
- `backend/migrations/199_channel_monitor_v2_fixed_rollups.sql`: merged the upstream v0.1.173 change into the fork release set.
- `backend/migrations/200_channel_monitor_v2_rollup_permissions.sql`: merged the upstream v0.1.173 change into the fork release set.
- `backend/migrations/201_channel_monitor_v2_refresh_5m.sql`: merged the upstream v0.1.173 change into the fork release set.
- `backend/migrations/202_channel_monitor_v2_full_table_permissions.sql`: merged the upstream v0.1.173 change into the fork release set.
- `backend/migrations/203_channel_monitor_v2_default_ignore_and_cache.sql`: merged the upstream v0.1.173 change into the fork release set.
- `backend/migrations/204_channel_monitor_hide_throughput.sql`: merged the upstream v0.1.173 change into the fork release set.
- `backend/migrations/205_channel_monitor_v2_reset_factory_cache_thresholds.sql`: merged the upstream v0.1.173 change into the fork release set.
- `backend/migrations/206_channel_monitor_v2_privacy_defaults.sql`: merged the upstream v0.1.173 change into the fork release set.
- `docs/channel-monitor-v2-safe-defaults.md`: adds the upstream operational safety notes for channel-monitor v2.
- `frontend/src/api/__tests__/channelMonitorV2.spec.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/api/admin/settings.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/api/channelMonitorV2.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/components/icons/Icon.vue`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/features/channel-monitor-v2/FilterMultiSelect.vue`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/features/channel-monitor-v2/MetricCell.vue`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/features/channel-monitor-v2/MonitorRankBadge.vue`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/features/channel-monitor-v2/MonitorSettingsPanel.vue`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/features/channel-monitor-v2/MonitorTrendChart.vue`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/features/channel-monitor-v2/RelayPulseMatrix.vue`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/features/channel-monitor-v2/__tests__/MetricCell.spec.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/features/channel-monitor-v2/__tests__/RelayPulseMatrix.spec.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/features/channel-monitor-v2/__tests__/designSystem.structure.spec.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/features/channel-monitor-v2/__tests__/monitorFormat.spec.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/features/channel-monitor-v2/__tests__/monitorZoom.spec.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/features/channel-monitor-v2/monitorFormat.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/features/channel-monitor-v2/monitorZoom.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/i18n/locales/en/admin/settings.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/i18n/locales/en/channelMonitorV2.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/i18n/locales/en/common.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/i18n/locales/en/index.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/i18n/locales/zh/admin/settings.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/i18n/locales/zh/channelMonitorV2.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/i18n/locales/zh/common.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/i18n/locales/zh/index.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/types/index.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/utils/featureFlags.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/views/admin/ChannelMonitorView.vue`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/views/admin/SettingsView.vue`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/views/admin/__tests__/ChannelMonitorView.duplicate.spec.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/views/admin/__tests__/ChannelMonitorView.grok.spec.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/views/admin/__tests__/SettingsView.spec.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/views/auth/EmailVerifyView.vue`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/views/auth/RegisterView.vue`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/views/auth/__tests__/EmailVerifyView.spec.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/views/auth/__tests__/RegisterView.spec.ts`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/views/user/ChannelStatusV1View.vue`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/views/user/ChannelStatusV2View.vue`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/views/user/ChannelStatusView.vue`: merged the upstream v0.1.173 change into the fork release set.
- `frontend/src/views/user/__tests__/ChannelStatusView.mode.spec.ts`: merged the upstream v0.1.173 change into the fork release set.
- `progress.md`: records this merge, verification evidence, changed-file inventory, and rollback point.
- `.superpowers/`: remains an existing untracked directory and is excluded from the merge and release.
- Rollback after publication with `git revert -m 1 v0.1.173-fy.1`, then push the revert commit and publish a follow-up fork tag.

## 2026-08-10 - Task: Redact upstream fallback URLs from client errors
### What was done
- Added client-only URL redaction for upstream error messages, replacing complete HTTP(S) URLs with `[upstream URL]`.
- Applied the redaction to OpenAI WebSocket fallback errors, error passthrough rules, and failover-exhausted gateway responses without changing scheduling or status-code behavior.
- Kept internal Ops upstream messages available for diagnostics while preserving the existing masking of sensitive query values.

### Testing
- `cd backend && go test ./internal/service -run 'Test(ResolveOpenAIWSFallbackErrorResponse|ApplyErrorPassthroughRule)' -count=1` passed.
- `cd backend && go test ./internal/handler -run 'Test(OpenAIGatewayHandler|GatewayHandler|.*Failover.*|.*Stream.*)' -count=1` passed.
- `cd backend && go test ./internal/service ./internal/handler -run '^$' -count=1` passed package compilation.
- `git diff --check` passed.

### Notes
- `backend/internal/service/gemini_messages_compat_service.go`: added the shared client-facing URL redaction helper while leaving the internal sanitizer behavior intact.
- `backend/internal/service/error_passthrough_runtime.go`: redacts upstream URLs in matched passthrough messages and configured custom messages.
- `backend/internal/service/openai_gateway_service.go`: separates redacted WebSocket fallback client messages from internal Ops messages.
- `backend/internal/handler/gateway_handler.go`: redacts failover passthrough messages for the generic gateway.
- `backend/internal/handler/openai_gateway_handler.go`: redacts failover passthrough messages for the OpenAI gateway.
- `backend/internal/handler/gemini_v1beta_handler.go`: redacts failover passthrough messages for the Gemini native gateway.
- `backend/internal/service/error_passthrough_runtime_test.go`: verifies passthrough URL redaction.
- `backend/internal/service/openai_ws_fallback_test.go`: verifies fallback URL redaction and internal diagnostic preservation.
- `.gitignore`: allows the upstream error redaction document to be versioned under `docs/`.
- `docs/UPSTREAM_ERROR_URL_REDACTION.md`: documents the client-visible behavior and verification commands.
- `progress.md`: records implementation, verification, changed files, and rollback instructions.
- Rollback with `git revert <commit-containing-this-task>` after commit, or before commit restore only the files listed above from `HEAD` and remove `docs/UPSTREAM_ERROR_URL_REDACTION.md`.

## 2026-08-10 - Task: Replace redacted URL placeholders with customer-safe errors
### What was done
- Replaced client-visible URL placeholders and transport diagnostics with `Service temporarily unavailable, please retry later`.
- Added detection for common DNS, TCP, proxy, TLS, timeout, and routing failure text while preserving ordinary upstream business errors.
- Kept internal WebSocket fallback diagnostics separate from the customer response.

### Testing
- `cd backend && go test ./internal/service -run 'Test(ResolveOpenAIWSFallbackErrorResponse|ApplyErrorPassthroughRule|SanitizeClientUpstreamErrorMessage)' -count=1` passed.
- `cd backend && go test ./internal/handler -run 'Test(OpenAIGatewayHandler|GatewayHandler|.*Failover.*|.*Stream.*)' -count=1` passed.
- `git diff --check` passed.

### Notes
- `backend/internal/service/gemini_messages_compat_service.go`: classifies client-visible infrastructure errors and returns a stable service-unavailable message.
- `backend/internal/service/error_passthrough_runtime_test.go`: verifies network diagnostics are hidden and ordinary business errors remain unchanged.
- `backend/internal/service/openai_ws_fallback_test.go`: verifies fallback clients receive only the stable customer message.
- `docs/UPSTREAM_ERROR_URL_REDACTION.md`: documents the final customer-facing behavior.
- `progress.md`: records this refinement, verification, files, and rollback instructions.
- Rollback this refinement by reverting its commit, or before commit restore the four files listed above from the prior working-tree state.

## 2026-08-11 - Task: Exempt hjt13845049131 from OpenAI CY history filtering
### What was done
- Confirmed `hjt13845049131@163.com` is user ID `55` through the existing admin users API.
- Added a fixed user-ID exemption before CY marker lookup, so this user is not skipped even when the account toggle is enabled and a historical marker exists.
- Updated the Chinese and English account-toggle descriptions and the CY filter documentation.

### Testing
- `cd backend && go test ./internal/service -run 'Test(OpenAICyberPolicyUserFilter|ContentModerationService_MarkCyberPolicyUser|RecordCyberPolicyEvent_)' -count=1` passed.
- `cd frontend && pnpm exec vitest run src/i18n/__tests__/localesMessageCompile.spec.ts` passed.
- `cd frontend && pnpm exec eslint src/i18n/locales/zh/admin/accounts.ts src/i18n/locales/en/admin/accounts.ts` passed.
- `cd frontend && pnpm exec vue-tsc --noEmit` passed.
- `git diff --check` passed.

### Notes
- `backend/internal/service/openai_cyber_policy_user_filter.go`: exempts user ID `55` from the protected-account filter without an extra database lookup.
- `backend/internal/service/openai_cyber_policy_user_filter_test.go`: verifies the explicit exemption bypasses both cache and database marker reads.
- `frontend/src/i18n/locales/zh/admin/accounts.ts`: documents the Chinese exemption.
- `frontend/src/i18n/locales/en/admin/accounts.ts`: documents the English exemption.
- `docs/OPENAI_CYBER_POLICY_ACCOUNT_FILTER.md`: documents the fixed user-ID exemption.
- `progress.md`: records implementation, verification, changed files, and rollback instructions.
- Rollback this change by reverting its commit, or before commit remove the user-ID exemption and the matching test/documentation lines from the files listed above.

## 2026-08-10 - Task: Add account-level OpenAI CY user scheduling filter
### What was done
- Added a permanent database marker for regular users with historical upstream `cyber_policy` events and backfilled existing flagged moderation logs.
- Added an opt-in OpenAI account setting that skips marked users before scoring and during sticky, previous-response, WebSocket, model-routing, and retry rechecks while preserving normal failover behavior.
- Added trusted user-role context, admin exemption, Redis positive/negative marker caching, fail-open error handling, admin UI controls, tests, and documentation.

### Testing
- `cd backend && go test ./internal/service -run 'Test(OpenAICyberPolicyUserFilter|ContentModerationService_MarkCyberPolicyUser|RecordCyberPolicyEvent_Marks|RecordCyberPolicyEvent_DoesNotMark)' -count=1` passed.
- `cd backend && go test ./internal/service ./internal/repository ./internal/server/middleware -run '^$' -count=0` passed package compilation.
- `cd frontend && pnpm typecheck` passed.
- `cd frontend && pnpm exec vitest run src/components/account/__tests__/EditAccountModal.spec.ts --reporter=dot` passed.
- `git diff --check` passed.

### Notes
- `backend/migrations/221_cyber_policy_user_marks.sql`: creates and backfills permanent CY user markers.
- `backend/internal/repository/user_repo.go`: adds marker persistence queries.
- `backend/internal/repository/cyber_policy_user_marker_cache.go`: adds Redis marker cache operations.
- `backend/internal/pkg/ctxkey/ctxkey.go`: adds trusted request user-role context key.
- `backend/internal/server/middleware/api_key_auth.go`: stores authenticated role in request context.
- `backend/internal/service/openai_cyber_policy_user_filter.go`: defines the account setting, marker interfaces, cache lookup, and fail-open filtering.
- `backend/internal/service/content_moderation.go`: records a regular user marker on upstream CY events.
- `backend/internal/service/openai_gateway_scheduling.go`: filters protected accounts before scoring and final DB rechecks.
- `backend/internal/service/openai_account_scheduler.go`: initializes filtering for advanced and image scheduling entrypoints.
- `backend/internal/service/openai_ws_forwarder_support.go`: rejects protected previous-response bindings.
- `backend/internal/service/openai_cyber_policy_user_filter_test.go`: covers marking, admin/default behavior, cache/DB failures, and scheduler rechecks.
- `frontend/src/components/account/EditAccountModal.vue`: adds and persists the per-account toggle.
- `frontend/src/i18n/locales/en/admin/accounts.ts`: adds English labels.
- `frontend/src/i18n/locales/zh/admin/accounts.ts`: adds Chinese labels.
- `docs/OPENAI_CYBER_POLICY_ACCOUNT_FILTER.md`: documents semantics, migration, fail-open behavior, and rollback.
- Rollback before deployment by reverting the listed code/docs changes and excluding migration 221; after migration, disable account toggles first and use a reviewed follow-up migration to remove the marker table if required. Preserve unrelated `.superpowers/`.

## 2026-08-10 - Task: Harden CY filter fail-open behavior and verification
### What was done
- Made marker-cache read errors immediately fail open instead of changing availability based on a database fallback during a cache outage.
- Initialized the request-local filter state in legacy model selection and previous-response resolver entrypoints so all direct scheduler callers share one marker lookup.
- Added SQL repository coverage for regular-user insert, admin/no-op insert, and marker lookup, and allowed the new documentation file through the existing docs ignore policy.

### Testing
- `cd backend && go test ./internal/service -run 'Test(OpenAICyberPolicyUserFilter|ContentModerationService_MarkCyberPolicyUser|RecordCyberPolicyEvent_Marks|RecordCyberPolicyEvent_DoesNotMark)' -count=1` passed.
- `cd backend && go test ./internal/repository -run TestUserRepositoryCyberPolicyUserMarker -count=1` passed.
- `cd backend && go test ./internal/service -run 'TestOpenAIGatewayService_SelectAccountWithScheduler_(DefaultDisabledUsesLegacyLoadAwareness|EnabledUsesAdvancedPreviousResponseRouting|PreviousResponseSticky|SessionSticky|SessionStickyBusyKeepsSticky|StickyWeightedPreviousRequiresMovableContext)|TestOpenAIGatewayService_SelectAccountWithLoadAwareness_DBFreshGroupRecheckWaitsOnValidAccount' -count=1` passed.
- `cd frontend && pnpm typecheck` and the focused `EditAccountModal.spec.ts` Vitest run passed.
- `git diff --check` passed.

### Notes
- `backend/internal/service/openai_cyber_policy_user_filter.go`: fail-open cache-read handling and direct-entrypoint state support.
- `backend/internal/service/openai_gateway_scheduling.go`: request-local state for legacy model/load scheduling.
- `backend/internal/service/openai_ws_forwarder_support.go`: request-local state for direct previous-response resolution.
- `backend/internal/service/openai_cyber_policy_user_filter_test.go`: validates cache outage behavior and final marker semantics.
- `backend/internal/repository/cyber_policy_user_marker_test.go`: validates marker SQL operations with a SQL mock.
- `.gitignore`: allows the new behavior documentation to be versioned.
- Rollback this hardening by reverting the files listed above; preserve the earlier task implementation and unrelated `.superpowers/` unless the whole feature is being rolled back.

## 2026-08-10 - Task: Record CY marker dependency on global risk control
### What was done
- Documented that new CY user markers are not written while the global risk-control switch is disabled.
- Recorded the accepted operating assumption that global risk control remains enabled for the foreseeable future; no scheduling or risk-control logic was changed.

### Testing
- Verified the limitation text states the current behavior and the required precaution before disabling global risk control.
- `git diff --check` passed.

### Notes
- `docs/OPENAI_CYBER_POLICY_ACCOUNT_FILTER.md`: added the known limitation and the operational prerequisite.
- `progress.md`: recorded the accepted limitation; this task made no code changes.
- Rollback by reverting this documentation-only entry and the corresponding known-limitation section; no runtime rollback is required.
## 2026-08-10 - Task: Enable Leo upstream model synchronization
### What was done
- Added Leo as an upstream-sync platform in the administrator model whitelist control.
- Added a Leo API-key `/v1/models` request path in Sub2, using the existing URL policy, proxy/TLS client, account header overrides, and redacted error handling.
- Added backend and frontend regression coverage plus an operator document; this change is prepared locally and has not been pushed or deployed.

### Testing
- `cd backend && go test ./internal/service -run 'Test(BuildUpstreamModels|FetchUpstreamSupportedModels)' -count=1` passed.
- `cd backend && go test ./internal/service -count=1` passed.
- `cd frontend && corepack pnpm@10.28.1 exec vitest run src/components/account/__tests__/ModelWhitelistSelector.spec.ts` passed (3 tests).
- `git diff --check` passed.

### Notes
- `backend/internal/service/upstream_models.go`: routes Leo accounts to a validated `/v1/models` request with Bearer authentication.
- `backend/internal/service/upstream_models_test.go`: verifies Leo request construction and OpenAI-compatible model response parsing.
- `frontend/src/components/account/ModelWhitelistSelector.vue`: enables the Leo sync action.
- `frontend/src/components/account/__tests__/ModelWhitelistSelector.spec.ts`: verifies the Leo sync action is visible.
- `docs/LEO_MODEL_SYNC.md`: documents the endpoint contract, security boundaries, verification, and rollback.
- `progress.md`: records this local implementation and verification.
- Rollback before commit: restore the four modified source/test files from `HEAD` and remove `docs/LEO_MODEL_SYNC.md`; no production service or database has been changed.

## 2026-08-10 - Task: Push Leo upstream model synchronization patch
### What was done
- Committed the verified Leo model synchronization change as `1cc9606` (`fix: sync Leo upstream models`).
- Pushed it to the remote branch `fix/leo-model-sync` only; no server, service, database, or deployment action was performed.

### Testing
- Confirmed the remote branch resolves to commit `1cc9606`.
- Confirmed the local working tree was clean after the code push.

### Notes
- `progress.md`: records the Git push and the deliberate no-deployment boundary.
- Rollback: delete the remote branch with `git push origin --delete fix/leo-model-sync` and revert commit `1cc9606` wherever it is merged; production remains unchanged until a separate deployment.

## 2026-08-11 - Task: Merge Leo model synchronization into the release branch
### What was done
- Reviewed `fix/leo-model-sync` and merged it into `codex/leo-video-channel` instead of releasing the branch directly, preserving the CY account protection introduced in `v0.1.173-fy.3`.
- Kept Leo model discovery limited to API-key accounts and the existing validated upstream request path.
- Resolved the append-only progress log conflict by retaining both branches' complete histories.

### Testing
- `cd backend && go test ./... -count=1` passed.
- `cd backend && go vet ./...` passed.
- `cd frontend && node node_modules/vitest/vitest.mjs run src/components/account/__tests__/ModelWhitelistSelector.spec.ts --reporter=default --maxWorkers=2 --minWorkers=1` passed (3 tests).
- `cd frontend && pnpm.exe build` passed; Vite transformed 1031 modules.

### Notes
- `backend/internal/service/upstream_models.go`: adds validated Leo `/v1/models` request construction with Bearer authentication.
- `backend/internal/service/upstream_models_test.go`: covers Leo request construction and OpenAI-compatible model parsing.
- `frontend/src/components/account/ModelWhitelistSelector.vue`: enables upstream model synchronization for Leo accounts.
- `frontend/src/components/account/__tests__/ModelWhitelistSelector.spec.ts`: verifies the Leo synchronization action is available.
- `docs/LEO_MODEL_SYNC.md`: documents behavior, security boundaries, verification, and rollback.
- `progress.md`: preserves both branch histories and records this release integration.
- Rollback after publication with `git revert -m 1 v0.1.173-fy.4`, then push the revert and publish a follow-up tag.

## 2026-08-11 - Task: Preserve Leo model names during Sub2 synchronization
### What was done
- Kept the existing `models: string[]` response for compatibility and added `model_details` entries containing each model ID and display name.
- Preserved model names while parsing OpenAI-compatible, Grok, and Leo model-list responses.
- Updated the account model selector to display `name (UUID)` while continuing to save and submit the UUID.
- Documented the catalog contract and fallback behavior for upstreams without names.

### Testing
- `cd backend && go test ./internal/service -run 'TestExtractUpstreamModelDetailsPreservesDisplayName|TestExtractUpstreamModelIDs|TestFetchUpstream' -count=1`: passed.
- `cd backend && go test ./internal/handler/admin -run 'TestAccountHandlerSyncUpstreamModels' -count=1`: passed.
- `cd frontend && node node_modules/vitest/vitest.mjs run src/components/account/__tests__/ModelWhitelistSelector.spec.ts --reporter=default --maxWorkers=2 --minWorkers=1`: passed (3 tests).
- `cd frontend && npm run build`: passed (`vue-tsc` and Vite production build).

### Notes
- `backend/internal/service/upstream_models.go`: adds named model detail parsing without changing the ID-only service wrapper.
- `backend/internal/service/upstream_models_test.go`: verifies display-name extraction and ID compatibility.
- `backend/internal/handler/admin/account_handler.go`: returns `model_details` alongside the legacy model list.
- `backend/internal/handler/admin/account_handler_available_models_test.go`: verifies IDs and names in the sync response.
- `frontend/src/api/admin/accounts.ts`: types the optional model detail payload.
- `frontend/src/components/account/ModelWhitelistSelector.vue`: renders names while preserving UUID values.
- `docs/LEO_MODEL_CATALOG.md`: documents the Sub2 response contract and endpoint behavior.
- `progress.md`: appends this implementation and verification record.
- Rollback: revert the task commit; existing account model whitelist strings require no data migration.


## 2026-08-13 - Task: Fix OpenAI pool reasoning rendering and merge upstream v0.1.175
### What was done
- Confirmed upstream release `v0.1.175` contains PR #5304 / commit `8aa425d22`, which accepts both Chat Completions `reasoning_content` and `reasoning` fields and converts them to standard Responses reasoning events.
- Merged upstream `v0.1.175` into the fork while retaining the fork's CY account filtering, administrator audit exemption, CY input summaries, and Leo model catalog behavior.
- Resolved four merge conflicts by composing both branches' required behavior; no database schema or deployment configuration changed.
- Added an operator document for the reasoning compatibility path. Upstream commits after the `v0.1.175` tag were deliberately excluded because they are not part of a formal upstream release.

### Testing
- `cd backend && go test ./internal/pkg/apicompat -run 'TestChatReasoningAlias|TestChatCompletions.*Reasoning|TestResponsesToChatCompletions.*Reasoning' -count=1` passed.
- `cd backend && go test ./internal/handler -run 'TestRunSecurityAudit|TestCachesSecurityAuditCompletion|TestOpenAIResponsesWebSocket' -count=1` passed.
- `cd backend && go test ./internal/service -run 'TestOpenAICyberPolicyUserFilter|TestOpenAIPoolMode|TestOpenAI.*Responses.*ChatCompletions|TestForwardResponses.*ChatCompletions' -count=1` passed.
- `cd backend && go test ./internal/service -run 'TestExtractUpstreamModelDetailsPreservesDisplayName|TestExtractUpstreamModelIDs|TestExtractGrokUpstreamModelIDs' -count=1` passed.
- `cd backend && go vet ./...` passed.
- `cd frontend && pnpm vitest run` passed, and `pnpm build` passed with Vite transforming 1031 modules.
- `cd backend && go test -p 1 ./... -count=1` passed. One earlier parallel full-suite run showed a load-sensitive server-timing test failure; that test passed 10 isolated repetitions and the complete serial suite.
- `git diff --cached --check` is required immediately before commit.
- A real production upstream replay was not available; protocol parity is established by deterministic stream and non-stream conversion fixtures.

### Notes
- `.gitignore`: allows the OpenAI reasoning compatibility document to be tracked.
- `DEV_GUIDE.md`: merges the upstream v0.1.175 development guidance update.
- `README.md`: merges the upstream v0.1.175 release documentation update.
- `README_CN.md`: merges the upstream v0.1.175 Chinese release documentation update.
- `README_JA.md`: merges the upstream v0.1.175 Japanese release documentation update.
- `backend/internal/handler/admin/backup_handler.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/handler/admin/channel_handler.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/handler/api_key_handler.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/handler/api_key_handler_validation_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/handler/failover_loop.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/handler/failover_loop_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/handler/openai_chat_completions.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/handler/openai_gateway_handler.go`: combines upstream WebSocket turn numbering with the fork's CY input-summary extraction.
- `backend/internal/handler/openai_gateway_handler_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/handler/security_audit_helper.go`: combines upstream per-turn WebSocket audit deduplication with the fork's administrator audit exemption.
- `backend/internal/handler/security_audit_helper_test.go`: retains both branches' audit tests and uses the shared evaluation counter.
- `backend/internal/pkg/apicompat/chatcompletions_anthropic_bridge.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/pkg/apicompat/chatcompletions_reasoning_alias_test.go`: covers stream and non-stream reasoning aliases plus field precedence.
- `backend/internal/pkg/apicompat/chatcompletions_responses_bridge.go`: accepts Chat Completions reasoning aliases and emits standard Responses reasoning events.
- `backend/internal/pkg/apicompat/testdata/issue5302/nonstream_reasoning.json`: adds the upstream regression fixture for Chat Completions reasoning alias conversion.
- `backend/internal/pkg/apicompat/testdata/issue5302/reasoning_content_precedence.json`: adds the upstream regression fixture for Chat Completions reasoning alias conversion.
- `backend/internal/pkg/apicompat/testdata/issue5302/stream_reasoning.json`: adds the upstream regression fixture for Chat Completions reasoning alias conversion.
- `backend/internal/pkg/apicompat/types.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/repository/backup_s3_store.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/repository/backup_s3_store_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/account_scheduling_threshold_eval.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/account_scheduling_threshold_eval_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/account_stats_pricing.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/account_stats_pricing_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/api_key_service.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/api_key_service_validation_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/backup_archive.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/backup_archive_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/backup_service.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/backup_service_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/billing_service.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/channel.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/channel_service.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/content_moderation.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/content_moderation_cyber_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/gateway_channel_restriction_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/gateway_forward.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/gateway_usage_billing.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/gemini_messages_compat_service.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/gemini_messages_compat_service_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/identity_service.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/identity_service_user_agent_validation_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_codex_fingerprint.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_codex_fingerprint_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_first_output_timeout_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_gateway_apikey_item_id_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_gateway_chat_completions.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_gateway_chat_completions_raw.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_gateway_chat_completions_raw_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_gateway_forward.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_gateway_grok_chat_bridge_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_gateway_grok_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_gateway_passthrough.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_gateway_response_handling.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_gateway_responses_empty_completed_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_gateway_scheduling.go`: combines upstream compatible-account eligibility checks with the fork's CY user account filter.
- `backend/internal/service/openai_gateway_service_codex_cli_only_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_gateway_service_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_gateway_upstream_errors.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_gateway_usage.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_gateway_usage_integrity.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_gateway_usage_integrity_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_images_responses.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_images_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_oauth_passthrough_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_oauth_service.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_privacy_service.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_profit_control_legacy_diagnostics_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_responses_item_id.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_silent_refusal.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_stream_read_error.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_subscription_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_upstream_client_error.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_upstream_client_error_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_visible_ttft_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_ws_v2/passthrough_relay.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/openai_ws_v2/passthrough_relay_internal_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/openai_ws_v2/passthrough_relay_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/pricing_service.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/ratelimit_service.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/ratelimit_service_403_html_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/response_model_billing_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/setting_features.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `backend/internal/service/setting_service_platform_threshold_test.go`: merges upstream v0.1.175 regression coverage into the fork release set.
- `backend/internal/service/upstream_response_model.go`: merges the upstream v0.1.175 backend behavior into the fork release set.
- `docs/OPENAI_REASONING_ALIAS_COMPAT.md`: documents the affected fallback path, protocol behavior, verification, and rollback.
- `frontend/src/api/admin/backup.ts`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/components/account/BulkEditAccountModal.vue`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/components/account/CreateAccountModal.vue`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/components/account/EditAccountModal.vue`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/components/admin/usage/UsageTable.vue`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/components/admin/usage/__tests__/UsageTable.spec.ts`: merges upstream v0.1.175 frontend regression coverage into the fork release set.
- `frontend/src/components/layout/AppSidebar.vue`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/constants/channel.ts`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/i18n/locales/en/admin/accounts.ts`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/i18n/locales/en/admin/channels.ts`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/i18n/locales/en/admin/overview.ts`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/i18n/locales/zh/admin/accounts.ts`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/i18n/locales/zh/admin/channels.ts`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/i18n/locales/zh/admin/overview.ts`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/views/admin/BackupView.vue`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/views/admin/ChannelsView.vue`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/views/admin/UsageView.vue`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/views/admin/__tests__/BackupView.spec.ts`: merges upstream v0.1.175 frontend regression coverage into the fork release set.
- `frontend/src/views/admin/__tests__/UsageView.spec.ts`: merges upstream v0.1.175 frontend regression coverage into the fork release set.
- `frontend/src/views/admin/__tests__/groupsImagePricing.spec.ts`: merges upstream v0.1.175 frontend regression coverage into the fork release set.
- `frontend/src/views/admin/groupsImagePricing.ts`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/views/admin/ops/components/OpsDashboardHeader.vue`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `frontend/src/views/admin/ops/utils/__tests__/opsFormatters.spec.ts`: merges upstream v0.1.175 frontend regression coverage into the fork release set.
- `frontend/src/views/admin/ops/utils/opsFormatters.ts`: merges the upstream v0.1.175 frontend behavior into the fork release set.
- `progress.md`: records the merge, conflict resolutions, test evidence, changed-file inventory, and rollback point.
- Rollback after publication with `git revert -m 1 <merge-commit>`, then push the revert and publish a follow-up fork tag; no database rollback is required. Preserve unrelated `.superpowers/` content.

## 2026-08-13 - Task: Publish OpenAI reasoning compatibility release
### What was done
- Pushed merge commit `f61031bea4a1a0061ac2f43cdedeaf54fb96c407` to `origin/codex/leo-video-channel`.
- Published tag and GitHub Release `v0.1.175-fy.1` from the verified merge commit.
- Confirmed the public Linux amd64 archive and checksum manifest are downloadable.

### Testing
- Both public Release asset URLs returned HTTP `200`.
- Downloaded `sub2api_0.1.175-fy.1_linux_amd64.tar.gz` and verified its SHA-256 is `320961c19a19e0a6d4ca840cf687dfd647b06ac0b49fc45d530b9978d5dd81db`, exactly matching `checksums.txt`.
- Verified the archive is `37,425,206` bytes and contains the single expected `sub2api` entry.
- Confirmed remote tag `v0.1.175-fy.1` resolves to merge commit `f61031bea4a1a0061ac2f43cdedeaf54fb96c407`.

### Notes
- `progress.md`: records the pushed commit, published Release, asset availability, checksum, archive content, and rollback point.
- The Release does not require a database migration or configuration change.
- Rollback source behavior with `git revert -m 1 f61031bea4a1a0061ac2f43cdedeaf54fb96c407` and publish a follow-up release; withdraw this release only after removing the GitHub Release and deleting remote tag `v0.1.175-fy.1`. Preserve unrelated `.superpowers/` content.

## 2026-08-13 - Task: Merge upstream v0.1.176 into the fork
### What was done
- Merged the formal upstream `v0.1.176` tag while preserving the fork's Leo video pricing, CY account filtering, administrator audit exemption, Leo model catalog, and other fork-only behavior.
- Resolved seven conflicts across group validation, Grok pricing, channel video pricing, usage billing, and the pricing editor; Leo video pricing remains model-specific while upstream generic group video pricing remains available to other platforms.
- Added upstream Grok 4.6 and JWT subscription-tier support, group model pricing and long-context controls, native `/x_search`, and the upstream billing, backup, cache, and account-usage fixes included in the tag.
- Documented the additive PostgreSQL migration and retained both complete `221_*` migration filenames so deployed migration history and checksums remain stable.

### Testing
- `cd backend && go test -tags=unit ./internal/service -run '^TestValidatePricingBillingMode$' -count=1`: passed, including Leo strict video pricing and generic video default pricing.
- `cd frontend && node node_modules/vitest/vitest.mjs run src/components/admin/channel/__tests__/PricingEntryCard.spec.ts --reporter=default --maxWorkers=2 --minWorkers=1`: passed (5 tests).
- `cd backend && go test -p 1 ./... -count=1`: passed.
- `cd frontend && node node_modules/vitest/vitest.mjs run --reporter=default --maxWorkers=2 --minWorkers=1`: passed; existing i18n and jsdom warnings remained non-fatal.
- `cd backend && go vet ./...`: passed.
- `cd frontend && npm run build`: passed (`vue-tsc` and Vite, 1031 modules); existing chunk warnings remained non-fatal.
- `git diff --cached --check` and conflict-marker checks are required immediately before the merge commit.

### Notes
- `.gitignore`: allows the v0.1.176 group pricing deployment document to be tracked.
- `README.md`: merges upstream v0.1.176 partner and release information.
- `README_CN.md`: merges upstream v0.1.176 Chinese partner and release information.
- `README_JA.md`: merges upstream v0.1.176 Japanese partner and release information.
- `assets/partners/logos/duckip.png`: adds the upstream DuckIP partner logo.
- `assets/partners/logos/haoai.svg`: removes the upstream-retired HaoAI partner logo.
- `assets/partners/logos/swiftprox.png`: adds the upstream SwiftProx partner logo.
- `backend/cmd/server/VERSION`: merges the version metadata present at the formal upstream tag; the fork Release workflow injects the fork tag version during packaging.
- `backend/cmd/server/wire_gen.go`: wires the upstream group pricing resolver dependencies.
- `backend/ent/group.go`: adds generated group model pricing and long-context fields.
- `backend/ent/group/group.go`: adds generated group field descriptors.
- `backend/ent/group/where.go`: adds generated group pricing predicates.
- `backend/ent/group_create.go`: adds generated create setters for group pricing fields.
- `backend/ent/group_update.go`: adds generated update setters for group pricing fields.
- `backend/ent/migrate/schema.go`: adds the group pricing columns to the generated schema.
- `backend/ent/mutation.go`: adds generated mutation support for group pricing fields.
- `backend/ent/runtime/runtime.go`: refreshes generated runtime schema bindings.
- `backend/ent/schema/group.go`: defines group model pricing and the default-enabled long-context switch.
- `backend/internal/handler/admin/group_handler.go`: retains Leo platform validation and exposes group pricing fields.
- `backend/internal/handler/dto/mappers.go`: maps new upstream account usage metadata.
- `backend/internal/handler/dto/types.go`: exposes new upstream account usage fields.
- `backend/internal/handler/gateway_handler.go`: merges upstream Grok gateway behavior.
- `backend/internal/handler/gateway_web_search.go`: composes shared web and x_search request handling.
- `backend/internal/handler/grok_audio.go`: merges upstream Grok audio billing behavior.
- `backend/internal/handler/grok_audio_billing_test.go`: covers upstream Grok audio billing changes.
- `backend/internal/handler/openai_x_search.go`: adds the native `/x_search` handler.
- `backend/internal/handler/openai_x_search_test.go`: covers native x_search request handling.
- `backend/internal/pkg/apicompat/chatcompletions_responses_bridge.go`: preserves x_search fields in Chat-to-Responses conversion.
- `backend/internal/pkg/apicompat/chatcompletions_responses_bridge_custom_tools_test.go`: expands custom-tool conversion coverage.
- `backend/internal/pkg/apicompat/chatcompletions_to_responses.go`: maps x_search filtering and tool choice.
- `backend/internal/pkg/apicompat/chatcompletions_x_search_test.go`: covers x_search protocol conversion.
- `backend/internal/pkg/apicompat/types.go`: adds x_search-compatible protocol fields.
- `backend/internal/pkg/xai/models.go`: adds Grok 4.6 model identifiers.
- `backend/internal/pkg/xai/models_test.go`: covers Grok 4.6 model recognition.
- `backend/internal/pkg/xai/oauth_test.go`: updates upstream OAuth fixtures.
- `backend/internal/pkg/xai/quota.go`: exposes JWT-derived subscription tiers.
- `backend/internal/pkg/xai/subscription_tier.go`: implements Grok JWT subscription-tier detection.
- `backend/internal/pkg/xai/subscription_tier_test.go`: covers Grok subscription-tier claims and fallbacks.
- `backend/internal/repository/api_key_repo.go`: preloads group pricing for API-key billing resolution.
- `backend/internal/repository/group_repo.go`: persists group pricing and long-context controls.
- `backend/internal/server/api_contract_test.go`: updates the gateway API contract for x_search.
- `backend/internal/server/routes/gateway.go`: registers the native x_search route.
- `backend/internal/server/routes/prompt_audit_route_coverage_test.go`: covers x_search prompt-audit routing.
- `backend/internal/service/account_test_service.go`: merges upstream Grok account test behavior.
- `backend/internal/service/admin_group.go`: validates and normalizes per-group model pricing.
- `backend/internal/service/admin_group_platform_cache_test.go`: covers channel cache invalidation after group platform changes.
- `backend/internal/service/admin_service.go`: invalidates affected channel caches after group changes.
- `backend/internal/service/backup_service.go`: adds the upstream scheduled-backup leader lock.
- `backend/internal/service/backup_service_test.go`: covers single-leader scheduled backups.
- `backend/internal/service/billing_search_audio_cost_test.go`: updates search and audio billing expectations.
- `backend/internal/service/billing_service.go`: uses upstream Grok 4.3/4.5/4.6 prices and 200k long-context rules without duplicate fork pricing.
- `backend/internal/service/billing_service_test.go`: covers new Grok and long-context pricing behavior.
- `backend/internal/service/channel.go`: retains video non-token interval semantics while merging upstream normalization.
- `backend/internal/service/channel_service.go`: preserves strict Leo video validation, allows generic video pricing elsewhere, and merges cache and normalization fixes.
- `backend/internal/service/channel_service_test.go`: covers pricing-key normalization plus explicit Leo and generic video validation boundaries.
- `backend/internal/service/gateway_usage_billing.go`: resolves group pricing before channel and built-in pricing.
- `backend/internal/service/grok_audio.go`: merges upstream Grok audio usage and pricing behavior.
- `backend/internal/service/grok_audio_test.go`: covers Grok audio usage changes.
- `backend/internal/service/grok_media.go`: merges upstream Grok media billing updates.
- `backend/internal/service/grok_model_quota_block.go`: adds model-scoped Grok quota blocking.
- `backend/internal/service/grok_oauth_service.go`: refreshes JWT-derived Grok subscription tiers.
- `backend/internal/service/grok_oauth_service_test.go`: covers tier refresh and stale-tier replacement.
- `backend/internal/service/grok_quota_fetcher.go`: derives Grok capacity and tiers from current account data.
- `backend/internal/service/grok_quota_fetcher_test.go`: covers tier and quota edge cases.
- `backend/internal/service/grok_quota_service.go`: carries current Grok tier metadata.
- `backend/internal/service/grok_upstream_failure.go`: classifies new Grok model quota failures.
- `backend/internal/service/grok_upstream_failure_test.go`: covers model-scoped Grok failures.
- `backend/internal/service/group.go`: adds group model pricing and long-context domain fields.
- `backend/internal/service/model_pricing_resolver.go`: implements Group to Channel to built-in pricing precedence.
- `backend/internal/service/model_pricing_resolver_test.go`: covers group pricing resolution and long-context selection.
- `backend/internal/service/openai_apikey_responses_probe.go`: avoids unsupported verdicts for incomplete probe responses.
- `backend/internal/service/openai_apikey_responses_probe_verdict_test.go`: covers complete, truncated, and failed probe verdicts.
- `backend/internal/service/openai_gateway_chat_completions_raw.go`: preserves x_search response sources.
- `backend/internal/service/openai_gateway_grok.go`: routes Grok 4.6 and x_search requests.
- `backend/internal/service/openai_gateway_grok_cache.go`: updates Grok cache handling for current model behavior.
- `backend/internal/service/openai_gateway_grok_chat_bridge.go`: preserves x_search fields through the Grok bridge.
- `backend/internal/service/openai_gateway_grok_chat_bridge_test.go`: covers Grok x_search bridge behavior.
- `backend/internal/service/openai_gateway_grok_test.go`: updates Grok gateway regression coverage.
- `backend/internal/service/openai_gateway_messages.go`: merges upstream message conversion behavior.
- `backend/internal/service/openai_gateway_record_usage_test.go`: covers group media pricing and new usage paths.
- `backend/internal/service/openai_gateway_usage.go`: composes group media pricing with retained Leo model, resolution, and duration normalization.
- `backend/internal/service/openai_ws_http_bridge.go`: merges upstream WebSocket HTTP bridge behavior.
- `backend/internal/service/wire.go`: registers the group model pricing resolver.
- `backend/migrations/221_group_model_pricing.sql`: adds and backfills additive group pricing columns; the existing `221_cyber_policy_user_marks.sql` remains unchanged.
- `docs/GROUP_MODEL_PRICING_V176.md`: documents startup migration, compatibility, verification, and rollback boundaries.
- `frontend/src/components/account/AccountUsageCell.vue`: displays current Grok subscription-tier usage data.
- `frontend/src/components/account/__tests__/AccountUsageCell.spec.ts`: covers current-tier account usage rendering.
- `frontend/src/components/admin/channel/PricingEntryCard.vue`: retains the Leo resolution editor and adds the upstream generic group video editor.
- `frontend/src/components/admin/channel/__tests__/PricingEntryCard.spec.ts`: distinguishes Leo fixed-resolution pricing from generic group video pricing.
- `frontend/src/components/common/PlatformTypeBadge.vue`: displays Grok subscription-tier badges.
- `frontend/src/components/common/__tests__/PlatformTypeBadge.grok.spec.ts`: covers new Grok tier badges.
- `frontend/src/composables/__tests__/useModelWhitelist.spec.ts`: covers Grok 4.6 whitelist options.
- `frontend/src/composables/useModelWhitelist.ts`: adds Grok 4.6 to model selection.
- `frontend/src/i18n/locales/en/admin/channels.ts`: adds English group video pricing labels.
- `frontend/src/i18n/locales/en/admin/overview.ts`: adds English Grok tier labels.
- `frontend/src/i18n/locales/zh/admin/channels.ts`: adds Chinese group video pricing labels.
- `frontend/src/i18n/locales/zh/admin/overview.ts`: adds Chinese Grok tier labels.
- `frontend/src/types/index.ts`: types new group and Grok usage fields.
- `frontend/src/utils/__tests__/accountUsageRefresh.spec.ts`: covers tier-aware incremental account usage refresh.
- `frontend/src/utils/accountUsageRefresh.ts`: refreshes account rows when Grok tier snapshots change.
- `frontend/src/views/admin/AccountsView.vue`: shows current Grok tier details in account administration.
- `frontend/src/views/admin/GroupsView.vue`: adds group model pricing and long-context controls while retaining the Leo-specific video editor.
- `frontend/src/views/admin/__tests__/AccountsView.sparkShadow.spec.ts`: covers updated account-tier display and refresh behavior.
- `progress.md`: records the merge scope, verification evidence, changed-file inventory, migration impact, and rollback point.
- Rollback source behavior after commit with `git revert -m 1 <merge-commit>`, push the revert, and publish a follow-up fork tag. Existing deployments may leave the additive group columns and `schema_migrations` row in place because the previous binary ignores them; do not drop the columns unless their data is backed up and a separate database change is explicitly approved. Preserve unrelated `.superpowers/` content.

## 2026-08-13 - Task: Publish v0.1.176-fy.1
### What was done
- Pushed merge commit `dc5188088bc6e7f30477a611bbd60c0ec7743bad` to `origin/codex/leo-video-channel`.
- Published annotated tag and GitHub Release `v0.1.176-fy.1` from that verified merge commit.
- Confirmed the Release workflow compiled, packaged, and published the Linux amd64 binary successfully.

### Testing
- GitHub Actions run `31662265668` completed with conclusion `success`; Linux build, package, and Release publication steps all passed.
- Public downloads for `sub2api_0.1.176-fy.1_linux_amd64.tar.gz` and `checksums.txt` both returned HTTP `200`.
- The archive is `37,464,294` bytes and its SHA-256 is `ec52f7085645a7d28a1956f43cb01ffd61074b24bbefd3932af51c5d222e0e43`, exactly matching `checksums.txt`.
- `tar -tzf` confirmed the archive contains the single expected `sub2api` entry.
- The remote annotated tag dereferences to `dc5188088bc6e7f30477a611bbd60c0ec7743bad`; the Release is public, non-draft, and non-prerelease.

### Notes
- `progress.md`: records the pushed merge commit, Release workflow result, public asset availability, checksum, archive content, and rollback point.
- Deploying this Release runs `221_group_model_pricing.sql` during startup; existing group long-context pricing remains enabled by default.
- Rollback source behavior with `git revert -m 1 dc5188088bc6e7f30477a611bbd60c0ec7743bad` and publish a follow-up release. Leave the additive group columns in place for binary rollback; withdraw this Release only after removing GitHub Release `v0.1.176-fy.1` and deleting its remote tag. Preserve unrelated `.superpowers/` content.

## 2026-08-14 - Task: Synchronize Leo image models and Seedance 2.5 capabilities
### What was done
- Changed Leo channel pricing synchronization to read the live authenticated model catalog, classify known video models separately from dynamic image models, remove raw UUID entries, and create one unsaved pricing row per newly discovered model so administrators only need to enter prices and save.
- Added Leo image-generation routing through the existing OpenAI-compatible image endpoint, including API-key authentication, configured base URL handling, dynamic public model names, reference-image fields, `n` task splitting through the existing account scheduler, and an explicit platform-neutral rejection for unsupported image edits.
- Added both Seedance 2.5 public model IDs with aligned 480p/720p, 4-30 second, aspect-ratio, media-reference, generated-audio, and per-resolution pricing constraints across server validation, channel pricing, the video page, and the customer API documentation.
- Kept future image-model additions catalog-driven; future video models still require a local capability specification before they can be exposed as video pricing entries.

### Testing
- `go test ./internal/service ./internal/handler ./internal/handler/admin`: passed.
- `go test ./internal/service -run "(SelectAccountWithSchedulerForImages|FetchUpstreamSupportedModels|Leo|OpenAIImages|VideoBillingResolution)"`: passed after the final model de-duplication and scheduler signature changes.
- `go test -tags unit ./internal/handler/admin -run "SyncPricingModels"`: passed.
- `go test ./cmd/server -run "^$"`: passed, confirming final dependency injection and server assembly compile.
- `pnpm.cmd exec vitest run src/components/admin/channel/__tests__/types.spec.ts src/views/user/__tests__/VideoGenerationView.spec.ts src/views/user/__tests__/VideoApiDocsView.spec.ts`: passed 50 tests; the final pricing-classification adjustment was rechecked with all 18 channel type tests passing.
- `pnpm.cmd run typecheck`: passed on the final TypeScript state.
- `pnpm.cmd run build`: passed; Vite transformed 1031 modules and completed the production build with only pre-existing chunking warnings.
- `git diff --check`: passed. Customer video pages and localized API documentation contain none of `LeoStudio`, `Leonardo`, `UUID`, `upstream_job_id`, or `account_id`.

### Notes
- `backend/cmd/server/wire_gen.go`: injects the account repository and account-test service into the channel handler.
- `backend/internal/handler/admin/account_handler_available_models_test.go`: verifies public image names replace raw UUIDs and carry media kinds.
- `backend/internal/handler/admin/channel_handler.go`: fetches the live Leo catalog for channel pricing and falls back to the local video catalog when no live account succeeds.
- `backend/internal/handler/admin/channel_handler_test.go`: verifies live image/video catalog synchronization and classification.
- `backend/internal/handler/gateway_handler.go`: includes the latest Leo video IDs in gateway model-list fallbacks.
- `backend/internal/handler/openai_images.go`: selects Leo image accounts and returns platform-neutral client errors.
- `backend/internal/handler/openai_images_split.go`: preserves Leo platform routing for every split `n` task.
- `backend/internal/server/routes/gateway.go`: enables the image-generation route for Leo groups.
- `backend/internal/service/account.go`: allows live-catalog Leo models through stale account mappings and enables Leo API-key image capability.
- `backend/internal/service/leo_account.go`: registers both Seedance 2.5 public IDs and generalizes the account mapping validation message.
- `backend/internal/service/leo_account_test.go`: verifies Seedance 2.5 appears in the default Leo model candidates.
- `backend/internal/service/leo_video_model_specs.go`: defines the Seedance 2.5 server-side capability contract.
- `backend/internal/service/leo_video_model_specs_test.go`: verifies accepted and rejected Seedance 2.5 parameter combinations.
- `backend/internal/service/openai_account_scheduler.go`: routes image selection to the required Leo platform while retaining existing account-pool scheduling.
- `backend/internal/service/openai_images.go`: forwards Leo image requests with the configured Bearer key/base URL and rejects unsupported edits before upstream access.
- `backend/internal/service/openai_images_test.go`: verifies Leo image capability, request forwarding, reference fields, `n`, authentication, and edit rejection.
- `backend/internal/service/upstream_models.go`: adds media kinds, public-name replacement, UUID filtering, and case-insensitive public model de-duplication.
- `backend/internal/service/upstream_models_test.go`: verifies mixed legacy/current catalog entries become a de-duplicated public image/video catalog.
- `backend/internal/service/video_billing_resolution.go`: restricts Seedance 2.5 pricing to 480p and 720p.
- `backend/internal/service/video_billing_resolution_test.go`: verifies both Seedance 2.5 aliases use the correct price tiers.
- `docs/LEO_MODEL_SYNC.md`: documents live catalog synchronization, automatic image pricing rows, and the remaining explicit-spec rule for new video models.
- `docs/LEO_VIDEO_MODEL_SPECS.md`: documents the Seedance 2.5 parameter and media-reference limits.
- `frontend/src/api/admin/accounts.ts`: adds media-kind metadata to synchronized account models.
- `frontend/src/api/admin/channels.ts`: exposes channel model-detail metadata to the pricing form.
- `frontend/src/components/account/CreateAccountModal.vue`: adds both Seedance 2.5 aliases to new Leo account mappings.
- `frontend/src/components/admin/channel/__tests__/types.spec.ts`: verifies Seedance 2.5 price tiers/order and dynamic image/video price-row classification.
- `frontend/src/components/admin/channel/types.ts`: creates video or image pricing rows from synchronized model kinds.
- `frontend/src/constants/channel.ts`: registers both Seedance 2.5 public IDs in the frontend video catalog.
- `frontend/src/i18n/locales/en/dashboard.ts`: adds customer-facing English Seedance 2.5 constraints without internal identifiers.
- `frontend/src/i18n/locales/zh/dashboard.ts`: adds customer-facing Chinese Seedance 2.5 constraints without internal identifiers.
- `frontend/src/views/admin/ChannelsView.vue`: passes synchronized model metadata into automatic pricing-row creation.
- `frontend/src/views/user/VideoApiDocsView.vue`: adds both Seedance 2.5 aliases to the customer parameter matrix and examples.
- `frontend/src/views/user/VideoGenerationView.vue`: limits Seedance 2.5 controls to its supported resolution, duration, ratio, and media combinations.
- `frontend/src/views/user/__tests__/VideoApiDocsView.spec.ts`: verifies the complete model list and blocks internal identifiers in rendered documentation.
- `frontend/src/views/user/__tests__/VideoGenerationView.spec.ts`: verifies Seedance 2.5 options on the generation page.
- `progress.md`: records this implementation, verification evidence, changed-file inventory, and rollback point.
- Rollback point: source baseline `9b329d435aa18fa48be5aadc32166f83fa8f8e55`. After this task is committed, use `git revert <task-commit>`; preserve the pre-existing `docs/LEO_VIDEO_CHANNEL.md`, earlier `progress.md` history, and `.superpowers/` content.
