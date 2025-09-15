# Testing Service Control Policy Implementation

This document provides step-by-step instructions for testing the DCE Service Control Policy (SCP) implementation to verify it prevents IAM privilege escalation.

## Prerequisites

- AWS Organizations set up with DCE deployed
- SCP enabled in Terraform configuration (`enable_scp = true`)
- Test account within the organization where SCP is attached
- Access to both AdminRole and principal user credentials

## Test Scenarios

### Test 1: Verify IAM Privilege Escalation Prevention

This test verifies that principal users cannot create new IAM roles with administrative access.

#### Setup
1. Deploy DCE with SCP enabled:
   ```hcl
   enable_scp = true
   scp_target_ids = ["123456789012"]  # Your test account ID
   ```

2. Apply Terraform configuration:
   ```bash
   terraform apply
   ```

#### Test Steps

**Step 1: Test as Principal User (Should Fail)**
1. Assume the principal role in the test account
2. Attempt to create a new IAM role with administrative access:
   ```bash
   aws iam create-role \
     --role-name TestAdminRole \
     --assume-role-policy-document '{
       "Version": "2012-10-17",
       "Statement": [{
         "Effect": "Allow",
         "Principal": {"Service": "ec2.amazonaws.com"},
         "Action": "sts:AssumeRole"
       }]
     }'
   ```
3. **Expected Result**: Command should fail with access denied error

**Step 2: Attempt to Attach Administrative Policy (Should Fail)**
1. Try to attach AWS managed administrative policy:
   ```bash
   aws iam attach-role-policy \
     --role-name TestAdminRole \
     --policy-arn arn:aws:iam::aws:policy/AdministratorAccess
   ```
2. **Expected Result**: Command should fail with access denied error

**Step 3: Test as AdminRole (Should Succeed)**
1. Assume the AdminRole in the test account
2. Repeat the same commands from Steps 1-2
3. **Expected Result**: Commands should succeed, demonstrating AdminRole can still manage IAM

### Test 2: Verify Service Restrictions

This test verifies that the SCP restricts access to unsupported AWS services.

#### Test Steps

**Step 1: Test Allowed Service (Should Succeed)**
1. As principal user, try to use an allowed service like S3:
   ```bash
   aws s3 ls
   ```
2. **Expected Result**: Command should succeed

**Step 2: Test Restricted Service (Should Fail)**
1. As principal user, try to use a restricted service like WorkSpaces:
   ```bash
   aws workspaces describe-workspaces
   ```
2. **Expected Result**: Command should fail with access denied error

### Test 3: Verify Principal Role/Policy Protection

This test verifies that the SCP protects DCE's own IAM resources.

#### Test Steps

**Step 1: Attempt to Modify Principal Role (Should Fail)**
1. As principal user, try to modify the principal role:
   ```bash
   aws iam put-role-policy \
     --role-name DCEPrincipal-namespace \
     --policy-name TestPolicy \
     --policy-document '{
       "Version": "2012-10-17",
       "Statement": [{
         "Effect": "Allow",
         "Action": "*",
         "Resource": "*"
       }]
     }'
   ```
2. **Expected Result**: Command should fail with access denied error

**Step 2: Attempt to Delete Principal Policy (Should Fail)**
1. Try to delete the principal policy:
   ```bash
   aws iam delete-policy \
     --policy-arn arn:aws:iam::ACCOUNT:policy/DCEPrincipalPolicy-namespace
   ```
2. **Expected Result**: Command should fail with access denied error

## Automated Testing Script

Create a test script to automate these validations:

```bash
#!/bin/bash
# test-scp.sh

set -e

ACCOUNT_ID="123456789012"
NAMESPACE="your-namespace"

echo "Testing SCP Implementation..."

# Test 1: IAM privilege escalation prevention
echo "Test 1: Attempting to create admin role as principal user..."
if aws iam create-role --role-name TestAdminRole --assume-role-policy-document file://trust-policy.json 2>/dev/null; then
    echo "❌ FAIL: Principal user was able to create IAM role"
    exit 1
else
    echo "✅ PASS: Principal user blocked from creating IAM role"
fi

# Test 2: Service restrictions
echo "Test 2: Testing service restrictions..."
if aws workspaces describe-workspaces 2>/dev/null; then
    echo "❌ FAIL: Principal user has access to restricted service"
    exit 1
else
    echo "✅ PASS: Principal user blocked from restricted service"
fi

# Test 3: Allowed service access
echo "Test 3: Testing allowed service access..."
if aws s3 ls 2>/dev/null; then
    echo "✅ PASS: Principal user can access allowed service"
else
    echo "❌ FAIL: Principal user blocked from allowed service"
    exit 1
fi

echo "All SCP tests passed! ✅"
```

## Troubleshooting

### Common Issues

1. **SCP not taking effect**: 
   - Verify SCP is attached to the correct organizational unit or account
   - Check that the account is within the organization structure

2. **Tests passing when they should fail**:
   - Confirm you're using principal user credentials, not AdminRole
   - Verify SCP policy syntax and variable substitution

3. **All operations failing**:
   - Check if SCP is too restrictive
   - Verify AdminRole can still perform operations

### Validation Commands

```bash
# Check SCP attachment
aws organizations list-policies-for-target --target-id 123456789012 --filter SERVICE_CONTROL_POLICY

# View SCP content
aws organizations describe-policy --policy-id p-xxxxxxxxxx

# Verify current role
aws sts get-caller-identity
```

## Security Considerations

- Always test in a non-production environment first
- Monitor CloudTrail logs for denied actions during testing
- Ensure AdminRole maintains necessary permissions for DCE operations
- Document any custom modifications to the SCP policy

## Expected Test Results Summary

| Test Scenario | Principal User | AdminRole |
|---------------|----------------|-----------|
| Create IAM Role | ❌ Denied | ✅ Allowed |
| Attach Admin Policy | ❌ Denied | ✅ Allowed |
| Access S3 | ✅ Allowed | ✅ Allowed |
| Access WorkSpaces | ❌ Denied | ✅ Allowed |
| Modify Principal Role | ❌ Denied | ✅ Allowed |

If all tests show expected results, the SCP implementation is working correctly and the IAM privilege escalation vulnerability has been successfully mitigated.
