const { loadYamlFile, getProgress, saveProgress } = require('./io')
const { calculateScore } = require('./scoring')
const { createLogger } = require('./logger')
const { TASKS_FILE_PATH, PROGRESS_FILE_PATH } = require('./constants')
const { prompt } = require('inquirer')

const logger = createLogger({})

async function askToBeScored (currentTask) {
  logger.info(`Would you like to be scored for task ${currentTask}? (If you choose no you cannot be scored on this task in the future)`)
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

function updateProgressWithNewTask (progress, newTask) {
  progress.current_task = newTask

  if (progress[newTask] === undefined) {
    progress[newTask] = {
      lastHintIndex: undefined,
      score: undefined
    }
  }

  return progress
}

function scoreTask (answer, progress, tasks, currentTask, log = logger) {
  if (answer === undefined) {
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
    throw new Error(
      `No changes made - unrecognised answer from scoring prompt: ${answer}`)
  }
}

async function startTask (newTask, taskspath = TASKS_FILE_PATH,
  progresspath = PROGRESS_FILE_PATH, log = logger) {
  const { tasks } = loadYamlFile(taskspath)

  if (newTask !== undefined && !tasks[newTask]) {
    log.warn('Cannot find task')
    return false
  }

  const progress = getProgress(progresspath)

  const newProgress = await processTask(newTask, tasks, progress, log, askToBeScored)

  if (newProgress !== false) {
    saveProgress(newProgress, progresspath)
  }
}

async function processTask (newTask, tasks, progress, log, prompter) {
  const currentTask = progress.current_task

  if (newTask === undefined && currentTask === undefined) {
    log.warn('Cannot end task - you have not started one')
    return false
  }

  // User is trying to switch to the task they are already on
  if (newTask === currentTask) {
    log.warn(`You are already on ${currentTask}`)
    return false
  }

  // User hasnt started a task yet
  if (currentTask === undefined) {
    log.info(`You are now on task ${newTask}`)
    return updateProgressWithNewTask(progress, newTask)
  }

  // user has started a task and previously either asked not to be scored or
  // was already scored
  if (newTask !== undefined && progress[currentTask].score !== undefined) {
    log.info(`You are now on task ${newTask}`)
    return updateProgressWithNewTask(progress, newTask)
  }

  // user must have started a task and hasn't yet been scored or skipped scoring
  const { answer } = await askToBeScored(currentTask)
  scoreTask(answer, progress, tasks, currentTask, log)

  if (newTask !== undefined) {
    log.info(`You are now on task ${newTask}`)
    return updateProgressWithNewTask(progress, newTask)
  }

  delete progress.current_task
  log.info(`You have ended task ${currentTask}`)
  return progress
}

function getCurrentTask (progresspath = PROGRESS_FILE_PATH) {
  const progress = getProgress(progresspath)
  return progress.current_task
}

module.exports = {
  startTask,
  getCurrentTask,
  updateProgressWithNewTask
}
