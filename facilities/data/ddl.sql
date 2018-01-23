CREATE DATABASE IS4B;

use IS4B;

CREATE TABLE BUILDING (
	id int NOT NULL auto_increment,
    name varchar(128) NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE MEETING_ROOM (
	id int NOT NULL auto_increment,
	buildingId int NOT NULL,
    floor smallint NOT NULL,
    name varchar(128) NOT NULL,
    capacity mediumint,
    hasProjector bool,
    hasWhiteboard bool,
    hasConferenceLine bool,
    PRIMARY KEY(id)
);
