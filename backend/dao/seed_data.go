package dao

import (
	"encoding/json"
	"strings"
	"time"

	"personal/blog/backend/model"
)

func seedSiteProfile() model.SiteProfile {
	return model.SiteProfile{
		Name:     "Vanyongqi",
		Headline: "写给后端工程师与大模型系统开发者的技术博客",
		Intro:    "这里记录我在工作中遇到的一些解决的问题和生活中的思考。",
		Location: "银河系 / 地球 / 北京",
		Domain:   "718614413.xyz",
		Email:    "718614413@qq.com",
		Motto:    "把服务做稳，把链路做清，把大模型工程做成系统。",
		TechStack: []string{
			"Golang", "Python", "MySQL", "Redis", "Doris", "Kafka", "LLM Serving",
		},
		Stats: []model.SiteStat{
			{Label: "线上链路", Value: "20+"},
			{Label: "模型场景", Value: "6"},
			{Label: "工程经验", Value: "8y"},
		},
		SocialLinks: []model.SocialLink{
			{Label: "GitHub", URL: "https://github.com/vanyongqi"},
			{Label: "Email", URL: "mailto:718614413@qq.com"},
			{Label: "Domain", URL: "https://718614413.xyz"},
		},
	}
}

func seedPosts() []model.Post {
	return []model.Post{
		{
			Slug:        "llm-apps-fail-in-the-backend",
			Title:       "大模型应用真正难的是后端链路，不是 Prompt",
			Summary:     "把检索、重排、缓存、限流、降级和 tracing 串起来后，你会发现大模型项目的主战场其实在后端系统。",
			Category:    "LLM Systems",
			ReadTime:    "9 分钟",
			CoverLabel:  "模型工程",
			ContentMarkdown: strings.TrimSpace(`
一个线上可用的大模型应用，至少包含 query rewrite、召回、重排、上下文压缩、模型调用、结果缓存和链路观测。只盯 Prompt，迟早会在延迟、成本和稳定性上失控。

## 先把链路拆开，再谈模型效果

把每一段耗时、失败率和命中率拆出来，你才能知道问题是在 embedding、向量检索、rerank，还是模型服务本身。

- 召回与重排必须分别打点
- 模型服务要有超时、熔断和 fallback
- 上下文拼接必须可回放、可审计

> 模型应用要稳定上线，靠的是后端工程，不是灵感。
`),
			Tags:        []string{"RAG", "Golang", "Observability"},
			Featured:    true,
			PublishedAt: mustParseDate("2026-03-22"),
		},
		{
			Slug:        "designing-an-inference-gateway",
			Title:       "高并发推理网关该怎么设计，才能扛住真实流量",
			Summary:     "从连接池、配额、超时、模型路由和观测性切入，讲清楚推理网关在生产环境里的关键设计点。",
			Category:    "Backend",
			ReadTime:    "8 分钟",
			CoverLabel:  "服务治理",
			ContentMarkdown: strings.TrimSpace(`
当你同时接多个模型服务商、多个模型版本和不同优先级流量时，网关就不只是代理层，而是整个模型平台的控制面。

## 四件必须先定义清楚的事

并发策略、超时策略、失败回退顺序和配额策略要在接入业务前就明确，否则流量一上来很快会进入不可解释状态。

- 调用级超时优先于网关总超时
- 租户配额要支持硬限额和软告警
- 模型路由策略要可热更新
- 请求与响应必须带 trace id
`),
			Tags:        []string{"Inference", "Gateway", "Golang"},
			Featured:    true,
			PublishedAt: mustParseDate("2026-03-14"),
		},
		{
			Slug:        "doris-for-event-replay-and-analysis",
			Title:       "用 Doris 做事件分析与质量回放，怎么兼顾速度和口径",
			Summary:     "埋点多、明细大、回放要求高时，Doris 很适合承担分析与排障的双重职责，但前提是口径和明细组织要先想清楚。",
			Category:    "Data Infra",
			ReadTime:    "7 分钟",
			CoverLabel:  "数据基础设施",
			ContentMarkdown: strings.TrimSpace(`
事件分析和质量回放其实是两类需求：前者看聚合口径，后者看原始链路。表设计如果只服务其中一种，另一种迟早会变得很难维护。

## 回放链路要保留什么

最少要保留 trace id、时间戳、租户、版本、上下游关键字段和决策结果，否则故障回放只能靠猜。

- 聚合字段和明细字段分层存储
- 高频筛选字段提前建好排序键
- 回放接口不要依赖线上热路径
`),
			Tags:        []string{"Doris", "ETL", "Analytics"},
			Featured:    false,
			PublishedAt: mustParseDate("2026-03-02"),
		},
	}
}

func seedProjects() []model.Project {
	return []model.Project{
		{
			Name:      "LLM Gateway",
			Summary:   "统一收口多模型服务、租户配额、降级与可观测性的推理网关，目标是让模型调用具备后端级别的稳定性。",
			Status:    "Production",
			Link:      "https://718614413.xyz",
			ImageURL:  "https://placehold.co/960x540/png?text=LLM+Gateway",
			Accent:    "ember",
			TechStack: []string{"Golang", "Redis", "Kafka", "OpenTelemetry"},
		},
		{
			Name:      "Retrieval Pipeline",
			Summary:   "围绕 embedding、召回、重排和上下文压缩搭建的 RAG 检索链路，用来验证模型应用在真实数据上的稳定性。",
			Status:    "Ongoing",
			Link:      "https://718614413.xyz",
			ImageURL:  "https://placehold.co/960x540/png?text=Retrieval+Pipeline",
			Accent:    "forest",
			TechStack: []string{"Python", "Milvus", "MySQL", "Celery"},
		},
		{
			Name:      "Event Replay Platform",
			Summary:   "针对行为埋点、模型调用和质量分析做的事件回放平台，用来排查线上口径、回放请求和验证版本差异。",
			Status:    "Prototype",
			Link:      "https://718614413.xyz",
			ImageURL:  "https://placehold.co/960x540/png?text=Event+Replay+Platform",
			Accent:    "ink",
			TechStack: []string{"Golang", "Doris", "ClickHouse", "S3"},
		},
	}
}

func seedVideos() []model.Video {
	return []model.Video{
		{
			Title:       "推理链路拆解",
			Description: "关于检索、重排、缓存和 tracing 的一段视频记录。",
			URL:         "https://www.youtube.com/watch?v=aircAruvnKk",
			ThumbnailURL: "https://i.ytimg.com/vi/aircAruvnKk/hqdefault.jpg",
			PublishedAt: mustParseDate("2026-03-24"),
		},
		{
			Title:       "推理网关治理",
			Description: "关于超时、配额、降级和模型路由的一段视频记录。",
			URL:         "https://www.youtube.com/watch?v=5MgBikgcWnY",
			ThumbnailURL: "https://i.ytimg.com/vi/5MgBikgcWnY/hqdefault.jpg",
			PublishedAt: mustParseDate("2026-03-16"),
		},
	}
}

func seedTimeline() []model.TimelineEntry {
	return []model.TimelineEntry{
		{Period: "2026", Title: "聚焦模型系统工程", Description: "把写作重心放到推理服务、RAG 链路和后端治理上。"},
		{Period: "2025", Title: "沉淀数据与检索基础设施", Description: "持续打磨分析链路、查询口径和数据回放能力。"},
		{Period: "2024", Title: "系统化记录后端经验", Description: "开始把线上故障、性能治理和架构判断沉淀成文章。"},
		{Period: "2023", Title: "搭建长期输出机制", Description: "围绕项目、知识库和博客建立个人技术表达体系。"},
	}
}

func mustJSON(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func mustParseDate(value string) time.Time {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		panic(err)
	}
	return parsed
}
