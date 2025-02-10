package main

const (
	CreateUserTable = `CREATE TABLE IF NOT EXISTS Users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	CreateDecisionsTable = `CREATE TABLE IF NOT EXISTS Decisions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    actor_id INT NOT NULL,
    recipient_id INT NOT NULL,
    liked BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_Decisions (actor_id, recipient_id),
    FOREIGN KEY (actor_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (recipient_id) REFERENCES Users(id) ON DELETE CASCADE
);`

	AddDummyUserData = `INSERT INTO Users (username, first_name, last_name) VALUES
                                                        ('user1', 'John', 'Doe'),
                                                        ('user2', 'Jane', 'Smith'),
                                                        ('user3', 'Alice', 'Johnson'),
                                                        ('user4', 'Bob', 'Brown'),
                                                        ('user5', 'Charlie', 'Davis'),
                                                        ('user6', 'David', 'Miller'),
                                                        ('user7', 'Emma', 'Wilson'),
                                                        ('user8', 'Frank', 'Moore'),
                                                        ('user9', 'Grace', 'Taylor'),
                                                        ('user10', 'Henry', 'Anderson');`

	AddDummyDecisionData = `INSERT INTO Decisions (actor_id, recipient_id, liked) VALUES
                                                        (1, 2, TRUE), (1, 3, TRUE), (1, 4, FALSE),
                                                        (2, 5, TRUE), (2, 6, FALSE), (2, 7, TRUE),
                                                        (3, 8, FALSE), (3, 9, TRUE), (3, 10, FALSE),
                                                        (4, 1, TRUE), (4, 2, FALSE), (4, 3, TRUE),
                                                        (5, 6, FALSE), (5, 7, TRUE), (5, 8, FALSE),
                                                        (6, 9, TRUE), (6, 10, FALSE), (6, 1, TRUE),
                                                        (1, 5, FALSE), (7, 3, TRUE), (7, 4, FALSE),
                                                        (8, 5, TRUE), (8, 6, FALSE), (8, 7, TRUE),
                                                        (9, 8, FALSE), (9, 9, TRUE), (9, 10, FALSE),
                                                        (10, 1, TRUE), (10, 2, FALSE), (10, 3, TRUE);`
)
