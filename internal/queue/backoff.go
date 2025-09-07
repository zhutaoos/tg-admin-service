package queue

import "time"

func ComputeBackoff(attempts int) time.Duration {
    if attempts <= 0 {
        return 2 * time.Second
    }
    // 2s -> 5s -> 15s -> 60s
    seq := []time.Duration{2 * time.Second, 5 * time.Second, 15 * time.Second, 60 * time.Second}
    if attempts-1 < len(seq) {
        return seq[attempts-1]
    }
    // 上限逐步扩大，最大5分钟
    d := seq[len(seq)-1] * time.Duration(attempts-len(seq)+1)
    if d > 5*time.Minute {
        d = 5 * time.Minute
    }
    return d
}

