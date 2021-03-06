# TODO List

- [ ] Provide hosted and downloadable version
- [ ] Have a demo page that shows randomly generated fake data
    - Maybe it just loops a set of pre-made conversations instead of full generation

- [ ] Maybe if conversation is just two people, use 2 series of data for line chart, one for each user
- [ ] Maybe find way to embed static dir in go files
- [ ] Maybe have ratelimit for IPs (flag/option to decide the limit)

## Done List

- [x] Have home page different from view page
- [x] Home page should be the "upload" page
- [x] Switch from form to upload and stats api requests
- [x] Don't store data in memory unless requested
- [x] Upload should have options such as "cache for x amount of time
- [x] Dockerize
- [x] Github actions
- [x] Readme badges
- [x] Home page should have upload/stat options, gh link, credits, and instructions
- [x] Any server-side errors should be prettily shown to the client instead of just server logging / alert()ing
- [x] View page should be suitable for printing to pdf or whatever
- [x] Should make it clear on the hosted version that data is transferred but not stored
- [x] If user downloads binary, make it clear that it does not request any resources from internet (user must download static dir tho)

### JSON WebSocket API

File upload progress response:

```json
{
  "progress": int (0-100)
}
```

Send file count and then raw json file bytes, receive json:

```json
{
  "error": string ("" if no error),
  "conversation_title": string,
  "messages_per_month": highcharts-object,
  "messages_per_user": array of highcharts-object,
  "messages_per_weekday": highcharts-object
}
```
