package sorts

func mustMatchHead(err error, head Name, list List) {
	if len(list) >= 1 && list[0] == head {
		return
	}
	panic(err)
}
