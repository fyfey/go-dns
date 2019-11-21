package main

type Zone struct {
	Name    string
	Records map[QueryType]map[string]*Resource
}

func NewZone(name string) *Zone {
	return &Zone{
		name,
		map[QueryType]map[string]*Resource{},
	}
}

func (z *Zone) AddRecord(t QueryType, name string, address string) *Zone {

	if _, ok := z.Records[t]; !ok {
		z.Records[t] = map[string]*Resource{}
	}
	z.Records[t][name] = &Resource{
		Name:        name,
		Class:       IN,
		Type:        t,
		TTL:         300,
		RDataLength: 4,
		RData:       address,
	}

	return z
}
