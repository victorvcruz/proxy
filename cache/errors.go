package cache

import "fmt"

type CacheNotFoundError struct{}

func (e *CacheNotFoundError) Error() string {
	return fmt.Sprintf("not found")
}
