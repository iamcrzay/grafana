package conditions

import (
	"testing"

	"github.com/grafana/grafana/pkg/components/simplejson"
	. "github.com/smartystreets/goconvey/convey"
)

func evalutorScenario(json string, reducedValue float64, datapoints ...float64) bool {
	jsonModel, err := simplejson.NewJson([]byte(json))
	So(err, ShouldBeNil)

	evaluator, err := NewAlertEvaluator(jsonModel)
	So(err, ShouldBeNil)

	return evaluator.Eval(&reducedValue)
}

func TestEvalutors(t *testing.T) {
	Convey("greater then", t, func() {
		So(evalutorScenario(`{"type": "gt", "params": [1] }`, 3), ShouldBeTrue)
		So(evalutorScenario(`{"type": "gt", "params": [3] }`, 1), ShouldBeFalse)
	})

	Convey("less then", t, func() {
		So(evalutorScenario(`{"type": "lt", "params": [1] }`, 3), ShouldBeFalse)
		So(evalutorScenario(`{"type": "lt", "params": [3] }`, 1), ShouldBeTrue)
	})

	Convey("within_range", t, func() {
		So(evalutorScenario(`{"type": "within_range", "params": [1, 100] }`, 3), ShouldBeTrue)
		So(evalutorScenario(`{"type": "within_range", "params": [1, 100] }`, 300), ShouldBeFalse)
		So(evalutorScenario(`{"type": "within_range", "params": [100, 1] }`, 3), ShouldBeTrue)
		So(evalutorScenario(`{"type": "within_range", "params": [100, 1] }`, 300), ShouldBeFalse)
	})

	Convey("outside_range", t, func() {
		So(evalutorScenario(`{"type": "outside_range", "params": [1, 100] }`, 1000), ShouldBeTrue)
		So(evalutorScenario(`{"type": "outside_range", "params": [1, 100] }`, 50), ShouldBeFalse)
		So(evalutorScenario(`{"type": "outside_range", "params": [100, 1] }`, 1000), ShouldBeTrue)
		So(evalutorScenario(`{"type": "outside_range", "params": [100, 1] }`, 50), ShouldBeFalse)
	})

	Convey("no_data", t, func() {
		So(evalutorScenario(`{"type": "no_data", "params": [] }`, 50), ShouldBeFalse)

		jsonModel, err := simplejson.NewJson([]byte(`{"type": "no_data", "params": [] }`))
		So(err, ShouldBeNil)

		evaluator, err := NewAlertEvaluator(jsonModel)
		So(err, ShouldBeNil)

		So(evaluator.Eval(nil), ShouldBeTrue)
	})
}
