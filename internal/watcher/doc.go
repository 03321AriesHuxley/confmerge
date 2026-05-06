// Package watcher provides lightweight file-change detection for confmerge's
// watch mode. It polls a list of config files at a configurable interval and
// emits Events whenever a file's modification time advances.
//
// Usage:
//
//	w := watcher.New(files, 500*time.Millisecond)
//	w.Start()
//	defer w.Stop()
//
//	db := watcher.NewDebouncer(200 * time.Millisecond)
//	for {
//		select {
//		case ev := <-w.Events:
//			_ = ev
//			db.Feed()
//		case <-db.Notify:
//			// re-run pipeline
//		}
//	}
package watcher
