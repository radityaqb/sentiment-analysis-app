from typing import Text
import tweepy
import json

# secrets
api_key = ""
api_key_secret = ""
access_token = ""
access_token_secret = ""

auth = tweepy.OAuthHandler(api_key, api_key_secret)

auth.set_access_token(access_token, access_token_secret)

api = tweepy.API(auth)

# result_search = api.search_30_day("production", "goto ipo", maxResults=10)
result_search = api.search_tweets("goto ipo", lang="en", count=100)
for r in result_search:
    print(r.user.screen_name, " : ", r.text)
# print(result_search)
