const test = require('ava')
const {parse, showUsage} = require('../lib/args')

test('parse migrate options', t => {
  const args = [ 'migrate' ]
  const parsed = parse(args)

  t.deepEqual({ command: 'migrate', options: {} }, parsed)
})

test('parse show-hints', t => {
  const args = [ 'show-hints', '--task', 'task-1' ]
  const parsed = parse(args)

  t.deepEqual({ command: 'show-hints', options: { task: 'task-1' } }, parsed)
})

test('parse next-hint', t => {
  const args = [ 'next-hint', '--task', 'task-1' ]
  const parsed = parse(args)

  t.deepEqual({ command: 'next-hint', options: { task: 'task-1' } }, parsed)
})

test('showUsage', t => {
  showUsage()
  showUsage(true)
  // just print so we can eyeball the output and make sure it doesnt error
  t.pass()
})

test('parse help options', t => {
  const args = [ 'help' ]
  const parsed = parse(args)

  t.deepEqual({ command: 'help', options: {} }, parsed)
})
