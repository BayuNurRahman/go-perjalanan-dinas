package main

import (
	"fmt"
	"log"

	"go-perjalanan-dinas/config"
)

// SyncUserRoles fixes inconsistencies where Role field doesn't match role_id
// This should be run once to ensure all users have consistent role assignments
func SyncUserRoles() error {
	config.LoadEnv()
	config.ConnectDatabase()

	// Check inconsistencies first
	var inconsistencies []map[string]interface{}
	config.DB.Raw(`
		SELECT u.id, u.email, u.role, u.role_id, r.name as actual_role
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.role != r.name OR (u.role_id IS NOT NULL AND u.role IS NULL)
	`).Scan(&inconsistencies)

	if len(inconsistencies) == 0 {
		fmt.Println("✓ No role inconsistencies found")
		return nil
	}

	fmt.Printf("Found %d users with role inconsistencies:\n", len(inconsistencies))
	for _, user := range inconsistencies {
		fmt.Printf("  - ID: %v, Email: %v, Current Role: %v → Should be: %v\n",
			user["id"], user["email"], user["role"], user["actual_role"])
	}

	// Fix the inconsistencies
	result := config.DB.Exec(`
		UPDATE users u
		SET role = r.name
		FROM roles r
		WHERE u.role_id = r.id
		AND u.role != r.name
	`)

	if result.Error != nil {
		return fmt.Errorf("failed to sync roles: %w", result.Error)
	}

	fmt.Printf("✓ Successfully synced %d users\n", result.RowsAffected)
	return nil
}

// Verify all users have consistent role assignments
func VerifyRoles() error {
	var users []map[string]interface{}
	result := config.DB.Raw(`
		SELECT u.id, u.email, u.role, u.role_id, r.name as actual_role
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		ORDER BY u.id
	`).Scan(&users)

	if result.Error != nil {
		return result.Error
	}

	fmt.Println("\nAll users role status:")
	fmt.Println("ID | Email | Current Role | RoleID | Actual Role | Status")
	fmt.Println(string([]rune("---+-------+--------------+--------+-------------+--------")))

	var mismatches int
	for _, user := range users {
		status := "✓ OK"
		if user["role"] != user["actual_role"] {
			status = "✗ MISMATCH"
			mismatches++
		}
		fmt.Printf("%v | %v | %v | %v | %v | %v\n",
			user["id"], user["email"], user["role"], user["role_id"], user["actual_role"], status)
	}

	if mismatches == 0 {
		fmt.Println("\n✓ All users have consistent role assignments")
	} else {
		fmt.Printf("\n✗ Found %d mismatches. Run SyncUserRoles() to fix.\n", mismatches)
	}

	return nil
}

func main() {
	if config.DB == nil {
		fmt.Println("Connecting to database...")
		config.LoadEnv()
		config.ConnectDatabase()
	}

	if config.DB == nil {
		log.Fatal("Failed to connect to database")
	}

	// First verify current state
	if err := VerifyRoles(); err != nil {
		log.Fatalf("Verification failed: %v\n", err)
	}

	// Then sync if needed
	if err := SyncUserRoles(); err != nil {
		log.Fatalf("Sync failed: %v\n", err)
	}

	// Verify again after sync
	if err := VerifyRoles(); err != nil {
		log.Fatalf("Verification after sync failed: %v\n", err)
	}
}
