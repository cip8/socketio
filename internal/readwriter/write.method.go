package readwriter

import (
	"io"

	errs "github.com/njones/socketio/internal/errors"
)

func (wtr *Writer) Bytes(p []byte) wtrErr {
	if wtr.err != nil {
		return wtr
	}

	_, wtr.err = wtr.w.Write(p)
	return onWtrErr{wtr}
}

func (wtr *Writer) Byte(p byte) wtrErr {
	if wtr.err != nil {
		return wtr
	}

	wtr.err = wtr.w.WriteByte(p)
	return onWtrErr{wtr}
}

func (wtr *Writer) String(str string) wtrErr {
	if wtr.err != nil {
		return wtr
	}

	return wtr.Bytes([]byte(str))
}

func (wtr *Writer) To(w io.WriterTo) wtrErr {
	if wtr.err != nil {
		return wtr
	}

	_, wtr.err = w.WriteTo(wtr.w)
	return onWtrErr{wtr}
}

func (wtr *Writer) Copy(src io.Reader) wtrErr {
	if wtr.err != nil {
		return wtr
	}

	_, wtr.err = io.Copy(wtr.w, src)
	return onWtrErr{wtr}
}

type onWtrErr struct{ *Writer }

func (e onWtrErr) OnErr(err errs.String) {
	if e.err != nil {
		e.err = err
	}
}
func (e onWtrErr) OnErrF(err errs.String, v ...interface{}) {
	if e.err != nil {
		e.err = err.F(v...)
	}
}
