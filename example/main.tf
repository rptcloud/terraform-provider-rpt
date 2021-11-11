terraform {
  required_providers {
    rpt = {
      source  = "local/cmarkulin/rpt"
      version = "0.54"
    }
  }
}

provider "rpt" {
}

module "hotel" {
  source = "./hotel"

}

resource "rpt_hotel" "sample" {
    name = "RPT Hotel Pitt"
    city = "Pittsburgh"
    state = "PA"
    rating = "3.9"
    photo = "https://upload.wikimedia.org/wikipedia/commons/thumb/3/3b/Vista_International_Hotel_Pittsburgh.jpg/480px-Vista_International_Hotel_Pittsburgh.jpg"
    description = "RPT Hotel at Pittsburgh"
}

output "hotel_data" {
  value = module.hotel.hotel
}
