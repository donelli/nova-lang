package shared_test

import (
	"nova-lang/pkg/shared"
	"testing"
)

func TestKeywordsCount(t *testing.T) {

	shared.LoadKeywords()

	if len(shared.KeywordsMap) != shared.KeywordCount {
		t.Errorf("KeywordsMap size is %d, expected %d", len(shared.KeywordsMap), shared.KeywordCount)
	}

}
