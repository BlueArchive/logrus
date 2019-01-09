package logrus

// A hook to be fired when logging on the logging levels returned from
// `Levels()` on your implementation of the interface. Note that this is not
// fired in a goroutine or a channel with workers, you should handle such
// functionality yourself if your call is non-blocking and you don't wish for
// the logging calls for levels returned from `Levels()` to block.
type Hook interface {
	Name() string
	Levels() []Level
	Fire(*Entry) error
}

// Internal type for storing the hooks on a logger instance.
type LevelHooks map[Level][]Hook

// Add a hook to an instance of logger. This is called with
// `log.Hooks.Add(new(MyHook))` where `MyHook` implements the `Hook` interface.
func (hooks LevelHooks) Add(hook Hook) {
	for _, level := range hook.Levels() {
		hooks[level] = append(hooks[level], hook)
	}
}

type HookFailure struct {
	Hook Hook
	Err error
}
type HookFailures []HookFailure

// Fire all the hooks for the passed level. Used by `entry.log` to fire
// appropriate hooks for a log entry.
// don't stop if a hook fails to fire - other hooks should get their chance
// return map of failed hook to err
func (hooks LevelHooks) Fire(level Level, entry *Entry) (hookFailures *HookFailures) {
	// aggregate failed hooks; we want all of them to fire
	for _, hook := range hooks[level] {
		if err := hook.Fire(entry); err != nil {
			if hookFailures == nil {
				hookFailures = &HookFailures{}
			}
			*hookFailures = append(*hookFailures, HookFailure{Hook:hook, Err:err})
		}
	}

	return hookFailures
}
