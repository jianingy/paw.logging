paw.logging
===========

A simple logging utility for golang


Basic Usage
-----------

        import "github.com/jianingy/paw.logging"

        func main() {
            LOG := NewPlainLogger("testing")

            LOG.DEBUG.Print("debug message")
            LOG.NOTICE.Print("notice message")
        }
