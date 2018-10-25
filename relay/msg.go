package relay

func (self *RelayACK) PassThrough() interface{} {
	switch self.Type {
	case RelayPassThroughType_Int64:
		return self.Int64
	case RelayPassThroughType_Int64Slice:
		return self.Int64Slice
	}

	return nil
}
