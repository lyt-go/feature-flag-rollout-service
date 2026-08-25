package store

// clone 返回指针指向值的副本。
//
// MemoryStore 在 map 中保存的是指针，若读方法直接把存储内部的指针返回给
// 调用方，调用方在未发起更新请求的情况下修改返回对象，就会写穿到存储里。
// 所有读方法返回 clone(p)，写方法写入 clone(p)，即可保证存储内对象与
// 返回给调用方的对象互为独立副本，调用方的改动不会影响存储。
func clone[T any](p *T) *T {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}
