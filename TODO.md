# TODO List

- [ ] Provide hosted and downloadable version
- [ ] Should make it clear on the hosted version that data is transferred but not stored (unless asked)
- [ ] If user downloads binary, make it clear that it does not request any resources from internet (user must download static dir tho)
	- Maybe find way to embed static dir in go files (is this now a part of go?)
- [ ] Any server-side errors should be prettily shown to the client instead of just server logging / alert()ing
- [ ] Home page should have upload/stat options, gh link, creits, and instructions
- [ ] View page should be suitable for printing to pdf or whatever
- [ ] If caching on upload, provide unique link, but make it clear that anyone can access data if they have link, e.g. /view?id=xyz
- [ ] Make every stat optional, have grid or something on upload page wth tickboxes for each
- [ ] Display "this stat takes ~Xms per 1000 messages" so users know how long wait will be
    - Check the timings aren't drastically different on hetzner server comapred to pc
- [ ] Binary should have flag or option that allows you to set request rate for api requests
- [ ] Binary should also have flag/options that allow diff features like caching data
- [ ] Have a demo page that shows randomly generated fake data
    - Maybe it just loops a set of pre-made conversations instead of full generation

## Done List

- [x] Have home page different from view page
- [x] Home page should be the "upload" page
- [x] Switch from form to upload and stats api requests
- [x] Don't store data in memory unless requested
- [x] Upload should have options such as "cache for x amount of time
- [x] Dockerize
- [x] Github actions
- [x] Readme badges

### New JSON API

- Error messages should be user-friendly if possible

#### POST `/api/upload?cache=true` 

- Must contain multipart form body
- Strict request rate limit
- If cache is false, the `id` is single use

```json
{
  "error": string ("" if no error),
  "id": string
}
```

#### GET `/api/stats?id=xxx`

```json
{
  "error": string ("" if no error),
  "conversation_title": string,
  "messages_per_month": highcharts-object,
  "messages_per_user": array of highcharts-object,
  "messages_per_weekday": highcharts-object
}
```
