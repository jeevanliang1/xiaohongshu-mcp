package xiaohongshu

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/go-rod/rod"
)

type SearchResult struct {
	Search struct {
		Feeds FeedsValue `json:"feeds"`
	} `json:"search"`
}

type SearchAction struct {
	page *rod.Page
}

func NewSearchAction(page *rod.Page) *SearchAction {
	pp := page.Timeout(60 * time.Second)

	return &SearchAction{page: pp}
}

func (s *SearchAction) Search(ctx context.Context, keyword string) ([]Feed, error) {
	page := s.page.Context(ctx)

	searchURL := makeSearchURL(keyword)
	page.MustNavigate(searchURL)
	page.MustWaitStable()

	page.MustWait(`() => window.__INITIAL_STATE__ !== undefined`)

	// 安全地提取搜索数据，避免循环引用
	result := page.MustEval(`() => {
			if (window.__INITIAL_STATE__) {
				const state = window.__INITIAL_STATE__;
				
				// 安全地提取 feeds 数据，避免循环引用
				function safeExtractFeeds(searchObj) {
					if (!searchObj || !searchObj.feeds) return null;
					
					const feeds = searchObj.feeds;
					if (!feeds._value || !Array.isArray(feeds._value)) return null;
					
					// 手动提取每个 feed 的必要字段，避免循环引用
					return feeds._value.map(feed => {
						if (!feed) return null;
						
						return {
							xsecToken: feed.xsecToken || "",
							id: feed.id || "",
							modelType: feed.modelType || "",
							index: feed.index || 0,
							noteCard: feed.noteCard ? {
								type: feed.noteCard.type || "",
								displayTitle: feed.noteCard.displayTitle || "",
								user: feed.noteCard.user ? {
									userId: feed.noteCard.user.userId || "",
									nickname: feed.noteCard.user.nickname || feed.noteCard.user.nickName || "",
									avatar: feed.noteCard.user.avatar || ""
								} : {},
								interactInfo: feed.noteCard.interactInfo ? {
									liked: feed.noteCard.interactInfo.liked || false,
									likedCount: feed.noteCard.interactInfo.likedCount || "0",
									sharedCount: feed.noteCard.interactInfo.sharedCount || "0",
									commentCount: feed.noteCard.interactInfo.commentCount || "0",
									collectedCount: feed.noteCard.interactInfo.collectedCount || "0",
									collected: feed.noteCard.interactInfo.collected || false
								} : {},
								cover: feed.noteCard.cover ? {
									width: feed.noteCard.cover.width || 0,
									height: feed.noteCard.cover.height || 0,
									url: feed.noteCard.cover.url || "",
									fileId: feed.noteCard.cover.fileId || "",
									urlPre: feed.noteCard.cover.urlPre || "",
									urlDefault: feed.noteCard.cover.urlDefault || "",
									infoList: feed.noteCard.cover.infoList || []
								} : {},
								video: feed.noteCard.video ? {
									capa: feed.noteCard.video.capa ? {
										duration: feed.noteCard.video.capa.duration || 0
									} : {}
								} : null
							} : {}
						};
					}).filter(feed => feed !== null);
				}
				
				const feeds = safeExtractFeeds(state.search);
				if (!feeds) return "";
				
				const searchData = {
					search: {
						feeds: {
							_value: feeds
						}
					}
				};
				
				return JSON.stringify(searchData);
			}
			return "";
		}`).String()

	if result == "" {
		return nil, fmt.Errorf("__INITIAL_STATE__ not found")
	}

	var searchResult SearchResult
	if err := json.Unmarshal([]byte(result), &searchResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal __INITIAL_STATE__: %w", err)
	}

	return searchResult.Search.Feeds.Value, nil
}

func makeSearchURL(keyword string) string {

	values := url.Values{}
	values.Set("keyword", keyword)
	values.Set("source", "web_explore_feed")

	return fmt.Sprintf("https://www.xiaohongshu.com/search_result?%s", values.Encode())
}
