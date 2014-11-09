/*
 * filename   : logging_test.go
 * created at : 2014-11-09 17:23:16
 * author     : Jianing Yang <jianingy.yang@gmail.com>
 */

package logging

import (
    "bytes"
    "testing"
	"text/template"
    "github.com/stretchr/testify/assert"
)

func TestDefaultLogging(t *testing.T) {
    LOG := NewPlainLogger("testing")
    out := new(bytes.Buffer)
    LOG.SetWriter(out)

    LOG.NOTICE.Print("notice message")
    assert.Contains(t, out.String(), "[NOTICE]")
    assert.Contains(t, out.String(), "notice message")
}

func TestLoggingLevel(t *testing.T) {
    LOG := NewPlainLogger("testing")
    out := new(bytes.Buffer)
    LOG.SetWriter(out)

    LOG.DEBUG.Print("debug message")
    LOG.NOTICE.Print("notice message")
    assert.Contains(t, out.String(), "[NOTICE]")
    assert.NotContains(t, out.String(), "[DEBUG]")
}

func TestDecreaseLoggingLevel(t *testing.T) {
    LOG := NewPlainLogger("testing")
    out := new(bytes.Buffer)
    LOG.SetWriter(out)

    LOG.DEBUG.Print("message")
    LOG.NOTICE.Print("message")
    LOG.ERROR.Print("error message")
    assert.Contains(t, out.String(), "[ERROR]")
    assert.Contains(t, out.String(), "[NOTICE]")
    assert.NotContains(t, out.String(), "[DEBUG]")

    out.Reset()

    LOG.SetOutputLevel(LevelError)
    LOG.DEBUG.Print("debug message")
    LOG.NOTICE.Print("notice message")
    LOG.ERROR.Print("error message")
    assert.NotContains(t, out.String(), "[NOTICE]")
    assert.NotContains(t, out.String(), "[DEBUG]")
    assert.Contains(t, out.String(), "[ERROR]")
}

func TestIncreaseLoggingLevel(t *testing.T) {
    LOG := NewPlainLogger("testing")
    out := new(bytes.Buffer)
    LOG.SetWriter(out)

    LOG.DEBUG.Print("debug message")
    LOG.NOTICE.Print("notice message")
    assert.Contains(t, out.String(), "[NOTICE]")
    assert.NotContains(t, out.String(), "[DEBUG]")

    out.Reset()

    LOG.SetOutputLevel(LevelDebug)
    LOG.DEBUG.Print("debug message")
    LOG.NOTICE.Print("notice message")
    assert.Contains(t, out.String(), "[NOTICE]")
    assert.Contains(t, out.String(), "[DEBUG]")
}

func TestLoggingFormat(t *testing.T) {
    LOG := NewPlainLogger("testing")
    out := new(bytes.Buffer)
    LOG.SetWriter(out)
    LOG.SetFormat(">>> {{.Level}} {{.Message}}")

    LOG.NOTICE.Print("notice message")
    assert.Equal(t, out.String(), ">>> NOTICE notice message\n")
}

func TestLoggingPrintf(t *testing.T) {
    LOG := NewPlainLogger("testing")
    out := new(bytes.Buffer)
    LOG.SetWriter(out)
    LOG.SetFormat(">>> {{.Level}} {{.Message}}")

    LOG.NOTICE.Printf("notice message: %d", 12306)
    assert.Equal(t, out.String(), ">>> NOTICE notice message: 12306\n")
}

func TestSetWriter(t *testing.T) {
    LOG := NewPlainLogger("testing")
    out := new(bytes.Buffer)
    dumb := new(bytes.Buffer)
    LOG.SetWriter(dumb)
    LOG.NOTICE.SetWriter(out)

    LOG.NOTICE.Print("notice message")
    LOG.ERROR.Print("error message")

    assert.Contains(t, out.String(), "[NOTICE]")
    assert.NotContains(t, out.String(), "[ERROR]")
}

func TestSetTemplate(t *testing.T) {
    LOG := NewPlainLogger("testing")
    out := new(bytes.Buffer)
    LOG.SetWriter(out)
    LOG.NOTICE.SetTemplate(template.Must(template.New("testing").Parse("<{{.Level}}> {{.Message}}")))

    LOG.NOTICE.Print("notice message")
    LOG.ERROR.Print("error message")

    assert.Contains(t, out.String(), "<NOTICE>")
    assert.Contains(t, out.String(), "[ERROR]")
}
