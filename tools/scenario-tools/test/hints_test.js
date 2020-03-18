const axios = require('axios')
const nock = require('nock')
const test = require('ava')
const { PROGRESS_SERVER_URL } = require('../lib/io')
const { showHint, nextHint } = require('../lib/hints')
const { fixture, testoutput, createSpyingLogger } = require('./helpers')

axios.defaults.adapter = require('axios/lib/adapters/http')

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

test.serial('nextHint progresses hints for the task', async t => {
  const spyinglogger = createSpyingLogger()
  const progresspath = testoutput('nexthinttest.json')
  const taskspath = fixture('tasks.yaml')
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

  // First POST will be the POST to save the initial progress so we start
  // start testing from index 1
  await nextHint(1, taskspath, progresspath, spyinglogger)
  t.deepEqual('{"name":"test-scenario","currentTask":null,"tasks":[{"id":1,"lastHintIndex":0,"score":null,"scoringSkipped":false}]}',
    posts[1], 'progress wasnt saved after showing first hint')
  t.false(spyinglogger.error.called, 'shouldnt have logged an error')
})
