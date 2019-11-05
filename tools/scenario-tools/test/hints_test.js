const { resolve, join } = require('path')
const test = require('ava')
const { spy } = require('sinon')
const { showHint } = require('../lib/hints')

function fixture (name) {
  return resolve(join(__dirname, 'fixtures', name))
}

test('showHint shows a valid hint for a valid task', t => {
  const spyinglogger = {
    warn: spy(),
    info: spy()
  }

  showHint('Task 1', 0, fixture('tasks.yaml'), spyinglogger)
  t.false(spyinglogger.warn.called, 'should not have displayed a warning')
  t.true(spyinglogger.info.called, 'should have logged a hint')
})

test('showHint warns for an invalid task', t => {
  const spyinglogger = {
    warn: spy(),
    info: spy()
  }

  showHint('Task 2', 0, fixture('tasks.yaml'), spyinglogger)
  t.true(spyinglogger.warn.called, 'should have displayed a warning')
  t.false(spyinglogger.info.called, 'should not have logged a hint')
})

test('showHint warns for a hint out of range for an valid task', t => {
  const spyinglogger = {
    warn: spy(),
    info: spy()
  }

  showHint('Task 1', 5, fixture('tasks.yaml'), spyinglogger)
  t.true(spyinglogger.warn.called, 'should have displayed a warning')
  t.false(spyinglogger.info.called, 'should not have logged a hint')
})
