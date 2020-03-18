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
  progress.currentTask = newTask

  if (progress.tasks.find(t => t.id === newTask) === undefined) {
    progress.tasks.push({
      id: Number(newTask),
      lastHintIndex: null,
      score: null,
      scoringSkipped: false
    })
  }

  return progress
}

function processResponse (answer, progress, tasks, log = logger) {
  if (answer === undefined) {
    throw new Error(
      'No changes made - expected an answer from the scoring prompt')
  }

  const currentTask = progress.currentTask

  if (answer === 'cancel') {
    log.info(`You cancelled... leaving you on ${currentTask}`)
    return false
  }

  if (answer === 'yes') {
    const score = calculateScore(progress, tasks)
    log.info(`Your score for task ${currentTask} was ${score}`)
    const task = progress.tasks.find(t => t.id === currentTask)
    task.score = score
    task.scoringSkipped = false
  } else if (answer === 'no') {
    log.info(`You chose not to be scored on task ${currentTask}`)
    const task = progress.tasks.find(t => t.id === currentTask)
    task.score = null
    task.scoringSkipped = true
  } else {
    throw new Error(
      `No changes made - unrecognised answer from scoring prompt: ${answer}`)
  }

  return progress
}

// This is the entrypoint for both `start_task` and `end_task` if no `newTask`
// argument was supplied we assume this was invoked as `end_task`
async function processTask (newTask, taskspath = TASKS_FILE_PATH,
  progresspath = PROGRESS_FILE_PATH, log = logger) {
  const { name, tasks } = loadYamlFile(taskspath)

  if (newTask !== undefined && !tasks[newTask]) {
    log.warn('Cannot find task')
    return false
  }

  const progress = getProgress(name, progresspath)

  let newProgress

  if (newTask !== undefined) {
    newProgress = await startTask(newTask, tasks, progress, log, askToBeScored)
  } else {
    newProgress = await endTask(tasks, progress, log, askToBeScored)
  }

  if (newProgress !== false) {
    saveProgress(newProgress, progresspath)
  }
}

async function endTask (tasks, progress, log, prompter) {
  const currentTask = progress.currentTask

  if (currentTask === null) {
    log.warn('Cannot end task - you have not started one')
    return false
  }

  const { answer } = await prompter(currentTask)
  const newProgress = processResponse(answer, progress, tasks, log)
  if (newProgress === false) return false

  progress.currentTask = null
  log.info(`You have ended task ${currentTask}`)
  return progress
}

async function startTask (newTask, tasks, progress, log, prompter) {
  const currentTask = progress.currentTask

  // User is trying to switch to the task they are already on
  if (newTask === currentTask) {
    log.warn(`You are already on task ${currentTask}`)
    return false
  }

  // User hasnt started a task yet
  if (currentTask === null) {
    log.info(`You are now on task ${newTask}`)
    return updateProgressWithNewTask(progress, newTask)
  }

  // user has started a task and previously either asked not to be scored or
  // was already scored
  if (progress.tasks.find(t => t.id === currentTask).score !== null) {
    log.info(`You are now on task ${newTask}`)
    return updateProgressWithNewTask(progress, newTask)
  }

  // user must have started a task and hasn't yet been scored or skipped scoring
  const { answer } = await prompter(currentTask)
  const newProgress = processResponse(answer, progress, tasks, log)
  if (newProgress === false) return false

  if (newTask !== null) {
    log.info(`You are now on task ${newTask}`)
    return updateProgressWithNewTask(newProgress, newTask)
  }
}

function getCurrentTask (progresspath = PROGRESS_FILE_PATH,
  taskspath = TASKS_FILE_PATH) {
  const { name } = loadYamlFile(taskspath)

  const progress = getProgress(name, progresspath)
  return progress.currentTask
}

module.exports = {
  endTask,
  getCurrentTask,
  processTask,
  processResponse,
  startTask,
  updateProgressWithNewTask
}
