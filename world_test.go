package main_test

import (
	. "github.com/dustinrohde/ecoscript"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("World", func() {
	var (
		world         *World
		entWalkable   *Entity
		entUnwalkable *Entity
		ent3          *Entity
	)

	BeforeEach(func() {
		world = NewWorld(3, 3, []string{"ground"})
		entWalkable = NewEntity("entity1", "1", &Attributes{
			Walkable: true,
		})
		entUnwalkable = NewEntity("entity2", "2", &Attributes{
			Walkable: false,
		})
		ent3 = NewEntity("entity3", "3", new(Attributes))
	})

	Describe("World#Add()", func() {
		Context("Entity is NOT walkable", func() {
			It("should inhabit AND occupy the Cell", func() {
				vec := Vec(0, 0, 0)
				exec, ok := world.Add(entUnwalkable, vec)
				Expect(ok).To(BeTrue())
				exec()

				cell := world.Cell(vec)
				Expect(len(cell.Entities())).To(Equal(1))
				Expect(cell.Entities()[0].ID()).To(Equal(entUnwalkable.ID()))
				Expect(cell.Occupied()).To(BeTrue())
				Expect(cell.Occupier().ID()).To(Equal(entUnwalkable.ID()))
			})
		})

		Context("Entity IS walkable", func() {
			It("should inhabit but NOT occupy the Cell", func() {
				vec := Vec(0, 0, 0)
				exec, ok := world.Add(entWalkable, vec)
				Expect(ok).To(BeTrue())
				exec()

				cell := world.Cell(vec)
				Expect(len(cell.Entities())).To(Equal(1))
				Expect(cell.Entities()[0].ID()).To(Equal(entWalkable.ID()))
				Expect(cell.Occupied()).To(BeFalse())
				Expect(cell.Occupier()).To(BeNil())
			})
		})
	})

	Describe("World#Remove()", func() {
		Context("Entity IS walkable", func() {
			It("should be removed from the Cell", func() {
				vec := Vec(0, 0, 0)
				exec, ok := world.Add(entWalkable, vec)
				Expect(ok).To(BeTrue())
				exec()

				cell := world.Cell(vec)
				Expect(len(cell.Entities())).To(Equal(1))
				Expect(cell.Entities()[0].ID()).To(Equal(entWalkable.ID()))

				exec, ok = world.Remove(entWalkable, vec)
				Expect(ok).To(BeTrue())
				exec()

				cell = world.Cell(vec)
				Expect(len(cell.Entities())).To(Equal(0))
			})
		})

		Context("Entity is NOT walkable", func() {
			It("should be removed from the Cell and its occupancy", func() {
				vec := Vec(0, 0, 0)
				exec, ok := world.Add(entUnwalkable, vec)
				Expect(ok).To(BeTrue())
				exec()

				cell := world.Cell(vec)
				Expect(len(cell.Entities())).To(Equal(1))
				Expect(cell.Entities()[0].ID()).To(Equal(entUnwalkable.ID()))
				Expect(cell.Occupied()).To(BeTrue())
				Expect(cell.Occupier().ID()).To(Equal(entUnwalkable.ID()))

				exec, ok = world.Remove(entUnwalkable, vec)
				Expect(ok).To(BeTrue())
				exec()

				cell = world.Cell(vec)
				Expect(len(cell.Entities())).To(Equal(0))
				Expect(cell.Occupied()).To(BeFalse())
				Expect(cell.Occupier()).To(BeNil())
			})
		})
	})

	Describe("World#Move()", func() {
		Context("Entity IS walkable", func() {
			It("should be moved from one Cell to another", func() {
				src := Vec(0, 0, 0)
				exec, ok := world.Add(entWalkable, src)
				Expect(ok).To(BeTrue())
				exec()

				cell := world.Cell(src)
				Expect(len(cell.Entities())).To(Equal(1))
				Expect(cell.Entities()[0].ID()).To(Equal(entWalkable.ID()))

				dst := Vec(1, 2, 0)
				exec, ok = world.Move(entWalkable, src, dst)
				Expect(ok).To(BeTrue())
				exec()

				cell = world.Cell(src)
				Expect(len(cell.Entities())).To(Equal(0))
				cell = world.Cell(dst)
				Expect(len(cell.Entities())).To(Equal(1))
				Expect(cell.Entities()[0].ID()).To(Equal(entWalkable.ID()))
			})
		})

		Context("Entity is NOT walkable", func() {
			It("should be moved from one Cell to another and occupancy changed", func() {
				src := Vec(0, 0, 0)
				exec, ok := world.Add(entUnwalkable, src)
				Expect(ok).To(BeTrue())
				exec()

				cell := world.Cell(src)
				Expect(len(cell.Entities())).To(Equal(1))
				Expect(cell.Entities()[0].ID()).To(Equal(entUnwalkable.ID()))
				Expect(cell.Occupied()).To(BeTrue())
				Expect(cell.Occupier().ID()).To(Equal(entUnwalkable.ID()))

				dst := Vec(1, 2, 0)
				exec, ok = world.Move(entUnwalkable, src, dst)
				Expect(ok).To(BeTrue())
				exec()

				cell = world.Cell(src)
				Expect(len(cell.Entities())).To(Equal(0))
				Expect(cell.Occupied()).To(BeFalse())
				Expect(cell.Occupier()).To(BeNil())

				cell = world.Cell(dst)
				Expect(len(cell.Entities())).To(Equal(1))
				Expect(cell.Entities()[0].ID()).To(Equal(entUnwalkable.ID()))
				Expect(cell.Occupied()).To(BeTrue())
				Expect(cell.Occupier().ID()).To(Equal(entUnwalkable.ID()))
			})
		})
	})
})
