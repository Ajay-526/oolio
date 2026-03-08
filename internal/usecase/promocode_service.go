package usecase

import (
	"bufio"
	"context"
	"fmt"
	"oolio/internal/adaptors/persistance/redis"
	"os"
	"strings"
)

type PromoService struct {
	promoCounts map[string]int
	redis       *redis.RedisService
}

func NewPromoService(files []string, redis *redis.RedisService) (*PromoService, error) {

	promoCounts := make(map[string]int)

	type result struct {
		counts map[string]int
		err    error
	}

	results := make(chan result)

	for _, file := range files {

		go func(file string) {

			localCounts := make(map[string]int)
			seenInFile := make(map[string]bool)

			f, err := os.Open(file)
			if err != nil {
				results <- result{err: err}
				return
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)

			for scanner.Scan() {

				code := strings.TrimSpace(scanner.Text())

				if code == "" {
					continue
				}

				seenInFile[code] = true
			}

			for code := range seenInFile {
				localCounts[code]++
			}

			results <- result{counts: localCounts}

		}(file)
	}

	// collect results
	for i := 0; i < len(files); i++ {

		res := <-results

		if res.err != nil {
			return nil, res.err
		}

		for code, count := range res.counts {
			promoCounts[code] += count
		}
	}

	return &PromoService{
		promoCounts: promoCounts,
		redis:       redis,
	}, nil
}
func (p *PromoService) ValidatePromo(ctx context.Context, code string) bool {

	if len(code) < 8 || len(code) > 10 {
		return false
	}

	cacheKey := fmt.Sprintf("promo:valid:%s", code)

	// 1. Check Redis
	val, err := p.redis.GetValue(ctx, cacheKey)

	if err == nil {
		return val == "true"
	}

	// 2 Validate from memory
	count, ok := p.promoCounts[code]

	valid := ok && count >= 2

	// 3 Store result in Redis
	if valid {
		p.redis.SetValue(ctx, cacheKey, "true")
	} else {
		p.redis.SetValue(ctx, cacheKey, "false")
	}

	return valid
}
