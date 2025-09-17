package queue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrChatRunnerLimit = errors.New("queue: chat runner limit reached")

// RunnerManager 按 chatID 动态管理 Worker/Mover 生命周期
type RunnerManager struct {
	rdb      *redis.Client
	cfg      *Config
	limiter  *Limiter
	tg       TelegramProvider
	registry BotRegistry
	failure  *FailureTracker

	mu      sync.Mutex
	runners map[int64]*runner
}

func NewRunnerManager(rdb *redis.Client, cfg *Config, limiter *Limiter, failure *FailureTracker, tg TelegramProvider, registry BotRegistry) *RunnerManager {
	return &RunnerManager{
		rdb:      rdb,
		cfg:      cfg,
		limiter:  limiter,
		failure:  failure,
		tg:       tg,
		registry: registry,
		runners:  make(map[int64]*runner),
	}
}

func (m *RunnerManager) EnsureRunner(ctx context.Context, chatID int64) error {
	if chatID == 0 {
		return fmt.Errorf("chatID cannot be zero")
	}
	m.mu.Lock()
	if r, ok := m.runners[chatID]; ok {
		r.touch()
		m.mu.Unlock()
		return nil
	}
	if m.cfg.MaxChatRunners > 0 && len(m.runners) >= m.cfg.MaxChatRunners {
		m.mu.Unlock()
		return ErrChatRunnerLimit
	}
	r := newRunner(m, chatID)
	m.runners[chatID] = r
	m.mu.Unlock()

	if err := r.start(ctx); err != nil {
		m.mu.Lock()
		delete(m.runners, chatID)
		m.mu.Unlock()
		return err
	}
	return nil
}

func (m *RunnerManager) Shutdown() {
	m.mu.Lock()
	runners := make([]*runner, 0, len(m.runners))
	for _, r := range m.runners {
		runners = append(runners, r)
	}
	m.runners = make(map[int64]*runner)
	m.mu.Unlock()

	for _, r := range runners {
		r.stop()
	}
}

func (m *RunnerManager) Run(ctx context.Context) {
	if m.cfg.IdleRunnerTTLMs <= 0 {
		<-ctx.Done()
		return
	}
	interval := time.Duration(m.cfg.IdleRunnerTTLMs/2) * time.Millisecond
	if interval <= 0 {
		interval = time.Minute
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.collectIdle()
		}
	}
}

func (m *RunnerManager) collectIdle() {
	ttlMs := int64(m.cfg.IdleRunnerTTLMs)
	if ttlMs <= 0 {
		return
	}
	now := time.Now().UnixMilli()
	var toStop []*runner

	m.mu.Lock()
	for chatID, r := range m.runners {
		last := r.lastActive.Load()
		if last == 0 {
			last = now
		}
		if now-last >= ttlMs {
			toStop = append(toStop, r)
			delete(m.runners, chatID)
		}
	}
	m.mu.Unlock()

	for _, r := range toStop {
		r.stop()
	}
}

func (m *RunnerManager) ensureGroup(ctx context.Context, chatID int64) {
	stream := streamReady(chatID)
	group := consumerGroup(chatID)
	_ = m.rdb.XGroupCreateMkStream(ctx, stream, group, "$").Err()
}

// runner 表示单个 chat 的工作实例
type runner struct {
	chatID     int64
	manager    *RunnerManager
	worker     *Worker
	mover      *Mover
	cancel     context.CancelFunc
	lastActive atomic.Int64
}

func newRunner(m *RunnerManager, chatID int64) *runner {
	r := &runner{chatID: chatID, manager: m}
	touch := func() { r.touch() }
	r.worker = NewWorker(m.rdb, m.cfg, m.limiter, m.failure, m.tg, m.registry, chatID, touch)
	r.mover = NewMover(m.rdb, m.cfg, chatID, touch)
	return r
}

func (r *runner) start(ctx context.Context) error {
	r.touch()
	r.manager.ensureGroup(ctx, r.chatID)
	runCtx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	go func() { _ = r.mover.Run(runCtx) }()
	go func() { _ = r.worker.Run(runCtx) }()
	return nil
}

func (r *runner) stop() {
	if r.cancel != nil {
		r.cancel()
	}
}

func (r *runner) touch() {
	r.lastActive.Store(time.Now().UnixMilli())
}
