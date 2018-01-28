use IS4B;

CREATE TABLE MEETING_ROOM_RESERVATION (
	id              int NOT NULL auto_increment,
	userId          int NOT NULL,
	buildingId      int NOT NULL,
    startTime       datetime NOT NULL,
    endTime         datetime NOT NULL,
    timeReserved    timestamp default CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);
