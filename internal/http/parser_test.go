package http

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func TestParser_IsRequestComplete(t *testing.T) {
	for _, tC := range isRequestCompleteTestCases {
		t.Run(tC.desc, func(subT *testing.T) {
			got, err := IsRequestComplete(tC.input)
			if (err != nil) != tC.wantErr {
				subT.Errorf("IsRequestComplete() error = %v, wantErr %v", err, tC.wantErr)
				return
			}

			if err != nil && err != tC.expectedErr {
				subT.Errorf("IsRequestComplete() error type mismatch expecting %s and got %s", tC.expectedErr.Error(), err.Error())
				return
			}

			if !reflect.DeepEqual(got, tC.expected) {
				subT.Errorf("IsRequestComplete() got = %v, want %v", got, tC.expected)
			}
		})
	}
}

func TestParser_ParseContentLength(t *testing.T) {
	for _, tC := range parseContentLengthTestCases {
		t.Run(tC.desc, func(subT *testing.T) {
			got, err := parseContentLength(tC.input)
			if (err != nil) != tC.wantErr {
				subT.Errorf("parseContentLength() error = %v, wantErr %v", err, tC.wantErr)
				return
			}

			if err != nil && err != tC.expectedErr {
				subT.Errorf("parseContentLength() error type mismatch expecting %s and got %s", tC.expectedErr.Error(), err.Error())
				return
			}

			if !reflect.DeepEqual(got, tC.expected) {
				subT.Errorf("parseContentLength() got = %v, want %v", got, tC.expected)
			}
		})
	}

	// After the basic test cases let's run a bunch of tests on random numbers
	for i := 0; i < 100; i++ {
		r := rand.Int63n(1000000)
		b := []byte(fmt.Sprintf("%d", r))
		got, err := parseContentLength(b)
		if err != nil {
			t.Errorf("parseContentLength() error type mismatch expecting nil and got %s", err.Error())
			return
		}

		if !reflect.DeepEqual(got, r) {
			t.Errorf("parseContentLength() got = %v, want %v", got, r)
		}
	}
}
