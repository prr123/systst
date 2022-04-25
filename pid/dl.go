//
func getUnitFileName() (unit string, err error) {
	libname := C.CString("libsystemd.so")
	defer C.free(unsafe.Pointer(libname))
	handle := C.dlopen(libname, C.RTLD_LAZY)
	if handle == nil {
		err = fmt.Errorf("error opening libsystemd.so")
		return
	}
	defer func() {
		if r := C.dlclose(handle); r != 0 {
			err = fmt.Errorf("error closing libsystemd.so")
		}
	}()

	sym := C.CString("sd_pid_get_unit")
	defer C.free(unsafe.Pointer(sym))
	sd_pid_get_unit := C.dlsym(handle, sym)
	if sd_pid_get_unit == nil {
		err = fmt.Errorf("error resolving sd_pid_get_unit function")
		return
	}

	var s string
	u := C.CString(s)
	defer C.free(unsafe.Pointer(u))

	ret := C.my_sd_pid_get_unit(sd_pid_get_unit, 0, &u)
	if ret < 0 {
		err = fmt.Errorf("error calling sd_pid_get_unit: %v", syscall.Errno(-ret))
		return
	}

	unit = C.GoString(u)
	return
}

//https://github.com/tiborvass/dl/blob/master/dl.go

// Open opens the shared library identified by the given name
// with the given flags. See man dlopen for the available flags
// and its meaning. Note that the only difference with dlopen is that
// if nor RTLD_LAZY nor RTLD_NOW are specified, Open defaults to
// RTLD_NOW rather than returning an error. If the name argument
// passed to name does not have extension, the default for the
// platform will be appended to it (e.g. .so, .dylib, etc...).
func Open(name string, flag int) (*DL, error) {
	if flag&RTLD_LAZY == 0 && flag&RTLD_NOW == 0 {
		flag |= RTLD_NOW
	}
	if name != "" && filepath.Ext(name) == "" {
		name = name + LibExt
	}
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	mu.Lock()
	handle := C.dlopen(s, C.int(flag))
	var err error
	if handle == nil {
		err = dlerror()
	}
	mu.Unlock()
	if err != nil {
		if runtime.GOOS == "linux" && name == "libc.so" {
			// In most distros libc.so is now a text file
			// and in order to dlopen() it the name libc.so.6
			// must be used.
			return Open(name+".6", flag)
		}
		return nil, err
	}
	return &DL{
		handle: handle,
	}, nil
}

