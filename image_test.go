// Copyright ©2024 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gifscaleissue

import (
	"bytes"
	"flag"
	"image"
	"image/gif"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/image/draw"
)

var update = flag.Bool("update", false, "regenerate golden image")

func TestGIFScale(t *testing.T) {
	f, err := os.Open("gopher-dance-long.gif")
	if err != nil {
		t.Fatalf("failed to read image data: %v", err)
	}
	g, err := gif.DecodeAll(f)
	if err != nil {
		t.Fatalf("failed to decode image data: %v", err)
	}
	rect := image.Rectangle{Max: image.Point{X: 72, Y: 72}}
	for i, frame := range g.Image {
		g.Image[i] = image.NewPaletted(rect, g.Image[0].Palette)
		draw.BiLinear.Scale(g.Image[i], rect, frame, frame.Bounds(), draw.Src, nil)
	}
	g.Config.Width = rect.Dx()
	g.Config.Height = rect.Dy()

	var buf bytes.Buffer
	err = gif.EncodeAll(&buf, g)
	if err != nil {
		t.Fatalf("unexpected error encoding image: %v", err)
	}

	path := "gopher-dance-long-scaled.gif"
	if *update {
		err = os.WriteFile(path, buf.Bytes(), 0o644)
		if err != nil {
			t.Fatalf("unexpected error writing golden file: %v", err)
		}
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected error reading golden file: %v", err)
	}
	if !bytes.Equal(buf.Bytes(), want) {
		t.Errorf("image mismatch:\n- want:\n+ got:\n%s", cmp.Diff(want, buf.Bytes()))
	}
}
