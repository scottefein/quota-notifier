# quota-notifier

1) Copy app.yaml.sample to app.yaml
2) Sign up for a Mailgun account
3) Set the environmental variables as follows

  PROJECT_ID: "<your project id>"
  THRESHOLD: .8
  MG_API_KEY: <Mailgun Private API Key>
  MG_DOMAIN: <The domain, either the sandbox one or one you've verified with mailgun>
  MG_FROM_EMAIL: <name@thedomainfromabove>
  MG_TO_EMAIL: <The email you want to send to, if you haven't verified the domain you'll need it to be an authorized recipient>
  MG_PUBLIC_API_KEY: <Mailgun Public API Key>
  MG_URL: "https://api.mailgun.net/v3"
  
4) Modify the main.go script "regions" variable to include all the regions that the project is running in
5) gcloud app deploy
