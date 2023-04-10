resource "kevel_channel" "example" {
  title = "My Channel"
  ad_types = [
    kevel_ad_type.example.id
  ]
}
