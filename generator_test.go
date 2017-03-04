package modestmock

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestGenerateMockCode(t *testing.T) {

	scenarios := []struct {
		interfaceName    string
		interfacePath    string
		expectedMockPath string
	}{
		{"Store", "simple.go", "simple_mock.go"},
		{"Returns", "returns.go", "returns_mock.go"},
	}

	for _, s := range scenarios {
		interfaceToMock := openTestFile(t, s.interfacePath)
		expectedMock := openTestFile(t, s.expectedMockPath)

		mock, err := New(strings.NewReader(interfaceToMock), s.interfaceName)

		if err != nil {
			t.Fatal("problem creating mock called", s.interfaceName, "from", s.interfacePath, err)
		}

		generatedCode, err := GenerateMockCode(mock)

		if err != nil {
			t.Log("Error generating mock", err)
		}

		assert.NoError(t, err)
		assert.Equal(t, expectedMock, generatedCode)
	}
}

func openTestFile(t *testing.T, path string) string {
	f, err := os.Open(fmt.Sprintf("testdata/%s", path))

	if err != nil {
		t.Fatal("problem opening", path, err)
	}

	b, err := ioutil.ReadAll(f)

	if err != nil {
		t.Fatal(err)
	}

	formatted, err := format.Source(b)

	if err != nil {
		t.Fatal("error formatting generated code, likely compile fail", err)
	}

	return string(formatted)
}
