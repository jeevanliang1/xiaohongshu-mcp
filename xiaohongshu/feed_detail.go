package xiaohongshu

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
)

// FeedDetailAction 表示 Feed 详情页动作
type FeedDetailAction struct {
	page *rod.Page
}

// NewFeedDetailAction 创建 Feed 详情页动作
func NewFeedDetailAction(page *rod.Page) *FeedDetailAction {
	return &FeedDetailAction{page: page}
}

// GetFeedDetail 获取 Feed 详情页数据
func (f *FeedDetailAction) GetFeedDetail(ctx context.Context, feedID, xsecToken string) (*FeedDetailResponse, error) {
	page := f.page.Context(ctx).Timeout(60 * time.Second)

	// 构建详情页 URL
	url := makeFeedDetailURL(feedID, xsecToken)

	// 导航到详情页
	page.MustNavigate(url)
	page.MustWaitDOMStable()
	time.Sleep(1 * time.Second)

	// 获取 window.__INITIAL_STATE__ 并转换为 JSON 字符串
	// 使用 Eval 而不是 MustEval，以便捕获错误而不是 panic
	resultVal, err := page.Eval(`() => {
		try {
			if (!window.__INITIAL_STATE__) {
				return "";
			}
			
			const state = window.__INITIAL_STATE__;
			
			// 安全地提取 note 数据，避免循环引用
			function safeExtractNoteData(noteObj) {
				if (!noteObj || !noteObj.note) return null;
				
				const note = noteObj.note;
				const noteDetailMap = {};
				
				// 遍历 noteDetailMap，安全提取每个 note 的数据
				if (note.noteDetailMap) {
					for (const [key, detail] of Object.entries(note.noteDetailMap)) {
						if (!detail || !detail.note) continue;
						
						const noteDetail = detail.note;
						const comments = detail.comments || { comments: [] };
						
						// 手动提取必要字段，避免循环引用
						noteDetailMap[key] = {
							note: {
								id: noteDetail.id || "",
								title: noteDetail.title || "",
								desc: noteDetail.desc || "",
								type: noteDetail.type || "",
								time: noteDetail.time || 0,
								lastUpdateTime: noteDetail.lastUpdateTime || 0,
								user: noteDetail.user ? {
									userId: noteDetail.user.userId || "",
									nickname: noteDetail.user.nickname || noteDetail.user.nickName || "",
									avatar: noteDetail.user.avatar || ""
								} : {},
								images: noteDetail.images ? (Array.isArray(noteDetail.images) ? noteDetail.images.map(img => ({
									url: img.url || "",
									width: img.width || 0,
									height: img.height || 0,
									fileId: img.fileId || "",
									urlPre: img.urlPre || "",
									urlDefault: img.urlDefault || ""
								})) : []) : [],
								video: noteDetail.video ? {
									media: noteDetail.video.media ? {
										stream: noteDetail.video.media.stream ? {
											h264: noteDetail.video.media.stream.h264 ? {
												url: noteDetail.video.media.stream.h264.url || ""
											} : {}
										} : {}
									} : {}
								} : null,
								interactInfo: noteDetail.interactInfo ? {
									liked: noteDetail.interactInfo.liked || false,
									likedCount: noteDetail.interactInfo.likedCount || "0",
									sharedCount: noteDetail.interactInfo.sharedCount || "0",
									commentCount: noteDetail.interactInfo.commentCount || "0",
									collectedCount: noteDetail.interactInfo.collectedCount || "0",
									collected: noteDetail.interactInfo.collected || false
								} : {},
								tagList: noteDetail.tagList || []
							},
							comments: {
								comments: Array.isArray(comments.comments) ? comments.comments.map(comment => ({
									id: comment.id || "",
									content: comment.content || "",
									user: comment.user ? {
										userId: comment.user.userId || "",
										nickname: comment.user.nickname || comment.user.nickName || "",
										avatar: comment.user.avatar || ""
									} : {},
									time: comment.time || 0,
									liked: comment.liked || false,
									likedCount: comment.likedCount || 0
								})) : []
							}
						};
					}
				}
				
				return {
					note: {
						noteDetailMap: noteDetailMap
					}
				};
			}
			
			const noteData = safeExtractNoteData(state);
			if (!noteData) {
				// 如果提取失败，尝试直接序列化（可能会失败）
				try {
					return JSON.stringify(state);
				} catch (e) {
					return "";
				}
			}
			
			return JSON.stringify(noteData);
		} catch (e) {
			return "";
		}
	}`)

	if err != nil {
		logrus.Errorf("执行 JavaScript 获取 Feed 详情失败: %v", err)
		return nil, fmt.Errorf("执行 JavaScript 失败: %w", err)
	}

	result := ""
	if resultVal != nil {
		result = resultVal.Value.String()
	}

	if result == "" {
		logrus.Warn("__INITIAL_STATE__ 未找到或为空")
		return nil, fmt.Errorf("__INITIAL_STATE__ not found or empty")
	}

	// 定义响应结构并直接反序列化
	var initialState struct {
		Note struct {
			NoteDetailMap map[string]struct {
				Note     FeedDetail  `json:"note"`
				Comments CommentList `json:"comments"`
			} `json:"noteDetailMap"`
		} `json:"note"`
	}

	if err := json.Unmarshal([]byte(result), &initialState); err != nil {
		return nil, fmt.Errorf("failed to unmarshal __INITIAL_STATE__: %w", err)
	}

	// 从 noteDetailMap 中获取对应 feedID 的数据
	noteDetail, exists := initialState.Note.NoteDetailMap[feedID]
	if !exists {
		return nil, fmt.Errorf("feed %s not found in noteDetailMap", feedID)
	}

	return &FeedDetailResponse{
		Note:     noteDetail.Note,
		Comments: noteDetail.Comments,
	}, nil
}

func makeFeedDetailURL(feedID, xsecToken string) string {
	return fmt.Sprintf("https://www.xiaohongshu.com/explore/%s?xsec_token=%s&xsec_source=pc_feed", feedID, xsecToken)
}
