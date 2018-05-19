package goinsta

import (
	"encoding/json"
	"fmt"
)

type Hashtags struct {
	inst *Instagram

	Name string `json:"name"`

	Sections []struct {
		LayoutType    string `json:"layout_type"`
		LayoutContent struct {
			// F*ck you instagram.
			// Why you do this f*cking horribly structure?!?
			// Media []Media IS EASY. CHECK IT!
			Medias []struct {
				Item Item `json:"media"`
			} `json:"medias"`
		} `json:"layout_content"`
		FeedType        string `json:"feed_type"`
		ExploreItemInfo struct {
			NumColumns      int  `json:"num_columns"`
			TotalNumColumns int  `json:"total_num_columns"`
			AspectRatio     int  `json:"aspect_ratio"`
			Autoplay        bool `json:"autoplay"`
		} `json:"explore_item_info"`
	} `json:"sections"`
	MediaCount          int     `json:"media_count"`
	ID                  int64   `json:"id"`
	MoreAvailable       bool    `json:"more_available"`
	NextID              string  `json:"next_max_id"`
	NextPage            int     `json:"next_page"`
	NextMediaIds        []int64 `json:"next_media_ids"`
	AutoLoadMoreEnabled bool    `json:"auto_load_more_enabled"`
	Status              string  `json:"status"`
}

// Sync updates Hashtags information preparing it to Next call.
func (h *Hashtags) Sync() error {
	insta := h.inst

	body, err := insta.sendSimpleRequest(urlTagSync, h.Name)
	if err == nil {
		var resp struct {
			Name       string `json:"name"`
			ID         int64  `json:"id"`
			MediaCount int    `json:"media_count"`
		}
		err = json.Unmarshal(body, &resp)
		if err == nil {
			h.Name = resp.Name
			h.ID = ID
			h.MediaCount = resp.MediaCount
		}
	}
	return err
}

// Next paginates over hashtag page (xd).
func (h *Hashtags) Next() bool {
	if h.err != nil {
		return false
	}
	insta := h.inst
	name := h.Name
	data, err := insta.prepareData(
		&reqOptions{
			Query: map[string]string{
				"max_id":     h.NextID,
				"rank_token": insta.rankToken,
				"page":       fmt.Sprintf("%d", h.NextPage),
			},
			Endpoint: fmt.Sprintf(urlTagContent, name),
			IsPost:   false,
		},
	)
	if err == nil {
		ht := &Hashtags{}
		err = json.Unmarshal(body, &ht)
		if err == nil {
			*h = ht
			h.inst = insta
			h.Name = name
			if !h.MoreAvailible {
				h.err = ErrNoMore
			}
		}
	}
	h.err = err
	return false
}

// TODO: func (h *Hashtags) Stories()
