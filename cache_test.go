package main

import (
  "testing"
)

func TestIdGenerator(t *testing.T) {
  cache := NewCacheManager(1000)

  point := cache.GenerateId(-31.123321, 24.666555)
  expect := "-31.123321~24.666555"

  if(point != expect) {
    t.Fatalf("Expected %s but got %s", expect, point)
  }
}

func TestPutAndGetFromCache(t *testing.T) {
    expected := 666
    cacheManager := NewCacheManager(20000)
    cacheManager.Put(expected, 0.00, 0.00)
    lastId := cacheManager.Get(0.00, 0.00)

    if(lastId != expected) {
      t.Fatalf("Expected %s but got %s", expected, lastId)
    }
}

func TestCacheMustClean(t *testing.T) {
    expected := 666
    cacheManager := NewCacheManager(0)
    cacheManager.Put(expected, 0.00, 0.00)
    cacheManager.Clean()
    lastId := cacheManager.Get(0.00, 0.00)

    if(lastId == expected) {
      t.Fatalf("Did not clean up cache")
    }
}
