package integration

import (
	"testing"

	"github.com/andrewneudegg/calc/pkg/display"
)

func TestInlineCommentsAreIgnored_SimpleAssign(t *testing.T) {
	r := display.NewREPL()

	v := r.EvaluateLine("x=1 // set x")
	if v.IsError() {
		t.Fatalf("unexpected error assigning with comment: %v", v)
	}

	v2 := r.EvaluateLine("x")
	if v2.IsError() || int(v2.Number) != 1 {
		t.Fatalf("expected x==1, got %+v", v2)
	}
}

func TestInlineCommentsAreIgnored_RateAssign(t *testing.T) {
	r := display.NewREPL()
	v := r.EvaluateLine("t= 99 gbp per day // Johnny's wage")
	if v.IsError() {
		t.Fatalf("unexpected error assigning rate with comment: %v", v)
	}
}
