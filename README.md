# ТРЕБОВАНИЕ К ПАРОЛЮ
## ДЛИНА ПАРОЛЯ ДОЛЖНА СОСТАВЛЯТЬ 8 СИМВОЛОВ И ВКЛЮЧАТЬ ЦИФРОВЫЕ, ПРОСТЫЕ СИМВОЛЫ ИЛИ ЛАТИНСКИЕ БУКВЫ
## ** ПОД ПРОСТЫМ Я ИМЕЮ ВВИДУ ЭТИ СИМВОЛЫ (!@#$%^&*()[]{};:'"<,.>?/\+_-=)


This is a web Platonus created using Go, HTML, and SQLite. The Platonus allows users to create posts and comments, associate categories with posts, and like/dislike posts and comments. The project also includes a filter mechanism that allows users to filter posts by categories, created posts, and liked posts.

## Installation

1. Build the Docker image:
```bash
docker build -t forum .
```

2. Run the Docker container:
```bash
docker run --name new-forum --rm -p8080:8080 forum
```

3. Open the forum in your web browser:
```
http://localhost:8080
```

## Usage

### Registration

To register as a new user on the Platonus, enter your email, username, and password. If the email or name is already taken, an error response will be returned. 
There is also a password requirement (read about password requirement in the top of the topic).

### Login

To access the forum, login using your credentials. Only registered users can create posts and comments, like/dislike posts and comments, and filter posts.

### Posts and Comments

Registered users can create posts and comments, and associate categories with posts. The posts and comments are visible to all users, including non-registered users.

### Likes and Dislikes

Only registered users can like or dislike posts and comments. The number of likes and dislikes is visible to all users.

### Filter

Users can filter posts by categories, created posts, and liked posts. Non-registered users can only see posts and comments.


- `database/sql` - database/sql provides a generic interface around SQL (or SQL-like) databases.
- `github.com/mattn/go-sqlite3` - sqlite3 driver for Go using database/sql.
- `golang.org/x/crypto/bcrypt` - package bcrypt implements Provos and Mazières's bcrypt adaptive hashing algorithm.
- `github.com/google/uuid` - the uuid package generates and inspects UUIDs based on RFC 4122 and DCE 1.1: Authentication and Security Services.


