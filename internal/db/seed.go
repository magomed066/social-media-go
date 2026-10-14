package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"social-media-app/internal/features/comments"
	"social-media-app/internal/features/posts"
	"social-media-app/internal/features/users"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Seed adds 100 development users, 200 posts, and 500 comments.
// Existing seed records are reused without modifying their contents.
// All seeded users have the hashed development password "seed-password".
func Seed(db *sql.DB) error {
	password, err := bcrypt.GenerateFromPassword([]byte("seed-password"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash seed password: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin seed transaction: %w", err)
	}
	defer tx.Rollback()

	// Serialize seed runs because posts and comments have no unique seed key.
	if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(736333100)`); err != nil {
		return fmt.Errorf("lock seed transaction: %w", err)
	}

	seedUsers := generateUsers(100)
	for _, user := range seedUsers {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO users (first_name, last_name, username, email, password)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT DO NOTHING
		`, user.Firstname, user.Lastname, user.Username, user.Email, password)

		if err != nil {
			return fmt.Errorf("seed user %q: %w", user.Username, err)
		}

		if err := tx.QueryRowContext(ctx, `
			SELECT id FROM users WHERE username = $1 AND email = $2
		`, user.Username, user.Email).Scan(&user.ID); err != nil {
			return fmt.Errorf("find seed user %q (username and email must match): %w", user.Username, err)
		}
	}

	seedPosts := generatePosts(200, seedUsers)
	for _, post := range seedPosts {
		err := tx.QueryRowContext(ctx, `
			SELECT id FROM posts WHERE user_id = $1 AND title = $2 ORDER BY id LIMIT 1
		`, post.UserID, post.Title).Scan(&post.ID)

		if errors.Is(err, sql.ErrNoRows) {
			err = tx.QueryRowContext(ctx, `
				INSERT INTO posts (user_id, title, content, tags)
				VALUES ($1, $2, $3, $4) RETURNING id
			`, post.UserID, post.Title, post.Content, pq.Array(post.Tags)).Scan(&post.ID)
		}

		if err != nil {
			return fmt.Errorf("seed post %q: %w", post.Title, err)
		}
	}

	for _, comment := range generateComments(500, seedUsers, seedPosts) {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO comments (user_id, post_id, content)
			SELECT $1, $2, $3
			WHERE NOT EXISTS (
				SELECT 1 FROM comments WHERE user_id = $1 AND post_id = $2 AND content = $3
			)
		`, comment.UserID, comment.PostID, comment.Content)

		if err != nil {
			return fmt.Errorf("seed comment on post %d: %w", comment.PostID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit seed transaction: %w", err)
	}

	return nil
}

func generateUsers(num int) []*users.User {
	data := make([]*users.User, num)

	firstNames := []string{"Alex", "Sam", "Jordan", "Taylor", "Morgan", "Casey", "Jamie", "Riley", "Avery", "Cameron"}
	lastNames := []string{"Smith", "Johnson", "Brown", "Davis", "Wilson", "Moore", "Taylor", "Anderson", "Thomas", "Jackson"}

	for i := 0; i < num; i++ {
		data[i] = &users.User{
			Firstname: firstNames[i%len(firstNames)],
			Lastname:  lastNames[(i/len(firstNames))%len(lastNames)],
			Username:  fmt.Sprintf("seed_user_%03d", i+1),
			Email:     fmt.Sprintf("seed_user_%03d@example.com", i+1),
		}
	}

	return data
}

func generatePosts(num int, seedUsers []*users.User) []*posts.Post {
	data := make([]*posts.Post, 0, num)

	if len(seedUsers) == 0 {
		return data
	}

	topics := []string{"Learning Go", "Building an API", "Working with PostgreSQL", "Writing tests", "Sharing progress"}

	for i := 0; i < num; i++ {
		topic := topics[i%len(topics)]
		data = append(data, &posts.Post{
			UserID:  seedUsers[i%len(seedUsers)].ID,
			Title:   fmt.Sprintf("[Seed %03d] %s", i+1, topic),
			Content: fmt.Sprintf("Today I am exploring %s. Here are my notes and lessons from building a social media application.", topic),
			Tags:    []string{"seed", "development"},
		})
	}

	return data
}

func generateComments(num int, seedUsers []*users.User, seedPosts []*posts.Post) []*comments.Comment {
	data := make([]*comments.Comment, 0, num)

	if len(seedUsers) == 0 || len(seedPosts) == 0 {
		return data
	}

	responses := []string{"Thanks for sharing your experience!", "This is a useful example.", "Looking forward to your next update.", "I am working on something similar.", "The explanation helped me understand this better."}

	for i := 0; i < num; i++ {
		data = append(data, &comments.Comment{
			UserID:  seedUsers[(i+1)%len(seedUsers)].ID,
			PostID:  seedPosts[i%len(seedPosts)].ID,
			Content: fmt.Sprintf("[Seed comment %03d] %s", i+1, responses[i%len(responses)]),
		})
	}

	return data
}
