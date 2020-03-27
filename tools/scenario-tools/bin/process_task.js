#!/usr/bin/env node

const { basename } = require('path')
const { createLogger } = require('../lib/logger')
const { cloneArray } = require('../lib/helpers')
const { processTask } = require('../lib/tasks.js')

require('../lib/error-handler')

// [ '/usr/bin/node', '/usr/bin/<bin_name>', ... ]
const args = cloneArray(process.argv)
const invokedBy = basename(args[1])

args.shift() // remove `node` from argv
args.shift() // remove `scenario.js` from argv

if ((invokedBy === 'start_task' && args.length === 2 && args[1] === '--debug') ||
    (invokedBy === 'end_task' && args.length === 1 && args[0] === '--debug')) {
  console.log('Setting global DEBUG flag')
  global.DEBUG = true
}

const logger = createLogger({})

if (invokedBy === 'start_task' && args.length < 1) {
  logger.error(
    'Please provide the task you wish to start. These are listed in the instructions')
  process.exit(1)
}

let taskId
if (args[0] !== undefined && args[0] !== '--debug') {
  taskId = Number(args[0])
}

// ignore the result - it will be true/false depending on whether the user
// actually switched task
logger.debug('Processing task', { task: taskId })
processTask(taskId).then(_ => {
  process.exit(0)
}, reason => {
  logger.error(reason.message)
  process.exit(1)
})
