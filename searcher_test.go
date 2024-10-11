package czdb

import "testing"

const ipv4Path = "/Users/spectator/Downloads/czdb/cz88_public_v4.czdb"
const ipv6Path = "/Users/spectator/Downloads/czdb/cz88_public_v6.czdb"
const ipv4 = "1.64.219.93"
const ipv6 = "2001:4:113:0:0:0:0:0"
const ipv6_2 = "2602:feda:30:cafe:68f1:6bff:fea6:2e36"

func TestIPV4Memory(t *testing.T) {

	searcher, err := NewDBSearcher(ipv4Path, "iLCLCRA5ijY5PqpkoU6/qQ==", QueryType_Memory)
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

func TestIPV6Memory(t *testing.T) {

	searcher, err := NewDBSearcher(ipv6Path, "iLCLCRA5ijY5PqpkoU6/qQ==", QueryType_Memory)
	if err != nil {
		t.Fatalf("NewDBSearcher error: %v", err)
		return
	}
	region, err := searcher.Search(ipv6)
	if err != nil {
		t.Fatalf("Search error: %v", err)
		return
	}

	t.Logf("Region: %v", region)
}
func TestIPV4Btree(t *testing.T) {

	searcher, err := NewDBSearcher(ipv4Path, "iLCLCRA5ijY5PqpkoU6/qQ==", QueryType_Btree)
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

func TestIPV6Btree(t *testing.T) {

	searcher, err := NewDBSearcher(ipv6Path, "iLCLCRA5ijY5PqpkoU6/qQ==", QueryType_Btree)
	if err != nil {
		t.Fatalf("NewDBSearcher error: %v", err)
		return
	}
	region, err := searcher.Search(ipv6)
	if err != nil {
		t.Fatalf("Search error: %v", err)
		return
	}

	t.Logf("Region: %v", region)
}
