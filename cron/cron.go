package cron

import (
	"context"
	kitContext "github.com/exluap/kit/context"
	kitLog "github.com/exluap/kit/log"
	"github.com/go-co-op/gocron"
	"sync"
	"time"
)

// Action each cron job must implement action to be run by cron manager
// ctxFn - function returning call context
type Action func(ctxFn func() context.Context)

// Cron allows setup cron parameters
type Cron interface {
	// Every indicates how often cron is executed
	Every(i time.Duration) Cron
	// At indicates cron runs once at the given time
	// Every call ignored
	At(t time.Time) Cron
	// Action takes action to run
	Action(a Action) Cron
	// UnderUser allows to specify a user under which action is executed
	UnderUser(userId, username string) Cron
}

// Manager allows manage all the cron jobs in centralized place
type Manager interface {
	// Add adds a ne cron
	Add(ctx context.Context, name string) Cron
	// Start starts all the jobs
	Start(ctx context.Context)
	// Stop stops all the jobs
	Stop(ctx context.Context)
}

type managerIml struct {
	sync.RWMutex
	items []*cronImpl
	lFn   kitLog.CLoggerFunc
}

// NewManager creates a new cron manager
func NewManager(lFn kitLog.CLoggerFunc) Manager {
	return &managerIml{
		items: []*cronImpl{},
		lFn:   lFn,
	}
}

func (m *managerIml) Add(ctx context.Context, name string) Cron {
	m.Lock()
	defer m.Unlock()
	c := newCron(name)
	m.items = append(m.items, c)
	return c
}

func (m *managerIml) Start(ctx context.Context) {

	l := m.lFn().C(ctx).Cmp("cron").Mth("start")

	m.RLock()
	defer m.RUnlock()
	for _, c := range m.items {
		l.F(kitLog.FF{"name": c.name})
		if err := c.start(); err != nil {
			l.E(err).St().Err()
		}
		l.Dbg("ok")
	}

}

func (m *managerIml) Stop(ctx context.Context) {

	l := m.lFn().C(ctx).Cmp("cron").Mth("stop")

	m.RLock()
	defer m.RUnlock()
	for _, c := range m.items {
		l.F(kitLog.FF{"name": c.name})
		c.stop()
		l.Dbg("ok")
	}

}

type cronImpl struct {
	sync.RWMutex
	scheduler        *gocron.Scheduler
	action           Action
	userId, username string
	name             string
}

func newCron(name string) *cronImpl {
	return &cronImpl{
		scheduler: gocron.NewScheduler(time.UTC),
		name:      name,
	}
}

func (c *cronImpl) Every(i time.Duration) Cron {
	c.Lock()
	defer c.Unlock()
	c.scheduler = c.scheduler.Every(i)
	return c
}

func (c *cronImpl) At(t time.Time) Cron {
	c.Lock()
	defer c.Unlock()
	c.scheduler = c.scheduler.Every(1).Day().StartAt(t)
	return c
}

func (c *cronImpl) Action(a Action) Cron {
	c.Lock()
	defer c.Unlock()
	c.action = a
	return c
}

func (c *cronImpl) UnderUser(userId, username string) Cron {
	c.Lock()
	defer c.Unlock()
	c.userId, c.username = userId, username
	return c
}

func (c *cronImpl) start() error {
	c.RLock()
	defer c.RUnlock()

	ctxFn := func() context.Context {
		return kitContext.NewRequestCtx().
			Job().
			WithNewRequestId().
			WithUser(c.userId, c.username).
			ToContext(context.Background())
	}

	_, err := c.scheduler.Tag(c.name).Do(c.action, ctxFn)
	if err != nil {
		return ErrCronStart(err, c.name)
	}

	c.scheduler.StartAsync()

	return nil
}

func (c *cronImpl) stop() {
	c.scheduler.Stop()
}
