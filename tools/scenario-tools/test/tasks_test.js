const { unlinkSync, readFileSync } = require('fs')
const { spy, fake } = require('sinon')
const test = require('ava')
const {
  endTask,
  getCurrentTask,
  processTask,
  processResponse,
  startTask,
  updateProgressWithNewTask
} = require('../lib/tasks')
const { fixture, testoutput, createSpyingLogger } = require('./helpers')

test('processTask warns for invalid task and returns false', async t => {
  const taskspath = fixture('tasks.yaml')
  const progresspath = testoutput('progress.json')
  const logger = createSpyingLogger()

  const result = await processTask('Invalid', taskspath, progresspath, logger)

  t.false(result, 'should have returned false')
  t.true(logger.warn.called, 'should have logged a warning')
})

test('processTask writes current_task to progress', async t => {
  const taskspath = fixture('tasks.yaml')
  const progresspath = testoutput('progress.json')
  const logger = createSpyingLogger()

  // Delete any testoutput from previous runs
  try { unlinkSync(progresspath) } catch {}

  await processTask('Task 1', taskspath, progresspath, logger)
  const progress = readFileSync(progresspath, 'utf-8')

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

test('processResponse throws when answer is not defined', t => {
  t.throws(() => processResponse(undefined), {
    message: 'No changes made - expected an answer from the scoring prompt'
  }, 'should have thrown an error')
})

test('processResponse returns false when answer is cancel', t => {
  const progress = {
    current_task: 1,
    1: { lastHintIndex: 0, score: undefined }
  }
  const result = processResponse('cancel', progress)

  t.false(result, 'should have returned false')
})

test('processResponse sets score to skip when answer is no', t => {
  const progress = {
    current_task: 1,
    1: { lastHintIndex: 0, score: undefined }
  }
  const tasks = {
    1: {
      hints: [
        { test: 'hint 1', penalty: 10 }
      ]
    }
  }
  const logger = createSpyingLogger()
  const result = processResponse('no', progress, tasks, logger)

  t.truthy(result, 'should have returned new progress')
  t.deepEqual(result, {
    current_task: 1,
    1: { lastHintIndex: 0, score: 'skip' }
  }, 'should have set score to skip')
  t.true(logger.info.called, 'should have logged a message')
  t.true(logger.info.calledWith('You chose not to be scored on task 1'),
    'should have told the user they chose not to be scored')
})

test('processResponse sets score when answer is yes', t => {
  const progress = {
    current_task: 1,
    1: { lastHintIndex: 0, score: undefined }
  }
  const tasks = {
    1: {
      hints: [
        { test: 'hint 1', penalty: 10 }
      ]
    }
  }
  const logger = createSpyingLogger()
  const result = processResponse('yes', progress, tasks, logger)

  t.truthy(result, 'should have returned new progress')
  t.deepEqual(result, {
    current_task: 1,
    1: { lastHintIndex: 0, score: 90 }
  }, 'should have set score to skip')
  t.true(logger.info.called, 'should have logged a message')
  t.true(logger.info.calledWith('Your score for task 1 was 90'),
    'should have told the user their score')
})

test('endTask warns if the user has not started a task', async t => {
  const progress = {
    current_task: undefined
  }
  const tasks = {
    1: {
      hints: [
        { test: 'hint 1', penalty: 10 }
      ]
    }
  }
  const logger = createSpyingLogger()
  const prompter = spy()

  const result = await endTask(tasks, progress, logger, prompter)
  t.false(result, 'should have returned false to indicate it did nothing')
  t.true(logger.warn.called, 'should have warned the user')
  t.true(logger.warn.calledWith('Cannot end task - you have not started one'))
})

test('endTask returns false if the user cancels', async t => {
  const progress = {
    current_task: 1,
    1: { lastHintIndex: 0, score: undefined }
  }
  const tasks = {
    1: {
      hints: [
        { test: 'hint 1', penalty: 10 }
      ]
    }
  }
  const logger = createSpyingLogger()
  const prompter = fake.returns({ answer: 'cancel' })

  const result = await endTask(tasks, progress, logger, prompter)
  t.false(result, 'should have returned false to indicate it did nothing')
})

test('endTask skips scoring and ends task when user says no to scoring', async t => {
  const progress = {
    current_task: 1,
    1: { lastHintIndex: 0, score: undefined }
  }
  const tasks = {
    1: {
      hints: [
        { test: 'hint 1', penalty: 10 }
      ]
    }
  }
  const logger = createSpyingLogger()
  const prompter = fake.returns({ answer: 'no' })

  const result = await endTask(tasks, progress, logger, prompter)
  t.deepEqual(result, {
    current_task: undefined,
    1: { lastHintIndex: 0, score: 'skip' }
  }, 'should have changed task progress correctly')
})

test('endTask sets score and ends task when user says yes to scoring', async t => {
  const progress = {
    current_task: 1,
    1: { lastHintIndex: 0, score: undefined }
  }
  const tasks = {
    1: {
      hints: [
        { test: 'hint 1', penalty: 10 }
      ]
    }
  }
  const logger = createSpyingLogger()
  const prompter = fake.returns({ answer: 'yes' })

  const result = await endTask(tasks, progress, logger, prompter)
  t.deepEqual(result, {
    current_task: undefined,
    1: { lastHintIndex: 0, score: 90 }
  }, 'should have changed task progress correctly')
})

test('startTask warns if user tries to start the same task', async t => {
  const progress = {
    current_task: 1,
    1: { lastHintIndex: 0, score: undefined }
  }
  const tasks = {
    1: {
      hints: [
        { test: 'hint 1', penalty: 10 }
      ]
    }
  }
  const logger = createSpyingLogger()
  const prompter = spy()

  const result = await startTask(1, tasks, progress, logger, prompter)

  t.false(result, 'should have returned false to indicate it did nothing')
  t.true(logger.warn.called, 'should have warned the user')
  t.true(logger.warn.calledWith('You are already on task 1'))
})

test('startTask initialises progress when user has no progress', async t => {
  const progress = { }
  const tasks = {
    1: {
      hints: [
        { test: 'hint 1', penalty: 10 }
      ]
    }
  }
  const logger = createSpyingLogger()
  const prompter = spy()

  const result = await startTask(1, tasks, progress, logger, prompter)

  t.truthy(result, 'should have returned progress')
  t.true(logger.info.called, 'should have messaged the user')
  t.true(logger.info.calledWith('You are now on task 1'))
  t.deepEqual(result, {
    current_task: 1,
    1: { lastHintIndex: undefined, score: undefined }
  }, 'should have initialised progress')
})

test('startTask updates current_task and doesnt rescore if already scored', async t => {
  const progress = {
    current_task: 1,
    1: { lastHintIndex: 0, score: 90 }
  }
  const tasks = {
    1: { hints: [{ test: 'hint 1', penalty: 10 }] },
    2: { hints: [{ test: 'hint 1', penalty: 10 }] }
  }
  const logger = createSpyingLogger()
  const prompter = spy()

  const result = await startTask(2, tasks, progress, logger, prompter)

  t.false(prompter.called, 'should not have prompted to do scoring')
  t.true(logger.info.called, 'should have messaged the user')
  t.true(logger.info.calledWith('You are now on task 2'))
  t.deepEqual(result, {
    current_task: 2,
    1: { lastHintIndex: 0, score: 90 },
    2: { lastHintIndex: undefined, score: undefined }
  }, 'should have only updated current task')
})

test('startTask updates current_task and skips scoring if answer is no', async t => {
  const progress = {
    current_task: 1,
    1: { lastHintIndex: 0, score: undefined }
  }
  const tasks = {
    1: { hints: [{ test: 'hint 1', penalty: 10 }] },
    2: { hints: [{ test: 'hint 1', penalty: 10 }] }
  }
  const logger = createSpyingLogger()
  const prompter = fake.returns({ answer: 'no' })

  const result = await startTask(2, tasks, progress, logger, prompter)

  t.true(logger.info.called, 'should have messaged the user')
  t.true(logger.info.calledWith('You are now on task 2'))
  t.deepEqual(result, {
    current_task: 2,
    1: { lastHintIndex: 0, score: 'skip' },
    2: { lastHintIndex: undefined, score: undefined }
  }, 'should have only updated current task')
})

test('startTask updates current_task and scores if answer is yes', async t => {
  const progress = {
    current_task: 1,
    1: { lastHintIndex: 0, score: undefined }
  }
  const tasks = {
    1: { hints: [{ test: 'hint 1', penalty: 10 }] },
    2: { hints: [{ test: 'hint 1', penalty: 10 }] }
  }
  const logger = createSpyingLogger()
  const prompter = fake.returns({ answer: 'yes' })

  const result = await startTask(2, tasks, progress, logger, prompter)

  t.true(logger.info.called, 'should have messaged the user')
  t.true(logger.info.calledWith('You are now on task 2'))
  t.deepEqual(result, {
    current_task: 2,
    1: { lastHintIndex: 0, score: 90 },
    2: { lastHintIndex: undefined, score: undefined }
  }, 'should have only updated current task')
})
