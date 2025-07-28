package tests_test

import (
	. "github.com/guidewire-oss/fern-ginkgo-client/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Adder", Ordered, Label("unit"), func() {
	Describe("Add", func() {

		It("adds two numbers", Label("add 1"), func() {
			sum := Add(2, 3)
			Expect(sum).To(Equal(5))
		})

		It("adds two numbers, where one is negative", Label("add 2"), func() {
			sum := Add(2, -3)
			Expect(sum).To(Equal(-1))
		})
	})
})
