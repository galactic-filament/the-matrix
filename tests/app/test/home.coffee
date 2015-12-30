supertest = require 'supertest'
expect = require('chai').expect

request = supertest 'http://ApiServer'
describe 'Homepage', ->
  it 'Should return standard greeting', (done) ->
    request
      .get '/'
      .end (err, res) ->
        expect(err).to.equal null
        expect(res.status).to.equal 200
        expect(res.text).to.equal 'Hello, world!'
        done()
describe 'Ping endpoint', ->
  it 'Should respond to standard ping', (done) ->
    request
      .get '/ping'
      .end (err, res) ->
        expect(err).to.equal null
        expect(res.status).to.equal 200
        expect(res.text).to.equal 'Pong'
        done()
describe 'Json reflection', ->
  it(
    'Should return identical Json in response as provided by request'
    (done) ->
      body = { greeting: 'Hello, world!' }
      request
        .post '/reflection'
        .send body
        .end (err, res) ->
          expect(err).to.equal null
          expect(res.status).to.equal 200
          expect(res.body.greeting).to.equal body.greeting
          done()
  )
describe 'Post creation endpoint', ->
  it "Should return the new post's id", (done) ->
    request
      .post '/posts'
      .send { body: 'Hello, world!' }
      .end (err, res) ->
        expect(err).to.equal null
        expect(res.status).to.equal 200
        expect(typeof res.body.id).to.equal 'number'
        done()
