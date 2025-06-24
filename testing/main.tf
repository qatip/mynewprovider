terraform {
  required_providers {
    mynewprovider = {
      source  = "custom/mynewprovider"
      version = "0.1.0"
    }
  }
}

provider "mynewprovider" {}

resource "mynewprovider_task" "example" {
  title       = "From custom provider"
  description = "Created using mynewprovider"
  completed   = false
}
