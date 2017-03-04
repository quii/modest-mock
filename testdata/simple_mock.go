package testdata

type StoreMock struct {
	Calls struct {
		Save []struct {
			firstname string
		}
	}
}

func (s *StoreMock) Save(firstname string) {
	call := struct{ firstname string }{firstname}
	s.Calls.Save = append(s.Calls.Save, call)

}
