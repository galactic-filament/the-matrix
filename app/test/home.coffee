supertest = require 'supertest'
expect = require('chai').expect

request = supertest 'http://ApiServer'
describe 'Homepage', ->
  it 'Should return standard greeting', (done) ->
    request
      .get '/'
      .end (err, res) ->
        expect(err).to.equal null
        expect(res.text).to.equal 'Hello, world!'
        done()
describe 'Ping endpoint', ->
  it 'Should respond to standard ping', (done) ->
    request
      .get '/ping'
      .end (err, res) ->
        expect(err).to.equal null
        expect(res.text).to.equal 'Pong'
        done()
