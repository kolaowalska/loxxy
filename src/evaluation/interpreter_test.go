package evaluation_test

import (
	"github.com/kolaowalska/loxxy/src/evaluation"
	"github.com/kolaowalska/loxxy/src/representation"
	"testing"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		expr    representation.Expr
		want    any
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := evaluation.Evaluate(tt.expr)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Evaluate() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Evaluate() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
