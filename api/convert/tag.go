package convert

import (
	"zendea/model"
)


func ToTag(tag *model.Tag) *model.TagResponse {
	if tag == nil {
		return nil
	}
	return &model.TagResponse{TagId: tag.ID, TagName: tag.Name}
}

func ToTags(tags []model.Tag) *[]model.TagResponse {
	if len(tags) == 0 {
		return nil
	}
	var responses []model.TagResponse
	for _, tag := range tags {
		responses = append(responses, *ToTag(&tag))
	}
	return &responses
}
