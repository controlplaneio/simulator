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
    return log.warn(`There are no more hints for ${task}`)
  }

  const hint = tasks[task].hints[index]

  log.info(hint.text)
  log.info(`This hint incurred a penalty of ${hint.penalty} to your score`)
}

async function showHints (task, taskspath = TASKS_FILE_PATH,
  progresspath = PROGRESS_FILE_PATH, log = logger) {
  const { name, tasks } = loadYamlFile(taskspath)
  const progress = await getProgress(name, progresspath)

  if (!tasks[task]) {
    return logger.warn('Cannot find task')
  }

  const taskProgress = progress.tasks.find(t => t.id === task)
  if (taskProgress === undefined || taskProgress.lastHintIndex === null) {
    return logger.info(`You have not seen any hints for ${task}`)
  }

  const lastHintIndex = taskProgress.lastHintIndex
  const hintcount = tasks[task].hints.length

  for (let i = 0; i <= lastHintIndex && i < hintcount; i++) {
    logger.info(tasks[task].hints[i].text)
  }

  if (lastHintIndex >= hintcount - 1) {
    return logger.info('You have seen all the hints for this task')
  }
}

async function nextHint (task, taskspath = TASKS_FILE_PATH,
  progresspath = PROGRESS_FILE_PATH, log = logger) {
  if (!task) {
    return logger.error('No task provided to nextHint')
  }
  const { name, tasks } = loadYamlFile(taskspath)

  const progress = await getProgress(name, progresspath)
  let hintIndex

  let taskProgress = progress.tasks.find(t => t.id === task)

  if (taskProgress === undefined) {
    taskProgress = {
      id: task,
      lastHintIndex: null,
      score: null,
      scoringSkipped: false
    }
    progress.tasks.push(taskProgress)
  }

  if (taskProgress.lastHintIndex === null) {
    hintIndex = taskProgress.lastHintIndex = 0
  } else {
    hintIndex = taskProgress.lastHintIndex = taskProgress.lastHintIndex + 1
  }

  if (hintIndex >= tasks[task].hints.length) {
    return logger.info('You have seen all the hints for this task')
  }

  await saveProgress(progress, progresspath)

  showHint(task, hintIndex, taskspath, log)
}

module.exports = {
  showHint,
  showHints,
  nextHint
}
