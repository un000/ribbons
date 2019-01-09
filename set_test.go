package ribbons

import (
	"encoding/json"
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRibbons(t *testing.T) {
	Convey("Check the main logic", t, func() {
		Convey("Check empty", func() {
			s := New()
			So(s.Has(2), ShouldBeFalse)
		})

		Convey("Check combinations of the added values", func() {
			s := New()
			So(s.Len(), ShouldEqual, 0)

			s.Add(128)
			s.Add(128)
			So(s.Has(128), ShouldBeTrue)
			So(s.Len(), ShouldEqual, 1)

			s.Add(129)
			So(s.Has(128), ShouldBeTrue)
			So(s.Has(129), ShouldBeTrue)
			So(s.Len(), ShouldEqual, 2)

			s.Add(2)
			So(s.Has(2), ShouldBeTrue)
			So(s.Has(128), ShouldBeTrue)
			So(s.Has(129), ShouldBeTrue)
			So(s.Len(), ShouldEqual, 3)

			s.Add(257)
			So(s.Has(2), ShouldBeTrue)
			So(s.Has(128), ShouldBeTrue)
			So(s.Has(129), ShouldBeTrue)
			So(s.Has(257), ShouldBeTrue)
			So(s.Len(), ShouldEqual, 4)

			So(s.List(), ShouldResemble, []uint64{2, 128, 129, 257})
		})

		Convey("Check that some values exist", func() {
			s := New()
			s.Add(5)
			s.Add(8)
			s.Add(32)
			s.Add(127)
			s.Add(128)
			s.Add(129)
			s.Add(256)
			s.Add(257)

			So(s.Has(0), ShouldBeFalse)
			So(s.Has(math.MaxUint64), ShouldBeFalse)

			for i := uint64(0); i < 512; i++ {
				if i == 5 || i == 8 || i == 32 || i == 127 || i == 128 || i == 129 || i == 256 || i == 257 {
					So(s.Has(i), ShouldBeTrue)
				} else {
					So(s.Has(i), ShouldBeFalse)
				}
			}
		})

		Convey("Check unmarshal", func() {
			Convey("When simple unmarshalling", func() {
				s := New()
				err := json.Unmarshal([]byte(`[1, 2, 0, 1000, 8, 5, 5, 10, 11, 24, 9]`), &s)
				So(err, ShouldBeNil)

				validRes := []uint64{0, 1, 2, 5, 8, 9, 10, 11, 24, 1000}
				So(s.List(), ShouldResemble, validRes)
				for _, v := range validRes {
					So(s.Has(v), ShouldBeTrue)
				}

			})

			Convey("Check initialized", func() {
				s, err := NewFromJSON([]byte(`[]`))
				So(err, ShouldBeNil)
				So(s.Initialized(), ShouldBeTrue)

				s, err = NewFromJSON([]byte(`null`))
				So(err, ShouldBeNil)
				So(s.Initialized(), ShouldBeFalse)
			})
		})

		Convey("Check delete", func() {
			s := New()
			err := json.Unmarshal([]byte(`[1, 2, 0, 1000, 8, 5, 5, 10, 11, 24, 9]`), &s)
			So(err, ShouldBeNil)

			validRes := []uint64{0, 1, 2, 5, 8, 9, 10, 11, 24, 1000}
			So(s.List(), ShouldResemble, validRes)

			s.Delete(0)
			So(s.List(), ShouldResemble, []uint64{1, 2, 5, 8, 9, 10, 11, 24, 1000})
			So(s.Has(0), ShouldBeFalse)

			s.Delete(1000)
			So(s.List(), ShouldResemble, []uint64{1, 2, 5, 8, 9, 10, 11, 24})
			So(s.Has(1000), ShouldBeFalse)

			s.Delete(8)
			So(s.List(), ShouldResemble, []uint64{1, 2, 5, 9, 10, 11, 24})
			So(s.Has(8), ShouldBeFalse)

			for _, v := range s.List() {
				s.Delete(v)
				So(s.Has(v), ShouldBeFalse)
			}

			So(s.Len(), ShouldEqual, 0)
			So(s.List(), ShouldResemble, []uint64{})

			s.Delete(0)
			s.Delete(100)
			So(s.Len(), ShouldEqual, 0)
		})

		Convey("Check add -> delete -> add minumum", func() {
			s := New()
			err := json.Unmarshal([]byte(`[0, 1]`), &s)
			So(err, ShouldBeNil)

			s.Delete(0)
			So(s.List(), ShouldResemble, []uint64{1})
			s.Add(0)
			So(s.List(), ShouldResemble, []uint64{0, 1})

			s.Delete(0)
			s.Delete(1)
			So(s.List(), ShouldResemble, []uint64{})
			s.Add(1)
			So(s.List(), ShouldResemble, []uint64{1})
		})

		Convey("Check sum", func() {
			Convey("When adding 2 sets", func() {
				s1 := New()
				err := json.Unmarshal([]byte(`[13, 10, 11, 24, 0, 9]`), &s1)
				So(err, ShouldBeNil)

				s2 := New()
				err = json.Unmarshal([]byte(`[5, 6, 7, 0]`), &s2)
				So(err, ShouldBeNil)

				s1.Or(&s2)
				So(s1.List(), ShouldResemble, []uint64{0, 5, 6, 7, 9, 10, 11, 13, 24})
			})

			Convey("When adding to the empty set", func() {
				s1 := New()
				err := json.Unmarshal([]byte(`[5, 6, 7, 0]`), &s1)
				So(err, ShouldBeNil)

				s2 := New()
				s2.Or(&s1)
				So(s2.List(), ShouldResemble, []uint64{0, 5, 6, 7})
			})

			Convey("When adding empty set", func() {
				s1 := New()
				err := json.Unmarshal([]byte(`[5, 6, 7, 0]`), &s1)
				So(err, ShouldBeNil)

				s2 := New()
				s1.Or(&s2)
				So(s1.List(), ShouldResemble, []uint64{0, 5, 6, 7})
			})
		})

		Convey("Check multiplication", func() {
			Convey("When multiplying 2 sets", func() {
				s1 := New()
				err := json.Unmarshal([]byte(`[13, 10, 11, 24, 0, 9]`), &s1)
				So(err, ShouldBeNil)

				s2 := New()
				err = json.Unmarshal([]byte(`[5, 11, 7, 0]`), &s2)
				So(err, ShouldBeNil)

				s1.And(&s2)
				So(s1.List(), ShouldResemble, []uint64{0, 11})
			})

			Convey("When multiplying 2 sets with 1 empty set", func() {
				s1 := New()
				err := json.Unmarshal([]byte(`[13, 10, 11, 24, 0, 9]`), &s1)
				So(err, ShouldBeNil)

				s2 := New()
				err = json.Unmarshal([]byte(`[]`), &s2)
				So(err, ShouldBeNil)

				s1.And(&s2)
				So(s1.List(), ShouldResemble, []uint64{})
			})

			Convey("When multiplying 2 empty sets", func() {
				s1 := New()
				s2 := New()

				s1.And(&s2)
				So(s1.List(), ShouldResemble, []uint64{})
			})

			Convey("When inverse", func() {
				s1 := New()
				err := json.Unmarshal([]byte(`[13, 10, 11]`), &s1)
				So(err, ShouldBeNil)

				s2 := New()
				err = json.Unmarshal([]byte(`[5, 10]`), &s2)
				So(err, ShouldBeNil)

				s1.AndNot(&s2)
				So(s1.List(), ShouldResemble, []uint64{11, 13})
			})
		})
	})
}
