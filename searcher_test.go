package czdb

import "testing"

const ipv4Path = "/Users/spectator/Downloads/cz88_public_v4.czdb"
const ipv4 = "1.64.219.93"

func TestIPV4(t *testing.T) {

	searcher, err := NewDBSearcher(ipv4Path, "UBN0Iz3juX2qjK3sWbwcHQ==", QueryType_Memory)
	if err != nil {
		t.Fatalf("NewDBSearcher error: %v", err)
		return
	}
	region, err := searcher.Search(ipv4)
	if err != nil {
		t.Fatalf("Search error: %v", err)
		return
	}

	t.Logf("Region: %v", region)
}
