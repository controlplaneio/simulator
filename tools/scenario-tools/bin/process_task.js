#!/usr/bin/env node

const { createLogger } = require('../lib/logger')
const { cloneArray } = require('../lib/helpers')
const { processTask } = require('../lib/tasks.js')

require('../lib/error-handler')

const args = cloneArray(process.argv)

args.shift() // remove `node` from argv
args.shift() // remove `scenario.js` from argv

if (args.length === 2 && args[1] === '--debug') {
  global.DEBUG = true
}

const logger = createLogger({})

if (process.argv0 === 'start_task' && args.length < 1) {
  logger.error(
    'Please provide the task you wish to start. These are listed in the instructions')
  process.exit(1)
}

const taskId = Number(args[0])

// ignore the result - it will be true/false depending on whether the user
// actually switched task
logger.debug('Processing task', { task: taskId })
processTask(taskId).then(_ => {
  process.exit(0)
}, reason => {
  logger.error(reason.message)
  process.exit(1)
})
