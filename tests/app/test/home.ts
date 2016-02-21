/// <reference path="../typings/tsd.d.ts" />
import * as supertest from "supertest";
import * as test from "tape";

let request = supertest("http://ApiServer");

interface PostCallback {
  (id: number): void;
}
let createPost = (t: test.Test, cb: PostCallback) => {
  let url = "/posts";
  request
    .post(url)
    .send({ body: "Hello, world!" })
    .end((err: Error, res: supertest.Response) => {
      t.equal(err, null, `POST ${url} err was not null`);
      t.equal(res.status, 200, `POST ${url} res.status was not 200`);
      t.equal(typeof res.body.id, "number", `POST ${url} body.id was not a number`);
      cb(res.body.id);
    });
};

test("Homepage Should return standard greeting", (t: test.Test) => {
  let url = "/";
  request
    .get(url)
    .end((err: Error, res: supertest.Response) => {
      t.equal(err, null, `GET ${url} err was not null`);
      t.equal(res.status, 200, `GET ${url} res.status was not 200`);
      t.equal(res.text, "Hello, world!", `GET  ${url} response body was not Hello, world!`);
      t.end();
    });
});
test("Ping endpoint Should respond to standard ping", (t: test.Test) => {
  let url = "/ping";
  request
    .get(url)
    .end((err: Error, res: supertest.Response) => {
      t.equal(err, null, `GET ${url} err was not null`);
      t.equal(res.status, 200, `GET ${url} res.status was not 200`);
      t.equal(res.text, "Pong", `GET ${url} response body was not Pong`);
      t.end();
    });
});
test("Json reflection Should return identical Json in response as provided by request", (t: test.Test) => {
  let url = "/reflection";
  let body = { greeting: "Hello, world!" };
  request
    .post(url)
    .send(body)
    .end((err: Error, res: supertest.Response) => {
      t.equal(err, null, `POST ${url} err was not null`);
      t.equal(res.status, 200, `POST ${url} res.status was not 200`);
      t.equal(res.body.greeting, body.greeting, `POST ${url} greeting did not match`);
      t.end();
    });
});
test("Post creation endpoint Should return the new post's id", (t: test.Test) => {
  createPost(t, (id: number) => t.end());
});
test("Post endpoint Should return a post", (t: test.Test) => {
  createPost(t, (id: number) => {
    let url = "/post/" + id;
    request
      .get(url)
      .end(function getPostEnd(err: Error, res: supertest.Response) {
        t.equal(err, null, `GET ${url} err was not null`);
        t.equal(res.status, 200, `GET ${url} res.status was not 200`);
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
        t.equal(err, null, `DELETE ${url} err was not null`);
        t.equal(res.status, 200, `DELETE ${url} res.status was not 200`);
        t.end();
      });
  });
});
test("Post endpoint Should update a post", (t: test.Test) => {
  createPost(t, (id: number) => {
    let url = "/post/" + id;
    let body = { body: "Jello, world!" };
    request
      .put(url)
      .send(body)
      .end(function updatePostEnd(err: Error, res: supertest.Response) {
        t.equal(err, null, `PUT ${url} err was not null`);
        t.equal(res.status, 200, `PUT ${url} res.status was not 200`);
        t.equal(res.body.body, body.body, `PUT ${url} request and response bodies did not match`);
        t.end();
      });
  });
});
