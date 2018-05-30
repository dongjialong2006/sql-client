package file

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/fortytw2/leaktest"
)

func TestFile(t *testing.T) {
	defer leaktest.Check(t)()
	RegisterFailHandler(Fail)
	RunSpecs(t, "File Suite")
	os.Exit(2)
}

var _ = Describe("File", func() {
	Specify("read file test", func() {
		tmp, err := ReadFile("")
		Expect(err).ShouldNot(BeNil())

		tmp, err = ReadFile("./tmp/config.ini")
		Expect(err).ShouldNot(BeNil())

		tmp, err = ReadFile("../../config/config.ini")
		if nil != err {
			tmp, err = ReadFile("./config/config.ini")
		}
		Expect(err).Should(BeNil())
		Expect(tmp).ShouldNot(BeEmpty())
	})

	Specify("generate new file test", func() {
		tmp, err := GenerateConfigFile("")
		Expect(err).ShouldNot(BeNil())

		tmp, err = GenerateConfigFile("./file/config.ini")
		Expect(err).ShouldNot(BeNil())

		tmp, err = GenerateConfigFile("../../config/config.ini")
		if nil != err {
			tmp, err = GenerateConfigFile("./config/config.ini")
		}
		Expect(err).Should(BeNil())
		Expect(tmp).ShouldNot(BeEmpty())

		info, err := os.Stat("../../config/.config.ini")
		if nil != err {
			info, err = os.Stat("./config/.config.ini")
			defer os.RemoveAll("./config/.config.ini")
		} else {
			defer os.RemoveAll("../../config/.config.ini")
		}
		Expect(err).Should(BeNil())
		Expect(info).ShouldNot(BeNil())
	})
})
