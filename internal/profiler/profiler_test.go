package profiler_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/confmerge/internal/profiler"
)

func TestTrack_RecordsStage(t *testing.T) {
	p := profiler.New()
	func() {
		defer p.Track("load")()
		time.Sleep(1 * time.Millisecond)
	}()

	stages := p.Stages()
	if len(stages) != 1 {
		t.Fatalf("expected 1 stage, got %d", len(stages))
	}
	if stages[0].Name != "load" {
		t.Errorf("expected stage name 'load', got %q", stages[0].Name)
	}
	if stages[0].Duration < time.Millisecond {
		t.Errorf("expected duration >= 1ms, got %s", stages[0].Duration)
	}
}

func TestTrack_MultipleStages(t *testing.T) {
	p := profiler.New()
	for _, name := range []string{"resolve", "load", "merge", "write"} {
		name := name
		func() { defer p.Track(name)() }()
	}

	stages := p.Stages()
	if len(stages) != 4 {
		t.Fatalf("expected 4 stages, got %d", len(stages))
	}
	if stages[2].Name != "merge" {
		t.Errorf("expected third stage 'merge', got %q", stages[2].Name)
	}
}

func TestTotal_IsPositive(t *testing.T) {
	p := profiler.New()
	time.Sleep(1 * time.Millisecond)
	if p.Total() <= 0 {
		t.Error("expected positive total duration")
	}
}

func TestPrint_ContainsStageNames(t *testing.T) {
	p := profiler.New()
	func() { defer p.Track("resolve")() }()
	func() { defer p.Track("merge")() }()

	var buf bytes.Buffer
	p.Print(&buf)
	out := buf.String()

	for _, want := range []string{"resolve", "merge", "total", "Pipeline profile:"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\ngot:\n%s", want, out)
		}
	}
}

func TestStages_ReturnsCopy(t *testing.T) {
	p := profiler.New()
	func() { defer p.Track("load")() }()

	s1 := p.Stages()
	s1[0].Name = "mutated"
	s2 := p.Stages()

	if s2[0].Name == "mutated" {
		t.Error("Stages() should return a copy, not a reference")
	}
}
