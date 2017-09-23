package main_test

import (
	. "github.com/dustinrohde/ecoscript"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("World", func() {
	var (
		world         *World
		orgWalkable   *Organism
		orgUnwalkable *Organism
		org3          *Organism
	)

	BeforeEach(func() {
		world = NewWorld(3, 3, []string{"ground"})
		orgWalkable = NewOrganism("organism1", "1", &Attributes{
			Walkable: true,
		})
		orgUnwalkable = NewOrganism("organism2", "2", &Attributes{
			Walkable: false,
		})
		org3 = NewOrganism("organism3", "3", new(Attributes))
	})

	Describe("World#Add()", func() {
		Context("Organism is NOT walkable", func() {
			It("should inhabit AND occupy the Cell", func() {
				vec := Vec(0, 0, 0)
				exec, ok := world.Add(orgUnwalkable, vec)
				Expect(ok).To(BeTrue())
				exec()

				cell := world.Cell(vec)
				Expect(len(cell.Organisms())).To(Equal(1))
				Expect(cell.Organisms()[0].ID()).To(Equal(orgUnwalkable.ID()))
				Expect(cell.Occupied()).To(BeTrue())
				Expect(cell.Occupier().ID()).To(Equal(orgUnwalkable.ID()))
			})
		})

		Context("Organism IS walkable", func() {
			It("should inhabit but NOT occupy the Cell", func() {
				vec := Vec(0, 0, 0)
				exec, ok := world.Add(orgWalkable, vec)
				Expect(ok).To(BeTrue())
				exec()

				cell := world.Cell(vec)
				Expect(len(cell.Organisms())).To(Equal(1))
				Expect(cell.Organisms()[0].ID()).To(Equal(orgWalkable.ID()))
				Expect(cell.Occupied()).To(BeFalse())
				Expect(cell.Occupier()).To(BeNil())
			})
		})
	})

	Describe("World#Remove()", func() {
		Context("Organism IS walkable", func() {
			It("should be removed from the Cell", func() {
				vec := Vec(0, 0, 0)
				exec, ok := world.Add(orgWalkable, vec)
				Expect(ok).To(BeTrue())
				exec()

				cell := world.Cell(vec)
				Expect(len(cell.Organisms())).To(Equal(1))
				Expect(cell.Organisms()[0].ID()).To(Equal(orgWalkable.ID()))

				exec, ok = world.Remove(orgWalkable, vec)
				Expect(ok).To(BeTrue())
				exec()

				cell = world.Cell(vec)
				Expect(len(cell.Organisms())).To(Equal(0))
			})
		})

		Context("Organism is NOT walkable", func() {
			It("should be removed from the Cell and its occupancy", func() {
				vec := Vec(0, 0, 0)
				exec, ok := world.Add(orgUnwalkable, vec)
				Expect(ok).To(BeTrue())
				exec()

				cell := world.Cell(vec)
				Expect(len(cell.Organisms())).To(Equal(1))
				Expect(cell.Organisms()[0].ID()).To(Equal(orgUnwalkable.ID()))
				Expect(cell.Occupied()).To(BeTrue())
				Expect(cell.Occupier().ID()).To(Equal(orgUnwalkable.ID()))

				exec, ok = world.Remove(orgUnwalkable, vec)
				Expect(ok).To(BeTrue())
				exec()

				cell = world.Cell(vec)
				Expect(len(cell.Organisms())).To(Equal(0))
				Expect(cell.Occupied()).To(BeFalse())
				Expect(cell.Occupier()).To(BeNil())
			})
		})
	})
})
