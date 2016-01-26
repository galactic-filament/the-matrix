/// <reference path="../typings/tsd.d.ts" />
import * as supertest from "supertest";
import * as test from "tape";

let request = supertest("http://ApiServer");

interface PostCallback {
  (id: number): void;
}
let createPost = (t: test.Test, cb: PostCallback) => {
  request
    .post("/posts")
    .send({ body: "Hello, world!" })
    .end((err: Error, res: supertest.Response) => {
      t.equal(null, err);
      t.equal(200, res.status);
      t.equal("number", typeof res.body.id);
      cb(res.body.id);
    });
};

test("Homepage Should return standard greeting", (t: test.Test) => {
  request
    .get("/")
    .end((err: Error, res: supertest.Response) => {
      t.equal(null, err);
      t.equal(200, res.status);
      t.equal("Hello, world!", res.text);
      t.end();
    });
});
test("Ping endpoint Should respond to standard ping", (t: test.Test) => {
  request
    .get("/ping")
    .end((err: Error, res: supertest.Response) => {
      t.equal(null, err);
      t.equal(200, res.status);
      t.equal("Pong", res.text);
      t.end();
    });
});
test("Json reflection Should return identical Json in response as provided by request", (t: test.Test) => {
  let body = { greeting: "Hello, world!" };
  request
    .post("/reflection")
    .send(body)
    .end((err: Error, res: supertest.Response) => {
      t.equal(null, err);
      t.equal(200, res.status);
      t.equal(body.greeting, res.body.greeting);
      t.end();
    });
});
test("Post creation endpoint Should return the new post's id", (t: test.Test) => {
  createPost(t, (id: number) => {
    t.end();
  });
});
test("Post endpoint Should return a post", (t: test.Test) => {
  createPost(t, (id: number) => {
    let url = "/post/" + id;
    request
      .get(url)
      .end(function getPostEnd(err: Error, res: supertest.Response) {
        t.equal(null, err, `GET ${url} err was not null`);
        t.equal(200, res.status, `GET ${url} res.status was not 200`);
        t.end();
      });
  });
});
test("Post endpoint Should delete a post", (t: test.Test) => {
  createPost(t, (id: number) => {
    let url = "/post/" + id;
    request
      .delete(url)
      .end(function deletePostEnd(err: Error, res: supertest.Response) {
        t.equal(null, err, `DELETE ${url} err was not null`);
        t.equal(200, res.status, `DELETE ${url} res.status was not 200`);
        t.end();
      });
  });
});
