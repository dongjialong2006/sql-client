package sql

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/fortytw2/leaktest"
)

func TestSecurity(t *testing.T) {
	defer leaktest.Check(t)()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Security Suite")
	os.Exit(2)
}

var _ = Describe("sql", func() {
	It("sql/Insert", func() {
		err := Create("./cache/", "agent.db", "running")
		Expect(err).Should(BeNil())
		tmp := "<policy>\"dvwdwewe\"</policy>"
		tmp = strings.Replace(tmp, "\"", "'", -1)
		fmt.Println(tmp)
		node := &Node{
			SN:        int64(3453464574),
			Pid:       "Some id",
			ErrOP:     "stop-on-error",
			Policy:    tmp,
			Timestamp: time.Now().UnixNano(),
		}
		err = Insert("running", node)
		Expect(err).Should(BeNil())
		fmt.Println("------------")
		nodes, err := Query("running", int64(3453464574), "Some id")
		Expect(err).Should(BeNil())
		fmt.Println(len(nodes))
		return
		/*
			err = Delete(pid, "running", false)
			Expect(err).Should(BeNil())
		*/
		time.Sleep(time.Second)
		node.Policy = "ergsdtyftdrgszes"
		node.Timestamp = time.Now().UnixNano()
		err = Update("running", node)
		Expect(err).Should(BeNil())

		nodes, err = Query("running", 0, "")
		Expect(err).Should(BeNil())
		fmt.Println(nodes[0])
	})
})
