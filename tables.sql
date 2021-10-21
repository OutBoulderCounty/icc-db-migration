DROP TABLE IF EXISTS options;
DROP TABLE IF EXISTS elements;
DROP TABLE IF EXISTS forms;

CREATE TABLE forms (
  id int NOT NULL AUTO_INCREMENT,
  name text,
  required boolean,
  live boolean,
  PRIMARY KEY (id)
);

CREATE TABLE elements (
  id int NOT NULL AUTO_INCREMENT,
  formID int NOT NULL,
  label text NOT NULL,
  type text NOT NULL,
  position int NOT NULL, -- index
  required boolean NOT NULL,
  priority int,
  search boolean NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE options (
  id int NOT NULL AUTO_INCREMENT,
  elementID int NOT NULL,
  name text NOT NULL,
  position int NOT NULL, -- index
  PRIMARY KEY (id)
);
