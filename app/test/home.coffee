supertest = require 'supertest'
expect = require('chai').expect

request = supertest 'http://ApiServer'
describe 'home', ->
  it 'should return standard greeting', (done) ->
    request
      .get '/'
      .end (err, res) ->
        expect(err).to.equal null
        expect(res.text).to.equal 'Hello, world!'
        done()
