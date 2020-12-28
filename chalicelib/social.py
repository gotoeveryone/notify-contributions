import os

from requests_oauthlib import OAuth1Session


class Twitter:
    def __init__(self):
        self.twitter = OAuth1Session(
            os.getenv('CONSUMER_KEY'),
            os.getenv('CONSUMER_SECRET'),
            os.getenv('ACCESS_TOKEN'),
            os.getenv('ACCESS_TOKEN_SECRET'),
        )

    def post(self, message: str):
        url = 'https://api.twitter.com/1.1/statuses/update.json'
        res = self.twitter.post(url, params={
            'status': message,
        })

        return {
            'code': res.status_code,
            'body': res.json(),
        }
