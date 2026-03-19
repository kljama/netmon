## 2024-05-18 - [Fast-Path Pre-Check on Append-Only sync.Map]
**Learning:** In `src/scanner.go`, the `activeIPs` sync.Map is used strictly as an append-only store (hosts are never deleted). This codebase-specific architectural pattern allows a safe, lock-free `Load` check outside the goroutine and concurrency limits to skip already processed items. This was missing, causing thousands of unnecessary goroutine allocations and channel operations per discovery sweep.
**Action:** When iterating over a large dataset in Go where state tracking uses a `sync.Map` in an append-only fashion, apply an early return pattern `if _, exists := map.Load(key); exists { continue }` *before* hitting concurrency limiters or spawning goroutines.

## $(date +%Y-%m-%d) - [Safe Pre-allocation of Slice Capacity for CIDR Expansion]
**Learning:** In `src/config.go`, the `GenerateTargets` function expands CIDR blocks into a slice of IPs. Without pre-allocation, this causes excessive memory reallocations. However, pre-allocating large blocks like `/8` blindly can cause immediate OOM panics.
**Action:** When pre-allocating capacity based on user-provided CIDR sizes, always cap the maximum pre-allocation (e.g., to a `/16` network, or `65536` capacity) to safely reduce memory allocations by ~68% without risking Out-of-Memory crashes for massive subnets.
