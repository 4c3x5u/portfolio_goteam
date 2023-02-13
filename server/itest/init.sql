-- Initialise the database for server use.

CREATE SCHEMA app;

CREATE TABLE app."user" (
  id       VARCHAR(15) PRIMARY KEY,
  password BYTEA       NOT NULL
);

CREATE TABLE app.board (
  id   INTEGER     PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  name VARCHAR(30) NOT NULL
);

CREATE TABLE app.user_board (
  id      INTEGER     PRIMARY KEY GENERATED ALWAYS      AS IDENTITY,
  userID  VARCHAR(15) NOT NULL    REFERENCES app."user",
  boardID SERIAL      NOT NULL    REFERENCES app.board,
  isAdmin BOOLEAN     NOT NULL
);

CREATE TABLE app."column" (
  id      INTEGER  PRIMARY KEY GENERATED ALWAYS     AS IDENTITY,
  tableID SERIAL   NOT NULL    REFERENCES app.board,
  "order" SMALLINT NOT NULL
);

CREATE TABLE app.task (
  id          INTEGER     PRIMARY KEY GENERATED ALWAYS        AS IDENTITY,
  columnID    SERIAL      NOT NULL    REFERENCES app."column",
  title       VARCHAR(50) NOT NULL,
  description TEXT
);

CREATE TABLE app.subtask (
  id     INTEGER     PRIMARY KEY GENERATED ALWAYS    AS IDENTITY,
  taskID SERIAL      NOT NULL    REFERENCES app.task,
  title  VARCHAR(50) NOT NULL,
  isDone BOOLEAN     NOT NULL
);
