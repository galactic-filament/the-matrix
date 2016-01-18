/// <reference path="../typings/tsd.d.ts" />
import supertest = require("supertest");
import chai = require("chai");

let expect = chai.expect;
let request = supertest("http://ApiServer");

interface PostCallback {
  (id: number): void;
}
let createPost = (cb: PostCallback) => {
  request
    .post("/posts")
    .send({ body: "Hello, world!" })
    .end((err: Error, res: supertest.Response) => {
      expect(err).to.equal(null);
      expect(res.status).to.equal(200);
      expect(typeof res.body.id).to.equal("number");
      cb(res.body.id);
    });
};

describe("Homepage", () => {
    it("Should return standard greeting", (done: MochaDone) => {
      request
        .get("/")
        .end((err: Error, res: supertest.Response) => {
          expect(err).to.equal(null);
          expect(res.status).to.equal(200);
          expect(res.text).to.equal("Hello, world!");
          done();
        });
    });
});
describe("Ping endpoint", () => {
  it("Should respond to standard ping", (done: MochaDone) => {
    request
      .get("/ping")
      .end((err: Error, res: supertest.Response) => {
        expect(err).to.equal(null);
        expect(res.status).to.equal(200);
        expect(res.text).to.equal("Pong");
        done();
      });
  });
});
describe("Json reflection", () => {
  it("Should return identical Json in response as provided by request", (done: MochaDone) => {
      let body = { greeting: "Hello, world!" };
      request
        .post("/reflection")
        .send(body)
        .end((err: Error, res: supertest.Response) => {
          expect(err).to.equal(null);
          expect(res.status).to.equal(200);
          expect(res.body.greeting).to.equal(body.greeting);
          done();
        });
  });
});
describe("Post creation endpoint", () => {
  it("Should return the new post's id", (done: MochaDone) => {
    createPost((id: number) => {
      done();
    });
  });
});
describe("Post endpoint", () => {
  it("Should return a post", (done: MochaDone) => {
    createPost((id: number) => {
      request
        .get("/post/" + id)
        .end((err: Error, res: supertest.Response) => {
          expect(err).to.equal(null);
          expect(res.status).to.equal(200);
          done();
        });
    });
  });
  it("Should delete a post", (done: MochaDone) => {
    createPost((id: number) => {
      request
        .delete("/post/" + id)
        .end((err: Error, res: supertest.Response) => {
          expect(err).to.equal(null);
          expect(res.status).to.equal(200);
          done();
        });
    });
  });
});
