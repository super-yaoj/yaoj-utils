package ratings

import (
	"fmt"
	"math"
	"sort"
)

type RatingRater interface {
	Rate(rating int)
	Rating() int
	// 此前参加了多少场比赛（不算当前这场）
	Count() int
}

// 注意 list 应按照比赛成绩从高到底排序
// implementing https://codeforces.com/blog/entry/20762
// https://codeforces.com/blog/entry/77890
func CalcRating[T RatingRater](list []T) error {
	var beginnerStage = []int{500, 350, 250, 150, 100, 50}

	// probability of A winning B
	Pwin := func(ratingA int, ratingB int) float64 {
		return 1 / float64(1+math.Pow(10, float64(ratingB-ratingA)/400))
	}
	// expect place before contest
	SeedOf := func(rating int) (res float64) {
		res = 1
		for _, ctst := range list {
			res += Pwin(ctst.Rating(), rating)
		}
		return res
	}

	ExpectRatingOf := func(rating int, rank int) int {
		midRank := math.Sqrt(SeedOf(rating) * float64(rank))
		var l, r = 1, 8000
		for l < r {
			mid := (l + r) / 2
			if SeedOf(mid) < midRank {
				r = mid
			} else {
				l = mid + 1
			}
		}
		return l
	}

	// Total sum should not be more than zero.
	var deltas = make([]int, len(list))
	var sumDelta int
	for i, ctst := range list {
		deltas[i] = (ExpectRatingOf(ctst.Rating(), i+1) - ctst.Rating()) / 2
		sumDelta += deltas[i]
	}
	inc := -sumDelta/len(list) - 1
	for i := range deltas {
		deltas[i] += inc
	}

	// Sum of top-4*sqrt should be adjusted to zero.
	sortedList := make([]int, 0, len(list))
	for i := range list {
		sortedList = append(sortedList, i)
	}
	sort.Slice(sortedList, func(i, j int) bool {
		return list[sortedList[i]].Rating() > list[sortedList[j]].Rating()
	})

	s := int(math.Floor(math.Min(float64(len(list)), 4*math.Sqrt(float64(len(list))))))
	var sumd int
	for _, id := range sortedList[:s] {
		sumd += deltas[id]
	}
	inc = int(math.Ceil(math.Min(math.Max(float64(-sumd/s), -10), 0)))
	for i := range deltas {
		deltas[i] += inc
	}

	// validation
	for i := range list {
		for j := i + 1; j < len(list); j++ {
			if list[i].Rating() > list[j].Rating() {
				if !(list[i].Rating()+deltas[i] >= list[j].Rating()+deltas[j]) {
					return fmt.Errorf("first rating invariant failed")
				}
			}
			if list[i].Rating() < list[j].Rating() {
				if !(deltas[i] >= deltas[j]) {
					return fmt.Errorf("second rating invariant failed")
				}
			}
		}
	}

	for i := range list {
		if cnt := list[i].Count(); cnt < len(beginnerStage) {
			list[i].Rate(list[i].Rating() + deltas[i] + beginnerStage[cnt])
		} else {
			list[i].Rate(list[i].Rating() + deltas[i])
		}
	}

	return nil
}