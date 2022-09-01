package gutils

func InUint32List(val uint32, list []uint32) bool {
	var exist bool
	for _, v := range list {
		if val == v {
			exist = true
			break
		}
	}
	return exist
}

func InUint64List(val uint64, list []uint64) bool {
	var exist bool
	for _, v := range list {
		if val == v {
			exist = true
			break
		}
	}
	return exist
}

func InUint16List(val uint16, list []uint16) bool {
	var exist bool
	for _, v := range list {
		if val == v {
			exist = true
			break
		}
	}
	return exist
}

func InUintList(val uint, list []uint) bool {
	var exist bool
	for _, v := range list {
		if val == v {
			exist = true
			break
		}
	}
	return exist
}

func InInt32List(val int32, list []int32) bool {
	var exist bool
	for _, v := range list {
		if val == v {
			exist = true
			break
		}
	}
	return exist
}

func InInt64List(val int64, list []int64) bool {
	var exist bool
	for _, v := range list {
		if val == v {
			exist = true
			break
		}
	}
	return exist
}

func InInt16List(val int16, list []int16) bool {
	var exist bool
	for _, v := range list {
		if val == v {
			exist = true
			break
		}
	}
	return exist
}

func InIntList(val int, list []int) bool {
	var exist bool
	for _, v := range list {
		if val == v {
			exist = true
			break
		}
	}
	return exist
}

func InFloat32List(val float32, list []float32) bool {
	var exist bool
	for _, v := range list {
		if val == v {
			exist = true
			break
		}
	}
	return exist
}

func InFloat64List(val float64, list []float64) bool {
	var exist bool
	for _, v := range list {
		if val == v {
			exist = true
			break
		}
	}
	return exist
}

func InStringList(val string, list []string) bool {
	var exist bool
	for _, v := range list {
		if val == v {
			exist = true
			break
		}
	}
	return exist
}
