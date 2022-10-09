package convert

import (
	"zendea/model"
)

func ToNode(node *model.Node) *model.NodeResponse {
	if node == nil {
		return nil
	}
	return &model.NodeResponse{
		NodeId:      node.ID,
		Name:        node.Name,
		Description: node.Description,
		TopicCount:  node.TopicCount,
	}
}

//func ToTags(tags []model.Tag) *[]model.TagResponse {
func ToNodes(nodes []model.Node) *[]model.NodeResponse {
	if len(nodes) == 0 {
		return nil
	}
	var ret []model.NodeResponse
	for _, node := range nodes {
		ret = append(ret, *ToNode(&node))
	}
	return &ret
}
