package ribbons

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBits(t *testing.T) {
	Convey("Check bits", t, func() {
		Convey("Check set", func() {
			u := uint64(0)

			u = setBit(u, 3)
			So(u, ShouldEqual, 8)

			u = setBit(u, 0)
			So(u, ShouldEqual, 9)
		})

		Convey("Check del", func() {
			u := uint64(8)

			u = delBit(u, 3)
			So(u, ShouldEqual, 0)

			u = delBit(u, 3)
			So(u, ShouldEqual, 0)
		})

		Convey("Check has", func() {
			u := uint64(8)

			So(hasBit(u, 3), ShouldBeTrue)
			So(hasBit(u, 0), ShouldBeFalse)
		})

		Convey("Check extract", func() {
			u := uint64(9)

			bits := extractToggledBits(u, 0)
			So(bits, ShouldResemble, []uint64{0, 3})
		})
	})
}
