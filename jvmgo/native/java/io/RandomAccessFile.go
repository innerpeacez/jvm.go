package io

import (
	. "github.com/zxh0/jvm.go/jvmgo/any"
	"github.com/zxh0/jvm.go/jvmgo/jvm/rtda"
	rtc "github.com/zxh0/jvm.go/jvmgo/jvm/rtda/class"
	"github.com/zxh0/jvm.go/jvmgo/util"
	"os"
	"syscall"
)

func init() {
	_raf(raf_initIDs, "initIDs", "()V")
	_raf(raf_open, "open", "(Ljava/lang/String;I)V")
	_raf(raf_close0, "close0", "()V")
	_raf(raf_write0, "write0", "(I)V")
	_raf(raf_writeBytes, "writeBytes", "([BII)V")
	_raf(raf_readBytes, "readBytes", "([BII)I")
	_raf(raf_read0, "read0", "()I")
	_raf(raf_seek0, "seek0", "(J)V")
	//TODO
	//_raf(raf_getFilePointer, "getFilePointer", "()J")
}

func _raf(method Any, name, desc string) {
	rtc.RegisterNativeMethod("java/io/RandomAccessFile", name, desc, method)
}

// private static native void initIDs();
// ()V
func raf_initIDs(frame *rtda.Frame) {
	//TODO
}

// private native void open(String name, int mode) throws FileNotFoundException;
// (Ljava/lang/String;)V
func raf_open(frame *rtda.Frame) {
	vars := frame.LocalVars()
	this := vars.GetThis()
	name := vars.GetRef(1)
	mode := vars.GetInt(2) //flag
	flag := 0

	if mode&1 > 0 {
		flag |= syscall.O_RDONLY
	}

	//write
	if mode&2 > 0 {
		flag |= syscall.O_RDWR | syscall.O_CREAT
	}

	if mode&4 > 0 {
		flag |= syscall.O_SYNC | syscall.O_CREAT
	}

	if mode&8 > 0 {
		flag |= syscall.O_DSYNC | syscall.O_CREAT
	}

	goName := rtda.GoString(name)
	goFile, err := os.OpenFile(goName, flag, 0660)
	if err != nil {
		frame.Thread().ThrowFileNotFoundException(goName)
		return
	}

	this.SetExtra(goFile)
}

// private native void close0() throws IOException;
// ()V
func raf_close0(frame *rtda.Frame) {
	vars := frame.LocalVars()
	this := vars.GetThis()

	goFile := this.Extra().(*os.File)
	err := goFile.Close()
	if err != nil {
		//TODO
		panic("IOException")
	}
}

// private native void writeBytes(byte b[], int off, int len) throws IOException;
// ([BIIZ)V
func raf_writeBytes(frame *rtda.Frame) {
	vars := frame.LocalVars()
	this := vars.GetThis()       // this
	byteArrObj := vars.GetRef(1) // b
	offset := vars.GetInt(2)     // off
	length := vars.GetInt(3)     // len

	goFile := this.Extra().(*os.File)

	jBytes := byteArrObj.Fields().([]int8)
	jBytes = jBytes[offset : offset+length]
	goBytes := util.CastInt8sToUint8s(jBytes)
	goFile.Write(goBytes)
}

// private native void write0(int b) throws IOException;
// (I)V
func raf_write0(frame *rtda.Frame) {
	vars := frame.LocalVars()
	this := vars.GetThis()
	intObj := vars.GetInt(1) // b

	goFile := this.Extra().(*os.File)
	//b := make([]byte, 4)
	//binary.BigEndian.PutUint32(b, uint32(intObj))
	_, err := goFile.Write([]byte{byte(intObj)})

	if err != nil {
		panic("IOException!" + err.Error())
	}
}

// private native int readBytes(byte b[], int off, int len) throws IOException;
// ([BII)I
func raf_readBytes(frame *rtda.Frame) {
	vars := frame.LocalVars()
	this := vars.GetThis()
	buf := vars.GetRef(1)
	off := vars.GetInt(2)
	_len := vars.GetInt(3)

	goFile := this.Extra().(*os.File)
	goBuf := buf.GoBytes()
	goBuf = goBuf[off : off+_len]

	n, err := goFile.Read(goBuf)
	if err == nil || n > 0 {
		frame.OperandStack().PushInt(int32(n))
	} else {
		//TODO
		panic("IOException!" + err.Error())
	}
}

// public native int read() throws IOException;
// ()I
func raf_read0(frame *rtda.Frame) {
	vars := frame.LocalVars()
	this := vars.GetThis()

	goFile := this.Extra().(*os.File)

	//b := make([]byte, 4)
	b := make([]byte, 1)
	_, err := goFile.Read(b)

	if err != nil {
		panic("IOException!" + err.Error())
	}
	//n := binary.BigEndian.Uint32(b)
	//frame.OperandStack().PushInt(int32(n))
	frame.OperandStack().PushInt(int32(b[0]))
}

// private native void seek0(long pos) throws IOException;
// (J)V
func raf_seek0(frame *rtda.Frame) {
	vars := frame.LocalVars()
	this := vars.GetThis()
	pos := vars.GetLong(1)

	goFile := this.Extra().(*os.File)

	if pos < 0 {
		//TODO
		panic("IOException! Negative seek offset")
	}

	if _, err := goFile.Seek(pos, os.SEEK_SET); err != nil {
		//TODO
		panic("IOException!" + err.Error())
	}
}
