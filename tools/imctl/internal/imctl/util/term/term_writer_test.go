// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package term

import (
	"bytes"
	"strings"
	"testing"
)

const test = "Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam Iam"

func TestWordWrapWriter(t *testing.T) {
	testcases := map[string]struct {
		input    string
		maxWidth uint
	}{
		"max 10":   {input: test, maxWidth: 10},
		"max 80":   {input: test, maxWidth: 80},
		"max 120":  {input: test, maxWidth: 120},
		"max 5000": {input: test, maxWidth: 5000},
	}
	for k, tc := range testcases {
		b := bytes.NewBufferString("")
		w := NewWordWrapWriter(b, tc.maxWidth)
		_, err := w.Write([]byte(tc.input))
		if err != nil {
			t.Errorf("%s: Unexpected error: %v", k, err)
		}
		result := b.String()
		if !strings.Contains(result, "Iam") {
			t.Errorf("%s: Expected to contain \"Iam\"", k)
		}
		if len(result) < len(tc.input) {
			t.Errorf(
				"%s: Unexpectedly short string, got %d wanted at least %d chars: %q",
				k,
				len(result),
				len(tc.input),
				result,
			)
		}
		for _, line := range strings.Split(result, "\n") {
			if len(line) > int(tc.maxWidth) {
				t.Errorf("%s: Every line must be at most %d chars long, got %d: %q", k, tc.maxWidth, len(line), line)
			}
		}
		for _, word := range strings.Split(result, " ") {
			if !strings.Contains(word, "Iam") {
				t.Errorf("%s: Unexpected broken word: %q", k, word)
			}
		}
	}
}

func TestMaxWidthWriter(t *testing.T) {
	testcases := map[string]struct {
		input    string
		maxWidth uint
	}{
		"max 10":   {input: test, maxWidth: 10},
		"max 80":   {input: test, maxWidth: 80},
		"max 120":  {input: test, maxWidth: 120},
		"max 5000": {input: test, maxWidth: 5000},
	}
	for k, tc := range testcases {
		b := bytes.NewBufferString("")
		w := NewMaxWidthWriter(b, tc.maxWidth)
		_, err := w.Write([]byte(tc.input))
		if err != nil {
			t.Errorf("%s: Unexpected error: %v", k, err)
		}
		result := b.String()
		if !strings.Contains(result, "Iam") {
			t.Errorf("%s: Expected to contain \"Iam\"", k)
		}
		if len(result) < len(tc.input) {
			t.Errorf(
				"%s: Unexpectedly short string, got %d wanted at least %d chars: %q",
				k,
				len(result),
				len(tc.input),
				result,
			)
		}
		lines := strings.Split(result, "\n")
		for i, line := range lines {
			if len(line) > int(tc.maxWidth) {
				t.Errorf("%s: Every line must be at most %d chars long, got %d: %q", k, tc.maxWidth, len(line), line)
			}
			if i < len(lines)-1 && len(line) != int(tc.maxWidth) {
				t.Errorf(
					"%s: Lines except the last one are expected to be exactly %d chars long, got %d: %q",
					k,
					tc.maxWidth,
					len(line),
					line,
				)
			}
		}
	}
}
