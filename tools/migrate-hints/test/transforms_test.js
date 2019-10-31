const test = require('ava')
const {transformV0ToV1} = require('../lib/transforms')

test('transforms v0 to v1 schema', t => {
  const original = [{
      general_overview: {
        objective: "test objective",
        "starting-point": "kubectl pod exec.",
        "num-hints": "2"
      }, {
      hints: {
        "hint-1": "test hint 1",
        "hint-2": "test hint 2"
      }
    }
  }]

  const expected = [{
      kind: "cp.simulator/scenario:1.0.0",
    }, {
      general_overview: {
        objective: "test objective",
        "starting-point": "kubectl pod exec."
      }, {
    hints: [ {
        "text": "test hint 1"
      }, {
        "text": "test hint 2"
      }
    }]
  }

  t.deepEqual(transformV0ToV1(original), expected)
})
