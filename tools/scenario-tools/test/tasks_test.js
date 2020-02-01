const { unlinkSync, readFileSync } = require('fs')
const test = require('ava')
const { updateProgressWithNewTask, startTask, getCurrentTask } = require('../lib/tasks')
const { fixture, testoutput, createSpyingLogger } = require('./helpers')

test('startTask warns for invalid task and returns false', async t => {
  const taskspath = fixture('tasks.yaml')
  const progresspath = testoutput('progress.json')
  const logger = createSpyingLogger()

  const result = await startTask('Invalid', taskspath, progresspath, logger)

  t.false(result, 'should have returned false')
  t.true(logger.warn.called, 'should have logged a warning')
})

test('startTask writes current_task to progress', async t => {
  const taskspath = fixture('tasks.yaml')
  const progresspath = testoutput('progress.json')
  const logger = createSpyingLogger()

  // Delete any testoutput from previous runs
  try { unlinkSync(progresspath) } catch {}

  const result = await startTask('Task 1', taskspath, progresspath, logger)
  const progress = readFileSync(progresspath, 'utf-8')

  t.true(result, 'should have returned true')
  t.deepEqual('{"current_task":"Task 1","Task 1":{}}', progress, 'should have written progress')
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

test('updateProgressWithNewTask sets current_task', t => {
  const progress = {}

  const result = updateProgressWithNewTask(progress, 1)

  t.is(result.current_task, 1, 'should have set current_task')
})

test('updateProgressWithNewTask initialises task progress', t => {
  const progress = {}

  const result = updateProgressWithNewTask(progress, 1)

  t.truthy(result[1], 'should have initialised an object for task progress')
  t.true(Object.prototype.hasOwnProperty.call(result[1], 'lastHintIndex'),
    'should have initialised a lastHintIndex property on task')
  t.true(Object.prototype.hasOwnProperty.call(result[1], 'score'),
    'should have initialised a score property on task')
})

test('updateProgressWithNewTask does not overwrite existing progress', t => {
  const progress = {
    current_task: 2,
    1: { lastHintIndex: 1, score: 'skip' },
    2: { lastHintIndex: undefined, score: undefined }
  }

  const result = updateProgressWithNewTask(progress, 1)

  t.deepEqual(result[1], { lastHintIndex: 1, score: 'skip' },
    'should not have changed existing task progress')
})
