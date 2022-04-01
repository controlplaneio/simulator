variable "default_tags" {
  description = "Default tags for all resources"
  type        = map(string)
}

variable "join_bucket_arn" {
  type = string
  description = "ARN for join S3 bucket" 
}
