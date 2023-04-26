resource "kevel_channel_site_map" "example" {
  channel_id = kevel_channel.example.id
  site_id    = kevel_site.example.id
  priority   = 10
}
