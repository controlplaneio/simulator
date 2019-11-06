const test = require('ava')
const { parse, showUsage } = require('../lib/args')

test('parse migrate with name option', t => {
  const args = ['migrate', '--name', 'container-ambush', '--type', 'legacy']
  const parsed = parse(args)

  t.deepEqual({
    command: 'migrate',
    options: {
      name: 'container-ambush',
      all: false,
      type: 'legacy'
    }
  }, parsed)
})

test('parse migrate throws for missing type', t => {
  const args = ['migrate', '--name', 'container-ambush']
  t.throws(() => parse(args))
})

test('parse migrate with all option', t => {
  const args = ['migrate', '--all', '--type', 'legacy']
  const parsed = parse(args)

  t.deepEqual({ command: 'migrate', options: { all: true, type: 'legacy' } }, parsed)
})

test('showUsage', t => {
  showUsage()
  showUsage(true)
  // just print so we can eyeball the output and make sure it doesnt error
  t.pass()
})

test('parse help options', t => {
  const args = ['help']
  const parsed = parse(args)

  t.deepEqual({ command: 'help', options: {} }, parsed)
})
