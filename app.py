from datetime import date, datetime, timedelta
import os

from bs4 import BeautifulSoup
from chalice import Chalice, Cron
import requests

from chalicelib.social import TwitterClient

app = Chalice(app_name='notify-github-contributions')
app.debug = True if os.getenv('DEBUG', False) else False


class Contribution:
    def __init__(self, user_name: str, base_date: date, yesterday: str, today: str):
        self.user_name = user_name
        self.base_date = base_date
        self.yesterday = int(yesterday)
        self.today = int(today)

    @property
    def base_date_string(self) -> str:
        return self.base_date.strftime('%Y-%m-%d')

    def difference_to_yesterday(self) -> int:
        return self.today - self.yesterday


def get_message(contribution: Contribution) -> str:
    return '\n'.join([
        f'Contribute count to GitHub on {contribution.base_date_string} by {contribution.user_name}: {contribution.today}.',
        f'Difference from yesterday: {contribution.difference_to_yesterday()}.',
    ])


def get_contributions(user_name: str, base_date: date) -> Contribution:
    res = requests.get(f'https://github.com/users/{user_name}/contributions')
    data = BeautifulSoup(res.content, 'html.parser')

    str_today = base_date.strftime('%Y-%m-%d')
    yesterday = base_date - timedelta(days=1)
    str_yesterday = yesterday.strftime('%Y-%m-%d')
    today_rect = data.select_one(f'rect[data-date="{str_today}"]')
    yesterday_rect = data.select_one(f'rect[data-date="{str_yesterday}"]')

    return Contribution(
        user_name=user_name,
        base_date=base_date,
        yesterday=yesterday_rect.attrs['data-count'],
        today=today_rect.attrs['data-count'],
    )


@app.schedule(Cron(0, 15, '*', '*', '?', '*'))  # UTC 15:00 -> JST 0:00
def lambda_handler(event, context={}):
    user_name = os.getenv('USER_NAME')
    today = datetime.now()
    contributions = get_contributions(user_name, today)

    message = get_message(contributions)

    if app.debug:
        app.log.info(message)
    else:
        client = TwitterClient()
        res = client.send(message)

        return {
            'code': res.status_code,
            'body': res.json(),
        }
