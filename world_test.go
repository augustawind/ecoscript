package main_test

import (
	. "github.com/dustinrohde/ecoscript"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("World", func() {
	var (
		world *World
		org1  *Organism
		org2  *Organism
		org3  *Organism
	)

	BeforeEach(func() {
		world = NewWorld(3, 3, []string{"ground"})
		org1 = NewOrganism("organism1", "1", new(Attributes))
		org2 = NewOrganism("organism2", "2", new(Attributes))
		org3 = NewOrganism("organism3", "3", new(Attributes))
	})

	Describe("adding an Organism to a World", func() {
		Context("if it's NOT walkable", func() {
			It("should inhabit AND occupy a Cell", func() {
				org1.Attrs.Walkable = false

				vec := Vec(0, 0, 0)
				exec, ok := world.Add(org1, vec)
				Expect(ok).To(BeTrue())
				exec()

				cell := world.Cell(vec)
				Expect(len(cell.Organisms())).To(Equal(1))
				Expect(cell.Organisms()[0].ID()).To(Equal(org1.ID()))
				Expect(cell.Occupied()).To(BeTrue())
				Expect(cell.Occupier().ID()).To(Equal(org1.ID()))
			})
		})
		Context("if it IS walkable", func() {
			It("should inhabit but NOT occupy a Cell", func() {
				org1.Attrs.Walkable = true

				vec := Vec(0, 0, 0)
				exec, ok := world.Add(org1, vec)
				Expect(ok).To(BeTrue())
				exec()

				cell := world.Cell(vec)
				Expect(len(cell.Organisms())).To(Equal(1))
				Expect(cell.Organisms()[0].ID()).To(Equal(org1.ID()))
				Expect(cell.Occupied()).To(BeFalse())
				Expect(cell.Occupier()).To(BeNil())
			})
		})
	})
})
