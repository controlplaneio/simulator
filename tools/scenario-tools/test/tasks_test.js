const axios = require('axios')
const nock = require('nock')
const { spy, fake } = require('sinon')
const test = require('ava')
const { PROGRESS_SERVER_URL } = require('../lib/io')
const {
  endTask,
  getCurrentTask,
  processTask,
  processResponse,
  startTask,
  updateProgressWithNewTask
} = require('../lib/tasks')
const { fixture, testoutput, createSpyingLogger } = require('./helpers')

axios.defaults.adapter = require('axios/lib/adapters/http')

function fakeTasks () {
  return {
    1: { hints: [{ test: 'hint 1', penalty: 10 }] },
    2: { hints: [{ test: 'hint 1', penalty: 10 }] }
  }
}

function onTask1 (score) {
  if (!score) {
    score = null
  }
  const scoringSkipped = false

  return {
    name: 'testing-scenario',
    currentTask: 1,
    tasks: [{ id: 1, lastHintIndex: 0, score: score, scoringSkipped }]
  }
}

test('processTask warns for invalid task and returns false', async t => {
  const taskspath = fixture('tasks.yaml')
  const progresspath = testoutput('progress.json')
  const logger = createSpyingLogger()

  const result = await processTask('Invalid', taskspath, progresspath, logger)

  t.false(result, 'should have returned false')
  t.true(logger.warn.called, 'should have logged a warning')
})

test.serial('processTask posts currentTask', async t => {
  const taskspath = fixture('tasks.yaml')
  const progresspath = testoutput('progress.json')
  const logger = createSpyingLogger()
  const posts = []
  nock.cleanAll()
  nock(PROGRESS_SERVER_URL)
    .persist()
    .get('/?scenario=test-scenario')
    .reply(404)
  nock(PROGRESS_SERVER_URL)
    .persist()
    .post('/', body => posts.push(JSON.stringify(body)))
    .reply(200)

  await processTask(1, taskspath, progresspath, logger)

  t.log(posts)
  t.deepEqual('{"name":"test-scenario","currentTask":1,"tasks":[{"id":1,"lastHintIndex":null,"score":null,"scoringSkipped":false}]}',
    posts[posts.length - 1], 'should have written progress')
  t.true(logger.info.called, 'should have logged an info message')
})

test.serial('getCurrentTask returns the current task', async t => {
  const taskspath = fixture('tasks.yaml')
  const progresspath = fixture('progress.json')
  nock.cleanAll()
  nock(PROGRESS_SERVER_URL)
    .get('/?scenario=test-scenario')
    .reply(200, { name: 'test-scenario', currentTask: 1, tasks: [] })

  const task = await getCurrentTask(progresspath, taskspath)
  t.is(1, task)
})

test.serial('getCurrentTask returns null when no current task', async t => {
  const taskspath = fixture('tasks.yaml')
  const progresspath = fixture('does-not-exist.json')
  nock.cleanAll()
  nock(PROGRESS_SERVER_URL)
    .persist()
    .get('/?scenario=test-scenario')
    .reply(404)
  nock(PROGRESS_SERVER_URL)
    .persist()
    .post('/')
    .reply(200)

  const task = await getCurrentTask(progresspath, taskspath)
  t.is(null, task)
})

test('updateProgressWithNewTask sets currentTask', t => {
  const progress = {
    name: 'testing-scenario',
    currentTask: null,
    tasks: []
  }

  const result = updateProgressWithNewTask(progress, 1)

  t.is(result.currentTask, 1, 'should have set currentTask')
})

test('updateProgressWithNewTask initialises task progress', t => {
  const progress = {
    name: 'testing-scenario',
    currentTask: null,
    tasks: []
  }

  const result = updateProgressWithNewTask(progress, 1)

  t.truthy(result.tasks[0], 'should have initialised an object for task progress')
  t.true(Object.prototype.hasOwnProperty.call(result.tasks[0], 'id'),
    'should have initialised an id property on task')
  t.true(Object.prototype.hasOwnProperty.call(result.tasks[0], 'lastHintIndex'),
    'should have initialised a lastHintIndex property on task')
  t.true(Object.prototype.hasOwnProperty.call(result.tasks[0], 'score'),
    'should have initialised a score property on task')
  t.true(Object.prototype.hasOwnProperty.call(result.tasks[0], 'scoringSkipped'),
    'should have initialised a scoringSkipped property on task')
})

test('updateProgressWithNewTask does not overwrite existing progress', t => {
  const progress = {
    name: 'testing-scenario',
    currentTask: 2,
    tasks: [
      { id: 1, lastHintIndex: 1, score: null, scoringSkipped: true },
      { id: 2, lastHintIndex: 0, score: null, scoringSkipped: false }
    ]
  }

  const result = updateProgressWithNewTask(progress, 1)

  t.deepEqual(result.tasks.find(t => t.id === 1), {
    id: 1,
    lastHintIndex: 1,
    score: null,
    scoringSkipped: true
  }, 'should not have changed existing task progress')
})

test('processResponse throws when answer is not defined', t => {
  t.throws(() => processResponse(undefined), {
    message: 'No changes made - expected an answer from the scoring prompt'
  }, 'should have thrown an error')
})

test('processResponse returns false when answer is cancel', t => {
  const progress = onTask1()
  const result = processResponse('cancel', progress)

  t.false(result, 'should have returned false')
})

test('processResponse sets scoringSkipped when answer is no', t => {
  const progress = onTask1(undefined)
  const tasks = fakeTasks()
  const logger = createSpyingLogger()
  const result = processResponse('no', progress, tasks, logger)

  t.truthy(result, 'should have returned new progress')
  t.deepEqual(result, {
    name: 'testing-scenario',
    currentTask: 1,
    tasks: [{ id: 1, lastHintIndex: 0, score: null, scoringSkipped: true }]
  }, 'should have set score to skip')
  t.true(logger.info.called, 'should have logged a message')
  t.true(logger.info.calledWith('You chose not to be scored on task 1'),
    'should have told the user they chose not to be scored')
})

test('processResponse sets score when answer is yes', t => {
  const progress = onTask1()
  const tasks = fakeTasks()
  const logger = createSpyingLogger()
  const result = processResponse('yes', progress, tasks, logger)

  t.truthy(result, 'should have returned new progress')
  t.deepEqual(result, {
    name: 'testing-scenario',
    currentTask: 1,
    tasks: [{ id: 1, lastHintIndex: 0, score: 90, scoringSkipped: false }]
  }, 'should have set score to skip')
  t.true(logger.info.called, 'should have logged a message')
  t.true(logger.info.calledWith('Your score for task 1 was 90'),
    'should have told the user their score')
})

test('endTask warns if the user has not started a task', async t => {
  const progress = {
    name: 'testing-scenario',
    currentTask: null,
    tasks: []
  }
  const tasks = fakeTasks()
  const logger = createSpyingLogger()
  const prompter = spy()

  const result = await endTask(tasks, progress, logger, prompter)
  t.false(result, 'should have returned false to indicate it did nothing')
  t.true(logger.warn.called, 'should have warned the user')
  t.true(logger.warn.calledWith('Cannot end task - you have not started one'))
})

test('endTask returns false if the user cancels', async t => {
  const progress = onTask1()
  const tasks = fakeTasks()
  const logger = createSpyingLogger()
  const prompter = fake.returns({ answer: 'cancel' })

  const result = await endTask(tasks, progress, logger, prompter)
  t.false(result, 'should have returned false to indicate it did nothing')
})

test('endTask skips scoring and ends task when user says no to scoring', async t => {
  const progress = onTask1()
  const tasks = fakeTasks()
  const logger = createSpyingLogger()
  const prompter = fake.returns({ answer: 'no' })

  const result = await endTask(tasks, progress, logger, prompter)
  t.deepEqual(result, {
    name: 'testing-scenario',
    currentTask: null,
    tasks: [{ id: 1, lastHintIndex: 0, score: null, scoringSkipped: true }]
  }, 'should have changed task progress correctly')
})

test('endTask sets score and ends task when user says yes to scoring', async t => {
  const progress = onTask1()
  const tasks = fakeTasks()
  const logger = createSpyingLogger()
  const prompter = fake.returns({ answer: 'yes' })

  const result = await endTask(tasks, progress, logger, prompter)
  t.deepEqual(result, {
    name: 'testing-scenario',
    currentTask: null,
    tasks: [{ id: 1, lastHintIndex: 0, score: 90, scoringSkipped: false }]
  }, 'should have changed task progress correctly')
})

test('startTask warns if user tries to start the same task', async t => {
  const progress = onTask1()
  const tasks = fakeTasks()
  const logger = createSpyingLogger()
  const prompter = spy()

  const result = await startTask(1, tasks, progress, logger, prompter)

  t.false(result, 'should have returned false to indicate it did nothing')
  t.true(logger.warn.called, 'should have warned the user')
  t.true(logger.warn.calledWith('You are already on task 1'))
})

test('startTask initialises progress when user has no progress', async t => {
  const progress = {
    name: 'testing-scenario',
    currentTask: null,
    tasks: []
  }
  const tasks = fakeTasks()
  const logger = createSpyingLogger()
  const prompter = spy()

  const result = await startTask(1, tasks, progress, logger, prompter)

  t.truthy(result, 'should have returned progress')
  t.true(logger.info.called, 'should have messaged the user')
  t.true(logger.info.calledWith('You are now on task 1'))
  t.deepEqual(result, {
    name: 'testing-scenario',
    currentTask: 1,
    tasks: [
      { id: 1, lastHintIndex: null, score: null, scoringSkipped: false }
    ]
  }, 'should have initialised progress')
})

test('startTask updates current_task and doesnt rescore if already scored', async t => {
  const progress = onTask1(90)
  const tasks = fakeTasks()
  const logger = createSpyingLogger()
  const prompter = spy()

  const result = await startTask(2, tasks, progress, logger, prompter)

  t.false(prompter.called, 'should not have prompted to do scoring')
  t.true(logger.info.called, 'should have messaged the user')
  t.true(logger.info.calledWith('You are now on task 2'))
  t.deepEqual(result, {
    name: 'testing-scenario',
    currentTask: 2,
    tasks: [
      { id: 1, lastHintIndex: 0, score: 90, scoringSkipped: false },
      { id: 2, lastHintIndex: null, score: null, scoringSkipped: false }
    ]
  }, 'should have only updated current task')
})

test('startTask updates current_task and skips scoring if answer is no', async t => {
  const progress = onTask1()
  const tasks = fakeTasks()
  const logger = createSpyingLogger()
  const prompter = fake.returns({ answer: 'no' })

  const result = await startTask(2, tasks, progress, logger, prompter)

  t.true(logger.info.called, 'should have messaged the user')
  t.true(logger.info.calledWith('You are now on task 2'))
  t.deepEqual(result, {
    name: 'testing-scenario',
    currentTask: 2,
    tasks: [
      { id: 1, lastHintIndex: 0, score: null, scoringSkipped: true },
      { id: 2, lastHintIndex: null, score: null, scoringSkipped: false }
    ]
  }, 'should have only updated current task')
})

test('startTask updates current_task and scores if answer is yes', async t => {
  const progress = onTask1()
  const tasks = fakeTasks()
  const logger = createSpyingLogger()
  const prompter = fake.returns({ answer: 'yes' })

  const result = await startTask(2, tasks, progress, logger, prompter)

  t.true(logger.info.called, 'should have messaged the user')
  t.true(logger.info.calledWith('You are now on task 2'))
  t.deepEqual(result, {
    name: 'testing-scenario',
    currentTask: 2,
    tasks: [
      { id: 1, lastHintIndex: 0, score: 90, scoringSkipped: false },
      { id: 2, lastHintIndex: null, score: null, scoringSkipped: false }
    ]
  }, 'should have only updated current task')
})
