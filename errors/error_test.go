package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestError(t *testing.T) {
	Convey("Given errors are created", t, func() {
		e1 := New(fmt.Errorf("External error")).WithStatus(http.StatusBadRequest).Add("Real error")

		var T = struct {
			R []string `json:"errors"`
			S int      `json:"status"`
		}{}

		Convey("The content is applied as expected", func() {
			So(len(e1.Reasons), ShouldEqual, 2)
			So(e1.Origin.Error(), ShouldEqual, "External error")
			So(e1.Reasons[1], ShouldEqual, "Real error")
			So(e1.Status, ShouldEqual, http.StatusBadRequest)
		})

		Convey("The error is marshalled to JSON as expected", func() {
			err := json.Unmarshal(e1.AsJSON(), &T)

			So(err, ShouldBeNil)
			So(len(T.R), ShouldEqual, 2)
			So(T.S, ShouldEqual, http.StatusBadRequest)
			So(T.R[1], ShouldEqual, "Real error")
		})

		Convey("The error can be modified", func() {
			e1.WithStatus(http.StatusOK).CausedBy(fmt.Errorf("Updated"))

			So(e1.Status, ShouldEqual, http.StatusOK)
			So(len(e1.Reasons), ShouldEqual, 3)
			So(e1.Origin.Error(), ShouldEqual, "Updated")

			Convey("And the JSON will be updated", func() {
				err := json.Unmarshal(e1.AsJSON(), &T)

				So(err, ShouldBeNil)
				So(len(T.R), ShouldEqual, 3)
				So(T.S, ShouldEqual, http.StatusOK)
				So(T.R[2], ShouldEqual, "Updated")
			})
		})
	})
}
