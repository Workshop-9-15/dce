resource "aws_organizations_policy" "dce_security_scp" {
  count = var.enable_scp ? 1 : 0
  
  name        = "DCE-Security-SCP-${var.namespace}"
  description = "DCE Service Control Policy to prevent IAM privilege escalation and restrict to supported services"
  type        = "SERVICE_CONTROL_POLICY"
  content     = templatefile("${path.module}/fixtures/policies/scp_policy.json", {
    admin_role_name       = "AdminRole"
    principal_role_name   = local.principal_role_name
    principal_policy_name = local.principal_policy_name
  })
  
  tags = var.global_tags
}

resource "aws_organizations_policy_attachment" "dce_security_scp_attachment" {
  count = var.enable_scp && length(var.scp_target_ids) > 0 ? length(var.scp_target_ids) : 0
  
  policy_id = aws_organizations_policy.dce_security_scp[0].id
  target_id = var.scp_target_ids[count.index]
}
