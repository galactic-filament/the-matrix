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
  exec "cd ./#{repoName} && #{cmd}", (err) -> cb err
withRepos = (repos, done, iterator) ->
  async.each repos, iterator, (err) -> done err

# derived commands
cloneRepo = (repoName, cb) ->
  runCmd "git clone https://github.com/ihsw/#{repoName}.git", cb
deleteRepo = (repoName, cb) -> runCmd "rm -rf ./#{repoName}", cb

# derived repo commands
buildRepo = (repoName, cb) -> repoCmd repoName, "./bin/build-images", cb
upWeb = (repoName, cb) -> repoCmd repoName, 'docker-compose up -d web', cb
stopWeb = (repoName, cb) -> repoCmd repoName, 'docker-compose stop web', cb

describe 'Arithmetic', ->
  before (done) ->
    withRepos repos, done, (repoName, eachNext) ->
      tasks = [
        (seriesNext) -> cloneRepo repoName, (err) -> seriesNext err
        (seriesNext) -> buildRepo repoName, (err) -> seriesNext err
        # (seriesNext) -> upWeb repoName, (err) -> seriesNext err
      ]
      async.series tasks, (err) -> eachNext err
  after (done) ->
    withRepos repos, done, (repoName, eachNext) ->
      tasks = [
        # (seriesNext) -> stopWeb repoName, (err) -> seriesNext err
        (seriesNext) -> deleteRepo repoName, (err) -> seriesNext err
      ]
      async.series tasks, (err) -> eachNext err
  it 'should add two numbers', -> expect(2+2).to.equal 4
