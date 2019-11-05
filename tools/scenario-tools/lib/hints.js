const { loadYamlFile, getProgress, saveProgress } = require('./io')
const { createLogger } = require('./logger')
const { TASKS_FILE_PATH, PROGRESS_FILE_PATH } = require('./constants')

const logger = createLogger({})

// All the optional parameters in these functions are to enable injection
// for testing

function showHint (task, index, taskspath = TASKS_FILE_PATH, log = logger) {
  const { tasks } = loadYamlFile(taskspath)

  if (!tasks[task]) {
    return log.warn('Cannot find task')
  }

  if (index >= tasks[task].hints.length) {
    return log.warn('There are no more hints for this task')
  }

  const hint = tasks[task].hints[index]

  log.info(hint)
  log.info('Use `scenario show-hints --task <task>` to reshow all the hints you have seen so far')
}

function showHints (task, taskspath = TASKS_FILE_PATH,
  progresspath = PROGRESS_FILE_PATH, log = logger) {
  const progress = getProgress(progresspath)
  const { tasks } = loadYamlFile(taskspath)

  if (!tasks[task]) {
    return logger.warn('Cannot find task')
  }

  if (!progress[task]) {
    return logger.info('You have not seen any hints for this task')
  }

  const lastSeenHintIndex = progress[task]

  for (let i = 0; i <= lastSeenHintIndex; i++) {
    logger.info(tasks[task].hints[i])
  }
}

function nextHint (task, taskspath = TASKS_FILE_PATH,
  progresspath = PROGRESS_FILE_PATH, log = logger) {
  if (!task) {
    return logger.error('No task provided to nextHint')
  }

  const progress = getProgress(progresspath)
  let hintIndex

  if (!progress[task]) {
    hintIndex = progress[task] = 0
  } else {
    hintIndex = progress[task]++
  }
  saveProgress(progress, progresspath)

  showHint(task, hintIndex, taskspath, log)
}

module.exports = {
  showHint,
  showHints,
  nextHint
}
