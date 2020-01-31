const { loadYamlFile, getProgress, saveProgress } = require('./io')
const { calculateScore } = require('./scoring')
const { createLogger } = require('./logger')
const { TASKS_FILE_PATH, PROGRESS_FILE_PATH } = require('./constants')
const { prompt } = require('inquirer')

const logger = createLogger({})

async function askToBeScored () {
  return prompt({
    type: 'expand',
    name: 'answer',
    choices: [{
      key: 'y',
      name: 'Yes',
      value: 'yes'
    }, {
      key: 'n',
      name: 'No',
      value: 'no'
    }, {
      key: 'c',
      name: 'Cancel',
      value: 'cancel'

    }]
  })
}

async function startTask (newTask, taskspath = TASKS_FILE_PATH,
  progresspath = PROGRESS_FILE_PATH, log = logger) {
  const { tasks } = loadYamlFile(taskspath)

  if (newTask !== undefined && !tasks[newTask]) {
    log.warn('Cannot find task')
    return false
  }

  const progress = getProgress(progresspath)
  const currentTask = progress.current_task

  if (newTask === undefined && currentTask === undefined) {
    log.warn('Cannot end task - you have not started one')
    return false
  }

  // User hasnt started a task yet
  if (currentTask === undefined) {
    progress.current_task = newTask
    progress[newTask] = {}
    saveProgress(progress, progresspath)
    log.info(`You are now on task ${newTask}`)

    return true
  }

  if (newTask === currentTask) {
    log.warn(`You are already on ${currentTask}`)
    return false
  }

  // user has started a task and previously either asked not to be scored or
  // was already scored
  if (newTask !== undefined && progress[currentTask].score !== undefined) {
    progress.current_task = newTask
    if (progress[newTask] === undefined) { progress[newTask] = {} }

    saveProgress(progress, progresspath)
    log.info(`You are now on task ${newTask}`)

    return true
  }

  // user must have started a task and hasn't yet been scored or skipped scoring
  logger.info(`Would you like to be scored for task ${currentTask}? (If you choose no you cannot be scored on this task in the future)`)
  const { answer } = await askToBeScored()
  if (answer === undefined) {
    // should never happen
    throw new Error(
      'No changes made - expected an answer from the scoring prompt')
  }

  if (answer === 'cancel') {
    log.info(`You cancelled... leaving you on ${currentTask}`)
    return false
  }

  if (answer === 'yes') {
    const score = calculateScore(progress, tasks)
    log.info(`Your score for task ${currentTask} was ${score}`)
    progress[currentTask].score = score
  } else if (answer === 'no') {
    log.info(`You chose not to be scored on ${currentTask}`)
    progress[currentTask].score = 'skip'
  } else {
    // should never happen
    throw new Error(
      `No changes made - unrecognised answer from scoring prompt: ${answer}`)
  }

  if (newTask !== undefined) {
    progress.current_task = newTask
    if (progress[newTask] === undefined) { progress[newTask] = {} }
    saveProgress(progress, progresspath)
    log.info(`You are now on task ${newTask}`)
  } else {
    delete progress.current_task
    saveProgress(progress, progresspath)
    log.info(`You have ended task ${currentTask}`)
  }

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
