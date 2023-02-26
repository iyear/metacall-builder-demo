package builder

var Configurations = map[Language]func(scripts, ports bool) []string{
	Python: confPy,
	// TODO(iyear): Add more languages
}

var confGen = func(name string, scripts, ports bool) []string {
	s := make([]string, 0)
	s = append(s, "-DOPTION_BUILD_LOADERS_"+name+"=On") // necessary for all languages
	if scripts {
		s = append(s, "-DOPTION_BUILD_SCRIPTS_"+name+"=On")
	}
	if ports {
		s = append(s, "-DOPTION_BUILD_PORTS_"+name+"=On")
	}
	return s
}

func confPy(scripts, ports bool) []string {
	return confGen("PY", scripts, ports)
}
