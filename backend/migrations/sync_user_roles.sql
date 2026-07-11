-- Migration: Sync User Role field with actual role_id
-- Purpose: Fix inconsistencies where denormalized Role field doesn't match the role_id foreign key
-- Example: User with role_id=3 (ADMIN_IT) but Role='EMPLOYEE' will be corrected to Role='ADMIN_IT'

UPDATE users u
SET role = r.name
FROM roles r
WHERE u.role_id = r.id
AND u.role != r.name;

-- Verify the sync by checking if any mismatches remain
-- SELECT u.id, u.email, u.role, u.role_id, r.name as actual_role
-- FROM users u
-- LEFT JOIN roles r ON u.role_id = r.id
-- WHERE u.role != r.name OR (u.role_id IS NOT NULL AND u.role IS NULL);
