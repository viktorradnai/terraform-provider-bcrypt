terraform {
    required_providers {
    bcrypt = {
      source  = "viktest.com/viktorradnai/bcrypt"
      version = ">= 0.0.0"
    }
  }
}

provider "bcrypt" {
}

resource "bcrypt_hash" "test" {
  cleartext = "hunter2"
}
