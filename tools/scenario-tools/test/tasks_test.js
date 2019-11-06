const { unlinkSync, readFileSync } = require('fs')
const test = require('ava')
const { startTask, getCurrentTask } = require('../lib/tasks')
const { fixture, testoutput, createSpyingLogger } = require('./helpers')

test('startTask warns for invalid task and returns false', t => {
  const taskspath = fixture('tasks.yaml')
  const progresspath = testoutput('progress.json')
  const logger = createSpyingLogger()

  const result = startTask('Invalid', taskspath, progresspath, logger)

  t.false(result, 'should have returned false')
  t.true(logger.warn.called, 'should have logged a warning')
})

test('startTask writes current_task to progress', t => {
  const taskspath = fixture('tasks.yaml')
  const progresspath = testoutput('progress.json')
  const logger = createSpyingLogger()

  // Delete any testoutput from previous runs
  try { unlinkSync(progresspath) } catch {}

  const result = startTask('Task 1', taskspath, progresspath, logger)
  const progress = readFileSync(progresspath, 'utf-8')

  t.true(result, 'should have returned true')
  t.deepEqual('{"current_task":"Task 1"}', progress, 'should have written progress')
  t.true(logger.info.called, 'should have logged an inffo message')
})

test('getCurrentTask returns the current task', t => {
  const progresspath = fixture('progress.json')

  const task = getCurrentTask(progresspath)
  t.is('Task 1', task)
})

test('getCurrentTask returns undefined when no current task', t => {
  const progresspath = fixture('does-not-exist.json')

  const task = getCurrentTask(progresspath)
  t.is(undefined, task)
})
