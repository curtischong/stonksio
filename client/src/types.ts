export interface Tweet {
  name: string;
  msg: string;
  timestamp: Date;
}

export interface ServerTweet {
  Id?: string;
  Username: string;
  UsrPicUrl?: string;
  Body: string;
  Timestamp: string;
}

export interface Price {
  y: number;
  x: string;
}

export interface ServerPrice {
  Asset: string;
  TradePrice: string;
  Timestamp: string;
}
