package main

import (
  "strconv"
  "time"
  "fmt"
)

type LastIdCached struct {
  LastId int
  LastRequest int64
}



type CacheManager struct {
  CacheMap map[string]LastIdCached
  aliveTime int64
}

func NewCacheManager(aliveTime int64) *CacheManager {
  cm := new(CacheManager)
  cm.CacheMap = make(map[string]LastIdCached)
  cm.aliveTime = aliveTime
  return cm
}

func (cm *CacheManager) GenerateId(lat, lon float64) string {
  latString := strconv.FormatFloat(lat, 'f', -1, 64)
  lonString := strconv.FormatFloat(lon, 'f', -1, 64)

  return latString + "~" + lonString
}

func (cm *CacheManager) Get(lat, lon float64) int {
  key := cm.GenerateId(lat, lon)
  lastIdCached := cm.CacheMap[key]

  if &lastIdCached == nil {
    return -1;
  }

  lastIdCached.LastRequest = time.Now().UnixNano() % 1e6 / 1e3
  return lastIdCached.LastId
}

func (cm *CacheManager) Put(lastId int, lat, lon float64) {
  if len(cm.CacheMap) > 1000 {
    return
  }

  key := cm.GenerateId(lat, lon)
  time := time.Now().UnixNano() % 1e6 / 1e3
  cm.CacheMap[key] = LastIdCached{lastId, time}
}

func (cm *CacheManager) Clean() {
  toRemove := make([]string, 0, 0)

  for key, value := range cm.CacheMap {
      elapsed := (time.Now().UnixNano() % 1e6 / 1e3) - value.LastRequest
      if(elapsed >= cm.aliveTime) {
        fmt.Println("Adding item to clean cache")
        toRemove = append(toRemove, key)
      }
  }

  for _, key := range toRemove {
    fmt.Println("Removing item from cache")
    delete(cm.CacheMap, key)
  }
}
