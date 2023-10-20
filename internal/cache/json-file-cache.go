package cache

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type CacheItem[T any, TFilter any] struct {
	// The data that is cached
	Data T `json:"data"`
	// The time when the data was cached
	CachedAt time.Time `json:"cached_at"`

	// Filter criteria, if any
	Filter TFilter `json:"filter,omitempty"`
}

// A cache that stores data in a JSON file
type JsonFileCache[T any, TFilter comparable] struct {
}

// Get the data from the cache
func (*JsonFileCache[T, TFilter]) Get(filter TFilter, filePath string, maxAge time.Duration) (*T, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("jsonFileCache.Get could not open file: %w", err)
	}
	defer file.Close()

	// Read the contents of the file
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("jsonFileCache.Get could not read file: %w", err)
	}

	// Unmarshal the JSON data into a CacheItem
	var cacheItems []CacheItem[T, TFilter]
	err = json.Unmarshal(fileContent, &cacheItems)
	if err != nil {
		return nil, fmt.Errorf("jsonFileCache.Get could not unmarshal file: %w", err)
	}

	// Find the cache item that matches the filter
	var cacheItem *CacheItem[T, TFilter]
	for _, item := range cacheItems {
		if item.Filter == filter {
			cacheItem = &item
			break
		}
	}

	// Check if the cache item was found
	if cacheItem == nil {
		return nil, fmt.Errorf("cache item not found")
	}

	// Check if the cache is still valid
	if time.Since(cacheItem.CachedAt) > maxAge {
		return nil, fmt.Errorf("cache is too old")
	}

	return &cacheItem.Data, nil
}

// Save the data to the cache file.  The cache file contains a list of CacheItem where the key is the filter.
func (*JsonFileCache[T, TFilter]) Save(filter TFilter, filePath string, data *T, maxAge time.Duration) error {
	// Create the cache item
	cacheItem := CacheItem[T, TFilter]{
		Data:     *data,
		CachedAt: time.Now(),
		Filter:   filter,
	}

	// Open the file
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("jsonFileCache.Save could not open file: %w", err)
	}
	defer file.Close()

	// Read the contents of the file
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("jsonFileCache.Save could not read file: %w", err)
	}

	cacheItemFound := false
	var cacheItems []CacheItem[T, TFilter]

	// Check if the file is empty
	if len(fileContent) != 0 {
		// Unmarshal the JSON data into a CacheItem
		err = json.Unmarshal(fileContent, &cacheItems)
		if err != nil {
			return fmt.Errorf("jsonFileCache.Save could not unmarshal file [%s]: %w", filePath, err)
		}

		// Find the cache item that matches the filter
		for i, item := range cacheItems {
			if item.Filter == filter {
				cacheItems[i] = cacheItem
				cacheItemFound = true
				break
			}
		}
	}
	if cacheItemFound {
		// If the cache item was found, replace it in the list

	} else {
		// If the cache item was not found, append it to the list
		cacheItems = append(cacheItems, cacheItem)
	}

	// Also do a cleanup of the cache items
	// Remove any cache items that are older than the max age
	var newCacheItems []CacheItem[T, TFilter]
	for _, item := range cacheItems {
		if time.Since(item.CachedAt) <= maxAge {
			newCacheItems = append(newCacheItems, item)
		}
	}

	// Replace the cache items with the new list
	cacheItems = newCacheItems

	// Marshal the data
	bytes, err := json.Marshal(cacheItems)
	if err != nil {
		return fmt.Errorf("jsonFileCache could not marshal file: %w", err)
	}

	// Write the file
	// seek to the beginning of the file
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}

	// Truncate the file to the length of the new bytes
	err = file.Truncate(int64(len(bytes)))
	if err != nil {
		return err
	}

	return nil
}

// Delete the data from the cache
func (*JsonFileCache[T, TFilter]) Delete(filePath string) error {
	// Open the file
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("jsonFileCache.Delete could not open file: %w", err)
	}
	defer file.Close()

	// Write the file
	// seek to the beginning of the file
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	empty := "[]"
	bytes := []byte(empty)

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}

	// Truncate the file to the length of the new bytes
	err = file.Truncate(int64(len(bytes)))
	if err != nil {
		return err
	}

	return nil
}
