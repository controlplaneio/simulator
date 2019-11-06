const test = require('ava')
const { transformV0ToV1, removeTaskPrefix, groupHintsByTask } = require('../lib/migrate')

test('groupHintsByTask', t => {
  const original = [
    { text: 'Task 1: test hint' },
    { text: 'Task 1: test hint' },
    { text: 'Task 2: test hint' }
  ]

  const expected = {
    'Task 1': {
      'sort-order': 1,
      'starting-point': '',
      hints: [{ text: 'test hint' }, { text: 'test hint' }]
    },
    'Task 2': {
      'sort-order': 2,
      'starting-point': '',
      hints: [{ text: 'test hint' }]
    }
  }

  t.deepEqual(expected, groupHintsByTask(original))
})

test('transforms v0 to v1 schema', t => {
  const original = [{
    general_overview: {
      objective: 'test objective',
      'starting-point': 'kubectl pod exec.',
      'num-hints': '2'
    }
  }, {
    hints: {
      'hint-1': 'Task 1: test hint 1',
      'hint-2': 'Task 1: test hint 2'
    }
  }
  ]

  const expected = {
    kind: 'cp.simulator/scenario:1.0.0',
    objective: 'test objective',
    'starting-point': 'kubectl pod exec.',
    tasks: {
      'Task 1': {
        'sort-order': 1,
        'starting-point': '',
        hints: [{ text: 'test hint 1' }, { text: 'test hint 2' }]
      }
    }
  }

  t.deepEqual(expected, transformV0ToV1(original))
})

test('remove task prefix removes the prefixes', t => {
  const original = {
    kind: 'cp.simulator/scenario:1.0.0',
    objective: 'test objective',
    'starting-point': 'kubectl pod exec.',
    tasks: {
      'Task 1': {
        'sort-order': 1,
        'starting-point': '',
        hints: [{ text: 'test hint 1' }, { text: 'test hint 2' }]
      },
      'Task 2': {
        'sort-order': 1,
        'starting-point': '',
        hints: [{ text: 'test hint 1' }, { text: 'test hint 2' }]
      }

    }
  }

  const expected = {
    kind: 'cp.simulator/scenario:1.0.0',
    objective: 'test objective',
    'starting-point': 'kubectl pod exec.',
    tasks: {
      '1': { // eslint-disable-line
        'sort-order': 1,
        'starting-point': '',
        hints: [{ text: 'test hint 1' }, { text: 'test hint 2' }]
      },
      '2': { // eslint-disable-line
        'sort-order': 1,
        'starting-point': '',
        hints: [{ text: 'test hint 1' }, { text: 'test hint 2' }]
      }

    }
  }

  t.deepEqual(expected, removeTaskPrefix(original))
})
