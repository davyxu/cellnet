package cellnet

/*
var result []string

func makeDispatcher() *DispatcherHandler {

	result = make([]string, 0)

	dispatcher := NewDispatcherHandler()

	dispatcher.AddCallback(1, func(interface{}) error {
		result = append(result, "A")
	}).Tag = "A"

	dispatcher.AddCallback(1, func(interface{}) {
		result = append(result, "B")
	}).Tag = "B"

	dispatcher.AddCallback(1, func(interface{}) {
		result = append(result, "C")
	}).Tag = "C"

	return dispatcher
}

type idmaker struct {
}

func (self idmaker) ContextID() uint32 {
	return 1
}

func TestDispatcherVisitRemoveExceptMid(t *testing.T) {

	dispatcher := makeDispatcher()

	dispatcher.VisitCallback(func(id uint32, ctx *CallbackContext) VisitOperation {

		if ctx.Tag != "B" {
			return VisitOperation_Remove
		}

		return VisitOperation_Continue

	})

	dispatcher.CallData(&idmaker{})

	t.Log(result)

	if len(result) != 1 || result[0] != "B" {

		t.Log("remove except b failed")
		t.FailNow()
	}

}

func TestDispatcherVisitRemoveExceptHead(t *testing.T) {

	dispatcher := makeDispatcher()

	dispatcher.VisitCallback(func(id uint32, ctx *CallbackContext) VisitOperation {

		if ctx.Tag != "A" {
			return VisitOperation_Remove
		}

		return VisitOperation_Continue

	})

	dispatcher.CallData(&idmaker{})

	t.Log(result)

	if len(result) != 1 || result[0] != "A" {

		t.Log("remove except a failed")
		t.FailNow()
	}

}

func TestDispatcherVisitRemoveExceptTail(t *testing.T) {

	dispatcher := makeDispatcher()

	dispatcher.VisitCallback(func(id uint32, ctx *CallbackContext) VisitOperation {

		if ctx.Tag != "C" {
			return VisitOperation_Remove
		}

		return VisitOperation_Continue

	})

	dispatcher.CallData(&idmaker{})

	t.Log(result)

	if len(result) != 1 || result[0] != "C" {

		t.Log("remove except c failed")
		t.FailNow()
	}

}

*/
