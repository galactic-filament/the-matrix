expect = require('chai').expect
exec = require('child_process').exec
async = require 'async'
assert = require 'assert'

repos = [
  'omega-jazz'
  'pho-sho'
  'go-home'
  'py-lyfe'
]

# utility commands
runCmd = (cmd, cb) -> exec cmd, (err, stdout, stderr) -> cb err
repoCmd = (repoName, cmd, cb) ->
  exec "cd ./repos/#{repoName} && #{cmd}", (err) -> cb err
withRepos = (repos, done, iterator) ->
  async.each repos, iterator, (err) -> done err

# derived commands
cloneRepo = (repoName, cb) ->
  runCmd(
    "git clone https://github.com/ihsw/#{repoName}.git ./repos/#{repoName}"
    cb
  )
deleteRepo = (repoName, cb) -> runCmd "rm -rf ./repos/#{repoName}", cb

# derived repo commands
buildImage = (repoName, cb) -> repoCmd repoName, "./bin/build-images", cb
upWeb = (repoName, cb) -> repoCmd repoName, 'docker-compose up -d base', cb
stopWeb = (repoName, cb) -> repoCmd repoName, 'docker-compose stop base', cb

describe 'Api Servers', ->
  before (done) ->
    withRepos repos, done, (repoName, eachNext) ->
      tasks = [
        (seriesNext) -> cloneRepo repoName, (err) -> seriesNext err
        (seriesNext) -> buildImage repoName, (err) -> seriesNext err
      ]
      async.series tasks, (err) -> eachNext err
  after (done) ->
    withRepos repos, done, (repoName, eachNext) ->
      tasks = [
        (seriesNext) -> deleteRepo repoName, (err) -> seriesNext err
      ]
      async.series tasks, (err) -> eachNext err
  it 'should run the test suite', (done) ->
    withRepos repos, done, (repoName, eachNext) ->
      tasks = [
        (seriesNext) -> upWeb repoName, (err) -> seriesNext err
        (seriesNext) -> stopWeb repoName, (err) -> seriesNext err
      ]
      async.series tasks, (err) -> eachNext err
