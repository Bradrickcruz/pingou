package commands

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestRunShutdownRunsStepsInOrder(t *testing.T) {
	var got []string
	step := func(name string) shutdownStep {
		return shutdownStep{name, func(context.Context) error { got = append(got, name); return nil }}
	}
	err := runShutdown(context.Background(), []shutdownStep{step("a"), step("b"), step("c")})
	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if want := []string{"a", "b", "c"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("ordem = %v, want %v", got, want)
	}
}

// Etapa que estoura o prazo não impede fechar banco e lock.
func TestRunShutdownKeepsGoingAfterTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	closed := false
	err := runShutdown(ctx, []shutdownStep{
		{"checks", func(ctx context.Context) error { <-ctx.Done(); return ctx.Err() }},
		{"database", func(context.Context) error { closed = true; return nil }},
	})
	if err == nil {
		t.Fatal("err = nil, want erro (prazo estourado)")
	}
	if !closed {
		t.Fatal("banco não foi fechado depois do estouro")
	}
}

// Lock deixado pelo próprio PID (container reiniciado, PID 1 de novo) é velho.
func TestCheckExistingLockIgnoresOwnPID(t *testing.T) {
	dir := t.TempDir()
	old := lockFile
	lockFile = dir + "/pingou.lock"
	t.Cleanup(func() { lockFile = old })

	if err := writeLockForTest(lockFile); err != nil {
		t.Fatal(err)
	}
	if err := checkExistingLock(); err != nil {
		t.Fatalf("checkExistingLock = %v, want nil", err)
	}
}

func writeLockForTest(path string) error {
	return os.WriteFile(path, []byte(fmt.Sprintf("%d\n%d\n", os.Getpid(), time.Now().Unix())), 0o644)
}
