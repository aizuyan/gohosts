package main

// 单个hosts信息，名称，内容，是否开启，语法校验是否通过
type hostsItem struct {
	HostsName string
	Content string
	Toggle bool
}

func (h *hostsItem) toggleHostsItemSwitch(flag bool) {
	h.Toggle = flag
}

func hostsNameToString(hosts []hostsItem) string {
	ret := ""
	for _, h := range hosts {
		if h.Toggle {
			ret += "√ "
		} else {
			ret += "× "
		}
		ret += h.HostsName + "\n"
	}

	return ret
}