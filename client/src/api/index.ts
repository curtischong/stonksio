import { getTweets, postTweet } from "./tweets";

const stonksAPI = {
  tweets: {
    get: getTweets,
    post: postTweet,
  },
};

export default stonksAPI;
