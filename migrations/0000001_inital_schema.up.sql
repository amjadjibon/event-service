CREATE TABLE events
(
    id       INT AUTO_INCREMENT PRIMARY KEY,
    title    VARCHAR(255),
    start_at DATETIME,
    end_at   DATETIME,
    CONSTRAINT chk_event_start_end CHECK (start_at <= end_at)
);

CREATE TABLE workshops
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    event_id    INT,
    title       VARCHAR(255),
    description TEXT,
    start_at    DATETIME,
    end_at      DATETIME,
    CONSTRAINT chk_workshop_start_end CHECK (start_at <= end_at),
    CONSTRAINT fk_workshop_event FOREIGN KEY (event_id) REFERENCES events (id)
);


CREATE TABLE reservations
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(255),
    email       VARCHAR(255),
    workshop_id INT,
    CONSTRAINT fk_reservation_workshop FOREIGN KEY (workshop_id) REFERENCES workshops(id),
    CONSTRAINT unique_user_workshop UNIQUE KEY (email, workshop_id)
);