const { unlinkSync, readFileSync } = require('fs')
const test = require('ava')
const { showHint, nextHint } = require('../lib/hints')
const { fixture, testoutput, createSpyingLogger } = require('./helpers')

test('showHint shows a valid hint for a valid task', t => {
  const spyinglogger = createSpyingLogger()

  showHint(1, 0, fixture('tasks.yaml'), spyinglogger)
  t.false(spyinglogger.warn.called, 'should not have displayed a warning')
  t.true(spyinglogger.info.called, 'should have logged a hint')
})

test('showHint warns for an invalid task', t => {
  const spyinglogger = createSpyingLogger()

  showHint(2, 0, fixture('tasks.yaml'), spyinglogger)
  t.true(spyinglogger.warn.called, 'should have displayed a warning')
  t.false(spyinglogger.info.called, 'should not have logged a hint')
})

test('showHint warns for a hint out of range for an valid task', t => {
  const spyinglogger = createSpyingLogger()

  showHint(1, 5, fixture('tasks.yaml'), spyinglogger)
  t.true(spyinglogger.warn.called, 'should have displayed a warning')
  t.false(spyinglogger.info.called, 'should not have logged a hint')
})

test('nextHint progresses hints for the task', t => {
  const spyinglogger = createSpyingLogger()
  const progresspath = testoutput('nexthinttest.json')
  const taskspath = fixture('tasks.yaml')

  // delete any test output hanging around
  try { unlinkSync(progresspath) } catch {}

  nextHint(1, taskspath, progresspath, spyinglogger)
  const progress = readFileSync(progresspath, 'utf-8')
  t.deepEqual('{"name":"test-scenario","currentTask":null,"tasks":[{"id":1,"lastHintIndex":0,"score":null,"scoringSkipped":false}]}',
    progress, 'progress wasnt saved after showing first hint')
  t.false(spyinglogger.error.called, 'shouldnt have logged an error')

  nextHint(1, taskspath, progresspath, spyinglogger)
  const progress2 = readFileSync(progresspath, 'utf-8')
  t.deepEqual('{"name":"test-scenario","currentTask":null,"tasks":[{"id":1,"lastHintIndex":1,"score":null,"scoringSkipped":false}]}',
    progress2, 'progress wasnt updated after showing second hint')
  t.false(spyinglogger.error.called, 'shouldnt have logged an error')

  nextHint(1, taskspath, progresspath, spyinglogger)
  const progress3 = readFileSync(progresspath, 'utf-8')
  t.deepEqual('{"name":"test-scenario","currentTask":null,"tasks":[{"id":1,"lastHintIndex":1,"score":null,"scoringSkipped":false}]}',
    progress3, 'progress should not have updated after showing second hint')
})
