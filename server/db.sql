CREATE TABLE "user" (
  id       VARCHAR(15) PRIMARY KEY,
  password BYTEA       NOT NULL
);

CREATE TABLE user_board (
  id      SERIAL      PRIMARY KEY,
  userID  VARCHAR(15) NOT NULL    REFERENCES "user",
  boardID SERIAL      NOT NULL    REFERENCES board,
  isAdmin BOOLEAN     NOT NULL
)

CREATE TABLE board (
  id   SERIAL      PRIMARY KEY,
  name VARCHAR(30) NOT NULL
);

CREATE TABLE column (
  id      SERIAL   PRIMARY KEY,
  tableID SERIAL   NOT NULL    REFERENCES board,
  order   SMALLINT NOT NULL
)

CREATE TABLE task (
  id          SERIAL      PRIMARY KEY,
  columnID    SERIAL      NOT NULL    REFERENCES column,
  title       VARCHAR(50) NOT NULL,
  description TEXT
)

CREATE TABLE subtask (
  id     SERIAL      PRIMARY KEY,
  taskID SERIAL      NOT NULL    REFERENCES task,
  title  VARCHAR(50) NOT NULL,
  isDone BOOLEAN     NOT NULL
)
