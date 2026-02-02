package test

import (
	"fmt"
	"math"
	"testing"
)

type Cluster struct {
	data    map[uint32]struct{}
	buckets [20]uint32
}

func (d *Cluster) String() string {
	return fmt.Sprintf("%v", d.buckets)
}

func (d *Cluster) Add(value uint32) {
	d.data[value] = struct{}{}
	per := int(math.Ceil(20 / float64(len(d.data))))
	tmps := map[uint32]int{}
	for pos := 0; pos < 20; pos++ {
		item := d.buckets[pos]
		if item <= 0 || tmps[item] >= per {
			item = value
		}
		if tmps[item] < per {
			tmps[item]++
			d.buckets[pos] = item
		}
	}
}

func (d *Cluster) Del(value uint32) {
	delete(d.data, value)
	if len(d.data) <= 0 {
		for pos := 0; pos < 20; pos++ {
			d.buckets[pos] = 0
		}
		return
	}

	per := int(math.Ceil(20 / float64(len(d.data))))
	tmps := map[uint32]int{}
	for _, val := range d.buckets {
		if val != value {
			tmps[val]++
		}
	}
	pos := 0
	for key := range d.data {
		for ; pos < 20; pos++ {
			if item := d.buckets[pos]; item == value {
				if tmps[key] >= per {
					break
				}
				tmps[key]++
				d.buckets[pos] = key
			}
		}
	}
}

func TestCluster(t *testing.T) {
	aa := &Cluster{
		data: make(map[uint32]struct{}),
	}

	aa.Add(123)
	t.Log(aa.String())
	aa.Add(654)
	t.Log(aa.String())
	aa.Add(666)
	t.Log(aa.String())
	aa.Add(999)
	t.Log(aa.String())
	t.Log("---------delete---------")
	aa.Del(654)
	t.Log(aa.String())
	aa.Add(654)
	t.Log(aa.String())
}
