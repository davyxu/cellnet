package cellnet

type Codec interface {
	Encode(interface{}) ([]byte, error)

	Decode([]byte, interface{}) error

	Name() string
}

var codecByName = map[string]Codec{}

func RegisterCodec(name string, c Codec) {

	if _, ok := codecByName[name]; ok {
		panic("duplicate codec: " + name)
	}

	//log.Infof("registere codec '%s'!", name)

	codecByName[name] = c
}

func FetchCodec(name string) Codec {
	if v, ok := codecByName[name]; ok {
		return v
	}

	return nil
}
