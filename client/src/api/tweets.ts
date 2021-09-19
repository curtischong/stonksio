import { Tweet, TweetReponse } from '../types';
import apiClient from './client';

const mapTweetResponse = (resp: TweetReponse): Tweet => {
  const ts = new Date(resp.Timestamp)
  return {
    name: resp.Username,
    msg: resp.Body,
    timestamp: `${ts.getHours()}:${ts.getMinutes()}`,
  }
}

export const getTweets = async (count: number) => {
  const response = await apiClient().get('/posts', {
    params: {
      count
    }
  });
  return response.data.map(mapTweetResponse);
};


export const postTweet = async (tweet: Tweet) => {
  const response = await apiClient().post(
    `/tweets`,
    tweet
  );
  const { data } = response;
  return data;
};
