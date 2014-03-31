package testlib

import (
	"fmt"
	"github.com/remogatto/imagetest"
	"github.com/remogatto/mandala/test/src/testlib"

	gl "github.com/remogatto/opengles2"
)

const (
	distanceThreshold = 0.05
)

func distanceError(distance float64, filename string) string {
	return fmt.Sprintf("Image differs by distance %f, result saved in %s", distance, filename)
}

func (t *TestSuite) TestPrint() {
	filename := "expected_hello_world.png"

	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)

		// Render an "Hello World" string
		text, err := world.font.Printf("%s", "Hello World!")
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

func (t *TestSuite) TestRotate() {
	filename := "expected_hello_world_rotated.png"

	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)

		// Render an "Hello World" string
		text, err := world.font.Printf("%s", "Hello World!")
		if err != nil {
			panic(err)
		}

		text.AttachToWorld(world)
		text.MoveTo(float32(world.width/2)+5, -5.0)
		text.Rotate(45)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		text.Draw()

		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
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
