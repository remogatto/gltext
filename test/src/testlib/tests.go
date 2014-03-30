package testlib

import (
	"bytes"
	"fmt"

	"github.com/remogatto/gltext"
	"github.com/remogatto/imagetest"
	"github.com/remogatto/mandala"
	"github.com/remogatto/mandala/test/src/testlib"

	gl "github.com/remogatto/opengles2"
)

const (
	distanceThreshold = 0.03
)

func distanceError(distance float64, filename string) string {
	return fmt.Sprintf("Image differs by distance %f, result saved in %s", distance, filename)
}

func (t *TestSuite) TestPrint() {
	// Load the font
	responseCh := make(chan mandala.LoadResourceResponse)
	mandala.ReadResource("raw/FreeSans.ttf", responseCh)
	response := <-responseCh
	fontBuffer := response.Buffer
	err := response.Error
	if err != nil {
		panic(err)
	}

	filename := "expected_hello_world.png"

	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)

		// Render an "Hello World" string
		sans, err := gltext.LoadTruetype(bytes.NewBuffer(fontBuffer), world, 40, 32, 127, gltext.LeftToRight)
		if err != nil {
			panic(err)
		}

		text, err := sans.Printf("%s", "Hello World!")
		if err != nil {
			panic(err)
		}

		text.AttachToWorld(world)
		text.MoveTo(float32(world.width/2), -10.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		text.Draw()

		t.renderState.window.SwapBuffers()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
	}

	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}

	t.True(distance < distanceThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}
