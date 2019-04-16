package structures

import (
	"fmt"
	"giveaway/instagram/structures/stories"
)

type PageInfo struct {
	HasNextPage bool   `json:"has_next_page"`
	EndCursor   string `json:"end_cursor"`
}

type ShortCodeMediaResponse struct {
	Data struct {
		ShortCodeMedia *struct {
			Typename   string `json:"__typename"`
			ID         string `json:"id"`
			ShortCode  string `json:"shortcode"`
			Dimensions struct {
				Height int `json:"height"`
				Width  int `json:"width"`
			} `json:"dimensions"`
			GatingInfo       interface{} `json:"gating_info"`
			MediaPreview     string      `json:"media_preview"`
			DisplayURL       string      `json:"display_url"`
			DisplayResources []struct {
				Src          string `json:"src"`
				ConfigWidth  int    `json:"config_width"`
				ConfigHeight int    `json:"config_height"`
			} `json:"display_resources"`
			AccessibilityCaption  string `json:"accessibility_caption"`
			IsVideo               bool   `json:"is_video"`
			ShouldLogClientEvent  bool   `json:"should_log_client_event"`
			TrackingToken         string `json:"tracking_token"`
			EdgeMediaToTaggedUser struct {
				Edges []interface{} `json:"edges"`
			} `json:"edge_media_to_tagged_user"`
			EdgeMediaToCaption struct {
				Edges []struct {
					Node struct {
						Text string `json:"text"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_media_to_caption"`
			CaptionIsEdited    bool `json:"caption_is_edited"`
			HasRankedComments  bool `json:"has_ranked_comments"`
			EdgeMediaToComment struct {
				Count    int64    `json:"count"`
				PageInfo PageInfo `json:"page_info"`
				Edges    []struct {
					Node struct {
						ID              string `json:"id"`
						Text            string `json:"text"`
						CreatedAt       int64  `json:"created_at"`
						DidReportAsSpam bool   `json:"did_report_as_spam"`
						Owner           struct {
							ID            string `json:"id"`
							IsVerified    bool   `json:"is_verified"`
							ProfilePicURL string `json:"profile_pic_url"`
							Username      string `json:"username"`
						} `json:"owner"`
						ViewerHasLiked bool `json:"viewer_has_liked"`
						EdgeLikedBy    struct {
							Count int64 `json:"count"`
						} `json:"edge_liked_by"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_media_to_comment"`
			CommentsDisabled     bool  `json:"comments_disabled"`
			TakenAtTimestamp     int64 `json:"taken_at_timestamp"`
			EdgeMediaPreviewLike struct {
				Count int64         `json:"count"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_media_preview_like"`
			EdgeMediaToSponsorUser struct {
				Edges []interface{} `json:"edges"`
			} `json:"edge_media_to_sponsor_user"`
			Location struct {
				ID            string `json:"id"`
				HasPublicPage bool   `json:"has_public_page"`
				Name          string `json:"name"`
				Slug          string `json:"slug"`
				AddressJSON   string `json:"address_json"`
			} `json:"location"`
			ViewerHasLiked             bool `json:"viewer_has_liked"`
			ViewerHasSaved             bool `json:"viewer_has_saved"`
			ViewerHasSavedToCollection bool `json:"viewer_has_saved_to_collection"`
			ViewerInPhotoOfYou         bool `json:"viewer_in_photo_of_you"`
			ViewerCanReShare           bool `json:"viewer_can_reshare"`
			Owner                      struct {
				ID                string `json:"id"`
				IsVerified        bool   `json:"is_verified"`
				ProfilePicURL     string `json:"profile_pic_url"`
				Username          string `json:"username"`
				BlockedByViewer   bool   `json:"blocked_by_viewer"`
				FollowedByViewer  bool   `json:"followed_by_viewer"`
				FullName          string `json:"full_name"`
				HasBlockedViewer  bool   `json:"has_blocked_viewer"`
				IsPrivate         bool   `json:"is_private"`
				IsUnpublished     bool   `json:"is_unpublished"`
				RequestedByViewer bool   `json:"requested_by_viewer"`
			} `json:"owner"`
			IsAd                       bool `json:"is_ad"`
			EdgeWebMediaToRelatedMedia struct {
				Edges []interface{} `json:"edges"`
			} `json:"edge_web_media_to_related_media"`
		} `json:"shortcode_media"`
	} `json:"data"`
	Status string `json:"status"`
}

type HashTagResponse struct {
	Data struct {
		HashTag *struct {
			ID                 string `json:"id"`
			Name               string `json:"name"`
			AllowFollowing     bool   `json:"allow_following"`
			IsFollowing        bool   `json:"is_following"`
			IsTopMediaOnly     bool   `json:"is_top_media_only"`
			ProfilePicURL      string `json:"profile_pic_url"`
			EdgeHashTagToMedia struct {
				Count    int64    `json:"count"`
				PageInfo PageInfo `json:"page_info"`
				Edges    []struct {
					Node struct {
						CommentsDisabled   bool   `json:"comments_disabled"`
						Typename           string `json:"__typename"`
						ID                 string `json:"id"`
						EdgeMediaToCaption struct {
							Edges []struct {
								Node struct {
									Text string `json:"text"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_caption"`
						ShortCode          string `json:"shortcode"`
						EdgeMediaToComment struct {
							Count int64 `json:"count"`
						} `json:"edge_media_to_comment"`
						TakenAtTimestamp int64 `json:"taken_at_timestamp"`
						Dimensions       struct {
							Height int `json:"height"`
							Width  int `json:"width"`
						} `json:"dimensions"`
						DisplayURL  string `json:"display_url"`
						EdgeLikedBy struct {
							Count int64 `json:"count"`
						} `json:"edge_liked_by"`
						EdgeMediaPreviewLike struct {
							Count int64 `json:"count"`
						} `json:"edge_media_preview_like"`
						Owner struct {
							ID string `json:"id"`
						} `json:"owner"`
						ThumbnailSrc       string `json:"thumbnail_src"`
						ThumbnailResources []struct {
							Src          string `json:"src"`
							ConfigWidth  int    `json:"config_width"`
							ConfigHeight int    `json:"config_height"`
						} `json:"thumbnail_resources"`
						IsVideo              bool   `json:"is_video"`
						AccessibilityCaption string `json:"accessibility_caption"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_hashtag_to_media"`
			EdgeHashTagToTopPosts struct {
				Edges []struct {
					Node struct {
						Typename           string `json:"__typename"`
						ID                 string `json:"id"`
						EdgeMediaToCaption struct {
							Edges []struct {
								Node struct {
									Text string `json:"text"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_caption"`
						ShortCode          string `json:"shortcode"`
						EdgeMediaToComment struct {
							Count int64 `json:"count"`
						} `json:"edge_media_to_comment"`
						TakenAtTimestamp int64 `json:"taken_at_timestamp"`
						Dimensions       struct {
							Height int `json:"height"`
							Width  int `json:"width"`
						} `json:"dimensions"`
						DisplayURL  string `json:"display_url"`
						EdgeLikedBy struct {
							Count int64 `json:"count"`
						} `json:"edge_liked_by"`
						EdgeMediaPreviewLike struct {
							Count int64 `json:"count"`
						} `json:"edge_media_preview_like"`
						Owner struct {
							ID string `json:"id"`
						} `json:"owner"`
						ThumbnailSrc       string `json:"thumbnail_src"`
						ThumbnailResources []struct {
							Src          string `json:"src"`
							ConfigWidth  int    `json:"config_width"`
							ConfigHeight int    `json:"config_height"`
						} `json:"thumbnail_resources"`
						IsVideo              bool   `json:"is_video"`
						AccessibilityCaption string `json:"accessibility_caption"`
					} `json:"node,omitempty"`
				} `json:"edges"`
			} `json:"edge_hashtag_to_top_posts"`
			EdgeHashTagToContentAdvisory struct {
				Count int64         `json:"count"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_hashtag_to_content_advisory"`
		} `json:"hashtag"`
	} `json:"data"`
	Status string `json:"status"`
}

type UserInfoResponse struct {
	LoggingPageID         string `json:"logging_page_id"`
	ShowSuggestedProfiles bool   `json:"show_suggested_profiles"`
	GraphQL               struct {
		User *struct {
			Biography              string `json:"biography"`
			BlockedByViewer        bool   `json:"blocked_by_viewer"`
			CountryBlock           bool   `json:"country_block"`
			ExternalURL            string `json:"external_url"`
			ExternalURLLinkShimmed string `json:"external_url_linkshimmed"`
			EdgeFollowedBy         struct {
				Count int64 `json:"count"`
			} `json:"edge_followed_by"`
			FollowedByViewer bool `json:"followed_by_viewer"`
			EdgeFollow       struct {
				Count int64 `json:"count"`
			} `json:"edge_follow"`
			FollowsViewer        bool   `json:"follows_viewer"`
			FullName             string `json:"full_name"`
			HasChannel           bool   `json:"has_channel"`
			HasBlockedViewer     bool   `json:"has_blocked_viewer"`
			HighlightReelCount   int    `json:"highlight_reel_count"`
			HasRequestedViewer   bool   `json:"has_requested_viewer"`
			ID                   string `json:"id"`
			IsBusinessAccount    bool   `json:"is_business_account"`
			IsJoinedRecently     bool   `json:"is_joined_recently"`
			BusinessCategoryName string `json:"business_category_name"`
			IsPrivate            bool   `json:"is_private"`
			IsVerified           bool   `json:"is_verified"`
			EdgeMutualFollowedBy struct {
				Count int64         `json:"count"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_mutual_followed_by"`
			ProfilePicURL          string      `json:"profile_pic_url"`
			ProfilePicURLHd        string      `json:"profile_pic_url_hd"`
			RequestedByViewer      bool        `json:"requested_by_viewer"`
			Username               string      `json:"username"`
			ConnectedFbPage        interface{} `json:"connected_fb_page"`
			EdgeFelixVideoTimeline struct {
				Count    int64         `json:"count"`
				PageInfo PageInfo      `json:"page_info"`
				Edges    []interface{} `json:"edges"`
			} `json:"edge_felix_video_timeline"`
			EdgeOwnerToTimelineMedia struct {
				Count    int64    `json:"count"`
				PageInfo PageInfo `json:"page_info"`
				Edges    []struct {
					Node struct {
						Typename           string `json:"__typename"`
						ID                 string `json:"id"`
						EdgeMediaToCaption struct {
							Edges []struct {
								Node struct {
									Text string `json:"text"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_caption"`
						Shortcode          string `json:"shortcode"`
						EdgeMediaToComment struct {
							Count int64 `json:"count"`
						} `json:"edge_media_to_comment"`
						CommentsDisabled bool  `json:"comments_disabled"`
						TakenAtTimestamp int64 `json:"taken_at_timestamp"`
						Dimensions       struct {
							Height int `json:"height"`
							Width  int `json:"width"`
						} `json:"dimensions"`
						DisplayURL  string `json:"display_url"`
						EdgeLikedBy struct {
							Count int64 `json:"count"`
						} `json:"edge_liked_by"`
						EdgeMediaPreviewLike struct {
							Count int64 `json:"count"`
						} `json:"edge_media_preview_like"`
						Location struct {
							ID            string `json:"id"`
							HasPublicPage bool   `json:"has_public_page"`
							Name          string `json:"name"`
							Slug          string `json:"slug"`
						} `json:"location"`
						GatingInfo   interface{} `json:"gating_info"`
						MediaPreview string      `json:"media_preview"`
						Owner        struct {
							ID       string `json:"id"`
							Username string `json:"username"`
						} `json:"owner"`
						ThumbnailSrc       string `json:"thumbnail_src"`
						ThumbnailResources []struct {
							Src          string `json:"src"`
							ConfigWidth  int    `json:"config_width"`
							ConfigHeight int    `json:"config_height"`
						} `json:"thumbnail_resources"`
						IsVideo              bool   `json:"is_video"`
						AccessibilityCaption string `json:"accessibility_caption"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_owner_to_timeline_media"`
			EdgeSavedMedia struct {
				Count    int64         `json:"count"`
				PageInfo PageInfo      `json:"page_info"`
				Edges    []interface{} `json:"edges"`
			} `json:"edge_saved_media"`
			EdgeMediaCollections struct {
				Count    int64         `json:"count"`
				PageInfo PageInfo      `json:"page_info"`
				Edges    []interface{} `json:"edges"`
			} `json:"edge_media_collections"`
		} `json:"user"`
	} `json:"graphql"`
	FelixOnBoardingVideoResources struct {
		Mp4    string `json:"mp4"`
		Poster string `json:"poster"`
	} `json:"felix_onboarding_video_resources"`
	ShowFollowDialog bool `json:"show_follow_dialog"`
}

type ShortCodeMediaLikersResponse struct {
	Data struct {
		ShortCodeMedia *struct {
			ID          string `json:"id"`
			ShortCode   string `json:"shortcode"`
			EdgeLikedBy struct {
				Count    int64    `json:"count"`
				PageInfo PageInfo `json:"page_info"`
				Edges    []struct {
					Node struct {
						ID                string `json:"id"`
						Username          string `json:"username"`
						FullName          string `json:"full_name"`
						ProfilePicURL     string `json:"profile_pic_url"`
						IsPrivate         bool   `json:"is_private"`
						IsVerified        bool   `json:"is_verified"`
						FollowedByViewer  bool   `json:"followed_by_viewer"`
						RequestedByViewer bool   `json:"requested_by_viewer"`
						Reel              struct {
							ID              string      `json:"id"`
							ExpiringAt      int64       `json:"expiring_at"`
							LatestReelMedia interface{} `json:"latest_reel_media"`
							Seen            interface{} `json:"seen"`
							Owner           struct {
								Typename      string `json:"__typename"`
								ID            string `json:"id"`
								ProfilePicURL string `json:"profile_pic_url"`
								Username      string `json:"username"`
							} `json:"owner"`
						} `json:"reel"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_liked_by"`
		} `json:"shortcode_media"`
	} `json:"data"`
	Status string `json:"status"`
}

type WatchedStoryEntry struct {
	Item     stories.StoryItem `json:"item" bson:"item"`
	Duration [2]int64          `json:"duration" bson:"duration"`
}

func MakeWatchedStoryRequest(belongs *Story, stories ...WatchedStoryEntry) map[string][]string {
	var m = make(map[string][]string)
	for _, s := range stories {
		m[fmt.Sprintf("%s_%s", s.Item.ID, belongs.ID)] = []string{fmt.Sprintf("%d_%d", s.Duration[0], s.Duration[1])}
	}
	return m
}

type HashTagStoriesResponse struct {
	Story  *Story `json:"story"`
	Status string `json:"status"`
}

type UserStoriesResponse struct {
	Broadcast interface{} `json:"broadcast"`
	Reel      Reel        `json:"reel"`
	Status    string      `json:"status"`
}
type FriendshipStatus struct {
	Following       bool `json:"following"`
	FollowedBy      bool `json:"followed_by"`
	Blocking        bool `json:"blocking"`
	Muting          bool `json:"muting"`
	IsPrivate       bool `json:"is_private"`
	IncomingRequest bool `json:"incoming_request"`
	OutgoingRequest bool `json:"outgoing_request"`
	IsBestie        bool `json:"is_bestie"`
}
type Reel struct {
	ID              int                 `json:"id"`
	LatestReelMedia int                 `json:"latest_reel_media"`
	ExpiringAt      int64               `json:"expiring_at"`
	Seen            int                 `json:"seen"`
	CanReply        bool                `json:"can_reply"`
	CanReshare      bool                `json:"can_reshare"`
	ReelType        string              `json:"reel_type"`
	User            stories.User        `json:"user"`
	Items           []stories.StoryItem `json:"items"`
	PrefetchCount   int                 `json:"prefetch_count"`
	HasBestiesMedia bool                `json:"has_besties_media"`
	MediaCount      int                 `json:"media_count"`
}

type Owner struct {
	Type               string `json:"type"`
	Pk                 string `json:"pk"`
	Name               string `json:"name"`
	ProfilePicURL      string `json:"profile_pic_url"`
	ProfilePicUsername string `json:"profile_pic_username"`
}

type Story struct {
	ID                  string              `json:"id"`
	LatestReelMedia     int                 `json:"latest_reel_media"`
	ExpiringAt          int64               `json:"expiring_at"`
	Seen                int                 `json:"seen"`
	CanReply            bool                `json:"can_reply"`
	CanReshare          bool                `json:"can_reshare"`
	ReelType            string              `json:"reel_type"`
	Owner               Owner               `json:"owner"`
	Items               []stories.StoryItem `json:"items"`
	PrefetchCount       int                 `json:"prefetch_count"`
	UniqueIntegerReelID int64               `json:"unique_integer_reel_id"`
	Muted               bool                `json:"muted"`
}

type HashTagSummary struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	AllowFollowing     bool   `json:"allow_following"`
	IsFollowing        bool   `json:"is_following"`
	IsTopMediaOnly     bool   `json:"is_top_media_only"`
	ProfilePicURL      string `json:"profile_pic_url"`
	EdgeHashTagToMedia struct {
		Count int `json:"count"`
	} `json:"edge_hashtag_to_media"`
}

type PostSummary struct {
	Typename   string `json:"__typename"`
	ID         string `json:"id"`
	Shortcode  string `json:"shortcode"`
	Dimensions struct {
		Height int `json:"height"`
		Width  int `json:"width"`
	} `json:"dimensions"`
	DisplayURL       string `json:"display_url"`
	DisplayResources []struct {
		Src          string `json:"src"`
		ConfigWidth  int    `json:"config_width"`
		ConfigHeight int    `json:"config_height"`
	} `json:"display_resources"`
	IsVideo               bool   `json:"is_video"`
	TrackingToken         string `json:"tracking_token"`
	EdgeMediaToTaggedUser struct {
		Edges []struct {
			Node struct {
				User struct {
					FullName      string `json:"full_name"`
					ID            string `json:"id"`
					IsVerified    bool   `json:"is_verified"`
					ProfilePicURL string `json:"profile_pic_url"`
					Username      string `json:"username"`
				} `json:"user"`
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"edge_media_to_tagged_user"`
	EdgeMediaToCaption struct {
		Edges []struct {
			Node struct {
				Text string `json:"text"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"edge_media_to_caption"`
	CaptionIsEdited    bool `json:"caption_is_edited"`
	HasRankedComments  bool `json:"has_ranked_comments"`
	EdgeMediaToComment struct {
		Count    int `json:"count"`
		PageInfo struct {
			HasNextPage bool        `json:"has_next_page"`
			EndCursor   interface{} `json:"end_cursor"`
		} `json:"page_info"`
		Edges []struct {
			Node struct {
				ID              string `json:"id"`
				Text            string `json:"text"`
				CreatedAt       int    `json:"created_at"`
				DidReportAsSpam bool   `json:"did_report_as_spam"`
				Owner           struct {
					ID            string `json:"id"`
					IsVerified    bool   `json:"is_verified"`
					ProfilePicURL string `json:"profile_pic_url"`
					Username      string `json:"username"`
				} `json:"owner"`
				ViewerHasLiked bool `json:"viewer_has_liked"`
				EdgeLikedBy    struct {
					Count int `json:"count"`
				} `json:"edge_liked_by"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"edge_media_to_comment"`
	CommentsDisabled     bool `json:"comments_disabled"`
	TakenAtTimestamp     int  `json:"taken_at_timestamp"`
	EdgeMediaPreviewLike struct {
		Count int           `json:"count"`
		Edges []interface{} `json:"edges"`
	} `json:"edge_media_preview_like"`
	EdgeMediaToSponsorUser struct {
		Edges []interface{} `json:"edges"`
	} `json:"edge_media_to_sponsor_user"`
	Location                   interface{} `json:"location"`
	ViewerHasLiked             bool        `json:"viewer_has_liked"`
	ViewerHasSaved             bool        `json:"viewer_has_saved"`
	ViewerHasSavedToCollection bool        `json:"viewer_has_saved_to_collection"`
	ViewerInPhotoOfYou         bool        `json:"viewer_in_photo_of_you"`
	ViewerCanReshare           bool        `json:"viewer_can_reshare"`
	Owner                      struct {
		ID                string `json:"id"`
		IsVerified        bool   `json:"is_verified"`
		ProfilePicURL     string `json:"profile_pic_url"`
		Username          string `json:"username"`
		BlockedByViewer   bool   `json:"blocked_by_viewer"`
		FollowedByViewer  bool   `json:"followed_by_viewer"`
		FullName          string `json:"full_name"`
		HasBlockedViewer  bool   `json:"has_blocked_viewer"`
		IsPrivate         bool   `json:"is_private"`
		IsUnpublished     bool   `json:"is_unpublished"`
		RequestedByViewer bool   `json:"requested_by_viewer"`
	} `json:"owner"`
	IsAd                       bool `json:"is_ad"`
	EdgeWebMediaToRelatedMedia struct {
		Edges []interface{} `json:"edges"`
	} `json:"edge_web_media_to_related_media"`
	EdgeSidecarToChildren struct {
		Edges []struct {
			Node struct {
				Typename   string `json:"__typename"`
				ID         string `json:"id"`
				Shortcode  string `json:"shortcode"`
				Dimensions struct {
					Height int `json:"height"`
					Width  int `json:"width"`
				} `json:"dimensions"`
				GatingInfo       interface{} `json:"gating_info"`
				MediaPreview     string      `json:"media_preview"`
				DisplayURL       string      `json:"display_url"`
				DisplayResources []struct {
					Src          string `json:"src"`
					ConfigWidth  int    `json:"config_width"`
					ConfigHeight int    `json:"config_height"`
				} `json:"display_resources"`
				AccessibilityCaption  string `json:"accessibility_caption"`
				IsVideo               bool   `json:"is_video"`
				ShouldLogClientEvent  bool   `json:"should_log_client_event"`
				TrackingToken         string `json:"tracking_token"`
				EdgeMediaToTaggedUser struct {
					Edges []interface{} `json:"edges"`
				} `json:"edge_media_to_tagged_user"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"edge_sidecar_to_children"`
}
