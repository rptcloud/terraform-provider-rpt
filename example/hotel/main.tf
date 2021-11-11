terraform {
  required_providers {
    rpt = {
      version = "0.54"
      source  = "local/cmarkulin/rpt"
    }
  }
}

variable "hotel_name" {
  type    = string
  default = "RPT Hotel Pitt"
}

/*
resource "rpt_hotel" "test_hotel"{
    name = "Updated in Terraform"
    state = "PA"
    city = "Pittsburgh"
    rating = "5.0"
    description = "This hotel was created in Terraform through a provider for this web app."
    photo = "https://dynamic-media-cdn.tripadvisor.com/media/photo-o/18/44/44/f0/exterior.jpg?w=900&h=-1&s=1"
}
*/


data "rpt_hotels" "all" {
    hotel_num = "3"
}

# Returns all hotels
output "hotel" {
  value = data.rpt_hotels.all
}
