#!/usr/bin/env node

const { createLogger } = require('../lib/logger')
const { getCurrentTask } = require('../.lib/tasks.js')
const { nextHint } = require('../.lib/hints.js')

const logger = createLogger({})
require('../lib/error-handler')

const task = getCurrentTask()
if (task === undefined) {
  logger.error('You have not started a task.  Please run `start_task` to select your task')
  process.exit(1)
}

nextHint(task)
process.exit(0)
