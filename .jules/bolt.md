## 2024-05-18 - [Fast-Path Pre-Check on Append-Only sync.Map]
**Learning:** In `src/scanner.go`, the `activeIPs` sync.Map is used strictly as an append-only store (hosts are never deleted). This codebase-specific architectural pattern allows a safe, lock-free `Load` check outside the goroutine and concurrency limits to skip already processed items. This was missing, causing thousands of unnecessary goroutine allocations and channel operations per discovery sweep.
**Action:** When iterating over a large dataset in Go where state tracking uses a `sync.Map` in an append-only fashion, apply an early return pattern `if _, exists := map.Load(key); exists { continue }` *before* hitting concurrency limiters or spawning goroutines.

## 2024-05-18 - [Optimized CIDR Expansion with net/netip and Slice Pre-allocation]
**Learning:** Expanding large CIDR ranges (like `/16` or `/8`) into a dynamically appended string slice using `net.IP` causes catastrophic performance degradation due to millions of unnecessary memory allocations and array copies.
**Action:** Always pre-calculate the required slice capacity and initialize it with `make([]string, 0, capacity)`. Use Go 1.18's `net/netip` package, which is significantly faster and more allocation-friendly than the older `net.IP` interface, especially for loops handling IP iteration and comparisons.
