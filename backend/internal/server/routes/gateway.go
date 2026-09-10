package routes

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/handler"
	pkghttputil "github.com/Wei-Shaw/sub2api/internal/pkg/httputil"
	"github.com/Wei-Shaw/sub2api/internal/pkg/requestmodel"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// RegisterGatewayRoutes 注册 API 网关路由（Claude/OpenAI/Gemini 兼容）
func RegisterGatewayRoutes(
	r *gin.Engine,
	h *handler.Handlers,
	apiKeyAuth middleware.APIKeyAuthMiddleware,
	apiKeyService *service.APIKeyService,
	subscriptionService *service.SubscriptionService,
	opsService *service.OpsService,
	settingService *service.SettingService,
	compositeResolver *service.CompositeRouteResolver,
	cfg *config.Config,
) {
	bodyLimit := middleware.RequestBodyLimit(cfg.Gateway.MaxBodySize)
	textBodyLimit := middleware.RequestBodyLimit(cfg.Gateway.TextMaxBodySize)
	clientRequestID := middleware.ClientRequestID()
	opsErrorLogger := handler.OpsErrorLoggerMiddleware(opsService)
	endpointNorm := handler.InboundEndpointMiddleware()
	compositeTarget := compositeTargetPlatformMiddleware(compositeResolver)
	compositeGeminiTarget := compositeGeminiTargetPlatformMiddleware(compositeResolver)

	// 未分组 Key 拦截中间件（按协议格式区分错误响应）
	requireGroupAnthropic := middleware.RequireGroupAssignment(settingService, middleware.AnthropicErrorWriter)
	requireGroupGoogle := middleware.RequireGroupAssignment(settingService, middleware.GoogleErrorWriter)

	// 分组级模型白名单准入：在 apiKeyAuth 之后、compositeTarget 之前，
	// 保证校验发生在合成路由改写与调度之前，且只看客户端书写的模型名。
	groupModelAllowlist := middleware.GroupModelAllowlist()

	isOpenAIResponsesCompatibleGatewayPlatform := func(c *gin.Context) bool {
		switch getGroupPlatform(c) {
		case service.PlatformOpenAI, service.PlatformGrok,
			service.PlatformKimi, service.PlatformZhipu, service.PlatformDeepseek, service.PlatformMiniMax:
			// 国产 OpenAI 兼容供应商与 openai/grok 一样经 OpenAI 网关转发。
			return true
		default:
			return false
		}
	}
	rejectLeoUnsupported := func(c *gin.Context, capability string) bool {
		if !service.IsMediaPlatform(getGroupPlatform(c)) {
			return false
		}
		service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"type":    "not_found_error",
				"message": capability + " is not supported for this platform",
			},
		})
		return true
	}
	countTokensHandler := func(c *gin.Context) {
		if rejectLeoUnsupported(c, "Token counting") {
			return
		}
		switch getGroupPlatform(c) {
		case service.PlatformOpenAI, service.PlatformKimi, service.PlatformZhipu, service.PlatformDeepseek, service.PlatformMiniMax:
			h.OpenAIGateway.CountTokens(c)
		case service.PlatformGrok:
			h.OpenAIGateway.GrokCountTokens(c)
		default:
			h.Gateway.CountTokens(c)
		}
	}
	codexModelsHandler := func(c *gin.Context) {
		dispatchCodexModelsGateway(c, h.OpenAIGateway.CodexModels, h.Gateway.CodexModels)
	}
	modelsHandler := func(c *gin.Context) {
		if c.Query("client_version") != "" {
			codexModelsHandler(c)
			return
		}
		h.Gateway.Models(c)
	}
	isOpenAIOnlyEndpointGatewayPlatform := func(c *gin.Context) bool {
		return getGroupPlatform(c) == service.PlatformOpenAI
	}
	imagesHandler := func(c *gin.Context) {
		switch getGroupPlatform(c) {
		case service.PlatformOpenAI:
			h.OpenAIGateway.Images(c)
		case service.PlatformGrok:
			h.OpenAIGateway.GrokImages(c)
		case service.PlatformLeo, service.PlatformOpenAIMedia:
			h.OpenAIGateway.Images(c)
		default:
			service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"type":    "not_found_error",
					"message": "Images API is not supported for this platform",
				},
			})
		}
	}
	videoGenerationHandler := func(c *gin.Context) {
		switch getGroupPlatform(c) {
		// Video status/content lookups below already allow Composite groups; keep
		// task creation aligned so composite keys that route to Grok accounts can
		// submit video generation jobs.
		case service.PlatformGrok, service.PlatformComposite:
			h.OpenAIGateway.GrokVideoGeneration(c)
		case service.PlatformLeo, service.PlatformOpenAIMedia:
			// OpenAI-compatible relays create through multipart/form-data
			// (Sora format); native clients keep the JSON contract.
			if strings.HasPrefix(c.ContentType(), "multipart/form-data") {
				h.OpenAIGateway.SoraVideoCompatCreate(c)
				return
			}
			h.OpenAIGateway.LeoVideoGeneration(c)
		default:
			service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"type":    "not_found_error",
					"message": "Videos API is not supported for this platform",
				},
			})
		}
	}
	videoStatusHandler := func(c *gin.Context) {
		// Legacy Leo async jobs use a local vidjob_ ID and are persisted before
		// the response is returned. They have no model on lookup, so use the ID
		// shape to keep composite-group polling on the Leo job service.
		platform := getGroupPlatform(c)
		requestID := videoLookupRequestID(c)
		if (platform == service.PlatformComposite || service.IsMediaPlatform(platform)) && isLeoVideoJobID(requestID) {
			h.OpenAIGateway.LeoVideoJob(c)
			return
		}
		// Unprefixed external media IDs retain the existing external lookup path;
		// provider-specific ID bindings must be registered before adding another
		// asynchronous media provider to a composite group.
		if platform == service.PlatformGrok || platform == service.PlatformComposite {
			h.OpenAIGateway.GrokVideoStatus(c)
			return
		}
		service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"type":    "not_found_error",
				"message": "Videos API is not supported for this platform",
			},
		})
	}
	// soraVideoStatusHandler serves GET /videos/{id}: local vidjob IDs answer
	// in the OpenAI video-object dialect for Sora-format relays, while the
	// native shape stays on /videos/jobs/{id} and /videos/generations/{id}.
	soraVideoStatusHandler := func(c *gin.Context) {
		platform := getGroupPlatform(c)
		requestID := videoLookupRequestID(c)
		if (platform == service.PlatformComposite || service.IsMediaPlatform(platform)) && isLeoVideoJobID(requestID) {
			h.OpenAIGateway.LeoVideoJobSora(c)
			return
		}
		if platform == service.PlatformGrok || platform == service.PlatformComposite {
			h.OpenAIGateway.GrokVideoStatus(c)
			return
		}
		service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"type":    "not_found_error",
				"message": "Videos API is not supported for this platform",
			},
		})
	}
	videoContentHandler := func(c *gin.Context) {
		platform := getGroupPlatform(c)
		requestID := videoLookupRequestID(c)
		if (platform == service.PlatformComposite || service.IsMediaPlatform(platform)) && isLeoVideoJobID(requestID) {
			h.OpenAIGateway.LeoVideoJobContent(c)
			return
		}
		// Unprefixed external media IDs retain the existing external lookup path;
		// provider-specific ID bindings must be registered before adding another
		// asynchronous media provider to a composite group.
		if platform == service.PlatformGrok || platform == service.PlatformComposite {
			h.OpenAIGateway.GrokVideoContent(c)
			return
		}
		service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"type":    "not_found_error",
				"message": "Videos API is not supported for this platform",
			},
		})
	}
	videoEditHandler := func(c *gin.Context) {
		if getGroupPlatform(c) == service.PlatformGrok {
			h.OpenAIGateway.GrokVideoEdit(c)
			return
		}
		service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"type": "not_found_error", "message": "Videos API is not supported for this platform"}})
	}
	videoExtensionHandler := func(c *gin.Context) {
		if getGroupPlatform(c) == service.PlatformGrok {
			h.OpenAIGateway.GrokVideoExtension(c)
			return
		}
		service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"type": "not_found_error", "message": "Videos API is not supported for this platform"}})
	}
	leoVideoOnly := func(next gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			platform := getGroupPlatform(c)
			if !service.IsMediaPlatform(platform) && platform != service.PlatformComposite {
				service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
				c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"type": "not_found_error", "message": "Videos API is not supported for this platform"}})
				return
			}
			next(c)
		}
	}
	// /responses/*subpath 的子路径会被转发到上游同名端点之后，因此在入口就拒掉
	// 不可转发的子路径，不让它进入调度与转发流程。可转发的判定见
	// service.IsForwardableOpenAIResponsesRequestPath 及 upstream_path_guard.go。
	guardResponsesSubpath := func(next gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			if !service.IsForwardableOpenAIResponsesRequestPath(c) {
				service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalPolicyDenied)
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"error": gin.H{
						"type":    "not_found_error",
						"message": "Unsupported responses subpath",
					},
				})
				return
			}
			if service.IsOpenAIResponsesInputTokensRequestPath(c) && isOpenAIResponsesCompatibleGatewayPlatform(c) {
				h.OpenAIGateway.ResponsesInputTokens(c)
				return
			}
			next(c)
		}
	}
	// API网关（Claude API兼容）
	gateway := r.Group("/v1")
	gateway.Use(bodyLimit)
	gateway.Use(clientRequestID)
	gateway.Use(opsErrorLogger)
	gateway.Use(endpointNorm)
	gateway.Use(gin.HandlerFunc(apiKeyAuth))
	gateway.GET("/sub2api/billing", h.Gateway.KeyBillingInfo)
	gateway.Use(groupModelAllowlist)
	gateway.Use(compositeTarget)
	gateway.Use(requireGroupAnthropic)
	{
		// /v1/messages: auto-route based on group platform
		gateway.POST("/messages", func(c *gin.Context) {
			if rejectLeoUnsupported(c, "Messages API") {
				return
			}
			if isOpenAIResponsesCompatibleGatewayPlatform(c) {
				h.OpenAIGateway.Messages(c)
				return
			}
			h.Gateway.Messages(c)
		})
		// /v1/messages/count_tokens: OpenAI bridges upstream, Grok estimates
		// locally, Anthropic-compatible platforms retain their existing path,
		// and Leo remains explicitly unsupported.
		gateway.POST("/messages/count_tokens", countTokensHandler)
		// Codex CLI / Codex app refresh their model picker from the provider's
		// /models endpoint with a client_version query and expect the ChatGPT
		// Codex manifest format; other clients keep the OpenAI-style list.
		gateway.GET("/models", modelsHandler)
		gateway.GET("/usage", h.Gateway.Usage)
		gateway.POST("/live", h.OpenAIGateway.Live)
		gateway.GET("/live/:call_id", h.OpenAIGateway.LiveSideband)
		// OpenAI Responses API: auto-route based on group platform
		gateway.POST("/responses", func(c *gin.Context) {
			if rejectLeoUnsupported(c, "Responses API") {
				return
			}
			if isOpenAIResponsesCompatibleGatewayPlatform(c) {
				h.OpenAIGateway.Responses(c)
				return
			}
			h.Gateway.Responses(c)
		})
		gateway.POST("/responses/*subpath", guardResponsesSubpath(func(c *gin.Context) {
			if rejectLeoUnsupported(c, "Responses API") {
				return
			}
			if isOpenAIResponsesCompatibleGatewayPlatform(c) {
				h.OpenAIGateway.Responses(c)
				return
			}
			h.Gateway.Responses(c)
		}))
		gateway.POST("/alpha/search", textBodyLimit, func(c *gin.Context) {
			if rejectLeoUnsupported(c, "Alpha Search API") {
				return
			}
			h.OpenAIGateway.AlphaSearch(c)
		})
		gateway.GET("/responses", func(c *gin.Context) {
			if rejectLeoUnsupported(c, "Responses WebSocket API") {
				return
			}
			h.OpenAIGateway.ResponsesWebSocket(c)
		})
		// OpenAI Chat Completions API: auto-route based on group platform
		gateway.POST("/chat/completions", func(c *gin.Context) {
			if rejectLeoUnsupported(c, "Chat Completions API") {
				return
			}
			if isOpenAIResponsesCompatibleGatewayPlatform(c) {
				h.OpenAIGateway.ChatCompletions(c)
				return
			}
			h.Gateway.ChatCompletions(c)
		})
		gateway.POST("/embeddings", textBodyLimit, func(c *gin.Context) {
			if !isOpenAIOnlyEndpointGatewayPlatform(c) {
				service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
				c.JSON(http.StatusNotFound, gin.H{
					"error": gin.H{
						"type":    "not_found_error",
						"message": "Embeddings API is not supported for this platform",
					},
				})
				return
			}
			h.OpenAIGateway.Embeddings(c)
		})
		gateway.POST("/images/generations", imagesHandler)
		gateway.POST("/images/edits", imagesHandler)
		gateway.POST("/images/generations/async", h.AsyncImage.Submit)
		gateway.POST("/images/edits/async", h.AsyncImage.Submit)
		gateway.GET("/images/tasks/:task_id", h.AsyncImage.Get)
		gateway.POST("/images/batches", h.BatchImage.Submit)
		gateway.GET("/images/batches", h.BatchImage.List)
		gateway.GET("/images/batches/models", h.BatchImage.Models)
		gateway.GET("/images/batches/:id", h.BatchImage.Get)
		gateway.GET("/images/batches/:id/items", h.BatchImage.Items)
		gateway.GET("/images/batches/:id/items/:custom_id/content", h.BatchImage.ItemContent)
		gateway.GET("/images/batches/:id/download", h.BatchImage.Download)
		gateway.POST("/images/batches/:id/cancel", h.BatchImage.Cancel)
		gateway.DELETE("/images/batches/:id", h.BatchImage.DeleteRecord)
		gateway.DELETE("/images/batches/:id/outputs", h.BatchImage.DeleteOutputs)
		// OpenAI-compatible clients may create through /videos; xAI receives the
		// canonical /videos/generations route inside the Grok media forwarder.
		gateway.POST("/videos", videoGenerationHandler)
		gateway.POST("/videos/generations", videoGenerationHandler)
		gateway.POST("/videos/uploads", leoVideoOnly(h.OpenAIGateway.LeoVideoUpload))
		gateway.GET("/videos/jobs", leoVideoOnly(h.OpenAIGateway.LeoVideoJobs))
		gateway.GET("/videos/jobs/:job_id/content", leoVideoOnly(h.OpenAIGateway.LeoVideoJobContent))
		gateway.GET("/videos/jobs/:job_id", leoVideoOnly(h.OpenAIGateway.LeoVideoJob))
		gateway.DELETE("/videos/jobs/:job_id", leoVideoOnly(h.OpenAIGateway.CancelLeoVideoJob))
		gateway.POST("/videos/edits", videoEditHandler)
		gateway.POST("/videos/extensions", videoExtensionHandler)
		gateway.GET("/videos/generations/:request_id/content", videoContentHandler)
		gateway.GET("/videos/edits/:request_id/content", videoContentHandler)
		gateway.GET("/videos/extensions/:request_id/content", videoContentHandler)
		gateway.GET("/videos/generations/:request_id", videoStatusHandler)
		gateway.GET("/videos/edits/:request_id", videoStatusHandler)
		gateway.GET("/videos/extensions/:request_id", videoStatusHandler)
		gateway.GET("/videos/:request_id", soraVideoStatusHandler)
		gateway.GET("/videos/:request_id/content", videoContentHandler)

		// xAI Voice APIs (Grok platform only): HTTP TTS/STT + Realtime WS.
		// Not part of the creation-center product surface — gateway relay only.
		voiceHandler := func(endpoint string) gin.HandlerFunc {
			return func(c *gin.Context) {
				if getGroupPlatform(c) != service.PlatformGrok {
					service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
					c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"type": "not_found_error", "message": "Voice API is not supported for this platform"}})
					return
				}
				h.OpenAIGateway.GrokVoice(c, endpoint)
			}
		}
		gateway.POST("/tts", voiceHandler("tts"))
		gateway.POST("/stt", voiceHandler("stt"))
		gateway.POST("/custom-voices", voiceHandler("custom-voices"))
		customVoicePathHandler := func(c *gin.Context) {
			if getGroupPlatform(c) != service.PlatformGrok {
				service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
				c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"type": "not_found_error", "message": "Voice API is not supported for this platform"}})
				return
			}
			h.OpenAIGateway.GrokVoice(c, grokCustomVoiceEndpoint(c))
		}
		gateway.GET("/custom-voices", voiceHandler("custom-voices"))
		gateway.GET("/custom-voices/:voice_id/audio", customVoicePathHandler)
		gateway.GET("/custom-voices/:voice_id", customVoicePathHandler)
		gateway.PATCH("/custom-voices/:voice_id", customVoicePathHandler)
		gateway.DELETE("/custom-voices/:voice_id", customVoicePathHandler)
		gateway.GET("/realtime", func(c *gin.Context) {
			if getGroupPlatform(c) != service.PlatformGrok {
				service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
				c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"type": "not_found_error", "message": "Realtime API is not supported for this platform"}})
				return
			}
			h.OpenAIGateway.GrokRealtime(c)
		})
		gateway.POST("/web_search", func(c *gin.Context) {
			if getGroupPlatform(c) != service.PlatformGrok {
				service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
				c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"type": "not_found_error", "message": "Web Search API is not supported for this platform"}})
				return
			}
			h.Gateway.WebSearch(c)
		})
		gateway.POST("/x_search", func(c *gin.Context) {
			if getGroupPlatform(c) != service.PlatformGrok {
				service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
				c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"type": "not_found_error", "message": "X Search API is not supported for this platform"}})
				return
			}
			h.Gateway.XSearch(c)
		})
	}

	// Gemini 原生 API 兼容层（Gemini SDK/CLI 直连）
	gemini := r.Group("/v1beta")
	gemini.Use(bodyLimit)
	gemini.Use(clientRequestID)
	gemini.Use(opsErrorLogger)
	gemini.Use(endpointNorm)
	gemini.Use(middleware.APIKeyAuthWithSubscriptionGoogle(apiKeyService, subscriptionService, cfg))
	gemini.Use(groupModelAllowlist)
	gemini.Use(compositeGeminiTarget)
	gemini.Use(requireGroupGoogle)
	{
		gemini.GET("/models", h.Gateway.GeminiV1BetaListModels)
		gemini.GET("/models/:model", h.Gateway.GeminiV1BetaGetModel)
		// Gin treats ":" as a param marker, but Gemini uses "{model}:{action}" in the same segment.
		gemini.POST("/models/*modelAction", h.Gateway.GeminiV1BetaModels)
	}

	// OpenAI Responses API（不带v1前缀的别名）— auto-route based on group platform
	responsesHandler := func(c *gin.Context) {
		if rejectLeoUnsupported(c, "Responses API") {
			return
		}
		if isOpenAIResponsesCompatibleGatewayPlatform(c) {
			h.OpenAIGateway.Responses(c)
			return
		}
		h.Gateway.Responses(c)
	}
	// 根路径别名共用中间件链：白名单准入在 apiKeyAuth 之后、compositeTarget
	// 之前，避免逐条路由手工维护链导致漏挂。
	rootRoute := func(method, path string, limit gin.HandlerFunc, handler gin.HandlerFunc) {
		r.Handle(method, path, limit, clientRequestID, opsErrorLogger, endpointNorm, gin.HandlerFunc(apiKeyAuth), groupModelAllowlist, compositeTarget, requireGroupAnthropic, handler)
	}
	rootRoute(http.MethodPost, "/responses", bodyLimit, responsesHandler)
	rootRoute(http.MethodPost, "/responses/*subpath", bodyLimit, guardResponsesSubpath(responsesHandler))
	rootRoute(http.MethodPost, "/alpha/search", textBodyLimit, func(c *gin.Context) {
		if rejectLeoUnsupported(c, "Alpha Search API") {
			return
		}
		h.OpenAIGateway.AlphaSearch(c)
	})
	rootRoute(http.MethodGet, "/responses", bodyLimit, func(c *gin.Context) {
		if rejectLeoUnsupported(c, "Responses WebSocket API") {
			return
		}
		h.OpenAIGateway.ResponsesWebSocket(c)
	})
	rootRoute(http.MethodGet, "/models", bodyLimit, modelsHandler)
	rootRoute(http.MethodPost, "/messages/count_tokens", bodyLimit, countTokensHandler)
	codexDirect := r.Group("/backend-api/codex")
	codexDirect.Use(bodyLimit, clientRequestID, opsErrorLogger, endpointNorm, gin.HandlerFunc(apiKeyAuth), groupModelAllowlist, compositeTarget, requireGroupAnthropic)
	{
		codexDirect.POST("/realtime/calls", h.OpenAIGateway.Live)
		codexDirect.GET("/:call_id", h.OpenAIGateway.LiveSideband)
		codexDirect.POST("/responses", responsesHandler)
		codexDirect.POST("/responses/*subpath", guardResponsesSubpath(responsesHandler))
		codexDirect.POST("/alpha/search", textBodyLimit, func(c *gin.Context) {
			if rejectLeoUnsupported(c, "Alpha Search API") {
				return
			}
			h.OpenAIGateway.AlphaSearch(c)
		})
		codexDirect.GET("/responses", func(c *gin.Context) {
			if rejectLeoUnsupported(c, "Responses WebSocket API") {
				return
			}
			h.OpenAIGateway.ResponsesWebSocket(c)
		})
		codexDirect.GET("/models", codexModelsHandler)
	}
	// OpenAI Chat Completions API（不带v1前缀的别名）— auto-route based on group platform
	rootRoute(http.MethodPost, "/chat/completions", bodyLimit, func(c *gin.Context) {
		if rejectLeoUnsupported(c, "Chat Completions API") {
			return
		}
		if isOpenAIResponsesCompatibleGatewayPlatform(c) {
			h.OpenAIGateway.ChatCompletions(c)
			return
		}
		h.Gateway.ChatCompletions(c)
	})
	rootRoute(http.MethodPost, "/embeddings", textBodyLimit, func(c *gin.Context) {
		if !isOpenAIOnlyEndpointGatewayPlatform(c) {
			service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"type":    "not_found_error",
					"message": "Embeddings API is not supported for this platform",
				},
			})
			return
		}
		h.OpenAIGateway.Embeddings(c)
	})
	rootRoute(http.MethodPost, "/images/generations", bodyLimit, imagesHandler)
	rootRoute(http.MethodPost, "/images/edits", bodyLimit, imagesHandler)
	rootRoute(http.MethodPost, "/images/generations/async", bodyLimit, h.AsyncImage.Submit)
	rootRoute(http.MethodPost, "/images/edits/async", bodyLimit, h.AsyncImage.Submit)
	rootRoute(http.MethodGet, "/images/tasks/:task_id", bodyLimit, h.AsyncImage.Get)
	r.GET("/internal/video-inputs/:token", h.OpenAIGateway.LeoVideoInputInternal)
	r.GET("/media/image-inputs/:token", h.OpenAIGateway.LeoImageInputPublic)
	r.HEAD("/media/image-inputs/:token", h.OpenAIGateway.LeoImageInputPublic)
	rootRoute(http.MethodPost, "/videos", bodyLimit, videoGenerationHandler)
	rootRoute(http.MethodPost, "/videos/generations", bodyLimit, videoGenerationHandler)
	rootRoute(http.MethodPost, "/videos/uploads", bodyLimit, leoVideoOnly(h.OpenAIGateway.LeoVideoUpload))
	rootRoute(http.MethodGet, "/videos/jobs", bodyLimit, leoVideoOnly(h.OpenAIGateway.LeoVideoJobs))
	rootRoute(http.MethodGet, "/videos/jobs/:job_id/content", bodyLimit, leoVideoOnly(h.OpenAIGateway.LeoVideoJobContent))
	rootRoute(http.MethodGet, "/videos/jobs/:job_id", bodyLimit, leoVideoOnly(h.OpenAIGateway.LeoVideoJob))
	rootRoute(http.MethodDelete, "/videos/jobs/:job_id", bodyLimit, leoVideoOnly(h.OpenAIGateway.CancelLeoVideoJob))
	rootRoute(http.MethodPost, "/videos/edits", bodyLimit, videoEditHandler)
	rootRoute(http.MethodPost, "/videos/extensions", bodyLimit, videoExtensionHandler)
	rootRoute(http.MethodGet, "/videos/generations/:request_id/content", bodyLimit, videoContentHandler)
	rootRoute(http.MethodGet, "/videos/edits/:request_id/content", bodyLimit, videoContentHandler)
	rootRoute(http.MethodGet, "/videos/extensions/:request_id/content", bodyLimit, videoContentHandler)
	rootRoute(http.MethodGet, "/videos/generations/:request_id", bodyLimit, videoStatusHandler)
	rootRoute(http.MethodGet, "/videos/edits/:request_id", bodyLimit, videoStatusHandler)
	rootRoute(http.MethodGet, "/videos/extensions/:request_id", bodyLimit, videoStatusHandler)
	rootRoute(http.MethodGet, "/videos/:request_id", bodyLimit, soraVideoStatusHandler)
	rootRoute(http.MethodGet, "/videos/:request_id/content", bodyLimit, videoContentHandler)

	rootVoiceHandler := func(endpoint string) gin.HandlerFunc {
		return func(c *gin.Context) {
			if getGroupPlatform(c) != service.PlatformGrok {
				service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
				c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"type": "not_found_error", "message": "Voice API is not supported for this platform"}})
				return
			}
			h.OpenAIGateway.GrokVoice(c, endpoint)
		}
	}
	rootRoute(http.MethodPost, "/tts", bodyLimit, rootVoiceHandler("tts"))
	rootRoute(http.MethodPost, "/stt", bodyLimit, rootVoiceHandler("stt"))
	rootRoute(http.MethodPost, "/custom-voices", bodyLimit, rootVoiceHandler("custom-voices"))
	rootCustomVoicePathHandler := func(c *gin.Context) {
		if getGroupPlatform(c) != service.PlatformGrok {
			service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"type": "not_found_error", "message": "Voice API is not supported for this platform"}})
			return
		}
		h.OpenAIGateway.GrokVoice(c, grokCustomVoiceEndpoint(c))
	}
	rootRoute(http.MethodGet, "/custom-voices", bodyLimit, rootVoiceHandler("custom-voices"))
	rootRoute(http.MethodGet, "/custom-voices/:voice_id/audio", bodyLimit, rootCustomVoicePathHandler)
	rootRoute(http.MethodGet, "/custom-voices/:voice_id", bodyLimit, rootCustomVoicePathHandler)
	rootRoute(http.MethodPatch, "/custom-voices/:voice_id", bodyLimit, rootCustomVoicePathHandler)
	rootRoute(http.MethodDelete, "/custom-voices/:voice_id", bodyLimit, rootCustomVoicePathHandler)
	rootRoute(http.MethodGet, "/realtime", bodyLimit, func(c *gin.Context) {
		if getGroupPlatform(c) != service.PlatformGrok {
			service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"type": "not_found_error", "message": "Realtime API is not supported for this platform"}})
			return
		}
		h.OpenAIGateway.GrokRealtime(c)
	})
	rootRoute(http.MethodPost, "/web_search", bodyLimit, func(c *gin.Context) {
		if getGroupPlatform(c) != service.PlatformGrok {
			service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"type": "not_found_error", "message": "Web Search API is not supported for this platform"}})
			return
		}
		h.Gateway.WebSearch(c)
	})
	rootRoute(http.MethodPost, "/x_search", bodyLimit, func(c *gin.Context) {
		if getGroupPlatform(c) != service.PlatformGrok {
			service.MarkOpsClientBusinessLimited(c, service.OpsClientBusinessLimitedReasonLocalFeatureGate)
			c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"type": "not_found_error", "message": "X Search API is not supported for this platform"}})
			return
		}
		h.Gateway.XSearch(c)
	})

	// Antigravity 模型列表
	r.GET("/antigravity/models", gin.HandlerFunc(apiKeyAuth), requireGroupAnthropic, h.Gateway.AntigravityModels)

	// Antigravity 专用路由（仅使用 antigravity 账户，不混合调度）
	antigravityV1 := r.Group("/antigravity/v1")
	antigravityV1.Use(bodyLimit)
	antigravityV1.Use(clientRequestID)
	antigravityV1.Use(opsErrorLogger)
	antigravityV1.Use(endpointNorm)
	antigravityV1.Use(middleware.ForcePlatform(service.PlatformAntigravity))
	antigravityV1.Use(gin.HandlerFunc(apiKeyAuth))
	antigravityV1.Use(groupModelAllowlist)
	antigravityV1.Use(requireGroupAnthropic)
	{
		antigravityV1.POST("/messages", h.Gateway.Messages)
		antigravityV1.POST("/messages/count_tokens", h.Gateway.CountTokens)
		antigravityV1.GET("/models", h.Gateway.AntigravityModels)
		antigravityV1.GET("/usage", h.Gateway.Usage)
	}

	antigravityV1Beta := r.Group("/antigravity/v1beta")
	antigravityV1Beta.Use(bodyLimit)
	antigravityV1Beta.Use(clientRequestID)
	antigravityV1Beta.Use(opsErrorLogger)
	antigravityV1Beta.Use(endpointNorm)
	antigravityV1Beta.Use(middleware.ForcePlatform(service.PlatformAntigravity))
	antigravityV1Beta.Use(middleware.APIKeyAuthWithSubscriptionGoogle(apiKeyService, subscriptionService, cfg))
	antigravityV1Beta.Use(groupModelAllowlist)
	antigravityV1Beta.Use(requireGroupGoogle)
	{
		antigravityV1Beta.GET("/models", h.Gateway.GeminiV1BetaListModels)
		antigravityV1Beta.GET("/models/:model", h.Gateway.GeminiV1BetaGetModel)
		antigravityV1Beta.POST("/models/*modelAction", h.Gateway.GeminiV1BetaModels)
	}

}

func dispatchCodexModelsGateway(c *gin.Context, openAIHandler, generatedHandler gin.HandlerFunc) {
	if getGroupPlatform(c) == service.PlatformOpenAI {
		openAIHandler(c)
		return
	}
	generatedHandler(c)
}

// getGroupPlatform extracts the group platform from the API Key stored in context.
func getGroupPlatform(c *gin.Context) string {
	apiKey, ok := middleware.GetAPIKeyFromContext(c)
	if !ok || apiKey.Group == nil {
		return ""
	}
	if apiKey.Group.Platform == service.PlatformComposite {
		if platform, ok := service.ResolvedTargetPlatformFromContext(c.Request.Context()); ok {
			return platform
		}
	}
	return apiKey.Group.Platform
}

func videoLookupRequestID(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if requestID := strings.TrimSpace(c.Param("request_id")); requestID != "" {
		return requestID
	}
	return strings.TrimSpace(c.Param("job_id"))
}

func isLeoVideoJobID(requestID string) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(requestID)), "vidjob_")
}

func compositeTargetPlatformMiddleware(resolver *service.CompositeRouteResolver) gin.HandlerFunc {
	if resolver == nil {
		resolver = service.NewCompositeRouteResolver(nil)
	}
	return func(c *gin.Context) {
		apiKey, ok := middleware.GetAPIKeyFromContext(c)
		if !ok || apiKey == nil || apiKey.Group == nil || apiKey.Group.Platform != service.PlatformComposite {
			c.Next()
			return
		}
		if c.Request == nil || c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		body, err := pkghttputil.ReadRequestBodyWithPrealloc(c.Request)
		if err != nil {
			status := http.StatusBadRequest
			message := "Failed to read request body"
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr) {
				status = http.StatusRequestEntityTooLarge
				message = "Request body is too large"
			}
			c.JSON(status, gin.H{"error": gin.H{"type": "invalid_request_error", "message": message}})
			c.Abort()
			return
		}

		// Live 入口按 session.model 分发并改写（与白名单准入、Live handler 的
		// 解析保持一致），避免顶层 model 别名与 session 模型不一致时改错对象。
		routePath := c.FullPath()
		model := requestmodel.FromBodyForRoute(routePath, c.GetHeader("Content-Type"), body)
		if model != "" {
			decision, err := resolver.Resolve(c.Request.Context(), apiKey.Group.ID, model, compositeRouteEndpointForPath(c.Request.URL.Path))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"type": "server_error", "message": "Failed to resolve composite model route"}})
				c.Abort()
				return
			}
			if decision.Matched {
				c.Request = c.Request.WithContext(service.WithCompositeRouteDecision(c.Request.Context(), decision))
				if upstreamModel := strings.TrimSpace(decision.UpstreamModel); upstreamModel != "" && upstreamModel != model && gjson.ValidBytes(body) {
					if _, modelPath := requestmodel.JSONModelPathForRoute(routePath, body); modelPath != "" {
						if rewritten, rewriteErr := sjson.SetBytes(body, modelPath, upstreamModel); rewriteErr == nil {
							body = rewritten
						}
					}
				}
			}
		}
		requestmodel.ResetRequestBody(c.Request, body)
		c.Next()
	}
}

func compositeGeminiTargetPlatformMiddleware(resolver *service.CompositeRouteResolver) gin.HandlerFunc {
	if resolver == nil {
		resolver = service.NewCompositeRouteResolver(nil)
	}
	return func(c *gin.Context) {
		apiKey, ok := middleware.GetAPIKeyFromContext(c)
		if ok && apiKey != nil && apiKey.Group != nil && apiKey.Group.Platform == service.PlatformComposite {
			model := compositeGeminiModelFromParams(c)
			if model != "" {
				decision, err := resolver.Resolve(c.Request.Context(), apiKey.Group.ID, model, service.CompositeRouteEndpointGemini)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"type": "server_error", "message": "Failed to resolve composite model route"}})
					c.Abort()
					return
				}
				if decision.Matched {
					c.Request = c.Request.WithContext(service.WithCompositeRouteDecision(c.Request.Context(), decision))
				}
			}
			if _, resolved := service.ResolvedTargetPlatformFromContext(c.Request.Context()); !resolved {
				c.Request = c.Request.WithContext(service.WithResolvedTargetPlatform(c.Request.Context(), service.PlatformGemini))
			}
		}
		c.Next()
	}
}

// grokCustomVoiceEndpoint derives the upstream Voice endpoint for the
// /custom-voices/:voice_id[/audio] routes.
//
// The /audio suffix must be decided from the matched route template, not from
// the raw URL path: a voice literally named "audio" makes GET
// /custom-voices/audio match /custom-voices/:voice_id, and a raw-path suffix
// check would rewrite it to custom-voices/audio/audio — turning a profile
// lookup into an audio download.
func grokCustomVoiceEndpoint(c *gin.Context) string {
	endpoint := "custom-voices/" + c.Param("voice_id")
	if strings.HasSuffix(c.FullPath(), "/:voice_id/audio") {
		endpoint += "/audio"
	}
	return endpoint
}

func compositeGeminiModelFromParams(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if model := strings.TrimSpace(c.Param("model")); model != "" {
		return model
	}
	modelAction := strings.TrimPrefix(strings.TrimSpace(c.Param("modelAction")), "/")
	if modelAction == "" {
		return ""
	}
	if idx := strings.LastIndex(modelAction, ":"); idx >= 0 {
		return strings.TrimSpace(modelAction[:idx])
	}
	return modelAction
}

func compositeRouteEndpointForPath(path string) string {
	switch {
	case strings.Contains(path, "/messages/count_tokens"):
		return service.CompositeRouteEndpointCountTokens
	case strings.Contains(path, "/messages"):
		return service.CompositeRouteEndpointMessages
	case strings.Contains(path, "/responses"),
		strings.Contains(path, "/alpha/search"),
		strings.Contains(path, "/realtime/calls"),
		strings.HasSuffix(strings.TrimRight(path, "/"), "/live"):
		return service.CompositeRouteEndpointResponses
	case strings.Contains(path, "/chat/completions"):
		return service.CompositeRouteEndpointChatCompletions
	case strings.Contains(path, "/embeddings"):
		return service.CompositeRouteEndpointEmbeddings
	case strings.Contains(path, "/images/"):
		return service.CompositeRouteEndpointImages
	case strings.Contains(path, "/v1beta/"):
		return service.CompositeRouteEndpointGemini
	default:
		return service.CompositeRouteEndpointAny
	}
}
