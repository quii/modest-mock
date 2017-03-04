package testdata

type StoreMock struct {
	Calls struct {
		Save []struct {
			firstname string
		}
	}
}

func (s *StoreMock) Save(firstname string) {

}
