CREATE TABLE operations (
  id          string         NOT NULL,
  namespace   string         NOT NULL,
  msg_id      string         NOT NULL,
  data_id     string,
  optype      string         NOT NULL,
  opstatus    string         NOT NULL,
  recipient   string,
  plugin      string         NOT NULL,
  backend_id   string         NOT NULL,
  created     int64          NOT NULL,
  updated     int64,
  error       string         NOT NULL,
);

CREATE UNIQUE INDEX operations_primary ON operations(id);
CREATE INDEX operations_created ON operations(created);