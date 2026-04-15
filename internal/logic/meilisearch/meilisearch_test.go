package meilisearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	ctx := context.Background()

	// 测试搜索功能
	req := &SearchReq{
		Query: "test",
		Limit: 10,
	}

	res, err := Search(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.GreaterOrEqual(t, res.Total, 0)

	t.Logf("搜索结果数量: %d", res.Total)
	if res.Total > 0 {
		t.Logf("第一个结果: %+v", res.Hits[0])
	}
}
