const { loadYamlFile, getProgress, saveProgress } = require('./io')
const { createLogger } = require('./logger')
const { TASKS_FILE_PATH, PROGRESS_FILE_PATH } = require('./constants')

const logger = createLogger({})

function startTask (task, taskspath = TASKS_FILE_PATH,
  progresspath = PROGRESS_FILE_PATH, log = logger) {
  const { tasks } = loadYamlFile(taskspath)

  if (!tasks[task]) {
    log.warn('Cannot find task')
    return false
  }

  const progress = getProgress(progresspath)
  progress.current_task = task
  saveProgress(progress, progresspath)

  log.info(`You are now on ${task}`)

  return true
}

function getCurrentTask (progresspath = PROGRESS_FILE_PATH) {
  const progress = getProgress(progresspath)
  return progress.current_task
}

module.exports = {
  startTask,
  getCurrentTask
}
