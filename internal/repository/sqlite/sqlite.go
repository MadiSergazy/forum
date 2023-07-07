package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zhayt/clean-arch-tmp-forum/config"
	"time"
)

const (
	userTable = `CREATE TABLE IF NOT EXISTS user(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			password TEXT UNIQUE,
			login TEXT UNIQUE,
			username TEXT UNIQUE,
			token TEXT DEFAULT NULL,
			tokenduration DATETIME DEFAULT NULL
			);`
	postTable = `CREATE TABLE IF NOT EXISTS post(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			author TEXT,
			title TEXT,
			description TEXT,
			like INTEGER DEFAULT 0,
			dislike INTEGER DEFAULT 0,
			date TEXT,
			FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE
			);`

	categoryTable = `CREATE TABLE IF NOT EXISTS category(
    	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    	title TEXT NOT NULL UNIQUE 
    );
		INSERT INTO category(title)
		VALUES ("IT"), ("Sport"), ("Education"), ("News"), ("Health");
`
	postCategoryTable = `CREATE TABLE IF NOT EXISTS post_category(
    		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    		post_id INTEGER NOT NULL,
    		category_id INTEGER NOT NULL,
    		FOREIGN KEY(post_id) REFERENCES post(id) ON DELETE CASCADE,
    		FOREIGN KEY(category_id) REFERENCES category(id) ON DELETE CASCADE
	);`

	commentTable = `CREATE TABLE IF NOT EXISTS comment(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			postId INTEGER,
			userId INTEGER,
			author TEXT,
			text TEXT,
			like INTEGER DEFAULT 0,
			dislike INTEGER DEFAULT 0,
			date TEXT,
			FOREIGN KEY(postId) REFERENCES post(id) ON DELETE CASCADE
			);`
	likeTable = `CREATE TABLE IF NOT EXISTS like(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			postId INTEGER,
			userId INTEGER,
			commentId INTEGER,
			active INTEGER DEFAULT 0,
			FOREIGN KEY(postId) REFERENCES post(id) ON DELETE CASCADE
		);`
	dislikeTable = `CREATE TABLE IF NOT EXISTS dislike(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			postId INTEGER,
			userId INTEGER,
			commentId INTEGER,
			active INTEGER DEFAULT 0,
			FOREIGN KEY(postId) REFERENCES post(id) ON DELETE CASCADE
		);`
	drop     = `DROP TABLE IF EXISTS schedule;`
	schedule = `
CREATE TABLE IF NOT EXISTS schedule (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    student_id INTEGER NOT NULL,
    class_name TEXT NOT NULL,
    day TEXT NOT NULL,
    time_start TEXT NOT NULL,
    time_end TEXT NOT NULL
);`
	materials = `CREATE TABLE IF NOT EXISTS materials (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		subtopic TEXT,
		content TEXT
	)`
)

func Dial(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("couldn't get db connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("couldn't ping db: %w", err)
	}
	return db, nil
}

func InnitDB(db *sql.DB) error {
	for _, table := range []string{userTable, postTable, commentTable, likeTable, dislikeTable, categoryTable, postCategoryTable, drop, schedule, materials} {
		if table == categoryTable {
			var count int
			err := db.QueryRow("SELECT COUNT(*) FROM category").Scan(&count)
			if err != nil {
				fmt.Printf(err.Error())
			}

			if count == 0 { // if category table is empty, insert data
				_, err = db.Exec(categoryTable)
				if err != nil {
					return fmt.Errorf("couldn't init db: error: %w table: %s", err, categoryTable)
				}
			}
			continue
		}

		_, err := db.Exec(table)
		if err != nil {
			return fmt.Errorf("couldn't init db: error: %w table: %s", err, table)
		}

	}
	sqlq := `-- Monday
INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'Math', 'Monday', '9:00', '10:30');

INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'Physics', 'Monday', '11:00', '12:30');

INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'Chemistry', 'Monday', '13:00', '14:30');

INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'History', 'Monday', '15:00', '16:30');

-- Tuesday
INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'English', 'Tuesday', '10:00', '11:30');

INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'Math', 'Tuesday', '12:00', '13:30');

INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'Biology', 'Tuesday', '14:00', '15:30');

-- Wednesday
INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'Art', 'Wednesday', '9:00', '10:30');

INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'Chemistry', 'Wednesday', '11:00', '12:30');

INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'Math', 'Wednesday', '13:00', '14:30');

-- Thursday
INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'History', 'Thursday', '10:00', '11:30');

INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'Physics', 'Thursday', '12:00', '13:30');

INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'English', 'Thursday', '14:00', '15:30');

INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'Math', 'Thursday', '16:00', '17:30');

-- Friday
INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'Biology', 'Friday', '9:00', '10:30');

INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'Math', 'Friday', '11:00', '12:30');

INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'Art', 'Friday', '13:00', '14:30');

INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'History', 'Friday', '15:00', '16:30');

INSERT INTO schedule (student_id, class_name, day, time_start, time_end)
VALUES (1, 'English', 'Friday', '17:00', '18:30');
`
	_, err = db.Exec(sqlq)
	if err != nil {
		return fmt.Errorf("couldn't init db: error: %w table: %s", err)
	}
	return nil
}
